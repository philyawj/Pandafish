# v0.10 Twitch session on the M2 MacBook Air

Recorded November 12, 2026, 17:32:12 to 19:00:00 UTC (09:32 to 11:00 PST).
The engine was the clean v0.10 release at
`1511e88b5e51752b775487822d27eebb2f61fa9c`, running alongside lichess-bot,
the browser, and Twitch streaming software.

All 40 source files are preserved byte-for-byte under [raw](raw/), with hashes
in [inventory.json](inventory.json). The original ignored session is
`local-games/v0.10-twitch-m2-20261112-093212`. The captured executable is retained
locally and hashed but remains gitignored. No Lichess API request was needed.

- 14 PGN records: 13 played games, 8 wins and 5 losses, plus one zero-move
  abandoned record. One win was the opponent's time forfeit.
- 550 completed legal searches, no fallback or missing final metrics.
- 510 searches have original PGN linkage, representing 509 unique positions;
  one position was searched by two processes. Forty searches remain unlinked
  in the original report. An additional audit finds individual matches for
  38 of those, leaving two ambiguous initial requests.
- The operator reported unintended simultaneous games. Trace-derived timing
  confirms 134 searches overlapping another process, a peak of two processes,
  and 168.5 seconds of overlapping request intervals. Configured concurrency
  was one. These intervals are not CPU utilization measurements.
- All 1,092 saved played plies were surveyed with unrestricted Stockfish 19
  at 100,000 nodes/search, including 548 moves by GoPandaFish. Twelve selected
  positions were rechecked at 500,000 nodes/search.

[Review and next-increment recommendations](REVIEW.md) explains the evidence,
ranked strength priorities, and why local tests should drive the next batch.
Check-aware quiescence is the preferred next experiment, with killer-move
ordering following as a separate efficiency experiment.
[The client investigation](UPSTREAM_CLIENT_NOTE.md) preserves exact error context
and an offline reproducer. It identified a duplicate-worker bug in the older
client that upstream had already fixed in PR #1262. The local checkout has now
been updated to version 2026.11.12.1 and verified offline.

| Record | Contents |
| --- | --- |
| [PROCESS.md](PROCESS.md) | Commands, provenance, processing fixes, and reproduction |
| [OPERATIONS.md](OPERATIONS.md) | Outcomes, controls, timing, and client incidents |
| [Survey](analysis/stockfish-n100000/README.md) | Complete 100,000-node comparison |
| [Confirmation](analysis/stockfish-n500000-confirmed/README.md) | Targeted 500,000-node follow-up |
| [Depth replay](analysis/engine-depth-replay/README.md) | Three development cases replayed with the captured engine |
| [Capture audit](capture-diagnostics-summary.json) | Compact linkage and concurrency diagnostics |
| [Detailed audit](capture-diagnostics.json) | Matching candidates and explicit overlap intervals |

The live results look encouraging relative to the earlier stream, but opponents,
positions, clocks, and workload differ. They are not a paired strength test or
an Elo estimate. Pandafish engine behavior is unchanged. The offline processor
received the small fixes documented in PROCESS.md; the companion client was
subsequently updated to include its existing upstream reconnect fix.
