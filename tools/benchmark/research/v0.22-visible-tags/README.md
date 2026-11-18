# Rejected occupied-ray tag revision

This [separate revision](../v0.22-next-assessment/SELECTION_VISIBLE.md) preserves
all 325 fixed-depth results and passes the same 100,650 pseudo-legal tag controls,
596 independent saved vectors, race/vet/build and original upstream root-package
race tests. All original runtime files remain intact. Source and exact hashes
are frozen in `candidate-source.json` / `source-hashes.json` over v0.21.

It fails the unchanged 15% performance gate. Three depth-four ratios are 0.96651,
0.96902 and 0.95970 (median **3.35% less time**). Depth-six is 1.00104, effectively
flat/slightly slower; newer depth-four is 0.96863. The timing driver exits 1 and
retains `timing_gate_passed: false`. No native decision suite, candidate warm
workload or strength screen follows. Unexecuted warm/profile drivers remain
source only. No precise strength inference is available.

Visible king-ray blockers avoid more impossible checks than full rays, but this
does not produce a sufficient net elapsed gain. Preserve the code as inactive
research rather than add another complexity layer to the playing engine. The
next reassessment isolates the cost of reading both wall and monotonic clocks
for an existing monotonic deadline, without reducing check frequency or changing
time management. The failed storage/tag versions are not bundled into it.

Restore with `tools/research replay --source` and this archive's JSON in a new
scratch destination. Driver paths identify the original scratch directories;
adapt them and use new output paths when reproducing.
