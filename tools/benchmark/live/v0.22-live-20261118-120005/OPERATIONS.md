# Live session operations

Engine: Pandafish v0.22; source revision `1953f2af97f956c4499cbc5b61c3c1c698856740`; clean engine build: true.

Session duration: 27h4m32s. Platform: darwin/arm64.

105 games against 83 distinct opponents: 42 wins, 6 draws, 52 losses, 5 unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.

| Time control | Games |
| --- | --- |
| 180+0 | 2 |
| 180+1 | 5 |
| 180+2 | 77 |
| 300+0 | 1 |
| 300+2 | 1 |
| 300+3 | 3 |
| 300+5 | 2 |
| 420+2 | 1 |
| 60+1 | 3 |
| 60+2 | 1 |
| 60+5 | 1 |
| 600+2 | 1 |
| 600+3 | 2 |
| 600+5 | 2 |
| 900+10 | 1 |
| 900+5 | 2 |

Maximum hard-budget overrun: 0.033167 ms. Maximum response minus search time: 2.172 ms. 113 clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals 28425.7 seconds; this is not measured CPU utilization.

Fixed requests: `map[go movetime 10000:103]`. Clock requests: 5057.

Completed searches: 5160; linked searches: 5084; unique linked positions: 5084; positions with repeated searches: 0; unlinked completed searches: 76. Timing totals include all completed searches, including repeated and unlinked attempts.

Longest response: 37461.449 ms at [https://lichess.org/game/AVlXnr6L](https://lichess.org/game/AVlXnr6L), zero-based ply 2. Request: `go wtime 898999 btime 900000 winc 10000 binc 10000`; allocated hard budget 37461.300 ms. Compare the request and budget before calling a long search a stall.

## Client errors

Counts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.

- control_stream_reconnect: 13
- http_rate_limit: 57
- opponent_daily_game_cap: 35
- other_client_error: 3

## Games

| Game | Bot color | Opponent | Clock | Result | Termination |
| --- | --- | --- | --- | --- | --- |
| [game](https://lichess.org/game/KgctDwEO) | Black | lamprook | 180+2 | draw | Normal |
| [game](https://lichess.org/game/v7c5pVV1) | Black | FenrirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/93JFs3RB) | White | LegoTechnicControlPl | 180+2 | loss | Normal |
| [game](https://lichess.org/game/48CHb0oJ) | Black | O-Rei-CX | 180+2 | win | Normal |
| [game](https://lichess.org/game/rMr32JUb) | Black | princhess_policy_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/DrjE4O2r) | White | TerritoryBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/FU3KLGg1) | Black | rust-in-pieces | 180+2 | loss | Normal |
| [game](https://lichess.org/game/MfmxQtF7) | White | Cloudia-Chess | 180+2 | win | Normal |
| [game](https://lichess.org/game/YmoztEDQ) | Black | chatgpt5_5 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/LW7mc1tY) | White | byte-knight | 180+2 | loss | Normal |
| [game](https://lichess.org/game/gSsa1QLN) | White | TroutBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/53HukkXX) | Black | Cactus_Bot | 600+2 | loss | Normal |
| [game](https://lichess.org/game/DRSiOoTf) | Black | nebubot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/n5N2xU29) | White | rudim-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/rofyW9Uw) | White | lamprook | 180+2 | loss | Normal |
| [game](https://lichess.org/game/bbUUoSNG) | White | Ballin-bot | 180+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/ZJE1kANQ) | White | smartysprout | 180+2 | win | Normal |
| [game](https://lichess.org/game/nViUqeab) | White | MutusBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/HVV4rIGQ) | Black | StockConda | 180+2 | win | Normal |
| [game](https://lichess.org/game/skkDmXAu) | White | SleepMindEngine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/PBPEKKp3) | Black | crocodile_bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/QpPao7Sr) | Black | BlueMoonBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/qtjfcmCI) | White | CookieBot711 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Fd9TKkJ2) | Black | Athena-ChessBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7zlgD0un) | White | Ballin-bot | 180+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/5R1Xa9ay) | White | incipit_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Mfi4JiuZ) | Black | FSChess-Engine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Bwo78ywt) | White | ratsu-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ZsXkGPEs) | Black | Ballin-bot | 180+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/KceAXtSE) | White | thunfisch-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/3Vs4DQ30) | Black | JongsungKim | 180+2 | loss | Normal |
| [game](https://lichess.org/game/RN2Chbl1) | Black | mate1-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ruzTzMtj) | Black | SeedChess | 180+2 | loss | Normal |
| [game](https://lichess.org/game/9FT7vIEz) | Black | LegoTechnicControlPl | 180+2 | loss | Normal |
| [game](https://lichess.org/game/7hFQYQ1d) | White | Casanchess-NNUE | 180+2 | loss | Normal |
| [game](https://lichess.org/game/8ojddlqN) | White | MutusBot | 180+1 | win | Normal |
| [game](https://lichess.org/game/mBm0HUvS) | Black | manu-bot | 600+5 | win | Normal |
| [game](https://lichess.org/game/t0mFV93D) | Black | Bytefish | 300+5 | win | Normal |
| [game](https://lichess.org/game/u4oAJWu0) | White | DingusBot | 180+0 | win | Normal |
| [game](https://lichess.org/game/Z4Zg7bYn) | White | Supra_engine | 900+5 | win | Normal |
| [game](https://lichess.org/game/PVGNjS34) | Black | BockFish | 300+2 | loss | Normal |
| [game](https://lichess.org/game/yRhGYK9D) | Black | patgro986_bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/AVlXnr6L) | White | myChessJava | 900+10 | loss | Normal |
| [game](https://lichess.org/game/9CA9Sz2w) | White | littleply | 180+2 | win | Normal |
| [game](https://lichess.org/game/TiOrLhnj) | Black | sunfish-nnue-engine | 300+3 | win | Normal |
| [game](https://lichess.org/game/SgEMWySj) | White | Cosmic_Engine | 180+1 | loss | Normal |
| [game](https://lichess.org/game/7RPYss9J) | White | Z-O_AspenLabs | 900+5 | loss | Normal |
| [game](https://lichess.org/game/bB4Gj7g3) | Black | jester7777 | 60+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/VJoWDHPC) | Black | manu-bot | 600+3 | win | Normal |
| [game](https://lichess.org/game/gX4ghUo0) | Black | Cosmic_Engine | 180+1 | loss | Normal |
| [game](https://lichess.org/game/gDJalA0A) | White | Joey_ChessEngine | 60+1 | win | Normal |
| [game](https://lichess.org/game/GmVBYy2I) | Black | KirinChessBot1500 | 60+5 | win | Normal |
| [game](https://lichess.org/game/IK3YVmhF) | White | Joey_ChessEngine | 60+1 | draw | Normal |
| [game](https://lichess.org/game/unOdJ6Ba) | White | manu-bot | 300+5 | win | Normal |
| [game](https://lichess.org/game/Y3KAsmwx) | Black | CarlsenStyleBot | 300+3 | win | Normal |
| [game](https://lichess.org/game/eWmuPNmR) | White | Z-O_AspenLabs | 600+5 | draw | Normal |
| [game](https://lichess.org/game/Ustisvmm) | White | Sloth_Engine | 180+0 | win | Normal |
| [game](https://lichess.org/game/YGbXhDju) | Black | Joey_ChessEngine | 60+1 | win | Normal |
| [game](https://lichess.org/game/rPaHjnoM) | Black | Botkingteam33 | 180+2 | draw | Insufficient material |
| [game](https://lichess.org/game/rLbugNb3) | Black | ZaatarTheCat | 300+3 | win | Normal |
| [game](https://lichess.org/game/mIlDvGLy) | Black | Jibbby | 180+2 | draw | Normal |
| [game](https://lichess.org/game/EnxvOcoj) | Black | SquareFish-900 | 180+2 | win | Normal |
| [game](https://lichess.org/game/HVVQVqOZ) | Black | rust-in-pieces | 180+2 | loss | Normal |
| [game](https://lichess.org/game/wwxRtPYp) | Black | karl-engine-v1 | 180+2 | win | Normal |
| [game](https://lichess.org/game/35s0eXdt) | White | Blujay-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/gNQRBta1) | Black | sseh-c | 180+2 | loss | Normal |
| [game](https://lichess.org/game/FlBFkERu) | Black | glorbo_evermorrow | 180+2 | loss | Normal |
| [game](https://lichess.org/game/p74AX63h) | White | Stronggoliath | 180+2 | loss | Normal |
| [game](https://lichess.org/game/2WGLvZUk) | Black | Cloudia-Chess | 180+2 | win | Normal |
| [game](https://lichess.org/game/vMJ64VjF) | Black | Cloudia-Chess | 180+2 | win | Normal |
| [game](https://lichess.org/game/iJ00aKHL) | Black | kbelludoo | 180+2 | win | Time forfeit |
| [game](https://lichess.org/game/Z7uKj4jC) | White | TheChessblazer | 180+2 | win | Normal |
| [game](https://lichess.org/game/yXWZOlKn) | Black | FenrirBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/2QjRexZ8) | Black | bernstein-4ply | 180+2 | win | Normal |
| [game](https://lichess.org/game/VMrpzeys) | Black | Blujay-bot | 180+2 | win | Normal |
| [game](https://lichess.org/game/95LVf0TO) | Black | TopasBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/Am5FLZVf) | White | MDBOT | 180+2 | loss | Normal |
| [game](https://lichess.org/game/4vMIlqpg) | Black | Alarm2001 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/dS5HniCN) | White | zion-bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/3GmZfeIU) | Black | MutusBot | 180+2 | win | Normal |
| [game](https://lichess.org/game/tppTJbt0) | White | Bongaclang | 180+2 | loss | Normal |
| [game](https://lichess.org/game/DijN8b0i) | White | Benzene-Phenol-02 | 180+2 | unfinished | Abandoned |
| [game](https://lichess.org/game/yhO4ggxY) | Black | ailedbot | 180+2 | win | Normal |
| [game](https://lichess.org/game/YMHob41k) | Black | NibbleKnight | 180+2 | loss | Normal |
| [game](https://lichess.org/game/5RbuISTD) | White | R0bspierre | 180+2 | win | Normal |
| [game](https://lichess.org/game/2PKiqrUW) | White | LouisChess48-6K | 180+2 | win | Normal |
| [game](https://lichess.org/game/rlaSxu4x) | White | walpurti | 180+2 | loss | Normal |
| [game](https://lichess.org/game/6Nup5qBd) | Black | MilanBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/3yHc95pS) | White | Ndr1sn | 180+2 | loss | Normal |
| [game](https://lichess.org/game/JHf8LHRa) | White | TroutBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/KtVulsnk) | White | SleepMindEngine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/23V05AWI) | Black | SleepMindEngine | 180+2 | loss | Normal |
| [game](https://lichess.org/game/66HTxTj6) | White | BockFish | 300+0 | loss | Normal |
| [game](https://lichess.org/game/E6gj3NUf) | Black | CassidyChess-Bot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/ybP5TfSN) | Black | Chessiverse_Raak | 180+2 | win | Normal |
| [game](https://lichess.org/game/59zOq6FN) | White | Quirkafleeg-ZX48K | 180+2 | win | Normal |
| [game](https://lichess.org/game/SLxbGbFt) | White | SardineBot | 180+2 | loss | Normal |
| [game](https://lichess.org/game/JJCmtTvg) | Black | Grindfish | 180+2 | loss | Normal |
| [game](https://lichess.org/game/50kq2J5b) | Black | baby_eubos | 180+2 | win | Normal |
| [game](https://lichess.org/game/SwnETCnu) | White | Cosmic_Engine | 180+1 | loss | Normal |
| [game](https://lichess.org/game/q2mkXBUe) | Black | spyros1bot | 420+2 | draw | Normal |
| [game](https://lichess.org/game/abAXTpmd) | Black | Angela_eng | 600+3 | win | Normal |
| [game](https://lichess.org/game/jmVmUcjT) | White | Lightning-762 | 180+2 | loss | Normal |
| [game](https://lichess.org/game/m6B9i3bN) | Black | TrinculoBot | 180+1 | win | Normal |
| [game](https://lichess.org/game/WwybQSjn) | White | Perl-GigaChess | 180+2 | win | Normal |
