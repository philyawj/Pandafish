package main

import "github.com/corentings/chess/v2"

// Only a child reached by an irreversible pawn/capture move resets history.
// Direct entry points (including FEN and repeated-position tests) do not reset.
// Keep counts for draw detection; reset only the cache fingerprint. Earlier
// positions cannot recur after a pawn move/capture, so their counts cannot
// affect descendant repetition results. Restore the prior key on return.
func (w *searcher) childSearch(p *chess.Position, d, a, b int, c bool, ply int) (int, bool) {
	if p.HalfMoveClock() != 0 {
		return w.searchAtPly(p, d, a, b, c, ply)
	}
	key := w.historyKey
	w.historyKey = 0
	defer func() { w.historyKey = key }()
	return w.searchAtPly(p, d, a, b, c, ply)
}
func (w *searcher) childQuiet(p *chess.Position, d, a, b int, c bool, ply int) (int, bool) {
	if p.HalfMoveClock() != 0 {
		return w.quietAtPly(p, d, a, b, c, ply)
	}
	key := w.historyKey
	w.historyKey = 0
	defer func() { w.historyKey = key }()
	return w.quietAtPly(p, d, a, b, c, ply)
}

func (w *searcher) childContactQuiet(p *chess.Position, d, a, b int, c bool, ply, budget int, pending *chess.Move) (int, bool) {
	if p.HalfMoveClock() != 0 {
		return w.quietContactAtPly(p, d, a, b, c, ply, budget, pending)
	}
	key := w.historyKey
	w.historyKey = 0
	defer func() { w.historyKey = key }()
	return w.quietContactAtPly(p, d, a, b, c, ply, budget, pending)
}
