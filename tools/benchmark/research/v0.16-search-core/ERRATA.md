# Preserved development failures and corrections

- `tests-initial.txt`: my new test helper initially tried to compare library
  `Move` values with `slices.Equal`, but the type is not comparable. I changed
  the assertions to compare UCI identity and relevant move tags. The same
  authoring pass corrected the color-mirrored queen fixture and made deadline
  injection expire after the child push even when no tactical continuation
  exists. `tests-second.txt` and the full race suite then passed. These were
  test-authoring corrections, not suppressed engine failures.
- `inputs-check.txt`: my first corpus command expected `stockfish-2.pgn`; the
  saved filename is `stockfish.pgn`. No inputs were written before the error.
  `inputs-check-corrected.txt` records the corrected 84-case corpus. Selection
  rules were unchanged and no candidate measurement preceded the correction.
- `artifact-audit-initial.json`: the first final-artifact audit found a link to
  the per-run README before I had written it. The runner saves raw PGNs/summary
  and the global index; I added the explanatory result report, then repeated
  the audit. No engine test, timing probe or game was rerun or altered.
- `release-audit-initial.json`: the release artifact helper overescaped the
  identity normalization and rejected `main.go`. The separate release validator
  had already verified the identity-only difference. I corrected that helper;
  no playing code or measured result changed.
