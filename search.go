package main

import (
	"slices"

	"github.com/corentings/chess/v2"
)

const searchDepth = 3
const searchInfinity = mateScore + 1

// search uses negamax: each side maximizes its own score. Negating a child's
// score gives the parent's perspective. Alpha-beta skips replies that cannot
// improve the result within the current window.
func search(position *chess.Position, depth, alpha, beta int) int {
	if depth <= 0 || position.Status() != chess.NoMethod {
		return evaluate(position, position.Turn())
	}
	best := -searchInfinity
	for _, move := range orderedMoves(position) {
		score := -search(position.Update(&move), depth-1, -beta, -alpha)
		best = max(best, score)
		alpha = max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best
}

// Searching captures first helps alpha-beta find useful bounds earlier.
// ValidMoves returns a copy, so reordering does not change the position.
func orderedMoves(position *chess.Position) []chess.Move {
	moves := position.ValidMoves()
	capture := func(move chess.Move) int {
		if move.HasTag(chess.Capture) || move.HasTag(chess.EnPassant) {
			return 1
		}
		return 0
	}
	slices.SortStableFunc(moves, func(a, b chess.Move) int {
		return capture(b) - capture(a)
	})
	return moves
}
