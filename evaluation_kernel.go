package main

import "github.com/corentings/chess/v2"

var evaluationPieceTypes = [12]chess.PieceType{
	chess.King, chess.Queen, chess.Rook, chess.Bishop, chess.Knight, chess.Pawn,
	chess.King, chess.Queen, chess.Rook, chess.Bishop, chess.Knight, chess.Pawn,
}

// Cache exactly the original integer formulas for every possible phase and
// piece/square. This changes no material value, placement term or rounding.
var basePieceScores = func() [25][12][64]int16 {
	var scores [25][12][64]int16
	for phase := 0; phase <= 24; phase++ {
		for index, kind := range evaluationPieceTypes {
			for sq := 0; sq < 64; sq++ {
				rank, sign := sq/8, 1
				if index >= 6 {
					rank, sign = 7-rank, -1
				}
				scores[phase][index][sq] = int16(sign * (pieceValues[kind] + placement(kind, sq%8, rank, phase)))
			}
		}
	}
	return scores
}()
