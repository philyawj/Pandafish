package main

import (
	"path/filepath"
	"testing"
)

func TestDepthReviewSelectionUsesPositionsAndAllObservedDepths(t *testing.T) {
	dir := t.TempDir()
	rows := []probeNative{}
	for _, position := range []string{"position startpos moves e2e4", "position startpos moves d2d4"} {
		row := probeNative{Case: probeCase{ID: "reused-id", Position: position}}
		for i, actor := range []string{"baseline", "candidate"} {
			row.Trials = append(row.Trials, probeTrial{Actor: actor, Query: probeQuery{Move: []string{"e7e5", "d7d5"}[i], Score: probeScore{Depth: 6 + i}}})
		}
		rows = append(rows, row)
	}
	path := filepath.Join(dir, "native.json")
	if err := writeNewJSON(path, rows); err != nil {
		t.Fatal(err)
	}
	selected, err := depthReviewSelection(path)
	if err != nil || len(selected) != 2 || len(selected[0].Depths) != 2 || selected[0].Depths[0] != 6 || selected[0].Depths[1] != 7 {
		t.Fatal(selected, err)
	}
	rows[0].Trials[0].Query.Error = "timeout"
	bad := filepath.Join(dir, "bad.json")
	if err := writeNewJSON(bad, rows); err != nil {
		t.Fatal(err)
	}
	if _, err := depthReviewSelection(bad); err == nil {
		t.Fatal("incomplete input accepted")
	}
}
