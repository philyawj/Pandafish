package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/corentings/chess/v2"
)

const maxSearchDepth = 64
const defaultMoveTime = 100 * time.Millisecond

type goLimits struct {
	moveTime         time.Duration
	clocks           [2]time.Duration
	incs             [2]time.Duration
	hasClock         [2]bool
	hasTime          bool
	depth            int
	movesToGo        int64
	infinite, ponder bool
	searchMoves      []string // nil means unrestricted; empty means no moves.
}

func parseGo(args []string) (goLimits, error) {
	limits := goLimits{}
	for i := 0; i < len(args); i++ {
		key := args[i]
		switch key {
		case "infinite":
			limits.infinite = true
		case "ponder":
			limits.ponder = true
		case "searchmoves":
			limits.searchMoves = append([]string{}, args[i+1:]...)
			return limits, nil
		case "movetime", "wtime", "btime", "winc", "binc", "movestogo", "depth":
			if i+1 == len(args) {
				return limits, fmt.Errorf("missing value for %s", key)
			}
			i++
			n, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil || n < 0 || n > int64((1<<63-1)/time.Millisecond) {
				return limits, fmt.Errorf("invalid value for %s", key)
			}
			switch key {
			case "depth":
				if n < 1 || n > maxSearchDepth {
					return limits, fmt.Errorf("depth must be 1..%d", maxSearchDepth)
				}
				limits.depth = int(n)
			case "movestogo":
				if n < 1 {
					return limits, fmt.Errorf("movestogo must be positive")
				}
				limits.movesToGo = n
			case "movetime":
				limits.moveTime, limits.hasTime = time.Duration(n)*time.Millisecond, true
			case "wtime", "btime":
				side := 0
				if key == "btime" {
					side = 1
				}
				limits.clocks[side], limits.hasClock[side] = time.Duration(n)*time.Millisecond, true
			case "winc", "binc":
				side := 0
				if key == "binc" {
					side = 1
				}
				limits.incs[side] = time.Duration(n) * time.Millisecond
			}
		}
	}
	return limits, nil
}

type timeBudget struct {
	soft, hard time.Duration
	timed      bool
}

// Explicit move time takes precedence over clock allocation, but is capped by
// the usable remaining clock when supplied. Depth is an additional limit.
func (limits goLimits) budget(turn chess.Color) (timeBudget, error) {
	if limits.infinite || limits.ponder {
		return timeBudget{}, nil
	}
	side := 0
	if turn == chess.Black {
		side = 1
	}
	allowance := defaultMoveTime
	if limits.hasTime {
		allowance = limits.moveTime
	}
	if limits.hasClock[side] {
		clock := limits.clocks[side]
		usable := clock - min(10*time.Millisecond, clock/10)
		if limits.hasTime {
			allowance = min(allowance, usable)
		} else {
			horizon := limits.movesToGo
			if horizon == 0 {
				horizon = 30
			}
			// Cap to half the usable clock before adding increment, avoiding
			// overflow and never relying on increment not yet credited.
			cap := usable / 2
			share := min(usable/time.Duration(horizon), cap)
			allowance = share + min(limits.incs[side]/4*3, cap-share)
		}
	} else if (limits.hasClock[0] || limits.hasClock[1]) && !limits.hasTime {
		return timeBudget{}, fmt.Errorf("missing clock for side to move")
	} else if !limits.hasTime && limits.depth > 0 {
		return timeBudget{}, nil
	}
	// Leave a small margin inside the client's allowance. The runner's
	// protocol grace is not available to the engine.
	hard := allowance - min(5*time.Millisecond, allowance/10)
	return timeBudget{soft: hard / 4 * 3, hard: hard, timed: true}, nil
}
