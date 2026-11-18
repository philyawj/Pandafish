# History workload correction before signed-revision measurements

Freeze this correction before preparing/running the signed revision's timing
driver. The positive-only version's original failed gate stays a failed result;
no native or match observations exist for it and no unchanged favorable rerun is
authorized. This correction applies prospectively to the separate signed policy.

The exact-kernel benchmark directly enters one depth with a fresh searcher for
each root. That measures identical work appropriately for byte-copy/sort/clock
kernels. Quiet history's intended lifetime spans **iterative passes within a go**;
the direct-root benchmark omits those passes. It remains a useful correctness and
cost control, but does not alone measure the new mechanism's production lifetime.
Using it as the sole 15% admissibility target would confuse cold single-depth
performance with the hypothesized useful ordering improvement.

Keep every direct-root result, including the original `timing_gate_passed` field.
Require all 325 scores and tagged-baseline selected-best membership, and no
aggregate elapsed regression on public depth-four median, depth-six or newer
depth-four. Expected node/order differences remain reported. Then measure three
alternating baseline/candidate pairs on the frozen 1,145-root / 20-game sequence,
using the production iterative-deepening entry point, per-request fresh history
and one persistent TT per game at depth four. No movetime deadline is reached;
the same root-depth/trajectory/history/budget applies to both actors.

The revised efficiency gate is **at least 15% median search elapsed reduction on
those three warm iterative pairs**, equal complete scores and independent tagged-
baseline membership for every selected move. This keeps the larger-gain target
but measures the actual lifetime under test. Do not choose one favorable pair,
replace games or lower the threshold after observing warm data. Parsing/setup
and MemStats overhead remain separately reported. All four actual-budget native
gates and the standard concurrent 150-game screen remain mandatory afterward.

This methodological reassessment is part of the research record. It establishes
neither a successful candidate nor playing strength. If it fails, preserve the
code/data and reassess another hypothesis without automatic broader confirmation.
