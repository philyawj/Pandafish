package chess

// searchSquareAttackedBy preserves the original geometric attack predicate.
// Rays are computed once and piece sets are selected directly. The initial
// geometric rejection is identical to the pinned implementation.
func searchSquareAttackedBy(board *Board, square Square, attacker Color) bool {
	occupied := ^board.emptySqs
	hv, diagonal := hvAttack(occupied, square), diaAttack(occupied, square)
	knight := bbKnightMoves[square]
	side, queens, rooks, bishops, knights, pawns, kings := board.blackSqs, board.bbBlackQueen, board.bbBlackRook, board.bbBlackBishop, board.bbBlackKnight, board.bbBlackPawn, board.bbBlackKing
	if attacker == White {
		side, queens, rooks, bishops, knights, pawns, kings = board.whiteSqs, board.bbWhiteQueen, board.bbWhiteRook, board.bbWhiteBishop, board.bbWhiteKnight, board.bbWhitePawn, board.bbWhiteKing
	}
	if ((hv|diagonal)&side)|(knight&side) == 0 {
		return false
	}
	if hv&(queens|rooks) != 0 || diagonal&(queens|bishops) != 0 || knight&knights != 0 {
		return true
	}
	var pawnAttacks bitboard
	if attacker == White {
		pawnAttacks = ((pawns &^ bbFileH &^ bbRank8) >> 9) | ((pawns &^ bbFileA &^ bbRank8) >> 7)
	} else {
		pawnAttacks = ((pawns &^ bbFileH &^ bbRank1) << 7) | ((pawns &^ bbFileA &^ bbRank1) << 9)
	}
	return pawnAttacks&bbForSquare(square) != 0 || bbKingMoves[square]&kings != 0
}
