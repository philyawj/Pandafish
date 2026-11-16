# v0.15 extended live test

One finished session on an **Apple M2 MacBook Air**, lasting 58h54m13s. It ran
in a terminal with occasional web browsing and YouTube music. There was no
Twitch stream. The original workload label was `terminal`; this operator
clarification supplements it without rewriting the manifest.

**265 completed games: 77 wins, 28 draws, 160 losses**, plus four unfinished
records. All 269 game IDs are distinct; there were 170 opponents, all BOT
accounts. The archive keeps every record, including the unfinished games.
Selected opponents, mixed clocks and a changing account rating prevent a
controlled strength or Elo inference.

This research is published with **v0.15.1**. The chess-playing policy remains
v0.15. A conservative reporting-tool repair recovers valid search/game links
without guessing the genuinely ambiguous games.

## Read the findings

- [Review and ranked next improvements](REVIEW.md): confirmed mistakes,
  counterfactuals, impact estimates and the recommended first experiment.
- [Testing volume](TESTING.md): why 50 games are a screen, uncertainty in the
  previous result, and a proposed larger confirmation stage.
- [Process and provenance](PROCESS.md): exact inputs, preservation, commands,
  deduplication, linkage repair and reproducible diagnostics.
- [Fresh-context handoff](HANDOFF.md): what to read and do next.
- [Client errors](CLIENT_NOTES.md): limits, reconnects and what was actually
  observed.
- [Original operations](OPERATIONS.md), [structured audit](session-audit.json),
  [file inventory](inventory.json) and [original session report](raw/summary.md).
- [Full 100,000-node survey](analysis/stockfish-n100000/README.md),
  [30 larger-budget checks](analysis/stockfish-n500000-confirmed/README.md),
  [captured-engine replays](analysis/engine-replay/README.md) and
  [repaired report](analysis/linker-fixed-report/README.md).

## Operational result

| Measure | Result |
| --- | --- |
| Saved / completed / unfinished games | 269 / 265 / 4 |
| Recorded searches / completed / legal | 12,091 / 12,091 / 12,091 |
| Missing final metrics / fallback moves | 0 / 0 |
| Original / repaired uniquely linked searches | 11,615 / 12,021 |
| Remaining unresolved search links | 70 |
| Overlapping recorded searches / peak processes searching | 0 / 1 |
| Own losses by time forfeit | 0 |
| Maximum hard-budget overrun | 1.132 ms |
| Response p95 / maximum | 12,211.459 / 37,462.405 ms |

The longest response followed a 15+10 clock allocation of 37,461.3 ms. It was
not a 37-second response to a 100 ms request. Timing measures end before
client submission and do not include network latency or CPU/temperature data.

The main rated 3+2 cohort has 189 games: 52 wins, 20 draws, 116 losses and one
unfinished. Its completed score is 33.0%. Other clocks and repeated opponents
remain separate in the [audit](session-audit.json).

## Offline referee result

Stockfish 19 at full strength surveyed all 24,208 played plies, including
12,082 GoPandaFish moves. It flagged 1,361 bot moves with at least 100 cp loss
and 129 newly allowed mate reports. Thirty selected distinct-game positions
were checked again at 500,000 nodes; all retained at least 100 cp loss or a
newly allowed mate report. They were chosen for investigation, not sampled to
estimate population accuracy. See the review for threshold changes and limits.

Captured v0.15 engine SHA-256:
`26f9af01670cfcb1bae9c261adf731cf7e2726034b717323130f8d0422abfdca`.
The original executable is retained locally under `raw/pandafish` and ignored
by Git. Its hashes and source revision are recorded; all text/JSON/PGN evidence
is preserved in this archive.
