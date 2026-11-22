# Full pin-context diagnosis: insufficient headroom

All six execution steps pass. This means the diagnostic runs correctly; it does
not mean the prospective implementation-opportunity gate passes.

The original 1,145 warm roots /43,118,858 nodes remain exact with zero changed
moves/scores. All 596 independent rule vectors pass. The two-transition-oracle
control compares 140,778 tags, with zero shadow contradictions, normally and
under race. No playing source is changed.

The compact warm record counts 126,057,881 tags and 42,791,870 constructed shadow
contexts. Only 51,047,094 ordinary own-king queries are eligible: 40.62% of
ordinary own queries, or 26.02% of all ordinary own/opponent attack queries.
Many existence checks search king moves, which this shortcut cannot skip.
15,558,661 of43,113,133 status calls are followed by generation (36.09%).

The root-distribution cost control measures62ns/context and about16ns/attack
query. The [automated assessment](assessment.json) combines these with the
earlier profile, explicitly treating counts and microcosts as a rough proxy.
Even an optimistic gross model credits only4.98% saving before context costs;
the tested full-context setup is materially more expensive. It does not support
the frozen8% net objective. These are not rigorous ceilings or whole-engine
speed measurements, and may double-credit overlapping status/tag work.

**Decision:** do not activate this full per-position pin-context/own-check
shortcut. Counts do not justify weakening the gate or promoting safe but costly
code. This rejects the measured construction/eligibility strategy, not every
shared attack representation. Instrumented elapsed ratio1.3634 is not evidence
of an engine regression: the diagnostic computes original results plus shadow
contexts, counters and assertions. No native/teacher/game stage runs.

The next bounded reassessment selects exact target-centric attack queries:
remove captured/moved squares from attacker results rather than copying and
editing six piece sets before every query. It uses existing pawn-origin masks
and virtual occupancy, with no eager per-node ray/pin context. This differs
materially from the rejected own-query-skipping study and leaves evaluation,
order, search horizon, TT/repetition and clocks unchanged. The alternative of
starting a new evaluator has a larger data/inference boundary; a broad state
rewrite has weaker current headroom evidence. No sweep of pin variants follows.

## Reproduction and lean evidence

Use the shared `research replay --base v0.23 --out NEW_SCRATCH`, copy
`prepare.go.txt` to a temporary `.go` path, and run it with the scratch/evidence
directories. Its five-file Go overlay is also directly replayable with
`research replay --source diagnostic-source.json --out NEW_SCRATCH`.
Build the root test executable and adapt scratch/output paths in the frozen
[workflow plan](workflow-plan.json), then use the shared runner.
`report.go.txt` computes the compact final proxy; it does not run engines.

The runner records source/input/output hashes, errors, logs and gates.
Temporary full warm rows are retained locally at the recorded `/tmp` path and
identified in `exactness.json`; their original exact moves/scores/nodes already
exist in the accepted v0.23 candidate rows. Only counters, adverse witnesses,
cost results and summaries are added here, avoiding another tracked full warm
trajectory. Go1.25/M2/8 CPUs/AC; background and thermal state unmeasured.
No other deliberate heavy work overlaps the warm/cost stages.

This directory remains uncommitted on the research branch. Cleanup PR #42 has
already merged; no next engine publication is authorized by that merge.
