package main

import "github.com/corentings/chess/v2"

// Static exchange evaluation estimates material gained by a legal move and
// least-valuable legal recaptures on its destination. Either side may stop.
// It is a capture-ordering/filter heuristic, not a tactical search score.
func staticExchange(position *chess.Position, move chess.Move) int {
	score, _ := (&searcher{}).exchange(position, move)
	return score
}

func captureValue(position *chess.Position, move chess.Move) int {
	if move.HasTag(chess.EnPassant) {
		return pieceValues[chess.Pawn]
	}
	return pieceValues[position.Board().Piece(move.S2()).Type()]
}

type exchangeBoard struct {
	pieces [64]chess.Piece
	kings  [3]chess.Square // Colors are NoColor, White, Black.
	sides  [3]uint64
	types  [7]uint64
}

func (worker *searcher) exchange(position *chess.Position, move chess.Move) (int, bool) {
	if worker.stopped() {
		return 0, false
	}
	board := newExchangeBoard(position)
	from, target := move.S1(), move.S2()
	piece := board.pieces[from]
	gain := [64]int{captureValue(position, move)}
	if move.Promo() != chess.NoPieceType {
		gain[0] += pieceValues[move.Promo()] - pieceValues[chess.Pawn]
		piece = chess.NewPiece(move.Promo(), piece.Color())
	}
	board.put(from, chess.NoPiece)
	board.put(target, piece)
	if move.HasTag(chess.EnPassant) {
		captured := target - 8
		if piece.Color() == chess.Black {
			captured = target + 8
		}
		board.put(captured, chess.NoPiece)
	}
	if piece.Type() == chess.King {
		board.kings[piece.Color()] = target
	}
	side, depth := piece.Color().Other(), 0
	for depth+1 < len(gain) {
		if worker.stopped() {
			return 0, false
		}
		attacker := board.leastAttacker(target, side)
		if attacker < 0 {
			break
		}
		depth++
		incoming := board.pieces[attacker]
		promotion := 0
		if incoming.Type() == chess.Pawn && (target.Rank() == chess.Rank1 || target.Rank() == chess.Rank8) {
			promotion = pieceValues[chess.Queen] - pieceValues[chess.Pawn]
			incoming = chess.NewPiece(chess.Queen, side)
		}
		gain[depth] = pieceValues[board.pieces[target].Type()] + promotion - gain[depth-1]
		board.put(chess.Square(attacker), chess.NoPiece)
		board.put(target, incoming)
		if incoming.Type() == chess.King {
			board.kings[side] = target
		}
		side = side.Other()
	}
	for depth > 0 {
		gain[depth-1] = -max(-gain[depth-1], gain[depth])
		depth--
	}
	return gain[0], !worker.stopped()
}
