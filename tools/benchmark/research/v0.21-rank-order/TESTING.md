# Rank-order candidate validation

The [plan](PLAN.md) and [profile-driven selection](../v0.21-next-assessment/SELECTION.md)
precede implementation. All source, tests and research drivers are Go. The
runtime changes are an additional private per-ply scratch buffer and the move
ordering computation; upstream/library code, values, search/TT policy and clocks
remain unchanged.

## Correctness and identity

`rank-controls.txt` in the assessment directory preserves initial controls.
The unchanged ordering helpers provide the reference stable result for legal
full/tactical moves on the existing seeded 500-position suite. New permutation
controls compare whole Move values against an independent reference sort over
2,100 length/tie patterns, including empty/single lists and nontrivial cycles.
Reuse controls include captures, EP, promotions, castling, mate, different plies
and repeated calls. Whole-value assignments preserve tags and PGN metadata.

`race.txt`, `vet.txt` and `build.txt` pass for the full module and native engine.
Existing mate, repetition, TT/history, perft, UCI/cancellation and move ownership
controls remain active. Library files do not change, so the new experiment does
not add a separate attack or transition implementation requiring another oracle.

`candidate-source.json` contains three Go overlays; `source-hashes.json` records
their bytes. Its `base_tag` field holds the immutable merged Skill 3 tooling
commit `6bfb4d565c2adf2a2ad59216a462b0f4f0c01503`, whose playing engine is v0.20.
The replay command accepts that Git revision. The measured engine is labelled
`v0.21-rank-order-dev`, native UCI name `Pandafish v0.20`, binary SHA-256
`d9f2aee7c30832e00690e1c24f1490e27b7dcfe0dda6810858d71dd4bb232f1f`.
Baseline SHA-256 is
`c4ac32d3360db5ff813f94ba2906e67422b0d5ad8b43380499ce6222a61e9a8d`,
built from released v0.20. Development labels do not identify a future tag.

## Selected fixed work and profiles

The saved Go timing driver records every command/exit and refuses old output
paths. Three alternating baseline/candidate depth-four pairs plus depth-three,
depth-six and newer-position probes preserve all 325 distinct scores, root moves
and node counts. `timing-analysis.json` retains all pairs and the frozen 15%
median gate. All pairs pass that threshold; the median reduction is 18.9%.
The deeper subset and newer roots improve as well. No favorable outlier is chosen.

The direct-root driver compares identical internal terminal returns; production
terminal/null behavior has separate native/unit controls. Profile capture uses
Go 1.25 and the installed local pprof tool renders it. Profiles are selected
workloads, not full-game strength measures. Adaptive benchmark iteration counts
differ (220 versus 261 on public inputs), so their partial cyclic case mixes and
aggregate B/op are not controlled allocation comparisons. Fixed-work paired
timings are the efficiency evidence; the candidate adds reusable rank scratch.

## Equal-time gates

Four suites contain 1,252 native actor requests: public 84 roots at three 100ms
trials and one 1000ms trial, newer 145 roots at one trial at each budget. Native
blocks allow three concurrent searches, alternate actor order and finish before
full-strength reference queries. Existing matching 100,000-node teacher caches
retain their source hashes and original observations. The game opponent is
Skill 3; the labelling teacher is separately full strength with Hash 64 MiB.

All legal production root moves remain eligible. Only teacher labelling of a
selected move uses `searchmoves`. The 25cp band accepts close alternatives;
mate-risk and fully finite CP cases remain separate. Finite subsets can differ
between budgets/candidates, so compare the paired actors within each report.
The 145 reused roots are inspected regression controls, not a blind holdout.

All four reports pass, including each newer root turn, with no new losing-mate
flag or response past allowance plus 50ms. The saved `choice-review.json`
independently recomputes raw scores, legal selections, cache chains, aggregate/
turn loss, severe cases and timings using the
[archived Go reviewer](../v0.20-next-assessment/review_choices.go.txt).
These diagnostics do not prove that every finite-reference position was savable.

## Match and final review boundary

The one standard screen runs 50 games each versus tagged v0.20, Stockfish Skill
3 and pinned Sacrifice, all three concurrent, with unchanged openings, 100ms
allowance, 50ms grace and a 400-ply cap. No other engine probes/heavy validation
overlap it. Full identities, commands, source revisions, raw games and timings
are retained. Final review legally replays the PGNs and independently checks
statistics, color pairs, options and attempted-response counts including any
forfeited attempt. Failures remain failures even when data integrity passes.

Reproduce in a new scratch destination:

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.21-rank-order/candidate-source.json \
  --out /tmp/pandafish-rank-replay --check
```

The base revision must exist locally. Saved driver paths refer to the original
scratch workspaces; adapt them and use new outputs when reproducing. Archive
replay checks source/tests/build, not bit-identical executables or game outcomes.
No larger confirmation, live game, experimental commit, push or PR is automatic.
