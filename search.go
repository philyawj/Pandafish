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
	contactPass, contactChild chess.SearchPosition
	positionBuffers           []*chess.SearchPosition
	moveBuffers               [][]chess.Move
	rankBuffers               [][]moveRank
	repetitions               map[uint64]int
	historyKey                uint64
	nodes                     uint64
	table                     *transpositionTable
	killers                   *killerTable
	ctx                       context.Context
	deadline                  time.Time
	now                       func() time.Time
	unfiltered                bool // Internal reference mode for broader-quiescence ablations.
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
	hint := uint16(0)
	if hit {
		hint = entry.move
	} else {
		hint = worker.table.boardHint(position)
	}
	picker := newMainMovePicker(worker, position, ply, hint)
	best, bestMove := -searchInfinity, uint16(0)
	moveIndex := 0
	for move, ok := picker.next(); ok; move, ok = picker.next() {
		if worker.stopped() {
			return 0, false
		}
		child := worker.childPosition(position, &move, ply+1)
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
		moveIndex++
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
	return worker.quietContactAtPly(position, depth, alpha, beta, checked, ply, 1, nil)
}

// One contact-threat event per q path uses a window-independent frontier.
// A pending move is revalidated on the actual board after a legal defense.
func (worker *searcher) quietContactAtPly(position *chess.Position, depth, alpha, beta int, checked bool, ply, budget int, pending *chess.Move) (int, bool) {
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
	if pending != nil && depth > 0 {
		if m, ok := position.LookupSearchMove(pending.S1(), pending.S2(), pending.Promo()); ok && m.HasTag(chess.Check) {
			if worker.contactChild.Update(position, &m).SearchStatus() == chess.Checkmate {
				return mateScore - ply - 1, !worker.stopped()
			}
		}
	}
	static := staticScore(position, position.Turn())
	if depth <= -checkEvasionDepth || (depth <= 0 && !checked) {
		return static, !worker.stopped()
	}
	best := -searchInfinity
	var threat chess.Move
	defending := false
	if !checked && depth >= 2 && budget > 0 {
		threat, defending = passedContactThreat(position, &worker.contactPass, &worker.contactChild)
	}
	if !checked && !defending {
		best = static
		if best >= beta {
			return best, !worker.stopped()
		}
		alpha = max(alpha, best)
	}
	moves := worker.searchMoves(position, ply, !checked && !defending)
	for _, move := range moves {
		if worker.stopped() {
			return 0, false
		}
		if !checked && !defending && !move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant) && move.Promo() == chess.NoPieceType {
			continue
		}
		// Equal-or-cheaper attackers usually cannot lose a material exchange.
		// Keep them without SEE, including recapture-promotion exceptions.
		if !worker.unfiltered && !checked && !defending && !move.HasTag(chess.Check) && move.Promo() == chess.NoPieceType &&
			captureValue(position, move) < pieceValues[position.Board().Piece(move.S1()).Type()] {
			gain, complete := worker.exchange(position, move)
			if !complete {
				return 0, false
			}
			if gain < 0 {
				continue
			}
		}
		nextBudget := budget
		var nextThreat *chess.Move
		if defending {
			nextBudget--
			nextThreat = &threat
		}
		score, complete := worker.childContactQuiet(worker.childPosition(position, &move, ply+1), depth-1, -beta, -alpha, move.HasTag(chess.Check), ply+1, nextBudget, nextThreat)
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
