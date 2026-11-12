# v0.9 positional-evaluation screen

This candidate adds original, small placement formulas to material evaluation.
Internal scores use centipawns. King placement blends home-rank preferences
with endgame activity using remaining non-pawn material. Search structure,
clock policy, quiescence cap, and random choice among exact ties are retained.
This does not add pawn-structure analysis or development-history penalties.

## Previous-release comparison

The previous-release matchup ran serially: one game at a time, alternating
searches between the two engines. It completed 100 games against v0.8 with
**66 wins, 25 draws, and 9 losses (78.5%)**, no forfeits, and no unfinished games.
This supports the increment under the recorded settings. It is a practical
screen, not a precise Elo estimate or proof of strength at every time control.

This stage was finalized separately before starting fresh serial external
comparisons. The [Stockfish and Sacrifice stage](../v0.9-dev-vs-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
uses the same candidate hash. No parallel external games enter either result.
The candidate is a development build, not a released v0.9 binary.

## Reproduce this stage

Build the tagged v0.8 baseline as in the [benchmark README](../../README.md).
From the repository root:

```sh
go run ./tools/benchmark --validate --race --version v0.9-dev \
  --stage previous \
  --previous tools/benchmark/bin/pandafish-v0.8 --previous-version v0.8 \
  --previous-revision 78598c2d839467b01ff4280ee3b4f48bb2c0ed8e \
  --revision 69cf5c2a2f9652c8cb6e5ccbb416125ca24f4e27-dirty \
  --machine 'Apple M2 MacBook Air; previous-release matchup serial; no stream'
```

The recorded binary hash identifies the actual candidate from local changes
on the stated base commit. Repeating the command need not reproduce identical
games because ties are random and search is timed. The ten-opening screening
suite uses 50 color pairs, 100 ms per move plus 50 ms protocol allowance, and
a 400-ply cap. The machine reported AC power during this run.

[summary.json](summary.json) preserves checks, identities, timings, and all
100 game records. [previous.pgn](previous.pgn) preserves the games, also copied
to combined [games.pgn](games.pgn). The source review's Lichess development and
holdout positions receive a separate fixed-depth spot check after all matches;
that check is diagnostic and does not measure centipawn loss.

## Additional validation

The [position audit](position-audit.json) supplies the original full position
history and requests `go depth 3` once per engine for each of the five public
development and four public holdout cases. It waits for `bestmove` and validates
legality. All 18 outputs were legal. Exact agreement with the single Stockfish
reference move was 2/5 for each version on development cases, and 0/4 for
v0.8 versus 2/4 for v0.9-dev on holdout cases. Equal alternatives and random
ties make exact agreement an incomplete measure; this is not evidence of a
particular centipawn-loss reduction or statistical generalization. No weights
were tuned after inspecting these outputs.

Go unit tests, race checks, vet, native ARM64 builds, and six final serial
integration-smoke games passed. The local lichess-bot wrapper validated both
private config copies and returned 12 legal plies per copy with the expected
v0.9-dev identity; the public example also loaded after substituting local
engine paths. No online games were started. The candidate's maximum response
in the v0.8 screen was 120.222 ms, with six responses past the 100 ms request
but none past the 150 ms deadline. The baseline maximum was 110.436 ms.
