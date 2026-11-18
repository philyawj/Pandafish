# Reference-level assessment after v0.20

**Recommend Skill 3 as the next routine reference.** Its clean 48% screen is
near even and provides a useful increase over Skill 2, which scores 62%. Skills
4 and 5 are useful harder references at 32% and 34%; the two-point inversion is
too small to rank their difficulty from these samples. Skills 8–10 are near the
floor and provide little resolution for routine incremental comparisons.
The actual default remains Skill 2 pending review. Timing failures in the upper
ladder remain a separate reliability concern; this report does not automatically
advance the target or treat the entire ladder as a clean pass.

The [complete table](README.md) links all nine summaries and PGNs. Each skill
gets 50 games, 25 per color, under the same paired-opening distribution. There
are **450 recorded games: 449 completed outcomes and one unfinished game**.
Three completed outcomes are timeout forfeits, which remain in the official
scores. All nine individual runner commands exit 0; saved failures still require
review because non-validation benchmark mode does not turn forfeits into a
process error.

## Release and conditions

[PR #36](https://github.com/philyawj/Pandafish/pull/36) is merged and
[v0.20](https://github.com/philyawj/Pandafish/releases/tag/v0.20) is published at
`c90ca44d71ce4c3897f43402976f4ca39ee1dc73`. This later benchmark tests an engine
and runner built from an isolated archive of that tag. The engine SHA-256 is
`c4ac32d3360db5ff813f94ba2906e67422b0d5ad8b43380499ce6222a61e9a8d`;
the pinned Stockfish 19 hash is
`dc2f18c34ae962dff591b66147d220ec06e61d756b93d8a4f5e04fd8e55c251f`.
These identities differ from the development/smoke labels, which remain preserved
in their original records. The development Skill 2 screen scored 67%; this
separate 50-game tagged sample scores 62%. Source validation preserves the
same playing policy through the version substitution; the sample difference
does not establish a five-point change in underlying strength.

The [frozen plan](PLAN.md), [method](method.json), [events](events.jsonl) and
[commands/exits](run-results.json) retain the experimental conditions. The
allowance is 100ms with 50ms protocol grace and a 400-ply cap. Stockfish uses
native UCI Skill Level 2–10, Threads 1, Hash 16 MiB, strength limiting disabled
and no Syzygy path. This does not calibrate Elo or Lichess GUI level numbers.

Three levels run concurrently in fixed mixed batches 2/5/8, 3/6/9 and 4/7/10 on
the same eight-logical-CPU ARM64 host on AC power. Active load falls at batch
tails. Each saved runner profile contains one matchup; the outer method/events
record the total concurrency. No other engine probes or heavy validation overlap
the games. No opening book, online dataset, live Lichess game or engine policy
change is introduced. All new programs and saved research source are Go.

## Preserved timeouts and unfinished game

The [forfeit details](forfeits.json) identify the last attempted responder:

| Skill | Game | Engine that forfeited | Response ms | Recorded result |
| --- | --- | --- | --- | --- |
| 6 | 45 | Stockfish | 169.346 | Pandafish win |
| 9 | 46 | Stockfish | 151.697 | Pandafish win |
| 9 | 47 | Pandafish | 163.416 | Pandafish loss |

All are recorded as `move response timeout`. Their causes are unresolved from
these timing records alone. Both engines were affected in batch two; the data
does not isolate engine computation, host scheduling or external host activity.
No deadline, time policy or game result was changed, and no failed game or level
was rerun. There are no recorded illegal-move or crash forfeits.

The official Skill 6 score is 18%, including the Stockfish forfeit win. Treating
that game's counterfactual outcome as unknown gives a descriptive 16–18% range
over all 50 games. Skill 9's only win is Stockfish's timeout; treating both
timeout-game outcomes as unknown gives 0–4%. These are sensitivity bounds,
not revised official scores or statistical confidence intervals. They leave the
upper-ladder difficulty conclusion unchanged but prevent interpreting the exact
18% and 2% as clean playing-strength estimates.

Skill 7 has one unfinished game at the 400-ply cap. Its 4 wins, 2 draws and 43
losses give 10.2% across 49 completed games. Across all 50, treating the unfinished
game as a loss through a win gives 10–12%. The original cap/result remains intact.

## Integrity, reporting and limits

The [report review](review.json) independently verifies statistics, options,
hashes, 25 color pairs per skill, response counts and all nine legal PGN checks.
`data_validation_passed` is true; `all_checks_passed` is false because of the
timeouts. The original collector finishes all matchups before failing its final
timing-count assertion: a forfeited attempt adds a timing record without adding
a played move. The [reporting notes](REPORTING_NOTES.md), original driver/log and
Go report-only continuation preserve that failure and its correction. The
continuation starts no new match and keeps the reliability rejection.

The postcheck verifies frozen source/binary identities, the release tag, equal
opening distributions, both colors and combined/individual PGN bytes. The final
inventory records every evidence file except itself. No result is discarded or
relabelled. The three recorded timeout failures are game results; the later
assertion failure is a reporting defect and cannot explain those timeouts.

There are only ten openings, repeated within 25 color pairs, at one short time
budget and one concurrency profile. Games and levels are correlated; 450 total
games do not supply 450 observations at any one skill. The 32%/34% inversion
illustrates why small adjacent-level differences should not drive selection.
This is a broad reference-selection screen with a timing qualification, not a
precise Elo estimate or a larger confirmation against the previous release.

The next decision is review of Skill 3 as the routine target. Keep previous-
release comparisons and retain Skill 2 as the historical anchor. Skills 4/5 can
serve as deliberate harder checks when useful; their extra games are not
automatic. Investigate the timeout outliers separately before using affected
upper-level rows for fine strength claims. No longer confirmation, additional
game, default change, second commit or PR follows automatically. Results remain
local for the requested check-in.
