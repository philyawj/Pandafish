# Prospective standard screen decision

All six actual-choice suites pass for the frozen binary. No surrogate exception,
threshold change, fitted coefficient or production root restriction is used.

| Suite | v0.22 mean finite loss | Candidate | Result |
| --- | --- | --- | --- |
| Public 100ms, three trials | 45.896 cp | 36.821 cp | Pass |
| Public 1000ms | 54.851 cp | 49.970 cp | Pass |
| Newer 100ms | 30.679 cp | 32.153 cp | Pass, small worsening retained |
| Newer 1000ms | 24.212 cp | 20.083 cp | Pass |
| Fresh 100ms, three trials | 27.008 cp | 25.932 cp | Pass |
| Fresh 1000ms | 32.744 cp | 32.120 cp | Pass |

Both required turns pass in newer/fresh suites. There are no new losing-mate
flags or timing overruns. Fresh roots come from the fixed 16-family rollout,
147 retained roots (74 White, 73 Black), with prior inspected boards excluded.
Fresh differences are small and do not establish a playing-strength gain.
All reference estimates, negative deltas and mate cases remain recorded.

Full normal/race, library race, vet and Linux build checks pass. Clean-overlay
tests/build pass after an export-format-only correction, with original failure
and full snapshot preserved. Generator controls match 596 independent rule
positions plus explicit/seeded cases, and generator cost meets its original
50% saving criterion. Initial quiet-mate, escape/capture/block and bounded-layer
controls pass. The capture-only exhaustive oracle's initial failure is preserved
and its independent complete-list minimax now uses the specified new horizon.

Frozen executable SHA-256:
`a0a02a403aa20560a8e9207d4abe51b2daa0c466f4cd71f7ada09214880bf6f1`.
Original full source snapshot SHA-256:
`6a33e02dc18c5bce42fa92a843093054bf4d2482e6d7131d83cdd24a066b401b`.
The development label is `v0.23-quiet-checks-dev`; UCI identity remains v0.22
until promotion is separately approved. The tagged baseline binary is
`e8d2f6d8e19207921419ebac1a44528aa6d7e92b993aa8ef18988bcbe6e4b003`.
Material/evaluation, TT/draw policy, main-search ordering and clocks remain
released v0.22. Neither earlier rejected v0.23 implementation is active.

Run **one** standard 50-game-per-opponent screen against tagged v0.22, Skill 3
and pinned Sacrifice b528f0ba, concurrent, 100ms plus 50ms grace, the unchanged
ten openings and 400-ply cap. Review the unchanged direct criterion and both
external rows together with functional evidence and sample limitations. Preserve
every game, result and failure. Do not rerun unchanged code for a favorable
score, launch a larger confirmation, start live games or publish automatically.
Only a candidate that merits promotion reaches the requested pre-commit review.
