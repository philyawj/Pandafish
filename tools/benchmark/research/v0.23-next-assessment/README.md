# Research toward v0.23

Investigation starts from released v0.22 after its separately merged live review.
No commit, push, PR or release is authorized before a viable-candidate checkpoint.
The [standard workflow](../../../../docs/RESEARCH_WORKFLOW.md) is linked from
root AGENTS.md; “Find the next improvement using the standard research workflow”
is sufficient for future requests with the same scope.

The [initial plan](PLAN.md) and [baseline diagnosis](DIAGNOSIS.md) distinguish
confirmed live errors from their possible causes. All six recorded choices
reproduce at matched depth; four change with another depth, and two persist.
Exact recovered leaves expose a castling/rook-activity tie and coarse valuations,
without proving that every error has the same cause. The live sample informs
selection but does not override earlier research or establish an Elo gain.

| Attempt | Result | Decision |
| --- | --- | --- |
| [Defender-aware checking-pressure evaluation](../v0.23-defended-checks/README.md) | Public controls improve; newer 100ms mean loss 30.679→36.115 cp, both turns fail. | Reject before fresh controls or games; source and exact cleanup preserved. |
| [Child-state reuse on current kernel](../v0.23-child-reuse/README.md) | All 1,145 warm roots exact in three pairs; median elapsed saving 9.137% versus frozen 15% target. | Reject larger-gain attempt before native/games; old v0.21 failure remains separate. |
| [Checking-only generator and bounded tactical layer](../v0.23-quiet-checks/README.md) | Exact generator controls pass; about 78% faster generation. All native gates pass; games score 50% / 50% / 77%. | Reject: the direct-release game gate fails. Preserve all 150 games. |
| [Phase-sensitive activity](../v0.23-phase-activity/README.md) | Older long-budget Black gate fails; a prospective stricter fresh control also fails, 37.30→40.98 cp loss at 100ms. | Reject before games; no changed coefficients or alternative final seed. |
| [Guarded pass probe with real-move verification](../v0.23-verified-null/README.md) | 9.89% fewer nodes but median 0.110% slower over three warm pairs. | Reject on both frozen cost targets, before native/games. |
| [Contextual pawn structure](../v0.23-pawn-structure/README.md) | Exact compaction lowers overhead to 3.95%; public 100ms loss worsens 45.706→57.881 cp. | Reject before later native suites, fresh collection or games. |
| [Capture-futility opportunity](../v0.23-delta-captures/README.md) | All 1,145 roots exact; eligible outer subtrees cover 4.689% of nodes and already end after one child. | Reject diagnosis before any playing pruning policy. |
| [Full-opponent-ply null probe](../v0.23-null-probe/README.md) | Full controls/replay pass; first warm pair saves only 5.01% nodes. | Shared runner rejects early node gate; no extra pairs/native/games. |
| [Learned piece–king context](../v0.23-king-context/PLAN.md) | Converged fit worsens validation 36.095→37.462 and historical holdout 44.043→56.179 cp, severe errors 3→6. | Reject without runtime activation; numerical failures and repaired fit preserved. |
| [Direct-leaf target alignment](../v0.23-direct-leaf-context/README.md) | 85 queries; median target difference 45.5cp and 7.971% clear-preference reversals both miss frozen opportunity gates. | Reject diagnosis before remaining collection or another fit. |
| [Root-parallel scheduling opportunity](../v0.23-parallel-opportunity/README.md) | All 1,145 roots exact; two-lane serial-first proxy saves 38.193% of node span. | Implement bounded root-only scheduling; proxy is neither elapsed nor strength evidence. |
| [Two-worker root search](../v0.23-root-parallel/README.md) | 25.86% median elapsed saving; fresh 1000ms fails. The one complete screen scores 50% direct / 47% Skill 3 / 80% Sacrifice. | Reject direct gate; preserve the pre-game exception, failed surrogate and all 150 legal games. |
| [Root PVS scout](../v0.23-root-scout/README.md) | All 1,145 values/choices match; only 4.123% fewer nodes. | Reject at the early cost gate; no later pairs/native/games. |
| [Virtual ordinary-move tags](../v0.23-virtual-tags/README.md) | Exact rules/tree controls; 4.543% median elapsed saving. | Reject larger-gain cost gate; no native/games. |
| [Losing-capture stages](../v0.23-capture-stages/README.md) | Exact warm values/choices; 5.944% fewer nodes. | Reject early cost gate; preserve allocation/control repairs. |
| [Mate-distance bounds](../v0.23-mate-bounds/README.md) | Exact warm values/choices; 3.386% fewer nodes. | Reject early cost gate; no native/games. |
| [Child-scout reduction opportunity](../v0.23-scout-reduction-study/README.md) | Unchanged-tree eligible outer scouts cover 66.436%, including wide-window parents. | Pass diagnosis; old narrow-parent diagnosis remains failed. |
| [Initial child-scout reductions](../v0.23-scout-reductions/README.md) | Independent low-material forced-mate control fails; an accidentally started cost run is interrupted. | Reject; preserve all partial records and repair sequencing tooling. |
| [Low-material-protected scout reductions](../v0.23-scout-material/README.md) | 45.448% nominal-depth elapsed saving, full-depth regret 1.014cp; four existing suites pass. Both fresh gates fail. | Reject; preserve the conditional scope decision, one-million-node mate confirmation and matched-depth diagnosis. No games. |
| [Last-two-plies protection](../v0.23-scout-frontier/README.md) | Rules/race/full-depth-four controls pass; first pair saves 12.128% nodes. | Reject early cost target; no native/games. |
| [Quiet-count scout reductions](../v0.23-quiet-scouts/README.md) | Controls and all six choice gates pass; 43.128% elapsed saving. Sole screen: 57% direct /45% Skill 3 /77% Sacrifice, all 150 legal without forfeits. | Revise: direct gain does not offset the adverse external signal. No unchanged rerun. |
| [Two-choice checking budget](../v0.23-check-budget/README.md) | Correctness/race pass; median elapsed ratio1.363142 exceeds1.25. | Reject before native/reference/fresh/games; 45-file archive sealed. |
| [One-choice checking budget](../v0.23-check-budget-one/README.md) | Cost/reference/existing/fresh100 pass; fresh1000 Black19.677→26.462cp fails. | Reject without games/exception;142-root cohort consumed. |
| [Static-margin diagnosis](../v0.23-static-cut-opportunity/README.md) | 20.808% opportunity;179 false predictions include one hidden losing mate. | Reject before implementation. |
| [Quiescence-verified static diagnosis](../v0.23-verified-static-opportunity/README.md) | Same hidden mate; stand-pat makes the proposed verification ineffective. | Reject and close this line; design mistake preserved. |
| [Coordinated exact state kernel](../v0.23-state-kernel/README.md) | Exact1,145 warm roots/43.1Mnodes,14.142% less time; all six choice gates pass. Sole screen51% direct /69.4% completed Skill3 /76% Sacrifice. | Recommend moderate functional promotion; strength inconclusive. One capped Stockfish game; all150 legal, no forfeits. Stop before commit. |

The [reassessment after evaluation](REASSESSMENT_AFTER_EVALUATION.md) explicitly
returns to pre-live recommendations. The current warm profile retains all ten
openings and 1,145 recorded roots from the earlier study. It identifies remaining
child-state allocation, but the new timing pairs again show that allocated bytes
are not proportional to elapsed headroom. Earlier clock, aspiration/history,
reuse, static safety bonuses and fitted models are not silently reenabled.

The third attempt follows the earlier checking-interface recommendation with
a new generation-cost boundary and one selective horizon layer. It does not
reactivate the old failed standalone mate proof, add unbounded quiet checks,
change material or relax draw/cache identity. The independent complete-list
oracle was updated for the new horizon; its original capture-only failures
remain in the candidate's control notes and test output.

Every rejected source, method, adverse gate, raw query and completed measurement
stays on this branch. Compiled binaries remain ignored. Fresh controls have a
[frozen collection plan](FRESH_CONTROL_PLAN.md); they are not silently recycled
as hidden holdouts after inspection. Standard strength screens, when earned,
remain 50 games each versus tagged v0.22, Skill 3 and pinned Sacrifice, concurrent.
No live Lichess play or automatic larger confirmation is part of this research.

The [process audit](../v0.23-process-audit/README.md) replaces copied drivers,
manual report arithmetic and module-overlay repair with shared Go commands.
Its [tooling flowchart](../../../../docs/RESEARCH_TOOLING.md) and the compact
[checkpoint](../../../../docs/RESEARCH_CHECKPOINT.md) guide future contexts.
Original observations and frozen gates remain unchanged. The audit recomputes
three preserved cost decisions without rerunning engines and retains its own
initial resume-identity bug and corrected regression tests.

After the first full screen failed, the fourth attempt returned to the older
phase-dependent activity recommendation. Its fresh failure then motivated a
separate selective-search attempt. The old null-move screen remains failed;
two full pass-probe plies, probe-only checking visibility, material/pawn guards
and isolated real-move verification define the new experiment. These safeguards
are hypotheses about risk and cost, not a retrospective diagnosis of the old
losses. Neither normal quiet-check expansion nor the phase tables remain active.
