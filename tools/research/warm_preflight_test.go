package main

import (
	"path/filepath"
	"testing"
)

func TestWarmPreflightRejectsFailedOrChangedSources(t *testing.T) {
	current := map[string]fileEntry{"engine.go": hash([]byte("current"))}
	for _, name := range []string{"passed", "failed", "stale"} {
		state := workflowState{Status: "passed", Inputs: map[string]fileEntry{"engine.go": current["engine.go"]}}
		if name == "failed" {
			state.Status = "execution-failed"
		}
		if name == "stale" {
			state.Inputs["engine.go"] = hash([]byte("old"))
		}
		path := filepath.Join(t.TempDir(), "state.json")
		if err := writeNewJSON(path, state); err != nil {
			t.Fatal(err)
		}
		err := validateWarmPreflight(path, current)
		if (err == nil) != (name == "passed") {
			t.Fatal(name, err)
		}
	}
}
