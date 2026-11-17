# Reassessment after the exact evaluation revision

**I recommend promoting the isolated exact-score kernel under the existing
practical gate.** It retains v0.17's values and fixed-depth behavior, removes
measured work and passes the standard screen. This recommendation does not
promote the fitted evaluator archived beside it.

The direct 57% score passes the practical screen, with zero forfeits or
unfinished games. Its wide nominal interval and opening concentration limit
the inference; the 73% Skill 2 sample is encouraging and Sacrifice's 51% is
near even. I can establish the functional speedup more firmly than the size
or generality of a strength improvement.

A broader-opening confirmation is justified before claiming a large strength
gain or advancing Skill 2. The existing optional 100 games per opponent on
50 opening positions would address part of this uncertainty. I have not run
it and do not add it as a new mandatory release gate after seeing results.

## The next larger project

I prioritize **decision-aligned evaluation data and model selection** next.
The original static fit improved held-out label error but failed short-budget
choices twice; reducing inference cost ruled out slowdown as a complete
explanation. More absolute-score samples with the same selection objective,
or another guessed term, would leave that failure mechanism unresolved.

[LEARNING_REVISION.md](LEARNING_REVISION.md) sets out a new approach: independent
families, response branches and quiet student leaves, multiple acceptable
moves, a branch/ranking objective with severe-error controls, fixed material,
explicit inference cost and one fresh holdout. A separate numeric collection
and selection plan must be frozen before generating or fitting that model.
The exact kernel offers a cheaper base without requiring a new board model.

A compact neural representation remains a candidate within that project, not
an automatic next implementation. The small classical corpus does not decide
classical versus neural evaluation. I would first establish a useful dataset,
decision objective and runtime budget, then select representation from evidence.
No production opening book or Lichess dataset is a prerequisite; offline legal
exploration and self-play can supply data while preserving independent controls.

The existing move-generation route already met its target; a full board rewrite
still adds a large rule/state validation burden without a better measured case.
Further pruning changes the search hypothesis and remains separate. I reassess
these choices after each result rather than committing to all of them in order.

## Check-in boundary

The current source and complete research results are local and uncommitted.
No PR, push, release, live client, next model run or larger confirmation follows
this screen automatically. I will check in before a PR. The initial failed
models, their exact overlays, labels, gates and every error remain intact.

## Subsequent authorization

I approved promotion as v0.18 and authorized PR/merge/release, then a fresh
ranked short research note and work on its top choice. That supersedes the
completed candidate's local-only checkpoint above. The next candidate remains
subject to a check-in before its PR; no larger confirmation or live play is
authorized by this promotion.
