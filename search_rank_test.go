package main

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestMoveOrderPermutationMatchesStableReference(t *testing.T) {
	base := chess.StartingPosition().ValidMoves()
	rng := rand.New(rand.NewPCG(21, 2101))
	for n := 0; n <= len(base); n++ {
		for trial := 0; trial < 100; trial++ {
			got := slices.Clone(base[:n])
			want := slices.Clone(got)
			priorities := map[string]int{}
			ranks := make([]moveRank, n)
			for i, m := range got {
				priority := rng.IntN(7) - 3 // Deliberately many ties and long cycles.
				priorities[m.String()] = priority
				ranks[i] = moveRank{index: i, priority: priority}
			}
			slices.SortStableFunc(want, func(a, b chess.Move) int {
				return priorities[b.String()] - priorities[a.String()]
			})
			slices.SortStableFunc(ranks, func(a, b moveRank) int {
				return b.priority - a.priority
			})
			applyMoveOrder(got, ranks)
			if !reflect.DeepEqual(got, want) {
				t.Fatal("permutation changed stable order or move contents", n, trial)
			}
		}
	}
	t.Log("Move bytes", reflect.TypeOf(chess.Move{}).Size(), "rank bytes", reflect.TypeOf(moveRank{}).Size())
}

func TestRankReuseKeepsExactMoveOrdering(t *testing.T) {
	worker := searcher{}
	positions := []*chess.Position{
		chess.StartingPosition(),
		evalFEN(t, "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"),
		evalFEN(t, "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1"),
		evalFEN(t, "r6k/1P6/8/8/8/8/8/7K w - - 0 1"),
		evalFEN(t, "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1"),
	}
	for trial := 0; trial < 10; trial++ {
		for _, position := range positions {
			for _, tactical := range []bool{false, true} {
				want := orderedMoves(position)
				if tactical {
					want = orderedTacticalMoves(position)
				}
				got := worker.searchMoves(position, 0, tactical)
				if !sameMoves(got, want) {
					t.Fatal("reused ranks changed ordering", position, tactical)
				}
				saved := slices.Clone(got)
				worker.searchMoves(chess.StartingPosition(), 30, false)
				if !reflect.DeepEqual(got, saved) {
					t.Fatal("other ply changed active move list")
				}
			}
		}
	}
}
