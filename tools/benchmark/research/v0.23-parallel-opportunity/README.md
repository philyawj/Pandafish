# Two-worker opportunity passes; implement root-only first

The unchanged released workload covers 1,145 roots and 43,118,858 nodes.
Every move, score, completed depth and node count matches the released rows.
The serial-first-PV plus two-lane remaining-root proxy saves **38.193%** of node
span. Non-overlapping late-child regions cover 97.241%, but the simpler root
proxy already clears its frozen 20% gate. No nested splitting is selected.

These are recorded-node scheduling estimates, not measured parallel elapsed
time or playing strength. Concurrent windows may search additional work;
locks, imbalance, allocation and scheduling can erase the opportunity. Preserve
all rows, source, tests and exactness checks before implementing the candidate.
The implementation must report its extra CPU use and validate all workers'
state, cancellation and cache synchronization without changing score policy.
