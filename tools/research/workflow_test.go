package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWorkflowChild(t *testing.T) {
	mode := os.Getenv("PANDA_WORKFLOW_TEST_CHILD")
	if mode == "" {
		return
	}
	if mode == "timeout" {
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}
	if err := os.WriteFile(os.Getenv("PANDA_WORKFLOW_TEST_OUTPUT"), []byte(`{"passed":true,"mean_loss_cp":{"baseline":10,"candidate":9}}`), 0644); err != nil {
		os.Exit(3)
	}
	if mode == "failure" {
		os.Exit(2)
	}
	os.Exit(0)
}
func workflowFixture(t *testing.T, mode string, expected bool) (string, string, string) {
	t.Helper()
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	if err := os.WriteFile(input, []byte("frozen"), 0644); err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(dir, "report.json")
	step := workflowStep{Name: "first", Command: []string{exe, "-test.run=^TestWorkflowChild$"}, Env: map[string]string{"PANDA_WORKFLOW_TEST_CHILD": mode, "PANDA_WORKFLOW_TEST_OUTPUT": report}, Outputs: []string{report}, Gates: []workflowGate{{File: report, Pointer: "/passed", Equals: expected, Brief: []string{"/mean_loss_cp"}}}}
	if mode == "timeout" {
		step.TimeoutSeconds = 1
	}
	plan := workflowPlan{Version: 1, Inputs: []string{input, exe}, Steps: []workflowStep{step, {Name: "second", Command: []string{exe, "-test.run=^TestWorkflowChild$"}}}}
	planPath := filepath.Join(dir, "plan.json")
	if err = workflowJSON(planPath, plan); err != nil {
		t.Fatal(err)
	}
	return planPath, filepath.Join(dir, "run"), input
}
func TestWorkflowGatesAndResume(t *testing.T) {
	for _, c := range []struct {
		mode   string
		gate   bool
		status string
	}{{"success", true, "passed"}, {"success", false, "gate-failed"}, {"failure", true, "execution-failed"}, {"timeout", true, "execution-failed"}} {
		t.Run(c.mode+c.status, func(t *testing.T) {
			plan, out, _ := workflowFixture(t, c.mode, c.gate)
			err := runWorkflow(plan, out, false)
			var state workflowState
			if e := load(filepath.Join(out, "state.json"), &state); e != nil {
				t.Fatal(e)
			}
			if state.Status != c.status || (err == nil) != (c.status == "passed") {
				t.Fatal(state.Status, err)
			}
			if c.status != "passed" && len(state.Results) != 1 {
				t.Fatal("advanced after failure")
			}
			if _, e := os.Stat(filepath.Join(out, "first.txt")); e != nil {
				t.Fatal("failure log missing", e)
			}
			resumeErr := runWorkflow(plan, out, true)
			if (resumeErr == nil) != (c.status == "passed") {
				t.Fatal("unsafe resume", resumeErr)
			}
		})
	}
}
func TestWorkflowRejectsIdentityChanges(t *testing.T) {
	for _, change := range []string{"input", "output", "plan", "environment"} {
		t.Run(change, func(t *testing.T) {
			plan, out, input := workflowFixture(t, "success", true)
			if err := runWorkflow(plan, out, false); err != nil {
				t.Fatal(err)
			}
			path := input
			if change == "output" {
				path = filepath.Join(out, "first.txt")
			}
			if change == "plan" {
				path = plan
			}
			if change == "environment" {
				t.Setenv("GOGC", "123")
			} else {
				f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					t.Fatal(err)
				}
				_, err = f.WriteString(" ")
				f.Close()
				if err != nil {
					t.Fatal(err)
				}
			}
			if err := runWorkflow(plan, out, true); err == nil {
				t.Fatal("accepted changed", change)
			}
		})
	}
}

func TestWorkflowMeasurementLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "measurement.lock")
	release, err := lockWorkflowFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if second, err := lockWorkflowFile(path); err == nil {
		second()
		release()
		t.Fatal("allowed concurrent workflow")
	}
	release()
	third, err := lockWorkflowFile(path)
	if err != nil {
		t.Fatal(err)
	}
	third()
}

func TestWorkflowPinsSymlinkTargetAndContent(t *testing.T) {
	dir := t.TempDir()
	a, b, link := filepath.Join(dir, "a"), filepath.Join(dir, "b"), filepath.Join(dir, "tool")
	for _, path := range []string{a, b} {
		if err := os.WriteFile(path, []byte("same binary bytes"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(a, link); err != nil {
		t.Skip(err)
	}
	if _, err := workflowIdentities([]string{link}); err == nil {
		t.Fatal("source snapshot followed a symlink")
	}
	first, err := workflowIdentities([]string{link}, true)
	if err != nil {
		t.Fatal(err)
	}
	if first[link] != hash([]byte("same binary bytes")) {
		t.Fatal("target content not hashed")
	}
	if err = os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(b, link); err != nil {
		t.Fatal(err)
	}
	second, err := workflowIdentities([]string{link}, true)
	if err != nil {
		t.Fatal(err)
	}
	if first["symlink:"+link] == second["symlink:"+link] {
		t.Fatal("retargeting not detected")
	}
}
func TestWorkflowPointer(t *testing.T) {
	got, err := workflowPointer(map[string]any{"a/b": map[string]any{"~": true}}, "/a~1b/~0")
	if err != nil || got != true {
		t.Fatal(got, err)
	}
	if _, err = workflowPointer(map[string]any{}, "/missing"); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatal(err)
	}
}
