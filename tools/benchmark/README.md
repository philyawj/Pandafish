# Local benchmarks

Track Pandafish against the previous release, pinned Sacrifice, and a
Stockfish skill ladder starting at Level 1 and progressing toward Level 20.
This Go program runs local UCI matches separately from the engine. It reuses the project's chess library.

[Results, run profiles, timings, and saved games](results/README.md)

[Live Lichess session archives and analyses](live/README.md) preserve real-game
evidence separately from controlled matches. Capture and analysis code lives
in `tools/lichess/`.

## Routine workflow after v0.13

Starting after v0.13, use **50 games per opponent with all three matchups
concurrent** to speed up development. That is 25 color pairs each against the
previous release, Stockfish Skill 2, and pinned Sacrifice: 150 games total.
The runner defaults to 25 pairs and Skill 2. A full three-opponent suite starts
all matchups together on hosts with at least four logical CPUs. Smaller hosts
fall back to serial execution; `--parallel-all` requires the concurrent profile
explicitly and rejects insufficient capacity.

The previous-release criterion remains a review requirement, rather than an
early scheduling gate in concurrent mode. A failed screen keeps all results
and returns an error after the matchups finish. Historical records, including
v0.13's 100-game serial screens and the subsequent concurrency experiment,
retain their original labels, sample sizes, and profiles.

This prioritizes faster feedback on noticeable improvements. The calibration
found a typical 6% to 8% search-throughput cost, but the 150-game suite finished
in about 8.5 minutes without protocol failures. That tradeoff is worthwhile
at this stage; it does not establish that parallel and serial strength scores
are interchangeable. Keep concurrency consistent for future comparisons.

When development reaches smaller gains, reconsider serial execution and larger
samples for the specific question. Use `--sequential` and, if justified,
`--pairs 50` for 100 games per opponent or a planned extended comparison.
Do not extend every promising screen automatically or rerun until favorable.
Timing budgets, openings, unit tests, race checks when applicable, and integration
checks stay unchanged. See the [calibration record](results/v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
for measurements and limits.

## Local development and live sessions

Local testing should drive the next batch of strength changes. Keep each
experiment focused, check saved positions and search behavior, then use the
[previous-release criterion](#routine-order-and-practical-gate) and concurrent
comparisons. A live stream is not required for every increment. Use live sessions
periodically after meaningful progress, or when clocks, UCI handling, or client
integration change. Controlled matches provide cleaner strength comparisons;
live games expose operational failures and new positions.

The [v0.10 review](live/v0.10-twitch-m2-20261112/REVIEW.md#current-strength-priorities)
records the current five strength priorities and the reasoning behind them.
They are hypotheses to test, not guaranteed gains or an approved release schedule.

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
The separate `tools/benchmark/` folder keeps it out of the engine executable while
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
mkdir -p tools/benchmark/bin
benchmark_bin="$PWD/tools/benchmark/bin"
go build -o "$benchmark_bin/benchmark" ./tools/benchmark

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

build_dir=$(mktemp -d)
git archive v0.5 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.5" .)

build_dir=$(mktemp -d)
git archive v0.6 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.6" .)

build_dir=$(mktemp -d)
git archive v0.7 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.7" .)

build_dir=$(mktemp -d)
git archive v0.8 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.8" .)

build_dir=$(mktemp -d)
git archive v0.9 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.9" .)

build_dir=$(mktemp -d)
git archive v0.10 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.10" .)

build_dir=$(mktemp -d)
git archive v0.11 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.11" .)

build_dir=$(mktemp -d)
git archive v0.12 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.12" .)

build_dir=$(mktemp -d)
git archive v0.13 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.13" .)

build_dir=$(mktemp -d)
git archive v0.14 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.14" .)
```

These commands assume a macOS/Linux shell. Temporary source directories can be
deleted after building. Keep the released binaries as baselines. The Go runner
accepts executable paths directly and contains no paths to the maintainer's machine.
Official Stockfish builds include their default neural network. Custom engine
builds requiring external data files need those files available independently;
only the executable is copied into the benchmark's temporary snapshot.
The pinned Sacrifice build below is one such engine.

## Build the pinned Sacrifice opponent

Starting with v0.9, every strength increment also compares against
[Sacrifice-Engine](https://github.com/piratesboys/Sacrifice-Engine-). The current
pin is package **0.2.0**, commit
`b528f0ba373910b72052d991700a990ca1cb07fd`, including philyawj's merged
[PR #1](https://github.com/piratesboys/Sacrifice-Engine-/pull/1).
Its UCI identity says only `Sacrifice-Engine`; the package version and source
revision identify the build. [Build provenance](opponents/sacrifice.json)
records the toolchain and hashes used for the first comparison.

Install Rust and Cargo, then run from the repository root:

```sh
git clone https://github.com/piratesboys/Sacrifice-Engine-.git tools/benchmark/bin/sacrifice-source
git -C tools/benchmark/bin/sacrifice-source checkout --detach b528f0ba373910b72052d991700a990ca1cb07fd
cargo build --release --locked --manifest-path tools/benchmark/bin/sacrifice-source/Cargo.toml
```

Use the native release executable, not `cargo run`. On macOS/Linux it is
`tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine`.
On Windows add `.exe`. This third-party Rust engine is an opponent only;
Pandafish and the match runner remain Go programs. The checkout, binary, and
large PGN asset stay in ignored `bin/`, not in source control.

**Keep the checkout at its build location.** This pinned engine embeds its
Cargo source path and reads `data/Tal.pgn` there. Copying its executable alone
is insufficient. The runner records that PGN's SHA-256 and checks it before
and after its matchup; pass `--sacrifice-data /path/to/data/Tal.pgn` if built
elsewhere. That flag verifies the asset, it does not redirect the engine's
embedded lookup. The engine's own Tal move-ordering feature stays enabled.

The default `--sacrifice-version` and `--sacrifice-revision` flags describe
this pin. They are provenance labels, not automatic source verification.
Confirm `git -C tools/benchmark/bin/sacrifice-source rev-parse HEAD` before
building. Keep this revision fixed across future increments. Deliberately
updating the opponent requires new labels, source/build provenance, and a
separate comparison series; do not silently replace it with upstream HEAD.

## Run the initial baseline

These commands reproduce the earlier 100-game baseline profile. From the
repository root:

```sh
tools/benchmark/bin/benchmark \
  --candidate tools/benchmark/bin/pandafish-v0.1 --version v0.1 --pairs 50 --stockfish-skill 1 --sequential

tools/benchmark/bin/benchmark \
  --candidate tools/benchmark/bin/pandafish-v0.2 --version v0.2 \
  --previous tools/benchmark/bin/pandafish-v0.1 --previous-version v0.1 --pairs 50 --stockfish-skill 1 --sequential
```

These play v0.1 against Stockfish, then v0.2 against v0.1 and Stockfish. Each
matchup has 50 color pairs, or 100 games. Each pair uses the same opening with
the engines swapping colors. The ten openings in `openings-screening-10-v1.txt`
cycle five times.

These commands set Stockfish's native UCI Skill Level 1. No Stockfish config
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

Build the v0.14 tagged baseline as above for the next engine change. Then run
one command from the repository root:

```sh
go run ./tools/benchmark --validate --race --version candidate-dev --pairs 25 \
  --previous tools/benchmark/bin/pandafish-v0.14 --previous-version v0.14 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine
```

Validation checks formatting, runs Go tests and vet, builds your working source
as `tools/benchmark/bin/pandafish-candidate`, and runs the staged comparisons
below. Add
`--race` when changing concurrent code. Omit `--candidate` in validation mode;
the fixed candidate path protects released baselines from accidental overwrites.
Any failed code check stops before matches. Forfeits make validation exit with
an error after saving the results. The previous-release screen can stop the
workflow in serial mode. In concurrent mode all matchups finish before a failed
previous-release decision returns an error. External scores remain evidence to
review rather than automatic release decisions. The local lichess-bot wrapper
check is separate.

### Routine order and practical gate

1. Run **50 games each against the previous release, Stockfish Skill 2, and
   pinned Sacrifice**, with all three matchups concurrent. Each matchup plays
   one game at a time; the two engines take turns searching.
2. Require a previous-release score strictly above 50% when unfinished games
   count as losses, with zero forfeits by either engine. This is a conservative
   screen, not a significance test or proof that the candidate is stronger.
3. In concurrent mode a failed previous screen returns an error after all
   matchups finish and preserves every result. In `--sequential` mode it stops
   before external matches. Do not discard or relabel a failed screen.
4. Review all three scores, code checks, timings, forfeits, and unfinished games
   before proposing a release. A narrow lead remains inconclusive.

The routine concurrent suite plans 150 games. Keep Skill 2 as the
current reference point; further ladder changes are deliberate decisions.
Sunfish remains a possible future opponent when adding it is worthwhile.

Full validation requires `--previous` and `--sacrifice`. Use `--stage previous`
for only the previous-release matchup or `--stage external` for external
comparisons; these partial stages do not automatically use three-way scheduling.
The latter skips the previous-release review and does not verify a prior result.
When splitting a run, preserve the same candidate hash/settings and link the
result folders. Direct match mode also supports historical two-opponent runs.

### Why routine comparisons are concurrent

The maintainer accepted a modest throughput cost to shorten the development
loop. The [v0.13 calibration](results/v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
compared unchanged opponents with published results and measured search work
under one-search and three-search loads. All 150 games completed in about
8 minutes 34 seconds with no forfeits or unfinished games. Typical probe
throughput fell about 6% to 8%; the diagnostic workload finished about three
times faster. These measurements support the practical choice on this host,
not universal capacity or statistical equivalence.

Use the same profile for future comparisons and record CPU/power/background
workload with `--machine`. Reconsider sequential matches when investigating
smaller gains, timing problems, or constrained hardware. Serial execution and
larger samples are available tools, not automatic follow-ups to every screen.

- Default full suite: all three matchups overlap on at least four logical CPUs.
- `--parallel-all`: explicitly require that full-suite profile and capacity.
- `--sequential`: all matchups serial; the previous-release gate runs first.
- `--parallel`: legacy mode, previous release first, then two external matchups
  concurrent when capacity permits. Do not combine scheduling flags.

The summary records `matchup_concurrency: 3` for all-three scheduling. Historical
external-only concurrency keeps its original metadata. Profiles distinguish
these modes, so historical serial results remain readable without relabeling.
Go checks finish before match play; tests are not run alongside the games.

Starting with v0.7, UCI command handling stays active while search runs. A
script must wait for `bestmove` before sending `quit` or closing stdin, which
cancel unfinished work. Sending `go depth N` followed immediately by `quit`
does not verify a completed depth-N search. The protocol tests exchange
commands and responses through pipes instead of preloading an entire session.

The default output is a short check/result summary, with detailed records in
JSON and PGN. Use `--verbose` only when you want progress after every color pair.
All moves and scores are handled by Go; no LLM is involved in playing games.
Automated reviews should read the summary first and inspect individual games
only to investigate a failure or a specific behavior.

No release tag is needed. A committed candidate should use a label containing
its revision, such as `candidate-dev-ab12cd3`; append `-dirty` if source is modified.
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

Follow [the release workflow](../../RELEASING.md) for the preparation, PR, merge,
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
checks and a short smoke test. Repeat the routine 50-game comparison when engine
behavior changes. A separate run labeled with the release version is optional
if you want results recorded for that exact released binary.

## Deciding whether a change helps

| Stage | Games per opponent | Purpose |
| --- | --- | --- |
| Unit tests and focused positions | No match required | Check intended behavior and known regressions |
| Smoke test (`--pairs 1`) | 2 | Check executable paths, options, and UCI integration |
| Routine screen (default `--pairs 25`) | 50 | Find noticeable changes and reliability problems quickly |
| Optional larger screen (`--pairs 50`) | 100 | Get more evidence when a close result matters |
| Optional extended comparison (`--pairs 500`) | 1,000 | Resolve small differences or support a planned skill-level advancement |

After v0.13, the routine is 50 games per opponent. A promising
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

Use the screening suite unchanged for routine comparisons. With 25 color pairs,
the first five openings appear three times and the remaining five twice; each
pair still swaps colors. Keep that order fixed for comparisons using this profile.
When more evidence is needed, explicitly request a broader, longer comparison:

```sh
go run ./tools/benchmark --validate --version candidate-dev \
  --previous tools/benchmark/bin/pandafish-v0.14 --previous-version v0.14 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --pairs 500 --openings tools/benchmark/openings-confirmation-50-v1.txt
```

That plays 1,000 games per opponent, or 3,000 total if the previous-release
gate passes. The separate confirmation
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
and hardware consistent. Use a stable power setting and little background
load, and keep the external-concurrency profile consistent. The report shows
comparison profiles and timing
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

Positions that timed out are retained in `BenchmarkBestMove`. To check
search cost separately from match strength, run from the repository root:

```sh
go test -run '^$' -bench '^BenchmarkBestMove$'
```

This reports average time per search on your machine, not a worst-case time
guarantee. Unit tests compare alpha-beta against exhaustive minimax, including
that position, and check queen-winning forks that need the third ply.
For v0.7, this benchmark deliberately retains a fixed depth-three search with
the same quiescence cap, so it measures computation cost separately from the
timed UCI search. The match report measures actual UCI response times.

The first v0.6 quiescence run (`v0.6-dev`) scored 83 wins, 10 draws, and 7
losses against v0.5, but two losses were move-response timeouts. The Stockfish
leg was stopped to address those failures, and the partial results are
preserved. A strength lead does not make a run with timing failures ready
for release. The revised candidate (`v0.6-quiescence-dev`) orders captures by
victim and attacker value, reuses root bounds while preserving exact best-move
ties, and avoids repeated evaluation work. The extra tactical search remains
capped at two plies; it does not implement clock management.

Focused tests cover poisoned captures, recaptures, quiet and capture
promotions, underpromotion to avoid stalemate, en passant, and check positions
requiring quiet king moves or blocks. Reference minimax includes the bounded
quiescence rules. Captured timeout positions remain in the performance suite;
the main search comparisons use small representative positions so routine
unit tests stay quick.

The revised candidate scored 80 wins, 16 draws, and 4 losses against v0.5 in
100 games, an 88% score, with no forfeits or unfinished games. This supports
the increment under the recorded settings rather than establishing a precise
Elo gain or reliability at every time control. Its slowest response in this
matchup was 110.0 ms. Six responses exceeded the requested 100 ms budget but
remained within the explicit 50 ms protocol allowance. Clock management is
still needed before expanding the search further.

Against Stockfish 19 at Skill Level 1, the revised candidate scored 12 wins,
1 draw, and 87 losses, a 12.5% score, with no forfeits or unfinished games.
This is a higher observed score than the v0.5 candidate's 4%, while the
previous-release comparison remains the clearest evidence for this increment.
The Stockfish target stays at Level 1. The slowest candidate response across
both final matchups was 134.1 ms. Eleven responses exceeded the requested
100 ms budget but stayed inside the 150 ms deadline including protocol
allowance. These measurements do not establish a worst-case search time.

The v0.7 candidate primarily adds time-control and UCI reliability. Functional
tests establish budget calculation, cancellation, retention of completed
iterations, and legal fallback behavior separately from playing strength.
Iterative deepening can improve strength by reaching more depth, but honoring
a deadline can also stop earlier than v0.6's fixed search. A higher match
score is not the definition of correct time control.

In 100 games against v0.6, `v0.7-dev` scored 63 wins, 21 draws, and 15 losses,
with one unfinished game at the 400-ply cap and no forfeits. The completed-game
score was 74.2%; treating the unfinished game as a loss through a win gives
73.5% to 74.5%. This supports the increment under these settings without
establishing a precise Elo gain or optimal clock allocation. The candidate's
slowest response in that matchup was 101.6 ms, with one response beyond the
100 ms request but within the 150 ms deadline including protocol allowance.

The fixed-movetime matches do not test game-clock allocation over a whole
game. A separate local UCI clock/increment sequence exercises successive
legal moves and clock consumption. The race and protocol tests also cover
`isready`, repeated `stop`, `quit`, held infinite/ponder results, `ponderhit`,
replacement searches, and new-game cancellation. Live-game capture was originally planned for v0.7.1 and shipped as the
separate v0.8 tooling increment.

Against Stockfish 19 at Skill Level 1, v0.7-dev scored 9 wins, 1 draw, and
90 losses, a 9.5% score, with no forfeits or unfinished games. The earlier
v0.6 candidate scored 12.5%; these screens do not establish a clear change
against Stockfish. The target remains Level 1. Across both v0.7 matchups, the
slowest candidate response was 101.9 ms, and two responses exceeded the 100 ms
request while staying within the 150 ms deadline. This is measured protocol
wall time, including position handling and scheduling, rather than a strict
bound on all future responses.

The clock-allocation formula is deliberately conservative and has not been
tuned through live play. The functional checks verify its rules and search
control; later recorded Lichess games can test how well that policy spends
the available game clock. A weaker result against some opponents or at some
time controls would need review without invalidating the independently tested
protocol improvements.

The v0.9 positional-evaluation candidate scored **66 wins, 25 draws, and
9 losses (78.5%)** against v0.8, with no forfeits or unfinished games. This
supports the increment under the serial screening profile; it does not give
a precise Elo estimate. The [previous-release record](results/v0.9-dev-vs-v0.8/README.md)
contains the candidate hash and complete games.

Against Stockfish 19 at native Skill Level 1, the same candidate scored
**47 wins, 10 draws, and 43 losses (52%)**, with no forfeits or unfinished
games. This is the first recorded routine screen above 50% at Level 1 and a
much higher observed score than earlier material-only candidates. The margin
above half is close in 100 games; it does not establish a reliable above-50%
expectation or trigger an automatic level advancement. Keep the Level 1 result
and use the progression policy below before choosing Level 2. The
[external comparison record](results/v0.9-dev-vs-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
contains Stockfish and pinned Sacrifice evidence. Against Sacrifice package
0.2.0 at `b528f0ba`, the candidate scored 4 wins, 17 draws, and 79 losses
(12.5%), with no forfeits or unfinished games. This establishes the first
controlled baseline against that opponent. Across all three matchups the
candidate's maximum response was 120.222 ms; nine responses exceeded the
100 ms request while remaining within the 150 ms deadline.

These earlier routine matchups used the serial profile. Starting after v0.13,
the [current routine](#routine-workflow-after-v013) uses three concurrent
50-game comparisons. Historical results retain their original profiles.

An earlier v0.3 screening run recorded the first local benchmark draw against
Stockfish 19 at Skill Level 1: a stalemate with Pandafish playing White in round
25. The [milestone PGN](milestones/first-stockfish-level-1-draw.pgn) is preserved
separately from routine results.

## Stockfish skill progression

The ladder began at Skill Level 1. Starting with the post-v0.13 workflow,
**Skill Level 2 is the standard reference**. The maintainer selected it after
the Skill 1 results and a 60% first 50-game screen at Skill 2. This is a practical
comparison point, not proof of a reliable above-50% rate at that level.

For later advancement, review results and reliability at the current level,
preserve them, and deliberately choose the next level. Score means wins plus
half the draws divided by completed games. Require the possible-score lower
bound to exceed 50% when unfinished games are present, and do not advance on
protocol failures or timing problems. Close results remain inconclusive; a
larger or serial comparison is optional when needed for a specific decision,
not an automatic 1,000-game requirement.

Record the source revision, Stockfish version and skill, sample size, settings,
W/D/L, unfinished games, and JSON/PGN links when documenting an advancement.
Keep previous levels' records. Continue comparing against the previous Pandafish
release after moving to a new Stockfish level.

Advance manually with `--stockfish-skill 3`, then 4, and so on through 20.
The flag accepts 1 through 20 and defaults to 2. Use `--stockfish-skill 1` to
reproduce earlier profiles. For example, a future Skill 3 comparison:

```sh
tools/benchmark/bin/benchmark \
  --candidate tools/benchmark/bin/pandafish-candidate --version candidate-dev \
  --previous tools/benchmark/bin/pandafish-v0.14 --previous-version v0.14 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --stockfish-skill 3
```

Skill is recorded in the options, matchup label, and folder name. The runner
does not automatically advance levels. Pandafish releases and ladder advances
are separate decisions: an improved release may still lose every game at the
current Stockfish level. Previous-version matches and focused positions remain
useful when the Stockfish score stays at zero. Skill levels are not Elo ratings,
and success at 100 ms does not establish equal strength at longer controls.

Sacrifice now supplies an independent comparison alongside Stockfish and the
previous Pandafish release. Later, add [Sunfish](https://github.com/thomasahle/sunfish)
when the engine is stronger and another comparison is worth the runtime.
Multiple opponents help detect gains that do not generalize. Pin each
opponent's version and settings and introduce updates as separate series.

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
`-run-2`, `-run-3`, etc. for repeats. Each contains combined `games.pgn` and
`summary.json`. New runs also save
`previous.pgn`, `stockfish.pgn`, and `sacrifice.pgn` for the opponents played.
JSON records binary identities and hashes, optional source revisions, runner
hash, Go version, platform, optional machine note, settings, opening hash,
outcomes, and move timings. Validation-mode runs also record the automated
checks and commands. Each PGN includes its opening moves.

After each completed matchup, the runner syncs the PGN and atomically saves
its summary. Parallel external matchups have independent PGN writers; only
the scheduler updates the shared summary. Completed matchups from an
interrupted run remain visible as a partial run. An incomplete matchup may
have additional PGNs but is excluded from the summarized scores. There is
no automatic resume. The report distinguishes partial and completed runs.

`results/README.md` is rebuilt after each matchup and at completion. To rebuild
it without playing, run `tools/benchmark/bin/benchmark --report-only`. Commit the
JSON, PGNs, and Markdown report when publishing a benchmark. Compiled binaries
stay gitignored in `bin/`; the benchmark has no Python dependency.

Use `--pairs 1 --results /tmp/pandafish-smoke` for a quick check outside the
tracked history. `go test ./...` runs both engine and runner unit tests;
`go test -race ./...` also checks concurrent code. Unit tests use a small Go
UCI helper, so they do not require Stockfish. They cover failures, timeouts,
blocked writes, snapshots, options, color scoring, score bounds, move caps,
opening legality, PGN output, checkpoints, and report history. These local
checks do not connect to Lichess.

### v0.8 recording and reporting

v0.8 adds observation and session tooling without changing the evaluator,
search rules, or time-allocation policy. It does not claim a strength gain.
Validation focuses on protocol and capture correctness: unit/race checks,
configuration loading, legal moves through the local client wrapper, complete
search metrics, explicit PGN linkage, and short local integration matches.
The release smoke comparison uses two games each against v0.7 and Stockfish
Skill Level 1. That sample checks integration rather than playing strength;
its records stay outside tracked history. No new routine strength tournament
is required for this tooling increment. Level 1 remains the target.

See the [session tool documentation](../lichess/README.md) for commands,
recording limits, and offline reporting. Future changes to move selection or
time allocation use the routine comparison workflow above with the latest
release as the baseline, currently v0.14.


## v0.10 transposition-table candidate

The candidate adds a fixed-size in-memory table, cached-move ordering, and
per-game reuse while preserving v0.9 evaluation and time control. Cache score
cutoffs require matching depth; deeper entries supply move hints only. This
keeps fixed-depth comparisons exact. Quiescence is not cached. Position keys
include the library's full 64-bit Zobrist key and a separate half-move clock.
Index collisions check the full key; ordinary 64-bit hash collisions remain
a theoretical limitation. Repetition history is not used by the current
search, so caching does not add draw awareness; revisit cache validity when
history-dependent draws are implemented. Mate values remain signs without
ply distance, so no distance normalization is needed in this increment.

Use the routine 100 serial games per opponent: v0.9 first, then Stockfish
Skill 1, then pinned Sacrifice. For this increment, the predeclared criterion
for an **additional** 100-game Skill 2 screen is at least 60% at Skill 1 with
zero forfeits and zero unfinished games. Run that additional comparison only
after Sacrifice finishes and after checking reliability of all routine
matchups. This is a screening threshold for trying the next opponent, not a
statistical confidence claim or permanent ladder advancement. Preserve the
Skill 1 series and record Skill 2 separately; do not rerun until favorable.
The candidate was reviewed in draft PR #17 and approved for v0.10 release.
Its recorded development labels and hashes are retained; release preparation
changes version text and documentation only.

Functional validation compares cached and uncached fixed-depth scores and
complete root tie sets, validates score bounds and depth isolation, checks
collisions and en passant/turn/castling/half-move context, and exercises
cancellation, per-search reuse, restricted roots, and UCI new-game clearing.
To inspect the small node-count checks without playing games:

```sh
go test -run 'TestCacheReuseReducesNodesAndPreservesRootTies|TestUCICacheLifecycle' -v .
```

Node reductions on these positions check reuse; they are not a prediction of
match strength or guaranteed completed-depth gains at every time control.
After v0.10 is reviewed, another recorded Lichess session is the proposed next
checkpoint to find real-game weaknesses and choose subsequent improvements.


The v0.10 candidate completed all 300 serial games without forfeits or
unfinished games: 37/39/24 against v0.9 (56.5%), 45/14/41 against Stockfish
Skill 1 (52%), and 0/27/73 against Sacrifice (13.5%). The previous-release lead
is modest; the external screens do not establish a strength gain. The fixed
position audit preserved all nine depth-three scores and completed a higher
depth on three of nine 100-ms probes, with six equal and none lower in that
single trial. The maximum candidate response across the routine matches was
98.730 ms, with none exceeding the requested 100 ms. These are measurements,
not timing guarantees. Skill 2 was skipped because 52% did not meet the
predeclared 60% trigger. See the [complete v0.10 candidate record](results/v0.10-dev-vs-v0.9-and-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
for counts, method, limitations, raw games, and the recorded decision.


## v0.11 check-aware quiescence candidate

The working candidate keeps the two-ply capture/promotion horizon and extends
checked nodes for at most four additional evasion plies. It still uses static
evaluation at the final safety cap. Evaluation, clocks, TT, and exact root ties
are unchanged. This addresses a concrete horizon limitation; a strength gain
must be measured rather than assumed. The routine previous-release gate uses
v0.10, followed by Stockfish Skill 1 and pinned Sacrifice, all serial.

[Candidate method, focused checks, and results](results/v0.11-dev-vs-v0.10-and-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
record the implementation limit, power conditions, and saved-position probes.
The candidate was reviewed for the v0.11 release. Original development labels
and hashes are retained; the release does not advance the Stockfish target.


## v0.12 quiet-move ordering candidate

Test two quiet beta-cutoff hints per main-search ply, after cached moves and
tactical priorities. Hints reset for each search and do not change evaluation,
quiescence, clocks, or exact root ties. Use v0.11 as the previous-release gate,
then Stockfish Skill 1 and pinned Sacrifice, with all comparisons serial.

[Fresh assessment, source profile, correctness/node checks, and results](results/v0.12-dev-vs-v0.11-and-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
separate deterministic search-work evidence from measured playing strength.


## Choosing the next strength experiment

The [post-v0.12 reference comparison and ranked candidates](research/v0.13-candidates/README.md)
assess selective quiescence, selective main search, and richer classical
scoring. The review includes a reproducible ordinary-exchange cutoff example
and separates measured efficiency from external strength. These are research
options. Recommendation 1 was subsequently selected for the local v0.13
experiment; recommendations 2 and 3 remain future alternatives. The maintainer
approved the selected candidate for v0.13.

The local [v0.13 selective-quiescence experiment](research/v0.13-candidates/selective-quiescence-audit/README.md)
records the SEE tests, broader-search ablation, and selected-position diagnostic.
Its diagnostic scores are separate from routine game results.


## v0.13 selective-quiescence candidate

The local candidate extends ordinary tactical continuation to a 16-ply safety
bound, retains four extra checked-node evasion plies, and filters negative
material exchanges for non-checking ordinary captures outside check. Evaluation,
main search, clocks, TT rules, and root tie handling are unchanged.

The serial 100-game screens scored 46/43/11 against v0.12 (67.5%), 58/11/31
against Skill 1 (63.5%), and 24/28/48 against pinned Sacrifice (38%), with no
forfeits or unfinished games. The external scores improve on the separate
v0.12 screens by 6 and 20.5 percentage points respectively. This is promising
screening evidence, not a precise Elo estimate or proof of gains at all controls.

[Method, checks, timing limits, interpretation, and raw records](results/v0.13-dev-vs-v0.12-and-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
preserve the exact tested source and selected-position ablation. The candidate
was approved for v0.13. Original development labels and hashes are preserved;
release preparation changes version text and documentation only.


## v0.14 benchmark workflow release

v0.14 adopts 50 games per opponent and three concurrent matchups as the routine
for faster iteration, with Stockfish Skill 2 as the reference. `--sequential`
retains the older previous-first schedule for questions about smaller gains
or timing. Small hosts fall back to serial mode unless concurrency is required
explicitly. The playing algorithm remains that of v0.13.

The [v0.13 concurrency experiment](results/v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
is evidence for the workflow, not a fresh strength test of v0.14. Its exact
source, labels, hashes, all 150 games, throughput probes, and comparisons are
preserved. Release validation uses code checks and short integration matches.
