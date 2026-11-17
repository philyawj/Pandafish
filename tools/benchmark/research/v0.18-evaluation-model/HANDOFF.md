# Current evaluation experiment

I published v0.17 via PR #32 and the v0.17 release at tagged commit
`76e73c8716d2ddcc826bdf349ff1c7a7ac71ecf3`. The next experiment is local on
`v0.18-evaluation-model`, with no commit/push/PR. No live games or larger
confirmation are authorized. The active bot/config are untouched.

I also clarified PR #32's opening paragraph: v0.16 is the before-state,
v0.17 introduces the local search API replacement. The requested detailed
note is `docs/CHESS_SEARCH_BACKEND.md`, linked from the main README. It remains
local alongside this experiment and should be retained through the next review.

The chosen hypothesis is a compact jointly calibrated tapered classical
residual over the unchanged material/placement score. The local chess package
adds allocation-free board features: 12 advantage-oriented inputs, with middle/
end weights. No bishop-pair, old piece-safety penalty or search/clock/table
change is bundled. Read ASSESSMENT.md, PLAN.json, CALIBRATION_PLAN.json,
DECISION_PLAN.json and ERRATA.md for frozen methods and preserved corrections.

The offline corpus has 686 played/strong-reference-reply rows and 639 distinct
fresh-FEN Stockfish queries at 50,000 nodes, one thread/Hash64/Skill20. All score
iterations are preserved. Opening families split six train/two validation/two
holdout; eligibility includes quiescent stability. 404 rows fit: 230 train,
83 validation, 91 holdout. Weights are frozen in model-weights.json. Exact
integer holdout MAE improves 219.38→203.95 cp (~7.04%), with >300 cp residuals
28→27. This is static label quality, not strength. No holdout retuning.

Independent feature checks initially found zero pawn masks from initialization
before upstream bbSquares init. Failed vectors/code are retained; direct-bit
initialization and stronger negative controls fix it, with all 360 independent
original/mirrored vectors now matching. Six hundred eighty-six active scalar/
perspective cases match independent integer arithmetic. Full/race/vet/build/
ARM64/offline-wrapper/upstream tests pass. Several old material-approximation
assertions were updated to exact independent minimax, retaining exact material
and tactical outcomes; all original failures remain. See ERRATA.md.

Fitted evaluator cost is about 791 ns versus 578 ns baseline, zero allocations.
This passes the <=1.75x static-cost gate. Search mechanics remain v0.17; score,
chosen move and pruning node counts may change with the evaluator. Cold fixed
probes complete and GHI controls pass. The scratch-copy harness initially
failed because v0.17 already contains third_party; it was resumed without
replacing the completed baseline probes.

The original slow-model UCI schedule completed all 2,046 requests. Its
prespecified 100ms move-quality gate FAILED (baseline 55.57 cp loss versus
64.10 candidate), while at 1000ms the candidate improved 68.0→56.33 cp. The
full raw source/probes/234 reference queries are preserved; no games ran for
that candidate. No weights were refit.

I completed one separately labeled exact-cost revision in
`fast-kernel/`. It precomputes all 19,200 original material/placement entries
and iterates value-copied piece bitboards. Tests prove original base scores on
596 boards and fitted scalar scores on all 686 corpus rows unchanged. Fitted
static cost falls from about 791ns to 258ns; original v0.17 base was about
578ns. Its PLAN.json was written before new native probes. Current source UCI
is `v0.18-evaluation-model-fast-dev`; old source is slow-model-source.json.

The exact-cost revision and ablation are complete: all checks pass, 2046
fast-model native requests and 336 fast-base requests are valid. Its 100ms
selected-move gate ALSO FAILS: baseline 55.63 cp loss, fitted fast model 66.58,
fast base 55.27. At 1000ms the model improves baseline/base 68.0 to 51.73.
The revision reuses 232 old reference queries unchanged and adds two new ones.
No weights/thresholds were changed, no 50-game matches ran for either model.

The final recommendation is REVISE the fitted evaluator, no promotion. Keep
its static/long-budget positive findings and both short-budget failures.
NEXT.md separates the exact base kernel as a potential performance candidate
that still needs its own standard game screen; do not attribute its gain to
learned terms. Revise future model data/objectives on new independent families
with student-budget decision controls, not these inspected gates.

The current working source remains the unpromoted fast fitted model for review.
Both slow/fast overlays and all data are preserved. No new PR, commit, push,
live game or larger confirmation was started. The backend note requested by
the user is complete and linked from README; PR #32's before/after wording is
updated. The final audit and clean-overlay replays are recorded in final-audit.json
and replay-checks.json. The released v0.17 baseline remains unchanged.
