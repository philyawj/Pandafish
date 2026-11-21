# Verified pass probe: reject on cost

The [prospective policy](PLAN.md) revisits the old selective-search thread with
two full pass-probe plies, checking visibility confined to that hypothetical
search, conservative material/pawn/draw guards and a separate real-move
verification. Released normal quiescence, evaluation and cache rules stay intact.

Initial full tests, focused guards and the full race suite pass. The controls
cover both-color quiet mate threats, draw priority, deadlines/cancellation,
parent FEN/hash/repetition restoration and no real-table/hint pollution. The
checking generator still matches 3,103 independent/seeded rule positions.

Three alternating, unprofiled fixed-depth-five pairs cover all 1,145 frozen
warm-table roots. Each candidate result is legal and complete. Actual nodes,
including both isolated searches, fall from 43,118,858 to 38,854,946 (9.89%).
Only 11 root move/score results differ. Elapsed candidate/baseline ratios are
1.005688 / 1.000668 / 1.001103: median **0.110% slower**. Both the 20% node and
15% elapsed-saving targets fail. The driver exits one; no native or match run
is earned and no strength conclusion follows.

Verification consumes the intended search saving on this workload. Dropping
guards or verification solely to raise the rate would be a different experiment,
not a repair of this failure. The old unverified null-move screen remains a
separate adverse result, not explained retroactively by these timings.
All source, tests, workload rows and logs remain here; runtime is restored to
released v0.22 before the pawn-structure attempt. The full source snapshot is
`1e43c2b73f9e086e26207f4ecaccfe092d543d9ea3ae0e4ab269a6f29d859374`.
