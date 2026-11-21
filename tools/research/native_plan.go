package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Only the already documented suite schedule is generated here. A fresh plan
// is separate: final controls are collected only after the existing gates pass.
func nativePlan(base, candidate, teacher, cache, evidence, profile, input, out string) error {
	if base == "" || candidate == "" || teacher == "" || evidence == "" || out == "" {
		return fmt.Errorf("baseline, candidate, teacher, evidence and new plan output required")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		return fmt.Errorf("refuse existing plan %s", out)
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	plan := workflowPlan{Version: 1, Inputs: append([]string{self, base, candidate, teacher, "tools/research/*.go"}, coreSourcePatterns()...)}
	type suite struct {
		name, input string
		ms, trials  int
		byTurn      bool
	}
	suites := []suite{{"public-100", "tools/benchmark/research/v0.17-search-moves/inputs.json", 100, 3, false}, {"public-1000", "tools/benchmark/research/v0.17-search-moves/inputs.json", 1000, 1, false}, {"newer-100", "tools/benchmark/research/v0.20-next-assessment/fresh-generation/inputs.json", 100, 1, true}, {"newer-1000", "tools/benchmark/research/v0.20-next-assessment/fresh-generation/inputs.json", 1000, 1, true}}
	if profile == "fresh" {
		if input == "" {
			return fmt.Errorf("fresh profile requires fixed input")
		}
		suites = []suite{{"fresh-100", input, 100, 3, true}, {"fresh-1000", input, 1000, 1, true}}
	} else if profile != "existing" {
		return fmt.Errorf("profile must be existing or fresh")
	}
	for _, s := range suites {
		plan.Inputs = append(plan.Inputs, s.input)
		dir := filepath.Join(evidence, s.name)
		args := []string{self, "choices", "--baseline", base, "--candidate", candidate, "--teacher", teacher, "--input", s.input, "--out", dir, "--milliseconds", fmt.Sprint(s.ms), "--trials", fmt.Sprint(s.trials), "--workers", "3"}
		if s.byTurn {
			args = append(args, "--by-turn")
		}
		if profile == "existing" && cache != "" {
			path := filepath.Join(cache, s.name)
			args = append(args, "--cache", path)
			plan.Inputs = append(plan.Inputs, filepath.Join(path, "method.json"), filepath.Join(path, "quality.json"), filepath.Join(path, "references.jsonl"))
		}
		if profile == "fresh" && s.ms == 1000 {
			args = append(args, "--cache", filepath.Join(evidence, "fresh-100"))
		}
		plan.Steps = append(plan.Steps, workflowStep{Name: s.name, Command: args, Outputs: []string{filepath.Join(dir, "*.json"), filepath.Join(dir, "*.jsonl")}, Gates: []workflowGate{
			{File: filepath.Join(dir, "quality.json"), Pointer: "/passed", Equals: true},
			{File: filepath.Join(dir, "brief.json"), Pointer: "/passed", Equals: true, Brief: []string{"/finite_cases", "/mean_loss_cp", "/by_turn_mean_loss_cp", "/severe_over_300", "/new_candidate_losing_mate_cases", "/timing_overruns", "/worst_regressions"}},
		}})
	}
	if err = writeNewJSON(out, plan); err != nil {
		return err
	}
	return printJSON(map[string]any{"plan": out, "steps": len(plan.Steps), "profile": profile})
}
