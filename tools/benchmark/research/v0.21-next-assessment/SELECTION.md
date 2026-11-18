# Profile-driven selection: preserve ordering with smaller sort records

The fresh v0.20 profiles change the first implementation choice. Search move
ordering has a larger isolated case than further temporary board simulation:

| Observation | Public roots | Newer development roots |
| --- | --- | --- |
| Searcher move generation/ordering, cumulative | 42.60% | 38.87% |
| Stable sorting of full Move records, cumulative | about 21.1% | 21.26% |
| Priority callback, cumulative | 14.81% | 12.29% |
| Search-specific move tagging, cumulative | 16.23% | 15.28% |
| Runtime bulk copy, flat across all callers | 10.75% | 15.28% |
| Runtime map lookup, cumulative across all callers | 14.20% | 11.96% |

These percentages overlap and cannot be added. Copy and map costs also have
other callers, including repetition, which remains unchanged. The profiles cover
selected depth-four roots rather than full-game time distributions. Allocation
profiles show fresh-table allocation and immutable child copies dominate bytes;
the driver creates a fresh table per root, whereas the production table persists
within a game. This does not justify changing TT lifetime, size or draw policy.

The current stable-sort comparator recalculates the same material-based priority
for every comparison and passes full `chess.Move` values containing PGN-tree
metadata. The desired final order depends only on a fixed integer priority and
the generation order for ties. The initial implementation will decorate each
move once with its original index and **the identical priority formula**, sort
those small records stably, then apply the resulting permutation in place.
Already ordered keys need no permutation. Per-ply rank storage remains owned by
the searcher; descendants cannot overwrite an ancestor's list or rank data.

No priority value, tie rule, move eligibility, legal generator, board transition,
exchange estimate, evaluation, TT, search horizon or clock changes. The original
reference ordering helpers remain comparison oracles. Required controls include
all-quiet ties, captures/EP, promotions/underpromotions, castling, empty/single
lists, cold/warm positions, nontrivial permutation cycles, repeated buffer use
and active-ply isolation. The existing 500-position ordering controls and 325
distinct exact-tree comparisons remain mandatory.

The original 15% median elapsed reduction gate applies before native probes.
If it fails, preserve the exact attempt and reassess using measured data. A
compact search-move representation or attack-only legality simulation remains
a separate possible revision, not part of this first implementation. No guessed
piece values, root preferences or cached attack state is introduced.
