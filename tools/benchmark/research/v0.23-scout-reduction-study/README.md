# Child-scout opportunity passes

The unchanged 1,145-root workload matches every released move, score, depth and
all 43,118,858 nodes. Nonoverlapping outer eligible scouts cover 28,646,303
nodes: **66.436%**, with 66.832% White and 65.957% Black. Wide-window parents
account for 63.791% of all nodes in these outer spans. There are 792,211 eligible
scout calls, including nested calls; the cost estimate does not double-count
nested spans or count a following original full-window re-search as saved work.

The frozen 30% overall/20% each-turn gates pass. These are opportunity upper
bounds, not measured savings: reduced work, verification and changed ordering
can erase the benefit. Instrumented elapsed is not performance evidence.

The earlier narrow-parent guarded-LMR diagnosis remains failed under its original
guards and cohorts. This measurement concerns the actual late child's scout
window, including wide-window parents, and explicitly accepts a different
selective risk. It does not reclassify the old observation. Preserve source,
raw per-root records, aggregate turn/category counts and all verification.
Restore uninstrumented v0.22 before implementing the separate candidate.
