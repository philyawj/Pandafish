package chess

import "math/bits"

// AppendSearchMoves appends legal moves in the same order as ValidMoves, using
// caller-owned storage without populating or mutating the position's cache.
// tacticalOnly retains captures, en passant and every promotion, including quiet
// underpromotions. A checked search must request all moves for legal evasions.
func (pos *Position) AppendSearchMoves(dst []Move, tacticalOnly bool) []Move {
	moves, _ := pos.searchMoves(dst, tacticalOnly, false)
	return moves
}

// HasLegalMove tests existence without materializing a complete move list.
func (pos *Position) HasLegalMove() bool {
	if pos.validMoves != nil {
		return len(pos.validMoves) != 0
	}
	_, found := pos.searchMoves(nil, false, true)
	return found
}

// SearchStatus has the same mate/stalemate semantics as Status, without
// allocating a legal move list merely to test whether any legal move exists.
func (pos *Position) SearchStatus() Method {
	if pos.HasLegalMove() {
		return NoMethod
	}
	if pos.inCheck {
		return Checkmate
	}
	return Stalemate
}

func (pos *Position) searchMoves(dst []Move, tacticalOnly, first bool) ([]Move, bool) {
	if pos.validMoves != nil {
		for _, m := range pos.validMoves {
			if !tacticalOnly || m.HasTag(Capture) || m.HasTag(EnPassant) || m.Promo() != NoPieceType {
				if first {
					return dst, true
				}
				dst = append(dst, m)
			}
		}
		return dst, false
	}
	allowed, targets := ^pos.board.whiteSqs, pos.board.blackSqs
	if pos.turn == Black {
		allowed, targets = ^pos.board.blackSqs, pos.board.whiteSqs
	}
	found := false
	for _, piece := range allPieces {
		if piece.Color() != pos.turn {
			continue
		}
		from := pos.board.bbForPiece(piece)
		for from != 0 {
			s1 := Square(bits.LeadingZeros64(uint64(from)))
			from &^= bbForSquare(s1)
			to := bbForPossibleMoves(pos, piece.Type(), s1) & allowed
			if tacticalOnly {
				mask := targets
				if piece.Type() == Pawn {
					if pos.enPassantSquare != NoSquare {
						mask |= bbForSquare(pos.enPassantSquare)
					}
					if pos.turn == White {
						mask |= bitboard(0xff)
					} else {
						mask |= bitboard(0xff00000000000000)
					}
				}
				to &= mask
			}
			for to != 0 {
				s2 := Square(bits.LeadingZeros64(uint64(to)))
				to &^= bbForSquare(s2)
				m := Move{s1: s1, s2: s2}
				promotion := (piece == WhitePawn && s2.Rank() == Rank8) || (piece == BlackPawn && s2.Rank() == Rank1)
				count := 1
				if promotion {
					count = len(promoPieceTypes)
				}
				for i := 0; i < count; i++ {
					if promotion {
						m.promo = promoPieceTypes[i]
					}
					m.tags = searchMoveTags(m, pos, !first)
					if m.HasTag(inCheck) {
						continue
					}
					if first {
						return dst, true
					}
					dst = append(dst, m)
					found = true
				}
			}
		}
	}
	if !tacticalOnly && !pos.inCheck && (pos.castleRights.CanCastle(pos.turn, KingSide) || pos.castleRights.CanCastle(pos.turn, QueenSide)) {
		castles := castleMoves(pos)
		if first && len(castles) > 0 {
			return dst, true
		}
		dst = append(dst, castles...)
		found = found || len(castles) > 0
	}
	return dst, found
}
