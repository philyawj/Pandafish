package main

import (
	"path/filepath"
	"testing"
)

func TestProbeExactScoreAndBoundControls(t *testing.T) {
	for _, tc := range []struct {
		line string
		ok   bool
		mate bool
	}{
		{"info depth 7 nodes 100 score cp -35 pv e2e4 e7e5", true, false},
		{"info depth 7 score cp 60 lowerbound nodes 100 pv e2e4", false, false},
		{"info depth 0 score mate 0 nodes 0", true, true},
		{"info string score cp 30", false, false},
	} {
		score, ok := parseProbeScore(tc.line)
		if ok != tc.ok {
			t.Fatal(tc.line, score, ok)
		}
		if tc.ok && tc.mate && (score.Mate == nil || *score.Mate != 0) {
			t.Fatal("terminal mate score lost")
		}
	}
}
func TestProbeFullHistoryAndTerminalMoveControls(t *testing.T) {
	g, err := probeGame(probeCase{Position: "position startpos moves e2e4 e7e5"})
	if err != nil {
		t.Fatal(err)
	}
	if !validProbeMove(g, "g1f3") || validProbeMove(g, "0000") || validProbeMove(g, "e2e4") {
		t.Fatal("live move controls")
	}
	g, err = probeGame(probeCase{Position: "position fen 7k/6Q1/5K2/8/8/8/8/8 b - - 0 1"})
	if err != nil {
		t.Fatal(err)
	}
	if !validProbeMove(g, "0000") || validProbeMove(g, "h8g8") {
		t.Fatal("terminal controls")
	}
}

func TestProbeCanonicalMirrorAndEnPassantControls(t *testing.T) {
	key := func(fen string) string {
		g, err := probeGame(probeCase{FEN: fen})
		if err != nil {
			t.Fatal(err)
		}
		return canonicalProbe(g.Position())
	}
	a := "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1"
	b := "4k3/8/8/8/3Pp3/8/8/4K3 b - d3 90 45"
	if key(a) != key(b) {
		t.Fatal("mirror/counter distinction")
	}
	if key(a) == key("4k3/8/8/3pP3/8/8/8/4K3 w - - 0 1") {
		t.Fatal("legal EP ignored")
	}
	if key("4k3/8/8/8/4P3/8/8/4K3 b - e3 0 1") != key("4k3/8/8/8/4P3/8/8/4K3 b - - 0 1") {
		t.Fatal("unusable EP retained")
	}
	if key("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1") == key("r3k2r/8/8/8/8/8/8/R3K2R w Qkq - 0 1") {
		t.Fatal("castling rights ignored")
	}
}

func TestProbeCanonicalAgainstArchivedKeys(t *testing.T) {
	checked := 0
	for _, name := range []string{"v0.19-decision-evaluation", "v0.19-independent-families", "v0.19-sharp-branches"} {
		var input struct {
			Cases []struct{ Key, Position, FEN string }
		}
		if err := load(filepath.Join("..", "benchmark", "research", name, "corpus.json"), &input); err != nil {
			t.Fatal(err)
		}
		for _, c := range input.Cases {
			if c.Key == "" {
				continue
			}
			g, err := probeGame(probeCase{Position: c.Position, FEN: c.FEN})
			if err != nil {
				t.Fatal(err)
			}
			if key := canonicalProbe(g.Position()); key != c.Key {
				t.Fatal("legacy canonical disagreement", c.FEN, key, c.Key)
			}
			checked++
		}
	}
	if checked < 1000 {
		t.Fatal("insufficient archived canonical controls", checked)
	}
}
