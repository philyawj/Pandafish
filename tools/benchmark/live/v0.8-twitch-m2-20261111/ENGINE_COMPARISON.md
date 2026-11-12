# Engine comparison and next increments

Reviewed after the first v0.8 stream. This follows the [session review](REVIEW.md).
The recommendation is to keep conventional alpha-beta search and the validated
time-control structure, add positional evaluation, then improve search reuse.
Sunfish is useful evaluation inspiration within that direction.

## Where the engines differ

| Engine | Evaluation | Search approach | Useful lesson for Pandafish |
| --- | --- | --- | --- |
| Pandafish v0.8 | Material, checkmate, stalemate; random choice among equal best scores | Iterative-deepening negamax alpha-beta, capture ordering, two-ply quiescence cap; interruptible timed search | The operational foundation works in the observed tests. Quiet positions still need better scoring. |
| [Sacrifice](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/evaluate.rs) | Centipawn material plus positional formulas, mobility, activity, king shelter, center and attack terms | [PVS alpha-beta](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs), iterative deepening, transposition table, killer/history ordering, aspiration windows, late-move reductions, quiescence | A practical example of several classical improvements working together. Separate those changes when testing our engine. |
| [Classic Sunfish](https://github.com/thomasahle/sunfish/blob/e396948acbd80cdcbad3dccfb07e878fba9df8f4/sunfish.py) | Material plus piece-square tables, incremental score updates, separate king tables for different material conditions | Iterative-deepening MTD-bi using repeated narrow bound searches, cached scores/moves, selective quiescence and reductions | Small positional scoring can give a compact engine useful preferences. Its current implementation is more involved than the famous line-count description. |
| [Stockfish](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp) | [NNUE evaluation](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/evaluate.cpp) with further scaling | Highly developed alpha-beta/PV search, iterative deepening, transposition caching, move ordering, reductions and pruning | Follow the separation of evaluation, search, timing, and measured testing. Neural evaluation and the full pruning machinery are later projects. |

Sacrifice combines more features than plain alpha-beta, which Pandafish already
has. PVS probes later moves with narrow windows and searches again when needed;
caching reuses work, ordering accelerates cutoffs, and reductions trade search
coverage for speed. Each deserves a separate experiment.

Its [Tal-PGN style data](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/pgn.rs)
influences root move ordering and is excluded from static evaluation. This is
not neural training. A game database or language change is unnecessary for
Pandafish's next improvement.

## Interpreting the stream observations

Early king walks, queen wandering before minor-piece development, and odd pawn
moves are consistent with a material-only evaluator. Many quiet alternatives
receive the same score until something becomes capturable or mate becomes
visible. Random tie breaking then has little positional guidance. This is an
explanation to test, not proof of a single cause for every bad move.

The [live review](REVIEW.md) confirms poor king, rook, and pawn decisions at a
higher analysis budget. Queen wandering remains an observer note rather than
a measured causal category.

Piece-square tables score placement, not repeated moves. They will not by
themselves prevent every queen excursion. Sacrifice's bonuses for pawn
advancement and queen activity do not directly penalize premature development;
copying those weights could reward behavior we want to investigate. Endgame
kings need activity, so a permanent king-centralization penalty would be
inappropriate.

## Recommended next steps

1. **First increment: positional evaluation.** Preserve search, clocks, the
   legal-move library, and random choice among true ties. Add modest piece-square
   terms for all pieces, with a small material-phase rule separating king safety
   from endgame activity. Encourage development and sensible central pawn
   placement without banning tactical king/queen moves. Internal scores are
   currently in pawn units and multiplied by 100 for UCI: migrate units,
   terminal-score dominance, quiescence stand-pat, ordering, and reporting
   consistently. Test color symmetry, phase transitions, and quiet-leaf scoring.
   Write our own small implementation; check licenses before copying tables.
2. **Second increment: a small transposition table.** Store position key, depth,
   bound type, score, and best move; try the hash move first. Test interruption,
   mate-score handling, collisions, and draw/history context. Compare completed
   depth and response tails under unchanged clock settings. This develops the
   existing alpha-beta architecture.
3. **Then follow evidence.** Persistent early queen/king problems could justify
   a separate phase-aware development or shelter term. Search bottlenecks could
   justify quiet-move ordering, then PVS. Defer aspiration windows, reductions,
   and null-move pruning until the simpler baseline is established. Extend
   quiescence or draw/conversion logic against reproducible cases.

The best fit is a small conventional Go engine with readable classical
scoring, using Stockfish's architectural discipline and Sunfish's compactness
as inspiration. This does not require adopting Sunfish's MTD-bi or reproducing
Stockfish's NNUE system. Sacrifice provides an intermediate classical example,
with more features than we need to add in one increment.

## Strength claims and acceptance

The friend's Level 5 result remains anecdotal here. Verify the difficulty
scale, Stockfish build, clocks, openings, colors, and sample size: Lichess
computer difficulty and native UCI Skill Level are different controls. The
[README example](https://lichess.org/NW5ZGVEZ) is a 5+3 human/bot draw; no comparable
Level 5 match series was verified. No Sacrifice or Sunfish matches were run.
Our roughly 10% figure comes from the [v0.7-dev benchmark](../../README.md),
whose search/evaluation policies v0.8 retained.

Test the first increment against the saved development/hold-out positions and
run 100 games against v0.8 plus 100 against Stockfish Skill Level 1. Later
increments use the latest released baseline. Keep the documented skill-ladder
promotion rule. A subsequent live session tests operational behavior and finds
new cases; the successful M2 stream does not prove every future clock scenario.

Sources are pinned above: Sacrifice `b528f0ba`, Sunfish `e396948a`, and Stockfish
`def601cc`. This is source inspection and planning; no engine code changed.
