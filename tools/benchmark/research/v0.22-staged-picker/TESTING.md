# Staged-picker controls and measurement contract

The [plan](PLAN.md) and [frozen mechanism](../v0.22-next-assessment/SELECTION_PICKER.md)
precede implementation. Five Go overlays/hashes over released v0.21 reproduce
the candidate. Only `search.go`, the new private picker/tests and two additive
local-library files change. All original pinned upstream runtime files remain
byte-identical. Material/evaluation, root/quiescence policy, clocks and score
cache rules remain unchanged.

The corrected direct-lookup/quiet-subset controls compare 406,674 legal/stale/
malformed cases with the unchanged original legal methods. Cold and warm cases
include castling, transit/self-checks, both EP directions, pinned EP, every
promotion/underpromotion, quiet checks/evasions, terminal states and deliberately
inconsistent supplied castling rights. Prefix and complete-cache ownership are
checked. Whole Move values and cached metadata are compared, not only strings.

10,300 complete staged lists exactly match the eager full-generation + stable
priorities + killers + table preference pipeline across 500 seeded positions,
focused fixtures, cold/warm caches and legal/invalid/duplicate/tactical killer
hints. Active-ply/outer-buffer growth cannot change pending parent stages.
Preferred-only control proves later buffers are not materialized and incurs zero
allocation. A staged startpos depth-four perft returns 197,281. Existing 37,795
immutable transition controls and all move-backend/cancellation/mate/history/GHI
controls remain active.

The initial inconsistent-rights failure, exact five-file source and logs stay
in [ERRATA.md](ERRATA.md) / `initial-source.json`; corrected controls precede any
timing. No failure is silently converted to a pass. Independent frozen vectors
pass all 596 move/tag/child-FEN/hash/status/depth-two-perft positions. Full race,
vet/build and pinned upstream root-package race checks pass. The unused upstream
rendering subpackages are not reintroduced.

The original Go direct-root driver compares three alternating public depth-four
pairs, public depth three, 12 deep depth-six roots and newer depth-four roots.
All 325 distinct scores/moves/nodes agree; every pair is retained. The 15% median
gate applies before native/game work. Three warm depth-five iterative pairs use
the same frozen 20 games / 1,145 candidate-turn positions, one TT per game and
fresh workers per root. All scores/moves/nodes must agree across actors and
repeat exactly per actor, with a 15% median elapsed gate. Parsing/MemStats/output
overhead remains separate; no instrumented-opportunity timing is an efficiency
claim. Saved drivers refuse old result paths.

After correctness/efficiency pass, require the existing four native suites and
their unchanged finite-reference loss/severe/mate/timing/both-turn gates. Cache
arguments use complete prior run directories; the earlier clock phase's metadata-
file error and reference-only repair remain a separate failure record. Production
roots stay unrestricted, close alternatives are accepted and teacher labelling
does not establish a unique acceptable move or savable finite position.

Only a surviving candidate reaches one standard 50-game screen against tagged
v0.21, Skill 3 and pinned Sacrifice, all three concurrent, 100ms plus 50ms grace,
ten unchanged openings and 400 plies. Preserve all results and independently
replay/account for the complete screen. No heavy validation overlaps timings,
native blocks or games. No larger confirmation or live games run automatically.
Stop before committing a viable result for the requested review.

The original four-gate prerequisite is not met. Public 100ms fails its mean-loss
tolerance; the other three suites pass. [The diagnosis](NATIVE_DIAGNOSIS.md),
independent `choice-review.json` and ten exact matched-depth comparisons precede
the [explicit research-screen exception](SCREEN_DECISION.md). The original
failure remains false; no native actor request is retried. `warm-review.json`
independently checks all 3,435 paired root observations, recorded-input hashes,
helper identity, actor repeatability and baseline-best membership through exact
choice equality. The raw analysis's generic pending-membership flag is an older
checkpoint, resolved by that separate review rather than rewritten.
