# Developer tools

These Go tools live in the engine's module but build as separate programs.
Run their commands from the repository root. `go build -o pandafish .` builds
only the engine, and `go test ./...` checks the root module's Go packages.
The separate chess module also requires `go -C third_party/chess test ./...`.
The research validation plan and benchmark validation include nested tests/vet;
research validation also includes races in both modules.

| Directory | Purpose | Status |
| --- | --- | --- |
| [benchmark](benchmark/README.md) | Controlled local UCI matches, validation, and recorded results | Implemented |
| [lichess](lichess/README.md) | Record sessions, archive raw data, and run offline Stockfish analysis | Implemented; first live session archived |
| [research](research/README.md) | Verify evidence, replay archived Go engines, check frozen decisions and mate proofs | Implemented in v0.19 |

All program/test/research source must be Go. The
[v0.19 summary](../docs/RESEARCH_V019.md) records retired historical script
workflows, their hashed release attachment and the limits of the Go replacements.
The [research tooling flowchart](../docs/RESEARCH_TOOLING.md) maps the shared
Go execution, reporting and evidence commands used by new investigations.
The [evidence-storage policy](../docs/EVIDENCE_STORAGE.md) keeps existing data
in place and defines lean records for new work without changing historical schemas.

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
