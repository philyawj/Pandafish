# Bounded quiet checking: reject after the standard screen

This candidate is **not recommended for promotion**. The single standard
concurrent screen finishes all 150 games cleanly but scores exactly 50% against
v0.22, failing the unchanged greater-than-50% direct criterion. No unchanged-code
retry or larger confirmation follows. The complete useful generator and failed
playing implementation remain archived; restore released v0.22 before continuing.

| Opponent | W / D / L | Unfinished | Forfeits | Score |
| --- | --- | --- | --- | --- |
| Tagged v0.22 | 21 / 8 / 21 | 0 | 0 | **50% — fail** |
| Stockfish Skill 3 | 23 / 4 / 23 | 0 | 0 | 50% |
| Pinned Sacrifice | 34 / 9 / 7 | 0 | 0 | 77% |

The [full screen](screen/v0.23-quiet-checks-dev-vs-v0.22-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba/README.md)
retains every PGN, move timing, identity and option. The driver exits one for
the direct strength gate after every concurrent matchup completes, not because
of a crash. Previous v0.22 external scores were 71%/70% in a separate screen;
the current 50%/77% is a mixed result, not a paired Elo estimate. It does not
justify promoting an apparent attacking benefit against Sacrifice alone.

## Mechanism and correctness

The additive checking-only API geometrically filters direct checks and
discovered-check blockers before full legal tagging, with checking castles and
conservative malformed/prechecked-position fallbacks. It matches all 596 frozen
independent move/tag cases plus 2,507 explicit/seeded controls. Its median cost
is 325.4ns versus 1472ns for all-quiets-and-filter, approximately 78% less.
Residual allocations remain; this is not a whole-engine speedup.

The search adds one layer of quiet checks at an unchecked main-search horizon,
after existing captures/promotions. Required evasions and ordinary capture
continuations keep their original bounds. No additional quiet checking layer
is added after a reply or capture. Material, placement, clocks, score-cache/
draw policy and main-search ordering remain v0.22. The earlier evaluator and
scratch child cells are inactive.

Quiet mates and knight forks become visible in the functional controls.
Escapes, captures, blocks, layer bounds, parent state and color symmetry pass.
Normal/race, library race, vet, Linux build and clean-overlay tests/build pass.
The original capture-only exhaustive-oracle failures and two tooling/export
failures remain in [control notes](CONTROL_NOTES.md), rather than disappearing
from the record.

## Actual-choice controls and the generalization limit

All six choice suites pass their original mean/severe/mate/time tolerances.
Both turns pass where required, with no new losing-mate flags or overruns:

| Suite | v0.22 mean finite loss | Candidate |
| --- | --- | --- |
| Public 100ms | 45.896 cp | 36.821 cp |
| Public 1000ms | 54.851 cp | 49.970 cp |
| Newer 100ms | 30.679 cp | 32.153 cp |
| Newer 1000ms | 24.212 cp | 20.083 cp |
| Fresh 100ms | 27.008 cp | 25.932 cp |
| Fresh 1000ms | 32.744 cp | 32.120 cp |

The fresh 147 roots come from 16 fixed, independently seeded rollout families,
74 White and 73 Black, with prior inspected boards excluded. No fitting uses
them. Differences on that sample are small. There are 2,428 native actor
requests across the six suites, but repeated trials/budgets are correlated and
are not independent strength observations. All teacher data remains finite
Stockfish estimates with a 25cp acceptable-alternative band, not a prescribed
root move or an oracle.

The favorable surrogate and functional results **do not transfer into a direct
match gain**. Adding a tactical layer changes the value/breadth tradeoff and
can lose useful main-search work at equal time. This screen does not isolate
how much is cost versus selective-horizon behavior. The fast API remains useful
research infrastructure, but that alone does not justify changing play.
Future checking-stage deferral or selectivity would need a separate cost/risk
diagnosis and new actual-budget evidence; none follows automatically here.

[Process and reproduction](PROCESS.md), the prospective
[screen decision](SCREEN_DECISION.md), original/compatible source manifests,
all quality reports and `screen-status.json` preserve the full chain. The
candidate binary SHA-256 is
`a0a02a403aa20560a8e9207d4abe51b2daa0c466f4cd71f7ada09214880bf6f1`;
the development label stays `v0.23-quiet-checks-dev`, not a released version.

Reassess a different pre-live direction: inexpensive phase-sensitive piece
activity through the already precomputed tables. Earlier board-feature fits
and king pressure had substantial inference cost; changing the table formulas
can test evaluation at the original hot-loop cost. This is a new hypothesis,
not permission to tune the rejected checking layer or assume the live sample
provides every best next step.
