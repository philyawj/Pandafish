# Explicit research-screen exception

This decision precedes the only staged-picker game screen. The original plan's
requirement that all four native gates pass is **not satisfied**. Public 100ms
remains failed: mean teacher loss 45.9552cp versus a 45.0448cp limit, with severe
cases increasing from one to two. No threshold changes or actor-query retries
are made. This is an exception to the planned screening prerequisite, not a
retroactive pass or an automatic promotion decision.

The independent saved-query review reproduces all four reports. Public 1000ms,
newer 100ms and newer 1000ms pass, including both required turn gates. Across
all 1,252 actor queries there are no illegal moves, query errors, new losing-mate
flags or timing overruns. All five changed short-budget roots are compared at
every observed completed depth: ten unrestricted iterative comparisons have
identical completed scores, moves and node counts. Both actors select Qd3 at
depth four and Qd4 at depth five on HhPI3vCx. This attributes the observed change
to completion depth on those inspected cases; it does not make Qd4 a good move
or prove that every future position is equivalent.

The implementation preserves complete move ordering and fixed-work results,
with 15.76% cold and 19.81% warm median search-time savings. The broad newer
short-budget mean improves 35.0840 to 31.7634cp, and the other longer-budget
means are within their limits. A single public horizon-sensitive case can
dominate the short-budget mean; a teacher surrogate alone cannot establish the
playing consequence of this exact-work efficiency change.

Run **one bounded standard research screen**, 50 games each against tagged
v0.21, Stockfish Skill 3 and pinned Sacrifice, all three matchups concurrent,
at 100ms plus 50ms grace with the unchanged ten openings and 400-ply cap.
Preserve all games, failures, identities and timings. No larger confirmation
or second favorable-score attempt follows automatically. Review the direct
comparison, external results and the failed surrogate together before presenting
a candidate. A positive screen does not erase the native failure or establish
a precise Elo improvement. No commit, PR or release is authorized for this
candidate until the next review checkpoint.
