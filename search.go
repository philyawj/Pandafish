package main

import (
	"slices"

	"github.com/corentings/chess/v2"
)

const searchDepth = 3
const searchInfinity = mateScore + 1
const quiescenceDepth = 2

// search uses negamax: each side maximizes its own score. Negating a child's
// score gives the parent's perspective. Alpha-beta skips replies that cannot
// improve the result within the current window.
func search(position *chess.Position, depth, alpha, beta int) int {
	return searchPosition(position, depth, alpha, beta, inCheck(position))
}

// Check tags on legal moves let descendants reuse the library's check result.
func searchPosition(position *chess.Position, depth, alpha, beta int, checked bool) int {
	if depth <= 0 {
		return quiescencePosition(position, quiescenceDepth, alpha, beta, checked)
	}
	if position.Status() != chess.NoMethod {
		return evaluate(position, position.Turn())
	}
	best := -searchInfinity
	for _, move := range orderedMoves(position) {
		score := -searchPosition(position.Update(&move), depth-1, -beta, -alpha, move.HasTag(chess.Check))
		best = max(best, score)
		alpha = max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best
}

// quiescence follows captures, promotions, and all legal check evasions.
// Away from check, the static score (stand pat) is an alternative to capturing.
// A hard depth cap bounds the work; its leaves can still contain tactics.
func quiescence(position *chess.Position, depth, alpha, beta int) int {
	return quiescencePosition(position, depth, alpha, beta, inCheck(position))
}

func quiescencePosition(position *chess.Position, depth, alpha, beta int, checked bool) int {
	if position.Status() != chess.NoMethod {
		return evaluate(position, position.Turn())
	}
	static := materialScore(position, position.Turn())
	if depth <= 0 {
		return static
	}
	best := -searchInfinity
	if !checked {
		best = static
		if best >= beta {
			return best
		}
		alpha = max(alpha, best)
	}
	for _, move := range orderedMoves(position) {
		if !checked && !move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant) && move.Promo() == chess.NoPieceType {
			continue
		}
		score := -quiescencePosition(position.Update(&move), depth-1, -beta, -alpha, move.HasTag(chess.Check))
		best = max(best, score)
		alpha = max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best
}

// The chess library exposes attacks, but not Position's check flag.
func inCheck(position *chess.Position) bool {
	pieces := position.Board().SquareMap()
	var king chess.Square
	for square, piece := range pieces {
		if piece.Type() == chess.King && piece.Color() == position.Turn() {
			king = square
			break
		}
	}
	for square, piece := range pieces {
		if piece.Color() != position.Turn() && slices.Contains(position.Board().AttacksFrom(square), king) {
			return true
		}
	}
	return false
}

// Search captures of valuable pieces first, preferring cheaper attackers.
// ValidMoves returns a copy, so reordering does not change the position.
func orderedMoves(position *chess.Position) []chess.Move {
	moves := position.ValidMoves()
	priority := func(move chess.Move) int {
		value := pieceValues[move.Promo()]
		if move.HasTag(chess.Capture) || move.HasTag(chess.EnPassant) {
			victim := pieceValues[position.Board().Piece(move.S2()).Type()]
			if move.HasTag(chess.EnPassant) {
				victim = pieceValues[chess.Pawn]
			}
			value += 100 + 10*victim - pieceValues[position.Board().Piece(move.S1()).Type()]
		}
		return value
	}
	slices.SortStableFunc(moves, func(a, b chess.Move) int {
		return priority(b) - priority(a)
	})
	return moves
}
