package main

import (
	"math/rand/v2"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestSearchUpdateMatchesOriginal(t *testing.T) {
	checked := 0
	check := func(p *chess.Position) {
		t.Helper()
		before, hash := p.String(), p.ZobristHash()
		pieces := p.Board().EvaluationPieces()
		cold := evalFEN(t, before)
		generated := cold.AppendSearchMoves(nil, false)
		want := p.ValidMoves()
		if !sameMoves(generated, want) {
			t.Fatal("move ordering/tags", before)
		}
		for i, move := range generated {
			for _, tag := range []chess.MoveTag{chess.Capture, chess.EnPassant, chess.Check, chess.KingSideCastle, chess.QueenSideCastle} {
				if move.HasTag(tag) != want[i].HasTag(tag) {
					t.Fatal("tag mismatch", before, move.String(), tag)
				}
			}
			a, b := cold.UpdateSearch(&move), p.Update(&want[i])
			if a.String() != b.String() || a.ZobristHash() != b.ZobristHash() || a.Board().EvaluationPieces() != b.Board().EvaluationPieces() {
				t.Fatal("child state/hash/pieces differ", before, move.String())
			}
			for square := chess.A1; square <= chess.H8; square++ {
				if a.Board().Piece(square) != b.Board().Piece(square) {
					t.Fatal("mailbox differs")
				}
			}
			if a.SearchStatus() != b.Status() || inCheck(a) != inCheck(b) {
				t.Fatal("child check/status differs", before, move.String())
			}
			if !sameMoves(a.AppendSearchMoves(nil, false), b.ValidMoves()) {
				t.Fatal("child legal moves differ", before, move.String())
			}
			checked++
		}
		if cold.String() != before || cold.ZobristHash() != hash || cold.Board().EvaluationPieces() != pieces {
			t.Fatal("parent mutated")
		}
		if a, b := cold.UpdateSearch(nil), p.Update(nil); a.String() != b.String() || a.ZobristHash() != b.ZobristHash() {
			t.Fatal("null transition differs")
		}
	}
	for _, fen := range []string{
		chess.StartingPosition().String(),
		"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"4k3/8/8/8/3Pp3/8/8/4K3 b - d3 0 1",
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
		"7k/8/8/8/8/8/1p6/R6K b - - 0 1",
		"k3r3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"7k/6Q1/5K2/8/8/8/8/8 b - - 0 1",
	} {
		check(evalFEN(t, fen))
	}
	rng := rand.New(rand.NewPCG(20, 2002))
	p := chess.StartingPosition()
	for i := 0; i < 1200; i++ {
		check(p)
		moves := p.ValidMoves()
		if len(moves) == 0 || i%100 == 99 {
			p = chess.StartingPosition()
		} else {
			m := moves[rng.IntN(len(moves))]
			p = p.Update(&m)
		}
	}
	if checked < 15000 {
		t.Fatal("insufficient transitions", checked)
	}
	t.Log("Matched original transitions:", checked)
}
