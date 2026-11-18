# Reassessment after the 47% staged-picker screen

The first staged picker saves 19.81% warm search time but fails the direct game
criterion at 47%; the strong 68% Skill 3 row does not overrule it. No unchanged
rerun or promotion follows. Original controls, code and all 150 games remain
archived. Before restoring the released baseline, profile that frozen staged
implementation on the existing 1,145-root depth-five warm pipeline.

The next conditional idea is **lazy legal validation within ordered stages**.
The first picker defers whole categories, but a reached tactical or quiet stage
still tags/validates its entire list before the first move is searched. That
includes moves discarded after an earlier move cuts off. Quiescence also eagerly
tags every tactical candidate before searching a prefix. A pseudo-legal ordered
stream could defer king-safety/check tagging until a move is actually consumed,
while preserving the original complete legal sequence. This differs from simply
rerunning or adding the old clock/history/storage code.

First measure opportunity without playing changes. Count generated legal moves
and moves never reached after cutoffs, separately for main stages and quiescence
(checked/un-checked). Retain scores/moves/nodes against the frozen warm run.
An opportunity must cover at least **30% of reached-stage legal records**, and
tagging/validation must remain a substantial profiled cost, before implementation.
Opportunity counts are not projected elapsed savings. If the gate fails, archive
the diagnosis and reassess a genuinely different tactical-horizon/evaluation
design. No automatic implementation follows a big percentage on a tiny subset.

Any implementation must compare complete ordered whole Moves, original rules,
special moves, stale hints, cached metadata and active-ply ownership; preserve
all cold and warm scores/moves/nodes. Require at least **10% additional warm
median elapsed saving over the rejected first staged picker**, measured in three
alternating pairs, as well as the original 15% total target against v0.21. This
prospective extra-gain gate prevents a trivial edit from buying another screen.
The original four native gates, standard 150-game concurrent screen and failure
preservation remain the default; any exception needs its own explicit review.
No original failure is reclassified. No cache/history policy, heuristic, clock,
book, live game or root-specific move restriction is included.

Outcome: 61,664,803 generated legal records, 15,577,149 unreached (about 25.26%),
below the frozen 30% criterion. No lazy playing implementation follows. The
exact counts and saved-result comparison are independently recomputed by
`review_lazy.go.txt`; no elapsed claim uses this instrumentation.

The first diagnostic test starts before the asynchronous archive extraction has
finished and fails to compile against the incomplete scratch workspace. Both
original and copied error logs are retained. After explicitly waiting for the
successful replay result, a new log records the first completed diagnostic.
There are no search observations in the failed invocation and no native actor
query or match is repeated. This is an orchestration failure, not a TT defect.
