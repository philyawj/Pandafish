# Exact exchange kernel: prospective revision

Baseline release is v0.21. Retain the functionally exact generation foundation
from the first evasion prototype, with all original failures preserved: first
staged picker 47% direct, evasion revision 9.11% extra warm versus a 10% target.
Neither is relabeled as promoted. This is a new exact search-kernel revision,
not an unchanged-code rerun. No clock/history/storage/aspiration/reuse heuristic
is included.

The [post-evasion profile](../v0.22-next-assessment/EVASION_CONTEXT_REVIEW.md)
attributes 14.98% sampled cumulative CPU to exchange evaluation, including 11.54%
to least-attacker selection. Current selection and king-safety checks repeatedly
scan all 64 squares. Replace only that computation with local piece/side bitsets,
fixed attack/ray masks and first occupied blockers. Preserve material values,
legal recapture tests, ascending-square tie behavior, promotions, EP, king state,
stop-or-exchange minimax and every cancellation observation point.

The old mailbox implementation remains a frozen Go test oracle. Compare exact
attacker sets, least legal attacker choices, parent restoration and full scores
over focused pins/xrays/king/EP/promotion/malformed-king controls and seeded legal
walks, not merely signs or selected root moves. Keep existing independent legal
exchange controls, frozen rules, exact complete move order and race checks.

Require all 325 cold and 1,145 warm results/nodes to remain exact. Retain the
15% total cold target versus v0.21, beneficial deep/newer work, and require
**10% additional warm median elapsed saving versus the first evasion revision**
in three alternating pairs. That baseline is already faster than the first
staged picker; this is a substantive new saving requirement before another
screen. Measure one fresh paired warm total comparison versus tagged v0.21 too.
No timing claim comes from the instrumented profile.

Only a surviving kernel reaches the original four native suites, followed by
one standard 50-game-per-opponent concurrent screen. Retain every failed gate;
exceptions require explicit documented review before games. No new search
eligibility, score/cache policy, clock or root restriction is introduced. No live
games, automatic larger confirmation or publication. Check in before committing
a viable result; archive failed code and continue after an unpromising result.
