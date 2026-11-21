package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Plans contain argv, never a shell program. Scientific gates are explicit,
// frozen inputs; the runner neither selects hypotheses nor authorizes publication.
type workflowPlan struct {
	Version int            `json:"version"`
	Inputs  []string       `json:"inputs"`
	Steps   []workflowStep `json:"steps"`
}
type workflowStep struct {
	Name           string            `json:"name"`
	Command        []string          `json:"command"`
	Env            map[string]string `json:"env,omitempty"`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty"`
	Outputs        []string          `json:"outputs,omitempty"`
	Gates          []workflowGate    `json:"gates,omitempty"`
}
type workflowGate struct {
	File    string   `json:"file"`
	Pointer string   `json:"pointer"`
	Equals  any      `json:"equals"`
	Brief   []string `json:"brief,omitempty"`
}
type workflowResult struct {
	Name    string               `json:"name"`
	Status  string               `json:"status"`
	Started string               `json:"started_utc"`
	Seconds float64              `json:"seconds"`
	Exit    int                  `json:"exit"`
	Error   string               `json:"error,omitempty"`
	Files   map[string]fileEntry `json:"files,omitempty"`
	Brief   map[string]any       `json:"brief,omitempty"`
}
type workflowState struct {
	Version     int                  `json:"version"`
	PlanSHA     string               `json:"plan_sha256"`
	Inputs      map[string]fileEntry `json:"inputs"`
	Environment map[string]string    `json:"environment"`
	Status      string               `json:"status"`
	Results     []workflowResult     `json:"results"`
}

func writeNewFile(path string, b []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
func writeNewJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return writeNewFile(path, append(b, '\n'))
}

func workflowIdentities(patterns []string, followLinks ...bool) (map[string]fileEntry, error) {
	identities := map[string]fileEntry{}
	for _, pattern := range patterns {
		paths, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(paths) == 0 {
			return nil, fmt.Errorf("input has no matches: %s", pattern)
		}
		for _, path := range paths {
			if strings.HasPrefix(path, "symlink:") {
				return nil, fmt.Errorf("reserved identity path: %s", path)
			}
			info, err := os.Lstat(path)
			if err != nil {
				return nil, err
			}
			if info.Mode()&os.ModeSymlink != 0 && len(followLinks) > 0 && followLinks[0] {
				link, err := os.Readlink(path)
				if err != nil {
					return nil, err
				}
				resolved, err := filepath.EvalSymlinks(path)
				if err != nil {
					return nil, err
				}
				resolved, err = filepath.Abs(resolved)
				if err != nil {
					return nil, err
				}
				identities["symlink:"+path] = hash([]byte(link + "\n" + resolved))
				info, err = os.Stat(path)
				if err != nil {
					return nil, err
				}
			}
			if !info.Mode().IsRegular() {
				return nil, fmt.Errorf("input is not a regular file: %s", path)
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			identities[path] = hash(b)
		}
	}
	return identities, nil
}
func workflowJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err = os.WriteFile(tmp, append(b, '\n'), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
func workflowPointer(v any, pointer string) (any, error) {
	if pointer == "" {
		return v, nil
	}
	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("invalid JSON pointer %q", pointer)
	}
	for _, part := range strings.Split(pointer[1:], "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		object, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("pointer does not name an object: %s", pointer)
		}
		v, ok = object[part]
		if !ok {
			return nil, fmt.Errorf("missing JSON pointer %s", pointer)
		}
	}
	return v, nil
}
func workflowStatus(out string) error {
	var state workflowState
	if err := load(filepath.Join(out, "state.json"), &state); err != nil {
		return err
	}
	steps := []map[string]any{}
	for _, r := range state.Results {
		steps = append(steps, map[string]any{"name": r.Name, "status": r.Status, "seconds": r.Seconds, "error": r.Error, "brief": r.Brief})
	}
	return printJSON(map[string]any{"status": state.Status, "steps": steps, "plan_sha256": state.PlanSHA})
}
func runWorkflow(planPath, out string, resume bool) error {
	if out == "" {
		return fmt.Errorf("new output directory required")
	}
	b, err := os.ReadFile(planPath)
	if err != nil {
		return err
	}
	var plan workflowPlan
	if err = json.Unmarshal(b, &plan); err != nil {
		return err
	}
	if plan.Version != 1 || len(plan.Steps) == 0 || len(plan.Inputs) == 0 {
		return fmt.Errorf("version 1, inputs and steps required")
	}
	names := map[string]bool{}
	for _, step := range plan.Steps {
		if !safePath(step.Name) || strings.Contains(step.Name, "/") || names[step.Name] || len(step.Command) == 0 || step.TimeoutSeconds < 0 {
			return fmt.Errorf("invalid or duplicate step %q", step.Name)
		}
		names[step.Name] = true
	}
	inputs, err := workflowIdentities(plan.Inputs, true)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	release, err := lockWorkflowFile(filepath.Join(os.TempDir(), "pandafish-research-"+hash([]byte(cwd)).SHA[:16]+".lock"))
	if err != nil {
		return err
	}
	defer release()
	environment := map[string]string{}
	for _, key := range []string{"PATH", "GOMAXPROCS", "GOGC", "GOMEMLIMIT", "GODEBUG", "GOFLAGS", "GOTOOLCHAIN", "GOOS", "GOARCH"} {
		environment[key] = os.Getenv(key)
	}
	state := workflowState{Version: 1, PlanSHA: hash(b).SHA, Inputs: inputs, Environment: environment, Status: "pending", Results: []workflowResult{}}
	statePath := filepath.Join(out, "state.json")
	if resume {
		var saved workflowState
		if err = load(statePath, &saved); err != nil {
			return err
		}
		if saved.PlanSHA != hash(b).SHA || !reflect.DeepEqual(saved.Inputs, inputs) || !reflect.DeepEqual(saved.Environment, environment) {
			return fmt.Errorf("plan or input identity changed; preserved run cannot resume")
		}
		state = saved
		if state.Version != 1 || len(state.Results) > len(plan.Steps) {
			return fmt.Errorf("invalid saved workflow state")
		}
		for i, result := range state.Results {
			if result.Name != plan.Steps[i].Name || len(result.Files) == 0 {
				return fmt.Errorf("incomplete saved step %d", i+1)
			}
			for path, want := range result.Files {
				data, e := os.ReadFile(path)
				if e != nil || hash(data) != want {
					return fmt.Errorf("saved output changed: %s", path)
				}
			}
			if result.Status != "passed" {
				return fmt.Errorf("step %s is %s; no implicit rerun or advance", result.Name, result.Status)
			}
		}
	} else {
		if err = os.Mkdir(out, 0755); err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(out, "plan.json"), b, 0644); err != nil {
			return err
		}
	}
	lock, err := os.OpenFile(filepath.Join(out, "running.lock"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("workflow already active or interrupted: %w", err)
	}
	lock.Close()
	defer os.Remove(filepath.Join(out, "running.lock"))
	if err = workflowJSON(statePath, state); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	for i := len(state.Results); i < len(plan.Steps); i++ {
		step := plan.Steps[i]
		fresh, e := workflowIdentities(plan.Inputs, true)
		if e != nil || !reflect.DeepEqual(fresh, inputs) {
			return fmt.Errorf("input changed before %s", step.Name)
		}
		start := time.Now()
		result := workflowResult{Name: step.Name, Status: "running", Started: start.UTC().Format(time.RFC3339Nano), Exit: -1}
		state.Status = "running"
		state.Results = append(state.Results, result)
		if err = workflowJSON(statePath, state); err != nil {
			return err
		}
		fmt.Printf("start %s (%d/%d)\n", step.Name, i+1, len(plan.Steps))
		logPath := filepath.Join(out, step.Name+".txt")
		log, e := os.OpenFile(logPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if e != nil {
			return e
		}
		stepCtx, cancel := context.WithCancel(ctx)
		if step.TimeoutSeconds > 0 {
			cancel()
			stepCtx, cancel = context.WithTimeout(ctx, time.Duration(step.TimeoutSeconds)*time.Second)
		}
		cmd := exec.CommandContext(stepCtx, step.Command[0], step.Command[1:]...)
		cmd.Env = os.Environ()
		for k, v := range step.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
		cmd.Stdout = log
		cmd.Stderr = log
		configureWorkflowProcess(cmd)
		runErr := cmd.Run()
		cancel()
		log.Close()
		result.Seconds = time.Since(start).Seconds()
		result.Status = "passed"
		result.Exit = 0
		if runErr != nil {
			result.Status = "execution-failed"
			result.Exit = -1
			if cmd.ProcessState != nil {
				result.Exit = cmd.ProcessState.ExitCode()
			}
			result.Error = runErr.Error()
		}
		result.Brief = map[string]any{}
		for _, gate := range step.Gates {
			var report any
			e := load(gate.File, &report)
			if e == nil {
				var got any
				got, e = workflowPointer(report, gate.Pointer)
				if e == nil && !reflect.DeepEqual(got, gate.Equals) {
					e = fmt.Errorf("gate %s%s = %v, expected %v", gate.File, gate.Pointer, got, gate.Equals)
				}
				for _, key := range gate.Brief {
					v, briefErr := workflowPointer(report, key)
					if briefErr == nil {
						result.Brief[gate.File+key] = v
					}
				}
			}
			if e != nil {
				if result.Status == "passed" {
					result.Status = "gate-failed"
					result.Error = e.Error()
				}
			}
		}
		patterns := append([]string{logPath}, step.Outputs...)
		for _, gate := range step.Gates {
			patterns = append(patterns, gate.File)
		}
		result.Files = map[string]fileEntry{}
		for _, pattern := range patterns {
			paths, e := filepath.Glob(pattern)
			if e != nil || len(paths) == 0 {
				if result.Status == "passed" {
					result.Status = "execution-failed"
					result.Error = "missing output: " + pattern
				}
				continue
			}
			for _, path := range paths {
				data, e := os.ReadFile(path)
				if e != nil {
					if result.Status == "passed" {
						result.Status = "execution-failed"
						result.Error = e.Error()
					}
					continue
				}
				result.Files[path] = hash(data)
			}
		}
		if fresh, e := workflowIdentities(plan.Inputs, true); e != nil || !reflect.DeepEqual(fresh, inputs) {
			result.Status = "identity-failed"
			result.Error = "input changed during step"
		}
		state.Results[i] = result
		state.Status = result.Status
		if err = workflowJSON(statePath, state); err != nil {
			return err
		}
		fmt.Printf("finish %s %s %.1fs\n", step.Name, result.Status, result.Seconds)
		if result.Status != "passed" {
			return fmt.Errorf("stopped at %s: %s", step.Name, result.Error)
		}
	}
	state.Status = "passed"
	return workflowJSON(statePath, state)
}
