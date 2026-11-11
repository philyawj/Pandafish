# Lichess session tool (planned)

This folder reserves the place for a future Go tool. No capture, launcher,
parser, or analysis commands have been implemented. Pandafish v0.7 and the
current bot configuration are unchanged by this folder setup.

## Intended operation

The recommended flow uses one terminal for a recorded session. A future
`record` command would launch lichess-bot as its child process and display
the client's output in that terminal. lichess-bot would continue to handle
all Lichess communication and start the normal Pandafish executable.

The tool would prepare an ignored session directory, snapshot the engine
binary, record its identity and hash, and configure recording paths for PGNs
and timestamped client logs. This keeps later engine rebuilds from changing
the build used by an active session. Private configuration stays local and
credentials are excluded from manifests and reports.

After the active games finish and the client exits cleanly, the tool would
finalize the session's files. A separate `report` command would then parse
those saved records into per-move JSONL and a readable summary. Reporting
would run offline, without the bot or a live collector still running.
Detailed Stockfish analysis could be a later, separate addition.

Proposed command shapes, not working commands yet:

```text
go run ./tools/lichess record [session and client options]
go run ./tools/lichess report [saved session directory]
```

A second terminal would be optional for inspecting logs. The initial design
does not require independently starting lichess-bot in one terminal and
attaching a collector from another. The recording command would own the
client's lifecycle, paths, and session metadata. Its shutdown behavior must
be verified when implemented so session finalization follows client shutdown.

## Planned engine reporting

Use the same Pandafish executable for normal and recorded play. A future
increment would add standard UCI search information and compact diagnostics;
the current engine does not yet expose the full planned metric set. Full
principal variations also require additional engine work.

Keep data capture and reporting separate from move selection. Measure the
overhead of instrumentation, preserve raw source records, and validate an
offline captured session before starting live collection. PGNs, logs, and a
version manifest are the source records; JSONL and summaries are derived.

Raw sessions belong under `local-games/` at the repository root and remain
gitignored. Curated regression positions or public summaries can be added
later. The [benchmark tool](../benchmark/README.md) remains the tool for
controlled strength comparisons.

## Current scope

This directory move does not start a live listener, change challenge rules,
enable matchmaking, or implement a new engine feature. The public example
still accepts casual standard games from humans, with outgoing matchmaking
and pondering disabled. Rated bot play requires later configuration work.
