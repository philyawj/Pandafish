# Stockfish analysis

Stockfish 19; 500000 nodes/search; one thread; 64 MiB hash; Skill Level 20; MultiPV 1.

Processed 12/12 played plies, including 12 by GoPandaFish.

| Classification | Bot moves |
| --- | --- |
| allowed_forced_mate | 1 |
| loss_100cp_plus | 6 |
| loss_300cp_plus | 3 |
| small_or_equal | 2 |

Mean centipawn loss among the 11 cp-to-cp comparisons: 225.1 cp. Mate transitions are separate, not assigned artificial centipawns. These finite-budget comparisons are estimates, not Lichess accuracy or an Elo estimate. Negative raw deltas indicate search disagreement and remain recorded.

Full per-game JSON includes root history, scores, PVs and raw Stockfish output. Scores favor the root mover. Each independent query starts with an empty hash and ucinewgame. Played moves use a restricted-root search at the same node budget unless the unrestricted best move already matches.

10 candidate positions are saved in candidates.json. confirmation.json selects at most one per game, up to 20, for a larger-budget check. It first prefers the earliest substantial deterioration from a roughly playable position (reference at least -100 cp, loss at least 200 cp with the played result at most -100 cp, or newly allowed mate), then fills remaining slots from the ranked candidates. Selection is based on this survey and is not a random sample.
