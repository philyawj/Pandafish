# Measure opportunity for two-worker search

Ten substantive attempts now lack a promotion recommendation: seven playing
prototypes, one fitted model and two opportunity diagnoses. Exact small-kernel
changes and evaluation/search heuristics have limited or adverse evidence.
Reassess a different architecture: bounded two-worker search under the existing
wall-clock allowance. This uses more CPU, not a claim of better per-core work.
On the eight-core M2, three concurrent matches have at most six active candidate
workers because each opponent waits for its turn. No live session is authorized.

Before implementation, measure the released warm depth-five workload without
changing its tree. Record every root child's actual nodes. Estimate the serial
first-PV plus two-lane remaining-root makespan using the recorded child costs.
Also count non-overlapping late-child regions at root/main depths >=3, protecting
the first child. This diagnoses possible nested split work if root-only splitting
has too little headroom. Do not double-count nested late-child regions.

Require all 1,145 root moves/scores/depths/nodes to match v0.22. No instrumented
elapsed time is performance evidence. A root-only node-makespan proxy saving
>=20% supports the simpler implementation; otherwise require >=50% outer late-
child work before investigating nested splitting. These are opportunity proxies,
not guaranteed speedups: stale windows, synchronization, imbalance and overhead
can erase them. Reject before parallel code if neither clears its frozen gate.

A surviving implementation must bound concurrency to two, preserve draw/cache
identity, material/evaluation, quiescence and time budgets, isolate worker state,
and synchronize shared tables without changing their score policy. Validate exact
values/legality/co-best choices, repetition, cancellation and complete worker
shutdown under race tests. Parallel ordering may change visited nodes and tied
moves, so do not relabel the old exact-tree protocol as passed. Require >=15%
median wall-time saving in fresh alternating pairs under an explicit parallel
protocol; report additional CPU/resource use. Then keep native/fresh quality
gates and the standard concurrent 150-game screen unchanged. Final seed, if
earned, is 2026112113/stream113. No automatic confirmation/publication follows.
