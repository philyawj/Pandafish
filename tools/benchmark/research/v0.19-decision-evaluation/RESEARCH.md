# My next priorities after v0.18

I released v0.18 through [PR #33](https://github.com/philyawj/Pandafish/pull/33)
and [its tag](https://github.com/philyawj/Pandafish/releases/tag/v0.18). The
search-core, legal-move/storage and exact-evaluation projects produced useful
engineering gains. The fitted evaluation model improved static error but failed
short-budget decisions twice. I therefore reassess the old architecture plan;
these are ranked experiments, not a fixed checklist or promised strength gains.

| Rank | Next experiment | Why I rank it here |
| --- | --- | --- |
| 1 | Decision-aligned compact evaluation | I will collect independent families and actual response/quiet-leaf alternatives, then select a bounded model on decision regret instead of static MAE alone. This addresses the observed failure directly and can use the cheaper v0.18 kernel. |
| 2 | Profile-driven legal tagging/attack summaries | The backend still spends work tagging/checking moves. I would profile the new release before choosing one additive API change; repeated information may be shareable without a board rewrite. |
| 3 | Compact neural representation using the new data | The small classical failure does not rule out NNUE. I need useful independent decision data, representation and inference budgets before choosing a network or incremental state. |
| 4 | Verified selective search | The efficient core makes careful reductions/history/aspiration worth reconsidering, but the earlier bundle failed. I need an isolated mechanism and tactical/zugzwang controls, not a blind retry. |
| 5 | Internal make/unmake board | I retain this as an option only if profiling identifies state updates as a large remaining cost. It still carries the broadest rule/hash/undo validation burden. |

I proceed with rank 1. I will reuse the corrected, independently checked
feature geometry as a bounded representation, not reuse the rejected weights.
The new experiment will compare alternative branches at quiet student leaves,
freeze family splits and fitting/gates before labels are inspected, and retain
multiple acceptable moves. A fresh holdout and actual equal-time search gate
must survive before any standard 50-game screen. Material, search, tables and
clocks remain v0.18. No book or Lichess dataset is needed for this collection.

I keep the 57% / 73% / 51% v0.18 screen and its opening limits intact. A broader
confirmation remains justified before a large strength claim or advancing the
Stockfish target; it has not run. I will reassess this ranking after the new
result and check in before its PR. No live games are authorized.
