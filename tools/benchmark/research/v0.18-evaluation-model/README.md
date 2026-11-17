## v0.19 tooling migration

This historical report retains its original measurements. Python/shell runners
were retired in v0.19; commands below describe the original experiment.
Current verification uses the [Go research tools](/tools/research/README.md).
Original scripts and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# v0.18 evaluation-model investigation

**Recommendation: revise the evaluation hypothesis; do not promote either
fitted-model variant.** Both fail the prespecified 100 ms selected-move gate,
although static held-out error and longer-budget decisions improve. No new
strength games, commit, push or PR were started for this experiment.

I started from [released v0.17](https://github.com/philyawj/Pandafish/releases/tag/v0.17)
after reassessing evaluation, further move-generation work, a board rewrite and
selective search in [ASSESSMENT.md](ASSESSMENT.md). The original material and
placement score stays intact; a compact, jointly fitted middle/end residual
adds safe pseudo-mobility, coordinated king-zone pressure/shelter and pawn
structure/advance/king-access features. The learned weights remain frozen.
No bishop-pair term, old piece-safety penalty, new pruning, clock or TT rule is
introduced. [ERRATA.md](ERRATA.md) preserves every failed prototype/check.

## Data and static result

The corpus contains 686 played/strong-reference-reply positions from the first
20 games per opponent in v0.17's screen. I preserved 639 distinct 50,000-node
fresh-FEN Stockfish queries and all score-bearing iterations. The teacher is a
finite reference, not ground truth or the sole acceptable move.

Opening families split six train/two validation/two holdout before fitting.
Eligibility excludes checked/draw-sensitive/large-value/unresolved-tactical
rows, with global board deduplication; every excluded row stays in the raw
corpus. [PLAN.json](PLAN.json) and the pre-fit refinement in
[CALIBRATION_PLAN.json](CALIBRATION_PLAN.json) record the exact procedure.
The 404 eligible rows split 230 / 83 / 91. Validation alone selects bounded,
integer-rounded joint weights; holdout is evaluated once after writing
[model-weights.json](model-weights.json). No holdout retuning occurred.

| Exact integer static metric | Baseline | Fitted model |
| --- | ---: | ---: |
| Training clipped MAE | 177.18 cp | 155.74 cp |
| Validation clipped MAE | 185.35 cp | 175.41 cp |
| Holdout clipped MAE | 219.38 cp | 203.95 cp |
| Holdout absolute errors above 300 cp | 28 | 27 |

The 7.04% holdout improvement passes the planned static gate. It does not
establish better search decisions or playing strength. There are only two
holdout opening families, source games are public/correlated, and the teacher
searches farther than the student. Mate labels remain separate from fitting.

## Correctness and cost

Independent Python-chess geometry matches all 360 original/color-mirrored
feature vectors after fixing the preserved initializer-order bug. All 686
active scalar/perspective cases match exact independent integer arithmetic.
Full repository/race/vet/native/ARM64/offline-wrapper/upstream tests pass.
Material values and base placement remain exact. Old material-approximation
assertions that no longer describe the positional model are replaced with
exact unpruned minimax, with tactical identities and material units retained.

The first fitted evaluator costs about 791 ns versus 578 ns for the released
base, with zero allocations. Its 2,046 native UCI requests complete correctly,
but it searches roughly 7% fewer nodes at equal time.

After its short-budget quality failure, I tested **one separately labeled
exact-cost revision** with the same frozen weights. A 38.4 KB immutable table
precomputes the original material/placement formulas for every phase/piece/
square. A value copy of existing piece bitboards avoids repeated mailbox/type/
formula work. All 19,200 entries, 596 old/new base evaluations and 686 fitted
scalar scores are exact. Fixed-depth nodes, scores and chosen moves match the
appropriate preserved slow variant. No values or coefficients change.

In that cost sample, base evaluation is about 48 ns and the fitted model about
258 ns, both allocation-free. These are component microbenchmarks, not total
engine speedups. The revised model completes all 2,046 native requests; a
separate fast-base-only ablation completes 336. Their source/provenance and
predeclared plan are in [fast-kernel](fast-kernel/PLAN.json).

## Actual selected-move quality

[DECISION_PLAN.json](DECISION_PLAN.json) was frozen before native/reference
results. At each budget, mean clamped finite-CP move loss must stay within
`max(10%, 5 cp)` of baseline, and severe positions may increase by at most two.
I score actual chosen moves, allow any good alternative, and retain mates
separately. All requested moves, exact iterations and failures stay in the
reference data. New revision queries reuse unchanged old queries byte-for-byte.

| Version / native comparison | 100 ms mean loss | 1,000 ms mean loss |
| --- | ---: | ---: |
| v0.17 baseline, first block | 55.57 cp | 68.00 cp |
| Slow fitted model | **64.10 cp — fails** | 56.33 cp |
| v0.17 baseline, revision block | 55.63 cp | 68.00 cp |
| Exact-cost fitted model | **66.58 cp — fails** | 51.73 cp |
| Exact-cost base only | 55.27 cp | 68.00 cp |

There are 67 finite-CP positions per budget; 17 mate-valued positions are kept
separately. The slow model changes severe short-budget positions from 1 to 3;
the faster model also has 3, while base-only has 1. At 1,000 ms, the faster
model has 2 versus baseline/base-only 3. These are public correlated diagnostic
estimates, not Elo or match-strength measurements.

Faster computation recovers node throughput, but **does not repair this model's
100 ms decision regression**. The base-only ablation has essentially unchanged
reference loss, so I cannot attribute its engineering gain to learned terms.
The long-budget improvement is encouraging evidence for a revised model or
objective; it does not override the failed standard-budget gate.

The two reference stages preserve all original 234 queries and reuse 232 in
the revision, adding only two genuinely new position/move queries. No weights,
thresholds or reference scores were changed to obtain a favorable outcome.

## Decision and next step

I stop this fitted-model candidate before the standard games. There are **no
new 50-game results or strength claim** for either model. Both sources, data,
failed gates and the exact-kernel/base-only controls remain preserved. The
released v0.17 tag and active bot/config remain unchanged; the current branch
is experimental and uncommitted.

[NEXT.md](NEXT.md) recommends separating the exact inference kernel as a
performance candidate, and revising evaluation data/objectives around actual
short-budget decisions on new independent families. Static prediction loss
alone is insufficient. I do not recommend guessing another weight, tuning to
these failed controls, or expanding into a board rewrite automatically.

The requested [local-library breakdown](../../../../docs/CHESS_SEARCH_BACKEND.md)
explains v0.17's original change, maintenance and future extension boundary.
PR #32 now explicitly describes v0.16 before / v0.17 after.

## Reproduction and archive audit

`python3 tools/benchmark/research/v0.18-evaluation-model/reproduce.py /tmp/new-eval-replay --variant fast`
restores the frozen fast source over tagged v0.17 in a new directory and runs
tests/build only. Use `--variant slow` for the initial fitted evaluator. Both
clean-overlay replays pass; [replay-checks.json](replay-checks.json) preserves
the exits and accompanying logs. No match or live connection is started.

[final-audit.json](final-audit.json) verifies the current 63 source/fixture
hashes, unchanged frozen weights, local document links and all 554 entries in
the six earlier archive inventories. [inventory.json](inventory.json) records
this investigation byte-for-byte, excluding itself and Python caches.
