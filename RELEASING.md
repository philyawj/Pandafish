# Releasing Pandafish

Use a branch and PR for each release. Keep source changes, tests, documentation,
and benchmark evidence together. Run commands from the repository root.

## Prepare and review

1. Implement the increment and add focused tests. For changes to move selection
   or time allocation, follow the
   [benchmark workflow](tools/benchmark/README.md#test-a-development-candidate):
   100 games against the previous release first. If the conservative screen
   passes, play 100 each against the current Stockfish skill and pinned
   Sacrifice version. Run all matchups serially, Stockfish before Sacrifice,
   to keep timed searches from competing for resources. Longer
   comparisons are optional when more evidence is needed.
   For tooling or recording changes that preserve those policies, use focused
   unit/race tests, offline capture/report checks, local wrapper validation,
   and a short real-engine smoke run. State that this validates integration
   rather than strength; a fresh 100-game comparison is not required.
2. Review check results, W/D/L, unfinished games, forfeits, and timings. Keep
   the candidate's original benchmark labels and hashes. Reference those
   results from the PR and release notes; do not rename a development run as
   a test of a released binary.
3. Set the release version in the UCI identity in `main.go`, its handshake test,
   and the example bot greetings. Update the main README's current version and
   completed milestone. Update benchmark examples to use the new release as
   the next baseline. Keep credentials and private configs out of Git.
   Keep the main README focused on functionality, setup, and short release
   milestones. Put benchmark scores, sample sizes, conditions, notable test
   games, and interpretation in `tools/benchmark/README.md` and its results report.
   Link to those documents from the main README instead of repeating evidence
   in each release milestone.
4. Run formatting, tests, vet, a build, and `git diff --check`. If only version
   text and documentation changed since benchmarking, use a short smoke test
   instead of another full match. For example, while preparing v0.9:

   ```sh
   go run ./tools/benchmark --validate --race --version v0.9-release-smoke \
     --previous tools/benchmark/bin/pandafish-v0.8 --previous-version v0.8 \
     --stage previous --pairs 1 --results /tmp/pandafish-release-smoke
   go run ./tools/benchmark --candidate tools/benchmark/bin/pandafish-candidate \
     --version v0.9-release-smoke --stage external \
     --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
     --pairs 1 --results /tmp/pandafish-release-smoke
   git diff --check
   ```

   These explicit stages check all three integrations without treating a
   two-game smoke score as a new strength gate. They supplement the already
   reviewed 100-game candidate comparisons.

   Also perform the local lichess-bot wrapper check when using that client.
   These checks do not start a live bot or play online games. If engine behavior
   changes after testing, repeat the routine comparison for the new candidate.
5. Inspect the staged files, commit, push, and open a PR describing the change
   and checks performed. Review and merge the PR into `main` before tagging.

## Tag and publish after merge

Confirm the PR is merged and the version is unused, then fast-forward local
`main` to the merged commit. For v0.9:

```sh
git switch main
git fetch origin
git merge --ff-only origin/main
git status --short
```

Check that the working tree is clean and HEAD is the intended release commit.
Then create and push the tag:

```sh
git tag -a v0.9 -m "Pandafish v0.9: positional evaluation"
git push origin v0.9
```

Publish a GitHub release for that tag, through the GitHub UI or
`gh release create`. Use the version alone as the release title, such as
`v0.9`. Describe the engine change, remaining limits,
and validation, linking the candidate benchmark evidence. An engine release
does not automatically advance the Stockfish skill target.

## Build the next baseline

Build from the tag rather than a development checkout. For v0.9:

```sh
mkdir -p tools/benchmark/bin
benchmark_bin="$PWD/tools/benchmark/bin"
build_dir=$(mktemp -d)
git archive v0.9 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.9" .)
```

The tagged executable becomes the previous version for v0.10 development.
Compiled binaries remain gitignored. Building this baseline does not require
relabeling or rerunning the development benchmarks. A separate benchmark of
the exact released binary is optional. Keep private bot greetings in sync
when choosing to run the new engine through lichess-bot.
