package main

import (
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

func TestTimeBudgets(t *testing.T) {
	for _, tt := range []struct {
		args  string
		side  chess.Color
		hard  time.Duration
		timed bool
	}{
		{"", chess.White, 95 * time.Millisecond, true},
		{"movetime 10", chess.White, 9 * time.Millisecond, true},
		{"movetime 0", chess.White, 0, true},
		{"depth 3", chess.White, 0, false},
		{"depth 3 movetime 100", chess.White, 95 * time.Millisecond, true},
		{"wtime 1000 btime 3000", chess.White, 29700 * time.Microsecond, true},
		{"wtime 1000 btime 3000 binc 100 movestogo 10", chess.Black, 369 * time.Millisecond, true},
		{"wtime 1000 movestogo 1", chess.White, 490 * time.Millisecond, true},
		{"wtime 100 movetime 1000", chess.White, 85 * time.Millisecond, true},
		{"wtime 1", chess.White, 27 * time.Microsecond, true},
		{"wtime 0 winc 10000", chess.White, 0, true},
		{"infinite movetime 10", chess.White, 0, false},
		{"ponder wtime 1000", chess.White, 0, false},
	} {
		t.Run(tt.args, func(t *testing.T) {
			limits, err := parseGo(strings.Fields(tt.args))
			if err != nil {
				t.Fatal(err)
			}
			budget, err := limits.budget(tt.side)
			if err != nil {
				t.Fatal(err)
			}
			if budget.hard != tt.hard || budget.timed != tt.timed || budget.soft != tt.hard/4*3 {
				t.Fatalf("budget = %+v, want hard %v, timed %v", budget, tt.hard, tt.timed)
			}
		})
	}
}

func TestInvalidGoLimits(t *testing.T) {
	for _, args := range []string{"movetime", "movetime -1", "winc -1", "depth 0", "depth 65", "movestogo 0", "btime bad", "movetime 9223372036854775807"} {
		if _, err := parseGo(strings.Fields(args)); err == nil {
			t.Errorf("accepted %q", args)
		}
	}
	limits, _ := parseGo(strings.Fields("wtime 100"))
	if _, err := limits.budget(chess.Black); err == nil {
		t.Error("accepted missing black clock")
	}
}

func TestGoLimitsSearchMovesAndLargeClocks(t *testing.T) {
	args := strings.Fields("movetime 1 movetime 2 searchmoves e2e4 d2d4")
	limits, err := parseGo(args)
	if err != nil {
		t.Fatal(err)
	}
	args[len(args)-1] = "a2a3"
	if limits.moveTime != 2*time.Millisecond || strings.Join(limits.searchMoves, " ") != "e2e4 d2d4" {
		t.Fatal(limits)
	}
	empty, _ := parseGo([]string{"searchmoves"})
	if empty.searchMoves == nil || len(empty.searchMoves) != 0 {
		t.Fatal("empty restriction lost")
	}
	large, err := parseGo(strings.Fields("wtime 9223372036854 winc 9223372036854 movestogo 1"))
	if err != nil {
		t.Fatal(err)
	}
	budget, err := large.budget(chess.White)
	if err != nil || budget.hard < 0 || budget.soft < 0 || budget.soft > budget.hard || budget.hard > large.clocks[0]/2 {
		t.Fatalf("overflowed budget: %+v, %v", budget, err)
	}
}
