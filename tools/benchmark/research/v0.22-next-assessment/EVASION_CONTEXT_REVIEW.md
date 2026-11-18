# Reassess mask-context timing after the warm failure

The first evasion revision fails the unchanged 10% extra warm target at 9.11%,
despite 24.34% total cold savings and exact results. No native/games run. Keep its
eight-file source, all three pairs and failed driver output unchanged.

Source review identifies avoidable context construction. `searchMoves` calls
the checker-mask helper before testing any king candidate. `HasLegalMove` uses
the same generator with `first:true`; it commonly returns on a legal king escape,
so that context is never used. Unchecked callers also enter the non-inlined
helper just to return inactive. The diagnostic already shows that existence
queries account for only 1,521,002 avoidable tags out of 56,645,308 tags, whereas
full generation provides the substantial opportunity.

Profile the frozen first evasion revision once on the same warm depth-five
trajectory. Preserve all rows and compare exact results; profile percentages
include setup/MemStats/runtime overhead and do not establish elapsed gains.
If the cost is meaningful, defer construction until the first occupied nonking
piece needs a mask. A king-first existence success then returns without creating
context; unchecked requests bypass it entirely. Later nonking/EP/fallback
semantics must remain identical. This does not disable evasion filtering or
change legal-existence order.

Any revision is frozen separately. It must still meet the original 10% warm
extra gain versus the first staged picker, 15% total target, all exact controls
and native gates before a game decision. The first evasion failure is not a
passing run or a strength test. No threshold reduction or unchanged-code rerun
is authorized by this review.

Outcome: the full profile attributes only 0.12s / 22.10s (**0.54%**) sampled CPU
to `searchEvasionMask`. The bounded context change has insufficient larger-gain
headroom; no playing context revision is implemented. Its raw profile, complete
rows and all-node report remain. Reassess the exchange evaluator, whose repeated
square scanning accounts for 14.98% cumulative sampled CPU in the same run.
