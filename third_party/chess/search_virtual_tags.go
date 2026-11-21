package chess

// Ordinary moves need only occupancy and the affected attacker sets, not a
// simulated Board/mailbox. Specials retain the original transition/tag path.
func searchMoveTags(move Move, pos *Position, classifyOpponent bool) MoveTag {
	piece, captured := pos.board.Piece(move.s1), pos.board.Piece(move.s2)
	if piece == NoPiece || move.promo != NoPieceType ||
		(piece.Type() == Pawn && captured == NoPiece && move.s2 == pos.enPassantSquare) ||
		((piece == WhiteKing && move.s1 == E1) || (piece == BlackKing && move.s1 == E8)) &&
			(move.s2 == C1 || move.s2 == G1 || move.s2 == C8 || move.s2 == G8) {
		return searchMoveTagsSimulated(move, pos, classifyOpponent)
	}
	var tags MoveTag
	if captured != NoPiece {
		tags |= Capture
	}
	from, to := bbForSquare(move.s1), bbForSquare(move.s2)
	occupied := (^pos.board.emptySqs &^ from) | to
	whiteKing, blackKing := pos.board.whiteKingSq, pos.board.blackKingSq
	if move.s1 == whiteKing {
		whiteKing = move.s2
	} else if move.s1 == blackKing {
		blackKing = move.s2
	}
	own, opponent := whiteKing, blackKing
	if pos.turn == Black {
		own, opponent = blackKing, whiteKing
	}
	if own != NoSquare && searchAttackedAfterOrdinary(pos.board, own, pos.turn.Other(), piece, captured, from, to, occupied) {
		tags |= inCheck
	}
	if classifyOpponent && opponent != NoSquare && searchAttackedAfterOrdinary(pos.board, opponent, pos.turn, piece, captured, from, to, occupied) {
		tags |= Check
	}
	return tags
}

func searchAttackedAfterOrdinary(board *Board, square Square, attacker Color, piece, captured Piece, from, to, occupied bitboard) bool {
	side, queens, rooks, bishops, knights, pawns, kings := board.blackSqs, board.bbBlackQueen, board.bbBlackRook, board.bbBlackBishop, board.bbBlackKnight, board.bbBlackPawn, board.bbBlackKing
	if attacker == White {
		side, queens, rooks, bishops, knights, pawns, kings = board.whiteSqs, board.bbWhiteQueen, board.bbWhiteRook, board.bbWhiteBishop, board.bbWhiteKnight, board.bbWhitePawn, board.bbWhiteKing
	}
	if captured.Color() == attacker {
		side &^= to
		switch captured.Type() {
		case King:
			kings &^= to
		case Queen:
			queens &^= to
		case Rook:
			rooks &^= to
		case Bishop:
			bishops &^= to
		case Knight:
			knights &^= to
		case Pawn:
			pawns &^= to
		}
	}
	if piece.Color() == attacker {
		side = (side &^ from) | to
		switch piece.Type() {
		case King:
			kings = (kings &^ from) | to
		case Queen:
			queens = (queens &^ from) | to
		case Rook:
			rooks = (rooks &^ from) | to
		case Bishop:
			bishops = (bishops &^ from) | to
		case Knight:
			knights = (knights &^ from) | to
		case Pawn:
			pawns = (pawns &^ from) | to
		}
	}
	hv, diagonal := hvAttack(occupied, square), diaAttack(occupied, square)
	knight := bbKnightMoves[square]
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
