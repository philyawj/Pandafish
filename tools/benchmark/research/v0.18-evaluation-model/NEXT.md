# Reassessment after the fitted-model gates

**Revise; do not promote the fitted evaluator or run its standard screen.**
The predeclared short-budget gate fails twice with the same frozen weights.
The exact-cost revision rules out the initial inference slowdown as a complete
explanation: it is cheaper than the old base evaluator, yet the short-budget
selected-move loss remains worse. Static holdout improvement is therefore an
insufficient selection objective for this model.

The model improves the 1,000 ms reference diagnostic. That supports continued
evaluation research, with limited confidence from public/correlated inputs;
it does not establish match strength or justify changing the routine time
control until the model looks favorable. No larger or longer match ran.

## Best next bounded choices

1. **Separate the exact base-evaluation kernel.** Its immutable table and
   piece-bitboard traversal preserve all old scores and tested fixed-depth
   nodes/moves, and materially reduce component cost. The base-only diagnostic
   is decision-neutral on the finite reference set. Give this its own clearly
   labeled runtime/load checks and standard 50-game-per-opponent screen before
   considering promotion. Do not ship the rejected fitted terms merely because
   this kernel is useful. It is a performance candidate, not a measured strength
   gain; that strength screen has not run.
2. **Revise the evaluation learning objective/data.** Keep new game/opening
   families independent of the control set already inspected here. Train and
   validate on response branches and student-search decision consequences at
   the intended budget, alongside quiet static labels. Allow multiple sound
   moves, retain tactical/endgame/mate controls, and account explicitly for
   compute cost. The existing label corpus has few independent families and
   only 404 eligible rows; it does not justify a general classical-versus-NNUE
   conclusion. A compact learned representation remains an alternative after
   a new data/runtime plan, rather than a drop-in network.

I would not immediately refit the regularization, hand-adjust bishop/rook
mobility or weaken the gate using these outcomes. That would turn the inspected
controls into tuning data. The current weights and holdout results are frozen.
A new model needs a new selection/validation plan.

Further legal-tag optimization and a full internal board rewrite remain lower
priority without a stronger measured case. The exact kernel already exposes
how the local library boundary can support efficient evaluation summaries.
No new feature, live game, larger confirmation, commit, push or PR follows
this report automatically. I will check in before the next PR.
