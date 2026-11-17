package chess

// EvaluationPieces returns a value copy of the existing piece bitboards, in
// King/Queen/Rook/Bishop/Knight/Pawn order, White then Black. It allocates no
// list/map and permits no mutation of board state through the returned value.
func (b *Board) EvaluationPieces() [12]uint64 {
	return [12]uint64{
		uint64(b.bbWhiteKing), uint64(b.bbWhiteQueen), uint64(b.bbWhiteRook), uint64(b.bbWhiteBishop), uint64(b.bbWhiteKnight), uint64(b.bbWhitePawn),
		uint64(b.bbBlackKing), uint64(b.bbBlackQueen), uint64(b.bbBlackRook), uint64(b.bbBlackBishop), uint64(b.bbBlackKnight), uint64(b.bbBlackPawn),
	}
}
