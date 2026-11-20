# Live session operations

Engine: Pandafish v0.22; source revision `1953f2af97f956c4499cbc5b61c3c1c698856740`; clean engine build: true.

Session duration: 6h4m52s. Platform: darwin/arm64.

33 games against 32 distinct opponents: 14 wins, 3 draws, 16 losses, 0 unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.

| Time control | Games |
| --- | --- |
| 180+1 | 1 |
| 180+2 | 28 |
| 300+3 | 1 |
| 480+2 | 1 |
| 60+1 | 1 |
| 600+5 | 1 |

Maximum hard-budget overrun: 0.029584 ms. Maximum response minus search time: 3.393 ms. 58 clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals 7628.7 seconds; this is not measured CPU utilization.

Fixed requests: `map[go movetime 10000:33]`. Clock requests: 1480.

Completed searches: 1513; linked searches: 1512; unique linked positions: 1512; positions with repeated searches: 0; unlinked completed searches: 1. Timing totals include all completed searches, including repeated and unlinked attempts.

Longest response: 23711.444 ms at [https://lichess.org/game/fAa1ypzA](https://lichess.org/game/fAa1ypzA), zero-based ply 2. Request: `go wtime 598999 btime 600000 winc 5000 binc 5000`; allocated hard budget 23711.300 ms. Compare the request and budget before calling a long search a stall.

## Client errors

Counts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.

- control_stream_reconnect: 1
- http_rate_limit: 6
- opponent_daily_game_cap: 17

## Games

| Game | Bot color | Opponent | Clock | Result | Termination |
| --- | --- | --- | --- | --- | --- |
| [game](https://lichess.org/game/Ah9LcQlH) | White | sseh-c | 180+2 | win | Normal |
| [game](https://lichess.org/game/u4xtX1iv) | White | Chessiverse_Pawnza | 180+2 | loss | Normal |
| [game](https://lichess.org/game/1TrbIGaE) | White | FenrirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/2V40jKiJ) | Black | MudskipperBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/JwEJzRsX) | Black | Noksiy | 180+2 | loss | Normal |
| [game](https://lichess.org/game/A4Hg9awN) | Black | Chessiverse_Ethan | 180+2 | win | Normal |
| [game](https://lichess.org/game/EDoYRsVz) | Black | TroutBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/D6rn4haK) | White | MilanBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/c4POr7th) | White | MarlinChessBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/g29qunfO) | Black | SimonEricAfonso | 180+2 | win | Normal |
| [game](https://lichess.org/game/XDKx7goI) | White | prokopakop | 180+2 | loss | Normal |
| [game](https://lichess.org/game/oRSwsRb3) | Black | LentusChess | 300+3 | win | Normal |
| [game](https://lichess.org/game/bcmp5xRI) | White | uSunfish-l7 | 480+2 | win | Normal |
| [game](https://lichess.org/game/RuCn73pd) | Black | bot64jaques | 180+2 | loss | Normal |
| [game](https://lichess.org/game/jMayAOLf) | White | kopyto_dev | 180+2 | draw | Normal |
| [game](https://lichess.org/game/ZiWwAAsF) | White | SJ_NorBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/9m86yw5x) | White | GambitOnlyBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/qRnP8DNI) | Black | Jaxakos13_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/MhLL6txE) | Black | lamprook | 180+2 | loss | Normal |
| [game](https://lichess.org/game/sG9ZxJFN) | White | princhess_policy_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/o8HptRoh) | White | Chessiverse_Luc | 180+2 | win | Normal |
| [game](https://lichess.org/game/U1Z9lTsh) | White | incipit_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/X4j9rEtC) | White | Ponytail_engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/fAa1ypzA) | White | manu-bot | 600+5 | win | Normal |
| [game](https://lichess.org/game/VWr9LKNk) | White | MostlyHuman1900 | 180+2 | win | Normal |
| [game](https://lichess.org/game/vrPZQy5O) | White | rudim-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/gwliYouc) | Black | NICE_BOT | 180+2 | draw | Normal |
| [game](https://lichess.org/game/bw0hYYqY) | White | Iskelt-Corax | 60+1 | loss | Normal |
| [game](https://lichess.org/game/faVjqhvv) | White | HardRok | 180+1 | loss | Normal |
| [game](https://lichess.org/game/otr1y1Bs) | White | AetherBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/5khQm35W) | White | NICE_BOT | 180+2 | win | Normal |
| [game](https://lichess.org/game/R5SDcDSV) | White | icarus-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/sIxuenPs) | Black | Angela_eng | 180+2 | win | Normal |
