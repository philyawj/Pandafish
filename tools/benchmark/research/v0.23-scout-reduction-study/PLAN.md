# Diagnose late child scouts, including wide-window parents

The old v0.22 guarded-LMR diagnosis remains rejected on its original cohorts,
guards and 30% criterion. It required the original parent window to be non-PV.
The released PVS loop also sends later children of wide-window parents through
an individual zero-window scout. That different child/parent boundary can expose
more work; it also introduces a different risk of overlooking a PV improvement.
This is not a retrospective pass or a lower threshold for the old study.

Measure unchanged tagged-v0.22 iterative search over all 1,145 warm roots and
ten openings. Eligible scouts: parent depth >=3, fourth/later move, ply >=1,
parent and move nonchecking, ordinary quiet move, no killer hint, and finite
scout bounds outside the mate band. Do not require a narrow original parent
window; every counted child still has an actual one-centipawn PVS scout.
Record non-overlapping outer scout node spans only, closing before any original
full-window re-search. Do not count the re-search as saved work. Report wide-
parent and existing narrow-parent categories separately where possible.

Require all baseline moves/scores/nodes exact, >=30% eligible outer node coverage
overall and >=20% for each root turn. Instrumented elapsed time is not performance
evidence. If insufficient, archive without reductions. If promising, freeze one-
ply reductions with full-depth null-window verification on every alpha improvement
or mate-band result; only then retain ordinary full-window PVS re-search rules.
Never reduce checks/captures/promotions/killers/first moves or change draw/cache
identity. Fixed-depth equality is not promised for that selective implementation.
Correctness, a frozen cost test, all native/fresh quality checks and the unchanged
single concurrent 150-game screen remain necessary before a recommendation.
