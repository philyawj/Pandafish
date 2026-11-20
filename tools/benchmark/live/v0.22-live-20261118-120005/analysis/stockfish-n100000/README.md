# Stockfish analysis

Stockfish 19; 100000 nodes/search; one thread; 64 MiB hash; Skill Level 20; MultiPV 1.

Processed 10324/10324 played plies, including 5159 by GoPandaFish.

| Classification | Bot moves |
| --- | --- |
| allowed_forced_mate | 33 |
| loss_100cp_plus | 367 |
| loss_300cp_plus | 63 |
| mate_score_comparison | 520 |
| missed_forced_mate | 10 |
| search_disagreement | 475 |
| small_or_equal | 3691 |

Mean centipawn loss among the 4596 cp-to-cp comparisons: 30.6 cp. Mate transitions are separate, not assigned artificial centipawns. These finite-budget comparisons are estimates, not Lichess accuracy or an Elo estimate. Negative raw deltas indicate search disagreement and remain recorded.

Full per-game JSON includes root history, scores, PVs and raw Stockfish output. Scores favor the root mover. Each independent query starts with an empty hash and ucinewgame. Played moves use a restricted-root search at the same node budget unless the unrestricted best move already matches.

473 candidate positions are saved in candidates.json. confirmation.json selects at most one per game, up to 20, for a larger-budget check. It first prefers the earliest substantial deterioration from a roughly playable position (reference at least -100 cp, loss at least 200 cp with the played result at most -100 cp, or newly allowed mate), then fills remaining slots from the ranked candidates. Selection is based on this survey and is not a random sample.
