# v0.10 transposition-table screen

The candidate retains v0.9's evaluator, alpha-beta search, two-ply quiescence
cap, clock policy, and exact root tie behavior. It adds a fixed 65,536-entry
per-process table, about 2 MiB, for completed main-search scores and move hints.
Scores are reused only at matching depth. Quiescence is uncached. The table
persists across search iterations and moves, clears on UCI `ucinewgame`, and
is not saved to disk. Each engine process owns a separate table.

## Method and acceptance criteria

Run 100 games against v0.9 alone. If the existing previous-release gate passes,
run 100 against Stockfish 19 Skill 1 and then 100 against Sacrifice 0.2.0
`b528f0ba`. All games and matchups are serial. The predeclared trigger for an
additional 100-game Skill 2 screen is at least 60% at Skill 1, with no forfeits
or unfinished games. Check all routine matchups for reliability and finish
Sacrifice before starting Skill 2. This trigger permits an additional screen;
it is not a confidence interval or automatic permanent ladder advancement.
Do not lengthen or repeat the run automatically after a favorable result.

The unchanged ten-opening suite uses 50 color pairs, 100 ms requested per move,
50 ms protocol allowance, and a 400-ply cap. Hardware is an Apple M2 MacBook
Air; macOS reported AC power at the start. Only one game runs at a time and
the two engines take turns searching. The exact candidate, opponent, opening,
runner, and Sacrifice PGN asset hashes are in [summary.json](summary.json).
This is a development build, not a released v0.10 binary.

## Reproduce

Build the tagged v0.9 baseline and the pinned Sacrifice opponent using the
[benchmark instructions](../../README.md), then run from the repository root:

```sh
go run ./tools/benchmark --validate --race --version v0.10-dev \
  --previous tools/benchmark/bin/pandafish-v0.9 --previous-version v0.9 \
  --previous-revision 88cc32fec528d763681f108e05f7078924071152 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 88cc32fec528d763681f108e05f7078924071152-dirty \
  --machine 'Apple M2 MacBook Air; all matchups serial; AC power reported at start' \
  --verbose
```

The revision is the base commit with local changes, not a claim that the base
contains this implementation. The saved SHA-256 identifies the actual binary.
Random ties and timed search mean rerunning need not reproduce identical games.
Keep these development labels/hashes when preparing a later release.

## Functional evidence

[Cache checks](cache-checks.json) record the small fixed-depth node-reuse checks.
These show reduced repeated work on the checked position, not a guaranteed
speedup or completed-depth gain everywhere. Tests also compare cached scores
against uncached alpha-beta, which remains checked against exhaustive minimax.
They cover bounds, horizon isolation, mate scores, index collisions, move
encoding/promotions, en passant, castling/turn/half-move context, interruption,
restricted roots, and UCI cache clearing. Root restrictions never enter the
shared score cache. Existing race/protocol tests cover cancellation and worker
replacement. Repetition-history draw evaluation remains future work and will
require revisiting cache validity.

## Results

The previous-release matchup completed with **37 wins, 39 draws, and 24 losses
(56.5%)**, zero forfeits, and zero unfinished games. This passes the practical
screen, but the modest lead in 100 games is not conclusive evidence of a general
strength gain. The cache's correctness and reduction of repeated work are
established separately by focused tests. Do not equate those functional checks
with guaranteed improvement against every opponent or at every time control.

| Opponent | W / D / L | Score | Unfinished / forfeits | Candidate max response |
| --- | --- | --- | --- | --- | --- |
| v0.9 | 37 / 39 / 24 | 56.5% | 0 / 0 | 97.314 ms |
| Stockfish 19 / Skill 1 | 45 / 14 / 41 | 52.0% | 0 / 0 | 98.730 ms |
| Sacrifice 0.2.0 b528f0ba | 0 / 27 / 73 | 13.5% | 0 / 0 | 97.378 ms |

All 300 games completed, with no forfeits by either engine and no unfinished
games. No candidate response exceeded the 100 ms request. These are measured
protocol response times, not a guarantee for future positions or host loads.
The opponent maxima were 97.318 ms (v0.9), 101.967 ms (Stockfish), and
100.321 ms (Sacrifice), all within the 150 ms deadline.

The Skill 1 score is the same observed 52% as v0.9's screening result, giving
no clear evidence of a gain against Stockfish. Sacrifice's 13.5% is similar to
v0.9's 12.5%; the sample does not establish a change against it either. The
candidate won fewer games against Sacrifice but drew more. Preserve the full
W/D/L rather than describing this as a demonstrated strength increase.

**No Skill 2 run was started.** The [recorded decision](skill-2-decision.json)
compares the 52% Skill 1 result with the predeclared 60% threshold. Skill 1
remains the routine target. Do not repeat this candidate's screen to seek a favorable threshold result.

## Saved-position audit and next checkpoint

The [position audit](position-audit.json) ran both engines on all nine saved
public live-game cases, serially after the routine matches. Each engine had
one fresh-process `go depth 3` probe and one separate fresh-process
`go movetime 100` probe per position. All 36 moves were legal, and all nine
fixed-depth scores agreed, as expected from unchanged evaluation and horizon.

At 100 ms, v0.10 completed one extra main-search ply on three positions and
the same depth on six; none completed less depth in this trial. This small,
single-trial cold-table audit demonstrates useful search reuse in these
examples without establishing a general speedup or precise strength gain.
Actual game processes also reuse their tables across moves. No code or weights
were tuned after inspecting the public holdout outputs.

The evidence supports a correct cache implementation and less repeated work,
with a modest positive previous-release screen. It supports a draft PR for
review, not an inflated claim that every opponent comparison improved.
The proposed next checkpoint after review is a recorded Lichess session to
check operational behavior and collect varied regression cases before choosing
further evaluation/search work. Keep the existing alpha-beta foundation.

Records: [summary.json](summary.json), [previous.pgn](previous.pgn),
[stockfish.pgn](stockfish.pgn), [sacrifice.pgn](sacrifice.pgn), and combined
[games.pgn](games.pgn). Binaries and private stream configs remain ignored.
