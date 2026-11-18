# Reassess selective search after exact-work/ordering attempts

Seven variants now fail their applicable larger-gain gates. The older architecture
notes defer coherent selective search until root/cache/move coordination and
backend work are measured; v0.16–v0.21 supplied that foundation. The earlier
pre-v0.15 PVS/LMR experiment remains a failed reference, not code to reactivate.

Before implementing a new policy, measure the cost of late quiet zero-window
probes in isolated tagged v0.21 iterative searches: public deep depth-six roots
and all newer depth-five roots. Candidate guards are parent depth at least three,
fourth/later move, ply at least two, original parent non-PV window, neither parent
nor move checking, ordinary quiet move, no killer hint and finite non-mate window.
Record eligible probes and non-overlapping outer probe node cost so nested spans
are not double-counted. No timing claim uses instrumentation or its concurrency.

If outer eligible probes account for at least 30% of visited nodes in both
cohorts, freeze a one-ply reduction with full-depth verification on any alpha
improvement or mate-band result. This changes selective-search policy; fixed-
depth node/score equality is no longer a promised invariant. Legality, rules,
history, mate/cancellation controls and actual-budget safety remain mandatory,
with full-width selected-move regret reviewed independently before a game screen.
Otherwise reassess another larger hypothesis rather than tune these guards to
make this opportunity pass. No failed aspiration/history/kernel is bundled.
