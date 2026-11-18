# Immutable search transition specialization

The original library board update scans all twelve piece bitboards. Each pass
clears the destination and dispatches through generic getters/setters, even for
an ordinary quiet move affecting one set. Search pays this cost while validating
candidate move tags and again when creating each searched child.

The additive `search_update.go` supplies three helpers:

- `searchPieceSet` selects an affected piece bitboard directly.
- `updateSearch` updates moving/captured piece sets and handles promotion,
  en passant and castling, then uses the original convenience-field refresh.
- `Position.UpdateSearch` preserves immutable board copies and the original
  castle/EP/counter/hash helpers. Nil delegates to the original null transition;
  no null-move search is introduced.

`searchMoveTags` retains the original tag logic and original attack functions,
using the specialized temporary board transition. The authored search generator
calls it. Main search, quiescence and root search use `UpdateSearch` for children.
Original public game/update/tag methods remain untouched and serve as oracles.

This does not introduce make/unmake, an undo stack, a different position layout,
new caches, omitted tags, shared attack rays, quiet-mate proof, evaluation values,
move ordering, search reductions, table policy or clock changes. The API requires
the generated legal moves used by search, matching the existing unvalidated
`Update` contract at its call sites.

Focused controls compare 37,795 legal child transitions against original methods
over seeded traversals and explicit castling, EP, pins and all promotions. They
check flags, FEN, bitboards, mailbox, hashes, check/status, subsequent legal moves,
parent immutability and nil delegation. Existing generation/perft, repetition,
mate, timeout and root-search tests continue to apply.

The original pinned Go runtime files remain byte-identical. Only authored local
additions and the three child-update call sites change. The MIT license and
original source hashes remain applicable. Future upstream updates must replay
the original/runtime-specific controls and compare these specialized transitions
against the new general update implementation; this additive boundary creates
a synchronization obligation even without a full board rewrite.

Initial and all three alternating depth-four timings preserve exact scores,
root choices and nodes. The median elapsed ratio is 0.78865, about 21.1% less
time for identical work (about 26.8% higher throughput). Depth three, the deep
depth-six subset and 145 newly generated depth-four controls show similar ratios
and exact tree identity. These selected-workload timings establish efficiency
within the controls; equal-time choice tests and standard match evidence are
separate gates and cannot be inferred from the speedup alone.
