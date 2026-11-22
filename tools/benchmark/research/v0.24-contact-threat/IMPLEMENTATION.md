# Activated response contract

The initial shadow finds53,508 certified threats at25,813,637 existing stand-pat
cutoffs;348,626 geometric guards fire. Its separate instrumented trajectory
matches all baseline moves/scores/nodes. These are repeated occurrences, not
independent mistakes or root losses. The11 training threats remain development
controls; no new labels are collected.

Activation uses the same rule in every window: outside check, at depth>=2 with
one event budget remaining, a proved hypothetical queen-contact mate removes
stand-pat and enables all legal responses, without SEE omissions. Every actual
response reaches a child that revalidates the original queen move under the new
board, including capture tags, own-king safety, check and no legal reply. Only
that real legal mate supplies a mate score. A capture of the queen, removal of
protection, countercheck, king move, block or immediate winning mate can refute
it. Unresolved children retain the ordinary q heuristic, with event budget zero.

The defense and verified reply fit two remaining ordinary plies. All q recursion
propagates the budget; main-search q entries begin with one. Pending context is
never cached. Existing repetition/terminal handling precedes revalidation, and
irreversible child history is reset/restored exactly as before. Source helpers
reuse independent scratch cells; actual search descendants retain normal
per-ply position ownership. No material/move-count model is reactivated.

Focused tests compare full legal defense enumeration and public immutable reply
transitions against the activated path. Positive forced-mate controls and their
color mirrors return mate-in-two losses; capture/countercheck controls refute
that classification. Full/narrow-window bounds, fixed caps, stalemate, history
restoration and cancellation pass. Certificate probes are rule queries, like
SEE work, rather than recursively counted search nodes; elapsed cost includes
them. Functional proof scope remains this declared subset and response horizon.
