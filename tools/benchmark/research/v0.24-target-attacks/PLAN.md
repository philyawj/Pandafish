# Exact target-centric attack kernel

Base: tagged v0.23 playing source; research branch after cleanup merge #42.
The preceding full pin-context diagnosis fails its prospective8% headroom
objective and remains preserved. No unchanged diagnostic is rerun.

Hypothesis: ordinary tag queries can preserve the released attack result with
less work by querying attackers at the king from original piece sets under the
virtual occupancy. For own-king legality, mask the captured destination from
the result; no mover belongs to the opponent's attacker sets. For opponent
check, mask the old mover origin, then add a direct attack from its new square.
Use existing pawn-origin masks instead of rebuilding whole pawn attack sets.
Retain the original geometric early rejection. Specials and malformed origins
keep the prior path. No per-node pin context or eager8-ray scan is added.

Before timing, all two-transition-oracle controls,596 independent vectors and
normal/race/vet checks in both modules must pass. Add focused tests comparing
the new attackers with the original virtual implementation for both colors,
captures, king moves, pins, prechecked/malformed positions and all specials.
Whole-tree values, choices and nodes must stay exact.

The exact-kernel objective is >=8% median warm saving (maximum ratio0.92) on
three alternating pairs using all1,145 frozen depth-five roots /43,118,858
nodes. Do not use diagnostic elapsed times or relax this target after results.
If it fails, preserve the implementation and stop later native/game work.
If it passes, the standard four existing and two fresh actual-budget gates and
the usual concurrent50 games each vs v0.23/Skill3/Sacrifice remain necessary.
Fresh data needs a frozen new seed and explicit exclusions for every consumed
cohort; no collection begins before its prerequisite gates.

Keep the original upstream runtime files unchanged. Evaluation, move order,
search/q policy, cache/history, repetition/mate scores and clock budgets stay
v0.23. No release/live play is authorized. Stop before publishing a recommended
playing candidate. Use shared Go execution and lean records; preserve failures.
