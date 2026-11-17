## Publication note for v0.19

The original local experiment report below retains its stop/check-in state.
These attempts were subsequently approved for publication together in v0.19.
Playing policy remains v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# Sharp-branch evaluation retry after v0.18

**I reject this fitted candidate for promotion.** I retained sharp/mate branches,
balanced both root turns, expanded negative alternatives and corrected an offline
target-scaling error before fitting. Validation improved, but the fresh holdout
failed: mean decision penalty rose **13.51%**, severe cases rose **3 to 6**, and
one new reference-mate-risk choice appeared. I never activated these weights in
the playing engine. Released v0.18 remains unchanged; no strength matches ran.

I preserve this as the third local post-v0.18 attempt, after the
[decision fit](../v0.19-decision-evaluation/README.md) and
[independent-family fit](../v0.19-independent-families/README.md).
The `v0.19-` directory names describe research toward a possible release, not
published versions. I have not committed, pushed, prepared a new PR or released
these attempts. I will include the failed findings in the next authorized PR;
I do not need a standalone release for every local experiment.

## My hypothesis and frozen design

The preceding fit excluded whole mate-containing roots from CP supervision.
It passed its frozen-leaf surrogate but failed actual short-budget controls and
one mate-survival guard. Its fresh sample also contained only White-to-move
roots. I tested whether retaining those difficult alternatives and balancing
turns would let the same compact model generalize safely.

I froze [PLAN.json](PLAN.json) before collection. I retained the same corrected
12 board features, 24 nonnegative bounded middle/end coefficients, original
material/base scores and regularization grid. I added a separate soft ordering
constraint between a reported losing mate and an alternative without that
report; I did not translate mate distances into numerical CP targets. Exact
terminal/repetition witness values remained fixed. Winning mates and multiple
close finite alternatives remained acceptable; I did not prescribe root moves.

Before the first fit or validation/holdout quality calculation, I found that
`clip(teacher_gap, +/-400) - base_gap` could demand the wrong correction
direction for large already-correct material margins. I froze the algebraic
repair, `clip(teacher_gap - base_gap, +/-400)`, in
[FIT_REFINEMENT.json](FIT_REFINEMENT.json), retained the original plan/source,
and added zero-residual and sign controls. This was an experimental offline
objective error; released v0.18 contains no fitted residual. I have not rewritten
earlier results or established this correction's isolated empirical effect.

I describe the mate conventions and their limits in
[METRIC_NOTES.md](METRIC_NOTES.md). The primary **capped decision penalty** is a
selection diagnostic, not a true CP loss or playing-strength estimate. A finite
reference score may still represent a lost position. “New mate risk” means a
selected losing-mate report when another labelled alternative lacks that report,
not proof that the original position was salvageable.

## My data and controls

I generated 80 families in 142 attempts, sampling adjacent even/odd plies from
offline tagged-v0.18/Stockfish rollouts. These are data-collection queries, not
strength games or live Lichess games. I blocked 140,130 canonical states from
16 prior/reserved inputs, including inspected query and witness paths.

Before reference branch labels, I traced every legal root alternative through
unchanged depth-three student search. The forest contains 29,357 records;
28,890 learnable vectors/scalars pass the independent python-chess oracle.
Whole-family transitive grouping yielded 78 components, retaining all 80
families. Both independent boundary audits reconstruct paths and report zero
blocked or cross-split positions. Mirrors and counter differences cannot create
new folds; the boundary starts at ply four.

| Coverage | Training | Validation | Holdout |
| --- | ---: | ---: | ---: |
| Raw roots, all retained | 504 | 186 | 186 |
| Roots with at least two alternatives | 494 | 182 | 184 |
| White / Black decision roots | 248 / 246 | 92 / 90 | 92 / 92 |
| Decision roots with a finite reference | 488 | 176 | 181 |
| Roots with losing-mate and non-losing-mate-reported alternatives | 37 | 9 | 11 |

Each labelled set contains the unique union of reference MultiPV3, the baseline
100ms move, its worst learnable depth-three witness and one seeded remaining
legal move, at most six choices. I retained terminal/single-choice controls and
all mate-containing roots. The selected alternatives are still a subset of all
legal actions; their labels use a finite 50,000-node reference budget.

The six boundary controls, three query controls and twelve model controls pass.
Positive/negative cases cover early transpositions, mirrors, legal en passant,
castling rights, component collisions, terminal null moves versus invalid live
null moves, legal timed fallbacks, close acceptable CP choices, winning mates,
unavoidable losses, unsafe ties, immutable terminal values, solver scale and
target direction. I retain the initial environment/import and logging failures
in [ERRATA.md](ERRATA.md).

The rollout logger initially stopped on a legal 25ms baseline fallback with no
completed scored depth. I repaired the logger, replayed its 1,392 recorded
queries without contacting an engine and appended only unfinished work. The
prefix bytes/hash and 22 completed families remain identical. I did not reroll
an unfavorable result, resplit families or extend a failed sample.

## My fit and holdout result

I fitted 4,484 informative finite-CP pairs and 44 active ordinal pairs, retaining
counts for constant/already-satisfied constraints. All five grid fits converged.
Validation selected regularization 1; I froze integer coefficients and source,
plan/refinement, corpus and membership hashes before opening holdout quality.
The zero model remained eligible. Full coefficients, all grid results and every
decision remain in [model-weights.json](model-weights.json) and
[fit-results.json](fit-results.json).

| Diagnostic | Unchanged baseline | Selected fit |
| --- | ---: | ---: |
| Training mean penalty | 51.1503 | 46.2874 |
| Validation mean penalty | 36.0952 | 33.7198 |
| Holdout mean penalty | 44.0435 | 49.9918 |
| Holdout severe roots, penalty >300 | 3 | 6 |
| Holdout reference-mate-risk roots | 3 | 4 |
| Holdout White mean penalty | 65.8098 | 71.3533 |
| Holdout Black mean penalty | 22.2772 | 28.6304 |
| Holdout CP-only mean regret, 159 fully finite roots | 32.1006 | 32.5692 |

The holdout required at least 5% improvement, no new risky root, at most one
additional severe root, and each turn within its declared tolerance. The fit
fails all but the White-turn tolerance. The new risky root is
`sharp-028-ply36`. I did not change its weights, bounds, split or gate afterward.

I stopped before Go activation, inference profiling, native equal-time tests and
the standard concurrent 50-game screen. Consequently I have **no measured
strength result**, candidate runtime-cost result or equal-time performance claim
for this fit. There is no new engine binary/source overlay to promote. Source
identity with tagged v0.18, rather than a new search test, establishes unchanged
playing behavior in this working tree.

## What I learned after stopping

I used post-result probes for diagnosis, without fitting another model or
querying a reference again. The original exploratory capacity output remains
saved. [The conservative rational-bound analysis](capacity-diagnosis-robust.json)
allows a full 2cp for integer truncation. Within this exact coefficient box,
15 training mate-pair orderings cannot be repaired; four dominating bad branches
span **two distinct training roots**. A finite/mate pair also has identical base
score and feature vector. Another alternative could still resolve that root.
These facts concern frozen witnesses and these bounds. They do not establish
that every linear evaluator, larger weights, neural evaluation or actual search
must fail. No strict dominance proof survives removal of the coefficient caps.

The initial six-branch probe of the two bounded-capacity roots found no immediate
mate. I then scanned **all 2,557 selected learnable training witnesses**, without
using validation/holdout to choose examples. I found **15 static-witness records
with a legal quiet mate in one**, representing 14 distinct leaf FENs, 13 roots
and 12 families. Python-chess verifies the moves. An untouched tagged-v0.18
scratch tree independently verifies every legal mating move and reproduces a
finite full-window quiescence score equal to stand pat in all 15 records:
[proof log](quiet-proof.txt), [probe data](training-quiet-threat-scan.json).

For example, the saved leaf after `sharp-009-ply37` / `a5b6` admits `e5h8#`;
released quiescence instead returns its finite static score of -133 for the
side to move. This proves a local capture-only horizon omission. It does not
prove a forced root loss, improved root choice, frequency in real games or a
future strength gain: the witness line itself depends on the old search.

The same scan finds 54 hypothetical opponent mates after a pass. Those are
threat diagnostics only. Passing is not a legal chess continuation, and a legal
quiet defense may exist; I cannot return a forced-loss score from that probe.

I therefore retain the corrected fitting/data tools but reject this selected
model. I now rank a **bounded exact quiet-mate proof experiment** ahead of another
24-weight retry. The older naive recursive quiet-check extension timed out;
that result remains a constraint. I set out the mechanism, cost/negative controls,
alternatives and stop rules in [NEXT.md](NEXT.md). I have not implemented it here.

## My preservation and replay boundary

I keep raw generation/label query iterations, all rollout attempts, the all-legal
forest and selected witnesses, both boundary audits, feature oracle/source,
initial failures and resume manifest, pre-correction objective, frozen fit and
post-result probes. [HANDOFF.md](HANDOFF.md) gives reproduction commands.
[final-audit.json](final-audit.json) independently checks decisions, frozen
hashes, source identity and all previous inventories. The inventory hashes all
current archive files except itself and disposable Python caches.

I did not change material, evaluation, search, quiescence, tables, clocks or
dependency versions in the playing tree. The diagnostic Go files run only in
scratch trees. I have not launched live games, run a larger confirmation,
committed, pushed or prepared a PR.
