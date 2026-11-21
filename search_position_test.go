package main

import (
	"context"
	"math/rand/v2"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

func assertSearchPosition(t *testing.T, got, want *chess.Position) {
	t.Helper()
	if got.String() != want.String() || got.ZobristHash() != want.ZobristHash() || got.Board().EvaluationPieces() != want.Board().EvaluationPieces() {
		t.Fatal("scratch child state/hash/pieces", got, want)
	}
	for sq := chess.A1; sq <= chess.H8; sq++ {
		if got.Board().Piece(sq) != want.Board().Piece(sq) {
			t.Fatal("scratch mailbox", sq)
		}
	}
	if got.SearchStatus() != want.Status() || !sameMoves(got.AppendSearchMoves(nil, false), want.ValidMoves()) || inCheck(got) != inCheck(want) {
		t.Fatal("scratch legal state/cache", got, want)
	}
}

func TestSearchPositionReuseMatchesImmutable(t *testing.T) {
	cell := new(chess.SearchPosition)
	checked := 0
	check := func(p *chess.Position) {
		t.Helper()
		fen, hash := p.String(), p.ZobristHash()
		for _, move := range p.AppendSearchMoves(nil, false) {
			got := cell.Update(p, &move)
			assertSearchPosition(t, got, p.Update(&move))
			assertSearchPosition(t, got, p.UpdateSearch(&move))
			// Populate the scratch cache; the next sibling must discard it.
			got.ValidMoves()
			checked++
		}
		assertSearchPosition(t, cell.Update(p, nil), p.Update(nil))
		if p.String() != fen || p.ZobristHash() != hash {
			t.Fatal("scratch mutated parent")
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
		"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1",
	} {
		check(evalFEN(t, fen))
	}
	rng := rand.New(rand.NewPCG(22, 2201))
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
		t.Fatal("insufficient scratch transitions", checked)
	}
	t.Log("Scratch transitions matched:", checked)
}

func TestSearchPositionIsolationAndAliasRejection(t *testing.T) {
	p := chess.StartingPosition()
	worker := searcher{}
	m := p.ValidMoves()[0]
	child := worker.childPosition(p, &m, 1)
	fen, hash := child.String(), child.ZobristHash()
	moves := child.ValidMoves()
	for _, next := range moves {
		worker.childPosition(child, &next, 2)
		worker.childPosition(child, &next, 90)
	}
	other := searcher{}
	other.childPosition(p, &m, 1)
	if child.String() != fen || child.ZobristHash() != hash || !sameMoves(child.ValidMoves(), moves) {
		t.Fatal("descendant/growth/worker overwrote active parent")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("source/destination alias accepted")
			}
		}()
		worker.positionBuffers[1].Update(child, &moves[0])
	}()
	if child.String() != fen || child.ZobristHash() != hash {
		t.Fatal("alias panic occurred after mutation")
	}
	// The same root may be read concurrently by private workers after warming
	// its general-purpose legal cache. Scratch storage is never shared.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := searcher{}
			for j := 0; j < 100; j++ {
				q := w.childPosition(p, &m, 1)
				if q.String() != child.String() {
					t.Error("independent worker state")
				}
			}
		}()
	}
	wg.Wait()
}

func TestSearchPositionPerft(t *testing.T) {
	w := searcher{}
	var count func(*chess.Position, int, int) uint64
	count = func(p *chess.Position, depth, ply int) uint64 {
		if depth == 0 {
			return 1
		}
		var n uint64
		for _, m := range p.AppendSearchMoves(nil, false) {
			n += count(w.childPosition(p, &m, ply+1), depth-1, ply+1)
		}
		return n
	}
	if n := count(chess.StartingPosition(), 4, 0); n != 197281 {
		t.Fatal("scratch perft", n)
	}
}

func TestSearchPositionCancellationCleanupAndReuse(t *testing.T) {
	g, err := parsePosition([]string{"startpos", "moves", "g1f3", "g8f6", "f3g1", "f6g8"})
	if err != nil {
		t.Fatal(err)
	}
	p := g.Position()
	fen, hash := p.String(), p.ZobristHash()
	w := searcher{table: newTable(), killers: new(killerTable), ctx: context.Background()}
	for _, k := range repetitionHistory(g) {
		w.pushKey(k)
	}
	saved, key := map[uint64]int{}, w.historyKey
	for k, n := range w.repetitions {
		saved[k] = n
	}
	base := time.Now()
	calls := 0
	w.deadline = base.Add(time.Second)
	w.now = func() time.Time {
		calls++
		if calls > 500 {
			return base.Add(2 * time.Second)
		}
		return base
	}
	if _, _, ok := w.root(p, rootMoves(p, nil), 4); ok {
		t.Fatal("expected interrupted recursive work")
	}
	if !reflect.DeepEqual(w.repetitions, saved) || w.historyKey != key || p.String() != fen || p.ZobristHash() != hash {
		t.Fatal("interruption leaked parent/history state")
	}
	w.now = nil
	w.deadline = time.Time{}
	_, score, ok := w.root(p, rootMoves(p, nil), 3)
	if !ok {
		t.Fatal("reuse incomplete")
	}
	fresh := searcher{table: newTable(), killers: new(killerTable)}
	for _, k := range repetitionHistory(g) {
		fresh.pushKey(k)
	}
	_, want, ok := fresh.root(p, rootMoves(p, nil), 3)
	if !ok || score != want || !reflect.DeepEqual(w.repetitions, saved) || w.historyKey != key {
		t.Fatal("reused worker changed score/history", score, want)
	}
}
