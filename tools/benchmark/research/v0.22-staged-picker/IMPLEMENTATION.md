# Deferred main-search move generation

The previous pipeline generated/tagged all legal moves, precomputed stable
priorities, moved killer hints and finally preferred the table hint. Only then
did it try the first move. When that move cut off, the remaining generation was
unused. The frozen warm diagnosis finds that early table/tactical/killer cutoffs
cover about 80% of main-search ordinary quiet records. That is an opportunity,
not a claimed speedup; actual paired tests include validation/staging overhead.

The new picker preserves the complete eager sequence:

| Stage | Work performed only when reached | Preserved order |
| --- | --- | --- |
| Table | Decode and freshly validate the captured table/board hint. | The same legal hint formerly moved to the front. |
| Tactics | Append legal captures/EP/promotions and apply existing stable priorities. | Original tactical order, excluding a previously tried hint. |
| Killers | Freshly validate the captured quiet killer 0/1 signatures. | Original quiet-killer precedence, excluding already tried hints. |
| Ordinary quiet moves | Append legal noncapture/nonpromotion moves. | Original generation order, excluding the preferred quiet hints. |

The table and killer priorities are captured before searching a child. This
avoids accidentally reading a newly changed descendant board hint when returning
to the parent. A capture that matches a stale quiet killer remains a capture;
the killer stage accepts only freshly validated quiet moves. Invalid/duplicate
hints add no move. Every legal move appears exactly where the eager pipeline
would have put it; there is no killer-before-capture policy or new priority.

`LookupSearchMove` is an additive method in the pinned chess extension. It checks
signature bounds, current-side occupancy, destination ownership, pawn promotion
requirements, pseudo-legal geometry and the existing full move-tag/own-king
legality predicate. Castling uses the pinned castle generator, including transit
attacks. A complete cached list supplies the original whole Move value/metadata.
It never constructs a full legal list solely to validate an ordinary hint.

`AppendSearchQuietMoves` shares the existing search generator. In uncached
generation it masks out captures, EP and every promotion before tagging. It
retains quiet checks/evasions and legal castles. In cached generation it copies
the quiet subset without changing cache contents. The original public full/
tactical/existence behavior remains available and is checked independently.

Each active ply retains its own move storage. A parent finishes its tactical
slice before reusing that ply's storage for the quiet slice. Descendants use
larger ply indices; growing the outer buffer slice does not overwrite a parent's
underlying stage. The picker itself stays on the stack, and a validated first
hint creates no later-stage move/rank buffers or allocation in its control.
Position updates stay immutable; no undo state or reused child position is added.

The pinned library's inconsistent-FEN castling-right behavior needs an explicit
compatibility path. If current-side rights exist without the king on its starting
square, the original generator can emit duplicate signatures. The picker retains
the eager path there, and direct lookup retains the original castle-list route.
The [failed initial control and correction](ERRATA.md) are preserved. This
experiment neither fixes upstream malformed-position semantics nor changes
standard legal castling rules.

Main search still tests terminal/draw conditions, probes matching-depth scores,
uses the same first/full and later PVS windows, updates the same hints and stores
the same bounds. Only when moves become available changes. No extra search node
is added or removed at fixed work. Root and quiescence generation/order/eligibility,
SEE, values, placement, mate/repetition identity, TT capacity/replacement and
time/cancellation policy remain v0.21. All earlier failed code is inactive.

Late nodes can require a second pseudo-destination pass and hints cost direct
validation. The benefit comes from avoiding later work at cutoffs and avoiding
full-list quiet rank/permutation work, not a universal cheaper lookup guarantee.
Frozen source, exact nodes and paired timings determine whether that trade pays.
Game strength remains a separate measured screen, with finite-sample limits.
The new extension could support later staged/checking-move research, but none
of those possible features is included or automatically approved here.
