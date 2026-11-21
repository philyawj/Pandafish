package chess

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

func TestVirtualTagsMatchBothTransitionOracles(t *testing.T) {
	rng := rand.New(rand.NewPCG(23, 1301))
	positions := []*Position{StartingPosition()}
	for _, fen := range []string{
		"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"4k3/P7/8/8/8/8/7p/4K3 w - - 0 1",
		"4k3/8/8/8/8/8/8/R3R2R w KQ - 0 1",
		"8/8/8/8/8/8/8/8 w - - 0 1",
		"4k3/8/8/8/8/8/8/3KK3 w - - 0 1",
	} {
		option, err := FEN(fen)
		if err != nil {
			t.Fatal(err)
		}
		positions = append(positions, NewGame(option).Position())
	}
	p := StartingPosition()
	for i := 0; i < 1500; i++ {
		positions = append(positions, p)
		moves := p.ValidMoves()
		if len(moves) == 0 || i%100 == 99 {
			p = StartingPosition()
		} else {
			p = p.Update(&moves[rng.IntN(len(moves))])
		}
	}
	compared := 0
	for _, p := range positions {
		before, hash := p.String(), p.ZobristHash()
		for _, piece := range allPieces {
			if piece.Color() != p.turn {
				continue
			}
			from := p.board.bbForPiece(piece)
			for from != 0 {
				s1 := Square(bits.LeadingZeros64(uint64(from)))
				from &^= bbForSquare(s1)
				to := bbForPossibleMoves(p, piece.Type(), s1)
				for to != 0 {
					s2 := Square(bits.LeadingZeros64(uint64(to)))
					to &^= bbForSquare(s2)
					promotions := []PieceType{NoPieceType}
					if piece.Type() == Pawn && (s2.Rank() == Rank1 || s2.Rank() == Rank8) {
						promotions = promoPieceTypes[:]
					}
					for _, promo := range promotions {
						m := Move{s1: s1, s2: s2, promo: promo}
						for _, classify := range []bool{false, true} {
							got, want := searchMoveTags(m, p, classify), searchMoveTagsSimulated(m, p, classify)
							original := moveTags(m, p)
							if !classify {
								original &^= Check
							}
							if got != want || got != original {
								t.Fatal(before, m, classify, got, want, original)
							}
							compared++
						}
					}
				}
			}
		}
		if before != p.String() || hash != p.ZobristHash() {
			t.Fatal("parent mutated")
		}
	}
	t.Log("pseudo-legal tag comparisons", compared)
}
