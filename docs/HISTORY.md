# Pandafish version history

Earlier milestones and their original research links are preserved below.
The [v1.0 milestone and development pause](../README.md#milestones-and-plans)
mark the current project state. Benchmark evidence remains in the linked reports.


**v0.1: Get games working (complete).** This release built a minimal Go engine
that speaks UCI and can complete games on Lichess through lichess-bot. Random
legal moves were enough to prove that connection worked. GoPandaFish's first
successful games included:

- [Casual correspondence / unlimited game](https://lichess.org/p3w5nk1V)
- [Casual timed game: 3+0 blitz](https://lichess.org/wb007ifr)

Both games reached checkmate. This milestone was about getting the engine
playing complete games, regardless of the result.

**v0.2: Choose moves for a reason (complete).** Pandafish now evaluates the
immediate material balance after each legal move, chooses the move with the
best score from the moving side's perspective, and breaks ties randomly. This
is a one-ply material evaluation. It does not consider the opponent's reply,
checkmate, or positional features.

**v0.3: Look ahead one reply (complete).** Depth-two minimax now accounts
for the opponent's best immediate material response. The v0.2 piece values and
random tie breaking are unchanged.

**v0.4: Know when the game is over (complete).** The evaluator recognizes
checkmate and stalemate, including an opponent's immediate mating reply.
Depth-two search, material values, and random tie breaking are unchanged.
Repetition and fifty-move draw detection inside the search remain future work.

**v0.5: Build a deeper search (complete).** Replace the fixed two-ply
calculation with recursive minimax using alpha-beta pruning at three plies.
The v0.4 material evaluator, mate/stalemate handling, and random tie breaking
remain unchanged. This release retained fixed-depth search and ignored
clock/depth limits.

**v0.6: Finish the exchange before counting the pieces (complete).**
Add bounded quiescence search for captures, promotions, and check evasions
beyond the three-ply boundary. Material values, mate/stalemate scoring, and
random tie breaking are unchanged. This release still used a fixed main
search depth and ignored clock/depth limits.

**v0.7: Know when to stop thinking (complete).** Honor UCI time and
depth limits with conservative clock allocation, iterative deepening, and
interruptible search. Keep the last completed iteration and a legal fallback.
Material scoring, mate/stalemate handling, and the two-ply quiescence cap are
unchanged.

**v0.8: Record evidence from real games (complete).** Add optional
engine traces and Go `record`/`report` tools for local Lichess sessions.
Keep the existing search and challenge policy. See the
[session tool documentation](../tools/lichess/README.md) for operation and limits.

**v0.9: Give quiet moves positional guidance (complete).** Add small
placement terms alongside material, with a phase-dependent king score.
Keep the existing alpha-beta search, timing rules, quiescence cap, and random
choice among equally best moves. Placement helps distinguish quiet choices;
it does not guarantee sensible opening play in every position.

**v0.10: Reuse searched positions (complete).** Cache completed main-search
results in a fixed-size transposition table and try cached moves first.
Keep v0.9's evaluation, time control, quiescence cap, and exact root ties.
The in-memory table is private to each process, reused during a game, and
cleared on `ucinewgame`. No learned data is saved to disk.

**v0.11: Resolve checks at the quiescence horizon (complete).** Keep the
two-ply capture/promotion horizon, but search required check evasions beyond it.
Allow at most four extra check-evasion plies, then use static evaluation at the
final safety cap. A non-checked position at the ordinary horizon still stops.
Evaluation, clocks, transposition caching, and random root ties are unchanged.

**v0.12: Reuse quiet cutoff moves (complete).** Keep two killer-move
hints per search ply. Try matching legal quiet moves early after the cached
move and tactical priorities. Reset hints for each `go` command; iterative
passes share them. Evaluation, search horizons, time allocation, and exact root
ties remain unchanged.

**v0.13: Settle longer exchanges (complete).** Continue selective
quiescence beyond the old two-ply boundary, with a 16-ply ordinary safety bound
and the existing four extra check-evasion plies. Estimate material exchanges
to skip losing non-checking captures. Keep checking captures, promotions,
and evasions, and preserve evaluation, main search, clocks, and random root ties.

**v0.14: Speed up the development benchmark workflow (complete).** Run the
three standard opponent matchups concurrently with smaller screening samples.
Keep serial mode available for finer comparisons later. Search, evaluation,
time allocation, and move selection are unchanged from v0.13. See the
[benchmark workflow](../tools/benchmark/README.md) for settings and evidence.

**v0.14.1: Preserve an unpromoted research experiment (complete).** Archive
selective-search code, validation records, and the decision to keep the released
search policy. See the [research report](../tools/benchmark/research/v0.15-selective-search/README.md)
for the experiment and the proposed next direction. Playing behavior is unchanged
from v0.14.

**v0.14.2: Preserve the positional-evaluation experiment (complete).** Archive
board-aware evaluation code, tests, diagnostics, and match records. Keep the
released playing behavior unchanged. See the
[research report](../tools/benchmark/research/v0.15-positional-evaluation/README.md)
for findings and replay instructions.

**v0.14.3: Preserve the null-move search experiment (complete).** Archive the
reference-engine comparison, guarded null-move code, tests, diagnostics, and
match records. Keep the released playing behavior unchanged. See the
[research report](../tools/benchmark/research/v0.15-null-move/README.md).

**v0.14.4: Preserve the king-safety evaluation experiment (complete).** Archive
the ranked positional hypotheses, exact code, tests, saved-decision diagnostics,
and match records. Keep released playing behavior unchanged. See the
[research report](../tools/benchmark/research/v0.15-king-danger/README.md).

**v0.15: Convert wins and recognize repetitions (complete).**
Prefer shorter forced mates, carry mate distance through quiescence and cached
scores, and evaluate threefold repetitions from known game/search history.
Report calculated signed mate counts through UCI. Keep the released positional
evaluation, search breadth and clocks. See the
[conversion diagnosis and validation](../tools/benchmark/research/v0.15-mate-conversion/README.md).

**v0.15.1: Review the extended live test (complete).** Preserve the session,
full-strength offline analysis, diagnostic replays and next-step recommendations.
Repair partial search/game linkage in the reporting tool. Engine playing
behavior remains v0.15. See the
[live research and handoff](../tools/benchmark/live/v0.15-live-20261113-074740/README.md).

**v0.15.2: Preserve the piece-safety experiment (complete).** Archive the
isolated candidate, quiet-leaf diagnosis, positive and negative controls, and
all failed and completed probes. The experiment is rejected; released playing
behavior remains v0.15. See the
[piece-safety research report](../tools/benchmark/research/v0.16-piece-safety/README.md).

**v0.15.3: Investigate larger strength opportunities (complete).** Preserve the
unpromoted bishop-pair source and games, cache/repetition controls, search
ablations, CPU/allocation profiles and a broader architecture assessment.
Released playing behavior remains v0.15. See the
[research and measured findings](../tools/benchmark/research/v0.15.3-architecture/README.md).

**v0.16: Coordinate the search core (complete).** Return one proven best move,
use full-depth PVS, recover draw-safe move hints and reversible cache contexts,
and sort only eligible tactical moves outside check. Preserve evaluation,
quiescence eligibility, mate/repetition correctness and time allocation. See the
[implementation, validation and measured limits](../tools/benchmark/research/v0.16-search-core/README.md).

**v0.17: Generate legal tactics with reusable storage (complete).** Add a
search API to the pinned chess library, generate captures/promotions directly,
check legal-move existence without constructing a list, and reuse move buffers
by active ply. Preserve fixed-depth scores, chosen moves and node counts,
with the original library methods and license intact. See the
[implementation, validation and measured limits](../tools/benchmark/research/v0.17-search-moves/README.md).

**v0.18: Compute the same evaluation with less work (complete).** Precompute
original material/placement formulas and evaluate value-copied piece bitboards.
Preserve exact scores and fixed-depth search behavior. Archive the unpromoted
fitted evaluator and its failed gates separately. See the
[implementation, tests and measured limits](../tools/benchmark/research/v0.18-exact-evaluation/README.md).

**v0.19: Preserve evaluation research and Go-only tooling (complete).** Archive
three rejected post-v0.18 fits, preserve their findings and retire non-Go
research runners with original evidence retained in the release attachment.
Add shared Go research checks and tracked agent guidance. Playing policy remains
v0.18. See the [research summary and next steps](RESEARCH_V019.md).

**v0.20: Remove repeated search-legality work (complete).** Specialize immutable
child transitions, compute attack rays once and skip opponent-check tags when
only legal-move existence is needed. Preserve fixed-depth scores, chosen moves,
node counts and search policy. Keep the unsuccessful quiet-mate and transition-
only experiments archived. See the
[implementation, validation and measured limits](../tools/benchmark/research/v0.20-search-legality/README.md).

**v0.21: Compute stable move order with less work (complete).** Calculate each
existing priority once, sort small rank records and apply their permutation to
the original moves. Preserve exact ordering, fixed-depth scores/moves/nodes,
values and search policy. See the
[implementation, validation and measured limits](../tools/benchmark/research/v0.21-rank-order/README.md).

**v0.22: Avoid unused legal work and repeated exchange scans (complete).**
Generate later move stages only when needed, prefilter impossible check evasions
and compute the existing legal exchange heuristic with local bitsets. Preserve
complete move order, fixed-work scores/moves/nodes, values and search policy.
Keep unsuccessful experiments and all original results archived. See the
[implementation, validation and measured limits](../tools/benchmark/research/v0.22-bitboard-exchange/README.md).

**v0.23: Reuse search state and avoid temporary tag copies (complete).** Reuse
worker-owned child positions and boards by search ply, and classify ordinary
moves from virtual occupancy and attacker sets. Preserve fixed-work moves,
scores, nodes and playing policy. Archive unsuccessful research and add a
repeatable Go research workflow. See the
[implementation and measured limits](../tools/benchmark/research/v0.23-state-kernel/README.md).


**v1.0: Stable engine milestone.** Promote bounded queen-contact threat defense,
preserve all research and document the development pause. See the
[v1.0 research and evidence](RESEARCH_V024.md) and
[release validation](../tools/research/v1.0-validation.md).
