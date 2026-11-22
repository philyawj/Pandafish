# Audit and assessment after released v0.23

This directory contains bounded unchanged-baseline diagnostics, not an engine
candidate or a strength experiment. See the
[repository audit](../../../../docs/REPOSITORY_AUDIT_V023.md) and
[eight ranked engine projects](../../../../docs/NEXT_ENGINE_ASSESSMENT_V023.md).

Baseline HEAD/local remote-tracking ref/tag:
`28948e7473251d950003ec1f966f0cf26369b817`. GitHub independently confirms the
published normal release, merged PR and current remote main in
[remote-state.json](remote-state.json). The source checkout has no engine/library
changes relative to the tag. Reports and this evidence are new uncommitted files.

## Repository measurement

[The Go diagnostic](repository_audit.go.txt) reads `git ls-files`, file bytes and
Git release differences. It records tracked sizes, SHA-256 exact duplicates,
snapshot payload repetition, missing inline-Markdown local paths, and lossless
gzip best-speed sizes using counting writers. It never changes tracked inputs
or creates compressed copies. It is an archived one-off diagnostic, not a new
production tool or a template for engine orchestration.

To repeat from the declared baseline in a scratch checkout with Go 1.25:

```sh
cp tools/benchmark/research/post-v0.23-audit/repository_audit.go.txt /tmp/pandafish-repository-audit.go
GOCACHE=/tmp/pandafish-v016-go-cache go run /tmp/pandafish-repository-audit.go /tmp/pandafish-repository-metrics.json
```

Copy this diagnostic from the assessment first if the scratch checkout contains
only the original tag. The measured tracked set is HEAD's checkout; after new
files are committed, totals will naturally differ. The release-only comparison
starts at pre-release main `8b0109d`, while the v0.22 comparison also includes
the intervening live-review PR. Compression is per-file, including already
compressed files, with no tar overhead. Counts use decimal bytes/MB; Git object
output separately uses MiB.

[repository-metrics.json](repository-metrics.json) and
[metrics-log.txt](metrics-log.txt) retain the results. The initial diagnostic
only measured duplicate/size/link information; the final diagnostic additionally
measures compression and distinguishes release-only changes. This is a method
extension, not a replacement of failed engine evidence.

The path check covers inline links and path existence, including repository-root
links. It excludes external URLs and heading-fragment validation. The five
missing paths are classified in the audit. New reports have a separate document
validation record. [Git object storage](git-object-storage.txt) is read-only
`git count-objects -vH`, not a remote pack or network-transfer measurement.

[Document validation](document-validation.json) checks the three new reports'
local inline links, six tables, heading separation, whitespace and eight
ranked candidates. Its [Go checker](document_check.go.txt) can be copied to a
scratch `.go` path and run from repository root. The existing `wording` command
passes on all three documents; the source-policy test also passes after adding
the diagnostic sources as Go text data. `git diff --check` and the explicit
engine/library comparison to v0.23 pass. All existing tracked files are unchanged.

## Current CPU profile

The existing `research_warm_test.go` has no per-root `ReadMemStats` calls. It
loads/replays all 20 frozen games, verifies their openings, reconstructs supplied
history and keeps per-game tables. It searches 1,145 candidate turns at complete
depth five. No new position collection, teacher queries or games occur.

Build and run from repository root:

```sh
GOCACHE=/tmp/pandafish-v016-go-cache go test -c -o /tmp/pandafish-v023-audit.test .
WARM_INPUT=tools/benchmark/research/v0.21-rank-order/screen/v0.21-rank-order-dev-vs-v0.20-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba \
WARM_OUTPUT=/tmp/pandafish-audit-profile-roots.json WARM_MODE=search WARM_DEPTH=5 \
/tmp/pandafish-v023-audit.test -test.run=^TestResearchWarmWorkload$ \
  -test.count=1 -test.v -test.cpuprofile=/tmp/pandafish-audit-cpu.pb.gz
```

The recorded run used this directory's `profile-roots.json` and
`baseline-cpu.pb.gz` output paths and redirected its console to
[profile-log.txt](profile-log.txt). Test binary SHA-256:
`b4ab35a49aec4477a01bbeda88ca4f6dd219be6c5c7e949bff7248375a9ed9e2`.
The ordinary tagged executable remains available with SHA-256
`88d78500a8996471322d8e282f69d704ab9aceca386ac1693cf2aca82d77e6be`.
Both build identities are Go 1.25.0, darwin/arm64. Build paths can affect hashes.

The profile duration is 18.96s, with 16.97s CPU samples. All roots, scores, moves
and 43,118,858 nodes match the original candidate through the existing tool:

```sh
go run ./tools/research warm-report \
  --baseline tools/benchmark/research/v0.23-state-kernel/warm/warm-candidate-pair-1.json \
  --candidate tools/benchmark/research/post-v0.23-audit/profile-roots.json \
  --out /tmp/pandafish-audit-exactness.json \
  --pairs 1 --roots 1145 --depth 5 --exact
```

The recorded comparison is [profile-exactness.json](profile-exactness.json).
Its timing ratio has no efficiency gate and is not used as speedup evidence:
one profiled diagnostic and an old unprofiled timing are different protocols.
CPU subtrees overlap; their cumulative percentages cannot be added.

`go tool pprof` failed twice with `go: no such tool "pprof"` before analysis.
The installed Go 1.25 module toolchain lacks that executable. The available
`/opt/homebrew/Cellar/go/1.24.5/libexec/pkg/tool/darwin_arm64/pprof` successfully
read the profile with `-top -nodecount=50` and `-top -cum -nodecount=50`, producing
[flat](cpu-flat.txt) and [cumulative](cpu-cumulative.txt) reports. The analyzer
version does not change the engine's Go 1.25 build. No profile was rerun.

Apple M2 / eight logical CPUs, AC reported by `pmset -g batt`; it simultaneously
reported a discharging battery estimate, so no stronger power-stability claim is
made. The sandbox initially denied `ps`; an approved read-only process snapshot
showed background system/browser work. No deliberate tests/builds/matches
overlapped the profile. Thermal state and continuous system load were unmeasured.

## Validation boundary

The audit ran root `go test ./...` and explicit nested chess `go test ./...`;
both passed, with some cached results. Vet passed in both modules. Root tests
reported engine 6.922s and research-tool 2.647s; other package results were cached.
The default package list contains five intended root-module packages. That list
emitted a nonfatal sandbox warning about writing module-stat cache metadata;
the tests and package output completed. No dependency/tool installation occurred.

[prior-evidence-audit.json](prior-evidence-audit.json) verifies all 56 existing
local inventories /3,710 entries, with the historical release asset explicitly
unchecked. The new directory is sealed separately after report verification.
No existing inventory, raw observation, source snapshot or failed gate changed.
No full race/release smoke or external wrapper was rerun for this documentation
assessment. The original v0.23 validation remains available in its release record.
