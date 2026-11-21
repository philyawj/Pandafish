# Learned, bucketed piece–king context

Both guarded null policies fail cost; further shallow pruning is not the next
default. Return to the pre-live recommendation for a different, inexpensive
contextual representation. Earlier hand-weighted maps and the old 24-coefficient
fits remain rejected. This is one new model design, not their reactivation.

Reuse the balanced v0.19 sharp-branch corpus and its unchanged reference labels,
base scores, legal witness FENs and repaired whole-family boundary. The old
training split alone supplies fitting constraints. Old validation and holdout
are now **inspected development controls**, not newly unseen evidence. Their
prior failure remains intact. This avoids another collector/teacher campaign
before testing whether the representation has signal. A new final actual-search
cohort, collected only after weights are frozen, remains necessary.

Representation: five non-king piece types, relative to both own and opposing
king. Orient ranks for the piece's color; use absolute file distance capped at
four and relative rank difference capped to -4..4. There are 2×5×45 = 450 signed
coefficients, tapered by remaining phase/24. This is a linear statistical fit
over nonlinear board interaction features, not NNUE or a blocker/attack oracle.
Use one frozen fit, no grid: coefficient bound ±32cp, zero group sum excluding
the impossible coincident-square bin, ridge 4, Huber delta 200 and 25cp finite
gap uncertainty. Optimize projected accelerated gradient, at most 1,500 steps,
with backtracking and explicit convergence reporting.

Finite pair target is clip(teacher gap - original base gap, -400, 400), avoiding
the earlier already-correct-material direction bug. Normalize pair weights per
root. Keep every losing-mate/non-losing alternative ordering constraint with a
25cp margin; do not convert mate distances into numerical CP labels. Terminal/
draw/nonlearnable witness values remain fixed. Constant constraints remain
reported. No root move is prescribed and close finite alternatives remain valid.

Quantize group coefficients with bounded zero-sum rounding. Require exact saved
data/decision verification, legal witness/base controls, feature color symmetry,
zero-residual/target-sign/ordinal/solver/bound controls and convergence before
reviewing the new development result. Require >=5% lower capped mean decision
penalty on each inspected development split, no increased severe count and no
new avoidable-loss-reported root. Keep complete metrics and coefficients even
if this fails; do not refit to those outcomes. These gates test a frozen-leaf
surrogate, not actual search strength.

A surviving model gets an original Go runtime with precomputed phase/king/
piece-square corrections. Keep material constants/base placement, search,
quiescence, TT policy and clocks v0.22. Require scalar-table/original-base
equivalence, symmetry, terminal/draw/restoration/full/race controls and <=15%
median warm slowdown in three alternating pairs. Then require the original four
native gates and fixed seed 2026112093/stream93 final controls, blocking the old
training forest/blocked states, public/newer, consumed v0.23 pools and live survey
boards. Only then run one standard concurrent 50-game-per-opponent screen.

Limits: old depth-three witness trees approximate a newer actual-search policy;
selected labels cover at most six alternatives per root. King-relative geometry
does not see blockers, pins, support or every quiet threat. Historical controls
are inspected, finite reference estimates are fallible, and capped penalty is
not Elo. No automatic larger confirmation or publication follows. If this model
fails, preserve it and reassess data/representation instead of choosing a new
ridge, cap or favorable subset from its observed outcomes.
