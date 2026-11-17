package main

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/corentings/chess/v2"
)

// Independently score every root alternative with its own full window and no
// cache. This enumerates ties for diagnostics; production only proves one best.
func fullWindowRoot(t *testing.T, p *chess.Position, moves []chess.Move, depth int, history []uint64) ([]string, int) {
	t.Helper()
	best, value := []string{}, -searchInfinity
	for _, m := range moves {
		w := searcher{}
		for _, k := range history {
			w.pushKey(k)
		}
		w.pushPosition(p)
		score, ok := w.childSearch(p.Update(&m), depth-1, -searchInfinity, searchInfinity, m.HasTag(chess.Check), 1)
		if !ok {
			t.Fatal("unlimited full-window reference interrupted")
		}
		score = -score
		if score > value {
			value, best = score, best[:0]
		}
		if score == value {
			best = append(best, m.String())
		}
	}
	return best, value
}

func TestOneBestRootContractAcrossOrders(t *testing.T) {
	for _, fen := range []string{
		"7k/5Q2/5K2/8/8/8/8/8 w - - 0 1",
		"7k/8/5Q2/5K2/8/8/8/8 b - - 0 1",
		"4k3/8/8/8/8/8/8/3QK3 w - - 0 1",
		"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
	} {
		p := evalFEN(t, fen)
		moves := rootMoves(p, nil)
		best, score := fullWindowRoot(t, p, moves, 3, nil)
		table := newTable()
		for i := 0; i < min(8, len(moves)); i++ {
			moves = append(moves[1:], moves[0])
			for _, tt := range []*transpositionTable{nil, table} {
				got, value, ok := (&searcher{table: tt, killers: new(killerTable)}).root(p, moves, 3)
				if !ok || value != score || len(got) != 1 || !slices.Contains(best, got[0]) {
					t.Fatalf("%s order%d best%v score%d want%v/%d", fen, i, got, value, best, score)
				}
			}
		}
	}
}

func TestWarmScoresRespectActualDrawHistories(t *testing.T) {
	// The same board has a losing value without history and an available drawing
	// repetition with history. Warm both directions, including across generations.
	for _, command := range []string{
		"fen 4k1n1/8/8/8/8/8/8/3QK1N1 w - - 0 1 moves g1f3 g8f6 f3g1 f6g8 g1f3 g8f6 f3g1",
		"fen 3qk1n1/8/8/8/8/8/8/4K1N1 b - - 0 1 moves g8f6 g1f3 f6g8 f3g1 g8f6 g1f3 f6g8",
	} {
		g, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		p := g.Position()
		history := repetitionHistory(g)
		run := func(tt *transpositionTable, h []uint64) int {
			w := searcher{table: tt}
			for _, k := range h {
				w.pushKey(k)
			}
			before := w.historyKey
			counts := maps.Clone(w.repetitions)
			s, ok := w.search(p, 2, -searchInfinity, searchInfinity, inCheck(p))
			if !ok || w.historyKey != before || !maps.Equal(counts, w.repetitions) {
				t.Fatal("history or completion changed")
			}
			return s
		}
		noHistory, draw := run(newTable(), nil), run(newTable(), history)
		if noHistory >= -500 || draw != 0 {
			t.Fatalf("fixture not a losing/drawing pair: %d/%d", noHistory, draw)
		}
		tt := newTable()
		for _, h := range [][]uint64{nil, history, nil, history} {
			tt.generation++
			want := noHistory
			if len(h) > 0 {
				want = draw
			}
			if got := run(tt, h); got != want {
				t.Fatalf("warm-history contamination: %d want%d", got, want)
			}
		}
	}
}

func TestChildHistoryResetAndUnwind(t *testing.T) {
	for _, c := range []struct {
		fen, move string
		reset     bool
	}{
		{"4k3/8/8/8/8/8/P7/3QK3 w - - 12 1", "a2a3", true},
		{"4k3/8/8/8/8/8/r7/R3K3 w - - 12 1", "a1a2", true},
		{"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1", "e5d6", true},
		{"4k3/8/8/8/8/8/8/4K2R w K - 12 1", "e1g1", false},
		{"4k3/8/8/8/8/8/8/4K2R w K - 12 1", "h1h2", false},
		{"4k3/8/8/8/8/8/8/3QK3 w - - 12 1", "d1d2", false},
	} {
		t.Run(c.move, func(t *testing.T) {
			p := evalFEN(t, c.fen)
			var child *chess.Position
			checked := false
			for _, m := range p.ValidMoves() {
				if m.String() == c.move {
					child = p.Update(&m)
					checked = m.HasTag(chess.Check)
				}
			}
			if child == nil || (child.HalfMoveClock() == 0) != c.reset {
				t.Fatal("invalid reset fixture")
			}
			w := searcher{table: newTable()}
			w.pushPosition(p)
			// A real ancestor, not a fabricated count of an unreachable child.
			w.pushKey(chess.StartingPosition().ZobristHash())
			key, counts := w.historyKey, maps.Clone(w.repetitions)
			expected := key
			if c.reset {
				expected = 0
			}
			observer := searcher{historyKey: expected, repetitions: maps.Clone(counts)}
			observer.pushPosition(child)
			expected = observer.historyKey
			for _, quiet := range []bool{false, true} {
				w.table.clear()
				w.now = nil
				w.deadline = time.Time{}
				var ok bool
				if quiet {
					_, ok = w.childQuiet(child, 2, -searchInfinity, searchInfinity, checked, 1)
				} else {
					_, ok = w.childSearch(child, 2, -searchInfinity, searchInfinity, checked, 1)
					if _, hit := w.table.probe(child, expected); !hit {
						t.Fatal("wrong child cache context")
					}
				}
				if !ok || w.historyKey != key || !maps.Equal(w.repetitions, counts) {
					t.Fatal("complete child leaked history")
				}
				// Expire after descent and the child push, not before entering the helper.
				calls := 0
				start := time.Unix(100, 0)
				w.deadline = start.Add(time.Second)
				w.now = func() time.Time {
					calls++
					if calls > 1 {
						return start.Add(time.Second)
					}
					return start
				}
				w.table.clear()
				if quiet {
					_, ok = w.childQuiet(child, 8, -searchInfinity, searchInfinity, checked, 1)
				} else {
					_, ok = w.childSearch(child, 8, -searchInfinity, searchInfinity, checked, 1)
				}
				if ok || w.historyKey != key || !maps.Equal(w.repetitions, counts) {
					t.Fatal("interrupted child leaked history")
				}
				if w.table.boardHint(child) != 0 {
					t.Fatal("incomplete child published hint")
				}
			}
		})
	}
}

func TestBoardHintsAreLegalOrderingOnly(t *testing.T) {
	p := chess.StartingPosition()
	tt := newTable()
	moves := p.ValidMoves()
	hint := encodeMove(moves[len(moves)-1])
	tt.storeHint(p, hint)
	if tt.boardHint(p) != hint {
		t.Fatal("missing board hint")
	}
	if _, hit := tt.probe(p, 123); hit {
		t.Fatal("hint became a score")
	}
	preferMove(moves, tt.boardHint(p))
	if encodeMove(moves[0]) != hint {
		t.Fatal("hint not ordered")
	}
	before := append([]chess.Move(nil), moves...)
	preferMove(moves, 65535)
	if !sameMoves(moves, before) {
		t.Fatal("illegal hint changed legal moves")
	}
	tt.hints = make([]boardHint, 1)
	tt.storeHint(p, hint)
	q := evalFEN(t, "4k3/8/8/8/8/8/P7/4K3 w - - 0 1")
	if tt.boardHint(q) != 0 {
		t.Fatal("hint slot collision accepted")
	}
	tt.clear()
	if tt.boardHint(p) != 0 {
		t.Fatal("new game retained hints")
	}
	empty := &transpositionTable{}
	empty.storeHint(p, hint)
	if empty.boardHint(p) != 0 {
		t.Fatal("empty hint table")
	}
	if bytes := uintptr(tableEntries) * unsafe.Sizeof(boardHint{}); bytes > 1<<20 {
		t.Fatalf("hint memory exceeds 1 MiB: %d", bytes)
	}
}

func TestTacticalOrderingPreservesEligibilityAndLibraryCache(t *testing.T) {
	for _, fen := range []string{
		"r3k2r/ppp2ppp/2n5/3pp3/3PP3/2N5/PPP2PPP/R3K2R w KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
		"8/k1P5/2K5/8/8/8/8/8 w - - 0 1",
		"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1",
	} {
		p := evalFEN(t, fen)
		before := append([]chess.Move(nil), p.ValidMovesUnsafe()...)
		want := orderedMoves(p)
		want = slices.DeleteFunc(want, func(m chess.Move) bool {
			return !m.HasTag(chess.Capture) && !m.HasTag(chess.EnPassant) && m.Promo() == chess.NoPieceType
		})
		got := orderedTacticalMoves(p)
		if !sameMoves(want, got) {
			t.Fatalf("tactical ordering/eligibility differs: %s", fen)
		}
		if !sameMoves(before, p.ValidMovesUnsafe()) {
			t.Fatal("library legal-move cache mutated")
		}
	}
}

func sameMoves(a, b []chess.Move) bool {
	return slices.EqualFunc(a, b, func(x, y chess.Move) bool {
		return x.String() == y.String() && x.HasTag(chess.Check) == y.HasTag(chess.Check) &&
			x.HasTag(chess.Capture) == y.HasTag(chess.Capture) && x.HasTag(chess.EnPassant) == y.HasTag(chess.EnPassant)
	})
}
