# Tagged v0.20 Stockfish Skill 2–10 screen

This plan precedes the ladder's games. v0.20 is released from merge commit
`c90ca44d71ce4c3897f43402976f4ca39ee1dc73` through
[PR #36](https://github.com/philyawj/Pandafish/pull/36).
The engine and benchmark runner are built from an isolated archive of that tag,
not from a development executable. The unchanged production target is Skill 2.

Run **50 games at each of Skills 2, 3, 4, 5, 6, 7, 8, 9 and 10**: 450 total.
Each level uses 25 color pairs, the existing ten-opening screening suite,
100ms per move, 50ms protocol grace and a 400-ply cap. Openings, strength limiting,
Threads 1, Hash 16 MiB and disabled Syzygy remain the standard runner settings.
Only Stockfish's native UCI Skill Level changes. This is not the Lichess GUI's
level numbering or an Elo calibration. No book, online game or new engine policy
is introduced.

Three levels run concurrently per batch, on the same host and AC power. Fixed
batches are **2/5/8, then 3/6/9, then 4/7/10**, so every batch includes a lower,
middle and higher skill rather than assigning the entire high end to the last
period. Individual runners each have one matchup; the outer Go driver records
the global maximum of three. Actual active load declines as a batch finishes,
as in the routine three-opponent screen. Retain start/end times, all individual
profiles, response times, binary/runner/opening hashes, options and every PGN.
Do not overlap other engine probes or heavy validation with the games.

Finish every requested level even after an unfavorable score or runner failure.
Do not retry unchanged code or extend the sample for a preferred result. Preserve
failures and incomplete rows visibly. After all games, legally replay each PGN,
independently recompute statistics/color pairing/options and aggregate one table.

The question is which next reference level provides a useful challenge after
the approved improvement. Interpret the curve broadly: 50 games and ten repeated
openings cannot resolve small differences, and adjacent levels need not produce
monotonic samples. Current-level reliability and the existing above-50% practical
advancement guidance still matter. Identify the first higher-level region where
performance falls toward/below even, and recommend a useful next target without
claiming a precise strength equivalence. Review the entire curve, not the single
best-looking level. No automatic target change, longer confirmation, commit,
second PR or live game follows this screen; report the complete table first.
