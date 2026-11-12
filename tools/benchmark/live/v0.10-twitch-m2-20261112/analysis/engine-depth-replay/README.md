# Captured-engine depth replay

Only the three deterministic development cases from the confirmed regression
set were used. The three public holdout cases were not used for these probes.
The executable is the original session copy, SHA-256
`8a5829be45932c8d6923b5a90fa2d151e7d74f9386564c55e501c059c4e7b1a9`.
No engine source, weights, or quiescence rule was changed.

Each fresh process searched only the played and Stockfish reference move at
the live completed depth, then one additional ply. A supplemental single-move
search at the extra depth checked whether the played move tied the selected
reference. There was no UCI time limit; a separate 20-second watchdog bounded
each probe. All nine probes completed to their requested depth with legal moves.

Scores below are Pandafish's own centipawns, not Stockfish's loss estimates.

| Game / decision | Live depth / move | Extra-depth two-choice result | Played-only score at extra depth |
| --- | --- | --- | --- |
| princhess_policy_bot, White move 24 | 7 / Nf4, +251 | depth 8: Rad1, +228 | Nf4, +224 |
| AetherBot, White move 15 | 6 / bxc4, -55 | depth 7: a4, +4 | bxc4, -2 |
| ShallowRedBot, Black move 19 | 6 / f6, -18 | depth 7: f6, +21 | f6, +21 |

The live-depth probes reproduced the recorded score and played choice. At one
extra ply, the reference strictly outranked the played move in two cases by
4 and 6 cp; this was not just random tie breaking. The third preference did
not change. None of these three live searches overlapped another process.

This is a restricted two-choice diagnostic, not a full-root correction or proof
that more depth solves every error. Other legal moves could outrank both.
The extra depth was also expensive: the first case's two-choice depth-eight
probe took about 11.7 seconds, beyond its live hard budget of about 4.565 seconds.
It would be inappropriate simply to force that depth during live play.

This supports testing better quiet-move ordering as a way to reach useful
main-search depth within unchanged deadlines. It does not prove that killer
moves will produce that speedup or a strength gain. The unchanged third case
keeps positional evaluation and tactical-horizon work on the later backlog.

- [inputs.json](inputs.json): cases, source links, live depth/score, and choices.
- [results.json](results.json): full two-choice UCI output and binary hashes.
- [played-only.json](played-only.json): exact played-move score checks.
- Diagnostic Go sources are preserved under [provenance](../../provenance/).
