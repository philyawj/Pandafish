# Same-horizon full-depth reference

Use the exact candidate checking horizon with reductions disabled in a separate
test worker. Reconstruct every one of the 1,145 warm histories through the shared
harness. The reference optimum is an unrestricted, fully completed iterative
depth-five result. Reevaluate the candidate's legal selected move with a fresh
table, full root window and unreduced depth five. No reference move is required.

Retain the earlier selective-review gate: mean finite positive regret <=25cp,
zero missed winning mates, new losing mates, unsupported mate claims or selected
values above the reference optimum. Compare mate bands, not unsupported exact
distance equality across selective trees. Check all positions, including ties.
This isolates reduction risk within the new horizon. It does not prove that
extra checks improve equal-time strength or replace the six native gates.
