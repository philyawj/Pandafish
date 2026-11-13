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
}

func (worker *searcher) exchange(position *chess.Position, move chess.Move) (int, bool) {
	if worker.stopped() {
		return 0, false
	}
	board := exchangeBoard{}
	for square := chess.A1; square <= chess.H8; square++ {
		piece := position.Board().Piece(square)
		board.pieces[square] = piece
		if piece.Type() == chess.King {
			board.kings[piece.Color()] = square
		}
	}
	from, target := move.S1(), move.S2()
	piece := board.pieces[from]
	gain := [64]int{captureValue(position, move)}
	if move.Promo() != chess.NoPieceType {
		gain[0] += pieceValues[move.Promo()] - pieceValues[chess.Pawn]
		piece = chess.NewPiece(move.Promo(), piece.Color())
	}
	board.pieces[from], board.pieces[target] = chess.NoPiece, piece
	if move.HasTag(chess.EnPassant) {
		captured := target - 8
		if piece.Color() == chess.Black {
			captured = target + 8
		}
		board.pieces[captured] = chess.NoPiece
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
		board.pieces[attacker], board.pieces[target] = chess.NoPiece, incoming
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

func (board *exchangeBoard) leastAttacker(target chess.Square, side chess.Color) int {
	best, value := -1, mateScore+1
	for from, piece := range board.pieces {
		if piece.Color() != side || !board.attacks(chess.Square(from), target) {
			continue
		}
		cost := pieceValues[piece.Type()]
		if piece.Type() == chess.King {
			cost = mateScore // Kings are considered after other attackers.
		}
		if cost >= value {
			continue
		}
		// Remove the origin before testing king safety, exposing pins/x-rays.
		captured, king := board.pieces[target], board.kings[side]
		board.pieces[from], board.pieces[target] = chess.NoPiece, piece
		if piece.Type() == chess.King {
			board.kings[side] = target
		}
		legal := !board.attacked(board.kings[side], side.Other())
		board.pieces[from], board.pieces[target], board.kings[side] = piece, captured, king
		if legal {
			best, value = from, cost
		}
	}
	return best
}

func (board *exchangeBoard) attacked(target chess.Square, side chess.Color) bool {
	for from, piece := range board.pieces {
		if piece.Color() == side && board.attacks(chess.Square(from), target) {
			return true
		}
	}
	return false
}

func (board *exchangeBoard) attacks(from, to chess.Square) bool {
	if from == to {
		return false
	}
	piece := board.pieces[from]
	x, y := int(from)%8, int(from)/8
	dx, dy := int(to)%8-x, int(to)/8-y
	switch piece.Type() {
	case chess.Pawn:
		direction := 1
		if piece.Color() == chess.Black {
			direction = -1
		}
		return abs(dx) == 1 && dy == direction
	case chess.Knight:
		return abs(dx)*abs(dy) == 2
	case chess.King:
		return abs(dx) <= 1 && abs(dy) <= 1
	case chess.Bishop:
		if abs(dx) != abs(dy) {
			return false
		}
	case chess.Rook:
		if dx != 0 && dy != 0 {
			return false
		}
	case chess.Queen:
		if dx != 0 && dy != 0 && abs(dx) != abs(dy) {
			return false
		}
	default:
		return false
	}
	stepX, stepY := 0, 0
	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}
	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}
	for x, y = x+stepX, y+stepY; x != int(to)%8 || y != int(to)/8; x, y = x+stepX, y+stepY {
		if board.pieces[y*8+x] != chess.NoPiece {
			return false
		}
	}
	return true
}
