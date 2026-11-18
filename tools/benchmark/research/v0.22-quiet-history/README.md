# Rejected positive-only quiet history

The [diagnosis and policy](../v0.22-next-assessment/SELECTION_HISTORY.md) select
per-search side/from/to history after nonfirst quiet cutoffs account for about
14% of main cutoffs. A bounded depth-squared positive reward orders the quiet
suffix before existing killer/TT precedence, without changing eligibility,
quiescence, values, table policy or clocks. Three Go overlays over v0.21 and
their hashes preserve the tested code. No prior failed kernel is active.

Bounds, eligibility, side/request isolation, stable order, cache/ply ownership,
killer/TT precedence and cancellation controls pass, as do full race/vet/build.
All 325 distinct completed fixed-depth scores agree; every selected root move
also happens to remain identical to baseline. Nodes are expected to change.

The attempt fails the unchanged 15% elapsed gate. Three public depth-four ratios
are 0.93941, 0.96603 and 0.97137; median **3.40% less time**, with 1,010,100 to
950,694 nodes (5.88% fewer). Depth-six uses 3,033,434 versus 3,658,933 nodes and
0.90471 elapsed ratio; newer depth-four uses 2,957,882 versus 3,354,194 nodes and
0.89980 ratio. These are useful functional reductions, but not the required
larger gain. The timing driver exits 1 and retains its failed gate. No warm
candidate workload, selected-move proof driver, native decision suites or
strength screen runs after it. Prepared helpers are source, not executed checks.

Archive without promotion. Positive rewards alone provide limited discrimination
between successful and repeatedly unsuccessful quiet moves. That is a hypothesis
for the next separate revision, not a proven causal diagnosis: add bounded
negative updates to earlier quiet moves after a completed quiet cutoff. Retain
all earlier results, apply the same gates and do not rerun unchanged code for a
favorable sample. Root choices remain unrestricted in every production search.
