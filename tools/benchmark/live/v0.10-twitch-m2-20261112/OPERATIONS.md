# Live session operations

Engine: Pandafish v0.10; source revision `1511e88b5e51752b775487822d27eebb2f61fa9c`; clean engine build: true.

Session duration: 1h27m48s. Platform: darwin/arm64.

14 games against 13 distinct opponents: 8 wins, 0 draws, 5 losses, 1 unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.

| Time control | Games |
| --- | --- |
| 180+0 | 2 |
| 180+2 | 11 |
| 60+2 | 1 |

Maximum hard-budget overrun: 0.054917 ms. Maximum response minus search time: 2.280 ms. 0 clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals 2310.0 seconds; this is not measured CPU utilization.

Fixed requests: `map[go movetime 10000:14]`. Clock requests: 536.

Completed searches: 550; linked searches: 510; unique linked positions: 509; positions with repeated searches: 1; unlinked completed searches: 40. Timing totals include all completed searches, including repeated and unlinked attempts.

Longest response: 9995.400 ms at [https://lichess.org/game/v0aclfVS](https://lichess.org/game/v0aclfVS), zero-based ply 1. Request: `go movetime 10000`; allocated hard budget 9995.000 ms. Compare the request and budget before calling a long search a stall.

## Client errors

Counts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.

- control_stream_reconnect: 2
- opponent_daily_game_cap: 2

## Games

| Game | Bot color | Opponent | Clock | Result | Termination |
| --- | --- | --- | --- | --- | --- |
| [game](https://lichess.org/game/dbkD6ZNt) | White | AetherBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/6j1YSLRZ) | White | Sacrifice-Engine | 180+0 | loss | Normal |
| [game](https://lichess.org/game/eNxSUcbC) | White | Vic98lg | 60+2 | loss | Normal |
| [game](https://lichess.org/game/JwrF9K8d) | White | O-Rei-CX | 180+2 | win | Normal |
| [game](https://lichess.org/game/dvyPuK12) | Black | KirinChessBot | 180+2 | win | Time forfeit |
| [game](https://lichess.org/game/y1xXiZbd) | Black | AetherBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/djj7sLqm) | White | BlueMoonBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/glYlALQg) | White | michaAmateurBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/AMmBtiCB) | White | LentusChess | 180+2 | win | Normal |
| [game](https://lichess.org/game/v0aclfVS) | Black | ShallowRedBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/yZ6wObvi) | White | PARTNER3615DIAGO | 180+2 | win | Normal |
| [game](https://lichess.org/game/3y5c9jT5) | White | princhess_policy_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/yQH0wqtl) | Black | uxugin-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/qAjzFwzX) | Black | Blunder-Engine | 180+0 | unfinished | Abandoned |
