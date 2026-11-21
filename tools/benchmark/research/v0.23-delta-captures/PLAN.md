# Measure guarded quiescence capture-futility opportunity

After six rejected attempts, return to search selectivity on released v0.22.
Earlier exact stand-pat relocation has a failed 5.6% node opportunity; it remains
separate. This candidate would make a heuristic decision about a capture's
maximum plausible immediate gain, not move the same exact cutoff. It is not a
TT, material, clock, ordering, evaluation or normal quiet-check change.

Before pruning anything, instrument the unchanged tree. Eligible moves are
non-checking, non-promotion ordinary captures in unchecked quiescence, with
finite non-mate alpha, phase at least eight, both kings valid, no advanced pawn
within two promotion steps and a halfmove clock below 80. Require static score
plus captured material plus a frozen 200cp positional margin <= alpha. Exclude
en passant initially. Existing SEE eligibility applies first. Root moves stay
unrestricted. The 200cp allowance is heuristic and cannot bound every later
tactic or positional change.

Count the full actually visited subtree of each outermost eligible capture;
do not double-count nested eligible subtrees. Retain all existing moves, scores,
depths and nodes on the 1,145-root warm depth-five workload. Require at least
20% of visited nodes to lie in these subtrees before implementing a playing
policy. Report this as potential work, not guaranteed elapsed savings or strength.

If opportunity passes, freeze this policy and add positive/negative controls
for checking captures, checks/evasions, promotions, advanced pawns, endgames,
mates, alpha boundary, repetition, cancellation and restoration. Compare an
unpruned reference for tactical controls without claiming general equivalence.
Require >=15% median elapsed saving in three alternating unprofiled warm pairs,
counting all actual work and retaining legal complete roots. Then require all
four original native quality/safety gates and new seed 2026112073 final controls
before the one standard concurrent 150-game screen. No retuning after final
controls, unchanged favorable rerun or automatic larger confirmation. Archive
a failure and reassess again; stop before commit at a recommended candidate.
