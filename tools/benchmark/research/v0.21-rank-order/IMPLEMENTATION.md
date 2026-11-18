# Compute ordering keys once and sort small records

The candidate changes the computation of the existing stable move order. It
does not change priorities or eligibility. Each generated move gets its original
index and the identical promotion/capture/EP priority once. A stable sort compares
those 16-byte records; the old comparator passed 128-byte Move records and
repeated the same map and board lookups at every comparison. If priorities are
already descending, the original list is already the required stable result.

Sorted ranks name the original source index for each destination. Disjoint
permutation cycles place the full Move records in order, using one saved Move
per cycle. Whole-value assignments preserve flags and PGN metadata. Completed
cycle indices are marked in private scratch data. Original ordering helpers and
their reference stable sort remain unchanged for differential validation.

Every active ply owns its rank buffer alongside its existing move buffer. The
searcher owns this storage; descendants and re-searches cannot overwrite an
ancestor's active list. This adds small reusable scratch storage. It is not a
score cache and does not change transposition entries, repetition identity,
draw counts, hints, killers, search depth or clocks. No library file, piece
value, static evaluator or exchange estimator changes.

Initial ordering controls and full race/vet/build checks pass. Randomized
permutation controls cover 2,100 lengths/tie patterns, with independent reference
sorting of complete Move values. Existing seeded legal-position controls compare
full/tactical ordering against the unchanged helpers, and buffer controls cover
repeated cold/warm use, different plies, empty lists, promotions, EP and castling.

All 325 distinct frozen fixed-depth comparisons retain scores, root moves and
nodes. Three alternating depth-four elapsed ratios are 0.79885, 0.81141 and
0.83991, median 0.81141: 18.9% less selected-workload time, exceeding the declared
15% gate. Depth-six time falls 23.6%, newer roots 18.7%, depth three 11.3%.
These are exact-work efficiency results, not established playing strength.
Native quality and the Skill 3 standard screen remain separate gates.
