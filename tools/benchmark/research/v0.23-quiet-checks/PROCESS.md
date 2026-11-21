# Method, source and reproduction

All work is offline. Baseline source is tagged v0.22 at
`1953f2af97f956c4499cbc5b61c3c1c698856740`; research branch starts from the later
documentation merge `8b0109d5fd578a397e2f7ae9857190bd94df8095`. The executable
and development labels identify actual tested artifacts; no release is implied.

## Source and independent controls

The full source is in `candidate-source.json`, SHA-256
`6a33e02dc18c5bce42fa92a843093054bf4d2482e6d7131d83cdd24a066b401b`.
Shared Go replay accepts `candidate-replay-source.json`. Its source-format proof
omits only unchanged nested module metadata already supplied by the base tag.
All runtime and test source is identical. The original refused replay is kept.

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.23-quiet-checks/candidate-replay-source.json \
  --out /tmp/pandafish-new-quiet-check-replay --check
```

Go normal/race and library race tests, vet, Linux build and clean-overlay checks
are separate from playing strength. `quiet_checks_test.go` compares the API with
596 frozen independent move/tag records and original complete legal generation,
plus explicit and seeded cases. Cold and cached lists, complete order, caller
prefixes and parent state are checked. FEN construction checks status without
populating the full move cache, so the cold generator/cost cases are uncached.
Positive/negative search controls cover quiet mates, escapes, captures, blocks
and the absence of a second quiet-check layer. Original capture-only reference
failures and the policy-specific independent-minimax update remain documented.

## Equal-time choice controls

The preserved assessment drivers run the existing Go `research choices` tool,
with tagged baseline and frozen candidate executables. Public roots receive
three 100ms trials and one 1000ms trial; newer roots receive one at each budget.
Three native requests run concurrently. Native work completes before referee
queries. These four suites contain 1,252 native actor requests in total.

The teacher is pinned Stockfish 19 at 100,000 nodes, Threads 1, Hash 64 MiB,
Skill 20, strength limiting off and MultiPV 1. Query histories, selected moves,
scores, bounds, PVs, cache provenance and raw output remain saved. Each reference
root/played-move request starts with fresh state/hash. Existing cached queries
are reused only at matching teacher hash and budget. This accepts any comparably
good move; no production root is forced to a teacher choice. CP losses use the
established 25cp uncertainty allowance and 1000cp cap; mates remain separate.

Existing mean tolerance is baseline plus max(10%,5cp), severe-case count at most
baseline plus two, no newly selected losing mate under the tool's finite-root
convention and no response past allowance plus 50ms. Newer/fresh turns also
pass separately. Small CP changes are estimates, not Elo or statistical proof.

## Fresh final controls

The frozen seed 2026112023 collector uses released v0.22 and Stockfish only,
not this candidate. Sixteen families and adjacent-turn sample plies are fixed
before generation. All raw queries, exclusions and paths remain saved. There
are 147 retained roots, 74 White and 73 Black; previously inspected canonical
boards and color mirrors are blocked, including all v0.22 live survey boards.
No fitting or candidate-dependent extension uses these families.

The changed generator is preserved in
`../v0.23-next-assessment/provenance/generate_fresh.go.txt`; its other Go files
come unchanged from `tools/research/` at the research base revision. Build the
temporary collector inside this module, under ignored `tmp/v023-generator`,
because the existing choice helpers import `internal/telemetry`. The first
failed compilation outside the module is retained. The actual invocation and
generation settings are in the logs and `fresh-generation/method.json`.

Fresh choices use three trials at 100ms and one at 1000ms, with both-turn gates.
That adds 1,176 native actor requests. Candidate source and binary are unchanged
between existing and fresh suites. These controls are inspected once; they are
development evidence afterward, not reusable hidden holdouts or an Elo test.

## Standard strength screen

`SCREEN_DECISION.md` precedes the only match screen. `screen_driver.go.txt`
checks all six quality gates and native binary identities before invoking the
unchanged Go benchmark. The exact command, start/end status and complete output
are separate JSON/text records. The profile is 50 games each against tagged
v0.22, Stockfish Skill 3 and pinned Sacrifice b528f0ba, all concurrent, 100ms
plus 50ms grace, ten standard openings and the 400-ply cap. Teacher Skill 20
analysis is distinct from the Skill 3 playing opponent.

No heavy agent validation overlaps timed choices or match play. The host is
an Apple M2 MacBook Air on AC power. External user activity is unmeasured;
privileged OS process/hardware queries are unavailable in the sandbox. Earlier
profile extraction uses Homebrew's installed `pprof` binary because the selected
Go toolchain lacks `go tool pprof`; those profiles remain diagnostic inputs in
the assessment directory, not match or unprofiled timing evidence.

Every failed implementation, infrastructure failure and adverse outcome remains
preserved. The rejected evaluator and child cells are inactive. No live play,
commit, PR, release or automatic larger match confirmation follows this work.
