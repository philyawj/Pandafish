# Preserved reporting findings

The original ladder driver and its frozen source hash remain unchanged throughout
the games. All nine individual benchmark commands finish before its final
aggregation starts. Raw runner summaries, timing observations, PGNs, per-level
logs and start/end events are the primary evidence.

The final timing-count assertion originally equated attempted responses with
successfully played plies after the opening. A forfeited response is recorded
as an additional attempt without a played move. The three observed timeout
games therefore require exactly one extra timing record each. The report-only
continuation corrects that count while retaining each forfeit, its awarded
result and the failed reliability status. It reads original saved observations
and runs legal PGN checks only after all timed matches finish. No new game,
reference query, clock or move-selection change is part of this correction.

The collector finishes all nine 50-game matchups, saves `run-results.json`,
validates the first four rows and then exits 2 with `panic: timing count` at the
Skill 6 forfeit. `driver-run.txt` preserves the actual stack/output. The separate
report continuation exits 1 as expected because timeouts remain a reliability
failure, while its data-validation flag is true and all nine PGN checks pass.
The completed-game total is 449 with one unfinished game; the continuation's
original inherited “Completed 450 games” status text counts 450 recorded game
objects and must be read with the table's explicit unfinished column.

`report_source.go.txt` derives the continuation from the frozen driver and saves
`report.go.txt`; both are Go. The collection entry point is replaced with saved
method/run loading, failed-attempt timing counts are corrected, PGN review output
gets separate filenames, and data integrity is reported separately from the
overall reliability gate. The original collector's output and exit remain
preserved. A timing-affected screen does not receive an all-clear status merely
because its records can be parsed and legally replayed.

The first small forfeit-reader compilation also failed because of an unused
`fmt` import. Its empty stdout and compilation message are retained separately;
the corrected Go source supplies the actor/response details. This diagnostic
program does not start engines or affect playing policy.
