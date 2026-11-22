# Repository audit after v0.23

Assessment date: 2026-11-21. Baseline: `28948e7473251d950003ec1f966f0cf26369b817`.

The most useful repository work is to separate current instructions from frozen
research, complete the shared validation boundary, and stop requiring a full
evidence collection in every checkout. The engine and three Go tools already
have sensible ownership. Moving their source into a new framework would add
migration work without addressing the largest costs.

The separate [ranked engine assessment](NEXT_ENGINE_ASSESSMENT_V023.md) recommends
a shared legality/attack-context investigation first. Repository cleanup and
engine candidates are separate review decisions.

## Scope and verified starting state

The local handoff requests an audit and 5–10 ranked alternatives,
superseding autonomous candidate implementation. This agrees with the explicit
scope exception in [the research workflow](RESEARCH_WORKFLOW.md). No prototype,
strength screen, training campaign, live game, commit or publication belongs to
this assessment. New reports and bounded unchanged-baseline diagnostics are in
scope. No historical files were changed or removed.

Local `main`, `origin/main` and peeled v0.23 all match the baseline. The GitHub
connector independently confirms current remote `main`, merged
[PR #41](https://github.com/philyawj/Pandafish/pull/41), and the normal, published
[v0.23 release](https://github.com/philyawj/Pandafish/releases/tag/v0.23).
[Saved remote fields](../tools/benchmark/research/post-v0.23-audit/remote-state.json)
record these checks. Initially only `handoff.md` was untracked. The available
tagged executable's SHA-256 matches the handoff's `88d78500…6be`.

Two handoff details need qualification:

- The release PR **changes 1,442 files**, adding 1,424 files containing
  312,973,122 bytes. “Adds 1,442” confuses changed and added files.
  Comparing v0.22 with v0.23 also includes the intervening live-review PR:
  1,984 changed / 1,965 added files, 510,516,565 added bytes.
- [The current checkpoint](RESEARCH_CHECKPOINT.md) opens with the old research
  branch, v0.22 baseline and pre-commit instruction, then ends with publication
  authorization. It is a frozen pre-release checkpoint presented as current.
  Historical “stop before commit” notes do not reopen the completed release.

## Evidence and limits

[Audit evidence and reproduction](../tools/benchmark/research/post-v0.23-audit/README.md)
contain the Go diagnostic, size/hash summaries, local-path findings, CPU sample,
exactness comparison and original-inventory audit. Counts describe tracked
baseline files, excluding ignored executables, local sessions and `.git`.
No clone-transfer benchmark, GitHub pack-size measurement or model-token-saving
measurement was performed.

Root `go test ./...`, explicit nested chess-module tests, and vet in both modules
pass. Some results use Go's existing test cache. The local evidence audit verifies
56 inventories / 3,710 entries; the external v0.19 attachment remains unverified
in this assessment. A passed inventory audit proves the listed local bytes, not
that every live/result file has an inventory or that an external asset is retrievable.
No new full race/release suite was necessary for report-only changes; v0.23's
original stronger validation remains separate evidence.

One unchanged-source warm CPU profile completes the existing 20-game,
1,145-root, depth-five workload. Shared `warm-report` verifies all moves, scores,
depths and 43,118,858 nodes against the frozen release candidate. There are no
per-root memory-stat reads. Sampling includes harness/runtime work; its elapsed
ratio is diagnostic, not a new throughput comparison. The installed Go 1.25
toolchain lacks the `pprof` executable; the installed Go 1.24.5 standalone
analyzer successfully reads the Go 1.25 profile. This tooling failure is retained
in the method. AC power was reported; system background work was present and
thermal conditions were uncontrolled.

## Layout, ownership and source boundaries

| Area | Actual responsibility | Assessment |
| --- | --- | --- |
| Root Go package | UCI, time allocation, search, evaluation, tables, recording | Small enough to navigate directly; no engine package split is justified by this audit. |
| `third_party/chess` | Pinned rules runtime and additive search APIs; separate module | Useful extension boundary; preserve original runtime identities and explicit nested checks. |
| `internal/telemetry` | Shared recording schema | Appropriate existing shared package. |
| `tools/benchmark` | Local matches, progress, validation, root Markdown report | Keep as the sole match runner. Defaults agree with current 25 pairs, Skill 3, 100ms + 50ms grace. |
| `tools/lichess` | Session recording and offline processing | Keep separate from strength experiments. No live client was started. |
| `tools/research` | Plans, runner, decisions, snapshots, replay, inventories | Main reusable research entry point; experimental fit/target commands have narrower historical contracts. |
| Research/live/results trees | Frozen source, observations, reports and failure history | Main checkout and navigation cost; need selective retrieval rather than another source reorganization. |

There are five root-module packages and one nested module. The tracked tree has
132 `.go` files, 299 archived `.go.txt` files and 372 Markdown files. Archived
Go source is evidence rather than an active program inventory. The old drivers
must not be copied as starting templates or restored as ordinary module packages.
The three ignored generator packages that blocked release formatting are already
retired; the current package list contains only the intended five packages.

The source policy is useful but not a full build-isolation rule: it skips local
`tmp`, `temp`, `bin` and hidden directories, while Go's package discovery can
still discover ordinary Go packages in some ignored directories. `.gitignore`
does not define the module boundary. Keep scratch research outside the module;
verify `go list ./...` before release rather than adding another blacklist.

## Ranked repository remedies

These are the original audit findings. The subsequent
[storage decision](EVIDENCE_STORAGE.md) defers the archive pilot and checkout-size
target. The cleanup validation record identifies the implemented guidance,
validation and link remedies; larger format/framework migrations remain deferred.

### 1. Make the current entry point reliable

**Known problem:** AGENTS.md sends new work to a checkpoint that still starts
from v0.22. The tooling guide's baseline examples also use v0.22. The 1,180-line
benchmark guide interleaves current commands with historical “next” directions,
including a paragraph that still says the old evaluation experiment is held
locally. The backend note has first-person v0.17 descriptions and says child
updates remain immutable allocations, without a current v0.23 overview.

**Proposed remedy:** replace the current checkpoint with a short dated pointer:
published v0.23, current assessment phase, no active candidate, baseline identity,
report links and next approved command when one exists. Preserve its previous
contents in history or a named frozen copy. Add a short current-state preface to
the backend note; retain the v0.17 narrative as historical. Put a current quick
start before historical benchmark milestones, and label historical commands.

**Benefit/cost:** high reduction in rereading and wrong-baseline work; small
documentation effort, no compute. **Risk/path:** mapped historical Markdown
requires the existing `map-docs` transition before the evidence audit. Do not
bulk rewrite frozen reports or silently retitle their old decisions. The new
assessment documents need no legacy remapping. These edits are proposed, not
performed here.

### 2. Put the real correctness boundary in the shared validation path

**Known problem:** [validation-plan](../tools/research/validation_plan.go) emits
root normal/race/vet and builds only. Root `./...` excludes the nested module.
[Benchmark validation](../tools/benchmark/checks.go) has the same nested-module
gap. The final v0.23 plan manually adds nested normal/race/vet and 596-vector
checks. The tooling map describes a shared fail-closed protocol, but a fresh
agent still needs to know about these manual additions.

**Proposed remedy:** extend the existing validation-plan with explicit chess
normal/race/vet steps and a documented opt-in backend-oracle stage with required
inputs/output. Check authored additive files' formatting separately from pinned
upstream formatting. Reuse this contract from benchmark validation or document
which entry point is complete; do not create a second runner. An optional small
Go CI workflow could call the same checks; no `.github` workflow currently exists.

**Benefit/cost:** prevents false confidence and repeated manual plan editing;
modest implementation, bounded local test time. **Risk/path:** retain existing
frozen plans; test failure sequencing and nested-test failure propagation before
using the revised builder. The independent vectors are frozen observations,
not a newly independent rules implementation.

### 3. Shrink the default checkout through verified cold evidence

This is the largest storage opportunity. The staged proposal below targets an
eventual **under-200-MB ordinary checkout**, with source, current controls and
compact evidence available offline. This is a proposed target, not a measured
post-migration result. Cold data remains fully recoverable.

### 4. Record consumed cohorts once

**Known problem:** `generate` automatically excludes old v0.19 witnesses, but
later public/newer/fresh/live exclusions remain caller-supplied flags. The guide
names older v0.23 cohorts and the checkpoint mentions later consumed seeds;
there is no single authoritative consumed-cohort list. A new seed alone does
not prove disjointness. `native-plan` uses established public/newer paths and
does not automatically select a new untouched final cohort.

**Proposed remedy:** a small versioned JSON cohort ledger listing input hashes,
seed/stream, role, exclusions and first inspection status. Let the existing
generator consume that ledger and preserve its resolved inputs in each method.
Do not introduce a database or collector framework.

**Benefit/cost:** fewer missed exclusions and long copied argument lists; modest
Go work. **Risk/path:** compare the resolved exclusion set with the last successful
state-kernel plan, then add all subsequent consumed inputs. Existing inspected
data remains development data. Never retrospectively relabel it untouched.

### 5. Repair navigation without altering archival meaning

The inline-Markdown path check finds five missing paths. One is actionable:
the [quiet-check report](../tools/benchmark/research/v0.23-quiet-checks/README.md)
links to a nonexistent matchup-level README; the actual report is
[at the screen root](../tools/benchmark/research/v0.23-quiet-checks/screen/README.md).
Four links in the exact copied pre-amendment candidate README resolve relative
to its old location; that frozen copy should retain its bytes and gain context
from a surrounding index. The checker does not validate URL availability,
reference-style links or heading fragments.

Add one compact attempt index with mechanism, outcome, decisive gate, source,
raw-data location and replay contract. Link it from the current quick start;
avoid maintaining several competing active rankings. Clarify that legacy
`policy` defaults intentionally describe v0.18→v0.19; they are not a current
release parity check. Keep historical commands explicit or require arguments
when ambiguity is otherwise likely. **Validation:** local links, tables,
wording, inventory/map transitions and a fresh-context dry read.

### 6. Consolidate genuine duplicated source when it next changes

`tools/benchmark/uci.go` and `tools/research/uci.go` are currently byte-identical
184-line process/handshake implementations. Sharing a narrow internal UCI
transport would avoid diverging timeout/cleanup fixes. Different score/probe
contracts should stay with their callers. This is a modest correctness/time
benefit, negligible disk saving and lower priority than validation and storage.
Move it during a relevant fix, with both tools' process-failure tests, rather
than making a general engine-client framework now. Historical driver duplication
is deliberately retained evidence and does not require active consolidation.

## Repository size and a concrete reduction plan

The tracked checkout contains **1,374,998,619 bytes: 1.375 GB / about 1.28 GiB**.
That is not the compressed Git repository size. Local object storage reports
184.26 MiB loose plus 2.44 MiB packed, roughly 187 MiB; this includes local
history and is not an exact remote clone-transfer size. The largest tracked
file is 25.9 MB. There is no demonstrated Git failure at this size.

GitHub's current repository guide recommends a maximum 10 GB compressed `.git`
size and keeping generated files outside Git. Pandafish is well below that
guideline; working-tree bulk and repeated evidence growth are the immediate
problems. [GitHub repository limits](https://docs.github.com/en/repositories/creating-and-managing-repositories/repository-limits)

| Tracked area | Files | Decimal MB | Independent gzip sizing, MB |
| --- | ---: | ---: | ---: |
| Research | 3,685 | 666.85 | 76.64 |
| Live evidence | 2,017 | 635.35 | 85.99 |
| Controlled results | 290 | 70.43 | 7.30 |

The Go sizing diagnostic compresses each file to a counting writer at gzip's
best-speed level, without modifying or storing compressed copies. Total evidence
falls from 1,372.62 MB to **169.93 MB**. This measures archive potential, not
expected improvement in Git's already compressed/delta-encoded objects.

Exact duplicates of at least 1 KiB offer only **2.24 MB** redundant bytes.
Source snapshots occupy 18.27 MB; their decoded payloads contain 16.27 MB of
copies versus 3.76 MB distinct content. Snapshot deduplication cannot explain
most of the repository. JSON whitespace accounts for 204.15 MB of the evidence,
but minifying frozen files would change their hashes while Git already compresses
whitespace well. None of these is the best first migration.

### A. First pilot: the oldest large live session

The v0.15 live session alone occupies **362.53 MB**. Package its original raw
files and full analysis in a content-hashed archive on release-asset or durable
object storage; keep its Markdown review, compact summaries, manifest and useful
PGNs in Git. Preserve original paths inside the archive. This offers up to
362.53 MB less checkout data before subtracting the small retained set. Extending
the policy to live evidence has a 635.35 MB gross ceiling; actual savings must
be reported from the approved file list.

No active Go test found in this audit directly loads a live archive. Offline
analysis and future cohort exclusion commands do consume those paths explicitly,
so they need a deliberate fetch/verify step or documented archive-cache paths.
Do not claim that an inventory audit can already verify arbitrary externalized
evidence: its current migration support covers the historical v0.19 attachment,
not general remote data.

**Pilot acceptance:** freeze the file list and original SHA-256/byte identities;
create an archive with a per-file manifest; verify an independent download and
safe extraction into a fresh scratch directory; reproduce the existing report
without engine queries; check every retained link and relevant exclusion input.
Only then remove the bulk files from the next tree through an explicit migration
record and replacement index. Keep a second archive copy. Missing external data
must report “not checked,” not pass. No source, raw failure or game is discarded.

### B. Stop multiplying new evidence

For future runs, retain one canonical raw journal plus compact summaries. The
73 `native.json` files occupy 44.92 MB while their 73 journals occupy 29.53 MB;
the code writes a sorted full JSON copy after the journal. Quality/checkpoint
files occupy another 10.57 MB. These are format opportunities, not permission
to delete old records: prove case/trial equality and deterministic reconstruction
before changing future readers. Keep interruption checkpoints until finalization,
then package them with the frozen run rather than copying every representation
into the main tree. Reference caches can point to hash-identified original
observations instead of embedding complete originals repeatedly.

Shared Go readers must continue to handle the existing schema. Regression tests
should compare normalized reconstructed rows and quality decisions, preserve
interrupted/failing runs, and require no rerun of engines. This is a narrow
evidence-format change, with meaningful future growth savings and modest code
cost. It does not require a new research framework.

### C. Separate test fixtures from full research archives

Root tests currently load v0.17 oracle inputs and v0.19 fit/leaves data.
Research tests load three historical corpora and original witness paths. The
generator also loads `blocked.json`, `unlabelled-roots.json` and
`forest-leaves.json`. Blindly externalizing all old research would break normal
tests or silently weaken exclusion coverage.

Extract compact, traceable Go-test fixtures into `testdata` only after recording
source hashes and a deterministic Go extraction method; retain the full original
archives separately. Preserve equivalent meaningful coverage, especially all
canonical-key and witness-history controls. Corpus-wide verification remains
an explicit evidence check after retrieval. Then move cold research bodies
behind the same verified archive mechanism. Research plus live evidence accounts
for over 1.30 GB of the checkout, so this route can plausibly reach the proposed
under-200-MB target. The exact retained dependency set is still open.

### D. Alternatives and their tradeoffs

| Option | Benefit | Limitation / recommendation |
| --- | --- | --- |
| Gzip cold files inside Git | Sizing suggests evidence can fit near 170 MB in the current tree; lossless original bytes can be verified after decompression. | Needs reader/mapping changes, reduces direct browsing, and old blobs remain in history. Packed clone size may not improve. Prefer archive assets for cold bulk. |
| Separate evidence repository | Independent checkout and growth policy; source repo stays small. | Moves storage elsewhere and adds cross-repository version/retrieval work. Useful if ordinary evidence browsing/history is essential. Preserve immutable IDs. |
| Git LFS | Keeps large payloads outside normal Git objects. | Normal LFS checkout still downloads payloads unless configured otherwise; storage/bandwidth policy and a client become dependencies. It also needs migration to remove old Git blobs. Lower fit for these cold, text-heavy runs. |
| Sparse checkout plus partial clone | Immediate optional local footprint relief without changing evidence. | Missing corpus dependencies must be selected; sparse checkout alone does not avoid object download. Full audits need hydration. Document an exact path set before advertising it. |
| Delete bulk from latest tree only | Smaller ordinary checkout, original history intact. | Full-history clones still retain old blobs; must first archive and verify evidence. This is the preferred first migration shape. |
| Rewrite all history | Can reduce future full-clone history size. | Breaks commit/tag identities, provenance and base-tag replay unless all references and archives are deliberately migrated. Keep as a separate last-resort decision; not justified by current ~187-MiB local object storage. |

Sparse checkout and blob-filtered clone have distinct roles in
[Git's sparse-checkout documentation](https://git-scm.com/docs/git-sparse-checkout)
and [clone documentation](https://git-scm.com/docs/git-clone).
Adding `.gitignore`, moving data to another tracked directory, or running local
garbage collection cannot remove evidence still reachable from history.

## Fresh-agent reproduction contract

A fresh clone can build the engine with Go 1.25. The fixed sandbox cache is
`GOCACHE=/tmp/pandafish-v016-go-cache`. Building the three existing tools once
avoids repeated `go run` compilation. Full local checks require both modules;
the independent oracle and warm workload require their explicit environment
variables because normal tests skip those optional helpers.

Archive replay reconstructs Go source from a base tag plus declared overlay;
it does not reproduce retired data collection or training. The exact v0.23
candidate overlay has 118 source/module files. The optional external wrapper
and pinned opponent binaries/assets are additional dependencies, not a reason
to start a live client. After storage migration, normal source/test use should
remain offline, while deep historical verification should say exactly which
archive must be fetched. Test that contract in a fresh checkout before claiming
the migration complete.

The recommended repository decision is a small current-guidance/validation fix
and a separately reviewed cold-live-evidence pilot. Keep all three existing
tools, all failures and all original release identities. Substantial cleanup
and engine implementation remain unstarted pending review.
