# Two-worker root search: reject after the complete screen

The [plan](PLAN.md) follows the released-tree scheduling diagnosis. Search the
first root move serially, then distribute remaining root moves across two private
workers sharing synchronized fixed-size score/hint tables. Only an exact strict
improvement can replace the incumbent; a cutoff bound cannot supply a false tie.
Every root joins its workers before returning. Repetition counts, history keys,
killers and reusable per-ply buffers are private to each worker.

Material, evaluation, ordinary child search, quiescence eligibility/horizons,
score-cache identity/depth/replacement and clock budgets remain v0.22. Scheduling
can change move order, tied choices and node counts. This is two-core wall-time
work, not a per-core efficiency improvement. Hosts with GOMAXPROCS below two use
the serial path; that fallback's locking overhead is not separately benchmarked.
There is no UCI Threads option or new live-play setting.

Normal/race, mate, repetition-context, parent-state, legal co-best, cancellation,
worker-join and team-reuse controls pass. The [explicit warm protocol](PARALLEL_PROTOCOL.md)
uses three alternating depth-five pairs, all ten openings and 1,145 actual game
roots with persistent per-game tables. Baseline nodes stay 43,118,858; candidate
nodes vary 44,867,189–44,945,632. Search elapsed ratios are 0.74600, 0.73091 and
0.74137: median **25.86% saving**. All baseline values match; post-timing serial
full-window verification independently validates all recorded selected moves.

| Equal-time suite | Mean finite loss, v0.22 → candidate | Result |
| --- | ---: | --- |
| Public 100ms, three trials | 45.8955 → 45.5970cp | Pass |
| Public 1000ms | 54.8507 → 54.8507cp | Pass |
| Newer 100ms | 30.6794 → 30.3282cp | Pass |
| Newer 1000ms | 24.2424 → 22.0682cp | Pass |
| Fresh 100ms, three trials | 27.6691 → 27.6691cp | Pass |
| Fresh 1000ms | 27.3971 → 32.9485cp | **Fail**, severe errors 2→3 |

The fresh cohort uses frozen seed 2026112113/stream113, all 16 families, 150
retained roots (75 per turn), and explicit prior-control/live-board exclusions.
No new losing-mate flags or timing overruns occur in any suite. The failed gate
is not changed. The [screen decision](SCREEN_DECISION.md) explicitly excepts it
for one research screen after all 60 matched-depth comparisons on 30 changed
positions pass. Those comparisons include improved/neutral choices, preserve
full game history and never force a production root move. Reaching a deeper
pre-existing bad horizon explains the inspected regressions; it does not make
the moves good or establish universal correctness or strength.

The standard concurrent screen completes all 150 games:

| Opponent | W / D / L | Score |
| --- | --- | ---: |
| Tagged v0.22 | 17 / 16 / 17 | **50%**, direct gate fails |
| Stockfish Skill 3 | 23 / 1 / 26 | 47% |
| Pinned Sacrifice | 38 / 4 / 8 | 80% |

Every game finishes legally, without forfeits. The independent legal audit runs
after the runner stops on its failed direct criterion. The earlier exception
only authorizes evidence gathering: it does not erase the fresh failure or
change the direct criterion. The mixed external rows and extra CPU use do not
justify promotion. Reject this version; no unchanged rerun or larger confirmation.
The source snapshots preserve the implementation for future research; restore
released single-worker playing code before trying a different mechanism.

[Infrastructure notes](INFRASTRUCTURE.md) preserve the source-replay cache-trim
failure and the pre-query zero-ply analysis-record failure. Original/corrected
plans, source snapshots, raw observations, step logs and identities remain here.
The production binary is unchanged across the collection/tooling repairs.
