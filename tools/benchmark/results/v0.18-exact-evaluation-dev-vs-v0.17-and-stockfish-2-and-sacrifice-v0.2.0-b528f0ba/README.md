# Exact evaluation kernel screen

I tested `v0.18-exact-evaluation-dev`, an unpromoted exact-score speedup over
released v0.17. The learned terms from the earlier failed model are absent.
The three matchups ran concurrently, 50 games each, 100 ms per move plus the
unchanged 50 ms wall allowance, with the public ten-opening screen.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.17 | 24 / 9 / 17 | 0 | 0 | 57% |
| Stockfish Skill 2 | 36 / 1 / 13 | 0 | 0 | 73% |
| Pinned Sacrifice | 21 / 9 / 20 | 0 | 0 | 51% |

All 150 games independently replay legally. The direct practical gate passes.
I recommend promotion of the isolated kernel under that gate. Its efficiency
is established more firmly than the magnitude of a strength gain: the direct
nominal pair interval includes equality, openings repeat, and the direct gain
clusters in the first two openings. Equal-opening sensitivity is about 52.08%
direct; it does not replace the official score. The 73% Skill 2 result is
encouraging and Sacrifice is close to even. Historical external scores are not
a simultaneous baseline control. I make no precise Elo or broad strength claim.

[summary.json](summary.json) records exact engine names, binary/asset hashes,
opponent options, source identities, settings, every result and response times.
[games.pgn](games.pgn), [previous.pgn](previous.pgn),
[stockfish.pgn](stockfish.pgn) and [sacrifice.pgn](sacrifice.pgn) retain every game.
The tested development label is preserved; no released v0.18 binary was tested.

[The research report](../../research/v0.18-exact-evaluation/README.md) contains
implementation, controls, exact-score/performance evidence, limitations,
[the independent game review](../../research/v0.18-exact-evaluation/screen-review.json)
and [next recommendations](../../research/v0.18-exact-evaluation/NEXT.md).
A broader-opening confirmation is justified before broad claims or raising
Skill 2; none has run. No live games, commit, push or new PR were started.
