## Publication note for v0.19

I preserve the original local experiment report below, including its stop/check-in
state. I subsequently received approval to publish these attempts together in
v0.19. I keep playing policy at v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# Decision-aligned evaluation revision after v0.18

**I reject this frozen fit for promotion and stop before activating it.** It
improves training/validation branch regret but fails the fresh holdout gate.
The stricter family audit also identifies early-transposition limitations in
the collection design. Released v0.18 remains the entire playing policy; no
candidate native-performance test, strength game, commit, push or PR follows.

I completed the prior promotion through [PR #33](https://github.com/philyawj/Pandafish/pull/33)
and [v0.18](https://github.com/philyawj/Pandafish/releases/tag/v0.18), at merge
`75f12bbdbeb0e38d148690ddc0f1d6405c5af3d3`. Its original 57% / 73% / 51%
screen and all failed-model archives remain intact. [RESEARCH.md](RESEARCH.md)
is the short ranked note I wrote before starting this top choice.

## What I revised

The earlier model selected lower absolute static error yet worsened actual
short-budget decisions. Here I keep its corrected 12-feature geometry but
replace its data and selection objective. I compare alternative root branches
at exact quiet-leaf witnesses of the released depth-three search. I fit
within-root value differences, allow a 25 cp finite-reference uncertainty band,
weight large disagreements and select rounded coefficients on validation regret
with a severe-error guard. The unchanged zero model is a valid selection.
Material values, base placement, search, quiescence, tables and clocks remain
unchanged. [PLAN.json](PLAN.json) was frozen before collection/labels/fitting.

## Corpus and independent controls

I generated 40 new move-prefix families in 52 attempts, using seeded offline
legal exploration and alternating tagged v0.18 / finite-node full-strength
Stockfish rollouts. These are training rollouts, not strength matches or a
production opening book. I sampled plies 12/20/28/36/48/64. The declared union
of shared boards from ply 8 onward yields 40 components; the frozen raw split
is 129 training / 39 validation / 43 holdout roots, 211 total.

At each playable root, alternatives combine three finite-reference candidates
and the unrestricted student's actual 100 ms move. Each restricted alternative
gets its own fresh-hash 50,000-node reference query, with all score-bearing
iterations/bounds and its last exact result retained. I force no production
root move and accept any comparably good alternative. The reference is a
finite estimate, not chess truth.

[queries.jsonl](queries.jsonl) preserves all 3,418 query records, including the
terminal sampler error, generation and student streams. [generation.json](generation.json),
[corpus-membership.json](corpus-membership.json) and [corpus.json](corpus.json)
preserve every attempted family, full rollout history, split and branch label.
One sampled terminal root has no alternatives and stays in the raw corpus.

The scratch-only Go tracer reconstructs a scalar witness through the unchanged
main/quiescence search, retaining full repetition history, SEE eligibility,
check evasions and safety caps. Controls cover three distinct start alternatives
and a history-sensitive repetition. All 693 reconstructed paths/FENs are legal
and match independent Python-chess replay. All 676 static-leaf feature vectors
match independent geometry, and their root-perspective original scalar scores
match the exact branch values. Seventeen terminal witnesses remain separate.
[extract-checks.json](extract-checks.json) and [leaf-oracle.json](leaf-oracle.json)
retain these results; no trace/feature code enters the released engine.

## Frozen fit and failed gate

Excluding entire roots with mate/nonstatic/incomplete alternative sets leaves
114 training, 34 validation and 38 holdout roots. I retain all 25 excluded raw
roots. The training set contributes 302 informative branch pairs. A bounded
24-coefficient middle/end ridge fit uses training-only RMS scaling and the
predeclared regularization grid. Validation selects regularization 0.01 and
writes [model-weights.json](model-weights.json) before the first holdout metric.
No holdout retuning or favorable rerun occurred.

| Frozen-leaf mean decision regret | Unchanged model | Selected fit |
| --- | ---: | ---: |
| Training, 114 roots | 34.99 cp | 29.89 cp |
| Validation, 34 roots | 60.26 cp | 50.79 cp |
| Holdout, 38 roots | 54.02 cp | 54.74 cp |

Severe (>300 cp) counts change 1→2 in training, 2→1 in validation and stay 1→1
in holdout. The holdout mean worsens **1.32%**, failing the required ≥5%
improvement. [fit-results.json](fit-results.json) retains the full grid and
all selected decisions, not only the chosen weights. This is a surrogate test;
actual model search and playing strength were never measured.

## Independence audit and preserved errors

The original collector incorrectly required a bestmove from a sampled checkmate.
I retain its terminal query, source and traceback. [resume-plan.json](resume-plan.json)
freezes completed IDs and membership hash; I resumed only unfinished labels,
without regenerating data or repeating completed queries. [ERRATA.md](ERRATA.md)
explains the correction.

The later [stricter audit](independence-audit.json) starts canonical grouping at
ply 4 instead of the declared ply 8. It finds two early positions shared by
training/validation families 16/19/21, and three early positions from two new
families related to reserved confirmation openings. These are **not identical
sampled roots**, but they weaken family-level separation and the claim that
prefix rejection fully protects reserved families. The held-out failure is
retained without repartitioning/refitting. Validation's apparent improvement
is not clean generalization evidence. [independence-controls.json](independence-controls.json)
checks an explicit early transposition, a distinct negative and color mirror.

Future model tests need fresh families and an untouched holdout with canonical
transposition exclusion/grouping from the family boundary, before labels.
Related confirmation families have now been inspected in this research; a
future model should use a freshly separated confirmation set. This does not
alter v0.18's earlier frozen source, screen or release decision.

## What I learned and recommend

The revised objective found a validation improvement but not a held-out one.
I cannot attribute failure solely to one coefficient, data size, the static
representation or the frozen-leaf approximation. A leaf chosen by the old
search can change when weights change, and only a small alternative subset is
labelled. Mate-labelled cases also stay outside the CP regression. A more
powerful model is not automatically a remedy for those limitations.

[NEXT.md](NEXT.md) reassesses the priorities. I would repair and validate data
independence before another learned fit, using new data rather than optimizing
this inspected holdout. If choosing a separate playing-code project now, fresh
v0.18 profiling of legal tagging/attack-summary reuse has a clearer engineering
boundary than a full board rewrite or an immediate NNUE jump. I do not continue
through the old ranking automatically.

I keep this entire revision local on `v0.19-decision-evaluation`. Production Go
source is byte-identical to tagged v0.18. The [final audit](final-audit.json)
checks source/provenance, all 783 earlier inventory entries and data consistency.
No live game, private config change, larger confirmation or next PR was started.
