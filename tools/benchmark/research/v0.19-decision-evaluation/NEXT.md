## Publication note for v0.19

The original local experiment report below retains its stop/check-in state.
These attempts were subsequently approved for publication together in v0.19.
Playing policy remains v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# Reassessment after the branch-objective fit

I keep the selected coefficients frozen and inactive. Holdout regret fails its
predeclared gate; the early-transposition audit further limits validation. I do
not relax the gate, change a weight or resplit the inspected data after seeing
it. This outcome rejects this fit, not every decision-oriented or neural model.

My next priorities are:

1. **Repair the data boundary before another learned fit.** Group and exclude
   canonical positions from the declared family boundary (including early move-
   order transpositions and mirrors), with synthetic positive/negative controls
   before collection. Keep all related games/branches together and use wholly
   fresh validation/holdout families. Increase independent coverage and verify
   the new data's branch signal before picking another representation. The
   inspected confirmation-related families are no longer an untouched model test.
2. **Reprofile tagged v0.18 legal tagging and attack/status work.** If I choose
   a separate playing-code project next, this is the first bounded candidate.
   Look for repeated check/legality information the additive backend can share,
   preserving exact moves/tags/state and fixed-depth behavior. The old profile
   supplies a reason to measure, not a guarantee of another large gain.
3. **Compare richer compact representations on useful data.** A small network
   or a more expressive classical model is plausible, but the current few
   hundred branches and incomplete scalar targets do not establish its design.
   Account for quiet threats/mates and frozen-leaf instability explicitly.
4. **Revisit bounded threat handling or verified selective search only with
   new mechanism evidence.** Keep mate/zugzwang controls and a strict compute
   budget. The old naive quiet-check extension and broad LMR bundle are not
   ready candidates to reactivate.

A full internal board rewrite remains lower priority: the existing additive
route delivered v0.17 and v0.18 without taking on rule/hash/undo replacement.
No opening book or Lichess dataset is a prerequisite for these projects.

The initial ranked research note remains the pre-experiment choice, and this
reassessment follows the result. I stop this bounded fit before activation,
runtime checks or standard games and will check in before a PR. The v0.18
release is complete; it is not held up by this unpromoted model investigation.
