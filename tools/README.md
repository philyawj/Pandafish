# Developer tools

These Go tools live in the engine's module but build as separate programs.
Run their commands from the repository root. `go build -o pandafish .` builds
only the engine, and `go test ./...` checks all implemented Go packages.

| Directory | Purpose | Status |
| --- | --- | --- |
| [benchmark](benchmark/README.md) | Controlled local UCI matches, validation, and recorded results | Implemented |
| [lichess](lichess/README.md) | Record sessions, archive raw data, and run offline Stockfish analysis | Implemented; first live session archived |

The benchmark's code, tests, opening suites, milestones, and results moved
together from `benchmarks/` to `tools/benchmark/`. Historical result contents
and metadata are preserved. Its compiled programs and engine baselines remain
ignored in `tools/benchmark/bin/`.

The Lichess tool uses the existing lichess-bot client and the
normal Pandafish executable. Raw session files will stay in the ignored
`local-games/` directory. Each tool has its own documentation and purpose;
the engine and session tool share the recording schema in `internal/telemetry/`.

Live-session data and written results live under
[tools/benchmark/live](benchmark/live/README.md). Tool source stays here in its
own program directories. Controlled matches remain in benchmark/results.
