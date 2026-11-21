package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWarmReportFrozenGates(t *testing.T) {
	for _, c := range []struct {
		name                   string
		ratio                  float64
		exact, change, illegal bool
		opportunity            float64
		want                   bool
	}{
		{"exact", .8, true, false, false, 0, true},
		{"timing", 1.1, false, false, false, 0, false},
		{"selective", .8, false, true, false, 0, true},
		{"exact-change", .8, true, true, false, 0, false},
		{"illegal", .8, false, false, true, 0, false},
		{"opportunity", .8, true, false, false, .2, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			a := warmReportRow{Game: 1, Ply: 8, Depth: 5, Best: "e2e4", FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", Nodes: 100, SearchMS: 10}
			b := a
			b.SearchMS *= c.ratio
			b.DeltaNodes = 10
			if c.change {
				b.Score = 1
			}
			if c.illegal {
				b.Best = "0000"
			}
			for name, row := range map[string]warmReportRow{"base": a, "candidate": b} {
				if err := workflowJSON(filepath.Join(dir, name+".json"), map[string]any{"Rows": []warmReportRow{row}}); err != nil {
					t.Fatal(err)
				}
			}
			o := warmReportOptions{Baseline: filepath.Join(dir, "base.json"), Candidate: filepath.Join(dir, "candidate.json"), Out: filepath.Join(dir, "report.json"), Pairs: 1, Roots: 1, Depth: 5, MaxTimeRatio: .9, Exact: c.exact, MinOpportunity: c.opportunity}
			if err := warmReport(o); err != nil {
				t.Fatal(err)
			}
			var report struct{ Passed bool }
			if err := load(o.Out, &report); err != nil {
				t.Fatal(err)
			}
			if report.Passed != c.want {
				t.Fatal("wrong gate", report.Passed, c.want)
			}
			if err := warmReport(o); err == nil {
				t.Fatal("overwrote report")
			}
		})
	}
}

func TestWarmValueModeAllowsOnlyNodeChanges(t *testing.T) {
	for _, change := range []string{"nodes", "score", "choice"} {
		t.Run(change, func(t *testing.T) {
			dir := t.TempDir()
			base := warmReportRow{Game: 1, Ply: 8, Depth: 5, Best: "e2e4", FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", Nodes: 100, SearchMS: 10}
			candidate := base
			candidate.Nodes, candidate.SearchMS = 80, 8
			if change == "score" {
				candidate.Score = 1
			}
			if change == "choice" {
				candidate.Best = "d2d4"
			}
			for name, row := range map[string]warmReportRow{"base": base, "candidate": candidate} {
				if err := writeNewJSON(filepath.Join(dir, name+".json"), map[string]any{"Rows": []warmReportRow{row}}); err != nil {
					t.Fatal(err)
				}
			}
			o := warmReportOptions{Baseline: filepath.Join(dir, "base.json"), Candidate: filepath.Join(dir, "candidate.json"), Out: filepath.Join(dir, "report.json"), Pairs: 1, Roots: 1, Depth: 5, Values: true}
			if err := warmReport(o); err != nil {
				t.Fatal(err)
			}
			var report struct{ Passed bool }
			if err := load(o.Out, &report); err != nil {
				t.Fatal(err)
			}
			if report.Passed != (change == "nodes") {
				t.Fatal(change, report.Passed)
			}
		})
	}
}
func TestInventoryRefusesRewriteAndSymlink(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "raw.json")
	if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := createInventory(dir); err != nil {
		t.Fatal(err)
	}
	if err := createInventory(dir); err == nil {
		t.Fatal("rewrote inventory")
	}
	var inventory map[string]fileEntry
	if err := load(filepath.Join(dir, "inventory.json"), &inventory); err != nil {
		t.Fatal(err)
	}
	if inventory["raw.json"] != hash([]byte("{}")) {
		t.Fatal("incorrect evidence hash")
	}
	linked := t.TempDir()
	if err := os.Symlink(path, filepath.Join(linked, "link.json")); err == nil {
		if err = createInventory(linked); err == nil {
			t.Fatal("accepted linked evidence")
		}
	}
}

func TestSourceOverlayPolicy(t *testing.T) {
	for path, want := range map[string]bool{"tools/research/context_fit.go": true, "go.mod": true, "third_party/chess/go.mod": true, "testdata/cases.json": true, "../secret.go": false, "/tmp/foreign.go": false, "helper.sh.txt": false, "third_party/other/go.mod": false} {
		if allowedOverlay(path) != want {
			t.Fatal(path, want)
		}
	}
}
