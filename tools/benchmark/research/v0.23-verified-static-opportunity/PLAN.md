# Tactical verification of static-margin candidates

The unverified static shortcut has20.808% outer descendant opportunity but hides
one original forced losing mate, so it is rejected and sealed. This follow-up
tests one changed mechanism: require a real quiescence search before any proposed
static-margin cut. Keep the predecessor's depth1–3, ply>=3, null-window, check,
mate-band, phase>=8 and150cp/ply static prerequisites. No LMR, extra quiet checks,
null moves, evaluator terms, new clocks or altered cache policy are active.

For a qualifying position, run released capture/evasion quiescence against
beta+150*depth, in the same path/history context. Pop the main node before q's
normal push and restore it afterward; a repeated position must not be counted
twice. Reject the proposed cut unless that actual q search also reaches the
threshold. This still cannot prove that a quiet tactic is absent.

Measure on the unchanged1,145-root tree. Save and reset instrumentation-only
verification node counts so the original released trajectory must still match
exactly. Charge every verification outside a successfully qualifying ancestor;
do not count nested avoided spans twice. Report raw outer saved descendants,
verification cost, net opportunity, false predictions and all hidden mate cases.
The prospective gate is net opportunity>=15% overall and>=10% each turn, no
hidden original forced losing mate, and exact released nodes/values/moves.
These are node proxies, not actual elapsed savings or playing strength.

One fixed diagnosis only; no threshold sweep, new framework or new data campaign.
If it fails, preserve evidence and reassess. If it passes, freeze a separate
implementation plan requiring correctness, realistic cost, selected-move regret,
existing/fresh actual-budget controls and the standard concurrent game screen.
No publication/live play; check in before committing a recommended candidate.
