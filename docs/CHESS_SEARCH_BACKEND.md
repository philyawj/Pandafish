# The local chess-library search backend introduced in v0.17

## Current state at v1.0

v1.0 adds a bounded protected queen-contact query and actual-defense verification
using the existing legal-move and child-state APIs. Original general rules APIs
remain intact. See the [engine guide](ENGINE.md) and
[implementation/evidence](../tools/benchmark/research/v0.24-contact-threat/README.md).
Active development is paused at the [v1.0 milestone](../README.md#milestones-and-plans).

The preceding v0.23 state is summarized below; the original v0.17 account remains
historical context rather than a current implementation specification.

The local pinned module remains the rules backend. Additive APIs now include
staged tactical/quiet generation, evasion prefiltering, virtual ordinary-move
tags and worker-owned per-ply child storage. Child updates copy board state into
reusable cells; they no longer require a new Position/Board allocation per edge.
This is still a copy-on-ply design, not make/unmake. General immutable APIs and
the original upstream runtime identities remain available.

Root `go test ./...` excludes this nested module; run its tests/races/vet
explicitly or use the current shared validation plan. Frozen independent vectors
need the optional oracle helper's explicit input/output. The
[post-v0.23 assessment](NEXT_ENGINE_ASSESSMENT_V023.md) originally ranked shared
legality/attack context first. The [completed investigation](RESEARCH_V024.md)
records its outcome and later candidates; it is not an active development queue.

The remainder records the v0.17 design and subsequent v0.18 follow-up. Its
allocation descriptions and future-work decisions describe those versions.

## Historical v0.17 introduction

I introduced this change in [PR #32](https://github.com/philyawj/Pandafish/pull/32)
and released it as [v0.17](https://github.com/philyawj/Pandafish/releases/tag/v0.17).
The PR's reference to **v0.16** describes the behavior before the change.
**v0.17** is the release that adds the local replacement and new search API.

## Why I made the change

Pandafish uses corentings/chess for board state, legal moves and game rules.
The library already uses bitboards and a mailbox. I did not need to introduce
bitboards to obtain this improvement.

Its general-purpose move API serves a different workload from an engine's
innermost search loop. A game client benefits from a complete cached legal-move
list and defensive copies that callers can safely change. Capture quiescence
usually needs a much smaller set: captures and promotions, with all legal
escapes when checked.

In v0.16, ordinary capture search still generated and tagged the complete legal
list before collecting its tactical subset. v0.16 already sorted only that
subset outside check; the remaining waste was generating and validating quiet
moves in the first place. Other search paths also allocated legal-list copies.
Testing whether *any* legal move existed could allocate a small list too.

The v0.16 profile put approximately 46% of sampled CPU under legal generation.
The timed kernel allocated about 71.1 MB and 154,761 objects per request in that
workload. These numbers made a focused search API a concrete engineering target.
They were not evidence that every part of the library needed replacing.

## What “local replacement” means

The root [go.mod](../go.mod) still names the original module/version, and adds:

```go
replace github.com/corentings/chess/v2 => ./third_party/chess
```

For this repository's builds, Go resolves that module's imports to the checked-in
local directory. This affects the engine and the developer tools that import
that module. It is not a second chess type system or a wrapper converting
positions between two libraries.

I copied the pinned **v2.6.1 root-package runtime**, its module metadata and MIT
license into [third_party/chess](../third_party/chess). The original runtime
files were preserved byte-for-byte. The unused rendering subpackages and
upstream test fixtures were not copied into the runtime directory.

The v0.17 addition is [the search API](../third_party/chess/search_moves.go),
plus [Pandafish's buffer/ordering integration](../search_moves.go). The existing
board, position update, hash, FEN/PGN and game-rule methods were retained.
This is an additive local fork of the library, not a new board implementation.

The original file hashes are recorded in
[upstream-source.json](../tools/benchmark/research/v0.17-search-moves/upstream-source.json).
The MIT license remains in the local directory and in the
[third-party notices](../THIRD_PARTY_NOTICES.md). A narrow Git attribute keeps
the original license's final blank line intact. Go's original module checksum
does not verify a local replacement; Git, these source hashes and the tests
provide that provenance.

## The new API

| Method | Purpose |
| --- | --- |
| `Position.AppendSearchMoves(dst, tacticalOnly)` | Append legal moves into caller-owned storage. In tactical mode, mask quiet destinations before expensive tagging/legality work; retain captures, en passant and every promotion, including quiet underpromotions. |
| `Position.HasLegalMove()` | Stop as soon as a legal move is found, without constructing a move list. |
| `Position.SearchStatus()` | Preserve mate/stalemate status using the existence check. |
| `searcher.searchMoves(...)` | Choose full or tactical generation, reuse the current ply's slice and apply the existing stable move ordering. |

The new generator uses the existing bitboards, attack helpers and move-tagging
logic. It visits squares in the original order. Full generation retains the
original legal castling routine. If a position already has a complete cached
legal list, the API can read it and append copies into the caller's storage.
It never sorts or edits the library's cached list.

An empty tactical list does **not** mean stalemate. The separate legal-move
existence check preserves that distinction. Main search requests all legal
moves, and checked quiescence requests every evasion, including quiet blocks
and king moves. The tactical-only request is used outside check.

## How the storage works

Each search worker owns one reusable move slice per active search ply. A
parent keeps iterating its own slice while its child uses the next ply's slice.
A later PVS re-search can reuse its slice after the earlier call returns.
Descendants cannot overwrite an ancestor's move list.

The slices grow as needed and are reused within the UCI search request. They
are not persisted to disk or shared globally. The separate score/hint tables
retain their existing game/search lifetimes.

Position updates remain immutable library updates: creating a child still
creates a new position and board copy. I did not introduce make/unmake, an undo
stack, new hash rules or a different game-state representation in v0.17.

## What this changes about play

At a fixed completed depth, this backend was required to return the same nodes,
score and chosen move as v0.16. All 192 paired fixed-depth probes met that
requirement. Stable ordering and legal search eligibility were preserved.

At equal time, the lower work cost lets the engine finish more search. That
can change the selected move and game outcome. It is the intended route from
an engineering improvement to a possible strength gain.

Material values, placement evaluation, PVS/root policy, capture filtering,
quiescence horizons, table policy, repetition/mate behavior and time management
were unchanged in the released v0.17 experiment.

## What the measurements showed

| Measurement | v0.16 baseline | v0.17 candidate |
| --- | ---: | ---: |
| Depth 4, 84 positions, identical 1,010,100 nodes | 4.891 s | 2.779 s |
| Depth 6, 12 positions, identical 3,658,933 nodes | 15.795 s | 9.579 s |
| Timed-kernel allocated bytes per request | 71.1 MB | 9.54 MB |
| Timed-kernel allocations per request | 154,761 | 73,499 |

The fixed-work cohorts were about 1.76x and 1.65x faster. Allocation per request
fell about 87%; that is cumulative allocation, not peak resident memory.
These are selected workload measurements, with timing variance and different
cohorts at the two depths. They are not universal performance bounds.

Native equal-time tests also improved completed depth overall. One of 504
three-search paired requests finished shallower and remains in the report;
its cause was not isolated. I do not assume every position gets faster.

The standard concurrent screen scored 59% against v0.16 and 53% each against
Stockfish Skill 2 and pinned Sacrifice, with no forfeits or unfinished games.
I approved promotion under the practical gate. The repeated openings, small
samples and direct 44% White / 74% Black split limit the strength inference.
A 53% sample does not establish a reliably above-50% rate or justify advancing
the Stockfish target. The six-game release smoke's two-game direct gate failed
at 50% draws; that integration-only outcome was retained separately.

[The full report](../tools/benchmark/research/v0.17-search-moves/README.md)
separates correctness, efficiency, game results and uncertainty, and links all
raw records and preserved failures.

## Why I considered this safer than a full rewrite

I could not avoid the full generation cost merely by calling
`ValidMovesUnsafe()`: that avoids a defensive copy, but its first call still
builds the complete legal list. A wrapper using only the public API has the
same limitation.

Adding the API inside the pinned package let me use its existing internal
bitboards and legality machinery. A new internal board would also require new
state updates, undo, special-move legality and hash maintenance. The additive
route delivered the measured target with a smaller correctness boundary.

Validation included the original upstream tests under the race detector,
500 seeded comparisons with unchanged library methods, special-move controls,
starting-position depth-four perft, 596 independent Python-chess move/tag,
child-FEN/hash/status and depth-two-perft cases, fixed-depth search equivalence,
UCI/wrapper tests and independent replay of all 150 screen games.

## The maintenance cost I took on

The local copy does not automatically receive upstream fixes. Updating it needs
an explicit version/source review and the same rule/state/search controls.
The import path staying the same does not make the local source identical to a
future upstream release.

The root `go test ./...` does not automatically run a separate nested module's
test suite. Pandafish's own tests exercise the new API, and the experiment also
replayed the pinned upstream tests/fixtures in a scratch module. That explicit
upstream replay and independent move/state oracles should remain part of future
library changes. The original runtime copy deliberately retains upstream
formatting; authored search API files are checked separately.

The first v0.18 prototype illustrates why the new additions need their own
oracles too: independent feature tests caught a pawn-mask initializer running
before an upstream lookup table's `init()` function. I corrected it before
calibration/activation and retained the failed vectors. That was an unreleased
addition, not a defect in the v0.17 search backend. The unchanged upstream tests
alone would not have covered it.

I need to keep three categories clear: original upstream files, local additive
APIs, and the engine code that consumes them. New work should make its boundary
and validation visible rather than silently editing unrelated game rules.

## How I intend to use it next

I plan to keep using this extension boundary for measured improvements. The
first v0.18 evaluation experiment tested a compact board-feature summary using
the existing bitboards, with no per-piece move-list allocation. I preserved that
unpromoted model and its failed short-budget gates as research data. The v0.18
follow-up isolates a value copy of piece bitboards and precomputes the original
material/placement formulas with exactly the same scores, without learned terms.
I approved promotion of that exact-score kernel after its standard screen; the
fitted values remain archived. The [v0.18 report](../tools/benchmark/research/v0.18-exact-evaluation/README.md)
separates the functional speedup, game results and limits. This is a separate
addition to the v0.17 backend, not a retroactive change to that release.

Other possibilities include sharing legal-status/attack information, reducing
repeated move-tag work, or improving caller-owned storage. Each needs evidence
that it removes meaningful cost and preserves rules. I have not committed to
implementing all of them.

An internal make/unmake position remains an option if profiling later justifies
its extra complexity. The v0.17 result met the bounded backend target, so I
have investigated the diagnosed evaluation/threat gaps. The first fitted model
failed its short-budget decision gate; I now recommend evaluating the exact
base-score kernel separately and revising the learned model on new data. I
completed that bounded kernel revision and approved it for v0.18; the fitted
model remains unpromoted.
If the added API
stabilizes, contributing it upstream could also reduce long-term divergence;
that would be a separate decision.
