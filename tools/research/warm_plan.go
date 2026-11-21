package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type warmPlanOptions struct {
	Parallel                                                     bool
	Values                                                       bool
	Preflight                                                    string
	BaselineTest, CandidateTest, Input, Evidence, Out, Reference string
	Depth, Roots                                                 int
	MaxTimeRatio, MaxNodeRatio                                   float64
	Exact                                                        bool
}

func createWarmPlan(o warmPlanOptions) error {
	if o.Parallel && (o.Exact || o.Reference != "" || o.MaxNodeRatio > 0) {
		return fmt.Errorf("parallel profile has separate value/co-best controls and reports node variation")
	}
	if o.BaselineTest == "" || o.CandidateTest == "" || o.Input == "" || o.Evidence == "" || o.Out == "" || o.Depth < 1 || o.Roots < 1 {
		return fmt.Errorf("test executables, input, evidence, new plan and positive depth/roots required")
	}
	if _, err := os.Stat(o.Out); !os.IsNotExist(err) {
		return fmt.Errorf("refuse existing plan %s", o.Out)
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	plan := workflowPlan{Version: 1, Inputs: append([]string{self, o.BaselineTest, o.CandidateTest, filepath.Join(o.Input, "summary.json"), filepath.Join(o.Input, "previous.pgn"), "tools/research/*.go"}, coreSourcePatterns()...)}
	if o.Preflight != "" {
		identities, err := workflowIdentities(coreSourcePatterns())
		if err != nil {
			return err
		}
		if err = validateWarmPreflight(o.Preflight, identities); err != nil {
			return err
		}
		plan.Inputs = append(plan.Inputs, o.Preflight)
		plan.Steps = append(plan.Steps, workflowStep{Name: "validated-source", Command: []string{self, "workflow", "--out", filepath.Dir(o.Preflight), "--status"}, Gates: []workflowGate{{File: o.Preflight, Pointer: "/status", Equals: "passed"}}})
	}
	if o.Reference != "" {
		plan.Inputs = append(plan.Inputs, o.Reference)
	}
	for pair := 1; pair <= 3; pair++ {
		order := []string{"baseline", "candidate"}
		if pair == 2 {
			order[0], order[1] = order[1], order[0]
		}
		for _, actor := range order {
			binary := o.BaselineTest
			if actor == "candidate" {
				binary = o.CandidateTest
			}
			name := fmt.Sprintf("warm-%s-pair-%d", actor, pair)
			path, err := filepath.Abs(filepath.Join(o.Evidence, name+".json"))
			if err != nil {
				return err
			}
			input, err := filepath.Abs(o.Input)
			if err != nil {
				return err
			}
			plan.Steps = append(plan.Steps, workflowStep{Name: name, Command: []string{binary, "-test.run=^TestResearchWarmWorkload$", "-test.count=1", "-test.v"}, Env: map[string]string{"WARM_INPUT": input, "WARM_OUTPUT": path, "WARM_MODE": "search", "WARM_DEPTH": fmt.Sprint(o.Depth)}, Outputs: []string{path}, TimeoutSeconds: 600})
		}
		// A deterministic node failure cannot be rescued by another timing pair.
		if o.MaxNodeRatio > 0 && pair < 3 {
			plan.Steps = append(plan.Steps, warmReportStep(self, o, pair, 0, fmt.Sprintf("node-gate-%d", pair)))
		}
	}
	plan.Steps = append(plan.Steps, warmReportStep(self, o, 3, o.MaxTimeRatio, "warm-gate"))
	if o.Parallel {
		input, err := filepath.Abs(o.Input)
		if err != nil {
			return err
		}
		output, err := filepath.Abs(filepath.Join(o.Evidence, "verified-choices.json"))
		if err != nil {
			return err
		}
		summary, err := filepath.Abs(filepath.Join(o.Evidence, "choice-validation.json"))
		if err != nil {
			return err
		}
		pattern, err := filepath.Abs(filepath.Join(o.Evidence, "warm-candidate-pair-*.json"))
		if err != nil {
			return err
		}
		plan.Steps = append(plan.Steps, workflowStep{Name: "choice-validation", Command: []string{o.CandidateTest, "-test.run=^TestResearchParallelChoices$", "-test.count=1", "-test.v"}, Env: map[string]string{"WARM_INPUT": input, "WARM_OUTPUT": output, "WARM_MODE": "search", "WARM_DEPTH": fmt.Sprint(o.Depth), "PARALLEL_VALIDATION_ROWS": pattern, "PARALLEL_VALIDATION_OUTPUT": summary}, Outputs: []string{output, summary}, Gates: []workflowGate{{File: summary, Pointer: "/passed", Equals: true, Brief: []string{"/roots", "/distinct_choices", "/full_window_verified"}}}, TimeoutSeconds: 600})
	}
	if err = writeNewJSON(o.Out, plan); err != nil {
		return err
	}
	return printJSON(map[string]any{"plan": o.Out, "steps": len(plan.Steps), "alternating_pairs": 3})
}

func validateWarmPreflight(path string, current map[string]fileEntry) error {
	var state workflowState
	if err := load(path, &state); err != nil {
		return err
	}
	if state.Status != "passed" {
		return fmt.Errorf("correctness preflight is not passed: %s", state.Status)
	}
	for path, want := range current {
		if state.Inputs[path] != want {
			return fmt.Errorf("correctness preflight source differs: %s", path)
		}
	}
	return nil
}
func warmReportStep(self string, o warmPlanOptions, pairs int, maxTime float64, name string) workflowStep {
	// Bracket expressions name only completed pairs in an early report.
	pattern := fmt.Sprintf("[1-%d]", pairs)
	report := filepath.Join(o.Evidence, name+".json")
	args := []string{self, "warm-report", "--baseline", filepath.Join(o.Evidence, "warm-baseline-pair-"+pattern+".json"), "--candidate", filepath.Join(o.Evidence, "warm-candidate-pair-"+pattern+".json"), "--pairs", fmt.Sprint(pairs), "--roots", fmt.Sprint(o.Roots), "--depth", fmt.Sprint(o.Depth), "--max-time-ratio", fmt.Sprint(maxTime), "--max-node-ratio", fmt.Sprint(o.MaxNodeRatio), "--out", report}
	if o.Exact {
		args = append(args, "--exact")
	}
	if o.Values {
		args = append(args, "--values")
	}
	if o.Parallel {
		args = append(args, "--parallel")
	}
	if o.Reference != "" {
		args = append(args, "--reference", o.Reference)
	}
	return workflowStep{Name: name, Command: args, Outputs: []string{report}, Gates: []workflowGate{{File: report, Pointer: "/passed", Equals: true, Brief: []string{"/median_time_ratio", "/node_ratios", "/issues"}}}}
}
