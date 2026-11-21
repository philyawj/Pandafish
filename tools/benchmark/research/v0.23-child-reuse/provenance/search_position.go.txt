package main

import "github.com/corentings/chess/v2"

// A child stays valid through its recursive calls, including deferred history
// cleanup and a possible PVS re-search. Descendants use larger ply indices.
// Pointer cells keep active ancestors stable if the outer slice grows.
func (worker *searcher) childPosition(parent *chess.Position, move *chess.Move, ply int) *chess.Position {
	for len(worker.positionBuffers) <= ply {
		worker.positionBuffers = append(worker.positionBuffers, nil)
	}
	if worker.positionBuffers[ply] == nil {
		worker.positionBuffers[ply] = new(chess.SearchPosition)
	}
	return worker.positionBuffers[ply].Update(parent, move)
}
