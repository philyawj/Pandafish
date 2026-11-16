# Client incidents and confirmed limits

The original [client log](raw/lichess-bot.log) is retained in full. The
[grouped incident extraction](client-incidents.json) preserves the messages;
[operations](operations.json) and [audit](session-audit.json) summarize them.
Counts refer to grouped ERROR/CRITICAL records and explicit control-stream
markers, not every wrapped traceback line or keyword occurrence.

| Category | Records | Interpretation |
| --- | --- | --- |
| HTTP rate limit | 104 | Request rejected or delayed by the service; not an engine crash |
| Opponent daily game cap | 98 | Challenged account had reached its bot-game limit |
| Own daily game cap | 2 | This account reached its bot-game limit |
| Control-stream reconnect marker | 6 | Connection interruption/recovery; preserve the associated traceback |
| Other client error | 2 | Opponent rating restrictions, not chess-engine failures |

Examples retained in the incident JSON include `bot.vsBot.day`, separate
`bot_is_rate_limited` and `opponent_is_rate_limited` fields, and server retry
timeouts. Those identify which account is limited. The two other errors
describe the opponent's rating requirements. Do not classify all 212 records
as 212 failed games or bypass the server's caps to increase collection volume.

The full log supplies the reconnect context. No raw server SSE payloads were
recorded, so it is not sufficient to reconstruct every delivered event.
No overlapping search intervals were found in 12,091 requests, with a peak
of one process searching. This does not measure whether two idle game objects
ever existed, but it supplies no recurrence of the v0.10 simultaneous-search
failure. There is no new demonstrated lichess-bot bug to fix from these counts.

The manifest identifies lichess-bot commit
`baa2b7e78320b99125a91c2b1805bb7a5c4fb6ea` with `client_dirty: true`.
The exact client working-tree patch was not captured. A commit hash alone
cannot reproduce that client state; this is a provenance limitation. The
previous [v0.10 investigation](../v0.10-twitch-m2-20261112/UPSTREAM_CLIENT_NOTE.md)
already documents discovery of an older duplicate-worker bug and the local
update to the upstream fix. Do not imply that this session discovered or
fixed it again. No companion client files were changed during this review.

## Reporting-tool defect, fixed separately

The original Pandafish report linked 11,615 searches. Its group intersection
required every request within one process/game boundary to match a saved game.
An extra final request whose returned move never appeared in the PGN could
therefore erase valid earlier links. This is a defect in
`tools/lichess/report.go`, not in move selection or lichess-bot.

The repaired matcher uses full history, FEN and the actual next move, intersects
only requests with matching evidence, and leaves truly unmatched requests
unlinked. Conflicting evidence does not justify a majority guess. Focused tests
cover an unplayed final request, conflicting game evidence, common openings,
identical games, transposed positions and `ucinewgame` boundaries.

On a temporary copy it recovers **406 links**, for **12,021 linked searches**.
Seventy remain unresolved: sixty requests from two indistinguishable identical
games, nine requests at saved final roots whose returned moves are absent from
the PGNs, and one other unmatched request. Their existence is not evidence of
an illegal move or a confirmed submission error. Preserve the uncertainty.

The [derived repaired report](analysis/linker-fixed-report/README.md) is separate
from the byte-preserved original. Full-game Stockfish analysis used PGN moves,
so its move evaluations did not depend on fixing this linker.
