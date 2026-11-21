# Static-margin diagnosis: reject the unverified shortcut

The [frozen plan](PLAN.md) tests an untouched v0.22 tree, not an active pruning
candidate. All1,145 moves/scores/depths and43,118,858 nodes match the released
reference. Qualifying outer descendant spans cover20.808% of work, with21.715%
White and19.712% Black. The opportunity thresholds pass, but the mandatory
hidden-mate safeguard fails. Instrumented elapsed is not performance evidence.

Across1,077,982 qualifications,179 original node results disprove the proposed
finite fail-high, including one forced losing mate. Every false prediction's
FEN, remaining depth, ply, original null window, proposed score and original
result is preserved in `study-summary.json`. The largest score error is998,694.
Static material advantage alone is insufficient even with the150cp/ply margin,
check/mate-window guards and phase>=8. Do not implement this unverified cutoff,
tune its margin or proceed to native tests/games. No fresh cohort is consumed.

This is the22nd unsuccessful substantive attempt:18 engine prototypes, one
offline model and three diagnoses. Only three full game screens have run; failed
early gates are not22 strength comparisons. Successful opportunity studies are
counted separately. Source, raw rows, workflow failure and exactness report stay
preserved. Restore released play before another hypothesis.

A bounded follow-up may test actual quiescence verification rather than a
hypothetical pass or another static threshold. It must charge verification work,
preserve repetition context and show that the same hidden mate is not concealed.
That is a new prospective diagnosis; this failed gate remains false.
