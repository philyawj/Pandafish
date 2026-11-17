# v0.16 coordinated search-core development screen

I tested the frozen `v0.16-search-core-dev` candidate: one proven best root,
full-depth PVS, separate legal board hints, reversible-context score identity
and tactical-only sorting. Evaluation/material, quiescence eligibility/horizons,
score-table replacement/exact-depth rules and time management stay unchanged.
The separate hint table adds about 1 MiB. Specialized move generation remains
a future experiment.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.15.3 | 27 / 11 / 12 | 0 | 0 | 65% |
| Stockfish 19 Skill 2 | 21 / 1 / 28 | 0 | 0 | 43% |
| Sacrifice v0.2.0 b528f0ba | 16 / 17 / 17 | 0 | 0 | 49% |

All three 50-game matchups ran concurrently, using 100 ms moves, 50 ms protocol
allowance, 25 color pairs and the unchanged ten-position screening suite. The
run completed in about 11 minutes. Every PGN move was independently validated
with python-chess; all 150 results and checkmate terminations check out. No
unfinished games or forfeits; maximum candidate response was about 99.7 ms.
The practical previous-release criterion passes.

**I recommend a separately approved broader confirmation before release
promotion.** The direct score is promising measured strength evidence alongside
verified work/depth improvements. Its nominal paired 95% interval is about
53.5–76.5%, but repeated public openings and shared-host timing limit that
calculation. The direct score is 50% White / 80% Black. Seven opening positions
score above 50%, two at 50%, and one below; broader color/opening checking is
still needed. I do not infer an Elo gain or a population score of 65%.

Sacrifice's 49% is encouraging beside the historical released-policy control
at 27%; this is not a simultaneous external baseline comparison. Stockfish's
43% is similar to the historical 44% control. This screen does not establish a
broad leap against all opponents or assign gains to individual components.

[Exact settings, binary identities/hashes and every game timing](summary.json),
[all PGNs](games.pgn), [v0.15.3 games](previous.pgn),
[Skill 2 games](stockfish.pgn), [Sacrifice games](sacrifice.pgn).

[Full hypothesis, implementation, controls, failures, raw probes and limitations](../../research/v0.16-search-core/README.md),
[independent screen review](../../research/v0.16-search-core/screen-review.json),
[proposed confirmation and reassessed next steps](../../research/v0.16-search-core/NEXT.md).

The original candidate identity and hashes are retained. No larger confirmation,
live game, commit, push, PR or release was started after the screen.

I subsequently approved promotion of this screened playing code as v0.16 and
authorized the PR, merge and release. The broader confirmation remains deferred;
the recommendation above records the pre-approval assessment. Original
development labels, binary hashes and game results are unchanged.
