# Rejected reusable child-storage attempt

Decision: archive without promotion. The implementation is correct on the saved
controls but fails the frozen 15% median fixed-work elapsed-reduction gate.
Baseline is tagged v0.21; no native decision suite or strength games run for this
attempt. No weights, search/TT policy, clocks or released identity change.

The [selection](../v0.22-next-assessment/SELECTION.md) diagnoses child-state heap
allocation after profiling 1,145 recorded roots with a persistent TT. An additive
`SearchPosition` cell holds a Position and Board per active ply. It copies parent
state, resets every field/cache and rejects aliasing; completed siblings reuse
the cell. Original immutable methods remain comparison oracles. Five exact Go
overlays and their hashes are preserved in `candidate-source.json` and
`source-hashes.json`; the normal engine source is restored before another attempt.

Controls pass: 38,617 scratch transitions against both immutable APIs, 37,795
unchanged transition controls, independent frozen child FEN/hash/tags and perft
on 596 positions, startpos depth-four perft, cache reuse, alias rejection,
ancestor/worker isolation, cancellation/history cleanup, full race/vet/build,
and original pinned upstream root-package race tests. All 21 original upstream
runtime files are byte-identical. Logs and the archived Go drivers preserve the
commands and outputs. No failed correctness assertion is discarded.

All 325 distinct public/deep/newer fixed-depth comparisons retain identical
scores, chosen moves and nodes. Three alternating public depth-four ratios are
0.96683, 0.98379 and 0.95624; median 0.96683 is only **3.3% less time**. Depth-six
and newer depth-four ratios are 0.94810/0.94826. `timing-analysis.json` explicitly
records `timing_gate_passed: false`, and the driver exits 1. That failed gate is
not weakened or rerun for a favorable measurement.

One diagnostic fixed-depth warm-table pipeline pair also retains every result:
1,145 roots / 9,874,565 nodes agree exactly. Search time falls 10,218.1 to
9,742.9ms (**4.65%**); measured allocated bytes fall 3,385,907,320 to 700,406,896
(**79.3%**, including position reconstruction/setup). Allocations fall 22,902,612
to 3,220,584. This is descriptive corroboration, not a substitute promotion gate.
The warm driver's `--fixed-only` mode is used; no candidate 100ms profile runs.
Per-root MemStats/output overhead and selected recorded positions limit these
measurements. Cumulative allocation is not peak memory or CPU cost.

The hypothesis that removing these allocations alone provides a larger elapsed
gain fails. Allocation dominance in bytes did not establish CPU dominance. Keep
this code inactive for a future independently scoped memory objective or as a
possible follow-up if changing bottlenecks justify its complexity. Current work
reassesses geometry-based move-tag/legality work separately; the storage change
is not bundled into it. No measured playing-strength claim follows.

Restore to a new scratch destination:

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.22-child-storage/candidate-source.json \
  --out /tmp/pandafish-child-storage-replay --check
```

The oracle test and warm-profile helper are retained as Go text under this
archive and the assessment directory. Their saved driver paths refer to the
original workspaces; use new destinations/output paths for reproduction.
