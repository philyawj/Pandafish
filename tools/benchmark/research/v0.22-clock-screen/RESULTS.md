# Small clock screen: reject promotion

The smaller functional objective passes, but the one standard game screen fails
the established previous-release gate. Archive this code without promotion.
The result does not establish a statistically reliable regression; it does not
support promotion under the declared practical process either. No favorable
unchanged retry or larger confirmation runs.

The exact old clock source/binary passes 325 cold score/move/node pairs. Three
new timed warm pairs retain all 1,145 complete results and 9,874,565 nodes per
actor, with median ratio 0.93306 (**6.69% less time**). Independent row/helper/
native-repair hashing review passes. The original larger 15% target stays failed;
the [explicit smaller scope](../v0.22-next-assessment/REASSESSMENT_SMALL_SCOPE.md)
does not rewrite it.

All 1,252 native requests pass. Independently recomputed finite paired mean
reference-choice loss (25cp band, clamped to [0,1000]) is:

| Suite | Finite cases | v0.21 | Candidate | Severe v0.21 / candidate |
| --- | --- | --- | --- | --- |
| Public 100ms, three trials | 67 | 40.0746 | 38.4925 | 1 / 1 |
| Public 1000ms | 67 | 52.5821 | 52.5821 | 3 / 3 |
| Newer 100ms | 131 | 37.7023 | 33.0229 | 4 / 3 |
| Newer 1000ms | 132 | 30.2500 | 30.1061 | 2 / 2 |

No new losing-mate flag or response overrun occurs; both newer turn gates pass.
These inspected positions are regression controls, not a blind strength sample.
Finite scores do not establish that every labelled position was savable. The
[cache-argument failure](ERRATA.md) and reference-only repair preserve all 504
original first-block queries byte-for-byte; no actor query is rerun.

The standard 50-game matchups run **all three concurrently**, 100ms plus 50ms
grace, ten unchanged openings and a 400-ply cap:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.21 | 15 / 16 / 19 | 0 | 0 | 46% |
| Stockfish Skill 3 | 23 / 1 / 26 | 0 | 0 | 47% |
| Pinned Sacrifice | 30 / 11 / 9 | 0 | 0 | 71% |

The runner exits 1 after preserving all three matchups. All 150 games legally
replay, and independent score/color-pair/settings/attempted-timing accounting
passes. The full run takes 628.34 seconds. Raw PGNs, binary/source/asset hashes,
commands, timing samples and the failed gate remain under
[the screen](screen/v0.22-clock-screen-dev-vs-v0.21-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba/summary.json).

The observed functional gain is real on the measured Apple M2 / Go 1.25
workloads; it does not imply a playing-strength gain. Clock-read cost is platform-
dependent and depth-only untimed work does not benefit. The small samples and
repeated openings give wide strength uncertainty, but there is no reason to
spend an automatic larger confirmation on this failed practical screen.

Keep the code inactive for a future explicitly scoped efficiency objective.
The next assessment measures deferred move generation: the current main search
constructs every legal move before trying table/cutoff priorities. Preserving the
exact staged order may avoid meaningful work when an early move already refutes
the node. No clock change or other failed mechanism is bundled into that attempt.
