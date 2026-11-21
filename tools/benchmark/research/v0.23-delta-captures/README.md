# Capture-futility opportunity: reject before pruning

The unchanged-tree diagnosis covers all 1,145 frozen warm roots. Shared
`warm-report` confirms every released move, score, depth and node count, with
zero issues and 43,118,858 total visited nodes. Outermost eligible captures cover
2,021,662 nodes, **4.689%**, below the prospectively frozen 20% opportunity gate.
There are 2,021,662 eligible captures and zero alpha-bound contradictions: each
eligible subtree already costs only one visited child. Baseline stand-pat cuts
already remove deeper continuation in these cases. The diagnosis's elapsed
time is instrumented and is not a measured performance comparison.

No playing capture-futility policy, actual-budget choices or games run. The
shared runner stops at `passed:false`; it does not rerun, lower the margin or
advance. Source, all rows, original test output, decision plan/state and exact
comparison remain here. The snapshot was captured after the process audit;
its additional shared optional harness was not part of the original diagnostic
execution and does not alter playing behavior. The executed diagnostic helper
and search overlay remain individually readable in `provenance`.

This is consistent with the older exact stand-pat relocation's limited node
opportunity, without relabeling that separate failure. It argues against another
immediate leaf-shortcut guess. Reassess broader selective search and the genuine
contextual-model project; do not reactivate the failed pawn/activity maps.
