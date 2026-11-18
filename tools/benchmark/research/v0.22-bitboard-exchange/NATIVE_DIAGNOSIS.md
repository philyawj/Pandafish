# Public 100ms prerequisite failure

The original 504 actor queries complete without errors, illegal moves, new
losing-mate flags or timing overruns. Fully finite mean loss is 40.0746cp for
v0.21 and 45.8955cp for the candidate, above the unchanged 45.0746cp limit.
Severe cases rise one to two within that separate tolerance. Keep
`public-100/quality.json` false and the native driver's failed exit. Do not rerun
these actor requests or alter a root move/threshold to obtain a pass.

As with the earlier independently recorded staged-picker failure, inspect **all
changed-choice roots**, including improvements, at every observed completed
depth using unrestricted fresh iterative searches of both actors. Require exact
score/move/node agreement on those matched-depth comparisons. This diagnosis
can identify a pre-existing horizon choice; it cannot make a worse choice good
or prove all future positions equivalent.

Complete the other three previously unexecuted native suites once as diagnostics,
preserving every pass/failure. This is not a retry or an automatic game decision.
A subsequent explicit review must weigh all original failures, exactness and
the 35.35% cold / 37.84% total warm efficiency evidence before any bounded game
exception. A prior exception on another binary does not authorize this one.
Legality, query, timing and new losing-mate failures remain hard blockers. No
screen, larger confirmation or publication follows automatically.
