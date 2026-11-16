# Start here in a fresh context

## Current state

v0.15.1 is a **research and reporting-tool release**. Playing behavior is still
v0.15: material plus small placement formulas, timed iterative-deepening
alpha-beta, bounded capture quiescence with exchange filtering, transposition
table, killer ordering, mate-distance scores and threefold-history handling.
No new evaluation hypothesis or quiet-check extension ships in this release.

One long v0.15 session was collected and processed, not several. It ran on an
M2 MacBook Air in a terminal, with light browsing/music and no stream. All raw
records are under this directory's `raw/`, with hashes in `inventory.json`.
The compiled captured engine remains ignored and available only locally.

269 unique games, 265 completed (77W/28D/160L), four unfinished. All 12,091
searches completed legally with metrics; no fallback or own time-forfeit loss.
No overlapping search requests. Client limit/reconnect incidents are explained
in [CLIENT_NOTES](CLIENT_NOTES.md), not treated as engine crashes.

## Read before changing the engine

1. [REVIEW.md](REVIEW.md): ranked hypotheses, uncertainty and exact next scope.
2. [Replay diagnostics](analysis/engine-replay/README.md): five persistent
   two-choice mistakes, one extra-ply correction, ties, and failed counterfactuals.
3. [TESTING.md](TESTING.md): 50-game screens are noisy; proposed confirmation
   stage and why it is not a new default yet.
4. [PROCESS.md](PROCESS.md), repository [benchmark guide](../../README.md)
   and [release workflow](../../../../RELEASING.md).

Do not recreate AGENTS.md or ENGINE_PLAN.md. Public checked-in workflow and
research documents are the source of truth. Keep benchmark evidence out of
the main README's functional release descriptions. Avoid em dashes in prose.

## Next engine candidate, not yet implemented or approved for release

Start from the tagged research release and try **one compact piece-safety
evaluation**. First trace quiet-leaf threats in 16G1f9rx (Nd4) and HhPI3vCx
(Qd6); preserve tactical counterexamples, both colors and mirrored controls.
Determine whether an attack/defender term actually explains their misvaluation.
Do not patch the root move by name or bundle king safety, pawn structure and
new pruning. Keep existing material, search and timing fixed. If the mechanism
does not help, record that early rather than force it through a match suite.

The other priorities are narrow pawn endings/passers, coherent king exposure,
phase-sensitive placement, then bishop pair/rook files. Gain guesses in the
review are low-confidence scenarios, not established Elo. Previous speculative
bundles and a one-line home-rank fix failed to establish broad improvement.
The quiet-check diagnostic timed out on both king-capture probes and is not
a validated feature to copy into production.

Split the 30 public confirmation cases before tuning a candidate: for example,
freeze alternating game IDs as development/holdout. They have already been
viewed, so call this a public regression holdout, not blind validation. Include
the older v0.10 regression set and synthetic controls. Check several legal
alternatives, fixed depth and equal time; Stockfish's move is not necessarily
unique. Preserve every result and report cost as well as chosen moves.

## Existing approved local workflow

Build the previous executable from tag v0.15.1, not from a modified checkout.
Run formatting, tests, vet, race checks, ARM64 build and an offline lichess-bot
wrapper sequence. Then use the normal routine for a real engine candidate:

```sh
go run ./tools/benchmark --validate --race \
  --version v0.16-piece-safety-dev \
  --previous tools/benchmark/bin/pandafish-v0.15.1 --previous-version v0.15.1 \
  --stockfish /absolute/path/to/stockfish --stockfish-skill 2 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --pairs 25 --parallel-all
```

This runs **50 games each, three matchups concurrent**. Pinned Sacrifice is
v0.2.0 at `b528f0ba373910b72052d991700a990ca1cb07fd`, including the earlier
philyawj contribution. Keep the previous-release criterion and report all
three results, including failures. Do not revert to serial mode because of
stale private instructions. Stockfish remains Skill 2; a live account rating
does not advance that benchmark target.

The proposed additional 200-game confirmation is documented but not silently
made mandatory or coded as a default. Decide the planned sample before running.
Fifty games cannot reliably distinguish the review's modest gain guesses.
No engine-strength match was run for this research-only release.

The next implementation should finish with a recommendation and check-in
before preparing its PR. Do not launch another live client automatically.
Existing private bot configuration and root engine remain the recorded v0.15
setup; this analysis did not change the playing profile or start online games.

Suggested next prompt: "Read the v0.15.1 live research HANDOFF and REVIEW.
Implement the narrow piece-safety hypothesis, validate the mechanism and cost,
run our approved local comparisons if justified, and check in before a PR."
