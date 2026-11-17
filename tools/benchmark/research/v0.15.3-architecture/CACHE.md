# What the cache investigation establishes

The cache is neither useless nor established as the main strength bottleneck.
The v0.15 draw-context protection greatly restricts transposition reuse, but
removing it is incorrect. A safe structural revision recovers useful work,
especially deeper in the tree. This is diagnostic efficiency evidence, not an
Elo result or a change shipped in v0.15.3.

## History of the current rules

[cache-history.txt](cache-history.txt) and [v015-history-diff.patch](v015-history-diff.patch)
record the source history. The cache existed before v0.15. Exact-depth score
reuse predates the repetition correction. v0.15 adds a count-set fingerprint
to distinguish positions whose scores depend on different repetition histories.
These are different decisions; the earlier cache introduction was not simply
an engine regression.

The score lookup requires the same board hash, halfmove clock and complete
count-set fingerprint. Search retains that fingerprint across pawn moves and
captures, even though earlier positions cannot recur after those moves. It
also uses the same restrictive identity for move hints. The UCI history loader
already trims to the reversible segment, so a cached search subtree can carry
irrelevant history that a later actual move request correctly omits.

The lookup policy is therefore conservative in two ways:

1. It rejects scores from genuinely different draw contexts, which is needed.
2. It also rejects board-derived ordering hints and some contexts that differ
   only in unreachable history, which is avoidable.

The search uses a table entry's score only at exactly the requested remaining
depth. Deeper entries can still supply move ordering. Changing this to a
minimum-depth contract is common in other engines, but changes this engine's
fixed-depth score contract. It must be an explicit design change with suitable
tests, not a one-character fix justified by speed alone.

## Measurements and their limits

In cold depth-four searches over the frozen 36 positions, the instrumented
released policy makes **56,145 probes, zero hits and zero score cutoffs**.
An unbounded shadow map sees 5,378 previously completed board/halfmove contexts;
all have different history fingerprints. Of those, 5,241 would meet a nominal
same-depth bound test if history were ignored. These are untrusted potential
reuses, not proven safe cutoffs or the hit rate of an equal-sized real table.
The diagnostic counters reproduce the uninstrumented scores, moves and nodes.

The table still helps ordering during iterative deepening and supports warm
requests. Repeating an identical root at identical depth is a particularly
favorable warm-cache test, not normal gameplay. Both a cold single-depth search
and a repeated-root unit test alone give an incomplete picture. Native UCI
warm/cold request pairs are saved separately; they use real per-request table
generation handling. The fixed-depth warm toy harness intentionally does not
advance generations and is not presented as a live-request reproduction.

Two conservative prototypes were tested separately and together:

- **Board-only move hints:** a separate fixed-size hint table, validated against
  generated legal moves, with score reuse still protected by the original
  history key. It adds about 1 MiB rather than changing score-table capacity.
- **Reversible-context fingerprint:** reset only the cache fingerprint when
  descending through a pawn move or capture; retain repetition counts and
  restore the previous fingerprint on return. Earlier positions are unreachable
  after these irreversible moves. Direct entry points and manual repetition
  controls do not reset history. Castling-right irreversibility is left
  conservative, so this is not the most permissive possible scheme.

Together, they reduce nodes about **4.5% at depth four** across 36 cases and
**16.1% at depth six** across the six prespecified deeper cases. The full-depth
scores and complete best sets match the baseline on those probes. All existing
root/internal tests and the explicit graph-history control pass. This supports
a future tested redesign; it does not certify every possible history or show
a measured playing-strength gain.

A 16-times-larger score table, 32 MiB instead of 2 MiB, changes none of the cold
fixed-depth node counts. Its unit-test failure is the deliberate resource-limit
assertion, not an observed illegal move. Increasing memory cannot repair an
identity rule that prevents a hit.

## Why board-only scores are rejected

[unsafe-board-scores-ghi.json](unsafe-board-scores-ghi.json) reproduces an actual
search-value contamination, rather than inserting a fabricated score:

| Same root board and halfmove clock | Score for side to move |
| --- | ---: |
| Fresh search without relevant earlier repetitions | -902 cp |
| Fresh search with the drawing history | 0 cp |
| Warm board-only cache seeded by the first search, then given drawing history | -902 cp |
| Warm original or history-safe prototype, given drawing history | 0 cp |

The losing side has a drawing knight repetition. A board-only parent score
prevents searching that continuation. Testing only whether the *cached score*
is zero would not fix this: the contaminated entry is a nonzero losing score.

The unsafe prototype also matches all cold shallow fixed-depth values in this
suite. That does not make it safe; the adversarial warm-history test is needed.
Its root history-isolation tests fail and those failures are retained.

This is the graph-history interaction problem, described in
[Kishimoto and Muller's AAAI paper](https://cdn.aaai.org/AAAI/2004/AAAI04-102.pdf).
[Stockfish's pinned search](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp)
also has explicit GHI/rule50 precautions. Neither source establishes that our
prototype implements a general solution, or licenses removing our safeguards.

## What not to conclude

The previous PVS/LMR strength experiment predates v0.15's history fingerprint.
The newer guard cannot explain why that older experiment failed its desired
strength goal. Likewise, a retrospective v0.14.4/v0.15.2 game comparison changes
mate distance and repetition behavior as well as cache context, so its Elo
cannot be attributed to the cache alone.

The evidence supports reclaiming safe reuse as part of a larger search-core
project. It does not support reverting repetition correctness, enlarging the
table as the main project, or promising that a cache patch alone closes the gap
to Sacrifice.
