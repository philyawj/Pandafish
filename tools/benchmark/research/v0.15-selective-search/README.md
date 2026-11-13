# Selective-search experiment: what we tried and why it was not promoted

This report preserves an unsuccessful strength hypothesis, including its code,
tests, measurements, interpretation, and proposed next step. It is published in
**v0.14.1 as a research archive**. The release keeps v0.14's playing algorithm.
The experimental PVS/LMR code is archived data, not part of the engine build.

The experiment identified itself as **v0.15-dev**. That label remains in its
original records; it was never promoted into a v0.15 release. A future candidate
should use a distinct label, such as `v0.15-evaluation-dev`.

## Decision at a glance

**Do not promote this selective-search candidate.** It reduced search work and
sometimes completed another iteration, but did not demonstrate the substantial
strength gain we wanted. Keep the evidence rather than turn a narrow pass into
a strength claim or run again until the score looks better.

| Opponent | Candidate W/D/L | Score | Same-profile external reference |
| --- | --- | --- | --- |
| v0.14 | 16/19/15 | 51% | Direct previous-version comparison |
| Stockfish Skill 2 | 22/8/20 | 52% | 60% |
| Sacrifice v0.2.0 b528f0ba | 8/18/24 | 34% | 34% |

All 150 games completed with zero forfeits and zero unfinished games. The
previous-release criterion technically passed, but only by one percentage
point. The external screen was down eight points at Skill 2 and flat against
Sacrifice. Fifty-game screens do not prove regression, equivalence, or a precise
Elo difference; they also give no convincing evidence of the larger gain sought.

[Full method, results, raw PGNs, timings, and checks](../../results/v0.15-dev-vs-v0.14-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
remain separate from this interpretation. No records were renamed to imply that
they tested the research release.

## Why this was the next hypothesis

The [post-v0.12 review](../v0.13-candidates/README.md) identified three larger
workstreams: stabilize ordinary tactical exchanges, add selective main search,
or improve classical evaluation. The first shipped in v0.13. v0.14 then made
testing faster without changing playing strength. Older live-game priorities
were written before the intervening check-evasion, ordering, and exchange work.
They were historical groundwork rather than a reason to keep making small
ordering changes indefinitely.

The next hypothesis was:

> Use narrow-window search and guarded reductions to spend less time on
> low-priority branches, finish useful deeper iterations within the same
> clock allowance, and improve game results.

This was chosen for search-depth upside. Another small ordering hint could
save a constant amount of work without changing completed depth. Reductions
can change the effective size of the tree more substantially. The improved
quiescence search provided a stronger tactical leaf policy on which to build.
That made selective main search a reasonable experiment, not a guaranteed win.

[Stockfish's pinned search](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp)
and [Sacrifice's pinned search](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs)
provide examples of narrow-window search, reductions, and verification in much
richer systems. Their feature combinations do not identify the Elo contribution
of one mechanism. The code here was original Go, with no copied tables or
attempt to reproduce every feature of either engine.

PVS was supporting machinery. The main strength hypothesis was the combination
with late-move reductions, rather than releasing PVS alone as another small
optimization.

## What changed in the experiment

Evaluation, quiescence, time allocation, transposition-table layout, and UCI
options stayed fixed. The root retained its comparison margin and random
selection among computed equal scores. There were no aspiration windows,
null-move searches, history tables, new dependencies, or new evaluation terms.

### Principal variation search

The released main search gave each move its current full alpha-beta window:

```go
score, complete := worker.searchAtPly(
    child, depth-1, -beta, -alpha, givesCheck, ply+1,
)
score = -score
```

PVS searched the first move normally. Later moves first asked a cheaper question:
can this move improve the current bound? A one-centipawn window answers that
without necessarily finding an exact value for every inferior alternative.
An improvement inside the original window then received a wider re-search.

This shortened excerpt shows the PVS-only mechanism; cancellation, caching,
and ordering are omitted here, but are present in the archived implementation:

```go
score, complete := worker.searchAtPly(
    child, depth-1, -alpha-1, -alpha, givesCheck, ply+1,
)
score = -score
if complete && score > alpha && score < beta {
    score, complete = worker.searchAtPly(
        child, depth-1, -beta, -alpha, givesCheck, ply+1,
    )
    score = -score
}
```

With reductions disabled, PVS should preserve the full-width result for the
same underlying leaf policy. Tests checked scores, fail-low/fail-high bounds,
and complete best-root-move sets instead of demanding one random tie choice.

### Guarded late-move reductions

For later PVS probes, a fixed eligibility rule allowed exactly one less ply:

- Remaining depth at least three, fourth move or later.
- Ordinary quiet move: no capture, en passant, or promotion.
- The side to move is not checked and the move does not give check.
- No matching transposition-table or killer-move hint.
- The current alpha is not mate-valued.
- Root moves are not reduced directly.

The first probe used `depth-2` rather than the normal child depth `depth-1`.
If it improved alpha, it had to repeat at full nominal depth before its score
could produce a cutoff. A full-window re-search followed an interior improvement.

```go
score, complete := worker.searchAtPly(
    child, depth-1-reduction, -alpha-1, -alpha, givesCheck, ply+1,
)
score = -score
if complete && reduction != 0 && score > alpha {
    score, complete = worker.searchAtPly(
        child, depth-1, -alpha-1, -alpha, givesCheck, ply+1,
    )
    score = -score
}
if complete && score > alpha && score < beta {
    score, complete = worker.searchAtPly(
        child, depth-1, -beta, -alpha, givesCheck, ply+1,
    )
    score = -score
}
```

These excerpts explain the mechanism; [the full experimental search](experimental-search.go.txt)
is the exact readable source. [The test source](experimental-selective_test.go.txt)
and [complete source snapshot](../../results/v0.15-dev-vs-v0.14-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/candidate-source.json)
include all surrounding details.

### What verification does and does not guarantee

Verification protects against a shallow probe that looks too optimistic.
It does not detect every quiet move that looks unpromising at reduced depth
but becomes good at full depth. Such a reduced fail-low is where the search
chooses to save work. Descendants of a full-nominal-depth re-search can still
use reductions themselves.

Consequently, LMR is a selective policy, not a guarantee of full-width minimax
parity. Ordering and cached hints can affect which lines are reduced and the
scores eventually returned. Keeping root comparison margins does not restore
an unchanged full-width best-move set. This distinction was documented rather
than hidden behind the phrase "exact root ties."

## What we checked before games

Three internal modes used fresh workers and tables: released-style full width,
PVS only, and PVS plus LMR. These were diagnostic controls, not UCI settings.
Production code was not switched between policies using the same score cache.

Focused checks covered:

- PVS parity with full-width scores, windows, and root tie sets.
- Protection of early moves, captures, en passant, promotions, checks, evasions,
  mate-valued alpha, and cached/killer hints.
- Real reduced fail-highs that required full-depth verification, including
  shallow optimism rejected by that verification.
- Cancellation during verification and wider re-search, with no incomplete
  parent cached and no unfinished iteration replacing a completed result.
- Root depth protection, existing UCI/fallback behavior, and tactical regressions.

Normal Go tests, race-enabled tests, vet, formatting, build, and whitespace
checks passed. The candidate was ARM64 and returned 12 legal plies through the
offline lichess-bot wrapper. Public config loading passed with resolved temporary
paths and a placeholder token. The active v0.14 executable was restored exactly;
private configs were unchanged. No live listener or games were started.

Two test adjustments are part of the record. The old killer ordering audit
assumed ordering-independent exact results; that is an unreduced-search
property, so it now isolated PVS-only ordering during the experiment. Separate
tests covered selective behavior instead. Also, the starting-position root
fixture exercised reductions without verification, while a direct full-window
search exercised five verified/refuted shallow fail-highs. Each fixture was
used for the behavior it actually exercised. No policy parameter was tuned to
improve the diagnostic or match score. These experimental test changes are
archived; the research release restores the released test suite.

## Diagnostic results

### Four cold depth-four root searches

| Position | Full-width nodes | PVS-only nodes | PVS plus LMR nodes |
| --- | --- | --- | --- |
| Starting position | 2,279 | 2,369 | 1,560 |
| Dense middlegame | 110,556 | 84,612 | 83,458 |
| Queen raid | 56,486 | 60,320 | 43,495 |
| Crowded middlegame | 53,735 | 59,634 | 20,618 |
| Total | 223,056 | 206,935 | 149,131 |

PVS alone preserved exact scores and complete root tie sets. The selected
policy visited about 33.1% fewer nodes than full width and 27.9% fewer than
PVS-only in this small set. The crowded-position score changed from 187 to
192 cp with LMR; the other three scores stayed the same. Changed selective
scores are retained, not treated as equivalent exact evaluations.

PVS alone was not uniformly cheaper: it used more nodes in three of the four
positions, despite a lower total. Aggregate savings are not a universal speed
promise. Nor does a smaller node count include every CPU cost.

### Thirteen varied positions at fixed depth and 100 ms

Each of three binaries received an independent cold-process depth-four probe
and a separate 100-ms request, serially. All 78 probes returned legal moves
without errors. All 13 PVS-only fixed-depth scores matched v0.14.

With LMR, five timed probes completed a deeper iteration, seven the same, and
one a shallower iteration. Two deeper cases and the shallower case were already
mate-scored by the referee. Only one of the 13 chosen moves differed from the
baseline. This showed some usable depth gains, not a broad tactical improvement.

The full-strength Stockfish 19 referee used 100,000 nodes per independent
unrestricted/restricted query, one thread, 64 MiB, cleared hash/new game, and
score/PV association checks. Nine centipawn-scored cases had summed signed
estimated loss of 1,560 cp for v0.14, 1,427 for PVS-only, and 1,428 for PVS plus
LMR at 100 ms. Four mate cases were excluded. These are selected finite-budget
estimates, affected by random ties; they are not Elo or an isolated proof of
LMR's contribution. No weights or thresholds were retuned from them.

## The game screen and why the candidate was not promoted

Run 50 games against each opponent with all three matchups concurrent. Use the
same ten-opening suite, color pairs, 100-ms request, 50-ms protocol allowance,
400-ply cap, Stockfish Skill 2, and pinned Sacrifice as the newer baseline
profile. Hardware was an Apple M2 MacBook Air. The match portion took about
543 seconds. Raw settings, source, hashes, and observations are preserved.

The earlier same-profile external reference used v0.13's playing algorithm,
which v0.14 retained. Its scores were 60% against Skill 2 and 34% against
Sacrifice. It is a separate 50-game screen, not a paired statistical control.

Before scores were reviewed, the descriptive visible-gain target was at least
+10 percentage points on one external opponent without a decline greater than
five points on the other. This was a way to distinguish the desired larger
improvement from another small engineering step, not a significance test.
The ordinary previous-release criterion also remained in effect.

The final 51% / 52% / 34% results did not meet that goal. The previous-release
criterion narrowly passed, but that alone did not justify promotion. Sacrifice
was flat, and Skill 2 was lower. There was no automatic longer screen, threshold
sweep, or favorable-result retry.

Ten of 8,187 candidate responses exceeded the requested 100 ms; the maximum
was 134.232 ms. No response from either side exceeded the 150-ms protocol
deadline. Reliability passed, so protocol failures do not explain this result.
That observation does not prove an exact cause for the strength outcome.

### What we learned

1. The search-work hypothesis was partly supported: fewer visits and some deeper
   completed iterations were measured.
2. The strength hypothesis was not supported convincingly by this screen. More
   nominal depth did not automatically produce better moves or game outcomes.
3. PVS-only parity was a useful control. It separated a correctness-preserving
   probing mechanism from a policy that deliberately changes coverage.
4. Reduced fail-low mistakes remain possible even with fail-high verification.
   The experiment did not identify one as the confirmed cause of the losses.
5. Evaluation quality could limit the value of extra depth, but this run did
   not isolate that as the cause either. That is the next hypothesis to test.
6. A fast routine screen can reject a weak promotion claim without pretending
   to establish a precise Elo difference. Keeping negative records is part of
   the workflow, not a reason to change the scoring rule afterward.

This result does **not** show that PVS/LMR is generally ineffective, that all
forms of this implementation are necessarily weaker, or that the mechanisms
must never return. It shows that this particular bounded policy did not earn
promotion under the current larger-gain priority. The source and controls are
available if future evidence makes revisiting it worthwhile.

## Why richer positional evaluation is next

The remaining major alternative directly changes how positions are judged.
Search chooses among the values supplied at its leaves. Exploring more of a
tree with limited positional distinctions may preserve the same preferences.
The diagnostic's largely unchanged move choices are consistent with that
possibility, but are not proof of causation.

The current evaluator in [evaluation.go](../../../../evaluation.go) includes
material, small placement formulas, and a material-phase blend for the king.
It is already positional; the proposal is to make it board-aware rather than
replace a material-only score that no longer exists.

| Current term | What it distinguishes | What it does not measure |
| --- | --- | --- |
| Material: 100/300/300/500/900 cp | Pawn, minor, rook, queen balance | Compensation and structural quality |
| Knight/bishop centrality | Preferred squares and some advancement | Whether the piece has useful squares or is blocked |
| Central pawn space and advancement | Some central expansion and advanced pawns | Doubled/isolated/passed pawns or pawn support |
| Rook on seventh rank | One useful activity pattern | Open files, blockage, or available routes |
| Small queen centrality | A mild square preference | Development context or safe activity |
| King file/rank plus phase | Home-side placement and endgame centralization | Whether shelter pawns exist or files near the king are open |

For example, the king's current middlegame formula grants 20 cp on flank files
regardless of the pawns actually protecting it. Two kings on the same square
with the same material phase receive the same king-placement term even when
one has a pawn shield and the other is exposed. That is a concrete limitation,
not a claim that every recorded king move had this cause.

The earlier observations of king wandering, queen activity before development,
and awkward pawn moves remain useful questions. Static bonuses should not be
invented merely to make those anecdotes disappear. Tactical necessities can
justify any of those moves. Use controlled positions and game evidence to test
which omissions matter.

## Proposed evaluation design, not yet implemented

Start from the released v0.14.1 playing policy, which is the same as v0.14.
Do not quietly combine new evaluation with the unpromoted reductions. Keeping
the search control stable makes the next result easier to interpret.

### 1. Keep the score contract and build a small board summary

Retain centipawns, terminal mate/stalemate handling, material values, side-to-move
negamax signs, and mate-score dominance. Build the needed board context once
per static evaluation: pawn counts by color/file, pawn locations, king squares,
piece counts, occupancy, and material phase. Use local arrays or existing board
access, with no new chess dependency or permanent evaluation cache initially.

A conceptual interface could be:

```go
type positionalTerms struct {
    material int
    middle   int
    ending   int
    phase    int
}

// Illustration only. No new evaluation terms ship in v0.14.1.
score := terms.material +
    (terms.middle*terms.phase + terms.ending*(24-terms.phase))/24
```

Accumulate a consistent signed White-minus-Black score, then apply the requested
perspective once. The existing phase uses knight/bishop = 1, rook = 2, queen = 4,
24 at the initial position, capped for promotions. Reuse that simple framework
unless a focused test establishes a reason to change it. The proposed full
middle/end blend is a design option; the current blend affects king placement.

### 2. Make pawn quality a property of the actual pawn arrangement

Useful initial definitions are:

- **Doubled pawns:** extra friendly pawns on a file beyond the first. Count extras
  once rather than penalizing both pawns as if each were an additional pawn.
- **Isolated pawns:** no friendly pawn on either adjacent file. This is a file
  relationship, not just an empty neighboring square. Doubled and isolated
  penalties can overlap, so keep their combined influence modest and visible.
- **Passed pawns:** no enemy pawn ahead on the same or adjacent files, using the
  correct forward direction for each color. Advancement and endgame phase can
  increase their value. A passed pawn blocked by another piece is not therefore
  unstoppable; avoid a bonus large enough to ignore material or tactical loss.

Connected/support terms, backward-pawn classification, promotion races, and
special endgame knowledge can wait until there is evidence for them. A coherent
pawn model does not require every textbook pawn feature in its first version.
Original coefficients should be documented before matches. Plausible values
are hypotheses, not tuned facts; do not repeatedly adjust them to win the same
50-game screen.

### 3. Score the king's real shelter and exposure

A small initial model could inspect friendly pawns one and two ranks in front
of the king across its file and adjacent files. Closer pawns can provide more
shelter credit, with the exact definition and coefficients frozen before play.
The same king square should then score differently when its shield is missing.

Also consider open or semi-open files around the king, with care about whether
enemy rooks/queens remain. Penalizing an open file beside an active endgame king
as strongly as beside a middlegame king would conflict with sensible endgame
activity. Scale safety terms by phase or remaining enemy attacking material;
retain the existing endgame centralization preference.

This is a board-derived safety estimate, not a rule that the king must never
move or must always castle. Do not add a generic check penalty as a substitute
for searching legal evasions. Confirmed mates still outrank every positional
preference, and the engine's legal-move layer already enforces king legality.

### 4. Distinguish active pieces from pieces on attractive-looking squares

Centrality is only a rough proxy for activity. A bishop blocked by its own
pawns and a bishop with useful diagonals can occupy the same square. Limited
mobility terms could count available or controlled squares for minor pieces,
then rooks if justified.

Define precisely whether the measure counts pseudo-legal attacks, destinations
excluding friendly occupancy, or safe mobility excluding enemy pawn attacks.
These are different features. An attack map is not a guarantee that every
counted move is legal or tactically safe. Avoid rewarding a queen merely for
having many destinations early; large queen activity weights could encourage
the very wandering we want to investigate.

Do not infer exact move history from pieces returning to their starting squares.
History-dependent evaluation would also complicate score-cache identity. Keep
the first model board-derived, with restrained activity terms and explicit
phase interaction. Bishop-pair, outpost, elaborate rook-file, and attack-network
bonuses are further candidates, not a required collection to add all at once.

### 5. Treat evaluation cost as part of the hypothesis

Static evaluation runs frequently, especially through quiescence. Generating
full legal move lists for both sides at every evaluation could consume the
extra useful work the engine currently completes. Start with reusable local
board summaries and allocation-light feature calculations; measure rather than
assume that more expensive scoring is free.

Keep a feature breakdown available to tests or offline diagnostics so a score
change can be explained: material, placement, pawns, shelter, activity, phase.
This need not add a public UCI option or unsolicited stdout. Incremental
updates, pawn hash tables, and other caching are later optimizations unless
profiling shows they are necessary for the chosen evaluation model.

## How to validate the next evaluation candidate

Use one coherent experiment, with components independently inspectable. A
reasonable first scope is board-aware pawn quality and king shelter under the
existing phase framework, with a restrained activity component only if its
cost and definitions are justified. Full new piece-square tables or learned
weights are not prerequisites. No implementation or weights are approved by
this research release alone.

Before games, add focused checks for:

- Color/rank mirroring, perspective sign, and material invariants.
- Phase limits, smooth intended middle/end transitions, and promoted material.
- Doubled/isolated/passed pawn definitions for both colors, including blocked
  passers and edge files without off-by-one wraparound.
- Equal-material king positions with different real shelter, and endgame king
  activity without an inappropriate permanent shelter penalty.
- Mobility definitions on blocked/open boards and their interaction with other
  terms. Test feature contributions separately when total-score terms overlap.
- Mate/stalemate dominance and existing tactical/forced-move regressions.
- No root mutation, legal UCI output, deadline/cancellation behavior, and cache
  assumptions appropriate to purely board-derived scores.

An exact chosen move is appropriate only when forced or uniquely best under
the test's stated contract. Quiet-position tests should often compare scores,
feature contributions, or an allowed set of moves. Making one arbitrary quiet
move the golden answer would overfit the test to the implementation.

Use varied positions beyond the old stream set. Those public examples have
already been inspected repeatedly and are not genuinely unseen holdouts now.
Keep a fresh validation subset separate from any positions used to choose
coefficients. Lichess exports/puzzles can supply future cases, with game IDs,
sampling rules, and analysis budgets recorded; no database/book/learning system
is required for the next increment. A finite Stockfish analysis is a useful
reference, not an infallible quiet-position oracle.

Measure feature overhead, nodes, completed timed depths, and move choices under
unchanged clocks. Component-disabled probes can identify whether a proposed
term has the intended effect without automatically multiplying full tournaments.
Freeze the candidate before the routine 50 games each against v0.14.1,
Stockfish Skill 2, and pinned Sacrifice, with all three matchups concurrent.
Preserve the source, hashes, profiles, raw games, and compact summaries.

Evaluate the complete result, not just the previous-version gate. With the
current priority, a narrow head-to-head pass and flat external scores should
not become another strength release. A larger visible-gain goal can be stated
before play, with sample-size limits acknowledged. Use larger or serial runs
only when they answer a specific decision, not to rescue an unfavorable screen.
No numerical strength improvement is guaranteed by adding plausible features.

## Reproduce or revisit the archived search

The research release deliberately removes the experimental search and test
changes from normal builds. The full original snapshot and its original binary
hash remain in the results directory. Use a temporary checkout of the recorded
base to replay it; do not assume the research release enables reductions.

From the repository root on a macOS/Linux shell:

```sh
repo_root="$PWD"
experiment_root=$(mktemp -d)
research_dir="$repo_root/tools/benchmark/research/v0.15-selective-search"
record_dir="$repo_root/tools/benchmark/results/v0.15-dev-vs-v0.14-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba"

git archive v0.14 | tar -x -C "$experiment_root"
cp "$research_dir/restore-snapshot.go.txt" /tmp/pandafish-restore-snapshot.go
go run /tmp/pandafish-restore-snapshot.go \
  "$record_dir/candidate-source.json" "$experiment_root"

(cd "$experiment_root" && go test ./...)
(cd "$experiment_root" && go build -o /tmp/pandafish-selective-replay .)
```

To replay the PVS-only control, overlay `pvs-only-source.json` onto a separate
copy of that restored candidate using the same helper, then build it separately.
Never share a score table between modes. For the node artifact, run inside the restored candidate checkout:

```sh
PANDAFISH_SELECTIVE_AUDIT=/tmp/node-audit.json \
  go test -run '^TestSelectiveNodeAudit$' -count=1 .
```
The positions, audit drivers, referee settings, match command, and their paths
are documented in the results report.

Compiler/platform differences, build metadata, random root ties, and scheduling
can change binary hashes, selected moves, and timings. Reproduction means
recovering the recorded policy and settings; bit-identical game outcomes are
not promised. Historical unfavorable scores stay in place even if a later run
or a later algorithm performs differently.

The release itself makes no new playing-strength claim. Its purpose is to keep
the investigation useful: what was expected, what happened, what remains
uncertain, and a concrete next hypothesis that can be tested independently.


## Research-release verification

[Release checks](research-release-checks.json) record the separate validation of
v0.14.1. The normal engine passed tests, race checks, vet, ARM64 build, and an
offline 12-ply lichess-bot wrapper check. The archived candidate was restored
into a temporary v0.14 checkout, its tests passed, and it built with the original
v0.15-dev identity. Production search/evaluation/timing source was checked
against v0.14 and remained unchanged.

Six short release smoke games completed without forfeits or unfinished games.
The two-game previous-version score returned a blocked strength decision; all
results were retained and no favorable-result retry was made. Those smoke
checks validate integration, not a strength gain. No live Lichess games or
account actions were performed. Local private config greetings were updated
for the research-release identity without changing credentials or policies.
