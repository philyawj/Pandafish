# Diagnose completed quiescence reuse

The guarded late-quiet opportunity fails before implementation. Reassess a
different exact mechanism: a separate bounded cache for completed quiescence
scores/bounds. Current main-search scores and board hints persist within a game,
but quiescence has no score cache. Its warm cumulative CPU share is about 43%.
Do not assume main-cache hit rates establish quiescence reuse.

Instrument isolated tagged v0.21 without reusing any score. Simulate a separate
65,536-slot cache, exact remaining quiescence/evasion depth, checked/unfiltered
mode, halfmove clock, the existing board hash plus repetition-count fingerprint,
mate-distance normalization and ordinary alpha-beta bound applicability. Keep
terminal and repetition handling ahead of hypothetical hits. Store only complete
results with the existing generation/depth replacement preference. Retain a
private diagnostic cache per persistent game table.

Run the same first 20 saved games / 1,145 roots through production iterative
depth-five search. Count eligible bound hits and non-overlapping descendant-node
costs they could avoid; do not double-count nested spans or the hit node itself.
No elapsed/allocated-byte claim uses this instrumented profile. Require exact
completion/trajectory; report all identity hashes through the frozen workload.

If potential avoided descendants exceed 20% of visited nodes, freeze an additive
quiescence-cache policy and correctness/performance gates before implementation.
Otherwise reassess rather than relax history identity or remaining-depth matching
to manufacture hits. No depth>= shortcut, history-free score key, altered
quiescence eligibility, new evaluation or failed earlier kernel is included.
