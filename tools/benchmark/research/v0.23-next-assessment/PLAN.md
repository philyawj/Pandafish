# Prospective v0.23 assessment

Baseline playing source is tagged v0.22, `1953f2af97f956c4499cbc5b61c3c1c698856740`;
branch starts after documentation merge `8b0109d5fd578a397e2f7ae9857190bd94df8095`.
No commit, publication or live play is authorized. Use the
[standard research workflow](../../../../docs/RESEARCH_WORKFLOW.md).

The latest live review moves king exposure/development diagnosis ahead of the
earlier checking-interface proposal. Prior v0.18/v0.19 fits already included
linear mobility, king-pressure and pawn features, and failed actual-budget or
holdout gates. Another weight tweak or reactivation is not a fresh candidate.

First reconstruct six independently confirmed development decisions with full
history and fresh v0.22 tables. Compare completed-depth choices at live depth
and the next depth, keeping timeouts as failures, then inspect finite-reference
continuations and static values. This does not claim exact live warm-state
reproduction. Existing live confirmations remain development evidence.

The leading implementation direction is nonlinear king exposure conditioned on
actual attacking resources and legal board context, rather than a constant
castling reward. Choose its exact scope only after diagnosis. A smaller search
omission discovered by that diagnosis can change the ranking. Preserve original
material, clocks and cache policy unless a separately documented hypothesis
explicitly requires otherwise.

Require focused positive/negative controls, color symmetry, material/terminal
invariants, no hot-path allocations, realistic inference cost and equal-time
decision evidence on existing and fresh families before games. Freeze any
selection/coefficients before final controls. Never use a selected teacher root
as the sole acceptable choice. A failed candidate remains archived and causes
another reassessment; it does not cause an automatic larger screen or stop.

Only a surviving candidate receives one standard concurrent 50-game screen
against v0.22, Skill 3 and pinned Sacrifice. Report all results and return for
approval before committing a candidate that merits promotion.
