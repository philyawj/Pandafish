package main

import (
	"context"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

type searchResult struct {
	move         string
	depth, score int
	metrics      telemetry.Metrics
}

func rootMoves(position *chess.Position, restriction []string) []chess.Move {
	moves := orderedMoves(position)
	if restriction != nil {
		moves = slices.DeleteFunc(moves, func(move chess.Move) bool {
			return !slices.Contains(restriction, move.String())
		})
	}
	return moves
}

func iterativeSearch(ctx context.Context, position *chess.Position, limits goLimits, budget timeBudget, start time.Time) searchResult {
	worker := searcher{ctx: ctx}
	if budget.timed {
		worker.deadline = start.Add(budget.hard)
	}
	return worker.iterate(position, rootMoves(position, limits.searchMoves), limits.depth, budget, start)
}

func (worker *searcher) iterate(position *chess.Position, moves []chess.Move, depthLimit int, budget timeBudget, start time.Time) (result searchResult) {
	result = searchResult{move: "0000"}
	reason := "depth"
	defer func() {
		if worker.ctx != nil && worker.ctx.Err() != nil {
			reason = "cancelled"
		} else if !worker.deadline.IsZero() && !worker.clock().Before(worker.deadline) {
			reason = "deadline"
		}
		result.metrics = telemetry.Metrics{
			Depth:     result.depth,
			Nodes:     worker.nodes,
			ElapsedMS: float64(time.Since(start)) / float64(time.Millisecond),
			Timed:     budget.timed,
			SoftMS:    float64(budget.soft) / float64(time.Millisecond),
			HardMS:    float64(budget.hard) / float64(time.Millisecond),
			Reason:    reason,
			Fallback:  result.depth == 0 && result.move != "0000",
			ScoreKind: "unknown",
		}
		if result.depth > 0 {
			result.metrics.ScoreKind, result.metrics.Score = "cp", result.score
			if result.score >= mateScore {
				result.metrics.ScoreKind, result.metrics.Score = "mate", 1
			}
			if result.score <= -mateScore {
				result.metrics.ScoreKind, result.metrics.Score = "mate", -1
			}
			result.metrics.PV = []string{result.move}
		}
	}()
	if len(moves) == 0 {
		reason = "no_moves"
		return result
	}
	result.move = moves[rand.IntN(len(moves))].String()
	if depthLimit == 0 {
		depthLimit = maxSearchDepth
	}
	for depth := 1; depth <= depthLimit; depth++ {
		reason = "soft_limit"
		if worker.stopped() || (budget.timed && !worker.clock().Before(start.Add(budget.soft))) {
			break
		}
		reason = "interrupted"
		worker.nodes++
		best, score, complete := worker.root(position, moves, depth)
		if !complete {
			break
		}
		result = searchResult{move: best[rand.IntN(len(best))], depth: depth, score: score}
		reason = "depth"
		// The previous iteration's choice is a useful first candidate, without
		// changing scores or random selection among equally best moves.
		i := slices.IndexFunc(moves, func(move chess.Move) bool { return move.String() == result.move })
		moves[0], moves[i] = moves[i], moves[0]
	}
	return result
}

func (worker *searcher) root(position *chess.Position, moves []chess.Move, depth int) ([]string, int, bool) {
	var best []string
	bestScore := -searchInfinity
	for _, move := range moves {
		if worker.stopped() {
			return nil, 0, false
		}
		// The one-point margin keeps tied root scores exact. Inferior moves
		// may return bounds and cannot join the best-move list.
		score, complete := worker.search(position.Update(&move), depth-1, -searchInfinity, 1-bestScore, move.HasTag(chess.Check))
		if !complete {
			return nil, 0, false
		}
		score = -score
		if score > bestScore {
			bestScore, best = score, best[:0]
		}
		if score == bestScore {
			best = append(best, move.String())
		}
	}
	return best, bestScore, !worker.stopped()
}
