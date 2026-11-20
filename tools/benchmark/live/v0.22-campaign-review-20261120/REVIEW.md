# Decision review and revised research priorities

The campaign supplies persistent errors worth investigating, but no isolated
v0.22 Elo estimate. The leading next question is why the engine prefers unsafe
king placement, weakening pawn moves and less useful piece activity in otherwise
contested positions. A contextual evaluation project moves ahead of the previous
mate-proof-first proposal, conditional on a root/leaf diagnosis establishing that
evaluation is responsible. No such implementation or experiment starts here.

## Survey and confirmation

The complete 100,000-node Stockfish survey covers 15,012 played plies, including
7,504 by Pandafish. Counts recomputed from the per-game records agree with all
three existing summaries:

| Classification | Bot moves |
| --- | --- |
| Finite CP loss below 100, excluding negative raw deltas | 5,327 |
| Negative raw CP delta / search disagreement | 671 |
| Finite CP loss 100–299 | 575 |
| Finite CP loss at least 300 | 107 |
| Newly found losing mate after the played move | 50 |
| Winning mate found only for the reference search | 11 |
| Other mate-score comparisons | 763 |
| Total | 7,504 |

The CP-only mean is 32.490 cp over 6,680 comparisons, including zero-clamped
negative raw deltas. Mate distances are not assigned artificial CP values.
This is a finite-budget loss surrogate, not Lichess accuracy. A single long
game can contribute many correlated positions: `rPaHjnoM` has 464 plies.

Before larger-budget queries, selection froze eight opening, eight middlegame
and eight late decisions from distinct games, with survey reference between
−100 and +300 cp and loss at least 150 cp. Four conversion opportunities, two
negative-delta controls and two agreement controls were added, each from a
different game. Phase boundaries and ordering are in [PROCESS.md](PROCESS.md).
This is a targeted diagnostic sample, not a random estimate of error frequency.

At **500,000 nodes per independent search**, all **24/24 CP deteriorations**
remain at least 150 cp, with a range of **151–659 cp**; seven are at least
300 cp. The mean of this deliberately selected subset is 285.083 cp. Complete
scores and raw UCI output are preserved, including changed reference moves.
Stockfish's chosen reference is not assumed to be the only acceptable move.

## Confirmed examples and what they establish

| Decision | Larger-budget evidence | Interpretation and limit |
| --- | --- | --- |
| `gNQRBta1`, 7...Kf8 | −67 reference versus −277 played; O-O is the reference | Material alone does not describe lost castling rights and exposure. This is a king-placement/development diagnostic, not proof of a missing constant bonus. |
| `Mido8Dw4`, 12.Kd2 | −43 versus −208; reference O-O | A second king walk instead of castling survives confirmation. Diagnose completed-depth evaluation and defensive continuations before attributing it to horizon or evaluation. |
| `IK3YVmhF`, 9.Nd2, 60+1 | −54 versus −713; played line begins Nd2 Bxh2+ Kh1 Bd6 Kg1 Ng4 | A concrete king attack is missed at live completed depth eight. The survey understated the damage. The line supports a tactical/king-exposure test; it does not prove a forced mate or a particular quiescence fix. |
| `6Nup5qBd`, 11...f6 | −88 versus −646; f6 exf6 Qxf6 Nd5 | Weakening a king near the center creates a large tactical deterioration. `50kq2J5b` has a separate confirmed ...f6 error, 224 cp. Shared pawn/king context is a useful hypothesis, not a causal finding. |
| `PV4oUeoe`, 9.Qf3, 60+1 | +11 versus −382; Qf3 Bxe5 Qxa8 b4 | A forcing material sequence does not make the queen excursion safe. Use the actual continuation rather than labeling the root as a generic hanging piece. |
| `u4xtX1iv`, 33.Ne1 | +19 versus −554; Ne1 Rde8 Qb7 Bd4 Nd3 g5 | Several important continuations are quiet activity/pressure moves. A capture-only description is insufficient; whether deeper main search or better leaf values solves this remains untested. |
| `jmVmUcjT`, 41.g3 | +240 versus +22; reference Ra6 | A rook/pawn ending loses its advantage despite live depth ten. Early development bonuses alone cannot cover the full sample. |

The [full confirmation table](CONFIRMED.md) includes all clocks, both colors,
remaining cases and legal continuation samples. These descriptions are
interpretations of confirmed lines, not instrumented attribution to a search
function, cache entry or evaluation term. Live warm-table state was not replayed.

## Controls change the mate interpretation

All four selected survey “missed forced mate” cases remain **winning mates
after the played move** at the larger budget:

| Game / decision | Reference / played mate distance at 500,000 nodes |
| --- | --- |
| `nViUqeab`, 22.Qxe5 | M6 / M10 |
| `m6B9i3bN`, 37...f5 | M8 / M9 |
| `vMJ64VjF`, 37...Bh4+ | M9 / M10 |
| `35s0eXdt`, 46.a5 | M7 / M9 |

These are slower conversions, not established missed wins. The third game
had earlier overlapping engine work; its later position remains valid offline
evidence, but its live throughput is not a serial timing control.

All **50** survey newly allowed-mate flags already have a reference CP score
of **−321 or worse**. None meets the roughly playable −100 cp criterion, so
the supplementary playable-mate-defence selection is empty and no searches
were issued for it. This does not make defensive mate recognition unnecessary:
earlier mistakes can cause the later losing positions, and finite CP is not a
proof of inevitable defeat. It does mean these flags cannot be counted as 50
independent losses a mate detector would have prevented.

The agreement controls `1TrbIGaE` and `2PKiqrUW` remain equal to the selected
reference. `ZiWwAAsF` changes from a negative-delta disagreement to exact
reference agreement. In `WwybQSjn`, the played Qh5+ still scores **better**
in its restricted search, +552 versus the unrestricted reference's +368.
That residual inconsistency is retained. Matching a single referee root move
is not an appropriate pass condition, and finite-budget reference search is
not an oracle over every alternative.

## Ranking after this result

This revises the [post-release ranking](../../research/v0.22-bitboard-exchange/NEXT.md).
It is a choice of the next question, not a fixed sequence that must all be run.

1. **Diagnose the persistent king-exposure/development decisions, then test a
   contextual evaluation representation if the diagnosis supports it.** Start
   with the independently confirmed king walks, ...f6 moves and Bxh2+ attack,
   plus quiet activity/endgame examples and safe castling/active-king negative
   controls. Compare recorded completed-depth choices with fresh root/history
   replays and inspect relevant leaf values. Distinguish mistakes corrected by
   another completed depth from stable undervaluation of king/pawn/piece context.
   This is the leading route toward a larger evaluation improvement, rather
   than another isolated safety coefficient or a reward for the teacher's root
   move. Preserve the old rejected models. A future model needs a strict Go-only
   data/inference pipeline, decision-level objectives, family-separated fitting,
   development and untouched holdouts, and an explicit inference-cost budget.
   These inspected positions are development evidence, not a future final
   holdout. Stop or rerank if the diagnosis chiefly identifies search omissions.
2. **A low-cost checking-move interface with a newly bounded tactical trigger.**
   The verified quiet-mate omissions from older research remain real, and the
   new Bxh2+ example supports investigating tactical reach. However, this
   campaign's mate-conversion flags provide weaker evidence of lost points
   than their survey labels suggested. First measure eligible checks and
   legal-response cost and prove which fresh errors the trigger could cover.
   Require both-turn escapes, blocks, captures, discovered-check, EP and
   promotion controls. The old failed proof remains inactive. Do not add all
   quiet checks recursively or assume mate recognition covers positional
   errors and long attacks.
3. **Reprofile before another search-efficiency or selectivity change.** The
   important decisions here complete depths seven through ten, so reaching
   useful extra depth could help. This review did not measure that benefit or
   reproduce warm-table choices; increased depth can also expose bad leaf
   preferences. Require fresh profile opportunity and root/leaf evidence before
   make/unmake, shared state, reductions or another cache-policy project.
   No result here implicates the earlier cache improvement as a regression.

A book or external Lichess training dataset is not a prerequisite established
by this evidence. Errors extend from early development into middlegames and
rook/pawn endings. A book could avoid some positions while leaving the general
decision problem intact. Additional representative data may be useful for a
contextual model, but the source should follow its objective and holdout plan.

Client overlap and ambiguous correlation deserve separate tooling follow-up;
they should not be bundled into an engine-strength experiment. There is no
reason to change time allocation in response to the expected laptop closure.
Any future candidate still needs correctness, actual equal-time quality and
the standard controlled previous-release / Stockfish Skill 3 / pinned Sacrifice
screen. No screen or larger confirmation of playing strength ran here.
