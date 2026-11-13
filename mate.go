package main

import "github.com/corentings/chess/v2"

// Legal material/placement scores are far below this band. The 1,024-ply
// margin exceeds the 64-ply main search plus bounded quiescence/evasions.
const mateBound = mateScore - 1024

func searchedEvaluation(position *chess.Position, ply int) int {
	score := evaluate(position, position.Turn())
	if score == -mateScore {
		return score + ply
	}
	if score == mateScore {
		return score - ply
	}
	return score
}

// Stored mate scores are relative to the cached position, not the caller's
// root. The same position can be reached at a different ply in another search.
func scoreToTable(score, ply int) int {
	if score >= mateBound {
		return score + ply
	}
	if score <= -mateBound {
		return score - ply
	}
	return score
}
func scoreFromTable(score, ply int) int {
	if score >= mateBound {
		return score - ply
	}
	if score <= -mateBound {
		return score + ply
	}
	return score
}

// UCI uses signed moves, rather than plies, until mate.
func mateMoves(score int) int {
	distance := (mateScore - abs(score) + 1) / 2
	if score < 0 {
		return -distance
	}
	return distance
}
