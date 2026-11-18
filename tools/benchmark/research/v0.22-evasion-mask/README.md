# Exact evasion-mask cost prototype: rejected before games

This revision retains the exact first staged picker and adds a geometric
check-evasion destination prefilter. [The plan](PLAN.md),
[implementation](IMPLEMENTATION.md) and [controls](TESTING.md) describe its
scope and the explicit post-diagnostic opportunity review. Eight Go overlays
over v0.21 and their hashes are frozen. Original upstream source is unchanged.

All focused controls, 596 independent vectors, root and upstream races,
vet/build and all 325 cold fixed-work results pass. Cold depth-four median time
falls 24.34% versus tagged v0.21; deep/newer ratios are 0.70540 / 0.74115.
All warm scores/moves/nodes also agree on 1,145 roots and 43,118,858 nodes per
actor, repeated exactly over three alternating pairs.

The required **10% additional warm saving over the rejected first staged picker
fails**. Ratios are 0.910674 / 0.902710 / 0.908855, median 0.908855 (**9.11%**).
The original driver returns exit 1 and `efficiency_gate_passed:false`. Preserve
every pair and the near miss. No native actor suite, game screen, strength
claim or larger confirmation follows this version.

The failure triggers a narrow source/profile reassessment: checker context is
currently constructed before the generator tries any king move, including a
legal-existence call that often returns immediately on its first king escape.
A later revision may defer that context until a nonking candidate actually needs
it. That is a code change requiring its own frozen source, unchanged extra-gain
gate and paired evidence; no unchanged favorable rerun is justified.

Restore using `tools/research replay --source` with `candidate-source.json` into
a new directory. Drivers preserve original commands/workspaces/settings. The
original oracle-path error and corrected first completed oracle remain visible.
The source is reusable functional research, not promoted playing code.
