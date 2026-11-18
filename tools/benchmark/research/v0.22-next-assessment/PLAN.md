# Next assessment after v0.21

Baseline: released v0.21, commit
`494604686a7afbda256da1625aa907baa6060084`, merged through PR #38. The rank
candidate's exact-work gain and clean 55%/50%/71% screen support its approved
promotion; no precise Elo gain is established. Skill 3 remains the routine
reference. Current research starts on `research/v0.22-next-candidate` and stops
at the next recommended candidate before a commit.

## Profile before selecting an implementation

The prior fixed-depth benchmark creates a fresh TT for every root. Production
keeps it within a game and reconstructs the supplied position/history before
each `go`. Measure both a search-only warm-table workload and the position-plus-
search pipeline before deciding that heap allocation, board copying or another
kernel is the best target. No source change to the playing engine precedes the
baseline observations.

Freeze the first 20 games (ten color pairs / all ten openings) from the v0.21
rank candidate's `previous.pgn` and matching summary. For each game, visit every
candidate turn after its supplied opening through the last played move. Preserve
the original legal trajectory regardless of the profiled engine's selected move.
These are offline recorded positions, not new games or training labels.

Use the production 100ms allowance and budget calculation, iterative deepening,
one new searcher per root and one persistent TT per game. The search-only mode
uses recorded immutable positions and equivalent repetition histories; pipeline
mode reparses the full supplied startpos history exactly as the UCI handler does.
Record every root's parsing/search times, nodes, depth, result, allocation deltas
and input identities. Profiles include measurement/setup overhead and do not
measure pipe input/output or simultaneous opponent load. Any timeout/failure is retained;
do not substitute games selected for better-looking costs.

The initial ranking is conditional: (1) reduce measured child-state allocation/
copying without make/unmake if it dominates the warm workload; (2) reduce repeated
position reconstruction or remaining tag/exchange work if the profile makes a
stronger case; (3) a different evaluation representation with new independent
data and actual-budget selection. Neither old fitted weights nor the failed
quiet-mate proof is reactivated. Reassess after every result.

## Candidate gates

Freeze a precise mechanism before editing runtime code. For an exact-kernel
candidate, preserve values, eligibility/order, search/TT policy, repetition,
draw identity and clocks. Require ownership/lifecycle, cold/warm, special-move,
history, cancellation, independent/original oracle and race controls appropriate
to that mechanism. All source, fixtures and research programs remain Go.

Retain scores/moves/nodes on the 325 existing public/deep/newer position/depth
comparisons, and on a fixed-depth warm sequence using persistent game tables.
Require at least 15% median elapsed reduction in three alternating depth-four
pairs versus v0.21, no contradictory deep result and a beneficial warm-workload
result before native probes. A narrower gain is archived or reassessed, rather
than automatically promoted from a microbenchmark.

Use the existing four native suites: public 84 roots at three 100ms trials and
one 1000ms trial, newer inspected 145 roots at one trial per budget, three native
searches at once, native work completed before teacher queries. Keep every
observation. Teacher: pinned full-strength Stockfish, 100,000 nodes, Threads 1,
Hash 64 MiB, MultiPV 1; cache reuse requires matching hashes/budgets. Production
roots stay unrestricted and close alternatives are acceptable.

On fully finite paired cases use the trial-mean
`clamp(root_reference_CP-selected_reference_CP-25,0,1000)`. Require mean loss at
most baseline plus max(5cp,10%), severe (>300) cases at most baseline plus two,
no new losing-mate flag under the existing finite/nonlosing-root convention,
and no response past allowance plus 50ms. Apply each newer turn's mean gate too.
Inspected inputs are regression controls, not a blind holdout. Finite labels do
not prove a position was savable; the reference is not a unique acceptable move.

A surviving mechanism reaches **50 games each versus tagged v0.21, Skill 3 and
pinned Sacrifice, all three concurrent**. Require the established previous-
release practical gate and review both externals, timings, forfeits and unfinished
games. Retain failed code as inactive research and continue to a substantive next
attempt. No favorable unchanged rerun, larger confirmation, live game, commit,
push, PR or release follows automatically. Check in at the next viable candidate.
