# Developer tools

These Go tools live in the engine's module but build as separate programs.
Run their commands from the repository root. `go build -o pandafish .` builds
only the engine, and `go test ./...` checks all implemented Go packages.

| Directory | Purpose | Status |
| --- | --- | --- |
| [benchmark](benchmark/README.md) | Controlled local UCI matches, validation, and recorded results | Implemented |
| [lichess](lichess/README.md) | Record Lichess sessions and report on saved games and diagnostics | Planned; no executable yet |

The benchmark's code, tests, opening suites, milestones, and results moved
together from `benchmarks/` to `tools/benchmark/`. Historical result contents
and metadata are preserved. Its compiled programs and engine baselines remain
ignored in `tools/benchmark/bin/`.

The planned Lichess tool will use the existing lichess-bot client and the
normal Pandafish executable. Raw session files will stay in the ignored
`local-games/` directory. Each tool has its own documentation and purpose;
shared helpers can be extracted if both programs need the same functionality.
