# v0.13 selective-quiescence experiment

Diagnostic of the development candidate later approved for v0.13. Original
labels and hashes are retained. Main search, material/placement
scoring, transposition rules, root ties, and time allocation are unchanged.
The candidate allows 16 ordinary tactical plies rather than two, retains the
four extra checked-node evasion plies, and filters negative material exchanges
for non-checking ordinary captures outside check. Checking captures,
promotions, and evasions remain eligible. No new UCI option is introduced.

## Exchange heuristic and tests

SEE copies the board into stack-local arrays and follows least-valuable legal
recaptures on the destination square, with the option to stop. It handles
exposed rays, pins, king safety, en passant removal, and promotions. It does not
branch over every equal-valued attacker or evaluate tactical compensation on
other squares. This is a selective-search heuristic, not a proof that discarded
moves are bad. Equal/cheaper attackers bypass the estimate, avoiding work on exchanges that
usually cannot lose material on the target. A recapture promotion can be an
exception; bypassing the filter retains rather than discards such a capture.

Seventeen hand-counted fixtures are checked against an independent material
oracle using every legal recapture, including equal-valued alternatives. The
13-position corpus exercises another 14 potential filters against that oracle.
Cancellation/deadline tests and both-color ordinary-exchange regressions are
included. These finite tests do not establish universal SEE correctness.

The earlier Rxd5/Qxd5/Bxd5 fixture is fixed by filtering the losing queen
recapture even with a two-ply allowance. Therefore that regression supports
selectivity, not a claim that the 16-ply bound alone produced the fix. The
broader unfiltered variant below separates those effects.

The dense exhaustive minimax test initially timed out when its unfiltered
reference was raised to 16 plies. Its oracle now uses four quiescence plies on
those dense positions, with the same capture eligibility as production. Sparse
reference cases still use the production bound. This was reference-tree
explosion, not a timed engine forfeit. The independent material oracle tests SEE
separately; the minimax reference tests search windows, signs, and root ties.

## Selected-position audit

[Positions](positions.json) combine seven additional public game positions from
the local chess-ai study and six previously confirmed stream positions. Source
hashes and links are recorded; study annotations are not best-move ground truth.
Three stream positions retain their previous holdout labels. This selected
corpus is diagnostic, not a representative strength sample or a tuning set.

Each variant receives a cold-process depth-three probe and a separate 100-ms
probe. Requests are serial. The referee is Stockfish 19, full strength, one
thread, 64 MiB, 100,000 nodes per unrestricted/restricted query, cleared hash
and new game each time. The saved analysis code associates the score with the
returned move/PV, rechecking mismatches. Mate scores remain separate. Signed
centipawn differences are finite-budget estimates, not exact regret.

All 78 probes returned legal moves with no watchdog errors or fallback moves.

| Variant | Depth-three nodes | Depth-three elapsed ms | 100-ms depth distribution |
| --- | --- | --- | --- |
| v0.12 | 34,167 | 135.274 | 3: 4, 4: 6, 5: 2, 7: 1 |
| Broader, unfiltered | 44,320 | 204.737 | 3: 6, 4: 5, 5: 1, 7: 1 |
| Broader, SEE filtered | 32,055 | 147.137 | 3: 3, 4: 8, 5: 1, 7: 1 |

At depth three, filtering reduces node visits 27.7% versus broader unfiltered
continuation, and elapsed time about 28.1%. Against v0.12 it visits 6.2% fewer
nodes but takes about 8.8% longer: SEE work is not counted as search nodes.
These are single serial measurements with scheduling noise, not stable NPS.

For the nine positions with centipawn scores on both unrestricted and played
queries, summed signed estimated loss at depth three is 2,179 cp for v0.12,
1,449 unfiltered, and 1,418 filtered. At 100 ms the sums are 1,484, 1,285, and
1,428 cp. These sums exclude four mate-scored positions and should not be
interpreted as Elo or a percentage strength gain.

The filtered candidate's 100-ms choices improve on v0.12 in three scored cases,
match three, and worsen three. It completes a deeper main-search iteration in
two cases, the same depth in nine, and a shallower iteration in two. This is
mixed evidence: filtering controls the broader tree, but does not establish a
large fixed-time gain. Match results must remain separate from this diagnostic.
No parameters were retuned from these results.

## Reproduction and records

- [Raw results with UCI output and binary hashes](results.json)
- [Compact aggregates](summary.json)
- [Candidate source snapshot](candidate-source.json), relative to v0.12
- [Unfiltered ablation patch](unfiltered.patch), relative to the candidate
- Audit sources: [driver](main.go.txt), [engine probes](engines.go.txt),
  [Stockfish client](stockfish.go.txt)

The `.go.txt` suffix keeps standalone diagnostic drivers out of `go test ./...`.
To reproduce, restore the candidate files over the recorded base revision in
a temporary checkout, build the baseline and candidate at the paths in the
driver, apply the zero-context patch with `git apply --unidiff-zero`, then build the
patched unfiltered source at
`/tmp/pandafish-v013-unfiltered`. Copy the three audit sources into a temporary
directory with `.go` extensions. From the checkout root, run:

```sh
go run /tmp/pandafish-v013-audit/engines.go \
  /tmp/pandafish-v013-audit/stockfish.go \
  /tmp/pandafish-v013-audit/main.go /tmp/pandafish-v013-audit.json
```

Adjust engine/Stockfish paths for the local installation. The driver reads
`testdata/quiescence-validation.json`, records the input and referee hashes,
uses a 20-second probe watchdog, and checkpoints after each position. Random
root ties and scheduling can change moves and timing. No online games are run.

## Match decision

Use the routine serial 100-game previous-release gate, then Skill 1 and pinned
Sacrifice if it passes. Before reviewing match scores, a visible external improvement was
described as at least 10 percentage points on one reference opponent versus
the v0.12 screen, without a drop greater than 5 points on the other. This is a
descriptive target, not significance or another automatic gate. Preserve flat
or negative outcomes. No automatic retries or longer runs.


The [completed routine game screen](../../../results/v0.13-dev-vs-v0.12-and-stockfish-1-and-sacrifice-v0.2.0-b528f0ba/README.md)
passed the previous-release gate and improved both external scores. It records
the full 300 games and interpretation separately from this selected diagnostic.
