# Decision after the legal tactical generation/storage screen

I approved promotion of the screened candidate as v0.17 and authorized its
PR, merge and release. The broader confirmation remains deferred. The earlier
recommendation and development-state notes below are retained as the original
assessment, superseded by this explicit release decision. I also authorized
reassessing and starting one next candidate after release, with review before
its next PR. No live games are authorized.

**I recommend promoting this candidate under the existing practical gate,
subject to review before another commit/PR.** Its fixed-work score/node/move
identity, independent rule/state oracles, roughly 1.65–1.76x fixed-work
throughput, large allocation reduction and positive direct screen agree.
This is a functional improvement with promising measured strength evidence.
It is not a precise Elo estimate or a confirmed above-50% external win rate.

The standard screen completes 59% against tagged v0.16, 53% Skill 2 and 53%
Sacrifice, with zero unfinished games/forfeits. The direct nominal pair interval
includes equality, and the lead is concentrated as Black (74%, versus 44%
White). Repeated public openings, finite timing variance and the single
shallower native timing pair remain limitations. The local pinned-library
copy adds a maintenance responsibility despite unchanged upstream methods.

A larger frozen-source confirmation over the public 50-position suite is
justified for stronger strength evidence, separately planned and approved.
It has not run. The standard gate does not by itself require it, and I am not
launching one automatically. If a new playing change is made, the present
screen no longer validates that changed source. Keep Skill 2 as the routine
opponent; do not advance the ladder from a 53% sample.

## Reassessed queue after this result

I now favor **coherent data-calibrated evaluation and threat modeling** as the
next major strength experiment, if this candidate is approved. Both bounded
search-efficiency projects have delivered their engineering targets and
positive direct screens. The API route achieved the requested throughput and
allocation gains without a full internal-position rewrite. I would not expand
that rewrite merely because move legality remains a large CPU share.

The known quiet-mate, king-attack and bad pawn-ending models remain unchanged.
First build a diverse position/strong-response corpus with game/opening-family
splits and runtime controls. Then choose a bounded tapered classical model or
a compact trained route using data quality, cost and validation evidence.
Do not start by reviving isolated bishop-pair/piece-safety weights or fitting
only the engine's own PV. A Lichess corpus or playing opening book is not a
prerequisite for constructing that experiment.

Keep further move-tag/legality optimization and make/unmake as alternatives if
new profiling identifies a concrete remaining opportunity worth the rewrite
risk. Keep selective-search additions separate. Reassess again after the next
result rather than committing to every item in sequence. No third candidate,
new live game, larger confirmation or PR has started.
