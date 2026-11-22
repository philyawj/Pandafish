package chess

import "math/bits"

// searchContactTargets is a geometric filter, not a legal threat or mate claim.
// Only empty king-neighbour squares protected by a pawn, knight or king qualify.
func (pos *Position) searchContactTargets(attacker Color) (bitboard, bitboard) {
	b := pos.board
	queens, pawns, knights, king, victim := b.bbWhiteQueen, b.bbWhitePawn, b.bbWhiteKnight, b.bbWhiteKing, b.bbBlackKing
	if attacker == Black {
		queens, pawns, knights, king, victim = b.bbBlackQueen, b.bbBlackPawn, b.bbBlackKnight, b.bbBlackKing, b.bbWhiteKing
	}
	if queens == 0 || bits.OnesCount64(uint64(king)) != 1 || bits.OnesCount64(uint64(victim)) != 1 {
		return 0, 0
	}
	targets := bbKingMoves[bits.LeadingZeros64(uint64(victim))] & b.emptySqs
	if targets == 0 {
		return 0, 0
	}
	protection := bbKingMoves[bits.LeadingZeros64(uint64(king))]
	if attacker == White {
		protection |= (pawns&^bbFileH&^bbRank8)>>9 | (pawns&^bbFileA&^bbRank8)>>7
	} else {
		protection |= (pawns&^bbFileH&^bbRank1)<<7 | (pawns&^bbFileA&^bbRank1)<<9
	}
	for knights != 0 {
		square := bits.LeadingZeros64(uint64(knights))
		knights &^= bbForSquare(Square(square))
		protection |= bbKnightMoves[square]
	}
	return queens, targets & protection
}

// SearchHasProtectedQueenContact tests geometry only. A true result requires
// legal move and mate verification; a false result excludes only this subset.
func (pos *Position) SearchHasProtectedQueenContact(attacker Color) bool {
	queens, targets := pos.searchContactTargets(attacker)
	if targets == 0 {
		return false
	}
	occupied := ^pos.board.emptySqs
	for queens != 0 {
		square := bits.LeadingZeros64(uint64(queens))
		queens &^= bbForSquare(Square(square))
		if (hvAttack(occupied, Square(square))|diaAttack(occupied, Square(square)))&targets != 0 {
			return true
		}
	}
	return false
}

// AppendSearchProtectedQueenContacts returns legal quiet queen checks within
// the protected-contact subset. Caller storage and cached metadata are preserved.
func (pos *Position) AppendSearchProtectedQueenContacts(dst []Move) []Move {
	queens, targets := pos.searchContactTargets(pos.turn)
	occupied := ^pos.board.emptySqs
	for queens != 0 && targets != 0 {
		from := Square(bits.LeadingZeros64(uint64(queens)))
		queens &^= bbForSquare(from)
		to := (hvAttack(occupied, from) | diaAttack(occupied, from)) & targets
		for to != 0 {
			square := Square(bits.LeadingZeros64(uint64(to)))
			to &^= bbForSquare(square)
			if m, ok := pos.LookupSearchMove(from, square, NoPieceType); ok && m.HasTag(Check) && !m.HasTag(Capture) {
				dst = append(dst, m)
			}
		}
	}
	return dst
}
