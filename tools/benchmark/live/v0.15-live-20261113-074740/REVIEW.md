# Review and next strength hypotheses

**Next recommendation: a small, board-aware piece-safety evaluation**, starting
with attacked, inadequately defended pieces at quiet search leaves. Establish
that it explains and improves saved decisions before paying for full matches.
This is an unvalidated hypothesis for the next engine increment, not a feature
included in v0.15.1 or a promise of a large strength gain.

The live test validates the observed integration much more strongly than it
identifies one winning evaluation formula. All 12,091 recorded searches
completed legally with final metrics, and there were no own time-forfeit
losses. There are substantial evaluation and tactical mistakes under ordinary
clocks. Another broad search optimization is not the first recommendation.

## What the larger sample establishes

The single 58h54m13s session contains 269 unique games, 265 completed, against
170 BOT accounts. It ran in a terminal on the M2 Air with light browsing/music,
without streaming. There is no second session to combine and no duplicate game
ID to discard. Identical openings or even identical games with different IDs
remain real occurrences, but they are correlated evidence, not independent
discoveries. Four unfinished records remain unfinished.

The engine won 77, drew 28 and lost 160. Five wins were opponent time forfeits;
none of its losses were recorded as time forfeits. The rated 3+2 subset scored
52 wins, 20 draws and 116 losses from 188 completed games, with one unfinished.
These results describe this opponent pool. They do not establish an Elo gain
over an earlier release, or imply that all losses were preventable.

Recorded search intervals show no overlap and a peak of one searching process.
This differs from the documented v0.10 duplicate-worker incident. The client
has reconnect and challenge-limit errors, discussed in [CLIENT_NOTES](CLIENT_NOTES.md),
but there is no evidence here that a new concurrency bug caused engine losses.

There were 9,675 searches slightly past the internal hard deadline; the largest
overrun was 1.132 ms. Count and size must be considered together. Response p95
was 12.211 s; maximum 37.462 s followed a 37.461 s budget. Only 132 clock-based
requests had less than ten seconds left, and 13 of those moves were significant
survey flags. Most flagged decisions occurred with more clock available.
No temperature, clock-frequency or CPU utilization telemetry was collected.

## Survey and confirmation

The referee was Stockfish 19, Skill 20, Elo limiting off, one thread, 64 MiB
hash, MultiPV 1. Each independent unrestricted/played-move query used the same
node budget and a fresh hash with full game history. This is not the Skill 2
benchmark opponent setting. Stored raw output and PVs permit rechecking each
score; finite-node scores can disagree and are not ground truth.

| Classification | Bot moves |
| --- | --- |
| Estimated loss 100 to 299 cp | 1,137 |
| Estimated loss at least 300 cp | 224 |
| Played search reports losing mate where reference does not | 129 |
| Reference reports winning mate but played search does not | 13 |
| Both searches report mate scores | 1,422 |
| Negative cp difference / search disagreement | 940 |
| Smaller or equal cp difference | 8,217 |
| Total | 12,082 |

The 1,361 cp flags are 11.3% of bot moves. The clamped cp-only mean was 40.4 cp
across 10,518 comparisons; mate transitions remain separate. This is neither
Lichess accuracy nor a directly comparable ACPL/Elo measurement. The earlier
v0.10 archive's 7.7% cp-flag rate used different opponents, clocks, positions
and a much smaller sample; the increase does not prove a version regression.

There are 109 survey flags where a roughly playable position deteriorates by
the investigation's stated threshold. Phase buckets are descriptive: before
ply 24 is opening; otherwise ply 60 onward or at most ten pieces is late;
the remainder is middle. They contain 31, 58 and 20 such flags respectively.
This suggests investigating middlegame decisions, rather than assuming every
problem is an opening or a mate-conversion failure.

Thirty cases from different games were frozen before their 500,000-node
results: eight playable-deterioration cases in each phase and six conversion
opportunities. Full history plus played move was deduplicated for selection.
All thirty retained at least 100 cp loss or an allowed-mate report: twenty
100 to 299 cp, nine at least 300 cp and one allowed mate. Some fell below the
original 200 cp selection threshold or had a worse reference score at the
larger budget. Do not say that every exact survey threshold was confirmed.
The 292.0 cp mean of the 29 cp cases reflects deliberately selected mistakes,
not session-wide play. Selection and all changed scores remain recorded.

## Decisions that deserve explanation

Scores below favor the side making the decision and come from the 500,000-node
check. Ply is zero-based. The reference is a finite-budget alternative, not
necessarily the only acceptable move.

| Game / ply | Played / reference | Reference / played | Estimated loss |
| --- | --- | --- | --- |
| [16G1f9rx](https://lichess.org/16G1f9rx) / 9 | Nd4 / O-O | -37 / -497 cp | 460 cp |
| [O4wRsvjl](https://lichess.org/O4wRsvjl) / 11 | Kf7 / Bf5 | -57 / -328 cp | 271 cp |
| [HhPI3vCx](https://lichess.org/HhPI3vCx) / 25 | Qd6 / Qe8 | -82 / -459 cp | 377 cp |
| [B1WTqhJy](https://lichess.org/B1WTqhJy) / 25 | Bxf3 / f6 | -15 / -271 cp | 256 cp |
| [DsqFd5Kj](https://lichess.org/DsqFd5Kj) / 63 | Rd5 / Rd1 | -7 / -291 cp | 284 cp |
| [L0mwQLiR](https://lichess.org/L0mwQLiR) / 64 | Nxh4 / Ra7 | +313 / -6 cp | 319 cp |
| [PNfQoYP0](https://lichess.org/PNfQoYP0) / 64 | h4 / Bf3 | +185 / 0 cp | 185 cp |
| [rpk4dbYK](https://lichess.org/rpk4dbYK) / 27 | Kxg5 / Kf6 | +106 cp / mate -11 | Mate transition |

Ra7 is a rook move, not a passed-pawn advance. Kxg5 captures a bishop, not a
pawn. That final game was actually drawn: allowing a referee mate report is
not the same as the opponent executing the mate in the live game.

## Replay evidence changes the recommendation

The exact recorded v0.15 binary was given the full history and just the two
listed root choices, at its live completed depth and one extra ply. A third
search valued the other move alone at the extra depth to distinguish ties.
All 24 searches completed legally. These cold-process, restricted-root tests
do not reconstruct the live warm table or full-root search.

Only L0mwQLiR switched from the played move to the reference at the extra ply.
B1WTqhJy already selected the reference at live depth in this diagnostic;
PNfQoYP0 had equal scores at the extra ply. The other five still chose the
played move. Examples include Qd6 valued at +159 cp versus Qe8 +93, and Kxg5
at +400 versus Kf6 +215. This is stronger motivation to inspect leaf scores
and threats than to demand one more ply everywhere. It does not prove that
an evaluation penalty alone fixes those tactical sequences.

Two small counterfactuals were retained as analysis data, without changing the
released engine:

- The king formula awards a shelter-side bonus even off the home rank. At full
  phase it gives e1/e8 -8 cp and f2/f7 0 cp. Restricting that bonus to the home
  rank changes O4wRsvjl to Bf5 at depth six, but it returns to Kf7 at depth
  seven; Kxg5 is unchanged. The incentive is real, but the one-line repair is
  not a robust solution. The full immediate static score already prefers Bf5
  in O4wRsvjl, so blaming that single term for the searched choice is unjustified.
- Allowing at most two additional quiet checks per quiescence path still chose
  Nd4 at depth seven. Both attempted rpk4dbYK probes exceeded the independent
  20-second watchdog before returning a move. No score or benefit is inferred
  from those timeouts. Naively adding quiet checks is not a validated fix.

See [the complete replay tables and source](analysis/engine-replay/README.md).

## Ranked next increments

The estimates below are **low-confidence planning guesses in head-to-head Elo
over v0.15**, not measurements, statistical bounds or estimates extracted from
the live win rate. The central values express priority, not numerical accuracy.
There is no calibrated model for these gains. Each could yield zero or a
regression; the earlier +75-Elo king-danger guess was not validated by matches.
Do not add the numbers or promote based on them.

| Rank | Hypothesis | Rough central guess | Planning scenarios | Evidence / confidence |
| --- | --- | --- | --- | --- |
| 1 | Board-aware piece safety | +40 Elo | -75 to +150 Elo | Frequent middlegame deterioration and persistent Qd6/Nd4 errors; moderate priority confidence, low gain confidence |
| 2 | Narrow pawn-ending and passer evaluation | +30 Elo | -50 to +100 Elo | Rd5 trades into a bad pawn ending; h4 misses a winning alternative; narrower but concrete cases |
| 3 | Coherent king exposure and shelter | +25 Elo | -75 to +125 Elo | Kf7/Kxg5 and a real placement discontinuity; previous king-danger attempt and new one-line probe caution against another broad bonus |
| 4 | Tapered minor/queen placement | +20 Elo | -50 to +100 Elo | Coarse phase-insensitive preferences, repeated development errors; indirect evidence rather than a demonstrated correction |
| 5 | Bishop-pair and rook-file evaluation | +15 Elo | -40 to +75 Elo | Simple missing interactions, smaller scope; opportunity not measured in this sample |

### 1. Piece safety: first experiment

The present evaluator adds independent material and square preferences. It
does not know whether a piece is attacked by a cheaper piece, whether its
defender is pinned, or whether moving it exposes a second target. Quiescence
settles many captures, but a quiet threat or forced defensive move can remain
beyond the boundary. The Qd6 line includes Rd1 chasing the queen, then Re1,
Bxf6 and Nd5; the Nd4 line exchanges pieces and loses material later. These
are evidence to inspect, not proof that all are simple hanging pieces.

Start with one small board summary: occupancy, pawn attacks and per-side
piece attacks along blocked rays. Score only clearly attacked, inadequately
defended non-king pieces, initially distinguishing pawn attacks and genuinely
undefended targets. Cap the term well below the value of the piece. Avoid a
blanket penalty for every attacked square or subtracting the full value of a
capture already handled by quiescence. Geometric defenders are not necessarily
legal defenders; pins, x-rays and king captures need deliberate handling or
an explicit conservative limitation. Start from clear cases, not a full SEE
search at every leaf.

Before tuning weights, trace the baseline PV through these decisions and
compare its first quiet leaves with the referee continuation. Show that the
new feature detects an actual threat there. Add mirrored, defended, pinned,
blocker and tactical-sacrifice controls. Run fixed-depth and equal-time
full-root checks, preserving all legal alternatives. If Nd4/Qd6 are not
explained, record that and use independently justified piece-safety fixtures;
do not force a hand-picked move with a root-only penalty. A term that only
matches Stockfish's move name without explaining the position is not enough.

Keep material values, main search, quiescence, table policy and clocks fixed.
Evaluate cost at the normal 100 ms benchmark budget. Only a mechanism that
survives these checks should advance to the routine three-opponent screen.
This scope was not in the old pawn/shield/minor-mobility bundle and is distinct
from the unpromoted king-zone attack bonus. Do not resurrect those bundles.

### 2. Pawn endings and promotion potential

After Rd5, c5xd5 e6xd5 leads to a pawn ending which the current material/square
score cannot assess adequately. At depth nine it still prefers Rd5 (+25)
over Rd1 (+22), while the larger-budget reference distinguishes them by nearly
three pawns. Inspect that exact ending before inventing a bonus.

Add only a narrow endgame component first: true passers on the same/adjacent
files, blocked forward squares, friendly/enemy king distance to the promotion
path, and support by another pawn or king. Taper it toward low material and
cap optimistic race bonuses. Side to move, opposition and access squares
matter; a geometric passer is not a guaranteed queen. Test both colors,
blocked and attacked passers, equal-material pawn races and deliberately bad
rook exchanges. PNfQoYP0 supplies a second conversion fixture. L0mwQLiR is an
activity/horizon example, not evidence of a passed-pawn move by itself.

The old combined evaluator's quadratic passer incentives were not validated
in isolation. Do not copy those weights or combine this with piece safety
in the first candidate. Expect benefits mainly in conversions, not every
opening. Fifty-move and repetition correctness must remain separate rules.

### 3. King safety without a generic attack-count gamble

Replace accidental shelter-side credit with board-aware shelter and exposure:
actual protecting pawns, open files near the king and credible enemy heavy-piece
access. Scale danger by opposing queens/rooks and phase so an active endgame
king is not always punished. Preserve legal king activity and distinguish
castling from merely standing on a flank file.

Use O4wRsvjl and rpk4dbYK with non-king controls and safe king-capture examples.
The previous king-zone candidate scored 50% against its prior release and did
not earn promotion; the current one-line change failed at the next depth.
A new candidate must explain why its model behaves differently. Do not add
unbounded checking extensions to mask an inadequate king score.

### 4. Phase-sensitive placement

Taper all pieces rather than only the king. Keep early minor development,
central influence and endgame activity coherent; avoid rewarding a wandering
queen merely for centralization. Base any queen term on exposure or phase,
not a fixed "never move the queen early" rule. A developed piece need not be
the right move when a tactic exists. Review repeated early king/queen/pawn
decisions from the audit, with tactical and opening holdouts. Use original
compact tables or formulas with documented provenance, not copied weights
whose licensing or board orientation is unclear. This is broader and harder
to attribute than one threat term, hence its lower priority.

### 5. Bishop pair and rook files

An original, modest bishop-pair bonus and rook activity on genuinely open or
semi-open files are cheap interactions the current seventh-rank-only rook
term misses. Implement and test them as separate ablations rather than one
large release bundle. Closed pawn positions, trapped bishops and a rook on
an unsafe open file are important negative controls. These are plausible
smaller gains with less direct causal evidence than the first three items.

## Other clear work, below the first five

- **Targeted quiet forcing moves:** quiet checks and defensive resources remain
  a search-coverage gap. The saved probe shows that a simple two-check extension
  can explode. Require bounded cost and a demonstrated tactical correction
  before revisiting; no positive gain is established here.
- **Draw rules:** fifty-move and insufficient-material recognition are separate
  correctness gaps. Audit actual endings and rule precedence, including mate,
  before assigning them a strength priority. This report does not claim that
  they explain the recorded 28 draws.
- **Principal-variation diagnostics:** a debug-only full PV/leaf breakdown would
  make causal diagnosis easier than the current one-move PV. Keep it off the
  normal fast path until its cost is measured. This is tooling, not Elo.
- **Move ordering / table refinements:** useful when a profile demonstrates
  wasted work, but lower priority than the persistent evaluation errors here.
- **Lichess database and additional opponents:** use later as a deduplicated,
  versioned source of training/holdout positions, not an opening-book shortcut
  presented as engine reasoning. Keep source games separate from evaluation
  holdouts. Sunfish remains a possible future benchmark; current controls stay
  the previous release, Skill 2 and pinned Sacrifice.

## Reference-engine lessons and limits

[Sacrifice's pinned evaluator](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/evaluate.rs)
combines material, square preferences, mobility, shelter, activity and separate
[attack pressure](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/attack.rs).
It is the closest source of small hand-written interaction ideas. Its strength
does not establish which term helps Pandafish, and geometric attacks require
care around pins and safety. Borrow the question to ask, not a bundle of weights.

[Sunfish's pinned source](https://github.com/thomasahle/sunfish/blob/e396948acbd80cdcbad3dccfb07e878fba9df8f4/sunfish.py)
illustrates how much compact material/square evaluation plus search can achieve.
It does not mean Pandafish should replace its working alpha-beta structure or
ignore the concrete board interactions its own games expose.

[Stockfish's pinned search](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp) and
[evaluation](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/evaluate.cpp)
combine a far more developed search with learned evaluation. Its full-strength
offline answers help locate failures, but do not identify one minimal patch or
justify importing NNUE now. Keep alpha-beta, controlled time management and
independent local comparisons. Confirm proposed features through ablations.

## Decision and next check-in

Publish this as **v0.15.1 research**, with the playing policy unchanged and the
reporting linkage repair covered by focused tests. Do not claim a strength
promotion. Start the next engine candidate with piece safety, stop early if it
cannot explain its target cases, and keep every failed result.

Fifty games per opponent remain useful triage but are inadequate to resolve
most gains estimated above. [TESTING.md](TESTING.md) recommends a separately
planned 200-game confirmation for promising strength candidates and explains
when even that is too small. This is a proposal, not a silently changed default.
The next engine change still needs implementation, validation and a user
check-in before a PR or engine release.
