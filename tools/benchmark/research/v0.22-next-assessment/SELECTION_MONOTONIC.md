# Reassessment: read the existing monotonic deadline directly

Storage reuse and both conservative tag versions fail the unchanged timing gate.
Their exact code, passing correctness controls and failed timings stay archived;
restore playing source to v0.21 before this attempt. Do not combine them.

The warm profiles attribute about 16.2%/21.4% cumulative CPU to time.Now, including
about 9.5%/15.3% runtime.walltime and 6.0%/5.3% runtime.nanotime. These are sampled,
overlapping selected-workload percentages with instrumentation/runtime overhead,
not an exact speedup ceiling. Every stopped check currently calls time.Now and
compares against a deadline built from the production start time. Go Time.Before
uses the embedded monotonic readings when both Time values have them; Go 1.25
time.Until directly reads the monotonic clock when its deadline has one, falling
back to time.Now subtraction otherwise. Inspect the installed standard-library
time.go as the implementation reference; this uses a public standard API.

Replace only the ordinary stopped deadline predicate with `time.Until(deadline)
<= 0`. Retain the injected clock's existing `!now().Before(deadline)` semantics
and every context/zero-deadline guard. Clock budgets, check frequency and sites,
soft iteration checks, cancellation and completed-result/fallback handling stay
unchanged. No interval batching, cached timer, goroutine timer flag or new clock
allocation is introduced. This is a cheaper observation of the same deadline,
not more available thinking time or permission to miss a deadline.

Controls cover zero/past/future monotonic and wall-only deadlines, far wall dates,
injected equality/before/after checks, clock call counts, canceled-context short
circuiting and the existing recursive deadline/interruption/UCI tests. Apply the
same 325 exact-tree comparisons and 15% median elapsed gate. For the warm exact
sequence use depth four plus an explicit 30s movetime allowance so every root
checks a deadline but completes normally; the old depth-only warm driver has no
deadline and cannot measure this mechanism. Freeze that helper before running.
Use both baseline/candidate with the same limits and unchanged budget formula.
The recorded 100ms warm pipeline and four native suites still use normal actual
budgets. Preserve all results and the same native/standard 150-game screen gates.
