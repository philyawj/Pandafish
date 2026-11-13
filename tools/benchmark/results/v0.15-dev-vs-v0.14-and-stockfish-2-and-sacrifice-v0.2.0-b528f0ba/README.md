# v0.15 selective-search experiment

**Do not promote this search candidate as a strength improvement.**
v0.14.1 archives the research while restoring v0.14's playing policy.
PVS and guarded late-move reductions reduced work and sometimes completed more
depth, but did not demonstrate the substantial playing-strength gain sought.
The narrow previous-release pass is insufficient for that goal.

## Assessment and scope

The [fresh assessment](../../research/v0.15-selective-search/README.md) reconciles
the main README, older live priorities, and the post-v0.12 reference comparison.
The older groundwork shipped in v0.11 through v0.13; v0.14 changed the workflow.
Selective main search was chosen as the larger search-depth hypothesis, with
richer classical evaluation as the alternative. Small ordering/memory tweaks
were deliberately deferred. No numerical Elo gain was promised.

PVS searches the first move normally, probes later moves with a one-centipawn
window, and widens a promising result when necessary. The selected LMR rule
reduces later probes by one ply only at remaining depth at least three and
move index at least three (fourth move or later). It excludes captures, en
passant, promotions, checking moves, all evasions, matching TT/killer hints,
and mate-valued alpha windows. Root moves are not reduced directly.

A reduced improvement must repeat at full nominal depth before a cutoff; an
interior improvement receives a full-window re-search. Descendants may still
be selective. A reduced fail-low can miss a quiet tactic. Full-width minimax
scores/tie sets and ordering independence are not guaranteed with LMR. Root
comparison margins and random ties among computed scores are retained.
Evaluation, quiescence, clocks, TT layout, and UCI options are unchanged.

## Functional checks and ablations

Formatting, full Go tests, `go test -race ./...`, vet, and build passed through
the Go validation command. Normal `go test ./...` and vet passed separately.
Focused tests cover PVS scores/windows and full root tie sets against the
legacy unreduced search, LMR guards, shallow fail-high refutations after
full-depth verification, cancellation during verification/wide re-search,
incomplete-parent cache safety, and root depth protection. Existing iteration,
fallback, UCI, and tactical tests also passed.

The old killer ordering audit initially failed because it assumed
ordering-independent exact scores with reductions enabled. It now isolates
unreduced PVS ordering; new tests check selective behavior separately. The
starting-position root depth-four fixture used reductions but no verification;
the direct full-window depth-four fixture exercised five verified/refuted
shallow fail-highs. The test now uses each context for the behavior it actually
exercises. No policy parameter was changed to improve the diagnostic or score.

[Four cold depth-four root ablations](node-audit.json):

| Mode | Nodes | Score/tie comparison |
| --- | --- | --- |
| Legacy full width | 223,056 | Reference |
| PVS only | 206,935 | Exact scores and complete root tie sets preserved |
| PVS plus LMR | 149,131 | Selective; one root score changed from 187 to 192 cp |

The selected policy visits about 33.1% fewer nodes than the full-width reference
and 27.9% fewer than PVS-only on those four positions. Node savings are not a
strength estimate. Reproduce with `PANDAFISH_SELECTIVE_AUDIT=/tmp/node-audit.json
go test -run '^TestSelectiveNodeAudit$' -count=1 .` from the captured source.

[Thirteen-position audit](position-audit.json): all 78 probes legal, no errors.
Each engine receives a cold depth-four probe and an independent 100-ms probe,
serially: tagged v0.14, PVS-only, and PVS-plus-LMR. Fresh hash/table per process.
All 13 PVS-only fixed-depth scores match v0.14. The selected policy completes a
deeper timed iteration in five cases, the same in seven, and a shallower one in
one. Two deeper cases and the shallower case were already mate-scored by the
referee. Only one selected move differs from v0.14 across the 13 timed cases.

The referee is Stockfish 19, full strength, 100,000 nodes per independent
unrestricted/restricted query, one thread and 64 MiB, with cleared hash/new game
and score/PV association checks. Nine centipawn-scored cases have summed signed
estimated loss 1,560 cp for v0.14, 1,427 for PVS-only, and 1,428 for selected LMR
at 100 ms. Four mate cases are excluded. Finite-budget scores and random ties
limit this selected diagnostic; it does not show a broad move-quality gain or
isolate the cause of the match result. No tuning followed these probes.

## Match method

Use the standard 50 games per opponent, with all three matchups concurrent:
25 color pairs, unchanged ten-opening suite, 100-ms request, 50-ms protocol
allowance, 400-ply cap. The first five openings repeat three times, the rest
twice; each pair swaps colors. The previous release is immutable tagged v0.14;
Stockfish 19 uses native Skill Level 2, Threads 1 and Hash 16 MiB; pinned
Sacrifice is v0.2.0 b528f0ba with the verified Tal PGN asset.

```sh
caffeinate -i go run ./tools/benchmark --validate --race --version v0.15-dev \
  --previous tools/benchmark/bin/pandafish-v0.14 --previous-version v0.14 \
  --previous-revision 895e8b41f3676dcbb342360f2a9f7af1a565ed37 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision 895e8b41f3676dcbb342360f2a9f7af1a565ed37-dirty \
  --machine 'Apple M2 MacBook Air; AC power reported at start; three concurrent matchups; no live bot' \
  --verbose
```

Hardware is an Apple M2 MacBook Air. AC power was reported before validation
and at completion, battery 61% to 57%, state charging. These are isolated
observations, not sustained temperature/frequency measurements. No other engine
experiment or live bot ran alongside matches. The match portion took about
543 seconds. No reruns, extended matches, or parameter sweeps were used.

The exact [candidate source](candidate-source.json) is relative to v0.14 commit
895e8b41f3676dcbb342360f2a9f7af1a565ed37. Audit and match binary hashes match:
`9fecf0414ee2da57b2f65145075d5464edd8f9be0908e3eba69a902d46a7492e`.
The [PVS-only source changes](pvs-only-source.json) apply over the candidate
with its mode fixed for that isolated process. Modes are internal ablations,
not runtime UCI options; no table is shared between different policies.

## Results and interpretation

| Opponent | W/D/L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| v0.14 | 16/19/15 | 0 | 0 | 51% |
| Stockfish Skill 2 | 22/8/20 | 0 | 0 | 52% |
| Sacrifice v0.2.0 b528f0ba | 8/18/24 | 0 | 0 | 34% |

All 150 games completed without forfeits or unfinished games. The practical
previous-release criterion passes by a narrow margin; this does not establish
improvement. The same newer-profile reference screen, using v0.13's playing
algorithm retained by v0.14, scored 60% against Skill 2 and 34% against Sacrifice.
The candidate is down eight percentage points against Skill 2 and flat against
Sacrifice. These separate small samples do not prove regression or equivalence.

Before scores were reviewed, the visible-gain goal was at least +10 percentage
points on one external opponent without a decline greater than five on the
other. This descriptive goal is not met. A modest node-count gain or higher
reported nominal depth is not a substitute for the requested strength gain.
The cause of the flat/lower scores is not isolated: selective fail-lows can
miss quiet tactics, extra depth may not change a move, and evaluation/sample
variation can matter. Do not invent a confirmed explanation from this screen.

Ten of 8,187 candidate responses exceeded 100 ms; the maximum was 134.232 ms.
No response from either side exceeded the 150-ms protocol deadline. Reliability
checks passed; observed timings are not guarantees.

## Decision and next direction

Keep this as a local experiment and preserve its unfavorable evidence. Do not
prepare a PR or release claiming a substantial strength improvement. There is
no automatic longer run to rescue the result. PVS remains a verified supporting
mechanism and LMR remains a future research option, but their combination here
does not justify promotion under the current priority.

The next preferred major hypothesis is richer classical evaluation, using the
released v0.14 playing policy as the control rather than stacking unproven
reductions. Start with a coherent board-aware pawn/king/activity weakness and
broad position checks, not isolated arbitrary bonuses or micro-optimizations.
That is a recommendation after this result, not another implementation started
by this experiment. No live session is needed to make this local decision.

The candidate passed an ARM64 build and 12 legal plies through the offline
lichess-bot wrapper, plus public-example config loading with temporary resolved
engine paths and a placeholder token. The active v0.14 executable was restored
byte-for-byte. Neither private config changed; their v0.14 greetings remain
accurate. No live listener, challenges, or online games were started. No engine-promotion PR or release was prepared at the end of the experiment.
The subsequent v0.14.1 research release preserves these records and archives
the source outside normal builds, without shipping the unpromoted policy.

## Saved records

- [Compact game results](compact-results.json) and [complete summary/timings](summary.json)
- [All games](games.pgn), [v0.14](previous.pgn), [Stockfish](stockfish.pgn), [Sacrifice](sacrifice.pgn)
- [Node ablations](node-audit.json), [position aggregates](position-summary.json),
  [raw position probes](position-audit.json), and [inputs](position-inputs.json)
- [Exact candidate source](candidate-source.json) and [PVS-only changes](pvs-only-source.json)
- Standalone position-audit Go source: [driver](position-audit-main.go.txt),
  [engine client](position-audit-engines.go.txt), [referee](position-audit-stockfish.go.txt)
- [Power preflight](preflight.json)

Copy standalone `.go.txt` files into a temporary directory with `.go` extensions.
From the repository root, run the three sources with an output JSON path. Build
the tagged baseline and captured candidate at the paths in the driver, and the
isolated PVS-only variant at `/tmp/pandafish-v015-pvs-only`. The driver reads
`testdata/quiescence-validation.json`, checkpoints after each case, and records
input, engine, and referee hashes. Adjust executable paths for the local system.


## Research archive

The [expanded report](../../research/v0.15-selective-search/README.md) includes
code examples, hypotheses, lessons, the non-promotion decision, detailed next
evaluation recommendations, and Go restoration commands. v0.14.1's engine keeps
the released v0.14 search; the original candidate source, tests, and hashes above
remain unchanged. Candidate tests are replayed from the archive rather than
compiled into the research release's normal test suite.
