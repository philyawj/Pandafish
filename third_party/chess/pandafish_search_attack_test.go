package chess

import (
	"math/rand/v2"
	"testing"
)

func mailboxAttackAbs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func mailboxAttackSign(v int) int {
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

// Independent mailbox geometry: no bitboards, rays, attack masks or move tags.
func mailboxAttackOracle(board *Board, target Square, color Color) bool {
	for from := A1; from <= H8; from++ {
		piece := board.Piece(from)
		if piece == NoPiece || piece.Color() != color || from == target {
			continue
		}
		dx, dy := int(target)%8-int(from)%8, int(target)/8-int(from)/8
		x, y := mailboxAttackAbs(dx), mailboxAttackAbs(dy)
		sliding := false
		switch piece.Type() {
		case Pawn:
			if x == 1 && ((color == White && dy == 1) || (color == Black && dy == -1)) {
				return true
			}
		case Knight:
			if (x == 1 && y == 2) || (x == 2 && y == 1) {
				return true
			}
		case King:
			if x <= 1 && y <= 1 {
				return true
			}
		case Bishop:
			sliding = x == y
		case Rook:
			sliding = dx == 0 || dy == 0
		case Queen:
			sliding = x == y || dx == 0 || dy == 0
		}
		if !sliding {
			continue
		}
		stepX, stepY := mailboxAttackSign(dx), mailboxAttackSign(dy)
		file, rank := int(from)%8+stepX, int(from)/8+stepY
		clear := true
		for file != int(target)%8 || rank != int(target)/8 {
			if board.Piece(Square(rank*8+file)) != NoPiece {
				clear = false
				break
			}
			file += stepX
			rank += stepY
		}
		if clear {
			return true
		}
	}
	return false
}

func TestSearchAttackMatchesOriginalAndMailbox(t *testing.T) {
	rng := rand.New(rand.NewPCG(20, 2202))
	position := StartingPosition()
	checked := 0
	for i := 0; i < 1000; i++ {
		for square := A1; square <= H8; square++ {
			for _, color := range []Color{White, Black} {
				got := searchSquareAttackedBy(position.board, square, color)
				original := isSquareAttackedBy(position.board, square, color)
				oracle := mailboxAttackOracle(position.board, square, color)
				if got != original || got != oracle {
					t.Fatal("attack disagreement", position, square, color, got, original, oracle)
				}
				checked++
			}
		}
		moves := position.ValidMoves()
		if len(moves) == 0 || i%100 == 99 {
			position = StartingPosition()
		} else {
			move := moves[rng.IntN(len(moves))]
			position = position.Update(&move)
		}
	}
	t.Log("Matched attack predicates:", checked)
}
