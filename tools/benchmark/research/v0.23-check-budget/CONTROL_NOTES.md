# Control development

Initial normal tests pass in 173.483 seconds. The complete-list, fixed-color
minimax oracle now includes the new checking horizon, independently of the
specialized generator and negamax implementation. Its larger unpruned tree
substantially raises validation cost. Preserve the initial validation logs,
including any race timeout, before repairing test execution cost.

The next focused control requires two quiet checking choices separated by a
forced evasion, in both colors, without prescribing a root move. A shallower
one-choice horizon must not prove that mate. Existing escape, capture and block
controls remain. No playing candidate changes are implied by a test-cost repair.

The first proposed two-choice fixture was incorrect: Kf7/Qb1 versus Kh8
already has mate in one (Qh1). The focused test correctly failed with score
999999 instead of expected 999997. A deterministic Go fixture finder then
identifies Kf6/Qa1 versus Kh8, with its color mirror. Complete-list minimax
independently verifies the resulting mate band; the shallow horizon remains
finite. The initial fixture source and finder remain as development evidence.

The memoized normal run passes in 39.884s, but its race run also reaches the
generic 300s wrapper limit without an assertion result. Preserve both failed
run states. A separate expanded validation plan grants this exhaustive reference
up to 12 minutes inside Go and 15 minutes in the runner, with JSON test events
for diagnosis. Source, assertions, scientific cost/choice gates and candidate
are unchanged. This is a test-runtime allowance, not a favorable measurement
retry. Performance measurements still cannot start until race/vet/build pass.

Expanded validation passes: normal 39.222s, race 319.928s, vet and both builds.
The JSON events distinguish the independent depth-seven mate test from the
larger start-position exhaustive comparison; both complete with their original
assertions. No production policy is relaxed to shorten validation. The final
117-file Go snapshot is `candidate-source.json`; the same-horizon reference is
replayed separately from that exact source.
