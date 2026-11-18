# Signed quiet-history revision

Baseline is tagged v0.21. Follow the [frozen signed-feedback policy](../v0.22-next-assessment/SELECTION_MALUS.md)
and original [history selection/gates](../v0.22-next-assessment/SELECTION_HISTORY.md).
The positive-only and four exact-kernel failures stay archived and inactive.
Before this revision's measurements, the [history workload correction](../v0.22-next-assessment/HISTORY_METHOD_CORRECTION.md)
replaces the cold-only efficiency gate with three warm iterative pairs and keeps
the 15% target. Cold score/membership and no-regression controls remain mandatory.
Preserve every failure. No larger confirmation, live game, commit, push or PR
follows automatically; check in at a viable candidate.
