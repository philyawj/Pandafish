# Bounded two-worker root search

The released-tree diagnosis retains all 1,145 results and shows a 38.193%
root-only node-span opportunity. Implement that simpler design, not nested
splitting. Keep material/evaluation, ordinary child search/quiescence, table
score identity/capacity/replacement, root eligibility and time budgets v0.22.

Search the prior PV/root-first move serially. Distribute remaining root moves
over two private worker states. Each snapshots the current incumbent and uses
the existing child window. Only a proven strict improvement of that snapshot
can update the incumbent; a cutoff bound cannot replace it as a false tie.
Exact equal improvements may prefer the earlier original index. Clone real
repetition context; isolate move/rank buffers and killers; share only properly
synchronized score/hint tables. Aggregate every worker's nodes, join all workers
before returning and preserve incomplete-depth/fallback behavior. Reuse lane
buffers across completed depths within one go command.

Default timed/depth production search uses at most two workers on hosts with
GOMAXPROCS>=2; direct zero-value/reference workers remain serial. Fake-clock
unit references stay serial; dedicated real-clock/race controls exercise the
parallel path and cancellation. Two candidate workers are additional CPU use,
not per-core efficiency. Three concurrently active candidate turns fit six
workers on the eight-core M2; opponents wait for their turn. No pondering/live
play is added.

Require legal exact-value/co-best controls, immutable parent/history restoration,
mate/draw/cache-context and cancellation/worker-shutdown checks under race tests.
Parallel scheduling may change node counts and tied moves. A separately declared
parallel warm protocol must require baseline score identity, legal choices and
independent full-window validation of every selected distinct move. Do not call
the old deterministic-node protocol passed. Require >=15% median wall-time saving
in three alternating unprofiled depth-five pairs, reporting actual node ratios.

Then retain the existing four native gates, new fixed final seed
2026112113/stream113 and one standard concurrent 50-game-per-opponent screen.
No automatic larger confirmation or publication. Present resource tradeoffs,
complete results and all rejected predecessors before the requested commit check.
