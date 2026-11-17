# Go research tools

I use `go run ./tools/research` from the repository root. This program needs
Go and Git; its checks do not install interpreters, collect fresh engine labels,
refit models or start games. The normal [benchmark](../benchmark/README.md) and
[session tools](../lichess/README.md) remain Go programs too.

```sh
go run ./tools/research audit
go run ./tools/research decisions
go run ./tools/research mates
go run ./tools/research policy --base v0.18 --version v0.19
```

`audit` verifies every research inventory. Original raw observations retain
their original hashes. The explicit [migration manifest](migration.json)
accounts for retired scripts and updated documentation without rewriting old
inventories. Missing or altered data fails. To verify every byte in the original
release attachment too, download it from [v0.19](https://github.com/philyawj/Pandafish/releases/tag/v0.19)
and pass `audit --legacy /absolute/path/to/v0.19-historical-research-evidence.tar.gz`.
The command reports whether the external evidence was actually checked.

`decisions` independently recomputes all baseline/candidate integer selections,
ties, CP-only loss, mate-risk and per-turn metrics for the sharp-branch attempt.
It checks them against the frozen results and fails on disagreement. It does
not train or optimize weights. `mates` checks every saved positive/negative
training proof record against the Go rules library, preserving the distinction
between an actual legal mate and an opponent threat after a hypothetical pass.
`policy` compares original engine/library files with the declared tag, allowing
only the UCI version substitution. `games --dir PATH` legally replays a benchmark
summary's PGNs and checks their saved results/lengths and mate/repetition claims.

I replay an archived Go candidate into a **new** directory with:

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.19-independent-families/candidate-source.json \
  --out /tmp/pandafish-replayed-candidate --check
```

The command extracts the declared base tag, overlays Go/module files and JSON fixtures,
then optionally tests and builds it. Existing destinations and unsafe paths
are rejected. Failed tests/builds retain their logs. No matches run.
For unchanged tagged code plus a saved Go proof test:

```sh
QUIET_PROOF_INPUT="$PWD/tools/benchmark/research/v0.19-sharp-branches/training-quiet-threat-scan.json" \
go run ./tools/research replay --base v0.18 \
  --test tools/benchmark/research/v0.19-sharp-branches/quiet_mate_proof_test.go.txt \
  --out /tmp/pandafish-quiet-proof --check
```

I check PR descriptions and release notes before publishing:

```sh
go run ./tools/research wording --file /tmp/pandafish-pr.md
go run ./tools/research wording --file /tmp/pandafish-release.md
```

I retired 131 historical Python/shell source files, including source stored as
`.txt`, rather than hide their language from GitHub. The hashed release
attachment preserves 1,148 original files and its own manifest. The three
post-v0.18 attempts were still local when archived; their scripts therefore
need that attachment, while older scripts also remain in earlier tags.

I replaced shared archive verification, Go-engine restoration, frozen decision
checks, mate-proof checks and the shell test fixture with Go. One-off old corpus
collectors, NumPy optimizers, external python-chess oracle drivers and historical
batch orchestrators are **retired**, not claimed as equivalent Go ports. Their
original code/results remain reproducible in the historical environment from
the attachment. I will implement future collection/fitting/oracles in Go under
a separately reviewed experiment design. I have not reproduced a new independent
rules oracle by merely calling the same Go library again.

I retain documentary shell command examples for invoking Go and Git. The optional
lichess-bot client is installed outside this repository and has its own Python
requirements; I do not vendor or rewrite that external service client here.
All repository program/test/research source is Go. The root source-policy test
rejects other source extensions, including archived script-name disguises.
