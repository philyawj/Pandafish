# Phase-sensitive piece activity at the existing table cost

Choose the earlier coherent-piece-square direction after the bounded checking
screen fails at 50% direct / 50% Skill 3 / 77% Sacrifice. Keep all three earlier
v0.23 attempts inactive. This is neither another attack-map bonus nor a retry of
old fitted weights. Existing placement uses phase only for the king; all other
pieces retain the same preferences from opening through ending.

Use original integer middle/end formulas for pawn advance, minor centrality,
rook activity and queen development, interpolated in the existing immutable
25-phase piece-square table. Preserve material 100/300/300/500/900, the complete
king formula, phase definition, table dimensions and hot static-score loop.
There is no new per-leaf feature extraction, model inference or board scan.
No fitting, coefficient sweep, specific root move, castling bonus or opening
book is used. Coordinate-only activity does not recognize actual passed pawns,
pins, safe mobility or pawn shelter; do not claim otherwise.

Freeze this initial map before testing:

| Piece | Middlegame | Ending |
| --- | --- | --- |
| Pawn | Original central space, advanced flank increment four | Three times relative rank squared, plus twice central-file measure |
| Knight | Original center times eight | Center times six |
| Bishop | Original center times four plus development | Center times six plus two times capped relative rank |
| Rook | Central-file measure times four, original seventh-rank twenty | Center times three plus relative rank times two and seventh-rank twenty-four |
| Queen | Center times two minus twice relative ranks beyond the second | Center times four |
| King | Original formula exactly | Original formula exactly |

Center/file measures and integer rounding remain original where shared. Verify
color/perspective symmetry, each phase endpoint/interpolation, original king and
material invariants, immutable precomputed-table exactness, zero allocations
and unchanged inference loop/cost. Use broad placement/ending controls, not only
the inspected castling witness. Existing material-only expectations that remain
valid are retained; diagnose any failed test before changing its assertion.

Use the four unchanged actual-choice suites and tolerances. A clear failure
rejects this map before games, without changing its constants. A surviving
frozen map receives new final-control seed 2026112043, excluding the consumed
v0.23 quiet-check control boards in addition to previous exclusions, followed
by one standard concurrent 50-game-per-opponent screen. Retain every outcome,
do not promote from a static preference and stop before committing a viable
candidate. Search/quiescence/TT/draw/time policy remains released v0.22.
