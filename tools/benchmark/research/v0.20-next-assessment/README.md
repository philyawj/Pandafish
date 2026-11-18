# Post-v0.19 assessment and rejected quiet-mate proof

This directory preserves the first attempt and shared development controls for
the next candidate. The base is released v0.19, commit
`e7bf00fe5b158e3d58ca77e83385ed918f839717`, whose playing policy is v0.18.
No release or live game is part of this investigation.

Read the initial [assessment](RESEARCH.md), frozen [choice plan](CHOICE_PLAN.md),
[quiet-proof result](QUIET_RESULT.md), and [preserved findings](ERRATA.md).
The [review findings](REVIEW_FINDINGS.md) retain utility/schema/format corrections
alongside the raw experiment failures.
The quiet-mate implementation is **rejected and absent from active play**.
The [transition-only attempt](../v0.20-search-transitions/RESULT.md) and
[search-legality revision](../v0.20-search-legality/README.md) follow separately.
The [combined assessment](../../../../docs/RESEARCH_V020.md) compares these results
with the earlier research and ranks further work.

## Evidence and identities

`quiet-candidate-source.json` preserves the initial timed v1 source;
`quiet-v1-failing-control-source.json` adds the discovered API regression control;
`quiet-corrected-source.json` preserves the corrected, tested v2. Corresponding
hash maps identify every overlaid Go file. All overlays declare v0.19 as the
base and can be restored with the Go research `replay` command. No archived
quiet-proof source is compiled into the current engine.

The baseline binary SHA-256 is
`fea4b1a87868b3c05fd5ac4f46aa52ab9afb745e259b967c7c469692cda8845c`.
Native method files preserve the exact candidate, input and teacher hashes.
CPU profiles, initial fixed-depth timings, all native/reference observations,
test logs and the EP defect remain retained. A corrected control is not a
successful rerun of a rejected strength candidate.

`fresh-generation/` contains every query, complete rollout history, retained and
excluded sample, input and method record for the 145 new development roots.
There are 73 White and 72 Black roots from 16 families. Prior sharp-batch paths,
blocked states and sample color/rank mirrors were excluded. These controls were
frozen before probing quiet v1 and reused for later backend comparisons. Once
inspected, they are regression controls, not a blind holdout or independent games.
No opening book, online game collection or learned values are introduced.

`choice-review.json` independently recomputes every completed choice suite in
all three attempt directories. The [Go review source](review_choices.go.txt)
checks legal histories/moves, native raw output, reference output or recursively
preserved cache observations, finite losses, mate-risk flags and time gates.
Legacy teacher caches preserve structured iterations rather than raw UCI lines;
their last exact iteration and source hashes are checked. Uncapturable EP FEN
fields are normalized only when no legal EP move exists. The original quiet
driver reported aggregate acceptance only; the separately frozen per-turn
criterion still rejects its Black regression. Original reports are not rewritten.

## Replay and driver boundary

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.20-next-assessment/quiet-corrected-source.json \
  --out /tmp/pandafish-quiet-replay --check
```

Use a new destination. Replay tests/builds the candidate without running games.
The reachable EP regression fixture is [saved separately](legal_ep_regression_test.go.txt)
for both v1/v2 restoration tests. The root [fixed-depth driver](probe_test.go.txt)
is a Go test with explicit output/input/depth environment variables. Ordinary
tests skip its timing workload. Terminal direct-root probes compare the same
internal return on both sides; production null/terminal behavior has separate
native and unit controls.

[Driver sources](drivers/) preserve the Go orchestration used for snapshots,
fixed-depth comparisons and upstream replay. Copy a `.go.txt` source to a `.go`
path before invoking it from the repository module. Original absolute temporary
paths are retained for provenance; adapt paths and use new result directories
for another reproduction. Timing drivers refuse existing JSON outputs, while
the original snapshot/upstream helpers can overwrite their output targets: use
scratch copies, never the evidence directory. The scoped legality-upstream
driver is the corrected command; the initial full-module failure and the
explicit successful root-runtime command remain documented separately.

No Python, shell or other language source is added. Documentary shell command
examples invoke Go; all implementation, tests, fixtures and research drivers
are Go. Each attempt has a byte/hash inventory, excluding only its inventory
itself. Original unsuccessful evidence stays intact.
