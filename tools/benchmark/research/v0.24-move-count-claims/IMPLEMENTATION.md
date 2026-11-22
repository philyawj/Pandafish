# Claim action and automatic termination

The new helper is consulted only after mate/stalemate handling. At100 halfmoves
it exposes a current claim. At99 it requires a legal quiet non-pawn move;
`AppendSearchQuietMoves` excludes captures, EP and all promotions. Castling can
qualify because it does not reset the halfmove clock. The rare99 check can
allocate; no claim cache or new backend legality code is introduced.

Main search starts with zero as an available action and keeps a better searched
continuation. Quiescence retains its finite heuristic boundary, applies the
claim floor and preserves complete check evasions. A positive bounded score is
still an estimate, not proof that a claimable game is won. At150 automatic zero
follows mate/stalemate checks in search/q; direct evaluation agrees. Root search
returns a legal compatibility move on an already automatically drawn board.

Root intended claims attach only to allowed legal moves whose child clock
reaches100. Existing claims may be taken before any move, including a resetting
one. Winning mates remain preferable, including a mating move that completes
100 or150 halfmoves. At an already reached150 position the game is over before
a future mate can be played. No material detector or neural model is included.

UCI supplies no draw-claim verb. The engine models claim availability and returns
a legal move; the game client is responsible for exercising a current claim.
The benchmark already claims when eligible. No guarantee is made that a generic
client which ignores claims will obtain the draw by executing a resetting move.
No client configuration or live behavior is changed by this research.

The released table already stores and checks exact halfmove counts. Warm tests
alternate the same board at0 and99 with the same supplied history, proving no
score reuse across clocks. Existing repetition identity, irreversible-history
reset and board-only hints remain intact. A full legal-action reference checks
focused main/root values; it shares the production q boundary and is not an
independent unlimited chess solver. Separate q tests check claim floors and
termination. Captures, EP, promotions, checked reset-only positions, searchmoves,
cancellation and both-turn mate precedence are explicit controls.
