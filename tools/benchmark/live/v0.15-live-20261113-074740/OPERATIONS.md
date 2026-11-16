# Live session operations

Engine: Pandafish v0.15; source revision `d7a61498a16c8dcdd5582b85406eb568917189d7`; clean engine build: true.

Session duration: 58h54m13s. Platform: darwin/arm64.

269 games against 170 distinct opponents: 77 wins, 28 draws, 160 losses, 4 unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.

| Time control | Games |
| --- | --- |
| 120+1 | 1 |
| 120+10 | 1 |
| 120+3 | 1 |
| 120+5 | 1 |
| 15+2 | 1 |
| 180+0 | 8 |
| 180+1 | 3 |
| 180+10 | 3 |
| 180+2 | 189 |
| 240+4 | 1 |
| 30+0 | 1 |
| 30+2 | 1 |
| 30+3 | 1 |
| 300+0 | 2 |
| 300+1 | 2 |
| 300+2 | 1 |
| 300+3 | 10 |
| 300+5 | 1 |
| 60+0 | 2 |
| 60+1 | 7 |
| 60+2 | 9 |
| 60+3 | 1 |
| 60+4 | 1 |
| 600+0 | 1 |
| 600+1 | 1 |
| 600+10 | 2 |
| 600+2 | 4 |
| 600+3 | 2 |
| 600+5 | 7 |
| 900+0 | 1 |
| 900+10 | 1 |
| 900+2 | 1 |
| 900+3 | 1 |

Maximum hard-budget overrun: 1.131750 ms. Maximum response minus search time: 10.767 ms. 132 clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals 65048.0 seconds; this is not measured CPU utilization.

Fixed requests: `map[go movetime 10000:266]`. Clock requests: 11825.

Completed searches: 12091; linked searches: 11615; unique linked positions: 11615; positions with repeated searches: 0; unlinked completed searches: 476. Timing totals include all completed searches, including repeated and unlinked attempts.

Longest response: 37462.405 ms at [https://lichess.org/game/XGkRZN2j](https://lichess.org/game/XGkRZN2j), zero-based ply 2. Request: `go wtime 898999 btime 900000 winc 10000 binc 10000`; allocated hard budget 37461.300 ms. Compare the request and budget before calling a long search a stall.

## Client errors

Counts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.

- control_stream_reconnect: 6
- http_rate_limit: 104
- opponent_daily_game_cap: 98
- other_client_error: 2
- own_daily_game_cap: 2

## Games

| Game | Bot color | Opponent | Clock | Result | Termination |
| --- | --- | --- | --- | --- | --- |
| [game](https://lichess.org/game/nzT8wmvw) | Black | Lightning-762 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/aYIPyAYl) | Black | Sacrifice-Engine | 180+0 | win | Normal |
| [game](https://lichess.org/game/TLkNQ4qP) | White | littleply | 300+3 | loss | Normal |
| [game](https://lichess.org/game/P1vTXlIe) | Black | Blunder-Engine | 180+0 | unfinished | Abandoned |
| [game](https://lichess.org/game/WfKXj0lh) | White | Blunder-Engine | 180+0 | unfinished | Abandoned |
| [game](https://lichess.org/game/mkEryUPf) | White | DrewsChessMachine | 300+3 | win | Normal |
| [game](https://lichess.org/game/z7FhDSVk) | Black | ailedbot | 180+2 | win | Normal |
| [game](https://lichess.org/game/5WC5FRtB) | Black | little-0 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/uOq0lhNd) | Black | uxugin-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/RML70jnQ) | White | sargon-1ply | 180+2 | win | Normal |
| [game](https://lichess.org/game/Ev21Chsu) | White | shrevenBOT | 180+2 | win | Normal |
| [game](https://lichess.org/game/K76aglQu) | Black | Hyperopic | 180+2 | loss | Normal |
| [game](https://lichess.org/game/yY7KYF4r) | White | Bongaclang | 180+2 | loss | Normal |
| [game](https://lichess.org/game/r60KNAff) | Black | remychessengine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/rpk4dbYK) | Black | Ballin-bot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/cBqH484w) | Black | HiusBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/WbBEK119) | Black | kolmato | 60+0 | loss | Normal |
| [game](https://lichess.org/game/fGPMHX2e) | White | FlowerLilac | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mWe70ObE) | White | thunfisch-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/r3PL3Tky) | Black | TopasBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/PxR2B1nJ) | Black | sunfish-nnue-engine | 180+0 | draw | Normal |
| [game](https://lichess.org/game/qfbjN5Ng) | White | colinbot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/cu5ff2fG) | White | NICE_BOT | 180+2 | win | Normal |
| [game](https://lichess.org/game/yPBz4ev0) | White | crocodile_bot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/qP6h6HDa) | White | bfiedler-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/FtRp2MlL) | Black | TerritoryBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/qaE8fqwI) | Black | MimirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/AfSNlfZC) | Black | uSunfish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7k20Z33B) | Black | Oxide-Eg | 600+2 | loss | Normal |
| [game](https://lichess.org/game/jtlVBF5N) | White | Chessiverse_Ethan | 180+2 | win | Normal |
| [game](https://lichess.org/game/4yxgnkGg) | Black | dala-1300 | 180+2 | win | Normal |
| [game](https://lichess.org/game/CBtJ5xEy) | Black | Cloudia-Chess | 180+2 | win | Normal |
| [game](https://lichess.org/game/K9VpY9Gn) | White | Angela_eng | 180+2 | win | Normal |
| [game](https://lichess.org/game/JarAPKZX) | White | chiral_ebt | 180+2 | loss | Normal |
| [game](https://lichess.org/game/MeUyNCid) | Black | chiral_ebt | 180+2 | loss | Normal |
| [game](https://lichess.org/game/B77ylgJp) | Black | uSunfish | 180+2 | draw | Normal |
| [game](https://lichess.org/game/zHzzp6IR) | White | little-0 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/CbHeGyMS) | White | notropis | 180+2 | loss | Normal |
| [game](https://lichess.org/game/zMAhZYmv) | Black | CookieBot711 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/WelgYsq1) | White | dala-1100 | 180+2 | win | Normal |
| [game](https://lichess.org/game/f3Z8Vm0P) | White | Vortex_Bot | 120+10 | loss | Normal |
| [game](https://lichess.org/game/Qg6L6Wm8) | White | Angela_eng | 900+2 | draw | Normal |
| [game](https://lichess.org/game/YzfaekGC) | Black | Tente-V2 | 180+2 | draw | Normal |
| [game](https://lichess.org/game/jXYhnq8u) | Black | ShallowRedBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/nBl4tfTd) | Black | GambitOnlyBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Lz2z3V5K) | Black | stickshark99 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mjlIrzUJ) | White | crushinator_bot | 300+3 | win | Normal |
| [game](https://lichess.org/game/YmgkAlpK) | White | Tente8x | 180+2 | loss | Normal |
| [game](https://lichess.org/game/AA4qRpWT) | White | Oxide-Eg | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ATzIvZ7w) | Black | matmoi | 180+2 | win | Normal |
| [game](https://lichess.org/game/Mt31Jibf) | White | crocodile_bot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/WfDjpws7) | White | Raven030 | 180+2 | win | Normal |
| [game](https://lichess.org/game/RSMJOMZk) | White | littleply | 180+2 | draw | Normal |
| [game](https://lichess.org/game/7x1tlyzh) | White | Hyperopic | 180+2 | loss | Normal |
| [game](https://lichess.org/game/IPKBGhp3) | Black | ConfidenceBuilder | 180+2 | win | Normal |
| [game](https://lichess.org/game/p15kqznG) | White | croco_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/YDlfsTjm) | White | Vic98lg | 60+1 | loss | Normal |
| [game](https://lichess.org/game/qhb4iAQB) | Black | shrevenBOT | 180+2 | win | Normal |
| [game](https://lichess.org/game/16G1f9rx) | Black | DarkOnBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ZVo6SSJu) | White | SquareFish-900 | 180+2 | win | Normal |
| [game](https://lichess.org/game/ZoOwfuZk) | White | ShallowRedBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/4R57Qmf5) | Black | croco_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/wuolaQaa) | Black | VietnameseCoffee | 180+2 | win | Normal |
| [game](https://lichess.org/game/QaHolDKW) | White | Ballin-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/ypnxqSZG) | White | O-Rei-CX | 180+2 | draw | Normal |
| [game](https://lichess.org/game/Wlj5PDq4) | Black | littleply | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mKw8R7C8) | White | felix_pufferlib | 300+1 | win | Normal |
| [game](https://lichess.org/game/SHb60gFb) | Black | AetherBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/qhUVejTs) | Black | PeachFruit | 180+2 | loss | Normal |
| [game](https://lichess.org/game/4kMsR96I) | White | Disrespect_Fish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/RZgMZce5) | White | Ferrite | 60+1 | loss | Normal |
| [game](https://lichess.org/game/doe0PNEA) | Black | incipit_bot | 180+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/ynCKmF0L) | Black | UnifiedChess-Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/kt2f2Ez6) | White | felix_pufferlib | 180+2 | win | Normal |
| [game](https://lichess.org/game/BefOoIGC) | White | PeachFruit | 180+2 | loss | Normal |
| [game](https://lichess.org/game/L0mwQLiR) | White | Oxide-Eg | 600+5 | draw | Normal |
| [game](https://lichess.org/game/3ZyPpSGl) | White | AggressiveStockfish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/WPUQlQyc) | Black | R0bspierre | 180+2 | win | Normal |
| [game](https://lichess.org/game/LJLDrQbM) | Black | eubos | 180+2 | loss | Normal |
| [game](https://lichess.org/game/3fFS9j9T) | White | Martuni | 180+2 | loss | Normal |
| [game](https://lichess.org/game/XGkRZN2j) | White | Glimbus | 900+10 | loss | Normal |
| [game](https://lichess.org/game/UeNjeyLW) | Black | colinbot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/pgi5eGLQ) | White | Ndr1sn | 180+2 | loss | Normal |
| [game](https://lichess.org/game/PGSZLKXg) | White | psion-chess | 180+0 | loss | Normal |
| [game](https://lichess.org/game/cHx1ti7d) | Black | bot_adario | 180+2 | loss | Normal |
| [game](https://lichess.org/game/0RV7CCdA) | White | SeedChess | 180+2 | loss | Normal |
| [game](https://lichess.org/game/BqbTfnng) | Black | grail-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ZATxhHXk) | Black | sseh-c | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ICEHQnX7) | White | Blujay-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/J0oJTUv4) | Black | sunfish-nnue-engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/1T8KKC69) | White | DarkElectronBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/H48lBmnp) | Black | Perl-GigaChess | 180+2 | win | Normal |
| [game](https://lichess.org/game/ooyg8SSZ) | Black | Bongaclang | 180+2 | loss | Normal |
| [game](https://lichess.org/game/kiOWTCJC) | Black | sunfish-nnue-engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/mLPtUMHP) | Black | NICE_BOT | 180+2 | loss | Normal |
| [game](https://lichess.org/game/9AprsVIr) | White | clessy-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/PNfQoYP0) | White | TerritoryBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/emMx3CVk) | White | FenrirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/hNr2PJbY) | White | Ballin-bot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/0coUnZSY) | White | AshNostromo | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ketFIvs9) | White | GeoTsiBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/T9eASxFQ) | White | Goopis | 180+2 | draw | Normal |
| [game](https://lichess.org/game/aOIecGbM) | Black | sunfish-nnue-engine | 180+1 | unfinished | Abandoned |
| [game](https://lichess.org/game/Yrz1bWLG) | White | AshNostromo | 180+2 | loss | Normal |
| [game](https://lichess.org/game/li4LOWY3) | Black | Chessiverse_Raak | 180+2 | win | Normal |
| [game](https://lichess.org/game/V962S41M) | Black | shrevenBOT | 180+2 | win | Normal |
| [game](https://lichess.org/game/zJegkegP) | Black | simple-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/zbCrEJbl) | Black | Blujay-bot | 300+0 | loss | Normal |
| [game](https://lichess.org/game/kqW5oiXD) | Black | glorbo_evermorrow | 180+2 | loss | Normal |
| [game](https://lichess.org/game/QlV5OJJg) | White | Styx-Chess | 60+1 | loss | Normal |
| [game](https://lichess.org/game/4RtvQXj3) | Black | mermaid_engine | 60+2 | win | Normal |
| [game](https://lichess.org/game/ITpXfIYb) | Black | rust-in-pieces | 30+2 | loss | Normal |
| [game](https://lichess.org/game/wIpdjT47) | White | NEWTSbot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/xa2sCgDA) | White | neryba | 180+2 | loss | Normal |
| [game](https://lichess.org/game/d9tFa7Y9) | Black | mate1-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/aaEEX0mp) | Black | TrinculoBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/WBYzfQmn) | White | HackerChessBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/TWyu8niJ) | Black | scipio-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/MEwNOoFD) | White | pi0w | 180+2 | loss | Normal |
| [game](https://lichess.org/game/JzWAukWt) | Black | Tigrament | 180+2 | loss | Normal |
| [game](https://lichess.org/game/U8obBbtt) | Black | CookieBot711 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7YECFY7h) | White | AshNostromo | 180+2 | loss | Normal |
| [game](https://lichess.org/game/I5e00b8C) | Black | KFCBot | 600+0 | loss | Normal |
| [game](https://lichess.org/game/vciXyGUA) | White | GeoTsiBot | 300+3 | win | Time forfeit |
| [game](https://lichess.org/game/R63ZRca9) | Black | Fobero | 60+2 | loss | Normal |
| [game](https://lichess.org/game/foyR30Fv) | White | uxugin-bot | 600+3 | win | Normal |
| [game](https://lichess.org/game/tOZZkKHM) | White | SpringerBot | 180+0 | win | Normal |
| [game](https://lichess.org/game/oEhEi7KY) | White | BockFish | 60+2 | loss | Normal |
| [game](https://lichess.org/game/OKOeNL7e) | Black | CassidyChess-Bot | 600+1 | loss | Normal |
| [game](https://lichess.org/game/3s6oiPMI) | White | kleiniBOT | 180+0 | loss | Normal |
| [game](https://lichess.org/game/UMtDZiNy) | Black | uundrstnd | 180+10 | loss | Normal |
| [game](https://lichess.org/game/kkjlnQkD) | White | MetisBot | 180+10 | loss | Normal |
| [game](https://lichess.org/game/nhcjG05R) | Black | HydraBot3141 | 60+1 | draw | Normal |
| [game](https://lichess.org/game/exSaaUfH) | Black | artidigm_research | 60+2 | draw | Normal |
| [game](https://lichess.org/game/99udi4R8) | White | zyrox-bot | 600+3 | draw | Normal |
| [game](https://lichess.org/game/kQ4iEA2q) | Black | manu-bot | 600+5 | win | Normal |
| [game](https://lichess.org/game/PyGK4Lnn) | Black | uundrstnd | 600+2 | loss | Normal |
| [game](https://lichess.org/game/HLk8Ajmi) | White | DrewsChessMachine | 300+0 | win | Normal |
| [game](https://lichess.org/game/uYLasbHI) | White | StargazerEngine | 120+3 | loss | Normal |
| [game](https://lichess.org/game/YfP3JSfu) | White | TheFinalChessCard | 600+10 | win | Normal |
| [game](https://lichess.org/game/ul58ftvF) | Black | ConsoleChessFastBot | 180+1 | win | Normal |
| [game](https://lichess.org/game/4asRCida) | Black | KirinChessBot2000 | 60+2 | loss | Normal |
| [game](https://lichess.org/game/Avap2OJY) | Black | entebot | 15+2 | loss | Normal |
| [game](https://lichess.org/game/eiKMw75Q) | Black | CCE_V2 | 60+4 | loss | Normal |
| [game](https://lichess.org/game/74OmsK55) | White | ScalaQueen | 900+0 | win | Time forfeit |
| [game](https://lichess.org/game/WFQzZjAI) | White | uxugin-bot | 120+5 | win | Normal |
| [game](https://lichess.org/game/UtQr1tHN) | Black | UnifiedChess-Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/h16hpLR8) | Black | chmod-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/zcdj0W2k) | White | Ursus_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/J0Nay5Rd) | Black | Tente8x | 180+2 | loss | Normal |
| [game](https://lichess.org/game/DsqFd5Kj) | Black | felix_pufferlib | 180+2 | win | Normal |
| [game](https://lichess.org/game/BLb1GuB7) | Black | Quirkafleeg-ZXNext | 300+3 | loss | Normal |
| [game](https://lichess.org/game/vBSvdO2f) | Black | To299d | 180+2 | loss | Normal |
| [game](https://lichess.org/game/DeozJiuS) | White | colinbot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/e8j054gd) | Black | Oxide-Eg | 900+3 | loss | Normal |
| [game](https://lichess.org/game/u0H60HmC) | Black | TilapiaBot | 600+5 | loss | Normal |
| [game](https://lichess.org/game/U2LaF6mc) | Black | ShallowRedBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/BGjcMkWb) | Black | ShallowRedBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/dXniMk2j) | Black | ratsu-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/NHnjB3DH) | White | littleply | 180+2 | loss | Normal |
| [game](https://lichess.org/game/VCbs3ECV) | White | smartysprout | 180+2 | win | Normal |
| [game](https://lichess.org/game/2koPBAP5) | Black | ShallowRedBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/CeUJN37G) | Black | RubyAI | 180+2 | loss | Normal |
| [game](https://lichess.org/game/cD7EwYD7) | Black | PeachFruit | 180+2 | loss | Normal |
| [game](https://lichess.org/game/MnaqVQc3) | Black | DevilfishZero | 180+2 | loss | Normal |
| [game](https://lichess.org/game/n4OYimuW) | Black | LeechGODestroy | 180+2 | loss | Normal |
| [game](https://lichess.org/game/BYybyvba) | White | BretemaBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/XzRGgjkv) | Black | SeedChess | 180+2 | loss | Normal |
| [game](https://lichess.org/game/6LAr0xIy) | Black | ConfidenceBuilder | 180+2 | win | Normal |
| [game](https://lichess.org/game/yIRBiE2X) | Black | TrinculoBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/mPUuYpZn) | Black | Chessiverse_Ethan | 180+2 | win | Normal |
| [game](https://lichess.org/game/HYhlfE4F) | Black | princhess_policy_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/Z7cAAmyS) | Black | laura-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/e6mgJNx7) | Black | Cloudia-Chess | 180+2 | win | Normal |
| [game](https://lichess.org/game/XNGMBc0p) | Black | croco_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/8mQaEKbe) | Black | io-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/FUnqAQXc) | White | TerritoryBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/9sfF61Hc) | White | DevilfishZero | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ZTKmWVYZ) | Black | NICE_BOT | 180+2 | loss | Normal |
| [game](https://lichess.org/game/63KWrCYk) | Black | MimirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/cid6fU2j) | White | psion-chess | 60+2 | loss | Normal |
| [game](https://lichess.org/game/BhrzqEje) | White | Oxide-Eg | 600+5 | loss | Normal |
| [game](https://lichess.org/game/HhPI3vCx) | Black | artidigm_research | 180+2 | loss | Normal |
| [game](https://lichess.org/game/pPG9ldxa) | White | ApexChess-Bot | 60+0 | win | Time forfeit |
| [game](https://lichess.org/game/PkwQ1Fkk) | White | TroutBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/19t5KUV4) | Black | SquareFish-900 | 180+2 | win | Normal |
| [game](https://lichess.org/game/bmjPgzjz) | Black | jpg-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/XHXhCNDe) | White | Flower-Queen | 180+2 | draw | Normal |
| [game](https://lichess.org/game/x89jO4W3) | White | HomoLudensBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/Qs1Xqlm2) | White | DarkElectronBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7UnWrqw3) | White | macbotitsmagic | 300+3 | loss | Normal |
| [game](https://lichess.org/game/rGmjbC5l) | Black | Waldemar-the-elf | 180+2 | win | Normal |
| [game](https://lichess.org/game/5FjSVrf9) | White | ShallowRedBot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/FvU0kOVi) | White | Grindfish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/jnsu7vCq) | Black | Disrespect_Fish | 60+1 | loss | Normal |
| [game](https://lichess.org/game/QkGDF0kv) | White | Casanchess-NNUE | 180+2 | loss | Normal |
| [game](https://lichess.org/game/GbDJ06Gh) | Black | Chessiverse_Pawnza | 180+2 | win | Normal |
| [game](https://lichess.org/game/FUyRo61O) | White | Bit-Fish8 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/IW8Rlhum) | White | FenrirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/rI0oJG0N) | White | CassidyChess-Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/QghUQNGi) | Black | prokopakop | 180+2 | loss | Normal |
| [game](https://lichess.org/game/v2I3Zvhk) | White | thunfisch-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/FG05Et8M) | White | Martuni | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ofogXLgg) | Black | Lightning-762 | 60+1 | loss | Normal |
| [game](https://lichess.org/game/P47w1ojK) | Black | chess_set_manager | 180+2 | loss | Normal |
| [game](https://lichess.org/game/uITt5KAi) | Black | MiccoEngine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/k9tE98QI) | Black | Cheszter | 180+2 | loss | Normal |
| [game](https://lichess.org/game/v3AGvAE8) | Black | Ballin-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/qSud3dNr) | Black | SeedChess | 180+2 | loss | Normal |
| [game](https://lichess.org/game/EP9iU6Hd) | Black | DeMen100ns-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/GkJnHETo) | White | littleply | 180+2 | draw | Normal |
| [game](https://lichess.org/game/d2PjsdwJ) | White | shrevenBOT | 180+2 | draw | Normal |
| [game](https://lichess.org/game/F3sIUo7d) | Black | SpringerBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/40CCYWYU) | Black | CarlsenStyleBot | 300+3 | loss | Normal |
| [game](https://lichess.org/game/KgdBoDhU) | Black | michaAmateurBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/Pds3R7kn) | Black | Angela_eng | 600+2 | loss | Normal |
| [game](https://lichess.org/game/Ebx1nIud) | Black | KirinChessBot1500 | 600+2 | win | Normal |
| [game](https://lichess.org/game/W0ErjnQY) | Black | Z-O_AspenLabs | 600+10 | loss | Normal |
| [game](https://lichess.org/game/yuWlaBRS) | Black | Cosmic_Engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/2F2UIN0j) | White | Was_Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/FurvigAM) | White | Gvait | 300+2 | loss | Normal |
| [game](https://lichess.org/game/tWJ6StAN) | Black | Zyqorath | 300+3 | win | Normal |
| [game](https://lichess.org/game/jOhf490d) | Black | myChessJava | 600+5 | loss | Normal |
| [game](https://lichess.org/game/B1WTqhJy) | Black | AnechkaChess | 30+0 | win | Time forfeit |
| [game](https://lichess.org/game/jCPLCkuL) | White | Jimmy_inc3710 | 240+4 | loss | Normal |
| [game](https://lichess.org/game/AQPxTfST) | Black | KirinChessBot1500 | 300+1 | win | Normal |
| [game](https://lichess.org/game/nw14BP8l) | White | Ponytail_engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/cnK9lkT9) | White | Botkingteam33 | 180+1 | win | Normal |
| [game](https://lichess.org/game/IaUg953e) | Black | uundrstnd | 180+10 | loss | Normal |
| [game](https://lichess.org/game/QcAZ6Mxb) | Black | Quirkafleeg-ZXNext | 300+3 | win | Normal |
| [game](https://lichess.org/game/XCwl7FvN) | White | ix-bot | 60+2 | draw | Normal |
| [game](https://lichess.org/game/9oRhkgb2) | White | Botkingteam33 | 60+1 | loss | Normal |
| [game](https://lichess.org/game/IVZNvGck) | White | sunfish-engine | 60+3 | loss | Normal |
| [game](https://lichess.org/game/EXBvLZKP) | White | Stepbot_Chess | 30+3 | loss | Normal |
| [game](https://lichess.org/game/utSAjyhY) | White | crushinator_bot | 300+5 | win | Normal |
| [game](https://lichess.org/game/HBykKGiX) | White | Rodent-IV | 120+1 | win | Time forfeit |
| [game](https://lichess.org/game/g5RSmCvy) | White | DrewsChessMachine | 300+3 | win | Normal |
| [game](https://lichess.org/game/JsvmD06t) | Black | fractionBOT | 60+2 | loss | Normal |
| [game](https://lichess.org/game/19ZuG5on) | Black | Quirkafleeg-ZXNext | 600+5 | loss | Normal |
| [game](https://lichess.org/game/uKoQGHOl) | White | KFCBot | 600+5 | draw | Normal |
| [game](https://lichess.org/game/AZsG8MxC) | Black | HardRok | 180+0 | win | Normal |
| [game](https://lichess.org/game/ILAhQR7F) | White | patgro986_bot | 60+2 | loss | Normal |
| [game](https://lichess.org/game/qX6GPM1t) | Black | psion-chess | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Vpqb3wpv) | White | Disrespect_Fish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Lmo91cuN) | Black | Chessiverse_Pawnza | 180+2 | win | Normal |
| [game](https://lichess.org/game/EmWSRK1u) | Black | PeachFruit | 180+2 | loss | Normal |
| [game](https://lichess.org/game/pZbbqbWo) | White | stickshark99 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/buYVo59V) | White | Goopis | 180+2 | draw | Normal |
| [game](https://lichess.org/game/9h2V1VXX) | White | uSunfish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/XkM9SOjv) | White | MostlyHuman1900 | 180+2 | win | Normal |
| [game](https://lichess.org/game/lmH485Oo) | White | ai-con | 180+2 | win | Normal |
| [game](https://lichess.org/game/GT8u2Nkb) | Black | StockConda | 180+2 | loss | Normal |
| [game](https://lichess.org/game/YKZH2bKm) | White | MDBOT | 180+2 | loss | Normal |
| [game](https://lichess.org/game/r0Ov9X8D) | Black | HomoLudensBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/TfxxR4P8) | Black | Jibbby | 180+2 | win | Normal |
| [game](https://lichess.org/game/WcdggQyT) | Black | abcd_engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/aTSlVmnQ) | White | JemBok | 180+2 | loss | Normal |
| [game](https://lichess.org/game/aubUxIkO) | White | JurdBird | 180+2 | win | Normal |
| [game](https://lichess.org/game/O4wRsvjl) | Black | To299d | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Cxu6BwTG) | Black | TikiFish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/RbSX4WaF) | White | Ballin-bot | 180+2 | draw | Normal |
| [game](https://lichess.org/game/Lj7IuLO1) | Black | MostlyHuman1900 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/LUQ8RErT) | Black | Chessiverse_Pawnza | 180+2 | win | Normal |
| [game](https://lichess.org/game/t9moKsKY) | White | UnifiedChess-Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/9omGrw4L) | White | Tigrament | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7vm186dX) | White | AshNostromo | 180+2 | loss | Normal |
| [game](https://lichess.org/game/CAsMVtgL) | Black | sunfish-engine | 180+2 | draw | Normal |
| [game](https://lichess.org/game/WKWFHDIF) | White | grail-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/9khJdBnw) | Black | O-Rei-CX | 180+2 | loss | Normal |
