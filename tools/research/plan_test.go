package main

import (
	"path/filepath"
	"slices"
	"testing"
)

func commandFlag(args []string, name string) string {
	i := slices.Index(args, name)
	if i < 0 || i+1 == len(args) {
		return ""
	}
	return args[i+1]
}
func TestNativePlanPreservesProtocols(t *testing.T) {
	for _, profile := range []string{"existing", "fresh"} {
		t.Run(profile, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "plan.json")
			if err := nativePlan("baseline", "candidate", "teacher", "cache", "evidence", profile, "fixed-input", path); err != nil {
				t.Fatal(err)
			}
			var plan workflowPlan
			if err := load(path, &plan); err != nil {
				t.Fatal(err)
			}
			count := 4
			if profile == "fresh" {
				count = 2
			}
			if len(plan.Steps) != count {
				t.Fatal("schedule")
			}
			for i, step := range plan.Steps {
				ms := "100"
				if i%2 == 1 {
					ms = "1000"
				}
				trials := "1"
				if i == 0 {
					trials = "3"
				}
				if commandFlag(step.Command, "--milliseconds") != ms || commandFlag(step.Command, "--trials") != trials || commandFlag(step.Command, "--workers") != "3" {
					t.Fatal("protocol drift", step)
				}
				if (profile == "fresh" || i >= 2) != slices.Contains(step.Command, "--by-turn") {
					t.Fatal("turn protection")
				}
				if len(step.Gates) != 2 {
					t.Fatal("missing gates")
				}
			}
			if profile == "existing" && !slices.Contains(plan.Inputs, filepath.Join("cache", "public-100", "references.jsonl")) {
				t.Fatal("reference identity missing")
			}
			if profile == "fresh" && commandFlag(plan.Steps[1].Command, "--cache") != filepath.Join("evidence", "fresh-100") {
				t.Fatal("fresh cache source")
			}
			if err := nativePlan("baseline", "candidate", "teacher", "cache", "evidence", profile, "fixed-input", path); err == nil {
				t.Fatal("plan overwritten")
			}
		})
	}
}
func TestWarmPlanAlternationAndEarlyNodeGate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	o := warmPlanOptions{BaselineTest: "base.test", CandidateTest: "candidate.test", Input: "trajectory", Evidence: "evidence", Out: path, Roots: 1145, Depth: 5, MaxTimeRatio: .85, MaxNodeRatio: .8, Exact: true}
	if err := createWarmPlan(o); err != nil {
		t.Fatal(err)
	}
	var plan workflowPlan
	if err := load(path, &plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps) != 9 {
		t.Fatal("missing early gates")
	}
	actors := []string{}
	for _, step := range plan.Steps {
		if len(step.Command) > 1 && step.Command[1] == "warm-report" {
			continue
		}
		actors = append(actors, step.Command[0])
		if step.Env["WARM_DEPTH"] != "5" || step.Env["WARM_MODE"] != "search" {
			t.Fatal("workload changed")
		}
	}
	if !slices.Equal(actors, []string{"base.test", "candidate.test", "candidate.test", "base.test", "base.test", "candidate.test"}) {
		t.Fatal("pair order", actors)
	}
	if commandFlag(plan.Steps[2].Command, "--max-time-ratio") != "0" || commandFlag(plan.Steps[2].Command, "--max-node-ratio") != "0.8" {
		t.Fatal("early gate changed")
	}
	if commandFlag(plan.Steps[8].Command, "--max-time-ratio") != "0.85" || !slices.Contains(plan.Steps[8].Command, "--exact") {
		t.Fatal("final gate changed")
	}
	if err := createWarmPlan(o); err == nil {
		t.Fatal("plan overwritten")
	}
}
