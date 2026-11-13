# v0.15 conversion correction: mate distance and repetition

Promoted in **v0.15** for demonstrated conversion corrections after review.
The original development labels, hashes and control records remain unchanged.
External matches did not establish a broad strength leap. The next work is the
[extended live investigation](../../live/v0.15-test-plan/README.md).

## Reproduced failure before implementation

Released v0.14.4 gives all forced mates the same score. At depth three in each
of two simple queen endings, its best list includes the unique mate in one and
six moves that delay mate. Random tie breaking therefore chooses the immediate
mate only one time in seven in those specific lists. This is an actual decision
failure, not an estimated Elo benefit or a claim that every delayed move draws.

The exact original output is in `baseline-choice.txt`. The positions are:

- `7k/5Q2/5K2/8/8/8/8/8 w - - 0 1`, immediate mate `f7g7`.
- `7k/8/5KQ1/8/8/8/8/8 w - - 0 1`, immediate mate `g6g7`.

These synthetic fixtures were chosen to isolate conversion. They are not saved
live-game examples and do not prove the cause of earlier match losses. The prior
king/queen decision audit did not establish a successful positional correction.
This is a stronger implementation target than another unverified weight guess.

## Sacrifice inspiration

[Sacrifice's pinned search](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs)
returns a losing mate score adjusted by ply in main search and quiescence. Its
`score_to_tt` / `score_from_tt` functions convert mate scores between root and
position coordinates. These are the specific concepts adopted here in original
Go. No Rust routines, tables, pruning parameters or external source is copied.

Sacrifice also has PVS, reductions, history and richer evaluation; those are not
part of this correction. Earlier failed experiments remain archived, not active.

## Initial mate-only scope and score contract

A mated node scores `-mateScore + ply` from the side to move. Negamax then prefers
shorter winning mates and longer forced losses. Root children start at ply one;
quiescence carries the same root ply instead of resetting the distance. Static
material/placement and terminal `evaluate` API values stay unchanged.

The transposition table stores distance relative to the cached position. Main
search converts the score on store/probe, including window-bound classification.
Ordinary centipawn values are unchanged. The mate band is 1,024 below the mate
constant, above all legal material/placement scores and comfortably wider than
the 64-ply main search plus bounded quiescence/evasions. No table size/layout,
identity, replacement, or equal-depth reuse policy changes.

Telemetry and UCI `score mate` now report a calculated signed move count, rounded
from plies. Earlier historical captures used mate signs only; those records keep
their original interpretation. This is a reporting-contract extension, not a
reinterpretation of old data. PV still contains just the selected root move.

The initial mate-only stage added no new positional term, material value, search
reduction/extension, clock policy, draw handling, book, dependency or UCI option.
Its repetition limitation motivated the pre-match scope amendment below. Equal scores at equal mate distance retain random tie breaking.

## Validation and decision before match scores

Require both colors to choose the unique immediate mate at depths one through
four with cold/warm tables. Independently exhaust the legal mate tree of a
losing fixture and compare longest resistance and exact distance. Cover positive
and negative cache coordinates/bound kinds across different root plies,
quiescence offsets, signed UCI move counts, cancellation, and root immutability.
Existing non-mate tests remain in place; update only the reporting assertion
that previously prohibited uncomputed mate distances.

Run normal/race tests, vet, ARM64/offline wrapper, and the standard 50 games each
against tagged v0.14.4, Stockfish Skill 2 and pinned Sacrifice, all three matchups
concurrent. Preserve all records even if the previous-release criterion fails.
The known mate-choice correction is a functional improvement; its overall Elo
and contribution to game results are unknown. Do not promise the large leap
sought or infer one from a handful of mate fixtures. Read complete scores and
reliability before recommending promotion. No retries or automatic longer match,
no candidate commit/PR until the completed assessment is reviewed.

## Evidence-driven amendment before routine matches

The initial conversion diagnostic used four preselected queen-versus-king
positions (two mirrored pairs), cold depth-five converters and Stockfish 19 at
20,000 nodes defending. It showed mate-only conversion in two of four cases,
while both policies repeated in the other two. Baseline drew all four. This was
an exploratory FEN-only diagnostic, not the routine strength screen. Its complete
record and mate-only source are retained separately; no outcomes are overwritten.

Add threefold recognition as the second half of this same conversion correction:
count the reversible UCI game history and current main/quiescence search path.
At a third occurrence, return zero; checkmate/stalemate terminal status is checked
first. Push/pop state unwinds after cancellation and cutoffs. Root positions are
counted once per iteration, not again when main search enters quiescence.
A bare FEN cannot supply history from before that FEN.

The UCI loop captures immutable board hashes before starting the worker. History
before the last pawn move or capture cannot recur and is dropped. Board identity
uses the chess library's Zobrist convention for pieces, turn, castling and
capturable en passant. It follows the same library identity convention used by
the local referee, rather than implementing a separate draw-law engine.

Cached scores now include a count-set history fingerprint in their existing key.
The same board with different prior occurrence counts must not share score bounds.
The fingerprint is updated on push/pop and combined with the board key; halfmove
context remains checked separately. Entry layout/size remain unchanged. This can
reduce transposition reuse and costs map work, so actual match play must judge the
tradeoff. As with the existing board hash, collisions are probabilistic.

A second, explicit full-history harness gave all three policies identical
initial FEN plus move histories: release, preserved mate-only, and combined.
Stockfish's defender queries retained the same FEN-only 20,000-node policy. This
one formal diagnostic converted 1/4, 2/4, and 4/4 respectively. Candidate mates
arrived in 15/7/15/7 plies. Release randomness changed one earlier result; both
runs stay in the record. This selected four-case result is a demonstrated
conversion improvement, not representative Elo or a guarantee against all
endgame positions. No routine match has been rerun for a favorable score.

The additional tests require a material-disadvantaged side to choose a drawing
repetition from actual UCI history, reject cached scores from a different history,
restore count state on draw/cancellation, trim irreversible history, and preserve
mate precedence. The final candidate label is `v0.15-conversion-dev`; the earlier
mate-only label/hash remains in its own control records. No new positional weights
or held experimental pruning are active. All scope changes occurred before the
routine 150-game screen. No further parameter or policy changes during play.

The [result directory](../../results/v0.15-conversion-dev-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
preserves exact final and mate-only source overlays, both diagnostic stages,
Go driver sources, all routine games and compact summaries.

## Completed assessment

The [full record](../../results/v0.15-conversion-dev-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
contains all 150 routine games, all 20 staged conversion trials, exact source and
controls. The routine screen scored 53% against v0.14.4, 50% against Skill 2 and
34% against Sacrifice, with no forfeits/unfinished games. The direct criterion
passed narrowly. Recommend taking the demonstrated conversion correction to
review; do not claim a broad strength leap, precise Elo, or proven absence of
regression. The correction was approved for the v0.15 release; the original
development records do not claim to test the version-renamed release binary.

## Release verification

[Release checks](release-checks.json) record separate version-only validation,
six concurrent smoke games, the ARM64 build and lichess-bot integration checks.
The playing implementation matches the tested development candidate; only its
UCI identity/handshake expectation and release/configuration text changed.
Offline preparation checks the configured recorder without launching the client.
The active bot config and private snapshot point at the released root executable;
they remain ignored. No live games are started as release verification.
