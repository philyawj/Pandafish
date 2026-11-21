package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestParallelWarmReportSeparatesSchedulingFromValues(t *testing.T) {
	for _, badScore := range []bool{false, true} {
		t.Run(fmt.Sprint(badScore), func(t *testing.T) {
			dir := t.TempDir()
			base := warmReportRow{Game: 1, Ply: 8, Depth: 5, Best: "e2e4", FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", Nodes: 100, SearchMS: 10}
			for pair := 1; pair <= 3; pair++ {
				candidate := base
				candidate.Nodes += uint64(pair)
				candidate.SearchMS = 7
				candidate.Best = "d2d4"
				if pair == 2 {
					candidate.Best = "g1f3"
				}
				if badScore && pair == 3 {
					candidate.Score = 1
				}
				for name, row := range map[string]warmReportRow{"base": base, "candidate": candidate} {
					if err := writeNewJSON(filepath.Join(dir, fmt.Sprintf("%s-%d.json", name, pair)), map[string]any{"Rows": []warmReportRow{row}}); err != nil {
						t.Fatal(err)
					}
				}
			}
			o := warmReportOptions{Baseline: filepath.Join(dir, "base-*.json"), Candidate: filepath.Join(dir, "candidate-*.json"), Out: filepath.Join(dir, "parallel.json"), Pairs: 3, Roots: 1, Depth: 5, Parallel: true, MaxTimeRatio: .85}
			if err := warmReport(o); err != nil {
				t.Fatal(err)
			}
			var report struct{ Passed bool }
			if err := load(o.Out, &report); err != nil {
				t.Fatal(err)
			}
			if report.Passed == badScore {
				t.Fatal("parallel value gate", report)
			}
			o.Out = filepath.Join(dir, "ordinary.json")
			o.Parallel = false
			if err := warmReport(o); err != nil {
				t.Fatal(err)
			}
			if err := load(o.Out, &report); err != nil {
				t.Fatal(err)
			}
			if report.Passed {
				t.Fatal("ordinary determinism gate was relaxed")
			}
		})
	}
}

func TestParallelWarmPlanRequiresChoiceValidation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	o := warmPlanOptions{BaselineTest: "base.test", CandidateTest: "candidate.test", Input: "trajectory", Evidence: "evidence", Out: path, Roots: 1145, Depth: 5, MaxTimeRatio: .85, Parallel: true}
	if err := createWarmPlan(o); err != nil {
		t.Fatal(err)
	}
	var plan workflowPlan
	if err := load(path, &plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps) != 8 || plan.Steps[7].Name != "choice-validation" {
		t.Fatal("missing independent choice check")
	}
	if commandFlag(plan.Steps[6].Command, "--parallel") == "" {
		found := false
		for _, arg := range plan.Steps[6].Command {
			if arg == "--parallel" {
				found = true
			}
		}
		if !found {
			t.Fatal("parallel report profile missing")
		}
	}
}
