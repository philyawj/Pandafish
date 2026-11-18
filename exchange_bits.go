package main

import (
	"math/bits"

	"github.com/corentings/chess/v2"
)

var exchangeGeometry = func() (g struct {
	rays           [64][8]uint64
	knights, kings [64]uint64
	pawns          [3][64]uint64
}) {
	for target := 0; target < 64; target++ {
		x, y := target%8, target/8
		for i, d := range [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
			for f, r := x+d[0], y+d[1]; f >= 0 && f < 8 && r >= 0 && r < 8; f, r = f+d[0], r+d[1] {
				g.rays[target][i] |= uint64(1) << uint(63-(r*8+f))
			}
		}
		for from := 0; from < 64; from++ {
			dx, dy := from%8-x, from/8-y
			bit := uint64(1) << uint(63-from)
			if abs(dx)*abs(dy) == 2 {
				g.knights[target] |= bit
			}
			if from != target && abs(dx) <= 1 && abs(dy) <= 1 {
				g.kings[target] |= bit
			}
			if abs(dx) == 1 && dy == -1 {
				g.pawns[chess.White][target] |= bit
			}
			if abs(dx) == 1 && dy == 1 {
				g.pawns[chess.Black][target] |= bit
			}
		}
	}
	return g
}()

func newExchangeBoard(position *chess.Position) (board exchangeBoard) {
	sets := position.Board().EvaluationPieces()
	pieces := [12]chess.Piece{chess.WhiteKing, chess.WhiteQueen, chess.WhiteRook, chess.WhiteBishop, chess.WhiteKnight, chess.WhitePawn, chess.BlackKing, chess.BlackQueen, chess.BlackRook, chess.BlackBishop, chess.BlackKnight, chess.BlackPawn}
	kinds := [6]chess.PieceType{chess.King, chess.Queen, chess.Rook, chess.Bishop, chess.Knight, chess.Pawn}
	for i, set := range sets {
		side, kind := chess.Color(i/6+1), kinds[i%6]
		board.sides[side] |= set
		board.types[kind] |= set
		for set != 0 {
			square := chess.Square(bits.LeadingZeros64(set))
			set &^= uint64(1) << uint(63-square)
			board.pieces[square] = pieces[i]
			if kind == chess.King {
				board.kings[side] = square
			}
		}
	}
	return board
}

func (board *exchangeBoard) put(square chess.Square, piece chess.Piece) {
	bit := uint64(1) << uint(63-square)
	if old := board.pieces[square]; old != chess.NoPiece {
		board.sides[old.Color()] &^= bit
		board.types[old.Type()] &^= bit
	}
	board.pieces[square] = piece
	if piece != chess.NoPiece {
		board.sides[piece.Color()] |= bit
		board.types[piece.Type()] |= bit
	}
}

func (board *exchangeBoard) attackers(target chess.Square, side chess.Color) uint64 {
	if target < chess.A1 || target > chess.H8 || side < chess.White || side > chess.Black {
		return 0
	}
	result := (exchangeGeometry.pawns[side][target] & board.types[chess.Pawn]) | (exchangeGeometry.knights[target] & board.types[chess.Knight]) | (exchangeGeometry.kings[target] & board.types[chess.King])
	occupied := board.sides[chess.White] | board.sides[chess.Black]
	for i, ray := range exchangeGeometry.rays[target] {
		hits := ray & occupied
		if hits == 0 {
			continue
		}
		first := hits & -hits
		if i%2 == 0 {
			first = uint64(1) << uint(63-bits.LeadingZeros64(hits))
		}
		sliders := board.types[chess.Bishop] | board.types[chess.Queen]
		if i < 4 {
			sliders = board.types[chess.Rook] | board.types[chess.Queen]
		}
		result |= first & sliders
	}
	return result & board.sides[side]
}

func (board *exchangeBoard) leastAttacker(target chess.Square, side chess.Color) int {
	best, value := -1, mateScore+1
	candidates := board.attackers(target, side)
	for candidates != 0 {
		from := chess.Square(bits.LeadingZeros64(candidates))
		candidates &^= uint64(1) << uint(63-from)
		piece := board.pieces[from]
		cost := pieceValues[piece.Type()]
		if piece.Type() == chess.King {
			cost = mateScore
		}
		if cost >= value {
			continue
		}
		captured, king := board.pieces[target], board.kings[side]
		board.put(from, chess.NoPiece)
		board.put(target, piece)
		if piece.Type() == chess.King {
			board.kings[side] = target
		}
		legal := !board.attacked(board.kings[side], side.Other())
		board.put(from, piece)
		board.put(target, captured)
		board.kings[side] = king
		if legal {
			best, value = int(from), cost
		}
	}
	return best
}

func (board *exchangeBoard) attacked(target chess.Square, side chess.Color) bool {
	return board.attackers(target, side) != 0
}
