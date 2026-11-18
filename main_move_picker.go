package main

import "github.com/corentings/chess/v2"

// The eager order is table hint, stable tactics, quiet killers, then remaining
// quiet generation order. Each stage materializes only after earlier stages fail.
type mainMovePicker struct {
	worker                         *searcher
	position                       *chess.Position
	ply, stage, index, killerIndex int
	hint                           uint16
	previous                       [3]uint16
	killers                        [2]uint16
	moves                          []chess.Move
}

func newMainMovePicker(worker *searcher, p *chess.Position, ply int, hint uint16) mainMovePicker {
	picker := mainMovePicker{worker: worker, position: p, ply: ply, hint: hint}
	if worker.killers != nil && ply >= 0 && ply < len(worker.killers) {
		picker.killers = worker.killers[ply]
	}
	// The pinned library trusts inconsistent FEN castling rights and can emit
	// duplicate signatures when a different piece occupies the king's origin.
	// Keep its original complete-list behavior for that compatibility corner.
	if eagerCastleCompatibility(p) {
		picker.moves = worker.searchMoves(p, ply, false)
		worker.killers.order(picker.moves, ply)
		preferMove(picker.moves, hint)
		picker.stage = 5
	}
	return picker
}

func eagerCastleCompatibility(p *chess.Position) bool {
	white := p.Turn() == chess.White
	if white && p.Board().Piece(chess.E1) == chess.WhiteKing {
		return false
	}
	if !white && p.Board().Piece(chess.E8) == chess.BlackKing {
		return false
	}
	rights := p.CastleRights()
	for i := 0; i < len(rights); i++ {
		if (white && (rights[i] == 'K' || rights[i] == 'Q')) || (!white && (rights[i] == 'k' || rights[i] == 'q')) {
			return true
		}
	}
	return false
}

func lookupHint(p *chess.Position, hint uint16) (chess.Move, bool) {
	if hint == 0 {
		return chess.Move{}, false
	}
	value := hint - 1
	return p.LookupSearchMove(chess.Square(value&63), chess.Square((value>>6)&63), chess.PieceType(value>>12))
}

func (picker *mainMovePicker) used(signature uint16) bool {
	for _, previous := range picker.previous {
		if signature == previous {
			return true
		}
	}
	return false
}

func (picker *mainMovePicker) next() (chess.Move, bool) {
	for {
		switch picker.stage {
		case 0:
			picker.stage = 1
			if move, ok := lookupHint(picker.position, picker.hint); ok {
				picker.previous[0] = picker.hint
				return move, true
			}
		case 1:
			if picker.moves == nil {
				picker.moves = picker.worker.searchMoves(picker.position, picker.ply, true)
			}
			for picker.index < len(picker.moves) {
				move := picker.moves[picker.index]
				picker.index++
				if !picker.used(encodeMove(move)) {
					return move, true
				}
			}
			picker.stage = 2
		case 2:
			for picker.killerIndex < len(picker.killers) {
				i := picker.killerIndex
				picker.killerIndex++
				hint := picker.killers[i]
				if hint == 0 || picker.used(hint) {
					continue
				}
				if move, ok := lookupHint(picker.position, hint); ok && isQuietMove(move) {
					picker.previous[i+1] = hint
					return move, true
				}
			}
			picker.stage = 3
			picker.moves = nil
			picker.index = 0
		case 3:
			if picker.moves == nil {
				picker.moves = picker.worker.quietSearchMoves(picker.position, picker.ply)
			}
			for picker.index < len(picker.moves) {
				move := picker.moves[picker.index]
				picker.index++
				if !picker.used(encodeMove(move)) {
					return move, true
				}
			}
			picker.stage = 4
		case 5:
			if picker.index < len(picker.moves) {
				move := picker.moves[picker.index]
				picker.index++
				return move, true
			}
			picker.stage = 4
		default:
			return chess.Move{}, false
		}
	}
}

func (worker *searcher) quietSearchMoves(p *chess.Position, ply int) []chess.Move {
	if len(worker.moveBuffers) <= ply {
		worker.moveBuffers = append(worker.moveBuffers, make([][]chess.Move, ply+1-len(worker.moveBuffers))...)
	}
	moves := p.AppendSearchQuietMoves(worker.moveBuffers[ply][:0])
	worker.moveBuffers[ply] = moves
	return moves
}
