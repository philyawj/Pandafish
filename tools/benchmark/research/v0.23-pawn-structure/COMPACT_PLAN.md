# Exact pawn-kernel compaction before decision measurements

The first two unprofiled pairs show 12.7–14.4% overhead. Before any native
decision or game observation, preserve that implementation and all three pairs.
Revise only its exact computation: derive occupied files and isolation with
bitboard fills/populations, and visit the sparse passed-pawn set rather than
every rank. Keep every coefficient, phase rounding, doubled-pawn definition,
immediate-blocker discount and material/search/cache/clock policy unchanged.

Require all scalar oracle comparisons and a separate frozen 1,145-root check
against the initial pawn binary to retain scores/moves/depths/nodes. Preserve
each source snapshot and component timing. Repeat three alternating warm pairs
against the release in a separate `compact-cost` directory under the original
15% slowdown limit. This revision cannot rescue or rerun a failed native/game
observation: neither has run. Continue only with the frozen compact binary.
