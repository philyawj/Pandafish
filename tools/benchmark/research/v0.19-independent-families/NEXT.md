## Publication note for v0.19

I preserve the original local experiment report below, including its stop/check-in
state. I subsequently received approval to publish these attempts together in
v0.19. I keep playing policy at v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# Reassessment after the repaired-boundary experiment

**I keep the boundary repair and revise the fitted evaluator.** Cleanly grouped
fresh data produces an 8.35% surrogate improvement and 19–21% better fresh
finite-reference actual-move loss. Public short-budget regressions and one
fresh long-budget mate-survival flag prevent promotion. No strength game ran.
These are separate conclusions: the pipeline repair works; the model has not
passed its whole acceptance plan or established playing strength.

My next priorities are:

1. **Retain sharp and mate-containing alternatives in the learning/validation
   design.** The flagged holdout move was already labelled mate -4, but its
   whole root was excluded from the CP surrogate. Twenty-two training roots
   containing mate labels were excluded too. A new experiment should separate
   finite-CP fitting from conservative outcome/order constraints, retaining
   dangerous alternatives rather than dropping an entire position. The current
   holdout stays frozen and inspected; use a fresh independent confirmation of
   a newly selected model. Do not fit this specific holdout failure.
2. **Broaden the tested alternative set on training data.** Three reference
   choices plus the old student's move is a useful start, not full action
   coverage. The all-legal forest now makes additional negative sampling
   possible. Predeclare random/low-value and training-model-preferred alternatives
   before selection. Include both odd/even sampled plies or same-component mirrored
   decisions: the current fresh roots are all White to move. Keep mate and
   terminal classes visible. This is a proposed
   mechanism, not a tested fix for the public regressions.
3. **Assess representation on those training-only counterexamples.** Fixed-depth
   value errors persist, so cheaper inference alone is not the next answer.
   The pair matrix has full rank and exact feature collisions impose only a
   tiny optimistic regret floor; neither fact proves linear adequacy. Compare
   context-aware classical features or a compact learned board representation
   only with explicit data and runtime budgets. No neural model is assumed to
   win automatically, and no old weight is hand-tuned against these controls.
4. **Reprofile the v0.18 backend if selecting a separate playing-code project.**
   Shared legal-tag/check/attack information remains the leading bounded
   engineering alternative. Measure it before implementation; keep exact rules,
   hashes, search and evaluation when testing that hypothesis. A full board
   rewrite or another broad pruning bundle still has no stronger measured case.

The exclusion policy intentionally makes fresh families different from the
previous common-opening controls, and the larger opening pool/data size changed
alongside the repair. I cannot attribute every gain/regression to one cause.
A future dataset should explicitly distinguish broad representative training,
known public regressions and truly untouched validation/holdout roles. Any old
families incorporated into training lose their status as unseen evidence.

I do not resplit, refit or weaken the failed gate to obtain a positive result.
The working playing source is restored to released v0.18; the exact experimental
source remains replayable. No new PR, commit, live game, larger confirmation or
next candidate follows automatically. I will check in before a PR.
