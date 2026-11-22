package main

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestValidationPlanOrdersChecksBeforeBuildsAndProtectsBaseline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := createValidationPlan("candidate", "candidate.test", path); err != nil {
		t.Fatal(err)
	}
	var plan workflowPlan
	if err := load(path, &plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps) != 8 || plan.Steps[0].Name != "normal-tests" || plan.Steps[1].Name != "race-tests" || plan.Steps[7].Name != "build-engine" {
		t.Fatal(plan.Steps)
	}
	for i, command := range [][]string{
		{"go", "-C", "third_party/chess", "test", "./..."},
		{"go", "-C", "third_party/chess", "test", "-race", "./..."},
		{"go", "-C", "third_party/chess", "vet", "./..."},
	} {
		if !slices.Equal(plan.Steps[i+3].Command, command) {
			t.Fatal("nested checks must precede builds", plan.Steps)
		}
	}
	if err := createValidationPlan("pandafish-v0.22", "candidate.test", filepath.Join(t.TempDir(), "bad.json")); err == nil {
		t.Fatal("tagged baseline destination accepted")
	}
	err := run([]string{"warm-plan", "--baseline-test", "base", "--candidate-test", "candidate", "--input", "input", "--evidence", "evidence", "--out", filepath.Join(t.TempDir(), "warm.json")})
	if err == nil || !strings.Contains(err.Error(), "preflight") {
		t.Fatal("CLI warm builder bypassed correctness prerequisite", err)
	}
}
