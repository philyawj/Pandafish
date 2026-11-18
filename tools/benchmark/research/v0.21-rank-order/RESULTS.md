# Results and recommendation

**Recommend promotion under the established practical gate, after review.**
The candidate removes repeated ordering computation and large comparator copies,
preserves exact measured trees, passes native gates and completes a clean
standard screen. No playing-code change has been committed. No precise Elo or
large general strength gain is established by these samples.

## Functional efficiency

All 325 distinct frozen position/depth combinations retain scores, root moves
and search nodes. Three alternating public depth-four elapsed ratios are
0.79885, 0.81141 and 0.83991. Their median is **0.81141: 18.9% less time** for
identical measured work (about 23.2% higher throughput). Every ratio clears the
predeclared 15% median target; no outlier is removed.

| Workload | Positions | Identical nodes | Baseline ms | Candidate ms | Elapsed ratio |
| --- | --- | --- | --- | --- | --- |
| Public depth three | 84 | 227,579 | 295.86 | 262.28 | 0.88651 |
| Public depth four, median pair | 84 | 1,010,100 | 1,483.82 | 1,203.99 | 0.81141 |
| Deep subset, depth six | 12 | 3,658,933 | 5,133.43 | 3,922.75 | 0.76416 |
| Newer roots, depth four | 145 | 3,354,194 | 5,282.75 | 4,297.51 | 0.81350 |

`timing-analysis.json`, every row file, command/exit record and driver are kept.
The deeper subset improves 23.6%; newer roots improve 18.7%; shallower work
improves 11.3%. These are selected local workloads and do not guarantee the same
speed ratio for every game or hardware target. Adaptive profile benchmark mixes
are not used as controlled allocation comparisons. The additional per-ply rank
storage is a real implementation cost; no allocation-reduction claim is made.

## Native equal-time controls

All four predeclared paired suites pass. Loss is the per-position trial mean of
`clamp(reference_root_CP-selected_move_CP-25,0,1000)` on fully finite cases.
Close alternatives are accepted; reference mate cases remain separate.

| Suite / budget | Finite roots | Baseline loss | Candidate loss | Severe baseline / candidate | New mate-risk flags | Overruns |
| --- | --- | --- | --- | --- | --- | --- |
| Public / 100ms | 67 | 43.0299 | 40.0746 | 1 / 1 | 0 | 0 |
| Public / 1000ms | 67 | 53.8507 | 52.5821 | 3 / 3 | 0 | 0 |
| Newer / 100ms | 131 | 38.4198 | 35.0840 | 4 / 3 | 0 | 0 |
| Newer / 1000ms | 132 | 30.2576 | 30.2500 | 2 / 2 | 0 | 0 |

Both newer root turns pass individually. At 100ms, White loss is 46.7727 to
41.1818 and Black 29.9385 to 28.8923. At 1000ms, White is 37.1940 to 37.6567
and Black 23.1077 to 22.6154. The slight White regression remains within the
frozen tolerance; the overall 1000ms newer result is essentially even. The
candidate does not uniformly improve every decision or subcategory.

All 1,252 native requests, exact reference/cache observations, input/binary hashes
and failures remain saved. `choice-review.json` independently recomputes the
reports. Finite subsets differ between budgets as mate-valued choices change;
compare paired actors within a suite. These are reused regression controls, not
a blind holdout or playing-strength estimate. Finite teacher scores are not
proof that a position could be saved. Production roots are unrestricted.

## Standard concurrent screen

The single screen runs **50 games per opponent, all three matchups concurrent**,
using the unchanged ten-opening suite, 100ms move allowance, 50ms protocol grace
and 400-ply cap. Full identities, options, source revisions, timings and PGNs
are in [the screen records](screen/README.md).

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.20 | 23 / 9 / 18 | 0 | 0 | 55% |
| Stockfish Skill 3 | 23 / 4 / 23 | 0 | 0 | 50% |
| Pinned Sacrifice | 28 / 15 / 7 | 0 | 0 | 71% |

The runner exits 0; the direct practical gate passes. All 150 games legally
replay, and independent accounting verifies all 75 color pairs, statistics,
Stockfish Skill 3 options, attempted-response counts and timing bounds. No
response exceeds the protocol allowance. The source/binary/cache audit agrees
with the frozen overlay and checks 63 other protected engine/library files
unchanged. Restored-candidate tests/build and Linux/ARM64 build pass.

The direct five-point lead is modest and remains uncertain in 50 games with ten
repeated openings. The 71% Sacrifice screen is encouraging but does not establish
a specific improvement over the earlier 58% v0.20 development sample; those are
different runs without a simultaneous external baseline control. Similarly,
50% versus the tagged ladder's 48% at Skill 3 is a small sample difference.
Keep the measured outcomes without treating either external difference as a
precise gain. Stockfish remains Skill 3; no further target advancement follows.

A broader-opening confirmation is justified before claiming a large/general
strength gain: the existing optional 100 games per opponent over 50 openings
addresses part of the sample/opening limitation. It has not run and is not
added as a new mandatory promotion gate after the result. No favorable retry,
live game, experimental commit, push, PR or release is automatic.

## Findings

Fresh profiles exposed an avoidable representation/computation cost: a generic
stable sort repeatedly copied 128-byte PGN-capable Move values and recalculated
fixed priorities. Sorting 16-byte rank records and applying one exact permutation
removes much of that cost without new chess knowledge. Successful prior backend
work did not imply that move ordering was already cheap; remeasurement changed
the ranking. The result supplies another verified efficiency improvement and a
practical promotion screen, with strength uncertainty explicitly preserved.
