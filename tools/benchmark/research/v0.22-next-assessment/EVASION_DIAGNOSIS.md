# Exact geometric evasion mask: prospective diagnosis

The [first staged picker](../v0.22-staged-picker/RESULTS.md) remains rejected at
47% direct. Its frozen source is still available for a genuinely revised exact
kernel, as with the separate transition/legality revisions before v0.20. Reuse
does not promote or erase the failed attempt. Lazy validation within reached
stages fails the separately declared 30% opportunity (about 25.2%); no playing
lazy iterator is implemented.

The next hypothesis restricts **pseudo-legal candidates while in check** before
full tagging. Currently the generator examines every piece/destination and only
then rejects moves leaving its king attacked. A single checker permits a king
move, capture of the checker or interposition on its line. Multiple checkers
normally permit king moves only. En passant must bypass the geometric mask
because its captured square differs from its destination and it changes two
occupancies. Every surviving candidate still uses the original full own-king
legality and check tagging. No attack cache, pin heuristic, reduced search or
changed legal order is introduced.

First add diagnostic-only checker geometry to the frozen staged source in a
scratch replay. Count all uncached tag validations (including existence checks)
and candidates the mask would skip, but **continue original tagging/search**.
Require zero legal candidates rejected and exact scores/moves/nodes on the same
1,145-root warm depth-five trajectory. Require at least **50% of all counted tag
validations** to be avoidable before a playing implementation. Keep checked and
existence/full subsets separate. These counts are opportunity, not timing gains.
Fall back to the original candidate set if the king is missing/nonunique, the
check flag has no corresponding checker, or geometry cannot establish a safe
mask. Original upstream runtime files remain unchanged.

If the opportunity passes, use the exact staged picker plus the evasion mask as
a new generation-kernel revision, with all prior source/results preserved. It
must save at least **10% additional warm median elapsed time over the rejected
first staged picker** in three alternating pairs, retain the original 15% total
target against v0.21, and preserve all cold/warm fixed-work results and complete
ordered Move oracles. This is a meaningful extra-gain requirement before any
new screen, not a trivial-edit retry. Use the original four actual-budget native
gates; any failed prerequisite remains explicit and cannot silently pass. One
standard 150-game concurrent screen may follow a surviving revision. No larger
confirmation, live game or publication follows automatically.
