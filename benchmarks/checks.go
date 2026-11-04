package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type checkRecord struct {
	Name      string   `json:"name"`
	Command   []string `json:"command"`
	Passed    bool     `json:"passed"`
	ElapsedMS int64    `json:"elapsed_ms"`
}

type commandRunner func(string, ...string) (string, error)

func runCommand(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	return string(output), err
}

func checkError(name, output string, err error) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 40 {
		lines = append([]string{"(last 40 output lines)"}, lines[len(lines)-40:]...)
	}
	return fmt.Errorf("%s failed: %w\n%s", name, err, strings.Join(lines, "\n"))
}

// validateCandidate delegates checks to Go tools and keeps successful output
// compact. Full game records remain available without needing to inspect them.
func validateCandidate(candidate string, race bool, run commandRunner, out io.Writer) ([]checkRecord, error) {
	dirs, err := run("go", "list", "-f", "{{.Dir}}", "./...")
	if err != nil {
		return nil, checkError("list packages", dirs, err)
	}
	var files []string
	for _, dir := range strings.Split(strings.TrimSpace(dirs), "\n") {
		matches, err := filepath.Glob(filepath.Join(dir, "*.go"))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no Go files found; run from the repository root")
	}
	if err := os.MkdirAll(filepath.Dir(candidate), 0755); err != nil {
		return nil, err
	}
	testArgs := []string{"test", "./..."}
	if race {
		testArgs = []string{"test", "-race", "./..."}
	}
	steps := []struct {
		name, executable string
		args             []string
	}{
		{"format", "gofmt", append([]string{"-l"}, files...)},
		{"tests", "go", testArgs},
		{"vet", "go", []string{"vet", "./..."}},
		{"build", "go", []string{"build", "-o", candidate, "."}},
	}
	var checks []checkRecord
	for _, step := range steps {
		started := time.Now()
		output, err := run(step.executable, step.args...)
		if err == nil && step.name == "format" && strings.TrimSpace(output) != "" {
			err = fmt.Errorf("files need gofmt")
		}
		checks = append(checks, checkRecord{Name: step.name, Command: append([]string{step.executable}, step.args...), Passed: err == nil, ElapsedMS: time.Since(started).Milliseconds()})
		if err != nil {
			return checks, checkError(step.name, output, err)
		}
		fmt.Fprintln(out, "PASS", step.name)
	}
	return checks, nil
}
