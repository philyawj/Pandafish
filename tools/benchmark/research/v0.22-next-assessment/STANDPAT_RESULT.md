# Stand-pat relocation opportunity rejected before implementation

The original 1,145-root / 43,118,858-node trajectory completes without a bound
contradiction. The conservative alpha>=0 predicate identifies 4,273,265 pre-SEE
opportunities, including 2,921,830 of 7,611,901 actual SEE calls (38.39%). Only
2,433,206 actual child/node visits would be avoided (5.64% of all nodes), because
many of the moves are already rejected by SEE. All source, counters and raw
observations remain preserved. Extra diagnostic board updates/static evaluations
make its elapsed time unsuitable as a speed measurement.

The paired criterion requires >=20% node coverage and >=30% SEE coverage. The
second condition passes; the first fails. No playing shortcut or incremental
evaluation code is implemented. A cheaper exact SEE/child-evaluation path remains
a possible future narrow objective, but the result does not justify claiming a
large search-work gain or relaxing draw/check guards. In particular, negative-
alpha cases cannot ignore a possible zero-score stalemate/repetition outcome.
The separately reassessed small clock assessment is not a claim that this
opportunity or any prior 15% research target passed.
