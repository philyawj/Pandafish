# Pandafish v0.8: first Twitch stress test

A complete capture from November 11, 2026, running on an Apple M2 MacBook Air
while the same laptop displayed games in a browser and streamed to Twitch.
The recording ran for 5 hours, 38 minutes, 5 seconds, ending November 12 UTC.
The engine came from clean release commit
`78598c2d839467b01ff4280ee3b4f48bb2c0ed8e`.

## Evidence

- [Original file inventory and SHA-256 hashes](inventory.json)
- [Capture manifest](raw/manifest.json), [hardware note](raw/hardware.txt),
  [redacted configuration](raw/config.redacted.yml), and [client log](raw/lichess-bot.log)
- [PGNs](raw/games/), [per-process protocol traces](raw/traces/),
  [captured search records](raw/moves.jsonl), and [capture summary](raw/summary.md)
- [Operational review](OPERATIONS.md), [structured operations](operations.json),
  and [classified client errors](client-incidents.json)
- [Stockfish survey](analysis/stockfish-n100000/README.md)
- [Final synthesis and recommendations](REVIEW.md)
- [500,000-node confirmation](analysis/stockfish-n500000-confirmed/README.md) and
  [regression candidates](analysis/stockfish-n500000-confirmed/regression-candidates.json)
- [Post-session client environment](provenance/client-environment.json)
- [Process and reproduction instructions](PROCESS.md)

The original ignored source remains at
`local-games/v0.8-twitch-m2-20261111-125756`. Every captured data file is retained
byte for byte. The executable is present in the local archive and has a recorded
hash; it remains gitignored. Neither that engine nor the original files are
needed to rerun the offline Stockfish analysis from a repository clone.

The capture contains 38 finished games against 34 bots: 5 wins, 2 draws, and
31 losses. One win was the opponent's time forfeit; there were no recorded
Pandafish time-forfeit losses. All 1,531 recorded searches completed and linked
to a PGN without capture warnings, illegal moves, or legal fallbacks.

This is a live operational/position dataset with mixed opponents and time
controls. It does not establish an Elo improvement over v0.7. The controlled
Stockfish Skill Level 1 comparison target is unchanged.

The complete survey analyzes all 3,069 played plies at 100,000 nodes/search.
Twenty selected positions were checked at 500,000 nodes/search. Read
[the review](REVIEW.md) for the findings and why positional evaluation is the
recommended next hypothesis to test. The original data remains unchanged.
