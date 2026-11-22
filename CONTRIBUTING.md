# Contributing

Bug reports, suggestions and focused pull requests are welcome.

The [README](README.md#milestones-and-plans) describes the v1.0 milestone and
development pause. Active engine development is paused while effort focuses on
other projects; bug reports and focused fixes remain welcome. The
[version history](docs/HISTORY.md) preserves earlier milestones.
Discuss larger changes in an issue before implementation.

Keep pull requests focused, describe the change and include its validation.
The README covers setup; [engine behavior](docs/ENGINE.md) and
[Lichess setup](docs/LICHESS.md) provide the detailed operating guides.
All repository program, test and research source must be Go. Follow the short
[agent guidance](AGENTS.md). Use neutral, factual engineering prose without
first-person voice or personal framing, and exclude its prohibited role word
from PR/release text.

For changes intended to improve playing strength, follow the
[benchmark workflow](tools/benchmark/README.md#test-a-development-candidate).
Include the candidate source revision, full command, comparison version,
settings, and results in your PR. Starting after v0.13, the routine screen is
50 games per opponent (`--pairs 25`), with the three matchups concurrent, to
speed up development. Use the documented defaults and preserve the run profile.
A screen is
useful evidence, but does not by itself establish a small strength improvement.
Run the Go
tests and include focused positions that exercise the behavior changed.

Before submitting a PR, format changed Go files and run tests, race checks and
vet in both modules, then build the engine and check whitespace:

```sh
go test ./...
go test -race ./...
go vet ./...
go -C third_party/chess test ./...
go -C third_party/chess test -race ./...
go -C third_party/chess vet ./...
go build -o pandafish .
git diff --check
```

Engine changes also need local UCI/client integration checks. No live games start
as part of validation. Keep credentials, private configuration, compiled binaries
and local handoff files out of commits. Inspect `git diff --cached` before publishing.

Follow the [release workflow](RELEASING.md) when preparing and publishing a version.
