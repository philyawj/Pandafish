# Lichess session recording

A Go tool to launch lichess-bot with a fixed Pandafish build, preserve local
records, and report on them offline. lichess-bot handles Lichess communication
and game scheduling. This tool does not implement a second bot client or use
Stockfish to choose live moves.

## Setup

Use Go 1.25 or newer, a working lichess-bot checkout with its Python environment,
a BOT account, and a private client configuration. This launcher supports
macOS and Linux and a native UCI engine without an interpreter or engine
command-line options. Build Pandafish v0.8 or newer to get trace support;
older release binaries do not write these traces.

Run from the Pandafish repository root:

```sh
go build -o pandafish .
go build -o bin/pandafish-lichess ./tools/lichess
```

The YAML dependency is used by the session tool only. The engine and benchmark
runner remain Go programs. No database or additional collector is required.

## Prepare without going online

This copies and probes the engine, writes a manifest and a redacted configuration,
and creates empty recording directories. It does not launch lichess-bot, contact
Lichess, or play games:

```sh
./bin/pandafish-lichess record \
  --client /absolute/path/to/lichess-bot \
  --config config-mybot.yml \
  --session local-games/setup-check \
  --prepare-only
```

`--client` and `--session` are required. Config paths are relative to the client
checkout unless absolute. The engine defaults to the executable named in that
config; `--engine /absolute/path/to/pandafish` overrides it. Python defaults to
`CLIENT/venv/bin/python`; use `--python` for another environment. Each session
must have a new directory. A preparation directory cannot later be resumed;
use a fresh name for the real session.

Load the prepared `config.redacted.yml` with the client's `load_config` and use
its UCI wrapper for an offline smoke check before collecting live games. It
contains a redacted token, so it is not a configuration for live play.

## Record a session

The following command starts the live client. Run it when ready to accept games:

```sh
./bin/pandafish-lichess record \
  --client /absolute/path/to/lichess-bot \
  --config config-mybot.yml \
  --session local-games/first-session
```

One terminal is enough. The launcher runs lichess-bot as its child, displays the
client's output, and saves timestamped console output. A second terminal is
optional for watching files. Do not start another copy of lichess-bot for the
same account alongside this session.

The launcher copies the engine into the session directory, records its UCI
identity, SHA-256, and available Go build revision/dirty flag, and points the
session config at that copy. Rebuilding the original engine does not affect an
active session. It also records the host platform, client Git revision/dirty flag, and redacted config
hash. Keep the session copy unchanged; `report` checks those two artifact hashes.

The original configuration is never rewritten. Only the engine executable,
PGN directory/grouping, and token delivery are changed for recording. Relative
resource paths continue to resolve from the client checkout. Challenge rules,
matchmaking, pondering, and other engine options are preserved. The current
GoPandaFish policy remains casual standard games from humans with outgoing
matchmaking and pondering disabled. Rated games against bots require a separate
policy change.

The token uses the same precedence as lichess-bot: a nonempty
`LICHESS_BOT_TOKEN`, then the config's token. Unset a stale environment value
before running if the token in the file should be used. The selected token is
passed through the child's environment, not command-line arguments. The runtime
config contains a placeholder for that token and is removed on normal shutdown
or a handled startup failure. The permanent snapshot redacts credential keys,
and known credential values are redacted from console capture, including values
split across writes. Treat all raw records as private and review before sharing;
a hard kill can leave the temporary runtime file behind, and arbitrary secrets
embedded in unrelated config fields cannot be identified automatically.

Press Ctrl+C once to request client shutdown. The launcher forwards an interrupt
and waits for the client to exit before finalizing the manifest. Whether active
games finish depends on the client's `quit_after_all_games_finish` setting; set
it to `true` in your private config for that behavior. A second Ctrl+C reaches
the client's force-quit handling. SIGTERM also requests the first graceful
interrupt. There is no automatic kill timeout, so a stuck client needs manual
attention. Abrupt termination leaves a partial session that can still be reported.

## Report after the session

```sh
./bin/pandafish-lichess report --session local-games/first-session
```

Reporting runs offline. It does not start the client, contact Lichess, or need an
API key. Run it again after parser changes; source records are preserved and
only the derived report files are replaced. Read `summary.md` first, then inspect
`moves.jsonl` or the source records for flagged issues.

```text
local-games/first-session/
  manifest.json          # Build/config identity, session status, client exit code
  pandafish              # Copied engine used throughout the session
  config.redacted.yml    # Recording config without credential values
  lichess-bot.log        # Timestamped client console output
  games/*.pgn            # Client's per-game PGNs
  traces/engine-*.jsonl   # Separate protocol trace per engine process
  moves.jsonl            # Derived search records and PGN linkage
  summary.json           # Compact machine-readable counts and warnings
  summary.md             # Human-readable report
```

All of `local-games/` is gitignored. Session configs and records use regular
0644 permissions; the copied engine is executable. Keep sessions in that ignored
directory. Curated regression positions or reviewed public summaries can be
added to the repository later.

## What the records mean

`record` sets `PANDAFISH_TRACE_DIR` for the client's engine processes. The normal
Pandafish executable writes a unique trace per process. Without that variable it
creates no files. Trace events preserve exact UCI input and output, sequence
numbers, UTC timestamps, and monotonic elapsed nanoseconds. There is no proxy or
second live capture service. Trace write failures do not suppress UCI moves;
reporting flags the resulting partial records. Node counting and final UCI search summaries also
work during ordinary UCI play.

Every returned move has a compact `info string pandafish` JSON summary, plus
standard UCI depth, nodes, elapsed milliseconds, material score when available,
and the chosen move as a one-move PV. Metrics include the allocated soft/hard
budgets, completed depth, total visited nodes across iterations including
quiescence and interrupted work, elapsed search time, stopping reason, and legal
fallback status. Abandoned searches remain identifiable from protocol events;
they do not have final metrics or a returned move. Score and depth always refer
to the last complete iteration. Material scores are centipawns from the root
side's perspective. Mate scores record a sign, not an uncomputed mate distance;
fallbacks do not invent a score. Full multi-move PVs remain future work.

The report validates returned move legality and counts completed, cancelled,
incomplete, and illegal searches, fallbacks, hard-budget overruns, missing
metrics, and unmatched searches. It summarizes completed depths, stop reasons,
response p95, and maximum hard-budget overrun. Response time is measured from command-loop
receipt of `go` to staging `bestmove`; search time ends before output and file
I/O. Neither includes client submission, acknowledgement, or network latency.
Small deadline overruns can occur while a node finishes. Use the raw clock
inputs and timing records to investigate rather than treating every overrun as
a lost game or time forfeit.

PGN linkage uses the recorded history, root FEN, returned move, and all searches
within the same process/game boundary. The engine does not receive Lichess game
IDs through UCI. Common opening prefixes can be ambiguous; later searches can
resolve them. Missing, truncated, and ambiguous records are flagged rather than
assigned a guessed game. Reconnected games may have several engine processes.
The client must actually save a PGN before linkage is possible. Preserve partial
sessions; startup validation processes with no searches are normal.

The client log is supporting evidence for disconnects, submission errors, and
process failures, not a structured network-latency measurement. The first report
does not classify every client error, calculate Elo or win probability, or run
Stockfish analysis. Keep controlled strength comparisons in the
[benchmark tool](../benchmark/README.md). A rating mixes opponents and versions
and does not isolate an engine improvement.

## Development checks

```sh
go test ./...
go test -race ./...
go vet ./...
git diff --check
```

Tests exercise offline preparation, immutable copying, policy preservation,
redaction across writes, graceful signal forwarding, client failure, exact
request/metric parsing, missing/truncated and duplicate events, PGN ambiguity,
and repeatable reporting. No test connects to Lichess. For engine changes, also
run the project's ARM64 build and local lichess-bot wrapper check. Capture a short
local wrapper session and confirm its searches and PGN link correctly before
starting live collection. Compare fixed-depth legality/scores and observe timing
with and without tracing; do not require a new strength tournament for recording
changes alone.
