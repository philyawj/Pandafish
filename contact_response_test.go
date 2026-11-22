package main

import (
	"context"
	"maps"
	"testing"

	"github.com/corentings/chess/v2"
)

func mirrorContactPosition(t *testing.T, p *chess.Position) *chess.Position {
	t.Helper()
	pieces := map[chess.Square]chess.Piece{}
	for sq, pc := range p.Board().SquareMap() {
		pieces[sq^56] = chess.NewPiece(pc.Type(), pc.Color().Other())
	}
	b, err := chess.NewBoard(pieces)
	if err != nil {
		t.Fatal(err)
	}
	return evalFEN(t, b.String()+" "+p.Turn().Other().String()+" - - 0 1")
}

// Full enumeration of actual defenses, public immutable transitions and the
// declared mate reply, with ordinary budget-zero q only for unresolved leaves.
func contactResponseReference(t *testing.T, p *chess.Position, depth int) int {
	t.Helper()
	var pass, reply chess.SearchPosition
	threat, ok := passedContactThreat(p, &pass, &reply)
	if !ok || p.Status() != chess.NoMethod {
		t.Fatal("reference requires live threat")
	}
	best := -searchInfinity
	for _, defense := range p.ValidMoves() {
		child := p.Update(&defense)
		score := 0
		switch child.Status() {
		case chess.Checkmate:
			score = -mateScore + 1
		case chess.Stalemate:
			score = 0
		default:
			mate := false
			for _, m := range child.ValidMoves() {
				if m.String() == threat.String() && m.HasTag(chess.Check) && child.Update(&m).Status() == chess.Checkmate {
					mate = true
					break
				}
			}
			if mate {
				score = mateScore - 2
			} else {
				var complete bool
				score, complete = (&searcher{}).quietContactAtPly(child, depth-1, -searchInfinity, searchInfinity, inCheck(child), 1, 0, nil)
				if !complete {
					t.Fatal("reference interrupted")
				}
			}
		}
		best = max(best, -score)
	}
	return best
}

func TestContactResponseSoundnessAndWindows(t *testing.T) {
	for _, c := range []struct {
		fen    string
		forced bool
	}{
		{"7k/p7/5KQ1/8/8/8/8/8 b - - 0 1", true},
		{"7k/5p2/5KQ1/8/8/8/8/8 b - - 0 1", false},
		{"7k/8/5KQ1/r7/8/8/8/8 b - - 0 1", false},
	} {
		base := evalFEN(t, c.fen)
		for _, p := range []*chess.Position{base, mirrorContactPosition(t, base)} {
			want := contactResponseReference(t, p, 4)
			if c.forced && want != -mateScore+2 || !c.forced && want <= -mateBound {
				t.Fatal("fixture/certificate", c.fen, want)
			}
			for _, window := range [][2]int{{-searchInfinity, searchInfinity}, {-1000, -999}, {100, 101}, {-mateScore, -mateScore + 1}, {-mateScore + 2, -mateScore + 3}, {-10, 10}} {
				w := searcher{}
				got, ok := w.quiet(p, 4, window[0], window[1], inCheck(p))
				if !ok || got <= window[0] && want > window[0] || got >= window[1] && want < window[1] || got > window[0] && got < window[1] && got != want {
					t.Fatal("window/response mismatch", p, window, got, want, ok)
				}
			}
		}
	}
}

func TestContactResponseStalemateBudgetAndHistory(t *testing.T) {
	p := evalFEN(t, "7k/8/5KQ1/8/8/8/8/8 b - - 0 1")
	if p.Status() != chess.Stalemate || quiescence(p, 4, -searchInfinity, searchInfinity) != 0 {
		t.Fatal("pass threat overrides stalemate")
	}
	p = evalFEN(t, "7k/p7/5KQ1/8/8/8/8/8 b - - 0 1")
	// Both actual reply plies must fit; depth one does not activate the event.
	for _, d := range []int{0, 1} {
		got, ok := (&searcher{}).quietContactAtPly(p, d, -searchInfinity, searchInfinity, false, 0, 1, nil)
		base, _ := (&searcher{}).quietContactAtPly(p, d, -searchInfinity, searchInfinity, false, 0, 0, nil)
		if !ok || got != base || got <= -mateBound {
			t.Fatal("budget/cap", d, got, base)
		}
	}
	w := searcher{}
	w.pushPosition(p)
	w.pushPosition(p)
	counts, key := maps.Clone(w.repetitions), w.historyKey
	if got, ok := w.quiet(p, 4, -searchInfinity, searchInfinity, false); !ok || got != 0 || !maps.Equal(counts, w.repetitions) || key != w.historyKey {
		t.Fatal("repetition/history", got, ok)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, ok := (&searcher{ctx: ctx}).quiet(p, 4, -searchInfinity, searchInfinity, false); ok {
		t.Fatal("cancelled threat search")
	}
}
