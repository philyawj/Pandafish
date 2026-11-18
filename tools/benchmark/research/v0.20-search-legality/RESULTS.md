# Legality-kernel results

**Recommendation: promote the candidate after review.** Immutable search
transition and legality specialization removes substantial work while preserving
every measured tree, passes the native gates, and scores 65% against the tagged
baseline in the standard screen. The candidate remains uncommitted pending review.

## Fixed-work efficiency

All 325 distinct position/depth comparisons retain exact scores, chosen moves
and search nodes. The repeated depth-four baseline/candidate elapsed ratios
are **0.61256, 0.61487 and 0.60983**, median 0.61256: **38.7% less time**, or about
63.2% more selected-workload throughput. The normalized ratio to the earlier
transition-only median is 0.77672, a further 22.3% reduction. This exceeds the
new revision's declared minimum additional 10% gain and total minimum 29% gain.
It is a substantive change after the predecessor's failure, not a favorable rerun.

| Workload | Positions | Identical search nodes | Baseline ms | Candidate ms | Elapsed ratio |
| --- | --- | --- | --- | --- | --- |
| Public depth three | 84 | 227,579 | 455.15 | 296.69 | 0.65185 |
| Public depth four, median pair | 84 | 1,010,100 | 2,407.17 | 1,474.53 | 0.61256 |
| Deep subset, depth six | 12 | 3,658,933 | 8,011.78 | 5,086.61 | 0.63489 |
| Newer development roots, depth four | 145 | 3,354,194 | 8,412.39 | 5,257.97 | 0.62503 |

`timing-analysis.json`, raw row files and `timing-checks.json` retain all pairs
and the initial exploratory timing. These are selected local fixed-work results,
not a hardware-independent guarantee or a measured Elo gain. Internal direct-root
terminal returns are compared identically; production terminal behavior is
covered by separate unit/native controls.

## Equal-time decisions

Every suite passes the frozen aggregate tolerance, severe-case tolerance,
mate-risk and response gates. Both newer root turns also pass individually.
Loss is the per-position trial mean of
`clamp(root_reference_CP - selected_reference_CP - 25, 0, 1000)` on fully finite
cases; reference mates are kept separately. Close alternative moves are allowed.

| Suite / budget | Finite cases | Baseline mean loss | Candidate mean loss | Severe baseline / candidate | New mate flags | Response overruns |
| --- | --- | --- | --- | --- | --- | --- |
| Public / 100ms | 67 | 41.8955 | 43.0299 | 1 / 1 | 0 | 0 |
| Public / 1000ms | 67 | 52.9104 | 53.8507 | 3 / 3 | 0 | 0 |
| Newer roots / 100ms | 130 | 40.4154 | 35.0923 | 4 / 3 | 0 | 0 |
| Newer roots / 1000ms | 132 | 29.6818 | 30.3409 | 2 / 2 | 0 | 0 |

Newer 100ms White loss is 42.7692 to 40.2462 and Black 38.0615 to 29.9385.
At 1000ms White is 35.4627 to 37.3582 and Black 23.7231 to 23.1077. The result
does **not** show uniformly better reference choices: three aggregate comparisons
are slightly worse, within the predeclared tolerance. The meaningful pre-game
finding is a substantial exact-tree speed gain without a failing decision gate.

Public 100ms uses three trials per actor; the other blocks use one. Inputs,
teacher hashes, exact queries and cache sources are frozen in each method record.
The independent [choice review](../v0.20-next-assessment/choice-review.json)
agrees with all accepted and rejected reports. Newer roots are reused development
controls, not an independent strength sample. Finite-reference scores and
mate-risk flags cannot prove an otherwise lost position was savable.

## Standard concurrent screen

The single declared screen uses 50 games per opponent, 100ms per move, the
unchanged ten-opening suite and all three matchups concurrent. Full identities,
settings and every PGN are retained in [the screen directory](screen/README.md).
The original development label and UCI name are preserved. The transition-only
candidate's failed 150-game screen remains archived separately.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.19 | 28 / 9 / 13 | 0 | 0 | 65% |
| Stockfish Skill 2 | 33 / 1 / 16 | 0 | 0 | 67% |
| Pinned Sacrifice | 21 / 16 / 13 | 0 | 0 | 58% |

All 150 games complete; the runner exits 0 and the practical previous-release
criterion passes. Independent saved-record accounting verifies all 75 color
pairs, settings, W/D/L and timing counts. No response exceeds the 150ms allowance
including protocol grace. Candidate maximum response is 117.38ms against v0.19,
116.72ms against Stockfish and 112.39ms against Sacrifice. These observations do
not constitute an absolute scheduling guarantee. Legal PGN replay separately
checks every move and saved result/length.

This is encouraging **measured screening strength**, alongside stronger
fixed-tree efficiency evidence. The 65% direct score is a materially useful
screening lead, not a precise Elo estimate. The ten openings repeat within 25
pairs; colors, games and rollout families are correlated. Selecting the revision
after the earlier 45% screen also adds selection uncertainty. Its separate
substantive speed threshold and complete preserved failure history reduce the
risk of misleading presentation; they do not remove that uncertainty.

The prior v0.18 promotion screen scored 73% against Skill 2 and 51% against
Sacrifice. Current 67%/58% samples do not establish a six-point decline or
seven-point gain: these are different runs without simultaneous external
baseline controls. The current direct baseline is v0.19, whose playing policy
is v0.18. No change of Stockfish target is justified from this screen alone.

A broader-opening confirmation is justified before claiming a large/general
strength gain or advancing the target: the existing optional 100 games per
opponent over 50 openings would address part of the opening/sample limitation.
It is **recommended for that question, not run automatically and not added as
a new mandatory promotion gate after the result**. No favorable rerun, live
game, version bump, commit, push or PR follows without the requested review.
