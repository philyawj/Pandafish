package main

import (
	"path/filepath"
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
	if len(plan.Steps) != 5 || plan.Steps[0].Name != "normal-tests" || plan.Steps[1].Name != "race-tests" || plan.Steps[4].Name != "build-engine" {
		t.Fatal(plan.Steps)
	}
	if err := createValidationPlan("pandafish-v0.22", "candidate.test", filepath.Join(t.TempDir(), "bad.json")); err == nil {
		t.Fatal("tagged baseline destination accepted")
	}
	err := run([]string{"warm-plan", "--baseline-test", "base", "--candidate-test", "candidate", "--input", "input", "--evidence", "evidence", "--out", filepath.Join(t.TempDir(), "warm.json")})
	if err == nil || !strings.Contains(err.Error(), "preflight") {
		t.Fatal("CLI warm builder bypassed correctness prerequisite", err)
	}
}
