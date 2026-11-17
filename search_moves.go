package main

import (
	"slices"

	"github.com/corentings/chess/v2"
)

// Each active ply owns its move slice. Descendants and PVS re-searches can
// reuse their own storage without overwriting any ancestor's iteration list.
func (w *searcher) searchMoves(p *chess.Position, ply int, tacticalOnly bool) []chess.Move {
	if len(w.moveBuffers) <= ply {
		w.moveBuffers = append(w.moveBuffers, make([][]chess.Move, ply+1-len(w.moveBuffers))...)
	}
	moves := p.AppendSearchMoves(w.moveBuffers[ply][:0], tacticalOnly)
	w.moveBuffers[ply] = moves
	priority := func(m chess.Move) int {
		value := pieceValues[m.Promo()]
		if m.HasTag(chess.Capture) || m.HasTag(chess.EnPassant) {
			victim := pieceValues[p.Board().Piece(m.S2()).Type()]
			if m.HasTag(chess.EnPassant) {
				victim = pieceValues[chess.Pawn]
			}
			value += 10000 + 10*victim - pieceValues[p.Board().Piece(m.S1()).Type()]
		}
		return value
	}
	slices.SortStableFunc(moves, func(a, b chess.Move) int { return priority(b) - priority(a) })
	return moves
}
