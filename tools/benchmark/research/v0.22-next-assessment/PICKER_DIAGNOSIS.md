# Reassess deferred main-search generation

The clock's small functional objective passes but its game screen fails (46%
versus v0.21). Restore v0.21 and keep every failure. Revisit the earlier backend
plan's staged move picker, which remains unimplemented: even when a table hint
or tactical/killer move refutes a node immediately, main search currently builds
and tags the entire legal list first. Earlier sort/temporary-board optimizations
reduce that work's cost but do not remove the unused generation itself.

Instrument the frozen 20-game depth-five baseline trajectory. For each main-search
generation, count full/quiet move records and capture its existing TT and killer
priorities before searching. On a cutoff, record whether it occurs at the table-
preferred first move, a capture/promotion or an existing quiet killer. Those
stages precede ordinary quiet moves in the current exact order. Count quiet
records whose generation could be deferred entirely; table-first cutoffs can
also avoid the tactical list. No score, ordering or returned node is changed.
Instrumented elapsed times are not performance evidence.

If at least half of generated main-search quiet records belong to these early-
stage cutoffs, freeze one staged-picker hypothesis: validated TT hint first,
existing stable tactics, validated quiet killers, then remaining quiet moves in
their original order. Add direct legal hint lookup and quiet-only generation
only as required by that one pipeline. Preserve every score, root choice and node
at fixed work; no killer-before-capture policy, selective pruning, clock change,
new cache capacity or failed earlier kernel is included. Otherwise reassess from
the data rather than assume that high first-move cutoff counts alone suffice.
