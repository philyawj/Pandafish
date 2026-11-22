# Contextual nonlinear evaluation: bounded representation diagnosis

Tagged v0.23 playing source is restored. The compact search's 46% direct screen
and dead-material correction's 50% screen remain failed and fully archived.
Their functional benefits do not establish playing strength. The next priority
is assessment project 2, with its larger learning campaign explicitly deferred
until this bounded capacity/cost test.

Freeze one sparse king/piece representation: 32 horizontally folded exact own-
king contexts × 10 relative non-king piece types × 64 oriented squares. Two
perspective accumulators each have 16 int32 units. Shared int16 input weights,
clipped-ReLU activations and two 16-unit side-to-move output heads fit within
1MiB of parameters. Moving a non-king piece updates only its removed/added
features; king movement refreshes that perspective. Current material/placement
stays the score control. Synthetic input/output weights use a fixed SplitMix64
seed of 2026112203, independent of labels. No weight/model download is needed.

Use only training pairs from the historical sharp corpus: all 15 robustly
infeasible strict-order cases plus the identical-base/feature counterexample.
Do not inspect historical validation/holdout metrics or collect labels. Require
non-colliding activated pair vectors and a bounded output-head feasibility
certificate giving all selected pairs at least 200cp corrected separation.
Output coefficients are integers in [-4096,4096], divided by 64; use at most
20,000 cyclic projection passes, a 232cp continuous target to cover rounding,
then verify the actual integer 200cp target. The cost path keeps its original
synthetic output heads; the feasibility head is a separate non-playing proof.
Input weights stay synthetic and fixed. This is a capacity witness, not a
trained/generalizing evaluator; no playing-strength claim follows. Preserve a
failed certificate without changing seed, width, cap, target or subset.

Before cost, require exact full/incremental integer accumulators and outputs,
color/rank and file-mirror invariance, all special transitions, independent
596-vector children, cache-cell ownership, missing-king conservatism and race/
vet/build controls. Scratch search computes and consumes the synthetic output
but always returns original v0.23 evaluation, preserving its complete tree.
No playing candidate is activated by this diagnostic.

Three alternating warm pairs must preserve all 1,145 moves/scores and 43,118,858
nodes and have median elapsed ratio <=1.10 (at most 10% overhead). This is a
prospective inference envelope, not an accepted playing cost or speed gain.
No time from full-recompute parity checks enters the cost measurement. Stop the
representation on failed capacity, parity or cost; do not fit another seed or
width. Only all passed diagnostic gates justify a separately bounded data/fit
plan with untouched later controls. No games, live play or research publication
are authorized at this stage.
