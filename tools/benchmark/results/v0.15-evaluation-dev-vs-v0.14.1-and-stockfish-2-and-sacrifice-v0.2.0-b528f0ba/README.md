# v0.15 positional-evaluation screen

## Decision: hold this candidate

The board-aware evaluator passed the direct previous-release criterion, but did
not deliver the substantial, consistent strength gain sought for v0.15. Keep
v0.14.1 as the active engine. This is a local experiment, not a release or a
claim that the new evaluator is generally stronger.

| Opponent | Wins | Draws | Losses | Score | Earlier released-policy score |
| --- | --- | --- | --- | --- | --- |
| Pandafish v0.14.1 | 18 | 21 | 11 | 57% | Not applicable |
| Stockfish 19 Skill 2 | 20 | 5 | 25 | 45% | 60% |
| Sacrifice v0.2.0 b528f0ba | 7 | 25 | 18 | 39% | 34% |

Each matchup had 50 completed games and zero forfeits. Score counts a draw as
half a win. The external reference comes from the
[v0.13 concurrency calibration](../v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md),
whose playing algorithm v0.14.1 retains. Those are separate runs, not simultaneous
paired control matches, so their differences are descriptive. Random ties,
Stockfish's weakened play, timing, and the small opening/sample set introduce
variation. A 57% direct score is promising screening evidence, not proof of a
strength gain. Conversely, 45% versus Stockfish does not prove a permanent
regression. It does make this a poor candidate to promote as a substantial win.

The predeclared descriptive goal was at least ten percentage points better on
one external opponent, without losing more than five on the other. Actual
changes were minus fifteen against Stockfish and plus five against Sacrifice.
No weights were changed after seeing these results. No longer match, favorable
retry, or second full-strength variant was run to rescue the outcome.

## Hypothesis and implementation

The [fixed plan and coefficients](../../research/v0.15-positional-evaluation/README.md)
add doubled/isolated/passed-pawn scoring, real king pawn shelter and nearby open
files, plus modest knight/bishop activity outside enemy pawn attacks. Material,
original placement, alpha-beta, quiescence, transposition-table policy, root tie
handling, and clocks are unchanged. The failed selective-search experiment is
not enabled. No database, learned weights, or queen-mobility reward was added.

`positional.go` builds one stack-local board summary per static evaluation.
Passed-pawn bonuses grow with advancement and endgame phase; blockers halve the
bonus. Shelter fades out in endings. Activity is pseudo-legal availability,
not a count of fully safe legal moves. Focused tests cover both colors,
hand-counted structures, phase, component accounting, symmetry, independent
attack-map checks, terminal dominance, root immutability, and zero allocations.
The en-passant test now compares quiescence with the positional leaf rather
than assuming a surviving advanced pawn has only its material value.

## Position diagnostics and cost

The 19 selected positions comprise 13 existing public/saved positions and six
synthetic quiet/endgame cases. They are not an unseen or representative strength
suite. Cold depth-four and independent 100 ms probes compare v0.14.1, an offline
pawn/king-only ablation, and the full candidate. All 114 probes returned legal
moves without errors. Fresh processes avoid table contamination between policies.

Stockfish 19 at full strength supplied a 100,000-node reference and restricted
move evaluations, using one thread, 64 MiB hash, and a fresh hash per query.
Signed estimated centipawn-loss sums below cover 15 centipawn-scored positions;
four mate-scored cases are excluded. Finite-node analysis is an estimate, not a
best-move oracle. These sums are diagnostic aggregates, not an Elo measure.

| Policy | Depth-four loss sum | 100 ms loss sum | Depth-four nodes | 100 ms nodes |
| --- | --- | --- | --- | --- |
| Released evaluator | 1691 | 1510 | 275002 | 374510 |
| Pawn/king-only | 879 | 550 | 247727 | 381833 |
| Full candidate | 657 | 637 | 268688 | 365425 |

The full timed candidate improved six centipawn-scored choices, matched eight,
and worsened one relative to the baseline. Across all 19 cases, completed depth
was higher once, equal fourteen times, and lower four times. The pawn/king-only
control was higher once, equal sixteen times, and lower twice. It was a position
diagnostic only: no full match result establishes its strength.

Isolated static-call samples were approximately 0.47-0.56 microseconds for the
released evaluator, 0.64-0.76 for pawn/king-only, and 0.83-0.97 for the full
candidate. All measured zero heap allocations. The full evaluation costs about
1.75-1.8 times as much per call. This is not a measurement of whole-search
slowdown: changed scores also change pruning and the searched tree. Better
selected-position judgments therefore did not guarantee better timed games.

## Reproducibility and timing

Base revision: `6d4e122a73e006b05738f212f58dc442d7c6321e` (v0.14.1).
Candidate identity: `Pandafish v0.15-evaluation-dev`.
Candidate SHA-256:
`7e94f9306d43d2c991a75e391b7253fd500e1af86c0543c2981cf433514c47cf`.
The position probes, matches, and wrapper used that same candidate binary hash.

```sh
caffeinate -i go run ./tools/benchmark --validate --race \
  --version v0.15-evaluation-dev \
  --previous tools/benchmark/bin/pandafish-v0.14.1 \
  --previous-version v0.14.1 \
  --previous-revision 6d4e122a73e006b05738f212f58dc442d7c6321e \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 6d4e122a73e006b05738f212f58dc442d7c6321e-dirty \
  --machine 'Apple M2 MacBook Air; AC power reported at start; three concurrent matchups; no live bot' \
  --verbose
```

Defaults were 25 opening pairs, 100 ms per move, 50 ms timeout grace, 400-ply
cap, and three concurrent matchups. Stockfish used Skill 2, one thread, 16 MiB
hash, UCI_LimitStrength false, and no Syzygy path. The exact opponent hashes,
options, opening hash, and response times are in `summary.json`. The match
window was approximately 527 seconds. AC power was observed before and after;
continuous temperatures or throttling were not recorded. No other engine
experiment ran alongside the match suite.

Across 7,641 candidate responses, ten exceeded the nominal 100 ms and none
exceeded the 150 ms timeout threshold. Maximum response was 143.315 ms; per-match
p95 was approximately 96.2-96.4 ms. These wall-clock observations include the
runner/engine exchange, not live-client submission or network latency.

## Checks and retained records

Formatting, normal Go tests, race-enabled Go tests, vet, and the ARM64 build
passed. The offline lichess-bot config/wrapper check returned 12 legal plies
with the candidate identity. The public example loaded with resolved local
paths and a placeholder token. No live games were played. Both private configs
were unchanged, and the root executable was restored byte-for-byte to v0.14.1.

- `summary.json`, `previous.pgn`, `stockfish.pgn`, `sacrifice.pgn`, `games.pgn`:
  all match records, original labels, hashes, and response data.
- `screen-summary.json`: compact scores, timing, and end power observation.
- `candidate-source.json`: exact tested Go files and fixture, relative to the base.
- `pawn-king-source.json`: ablation overrides relative to that candidate.
- `position-inputs.json`, `position-audit.json`, `position-summary.json`:
  input cases, raw position/referee records, and compact diagnostic totals.
- `position-audit-*.go.txt`: archived Go audit drivers, excluded from normal builds.
- `evaluation-cost.txt`, `preflight.json`, `checks.json`: raw call-cost samples,
  initial power observation, and final verification record.

## What to try next

Keep the released search as the control. Isolate pawn structure and king shelter
from minor mobility before considering more evaluation complexity. The ablation
has a lower per-call cost and lost less depth in the selected timed probes,
which gives it a concrete hypothesis to test, not permission to assume a gain.
Review a small set of changed decisions with component breakdowns to investigate
possible overvalued passers or shelter incentives. The present data does not
identify either as the cause of the Stockfish result.

Predeclare the next change and run the same routine screen once. Do not sweep
weights against these openings or treat this selected-position set as a training
objective. A genuinely broader positional test corpus could improve diagnostics;
it should supplement, not replace, timed game comparisons. No new implementation,
publication, or release follows automatically from this report.
