# Reassessment: rule out impossible king-attack changes before simulation

The [child-storage attempt](../v0.22-child-storage/README.md) fails its frozen
timing gate (3.3% median reduction; 4.65% on the warm sequence). Its five-file Go
overlay is archived, and source is restored to v0.21 before this separate attempt.
Warm byte allocation falls 79.3%, illustrating that allocation bytes alone do not
predict useful elapsed gains. No native or strength screen is spent on it.

Remaining move generation/tagging accounts for about 28–32% cumulative CPU in
the warm profiles; tagging itself about 20–23%. It simulates the entire temporary
board and checks both kings for every pseudo-legal move. Most quiet moves cannot
change either king's attack status. Select a conservative geometry precheck,
retaining the existing board simulation/attack helpers whenever uncertain.

Once per uncached generation, cache full geometric rays from each king and
potential checking destinations by moving piece type. Away from check, moving
a non-king whose origin is outside its king's rank/file/diagonals cannot expose
that king. No own-king simulation is needed in that case. If the opponent king
was not already attacked, a move cannot check when its origin is outside the
opponent's rays and its destination is outside the new piece's geometric attack
origins. Promotions use their promoted piece type. EP and castling remain
conservative simulation cases; checked parents and already-attacked opponent
kings retain the required tests. The context computes the latter condition once
rather than assuming every FEN is a reachable legal position. Missing kings keep
the existing guard semantics. Whenever either test remains possible, simulate
using the unchanged update and test exactly those required king predicates.

Initialize geometry from integer square coordinates only, avoiding dependencies
on the original package's later init functions. No pin cache, attack-map cache,
move omission, guessed priority, new state representation or pruning is added.
The old searchMoveTags and original library methods remain independent oracles.
Immutable child allocation remains baseline v0.21.

Require direct tag equivalence on every pseudo-legal destination (including
illegal self-checks) in seeded positions, every square/promotion in targeted
special/invalid-FEN controls where meaningful, both existence/full modes,
parent immutability and positive skip/conservative fallback controls. Keep all
existing independent move/state/perft, transition and race controls. The same
15% median timing, 325 exact-tree, persistent-TT sequence, four equal-time native
and one standard concurrent 150-game screen gates apply before recommendation.
