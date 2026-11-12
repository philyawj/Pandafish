> Preliminary diagnostic run before the returned-move/PV consistency fix. Superseded by `../stockfish-n500000-confirmed/`. Preserved as an execution record; exclude from final findings.

# Stockfish analysis

Stockfish 19; 500000 nodes/search; one thread; 64 MiB hash; Skill Level 20; MultiPV 1.

Processed 20/20 played plies, including 20 by GoPandaFish.

| Classification | Bot moves |
| --- | --- |
| allowed_forced_mate | 6 |
| loss_100cp_plus | 7 |
| loss_300cp_plus | 2 |
| mate_score_comparison | 5 |

Mean centipawn loss among the 9 cp-to-cp comparisons: 299.1 cp. Mate transitions are separate, not assigned artificial centipawns. These finite-budget comparisons are estimates, not Lichess accuracy or an Elo estimate. Negative raw deltas indicate search disagreement and remain recorded.

Full per-game JSON includes root history, scores, PVs and raw Stockfish output. Scores favor the root mover. Each independent query starts with an empty hash and ucinewgame. Played moves use a restricted-root search at the same node budget unless the unrestricted best move already matches.

15 candidate positions are saved in candidates.json. confirmation.json selects at most one per game, up to 20, for a larger-budget check. It first prefers the earliest substantial deterioration from a roughly playable position (reference at least -100 cp, loss at least 200 cp with the played result at most -100 cp, or newly allowed mate), then fills remaining slots from the ranked candidates. Selection is based on this survey and is not a random sample.
