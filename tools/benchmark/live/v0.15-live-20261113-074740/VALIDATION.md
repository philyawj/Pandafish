# v0.15.1 release validation

The playing algorithm is unchanged from v0.15. Only its UCI identity/test and
public example greeting version change. The functional Go patch is confined to
`tools/lichess/report.go`, with focused regression tests in `report_test.go`.
The archived counterfactual overlays are data, not part of the build.

## Completed checks

- `go test ./...`: all four packages passed.
- `go test -race ./...`: all four packages passed.
- `go vet ./...`: passed.
- `gofmt` on changed Go files and `git diff --check`: passed.
- Native `GOOS=darwin GOARCH=arm64 go build`: Mach-O ARM64 executable.
- The companion lichess-bot configuration validator loaded the public example
  with a resolved temporary executable path. Its actual UCI wrapper verified
  v0.15.1 identity, 12 legal plies, mate-in-one distance and a full-history
  threefold-saving move. No online client was started.
- `record --prepare-only` created an offline session using the release binary
  and a placeholder-token example config. No API connection or challenges.
- Real Stockfish processor preview: two plies from one original PGN at
  1,000 nodes, in a separate temporary fixture. This checks archive/process
  integration, not completeness or strength. Its [summary](provenance/offline-processor-smoke-summary.json)
  explicitly records the partial preview.
- Full-session report on a temporary copy: 12,091 completed searches,
  12,021 unique links, 70 unresolved, no illegal moves or missing metrics.
- All 547 original inventory hashes verified; credential-pattern scan found
  no matches in proposed archive text/JSON/PGNs. The raw executable remains
  ignored. Authored Markdown links and table column counts checked.

## Six-game concurrent smoke

The [saved summary](analysis/release-smoke/v0.15.1-release-smoke-vs-v0.15-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json) and adjoining PGNs preserve all six games:

| Opponent | W / D / L | Unfinished | Forfeits |
| --- | --- | --- | --- |
| v0.15 | 0 / 0 / 2 | 0 | 0 |
| Stockfish Skill 2 | 1 / 1 / 0 | 0 | 0 |
| Pinned Sacrifice | 0 / 2 / 0 | 0 | 0 |

This uses one color pair per opponent, three matchups concurrent, 100 ms moves
and the normal 50 ms protocol allowance. The runner returned exit status 1
because the two-game previous-release score blocked its screen. All six games
completed without forfeits; that verifies integration, not a favorable strength
result. There was no rerun for a better score. Identity-only source changes
and random ties/time-limited search make this far too small to infer a policy
regression or improvement. No new routine strength benchmark was required.

The release-smoke executable hash is in
[processing provenance](provenance/processing-environment.json). It was built
before committing, so its build identity is a dirty development tree, not a
claim to be the subsequently tagged binary. The active ignored root engine
still has the exact captured v0.15 hash. Private configs were not rewritten.
