## Publication note for v0.19

I preserve the original local experiment report below, including its stop/check-in
state. I subsequently received approval to publish these attempts together in
v0.19. I keep playing policy at v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# Current handoff: repaired boundary, unpromoted model

I completed the requested top recommendation after v0.18: repair early family
independence, collect fresh larger data and repeat the fixed objective/gates.
All work is local and uncommitted on `v0.19-decision-evaluation`. The prior
failed batch remains byte-preserved in `../v0.19-decision-evaluation/`.
The user accepted archiving these bounded failures with a later approved PR,
without creating a separate release for every experiment. Check in before PR.

v0.18 remains published through PR #33, tagged at
`75f12bbdbeb0e38d148690ddc0f1d6405c5af3d3`. No new release/strength screen,
live Lichess client, private config change or larger confirmation was run.
The working playing source is restored exactly to v0.18 after this candidate's
actual-choice gates fail; restored normal tests pass. Candidate-source.json plus
reproduce.py fully restore the model in a NEW directory; clean replay tests and
build pass. Candidate binary remains ignored under tools/benchmark/bin.

Read README.md, PLAN.json (effective precollection plan), PLAN-pre-design.json,
ENGINE_PLAN.json, ERRATA.md and NEXT.md. The pre-native mate-guard clarification
has its original ENGINE_PLAN-pre-mate.json retained. No gate/weight retuning.

Boundary repair: canonical placement/turn/castling/legalEP, clocks ignored,
mirrors grouped; include all paths fromply4 onward. Freeze40,058 prior/reserved
states. Six boundary tests pass before generation, two logger regressions later
verify terminal nulls/live null rejection. Generate80 fresh rollouts in86attempts,
417sample roots. Before labels/splits, trace all14,167 legal alternatives and
independently verify paths/FEN plus13,864 static vectors/scalars. Group all
rollout AND witness paths into76 components. Independent audit before and after
labels: zero blocked or cross-split states. Raw splits246/82/89.

Selected teacher/baseline alternatives reuse witnesses. After unchanged
mate/nonstatic/root exclusions:223train/74validation/81holdout,567trainingpairs.
Frozen ridge1 weights improve holdout surrogate61.12→56.02cp (8.35%), severe4/4.
Scalar13,495 uniqueFEN/bothperspective fixtures, normal/race/vet/native/ARM64,
offline wrapper, original upstream race, 192 completed/legal fixed probes and
GHI pass. Inference257.7ns/op,0alloc, within400ns. Two material-approximation
assertions become exact independent minimax; capture/material/EP gains retained,
original failures/contracts preserved in ERRATA.md.

All2,758 native requests valid:2,730nonterminal+28expectednulls, no timing overrun.
Depth decreases about0.07–0.12plies. Actual finite-reference quality:
public100ms55.27→68.04cp FAIL; public1000ms68→66.12pass;
fresh100ms78.5→63.24pass; fresh1000ms68.89→54.66 finiteCPpass, but new mate
survival guard FAILS. Reference477queries,206unchanged cached; all raw scores and
mate cases retained. No standard50game screen follows because overallgatefails.

The fresh flag is already about-1044cp lost, candidateRg1 refmate-4 vsbaselinec5
finite-1141cp. This is not proof baseline can save it or model turned equality
into a loss. Public100msmean independently fails. Rg1 was labelled in corpus,
but the entire holdout root was excluded from CP surrogate for its mate label;
the same filter excludes22training mate-containing roots and7validation roots.
Do not fit this holdout: the concrete next idea is retain sharp/mate outcome
constraints in a new learning plan. Saved same-depth public failures show cost
alone is not the explanation. Failure-diagnosis.json contains all evidence.

Current recommendation: retain the corrected pipeline, revise the model with
sharp/mate branch supervision and broader training-only negatives, fresh holdout.
Representation/runtime and v0.18 legal-tag/check-summary profiling are alternatives,
not an automatic checklist. No guessed weights, board rewrite, book or Lichess
dataset is needed to proceed. Final audit/inventory preserve all earlier results
and the current source/candidate distinction.

Coverage-diagnosis.json records another limit: all fresh sample plies are even,
so all417 raw fresh roots are White to move. Role alternation is not root-turn
balance. Scalar/color symmetry checks do not replace Black decision validation.
Plan odd/even or same-component mirrored decisions in the next dataset; do not
resplit/refit this inspected batch. Disjoint components are not an IID claim.
