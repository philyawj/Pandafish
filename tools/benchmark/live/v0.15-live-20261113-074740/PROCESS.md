# Processing and preservation

## Inputs and provenance

One session: `local-games/v0.15-live-20261113-074740`. It began at
2026-11-13 15:47:40 UTC and finished at 2026-11-16 02:41:53 UTC, still
November 15 in the operator's local Pacific time. It was terminal-only with
light browsing and YouTube music, not streamed. The original manifest's
`terminal` label is retained. No second session was merged.

| Input | Identity |
| --- | --- |
| Engine source | clean v0.15, `d7a61498a16c8dcdd5582b85406eb568917189d7` |
| Captured executable SHA-256 | `26f9af01670cfcb1bae9c261adf731cf7e2726034b717323130f8d0422abfdca` |
| Source inventory SHA-256 | `f137e4228b61907f58a8b2a397cd3a27d3f944cf6d5a55a1840dd385ac9442bb` |
| Full-strength referee | Stockfish 19, SHA-256 `dc2f18c34ae962dff591b66147d220ec06e61d756b93d8a4f5e04fd8e55c251f` |
| Client commit | `baa2b7e78320b99125a91c2b1805bb7a5c4fb6ea`, dirty working tree |

The exact client patch was not recorded. Source hashes do not reconstruct
that dirty state. The manifest, redacted config, original PGNs, complete log,
270 traces, moves and original reports are copied under `raw/`. All 547
inventory entries were verified against their SHA-256 hashes. The captured
binary stays locally available but ignored; no executable is committed.
The separate derived inventory covers research outputs and documentation.

No raw bytes were rewritten to improve counts or fix stale report wording.
In particular the original report's "mate sign" sentence predates the v0.15
mate-distance correction; raw UCI output and the versioned format are decisive.
The current formatter now describes that version distinction correctly.

## Deduplication and full survey

The archive has 269 PGN records and 269 unique Lichess game IDs. There are no
duplicate IDs to remove. Retain repeated opponents, positions and games with
different IDs as occurrences, and expose their correlation instead of silently
discarding them. The candidate set contains 1,503 occurrences and 1,461 unique
first-four-FEN-fields plus played-move keys. That is a different deduplication
question from unique games or full-history diagnostic cases.

The existing processed survey was reused rather than issuing another 24,208
positions' worth of referee work. It uses full-strength Stockfish, one thread,
64 MiB hash, Skill 20, Elo limiting off, MultiPV 1, 100,000 nodes per query.
Scores compare the unrestricted reference and the actually played root move
at equal budgets. Each query starts with fresh hash/search state. CP and mate
remain separate. Bound/reference consistency handling and raw output live in
the per-game JSON and [method](analysis/stockfish-n100000/method.json).

The repeatable first pass is the existing Go command, from the repo root:

```sh
go run ./tools/lichess process \
  --session local-games/SESSION-NAME \
  --archive tools/benchmark/live/SESSION-NAME \
  --bot GoPandaFish --stockfish /absolute/path/to/stockfish --nodes 100000
```

`process` archives records, extracts operational incidents and checkpoints the
offline survey. Read its completion counts before requesting a review. Preserve
the survey even when a candidate set is uninteresting. Do not run a new tool
version over an old raw archive expecting its inventory to stay unchanged;
use a fresh archive destination or a temporary session copy for derived reports.
No additional Lichess API requests or live games were needed for this review.

## Frozen larger-budget selection

The old automatic 20-position confirmation list contains repeated decisions.
It remains preserved. For this review, freeze 24 playable-deterioration cases
(eight in each phase) plus six conversion opportunities, earliest ply then
survey severity/game ID within each group. Require one case per game and
deduplicate full UCI history plus played move. This yields the saved
[selection](analysis/stockfish-n500000-confirmed/selection.json) and
[focus input](analysis/stockfish-n500000-confirmed/focus.json), before any
500,000-node results. Phase thresholds are documented in the review.

The completed confirmation command was:

```sh
archive="tools/benchmark/live/v0.15-live-20261113-074740"
go run ./tools/lichess process \
  --session local-games/v0.15-live-20261113-074740 --archive "$archive" \
  --bot GoPandaFish --stockfish /opt/homebrew/bin/stockfish --nodes 500000 \
  --focus "$archive/analysis/stockfish-n500000-confirmed/focus.json" \
  --profile stockfish-n500000-confirmed
```

This used the unchanged source records before the linker repair. Reproduce a
new analysis in a scratch destination, retaining the frozen focus input; do
not overwrite these reference outputs. The 30 selected cases are targeted,
not an independent accuracy test or a benchmark sample.

## Conservative audit and linkage repair

The standalone Go audit uses complete histories, FENs and returned moves to
recover individually unique matches, then consistent process-boundary evidence.
It also computes cohort/phase counts, repeated decisions and search interval
overlap. [audit.go.txt](provenance/audit.go.txt) is executable source archived
as data. Run on a scratch copy to preserve derived output:

```sh
archive="$PWD/tools/benchmark/live/v0.15-live-20261113-074740"
audit_dir=$(mktemp -d)
cp "$archive/provenance/audit.go.txt" "$audit_dir/audit.go"
cp -R "$archive" "$audit_dir/session"
go run "$audit_dir/audit.go" "$audit_dir/session"
```

It finds 406 recoverable links beyond the original 11,615. The production
reporter was then fixed to ignore unmatched requests when intersecting game
evidence, while keeping genuine ambiguity and conflicting evidence visible.
The full-session check runs against a temporary copy:

```sh
derived_dir=$(mktemp -d)
cp -R "$archive/raw" "$derived_dir/session"
go run ./tools/lichess report --session "$derived_dir/session"
```

Only the resulting `moves.jsonl`, `summary.json` and `summary.md` were copied
to `analysis/linker-fixed-report/`. This yields 12,021 links and 70 unresolved
requests, agreeing with the independent audit. Original raw counts remain
11,615 and 476. Some PGN moves have no uniquely identified capture; that does
not prevent their offline evaluation. No timestamp guess is used to resolve
the two identical games. Tests cover trailing unplayed requests, contradictory
process evidence, identical histories, transpositions and game boundaries.

## Decision probes and sampling analysis

[Replay instructions](analysis/engine-replay/README.md) restore the exact
diagnostic driver, evaluator breakdown and two source overlays. The normal
project build excludes all archived `.go.txt` and overlay JSON. The recorded
binary was used for baseline probes. Neither counterfactual was promoted or
given a strength benchmark.

The Go sampling script reads the original v0.15 candidate versus v0.14.4
summary. It uses deterministic PCG seed 151/20261115 and 100,000 bootstrap
replicates for color pairs and opening clusters; it does not run games or
invoke an LLM:

```sh
cp "$archive/provenance/statistics.go.txt" "$audit_dir/statistics.go"
go run "$audit_dir/statistics.go" \
  tools/benchmark/results/v0.15-conversion-dev-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json
```

The output is [benchmark-sampling.json](analysis/benchmark-sampling.json).
No ongoing agent analysis is required during long runs: Go programs capture
and checkpoint the data, then a reviewer reads compact summaries and selected
positions. Raw data remain available for independent interpretation.

## Release verification

See [VALIDATION.md](VALIDATION.md) for actual checks and their outcomes. No
engine policy changed and no full strength match was rerun. Research/tooling
smoke checks must not be described as a new strength gain. The active root
engine and private bot config remain the v0.15 live setup; offline verification
does not start lichess-bot's online client or send challenges.
