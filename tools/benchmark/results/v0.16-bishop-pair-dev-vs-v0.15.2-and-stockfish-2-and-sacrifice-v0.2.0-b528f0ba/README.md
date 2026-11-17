# v0.16 bishop-pair development screen

The candidate adds one untuned 25 cp bonus for a side with bishops on both
square colors, inside the existing evaluation traversal. Material values,
placement, search, quiescence, TT and timing policy remain unchanged.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.15.2 | 26 / 5 / 19 | 0 | 0 | 57% |
| Stockfish 19 Skill 2 | 21 / 2 / 27 | 0 | 0 | 44% |
| Sacrifice v0.2.0 b528f0ba | 7 / 11 / 32 | 0 | 0 | 25% |

All three 50-game matchups ran concurrently, using 100 ms moves, 50 ms protocol
grace, 25 color pairs and the standard ten-position book. No forfeits or
unfinished games; no candidate response over 100 ms. The practical previous
release gate passes, but its paired 95% score interval is about 44% to 70%.
The screen does not establish a strength gain or external regression against
a simultaneously tested baseline. Do not pool opponents or extrapolate Elo.

**Recommendation: hold, not promote.** A separately approved 200-game
confirmation per opponent is justified for the frozen source. No larger
confirmation, new commit, push, PR or live client has been started.

[Exact settings, binary identities/hashes, checks and all game metrics](summary.json),
[all PGNs](games.pgn), [previous games](previous.pgn),
[Skill 2 games](stockfish.pgn), [Sacrifice games](sacrifice.pgn).

[Fresh selection rationale](../../research/v0.16-next-assessment/README.md) and
[full hypothesis, source, controls, cost, regressions, uncertainty and replay](../../research/v0.16-bishop-pair/README.md)
separate the functional exchange correction from measured strength.
