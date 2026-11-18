# Rejected monotonic deadline optimization

The [selection](../v0.22-next-assessment/SELECTION_MONOTONIC.md) isolates a public
Go API change: read an existing monotonic deadline through time.Until rather than
read wall plus monotonic time for Time.Before. Every check site, allowance,
injected clock and cancellation guard stays intact. All two-file source overlays
and hashes over v0.21 are retained; no earlier failed code is bundled.

Corrected deadline controls and full race/vet/build pass. The initial zero-time
fixture failure remains preserved in [ERRATA.md](ERRATA.md), its original Go
source and raw failed log. All 325 distinct fixed-depth score/move/node
comparisons agree. Three depth-four ratios are 0.93045, 0.92912 and 0.93135,
median **6.95% less time**. Depth-six is 0.92772; newer depth-four is 0.89227.
The improvement is functional, but misses the unchanged 15% gate: the timing
driver exits 1 with `timing_gate_passed: false`. No native decision probes,
candidate warm workload or strength screen follows. The prepared warm helper
and driver are unexecuted source, not measurement evidence.

Archive without promotion. Reconsider this small, exact optimization under a
future explicitly scoped efficiency goal; do not relabel it as a larger gain or
relax this investigation's gate after observing its result. Current work pivots
to a search-ordering hypothesis, with a fresh cutoff-opportunity diagnosis before
implementation. The previous sampled walltime percentage overstated the useful
gain on these unprofiled fixed-work cohorts.

Restore the two-file JSON through `tools/research replay --source` into a new
scratch directory. Time.Until's inspected Go 1.25 implementation uses runtimeNano
for monotonic deadlines and a Time.Sub/Now fallback otherwise; this is an API-level
optimization, not a standard-library or runtime patch. No upstream/library file
changes in this attempt.

After the subsequent larger-work attempts also fail, a separately documented
[small-scope reassessment](../v0.22-clock-screen/PLAN.md) evaluates this exact
source under a different implementation-cost objective. The original 15% result
and fixture failure remain unchanged; the later phase must establish its own
previously unexecuted timed-warm, native and game gates. This initial archive's
rejection concerns the larger-gain claim, not a claim of failed correctness.
