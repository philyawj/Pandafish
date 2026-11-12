# v0.12 quiet-move ordering candidate

The candidate remembers two quiet beta-cutoff moves at each main-search ply.
It tries matching legal hints after captures/promotions; a transposition-table
move still has first priority. Hints do not create moves or supply scores.
Captures, en passant, and promotions are not stored. Each search has a fresh
256-byte hint table; iterative-deepening passes share it. Hints are not saved
between go commands or used in quiescence.

Evaluation, alpha-beta windows, exact root ties, v0.11 quiescence limits,
clock budgets, and transposition-score rules are unchanged. UCI exposes no new
options. The tested candidate identifies as v0.12-dev. It was reviewed for release as
v0.12; original development labels and binary hashes are preserved.

## Fresh assessment and baseline profile

Quiet ordering was preferred over adding evaluation terms or draw awareness
because it can reduce search work without changing the fixed-depth value.
The current engine already orders cached moves and tactical moves, but has no
quiet cutoff hints across sibling nodes. The saved stream cases also included
quiet resources that required greater depth in earlier restricted probes.
Those observations motivate this experiment; they do not prove a strength gain.
Draw awareness remains useful but requires history-sensitive TT handling, and
there is no evidence it is the dominant cause of the saved live errors.

A November 12 source cross-check also keeps the comparison precise:
[Stockfish scores quiet moves using several history tables](https://github.com/official-stockfish/Stockfish/blob/master/src/movepick.cpp),
while [Sunfish's variable named killer is retrieved from its position-keyed move table](https://github.com/thomasahle/sunfish/blob/master/sunfish.py).
Neither should be described as this exact two-hints-per-ply implementation.
Pandafish starts with the smaller experiment; history ordering remains a
separate candidate if the simpler approach proves useful.

An isolated tagged v0.11 source copy profiled ten cold-table depth-four roots:
start position, three existing timeout regressions, and all six confirmed public
stream cases. It recorded 25,938 explicit main-search beta cutoffs. Of these,
3,480 (13.4%) came from quiet moves; 3,206 of those (92.1%) repeated a quiet move
already seen at the same remaining depth. At a fixed root depth that corresponds
to the same search ply. This count excludes TT and quiescence cutoffs. Global
repetition within a tree is an opportunity signal, not a prediction that two
recent slots will catch all repeats or that all positions will benefit.

[Profile results](baseline-profile.json) and [instrumentation/test source](baseline-profile-source.json)
preserve the investigation. To reproduce, unpack tag v0.11 into a temporary
checkout, apply the saved search patch, write the saved test source to
ordering_profile_test.go, adjust its explicit fixture/output paths, and run
`go test -run '^TestOrderingProfile$' -count=1 .`. Instrumentation only adds
counters; it is outside production sources and provides no timing claim.

## Method fixed before match results

Use the routine 100-game v0.11 gate first. Continue only if it passes with
zero forfeits and a score strictly above 50% treating unfinished games as losses.
Then play 100 against Stockfish 19 Skill 1 and 100 against Sacrifice 0.2.0
b528f0ba, serially in that order. No automatic retry, longer confirmation,
or additional Skill 2 run is scheduled for this experiment.

Use the unchanged ten-opening suite, 50 color pairs, 100 ms requested per move,
50 ms protocol allowance, and 400-ply cap. Hardware is an Apple M2 MacBook Air,
AC power reported at start, with all matchups serial and no live bot. Preserve
available [power observations](power-observations.jsonl); they are not continuous
CPU/thermal measurements. The earlier v0.11 run began on battery and switched
to AC, limiting direct interpretation of separate external scores.

```sh
caffeinate -i go run ./tools/benchmark --validate --race --version v0.12-dev \
  --previous tools/benchmark/bin/pandafish-v0.11 --previous-version v0.11 \
  --previous-revision 1a18aaf2a26eb97cc63df0410d8df719f3a36b5e \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 1a18aaf2a26eb97cc63df0410d8df719f3a36b5e-dirty \
  --machine 'Apple M2 MacBook Air; AC power reported at start; all matchups serial; no live bot' \
  --verbose
```

The revision identifies the base plus local changes; the actual binary hash
belongs to the generated summary. [The source snapshot](candidate-source.json)
preserves the engine and focused tests, including newly added files. The two
slots and per-search lifetime were chosen before matches; no parameter sweep
or tuning to scores is planned. Keep unfavorable or blocked results too.

## Correctness and efficiency checks

Tests cover hint replacement, ply isolation, bounds, captures/promotions/en
passant exclusion, stale hints, preserved legal move sets, tactical/hash priority,
cancellation, alpha-beta windows, exact root scores and all equally best moves.
A nil hint table supplies an unordered reference with otherwise identical
search and transposition behavior. The existing exhaustive reference still
independently checks the underlying fixed-depth search.

[Four deterministic node comparisons](node-audit.json) preserved exact scores
and best-move sets at depth four while reducing nodes on all four cases:
216,241 to 138,879 in total (35.8%). Individual reductions range from 2.8% to
51.1%. This measures visited nodes on those positions, not a general runtime
speedup, timed depth gain, or strength estimate. Reproduce the artifact with
`PANDAFISH_NODE_AUDIT=/tmp/node-audit.json go test -run '^TestKillerNodeAudit$' -count=1 -v .`.

[Saved-position probes](position-audit.json) and [Go source](position-audit.go.txt)
compare both engines on the six public confirmed stream cases, including the
three public holdouts. Each gets an independent cold-process depth-three and
100-ms request, serially. All 24 probes were legal, all six fixed-depth scores
matched, and completed timed depths were equal in all six cases. These single
samples do not establish a general timed depth improvement. No engine code was
tuned after inspecting them.

The candidate is native ARM64 and passed the local lichess-bot config/UCI
wrapper with 12 legal plies. The active v0.11 release executable was restored
afterward; both private configs and greetings remain v0.11. Public example
greetings match v0.12-dev and loaded with engine paths resolved. No listener or
online game was started.

## Match results

All 300 games completed serially, with zero forfeits or unfinished games.

| Opponent | Wins / draws / losses | Score |
| --- | --- | --- |
| v0.11 | 34 / 49 / 17 | 58.5% |
| Stockfish 19 Skill 1 | 53 / 9 / 38 | 57.5% |
| Sacrifice 0.2.0 b528f0ba | 1 / 33 / 66 | 17.5% |

The previous-release screen passed. Its 8.5-point lead is promising screening
evidence, not a precise Elo estimate or proof of a general gain. The earlier
v0.11-dev external scores were 59.5% at Skill 1 and 18% against Sacrifice.
The present scores are essentially flat relative to those separate runs;
randomized ties, repeated openings, and different power histories prevent
attributing their small differences entirely to this increment. Keep these
results rather than rerunning until the external scores increase.

**Assessment:** the four deterministic node checks show reduced work with
unchanged exact results on those positions. The head-to-head screen supports
keeping this small ordering experiment for review. Neither the six timed
saved-position probes nor the external matches demonstrate a broad depth or
strength increase. This is a useful efficiency increment, not a breakthrough.
No extra Skill 2 screen, automatic longer confirmation, or live game was run.

### Timing limitation

[Timing review](timing-review.json) finds 29 of 15,816 candidate responses past
the requested 100 ms, with a maximum of 132.077 ms. All remained within the
runner's 150-ms request-plus-protocol allowance. Opponents also had long
responses, including the v0.11 baseline reaching 121.234 ms in this run.
These are end-to-end local protocol measurements, including position submission,
not engine-only deadline metrics. No scheduling/GC/utilization trace was saved
to establish the cause. Preserve this limitation; do not claim perfect wall-clock
deadline compliance or widen budgets to hide it. The separate v0.11 screen's
lower maximum cannot establish a causal timing regression across changed runs.

Formatting, Go race tests, vet, build, local config/wrapper checks, and the
focused behavior tests passed. The user approved promoting this efficiency increment to v0.12 after reviewing
these results. Release identity and integration checks are separate from the
development-strength games.

- [Complete settings, identities, checks, and results](summary.json)
- [Previous-release games](previous.pgn)
- [Stockfish games](stockfish.pgn)
- [Sacrifice games](sacrifice.pgn)
- [Supplemental verification](verification.json)

## Release preparation

Only version identity and documentation changed after the screen. Formatting,
Go race tests, vet, native ARM64 build, and a 12-ply local lichess-bot wrapper
check passed for the v0.12 release identity. Six additional serial integration
smoke games (two each against v0.11, Skill 1, and pinned Sacrifice) completed
with zero forfeits or unfinished games outside tracked strength history.
Both private config greetings were updated to v0.12 and loaded successfully;
credentials, challenge policies, and 0644 permissions were preserved. No online
game or listener was started. The development games and their timing limitation
remain unchanged, and are not relabeled as tests of the released executable.
