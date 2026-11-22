# Pandafish

A bare-bones chess engine written in Go. Pandafish v0.23 speaks UCI over
stdin/stdout and uses time-aware iterative deepening with alpha-beta pruning
and bounded quiescence. Standard chess only.

All repository program, test and research source is Go. Start with
[AGENTS.md](AGENTS.md) for project guidance. The
[v0.19 research summary](docs/RESEARCH_V019.md) explains the rejected evaluation
attempts, Go-only tooling migration and ranked next experiments.
The [v0.20 investigation](docs/RESEARCH_V020.md) records the promoted search-
legality kernel, rejected predecessors and reassessed next priorities.
The [v0.21 investigation](tools/benchmark/research/v0.21-rank-order/README.md)
records precomputed stable ordering and the next profile-driven priorities.
The [v0.22 investigation](docs/RESEARCH_V022.md) records deferred legal generation,
exact exchange computation and the preserved unsuccessful experiments.
The [v0.23 investigation](docs/RESEARCH_V023.md) records reusable search state,
virtual move tags, preserved unsuccessful attempts and research-tooling guidance.
The [current checkpoint](docs/RESEARCH_CHECKPOINT.md),
[repository audit](docs/REPOSITORY_AUDIT_V023.md),
[ranked next projects](docs/NEXT_ENGINE_ASSESSMENT_V023.md) and
[evidence-storage policy](docs/EVIDENCE_STORAGE.md) describe the post-release
work. Existing data stays in place; new diagnostics use leaner records.

The engine code is in `main.go`, `limits.go`, `iterative.go`, `search.go`,
`transposition.go`, `search_context.go`, `search_moves.go`, `main_move_picker.go`,
`killers.go`, `exchange.go`, `exchange_bits.go`, `search_position.go`,
`evaluation.go`, `mate.go`, and `repetition.go`.
A pinned local copy of [corentings/chess](https://github.com/corentings/chess)
with an additive search API handles chess rules,
including castling, en passant, and promotions. The
[local search-backend note](docs/CHESS_SEARCH_BACKEND.md) explains the v0.17
library replacement, its benefits, maintenance cost and future uses.
The engine assigns values of
1 to pawns, 3 to knights and bishops, 5 to rooks, and 9 to queens. It picks a
legal move by completing successively deeper searches, assuming each side
chooses its best continuation. It keeps the last fully completed result when
time runs out, with a legal fallback if no depth completes. Alpha-beta
pruning skips branches that cannot improve a search result, with captures
searched first to help pruning. A completed search returns one proven best
move with an exact score; it does not enumerate every tied move. Later moves
use a narrow PVS probe and receive a full-window search when needed. The
unsearched emergency fallback is selected randomly from legal root moves.
Checkmate scores above any material advantage, being
checkmated scores below any material disadvantage, and stalemate scores zero.
Shorter forced mates score higher, and longer forced losses score less badly.
Threefold repetition is scored as a draw using supplied game history and the
current search path. Other draw rules have no separate search evaluation. Modest placement scores
encourage central minor pieces, central pawn space, and active rooks on the
seventh rank. King placement blends home-rank safety with endgame activity
according to remaining non-pawn material. These are simple preferences, not
full development, pawn-structure, or king-safety analysis. Internal and UCI
scores use centipawns (100 per pawn).

Each engine process uses a fixed 65,536-entry transposition table (about
2 MiB), plus a separate fixed move-hint table (about 1 MiB). It caches completed
main-search scores at matching depths and tries stored moves first. Board-only
hints can order freshly validated legal moves across different draw contexts;
they never supply scores. Quiet moves that caused a beta cutoff can also be tried
early at the same search ply, after captures/promotions and any cached move.
Two such hints per ply are shared across iterative-deepening passes and reset
for each search. They affect ordering rather than evaluation.
The transposition table persists within a game, clears on `ucinewgame`,
and is never saved to disk. Quiescence scores are not cached. Search still
uses a repetition-count fingerprint to distinguish draw-sensitive score contexts.
After a pawn move or capture, unreachable earlier history is omitted from that
cache identity while draw counts remain intact.
It still has no separate fifty-move draw evaluation.

Quiescence means quietness. At the normal search boundary, it examines
captures and promotions, and all legal escapes when the king is in check.
Away from check it generates legal tactical moves directly into reusable
per-ply storage and can keep the current static score instead of making a
capture. This helps avoid counting a captured piece before considering the
recapture. A static exchange estimate skips losing ordinary captures outside
check, while checking captures, promotions, and required evasions remain eligible.
The estimate follows legal least-valuable recaptures on the destination square;
it is a material heuristic and can miss sacrifices with other tactical benefits.
Captures and promotions have a 16-ply safety bound. Required check evasions can
extend for at most four further plies. At that final safety cap, static evaluation
is used even if checks remain, so long tactical sequences can still be missed.

Clone this repository and build with Go 1.25 or newer:

```sh
git clone https://github.com/philyawj/Pandafish.git
cd Pandafish
go build -o pandafish .
go test ./...
```

Build on the machine that will run the engine. On Apple silicon, this produces
a native ARM64 executable. On Windows, use `go build -o pandafish.exe .`.

Try it:

```sh
./pandafish
```

Then enter:

```text
uci
isready
position startpos moves e2e4
go movetime 100
```

Wait for `bestmove` followed by the selected legal black move, then enter
`quit`. Sending `quit` or closing stdin immediately after `go` cancels the
search. When no legal move exists, the engine returns `bestmove 0000`.

The engine accepts `position startpos` or `position fen ...`, optionally
followed by `moves ...`. Search accepts `movetime`, `wtime`, `btime`, `winc`,
`binc`, `movestogo`, and `depth` (1 through 64). Times are milliseconds.
`searchmoves` restricts root choices and must come last. A depth-only request
has no time deadline; time and depth limits otherwise apply together.

Bare `go` uses a 100 ms allowance. Explicit `movetime` takes precedence over
clock allocation, capped to the usable side-to-move clock when supplied.
With game clocks, reserve the smaller of 10 ms or 10% of the remaining time, divide the
usable clock by `movestogo` (30 if absent), and add 75% of the increment.
Spending is capped to half the usable clock, so uncredited increment cannot
exhaust the current clock.

Inside the resulting allowance, reserve the smaller of 5 ms or 10% for
scheduling and communication. The search stops starting new iterations at
75% of that hard budget and checks cancellation inside both main search and
quiescence. Scheduling can still affect observed response time; the engine
makes no absolute real-time guarantee.

`isready` remains responsive during search. `stop` returns the last completed
result or legal fallback. Position changes, new games, and replacement
searches cancel old work. `go infinite` holds its result until `stop`.
`go ponder` also waits; `ponderhit` restarts a fresh timed search using its
original limits. Pondering remains disabled in the example configuration.
Invalid recognized search limits produce an `info string` error; unknown
commands and unsupported engine options are ignored.

To play through [lichess-bot](https://github.com/lichess-bot-devs/lichess-bot),
install that client and its Python dependencies following its setup instructions.
Use a Lichess BOT account and an
[API token with the Play bot moves permission](https://lichess.org/account/oauth/token/create?scopes[]=bot:play).
Pandafish itself needs no credentials; lichess-bot handles the connection.

Copy the supplied [example config](lichess-bot-config-example.yml) from this
repository's root:

```sh
cp lichess-bot-config-example.yml lichess-bot-config.yml
```

Edit **only your local copy**: replace `YOUR_LICHESS_BOT_API_TOKEN` and review the
three numbered engine lines: `dir` (binary directory), `name` (binary filename),
and `working_dir` (normally the same directory). Use absolute paths.
`lichess-bot-config.yml` is gitignored; the example must always contain placeholders.

The example accepts incoming **casual standard chess challenges from humans**,
with one active game at a time. It accepts all BOT-supported time controls,
including correspondence and unlimited games. Rated games and bot challengers
are declined; automatic outgoing challenges and pondering are disabled.
Lichess does not allow BOT accounts to play UltraBullet.

From your lichess-bot checkout, run the following on macOS/Linux, replacing the
config path with your own. This assumes its Python environment is named `venv`:

```sh
env -u LICHESS_BOT_TOKEN ./venv/bin/python lichess-bot.py --config /absolute/path/to/Pandafish/lichess-bot-config.yml
```

The `env` prefix prevents another account's environment token from overriding
the config. On Windows, use `venv\Scripts\python.exe` and clear any existing
`LICHESS_BOT_TOKEN` environment variable first. Stop the client with Ctrl+C.

The example keeps opening books, online moves, and tablebases disabled so every
move comes from Pandafish. Keep `uci_options` empty: the engine has no Hash,
Threads, or other configurable engine options. Casual play follows the guidance
in the [lichess-bot default config](https://github.com/lichess-bot-devs/lichess-bot/blob/master/config.yml.default)
for engines that do not try to win.

Contributions are welcome. Keep changes focused and the engine simple. Before
opening a pull request, format changed Go files and run:

```sh
go test ./...
go vet ./...
go build -o pandafish .
git diff --check
```

For engine changes, also check that lichess-bot can initialize the rebuilt
executable and obtain legal moves through its UCI wrapper. Its built-in config
validation checks configuration and paths; a move-sequence smoke test checks
engine integration. These local checks do not require a live Lichess game.
Include the checks you ran in your PR. Never commit tokens, private configs,
or compiled binaries. Before committing, inspect `git diff --cached`.

Pandafish is licensed under the [MIT License](LICENSE). The chess library's
license is reproduced in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
Include both notices when distributing compiled Pandafish binaries.

## Local benchmarks

[Developer tools](tools/README.md) describes the shared tools layout and the
[Lichess session recorder](tools/lichess/README.md). The recorder launches
lichess-bot with a fixed engine copy and reports on saved local records.
It also supports offline Stockfish processing.
[Recorded live sessions](tools/benchmark/live/README.md) preserve raw evidence,
analysis, and recommendations separately from controlled matches.

[Benchmark documentation](tools/benchmark/README.md) covers the validation workflow,
playing-strength comparisons, results, and their interpretation.
See the [release workflow](RELEASING.md) for preparing and publishing a version.

## Milestones and plans

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
[session tool documentation](tools/lichess/README.md) for operation and limits.

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
[benchmark workflow](tools/benchmark/README.md) for settings and evidence.

**v0.14.1: Preserve an unpromoted research experiment (complete).** Archive
selective-search code, validation records, and the decision to keep the released
search policy. See the [research report](tools/benchmark/research/v0.15-selective-search/README.md)
for the experiment and the proposed next direction. Playing behavior is unchanged
from v0.14.

**v0.14.2: Preserve the positional-evaluation experiment (complete).** Archive
board-aware evaluation code, tests, diagnostics, and match records. Keep the
released playing behavior unchanged. See the
[research report](tools/benchmark/research/v0.15-positional-evaluation/README.md)
for findings and replay instructions.

**v0.14.3: Preserve the null-move search experiment (complete).** Archive the
reference-engine comparison, guarded null-move code, tests, diagnostics, and
match records. Keep the released playing behavior unchanged. See the
[research report](tools/benchmark/research/v0.15-null-move/README.md).

**v0.14.4: Preserve the king-safety evaluation experiment (complete).** Archive
the ranked positional hypotheses, exact code, tests, saved-decision diagnostics,
and match records. Keep released playing behavior unchanged. See the
[research report](tools/benchmark/research/v0.15-king-danger/README.md).

**v0.15: Convert wins and recognize repetitions (complete).**
Prefer shorter forced mates, carry mate distance through quiescence and cached
scores, and evaluate threefold repetitions from known game/search history.
Report calculated signed mate counts through UCI. Keep the released positional
evaluation, search breadth and clocks. See the
[conversion diagnosis and validation](tools/benchmark/research/v0.15-mate-conversion/README.md).

**v0.15.1: Review the extended live test (complete).** Preserve the session,
full-strength offline analysis, diagnostic replays and next-step recommendations.
Repair partial search/game linkage in the reporting tool. Engine playing
behavior remains v0.15. See the
[live research and handoff](tools/benchmark/live/v0.15-live-20261113-074740/README.md).

**v0.15.2: Preserve the piece-safety experiment (complete).** Archive the
isolated candidate, quiet-leaf diagnosis, positive and negative controls, and
all failed and completed probes. The experiment is rejected; released playing
behavior remains v0.15. See the
[piece-safety research report](tools/benchmark/research/v0.16-piece-safety/README.md).

**v0.15.3: Investigate larger strength opportunities (complete).** Preserve the
unpromoted bishop-pair source and games, cache/repetition controls, search
ablations, CPU/allocation profiles and a broader architecture assessment.
Released playing behavior remains v0.15. See the
[research and measured findings](tools/benchmark/research/v0.15.3-architecture/README.md).

**v0.16: Coordinate the search core (complete).** Return one proven best move,
use full-depth PVS, recover draw-safe move hints and reversible cache contexts,
and sort only eligible tactical moves outside check. Preserve evaluation,
quiescence eligibility, mate/repetition correctness and time allocation. See the
[implementation, validation and measured limits](tools/benchmark/research/v0.16-search-core/README.md).

**v0.17: Generate legal tactics with reusable storage (complete).** Add a
search API to the pinned chess library, generate captures/promotions directly,
check legal-move existence without constructing a list, and reuse move buffers
by active ply. Preserve fixed-depth scores, chosen moves and node counts,
with the original library methods and license intact. See the
[implementation, validation and measured limits](tools/benchmark/research/v0.17-search-moves/README.md).

**v0.18: Compute the same evaluation with less work (complete).** Precompute
original material/placement formulas and evaluate value-copied piece bitboards.
Preserve exact scores and fixed-depth search behavior. Archive the unpromoted
fitted evaluator and its failed gates separately. See the
[implementation, tests and measured limits](tools/benchmark/research/v0.18-exact-evaluation/README.md).

**v0.19: Preserve evaluation research and Go-only tooling (complete).** Archive
three rejected post-v0.18 fits, preserve their findings and retire non-Go
research runners with original evidence retained in the release attachment.
Add shared Go research checks and tracked agent guidance. Playing policy remains
v0.18. See the [research summary and next steps](docs/RESEARCH_V019.md).

**v0.20: Remove repeated search-legality work (complete).** Specialize immutable
child transitions, compute attack rays once and skip opponent-check tags when
only legal-move existence is needed. Preserve fixed-depth scores, chosen moves,
node counts and search policy. Keep the unsuccessful quiet-mate and transition-
only experiments archived. See the
[implementation, validation and measured limits](tools/benchmark/research/v0.20-search-legality/README.md).

**v0.21: Compute stable move order with less work (complete).** Calculate each
existing priority once, sort small rank records and apply their permutation to
the original moves. Preserve exact ordering, fixed-depth scores/moves/nodes,
values and search policy. See the
[implementation, validation and measured limits](tools/benchmark/research/v0.21-rank-order/README.md).

**v0.22: Avoid unused legal work and repeated exchange scans (complete).**
Generate later move stages only when needed, prefilter impossible check evasions
and compute the existing legal exchange heuristic with local bitsets. Preserve
complete move order, fixed-work scores/moves/nodes, values and search policy.
Keep unsuccessful experiments and all original results archived. See the
[implementation, validation and measured limits](tools/benchmark/research/v0.22-bitboard-exchange/README.md).

**v0.23: Reuse search state and avoid temporary tag copies (complete).** Reuse
worker-owned child positions and boards by search ply, and classify ordinary
moves from virtual occupancy and attacker sets. Preserve fixed-work moves,
scores, nodes and playing policy. Archive unsuccessful research and add a
repeatable Go research workflow. See the
[implementation and measured limits](tools/benchmark/research/v0.23-state-kernel/README.md).

**Next:** Reassess [the documented research priorities](docs/RESEARCH_V023.md)
before starting another experiment. No further candidate is part of this release.
