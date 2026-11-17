# What I would change in the learned-evaluation approach

I keep the original model frozen and unpromoted. Its 7.04% static holdout
improvement is real within that finite corpus, but it did not meet the
short-budget decision gate. Making inference cheaper did not repair that
failure. The current exact-score candidate separates the reusable computation
improvement from the learned values.

## What the saved decisions show

The [descriptive comparison](prior-model-diagnosis.json) uses the already
inspected fast-model and fast-base results, not new validation data. At 100 ms,
the model improves reference loss on 16 finite-CP positions, worsens it on 19
and ties on 32. The three largest regressions contribute 14.31 cp to the
67-position mean, exceeding the overall 11.31 cp deterioration versus fast
base. Gains elsewhere partly offset those misses. At 1,000 ms, it improves
21, worsens 12 and ties 34, with a better aggregate result.

For example, in the saved Ruy Lopez position after `...Bd6` and White castles,
the model repeatedly chooses `...Nd4`, while exact base castles. The finite
reference estimates a 423 cp loss for the former and zero for the latter. Its
saved response includes `Nxd4 exd4 e5 Bxe5 Re1`. In the saved English position,
`Qc3` versus `Qb3` remains a roughly 294 cp reference-loss difference at both
budgets. I retain the full positions and actual moves in the JSON. These are
examples of branch-consequence failures, not instructions to hard-code a move.

The reference is a finite-node search. Root and move-restricted searches can
finish different depths and give inconsistent estimates; even the nominal
reference best move sometimes scores differently in its separate query. I do
not treat these numbers as exact chess truth, remove a bad case after seeing
it, or assume only a reference PV is acceptable.

## Why static calibration was insufficient

I trained on quiet played positions and one stronger-reference reply, with
only 404 eligible rows and ten opening families. A quiet static residual can
predict those labels more closely while changing a shallow search's preference
between response branches in the wrong direction. The engine also evaluates
many descendant positions outside that training distribution. The unchanged
search does not make the new evaluator's errors harmless.

The validation objective selected mean absolute static error. It did not
select coefficients for preserving good alternatives or avoiding a small
number of serious short-budget mistakes. Feature correlation and the small
family count further limit what the jointly fitted coefficients establish.
The fitted bishop and rook mobility weights reach their declared bounds in
both phases. That is a useful diagnostic clue: the constrained static objective
wants more of those correlated signals. It does not establish that a larger
weight would improve decisions. Teacher/student value calibration, limited data
and feature correlation are possible explanations that need fresh branch data.

The present results do not identify one coefficient as the sole cause, prove
that pseudo-mobility is fundamentally wrong, or rule out a compact neural
model. I should not infer a weight fix from one attractive failure line.

## The next model experiment I recommend

1. Collect new independent opening/game families, with all outcomes and a
   declared sampling schedule. Generate them offline through legal exploration
   and self-play if needed; no production opening book or Lichess dataset is
   required. Do not repurpose the reserved confirmation suite as training data.
2. Group whole games, transpositions, color mirrors and paired response
   branches together. Exclude overlap with the inspected controls and freeze
   train/validation/holdout membership before fitting. Use more independent
   families than the two validation and two holdout families here.
3. At the intended student budget, collect actual candidate alternatives,
   opponent responses and quiet search leaves. Preserve teacher score streams,
   bounds, mate cases and search budgets. Keep fresh-FEN static labels separate
   from history-sensitive decision labels.
4. Fit branch-value differences or a ranking loss alongside static residuals.
   Give little or no preference between alternatives within a declared finite-
   reference uncertainty band. Explicitly penalize severe decision errors;
   validate the student's actual search consequences rather than selecting on
   static MAE alone. Keep material fixed and the feature/runtime budget bounded.
5. Freeze one selected integer model on validation, then unlock a fresh holdout
   once. Require short-budget quality, tactical/mate/draw controls and equal-
   time cost gates before the standard concurrent strength screen. Keep a
   base-only implementation control, so I can attribute gains correctly.

Before collecting/fitting that dataset I need a separate numeric plan for its
family counts, sample budgets, uncertainty band, objective, regularization and
gates. I have not collected it or fit a replacement in this revision. Merely
shrinking the old weights, choosing another regularization after inspecting
these failures, or testing until one passes would reuse the failed controls
for selection. I do not recommend that approach.

## Priority after this revision

The exact base kernel is the immediate bounded candidate because its functional
boundary is already established: the same evaluation and fixed-depth search,
computed with less work. The standard game screen decides whether it merits a
promotion recommendation. It is not evidence that the learned terms work.

A decision-aligned dataset/model is the next larger evaluation project. A
compact NNUE representation remains a possible follow-up if the data supports
it; a board rewrite and additional selective pruning remain separate hypotheses.
I will reassess their priority from the screen rather than run a fixed checklist.
