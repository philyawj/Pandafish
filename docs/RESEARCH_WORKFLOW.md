# Engine improvement research workflow

“Find the next improvement” requests use this workflow unless the request states
a different scope. Begin from the latest released engine on a research branch.
Review recent results alongside earlier failed attempts and rank opportunities
again. Prefer a substantive, evidence-supported opportunity; the previous ranking
is not a checklist and a familiar feature is not automatically a new hypothesis.
Treat live games as one evidence source alongside controlled screens, current
profiles and earlier research directions; live-derived ideas do not automatically
outrank those alternatives.

Follow [the tooling map and flowchart](RESEARCH_TOOLING.md). Reuse the shared Go
plans, runner, summaries, snapshots and inventories instead of creating duplicate
per-attempt drivers. Read the current compact checkpoint before raw archives.
Reserve model reasoning for selection, diagnosis and review; do not use model
agents for polling or arithmetic. Keep all scientific gates and adverse evidence.

1. State one hypothesis, its mechanism, correctness boundary, cost/quality
   criteria and reasons for choosing it before measuring its outcome. Diagnose
   saved mistakes before assigning a cause. Separate development cases from
   untouched final controls; accept comparably good moves rather than requiring
   a particular reference move.
2. Implement and validate correctness, realistic computational cost and actual
   equal-time decisions. Use Go for all source, including research drivers.
   Preserve raw measurements, source identities, controls, adverse outcomes and
   uncertainty. Functional improvements are distinct from measured strength.
3. Archive an unpromising implementation and its findings without promoting it.
   Restore the released playing code before a different hypothesis. Keep every
   attempt on the research branch, reassess the ranking and continue without
   requesting approval merely because an attempt failed. Do not rerun unchanged
   candidates to obtain a preferred score or relabel failed gates as passed.
4. A promising candidate receives the standard 50-game-per-opponent screen
   against the tagged previous release, Stockfish Skill 3 and pinned Sacrifice,
   with all three matchups concurrent. Preserve every game and failure. Do not
   automatically start a larger confirmation; recommend it only when justified
   by the reviewed result. Follow the [benchmark guide](../tools/benchmark/README.md)
   for exact profiles and the previous-release criterion.
   For an exact-value efficiency change with a failed decision surrogate,
   follow the tooling guide's matched-depth review and explicit pre-game
   exception procedure before deciding whether one research screen is justified.
   The original gate remains false and must be disclosed in the final review.
5. When a candidate merits promotion, present the implementation, tests, equal-
   time evidence, complete screen, rejected predecessors and material limits.
   Stop before committing, pushing, opening a PR, merging or releasing unless
   those actions have separate explicit authorization. Do not start live games.

This authorizes continued local research within the requested investigation,
not an unlimited publication or deployment sequence. Reassess after every
meaningful result. If progress requires unavailable information or resources,
state the concrete blocker and completed evidence rather than inventing a
promotion recommendation. Publication follows [RELEASING.md](../RELEASING.md).

A sufficient future request is: **“Find the next improvement using the standard
research workflow.”** A target such as “for v0.24” may be included. The checkpoint
remains a recommended candidate before commit unless the request changes it.
