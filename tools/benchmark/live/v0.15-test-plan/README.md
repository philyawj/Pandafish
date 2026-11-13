# Extended v0.15 live investigation

## Purpose and scope

Collect roughly **200 completed Lichess games**, possibly more, across multiple
sessions on one fixed Pandafish v0.15 build. This is intentionally larger than
the initial 50-game suggestion, to gather varied real positions and recurring
failures at realistic game clocks. Record the actual total and why collection
ended; 200 is a planning target, not a claim that data has already been collected.

v0.15 corrects mate-distance scoring and recognizes threefold history. Its local
previous-release screen passed narrowly, while external scores did not establish
a broad strength leap. Recent speculative search/evaluation attempts did not earn
promotion. The next improvements should come from diagnosed game failures and
focused regression fixtures rather than feature lists or estimated Elo guesses.

The collection hardware is the **Apple M2 MacBook Air**. Use AC power and keep
it awake. Record browser/streaming and other competing workloads accurately in
each session. Keep engine build and challenge settings consistent throughout;
separate any intentional change into a new labeled session and analysis group.
Do not start another engine version, Stockfish analysis or local benchmark while
collecting games on this machine.

## Playing profile

The private GoPandaFish config in the companion lichess-bot checkout:

- Accept rated/casual **standard chess** from humans and bots at all supported
  time controls. Unsupported variants are rejected.
- Send outgoing rated standard challenges to bots at **3+2** only, with broad
  opponent ratings and a one-minute idle delay.
- Play one game at a time, with matchmaking during games and pondering disabled.
- Use Pandafish's own moves, with books, external move sources and tablebases off.
- Finish active games on graceful shutdown. Long incoming games can make stopping
  take longer than expected.
- State v0.15 and the data-collection purpose in factual greetings/outros, thanking
  opponents for helping the engine's development.

The active configuration is `/Users/philyawj/projects/lichess-bot/config-gopandafish.yml`.
This repo's ignored `config-gopandafish.yml` is a separate private snapshot.
Both were updated for this release; neither is committed. The public example
keeps its conservative casual-human/no-outgoing defaults.

The [lichess-bot documentation](https://github.com/lichess-bot-devs/lichess-bot/wiki/Configure-lichess-bot#note-on-daily-game-limits)
currently documents 100 bot-vs-bot games per day, with human games not subject to
that bot-vs-bot limit. Expect a couple hundred bot games to span days. Session
restarts do not reset Lichess limits. Retain rate-limit messages as operational
evidence and allow the client/server to enforce limits. No automatic game-count
stop or campaign-resume feature is being added.

## Start each session

The engine and recorder are built during release setup. To rebuild later, use
the released v0.15 source, then `go build -o pandafish .` and
`go build -o bin/pandafish-lichess ./tools/lichess`. Keep the engine build identical
across the campaign; do not rebuild between sessions without recording the change.

On this Mac, run from a single terminal:

```sh
cd /Users/philyawj/projects/Pandafish
session="local-games/v0.15-live-$(date +%Y%m%d-%H%M%S)"
env -u LICHESS_BOT_TOKEN caffeinate -i ./bin/pandafish-lichess record \
  --client /Users/philyawj/projects/lichess-bot \
  --config config-gopandafish.yml \
  --engine "$PWD/pandafish" \
  --session "$session" \
  --hardware "Apple M2 MacBook Air" \
  --workload "browser and terminal"
```

If streaming, change `--workload` to describe it before starting. The launcher
starts lichess-bot, displays its output and collects the data. Do not separately
start lichess-bot in another terminal. `caffeinate -i` prevents idle system sleep;
keep the laptop open and powered. This command starts online play, so run it only
when ready. Release verification uses offline preparation, not this live command.

Each launch needs a fresh session directory. The recorder copies/hashes the
engine and points its runtime config at that fixed session copy; rebuilding the
root binary does not change an already-running session. A prepared-only directory
also cannot be reused for a live session. Unsetting `LICHESS_BOT_TOKEN` selects the
private config's token; never paste it into the command or public records.

## Stop, check, and start another session

Press **Ctrl+C once** in the recording terminal. Wait for lichess-bot to finish
active games and for the recorder to finalize and print its report. A second
interrupt force-quits; use that only when necessary and retain the partial data.
Closing the window is not the normal shutdown procedure.

After the command returns:

```sh
cat "$session/summary.md"
```

The recorder automatically generates this report. To regenerate it offline:

```sh
./bin/pandafish-lichess report --session "$session"
```

Keep every session, including failures, reconnects and warnings. Record the
session path. To continue later, repeat the complete start block; the timestamp
creates a new directory. **Do not append to or overwrite the previous session.**
There is no single session-resume flag and no automatic cross-session pooling.
If using a new terminal, the old `$session` variable is gone; assign the actual
saved path before reporting or processing it:

```sh
session="local-games/v0.15-live-YYYYMMDD-HHMMSS"
```

Several shorter sessions totaling a couple hundred completed games satisfy the
plan. They need not be one uninterrupted run. Avoid overlapping account clients.
The test plan stays here; actual session archives get their own sibling folders.

## Offline processing and preservation

When no live session or benchmark is running on this machine, process each
finished session individually:

```sh
./bin/pandafish-lichess process \
  --session "$session" \
  --archive "tools/benchmark/live/${session##*/}" \
  --bot GoPandaFish \
  --stockfish /opt/homebrew/bin/stockfish \
  --nodes 100000
```

This archives all source files with SHA-256 hashes and analyzes PGN moves offline
with full-strength Stockfish, not Skill 2. The captured executable remains local,
ignored and hashed; raw game/trace/log data are retained after secret review.
The tool preserves originals and writes analysis separately. It does not need to
fetch games from Lichess when the client saved the PGNs. Missing games/linkage
are reported, not silently reconstructed or guessed.

A couple hundred games can take substantial offline analysis time. Optionally
add `--limit 20` for a small processing pilot, then rerun the same command without
that option to resume all plies from matching checkpoints. This limits plies,
not games. Reusing a matching archive is supported for processing; reusing a
recording directory for another live run is not. Do not change analysis settings
in a cached profile: select a new `--profile` when changing node budget/method.

The tool selects up to 20 confirmation candidates per session. Recheck important
flags at 500,000 nodes using the documented `--focus`/`--profile` workflow in the
[recorder README](../../../lichess/README.md#archive-and-analyze-a-finished-session).
Finite-node reference results can change. Retain both initial flags and confirmed
results, including flags that disappear.

## Campaign synthesis and development loop

1. Inventory every session's engine/config hashes, client revision, workload,
   completed/unfinished games and recording warnings. Verify v0.15 identity and
   compare hashes before pooling; investigate inconsistencies instead of assuming
   all sessions used the same build.
2. Deduplicate by Lichess game ID across sessions, including reconnects or resumed
   correspondence games. Preserve all raw records and link them to the canonical
   game. The existing tool reports sessions individually; campaign deduplication
   and synthesis are separate analysis work when records are available.
3. Analyze all retained unique games, including losses and uneventful draws.
   Separate the main standard 3+2 bot cohort from incoming humans, other clocks,
   rated/casual games and overlapping-process incidents. Repeated opponents and
   shared openings remain correlated; a larger live sample is not a paired Elo
   benchmark. Do not interpret account-rating changes as isolated engine strength.
4. Select roughly 20-30 consequential positions across games/phases, prioritizing
   deterioration from playable positions rather than only already-lost mates.
   Confirm references at a larger budget and preserve comparable good decisions.
   Replay older confirmed cases on v0.15 to discard failures already fixed.
5. Classify confirmed failures: tactical continuations, quiet positional choices,
   mate/promotion conversion, repetition/draw handling, or clocks/client operation.
   Use equal-budget root choices, controlled deeper searches, restricted-move
   comparisons and leaf-score breakdowns to distinguish hypotheses. A bad king
   move does not automatically imply missing king safety; a deeper correct move
   does not by itself identify the next search feature.
6. Publish a compact campaign review with counts, representative game/ply/FEN
   evidence, uncertainty and a ranked next increment. Keep detailed results in
   `tools/benchmark/live`, with only direction and links in the main README.
7. Turn a reproduced failure into a focused test, implement one supported
   correction, and validate with separate development/validation positions plus
   the normal 50-game concurrent local comparisons. Keep public validation cases
   out of weight tuning; they are not secret holdouts. Gather fresh live evidence
   periodically after changes.

The intended output is a useful failure collection and tested next hypothesis,
not a promise that 200 games guarantee a particular Elo gain or an automatic
training/database integration. Search and evaluation remain complementary; let
the evidence guide which needs work.

## What to provide when collection ends

Say: **"The v0.15 live campaign is finished. Process all sessions below, preserve
all raw records, deduplicate games across sessions, run offline Stockfish analysis,
and synthesize the next engine improvement."** Include every session directory,
any known errors/overlap, and whether streaming or other workloads changed.
You may process offline between sessions or leave that analysis for the completed
campaign. Either way, keep the original session directories intact.
