# Client investigation summary for a possible future report

Status: historical client bug reproduced; existing upstream fix found. Keep
this evidence locally in the archive. Do not file a duplicate report now.

With concurrency one, lichess-bot 2026.8.9.2 (`df7e730`) started a second handler
for an already-running game after a control-stream premature-end error. Two
engine processes searched glYlALQg ply 18 concurrently. A handler then freed
the unique game ID and the client accepted LentusChess while the duplicate
handler continued playing. This explains how the logical slot count could
understate running work. Raw SSE payloads were not logged, so exact event
replay remains inferred rather than directly captured.

Expected: one worker owns each game, and handler exit cannot free a slot while
another worker for that game is still running. Actual installed function probe:
two starts enqueue two jobs but leave one logical ID. Pinned upstream
`baa2b7e` probe: two starts enqueue one job, including for a pre-reserved game.
The watcher test deliberately simulates event replay; no HTTP or worker runs.

Both incidents logged `Control stream error, reconnecting:` followed by
`urllib3.exceptions.ProtocolError: Response ended prematurely` and
`requests.exceptions.ChunkedEncodingError: Response ended prematurely`.
UTC timestamps: 17:47:21.533 and 18:18:21.653 on November 12, 2026.
Exact full traces, acceptance context, line ranges, and hashes are in this
folder; complete source evidence remains in `../raw/`.

[PR #1262](https://github.com/lichess-bot-devs/lichess-bot/pull/1262) already fixes
this pattern through separate worker ownership and in-process stream recovery.
The local client has now been updated to 2026.11.12.1 (`baa2b7e`), including
that fix. Its 23 selected upstream offline tests and the 12-ply Pandafish v0.10
wrapper check passed; private configs are unchanged. No custom fix or live
check was made. If it recurs afterward, supply the new revision/dependencies and a
fresh trace, reference #1262/#1266, and demonstrate which invariant still fails.
See [the full investigation](../UPSTREAM_CLIENT_NOTE.md) for reproduction scope,
source hashes, uncertainty, and the proposed regression matrix.
