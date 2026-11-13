package main

import "github.com/corentings/chess/v2"

// Only reversible history can recur. Capture immutable hashes before starting
// the worker, so protocol updates cannot race with history collection.
func repetitionHistory(game *chess.Game) []uint64 {
	positions := game.Positions()
	end := len(positions) - 1
	history := make([]uint64, 0, min(end, game.Position().HalfMoveClock()))
	for _, p := range positions[max(0, end-game.Position().HalfMoveClock()):end] {
		history = append(history, p.ZobristHash())
	}
	return history
}

// A count-set fingerprint distinguishes draw-sensitive cache contexts. Order
// does not affect threefold eligibility. Like the board hash it is probabilistic.
func repetitionToken(key uint64, count int) uint64 {
	if count == 0 {
		return 0
	}
	x := key + uint64(count)*0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}
func (w *searcher) pushKey(key uint64) bool {
	if w.repetitions == nil {
		w.repetitions = make(map[uint64]int)
	}
	n := w.repetitions[key]
	w.historyKey ^= repetitionToken(key, n) ^ repetitionToken(key, n+1)
	w.repetitions[key] = n + 1
	return n+1 >= 3
}
func (w *searcher) pushPosition(p *chess.Position) bool { return w.pushKey(p.ZobristHash()) }
func (w *searcher) popPosition(p *chess.Position) {
	key := p.ZobristHash()
	n := w.repetitions[key]
	w.historyKey ^= repetitionToken(key, n) ^ repetitionToken(key, n-1)
	if n == 1 {
		delete(w.repetitions, key)
	} else {
		w.repetitions[key] = n - 1
	}
}
