# Validation and reproduction

All repository source, research programs and test fixtures are Go. Tests target
the changed runtime and its required ownership/rule boundaries. Original library
methods remain comparison oracles; the new mailbox attack oracle independently
implements geometry without shared bitboards/rays/masks.

## Correctness

- `transition-controls.txt`: 37,795 immutable legal child transitions agree with
  original state, tags, mailbox, piece sets, hashes, status and descendant moves.
  Special fixtures cover both-side castles, EP including a pin, promotions,
  captures, mate, parent immutability, cold/warm generation and nil transitions.
- `attack-oracles.txt`: 128,000 square/color attack predicates agree with both
  the original implementation and the independent Go mailbox oracle.
- `race.txt`, `vet.txt`: full root module checks pass, including existing perft,
  move storage/order, mate/repetition, TT context and cancellation controls.
- `upstream-race.txt`, `upstream-check.json`: the initial original whole-module
  replay fails only in the unvendored image dependency; its root suite passes.
  `upstream-root-race.txt` separately records `go test -race .` passing in the
  restored full upstream workspace. The failure is preserved, not relabelled.

Native and Linux/ARM64 builds pass. The downloaded Go 1.25 compiler builds all
tested binaries. The installed local Go pprof executable only renders captured
profiles because the downloaded toolchain lacks that tool; it does not compile
the engine or supply timing results. Original runtime/license hash checks and
the final archive replay are separate evidence checks.

## Fixed work and native choices

`timing-checks.json` records every command/exit. Three alternating depth-four
pairs and depth-three/depth-six/new-position blocks are retained, alongside the
initial exploratory comparison. All 325 distinct position/depth combinations
agree in scores, chosen moves and search nodes. Repeated depth-four timings are
replicates, not more distinct positions. `timing-analysis.json` retains every
ratio and the frozen extra-gain criterion; no favorable outlier is selected.

Public 84-position 100ms requests use three trials per actor; public 1000ms
and both newer 145-position budgets use one. Each native block allows three
searches concurrently, alternates actor order and finishes before teacher work.
The newer inputs were frozen before the first quiet-proof probe and are now
inspected development controls. All legal moves are acceptable production
choices. Stockfish `searchmoves` is used only to label each actor's selected
move; the candidate search is unrestricted. The 25cp band accepts close
alternatives, and finite-reference/mate limitations remain explicit.

Each choice directory preserves `method.json`, native UCI output, teacher
queries/cache provenance and its quality report. The independent
[Go review program](../v0.20-next-assessment/review_choices.go.txt) recomputes
those reports and validates raw observations. Full-strength teacher settings
are 100,000 nodes, Threads 1, Hash 64 MiB, MultiPV 1, strength limiting disabled.
The standard game opponent is separately set to **Skill 2**.

## Reproduce in new destinations

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.20-search-legality/candidate-source.json \
  --out /tmp/pandafish-legality-replay --check
go run ./tools/research games --dir \
  tools/benchmark/research/v0.20-search-legality/screen/v0.20-search-legality-dev-vs-v0.19-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba
```

Build the baseline from tag v0.19, never from this working tree. Timing drivers
and the fixed-depth Go test are preserved in the
[shared assessment](../v0.20-next-assessment/README.md). They need new result
directories and adapted temporary workspace paths; never overwrite old evidence.
The shared tools guide documents `choices` and `generate` flags. Native timing
blocks and game screens must not overlap other engine probes on the host.

The routine screen uses `--pairs 25 --parallel-all`: 50 games each against the
tagged baseline, Skill 2 and pinned Sacrifice. All three matchups run concurrently
with the unchanged ten-opening suite, 100ms allowance, 50ms grace and 400-ply cap.
Results preserve development labels, binary/compiler identities and opening
hashes. Legal PGN replay checks integration/rules, not independent playing
strength. No live client is launched. No larger confirmation is automatic.

## Final review

`final-checks.json` and `final-checks-run.txt` record the final commands and exit
expectations. All actual checks pass; the archived quiet-v1 EP control exits 1
as explicitly expected, while corrected v2 exits 0. `screen-pgn-check.json`
records all 150 games legally replayed. `screen-review.json` independently
recomputes settings, all 75 color pairs, statistics, time counts and the practical
gate, retaining the failed predecessor in the same review.

`archive-replay.json`, `archive-replay-tests.txt` and `archive-replay-build.txt`
show the frozen candidate restored into a new scratch directory and tested/built.
`final-tests.txt`, `final-vet.txt`, `final-native-build.txt` and
`final-arm64-build.txt` pass on the completed working tree. Earlier full race and
upstream controls remain valid: playing Go source still matches the frozen
overlay. The final source/cache checks and inventories supplement that identity.

All UCI integration in this investigation uses the Go harness. The optional
external lichess-bot wrapper is not invoked; no client or live Lichess session
starts. This is a candidate check-in before a commit, not release validation of
a future renamed/tagged binary.
