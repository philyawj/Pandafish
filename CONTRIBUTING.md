# Contributing

I welcome GitHub issues and pull requests! Feel free to report bugs, suggest
ideas, or help improve Pandafish.

The [README](README.md#milestones-and-plans) outlines the direction I'm currently
taking. For larger changes, please open an issue first so we can discuss the idea.

Keep pull requests focused, describe what changed, and mention how you tested it.
The README also covers setup and the checks to run before submitting a PR.

For changes intended to improve playing strength, follow the
[benchmark workflow](tools/benchmark/README.md#test-a-development-candidate).
Include the candidate source revision, full command, comparison version,
settings, and results in your PR. A 100-game screening run is useful evidence,
but does not by itself establish a small strength improvement. Run the Go
tests and include focused positions that exercise the behavior changed.

Maintainers follow the [release workflow](RELEASING.md) when preparing and
publishing a version.
