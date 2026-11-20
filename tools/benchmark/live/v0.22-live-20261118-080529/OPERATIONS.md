# Live session operations

Engine: Pandafish v0.22; source revision `1953f2af97f956c4499cbc5b61c3c1c698856740`; clean engine build: true.

Session duration: 3h54m31s. Platform: darwin/arm64.

19 games against 19 distinct opponents: 5 wins, 2 draws, 12 losses, 0 unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.

| Time control | Games |
| --- | --- |
| 180+2 | 16 |
| 180+3 | 1 |
| 60+1 | 1 |
| 900+3 | 1 |

Maximum hard-budget overrun: 8.420375 ms. Maximum response minus search time: 3.508 ms. 0 clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals 4424.3 seconds; this is not measured CPU utilization.

Fixed requests: `map[go movetime 10000:19]`. Clock requests: 815.

Completed searches: 834; linked searches: 833; unique linked positions: 833; positions with repeated searches: 0; unlinked completed searches: 1. Timing totals include all completed searches, including repeated and unlinked attempts.

Longest response: 32212.466 ms at [https://lichess.org/game/nVVSe2oF](https://lichess.org/game/nVVSe2oF), zero-based ply 2. Request: `go wtime 898999 btime 900000 winc 3000 binc 3000`; allocated hard budget 32211.300 ms. Compare the request and budget before calling a long search a stall.

## Client errors

Counts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.

- control_stream_reconnect: 5
- opponent_daily_game_cap: 9
- other_client_error: 1

## Games

| Game | Bot color | Opponent | Clock | Result | Termination |
| --- | --- | --- | --- | --- | --- |
| [game](https://lichess.org/game/Mido8Dw4) | White | justbot-v1 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/PzEtqL6B) | Black | GambitOnlyBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Gl7kmsLj) | White | zyrox-bot | 180+3 | win | Normal |
| [game](https://lichess.org/game/EQggxkZN) | Black | stickshark99 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mmMvJUHm) | Black | Blujay-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/PV4oUeoe) | White | NichtEliasBot | 60+1 | loss | Normal |
| [game](https://lichess.org/game/EpCUge8J) | Black | kolmato | 180+2 | loss | Normal |
| [game](https://lichess.org/game/LP57FP1I) | Black | mate1-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mQQjCDiw) | White | dala-1100 | 180+2 | win | Normal |
| [game](https://lichess.org/game/mMZ5Ne11) | White | nebubot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/8a73DmQx) | Black | colinbot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/SODGjjrZ) | Black | coda_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/nVVSe2oF) | White | Oxide-Eg | 900+3 | draw | Normal |
| [game](https://lichess.org/game/XgOBT8eH) | White | NoHopeForMan | 180+2 | loss | Normal |
| [game](https://lichess.org/game/y2IkeSMw) | Black | icarus-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/E3RvY9Nk) | White | Tente-V2 | 180+2 | win | Normal |
| [game](https://lichess.org/game/SRn3UtST) | White | dala-1300 | 180+2 | draw | Normal |
| [game](https://lichess.org/game/9CDIInpP) | White | notropis | 180+2 | win | Normal |
| [game](https://lichess.org/game/2byZKYLT) | Black | PZChessBot | 180+2 | loss | Time forfeit |
