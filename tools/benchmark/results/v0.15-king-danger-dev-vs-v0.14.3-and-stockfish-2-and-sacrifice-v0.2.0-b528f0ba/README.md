# v0.15 threat-aware king safety

## Decision: do not promote

The isolated king-danger term did not pass the direct previous-release criterion.
Keep v0.14.3 as the released playing policy. The candidate is preserved in the v0.14.4 research archive and
is excluded from the normal engine build. Its initial positional hypothesis was plausible, but these results
do not support adding the complexity or claiming the estimated gain.

| Opponent | W / D / L | Score | Unfinished | Forfeits |
| --- | --- | --- | --- | --- | --- |
| v0.14.3 | 13 / 24 / 13 | 50% | 0 | 0 |
| Stockfish Skill 2 | 21 / 4 / 25 | 46% | 0 | 0 |
| Sacrifice v0.2.0 b528f0ba | 7 / 21 / 22 | 35% | 0 | 0 |

Each matchup had 50 games, with draws worth half a win. All three matchups ran
concurrently. Exactly 50% against v0.14.3 fails the greater-than-50% criterion;
the command returned exit status 1 for that score decision, not a crashed engine
or unfinished suite. All 150 games and the original labels remain preserved.

Historical released-policy scores in the
[concurrent calibration](../v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
were 60% against Skill 2 and 34% against Sacrifice. This candidate's 46%/35%
provides no substantial external gain. These are separate screens, not paired
simultaneous controls. Fifty-game samples and random ties do not prove permanent
regression, equivalence, or a precise Elo difference. No rerun, weight tuning or
longer screen was used to obtain a favorable result.

## Hypothesis and implementation

The [pre-match report](../../research/v0.15-king-danger/README.md) documents five
ranked positional improvements, rough initial Elo guesses with wide scenario
ranges, reference-engine sources, and the fixed policy. The +75-Elo central
planning guess for king safety is not supported by this completed screen. Do
not cite it as a measured improvement or a calibrated prediction.

Released king placement sees king square and phase, not enemy attack pressure.
The candidate adds original Go geometric attack maps with blockers, counts
weighted attacks on adjacent king squares, increases the penalty for combined
pressure, caps it at 240 cp, and tapers by material phase. Only queen-backed
attacks are included. Pins, safe checks, defenders, queenless attack danger and
actual mating threats are not modeled. The term is a preference, not a tactical
oracle. It does not change material values or existing placement formulas.

Search, quiescence, move ordering, clocks, root ties, table layout and dependencies
remain unchanged. No held PVS/reductions, pawn/shelter/mobility bundle, or null
pruning is active. No new UCI option or runtime evaluation switch is added.

## Actual saved-decision audit

The extraction rule was fixed before referee scores: first eight king decisions
and first eight queen decisions by the released v0.14.2 side in the prior local
previous-version PGN, at plies 12-48, at most one of each kind per game. This
selected diagnostic is neither representative nor a hidden evaluation corpus.
The original PGN hash, exact positions, game/ply references, full UCI histories
and recorded moves are in `position-inputs.json`.

Stockfish 19 at full strength analyzed the unrestricted position, the original
recorded move, and every distinct probe move at 100,000 nodes per independent
query, with one thread, 64 MiB hash, MultiPV 1 and a cleared hash/new game each
query. The parser associates exact scores with their reported PV moves.
Six recorded decisions had estimated losses of at least 100 cp. These are
finite-budget flags, not confirmed error causes or a comprehensive game audit.
They were not rechecked at a larger budget, and this work did not analyze every
move of every game or draw conclusions about all conversion failures.

| Saved decision | Recorded move | Stockfish reference | Estimated loss |
| --- | --- | --- | --- |
| Round 1, ply 34 | d8e7 | a7a5 | 188 cp |
| Round 2, ply 23 | d1d3 | d5f6 | 108 cp |
| Round 4, ply 35 | e1f1 | b1c1 | 104 cp |
| Round 5, ply 20 | e8c8 | e6d7 | 157 cp |
| Round 6, ply 29 | d1e2 | c2c3 | 120 cp |
| Round 10, ply 13 | e1c1 | e5e3 | 291 cp |

The candidate retained most of these choices. The castling choice in round 10
remained unchanged in both candidate probes, so this implementation does not
solve that conspicuous example. Do not assume that a king move is a king-safety
evaluation error: legal tactical resources, development and horizon effects can
change the better alternative.

All 64 cold-process released/candidate probes were legal and completed without
errors. Each engine had separate depth-four and 100 ms probes. The finite-budget
signed centipawn-loss sums are diagnostic aggregates, not an Elo measure:

| Policy | Depth-four estimated loss sum | 100 ms estimated loss sum | Timed completed depths |
| --- | --- | --- | --- |
| Released v0.14.3 | 1132 | 1380 | depth 4: 6; depth 3: 10 |
| King-danger candidate | 1186 | 1309 | depth 4: 6; depth 3: 10 |

The candidate was slightly worse at fixed depth and slightly better timed on
this selected set. It did not present a convincing diagnostic strength gain
before match play. The weights were held fixed through all 150 games.

## Cost, timing, and verification

Static-call controls on starting/dense positions measured about 0.56-0.67
microseconds released versus 1.44-1.50 candidate, with zero heap allocations for
both. This is a roughly 2.2-2.7 times call-cost increase, not an equivalent
measurement of total search slowdown. Additional evaluation work can change the
searched tree as well as time per node. Better board awareness must earn its cost;
this candidate has not done so in the present screen.

The match window was approximately 527 seconds on the M2 MacBook Air. Early and
final power observations reported AC. Continuous temperature/throttling data
were not collected. Across 7,927 candidate responses, two exceeded nominal
100 ms and none exceeded the 150 ms timeout threshold. Maximum was 113.303 ms;
matchup p95 was about 96.3-96.7 ms. There were zero forfeits or unfinished games.
These are local runner measurements, not network/client submission latencies.

Focused tests independently compare attack maps with the chess library, cover
blockers/edges, hand-count pressure, lone/combined attacks, cap/phase behavior,
queenless exclusion, material/component invariants, mirror/perspective symmetry,
root immutability and zero allocations. Normal tests, race tests, vet, ARM64 build,
whitespace checks and the offline lichess-bot 12-ply wrapper check are recorded in
`checks.json`. The public example loaded with placeholder credentials and local
resolved paths. The active root binary was restored byte-for-byte to v0.14.1;
private configs were unchanged. No live client or online game was started.

## Reproduce and inspect

Base: `1470b161423263517c567e3fed934234487e553d` (v0.14.3).
Tested identity: `Pandafish v0.15-king-danger-dev`.
Tested SHA-256: `b02aeef2de5dc1be28b8cf42da17505e37334db4462bf91ececaee51e3471842`.
The position probes, matches and wrapper used this same engine hash.

```sh
caffeinate -i go run ./tools/benchmark --validate --race \
  --version v0.15-king-danger-dev \
  --previous tools/benchmark/bin/pandafish-v0.14.3 --previous-version v0.14.3 \
  --previous-revision 1470b161423263517c567e3fed934234487e553d \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 1470b161423263517c567e3fed934234487e553d-dirty \
  --parallel-all --machine 'Apple M2 MacBook Air; three concurrent matchups; no live bot' \
  --verbose
```

Defaults: 25 color pairs per opponent, 100 ms per move, 50 ms timeout grace,
400-ply cap. `summary.json` stores exact opponent identities, hashes, options,
opening hash and all response measurements. `screen-summary.json` is compact.
All original games are in per-opponent PGNs and combined `games.pgn`.

`candidate-source.json` overlays five Go files onto a temporary `git archive
v0.14.3` checkout. Use the shared
[restore helper](../../research/v0.15-selective-search/restore-snapshot.go.txt),
then test and build that temporary checkout. Referee drivers and extraction code
are archived as `.go.txt`, excluded from normal engine builds. Copy the three
`position-audit-*.go.txt` drivers to temporary `.go` files and `go run` them from
the restored root with an output JSON path. The driver reads the original
research fixture path. `evaluation-cost.txt` preserves raw cost samples;
`position-audit.json` and `position-summary.json` retain all positive and negative
position evidence. Reproduction recovers policy/settings, not identical random
moves or scheduling. A hash manifest covers the retained files.

## Next assessment

Do not work through the five-item list mechanically or add more king weights to
rescue this score. The ranking was a set of hypotheses, not evidence of missing
Elo waiting to be collected. A real static-score blind spot alone did not identify
the dominant source of lost games.

Before implementing the second positional candidate, diagnose the flagged moves
with larger-budget reference confirmation, explicit tactical lines and controlled
released-engine root/leaf comparisons. Determine whether the engine misvalues a
quiet position, misses a tactical continuation, or repeats an otherwise winning
position. Then require a focused fixture showing the proposed correction changes
the relevant decision. Loose-piece safety is a candidate only if that diagnosis
supports it. Mate distance and repetition remain separate conversion corrections;
this term has not addressed them. Check in before any candidate commit or PR.
