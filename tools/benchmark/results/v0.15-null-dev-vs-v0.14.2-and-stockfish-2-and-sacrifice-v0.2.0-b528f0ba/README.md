# v0.15 guarded null-move screen

## Decision: do not promote

The candidate passed the direct previous-release criterion but did not deliver
consistent external strength. Keep the released playing behavior. Source and
all 150 games are preserved in the v0.14.3 research release. The candidate itself
was not promoted, and its search policy is excluded from the normal build.

| Opponent | W / D / L | Score | Unfinished | Forfeits |
| --- | --- | --- | --- | --- |
| v0.14.2 | 18 / 20 / 12 | 56% | 0 | 0 |
| Stockfish 19 Skill 2 | 16 / 4 / 30 | 36% | 0 | 0 |
| Sacrifice v0.2.0 b528f0ba | 3 / 17 / 30 | 23% | 0 | 0 |

Draws count as half a win. Each opponent had 50 games. Historical released-policy
external scores were 60% against Skill 2 and 34% against Sacrifice in the
[concurrent calibration](../v0.13-parallel-check-vs-v0.12-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md).
These separate screens do not establish a precise Elo change or a causal
explanation. They provide no basis for claiming the substantial gain sought.
No favorable retry, weight tuning, or automatic longer match followed.

## Hypothesis and fresh reference review

The [pre-match assessment](../../research/v0.15-null-move/README.md) compares pinned
Stockfish, Sunfish and Sacrifice sources, ranks alternatives, and explains why
null-move pruning took priority over another evaluation variant. Stockfish and
Sunfish use guarded null probes. Sacrifice's reviewed search instead supplies
examples of PVS, reductions, history ordering and mate-distance handling.
Reference-engine strength does not establish the value of one borrowed concept
in Pandafish.

Keep released evaluation, quiescence, ordering, clocks and table layout. A
nonterminal node at depth at least three may temporarily pass if it is not in
check, has friendly non-pawn material, has static score at least beta, has a
non-mate beta, and did not immediately follow a pass. Search depth is reduced by
two additional plies. A non-mate fail-high returns beta only. The hypothetical
subtree cannot read or write real-position table entries or killer hints.
Cancelled work is not a score. The library pass clears en passant and changes
turn while preserving the board/castling rights and immutable parent.

This is an original small heuristic, not a reproduction of Stockfish's complete
node guards, reductions or verification. Excluding pawn-only endings does not
eliminate non-pawn zugzwang or every tactical false cutoff. The match results do
not identify which mechanism caused losses. Do not blame zugzwang, evaluation,
or timing without a controlled replay of actual failing positions.

## Diagnostics and checks

The fixed 19 saved/synthetic cases produced identical depth-four scores with
pruning disabled or enabled. Nodes fell from 387,832 to 279,122 (28.0%). Separate
100 ms iterative probes reached greater depth in four cases and equal depth in
fifteen, with no shallower case. Two timed scores differed. These selected cases
were too narrow to establish strength or rule out tactical errors elsewhere.
Reference and candidate tables were fresh and were never shared across policies.

Focused tests cover eligibility for both sides, shallow/check/pawn-only/mate
exclusions, pass state, immutable parents, pass-bound behavior, table isolation,
cancellation and node savings. Full tests, race tests, vet, ARM64 build and the
offline 12-ply wrapper check are recorded in `checks.json`. The wrapper used the
same candidate hash as match play, validated the resolved public placeholder
example, and restored the active v0.14.1 binary. Both private configs were
unchanged. No online game or live listener was used.

## Reproduction and retained data

```sh
caffeinate -i go run ./tools/benchmark \
  --candidate tools/benchmark/bin/pandafish-candidate --version v0.15-null-dev \
  --previous tools/benchmark/bin/pandafish-v0.14.2 --previous-version v0.14.2 \
  --previous-revision cd763f7920de0e9d523a99a7319cb65a6c36318f \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision cd763f7920de0e9d523a99a7319cb65a6c36318f-dirty \
  --parallel-all --machine 'Apple M2 MacBook Air; three concurrent matchups; no live bot' \
  --verbose
```

The Go validation preflight passed format, race tests, vet and build. Match play
used the separately built same-source candidate. Defaults were 25 opening pairs,
100 ms per move, 50 ms timeout grace and 400-ply cap. All three matchups ran
concurrently. The complete window was about 587 seconds. The final power
observation reported AC; continuous thermal or power data were not collected.

Across 8,434 candidate responses, 56 exceeded 100 ms and none exceeded 150 ms.
Maximum response was 131.688 ms; matchup p95 was about 96.4-97.1 ms. These runner
measurements include the local exchange, not live network or client submission.
No timeout forfeits occurred in this screen.

`summary.json` preserves exact identities, options, hashes, timings and games.
`previous.pgn`, `stockfish.pgn`, `sacrifice.pgn` and combined `games.pgn` preserve
all 150 games. `screen-summary.json` is the compact view. `candidate-source.json`
is a complete four-file Go overlay on tagged v0.14.2; replay with the shared
[restore helper](../../research/v0.15-selective-search/restore-snapshot.go.txt)
into a temporary `git archive v0.14.2` checkout, then run its tests and build.
`position-audit.json` contains the raw fixed-depth/timed control measurements;
reproduce with `PANDAFISH_NULL_AUDIT=/tmp/null-audit.json go test -run '^TestNullAudit$' -count=1 .`
in a restored checkout. Random ties, clocks and toolchains can change outcomes.
The hash manifest covers retained result files.

## What we learned and next recommendation

Two search-saving experiments and a richer evaluator have now failed to show the
large, consistent gain sought. That does not mean little strength remains. It
means feature presence in stronger engines and improvements on selected positions
are insufficient grounds for promotion. Do not automatically add another pruning
rule or combine the unsuccessful changes.

Use the saved local losses and drawn winning positions for a small, explicit
failure audit before the next implementation. Compare pruning enabled/disabled
at the same fixed depth and time, then use full-strength Stockfish on selected
critical positions. Distinguish selective-search errors, horizon limits, poor
static preferences and conversion failures rather than assuming one explanation.
This diagnostic work has not yet been performed on these match losses.

Mate-distance scoring and repetition-aware conversion are concrete remaining
gaps. Current mate scores have no distance and search has no repetition draw
score. They deserve focused fixtures and controlled conversion checks; correct
transposition normalization and history context are part of the work. Their
existence does not prove they caused this screen's result or promise a large Elo
gain. A smaller pawn/king evaluator remains an alternative after that assessment.
