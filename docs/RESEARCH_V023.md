# Research toward v0.23

The investigation starts from tagged v0.22 after its separately merged live
review. [The assessment index](../tools/benchmark/research/v0.23-next-assessment/README.md)
links every plan, source snapshot, raw result, failure and decision.
Twenty-three substantive attempts are unsuccessful: eighteen engine prototypes,
one offline fitted model and four diagnoses. Intermediate solver, oracle,
compile and infrastructure failures remain within their respective attempts.

| Unsuccessful attempt | Decisive result |
| --- | --- |
| Defender-aware checking pressure | Newer short-budget mean loss 30.679→36.115cp. |
| Child-state reuse | Exact results, 9.137% median saving misses the frozen 15% objective. |
| Checking-only generator and one quiet-check layer | All decision gates pass; sole screen 50% direct / 50% Skill 3 / 77% Sacrifice fails direct. |
| Phase-sensitive activity | Long-budget Black and prospective fresh-improvement gates fail. |
| Verified null move | 9.89% fewer nodes, 0.110% slower; cost targets fail. |
| Contextual pawn structure | Public loss 45.706→57.881cp despite exact compaction. |
| Capture-futility diagnosis | Only 4.689% eligible outer node opportunity. |
| Full-opponent-ply null probe | Only 5.01% fewer nodes; early stop. |
| Offline learned piece–king context | Validation/holdout worsen; holdout severe errors 3→6; no runtime activation. |
| Direct-leaf target diagnosis | Both signal gates fail after the fixed 85-query study. |
| Two-worker root search | 25.86% wall-time saving with extra CPU; fresh gate fails; sole screen 50% / 47% / 80% fails direct. |
| Root PVS scout | Exact values/choices; only 4.123% fewer nodes. |
| Virtual ordinary-move tags | Exact rules/tree; only 4.543% median elapsed saving. |
| Losing-capture stages | Exact values/choices; only 5.944% fewer nodes. |
| Mate-distance bounds | Exact values/choices; only 3.386% fewer nodes. |
| Initial late-child reductions | Independent sparse-mate control fails; accidentally started timing is interrupted and preserved. |
| Material-guarded reductions | 45.448% nominal-depth elapsed saving and low regret; both fresh gates fail. Conditional remaining-control decision does not lead to games. |
| Last-two-main-plies protection | Correctness/full-depth-four controls pass; only 12.128% fewer nodes. |
| Quiet-count scouts | All six choice gates and direct screen pass at 57%; adverse 45% Skill 3 row prevents recommendation. |
| Two-choice checking budget | Correctness/race pass; 36.314% nominal-depth elapsed increase exceeds25% budget. No native/games. |
| One-choice checking budget | 28.574% elapsed saving, low same-horizon regret, existing/fresh100 pass; fresh1000 Black19.677→26.462cp fails. No games. |
| Static-margin diagnosis | 20.808% outer work opportunity, but179 false predictions include one concealed losing mate. No pruning activated. |
| Quiescence-verified static diagnosis | Same hidden mate remains: stand-pat makes every proposed verification immediately accept. Design mistake preserved; line closed. |

The first three screens retain all 450 legal, finished games, without forfeits.
No unfavorable outcome is rerun. Successful opportunity studies are separate:
root scheduling estimates 38.193% node-span saving, while late child scouts cover
66.436% of unchanged warm nodes. Those proxies are not elapsed or strength results.

The older reduction diagnosis guarded on the original parent's narrow window.
Later children of wide-window parents also receive individual one-centipawn PVS
scouts; measuring that boundary reveals substantially more work. This changes the
selective risk and does not retroactively pass the old guarded diagnosis.

The [quiet-scout foundation](../tools/benchmark/research/v0.23-quiet-scouts/README.md)
counts previously searched ordinary quiet choices rather than all moves.
Captures/promotions/checks cannot make the first quiet appear late. First three
ordinary quiets, castles, checks, tactics, killer hints, mate windows and sparse
material stay full depth. Improving or mate-reporting reduced scouts receive
full-depth verification before normal PVS re-search. Evaluation, quiescence,
cache identity/depth/replacement and clock formulas remain released v0.22.

Normal/race/vet, independent mate and capture-before-quiet controls pass.
Three alternating 1,145-root nominal-depth-five pairs save 43.128% elapsed and
47.122% nodes. Selective search changes fixed-depth values; equivalent full-width
depth is not claimed. Original-tag full-depth review of every selected choice
reports mean finite regret 0.996cp over 1,103 roots, maximum 273cp, with zero
mate-risk flags or value contradictions. All four existing equal-time suites
pass, as do both fresh suites (seed2026112121/stream121, 130 roots). Its sole
screen scores 57% direct /45% Skill 3 /77% Sacrifice, with all 150 legal and no
forfeits. The previous v0.22 screen recorded 71% against Skill 3. These small
separate samples do not isolate Elo, but the adverse external signal outweighs
a promotion recommendation. No unchanged screen is repeated.

The deeper [coupled hypothesis](../tools/benchmark/research/v0.23-check-budget/README.md)
spends late-quiet savings on checks in the first three q plies but fails cost.
The [narrower revision](../tools/benchmark/research/v0.23-check-budget-one/README.md)
passes cost and same-horizon reference checks, but its failed fresh long-budget
Black control stops it without games or an exception. Neither component is
independently promoted. The two subsequent static-margin diagnoses stop before
implementation: substantial work opportunity does not establish tactical safety.
The attempted q verification is ineffective because stand-pat already accepts
the qualifying static score; that design mistake is recorded explicitly.

The [current exact state kernel](../tools/benchmark/research/v0.23-state-kernel/PLAN.md)
combines per-ply child storage and virtual ordinary-move tags as one state-handling
project. Its earlier components' separate15% cost gates stay failed. This new
combined candidate prospectively targets12% saving after reassessing the cost
of rejecting moderate exact improvements without actual-budget evidence.
Normal/race/vet pass in both modules;596 independent vectors exercise actual
scratch children. Three alternating warm pairs retain all1,145 values, moves
and43,118,858nodes exactly and save14.142% elapsed. All six native/fresh gates
pass with no exceptions/new mate flags/overruns. Fresh seed23/stream123 retains
148 roots,74 per turn; finite loss improves42.325→40.184cp at100ms and
30.630→28.197cp at1000ms. Its sole concurrent screen completes:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.22 | 20 / 11 / 19 | 0 | 0 | 51% |
| Stockfish Skill3 | 34 / 0 / 15 | 1 | 0 | 69.4% |
| Pinned Sacrifice | 34 / 8 / 8 | 0 | 0 | 76% |

The Stockfish game reaches the400-ply cap, giving a68–70% possible50-game score;
it is not a shutdown/timeout. All150 records replay legally,149 finish, no
forfeits; candidate responses remain below100ms. Across this investigation's
four screens,600 records are preserved,599 finished and one capped.
**Recommend the exact-state kernel for its moderate functional efficiency gain,
with playing-strength improvement inconclusive.** The practical direct gate
passes by one game; repeated openings and candidate selection limit inference.
The external rows do not show the earlier selective candidate's major adverse
signal, but comparison with v0.22's separate71%/70% screen is not paired proof.
Broader-opening confirmation is justified before claiming an Elo/strength gain
or increasing Skill3; no larger run starts automatically. Stop before commit.

A predecessor's short-budget mate flag persists at one million reference nodes.
At matched depth, both engines prefer its bad move at six and the more resistant
move at five. This exposes a pre-existing quiet-check horizon, without erasing
the equal-time disadvantage or proving a finite reference alternative saveable.
Its original failed gates and explicit scope reassessment remain. No FEN/root
exception is added. The separate long-budget failure stops that predecessor.

Live mistakes now serve as regression controls rather than the candidate agenda.
The repeated failures show a plateau in the current research approach, not a
proven engine-strength ceiling. Most attempts stop before games; there have not
been23 full strength comparisons. Cost-size failures also do not imply that a
prototype has no functional value. Parallel/search/backend work revisits earlier
controlled recommendations rather than arising solely from live losses.
Inspected
controls remain development data; final cohorts have fixed new seeds and
canonical exclusions. Comparable alternatives and co-best moves are accepted.
No opening book, online move dataset or live game is introduced.

The next ranking is conditional on review of this candidate, not an automatic
new experiment: first profile the accepted kernel on a fixed production-style
workload, separating sampling from expensive memory-stat reads, and pursue a
larger state/undo project only if remaining headroom supports it. A genuinely
different contextual evaluation representation remains a larger research option
with a separately budgeted data/validation plan; the failed linear fits stay
inactive. Further selectivity ranks behind a demonstrated cost-and-safety
opportunity rather than another margin/reduction sweep. Live losses remain
regression inputs. No book or online dataset is established as a prerequisite.

The [standard workflow](RESEARCH_WORKFLOW.md), [Go tooling flowchart](RESEARCH_TOOLING.md)
and [compact checkpoint](RESEARCH_CHECKPOINT.md) guide future contexts.
Validation/warm/native plan builders share one fail-closed runner with hashes,
logs, gates and resumability checks. Warm CLI requires passed source-matching
validation. Fixed larger-budget mate confirmation and matched-depth review
preserve old reports. The [efficiency log](../tools/benchmark/research/v0.23-next-assessment/EFFICIENCY_LOG.md)
records saved orchestration, stopped work, reusable proofs and the real
sequencing failure/repair. No per-step token measurement exists; upfront tool
development and remaining research reasoning are explicit costs.

Promotion, PR, merge and v0.23 publication are approved after the reviewed
checkpoint. [Release validation](../tools/research/v0.23-validation.md) records
the version-only checks separately. No larger confirmation, next experiment or
live session starts automatically.
