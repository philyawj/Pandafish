# Releasing Pandafish

Use a branch and PR for each release. Keep source changes, tests, documentation,
and benchmark evidence together. Run commands from the repository root.

## Prepare and review

1. Implement the increment and add focused tests. Follow the
   [benchmark workflow](benchmarks/README.md#test-a-development-candidate):
   100 games against the previous release and 100 against the current Stockfish
   skill. Longer comparisons are optional when more evidence is needed.
2. Review check results, W/D/L, unfinished games, forfeits, and timings. Keep
   the candidate's original benchmark labels and hashes. Reference those
   results from the PR and release notes; do not rename a development run as
   a test of a released binary.
3. Set the release version in the UCI identity in `main.go`, its handshake test,
   and the example bot greetings. Update the main README's current version and
   completed milestone. Update benchmark examples to use the new release as
   the next baseline. Keep credentials and private configs out of Git.
4. Run formatting, tests, vet, a build, and `git diff --check`. If only version
   text and documentation changed since benchmarking, use a short smoke test
   instead of another full match. For example, while preparing v0.3:

   ```sh
   go run ./benchmarks --validate --race --version v0.3-release-smoke \
     --previous benchmarks/bin/pandafish-v0.2 --previous-version v0.2 \
     --pairs 1 --results /tmp/pandafish-release-smoke
   git diff --check
   ```

   Also perform the local lichess-bot wrapper check when using that client.
   These checks do not start a live bot or play online games. If engine behavior
   changes after testing, repeat the routine comparison for the new candidate.
5. Inspect the staged files, commit, push, and open a PR describing the change
   and checks performed. Review and merge the PR into `main` before tagging.

## Tag and publish after merge

Confirm the PR is merged and the version is unused, then fast-forward local
`main` to the merged commit. For v0.3:

```sh
git switch main
git fetch origin
git merge --ff-only origin/main
git status --short
```

Check that the working tree is clean and HEAD is the intended release commit.
Then create and push the tag:

```sh
git tag -a v0.3 -m "Pandafish v0.3: depth-two minimax"
git push origin v0.3
```

Publish a GitHub release for that tag, through the GitHub UI or
`gh release create`. Describe the engine change, remaining limits,
and validation, linking the candidate benchmark evidence. An engine release
does not automatically advance the Stockfish skill target.

## Build the next baseline

Build from the tag rather than a development checkout. For v0.3:

```sh
mkdir -p benchmarks/bin
benchmark_bin="$PWD/benchmarks/bin"
build_dir=$(mktemp -d)
git archive v0.3 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.3" .)
```

The tagged executable becomes the previous version for v0.4 development.
Compiled binaries remain gitignored. Building this baseline does not require
relabeling or rerunning the development benchmarks. A separate benchmark of
the exact released binary is optional. Keep private bot greetings in sync
when choosing to run the new engine through lichess-bot.
