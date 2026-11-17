# v0.17 legal tactical generation/storage development screen

I tested one frozen additive-library-API candidate against tagged v0.16. It
masks quiet destinations before tactical legality/tagging work, tests legal-move
existence without a list, and reuses caller-owned move storage by ply. Existing
library runtime files and license are unchanged. Evaluation, search/PVS/root
policy, horizons, tables and time allocation remain v0.16.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.16 | 24 / 11 / 15 | 0 | 0 | 59% |
| Stockfish Skill 2 | 25 / 3 / 22 | 0 | 0 | 53% |
| Sacrifice v0.2.0 b528f0ba | 22 / 9 / 19 | 0 | 0 | 53% |

All three 50-game matchups ran concurrently at 100 ms requested time plus the
unchanged 50 ms protocol allowance, 25 color pairs over the standard ten
positions. Every game completes; all 150 PGNs independently replay legally
and match the saved results/checkmate terminations. Maximum candidate response
is about 98.2 ms. The practical previous-release criterion passes.

**I recommend promotion under the practical gate, pending review.** The direct
59% has a nominal paired 95% interval of about 47.5–70.5%, with repeated openings
limiting independence. The color split is 44% White / 74% Black. Both external
53% samples also have broad uncertainty; no Elo or reliably above-50% claim
follows. Historical external scores are context, not simultaneous baseline
controls. A larger confirmation is separately planned if stronger evidence is
needed; it has not run. No new commit, push, PR or live game was started.

[Exact identities/settings/timings](summary.json), [all PGNs](games.pgn),
[v0.16 games](previous.pgn), [Skill 2 games](stockfish.pgn),
[Sacrifice games](sacrifice.pgn).

[Full hypothesis, implementation, controls, preserved failures and limits](../../research/v0.17-search-moves/README.md),
[independent review](../../research/v0.17-search-moves/screen-review.json),
[next-step reassessment](../../research/v0.17-search-moves/NEXT.md).

I subsequently approved promotion of this screened playing code as v0.17 and
authorized its PR/merge/release. No broader confirmation was run. The original
development labels, hashes and game outcomes remain unchanged.
