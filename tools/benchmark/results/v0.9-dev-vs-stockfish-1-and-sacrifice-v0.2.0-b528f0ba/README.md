# v0.9 external comparisons

These fresh comparisons run **serially**, Stockfish first and Sacrifice second,
using the same candidate SHA-256 as the passing
[100-game v0.8 screen](../v0.9-dev-vs-v0.8/README.md). No parallel external
games count toward these results. The previous-release stage is recorded
separately; `--stage external` does not itself verify a previous result.

From the repository root, with the candidate and opponents already built:

```sh
tools/benchmark/bin/benchmark \
  --candidate tools/benchmark/bin/pandafish-candidate --version v0.9-dev \
  --stage external \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 69cf5c2a2f9652c8cb6e5ccbb416125ca24f4e27-dirty \
  --machine 'Apple M2 MacBook Air; all matchups serial; no stream' --verbose
```

See the [benchmark README](../../README.md) for repeatable baseline and opponent
builds. Each matchup uses 50 color pairs (100 games), the unchanged ten-opening
screening suite, 100 ms per move with 50 ms protocol allowance, and a 400-ply
cap. Both sides search on alternating turns, with one game running at a time.

Stockfish is **19 at native UCI Skill Level 1**, one thread and 16 MiB hash.
Sacrifice is package **0.2.0**, revision
`b528f0ba373910b72052d991700a990ca1cb07fd`, including philyawj's merged
[PR #1](https://github.com/piratesboys/Sacrifice-Engine-/pull/1).
Its Tal PGN move-ordering feature is enabled. Source/toolchain provenance is
in [opponents/sacrifice.json](../../opponents/sacrifice.json); executable and
PGN hashes are recorded in [summary.json](summary.json). Keep its source
checkout and `data/Tal.pgn` at their build location.

These 100-game screens measure outcomes under this profile, not precise Elo
or performance at every time control. Candidate scores use wins plus half
the draws; unfinished games remain separate. The Stockfish target remains
Level 1 pending a deliberate, supported advancement.

## Results

| Opponent | W / D / L | Score | Unfinished / forfeits | Candidate max response |
| --- | --- | --- | --- | --- | --- |
| Stockfish 19 / Skill 1 | 47 / 10 / 43 | 52.0% | 0 / 0 | 113.389 ms |
| Sacrifice 0.2.0 b528f0ba | 4 / 17 / 79 | 12.5% | 0 / 0 | 99.774 ms |

Pandafish's Stockfish result is the first routine screen above 50% at Level 1.
The margin is close in 100 games, so this alone does not establish a reliable
above-50% expectation or justify an automatic ladder advancement. The much
clearer evidence for this increment is the 78.5% previous-release screen.
Keep Stockfish at Level 1 until a deliberate advancement decision under the
existing progression policy. No longer run starts automatically.

Sacrifice remains substantially stronger under these settings. This is the
first controlled comparison against it, establishing a baseline for later
increments. It does not verify the anecdotal claim about beating Stockfish
Level 5, which may refer to a different difficulty scale or time control.

Three candidate responses against Stockfish exceeded the 100 ms request but
stayed inside the 150 ms deadline. None against Sacrifice exceeded 100 ms.
Opponent maximum responses were 130.021 ms (Stockfish) and 128.210 ms
(Sacrifice). There were no forfeits by either side in either matchup. These
are measured protocol response times, not worst-case search-time guarantees.

[summary.json](summary.json) contains identities, settings, all game outcomes,
and individual timings. [stockfish.pgn](stockfish.pgn) and
[sacrifice.pgn](sacrifice.pgn) contain the separate matchups;
[games.pgn](games.pgn) combines them.
