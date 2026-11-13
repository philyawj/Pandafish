package main

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

func TestImmediateMateIsPreferredForBothColors(t *testing.T) {
	for _, c := range []struct{ fen, move string }{
		{"7k/5Q2/5K2/8/8/8/8/8 w - - 0 1", "f7g7"},
		{"7k/8/5KQ1/8/8/8/8/8 w - - 0 1", "g6g7"},
		{"8/8/8/8/8/5k2/5q2/7K b - - 0 1", "f2g2"},
		{"8/8/8/8/8/5kq1/8/7K b - - 0 1", "g3g2"},
	} {
		p := evalFEN(t, c.fen)
		before := p.String()
		w := searcher{table: newTable()}
		for depth := 1; depth <= 4; depth++ {
			for repeat := 0; repeat < 2; repeat++ {
				best, score, ok := w.root(p, rootMoves(p, nil), depth)
				if !ok || score != mateScore-1 || len(best) != 1 || best[0] != c.move {
					t.Fatalf("%s depth%d best%v score%d complete%v", c.fen, depth, best, score, ok)
				}
			}
		}
		result := iterativeSearch(context.Background(), p, goLimits{depth: 3}, timeBudget{}, time.Now())
		if result.move != c.move || result.metrics.ScoreKind != "mate" || result.metrics.Score != 1 {
			t.Fatalf("wrong mate result %+v", result)
		}
		if p.String() != before {
			t.Fatal("root mutated")
		}
	}
}

func TestMateScoresAndTableCoordinates(t *testing.T) {
	for _, score := range []int{-mateScore + 10, mateScore - 10, -1200, 0, 1200} {
		for ply := 0; ply < 85; ply++ {
			if got := scoreFromTable(scoreToTable(score, ply), ply); got != score {
				t.Fatalf("round trip %d at%d ->%d", score, ply, got)
			}
		}
	}
	p := evalFEN(t, "7k/5Q2/5K2/8/8/8/8/8 w - - 0 1")
	w := searcher{table: newTable()}
	for _, ply := range []int{0, 5, 2, 17} {
		score, ok := w.searchAtPly(p, 3, -searchInfinity, searchInfinity, false, ply)
		if !ok || score != mateScore-ply-1 {
			t.Fatalf("cross-ply win%d score%d", ply, score)
		}
	}
	// Negative mate bounds use the opposite adjustment. Keep the bound kind
	// unchanged when normalizing a score at one root and probing at another.
	table := newTable()
	table.store(p, 3, scoreToTable(-mateScore+7, 4), scoreToTable(-mateScore+8, 4), scoreToTable(-mateScore+9, 4), 0)
	entry, ok := table.probe(p)
	if !ok || entry.bound != upperBound || scoreFromTable(int(entry.score), 9) != -mateScore+12 {
		t.Fatal("negative mate normalization/bound changed")
	}
	mated := evalFEN(t, "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1")
	if score, ok := (&searcher{}).quietAtPly(mated, quiescenceDepth, -searchInfinity, searchInfinity, true, 7); !ok || score != -mateScore+7 {
		t.Fatal("quiescence lost mate ply", score, ok)
	}
	for _, c := range []struct{ score, moves int }{{mateScore - 1, 1}, {mateScore - 3, 2}, {-mateScore + 2, -1}, {-mateScore + 4, -2}} {
		if got := mateMoves(c.score); got != c.moves {
			t.Fatalf("UCI mate count %d want%d", got, c.moves)
		}
	}
}

// Mate-only exhaustive reference: no alpha-beta, quiescence, static score or
// production mate helpers. Unknown at the horizon is neither a win nor a loss.
type mateProof struct{ outcome, distance int }
type mateKey struct {
	hash  uint64
	depth int
}

func proveMate(p *chess.Position, depth int, cache map[mateKey]mateProof) mateProof {
	if p.Status() == chess.Checkmate {
		return mateProof{-1, 0}
	}
	if p.Status() != chess.NoMethod || depth == 0 {
		return mateProof{}
	}
	key := mateKey{p.ZobristHash(), depth}
	if r, ok := cache[key]; ok {
		return r
	}
	shortest, longest, allLost := 10000, 0, true
	for _, m := range p.ValidMoves() {
		r := proveMate(p.Update(&m), depth-1, cache)
		if r.outcome == -1 {
			shortest = min(shortest, r.distance+1)
		}
		if r.outcome != 1 {
			allLost = false
		}
		longest = max(longest, r.distance+1)
	}
	r := mateProof{}
	if shortest < 10000 {
		r = mateProof{1, shortest}
	} else if allLost {
		r = mateProof{-1, longest}
	}
	cache[key] = r
	return r
}

func TestDelayForcedMateMatchesIndependentProof(t *testing.T) {
	p := evalFEN(t, "7k/8/5Q2/5K2/8/8/8/8 b - - 0 1")
	cache := map[mateKey]mateProof{}
	proof := proveMate(p, 7, cache)
	if proof.outcome != -1 {
		t.Fatalf("fixture not proven lost: %+v", proof)
	}
	var want []string
	for _, m := range p.ValidMoves() {
		r := proveMate(p.Update(&m), 6, cache)
		if r.outcome == 1 && r.distance+1 == proof.distance {
			want = append(want, m.String())
		}
	}
	best, score, ok := (&searcher{table: newTable()}).root(p, rootMoves(p, nil), 7)
	slices.Sort(best)
	slices.Sort(want)
	if !ok || !slices.Equal(best, want) || score != -mateScore+proof.distance {
		t.Fatalf("delay best%v want%v score%d proof%+v", best, want, score, proof)
	}
	t.Logf("proven losing distance%d; longest resistance%v", proof.distance, want)
}

func TestSearchedMateUCIOutput(t *testing.T) {
	p := evalFEN(t, "7k/5Q2/5K2/8/8/8/8/8 w - - 0 1")
	r := iterativeSearch(context.Background(), p, goLimits{depth: 3}, timeBudget{}, time.Now())
	var out strings.Builder
	reportSearch(&out, r)
	if !strings.Contains(out.String(), " score mate 1") || strings.Contains(out.String(), "score cp") {
		t.Fatal(out.String())
	}
}
