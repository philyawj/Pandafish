# lichess-bot reconnect investigation

**The installed client has a reproducible duplicate-worker bug. Current
upstream already contains its fix.** The evidence matches
[issue #1266](https://github.com/lichess-bot-devs/lichess-bot/issues/1266), fixed by
[PR #1262](https://github.com/lichess-bot-devs/lichess-bot/pull/1262) on November 2,
2026. A new duplicate report or replacement fix is not recommended on this
evidence. After investigation, the local client was updated to include that
existing fix. No new client-source patch was written.

## Versions and exact errors

The stream used client revision `df7e730de58cc3ef2f1415a0dc2eeda842d39167`,
version 2026.8.9.2, and Pandafish v0.10. Post-session inspection found no tracked
client changes, only an untracked `venv/`. This was not an independent
capture-time source snapshot.

Both errors contain this chain:

```text
Control stream error, reconnecting:
urllib3.exceptions.ProtocolError: Response ended prematurely
requests.exceptions.ChunkedEncodingError: Response ended prematurely
```

Complete timestamped tracebacks and contiguous acceptance/slot context are
preserved without reformatting in [excerpt 1](client-investigation/error-context-1.log.txt)
and [excerpt 2](client-investigation/error-context-2.log.txt).
[The index](client-investigation/excerpt-index.json) records original line ranges
and hashes; the entire [original log](raw/lichess-bot.log) is preserved too.
The traceback cannot distinguish a server, network, or local connection
interruption as the reason the underlying connection ended.

## Timeline

Times are UTC; subtract eight hours for the stream's PST clock. Line numbers
refer to the original `raw/lichess-bot.log`.

| Time | Lines | Observed event |
| --- | --- | --- |
| 17:47:21.246 | 732 to 736 | Sacrifice's handler frees its slot; Vic98lg challenge eNxSUcbC is accepted. |
| 17:47:21.533 | 737 to 759 | First premature-end traceback in `watch_control_stream`, line 121. |
| 17:47:21.565 to 17:47:22.823 | 760 to 789 | Vic98lg is queued and started, then logged as used again; Sacrifice 6j1YSLRZ also starts again. Logical count reaches two. |
| 17:47:32.883 | 835 to 836 | A Vic98lg handler frees its entry; traces still contain later searches individually matching that game while Sacrifice searches. |
| 18:18:20.065 | 2961 to 2963 | Original michaAmateurBot handler searches glYlALQg, move 10. |
| 18:18:21.653 | 2964 to 2986 | Second premature-end traceback. |
| 18:18:23.036 to 18:18:23.241 | 2987 to 2997 | glYlALQg starts again at the same ply. Two recorded processes overlap for 2.980 seconds and both return a2a3. |
| 18:18:26.635 to 18:18:26.844 | 3009 to 3017 | Logical count becomes zero; LentusChess AMmBtiCB is accepted and starts. The other glYlALQg process keeps searching. |
| 18:20:32.128 | 3430 onward | glYlALQg completion leaves LentusChess active. |
| 18:22:16.024 | 3527 onward | ShallowRedBot starts with LentusChess active; logical count is two. |
| 18:24:53.045 | 3902 onward | LentusChess completion leaves one game active. |

The last three transitions have exact [context 3](client-investigation/error-context-3.log.txt),
[context 4](client-investigation/error-context-4.log.txt), and
[context 5](client-investigation/error-context-5.log.txt).
[Capture intervals](capture-diagnostics-summary.json) prove overlapping work,
not CPU utilization.

## Confirmed behavior and remaining uncertainty

The installed `start_game_thread` always enqueues a new `play_game` job, even
when that ID has a worker. `active_games` counts unique IDs, so two jobs for one
ID count as one. The pool has `max_games + 1` workers. A `local_game_done` event
pops the ID and can expose a free slot while a duplicate worker still runs.
The old `final_queue_entries` can log "Game over" for a disconnected handler
without a terminal game state. That message alone does not prove completion.

The [offline reproducer](client-investigation/reproduce.py.txt) executes actual
installed scheduling and watcher function bodies with a fake pool and stream.
It reproduces two queued jobs for one logical ID, including when a challenge
reserved the ID before its first start. Completion removal is explicitly
modeled, not run through a live main loop. No workers execute and HTTP is blocked.
A synthetic reconnect deliberately replays a start, which the watcher forwards
twice. This establishes a feasible path, not the real server's exact payloads.

The INFO log lacks raw SSE events, so every delivery and worker exit cannot be
reconstructed. Repeated starts, duplicate same-ply searches, slot transitions,
and the upstream report strongly support the duplicate-handler explanation.
No Pandafish TT failure or engine crash is established. Two daily challenge-cap
rejections elsewhere in the log are separate incidents.

## Existing fix and current-upstream check

Read-only research on November 12 checked upstream
`baa2b7e78320b99125a91c2b1805bb7a5c4fb6ea`. PR #1262 introduced:

1. Separate `games_in_progress` ownership to suppress duplicate starts, cleared
   on completion. `active_games` cannot be the sole guard because it includes
   reserved challenges and games discovered at startup.
2. Reconnection inside the existing `play_game` process, retaining its engine.
3. Stream reopening after rate limits or failed moves while the game is active.

The same offline probe using the selected function bodies from pinned current
upstream queues one job for two starts and still starts a pre-reserved game's
first worker. It verifies the scheduling guard, not the newer client's complete
integration or reconnect lifecycle. [Results](client-investigation/reproduction-results.json)
record source hashes/assertions; [research metadata](client-investigation/upstream-check.json)
records the pinned source and related PRs. Older issue #1163 had a different
communication-error fix in #1164. Overflow proposal #1241 closed unmerged;
do not implement its abort policy based on this investigation.

### Completed local update

On November 12, the local checkout was fast-forwarded to the latest upstream
revision fetched that day, `baa2b7e78320b99125a91c2b1805bb7a5c4fb6ea`, version
2026.11.12.1. The fix merge commit `0c805828da75f923e3975426302814273b50d8c2`
is an ancestor of the installed HEAD. No custom fix was added.

All 23 selected upstream offline tests passed, including duplicate starts,
reserved challenges, reconnect/rate-limit/failed-move recovery, correspondence,
event handling, and config validation. The active config loaded successfully
and its Pandafish v0.10 engine returned 12 legal plies through the UCI wrapper.
Both private configs retained identical bytes and regular 0644 permissions.
[Update validation](client-investigation/update-validation.json) records the
versions and commands. No live listener, API game, or account action ran.
Historical reproduction results still describe the old installed source;
rerunning its old-side assertion now requires a separate historical checkout.

## If an upstream contribution becomes necessary

Only pursue a new report/fix if this recurs on an updated client. Supply exact
revision/dependencies, one-account/one-client setup, concurrency/ponder settings,
sanitized contiguous errors, slot timeline, game IDs, duplicate-worker evidence,
and a small offline reproducer. Test duplicate starts, reserved challenges,
normal completion, failure, correspondence recovery, and reconnect/rate-limit
paths. A stale completion must not release another worker's ownership. Do not
blindly guard membership in `active_games` or reduce the pool size without
checking correspondence work.

This bundle explains the installed-client failure and supports comparing a
future recurrence with the existing fix. No upstream issue, comment, PR,
account action, or custom client-source modification was made. The companion
checkout update and offline validation described above are complete.
