# Releasing Pandafish

Use a branch and PR for each release. Keep source changes, tests, documentation,
and benchmark evidence together. Run commands from the repository root.

## Prepare and review

1. Implement the increment and add focused tests. For changes to move selection
   or time allocation, follow the
   [benchmark workflow](tools/benchmark/README.md#test-a-development-candidate):
   starting after v0.13, run 50 games each against the previous release,
   Stockfish Skill 3, and pinned Sacrifice, with all three matchups concurrent.
   Review the previous-release criterion after play and keep all results even
   when it fails. Defaults match this routine on hosts with sufficient capacity.
   Use `--sequential` or larger samples later when investigating smaller gains
   or timing, with a specific reason rather than automatic escalation.
   For tooling or recording changes that preserve those policies, use focused
   unit/race tests, offline capture/report checks, local wrapper validation,
   and a short real-engine smoke run. State that this validates integration
   rather than strength; a fresh routine match comparison is not required.
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
   instead of another full match.

   Before the evidence audit, record authorized edits to historically mapped
   Markdown guides with `go run ./tools/research map-docs --file PATH.md --record NEW_TRANSITION.json`.
   Repeat `--file` for each changed mapped guide. Original hashes and raw evidence
   stay intact; do not regenerate old inventories to hide changed observations.

   For example, while preparing v0.23:

   ```sh
   go run ./tools/benchmark --validate --race --version v0.23-release-smoke \
     --previous tools/benchmark/bin/pandafish-v0.22 --previous-version v0.22 \
     --sequential --stage previous --pairs 1 --results /tmp/pandafish-release-smoke
   go run ./tools/benchmark --candidate tools/benchmark/bin/pandafish-candidate \
     --version v0.23-release-smoke --stage external --parallel \
     --stockfish-skill 3 \
     --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
     --pairs 1 --results /tmp/pandafish-release-smoke
   git diff --check
   ```

   Also smoke-test the full default concurrent schedule with `--pairs 1`.
   Its two-game previous screen may return a blocked decision; inspect all six
   saved games for integration failures instead of rerunning for a better score.

   These explicit stages check all three integrations without treating a
   two-game smoke score as a new strength gate. They supplement the already
   reviewed candidate comparisons (100 games per opponent through v0.13;
   50 per opponent afterward).

   Also perform the local lichess-bot wrapper check when using that client.
   These checks do not start a live bot or play online games. If engine behavior
   changes after testing, repeat the routine comparison for the new candidate.
5. Inspect the staged files, commit, push, and open a PR describing the change
   and checks performed. Review and merge the PR into `main` before tagging.

## Tag and publish after merge

Confirm the PR is merged and the version is unused, then fast-forward local
`main` to the merged commit. For v0.23:

```sh
git switch main
git fetch origin
git merge --ff-only origin/main
git status --short
```

Check that the working tree is clean and HEAD is the intended release commit.
Then create and push the tag:

```sh
git tag -a v0.23 -m "v0.23: reuse search state and avoid temporary tag copies"
git push origin v0.23
```

Publish a GitHub release for that tag, through the GitHub UI or
`gh release create`. Use the version alone as the release title, such as
`v0.23`. Describe the engine or tooling change, remaining limits,
and validation, linking the candidate benchmark evidence. An engine release
does not automatically advance the Stockfish skill target.

## Benchmark evidence in GitHub releases

For every release with strength results, include a compact table in its GitHub
release description, plus links to the complete settings, identities, hashes,
and PGNs. Copy values from the saved Go summary, identify the tested candidate,
and state the games per opponent. Use this column layout and matching separator:

```markdown
| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Previous version | ... | ... | ... | ... |
```

Completed score excludes unfinished games. If any games are unfinished, also
state the possible score when they count as losses through wins. Keep timing
and uncertainty caveats relevant to the result. A tooling-only release should
say when no new strength benchmark was run, link its existing baseline, and
describe smoke checks separately. If benchmarks were recorded after an older
release, date and label them as retrospective evidence. Preserve the original
development labels and hashes; do not imply that they tested the tagged binary.

Check the table's column counts and rendered Markdown before publishing.
Put summaries in release notes and link the detailed records in the repository;
benchmark files do not need to be attached again as release assets.

## Build the next baseline

Build from the tag rather than a development checkout. For v0.23:

```sh
mkdir -p tools/benchmark/bin
benchmark_bin="$PWD/tools/benchmark/bin"
build_dir=$(mktemp -d)
git archive v0.23 | tar -x -C "$build_dir"
(cd "$build_dir" && go build -o "$benchmark_bin/pandafish-v0.23" .)
```

The tagged executable becomes the previous version for v0.24 development.
Compiled binaries remain gitignored. Building this baseline does not require
relabeling or rerunning the development benchmarks. A separate benchmark of
the exact released binary is optional. Keep private bot greetings in sync
when choosing to run the new engine through lichess-bot.


## Tooling releases

v0.8 established the precedent for a minor release that adds developer tooling
without changing playing strength. v0.14 follows it for the faster benchmark
workflow. State the unchanged engine policy explicitly, distinguish workflow
calibration from new engine-strength evidence, and link the original records.


## Research releases

v0.14.1 through v0.14.4 use patch versions to publish an unpromoted experiment as research.
Document its hypothesis, code, controls, validation, results, limits, decision,
and concrete next direction. Preserve raw records and original candidate labels.
Archive experimental source/tests as data, with a reproducible restoration path,
and restore the released playing policy before shipping a research-only release.
Do not silently include failed experimental search in the executable or claim
its measurements as a new strength gain. Verify both the archived replay and
the released engine. Version-only smoke games check integration, not strength;
keep their evidence separate from the original experiment.


v0.15.1 applies the research-patch convention to a completed live investigation.
Preserve original session data, analysis methods, confirmed positions, negative
probes, uncertainty and a fresh-context handoff. A separately tested reporting
repair can accompany it while engine playing policy stays unchanged. State
that no new strength benchmark was required; distinguish original live engine
identity, current tool version and integration smoke checks.


v0.15.2 preserves the rejected narrow piece-safety experiment as research data.
The term fails its mechanism/cost assessment before strength matches, so the
release makes no match or Elo claim. Playing policy remains v0.15.1 (and v0.15).
Keep the original candidate label and all diagnostic failures; distinguish the
six-game release integration smoke from a strength screen. Reassess the next
hypothesis from the result, rather than following the old ranking as a checklist.


v0.15.3 archives the unpromoted bishop-pair screen and a larger architecture
investigation. Archive-only treatment and unchanged play were explicitly
selected. Keep all diagnostic prototypes and intentional/real failures
outside the normal build. Distinguish original bishop-pair results, retrospective
v0.15.2 controls, architecture position probes and release smoke games. No
prototype efficiency number is an Elo claim, and no optional confirmation
runs automatically. Future work starts from the tagged research release and
follows its measured larger-project plan rather than a fixed feature checklist.

## Source and publication rules

All repository engine, tool, test and research source must be Go, including
archived source and subprocess fixtures. Follow [AGENTS.md](AGENTS.md).
Use neutral, factual engineering prose without first-person voice or personal
framing. Never use the word "maintainer" in PR/release titles or descriptions.
Before publishing, run `go run ./tools/research wording --file PATH` on both
prepared bodies and check the titles too.

v0.19 publishes the [completed evaluation investigation](docs/RESEARCH_V019.md)
and Go-only migration with unchanged v0.18 playing policy. Original scripts and
unedited evidence remain in the hashed release attachment. Raw results stay
in Git, with [Go research commands](tools/research/README.md) for
supported replay/verification. Retired collectors/trainers are documented as
retired; verifying their frozen output does not reproduce fitting.
No new strength screen is required; release smoke games remain integration
evidence. Preserve every smoke failure; do not rerun for a favorable score.

v0.20 promotes the [search-legality kernel](tools/benchmark/research/v0.20-search-legality/README.md)
after its standard 50-game concurrent screens. The rejected quiet-mate and
transition-only candidates remain archived with their failures. Preserve all
development labels, hashes and games; version-only integration checks do not
replace or relabel that strength screen. The later tagged Skill 2–10 benchmark
is a separate reference-level assessment, not an automatic target advancement.

After reviewing that complete ladder, Skill 3 becomes the benchmark default in
a separate tooling/data PR without a release. The clean Skill 2/3 rows support
the reviewed choice; affected upper-level rows and all reporting failures remain
preserved. Future routine screens use Skill 3, while historical release evidence
keeps its recorded skill/options. At that adoption, engine playing behavior
remains tagged v0.20.

v0.21 promotes [precomputed stable ordering](tools/benchmark/research/v0.21-rank-order/README.md)
after its clean 55%/50%/71% standard screen. Original priorities, stable ties,
values, search/TT policy and clocks are unchanged. Preserve the development
binary identity and qualify the small direct lead; version-only smoke checks
are integration evidence. The current Stockfish reference remains Skill 3.

v0.22 promotes the [revised exact kernel](tools/benchmark/research/v0.22-bitboard-exchange/RESULTS.md)
after its clean 64%/71%/70% standard screen. Deferred generation, exact evasion
prefiltering and bitset exchange computation preserve fixed-work search results.
The public 100ms teacher-loss gate remains failed; its explicit pre-game review
exception and all earlier failed variants stay archived. Preserve the original
development label/hash. Version-only release smokes are integration evidence,
not another strength screen. Skill 3 remains the standard. No further experiment
starts as part of this release.


v0.23 promotes the [exact state kernel](tools/benchmark/research/v0.23-state-kernel/README.md)
for a 14.142% measured same-work saving. The original screen scores 51% direct,
69.4% over 49 completed Skill 3 games (one 400-ply cap; 68–70% possible overall),
and 76% against Sacrifice, with zero forfeits and 150 legal records. Playing-
strength improvement remains inconclusive. Preserve all 23 unsuccessful attempts,
the original development label/hash and every failed control. The
[Go workflow and flowchart](docs/RESEARCH_TOOLING.md) and efficiency observations
are part of the release. Version-only smokes are integration evidence. Skill 3
stays the standard; no next candidate or live session starts automatically.
