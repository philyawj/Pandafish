package main

import (
	"slices"

	"github.com/corentings/chess/v2"
)

type moveRank struct {
	index    int
	priority int
}

// Each active ply owns its move slice. Descendants and PVS re-searches can
// reuse their own storage without overwriting any ancestor's iteration list.
func (w *searcher) searchMoves(p *chess.Position, ply int, tacticalOnly bool) []chess.Move {
	if len(w.moveBuffers) <= ply {
		w.moveBuffers = append(w.moveBuffers, make([][]chess.Move, ply+1-len(w.moveBuffers))...)
	}
	if len(w.rankBuffers) <= ply {
		w.rankBuffers = append(w.rankBuffers, make([][]moveRank, ply+1-len(w.rankBuffers))...)
	}
	moves := p.AppendSearchMoves(w.moveBuffers[ply][:0], tacticalOnly)
	w.moveBuffers[ply] = moves
	ranks := w.rankBuffers[ply][:0]
	ordered := true
	for i := range moves {
		m := &moves[i]
		value := pieceValues[m.Promo()]
		if m.HasTag(chess.Capture) || m.HasTag(chess.EnPassant) {
			victim := pieceValues[p.Board().Piece(m.S2()).Type()]
			if m.HasTag(chess.EnPassant) {
				victim = pieceValues[chess.Pawn]
			}
			value += 10000 + 10*victim - pieceValues[p.Board().Piece(m.S1()).Type()]
		}
		if i > 0 && value > ranks[i-1].priority {
			ordered = false
		}
		ranks = append(ranks, moveRank{index: i, priority: value})
	}
	w.rankBuffers[ply] = ranks
	if !ordered {
		slices.SortStableFunc(ranks, func(a, b moveRank) int { return b.priority - a.priority })
		applyMoveOrder(moves, ranks)
	}
	return moves
}

// Each rank names the original source for its destination. Disjoint cycles
// need one saved move each; fixed indices mark completed cycles in scratch data.
func applyMoveOrder(moves []chess.Move, ranks []moveRank) {
	for i := range moves {
		if ranks[i].index == i {
			continue
		}
		saved := moves[i]
		for dst := i; ; {
			src := ranks[dst].index
			ranks[dst].index = dst
			if src == i {
				moves[dst] = saved
				break
			}
			moves[dst] = moves[src]
			dst = src
		}
	}
}
