# Compact search values with the existing rules backend

The general rules library Move remains the PGN/tree-bearing128-byte value.
The additive SearchMove carries only origin, destination, promotion and tags
in6bytes. CompactMove explicitly projects these fields; SearchMove.Move creates
a metadata-free general move when an old boundary requires one. Full cached
API results keep their original ownership/metadata.

One `searchCandidates` routine retains the original piece/square/promotion order,
evasion masks, tactical/quiet eligibility, tags and castle path. Its output
adapter either appends full Move values or compact records. The full cached
branch appends the original complete values; compact callers request projection.
This avoids a second legal generator and changes neither rules nor order.

Worker-owned compact buffers replace full values in the hot generation, stable
rank permutation, staged picker, killers/encoding and main/q iteration paths.
Parents and descendants retain distinct ply buffers. Full diagnostic/root
compatibility interfaces remain available. The existing root frontend and
ordering contract are unchanged; the production search still returns one exact
best move, with the same tied-choice policy.

`SearchPosition.UpdateCompact` converts the small value at the existing state
transition boundary, reusing v0.23's board/position cells. Board updates, castling
rights, EP, clocks and incremental hashes still use the original additive update
implementation. There is no new make/unmake, global pool or shared mutable state.

The foundation uses a small internal candidate while tagging/existence-testing.
A rights/occupancy-only prefilter avoids calling castleMoves when neither path
can possibly emit a castle. The unchanged castle routine still decides attacked
transit squares, tags and malformed-rights compatibility. No target-centric
attack rewrite or eager pin context is included; those prior outcomes stay failed.

Correctness controls cover compact/full ordered lists, legal/tactical/quiet
subsets, cached metadata, compact lookup/child FEN/hash, parent immutability,
active-ply ownership, preferred-only zero allocation, blocked-castle allocation,
special moves and malformed castling rights. The independent oracle explicitly
checks compact lists, compact child storage and compact depth-two perft on all
596 frozen cases. Existing minimax, mate, repetition, cancellation, UCI and
whole-move staging controls remain enabled. Both modules pass normal/race/vet.

The initial active-ply control compared a nil clone against an empty non-nil
slice with reflect.DeepEqual. Its source and failed output remain archived.
Element equality repairs that harness distinction; no playing implementation
changes to make that control pass. The corrected ownership control still detects
any changed parent element. No cost/native run precedes passed validation.

Search values and node counts must remain exact. Three alternating1,145-root
warm pairs pass with median elapsed ratio0.918382; this establishes an8.162%
selected-workload functional saving, not a universal bound or Elo. All four
existing and both fresh actual-budget gates pass without exceptions/new losing-
mate flags/overruns. Public100 mean loss is slightly worse within its frozen
tolerance, so passing aggregate gates does not mean every choice improves.
The standard screen and its eventual review remain separate strength evidence.
