# Child-state reuse on the released v0.22 kernel

This prospectively reassesses the pre-live storage thread after the contextual
evaluator fails. The old v0.21 attempt remains rejected: 3.3% cold median and
4.65% warm saving did not meet its original objective. Reusing its sound
ownership mechanism is not a claim that the old gate passed.

New evidence is the restored v0.22 depth-five warm profile: 1,145 roots,
43,118,858 nodes, 11.92 GB measured allocation. Child state accounts for 97.4%
of sampled allocated bytes; UpdateSearch has 11.6% cumulative sampled CPU,
with considerable overlapping runtime/GC work. These percentages are not
summed or treated as guaranteed elapsed savings. The prior 35–38% exact-kernel
improvements changed the relative bottlenecks. Only new unprofiled pairs decide.

Rebase the archived per-active-ply SearchPosition cells onto v0.22, retaining
all current staged generation, evasion, exchange and immutable transition APIs.
Do not overlay old v0.21 search files wholesale. A child survives recursion,
history cleanup and both PVS probes; descendants use distinct cells. Game-root
positions and other workers cannot alias scratch. Search scores, choices, nodes,
ordering, horizons, material, TT/draw policy and clocks must remain exact.

Before native queries, require all existing storage/oracle/special-move/race
controls, all 325 fixed-work comparisons and 1,145 warm roots to agree exactly.
Run three alternating unprofiled depth-five warm pairs without per-root
MemStats; require **at least 15% median warm elapsed saving** against tagged
v0.22. Require a beneficial cold comparison and no contradictory deep result.
The warm target emphasizes current production table lifetime and deeper work;
it is a new objective, not a retroactive change to the failed old cold target.
No favorable repetition or threshold adjustment follows a failure.

A surviving revision receives the four actual-choice suites, fresh frozen
controls if justified, and one standard concurrent 50-game-per-opponent screen.
Preserve every failed check and keep the evaluator inactive. If this context
still gives insufficient elapsed benefit, archive it again and move to the
checking-only interface rather than following allocation bytes indefinitely.
