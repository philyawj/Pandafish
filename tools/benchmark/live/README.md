# Live Lichess sessions

Recorded real-game sessions and their offline analyses belong here. Capture and
analysis programs live in [tools/lichess](../../lichess/README.md). Controlled
engine matches remain in [results](../results/README.md).

| Session | Engine | Hardware/workload | Games | Status |
| --- | --- | --- | --- | --- |
| [First Twitch stress test](v0.8-twitch-m2-20261111/README.md) | v0.8 | M2 MacBook Air, browser and Twitch streaming | 38 | Complete survey, targeted confirmation, and review |

These games have selected opponents, mixed clocks, and an account rating that
changes across games. They are operational and position evidence, not a paired
strength benchmark or an Elo estimate. Preserve unsuccessful sessions too.

Each session includes a file inventory, original capture data under `raw/`,
operational results, analysis settings, per-game Stockfish output, and a written
review. Raw records are copied without rewriting them. The captured executable
is also retained locally and hashed, but compiled executables are not tracked
in Git. A clone can analyze the archived records without that executable.

Do not put an unredacted API config in an archive. `process` rejects obvious
credential patterns and unredacted credential keys in YAML before copying.
That check is a guardrail rather than a guarantee that arbitrary private text
is safe to share; review the inventory and records before committing them.
The original private session remains in ignored `local-games/`.

Archive files opt out of line-ending conversion so hashes survive checkout.
Raw archive paths disable Git whitespace-error checking because original
terminal output contains padding. Do not strip or reformat that evidence.
Authored code and documentation retain normal whitespace checks.
