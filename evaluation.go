package main

import "github.com/corentings/chess/v2"

// All evaluation/search scores are centipawns. Mate exceeds any legal material
// balance plus positional terms, including positions with promoted pieces.
const mateScore = 1000000

var pieceValues = map[chess.PieceType]int{
	chess.Pawn: 100, chess.Knight: 300, chess.Bishop: 300, chess.Rook: 500, chess.Queen: 900,
}

func evaluate(position *chess.Position, perspective chess.Color) int {
	switch position.Status() {
	case chess.Checkmate:
		if position.Turn() == perspective {
			return -mateScore
		}
		return mateScore
	case chess.Stalemate:
		return 0
	}
	return staticScore(position, perspective)
}

// Original, small placement formulas rather than copied/tuned engine tables.
// Relative rank is zero on each side's home rank. The king's middlegame and
// endgame preferences blend using remaining non-pawn material (24 at start).
func placement(piece chess.PieceType, file, rank, phase int) int {
	center := 6 - (abs(2*file-7)+abs(2*rank-7))/2 // -1 on corners; 5 on central squares.
	switch piece {
	case chess.Pawn:
		// Center pawns get space; flank pawns are not rewarded for rushing forward.
		central := max(0, 3-abs(2*file-7))
		return central*min(rank, 3)*4 + max(0, rank-3)*6
	case chess.Knight:
		return center * 8
	case chess.Bishop:
		return center*4 + min(rank, 3)*3
	case chess.Rook:
		if rank == 6 {
			return 20
		}
		return 0
	case chess.Queen:
		return center * 2
	case chess.King:
		// Prefer the home rank and shelter-side files while material remains.
		shelter := 0
		if file <= 2 || file >= 5 {
			shelter = 20
		}
		middle := shelter - 12*rank - center*4
		ending := center * 8
		return (middle*phase + ending*(24-phase)) / 24
	}
	return 0
}
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func phaseWeight(piece chess.PieceType) int {
	switch piece {
	case chess.Knight, chess.Bishop:
		return 1
	case chess.Rook:
		return 2
	case chess.Queen:
		return 4
	}
	return 0
}
func staticScore(position *chess.Position, perspective chess.Color) int {
	board := position.Board()
	phase := 0
	for square := chess.A1; square <= chess.H8; square++ {
		phase += phaseWeight(board.Piece(square).Type())
	}
	phase = min(24, phase)
	score := 0
	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece == chess.NoPiece {
			continue
		}
		file, rank := int(square)%8, int(square)/8
		if piece.Color() == chess.Black {
			rank = 7 - rank
		}
		value := pieceValues[piece.Type()] + placement(piece.Type(), file, rank, phase)
		if piece.Color() == perspective {
			score += value
		} else {
			score -= value
		}
	}
	return score
}

// Retained separately for material invariants and diagnostics.
func materialScore(position *chess.Position, perspective chess.Color) int {
	score := 0
	for square := chess.A1; square <= chess.H8; square++ {
		piece := position.Board().Piece(square)
		value := pieceValues[piece.Type()]
		if piece.Color() == perspective {
			score += value
		} else {
			score -= value
		}
	}
	return score
}
