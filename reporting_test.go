package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/philyawj/Pandafish/internal/telemetry"
)

func TestRecordedSearchMetricsAndProtocol(t *testing.T) {
	dir := t.TempDir()
	trace, err := telemetry.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	// Registered first so the UCI loop is joined before the file is closed.
	t.Cleanup(func() {
		if err := trace.Close(); err != nil {
			t.Error(err)
		}
	})
	session := newRecordedUCISession(t, trace)
	session.send("position startpos moves e2e4")
	session.send("go depth 2 searchmoves e7e5")
	var metrics telemetry.Metrics
	for {
		select {
		case line := <-session.lines:
			if strings.HasPrefix(line, telemetry.Prefix) {
				if err = json.Unmarshal([]byte(strings.TrimPrefix(line, telemetry.Prefix)), &metrics); err != nil {
					t.Fatal(err)
				}
			}
			if strings.HasPrefix(line, "bestmove ") {
				if line != "bestmove e7e5" || metrics.Depth != 2 || metrics.Nodes == 0 || metrics.Fallback || metrics.ScoreKind != "cp" || metrics.Reason != "depth" || metrics.Timed || !slices.Equal(metrics.PV, []string{"e7e5"}) {
					t.Fatalf("move %s metrics %+v", line, metrics)
				}
				// Inspect only after a command-loop barrier to avoid racing the trace writer.
				session.send("isready")
				if session.read() != "readyok" {
					t.Fatal("missing ready")
				}
				files, err := filepath.Glob(filepath.Join(dir, "engine-*.jsonl"))
				if err != nil || len(files) != 1 {
					t.Fatal(files, err)
				}
				raw, err := os.ReadFile(files[0])
				if err != nil {
					t.Fatal(err)
				}
				scanner := bufio.NewScanner(strings.NewReader(string(raw)))
				var in, out []string
				for scanner.Scan() {
					var e telemetry.Event
					if err = json.Unmarshal(scanner.Bytes(), &e); err != nil {
						t.Fatal(err)
					}
					if e.Kind == "in" {
						in = append(in, e.Line)
					}
					if e.Kind == "out" {
						out = append(out, e.Line)
					}
				}
				if !slices.Contains(in, "go depth 2 searchmoves e7e5") || !strings.Contains(strings.Join(out, ""), "bestmove e7e5") {
					t.Fatal("request or result missing from trace")
				}
				return
			}
		case <-time.After(3 * time.Second):
			t.Fatal("timed out")
		}
	}
}
func TestReportingDoesNotChangeFixedDepthScores(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen " + queenRaidPosition))
	if err != nil {
		t.Fatal(err)
	}
	best, score, complete := (&searcher{}).root(game.Position(), rootMoves(game.Position(), nil), 2)
	if !complete {
		t.Fatal("reference interrupted")
	}
	result := iterativeSearch(context.Background(), game.Position(), goLimits{depth: 2}, timeBudget{}, time.Now())
	if result.score != score || !slices.Contains(best, result.move) || result.metrics.Score != score {
		t.Fatalf("result %+v, reference %v %d", result, best, score)
	}
	var output strings.Builder
	reportSearch(&output, result)
	if !strings.Contains(output.String(), "info depth 2 nodes ") || !strings.Contains(output.String(), " score cp ") {
		t.Fatal(output.String())
	}
}
func TestFallbackAndMateReporting(t *testing.T) {
	game, err := parsePosition(strings.Fields("startpos"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := iterativeSearch(ctx, game.Position(), goLimits{searchMoves: []string{"e2e4"}}, timeBudget{}, time.Now())
	if !result.metrics.Fallback || result.metrics.ScoreKind != "unknown" || result.metrics.Reason != "cancelled" {
		t.Fatal(result.metrics)
	}
	var output strings.Builder
	reportSearch(&output, result)
	if strings.Contains(output.String(), "score cp") {
		t.Fatal("fallback given a searched score")
	}
	result.metrics = telemetry.Metrics{Depth: 1, ScoreKind: "mate", Score: 1}
	output.Reset()
	reportSearch(&output, result)
	if strings.Contains(output.String(), "score cp") || strings.Contains(output.String(), "score mate") {
		t.Fatal("invented mate distance")
	}
}
