# Methods and scope

I selected the candidate after profiling tagged v0.16 and reassessing both the
backend and evaluation/threat projects in [ASSESSMENT.md](ASSESSMENT.md).
The release is published at `4ae65b14f69d95a1ad498baa759fec5e6bd0a1a6`.
The original v0.16 screen and deferred confirmation are unchanged.

## Additive library API and storage

`third_party/chess` contains the pinned v2.6.1 root-package runtime and license,
byte-preserved, plus `search_moves.go`. Unused rendering subpackages and the
upstream fixtures/test source are omitted from the runtime copy. No original
runtime method, board representation, state update, hash or parser is edited.
`upstream-source.json` verifies every copied original file. The root `go.mod`
uses a local replacement for this specific library. No new dependency version,
third-party engine code or neural weight is introduced.

The new generator visits existing bitboards in ascending square order, matching
the original move order. It masks non-tactical destinations before tagging and
legality testing, while keeping pawn promotions (all four types), en passant
and checking captures. Full generation includes the original legal castling
routine. Cached full legal lists are read-only inputs; caller-owned output
slices receive copies. Empty tactical lists are not used to infer stalemate.
`HasLegalMove` tests existence without constructing a list; `SearchStatus`
retains exactly the library's mate/stalemate semantics. Other draw rules keep
the released search treatment.

A searcher owns one reusable move slice per active ply. Descendants/re-searches
reuse their own storage and cannot overwrite an ancestor's loop. Capacity grows
on demand, and the searcher dies after the UCI request. This is not an internal
make/unmake rewrite. Position updates remain immutable library updates.
Main search requests all moves; checked quiescence requests all evasions;
nonchecked quiescence requests legal captures and all promotions only. Ordering,
SEE eligibility, search/PVS depth, score/hint tables, killers, repetition,
evaluation/material and time management remain unchanged.

## Controls and independent validation

- Existing repository tests, full race suite, vet, native/ARM64 build and offline
  public-config lichess-bot wrapper run before timing probes.
- New tests compare fresh/cached full and tactical move order, tags, status,
  child FEN/hash, immutability and per-ply ownership against the unchanged
  library methods over fixtures and 500 seeded legal positions. Fixtures cover
  castling, legal/pinned EP, all promotions, quiet evasions, mate and stalemate.
  Starting-position depth-four perft is 197,281.
- The original upstream root-package tests/fixtures are copied into a scratch
  module containing the new API and replayed under the race detector. Their
  source is from the pinned module cache, not modified or added to the runtime.
- An independent python-chess corpus contains the original 84 FENs and 512
  seeded legal random positions (596 total, 588 unique). Every legal move and
  tactical eligibility, check/capture/EP/castle tag, child FEN, Polyglot hash,
  status and depth-two perft is checked. FEN comparison uses the full EP target
  convention. Every expected row and any failed case is preserved.

The fixed-depth corpus and deep subset are byte-identical to v0.16's frozen
84/12 inputs. I use the same cold-table driver, full histories, depths 3/4 and
5/6 and per-case ten-second safeguard. This backend must retain **node counts,
exact scores and chosen moves**, not just membership in some reference move
set. All 192 pairs pass. Parsed roots are reused across depths; cold means
fresh tables, not a completely cold legal-move cache. The cases are public and
correlated, not blind generalization data.

## Performance

Baseline and candidate use isolated native UCI binaries. Serial 100 ms tests
have three trials with reversed actor/case order; 1,000 ms tests have one trial.
The three-search load has two reversed actor blocks with three independent
processes each, six trials per actor. Actual parent/next/cold-next UCI episodes
exercise real cache generations. The schedule is identical to the v0.16
experiment and retains all 2,046 requests, including expected terminal nulls.
Requested times, engine margins, collectors and watchdogs match. Raw collector
serialization is part of the shared-host load. I run no other engine probes
concurrently with timed blocks or a strength screen.

The profiles use 60 iterations of the same 24 neutral-position, 100 ms kernel,
with allocation-stack sampling disabled and normal GC. Fixed-depth wall times
are one run per case/actor; memory figures are timed-kernel means across these
60 requests. Per-request work differs because the candidate can search farther.
CPU percentages overlap and use different profile totals. They do not measure
strength or a universal throughput bound. Background desktop/thermal/power
conditions and finite timing noise are limitations.

If timing and correctness pass, I run the standard concurrent 50 games each
against tagged v0.16, Stockfish Skill 2 and pinned Sacrifice. The source remains
frozen and every result is retained. A failed practical gate remains a failure;
no rerun or larger confirmation follows automatically. The initial candidate stopped before commit/push/PR; I subsequently authorized
promotion as v0.17, with a PR/merge/release and reassessment of one next candidate.
No live games are authorized.

## Provenance and restoration

`checks.json`, `measurement-checks.json`, `profile-command.json`, source and
binary hashes, raw logs and `ERRATA.md` preserve every command/result, including
the initial formatting scope failure. The development UCI name is distinct
from the released baseline. `candidate-source.json` freezes the complete overlay
on tagged v0.16, including the local library and license. `reproduce.py` restores
and tests it in a new directory without games. Original release source/evidence
and the active bot executable/config remain unchanged.
