# Retrospective released-policy architecture control

This tests the **exact tagged v0.15.2 executable**, not the bishop-pair candidate,
a cache prototype or the v0.15.3 research release. Its playing behavior is v0.15.
All three 50-game matchups run concurrently with the standard 100 ms / 50 ms
grace / 400-ply / ten-opening profile and color-swapped pairs.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.14.4 | 21 / 14 / 15 | 0 | 0 | 56% |
| Stockfish 19 Skill 2 | 21 / 2 / 27 | 0 | 0 | 44% |
| Sacrifice v0.2.0 b528f0ba | 7 / 13 / 30 | 0 | 0 | 27% |

The practical previous-version gate passes. Fifty games are not enough to
assign a reliable Elo gain. Mate-distance and repetition behavior also changed
between these tags, so this is not a causal cache-only strength comparison.
The external results provide current baseline context for the archived bishop
screen, not proof from subtracting historical point estimates.

**Provenance correction:** the command supplied a wrong descriptive old-revision
abbreviation. The old executable was built from tag v0.14.4 at
`3f3ca331cbe8aac9c7c2194b9314b33b3da34a81`. The original summary and command are
retained; [the companion provenance](../../research/v0.15.3-architecture/control-provenance.json)
records the actual commit/hash and erroneous label. No game was rerun or changed.

[Original complete summary](summary.json), [combined PGN](games.pgn),
[previous-version games](previous.pgn), [Skill 2 games](stockfish.pgn),
[Sacrifice games](sacrifice.pgn).
[Full architecture investigation](../../research/v0.15.3-architecture/README.md)
and [next larger projects](../../research/v0.15.3-architecture/NEXT.md).
