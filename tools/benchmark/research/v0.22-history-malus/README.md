# Rejected signed quiet-history revision

The [signed policy](../v0.22-next-assessment/SELECTION_MALUS.md) rewards completed
quiet cutoffs and penalizes earlier quiet moves. Every move stays eligible;
evaluation, quiescence, TT/repetition policy and clocks remain v0.21. Three exact
Go overlays and hashes preserve this version separately from positive-only
history and the four earlier exact-kernel attempts.

All controls and full race/vet/build pass. All 325 cold fixed-depth scores/root
moves remain identical. Public median elapsed ratio is 0.97066; deep/newer ratios
0.92206/0.88613. The original cold 15% field remains false and its driver exits 1.
The [prospective workload correction](../v0.22-next-assessment/HISTORY_METHOD_CORRECTION.md)
was frozen before these measurements and requires a 15% median gain in three
production-style warm iterative pairs, with no cold regression at the specified
depth-four/deep/newer gates. The additional depth-three diagnostic is 0.8% slower
and remains preserved.

That corrected gate also fails. Each warm pair visits all 1,145 roots from the
same 20 recorded games. Scores and root choices match baseline everywhere; each
actor repeats its own results/nodes exactly across all pairs. Candidate nodes
fall 9,874,565 to 9,318,451 (**5.63%**), but elapsed ratios are 1.00205, 1.00880 and
1.00011: median **0.20% slower**. Allocation and parse costs remain in the raw
reports. `warm-analysis.json` records `efficiency_gate_passed: false`; the driver
exits 1. No pair is selected or replaced. Its generic membership-pending field
does not establish an executed independent proof; all actual selected moves and
scores are identical here and no separate proof driver runs after failure.

The reduced tree does not pay for sorting/feedback overhead in this workload.
This also shows that correcting the benchmark lifetime does not rescue the
ordering hypothesis. Archive without promotion. No native quality suite or
strength games run, so no strength claim follows. Next assess full-depth
aspiration windows against iteration-score stability before implementation.
No failed history/kernel code remains in that next candidate.
