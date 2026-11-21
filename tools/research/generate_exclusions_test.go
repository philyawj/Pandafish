package main

import (
	"path/filepath"
	"testing"
)

func TestAdditionalExclusionsPreserveCanonicalBoundary(t *testing.T) {
	dir := t.TempDir()
	fen := "4k3/8/8/8/3P4/8/8/4K3 w - - 0 1"
	cases := filepath.Join(dir, "cases.json")
	analysis := filepath.Join(dir, "analysis.json")
	metadata := filepath.Join(dir, "metadata.json")
	if err := workflowJSON(cases, map[string]any{"Cases": []probeCase{{FEN: fen}}}); err != nil {
		t.Fatal(err)
	}
	if err := workflowJSON(analysis, map[string]any{"game": "test", "moves": []map[string]string{{"fen": fen}}}); err != nil {
		t.Fatal(err)
	}
	if err := workflowJSON(metadata, map[string]string{"status": "complete"}); err != nil {
		t.Fatal(err)
	}
	blocked := map[string]bool{}
	sources, err := blockAdditionalProbeSources(blocked, []string{cases}, []string{analysis, metadata})
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) != 3 || len(blocked) != 1 || sources[2]["ignored_non_game_metadata"] != true {
		t.Fatal(sources, blocked)
	}
	g, err := probeGame(probeCase{FEN: "4k3/8/8/3p4/8/8/8/4K3 b - - 0 1"})
	if err != nil {
		t.Fatal(err)
	}
	if !blocked[canonicalProbe(g.Position())] {
		t.Fatal("color mirror escaped exclusion")
	}
	if _, err = blockAdditionalProbeSources(blocked, []string{filepath.Join(dir, "missing*")}, nil); err == nil {
		t.Fatal("empty exclusion accepted")
	}
}

func TestAdditionalExclusionsEmptyCompletedGame(t *testing.T) {
	for _, complete := range []bool{false, true} {
		path := filepath.Join(t.TempDir(), "abandoned.json")
		if err := workflowJSON(path, map[string]any{"game": "zero-ply", "moves": nil, "complete": complete}); err != nil {
			t.Fatal(err)
		}
		blocked := map[string]bool{}
		sources, err := blockAdditionalProbeSources(blocked, nil, []string{path})
		if !complete {
			if err == nil {
				t.Fatal("incomplete empty analysis accepted")
			}
			continue
		}
		if err != nil || len(sources) != 1 || len(blocked) != 0 || sources[0]["empty_completed_game"] != true {
			t.Fatal(sources, blocked, err)
		}
	}
}
