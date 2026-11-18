# Rejected one-probe aspiration attempt

The [fixed selection](../v0.22-next-assessment/SELECTION_ASPIRATION.md) tests one
±50cp full-depth probe from depth three, accepting only exact in-window results
and retrying once at the original full window on a bound. Prior mate scores and
early depths bypass it. [Implementation](IMPLEMENTATION.md), three frozen Go
overlays/hashes, opportunity traces and all tests remain preserved over v0.21.

Bound/equality recovery, exact-score/best-set, bypass, canceled probe/retry and
last-completed-result controls pass, along with full race/vet/build. All **554**
iterative position/depth comparisons preserve scores and selected root moves.
Public depth-five ratios are 0.98644, 0.98429 and 0.98007 (median **1.57% less
time**); deep depth-six is 0.97535 and newer depth-four/five 0.98449/0.97516.
The public depth-three diagnostic is 2.6% slower and remains preserved.
`timing_gate_passed: true` here means the preregistered iterative no-regression
stage, **not** the later 15% warm efficiency gate.

That decisive warm gate fails. All three pairs retain identical complete scores,
root choices and repeated actor results on 1,145 recorded roots. Nodes fall
9,874,565 to 9,682,297 (1.95%), while elapsed ratios are 0.97952, 0.97596 and
0.97905: median **2.10% less time**, below 15%. The warm driver exits 1 and
retains `efficiency_gate_passed: false`. No independent membership driver, native
decision suite or strength screen runs after failure. Identical actual choices/
scores are visible in the saved reports; generic pending fields do not claim
unexecuted proofs. No elapsed number comes from the opportunity instrumentation.

Iteration-score stability alone did not establish a useful aspiration opportunity.
The existing first-best root/PVS/table coordination already avoids much of the
work a narrower root window could remove. That is an inference from the limited
node reduction, not an isolated causal percentage. Archive without promotion;
retain the exact source and every favorable/unfavorable observation. No strength
claim follows. The next assessment examines narrowly guarded, verified late
quiet reductions, as separate selective-search research rather than another
exact-work kernel. No aspiration/history/other failed code is bundled into it.
