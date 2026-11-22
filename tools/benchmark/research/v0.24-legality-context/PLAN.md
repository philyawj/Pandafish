# Frozen first diagnosis: shared legality context

Base playing source: tagged v0.23, `28948e7473251d950003ec1f966f0cf26369b817`.
Research branch starts after cleanup PR #42, merged at
`68ef54b08ef9b63d66332243d148ed2044c80e38`. No playing source changes are active.

The [ranked assessment](../../../../docs/NEXT_ENGINE_ASSESSMENT_V023.md) selected
shared legality/attack information ahead of a larger state rewrite or evaluator.
The current sample attributes 27.52% cumulative CPU to library generation,
11.79% to tags and 7.07% to status, with overlap. The mechanism considered here
can skip repeated own-king attack queries for ordinary unpinned nonking moves
when the king is not checked. Special moves, checked nodes, malformed king
configurations and pinned origins retain the released predicate. Opponent-check
classification remains unchanged.

## Bounded measurement

Use a scratch replay of v0.23, adding diagnostic hooks only. Every actual search
decision continues to use the original tag/status/generation result. Count
contexts, tag categories, own/opponent attack-query opportunities and status
queries followed by generation. A cached shadow context is keyed by the actual
Position cell and board hash; no global context enters normal engine builds.
Compare every proposed skipped own-check result with the released predicate.
Independent vectors and the existing two-transition-oracle controls remain
separate checks. No full new legal generator or playing policy is implemented.

One existing warm depth-five workload: all 20 games /10 openings /1,145 roots,
persistent per-game tables and expected 43,118,858 nodes. The shared warm-report
must find identical moves/scores/nodes to the saved accepted candidate. Diagnostic
elapsed time is not speed evidence. A bounded cost sample on the frozen warm
root FENs separately compares context construction with the own-check work it
could remove; this is an inference input, not a warm speedup or strength result.

Retain one compact counter report and any contradictory FEN/move witnesses,
source overlay, commands/hashes, normal/race control results and the shared
exactness report. Full warm rows are temporary comparison material, referenced
by hash; the original exact rows already remain in Git. No new teacher data,
strength matches, sweep or live play is part of this diagnosis.

## Prospective continuation and stop conditions

Zero legality/tag contradictions are required. Root coverage and exact whole-tree
results must pass. Context setup and residual opponent/special work must be
accounted for; cumulative profile shares must not be added. Continue toward a
playing implementation only if measured category/cost evidence plausibly supports
at least 8% net whole-work saving. Counts alone cannot pass that gate. If the
bounded measurements cannot establish it, retain the result as insufficient
headroom evidence and reassess; do not silently relax the threshold or claim a
candidate gain. Preserve all failed or inconclusive controls.

Publication of any playing candidate is not authorized by cleanup publication.
The initial request starts this project; later strength/release claims require
their normal independent gates and separate publication decision.
