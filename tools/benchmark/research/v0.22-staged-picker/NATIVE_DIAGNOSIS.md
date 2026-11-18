# Public short-budget surrogate failure and follow-up scope

The first native suite completes all original 504 queries with no errors,
overtimes or new losing-mate flags. Its fully finite mean loss is 40.0448cp for
v0.21 and 45.9552cp for the candidate, above the frozen 45.0448cp limit. Severe
cases rise 1 to 2 (within the separate tolerance). Preserve `passed: false` and
the failed native-driver exit; do not rerun trials to obtain a favorable result.

Only five roots select different moves across the two actors/trials. All five,
including improvements, are retained in `changed-choice-inputs.json`. The largest
loss difference is HhPI3vCx: Qd3 to Qd4, 455cp worse under the fixed teacher, with
completed depths four/five observed. Other changes include a 42cp loss, two
small losses and an improvement. These are inspected regression cases, not a
blind sample; the teacher is not presumed to identify a unique acceptable move.

Before deciding whether to reject or justify a research game screen, compare
unrestricted eager/candidate iterative searches at **every observed completed
depth** for all five changed roots. Require exact scores/moves/nodes to distinguish
a staged-order defect from the baseline's own depth-dependent horizon choice.
Do not force Qd3/Qd4 or modify runtime behavior for a particular root.

Complete the other three previously unexecuted native suites as **diagnostics**,
retaining every pass/failure. This is not a retry of the failed 100ms block or
a claim that its gate passed. No game screen or promotion follows automatically.
A subsequent explicit review must state whether the isolated surrogate result
justifies rejection or a bounded research screen in light of exact correctness,
larger efficiency and the remaining budgets. Any original failed gate remains
visible even if a later decision uses other evidence. Timing/mate/legality
failures remain hard blockers.
