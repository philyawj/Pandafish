package chess

// LookupSearchMove matches a signature under the existing generator's rules
// without constructing a legal list.
// It returns the same tagged Move as full generation, including cached metadata
// when a complete legal cache already exists. Parent state/cache stays unchanged.
func (pos *Position) LookupSearchMove(from, to Square, promo PieceType) (Move, bool) {
	if from < A1 || from > H8 || to < A1 || to > H8 || from == to {
		return Move{}, false
	}
	if pos.validMoves != nil {
		for _, m := range pos.validMoves {
			if m.s1 == from && m.s2 == to && m.promo == promo {
				return m, true
			}
		}
		return Move{}, false
	}
	// Preserve the pinned generator's handling of supplied castling rights,
	// including inconsistent FENs. Standard positions take the same castle path.
	castleSignature := (from == E1 && (to == C1 || to == G1)) || (from == E8 && (to == C8 || to == G8))
	if castleSignature && promo == NoPieceType && !pos.inCheck {
		for _, m := range castleMoves(pos) {
			if m.s1 == from && m.s2 == to {
				return m, true
			}
		}
	}
	piece := pos.board.Piece(from)
	if piece == NoPiece || piece.Color() != pos.turn {
		return Move{}, false
	}
	own := pos.board.whiteSqs
	if pos.turn == Black {
		own = pos.board.blackSqs
	}
	if own&bbForSquare(to) != 0 {
		return Move{}, false
	}
	promotes := (piece == WhitePawn && to.Rank() == Rank8) || (piece == BlackPawn && to.Rank() == Rank1)
	if promotes {
		switch promo {
		case Queen, Rook, Bishop, Knight:
		default:
			return Move{}, false
		}
	} else if promo != NoPieceType {
		return Move{}, false
	}
	if bbForPossibleMoves(pos, piece.Type(), from)&bbForSquare(to) == 0 {
		return Move{}, false
	}
	m := Move{s1: from, s2: to, promo: promo}
	m.tags = searchMoveTags(m, pos, true)
	if m.HasTag(inCheck) {
		return Move{}, false
	}
	return m, true
}
