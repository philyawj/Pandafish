package main

import (
	"context"
	"slices"
	"time"

	"github.com/corentings/chess/v2"
)

const searchDepth = 3
const searchInfinity = mateScore + 1
const quiescenceDepth = 16
const checkEvasionDepth = 4

// A zero-value searcher is unlimited, useful for fixed-depth comparisons.
// Cancellation is distinct from a score, so unfinished nodes never affect ties.
type searcher struct {
	moveBuffers [][]chess.Move
	repetitions map[uint64]int
	historyKey  uint64
	nodes       uint64
	table       *transpositionTable
	killers     *killerTable
	ctx         context.Context
	deadline    time.Time
	now         func() time.Time
	unfiltered  bool // Internal reference mode for broader-quiescence ablations.
}

func (worker *searcher) clock() time.Time {
	if worker.now != nil {
		return worker.now()
	}
	return time.Now()
}

func (worker *searcher) stopped() bool {
	return (worker.ctx != nil && worker.ctx.Err() != nil) ||
		(!worker.deadline.IsZero() && !worker.clock().Before(worker.deadline))
}

func search(position *chess.Position, depth, alpha, beta int) int {
	score, _ := (&searcher{}).search(position, depth, alpha, beta, inCheck(position))
	return score
}

// Negamax flips each child's score and window to the parent's perspective.
func (worker *searcher) search(position *chess.Position, depth, alpha, beta int, checked bool) (int, bool) {
	return worker.searchAtPly(position, depth, alpha, beta, checked, 0)
}

func (worker *searcher) searchAtPly(position *chess.Position, depth, alpha, beta int, checked bool, ply int) (int, bool) {
	if worker.stopped() {
		return 0, false
	}
	if depth <= 0 {
		return worker.quietAtPly(position, quiescenceDepth, alpha, beta, checked, ply)
	}
	worker.nodes++
	if position.SearchStatus() != chess.NoMethod {
		return searchedEvaluation(position, ply), !worker.stopped()
	}
	if worker.pushPosition(position) {
		worker.popPosition(position)
		return 0, !worker.stopped()
	}
	defer worker.popPosition(position)
	entry, hit := worker.table.probe(position, worker.historyKey)
	// Equal depth preserves the exact fixed-depth value. Deeper entries still
	// provide move hints. Quiescence has a separate horizon and is not cached.
	if hit && int(entry.depth) == depth {
		score := scoreFromTable(int(entry.score), ply)
		if entry.bound == exactBound || (entry.bound == lowerBound && score >= beta) || (entry.bound == upperBound && score <= alpha) {
			return score, !worker.stopped()
		}
	}
	originalAlpha := alpha
	moves := worker.searchMoves(position, ply, false)
	worker.killers.order(moves, ply)
	if hit {
		preferMove(moves, entry.move)
	} else {
		preferMove(moves, worker.table.boardHint(position))
	}
	best, bestMove := -searchInfinity, uint16(0)
	for moveIndex, move := range moves {
		if worker.stopped() {
			return 0, false
		}
		child := position.UpdateSearch(&move)
		var score int
		var complete bool
		if moveIndex == 0 {
			score, complete = worker.childSearch(child, depth-1, -beta, -alpha, move.HasTag(chess.Check), ply+1)
		} else {
			// Scout later moves with a null window. Re-search an in-window
			// improvement at full depth; this introduces no move reductions.
			score, complete = worker.childSearch(child, depth-1, -alpha-1, -alpha, move.HasTag(chess.Check), ply+1)
			if complete && -score > alpha && -score < beta {
				score, complete = worker.childSearch(child, depth-1, -beta, -alpha, move.HasTag(chess.Check), ply+1)
			}
		}
		if !complete {
			return 0, false
		}
		score = -score
		if score > best {
			best, bestMove = score, encodeMove(move)
		}
		alpha = max(alpha, score)
		if alpha >= beta {
			worker.killers.remember(ply, move)
			break
		}
	}
	if worker.stopped() {
		return 0, false
	}
	worker.table.store(position, depth, scoreToTable(best, ply), scoreToTable(originalAlpha, ply), scoreToTable(beta, ply), bestMove, worker.historyKey)
	worker.table.storeHint(position, bestMove)
	return best, true
}

func quiescence(position *chess.Position, depth, alpha, beta int) int {
	score, _ := (&searcher{}).quiet(position, depth, alpha, beta, inCheck(position))
	return score
}

// Continue tactical exchanges until quiet or the safety bound is reached.
// Non-checking ordinary captures may be filtered by their material exchange.
// Checking captures, promotions, and every check evasion remain eligible.
func (worker *searcher) quiet(position *chess.Position, depth, alpha, beta int, checked bool) (int, bool) {
	return worker.quietAtPly(position, depth, alpha, beta, checked, 0)
}

func (worker *searcher) quietAtPly(position *chess.Position, depth, alpha, beta int, checked bool, ply int) (int, bool) {
	if worker.stopped() {
		return 0, false
	}
	worker.nodes++
	if position.SearchStatus() != chess.NoMethod {
		return searchedEvaluation(position, ply), !worker.stopped()
	}
	if worker.pushPosition(position) {
		worker.popPosition(position)
		return 0, !worker.stopped()
	}
	defer worker.popPosition(position)
	static := staticScore(position, position.Turn())
	if depth <= -checkEvasionDepth || (depth <= 0 && !checked) {
		return static, !worker.stopped()
	}
	best := -searchInfinity
	if !checked {
		best = static
		if best >= beta {
			return best, !worker.stopped()
		}
		alpha = max(alpha, best)
	}
	moves := worker.searchMoves(position, ply, !checked)
	for _, move := range moves {
		if worker.stopped() {
			return 0, false
		}
		if !checked && !move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant) && move.Promo() == chess.NoPieceType {
			continue
		}
		// Equal-or-cheaper attackers usually cannot lose a material exchange.
		// Keep them without SEE, including recapture-promotion exceptions.
		if !worker.unfiltered && !checked && !move.HasTag(chess.Check) && move.Promo() == chess.NoPieceType &&
			captureValue(position, move) < pieceValues[position.Board().Piece(move.S1()).Type()] {
			gain, complete := worker.exchange(position, move)
			if !complete {
				return 0, false
			}
			if gain < 0 {
				continue
			}
		}
		score, complete := worker.childQuiet(position.UpdateSearch(&move), depth-1, -beta, -alpha, move.HasTag(chess.Check), ply+1)
		if !complete {
			return 0, false
		}
		score = -score
		best, alpha = max(best, score), max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best, !worker.stopped()
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
			value += 10000 + 10*victim - pieceValues[position.Board().Piece(move.S1()).Type()]
		}
		return value
	}
	slices.SortStableFunc(moves, func(a, b chess.Move) int {
		return priority(b) - priority(a)
	})
	return moves
}

// The library still generates the full legal list. Copy eligible moves before
// sorting so its cached list is never mutated. Specialized generation is a
// separate experiment; eligibility and stable tactical ordering stay unchanged.
func orderedTacticalMoves(position *chess.Position) []chess.Move {
	var moves []chess.Move
	for _, m := range position.ValidMovesUnsafe() {
		if m.HasTag(chess.Capture) || m.HasTag(chess.EnPassant) || m.Promo() != chess.NoPieceType {
			moves = append(moves, m)
		}
	}
	priority := func(move chess.Move) int {
		value := pieceValues[move.Promo()]
		if move.HasTag(chess.Capture) || move.HasTag(chess.EnPassant) {
			victim := pieceValues[position.Board().Piece(move.S2()).Type()]
			if move.HasTag(chess.EnPassant) {
				victim = pieceValues[chess.Pawn]
			}
			value += 10000 + 10*victim - pieceValues[position.Board().Piece(move.S1()).Type()]
		}
		return value
	}
	slices.SortStableFunc(moves, func(a, b chess.Move) int {
		return priority(b) - priority(a)
	})
	return moves
}
