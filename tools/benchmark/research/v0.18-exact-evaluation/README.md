# Isolated exact evaluation kernel

I selected this bounded revision after reviewing the two failed fitted-model
short-budget gates. I preserve that model in
[the original investigation](../v0.18-evaluation-model/README.md) and remove its
learned terms and feature extractor from the current build. This candidate
computes **exactly v0.17's evaluation**, with less repeated work. The screened executable remains named
`v0.18-exact-evaluation-dev`. **I recommended promoting
this exact-score kernel under the existing practical gate, while retaining
the fitted model as unpromoted research.**

[PLAN.json](PLAN.json) freezes the scope and gates before these measurements.
The prior base-only ablation was already inspected and supported choosing this
revision; I do not present it as fresh independent validation. No material,
placement, search, quiescence, transposition-table, history, clock, dependency
or position-update policy changes here. I restored the original v0.17 search
score tests; only the UCI handshake test changes for the development identity.

## Implementation and correctness

`Board.EvaluationPieces()` returns a value copy of the existing 12 colored
piece bitboards. No caller can mutate board state through that value. A
38,400-byte immutable table precomputes the original integer material/placement
formulas for 25 phases, 12 piece kinds/colors and 64 squares. Evaluation counts
phase with population counts and visits occupied bits with leading-zero counts.
Black rank/sign conversion and the original per-king tapered rounding remain
exact. Both colors and the promoted-material phase cap are preserved.

All 19,200 entries match the literal old formula. The independent original
routine matches on 596 positions/both perspectives, with additional synthetic
promoted-material, empty-material, ownership and zero-allocation controls.
The full race suite, vet, native/ARM64 builds, offline public-config wrapper,
history-contamination controls and original root tests pass. Every original
pinned library runtime file remains byte-identical; the only new library file
is the value-copy getter. The existing v0.17 search API is unchanged.

All 192 fixed-depth pairs match **nodes, scores and selected moves** exactly.
The same public full-history inputs and original driver are retained.

| Fixed-work cohort | v0.17 | Exact kernel | Node count, both |
| --- | ---: | ---: | ---: |
| Depth 4, 84 positions | 2.785 s | 2.450 s | 1,010,100 |
| Depth 6, 12 positions | 9.580 s | 8.150 s | 3,658,933 |

These runs reduce elapsed time about 12% and 15%, respectively. They are
selected workload observations, not universal speed bounds or a strength gain.
[analysis.json](analysis.json) preserves every comparison/group and timing.
[checks.json](checks.json), [source-hashes.json](source-hashes.json) and
[provenance.json](provenance.json) preserve source, commands and binary hashes.
[ERRATA.md](ERRATA.md) records the premature analysis invocation; no engine
measurement was rerun or replaced.

## Equal-time and strength checks

All 2,046 native requests complete correctly: 2,018 legal nonterminal results
and 28 expected terminal nulls. Every response remains within the requested
budget plus the existing 50 ms allowance. The collector schedule, margins and
host conditions are recorded; no other engine probes overlap the timing blocks.

| Equal-time cohort | v0.17 mean depth | Exact kernel mean depth | Candidate node increase |
| --- | ---: | ---: | ---: |
| Serial 100 ms, 252 pairs | 5.115 | 5.214 | 17.2% |
| Serial 1,000 ms, 84 pairs | 6.750 | 6.833 | 15.7% |
| Three-search 100 ms, 504 pairs | 5.083 | 5.206 | 18.0% |

The respective deeper/equal/shallower counts are 27/223/2, 8/75/1 and
62/442/0. Warm requests are separate: 183 per actor, including 14 terminal
nulls each. One warm candidate request is shallower than its paired baseline.
[shallower-timing-pairs.json](shallower-timing-pairs.json) preserves all four
cross-actor depth regressions. A lower cost does not guarantee a deeper result
on every timed request. I have not isolated the causes of these outliers;
finite timing noise and depth boundaries are possible factors, not a diagnosis.

After reviewing correctness, fixed-work and equal-time gates, I started exactly
one standard screen: 50 games each versus tagged v0.17, Stockfish Skill 2 and
pinned Sacrifice, with all three matchups concurrent. The source/binaries are
frozen and [screen-command.json](screen-command.json) records the command and
status. The screen completed with exit 0 in about 10.4 minutes:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.17 | 24 / 9 / 17 | 0 | 0 | 57% |
| Stockfish Skill 2 | 36 / 1 / 13 | 0 | 0 | 73% |
| Pinned Sacrifice | 21 / 9 / 20 | 0 | 0 | 51% |

All 150 PGNs independently replay legally, with matching results and mate
outcomes. [screen-review.json](screen-review.json) retains color/opening
breakdowns, timing, every game and illustrative paired intervals.
[The complete match artifacts](../../results/v0.18-exact-evaluation-dev-vs-v0.17-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
retain the original labels, opponent identities/options, hashes and PGNs.

The direct practical gate passes (>50% with unfinished as losses, zero
forfeits). White/Black direct scores are 58%/56%, but the gain clusters strongly
by opening: the first two openings contribute 12 wins in 12 games. The nominal
pair interval is 42.25–71.75%, and repeated public openings/shared-host play
limit independence. In a post-result sensitivity summary, weighting each of
the ten opening families equally gives about 52.08% direct rather than the
official 57%; I do not replace the recorded score or invent a new gate.
[opening-sensitivity.json](opening-sensitivity.json) retains all three values.

The 73% Skill 2 result is encouraging; the 51% Sacrifice result is near even.
Historical v0.17 external scores are context, not a simultaneous external
baseline control, so I do not claim a causal 20-point Skill 2 improvement,
precise Elo, or a large general strength gain. No Stockfish target advances.

I recommend promotion of the **isolated exact-score kernel** under the existing
practical gate: fixed-work efficiency is established and the routine strength
screen supports it. I do not recommend promoting the rejected fitted values.
A separate 100-game-per-opponent / 50-opening confirmation is justified before
claiming broad strength or raising the Skill target, given opening sensitivity
and the unusually high Skill 2 sample. It is optional under the current release
policy and has not run. No larger confirmation, live game, commit, push or PR
was started automatically.

## What I changed in the research approach

[LEARNING_REVISION.md](LEARNING_REVISION.md) analyzes the already inspected
model failures and gives the next data/objective plan. I do not hand-adjust
weights or relax a gate based on those examples. The descriptive comparison
finds short-budget losses concentrated in a few large branch mistakes; mean
static error alone did not select a useful short-budget evaluator.

I separate two questions: whether computing the same evaluation faster helps
play, and whether learned values improve decisions. This candidate tests the
first. A new independently validated model is still needed to answer the
second. No book, external dataset, broad board rewrite or new pruning is
bundled into this revision.

[reproduce.py](reproduce.py) restores the complete frozen source overlay over
tagged v0.17 in a **new** directory, then runs tests/build without games. The
original failed model variants have their own separate restoration script.
The [backend note](../../../../docs/CHESS_SEARCH_BACKEND.md) explains the local
library boundary and why this additive getter is a small extension of it.

## Reassessment and handoff

[NEXT.md](NEXT.md) prioritizes a new decision-aligned evaluation dataset/model
as the next larger project, with the exact kernel as a cheaper inference
foundation. I do not automatically follow the old hypothesis ranking. The
larger board rewrite and additional pruning lack a stronger measured case at
this point. Functional efficiency and the finite strength screen remain
separate conclusions.

The clean tagged-source restoration passes tests and build; its logs and
[archive-replay-checks.json](archive-replay-checks.json) are retained. The final
[source/archive audit](final-audit.json) verifies frozen code/weights, all
699 earlier inventory entries and local document links. The original screen and failed models remain byte-preserved; the active
bot/config remain unchanged.

## Subsequent promotion as v0.18

I approved PR/merge/release as v0.18 after reviewing the 57% / 73% / 51% screen.
The release changes only UCI identity and its handshake test from the screened
Go source; [release-source-identity.json](release-source-identity.json) verifies
that boundary. Full race/vet/native/ARM64/offline wrapper checks pass again.
[release-provenance.json](release-provenance.json) identifies the release build.

The separate six-game concurrent integration smoke finishes legally without
forfeits or unfinished games: 1/0/1 against v0.17, 1/0/1 against Skill 2 and
2/0/0 against Sacrifice. Its two-game direct score is 50%, so the runner returns
exit 1. I retain that outcome without rerunning it; it is integration evidence,
separate from the passed 50-game strength screen.

I authorized reassessing and starting one next candidate after publication,
with a fresh short ranked research note and a check-in before its PR. No larger
confirmation or live game follows this promotion. The archived failed model
reports describe their original development stages, not the released policy.
