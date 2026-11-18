# Search-legality kernel revision

The transition-only candidate establishes an exact-tree 21% elapsed reduction
but fails its 50-game previous-release criterion. Its complete result is archived
separately. This revision requires a substantive additional measured gain before
any new strength screen; it is not an unchanged-code retry.

Hypothesis: search-specific attack tests can reuse each computed ray and avoid
generic piece lookups, while legal-move existence needs no opponent-check tag.
Retain the original predicate's geometric fast rejection, compute each sliding
ray once, select piece sets directly, and compute full check tags only when the
caller actually returns moves. Own-king legality is always checked. Keep the
transition specialization as the previously isolated functional foundation.

No evaluation, root preference, quiescence eligibility, depth, table, clock or
opening change is allowed. Quiet-mate proof remains archived and absent. Original
upstream runtime stays untouched as a reference, and an independent Go mailbox
attack oracle must agree over both colors/all squares on seeded legal positions.
All transition, tag, child, hash, perft, mate/repetition and timeout controls apply.

Require exact scores, moves and nodes on the same 84 public depth-three/four
positions, deep depth-six subset and 145 development regression positions. Run
three alternating depth-four baseline/candidate timing pairs. Require at least
10% additional elapsed reduction versus the transition-only median (and at
least 29% total reduction versus v0.19); reject the revision before a screen if
that threshold fails. Keep every initial timing and result. Native public/fresh
100ms/1000ms gates remain unchanged, including each fresh root turn.

Only after all gates pass will another standard 50 games per opponent run,
three concurrent, versus tagged v0.19, Skill 2 and pinned Sacrifice. Preserve all
earlier screens alongside it. Selecting a candidate after earlier screen failures
introduces selection uncertainty; a positive screen is provisional evidence,
not a precise Elo estimate. No automatic larger confirmation, live games,
commit, push or PR. Check in only when a valid candidate warrants review.
