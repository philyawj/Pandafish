# Contributing

Bug reports, suggestions and focused pull requests are welcome.

The [README](README.md#milestones-and-plans) outlines the project direction.
Discuss larger changes in an issue before implementation.

Keep pull requests focused, describe the change and include its validation.
The README also covers setup and the checks to run before submitting a PR.
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

Follow the [release workflow](RELEASING.md) when preparing and publishing a version.
