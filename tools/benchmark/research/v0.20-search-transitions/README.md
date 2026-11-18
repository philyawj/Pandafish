# Archived transition-only candidate

Promotion is rejected. Read the frozen [plan](PLAN.md),
[implementation](IMPLEMENTATION.md), [result and reassessment](RESULT.md), and
[environment findings](ERRATA.md).

The candidate removes all-piece update scans while retaining immutable state,
original attack/tag calculation, and exact search policy. The median selected
fixed-work elapsed reduction is 21.1%; all 325 distinct score/move/node
comparisons agree. Public and newer-position native gates pass at both budgets.
The standard 50-game-per-opponent concurrent screen nevertheless scores 45%
against tagged v0.19 and fails the practical criterion. All 150 games, runner
exit 1, source, hashes, CPU profiles, tests and native/reference queries remain.
No larger confirmation or unchanged retry follows.

The [screen](screen/README.md) preserves full settings and results.
`screen-pgn-check.json` records legal replay of every saved game. The candidate
label is `v0.20-search-transitions-dev`; its native UCI name remains v0.19.
Its binary SHA-256 is
`89d254677d32f4b3ab86106eb174096c381578414615c5042a75d1ef913defb9`.
Source is frozen in `candidate-source.json` and `source-hashes.json`, based on
tagged v0.19. The complete unsuccessful candidate is independently recoverable:

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.20-search-transitions/candidate-source.json \
  --out /tmp/pandafish-transition-replay --check
```

The subsequent [legality-kernel revision](../v0.20-search-legality/README.md)
reuses the functionally verified transitions but removes additional attack/tag
work under a new, higher cost gate. Its results do not replace this failed
screen. No quiet-mate proof, changed evaluation, new pruning or different clocks
are part of either backend attempt.
