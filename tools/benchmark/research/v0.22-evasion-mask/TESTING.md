# Controls before cost and playing evidence

Initial focused controls pass 6,036 independent full/tactical ordered whole-Move
comparisons over cold/warm caches, seeded check-seeking walks and special FENs.
Explicit positive controls retain white/black EP capture of a checking pawn,
EP interposition and an ordinary bishop block. Original full legal methods are
the expected result. Existing 406,674 lookup/subset controls and 10,300 exact
staged/eager lists still pass, alongside buffer isolation, zero-allocation first
hints and startpos depth-four perft.

The nested package's focused tests cover single sliding/non-sliding checks,
multiple checkers, inactive unchecked positions, missing/nonunique kings and an
inconsistent check flag. These run in the pinned upstream replay's full root-
package race suite. The original runtime hashes are unchanged. Root race,
vet/build and all 596 independent frozen move/tag/child-state/status/perft
vectors pass. Source and tests are frozen before timings.

The initial oracle invocation names nonexistent `oracle-cases.json`; it fails
before reading any observation. Original and copied error logs remain in
`oracle.txt` and `oracle-input-path-failure.txt`. The corrected invocation reads
the existing frozen `oracle-inputs.json`, saves `oracle-correct-input.txt` and
the first completed `oracle-results.json`. This is a path error, not a failed
oracle assertion or a rerun of actor queries/games. No runtime code changes.

Cold timing uses the original 325 distinct position/depth comparisons and three
alternating public depth-four pairs versus v0.21. Warm measurement compares the
new revision with the frozen first staged picker in three alternating depth-five
1,145-root persistent-table trajectories. All moves/scores/nodes must agree and
repeat per actor. The **10% extra** warm gate is unchanged after the opportunity
review; original total 15% and deeper/newer requirements also remain. A separate
comparison with the saved v0.21 rows establishes exact-work total identity, but
historical elapsed ratios are descriptive rather than new simultaneous pairs.

No heavy work overlaps fixed timings, native actor blocks or games. Preserve
every failed gate and all source/commands. A failed cost attempt is archived
before native/play work. No live play or automatic larger confirmation occurs.
