package main

import (
	"math/rand/v2"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestExchangeBitsetsMatchMailboxReference(t *testing.T) {
	rng := rand.New(rand.NewPCG(22, 2271))
	scores, attacks, least := 0, 0, 0
	check := func(position *chess.Position, full bool) {
		t.Helper()
		board := newExchangeBoard(position)
		old := referenceExchangeBoard{}
		for s := chess.A1; s <= chess.H8; s++ {
			piece := position.Board().Piece(s)
			old.pieces[s] = piece
			if piece.Type() == chess.King {
				old.kings[piece.Color()] = s
			}
		}
		if board.pieces != old.pieces || board.kings != old.kings {
			t.Fatal("bitset constructor differs", position.String())
		}
		if full {
			for target := chess.A1; target <= chess.H8; target++ {
				for _, side := range [2]chess.Color{chess.White, chess.Black} {
					expected := uint64(0)
					for from, piece := range old.pieces {
						if piece.Color() == side && old.attacks(chess.Square(from), target) {
							expected |= uint64(1) << uint(63-from)
						}
					}
					if got := board.attackers(target, side); got != expected || board.attacked(target, side) != (expected != 0) {
						t.Fatal("attack set differs", position.String(), target, side, got, expected)
					}
					attacks++
					before, oldBefore := board, old
					got, want := board.leastAttacker(target, side), old.leastAttacker(target, side)
					if got != want || board != before || old != oldBefore {
						t.Fatal("least attacker/restoration differs", position.String(), target, side, got, want)
					}
					least++
				}
			}
		}
		before, hash := position.String(), position.ZobristHash()
		for _, m := range position.ValidMoves() {
			got, want := staticExchange(position, m), referenceStaticExchange(position, m)
			if got != want {
				t.Fatal("exchange score differs", position.String(), m, got, want)
			}
			scores++
		}
		if position.String() != before || position.ZobristHash() != hash {
			t.Fatal("exchange mutated parent")
		}
	}
	for _, fen := range []string{
		"4k3/4b3/3p4/8/8/Q7/8/4R1K1 w - - 0 1",
		"4k3/4r3/4p3/8/8/1Q6/8/4R1K1 w - - 0 1",
		"7k/7p/8/8/8/8/8/1B4KR w - - 0 1",
		"3r2k1/8/8/3pP3/8/8/8/6K1 w - d6 0 1",
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
		"k7/8/8/8/8/7R/6p1/K6r w - - 0 1",
		"4k3/8/8/8/8/8/4K3/8 w KQ - 0 1",
		"4k3/8/8/8/8/8/8/8 w - - 0 1",
		"4k3/8/8/8/8/8/8/3KK3 w - - 0 1",
	} {
		check(evalFEN(t, fen), true)
	}
	position := chess.StartingPosition()
	for i := 0; i < 1500; i++ {
		check(position, i < 250)
		moves := position.ValidMoves()
		if len(moves) == 0 || i%150 == 149 {
			position = chess.StartingPosition()
		} else {
			m := moves[rng.IntN(len(moves))]
			position = position.Update(&m)
		}
	}
	t.Log("exact full scores", scores, "attack sets", attacks, "least-attacker/restoration comparisons", least)
}
