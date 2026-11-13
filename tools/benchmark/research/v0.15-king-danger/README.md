# Threat-aware king safety: experiment and findings

Published in **v0.14.4 as a research archive**. Experimental source and tests are
replayable data outside the normal build. Released playing behavior stays the
same as v0.14.3.

## Decision before tests and match scores

Start from v0.14.3, with its released search and evaluation. None of the archived
PVS/reduction, board-aware evaluation bundle, or null-move policy is enabled.
This attempt changes what the engine values, not its search depth or pruning.
No commit or PR until the completed assessment is reviewed.

## Five positional improvements, ranked

These are low-confidence planning estimates of head-to-head Elo gain over the
current engine, not measured effects, confidence intervals, or forecasts derived
from reference-engine strength. The central guesses only express priority; the
wide scenario ranges explicitly allow no gain or a regression. There is no
empirical basis for a precise estimate before controlled testing. Do not add the
numbers together. Fifty-game screens cannot reliably resolve these differences.

| Rank | Improvement | Rough central guess | Plausible planning scenarios | Why consider it |
| --- | --- | --- | --- | --- |
| 1 | Threat-aware king safety | +75 Elo | -75 to +200 Elo | The current score has no enemy attack pressure around the king; test a direct board interaction rather than another generic placement preference. |
| 2 | Piece threats and loose-piece safety | +60 Elo | -75 to +175 Elo | Distinguish vulnerable, attacked and poorly defended pieces in quiet leaves; avoid double-counting captures already resolved by quiescence. |
| 3 | Coherent middlegame/endgame piece-square evaluation | +50 Elo | -50 to +150 Elo | Current formulas are coarse, and only the king is tapered. Improve phase-sensitive development, queen placement and endgame activity without counting pseudo-mobility. |
| 4 | Narrow pawn structure and supported passers | +40 Elo | -60 to +125 Elo | Isolate weaknesses and realistic promotion potential; the archived combined evaluator did not validate its weights or advanced-passer incentives. |
| 5 | Rook files and minor-piece coordination | +20 Elo | -40 to +80 Elo | Open/semi-open files, bishop pair and useful minor-piece posts are missing. These are probably narrower gains than king danger or piece threats. |

The recommendation is number one, not a claim that a 75-Elo gain is likely or
that a large improvement must exist in the first weights tried. Promotion still
requires actual game evidence. Material, king placement and legal tactics remain
important; this term supplements them rather than declaring every attack a mate.

Mate-distance scoring and repetition recognition are separate conversion/search
correctness topics, not positional bonuses. Do not imply that this evaluation
term fixes equal mate scores or recognizes repeated game history.

## Reference ideas and limits

- [Sunfish's pinned evaluation](https://github.com/thomasahle/sunfish/blob/e396948acbd80cdcbad3dccfb07e878fba9df8f4/sunfish.py)
  demonstrates compact piece-square evaluation and an endgame king preference.
  It does not establish the Elo value of our new king term or implement this
  same explicit king-zone model.
- [Sacrifice's evaluation](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/evaluate.rs)
  includes pawn shelter and separate [attack pressure](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/attack.rs).
  Its attack code derives pressure from legal destinations of the side to move;
  our geometric maps instead score attacks symmetrically for both colors. Do not
  describe its legal-move destinations as a complete symmetric defender map.
- [Stockfish's pinned evaluation](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/evaluate.cpp)
  uses learned neural evaluation, not this small handcrafted model. It supports
  board-aware evaluation as a direction, not our particular coefficients. NNUE,
  training data and learned weights are outside this increment.

This is original Go code, not copied tables or evaluation routines. The previous
pawn-shelter/mobility bundle did not test enemy king-zone pressure in isolation.

## Observations and diagnosis limits

The released `placement(King, ...)` depends only on square and remaining
non-pawn material. It cannot distinguish identical king placement under different
enemy attack pressure. Saved games and earlier live notes motivate checking this
blind spot, but outcomes or an early king move do not prove the king term caused
a loss. Some older live flags changed at larger Stockfish budgets. Keep that
caution when using examples.

Before match play, inspect selected saved positions with unrestricted/restricted
full-strength Stockfish and cold released/candidate probes at fixed depth and
100 ms. Preserve negative examples and changed depth as well as favorable moves.
The position set is selected diagnostic data, not a hidden test set, and will
not be used for weight tuning. The fixture's provenance and selection rule must
be explicit in the result record.

## Fixed initial policy

Score geometric attacks into the eight neighboring squares of each king.
Sliding attacks stop at and include the first occupied square. Pawn diagonals
and knight jumps use their normal geometry. Friendly defended blockers are
included; pins and full legal move constraints are not modeled.

Count each attacking piece once, plus weighted attacked zone squares: pawn 1,
knight/bishop 2, rook 3, queen 5. Penalty is units squared divided by four, halved
for a lone attacker, capped at 240 cp, then tapered by the existing material phase
out of 24. Only a side that still has a queen contributes attack pressure in
this first version. Queenless king attacks and detailed safe checks are known
omissions. No pawn-shield, mobility, passer, bishop-pair, material-value or new
placement term is added. No weights are fitted to the diagnostic set.

Keep search, quiescence, move ordering, clocks, root ties and table layout fixed.
All inputs are board-derived, so the transposition identity remains applicable.
The static score remains antisymmetric across perspectives and color/rank mirrors.
Terminal scoring remains dominant. Geometric attacks are an evaluation hint, not
a proof of danger, legal checking moves or a mating attack.

## Validation and promotion

Run geometry/blocker/edge tests against the independent library attack maps,
hand-counted pressure and tapering checks, color/perspective symmetry, terminal
and material invariants, root immutability, zero-allocation and static-call cost
measurements. Then normal/race tests, vet, ARM64 build and offline wrapper.

Run the normal 50 games each against tagged v0.14.3, Stockfish Skill 2 and pinned
Sacrifice, all three matchups concurrent. Read compact summaries. Preserve all
150 games even if the previous-release criterion fails. A small gate pass alone
is not the substantial improvement sought. Use historical 60%/34% external scores
as context, not simultaneous controls or a numerical promotion rule. No favorable
rerun, longer screen or parameter sweep to rescue a result. Check in before any
candidate commit or PR, whether the result is promising or negative.

## Pre-match observations

A deterministic extraction selected the first eight king and first eight queen
decisions by released v0.14.2 at plies 12-48 from the saved previous-version
match PGN, at most one decision of each kind per game. Selection used no referee
score or new-candidate output. Six recorded decisions were at least 100 cp worse
than the unrestricted reference at 100,000 nodes. These are finite-budget flags,
not confirmed causes or a representative error rate. In particular, castling
long at round 10 ply 13 remained the choice of both engines in both probes, with
a 291-cp estimated loss. This candidate does not solve that example.

All 64 released/candidate probes were legal. The summed signed estimated
centipawn losses at depth four were 1,132 released versus 1,186 candidate; at
100 ms they were 1,380 versus 1,309. Both timed policies completed depth four in
six cases and depth three in ten. These are different diagnostics, not direct
Elo evidence. The mixed position results provide no pre-match strength claim.
Weights remain fixed for the complete game screen.

The static-call control measured about 0.56-0.67 microseconds for the released
evaluator and 1.44-1.50 for the candidate on starting/dense positions. Both
allocated zero heap bytes. That extra cost is real; these call timings are not
a measurement of total search slowdown. Timed game results must justify it.

## Completed screen: hold

[Full results and assessment](../../results/v0.15-king-danger-dev-vs-v0.14.3-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md).
The candidate scored 50% against v0.14.3, 46% against Stockfish Skill 2 and 35%
against pinned Sacrifice, with all 150 games finished and no forfeits. The direct
criterion failed at exactly 50%; the runner's exit status 1 records that decision.
Do not promote. The initial +75-Elo guess is not supported by this screen. The
static gap exists, but neither selected-position probes nor games demonstrated
that this term earns its cost. Next, diagnose actual flagged decisions more
deeply before choosing another positional bonus. No commit, push, PR or release.

## Replay and release verification

Restore the exact five-file candidate overlay onto tagged v0.14.3:

```sh
repo_root="$PWD"
replay_root=$(mktemp -d)
record_dir="$repo_root/tools/benchmark/results/v0.15-king-danger-dev-vs-v0.14.3-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba"
git archive v0.14.3 | tar -x -C "$replay_root"
cp tools/benchmark/research/v0.15-selective-search/restore-snapshot.go.txt /tmp/pandafish-restore-snapshot.go
go run /tmp/pandafish-restore-snapshot.go "$record_dir/candidate-source.json" "$replay_root"
(cd "$replay_root" && go test ./...)
(cd "$replay_root" && go build -o /tmp/pandafish-king-danger-replay .)
```

The original fixture remains at this report's `position-inputs.json`. The
results directory preserves the Go extraction/referee drivers, original labels,
hashes and all 150 games. Reproduction recovers policy, not identical scheduling
or random choices. [Release checks](research-release-checks.json) separately
record normal/race tests, vet, ARM64 build, offline wrapper, replay tests/build,
and six concurrent integration smoke games of the restored release engine.
Smoke scores are not strength evidence. No experimental evaluation ships.
