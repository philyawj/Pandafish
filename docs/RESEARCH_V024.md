# v1.0 research and preserved post-v0.23 investigation

Research baseline: **v0.23**, commit
`28948e7473251d950003ec1f966f0cf26369b817`. **v1.0 promotes the protected
queen-contact candidate and marks the development pause.** Cleanup
[PR #42](https://github.com/philyawj/Pandafish/pull/42) is merged without a release;
`handoff.md` is absent from that commit and stays outside the v1.0 commit.
The research is approved for publication with v1.0. Directory versions are
research labels, not intermediate engine releases. Historical reports retain
their original stopping-point/publication wording, labels, hashes and inventories.
[Release validation](../tools/research/v1.0-validation.md) records the version-only
source boundary and integration checks. No live play occurs.

The [repository audit](REPOSITORY_AUDIT_V023.md),
[ranked assessment](NEXT_ENGINE_ASSESSMENT_V023.md) and
[storage policy](EVIDENCE_STORAGE.md) remain the original post-release decisions.
Existing evidence stays in place. New diagnoses use compact records and source
hashes; existing native/game schemas remain supported.

## Preserved investigations

| Investigation | Result | Decision |
| --- | --- | --- |
| [Shared legality context](../tools/benchmark/research/v0.24-legality-context/README.md) | Exact diagnostic/oracles; optimistic gross proxy 4.98% before context cost. | Does not justify the frozen 8% net objective. No full pin-context playing candidate. |
| [Target-centric attacks](../tools/benchmark/research/v0.24-target-attacks/README.md) | Exact warm work, 1.294% median saving. | Fails 8% cost objective; no native/games. |
| [Compact representation opportunity](../tools/benchmark/research/v0.24-compact-moves-opportunity/README.md) | General Move 128 bytes versus 6-byte candidate; substantial move/castle allocation. | Justifies a representation study, not an elapsed or strength claim. |
| [Temporary compact generation](../tools/benchmark/research/v0.24-compact-generation/README.md) | Exact warm work, 2.166% median saving. | Fails its original 8% objective; preserved as a foundation without relabeling it passed. |
| [Coordinated compact search](../tools/benchmark/research/v0.24-compact-search/README.md) | 8.162% exact-work saving, all six quality gates pass; screen 46% direct /54% Skill 3 /76% Sacrifice. | Direct gate fails; no unchanged retry or promotion. |
| [Conservative dead material](../tools/benchmark/research/v0.24-draw-semantics/README.md) | Functional draw correction, 0.694% overhead within 5% ceiling, all six gates pass; screen 50% /59% /77%. | Direct gate fails; restore v0.23. No unchanged retry. |
| [Contextual capacity](../tools/benchmark/research/v0.24-context-capacity/README.md) | 655,488 parameter bytes, 54,191 exact incremental transitions; required old preference pair has identical witness FENs. | Capacity gate fails before fitting/cost. Any board-only evaluator must give those identical inputs the same value; the target contract needs redesign. |
| [Move-count claims](../tools/benchmark/research/v0.24-move-count-claims/README.md) | Optional 50-move claim actions and automatic 75-move termination; tests, 0.117% observed cost, exact warm trajectory and all six gates pass; screen 48% /63% /70%. | Direct gate fails; all 150 games legally replay. No unchanged retry. |
| [Protected queen-contact threat](../tools/benchmark/research/v0.24-contact-threat/README.md) | Legal mating certificates select complete actual-defense verification; 5.180% cost and all six quality gates pass; screen 54% direct /66% Skill 3 /77% Sacrifice. | Promoted in v1.0 after all 150 games legally replay. Narrow direct evidence; no Elo claim. |

Every failed gate, adverse row, initial diagnostic error and source version
remains preserved. Playing source is restored between different mechanisms.
No coefficient, width, seed, subset or unchanged screen is swept until favorable.
The contextual collision is a frozen-witness target contradiction, not evidence
that neural evaluation cannot improve actual search.

## Preserved move-count contract

The candidate models a current draw claim at 100 halfmoves and an intended-move
claim at 99 when an allowed legal non-pawn noncapture move qualifies. A claim is
an available zero-valued action; a better winning continuation remains eligible.
At 150 halfmoves termination is automatic, with checkmate precedence. Root
restrictions, resets, both turns, check evasions, history restoration and clock-
qualified cache values are tested. The existing table already distinguishes
halfmove counts. No rejected material detector, compact generator or model is
included. Existing repetition policy and other search/evaluation rules remain.

UCI has no claim command. The game client must exercise an available claim;
returning a legal resetting move does not itself claim a draw. The benchmark
already claims when eligible. This research changes no live-client configuration.
The finite q boundary remains heuristic, so positive scores are not proofs that
an available draw can be converted into a win.

All 1,145 fixed warm moves/scores and 43,118,858 nodes remain unchanged in three
pairs. Both modules pass normal/race/vet/builds and 596 independent move/state
vectors. Existing/fresh quality gates pass without new mate flags or overruns.
The newer 1000ms suite worsens from 23.598 to 24.833cp mean finite loss and retains
a 163cp regression. Its depth-5/depth-6 review matches moves, scores and nodes
at each depth; the timed loss remains adverse. Fresh controls are not strength
measurements, and fixed-work equality is local to the sampled trajectory.

The move-count screen finishes and legally replays but fails its direct gate at48%.
Only one candidate turn reaches halfmove99, against Sacrifice. The next bounded
reassessment activates protected queen-contact defense verification. It passes
soundness, window/history, cost, all six native gates and its single screen at
54% direct /66% Skill3 /77% Sacrifice. All150 games finish and legally replay,
without forfeits. All87 source files match before/after games. It merits promotion
review for bounded tactical reliability; general strength/Elo remains uncertain.
Its candidate preserves the original evaluator and adds no rejected draw model.
The greater-than-50% direct gate passes unchanged. Publication of v1.0 is now
authorized separately. No automatic larger confirmation or live game follows.
Further engine development is paused to focus on other efforts. A small screen
lead cannot establish a precise Elo gain or attribute wins to the rule change.
[The checkpoint](RESEARCH_CHECKPOINT.md) contains the current command and all
consumed fresh seeds; later collections must exclude these cohorts.
