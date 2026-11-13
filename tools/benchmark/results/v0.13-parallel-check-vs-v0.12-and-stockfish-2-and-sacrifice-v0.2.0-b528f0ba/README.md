# v0.13 three-way concurrency check

The maintainer requested 50 games each against v0.12, Stockfish Skill 2, and
pinned Sacrifice, starting all three matchups together. This tests a faster
workflow after v0.13; it is not a new engine version. The released engine and
private bot configuration were not changed. No online games were played.

## Method

Candidate: tagged v0.13 at 3228cf76550142cd4a2cce50a146f621fa6920e8,
built from the tag, identifying as Pandafish v0.13. Its SHA-256 is
`d802d96a2ccfbaeca419175c528a5c02893231ae6c7a6ec01c87b8ec4dd792d8`.
Use the same immutable v0.12 baseline and pinned Sacrifice v0.2.0 b528f0ba as
in the published v0.13 screen. Stockfish's binary is unchanged, but its Skill
Level is explicitly 2 instead of 1. Threads 1, Hash 16 MiB, and Elo limiting
disabled are unchanged. Summary JSON records actual identities, options, asset
hashes, runner hash, and settings.

Each matchup has 25 color pairs using the existing ten-opening suite, 100 ms
requested per move, 50 ms protocol allowance, and 400-ply cap. Every pair swaps
colors. With 25 pairs, the first five openings appear three times and the rest
twice. A maximum of three games run concurrently, one per matchup; each game's
engines alternate searching. Finished matchups are not replaced by extra games.

Hardware: Apple M2 MacBook Air. Preflight and final observations reported AC
power while the battery discharged from 72% to 65%. These observations do not
measure sustained temperatures, CPU placement, or frequency. `caffeinate` was
active and no other engine experiments or live bot ran during matches.

```sh
caffeinate -i tools/benchmark/bin/benchmark \
  --candidate tools/benchmark/bin/pandafish-v0.13 --version v0.13-parallel-check \
  --previous tools/benchmark/bin/pandafish-v0.12 --previous-version v0.12 \
  --previous-revision a7f6101819697cc0e756ce233a8468c9d41183e7 \
  --revision 3228cf76550142cd4a2cce50a146f621fa6920e8 \
  --stockfish-skill 2 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --parallel-all --pairs 25 \
  --machine 'Apple M2 MacBook Air; AC power reported at start; three matchups concurrent; no live bot' \
  --verbose
```

`--parallel-all` is the new opt-in runner mode. It starts the previous-version
matchup alongside the external comparisons, keeps their PGNs separate, and
updates shared summaries from one coordinator. The previous-release criterion
is still reviewed. If it fails, all concurrent results are retained and the
command returns an error after finishing; complete data is not confused with
a passing decision. Existing serial and external-only `--parallel` modes remain.

The [runner source snapshot](runner-source.json) is relative to the released
base commit above. Rebuild with `go build -o tools/benchmark/bin/benchmark
./tools/benchmark` before running this experiment on that source.

## Checks

Formatting, `go test ./...`, `go test -race ./...`, `go vet ./...`, runner build,
and whitespace checks passed. Focused tests cover all-three overlap, successful
and failed previous-release decisions, result ordering/checkpointing, draining
peer workers after an error, and report metadata. A six-game real-engine smoke
outside tracked history completed without forfeits or unfinished games. Its
previous-version pair drew twice, correctly returning a blocked decision after
retaining every result. That tiny smoke is integration evidence, not strength.
See [check record](checks.json). No engine behavior or dependency changed, so
this runner-only task did not rebuild the active engine or repeat its wrapper.

## Game results and published comparison

The published values were checked against [PR #21](https://github.com/philyawj/Pandafish/pull/21)
and the [v0.13 release](https://github.com/philyawj/Pandafish/releases/tag/v0.13).
Those games used v0.13-dev; release preparation changed comments, version text,
and documentation, not the search algorithm. Original labels/hashes remain.

| Opponent | Published serial 100 games | Published first 50 games | Concurrent 50 games |
| --- | --- | --- | --- |
| v0.12 | 46/43/11 (67.5%) | 20/23/7 (63%) | 25/23/2 (73%) |
| Sacrifice v0.2.0 b528f0ba | 24/28/48 (38%) | 10/14/26 (34%) | 9/16/25 (34%) |
| Stockfish Skill 1 | 58/11/31 (63.5%) | 30/8/12 (68%) | Not run |
| Stockfish Skill 2 | Not run | Not run | 26/8/16 (60%) |

All new 150 games finished, with zero forfeits and zero unfinished games.
The previous-release criterion passed. First-50 prefixes were selected before
reviewing scores to match the new 25-pair opening sequence, not to select a
favorable historical result. All historical records remain unchanged.

The v0.12 score is higher, not lower; Sacrifice is four percentage points below
the full published screen and exactly matches its first-50 score. These small
screens with repeated openings and random ties do not establish statistical
equivalence or prove concurrency has no strength effect. The v0.12 prefix
comparison differs by 10 points, illustrating the variation still present.
Stockfish Skill 2 is a new baseline: its 60% score cannot be compared with
Skill 1 to isolate the effect of concurrency. No serial Skill 2 control was run.

| Opponent | Candidate responses | Above 100 ms | Above 150 ms | Candidate maximum ms |
| --- | --- | --- | --- | --- |
| v0.12 | 2,754 | 1 | 0 | 100.058 |
| Stockfish Skill 2 | 2,396 | 3 | 0 | 113.565 |
| Sacrifice | 2,705 | 4 | 0 | 112.118 |

Eight of 7,855 candidate responses exceeded the requested 100 ms, with none
past the 150-ms protocol deadline. No opponent response exceeded that deadline
either. There is no observed deadline/reliability regression in this run.
This is an observation, not a real-time guarantee.

The matches took about 513.5 seconds (8 minutes 34 seconds), using creation
and final summary-write timestamps. Recorded response times sum to about
1,447 seconds across games and overlap in wall time. That sum is not a fresh
measurement of how long the same games would take serially. No full serial
rerun or favorable-result retry was made.

## Search-throughput calibration

Before matches, the same released v0.13 binary ran cold-process probes on six
saved stream positions. Each mode has 108 depth-four and 108 independent
100-ms probes. Three workers run sequentially or concurrently; their position
orders are offset. Blocks run serial/parallel/parallel/serial to reduce a simple
order/warm-up bias. All 432 probes returned legal moves without fallbacks.

| Measurement | Serial | Three searches | Parallel change |
| --- | --- | --- | --- |
| Aggregate depth-four nodes/sec | 153,251 | 149,544 | -2.4% |
| Aggregate timed nodes/sec | 164,612 | 156,627 | -4.9% |
| Median per-position depth-four NPS ratio | 1.000 | 0.923 | -7.7% |
| Median per-position timed NPS ratio | 1.000 | 0.940 | -6.0% |
| Timed probes completing depth 4 / depth 3 | 69 / 39 | 66 / 42 | Small distribution shift |

Per-position median ratios are more informative about typical throughput here:
a few slow serial probes affect the aggregate rates. Contention is measurable,
roughly 6% to 8% in typical probe throughput, rather than zero. The diagnostic
workload took 24.15 seconds serially versus 8.15 seconds concurrently, about a
threefold wall-time improvement. This is homogeneous Pandafish load, not the
exact mixed-engine memory/cache load or a sustained thermal certification.
No parameters or time allowances were tuned from these results.

## Recommendation

Three concurrent matchups look worthwhile for this stage of development.
The unchanged-opponent results do not show a broad strength collapse, timings
show no new protocol failures, and the measurable throughput cost is modest
relative to the wall-time saving. This supports practical confidence in the
faster workflow on this host, not proof that serial and parallel scores are
interchangeable. Keep this profile as a new comparison baseline and use the
same concurrency for future comparisons. Skill 2 has a promising first screen;
its exact strength and long-control performance remain uncertain.

Recommend 50 games each in three-way concurrency for the next development
increment, approved by the maintainer after reviewing this experiment. Serial comparisons
remain useful when debugging timing problems, using a constrained host, or
investigating a small result. The experiment itself used v0.13 and the opt-in scheduling mode. The
subsequent v0.14 tooling release adopts the faster profile and Skill 2 as
defaults without changing the playing algorithm. Its short release checks
remain separate from this 150-game record.

## Records and reproduction

- [Compact comparison](comparison.json), [full summary/timings](summary.json),
  and [published PR/release snapshots](published-evidence.json)
- [All games](games.pgn), [v0.12](previous.pgn), [Stockfish Skill 2](stockfish.pgn),
  and [Sacrifice](sacrifice.pgn)
- [Throughput aggregates](throughput-summary.json) and [raw probes](throughput-audit.json)
- [Power and runtime observations](run-observations.json), [checks](checks.json),
  and [runner source snapshot](runner-source.json)
- Standalone Go source: [throughput driver](throughput-main.go.txt),
  [engine client](throughput-engines.go.txt), and [comparison](comparison.go.txt)

Copy the `.go.txt` audit sources to a temporary directory with `.go` extensions.
From the repository root, run the throughput driver with the engine-client
source and an output JSON path. It reads `testdata/quiescence-validation.json`
and uses the tagged binary path shown in the match command. Run the standalone
comparison with three arguments: the published v0.13 development summary, the
new summary, and an output JSON path. Its scoring function is an exact copy of
the tested runner's `summarize`; prefix selection and timing aggregation run in
Go without per-game LLM analysis. Hashes identify both comparison inputs.


## Decision after review

The maintainer accepted the development-speed gains and made 50 games per
opponent with all three matchups concurrent the standard after v0.13. At a
later stage, when smaller gains matter more, revisit serial testing and sample
size for the particular question. Keep historical records and concurrency
profiles intact. v0.14 publishes the tooling and policy; it makes no additional
engine-strength claim. Release-only identity/documentation changes do not
relabel this v0.13 run or its binary hashes.


Release validation is recorded in [v0.14 checks](release-v0.14-checks.json).
Full tests/race checks/vet, ARM64 build, offline 12-ply wrapper, and config
loading passed. The final default-schedule smoke completed all six games with
no forfeits or unfinished games, using three-way scheduling and default Skill 2.
Its two-game previous-version score blocked the strength criterion; all six
results were correctly retained. This integration sample does not measure a
strength change and was not retried for a favorable result. Search/evaluation
sources are unchanged; v0.14 is a tooling release following v0.8's precedent.
