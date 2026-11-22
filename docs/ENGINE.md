# Engine behavior and limits

Pandafish v1.0 implements standard chess and UCI. Strength remains limited by
its simple evaluator and bounded search.

The engine code is in `main.go`, `limits.go`, `iterative.go`, `search.go`,
`transposition.go`, `search_context.go`, `search_moves.go`, `main_move_picker.go`,
`killers.go`, `exchange.go`, `exchange_bits.go`, `search_position.go`,
`evaluation.go`, `mate.go`, `repetition.go`, and `contact_threat.go`.
A pinned local copy of [corentings/chess](https://github.com/corentings/chess)
with an additive search API handles chess rules,
including castling, en passant, and promotions. The
[local search-backend note](CHESS_SEARCH_BACKEND.md) explains the v0.17
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
Away from check it normally generates legal tactical moves directly into reusable
per-ply storage and can keep the current static score instead of making a
capture. This helps avoid counting a captured piece before considering the
recapture. A static exchange estimate skips losing ordinary captures outside
check, while checking captures, promotions, and required evasions remain eligible.
The estimate follows legal least-valuable recaptures on the destination square;
it is a material heuristic and can miss sacrifices with other tactical benefits.
Captures and promotions have a 16-ply safety bound. Required check evasions can
extend for at most four further plies. At that final safety cap, static evaluation
is used even if checks remain, so long tactical sequences can still be missed.


## Protected queen-contact threats in v1.0

A legal quiet queen move adjacent to the king, protected by an attacker pawn,
knight or king, can select a specific mating threat. A hypothetical pass only
selects the threat. Search then considers actual legal defenses and revalidates
the proposed mating move after each defense. Stand-pat is suppressed during this
verification, and defenses are not omitted by the exchange heuristic. One event
per quiescence path is allowed, with two remaining ordinary plies. Capture,
countercheck or another legal response can refute the threat. This bounded
mechanism does not solve every quiet threat. Material/placement values and the
other released search rules remain unchanged. See the
[implementation and limits](../tools/benchmark/research/v0.24-contact-threat/README.md).

## UCI commands and time allocation

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
