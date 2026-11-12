package main

import (
	"github.com/corentings/chess/v2"
	"strings"
	"testing"
)

func evalFEN(t *testing.T, fen string) *chess.Position {
	t.Helper()
	g, err := parsePosition(strings.Fields("fen " + fen))
	if err != nil {
		t.Fatal(err)
	}
	return g.Position()
}
func TestPlacementPreferences(t *testing.T) {
	for _, c := range []struct {
		name                                        string
		piece                                       chess.PieceType
		goodFile, goodRank, badFile, badRank, phase int
	}{
		{"central knight", chess.Knight, 3, 3, 0, 2, 24},
		{"bishop development", chess.Bishop, 2, 3, 2, 0, 24},
		{"central pawn space", chess.Pawn, 4, 3, 4, 1, 24},
		{"middlegame king shelter", chess.King, 6, 0, 4, 3, 24},
		{"endgame king activity", chess.King, 4, 3, 6, 0, 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			if placement(c.piece, c.goodFile, c.goodRank, c.phase) <= placement(c.piece, c.badFile, c.badRank, c.phase) {
				t.Fatal("positional preference reversed")
			}
		})
	}
	if placement(chess.Pawn, 0, 3, 24) != placement(chess.Pawn, 0, 1, 24) {
		t.Fatal("early flank-pawn rush received a progress bonus")
	}
}
func TestEvaluationColorSymmetry(t *testing.T) {
	pairs := [][2]string{
		{"4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", "4k3/8/8/3n4/8/8/8/4K3 b - - 0 1"},
		{"r2qk2r/8/8/8/4P3/2N5/8/R2QK2R w KQkq - 0 1", "r2qk2r/8/2n5/4p3/8/8/8/R2QK2R b KQkq - 0 1"},
	}
	for _, pair := range pairs {
		a, b := evalFEN(t, pair[0]), evalFEN(t, pair[1])
		if evaluate(a, chess.White) != -evaluate(a, chess.Black) {
			t.Fatal("perspective sign mismatch")
		}
		if evaluate(a, chess.White) != evaluate(b, chess.Black) {
			t.Fatal("color/rank mirror changed score")
		}
	}
	if evaluate(chess.StartingPosition(), chess.White) != 0 {
		t.Fatal("starting position is not symmetric")
	}
}
func TestQuietLeafUsesPlacement(t *testing.T) {
	central := evalFEN(t, "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1")
	edge := evalFEN(t, "4k3/8/8/8/N7/8/8/4K3 w - - 0 1")
	if materialScore(central, chess.White) != 300 || materialScore(edge, chess.White) != 300 {
		t.Fatal("material units changed")
	}
	if quiescence(central, 0, -searchInfinity, searchInfinity) <= quiescence(edge, 0, -searchInfinity, searchInfinity) {
		t.Fatal("quiet leaf ignored knight placement")
	}
}
func TestKingPhaseAndMateDominance(t *testing.T) {
	middle := placement(chess.King, 4, 3, 24)
	ending := placement(chess.King, 4, 3, 0)
	mixed := placement(chess.King, 4, 3, 12)
	if !(middle < mixed && mixed < ending) {
		t.Fatal("king phase interpolation is not bounded")
	}
	// Even an overestimate of every square's queen/material bonus is below mate.
	if mateScore <= 64*(pieceValues[chess.Queen]+200) {
		t.Fatal("positional/material evaluation could outweigh mate")
	}
}

func TestQuietMoveSelectionUsesPlacementForBothColors(t *testing.T) {
	for _, c := range []struct {
		fen, good, bad string
	}{
		{"7k/8/8/8/8/8/8/KN6 w - - 0 1", "b1c3", "b1a3"},
		{"kn6/8/8/8/8/8/8/7K b - - 0 1", "b8c6", "b8a6"},
	} {
		position := evalFEN(t, c.fen)
		moves := rootMoves(position, []string{c.good, c.bad})
		best, _, complete := (&searcher{}).root(position, moves, 1)
		if !complete || len(best) != 1 || best[0] != c.good {
			t.Fatalf("quiet choices %v, want uniquely best %s", best, c.good)
		}
	}
}
