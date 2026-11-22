# Move-count claims: reject after the standard screen

The [frozen plan](PLAN.md) implements an available 50-move claim separately
from automatic 75-move termination. At 99 halfmoves, a legal intended non-pawn
noncapture move may support the claim; at 100 it is already available. Search
keeps a better winning continuation instead of forcing every claimable position
to zero. Mate takes precedence when the last move reaches 150; a position already
at 150 cannot continue toward a future mate. Root intended claims attach to
qualifying allowed moves. [The implementation note](IMPLEMENTATION.md) states
its finite-q and UCI/client boundaries. Existing repetition policy remains.

No rejected compact representation, dead-material detector or neural model is
reactivated. Score caches already distinguish halfmove clocks; their existing
protection is validated. Material/placement, legal generation/order, PVS, SEE,
repetition identities, horizons and time allocation remain v0.23 apart from the
specified move-count actions and terminal result.

## Functional and cost evidence

All eight root/nested normal/race/vet/build steps pass, as do 596 independent
move/state vectors. Focused controls cover both-turn mate precedence, 98/99/100/
149/150 boundaries, capture-only and pawn-only positions, EP/promotion/castling
resets, root restrictions, quiescence claim floors, history restoration, clock-
qualified score reuse and cancellation. A complete legal-action reference checks
main/root values at the declared q boundary; it is not an unlimited draw solver.

Three alternating warm pairs preserve every one of 1,145 root moves/scores and
43,118,858 nodes. Median ratio **1.001171** passes the prospective 5% cost ceiling.
The 0.117% observed overhead is within ordinary timing uncertainty; no precisely
established nonzero penalty or efficiency gain is claimed. Fixed-work equality
on this workload does not imply semantic equality at the new rule boundaries.

## Existing actual-budget controls

| Suite | Mean finite reference loss, baseline → candidate (cp) | Gate |
| --- | ---: | --- |
| Public 100ms | 45.896 → 45.896 | Pass |
| Public 1000ms | 54.851 → 54.851 | Pass |
| Newer 100ms | 30.328 → 30.328 | Pass |
| Newer 1000ms | 23.598 → 24.833 | Pass |

No new losing-mate flags or timing overruns. The newer long-budget suite worsens,
including a **163cp** White-root regression: baseline completes depth 6 and the
candidate depth 5. [Matched-depth review](newer-1000-depth-review/summary.json)
confirms the same move/score/node results at each depth (197,214 nodes /356cp
at depth 5, 1,325,093 nodes /363cp at depth 6). This explains the observed depth
boundary, not every timing cause, and does not remove the original timed loss
or turn it into an improvement. No timed actor query is repeated.

## Fresh controls and pending decision

The single frozen seed 2026112204/stream 126 collects 129 roots, 65 White and
64 Black, across all 16 families. The collector passes with 560 explicitly
resolved pre-existing input paths. No favorable extension/replacement occurs.
Exclusions include historical/public/newer/prior-fresh, warm/oracle, live analysis,
both preceding complete screens and draw-policy controls. Temporary adapters
avoid tracked duplicate trajectories. This cohort is consumed regardless of
outcome. Both fresh gates pass with unchanged mean finite loss: 38.971cp at
100ms and 29.236cp at 1000ms, with no new mate flags or overruns.

## Sole standard screen: reject

| Opponent | W / D / L | Unfinished | Forfeits | Score |
| --- | --- | --- | --- | --- |
| Tagged v0.23 | 23 / 2 / 25 | 0 | 0 | 48% |
| Stockfish Skill 3 | 31 / 1 / 18 | 0 | 0 | 63% |
| Pinned Sacrifice | 27 / 16 / 7 | 0 | 0 | 70% |

The original greater-than-50% direct gate fails. No unchanged retry or promotion
occurs. All 150 games finish and legally replay, with all 85 source files and
binary identities unchanged. [Complete report](screen/README.md) and
[rule exposure](screen-rule-exposure.json) retain adverse results. Only one
candidate turn reaches halfmove99, against Sacrifice; both recorded fifty-move
adjudications there are valid. No candidate root in the direct or Skill 3 match
reaches halfmove90. Exposure is limited, and no win attribution or Elo gain is
established. The functional claim tests remain valid despite failed promotion.

Restore v0.23 source before another mechanism. Reassess tactical frontier
semantics with one bounded unchanged-tree diagnostic; another draw-rule variant
or unchanged screen is not justified. All research remains uncommitted. Cleanup
PR #42 remains merged without a release; no live play occurs.
