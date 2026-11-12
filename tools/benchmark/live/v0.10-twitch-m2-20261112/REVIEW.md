# Review and recommendations

v0.10 completed the observed live workload without illegal engine moves,
missing final search metrics, or legal fallbacks. The new stream also contains
client reconnect and concurrency problems, so it is not a clean isolated speed
comparison. The client investigation identified an older lichess-bot duplicate-worker bug;
the local checkout now includes its existing upstream fix. The next strength
experiment should first address check-aware quiescence, followed by quiet-move
ordering. Local testing should drive the next batch of engine changes.

## What was recorded

The clean v0.10 release ran on the M2 MacBook Air with the browser and Twitch
software for 1h27m48s. The 14 PGN records contain 13 played games against 12
distinct opponents: 8 wins and 5 losses. Seven wins have normal termination;
one is the opponent's time forfeit after two plies. The fourteenth record has
zero moves, result `*`, and termination `Abandoned`. It should not be described
as the engine failing to finish a played endgame. All recorded opponents are
BOT accounts. Most records use 3+2; two use 3+0 and one uses 1+2.

All 550 engine searches returned legal moves and complete metrics, with no
fallback. The original report links 510 search attempts to 509 unique PGN
positions. One position has two completed attempts. Forty searches retain
missing/ambiguous linkage in the original evidence; they were not discarded
or relabeled to make the report look complete.

The audit individually matches 38 of those 40 searches by full move history,
root FEN, and the actual next PGN move. The two initial-position requests remain
ambiguous. The conservative original matcher requires consistency across a
process/game group, so an inconsistent initial request hides otherwise matching
later positions. The individual matches are diagnostic candidates, not a claim
that we identified an authoritative submission for every request. Stockfish
analyzes the actual PGN moves independently of capture metadata, covering all
548 moves by the GoPandaFish account.

## Timing and unintended concurrent play

The largest search-budget overrun was 0.054917 ms. The headline count of 440
budget crossings does not represent 440 time forfeits. The longest response,
9995.400 ms, followed `go movetime 10000` with a 9995 ms engine hard budget.
Fourteen requests used that special ten-second first-move allowance; the other
536 used clocks. No supplied own clock fell below ten seconds. Response p95
was 7273.105 ms, and the largest response-minus-search interval was 2.280 ms.
These do not measure network submission or acknowledgement time.

The configuration and recorded workload text say one active game, but the
operator reported a period of simultaneous games. The traces confirm 134
searches overlapping a different process, a peak of two processes, and about
168.5 seconds of overlapping request intervals. The overlap is intermittent,
not continuous between the earliest and latest timestamps.

| Approximate UTC envelope | Processes / associated games | Summed overlapping seconds |
| --- | --- | --- |
| 17:47:23 to 17:47:32 | Two initial unlinked requests | 8.748 |
| 17:47:33 to 17:50:35 | Individually matched Vic98lg searches and Sacrifice | 48.196 |
| 18:18:23 to 18:18:26 | Two attempts at michaAmateurBot ply 18 | 2.980 |
| 18:18:27 to 18:20:30 | michaAmateurBot and LentusChess | 65.541 |
| 18:22:22 to 18:24:52 | LentusChess and ShallowRedBot | 43.022 |

The log has control-stream reconnect tracebacks at 17:47:21 and 18:18:21 UTC,
plus two rejected challenges to opponents that had reached their daily bot cap.
The reconnects are separate from those daily-cap responses.
The [client investigation](UPSTREAM_CLIENT_NOTE.md) preserves the exact
tracebacks and slot timeline. An offline probe reproduces duplicate scheduling
on the installed older client; current upstream suppresses the duplicate with
the already-merged PR #1262 fix. The real SSE payloads were not captured.

Outstanding search requests are not measured CPU utilization. No temperature,
CPU-frequency, utilization, or power-source samples were captured. The stream
supports observed engine deadline/legality behavior under this mixed workload;
it cannot quantify the causal cost of streaming or prove a TT speedup from
live depth/node totals. The median recorded depth increased from 6 in v0.8 to
7 here, but positions, clocks, evaluation, and concurrent work differ.

## What the offline analysis shows

The full survey uses the same Stockfish 19 binary as the v0.8 review, one
thread, 64 MiB hash, MultiPV 1, Skill 20 with Elo limiting disabled, and
100,000 nodes per independent query. This is an unrestricted referee, not the
Skill 1 benchmark opponent. Best and played moves use the same root and budget,
with full history and fresh search/hash state. Fourteen returned-reference/PV
mismatches were independently re-scored; every stored score belongs to its move.

| Survey classification | v0.8 | v0.10 |
| --- | --- | --- |
| Saved played plies | 3,069 | 1,092 |
| GoPandaFish moves | 1,531 | 548 |
| Estimated loss 100 to 299 cp | 229 | 37 |
| Estimated loss at least 300 cp | 49 | 5 |
| Played search reports losing mate where reference does not | 30 | 4 |
| Reference reports winning mate but played search does not | 2 | 0 |
| Both searches report mate scores | 242 | 146 |
| Negative cp delta / search disagreement | 98 | 15 |
| Smaller or equal cp difference | 881 | 341 |
| Mean clamped loss, cp-to-cp subset | 61.3 cp / 1,257 cases | 34.9 cp / 398 cases |

At least-100-cp flags fell from 18.2% to 7.7% of the respective bot-move sets.
This is encouraging descriptive evidence, not an isolated v0.9/v0.10 strength
gain, standard ACPL comparison, or Elo estimate. Opponents and positions changed,
the account's rating affects matchmaking, and mate comparisons are excluded
from the cp-only mean. Keep both favorable and unfavorable examples.

## Larger-budget confirmation

Twelve selected positions from twelve games were rechecked at 500,000 nodes.
Selection prioritizes deterioration from playable positions, then fills from
other candidates. It is targeted and is not a representative accuracy sample.

Nine cp-to-cp cases still lost at least 100 cp: six between 100 and 299, three
at least 300. One promotion case still allowed a reported mate, but its reference
position was already about -11 pawns, so avoiding that mate report would not
establish that the game could be saved. Two survey flags changed substantially:
the Vic98lg king move became the reference itself, and a pawn move in a heavily
winning position fell to 51 cp loss. Do not treat all 100,000-node flags as facts.

Six confirmed cp cases have reference scores between -150 and +200 cp. They
form the saved [regression set](analysis/stockfish-n500000-confirmed/regression-candidates.json):
three development and three public holdout positions, split deterministically.
They are not hidden or statistically independent. Stockfish references can
have equally good alternatives; do not assert a unique move without checking.

### Development examples

Scores favor the player to move and come from the 500,000-node check.

| Game / decision | Played / reference | Reference / played score | Estimated loss |
| --- | --- | --- | --- |
| [AetherBot](https://lichess.org/dbkD6ZNt), White move 15 | bxc4 / a4 | -107 / -326 cp | 219 cp |
| [princhess_policy_bot](https://lichess.org/3y5c9jT5), White move 24 | Nf4 / Rad1 | +168 / +18 cp | 150 cp |
| [ShallowRedBot](https://lichess.org/v0aclfVS), Black move 19 | f6 / Qxd3 | 0 / -134 cp | 134 cp |

The first two candidates are quiet main-search resources. The third keeps
pawn structure and queen-exchange judgment on the evaluation/tactics backlog.
These examples do not prove that a particular new bonus will choose the reference.

## Captured-engine replay and recommended next increment

The exact captured v0.10 executable was replayed on only the three development
cases. The live-depth two-choice searches reproduced the recorded score and
played choice. At one extra ply, it strictly preferred the reference in two
cases by 4 and 6 of its own cp; a supplemental played-only search ruled out
random ties. The third choice was unchanged. All nine probes completed legally,
and none of these three original live searches overlapped another process.

This is a narrow two-choice diagnostic, not a full-root correction or a tested
new algorithm. The small score differences also show horizon sensitivity. The
extra ply can be expensive: one probe took 11.7 seconds against a live budget
of about 4.565 seconds. Simply demanding another ply would violate the timing
policy. See [the replay record](analysis/engine-depth-replay/README.md).

## Current strength priorities

These are ranked experiments, not promised gains or five approved releases.
The original depth replay motivated killer-move ordering, but it neither tested
that algorithm nor established it as the highest-value change. A source review
also identified the explicit static-evaluation cutoff while in check. Addressing
that tactical horizon is now the preferred first experiment for v0.11.

1. **Check-aware quiescence.** At the current two-ply cap, a nonterminal checked
   position can receive a static score before its forced response is searched.
   Start with focused positions that expose the distortion, then test bounded
   check-evasion handling while preserving deadlines and cancellation. Do not
   allow unbounded extensions. This is a concrete limitation, but it has not
   been isolated as the cause of the selected live losses or measured as a
   strength gain. The engine's returned moves remain legal.
2. **Killer-move ordering for quiet moves.** Remember a small number of legal
   quiet moves that caused beta cutoffs and try them early at the same search
   ply, after cached moves and tactical priorities. Keep evaluation, TT,
   quiescence, and clocks fixed for this separate experiment. Compare exact
   fixed-depth scores/best-move sets, node counts, and timed completed depths.
3. **Repetition and fifty-move draw awareness.** Check the current history and
   draw handling first, then add missing search-level draw handling with focused
   tests. Avoid throwing away winning positions and recognize drawing escapes
   when behind. Account for history-dependent draw scores before reusing TT
   entries; a board position alone does not describe its repetition history.
4. **Small history-based quiet-move ordering.** If killer ordering helps, test
   whether a bounded score for quiet moves that cause cutoffs improves ordering
   across positions. Add it separately so its benefit and overhead can be
   measured. Keep reset behavior and memory use explicit.
5. **Targeted pawn and king evaluation.** Use confirmed recurring errors and
   saved positions to choose one small term at a time, such as pawn weaknesses
   or king safety. Do not add a collection of plausible bonuses based only on
   stream impressions. Require broader local match evidence to limit overfitting.

Keep alpha-beta as the foundation. PVS, aspiration windows, reductions, books,
and database training can wait until simpler experiments justify more complexity.
No engine increment in this list has been implemented by this review.

## Local testing drives the next batch

Use this workflow for the next engine changes:

1. Make one focused engine change and test its intended behavior.
2. Check saved positions and search efficiency.
3. Run the usual serial 100-game comparisons against the previous version, then the current Stockfish target and pinned Sacrifice.
4. Use live Lichess sessions periodically after meaningful progress, or when changing clocks, UCI handling, or client integration.

For step 3, use the documented previous-version gate before external opponents,
and run Stockfish before Sacrifice. A 100-game screen is useful evidence, not
proof of a small gain. Longer tests need a specific decision to resolve.
The current saved-position cases are public; avoid tuning only to these examples.

One change at a time makes results easier to explain. It does not require a
stream after every change. Several independently validated increments can
precede another live session. Local matches provide cleaner strength comparisons;
live games test integration and expose positions missing from our local tests.

## Client discovery and completed update

This investigation connected the two premature-end errors to repeated game
handlers, duplicate same-ply searches, and a logical slot being freed while a
second handler kept playing. The installed August client reproduced duplicate
worker scheduling offline. Upstream had already fixed that pattern in
[PR #1262](https://github.com/lichess-bot-devs/lichess-bot/pull/1262), merged
November 2. Raw SSE payloads were not captured, so the exact server-event replay
remains inferred; the scheduling defect itself is reproduced.

After the investigation, the local companion checkout was fast-forwarded from
`df7e730de58cc3ef2f1415a0dc2eeda842d39167` to
`baa2b7e78320b99125a91c2b1805bb7a5c4fb6ea`, version **2026.11.12.1**, the latest
upstream revision fetched November 12. It includes the existing fix rather than
a new local patch. All 23 selected upstream offline reconnect, worker ownership,
correspondence, event, and config tests passed. Pandafish v0.10 passed the local
config/UCI wrapper check with 12 legal plies. Both private configs and their
0644 permissions were unchanged. No listener or online game was started.

See [the investigation](UPSTREAM_CLIENT_NOTE.md) for exact error excerpts,
reproduction scope, update evidence, and what a future report would need.
A new upstream contribution is warranted only if the issue recurs on updated
code. Pandafish engine behavior is unchanged; offline processing was corrected
to retain duplicate and unlinked attempts instead of aborting or dropping them.
