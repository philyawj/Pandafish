# Review and recommendations

The first live stress test supports v0.8's capture and search-control work.
It also gives concrete targets for improving chess decisions. The next engine
increment should test positional evaluation, starting with modest piece-square
terms. This is a hypothesis supported by the examples below, not a demonstrated
strength gain or a reason to skip controlled comparisons.

## What was actually tested

The M2 MacBook Air ran the engine, lichess-bot, game browser, and Twitch stream
for 5h38m5s. All 38 games finished. Pandafish scored 5 wins, 2 draws, and 31
losses against 34 BOT accounts. Four wins had normal termination; one was the
opponent's time forfeit. No recorded loss was a time forfeit.

All 1,531 engine searches completed, supplied final metrics, returned legal
moves, and linked to the saved PGNs. No fallback or capture warning occurred.
The 39 engine processes include a startup validation process; this is not
an extra game or evidence of a crash.

The headline count of 1,263 hard-budget overruns is misleading without their
size. The largest was 0.103834 ms. The longest response was 36.772 seconds in a
15+10 game with a 36.772-second hard budget. Response p95 was 9.995 seconds,
reflecting mixed clocks and the client's 38 explicit ten-second first-move
requests. This does not look like an unexplained search stall. The installed
client deliberately uses a special first-move allowance; it is not evidence
that the normal clock allocator ignored a clock.

There were 89 requests with less than ten seconds left on the moving side's
supplied clock. The maximum response-minus-search interval was 6.085 ms. This
interval is measured inside the engine trace, not network or submission time.
No CPU frequency, temperature, utilization, or memory-pressure samples were
captured. We cannot attribute performance changes to thermal throttling or
quantify the causal cost of streaming from this dataset.

The 17 grouped client ERROR records all report opponents that had already
reached their daily bot-game cap. These explain some matchmaking delays, not
engine failures or GoPandaFish having exhausted its own quota. Keyword counts
of rate-limit text would overcount wrapped fields and repeat messages.

## What Stockfish adds

The survey covers all 3,069 played plies, including all 1,531 Pandafish moves,
at 100,000 nodes per independent search. The referee is unrestricted Stockfish
19, one thread, 64 MiB hash, MultiPV 1. It is separate from the controlled
Skill Level 1 matches. Each query keeps full game history and resets search
state/hash. Full output and binary/settings hashes are preserved.

The corrected survey classifies Pandafish's moves as follows:

| Survey evidence | Moves |
| --- | --- |
| Estimated cp loss 100 to 299 | 229 |
| Estimated cp loss at least 300 | 49 |
| Played search reports losing mate where reference search does not | 30 |
| Reference search reports winning mate but played search does not | 2 |
| Both searches report mate scores | 242 |
| Negative raw cp delta, recorded as search disagreement | 98 |
| Smaller or equal cp difference | 881 |

The mean clamped loss for the 1,257 cp-to-cp comparisons is 61.3 cp. It excludes
mate comparisons and is not Lichess accuracy, a whole-session ACPL standard,
or an Elo estimate. The 100/300 cp buckets are local triage labels. A search
not finding a mate is not proof that one does not exist.

The live engine reported a zero material score in 453 moves. In 127 of those,
the survey estimates a loss of at least 100 cp relative to its reference move.
This is evidence that material equality can hide important distinctions. It
does not isolate evaluation from search horizon, move ordering, or finite
reference-search error.

### Why the extra consistency check mattered

Stockfish can return a move from an interrupted bound-only iteration while
its last exact-scored PV belongs to a different move. A post-run audit found
39 such reference mismatches, including seven reused played comparisons. The
processor now independently scores the returned reference in these cases.
Original unrestricted output remains in `selection_output`; only affected
reference moves required new queries. Every final score now belongs to its
stored PV/root move. A regression test covers this behavior.

The initial summary and early targeted pass remain clearly labeled as
superseded diagnostics. Current findings use the repaired survey and
[the final 500,000-node confirmation](analysis/stockfish-n500000-confirmed/README.md).

## Higher-budget confirmation

Twenty selected positions from twenty different games were checked at 500,000
nodes/search. Selection favored the first substantial deterioration from a
roughly playable position, then filled remaining slots from other candidates.
This is a targeted diagnostic sample, not an independent accuracy estimate.

- All nine selected cp-to-cp cases still lost at least 100 cp. Two exceeded
  300 cp; their mean estimated loss was 299.1 cp in this selected subset.
- Six selected positions still produced a losing mate score for the played
  move and a cp score for the reference. Those reference positions were
  already heavily losing, so they are not six examples of saving a game.
- Five mate-related cases changed interpretation: four also had a losing mate
  on the reference search, and one still mated after the played move. They
  are not confirmed as different mate-versus-cp outcomes. This is why the
  100,000-node screen should not be converted directly into prose claims.

### Development examples

Scores below favor the player to move, in pawns, from the 500,000-node check.
The named alternatives are reference moves, not proven unique choices.

| Game / decision | Played / reference | Live engine | Reference / played score | Estimated loss |
| --- | --- | --- | --- | --- |
| [Moteur441](https://lichess.org/game/mEY6KGOi), White move 6 | Ke2 / exd4 | -1.00, depth 6 | -0.46 / -4.52 | 406 cp |
| [R0bspierre](https://lichess.org/game/Lkfnlfui), White move 11 | Rg2 / e4 | 0.00, depth 5 | -0.67 / -3.53 | 286 cp |
| [LevyMaia](https://lichess.org/game/un7IhthI), Black move 8 | c5 / e6 | 0.00, depth 6 | -0.34 / -3.16 | 282 cp |

The first example places the king on e2 instead of capturing the central pawn
on d4. The second prefers another rook move over a central pawn move. These
are useful evaluation/search targets because live material scores did not
capture the size of the deeper reference differences. They do not prove that
a particular table bonus will select the reference move.

## Recommended order of engine work

1. **Positional evaluation first.** Add small piece-square terms for all
   pieces so development, centralization, pawn progress, and king placement
   can distinguish materially tied choices. Use consistent centipawn units,
   preserve material dominance and mate/stalemate handling, and retain random
   selection only among truly equal results. Keep king/endgame treatment
   explicit and small. Test color symmetry, unit scaling, and terminal-score
   dominance. Do not train tables from these nine positions or add an opening
   book to conceal the opening choices.
2. **Transposition caching and search efficiency next.** Repeated positions
   currently cost fresh work. A bounded table with tested depth and bound
   semantics could improve completed depth under the same clock policy.
   That is an efficiency hypothesis, not a strength conclusion from this
   session. The current metrics do not separate completed-iteration work from
   discarded partial work, so record that distinction if needed to evaluate it.
3. **Then target quiescence/check and draw behavior with specific tests.** The
   two-ply tactical extension remains a deliberate limit. Before changing it,
   reproduce a concrete failure and separate horizon effects from evaluation.
   This review found no cp-scored live result concealing the selected side's
   mate within the nominal completed-depth window among the checked short
   mate cases; this does not prove general tactical correctness. Add draw and
   conversion tests rather than assuming more depth solves every endgame.
4. **Keep live play as evidence alongside controlled benchmarks.** For the
   next engine candidate, run the routine 100 games against v0.8 and 100
   against the current Stockfish Skill Level 1 target. Preserve profile and
   opening settings. A higher live account rating mixes opponents, clocks,
   versions, and provisional-rating adjustment, so it is not the acceptance
   test for the change.

[Regression candidates](analysis/stockfish-n500000-confirmed/regression-candidates.json)
contain the nine confirmed cp cases: five development positions and four
public hold-out positions, split deterministically before engine changes.
This is not a hidden or statistically independent test set. Avoid using the
hold-out group to tune individual bonuses. Check alternatives and ties before
turning a Stockfish reference move into an exact-move unit test. v0.8 is not
expected to pass an improvement suite retroactively.

## Workflow changes made after this session

The Go recorder can now save hardware/workload details itself and automatically
produce the capture report on client exit. One subsequent `process` command
archives all raw records, validates hashes/credentials, summarizes operations,
and runs resumable per-game Stockfish comparisons. An optional focus command
runs the higher-budget check. Compact summaries come first; individual JSON
and PGNs are opened to investigate particular cases. There is no LLM call per
move and no Lichess API dependency for this complete archive.

See [PROCESS.md](PROCESS.md) for the actual session history, all commands,
method limitations, and the repeatable future-session workflow. The engine
itself has not changed as part of this analysis.

The follow-up [engine comparison](ENGINE_COMPARISON.md) reviews Sacrifice,
Sunfish, and Stockfish and proposes positional evaluation followed by a small
transposition table as separate increments.
