# Fresh assessment after the v0.20 ladder

[PR #37](https://github.com/philyawj/Pandafish/pull/37) merges the complete Skill
2–10 data and adopts Skill 3 as the routine reference without a release. That
decision is complete. The playing baseline remains tagged v0.20.

This research then starts from merged commit
`6bfb4d565c2adf2a2ad59216a462b0f4f0c01503` on
`research/v0.21-next-candidate`. Read the initial [plan](PLAN.md), profile-driven
[selection](SELECTION.md), the [candidate report](../v0.21-rank-order/README.md),
and [reassessment](NEXT.md). No experimental code is committed.

The first implementation precomputes identical priorities and sorts small
stable rank records. It passes the cost, correctness, native and standard-screen
gates and is recommended for review. There is no rejected new runtime variant
in this cycle; all initial measurements and controls remain. Historical failed
quiet-proof, transition-only and evaluation attempts stay in their existing
archives rather than being reactivated or relabelled.

This directory retains both original CPU/allocation profiles, rendered reports,
line-level attribution, initial ordering tests and Go preparation/review source.
The baseline profile uses the existing fixed-depth driver from tagged v0.20;
the candidate timing driver is derived from the prior archived Go driver with
the new baseline path and predeclared 15% threshold. The saved preparation source
records those substitutions. These profiles are selected-workload evidence;
their percentages overlap and fresh-table allocation is a driver artifact.

The source reviewer checks the exact active overlay, protected source boundary,
binary hashes and reference caches. The screen reviewer includes the failed-
attempt timing convention learned from the ladder and verifies Skill 3 options.
Both are standalone Go sources stored as `.go.txt`; copy to `.go` before running
from the module. They start no matches. Their original temporary paths describe
this run; use adapted scratch paths and new output directories for reproduction.
The choice reviewer is the unchanged
[archived Go implementation](../v0.20-next-assessment/review_choices.go.txt).

Both new research directories receive byte/hash inventories after final review.
No original research inventory, raw failure or tagged release is rewritten.
Future work remains Go-only; no larger confirmation or live game is automatic.
