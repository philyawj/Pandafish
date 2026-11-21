# Go research tools

Run individual commands with `go run ./tools/research` from the repository root.
Reusable workflow plans use the persistent built executable described in the
[tooling guide](../../docs/RESEARCH_TOOLING.md). This program needs
Go and Git. Verification/replay commands do not collect labels or start games;
the explicit `choices` and `generate` commands run offline UCI queries with
supplied engines. Explicit experimental `context-fit` and `context-targets`
commands execute their frozen offline fitting/label protocols, outside the
routine workflow. No command starts live games or publishes code.
The normal [benchmark](../benchmark/README.md) and
[session tools](../lichess/README.md) remain Go programs too.

For new engine investigations, start with the
[tooling map and flowchart](../../docs/RESEARCH_TOOLING.md). Shared `workflow`,
`validation-plan`, `warm-plan`, `warm-report`, `native-plan`, `snapshot` and `inventory` commands
replace duplicated per-attempt orchestration. `generate` accepts fixed seed/
stream and repeatable exclusion flags. `choices` also writes a compact
`brief.json`; raw records and scientific gates remain unchanged.

`validation-plan` orders normal/race/vet checks before test/engine builds.
`warm-plan` requires its passed source-matching state through `--preflight-state`.
`depth-review` reconstructs changed native choices at observed depths; its
`--flags-only` option selects saved new mate flags. `confirm-mates` checks every
such flag at a fixed larger reference budget using both actors' complete move
union and an unrestricted root. These diagnostics preserve the original failed
gates and do not repeat timed actor queries. Finite reference reports are estimates.

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

Replay an archived Go candidate into a **new** directory with:

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

Check PR descriptions and release notes before publishing:

```sh
go run ./tools/research wording --file /tmp/pandafish-pr.md
go run ./tools/research wording --file /tmp/pandafish-release.md
```

## Offline candidate choices and development inputs

The [post-v0.19 investigation](../../docs/RESEARCH_V020.md) adds two Go commands.
They require existing native baseline/candidate binaries and pinned Stockfish.
They refuse an existing output directory and preserve method hashes, raw UCI
observations and failures. No production root move is forced, and close
reference alternatives are accepted by a declared 25cp band.

```sh
go run ./tools/research choices \
  --baseline tools/benchmark/bin/pandafish-v0.19 \
  --candidate /absolute/path/to/candidate \
  --teacher /opt/homebrew/bin/stockfish \
  --out /tmp/pandafish-new-choice-run \
  --milliseconds 100 --trials 3 --workers 3
go run ./tools/research generate \
  --baseline tools/benchmark/bin/pandafish-v0.19 \
  --teacher /opt/homebrew/bin/stockfish \
  --out /tmp/pandafish-new-development-inputs
```

`choices` defaults to the existing 84 public roots. Use `--input PATH` for
another frozen input, `--by-turn` to require each root turn's tolerance, and
`--cache PATH` to reuse matching completed reference observations with pinned
teacher hash and 100,000-node budget. Timed native work finishes before teacher
queries. A completed quality rejection is saved as `passed:false`; it is not a
process failure. Inspect the report, not only the command's exit status.
Mate-risk diagnostics remain separate from fully finite CP loss. These are
development decision controls, not a strength benchmark or a proof of safety.

`generate` implements the fixed 16-family, seed-2026111728 rollout plan in the
[choice protocol](../benchmark/research/v0.20-next-assessment/CHOICE_PLAN.md).
It excludes sampled prior/duplicate/color-mirror states and reconstructs earlier
sharp-batch witness paths. It is a narrow offline development collector, not a
replacement for retired training pipelines or an independent holdout after use.
The canonical-state tests also compare 1,504 saved historical keys.

Run native timing blocks and the normal game screen separately. Do not overlap
CPU-heavy validation with either. This tooling does not change engine policy.

The migration retires 131 historical Python/shell source files, including source
stored as `.txt`. No language is hidden from GitHub. The hashed release
attachment preserves 1,148 original files and its own manifest. The three
post-v0.18 attempts were still local when archived; their scripts therefore
need that attachment, while older scripts also remain in earlier tags.

Shared archive verification, Go-engine restoration, frozen decision checks,
mate-proof checks and the shell test fixture now use Go. One-off old corpus
collectors, NumPy optimizers, external python-chess oracle drivers and historical
batch orchestrators are **retired**, with no claim of equivalent Go ports. Their
original code/results remain reproducible in the historical environment from
the attachment. Future collection/fitting/oracles must use Go under a separately
reviewed experiment design. Calling the same Go rules library does not provide
a new independent rules oracle.

Documentary shell command examples remain for invoking Go and Git. The optional
lichess-bot client is installed outside this repository and has its own Python
requirements; that external service client is neither vendored nor rewritten here.
All repository program/test/research source is Go. The root source-policy test
rejects other source extensions, including archived script-name disguises.
Publication text must use neutral, factual engineering prose without first-person
voice or personal framing. The `wording` command rejects first-person pronouns
and the prohibited role word specified in [AGENTS.md](../../AGENTS.md).
