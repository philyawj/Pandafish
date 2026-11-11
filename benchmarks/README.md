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

build_dir=$(mktemp -d)
git archive v0.3 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.3" .)

build_dir=$(mktemp -d)
git archive v0.4 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.4" .)

# After v0.5 has been tagged:
build_dir=$(mktemp -d)
git archive v0.5 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.5" .)
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

The runner does not fetch releases or inspect tags. In direct match mode you
supply executables; validation mode builds the working candidate automatically.
`--candidate` and `--previous` select binaries in direct match mode;
`--version` and `--previous-version` label their results. Labels do not verify
source or change an engine's UCI name. The runner records the actual UCI name
and binary SHA-256 alongside each label.

After v0.5 is released, build its tagged baseline as above for work toward
v0.6. Then run one command from the repository root:

```sh
go run ./benchmarks --validate --version v0.6-dev \
  --previous benchmarks/bin/pandafish-v0.5 --previous-version v0.5
```

Validation checks formatting, runs Go tests and vet, builds your working source
as `benchmarks/bin/pandafish-candidate`, and plays 100 games per opponent. Add
`--race` when changing concurrent code. Omit `--candidate` in validation mode;
the fixed candidate path protects released baselines from accidental overwrites.
Any failed code check stops before matches. Forfeits make validation exit with
an error after saving the results. A low match score is evidence to review,
not a code-validation failure. The local lichess-bot wrapper check is separate.

The default output is a short check/result summary, with detailed records in
JSON and PGN. Use `--verbose` only when you want progress after every color pair.
All moves and scores are handled by Go; no LLM is involved in playing games.
Automated reviews should read the summary first and inspect individual games
only to investigate a failure or a specific behavior.

No release tag is needed. A committed candidate should use a label containing
its revision, such as `v0.6-dev-ab12cd3`; append `-dirty` if source is modified.
Optional `--revision` and `--previous-revision` fields also record source
revisions. `--machine` records a hardware and power-setting note. A binary hash
identifies executable bytes, but does not reconstruct uncommitted source.

Before playing, the runner copies all executables to a temporary directory,
identifies and hashes those copies, and uses them for every game. Rebuilding an
original executable after snapshotting cannot change an active match. Temporary
copies are removed when the runner exits normally or returns an error. Rebuild
files between runs; avoid changing an executable while its initial copy is made.

After reviewing and releasing a candidate, build its baseline from the release
tag and use it as `--previous` for the next cycle. Preserve released binaries
rather than overwriting them with development builds.

### Moving a candidate to release

Follow [the release workflow](../RELEASING.md) for the preparation, PR, merge,
tag, and publishing steps.

When the checks and routine matches support an improvement, open a PR with the
engine change and its benchmark evidence. Prepare the release version in the
source, UCI identity, example greetings, and relevant documentation. After the
PR is reviewed and merged, tag the merged commit and publish the release.

Keep development-run labels, hashes, and results as recorded. A benchmark label
is not promoted into a release tag. Link the candidate's results from the PR
or release notes instead of relabeling them as a test of the released binary.
Build the next baseline from the release tag.

If only version text and documentation change after testing, rerun the code
checks and a short smoke test. Repeat the 100-game comparison when engine
behavior changes. A separate run labeled with the release version is optional
if you want results recorded for that exact released binary.

## Deciding whether a change helps

| Stage | Games per opponent | Purpose |
| --- | --- | --- |
| Unit tests and focused positions | No match required | Check intended behavior and known regressions |
| Smoke test (`--pairs 1`) | 2 | Check executable paths, options, and UCI integration |
| Screening (default `--pairs 50`) | 100 | Find large changes and reliability problems |
| Optional extended comparison (`--pairs 500`) | 1,000 | Resolve small differences or support a planned skill-level advancement |

The routine for early engine development is 100 games per opponent. A promising
screening result does not automatically trigger a longer run, and a longer run
is not required for every release. Use focused tests and the short matches to
evaluate large functional improvements such as adding an opponent reply.
Long matches gather statistical evidence, rather than stress-testing an engine.

Separate functional correctness from measured playing strength. Focused tests
can verify a chess rule or intended decision, while a match against the prior
release measures whether the full engine benefits under the tested conditions.
Report both the result and its limits when documenting a change.

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

Use the screening suite unchanged for routine comparisons. When more evidence
is needed, explicitly request a broader, longer comparison:

```sh
go run ./benchmarks --validate --version v0.6-dev \
  --previous benchmarks/bin/pandafish-v0.5 --previous-version v0.5 \
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

## Recorded development milestones

The [results report](results/README.md) contains the complete counts and run
conditions. Development labels and binary hashes remain as originally recorded.

The v0.3 candidate scored 40 wins and 60 draws against v0.2 in 100 games, and
lost all 100 games against Stockfish Skill Level 1. Both matchups had no
forfeits or unfinished games.

The v0.4 candidate scored 44 wins, 39 draws, and 17 losses against v0.3 in 100
games, a 63.5% score. This is a measured improvement under the recorded
conditions, supported by focused tests showing that both colors choose mate
over material, avoid immediate mating replies, and avoid stalemate when ahead.
The sample supports this early increment without establishing a precise Elo
gain or stronger play against every opponent or time control. The candidate
lost all 100 games against Stockfish Level 1. Neither matchup had forfeits or
unfinished games.

An earlier v0.5 knight-only experiment (recorded as `v0.5-dev`) scored
45 wins, 20 draws, and 35 losses against v0.4, a 55% score in 100 games.
This is a modest lead in one screening run, not conclusive evidence of a
general strength improvement. Its focused tests checked central knight
development without overlooking pawn captures or immediate material losses.

Against Stockfish 19 at Skill Level 1, that experiment scored 2 wins and 98 losses.
These are the first wins in the recorded routine Stockfish benchmarks: White
in round 65 and Black in round 98, both by checkmate. Both matchups completed
without forfeits or unfinished games. The candidate's largest measured move
response was 38.4 ms, within the requested 100 ms budget in this run.
The [saved results and games](results/README.md) preserve the full evidence.
The Stockfish target remains Level 1; a 2% score does not justify advancing.
The knight bonus was set aside in favor of deeper alpha-beta search. Its
results remain preserved as an experiment, not evidence for the replacement
candidate. The alpha-beta candidates use the unchanged v0.4 evaluator so the
search change can be assessed separately.

The first alpha-beta run (`v0.5-alpha-beta-dev`) won 99 games against v0.4 but
lost one on a move-response timeout. The remaining Stockfish matchup was
stopped to address that failure; the partial results are preserved. Captures
are now searched first to improve pruning without changing minimax scores.
The corrected candidate is recorded as `v0.5-alpha-beta-ordered-dev`.
It scored 95 wins, 4 draws, and 1 loss against v0.4 in 100 games, a 97% score,
with no forfeits or unfinished games. Together with the focused tactical
tests, this supports the deeper-search increment under the recorded settings.
It does not establish a precise Elo gain or strength at other time controls.
Against Stockfish 19 at Skill Level 1, it scored 4 wins and 96 losses, with no
forfeits or unfinished games. Its slowest measured response across both
matchups was 65.3 ms, within the requested 100 ms budget. This supports the
candidate's reliability in these matches, not a guarantee for every position:
the engine still uses fixed depth and ignores UCI clock limits. The Stockfish
target remains Level 1. The difference between four wins here and two in the
knight experiment is too small to establish a reliable gain against Stockfish.

The exact position that timed out is retained in `BenchmarkBestMove`. To check
search cost separately from match strength, run from the repository root:

```sh
go test -run '^$' -bench '^BenchmarkBestMove$'
```

This reports average time per search on your machine, not a worst-case time
guarantee. Unit tests compare alpha-beta against exhaustive minimax, including
that position, and check queen-winning forks that need the third ply.

An earlier v0.3 screening run recorded the first local benchmark draw against
Stockfish 19 at Skill Level 1: a stalemate with Pandafish playing White in round
25. The [milestone PGN](milestones/first-stockfish-level-1-draw.pgn) is preserved
separately from routine results.

## Stockfish skill progression

Start at Skill Level 1. For a planned advancement, use an extended comparison
with enough evidence to resolve the result. Once it scores strictly above 50%
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
  --candidate benchmarks/bin/pandafish-candidate --version v0.6-dev \
  --previous benchmarks/bin/pandafish-v0.3 --previous-version v0.3 \
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
outcomes, and move timings. Validation-mode runs also record the automated
checks and commands. Each PGN includes its opening moves.

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
