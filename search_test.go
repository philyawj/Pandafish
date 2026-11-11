package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

// Exhaustive minimax uses a fixed White perspective, with no pruning.
// This checks the negamax sign changes and alpha-beta bounds independently.
func exhaustiveScore(position *chess.Position, depth int) int {
	if depth == 0 || position.Status() != chess.NoMethod {
		return evaluate(position, chess.White)
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

func TestSearchMatchesExhaustiveMinimax(t *testing.T) {
	for _, command := range []string{
		"startpos",
		"startpos moves e2e4 e7e5 g1f3",
		"fen 4k3/1q6/8/5N2/8/8/8/K7 w - - 0 1",
		"fen k7/8/8/8/5n2/8/1Q6/4K3 b - - 0 1",
		"fen 7k/P7/8/8/8/8/7p/K7 w - - 0 1",
		"fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		"fen k7/8/8/3pP3/8/8/8/7K w - d6 0 1",
		"fen " + timeoutPosition,
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

func BenchmarkBestMove(b *testing.B) {
	game, err := parsePosition(strings.Fields("fen " + timeoutPosition))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		bestMove(game, nil)
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
		for _, check := range []struct{ depth, want int }{{2, -6}, {3, 3}} {
			if got := search(game.Position(), check.depth, -searchInfinity, searchInfinity); got != check.want {
				t.Fatalf("depth %d score = %d, want %d", check.depth, got, check.want)
			}
		}
	}
}
