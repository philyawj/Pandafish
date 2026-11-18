# Recommendation: promote the revised kernel after review

The revised exact kernel meets the standard direct criterion and produces strong
external rows. Correctness and efficiency are independently supported; the small
screen supports practical promotion, not a precise Elo claim. The public 100ms
surrogate failure remains a disclosed limitation, with its pre-game exception.
No commit, PR or release has been prepared for this candidate.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.21 | 26 / 12 / 12 | 0 | 0 | 64% |
| Stockfish Skill 3 | 35 / 1 / 14 | 0 | 0 | 71% |
| Pinned Sacrifice | 27 / 16 / 7 | 0 | 0 | 70% |

All three matchups run concurrently: 50 games each, 100ms plus 50ms grace,
400-ply cap, ten unchanged openings, Apple M2 / Go 1.25. Runtime is 664.75s;
the runner returns exit 0. The tested label is `v0.22-bitboard-exchange-dev`,
UCI identity v0.21, executable SHA-256
`c5486656b922db819cee1fa5e2851e9ad1b744806cb9d84a301df6a64d657505`.
The [full original summary](screen/v0.22-bitboard-exchange-dev-vs-v0.21-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba/summary.json),
adjacent PGNs, `screen-command.json` and `screen-status.json` preserve options,
all engine/opening/runner identities, source revision, timing policy and host.
No heavy validation or other native/timing batch overlaps play.

Independent legal replay checks all 150 games. Independent screen accounting
recomputes each result, all 75 color pairs, opening identity, response counts,
limits and the direct criterion. No response exceeds the 150ms allowance;
candidate maximum is 98.414ms rounded up, opponent maximum 102.016ms. All values
remain in `screen-review.json`; `screen-strata.json` retains every opening family.
The direct score is 70% White / 58% Black, with sixteen color pairs won, four
drawn and five lost. Both colors contribute to the direct lead.

## What changed functionally

The full cold search retains all 325 distinct score/move/node comparisons and
saves **35.35% median elapsed time** versus v0.21. Three warm pairs save
**14.65% additionally** over the rejected first evasion prototype. A new total
warm pair saves **37.84%** versus tagged v0.21, with identical results at all
1,145 roots and 43,118,858 nodes per actor. Independent saved-pair/source reviews
verify those claims. The new exchange code exactly matches 46,794 full mailbox
scores and 33,152 attacker/least-attacker/restoration comparisons. Root/upstream
race, vet/build, independent rules and scratch replay/cross-build checks pass.

The actual-budget data separately supports more work within unchanged clocks:
candidate completed depth is higher in **184 of 626 paired trials**, equal in
442 and lower in zero. These are repeated/inspected trials, not 626 independent
positions. Among pairs where both actors stop at their deadline, each suite's
median node-rate ratio is **1.51–1.64**. This is selected cold-native throughput,
not the warm trajectory or a universal bound. Full response distributions,
counts, early-finish exclusions and every lower-depth case are retained in
`native-performance-review.json`; no lower-depth cases or fallbacks are observed.

## What remains failed and uncertain

Public 100ms mean teacher loss is **45.8955cp**, above its **45.0746cp** limit;
its original report and driver remain failed. The other three suites pass,
with both newer turns passing. All 1,252 actor queries have no errors, overruns
or new losing-mate flags under the existing convention. All six changed roots
agree at matched depths; the dominant regression is the baseline's own deeper
Qd4 choice. [The explicit decision before games](SCREEN_DECISION.md) explains
why this substantive efficiency revision merits one bounded research screen.
Neither that exception nor a positive match retroactively passes the surrogate.

The first staged picker remains rejected at 47% direct despite its 68% Skill 3
row. The first evasion implementation remains rejected at 9.11% extra warm
saving versus its 10% gate. The new bitset kernel is separately frozen and adds
14.65% over that faster evasion foundation before this new screen. This is not
an unchanged-code retry. All older storage, tags, clock, history, aspiration,
reuse and diagnostic failures remain archived and inactive as documented.

The finite teacher is not uniquely correct, and finite scores do not prove a
position was savable. More completed depth can expose a worse horizon choice.
The ten repeated openings, concurrent scheduling, inspected regression roots
and selection after failed attempts limit the strength inference. A 64% sample
does not establish an exact long-run win rate. Absolute external scores from
different runs are not controlled baseline-vs-external comparisons. Skill 3
remains the standard; no automatic ladder advancement is justified.

## Review and next decision

Recommend promotion for the substantial exact-computation gain and clean
64% / 71% / 70% practical screen, with the failed surrogate explicitly retained.
The [next-step assessment](NEXT.md) ranks a new cost-bounded tactical-horizon
design ahead of another small kernel tweak, conditional on fresh diagnosis.
A broader-opening confirmation (100 games per opponent over the existing 50
openings) is justified before claiming a general/large strength gain or raising
the reference target, especially after candidate selection. It is optional and
has **not** run; no larger confirmation starts automatically. No live games run.
Stop at this review checkpoint before committing, pushing or preparing a PR.
