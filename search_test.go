package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

// Exhaustive minimax uses a fixed White perspective, with no pruning.
// This checks the negamax sign changes and alpha-beta bounds independently.
func exhaustiveScore(position *chess.Position, depth int) int {
	if position.Status() != chess.NoMethod {
		return evaluate(position, chess.White)
	}
	if depth == 0 {
		return exhaustiveQuietScore(position, quiescenceDepth)
	}
	best := searchInfinity
	if position.Turn() == chess.White {
		best = -searchInfinity
	}
	for _, move := range position.ValidMoves() {
		score := exhaustiveScore(position.Update(&move), depth-1)
		if position.Turn() == chess.White {
			best = max(best, score)
		} else {
			best = min(best, score)
		}
	}
	return best
}

// Reference quiescence uses fixed-color minimax without pruning or ordering.
func exhaustiveQuietScore(position *chess.Position, depth int) int {
	static := evaluate(position, chess.White)
	if depth <= 0 || position.Status() != chess.NoMethod {
		return static
	}
	checked := inCheck(position)
	best := static
	if checked {
		best = searchInfinity
		if position.Turn() == chess.White {
			best = -searchInfinity
		}
	}
	for _, move := range position.ValidMoves() {
		if !checked && !move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant) && move.Promo() == chess.NoPieceType {
			continue
		}
		score := exhaustiveQuietScore(position.Update(&move), depth-1)
		if position.Turn() == chess.White {
			best = max(best, score)
		} else {
			best = min(best, score)
		}
	}
	return best
}

func TestSearchMatchesExhaustiveMinimax(t *testing.T) {
	for _, command := range []string{
		"startpos",
		"fen 4k3/1q6/8/5N2/8/8/8/K7 w - - 0 1",
		"fen k7/8/8/8/5n2/8/1Q6/4K3 b - - 0 1",
		"fen k7/8/8/3pP3/8/8/8/7K w - d6 0 1",
		"fen 7k/5K2/6Q1/8/8/8/8/8 b - - 0 1",
		"fen 7k/6Q1/5K2/8/8/8/8/8 b - - 0 1",
	} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		before := game.FEN()
		beforeMoves := len(game.Moves())
		for depth := 0; depth <= searchDepth; depth++ {
			t.Run(fmt.Sprintf("%s/depth%d", command, depth), func(t *testing.T) {
				position := game.Position()
				want := exhaustiveScore(position, depth)
				if position.Turn() == chess.Black {
					want = -want
				}
				for _, window := range [][2]int{
					{-searchInfinity, searchInfinity},
					{want - 1, want + 1}, // Exact result inside the window.
					{want, want + 1},     // Fail low: result is an upper bound.
					{want - 1, want},     // Fail high: result is a lower bound.
				} {
					alpha, beta := window[0], window[1]
					got := search(position, depth, alpha, beta)
					switch {
					case want <= alpha:
						if got < want || got > alpha {
							t.Fatalf("fail-low score = %d, exact = %d, alpha = %d", got, want, alpha)
						}
					case want >= beta:
						if got < beta || got > want {
							t.Fatalf("fail-high score = %d, exact = %d, beta = %d", got, want, beta)
						}
					default:
						if got != want {
							t.Fatalf("search score = %d, exhaustive score = %d", got, want)
						}
					}
				}
			})
		}
		if game.FEN() != before || len(game.Moves()) != beforeMoves {
			t.Fatal("search changed the root game")
		}
	}
}

// A dense middlegame that exceeded the local match deadline before move ordering.
const timeoutPosition = "r2k2r1/5p2/Ppp1q3/3Np1Pp/2P1P1n1/Q5PB/P2BKP1N/1R5R w - - 0 28"

// These positions exceeded the deadline before value-based capture ordering.
const queenRaidPosition = "r1b1kbnr/ppQnpppp/2p5/3q4/4p3/2B4P/PPP2PP1/RN2KBNR b KQkq - 0 9"
const crowdedPosition = "r1b1kb1r/1p3pp1/p1pqpn2/3p2Bp/3P4/PRN2Q2/2P2PPP/4K1NR b Kkq - 3 15"

func BenchmarkBestMove(b *testing.B) {
	for _, tt := range []struct{ name, fen string }{
		{"dense-middlegame", timeoutPosition},
		{"queen-raid", queenRaidPosition},
		{"crowded-middlegame", crowdedPosition},
	} {
		b.Run(tt.name, func(b *testing.B) {
			game, err := parsePosition(strings.Fields("fen " + tt.fen))
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			for range b.N {
				bestMove(game, nil)
			}
		})
	}
}

func TestBestMoveMatchesFullWindowSearch(t *testing.T) {
	for _, command := range []string{
		"startpos",
		"fen 3r2k1/8/8/3p4/8/8/8/3Q2K1 w - - 0 1",
		"fen 3q2k1/8/8/8/3P4/8/8/3R2K1 b - - 0 1",
		"fen 7k/8/5KQ1/8/8/8/8/1q6 w - - 0 1",
	} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		best := -searchInfinity
		scores := make(map[string]int)
		for _, move := range game.ValidMoves() {
			uci := move.String()
			score := -search(game.Position().Update(&move), searchDepth-1, -searchInfinity, searchInfinity)
			scores[uci] = score
			best = max(best, score)
		}
		// Random tie breaking may choose any truly best move, but an inferior
		// move must never join the tie because a cutoff returned a bound.
		for range 20 {
			move := bestMove(game, nil)
			if score, ok := scores[move]; !ok || score != best {
				t.Fatalf("%s: move %s scores %d, best is %d", command, move, score, best)
			}
		}
	}
}

func TestBestMoveFindsThirdPlyQueenCapture(t *testing.T) {
	for _, tt := range []struct{ fen, moves, want string }{
		{"4k3/1q6/8/5N2/8/8/8/K7 w - - 0 1", "f5d6 f5h4", "f5d6"},
		{"k7/8/8/8/5n2/8/1Q6/4K3 b - - 0 1", "f4d3 f4h5", "f4d3"},
	} {
		game, err := parsePosition(strings.Fields("fen " + tt.fen))
		if err != nil {
			t.Fatal(err)
		}
		if got := bestMove(game, strings.Fields("searchmoves "+tt.moves)); got != tt.want {
			t.Fatalf("move = %s, want queen-winning fork %s", got, tt.want)
		}
		// Check evasions and the following capture now extend depth two as well.
		for _, check := range []struct{ depth, want int }{{2, 3}, {3, 3}} {
			if got := search(game.Position(), check.depth, -searchInfinity, searchInfinity); abs(got-check.want*100) > 90 {
				t.Fatalf("depth %d score = %d, want about %d cp", check.depth, got, check.want*100)
			}
		}
	}
}

// Keep v0.6's fixed-depth regression cases independent of timed search.
func bestMove(game *chess.Game, args []string) string {
	limits, err := parseGo(args)
	if err != nil {
		panic(err)
	}
	moves := rootMoves(game.Position(), limits.searchMoves)
	if len(moves) == 0 {
		return "0000"
	}
	best, _, _ := (&searcher{}).root(game.Position(), moves, searchDepth)
	return best[rand.IntN(len(best))]
}
