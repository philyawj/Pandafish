package chess

// searchPieceSet addresses only the affected set. The original general board
// update remains untouched and serves as an independent transition oracle.
func (b *Board) searchPieceSet(piece Piece) *bitboard {
	switch piece {
	case WhiteKing:
		return &b.bbWhiteKing
	case WhiteQueen:
		return &b.bbWhiteQueen
	case WhiteRook:
		return &b.bbWhiteRook
	case WhiteBishop:
		return &b.bbWhiteBishop
	case WhiteKnight:
		return &b.bbWhiteKnight
	case WhitePawn:
		return &b.bbWhitePawn
	case BlackKing:
		return &b.bbBlackKing
	case BlackQueen:
		return &b.bbBlackQueen
	case BlackRook:
		return &b.bbBlackRook
	case BlackBishop:
		return &b.bbBlackBishop
	case BlackKnight:
		return &b.bbBlackKnight
	case BlackPawn:
		return &b.bbBlackPawn
	}
	return nil
}

func (b *Board) updateSearch(move *Move) {
	piece, captured := b.Piece(move.s1), b.Piece(move.s2)
	from, to := bbForSquare(move.s1), bbForSquare(move.s2)
	if set := b.searchPieceSet(captured); set != nil {
		*set &^= to
	}
	if set := b.searchPieceSet(piece); set != nil {
		*set = (*set &^ from) | to
	}
	b.mailbox[move.s1], b.mailbox[move.s2] = NoPiece, piece
	if move.promo != NoPieceType && piece != NoPiece {
		promoted := NewPiece(move.promo, piece.Color())
		*b.searchPieceSet(piece) &^= from | to
		*b.searchPieceSet(promoted) |= to
		b.mailbox[move.s2] = promoted
	}
	if move.HasTag(EnPassant) {
		capture := move.s2 + 8
		victim := WhitePawn
		if piece.Color() == White {
			capture = move.s2 - 8
			victim = BlackPawn
		}
		*b.searchPieceSet(victim) &^= bbForSquare(capture)
		b.mailbox[capture] = NoPiece
	}
	castleFrom, castleTo, rook := NoSquare, NoSquare, NoPiece
	switch {
	case piece.Color() == White && move.HasTag(KingSideCastle):
		castleFrom, castleTo, rook = H1, F1, WhiteRook
	case piece.Color() == White && move.HasTag(QueenSideCastle):
		castleFrom, castleTo, rook = A1, D1, WhiteRook
	case piece.Color() == Black && move.HasTag(KingSideCastle):
		castleFrom, castleTo, rook = H8, F8, BlackRook
	case piece.Color() == Black && move.HasTag(QueenSideCastle):
		castleFrom, castleTo, rook = A8, D8, BlackRook
	}
	if rook != NoPiece {
		set := b.searchPieceSet(rook)
		*set = (*set &^ bbForSquare(castleFrom)) | bbForSquare(castleTo)
		b.mailbox[castleFrom], b.mailbox[castleTo] = NoPiece, rook
	}
	b.calcConvienceBBs(move)
}

// UpdateSearch preserves Update's immutable state/hash semantics for generated
// legal moves, while changing only affected piece sets. Nil retains Update's
// existing null-transition behavior; no null-move search is introduced.
func (pos *Position) UpdateSearch(move *Move) *Position {
	if move == nil {
		return pos.Update(nil)
	}
	count := pos.moveCount
	if pos.turn == Black {
		count++
	}
	rights := pos.updateCastleRights(move)
	piece := pos.board.Piece(move.s1)
	halfMove := pos.halfMoveClock + 1
	if piece.Type() == Pawn || move.HasTag(Capture) {
		halfMove = 0
	}
	board := pos.board.copy()
	board.updateSearch(move)
	child := &Position{board: board, turn: pos.turn.Other(), castleRights: rights, enPassantSquare: pos.updateEnPassantSquare(move), halfMoveClock: halfMove, moveCount: count, inCheck: move.HasTag(Check)}
	child.hash = pos.updateHash(move, rights, child.enPassantSquare)
	return child
}

// Full tagging matches moveTags. An existence-only caller discards the move and
// needs only own-king legality, so opponent-check classification is optional.
func searchMoveTagsSimulated(move Move, pos *Position, classifyOpponent bool) MoveTag {
	var tags MoveTag
	piece := pos.board.Piece(move.s1)
	if pos.board.isOccupied(move.s2) {
		tags |= Capture
	} else if move.s2 == pos.enPassantSquare && piece.Type() == Pawn {
		tags |= EnPassant
	}
	if (piece == WhiteKing && move.s1 == E1) || (piece == BlackKing && move.s1 == E8) {
		switch move.s2 {
		case C1, C8:
			tags |= QueenSideCastle
		case G1, G8:
			tags |= KingSideCastle
		}
	}
	local := move
	local.tags = tags
	board := *pos.board
	board.updateSearch(&local)
	if king := board.kingSquare(pos.turn); king != NoSquare && searchSquareAttackedBy(&board, king, pos.turn.Other()) {
		tags |= inCheck
	}
	if classifyOpponent {
		if king := board.kingSquare(pos.turn.Other()); king != NoSquare && searchSquareAttackedBy(&board, king, pos.turn) {
			tags |= Check
		}
	}
	return tags
}
