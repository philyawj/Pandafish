# Align supervision with the represented witness board

The ninth attempt improves training but worsens both inspected development
splits and adds a reported mate-risk choice. Its coefficients and all failures
stay frozen. Do not tune its ridge, cap or buckets against those outcomes.

The earlier fitted-model reports already warn that frozen student witnesses do
not represent the strong opponent's actual continuation. The last model labels
a later witness board using the strong reference's value of a root alternative.
First diagnose whether those targets differ substantially; do not assign this
as the proven cause of past failures without measurement.

Before querying, select eight White and eight Black training roots by SHA256 of
`v023-direct-leaf-audit/` plus root ID, with distinct saved family components.
Require at least two learnable witnesses per selected root. Include every
selected root's labeled alternative, not only earlier regressions or preferred
moves. Reconstruct the complete original history plus witness path and verify
the exact final FEN. Query static witnesses unrestricted at pinned full-strength
Stockfish, 50,000 nodes, Threads 1, Hash 64MiB, MultiPV1, three workers. Fixed
terminal/repetition witnesses retain their known outcome rather than a fitted
score. Preserve every UCI query/error and both label sets.

An alignment signal requires at least 32 finite comparisons and either median
absolute root-versus-leaf CP difference >=75, at least 10% opposite preferences
among pairs whose gaps are >25cp under both labels, or a changed losing-mate
classification. This detects target mismatch, not its isolated effect or a
perfect teacher. Finite-node horizon differences also contribute.

Only if that gate passes, extend the append-only labeling collection to all
remaining learnable witnesses. Reuse completed identical queries with matching
teacher/source/settings identities, never select favorable re-queries. Retain
original splits/family grouping. Existing boards/splits are development data;
new labels do not make them an unseen strength holdout.

Keep the same 450-dimensional king-context representation, cap ±32, group-zero
constraints, phase taper, ridge 4, band25, Huber200, residual clamp400, ordinal
margin25, quantization and repaired solver. Change supervision to the actual
witness position's root-perspective value, preserving fixed outcome scores.
Fit training only, once. Require >=5% lower new-label decision penalty on both
development splits, no additional severe count or new avoidable-loss-reported
root. Also retain original-root-label metrics as diagnostics; they are not an
equivalent value target. Do not refit to either result.

A surviving fit gets scalar/runtime/correctness and <=15% warm-overhead gates,
the original four native gates, explicit actual-search controls for all three
new mate-risk roots from the prior model, new seed 2026112103/stream103 final
families and the standard concurrent 50 games per opponent. Keep material/base
placement, search/quiescence, TT policy and clocks v0.22. No automatic larger
confirmation or publication. Check in before commit at a recommended candidate.
