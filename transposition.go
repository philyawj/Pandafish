package main

import "github.com/corentings/chess/v2"

const tableEntries = 1 << 16 // Fixed allocation per engine process, about 2 MiB.

type bound uint8

const (
	noBound bound = iota
	exactBound
	lowerBound
	upperBound
)

type tableEntry struct {
	key        uint64
	halfmove   int
	score      int32
	depth      int16
	move       uint16
	generation uint16
	bound      bound
}

type transpositionTable struct {
	entries    []tableEntry
	generation uint16
}

func newTable() *transpositionTable {
	return &transpositionTable{entries: make([]tableEntry, tableEntries)}
}

func (table *transpositionTable) clear() {
	clear(table.entries)
	table.generation = 0
}

func (table *transpositionTable) probe(position *chess.Position, history ...uint64) (tableEntry, bool) {
	if table == nil {
		return tableEntry{}, false
	}
	key := position.ZobristHash()
	if len(history) > 0 {
		key ^= history[0]
	} // Pieces, turn, castling, capturable en passant.
	entry := table.entries[key%uint64(len(table.entries))]
	return entry, entry.bound != noBound && entry.key == key && entry.halfmove == position.HalfMoveClock()
}

func (table *transpositionTable) store(position *chess.Position, depth, score, alpha, beta int, move uint16, history ...uint64) {
	if table == nil {
		return
	}
	key := position.ZobristHash()
	if len(history) > 0 {
		key ^= history[0]
	}
	entry := &table.entries[key%uint64(len(table.entries))]
	// Prefer deeper work within a search; allow replacement in a later search.
	if entry.bound != noBound && entry.generation == table.generation && int(entry.depth) > depth {
		return
	}
	kind := exactBound
	if score <= alpha {
		kind = upperBound
	} else if score >= beta {
		kind = lowerBound
	}
	*entry = tableEntry{key: key, halfmove: position.HalfMoveClock(), depth: int16(depth), score: int32(score), move: move, generation: table.generation, bound: kind}
}

func encodeMove(move chess.Move) uint16 {
	return 1 + uint16(move.S1()) + uint16(move.S2())<<6 + uint16(move.Promo())<<12
}

// Match the hint against generated legal moves rather than trusting cached tags.
func preferMove(moves []chess.Move, hint uint16) {
	for i, move := range moves {
		if encodeMove(move) == hint {
			copy(moves[1:i+1], moves[:i])
			moves[0] = move
			return
		}
	}
}
