package main

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestSearchMoveBackendMatchesUnchangedLibrary(t *testing.T) {
	fixtures := []string{
		chess.StartingPosition().String(),
		"r3k2r/ppp2ppp/2n5/3pp3/3PP3/2N5/PPP2PPP/R3K2R w KQkq - 0 1",
		"r3k2r/ppp2ppp/2n5/3pp3/3PP3/2N5/PPP2PPP/R3K2R b KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"k3r3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
		"8/k1P5/2K5/8/8/8/8/8 w - - 0 1",
		"7r/8/8/8/8/6k1/8/6BK w - - 0 1",
		"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1",
		"7k/6Q1/5K2/8/8/8/8/8 b - - 0 1",
	}
	check := func(p *chess.Position) {
		t.Helper()
		fen, hash := p.String(), p.ZobristHash()
		// Force a cold position even if game traversal already populated its cache.
		p = evalFEN(t, fen)
		got := p.AppendSearchMoves(nil, false)
		tactical := p.AppendSearchMoves(nil, true)
		status, exists := p.SearchStatus(), p.HasLegalMove()
		want := p.ValidMoves()
		if !sameMoves(got, want) || status != p.Status() || exists != (len(want) > 0) {
			t.Fatalf("full moves/status mismatch %s", fen)
		}
		filtered := slices.DeleteFunc(append([]chess.Move(nil), want...), func(m chess.Move) bool {
			return !m.HasTag(chess.Capture) && !m.HasTag(chess.EnPassant) && m.Promo() == chess.NoPieceType
		})
		if !sameMoves(tactical, filtered) {
			t.Fatalf("tactical mismatch %s", fen)
		}
		// Warm-cache append must retain a caller's prefix and preserve the cache.
		prefix := []chess.Move{{}}
		appended := p.AppendSearchMoves(prefix, true)
		if appended[0].String() != prefix[0].String() || !sameMoves(appended[1:], filtered) || !sameMoves(p.ValidMovesUnsafe(), want) {
			t.Fatal("append/cache ownership")
		}
		for i, m := range got {
			if m.HasTag(chess.KingSideCastle) != want[i].HasTag(chess.KingSideCastle) || m.HasTag(chess.QueenSideCastle) != want[i].HasTag(chess.QueenSideCastle) {
				t.Fatal("castle tags")
			}
			a, b := p.Update(&m), p.Update(&want[i])
			if a.String() != b.String() || a.ZobristHash() != b.ZobristHash() {
				t.Fatal("child state/hash differs")
			}
		}
		w := searcher{}
		if !sameMoves(w.searchMoves(p, 0, false), orderedMoves(p)) || !sameMoves(w.searchMoves(p, 0, true), orderedTacticalMoves(p)) {
			t.Fatal("stable sorted order differs")
		}
		if p.String() != fen || p.ZobristHash() != hash {
			t.Fatal("position mutated")
		}
	}
	for _, fen := range fixtures {
		check(evalFEN(t, fen))
	}
	random := rand.New(rand.NewPCG(17, 1601))
	p := chess.StartingPosition()
	for i := 0; i < 500; i++ {
		check(p)
		moves := p.ValidMoves()
		if len(moves) == 0 || i%80 == 79 {
			p = chess.StartingPosition()
			continue
		}
		m := moves[random.IntN(len(moves))]
		p = p.Update(&m)
	}
}

func TestSearchMoveBuffersRespectActivePly(t *testing.T) {
	p := chess.StartingPosition()
	w := searcher{}
	parent := w.searchMoves(p, 0, false)
	saved := append([]chess.Move(nil), parent...)
	m := parent[0]
	child := p.Update(&m)
	for i := 0; i < 4; i++ {
		w.searchMoves(child, 1, i%2 == 0)
		w.searchMoves(child, 30, false)
	}
	if !sameMoves(parent, saved) {
		t.Fatal("descendant/re-search overwrote parent moves")
	}
	// No captures is not stalemate, and a checked quiet block remains eligible.
	quiet := evalFEN(t, "7k/8/8/8/8/8/P7/7K w - - 0 1")
	if len(w.searchMoves(quiet, 1, true)) != 0 || quiet.SearchStatus() != chess.NoMethod {
		t.Fatal("empty tactics treated as stalemate")
	}
	checked := evalFEN(t, "7r/8/8/8/8/6k1/8/6BK w - - 0 1")
	moves := w.searchMoves(checked, 1, false)
	if len(moves) != 1 || moves[0].String() != "g1h2" {
		t.Fatal("quiet evasion omitted")
	}
}

func TestSearchLegalGenerationPerft(t *testing.T) {
	var count func(*chess.Position, int) uint64
	count = func(p *chess.Position, depth int) uint64 {
		if depth == 0 {
			return 1
		}
		var total uint64
		for _, m := range p.AppendSearchMoves(nil, false) {
			total += count(p.Update(&m), depth-1)
		}
		return total
	}
	if got := count(chess.StartingPosition(), 4); got != 197281 {
		t.Fatalf("start depth4 perft %d", got)
	}
}
