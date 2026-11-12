# Capture and offline analysis process

## What happened

The operator ran the released v0.10 engine with the Go recorder and the private
GoPandaFish stream configuration on the Apple M2 MacBook Air, alongside Twitch
and the browser. Intended concurrency was one, pondering off. Incoming standard
chess accepted humans/bots, rated/casual, all supported clocks; outgoing
matchmaking selected bots at rated 3+2. The manifest records normal client exit
and a finished session. The operator reported unexpected simultaneous games.

The original ignored session is `local-games/v0.10-twitch-m2-20261112-093212`.
The existing Go `process` command copied all 40 files, including the executable,
byte-for-byte into `raw/`. The inventory hashes and sizes cover everything;
the executable remains locally retained and gitignored. No Lichess API game
export was necessary. All 14 PGNs, traces, log, report, manifest, and redacted
config already existed locally. Private API configuration was not copied.

The first processing attempt stopped on a duplicate linked search at
`glYlALQg/18`. This was a useful recording anomaly, not a reason to discard a
search. The offline processor was corrected to retain multiple attempts per
position and unlinked completed attempts, attach all duplicates to PGN analysis
without choosing a canonical submission, and include all attempts in operational
timing. It also recognizes the explicit unlabelled control-stream reconnect
marker once per traceback. Engine and lichess-bot sources were unchanged.

Focused tests cover duplicate/unlinked searches, cached analysis resumption,
raw preservation, longest unlinked response, and reconnect counting without
counting chat quotations or every traceback line as a separate error.
`go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check`
passed. The session tool was rebuilt and exercised against real Stockfish on
a one-ply, 1,000-node fixture outside tracked history. No live check was run.

A 20-ply pilot preceded the complete resumed 100,000-node survey. A second
profile rechecked 12 selected positions at 500,000 nodes. All execution was
serial. The survey covers 1,092 played plies, including 548 GoPandaFish moves.
The [review](REVIEW.md) explains selection, scores, confirmation, and limits.
The 100,000-node analysis is an unrestricted referee, separate from Skill 1
strength matches. Do not describe it as Stockfish Level 1 analysis or standard
Lichess accuracy. No per-move LLM judgment was used to calculate scores.

## Reproduce the analyses

From the repository root, with Go 1.25+ and Stockfish installed:

```sh
go run ./tools/lichess process \
  --session tools/benchmark/live/v0.10-twitch-m2-20261112/raw \
  --archive tools/benchmark/live/v0.10-twitch-m2-20261112 \
  --bot GoPandaFish --stockfish stockfish --nodes 100000

go run ./tools/lichess process \
  --session tools/benchmark/live/v0.10-twitch-m2-20261112/raw \
  --archive tools/benchmark/live/v0.10-twitch-m2-20261112 \
  --bot GoPandaFish --stockfish stockfish --nodes 500000 \
  --profile stockfish-n500000-confirmed \
  --focus tools/benchmark/live/v0.10-twitch-m2-20261112/analysis/stockfish-n100000/confirmation.json
```

The original executable used was `/opt/homebrew/bin/stockfish`, Stockfish 19,
SHA-256 `dc2f18c34ae962dff591b66147d220ec06e61d756b93d8a4f5e04fd8e55c251f`.
Each profile's `method.json` records identity and settings. Matching checkpoints
resume; another binary/settings combination needs a different profile name.
Queries use one thread, 64 MiB hash, MultiPV 1, Skill 20, no Elo limit, full
position history, and fresh hash/game state per independent query. Played and
reference moves are scored from the same root at equal budgets. A reference
whose returned move differs from its PV is independently re-scored.

The committed records omit the executable as listed in the inventory. The
processor supports that omission for published archives. Analysis does not need
the captured engine. Exact engine replays do require the locally retained binary
with the matching hash. Preserve existing generated records when experimenting;
write supplementary diagnostics to another output path.

## Supplementary diagnostics

[Source provenance](provenance/processing-environment.json) records the dirty
processor base, binary/source hashes, post-session client inspection, and helper
hashes. [The processor patch](provenance/offline-processor.patch) preserves the
changes used. Diagnostic Go helpers are stored as `.go.txt` to keep historical
sources outside the build. To repeat the linkage/overlap audit:

```sh
cp tools/benchmark/live/v0.10-twitch-m2-20261112/provenance/session-audit.go.txt /tmp/pandafish-session-audit.go
go run /tmp/pandafish-session-audit.go \
  tools/benchmark/live/v0.10-twitch-m2-20261112/raw \
  /tmp/pandafish-capture-diagnostics.json
```

It matches full move histories, FENs, and the next PGN move individually and
computes intervals from each request timestamp plus monotonic response duration.
It does not rewrite the original 40 unlinked records or measure CPU usage.

The [depth replay](analysis/engine-depth-replay/README.md) uses only the three
development cases. Copy `provenance/depth-replay.go.txt` and
`provenance/played-depth-replay.go.txt` individually to temporary `.go` files
and run each with an output JSON path. Their fixture and captured binary paths
are explicit in the source. Each starts fresh engine processes, restricts root
choices, and has a 20-second watchdog. This is diagnostic evidence, not the
next engine implementation or a full-root strength match.

## Client error investigation

[The client note](UPSTREAM_CLIENT_NOTE.md) records exact errors, original line
ranges, UTC/PST interpretation, worker/slot evidence, confirmed source behavior,
and remaining inference. Exact contiguous excerpts retain terminal padding.
Read-only upstream research found the matching already-fixed issue #1266 and
PR #1262. Pinned metadata is under `client-investigation/`; full third-party
source was fetched only to a temporary file, not vendored into this repository.

To repeat the narrow offline probe, fetch `source_url` in
`client-investigation/upstream-check.json` to `/tmp/pandafish-lichess-bot-upstream.py`,
then run the archived script against an installed companion checkout:

```sh
env PYTHONDONTWRITEBYTECODE=1 \
  /Users/philyawj/projects/lichess-bot/venv/bin/python \
  tools/benchmark/live/v0.10-twitch-m2-20261112/client-investigation/reproduce.py.txt \
  --client /Users/philyawj/projects/lichess-bot \
  --upstream-source /tmp/pandafish-lichess-bot-upstream.py \
  --output /tmp/pandafish-client-reproduction.json
```

This historical probe expects the installed old revision and pinned newer
source. Adjust the local checkout/interpreter paths on another machine. If the
client has since been updated, use a separate historical checkout for the old
side. The script mocks jobs and streaming, blocks HTTP requests, and imports
only selected function bodies from the fetched source with installed globals.
It does not run online or test the entire newer client. The duplicate scheduling
assertion fails on the installed old code and passes on current upstream code;
first starts still work for reserved challenge IDs. The script asserts these
expected differing outcomes, so a successful script exit records both results.

The initial analysis did not modify the client. After the cause was verified,
the user authorized updating the companion checkout. It was fast-forwarded to
2026.11.12.1 (`baa2b7e`), including #1262. All 23 selected upstream offline tests
and the 12-ply Pandafish v0.10 wrapper check passed. Private configs were verified
unchanged by hashes and permissions. The [update record](client-investigation/update-validation.json)
contains commands and identities. No custom client fix, upstream issue/PR,
listener, or account action was performed. A new upstream contribution needs a
reproduction on the updated version.

## Future sessions

Use the existing [session workflow](../../../lichess/README.md) to prepare,
record, stop, and process a fresh named session. One `process` command archives,
validates, resumes Stockfish, and writes reports; do not manually export every
game or hand-score positions. Preserve failed/interrupted records too. Keep
machine/workload notes and the exact engine/client identities. Treat intended
concurrency as a configuration, and inspect the actual intervals before claiming
an isolated live performance comparison.
