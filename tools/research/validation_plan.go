package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func createValidationPlan(candidate, testBinary, out string) error {
	if candidate == "" || testBinary == "" || out == "" {
		return fmt.Errorf("candidate/test destinations and new plan required")
	}
	if regexp.MustCompile(`^pandafish-v[0-9]+\.[0-9]+(\.[0-9]+)?(\.exe)?$`).MatchString(filepath.Base(candidate)) {
		return fmt.Errorf("research validation cannot overwrite a tagged-baseline-shaped executable")
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	plan := workflowPlan{Version: 1, Inputs: append([]string{self, "tools/research/*.go", "testdata/*.json"}, coreSourcePatterns()...), Steps: []workflowStep{
		{Name: "normal-tests", Command: []string{"go", "test", "./..."}, TimeoutSeconds: 180},
		{Name: "race-tests", Command: []string{"go", "test", "-race", "./..."}, TimeoutSeconds: 300},
		{Name: "vet", Command: []string{"go", "vet", "./..."}, TimeoutSeconds: 180},
		{Name: "build-warm-test", Command: []string{"go", "test", "-c", "-o", testBinary, "."}, Outputs: []string{testBinary}, TimeoutSeconds: 180},
		{Name: "build-engine", Command: []string{"go", "build", "-buildvcs=false", "-o", candidate, "."}, Outputs: []string{candidate}, TimeoutSeconds: 180},
	}}
	if err = writeNewJSON(out, plan); err != nil {
		return err
	}
	return printJSON(map[string]any{"plan": out, "steps": len(plan.Steps), "preflight_state": "Pass the completed workflow state.json to warm-plan --preflight-state."})
}
