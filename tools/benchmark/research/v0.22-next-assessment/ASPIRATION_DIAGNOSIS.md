# Aspiration opportunity diagnosis

Signed history fails its corrected warm gate despite fewer nodes. Restore v0.21
and reassess a different main-search mechanism: a narrow initial root window
around the previous completed iteration's score, followed by an ordinary full-
window retry on any bound failure. This retains full depth and avoids selective
move omissions/reductions. No history or prior kernel is bundled.

Before implementation, instrument only completed iteration scores in an isolated
tagged baseline. Use all public 84 and newer 145 roots at depth five, retaining
every depth/score, identity and completion. No timing claim uses instrumentation.
Record the fraction of finite adjacent depth-three-through-five score changes
inside ±50cp. If at least 60% on both cohorts fit strictly inside that fixed band,
freeze a one-probe/one-full-retry policy and appropriate iterative workload gates.
Otherwise reassess instead of tuning the band to make these observations pass.

The previous direct-root helper does not execute iterative root-window logic.
Correctness/performance controls for an aspiration candidate must invoke the
actual iterative entry point in both actors, plus keep the old full-root path as
an unchanged oracle. Require exact completed scores and baseline selected-best
membership, not identical nodes or a unique tied root move. Native quality,
cancellation/budget and the standard concurrent screen remain necessary.
