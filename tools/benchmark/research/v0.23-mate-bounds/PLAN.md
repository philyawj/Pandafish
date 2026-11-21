# Use legal mate-distance bounds to stop impossible improvements

Capture-stage ordering preserves all warm values/choices but saves only 5.944%
nodes, so archive it and restore v0.22. Further conservative ordering layers
have limited measured headroom. Test an exact search bound already implied by
the released mate-distance convention, rather than another evaluation fit.

After terminal/repetition handling at ply p, no score lies below -mateScore+p
or above mateScore-p-1. Intersect the alpha/beta window with those limits and
stop when it closes. The original terminal evaluator, repetition identity,
mate/table normalization, search eligibility and clocks remain unchanged.
At the root, a proven mate in one is unbeatable; return the first proven move
without searching other roots. Never score a hypothetical threat as mate.

Use positive mates for both colors, losing/delayed mates, stalemate/repetition,
ordinary centipawn windows, cross-ply table bounds, parent/history restoration
and cancellation controls. Require all 1,145 warm scores/first choices exact,
deterministic candidate counts, >=15% node saving after each pair and >=15%
median search elapsed saving over three alternating pairs. Preserve failures
and stop early if the node opportunity is insufficient.

If promising: four existing native suites, new final seed 2026112117/stream117
with consumed/live exclusions, then the standard concurrent 150-game screen.
Any surrogate exception needs the documented independent matched-depth review
and separate pre-game decision. No automatic larger confirmation/publication;
check in before committing a viable candidate.
