package chess

// SearchPosition is caller-owned scratch for one active search position. Update
// invalidates its previous result. Use a distinct cell for every active ancestor
// and worker; do not retain a result in a game tree or share it concurrently.
// Do not copy a cell after use. The immutable Position.Update and UpdateSearch
// APIs remain available when the child must outlive the next scratch update.
type SearchPosition struct {
	position Position
	board    Board
}

// Update copies parent state into the cell and applies a generated legal move.
// It has UpdateSearch's state/hash semantics, including its nil transition.
// Aliasing the parent with the destination panics before changing either state.
func (dst *SearchPosition) Update(parent *Position, move *Move) *Position {
	if parent == &dst.position || parent.board == &dst.board {
		panic("chess: search child storage aliases parent")
	}
	count := parent.moveCount
	if parent.turn == Black {
		count++
	}
	dst.board = *parent.board
	if move == nil {
		dst.position = Position{board: &dst.board, turn: parent.turn.Other(), castleRights: parent.castleRights, enPassantSquare: NoSquare, halfMoveClock: parent.halfMoveClock + 1, moveCount: count}
		dst.position.hash = dst.position.computeHash()
		return &dst.position
	}
	rights := parent.updateCastleRights(move)
	halfMove := parent.halfMoveClock + 1
	if parent.board.Piece(move.s1).Type() == Pawn || move.HasTag(Capture) {
		halfMove = 0
	}
	dst.board.updateSearch(move)
	dst.position = Position{board: &dst.board, turn: parent.turn.Other(), castleRights: rights, enPassantSquare: parent.updateEnPassantSquare(move), halfMoveClock: halfMove, moveCount: count, inCheck: move.HasTag(Check)}
	dst.position.hash = parent.updateHash(move, rights, dst.position.enPassantSquare)
	return &dst.position
}
