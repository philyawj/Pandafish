# Isolated correction and policy boundary

Additive `SearchDeadMaterial` uses existing piece bitboards without allocation.
Pawns, rooks and queens immediately reject the subset. Exactly one king per
colour is required. With knights, only a single knight and no bishops qualifies.
Without knights, every bishop must occupy the same square colour. Bare kings
are included. A king and one minor cannot construct mate; bishops confined to
one colour cannot supply the missing opposite-colour control for mate with
only kings remaining. Two knights, opposite-colour bishops and mixed minors
remain live even when a win cannot be forced. No general dead-position solver
is claimed.

Main search and quiescence return zero after mate/stalemate checks, before
repetition push, score lookup or stand-pat. Direct evaluation agrees. The root
still returns a legal move; legal transitions from this subset cannot acquire
new mating material. Cache keys remain unchanged because this draw is entirely
board-dependent. Existing repetition/history and clock rules stay intact.

The baseline policy table reproduces the material mismatch and separately the
move-count mismatch. The initial table placed a bishop on a diagonal checking
the other king and used both root turns, creating invalid controls. Its original
source and rows remain explicit. Corrected controls keep both kings unattacked
for both turns. The correction changes no50/75-move policy: implementing an
optional claim requires an explicit client/search action contract, beyond this
isolated material hypothesis. No claim is made that the known long game was
caused by these rules or that this correction would have won it.

Tests compare the new bitboard predicate with the existing map-based game
material implementation across material/square combinations and separately
assert the accepted subset and live negatives. They check promotions, captures,
en passant, castling, root/q/main values, poisoned score entries, history unwind,
cancellation, zero allocations and mate precedence. The old independent rules
vectors still check move legality/state, not an independent draw oracle.

The first normal suite reports three old expectations for KN/K: two placement
checks and a queen-capture continuation. The original source snapshot and
failed workflow stay preserved. Placement controls now contain equal pawns to
keep the positions live; the existing queen-fork control now explicitly expects
zero after the last queen is captured. Runtime logic is unchanged by these test
repairs. A passed intermediate validation is also retained. A remaining
mixed-minor negative fixture checked one king for one root turn; final policy
vectors move that king to g8 and explicitly reject any unexpected check in
both-turn controls. Both earlier tables/source versions stay marked invalid.
The scratch baseline table command initially ran before replay extraction had
finished and failed to compile against the incomplete directory; after replay
completed, the same command passes. This is an orchestration error, not a rules
or candidate measurement.

The warm comparison is allowed to change values/nodes only through the explicit
material-draw terminal and its propagated ancestors. It is not an exact-value
speed rewrite. Separate live truth-table controls must stay exact; any changed
warm rows require a witness review before later gates. Static material/placement
formulas, move eligibility/order, cache/history and clock allocation are intact.
