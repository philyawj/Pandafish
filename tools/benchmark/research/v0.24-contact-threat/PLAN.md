# Protected queen-contact threat: bounded first diagnosis

Start from exact tagged v0.23 after all prior failures are archived. Three
post-release screens fail their original direct gate; no draw/compact candidate
is reactivated. Weekly allowance is constrained, so one Go-owned diagnosis must
earn further work before another candidate or screen.

Hypothesis: at stand-pat beta cutoffs, a cheap geometric guard can identify a
specific opponent quiet mating threat and justify searching actual defenses.
Restrict the certificate to a legal quiet queen move adjacent to the defender
king, on a square protected by an attacker pawn, knight or king. Final legal
move lookup and a real no-legal-reply check prove mate. A hypothetical pass only
selects a threat; it never supplies a forced-loss score.

First test the primitive against the fixed training-only quiet-threat survey,
plus explicit legal escapes, captures, blocks, pinned attackers, promotion/EP,
stalemate and color mirrors. Require zero false mating certificates and at least
eight distinct saved hypothetical-opponent-mate witnesses covered. This is a
bounded subset, not a complete threat/mate detector. No new labels or games.

Only if coverage/soundness pass, shadow the unchanged v0.23 warm tree to count
eligible beta cutoffs, candidate guards and certified threats. Require exact
1,145 moves/scores and43,118,858 nodes. A later candidate must fit a prospective
10% total changed-tree cost envelope; stop early if recognition alone exceeds
that envelope. No width/guard/subset/threshold sweep follows a failure.

A candidate is permitted only with an explicit response contract: suppress the
unsafe stand-pat cutoff, search all legal defenses without SEE omissions, and
revalidate the declared mating move after each actual defense. Complete legal
responses, repetition/history and mate precedence are mandatory. One threat
verification event per q path is the frozen local budget; original depth caps
remain. No blanket checking search, heuristic forced-loss score or accept-from-
stand-pat verification is allowed. Later unchanged native/fresh gates and one
standard screen are required. Stop before publication; no live games.

## Pre-activation window-consistency review

The initial shadow samples existing beta cutoffs only; that measures recognition
cost/opportunity and is not an activated q-value definition. A playing candidate
must apply the same threat/response contract in full and narrow windows.
Otherwise retaining unverified stand-pat in a full window but rejecting it in a
scout window would contradict PVS/cache value assumptions. Therefore activation
must test the same geometric/certificate guard at every nonchecked q node with
at least two remaining ordinary plies and unused threat budget. Suppress stand-
pat whenever the certificate exists, regardless of beta. Defense plus revalidated
mate must fit those two plies. The single-event budget propagates through all q
children; pending threat context is not cached. Require full/narrow-window and
explicit-response reference agreement before any native measurement.

This design condition is recorded before reviewing shadow timing. The original
10% envelope is unchanged; candidate cost includes the wider consistent guard
and actual defenses, not just the initial beta-cutoff shadow. A failed initial
recognition gate remains failed. The early harness duplicate-declaration build
failure and corrected source are preserved; it changes no threat criteria.

Before activated candidate measurement, add a generous early node-growth ceiling
of1.25 after the first warm pair. This only stops obviously expensive work before
further pairs; the original elapsed ceiling1.10 remains mandatory. No timed
candidate outcome has been inspected when adding this allowance-saving gate.
