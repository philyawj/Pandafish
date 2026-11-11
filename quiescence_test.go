package main

import (
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestQuiescence(t *testing.T) {
	for _, tt := range []struct {
		name, fen string
		want      int
	}{
		{"quiet position", "7k/8/8/8/8/8/P7/7K w - - 0 1", 1},
		{"white recaptures queen", "6k1/8/8/8/3q4/8/8/3R2K1 w - - 0 1", 5},
		{"black recaptures queen", "3r2k1/8/8/3Q4/8/8/8/6K1 b - - 0 1", 5},
		{"white avoids poisoned pawn", "3r2k1/8/8/3p4/8/8/8/3Q2K1 w - - 0 1", 3},
		{"black avoids poisoned pawn", "3q2k1/8/8/8/3P4/8/8/3R2K1 b - - 0 1", 3},
		{"capture and recapture", "r6k/q7/8/8/8/8/8/R6K w - - 0 1", -5},
		{"en passant", "k7/8/8/3pP3/8/8/8/7K w - d6 0 1", 1},
		{"quiet white promotion", "7k/P7/8/8/8/8/8/K7 w - - 0 1", 9},
		{"quiet black promotion", "k7/8/8/8/8/8/p7/7K b - - 0 1", 9},
		{"capture promotion", "r6k/1P6/8/8/8/8/8/7K w - - 0 1", 9},
		{"rook promotion avoids stalemate", "8/k1P5/2K5/8/8/8/8/8 w - - 0 1", 5},
		{"quiet king evasions lose queen", "k6r/p7/8/8/8/8/1b6/Q6K w - - 0 1", -9},
		{"only escape is quiet block", "7r/8/8/8/8/6k1/8/6BK w - - 0 1", -5},
		{"mate", "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1", -mateScore},
		{"stalemate", "7k/5K2/6Q1/8/8/8/8/8 b - - 0 1", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(strings.Fields("fen " + tt.fen))
			if err != nil {
				t.Fatal(err)
			}
			before := game.FEN()
			if got := quiescence(game.Position(), quiescenceDepth, -searchInfinity, searchInfinity); got != tt.want {
				t.Fatalf("quiescence score = %d, want %d", got, tt.want)
			}
			if game.FEN() != before || len(game.Moves()) != 0 {
				t.Fatal("quiescence changed the root game")
			}
		})
	}
}

func TestQuiescenceMatchesExhaustiveOnTimeoutPositions(t *testing.T) {
	for _, fen := range []string{timeoutPosition, queenRaidPosition, crowdedPosition} {
		game, err := parsePosition(strings.Fields("fen " + fen))
		if err != nil {
			t.Fatal(err)
		}
		want := exhaustiveQuietScore(game.Position(), quiescenceDepth)
		if game.Position().Turn() == chess.Black {
			want = -want
		}
		if got := quiescence(game.Position(), quiescenceDepth, -searchInfinity, searchInfinity); got != want {
			t.Errorf("%s: quiescence score = %d, exhaustive score = %d", fen, got, want)
		}
	}
}

func TestQuiescenceCannotStandPatInCheck(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen k6r/p7/8/8/8/8/1b6/Q6K w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	// Static material is zero, but every legal escape lets the bishop take
	// our queen. A stand-pat cutoff at beta = 0 would hide that loss.
	if got := quiescence(game.Position(), quiescenceDepth, -searchInfinity, 0); got != -9 {
		t.Fatalf("checked score = %d, want -9", got)
	}
}

func TestQuiescenceDepthCap(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 3r2k1/8/8/3Q4/8/8/8/6K1 b - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	// At the hard cap we return material even if captures remain.
	if got := quiescence(game.Position(), 0, -searchInfinity, searchInfinity); got != -4 {
		t.Fatalf("capped score = %d, want -4", got)
	}
}

func TestInCheck(t *testing.T) {
	for _, tt := range []struct {
		fen  string
		want bool
	}{
		{"7k/8/8/8/8/8/P7/7K w - - 0 1", false},
		{"k6r/8/8/8/8/8/8/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/5n2/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/6p1/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/1b6/K7 w - - 0 1", true},
		{"7k/8/8/8/8/8/8/K6R b - - 0 1", true},
	} {
		game, err := parsePosition(strings.Fields("fen " + tt.fen))
		if err != nil {
			t.Fatal(err)
		}
		if got := inCheck(game.Position()); got != tt.want {
			t.Errorf("%s: check = %v, want %v", tt.fen, got, tt.want)
		}
	}
}
