# Capture and analysis process

## What happened during this session

1. Built Pandafish v0.8 and the Go Lichess session tool on macOS ARM64.
2. Configured the installed lichess-bot checkout for standard chess: rated or
   casual incoming challenges from humans/bots at all BOT-supported clocks;
   automatic outgoing challenges to bots at rated 3+2 only. One concurrent
   game, no pondering, one-minute idle matchmaking delay, broad opponent
   ratings, temporary decline filtering, no external move sources.
3. Loaded both private configs through the installed validator. Ran 34
   synthetic incoming-rule cases, mock outgoing selection/retry/backoff,
   and a 12-ply local UCI wrapper check. These were offline checks.
4. The user manually launched `record` with `caffeinate -i`, displayed the
   account's TV and terminal side by side, and streamed from the M2 Air.
   `record` copied the engine, recorded its hash/build identity, redirected
   PGNs into the session, captured console output, and enabled engine traces.
5. After play, the user stopped the client, added `hardware.txt`, and manually
   ran `report`. The manifest records a normal client exit, 38 finished games,
   39 engine processes, and 1,531 fully linked searches. The extra process is
   consistent with startup engine validation; process count is not game count.
6. Archived every data file and the local executable without altering source
   bytes. The credential-pattern/YAML check passed. `inventory.json` records
   hashes and sizes, including the locally retained executable.
7. Grouped wrapped ERROR records from the client log, calculated game outcomes
   and clock/depth/response evidence, and preserved that derivation in Go.
8. Ran a 20-ply Stockfish pilot at 100,000 nodes/search, then resumed all games.
   The pilot took about eight seconds. No Lichess API request was needed:
   complete local PGNs and protocol records were already available.

The old client has two config traps worth preserving: omit
`opponent_rating_difference` instead of setting it to null when matchmaking is
enabled, and omit `challenge_days` instead of setting an empty list. The latter
passes config validation but fails when matchmaking calls `random.choice`.
Those issues were fixed before the stream started.

The capture marked the client checkout dirty. A post-session inspection found
no tracked source changes and only an untracked `venv/` directory.
`provenance/client-environment.json` records the revision, status, Python and
installed package versions inspected after the session; it is explicitly not
a capture-time environment snapshot. No temperature, CPU-frequency, or resource
utilization samples were recorded, so this archive cannot quantify thermal
throttling or the streaming workload's causal performance impact.

## Reproduce the complete first-pass analysis

From the repository root, with Go 1.25+ and Stockfish installed:

```sh
go run ./tools/lichess process \
  --session tools/benchmark/live/v0.8-twitch-m2-20261111/raw \
  --archive tools/benchmark/live/v0.8-twitch-m2-20261111 \
  --bot GoPandaFish --stockfish stockfish --nodes 100000
```

The initial analysis used `/opt/homebrew/bin/stockfish` (Stockfish 19). The
binary hash and complete settings are in `analysis/stockfish-n100000/method.json`.
An installed Stockfish with the same name but a different hash is a distinct
method. Use `--profile` with a new descriptive name instead of overwriting a
profile with different settings. Existing matching checkpoints are reused.
A source/archive hash mismatch stops processing. No network access is needed.

The survey covers every played move by both players. Each comparison uses
unrestricted full-strength Stockfish and, when needed, a second search
restricted to the actual played move, at the same root position and node
budget. Full move history preserves repetition context. Every independent
search clears the hash and starts a new UCI game. Settings are one thread,
64 MiB hash, Skill Level 20, MultiPV 1, and no Elo limiting. This is a referee
analysis, separate from the Skill Level 1 strength matches.

Scores favor the root mover. Positive cp values favor that player; mate scores
retain their type/sign/distance and are never converted into invented cp loss.
Finite-budget searches can disagree, so negative raw deltas are preserved and
reported separately. Mate-transition labels describe what the two finite
searches reported; a played search failing to find a mate does not prove that
no mate exists. The 100/300 cp buckets are local triage labels, not
Lichess accuracy or official mistake classifications.

The approach was informed by `chess-ai/scripts/analyze_game.py`, which records
exact engine settings, candidate lines, and before/after evaluations. This Go
workflow uses independent unrestricted/restricted searches at the same root
rather than subtracting an after-position search with a different root. It uses
one PV and a smaller initial budget so all games can be screened before deeper
follow-up. Both methods are finite-search estimates.

## Larger-budget confirmation

```sh
go run ./tools/lichess process \
  --session tools/benchmark/live/v0.8-twitch-m2-20261111/raw \
  --archive tools/benchmark/live/v0.8-twitch-m2-20261111 \
  --bot GoPandaFish --stockfish stockfish --nodes 500000 \
  --profile stockfish-n500000-confirmed \
  --focus tools/benchmark/live/v0.8-twitch-m2-20261111/analysis/stockfish-n100000/confirmation.json
```

This checks a deterministic selection of at most 20 bot positions, at most one
per game. It prioritizes the first substantial deterioration from a roughly
playable position, then fills from the survey's other candidate positions.
The focused sample is deliberately selected and does not estimate session-wide
accuracy. The survey JSON retains every candidate, not just the selected ones.

## Future sessions with fewer manual steps

Build the current tool once:

```sh
go build -o bin/pandafish-lichess ./tools/lichess
```

Start a named session, using your own client path, private config, and hardware:

```sh
caffeinate -i ./bin/pandafish-lichess record \
  --client /absolute/path/to/lichess-bot --config config-mybot.yml \
  --session local-games/vNEXT-stream-01 \
  --hardware 'Apple M2 MacBook Air' \
  --workload 'game browser and Twitch streaming'
```

On platforms without `caffeinate`, omit that wrapper. Unset a stale
`LICHESS_BOT_TOKEN` environment variable if the config's token should be used.
The tool now writes the hardware note and automatically generates the capture
report on client exit. The session flag can be omitted for an automatically
named directory; use the path printed by the launcher afterward. There is no
need to manage a shell variable, echo a hardware file, or run `report` manually.

Then use one offline command to archive, summarize operations, and analyze:

```sh
./bin/pandafish-lichess process \
  --session local-games/vNEXT-stream-01 \
  --archive tools/benchmark/live/vNEXT-stream-01 \
  --bot YourBot --stockfish stockfish --nodes 100000
```

Review the compact outputs and run the optional higher-budget confirmation
command before writing recommendations. Keep the automated evidence separate
from human interpretation. Turn understood failures into focused regression
positions, keep a public hold-out set for later evaluation, and validate actual
engine changes through the controlled benchmark workflow. Retain unfavorable
sessions. Do not merge live scores into paired-match benchmark rows.

### Returned move and scored PV consistency

Stockfish can stop a node-limited search during an aspiration-window iteration.
Its final `bestmove` can then differ from the last PV with an exact score.
The processor never attaches that earlier move's score to the returned move.
It independently re-searches the selected reference move when needed, retaining
all unrestricted selection output in `selection_output` and setting
`reference_rechecked`. The comparison score and PV then belong to that move.
A regression test covers this case. Matching cached results from the initial
implementation are repaired using their original selection output; original
capture data is untouched.

The initial full survey finished in 875.8 seconds after the 20-ply pilot.
A consistency audit identified 39 returned-reference/PV mismatches, including
seven comparisons that had reused the unrestricted result for the played move.
Only those selected reference moves needed new queries. Original unrestricted
outputs remain embedded in the repaired results. The initial summary is kept
under provenance with a superseded label. The early focused run is also marked
preliminary; final recommendations use the repaired survey and the separate
stockfish-n500000-confirmed profile.

## Final verification

All 85 original files and their archive copies matched their inventory hashes.
The final survey contains 3,069 legal scored reference/played PV pairs and the
focused run contains 20. Every returned reference agrees with its scored PV,
and every played comparison starts with the actual PGN move. All 1,531 bot
moves retain their matching live metrics. The one process without a search
contains exactly `uci`, `isready`, and `quit`, confirming startup validation.
The 17 grouped client errors are opponent daily-game-cap responses. Go tests,
race checks, vet, formatting, documentation links, and whitespace checks ran.
No live listener, challenge, account change, or game download was needed for
the analysis. No engine move-selection or clock-allocation code changed.

Depth, nodes, and time attached to a Stockfish score come from its last exact-
scored PV. The raw output also retains the final, potentially bounded
iteration, so reported score-iteration nodes need not equal the node budget.
