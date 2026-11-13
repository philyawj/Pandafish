# v0.15 conversion correction

## Decision: approved for v0.15

This candidate has a demonstrated functional benefit: prefer immediate forced
mates, avoid known winning-side repetitions, and choose an available drawing
repetition when behind. It narrowly passed the direct previous-release screen.
It does **not** demonstrate the broad strength leap sought. Recommend proceeding
to review for those conversion corrections, with the external uncertainty stated
explicitly. The correction was approved for v0.15. These records retain the
tested development identity/hash; they are not a test of the version-renamed
release binary.

| Opponent | W / D / L | Score | Unfinished | Forfeits |
| --- | --- | --- | --- | --- | --- |
| v0.14.4 | 20 / 13 / 17 | 53% | 0 | 0 |
| Stockfish Skill 2 | 24 / 2 / 24 | 50% | 0 | 0 |
| Sacrifice v0.2.0 b528f0ba | 10 / 14 / 26 | 34% | 0 | 0 |

All three 50-game matchups ran concurrently. Draws count as half a win. The 53%
previous-version result passes the greater-than-50% criterion, but the small
margin is not convincing Elo evidence. All 150 games finished without forfeits.
Historical released-policy external scores in the
[concurrent calibration](../v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
were 60% against Skill 2 and 34% against Sacrifice. The new 50%/34% screen is down
ten points/flat, so do not present external strength as improved or regression
as ruled out. These are separate 50-game screens, not simultaneous paired
controls; a precise strength difference cannot be inferred from them.

No routine match was rerun, enlarged or selected for a favorable score. The
recommendation is based on demonstrated conversion behavior plus a narrow gate
pass, not proof that the engine is generally stronger or closer to Stockfish 20.
The Stockfish target stays Skill 2. User review is required before publication.

## Actual failure and Sacrifice inspiration

Before changing code, depth-three searches in two synthetic queen endings gave
seven equally best moves: the unique immediate mate and six moves that delayed
mate. The original root lists are retained in the
[baseline evidence](../../research/v0.15-mate-conversion/baseline-choice.txt).
This meant a one-in-seven immediate-mate choice in those specific lists, not a
measured conversion rate over all chess positions.

[Sacrifice's pinned search](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs)
uses ply-adjusted terminal mates and root/position score normalization in its
transposition table. Those are the specific ideas adopted in original Go. No
Rust code, weights, PVS, reductions or rich positional terms were imported.
The [pre-match plan and amendment](../../research/v0.15-mate-conversion/README.md)
document the scope, tests, initial evidence, and why repetition was added before
routine match play.

The first mate-only version fixed the immediate-mate choices, but repeated in
two selected queen endings. That supplied a second observed conversion failure:
recognize threefold history so a known repetition scores zero, below a winning
continuation. This is an evidence-driven conversion increment, not another
collection of guessed positional weights. It does not establish the causes of
all earlier live-game or match losses.

## Selected endgame diagnostics, including unfavorable controls

Two queen-versus-king positions and their color/rank mirrors were fixed before
outcomes. Converters searched cold depth five per move; Stockfish 19 defended at
20,000 nodes, Skill 20, one thread, 64 MiB hash, fresh hash for each FEN query.
The game referee claimed eligible threefold/50-move draws, with a 100-ply cap.
A full-strength defender in a lost queen ending is not a full-game strength test
against Stockfish 20.

The exploratory FEN-only stage converted 0/4 released and 2/4 mate-only. It showed
that repetitions remained, but supplied no prior game history to either engine.
Its complete original record remains `mate-only-conversion.json`.

The explicit full-history stage then gave every converter the same initial FEN
plus played move history; defender queries retained the same FEN-only policy.
All three policies were run once per fixture, serially, before routine games:

| Policy | Wins / draws / losses | Mating plies when converted |
| --- | --- | --- |
| Released v0.14.4 | 1 / 3 / 0 | 29 |
| Mate-distance-only control | 2 / 2 / 0 | 7, 7 |
| Mate distance plus repetition | 4 / 0 / 0 | 15, 7, 15, 7 |

The released result changed between stages due to random ties; preserve both
rather than choosing the worse baseline. These are four selected examples with
one trial per policy, not a statistically representative endgame suite. They
establish the intended behavior in those examples and supply regression targets,
not an overall Elo estimate. No claimed large gain is attached to them.

`full-history-conversion.json` retains every FEN, move, score/output, identity and
binary hash across all 12 formal trials. The exploratory stage retains all eight
trials separately. The Go drivers are archived as `.go.txt`. Build the preserved
mate-only control into `/tmp/pandafish-mate-only` before replaying that driver.
The formal driver sends full history to all converters, including the release.

## Implementation and tests

`mate.go` represents a losing mate as `-mateScore + rootPly`. Negamax prefers
shorter wins and longer resistance. Root children start at ply one, and main and
quiescence search share the same root-ply coordinate. The transposition table
stores mate distance relative to the cached position, then restores it at the
caller's ply; bound classification uses converted score/window coordinates.
Ordinary centipawn material/placement values stay unchanged. UCI/telemetry now
report calculated signed moves to mate; old sign-only captures retain their
historical interpretation.

`repetition.go` captures the known reversible UCI history as immutable hashes
before launching the search goroutine. It tracks occurrence counts through main
and quiescence paths with balanced push/pop, and returns zero at a third
occurrence after terminal-status precedence. Root positions are counted once,
not again at the main/quiescence transition. A standalone FEN cannot provide
history before that FEN. Repetition identity uses the existing chess library's
Zobrist convention, also used by the local referee.

Draw-sensitive cache bounds must not cross histories. A count-set fingerprint
is included in the existing key, with the existing halfmove check retained.
This preserves the fixed table entry layout/size while distinguishing prior
occurrence counts; it can reduce cache reuse. Hash collisions remain probabilistic,
as with the existing board key. Game history before the last pawn move/capture
cannot recur and is dropped. This does not add fifty-move evaluation, a complete
draw-rule engine, or any new positional bonus.

Tests require both colors to take unique immediate mates at depths one through
four with cold/warm tables, preserve score coordinates/bound kinds across root
plies, carry quiescence offsets, emit signed UCI mate counts, and match an
independent exhaustive mate-only proof of longest resistance. A material-losing
side must find a drawing repetition from actual known history. Further tests
cover cache context separation, count-set order independence, cancellation and
path restoration, irreversible-history trimming and mate precedence. All prior
non-mate tests remain; the old reporting assertion was updated because a mate
distance is now computed rather than invented.

Normal tests, race tests, vet, ARM64 build, and the offline lichess-bot wrapper
passed. The wrapper checked 12 legal plies, then parsed `score mate 1` for the
immediate-mate fixture and zero for the full-history drawing move. Both private
configs were unchanged; the root executable was restored byte-for-byte to its
original v0.14.1 identity after checking. The public placeholder example loaded
with local resolved paths. No live listener or online game was used.

## Routine reproduction and retained records

Base: `3f3ca331cbe8aac9c7c2194b9314b33b3da34a81` (v0.14.4).
Candidate label: `v0.15-conversion-dev`.
Candidate SHA-256: `ab78d33434e380eb7d093c21e61772aef68fd82de17ff2b30e121509a2a20537`.
The formal combined-policy diagnostic, matches and wrapper used this same hash.
The earlier mate-only source/hash remain separate.

```sh
caffeinate -i go run ./tools/benchmark --validate --race \
  --version v0.15-conversion-dev \
  --previous tools/benchmark/bin/pandafish-v0.14.4 --previous-version v0.14.4 \
  --previous-revision 3f3ca331cbe8aac9c7c2194b9314b33b3da34a81 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 3f3ca331cbe8aac9c7c2194b9314b33b3da34a81-dirty \
  --parallel-all --machine 'Apple M2 MacBook Air; three concurrent matchups; no live bot' \
  --verbose
```

Defaults were 25 opening pairs, 100 ms per move, 50 ms timeout grace and 400-ply
cap. Stockfish used Skill 2, one thread, 16 MiB hash, Elo limiting off, no Syzygy.
Pinned Sacrifice remains v0.2.0 b528f0ba. Exact options/hashes, opening hash, move
records and response times are in `summary.json`; the compact view is
`screen-summary.json`. Per-opponent and combined PGNs contain all 150 games.

The match window was approximately 582 seconds on the M2 MacBook Air, with AC
power observed early and after play. Continuous thermal data were not collected.
All 8,325 candidate responses stayed below 100 ms; maximum was 99.356 ms and
matchup p95 about 96.5-96.6 ms. These local measurements do not include live-client
submission or network latency. They validate this run's operation, not every
future timing condition.

`candidate-source.json` is the exact eleven-file Go overlay onto tagged v0.14.4.
`mate-only-source.json` preserves the eight-file intermediate control. Restore
either onto separate temporary tag checkouts with the shared
[Go restore helper](../../research/v0.15-selective-search/restore-snapshot.go.txt),
then test/build them separately. The final source includes all focused tests.
Replay restores the policy, not guaranteed compiler hashes or random outcomes.
`checks.json`, power observations and a SHA-256 manifest preserve verification
and provenance alongside the raw diagnostics and matches.

## Assessment and next direction

Unlike the earlier unsuccessful positional/search hypotheses, this has concrete
before/after decision evidence. It is worth reviewing as a conversion correction.
Its full-game evidence is narrow and externally mixed; do not market it as the
big leap requested or rule out a strength regression from these samples.

Keep the positional backlog as hypotheses. Next work should link a confirmed
saved tactical/quiet-position failure to a focused test before adding weights.
Sacrifice's piece threats, piece-square scores and material distinctions remain
useful references, but their presence in that engine is not an Elo estimate for
Pandafish. A broader held-out conversion suite would strengthen this specific
claim; it is not being run automatically to improve a match score. The next approved work is the
[extended v0.15 live investigation](../../live/v0.15-test-plan/README.md).
