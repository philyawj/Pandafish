# Changed-value witness review

The cost gate passes with median ratio1.006937. Exact node count differs by
0.01937%. The first policy reviewer incorrectly requires every live root's
main value to stay identical, and fails on opposite-colour bishops6→5cp. That
root may reach a dead-material descendant; preserving all live-root values
would prevent propagation of the intended rule correction. Keep the failed
reviewer/log and both values. No failed cost threshold is relaxed.

Before native work, shadow both original baseline and frozen candidate in
scratch replays. Invoke the same material predicate before both original
main/q terminal sites; baseline always continues its original computation.
Record hit counts and full FEN witnesses, using existing warm rows to identify
roots. Instrumented times are excluded from performance claims. Require each
instrumented result to match its own original move/score/node trajectory.
Require a concrete dead descendant for the changed opposite-bishop control and
for every changed warm root, from at least one actor. If a root lacks a direct
witness, investigate cache ancestry instead of declaring it explained.

The changed policy control must remain live itself, with a legally reachable
proved draw descendant. Unrelated pawn/rook/count/mate controls must stay exact.
Record every warm change rather than only beneficial rows. This is a semantic
causality check, not a reference-engine quality estimate or a modified native
gate. No native queries begin until it passes.
