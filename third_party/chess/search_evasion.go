package chess

import "math/bits"

// Lines include the destination and exclude the origin. Integer shifts avoid
// depending on the upstream square-bit lookup's later init function.
var searchEvasionLines = func() (lines [64][64]bitboard) {
	for from := 0; from < 64; from++ {
		for to := 0; to < 64; to++ {
			df, dr := to%8-from%8, to/8-from/8
			if from == to || (df != 0 && dr != 0 && df*df != dr*dr) {
				continue
			}
			stepF, stepR := 0, 0
			if df > 0 {
				stepF = 1
			} else if df < 0 {
				stepF = -1
			}
			if dr > 0 {
				stepR = 1
			} else if dr < 0 {
				stepR = -1
			}
			for f, r := from%8+stepF, from/8+stepR; ; f, r = f+stepF, r+stepR {
				square := r*8 + f
				lines[from][to] |= bitboard(1) << uint(63-square)
				if square == to {
					break
				}
			}
		}
	}
	return lines
}()

var searchEvasionPawnOrigins = func() (origins [3][64]bitboard) {
	for _, attacker := range [2]Color{White, Black} {
		for king := 0; king < 64; king++ {
			rank := king/8 + 1
			if attacker == White {
				rank = king/8 - 1
			}
			if rank < 0 || rank >= 8 {
				continue
			}
			for _, df := range [2]int{-1, 1} {
				file := king%8 + df
				if file >= 0 && file < 8 {
					origins[attacker][king] |= bitboard(1) << uint(63-(rank*8+file))
				}
			}
		}
	}
	return origins
}()

// searchEvasionMask is only a geometric prefilter for ordinary nonking moves.
// En passant bypasses it; all retained candidates still need full legal tagging.
func searchEvasionMask(pos *Position) (bitboard, bool) {
	if !pos.inCheck {
		return 0, false
	}
	kings := pos.board.bbWhiteKing
	if pos.turn == Black {
		kings = pos.board.bbBlackKing
	}
	if bits.OnesCount64(uint64(kings)) != 1 {
		return 0, false
	}
	king := Square(bits.LeadingZeros64(uint64(kings)))
	queens, rooks, bishops := pos.board.bbBlackQueen, pos.board.bbBlackRook, pos.board.bbBlackBishop
	knights, pawns, enemyKings := pos.board.bbBlackKnight, pos.board.bbBlackPawn, pos.board.bbBlackKing
	enemy := Black
	if pos.turn == Black {
		queens, rooks, bishops = pos.board.bbWhiteQueen, pos.board.bbWhiteRook, pos.board.bbWhiteBishop
		knights, pawns, enemyKings = pos.board.bbWhiteKnight, pos.board.bbWhitePawn, pos.board.bbWhiteKing
		enemy = White
	}
	occupied := ^pos.board.emptySqs
	checkers := (hvAttack(occupied, king) & (queens | rooks)) | (diaAttack(occupied, king) & (queens | bishops)) |
		(bbKnightMoves[king] & knights) | (bbKingMoves[king] & enemyKings) | (searchEvasionPawnOrigins[enemy][king] & pawns)
	if checkers == 0 {
		return 0, false
	}
	if checkers&(checkers-1) != 0 {
		return 0, true
	}
	checker := Square(bits.LeadingZeros64(uint64(checkers)))
	return searchEvasionLines[king][checker] | bbForSquare(checker), true
}
