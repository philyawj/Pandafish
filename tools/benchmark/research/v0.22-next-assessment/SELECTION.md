# Selection: reusable child state by active search ply

Baseline observations precede runtime edits. Both modes complete the same 1,145
candidate-turn roots from 20 saved games. Full position reconstruction takes
562.4ms across 106,404.4ms of search, approximately 0.53% of the measured pipeline.
That does not support prioritizing protocol-history caching here. Neither mode
has a search response above 150ms. Their timed node counts/depths differ, so
allocation totals are descriptive, not equal-work comparisons.

In search-only allocation sampling, immutable child boards and positions account
for 97.5% of bytes (about 33,075 MiB total sampled allocation); full-pipeline
sampling attributes about 96.0% to those children. Keeping the TT per game removes
its former benchmark artifact but does not remove the child-allocation target.
CPU cumulative UpdateSearch costs are 8.7%/8.4%; allocation and garbage-collector
work also occur outside that call. These overlapping percentages are not summed.
Measurement includes PGN setup, profiles, per-root MemStats and output overhead;
runtime/scheduler percentages cannot all be attributed to child allocation.
Untimed equal-tree pairs, without profiling or per-root MemStats, must establish
the actual elapsed gain.

Implement one mechanism: an additive caller-owned child-state cell containing a
Position and Board, privately retained per active search ply. Copy the parent
board and run the existing affected-piece update/hash rules into that cell.
Overwrite a sibling only after its recursive calls and deferred history cleanup
return. Re-search retains its child until the full-window call finishes. Fully
reset cached legal moves and every Position field. Reject source/destination
aliasing before writing. Do not expose this storage to the game tree, root
protocol position, TT, another worker or asynchronous readers.

The original immutable Update and UpdateSearch remain untouched comparison
oracles. No undo stack, in-place parent mutation, changed move/tag rules, guessed
values, search pruning, TT lifetime/policy, history identity or clock change is
included. Coallocation alone would retain a per-node heap object; per-ply storage
tests the larger allocation opportunity while retaining value-copy transitions.

Required controls cover all original transition cases, independent frozen child
FEN/hash/perft vectors, legal-cache reuse, null transition, alias rejection,
ancestor/worker isolation, cancellation cleanup, repeated PVS use and startpos
perft. Apply the frozen 325 exact comparisons and a persistent-table sequence
before native and game gates. If the 15% fixed-work target fails, preserve the
attempt and reassess rather than weakening the gate after observing results.

Remaining attack/tag work is the next exact-kernel alternative. Clock reads are
visible but changing check frequency would change timing responsiveness and is
outside this experiment. A new evaluation representation remains separate work.
