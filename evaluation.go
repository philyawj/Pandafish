package main

import "github.com/corentings/chess/v2"

// A mating result outweighs any material balance in a legal position.
const mateScore = 10000

// Piece values are measured in pawns. Kings are never captured in legal chess.
var pieceValues = map[chess.PieceType]int{
	chess.Pawn: 1, chess.Knight: 3, chess.Bishop: 3, chess.Rook: 5, chess.Queen: 9,
}

// evaluate scores checkmate and stalemate before comparing material.
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
