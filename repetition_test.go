package main

import (
	"context"
	"github.com/corentings/chess/v2"
	"strings"
	"testing"
	"time"
)

func TestThirdOccurrenceScoresDrawAndUnwinds(t *testing.T) {
	p := evalFEN(t, "4k3/8/8/8/8/8/8/3QK3 w - - 0 1")
	w := searcher{table: newTable()}
	w.pushPosition(p)
	w.pushPosition(p)
	before := w.historyKey
	for _, depth := range []int{0, 3} {
		score, ok := w.search(p, depth, -searchInfinity, searchInfinity, false)
		if !ok || score != 0 {
			t.Fatalf("third repetition depth%d score%d", depth, score)
		}
		if w.historyKey != before || w.repetitions[p.ZobristHash()] != 2 {
			t.Fatal("draw check leaked path state")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	w.ctx = ctx
	if _, ok := w.search(p, 3, -searchInfinity, searchInfinity, false); ok {
		t.Fatal("cancel became draw score")
	}
	if w.historyKey != before {
		t.Fatal("cancel changed history")
	}
}
func TestHistoryContextSeparatesCachedScores(t *testing.T) {
	p := chess.StartingPosition()
	a, b := searcher{}, searcher{}
	a.pushPosition(p)
	b.pushPosition(p)
	b.pushPosition(p)
	if a.historyKey == b.historyKey {
		t.Fatal("counts not represented")
	}
	table := newTable()
	table.entries = make([]tableEntry, 1)
	table.store(p, 3, 300, -500, 500, 0, a.historyKey)
	if _, ok := table.probe(p, b.historyKey); ok {
		t.Fatal("draw context reused cached score")
	}
	if _, ok := table.probe(p, a.historyKey); !ok {
		t.Fatal("original context lost")
	}
	a.popPosition(p)
	if a.historyKey != 0 || len(a.repetitions) != 0 {
		t.Fatal("pop did not restore empty state")
	}
	// The draw rule depends on counts, not the order of the earlier positions.
	q := evalFEN(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	a.pushPosition(p)
	a.pushPosition(q)
	a.pushPosition(p)
	c := searcher{}
	c.pushPosition(q)
	c.pushPosition(p)
	c.pushPosition(p)
	if a.historyKey != c.historyKey {
		t.Fatal("same count set hashed differently")
	}
}
func TestKnownHistoryFindsDrawingMove(t *testing.T) {
	g, err := parsePosition(strings.Fields("fen 4k1n1/8/8/8/8/8/8/3QK1N1 w - - 0 1 moves g1f3 g8f6 f3g1 f6g8 g1f3 g8f6 f3g1"))
	if err != nil {
		t.Fatal(err)
	}
	history := repetitionHistory(g)
	if len(history) != 7 {
		t.Fatalf("history size%d", len(history))
	}
	r := iterativeSearchWithHistory(context.Background(), g.Position(), goLimits{depth: 2}, timeBudget{}, time.Now(), newTable(), history)
	if r.score != 0 || r.move != "f6g8" {
		t.Fatalf("missed available repetition: %+v", r)
	}
	// A pawn move makes all earlier positions unreachable by repetition.
	reset, err := parsePosition(strings.Fields("startpos moves g1f3 g8f6 f3g1 f6g8 e2e4"))
	if err != nil {
		t.Fatal(err)
	}
	if len(repetitionHistory(reset)) != 0 {
		t.Fatal("irreversible history retained")
	}
}
func TestMatePrecedesRepetition(t *testing.T) {
	p := evalFEN(t, "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1")
	w := searcher{}
	w.pushPosition(p)
	w.pushPosition(p)
	if score, ok := w.search(p, 3, -searchInfinity, searchInfinity, true); !ok || score != -mateScore {
		t.Fatal("mate replaced by draw")
	}
}
