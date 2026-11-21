# One bounded research screen with a failed surrogate

This decision precedes this binary's only game screen. The original requirement
that all six native/fresh gates pass is **not satisfied**. Fresh 1000ms remains
failed: mean finite loss 27.3971→32.9485cp, severe errors 2→3, with Black
31.6765→38.0882cp and White 23.1176→27.8088cp. Its false fields, runner failure,
all queries and frozen thresholds remain unchanged. No actor-query timing retry,
coefficient revision, root restriction in production or alternative cohort is used.

The other five suites pass; no suite has an illegal move, query error, timing
overrun or new losing-mate flag. Fresh 100ms quality is unchanged. Current
hardware shows 25.86% median fixed-work elapsed saving over three alternating
pairs, about 4.1% more searched nodes and exact baseline scores at all 1,145 roots.
A separate serial full-window pass validates all selected moves. The candidate
uses two workers versus one: the gain spends more CPU and is not per-core speed.

The initial interpretation was immediate archival on the fresh failure. Review
of the prior staged-picker and v0.22 decisions identifies a relevant distinction:
an exact-value efficiency change can reach a worse pre-existing search horizon.
The three largest fresh regressions all change from completed depth six to seven.
Before choosing any match outcome, a shared Go depth-review command selects
every changed choice across all six suites, including improvements and neutral
changes. At their observed completed depths, 60 unrestricted fresh-table native
comparisons on 30 distinct positions have equal scores and valid selected moves.
Any differing co-best choice needs a separate serial full-window proof. This
supports a depth/horizon explanation on these cases; it does not make the bad
moves good, establish universal equivalence or erase the failed surrogate.

As in the documented v0.22 investigation, the substantive speedup and these
diagnostics justify **one standard research screen** as an explicit exception
to the surrogate prerequisite. This is a revised evidence-gathering decision,
not a claim that the original plan passed. The unexecuted original screen plan
is retained alongside the separately named exception plan. The failure remains
part of any eventual promotion review.

Run 50 games each against tagged v0.22, Stockfish Skill 3 and pinned Sacrifice,
all three matchups concurrent, at 100ms plus 50ms grace, with the unchanged ten
openings and 400-ply cap. Require the unchanged direct criterion, zero forfeits
and review both external rows, shutdown/timing behavior, resource use and this
failed gate together. No second unchanged screen, automatic larger confirmation,
live game or publication follows. Stop before commit if review recommends promotion.
