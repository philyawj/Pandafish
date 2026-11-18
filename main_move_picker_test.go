package main

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/corentings/chess/v2"
)

var pickerFixtures = []string{
	chess.StartingPosition().String(), queenRaidPosition,
	"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
	"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
	"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
	"4k3/8/8/8/3Pp3/8/8/4K3 b - d3 0 1",
	"k3r3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
	"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
	"7k/8/8/8/8/8/1p6/R6K b - - 0 1",
	"7r/8/8/8/8/6k1/8/6BK w - - 0 1",
	"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1",
	"7k/6Q1/5K2/8/8/8/8/8 b - - 0 1",
	"4k3/8/8/8/8/8/4K3/8 w KQ - 0 1",
	"4k3/8/8/8/8/8/4K3/4Q3 w KQ - 0 1",
	"4q3/4k3/8/8/8/8/8/4K3 b kq - 0 1",
}

func TestSearchHintAndQuietSubsetAgainstOriginal(t *testing.T) {
	rng := rand.New(rand.NewPCG(22, 2251))
	checked := 0
	check := func(fen string) {
		t.Helper()
		for _, warm := range []bool{false, true} {
			p := evalFEN(t, fen)
			reference := evalFEN(t, fen)
			want := reference.ValidMoves()
			if warm {
				want = p.ValidMoves()
			}
			before, hash := p.String(), p.ZobristHash()
			cache := slices.Clone(want)
			legal := map[uint16]chess.Move{}
			quiet := []chess.Move{}
			for _, m := range want {
				legal[encodeMove(m)] = m
				if isQuietMove(m) {
					quiet = append(quiet, m)
				}
				got, ok := p.LookupSearchMove(m.S1(), m.S2(), m.Promo())
				if !ok || !reflect.DeepEqual(got, m) {
					t.Fatal("direct legal lookup differs", fen, warm, m, got, ok)
				}
				checked++
			}
			for i := 0; i < 300; i++ {
				hint := uint16(rng.IntN(65536))
				got, ok := lookupHint(p, hint)
				expected, exists := legal[hint]
				if ok != exists || (ok && !reflect.DeepEqual(got, expected)) {
					t.Fatal("stale/malformed lookup", fen, warm, hint, got, ok, exists)
				}
				checked++
			}
			prefix := []chess.Move{{}}
			got := p.AppendSearchQuietMoves(prefix)
			if !reflect.DeepEqual(got[0], chess.Move{}) || !sameMoves(got[1:], quiet) {
				t.Fatal("quiet subset/prefix", fen, warm)
			}
			if p.String() != before || p.ZobristHash() != hash {
				t.Fatal("lookup/quiet generation mutated parent")
			}
			if warm && !reflect.DeepEqual(p.ValidMoves(), cache) {
				t.Fatal("cached move metadata changed")
			}
			for _, bad := range [][3]int{{-1, 0, 0}, {0, 64, 0}, {0, 0, 0}, {0, 16, 15}} {
				if _, ok := p.LookupSearchMove(chess.Square(bad[0]), chess.Square(bad[1]), chess.PieceType(bad[2])); ok {
					t.Fatal("invalid hint accepted", bad)
				}
			}
		}
	}
	for _, fen := range pickerFixtures {
		check(fen)
	}
	p := chess.StartingPosition()
	for i := 0; i < 600; i++ {
		check(p.String())
		moves := p.ValidMoves()
		if len(moves) == 0 || i%100 == 99 {
			p = chess.StartingPosition()
		} else {
			m := moves[rng.IntN(len(moves))]
			p = p.Update(&m)
		}
	}
	t.Log("Lookup/subset controls:", checked)
}

func TestStagedPickerMatchesEveryEagerMove(t *testing.T) {
	rng := rand.New(rand.NewPCG(22, 2252))
	comparisons := 0
	check := func(fen string) {
		t.Helper()
		for _, warm := range []bool{false, true} {
			p := evalFEN(t, fen)
			if warm {
				p.ValidMoves()
			}
			base := (&searcher{}).searchMoves(p, 3, false)
			for trial := 0; trial < 10; trial++ {
				hint := uint16(rng.IntN(65536))
				killers := new(killerTable)
				if trial == 0 {
					hint = 0
				}
				if len(base) > 0 {
					if trial%3 == 1 {
						hint = encodeMove(base[trial%len(base)])
					}
					killers[3] = [2]uint16{encodeMove(base[(trial+1)%len(base)]), encodeMove(base[(trial+2)%len(base)])}
					if trial%3 == 2 {
						killers[3][1] = killers[3][0]
					}
				}
				ref := searcher{killers: killers}
				want := slices.Clone(ref.searchMoves(p, 3, false))
				ref.killers.order(want, 3)
				preferMove(want, hint)
				worker := searcher{killers: killers}
				picker := newMainMovePicker(&worker, p, 3, hint)
				got := []chess.Move{}
				for m, ok := picker.next(); ok; m, ok = picker.next() {
					got = append(got, m)
					if len(got) > len(want) {
						t.Fatal("picker duplicate/loop", fen)
					}
				}
				if !reflect.DeepEqual(got, append([]chess.Move{}, want...)) {
					t.Fatal("staged whole-Move order differs", fen, warm, trial)
				}
				comparisons++
			}
		}
	}
	for _, fen := range pickerFixtures {
		check(fen)
	}
	p := chess.StartingPosition()
	for i := 0; i < 500; i++ {
		check(p.String())
		moves := p.ValidMoves()
		if len(moves) == 0 || i%80 == 79 {
			p = chess.StartingPosition()
		} else {
			m := moves[rng.IntN(len(moves))]
			p = p.Update(&m)
		}
	}
	t.Log("Whole ordered lists matched:", comparisons)
}

func TestStagedPickerDefersBuffersAndIsolatesActivePly(t *testing.T) {
	p := chess.StartingPosition()
	first := p.AppendSearchMoves(nil, false)[0]
	w := searcher{}
	picker := newMainMovePicker(&w, p, 0, encodeMove(first))
	if got, ok := picker.next(); !ok || !reflect.DeepEqual(got, first) || len(w.moveBuffers) != 0 || len(w.rankBuffers) != 0 {
		t.Fatal("hint materialized later-stage buffers")
	}
	if allocs := testing.AllocsPerRun(1000, func() {
		worker := searcher{}
		p := newMainMovePicker(&worker, p, 0, encodeMove(first))
		m, ok := p.next()
		if !ok || m.S1() != first.S1() {
			panic("hint control")
		}
	}); allocs != 0 {
		t.Fatal("preferred-only picker allocated", allocs)
	}
	p = evalFEN(t, queenRaidPosition)
	ref := searcher{}
	want := slices.Clone(ref.searchMoves(p, 1, false))
	w = searcher{}
	picker = newMainMovePicker(&w, p, 1, 0)
	got := []chess.Move{}
	for m, ok := picker.next(); ok; m, ok = picker.next() {
		got = append(got, m)
		child := newMainMovePicker(&w, chess.StartingPosition(), 90, 0)
		for _, ok := child.next(); ok; _, ok = child.next() {
		}
	}
	if !reflect.DeepEqual(got, append([]chess.Move{}, want...)) {
		t.Fatal("descendant buffer growth changed parent stage")
	}
}

func TestStagedPickerPerft(t *testing.T) {
	w := searcher{}
	var count func(*chess.Position, int, int) uint64
	count = func(p *chess.Position, depth, ply int) uint64 {
		if depth == 0 {
			return 1
		}
		picker := newMainMovePicker(&w, p, ply, 0)
		var n uint64
		for m, ok := picker.next(); ok; m, ok = picker.next() {
			n += count(p.UpdateSearch(&m), depth-1, ply+1)
		}
		return n
	}
	if got := count(chess.StartingPosition(), 4, 0); got != 197281 {
		t.Fatal("staged perft", got)
	}
}
