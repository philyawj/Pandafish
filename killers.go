package main

import "github.com/corentings/chess/v2"

// Two quiet cutoff hints per search ply. Each go command owns a fresh table;
// iterative-deepening passes share it. Hints never supply scores or new moves.
type killerTable [maxSearchDepth][2]uint16

func isQuietMove(move chess.Move) bool {
	return !move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant) && move.Promo() == chess.NoPieceType
}

func (table *killerTable) remember(ply int, move chess.Move) {
	if table == nil || ply < 0 || ply >= len(table) || !isQuietMove(move) {
		return
	}
	hint := encodeMove(move)
	if table[ply][0] != hint {
		table[ply] = [2]uint16{hint, table[ply][0]}
	}
}

// orderedMoves puts all captures and promotions first. Keep those priorities,
// then move matching legal quiet hints ahead of the remaining quiet moves.
func (table *killerTable) order(moves []chess.Move, ply int) {
	if table == nil || ply < 0 || ply >= len(table) {
		return
	}
	start := 0
	for start < len(moves) && !isQuietMove(moves[start]) {
		start++
	}
	for _, hint := range table[ply] {
		if hint == 0 {
			continue
		}
		for i := start; i < len(moves); i++ {
			if isQuietMove(moves[i]) && encodeMove(moves[i]) == hint {
				move := moves[i]
				copy(moves[start+1:i+1], moves[start:i])
				moves[start] = move
				start++
				break
			}
		}
	}
}
