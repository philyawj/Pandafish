# Search-only legal-state kernel

The revised kernel retains immutable position copies and changes the cost of
state/tag calculation. Root policy, scores, legal move order, quiescence,
repetition, table behavior and time allocation remain released v0.19 policy.
The quiet-mate experiment remains archived and absent from the executable.

## Affected-piece transitions

`search_update.go` updates the moving/captured piece sets directly rather than
scanning and rewriting all twelve. Promotion, en passant and castling retain
their special state changes. Original convenience, castling-right, counter and
hash helpers remain in use. `Position.UpdateSearch` returns an immutable child;
it does not introduce make/unmake or an undo stack. Three engine child-update
sites and the authored generator's temporary legality board use this path.

## Needed attack/tag information

`search_attack.go` retains the original geometric rejection and attack semantics,
but computes each sliding ray once and selects piece sets directly. It avoids
repeated generic `NewPiece`/piece-set lookups and repeated ray calculations.

`AppendSearchMoves` still classifies opponent check for every returned legal move.
`HasLegalMove` only needs own-king legality and discards the candidate move, so
that path skips opponent-check classification. Required own-king checks, full
move tags and original castling validation remain intact. There is no attack
cache or speculative score, and no shared mutable board state.

## Boundary and controls

Original upstream runtime, public update/tag functions, module metadata and
license remain untouched. New additions have to stay synchronized with future
upstream changes; unchanged methods remain comparison oracles. The attack
predicate additionally agrees with an independent Go mailbox implementation
that uses no bitboards, ray helpers, move tags or attack masks.

Controls cover 128,000 attack predicates (all squares and both colors), 37,795
immutable child transitions, special moves, check/status/tags, mailbox/piece sets,
hashes, perft, parent ownership, repetition, mate distance and cancellation.
The full root race/vet/build checks and scoped upstream runtime race suite pass.
The unused upstream rendering dependency failure is preserved separately.

All 325 distinct position/depth comparisons retain scores, chosen moves and
nodes. Three alternating timing pairs give a depth-four median elapsed ratio
of 0.61256 versus v0.19, about 38.7% less time for identical work (about 63.2%
higher throughput). This is about 22.3% less time than the transition-only
revision's normalized median, exceeding the predeclared additional 10% gate.
Depth-three, depth-six and the 145-position development cohort also improve.

These are selected-workload engineering results. They do not establish playing
strength. Equal-time choices and the standard screen are separate acceptance
gates. The earlier transition-only screen failed despite its real efficiency
gain, so neither speed nor an isolated favorable game result is sufficient.
