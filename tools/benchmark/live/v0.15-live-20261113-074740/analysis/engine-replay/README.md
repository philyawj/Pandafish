# Captured-engine decision replay

Eight selected 500,000-node referee cases were replayed through the exact
captured v0.15 binary. [inputs.json](inputs.json) preserves full histories,
root FENs, live requests, elapsed time and completed depth. The selection is
diagnostic, not a random sample or strength test.

Each probe starts a fresh process, verifies move legality and has an independent
20-second watchdog. Two-choice searches use the played and referee moves at
live depth and depth plus one; the other move is then searched alone at the
extra depth. Scores below are Pandafish's own cp scores, not Stockfish scores.

| Game | Live depth | Pair at live depth | Pair at extra ply | Other move alone at extra ply |
| --- | --- | --- | --- | --- |
| 16G1f9rx | 7 | Nd4 -12 | Nd4 -19 | O-O -25 |
| O4wRsvjl | 6 | Kf7 -25 | Kf7 -11 | Bf5 -13 |
| HhPI3vCx | 6 | Qd6 +161 | Qd6 +159 | Qe8 +93 |
| B1WTqhJy | 4 | f6 -8 | f6 -8 | Bxf3 -26 |
| DsqFd5Kj | 8 | Rd5 +14 | Rd5 +25 | Rd1 +22 |
| L0mwQLiR | 7 | Nxh4 +128 | Ra7 +59 | Nxh4 +16 |
| PNfQoYP0 | 6 | h4 +112 | h4 +112 | Bf3 +112 |
| rpk4dbYK | 6 | Kxg5 +400 | Kxg5 +400 | Kf6 +215 |

All 24 probes completed legally. Five cases retained the played move without
an equal alternative at the extra ply, one switched, one already selected the
reference at live depth, and one tied. Live root ordering, a warm transposition
table and other root choices can change results. These paired cold probes do
not establish full-root corrections or failures to reproduce the live engine.
Extra-ply searches also need more time; HhPI3vCx's pair took about 9.67 seconds.
Do not assume every extra ply fits its live budget.

[results.json](results.json) includes all UCI output, requests, identities,
binary hashes and legality results. The copied released evaluator produces
[evaluation-breakdown.json](evaluation-breakdown.json), separating immediate
material and placement contributions. Child static scores are not searched
values. The geometric passer list is diagnostic, not a validated passer term.

## Counterfactuals, not promoted candidates

[home-rank-source.json](home-rank-source.json) overlays two files on v0.15:
one king-placement condition and a diagnostic identity. At full phase the
original king formula scores e1/e8 -8 and f2/f7 0. Restricting the shelter-side
bonus to the home rank makes Bf5 preferred in O4wRsvjl at depth six but still
Kf7 at depth seven. Kxg5 remains preferred at both depths. See
[home-rank-results.json](home-rank-results.json). This is a real incentive
with an inconclusive one-line fix, not an established strength gain.

[qchecks-source.json](qchecks-source.json) overlays a diagnostic identity and
bounded quiet-check quiescence, allowing at most two such additions along a
path while keeping existing caps. Nd4 still wins the pair at depth seven.
Both rpk4dbYK probes failed to return `bestmove` before the 20-second watchdog;
the saved error is `missing UCI response: bestmove`. They have no inferred
score. See [qchecks-results.json](qchecks-results.json). This expansion is not
a safe or demonstrated solution merely because it is bounded.

Neither overlay is compiled by the repository's normal build or shipped in
v0.15.1. No head-to-head match was run for them. Preserve these negative probes
to avoid rediscovering them as apparently validated features.

## Reproduce offline

Run from the repository root, using a scratch copy so original outputs survive:

```sh
archive="$PWD/tools/benchmark/live/v0.15-live-20261113-074740"
replay_dir=$(mktemp -d)
cp "$archive/provenance/replay-main.go.txt" "$replay_dir/main.go"
cp "$archive/provenance/replay-engines.go.txt" "$replay_dir/engines.go"
cp -R "$archive" "$replay_dir/session"
go run "$replay_dir/main.go" "$replay_dir/engines.go" "$replay_dir/session"
```

This requires the ignored `raw/pandafish` executable on the operator's machine.
For a clone, build tag v0.15 in a temporary checkout and put the result there.
The source policy can be reproduced; a local rebuild may have a different
binary hash or timing. Never substitute the current modified candidate and
call it the original captured binary.

For a counterfactual, extract tag v0.15 into another scratch directory and
write the named `files` from the corresponding overlay JSON over that source.
Build the home-rank binary to `/tmp/pandafish-v015-home-rank` or the quiet-check
binary to `/tmp/pandafish-v015-qchecks`, as the saved diagnostic drivers expect.
Replace the copied `main.go` above with `home-rank-main.go.txt` or
`qchecks-main.go.txt` from provenance, keeping `engines.go`, then run against
the scratch session. Do not apply overlays in the working repository. The
drivers and helper are archived as data; results include their failure cases.
