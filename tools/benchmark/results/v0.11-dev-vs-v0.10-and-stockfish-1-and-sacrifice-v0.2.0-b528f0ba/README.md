# v0.11 check-aware quiescence candidate

This candidate changes only the checked-node cutoff at the quiescence horizon.
Captures and promotions retain the normal two-ply allowance. A checked node can
continue past it, searching all legal evasions for at most four additional
plies. Non-checked nodes stop once the ordinary allowance is exhausted. The
final six-ply safety limit still returns static evaluation even if a checking
sequence remains. This is bounded horizon handling, not an unlimited tactical
search or a claim that every check sequence is resolved.

Evaluation, clocks, alpha-beta bounds, TT rules, and random exact root ties are
unchanged. The candidate identifies itself as Pandafish v0.11-dev. This
development build was reviewed for the v0.11 release; these results retain
their original labels and hashes.

## Method declared before results

Run the routine 100 games against v0.10 first. Continue only if the documented
previous-release gate passes, then run 100 against Stockfish 19 Skill 1 and
100 against pinned Sacrifice 0.2.0 b528f0ba, in that order. All games/matchups
are serial. No additional Skill 2 run is scheduled for this experiment. Keep
unfavorable results and do not rerun solely to obtain a passing score.

Settings: 50 color pairs, the unchanged ten-opening suite, 100 ms requested per
move, 50 ms protocol allowance, 400-ply unfinished cap. Hardware: Apple M2
MacBook Air, battery power at start (27%), no live bot. Prior v0.10 external
results started on AC power, so do not attribute differences between those
separate runs entirely to the engine. Within this run the engines share the
same host and alternate searching. AC power was observed early in the previous-release matchup after battery
power was recorded before validation.
[Power observations](power-observations.jsonl) preserve available samples;
they are not continuous measurements. This power transition limits the
consistency of the run and must remain visible when interpreting its score.

```sh
caffeinate -i go run ./tools/benchmark --validate --race --version v0.11-dev \
  --previous tools/benchmark/bin/pandafish-v0.10 --previous-version v0.10 \
  --previous-revision 1511e88b5e51752b775487822d27eebb2f61fa9c \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision a062c596ec5826eba58ce23a7c2aee1cfa4e57a4-dirty \
  --machine 'Apple M2 MacBook Air; serial matchups; battery power at start (27%); no live bot' \
  --verbose
```

The revision names the base plus local changes, not a commit containing v0.11.
Binary hashes and check commands belong to the generated summary. The engine
source/test patch is preserved in [candidate-source.json](candidate-source.json). The ordinary horizon
and four-ply extra budget were fixed before inspecting match results; there
was no sweep over extension values.

## Focused and saved-position checks

New tests reproduced the old cutoff before the search change: a checked side
scored about -4 pawns despite being able to capture the checking queen and
retain a rook. Tests now cover both colors, a forced quiet block and
countercheck, a capture/recapture that gives check at the normal horizon, the
unchanged non-checked cutoff, the ultimate safety cap, and a deadline during
an extended evasion. The independent exhaustive reference includes the bounded
extension and continues to check exact scores and alpha-beta windows. Existing
TT, root-tie, protocol, and cancellation tests remain part of validation.

The ARM64 candidate passed the updated local lichess-bot config/UCI wrapper
check with 12 legal plies. The original active v0.10 executable was restored
after the offline check; private configs and greetings were not changed. The public example greetings now
match v0.11-dev and loaded successfully after resolving its three engine paths.

[Saved-position probes](position-audit.json) compare v0.10 and the candidate on
all six public confirmed stream cases, including the three public holdouts.
Each gets a cold-process depth-three search and a separate 100-ms search,
serially. The fixture path and hash, engine hashes, legal moves, scores, nodes,
and depths are preserved, along with [the Go source](position-audit.go.txt).
No code was tuned after inspecting these outputs.

All 24 probes returned legal moves. Chosen moves were the same for both engines
on each request; completed timed depth was also the same on all six cases.
Scores and node counts changed on some cases. This small single-trial check
shows no correction of those six selected choices and no measured timed-depth
gain. It is integration/behavior evidence, not a strength estimate.

## Match results

All 300 games completed serially, with zero forfeits or unfinished games.

| Opponent | Wins / draws / losses | Score |
| --- | --- | --- |
| v0.10 | 30 / 51 / 19 | 55.5% |
| Stockfish 19 Skill 1 | 56 / 7 / 37 | 59.5% |
| Sacrifice 0.2.0 b528f0ba | 2 / 32 / 66 | 18% |

The previous-release screen passed. Its 5.5-point lead is modest and is not a
statistical demonstration of a small strength gain. Historical v0.10-dev
scored 52% against Skill 1 and 13.5% against this Sacrifice pin; this candidate's
external scores are encouraging, but separate randomized runs, repeated opening
positions, and different power conditions prevent isolating the size of a gain.
No additional Skill 2 run or automatic longer confirmation was performed.

Candidate maximum recorded response was about 98.1 ms across the 100-ms matches.
The full summary retains opponent timings too. The targeted tests demonstrate
improved horizon handling; the six saved real-game cases retained the same
choices and completed timed depths. These are distinct pieces of evidence,
not a claim that all previously observed mistakes are fixed.

**Assessment:** a reasonable small candidate to keep and review, supported by
the focused correction, a passing previous-release screen, and reliable local
play. The data suggests a modest positive change rather than a dramatic or
precisely measured gain. No live Lichess test is needed to repeat these local
checks. The user approved promoting this increment to v0.11 after reviewing the
local evidence. Release smoke checks are separate from these development games.

- [Full settings, identities, checks, timings, and results](summary.json)
- [Previous-release games](previous.pgn)
- [Stockfish games](stockfish.pgn)
- [Sacrifice games](sacrifice.pgn)
- [Supplemental verification](verification.json)

## Release preparation

Only version identity and documentation changed after the strength screen.
Formatting, Go race tests, vet, native ARM64 build, and a 12-ply local
lichess-bot wrapper check passed for the v0.11 release identity. A separate
six-game serial smoke (two each against v0.10, Skill 1, and Sacrifice) completed
with no forfeits or unfinished games. Those integration records remain outside
tracked strength history. Both private config greetings were updated to v0.11
and loaded successfully; credentials, challenge policies, and permissions were
preserved. No listener or online game was started. These checks do not relabel
the development games as results from the released binary.
