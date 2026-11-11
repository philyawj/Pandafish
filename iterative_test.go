package main

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

func TestIterativeMatchesFixedDepth(t *testing.T) {
	for _, command := range []string{"startpos", "fen " + queenRaidPosition, "fen 4k3/1q6/8/5N2/8/8/8/K7 w - - 0 1"} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		moves := rootMoves(game.Position(), nil)
		best, score, complete := (&searcher{}).root(game.Position(), moves, 3)
		if !complete {
			t.Fatal("unlimited reference interrupted")
		}
		before := game.FEN()
		got := iterativeSearch(context.Background(), game.Position(), goLimits{depth: 3}, timeBudget{}, time.Now())
		if got.depth != 3 || got.score != score || !slices.Contains(best, got.move) {
			t.Fatalf("%s: result %+v, best %v / %d", command, got, best, score)
		}
		if game.FEN() != before {
			t.Fatal("search changed root")
		}
	}
}

func TestCancellationAndFallback(t *testing.T) {
	position := chess.StartingPosition()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	worker := searcher{ctx: ctx}
	if _, complete := worker.search(position, 3, -searchInfinity, searchInfinity, false); complete {
		t.Error("main search ignored cancellation")
	}
	if _, complete := worker.quiet(position, 2, -searchInfinity, searchInfinity, false); complete {
		t.Error("quiescence ignored cancellation")
	}
	result := iterativeSearch(ctx, position, goLimits{searchMoves: []string{"e2e4"}}, timeBudget{}, time.Now())
	if result.move != "e2e4" || result.depth != 0 {
		t.Fatalf("fallback = %+v", result)
	}
}

func TestIncompleteIterationDoesNotReplaceCompletedResult(t *testing.T) {
	position := chess.StartingPosition()
	start := time.Unix(100, 0)
	moves := rootMoves(position, []string{"e2e4", "d2d4"})
	// Count checks for one completed iteration, then expire during the next.
	calls := 0
	counting := searcher{deadline: start.Add(time.Second), now: func() time.Time { calls++; return start }}
	best, score, complete := counting.root(position, moves, 1)
	if !complete {
		t.Fatal("reference incomplete")
	}
	limit := calls + 4 // iteration check, complete depth one, then part of depth two.
	calls = 0
	worker := searcher{deadline: start.Add(time.Second), now: func() time.Time {
		calls++
		if calls >= limit {
			return start.Add(time.Second)
		}
		return start
	}}
	got := worker.iterate(position, moves, 3, timeBudget{}, start)
	if got.depth != 1 || got.score != score || !slices.Contains(best, got.move) {
		t.Fatalf("partial iteration replaced result: %+v", got)
	}
}

func TestDeadlineInsideQuiescence(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen " + queenRaidPosition))
	if err != nil {
		t.Fatal(err)
	}
	start := time.Unix(100, 0)
	calls := 0
	worker := searcher{deadline: start.Add(time.Second), now: func() time.Time {
		calls++
		if calls > 3 {
			return start.Add(time.Second)
		}
		return start
	}}
	if _, complete := worker.quiet(game.Position(), quiescenceDepth, -searchInfinity, searchInfinity, inCheck(game.Position())); complete {
		t.Fatal("quiescence converted interrupted work into a score")
	}
}

func TestTimedUCIMove(t *testing.T) {
	session := newUCISession(t)
	started := time.Now()
	session.send("go movetime 20")
	move := session.read()
	legal := false
	for _, candidate := range chess.StartingPosition().ValidMoves() {
		if move == "bestmove "+candidate.String() {
			legal = true
		}
	}
	if !legal {
		t.Fatal(move)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("20 ms request took %v", elapsed)
	}
}

func TestSoftTargetKeepsCompletedIteration(t *testing.T) {
	position := chess.StartingPosition()
	moves := rootMoves(position, []string{"e2e4", "d2d4"})
	start := time.Unix(100, 0)
	calls := 0
	counting := searcher{deadline: start.Add(time.Second), now: func() time.Time { calls++; return start }}
	best, score, _ := counting.root(position, moves, 1)
	checks := calls
	calls = 0
	worker := searcher{deadline: start.Add(time.Second), now: func() time.Time {
		calls++
		if calls > checks+2 {
			return start.Add(750 * time.Millisecond)
		}
		return start
	}}
	got := worker.iterate(position, moves, 3, timeBudget{timed: true, soft: 750 * time.Millisecond, hard: time.Second}, start)
	if got.depth != 1 || got.score != score || !slices.Contains(best, got.move) {
		t.Fatalf("soft target discarded completed work: %+v", got)
	}
}
