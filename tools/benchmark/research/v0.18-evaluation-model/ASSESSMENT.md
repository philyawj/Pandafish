# Fresh choice after v0.17

I published v0.17 through PR #32 and its tagged release. The backend candidate
met its engineering target (1.65–1.76x fixed-work throughput, much lower
allocation, exact node/score/move parity) and passed the practical screen:
59% against v0.16, 53% each against Skill 2 and pinned Sacrifice. The larger
confirmation remains deferred; the strength claims retain their sample limits.

I reassessed the queue, including alternatives to evaluation:

| Opportunity | Current evidence | Choice |
| --- | --- | --- |
| Coherent calibrated evaluation/threat model | Quiet attack and pawn-ending failures remain; both backend projects met their targets. A bounded model can use the existing bitboards without per-piece legal-move reconstruction. | Select a compact, jointly calibrated tapered classical residual model, with runtime and held-out quality gates first. |
| Further legal-move/position rewrite | Move tagging is still expensive, but immutable position updates are a small CPU share and the API route already cut allocation sharply. A full rewrite has substantially greater correctness/maintenance cost. | Defer pending a concrete stronger opportunity. |
| Selective search/history/aspiration | Potentially valuable, but the earlier broad LMR bundle failed; no new isolated evidence identifies a safer large gain than the diagnosed model gap. | Keep separate and reassess later. |
| Compact trained/NNUE evaluator | Plausible larger route, but requires representation/incremental-update, training/data and inference decisions. | Retain as an alternative; first test whether cheap interpretable features explain independent labels. |

CPU share measures cost, not a feature's strength value. The current tagged
playing code differs only in identity from the profiled v0.17 candidate, so I
reuse that verified profile as context rather than automatically repeat it.
I will measure the new evaluator's actual cost.

The selected model leaves the original material and placement score intact,
adding a 12-feature White-minus-Black vector with jointly calibrated middle/
endgame weights. Features cover safe pseudo-mobility, coordinated king-zone
pressure and shelter, passed-pawn advance/king access and pawn defects. No
isolated bishop-pair bonus, old piece-safety penalty, new search or clock rule
is reintroduced. Pseudo-attacks are a model input, not legal-threat proof.

I will freeze corpus selection and opening-family splits before calibration,
include stronger-reference reply positions as well as played positions, retain
all finite-node teacher iterations and failures, and use validation data only
for regularization selection. Holdout quality, positive/negative controls and
equal-time cost must survive before any standard strength screen. A teacher's
reference move is not the only acceptable root move. Static-model loss and a
few corrected positions do not establish strength; neither does another speed
number. No live games or new PR is started automatically.
