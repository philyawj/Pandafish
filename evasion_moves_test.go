package main

import (
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestEvasionPrefilterAgainstOriginalMoves(t *testing.T) {
	fixtures := []string{
		"4r1k1/8/8/8/8/8/3B4/4K3 w - - 0 1",
		"4r1k1/8/8/8/8/2N5/8/4K3 w - - 0 1",
		"4r1k1/8/8/8/1b6/8/8/1N2K3 w - - 0 1",
		"6k1/8/8/3pP3/4K3/8/8/8 w - d6 0 1",
		"5bk1/8/8/3pP3/8/K7/8/8 w - d6 0 1",
		"8/8/8/3k4/3pP3/8/8/6K1 b - e3 0 1",
		"k6r/6P1/8/8/8/8/8/7K w - - 0 1",
		"4k3/8/8/8/8/8/4K3/8 w KQ - 0 1",
		"4k3/8/8/8/8/8/4K3/4Q3 w KQ - 0 1",
	}
	comparisons := 0
	check := func(fen string) {
		t.Helper()
		for _, warm := range []bool{false, true} {
			p, original := evalFEN(t, fen), evalFEN(t, fen)
			want := original.ValidMoves()
			if warm {
				want = p.ValidMoves()
			}
			before, hash := p.String(), p.ZobristHash()
			for _, tactical := range []bool{false, true} {
				expected := []chess.Move{}
				for _, move := range want {
					if !tactical || !isQuietMove(move) {
						expected = append(expected, move)
					}
				}
				got := p.AppendSearchMoves(nil, tactical)
				if !sameMoves(got, expected) {
					t.Fatal("evasion ordered moves differ", fen, warm, tactical, got, expected)
				}
				for i := range got {
					if !reflect.DeepEqual(got[i], expected[i]) {
						t.Fatal("evasion metadata differs", fen, warm, tactical)
					}
				}
				comparisons++
			}
			if p.HasLegalMove() != (len(want) > 0) || p.String() != before || p.ZobristHash() != hash {
				t.Fatal("evasion existence/ownership", fen)
			}
		}
	}
	for _, fen := range fixtures {
		check(fen)
	}
	rng := rand.New(rand.NewPCG(22, 2261))
	p := chess.StartingPosition()
	for i := 0; i < 1500; i++ {
		check(p.String())
		moves := p.ValidMoves()
		if len(moves) == 0 || i%150 == 149 {
			p = chess.StartingPosition()
		} else {
			// Prefer checks to exercise both single and discovered-check evasions.
			candidates := moves
			checking := []chess.Move{}
			for _, move := range moves {
				if move.HasTag(chess.Check) {
					checking = append(checking, move)
				}
			}
			if len(checking) > 0 && i%3 != 0 {
				candidates = checking
			}
			m := candidates[rng.IntN(len(candidates))]
			p = p.Update(&m)
		}
	}
	t.Log("full/tactical exact evasion comparisons", comparisons)
}

func TestEvasionEnPassantAndBlocksRemainLegal(t *testing.T) {
	for _, c := range []struct{ fen, move string }{
		{"6k1/8/8/3pP3/4K3/8/8/8 w - d6 0 1", "e5d6"},
		{"5bk1/8/8/3pP3/8/K7/8/8 w - d6 0 1", "e5d6"},
		{"8/8/8/3k4/3pP3/8/8/6K1 b - e3 0 1", "d4e3"},
		{"4r1k1/8/8/8/8/8/3B4/4K3 w - - 0 1", "d2e3"},
	} {
		p := evalFEN(t, c.fen)
		found := false
		for _, m := range p.AppendSearchMoves(nil, false) {
			if m.String() == c.move {
				found = true
			}
		}
		if !found {
			t.Fatal("legal evasion lost", c)
		}
	}
}
