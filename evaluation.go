package main

import "github.com/corentings/chess/v2"

// Piece values are measured in pawns. Kings are never captured in legal chess.
var pieceValues = map[chess.PieceType]int{
	chess.Pawn: 1, chess.Knight: 3, chess.Bishop: 3, chess.Rook: 5, chess.Queen: 9,
}

// evaluate returns own material minus the opponent's material.
func evaluate(position *chess.Position, perspective chess.Color) int {
	score := 0
	for _, piece := range position.Board().SquareMap() {
		value := pieceValues[piece.Type()]
		if piece.Color() == perspective {
			score += value
		} else {
			score -= value
		}
	}
	return score
}
