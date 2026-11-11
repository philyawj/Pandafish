package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidationDelegatesChecksAndSummarizes(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var calls []string
	run := func(name string, args ...string) (string, error) {
		command := strings.Join(append([]string{name}, args...), " ")
		calls = append(calls, command)
		if strings.HasPrefix(command, "go list ") {
			return dir + "\n", nil
		}
		if name == "gofmt" {
			return "", nil
		}
		return "verbose successful tool output\n", nil
	}
	var output bytes.Buffer
	checks, err := validateCandidate(filepath.Join(dir, "bin", "candidate"), true, run, &output)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 4 || len(calls) != 5 {
		t.Fatalf("missing checks: %v", calls)
	}
	if !strings.Contains(calls[2], "go test -race ./...") || !strings.Contains(calls[4], "go build -o ") {
		t.Fatalf("wrong tools: %v", calls)
	}
	if output.String() != "PASS format\nPASS tests\nPASS vet\nPASS build\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
	for _, check := range checks {
		if !check.Passed {
			t.Fatal("successful check marked failed")
		}
	}
}

func TestValidationStopsBeforeBuildOnFailure(t *testing.T) {
	for _, failure := range []string{"format", "tests"} {
		t.Run(failure, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "main.go"), nil, 0644); err != nil {
				t.Fatal(err)
			}
			built := false
			run := func(name string, args ...string) (string, error) {
				if name == "go" && args[0] == "list" {
					return dir, nil
				}
				if name == "gofmt" && failure == "format" {
					return "main.go\n", nil
				}
				if name == "go" && args[0] == "test" {
					return strings.Repeat("failure detail\n", 100), fmt.Errorf("test failure")
				}
				if name == "go" && args[0] == "build" {
					built = true
				}
				return "", nil
			}
			var output bytes.Buffer
			_, err := validateCandidate(filepath.Join(dir, "candidate"), false, run, &output)
			if err == nil || built {
				t.Fatal("validation proceeded after failure")
			}
			if len(strings.Split(err.Error(), "\n")) > 43 {
				t.Fatal("failure output was not bounded")
			}
		})
	}
}
