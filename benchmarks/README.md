# Local benchmarks

Track Pandafish against the previous release and a Stockfish skill ladder,
starting at Level 1 and progressing toward Level 20. This Go program runs local
UCI matches separately from the engine. It reuses the project's chess library.

[Results, run profiles, timings, and saved games](results/README.md)

## Why this exists

This benchmark runner was created without knowing how Fishtest works. After
building it, we realized it follows a similar methodology: compare a candidate
engine against a baseline through controlled matches and use the results to
evaluate changes.

[Fishtest](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-FAQ.html)
distributes games across volunteer computers and uses statistical tests to
assess small strength improvements. Pandafish's runner runs locally with fixed
game counts and reports scores without a statistical pass/fail verdict.

For now, it lives in the same repository so engine changes, benchmark tests,
instructions, and supporting results can evolve together in one pull request.
The separate `benchmarks/` folder keeps it out of the engine executable while
we build out Pandafish. A separate project would make more sense if the runner
eventually serves other engines or needs distributed workers.

## Build

Use Go 1.25 or newer and Git. [Download Stockfish](https://stockfishchess.org/download/)
and make its executable available as `stockfish`, or pass its path with
`--stockfish /absolute/path/to/stockfish`. Keep the same Stockfish version for
comparisons. The initial baseline uses Stockfish 19.

From the repository root, build the runner and released baselines. Temporary
source checkouts keep your current source untouched:

```sh
mkdir -p benchmarks/bin
benchmark_bin="$PWD/benchmarks/bin"
go build -o "$benchmark_bin/benchmark" ./benchmarks

build_dir=$(mktemp -d)
git archive v0.1 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.1" .)

build_dir=$(mktemp -d)
git archive v0.2 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.2" .)
```

These commands assume a macOS/Linux shell. Temporary source directories can be
deleted after building. Keep the released binaries as baselines. The Go runner
accepts executable paths directly and contains no paths to the maintainer's machine.
Official Stockfish builds include their default neural network. Custom engine
builds requiring external data files need those files available independently;
only the executable is copied into the benchmark's temporary snapshot.

## Run the initial baseline

From the repository root:

```sh
benchmarks/bin/benchmark \
  --candidate benchmarks/bin/pandafish-v0.1 --version v0.1

benchmarks/bin/benchmark \
  --candidate benchmarks/bin/pandafish-v0.2 --version v0.2 \
  --previous benchmarks/bin/pandafish-v0.1 --previous-version v0.1
```

These play v0.1 against Stockfish, then v0.2 against v0.1 and Stockfish. Each
matchup has 50 color pairs, or 100 games. Each pair uses the same opening with
the engines swapping colors. The ten openings in `openings-screening-10-v1.txt`
cycle five times.

The runner sets Stockfish's native UCI Skill Level 1 itself. No Stockfish config
file is needed. This is separate from Lichess's numbered computer opponents.
`UCI_LimitStrength` is disabled so its Elo setting cannot override the skill.
[Stockfish documents these controls here](https://official-stockfish.github.io/docs/stockfish-wiki/Stockfish-FAQ.html#how-do-skill-level-and-uci_elo-work).

## Test a development candidate

The runner does not fetch releases, inspect tags, or build engines. You supply
executables and labels. `--candidate` and `--previous` select binaries;
`--version` and `--previous-version` label their results. Labels do not verify
source or change an engine's UCI name. The runner records the actual UCI name
and binary SHA-256 alongside each label.

For work toward v0.3, build the released v0.2 baseline as above. Then run:

```sh
go test ./...
go vet ./...
go build -o benchmarks/bin/benchmark ./benchmarks
go build -o benchmarks/bin/pandafish-candidate .

benchmarks/bin/benchmark \
  --candidate benchmarks/bin/pandafish-candidate --version v0.3-dev \
  --previous benchmarks/bin/pandafish-v0.2 --previous-version v0.2
```

This builds your working source, including uncommitted changes, and plays 100
games against each opponent. Rebuild the candidate after each source change.
No release tag is needed. A committed candidate should use a label containing
its revision, such as `v0.3-dev-ab12cd3`; append `-dirty` if source is modified.
Optional `--revision` and `--previous-revision` fields also record source
revisions. `--machine` records a hardware and power-setting note. A binary hash
identifies executable bytes, but does not reconstruct uncommitted source.

Before playing, the runner copies all executables to a temporary directory,
identifies and hashes those copies, and uses them for every game. Rebuilding an
original executable after snapshotting cannot change an active match. Temporary
copies are removed when the runner exits normally or returns an error. Rebuild
files between runs; avoid changing an executable while its initial copy is made.

After reviewing and releasing a candidate, build its baseline from the release
tag. For the next cycle use released v0.3 as `--previous`, with a `v0.4-dev`
candidate. Preserve released binaries rather than overwriting them with development builds.

## Deciding whether a change helps

| Stage | Games per opponent | Purpose |
| --- | --- | --- |
| Unit tests and focused positions | No match required | Check intended behavior and known regressions |
| Smoke test (`--pairs 1`) | 2 | Check executable paths, options, and UCI integration |
| Screening (default `--pairs 50`) | 100 | Find large changes and reliability problems |
| Confirmation (`--pairs 500`) | 1,000 | Gather evidence before a release or skill-level advancement |

Choose the sample size before starting. These are practical starting points,
not guarantees of significance. Keep unfavorable and inconclusive results;
do not repeatedly rerun until a favorable score appears. A small lead is not
automatically an improvement. The runner does not issue a statistical pass/fail
verdict. Engine testing needs to account for paired results and uncertainty;
[Stockfish's testing documentation](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-Mathematics.html)
describes a more advanced system than this runner implements.

The filenames identify the purpose, position count, and suite version:

| File | Use |
| --- | --- |
| `openings-screening-10-v1.txt` | Quick screening: 10 positions, used by default |
| `openings-confirmation-50-v1.txt` | Broader confirmation: 50 positions, selected explicitly |

Use the screening suite unchanged for quick comparisons. For a broader confirmation run:

```sh
benchmarks/bin/benchmark \
  --candidate benchmarks/bin/pandafish-candidate --version v0.3-dev \
  --previous benchmarks/bin/pandafish-v0.2 --previous-version v0.2 \
  --pairs 500 --openings benchmarks/openings-confirmation-50-v1.txt
```

That plays 1,000 games per opponent, or 2,000 total. The separate confirmation
suite contains 50 distinct legal opening positions, played with both colors.
It is a public, hand-selected variety suite, not a calibrated or hidden test
set. Compare candidate and baseline under the same suite and settings. More
repetitions of identical games do not add information for deterministic engines.
Do not tune only to these openings; add focused tactical and endgame tests for
the behavior changed. Keep existing suites fixed and add separately named
versions when expanding coverage.

Investigate any illegal moves, crashes, or timeouts before interpreting strength.
The completed-game score excludes unfinished games. The possible-score range
shows all unfinished games counted as losses through all counted as wins.
It is an outcome bound, not a confidence interval. Do not silently count
unfinished games as draws or ignore a rising unfinished count.

Keep Stockfish version, skill, opening hash, time settings, move cap, compiler,
and hardware consistent. Run matches sequentially with a stable power setting
and little background load. The report shows comparison profiles and timing
measurements. Its fixed move allowance does not simulate a tournament's game
clock; test game-clock handling separately before tournament use.

## Stockfish skill progression

Start at Skill Level 1. After a confirmation match scores strictly above 50%
at the current level, preserve that result and advance one level. Score means
wins plus half the draws divided by games, rather than the fraction won outright.
If games are unfinished, require the possible-score lower bound to exceed 50%
for this decision. Do not advance on results containing protocol failures or
timing problems. Close results remain inconclusive without statistical analysis.

Record the source revision, Stockfish version and skill, sample size, settings,
W/D/L, unfinished games, and JSON/PGN links when documenting an advancement.
Keep previous levels' records. Continue comparing against the previous Pandafish
release after moving to a new Stockfish level.

Advance manually with `--stockfish-skill 2`, then 3, and so on through 20.
The flag accepts 1 through 20 and defaults to 1. For example:

```sh
benchmarks/bin/benchmark \
  --candidate benchmarks/bin/pandafish-candidate --version v0.3-dev \
  --previous benchmarks/bin/pandafish-v0.2 --previous-version v0.2 \
  --stockfish-skill 2
```

Skill is recorded in the options, matchup label, and folder name. The runner
does not automatically advance levels. Pandafish releases and ladder advances
are separate decisions: an improved release may still lose every game at the
current Stockfish level. Previous-version matches and focused positions remain
useful when the Stockfish score stays at zero. Skill levels are not Elo ratings,
and success at 100 ms does not establish equal strength at longer controls.

As Pandafish becomes stronger, expand the benchmarks to multiple independent
engines, for example [Sunfish](https://github.com/thomasahle/sunfish), alongside
Stockfish and previous Pandafish versions. This will help detect improvements
that work against one opponent but regress against others. Pin each opponent's
version and settings and introduce it as a separate comparison series. Additional
engine support is future work; the current runner supports Stockfish and one
previous Pandafish executable.

## Timing and game rules

- Each engine receives `go movetime 100`, requesting 100 ms per move. The hard
  wall-clock deadline is 150 ms by default: request plus a 50 ms allowance for
  position processing, scheduling, and protocol overhead. A move received after
  that deadline forfeits. The deadline also bounds a blocked command write.
- Every move attempt records elapsed wall time, including a failed attempt.
  The report shows maximum times and counts above the requested budget for
  each side. Exceeding the request while staying within the explicit allowance
  is diagnostic information, not a forfeit or a measurement of pure search CPU time.
- `--move-ms` and `--grace-ms` change the timing profile. Grace may be zero.
  Keep both fixed between comparisons. The former two-second watchdog results
  have been removed; the saved baseline starts fresh under `wall-clock-v2`.
- Stockfish uses one thread, 16 MiB of hash, the selected skill, Elo limiting off,
  and no Syzygy tablebases. Every game starts fresh. No pondering or opening
  book is requested beyond the opening moves supplied by the runner.
- Threefold repetition and fifty-move draw claims are taken when available
  at the current position. The chess library recognizes automatic game endings.
- Missing or illegal moves, crashes during play, and move timeouts forfeit.
  Startup/configuration errors abort the run instead of awarding games.
- The 400-ply cap marks games unfinished. `--max-plies` changes it. A terminal
  result or valid draw claim at the cap takes precedence over the cap.
- Completed score is `(wins + draws / 2) / completed games`. Possible score is
  `(wins + draws / 2)` through `(wins + draws / 2 + unfinished)`, divided by all
  games played. Forfeits count as losses for the offending engine.

## Saved history and checks

Results go to matchup folders such as `v0.2-vs-v0.1-and-stockfish-1`, with
`-run-2`, `-run-3`, etc. for repeats. Each contains `games.pgn` and `summary.json`.
JSON records binary identities and hashes, optional source revisions, runner
hash, Go version, platform, optional machine note, settings, opening hash,
outcomes, and move timings. Each PGN includes its opening moves.

After each completed matchup, the runner syncs the PGN and atomically saves
its summary before starting the next opponent. Completed matchups from an
interrupted run remain visible as a partial run. An incomplete matchup may
have additional PGNs but is excluded from the summarized scores. There is
no automatic resume. The report distinguishes partial and completed runs.

`results/README.md` is rebuilt after each matchup and at completion. To rebuild
it without playing, run `benchmarks/bin/benchmark --report-only`. Commit the
JSON, PGNs, and Markdown report when publishing a benchmark. Compiled binaries
stay gitignored in `bin/`; the benchmark has no Python dependency.

Use `--pairs 1 --results /tmp/pandafish-smoke` for a quick check outside the
tracked history. `go test ./...` runs both engine and runner unit tests;
`go test -race ./...` also checks concurrent code. Unit tests use a small Go
UCI helper, so they do not require Stockfish. They cover failures, timeouts,
blocked writes, snapshots, options, color scoring, score bounds, move caps,
opening legality, PGN output, checkpoints, and report history. These local
checks do not connect to Lichess.
