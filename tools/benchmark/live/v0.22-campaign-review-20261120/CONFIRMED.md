# Larger-budget confirmation

Stockfish 19, 500,000 nodes per independent query. Scores favor the root mover;
CP and signed mate distances are separate. Ply indices in JSON are zero-based.
The survey column shows the original 100,000-node CP loss, where finite.
This targeted selection is not a prevalence estimate or a unique-best-move test.

| Group | Game / move | Clock | Played | Reference / played score | CP loss, survey → confirmation | Live depth |
| --- | --- | --- | --- | --- | --- | --- |
| deterioration_opening | [wwxRtPYp](https://lichess.org/game/wwxRtPYp) 4... | 180+2 | Qe7 | -18 / -206 | 172 → 188 | 8 |
| deterioration_opening | [y2IkeSMw](https://lichess.org/game/y2IkeSMw) 5... | 180+2 | b5 | -67 / -218 | 165 → 151 | 8 |
| deterioration_opening | [YGbXhDju](https://lichess.org/game/YGbXhDju) 6... | 60+1 | g6 | 0 / -356 | 301 → 356 | 8 |
| deterioration_opening | [gNQRBta1](https://lichess.org/game/gNQRBta1) 7... | 180+2 | Kf8 | -67 / -277 | 180 → 210 | 8 |
| deterioration_opening | [U1Z9lTsh](https://lichess.org/game/U1Z9lTsh) 8. | 180+2 | e4 | -15 / -198 | 175 → 183 | 8 |
| deterioration_opening | [IK3YVmhF](https://lichess.org/game/IK3YVmhF) 9. | 60+1 | Nd2 | -54 / -713 | 330 → 659 | 8 |
| deterioration_opening | [PV4oUeoe](https://lichess.org/game/PV4oUeoe) 9. | 60+1 | Qf3 | 11 / -382 | 270 → 393 | 7 |
| deterioration_opening | [n5N2xU29](https://lichess.org/game/n5N2xU29) 9. | 180+2 | Ng5 | -70 / -232 | 161 → 162 | 7 |
| deterioration_middlegame | [6Nup5qBd](https://lichess.org/game/6Nup5qBd) 11... | 180+2 | f6 | -88 / -646 | 411 → 558 | 7 |
| deterioration_middlegame | [u4oAJWu0](https://lichess.org/game/u4oAJWu0) 12. | 180+0 | e5 | 307 / 126 | 161 → 181 | 7 |
| deterioration_middlegame | [Mido8Dw4](https://lichess.org/game/Mido8Dw4) 12. | 180+2 | Kd2 | -43 / -208 | 160 → 165 | 8 |
| deterioration_middlegame | [50kq2J5b](https://lichess.org/game/50kq2J5b) 12... | 180+2 | f6 | -97 / -321 | 245 → 224 | 7 |
| deterioration_middlegame | [DrjE4O2r](https://lichess.org/game/DrjE4O2r) 13. | 180+2 | Ng5 | 223 / 19 | 194 → 204 | 7 |
| deterioration_middlegame | [MhLL6txE](https://lichess.org/game/MhLL6txE) 13... | 180+2 | gxh3 | 185 / -30 | 217 → 215 | 8 |
| deterioration_middlegame | [2QjRexZ8](https://lichess.org/game/2QjRexZ8) 13... | 180+2 | Bxe2 | 66 / -112 | 158 → 178 | 7 |
| deterioration_middlegame | [p74AX63h](https://lichess.org/game/p74AX63h) 14. | 180+2 | Bxh7 | 269 / -185 | 426 → 454 | 7 |
| deterioration_late | [7RPYss9J](https://lichess.org/game/7RPYss9J) 31. | 900+5 | Qb3 | 180 / 0 | 212 → 180 | 9 |
| deterioration_late | [Y3KAsmwx](https://lichess.org/game/Y3KAsmwx) 32... | 300+3 | Qd8 | 0 / -179 | 170 → 179 | 9 |
| deterioration_late | [u4xtX1iv](https://lichess.org/game/u4xtX1iv) 33. | 180+2 | Ne1 | 19 / -554 | 515 → 573 | 7 |
| deterioration_late | [95LVf0TO](https://lichess.org/game/95LVf0TO) 34... | 180+2 | Rh8 | 81 / -494 | 513 → 575 | 7 |
| deterioration_late | [5khQm35W](https://lichess.org/game/5khQm35W) 35. | 180+2 | Rh7 | 263 / 42 | 192 → 221 | 8 |
| deterioration_late | [53HukkXX](https://lichess.org/game/53HukkXX) 40... | 600+2 | Rxg3 | 332 / 144 | 164 → 188 | 10 |
| deterioration_late | [jmVmUcjT](https://lichess.org/game/jmVmUcjT) 41. | 180+2 | g3 | 240 / 22 | 168 → 218 | 10 |
| deterioration_late | [YMHob41k](https://lichess.org/game/YMHob41k) 41... | 180+2 | h3 | 138 / -89 | 247 → 227 | 7 |
| conversion_mate | [nViUqeab](https://lichess.org/game/nViUqeab) 22. | 180+2 | Qxe5 | M6 / M10 | mate comparison | 8 |
| conversion_mate | [m6B9i3bN](https://lichess.org/game/m6B9i3bN) 37... | 180+1 | f5 | M8 / M9 | mate comparison | 8 |
| conversion_mate | [vMJ64VjF](https://lichess.org/game/vMJ64VjF) 37... | 180+2 | Bh4+ | M9 / M10 | mate comparison | 8 |
| conversion_mate | [35s0eXdt](https://lichess.org/game/35s0eXdt) 46. | 180+2 | a5 | M7 / M9 | mate comparison | 8 |
| disagreement_control | [WwybQSjn](https://lichess.org/game/WwybQSjn) 12. | 180+2 | Qh5+ | 368 / 552 | 0 → 0 | 8 |
| disagreement_control | [ZiWwAAsF](https://lichess.org/game/ZiWwAAsF) 20. | 180+2 | Qb2 | 1875 / 1875 | 0 → 0 | 7 |
| agreement_control | [1TrbIGaE](https://lichess.org/game/1TrbIGaE) 7. | 180+2 | dxc6 | -75 / -75 | 0 → 0 | 8 |
| agreement_control | [2PKiqrUW](https://lichess.org/game/2PKiqrUW) 7. | 180+2 | fxe5 | 799 / 799 | 0 → 0 | 8 |

## Legal continuation samples

These are referee lines, not claims that the played continuation was forced or
that every alternative is worse. Complete PVs and raw UCI output are preserved
in [confirmation-n500000.json](confirmation-n500000.json).

- **wwxRtPYp, zero-based ply 7:** played Qe7; reference `Nc6 cxd4 cxd4 Nxd4 Nf6 Nc3 Bb4 Bc4 Bxc3+ bxc3`. After the played move: `Qe7 Be2 Nc6 O-O Qd8 cxd4 cxd4 Nxd4 Nf6 Bf4`.
- **y2IkeSMw, zero-based ply 9:** played b5; reference `e6 Bxc4 Nf6 Qa4+ Bd7 Qb3 exd5 exd5 Qe7+ Kd1`. After the played move: `b5 a4 e6 axb5 Bb4+ Nc3 exd5 exd5 Nf6 Bxc4`.
- **YGbXhDju, zero-based ply 11:** played g6; reference `d5 fxe5 Bxf5 exf5 Qh4+ Ke2 Qg4+ Kd2 Qf4+ Ke1`. After the played move: `g6 fxe5 gxf5 Bf4 d6 exd6 cxd6 Nc3 Nf6 exf5`.
- **gNQRBta1, zero-based ply 13:** played Kf8; reference `O-O O-O b6 Ba3 Re8 cxd5 exd5 Qc2 h6 c4`. After the played move: `Kf8 cxd5 exd5 c4 Be6 cxd5 Bxd5 O-O Nb4 Ba3`.
- **U1Z9lTsh, zero-based ply 14:** played e4; reference `Qa4+ Nd7 d4 cxd4 Qxd4 Nxe5 Qxe5 Ne7 Qh5 Nf5`. After the played move: `e4 d4 Ne2 Nd7 Nc4 h5 a4 a5 f4 Qc7`.
- **IK3YVmhF, zero-based ply 16:** played Nd2; reference `Ne5 O-O Nxc6 Qe8 Ne5 Ba6 Re1 c5 f4 exf3`. After the played move: `Nd2 Bxh2+ Kh1 Bd6 Kg1 Ng4 g3 h5 Ncxe4 dxe4`.
- **PV4oUeoe, zero-based ply 16:** played Qf3; reference `Qe2 Qc7 f4 Bb7 Bb2 Nc6 Ne4 Nxe4 Bxe4 Nxe5`. After the played move: `Qf3 Bxe5 Qxa8 b4 Qf3 Qc7 Qh3 bxc3 dxc3 Nd5`.
- **n5N2xU29, zero-based ply 16:** played Ng5; reference `cxd4 cxd4 Bc2 Be7 d3 O-O Bb3 Nd7 Bd2 Nc5`. After the played move: `Ng5 b5 b3 Be7 Nc2 O-O Bb2 h6 Nf3 Be6`.
- **6Nup5qBd, zero-based ply 21:** played f6; reference `h5 Bc4 Nge7 Qf2 Nf5 Rhe1 d6 Kb1 Na5 Bd3`. After the played move: `f6 exf6 Qxf6 Nd5 Qg6 Nxc7 Qxg3 hxg3 Rb8 Bd6+`.
- **u4oAJWu0, zero-based ply 22:** played e5; reference `d3 Qc7 g3 a5 Bf4 Qa7 Qxc6 O-O Qc2 Ba6`. After the played move: `e5 Nd5 O-O a5 Re1 h6 h3 Qc7 Qe4 Bb7`.
- **Mido8Dw4, zero-based ply 22:** played Kd2; reference `O-O Nxe4 Qxe4 Qxc3 Rae1 O-O Re3 Qa5 Ne5 f6`. After the played move: `Kd2 c5 dxc5 O-O Bd3 Bd7 Qe5 Bc6 Nd4 Rad8`.
- **50kq2J5b, zero-based ply 23:** played f6; reference `a5 O-O Ba6 Qc2 h6 d5 cxd5 Rb1 Na7 Ba3`. After the played move: `f6 exf6 gxf6 Bh6 Qd5 Qc2 Qh5 Qd2 Nc3 Qxc3`.
- **DrjE4O2r, zero-based ply 24:** played Ng5; reference `c4 Rd8 Rfc1 a6 Bxg7 Qxg7 Nc3 Kc8 b4 g5`. After the played move: `Ng5 h6 Nf3 g5 Bxg7 gxf4 Bxf8 Rxf8 a4 a5`.
- **MhLL6txE, zero-based ply 25:** played gxh3; reference `g3 Nxe4 fxe4 Ra2 Na5 Qb5 b6 f4 exf3 gxf3`. After the played move: `gxh3 Rxh3 Na5 Qd1 Rf7 Rh6 Qe7 Qh5 Be6 Qh4`.
- **2QjRexZ8, zero-based ply 25:** played Bxe2; reference `Kb8 h3 Be6 Qe3 Rg8 Nd2 Ne7 Ng3 Qc7 a3`. After the played move: `Bxe2 Qxe2 exd4 cxb4 Bxb4 a3 Bd6 Nd2 Kb8 Nf3`.
- **p74AX63h, zero-based ply 26:** played Bxh7; reference `O-O Qb6 Qe3 Rg8 a3 Rg7 b4 O-O-O Kh1 Bd6`. After the played move: `Bxh7 Bf8 Qh5 Qf6 Nb5 O-O-O Qg5 Qxg5 Nxg5 f6`.
- **7RPYss9J, zero-based ply 60:** played Qb3; reference `Qf3 Qe8 d5 cxd5 Nd4 Qg6 c6 g4 hxg4 f4`. After the played move: `Qb3 Qe8 Qf3 Be6 Nc3 Bc4 Ra1 Kg8 Ra6 Qd7`.
- **Y3KAsmwx, zero-based ply 63:** played Qd8; reference `Ke7 Qe5+ Qe6 Qc7+ Qd7`. After the played move: `Qd8 Bg6 h5 Bg5 Qd1+ Kh2 Qd7 Bxf6 gxf6 Qxf6+`.
- **u4xtX1iv, zero-based ply 64:** played Ne1; reference `Kg2 Rde8 Qg4 Bd4 Nxd4 cxd4 Qxd4 Qxd4 Rxd4 Re2`. After the played move: `Ne1 Rde8 Qb7 Bd4 Nd3 g5 Kg2 Re7 Rxd4 cxd4`.
- **95LVf0TO, zero-based ply 67:** played Rh8; reference `Re7 Qxh6 Qg7 Qg5 c3 h5 Kh8 Rd3 Rb7 Rc1`. After the played move: `Rh8 Bxf5 Qf6 Qg4+ Kf7 Rd7+ Be7 Re1 exf5 Qh5+`.
- **5khQm35W, zero-based ply 68:** played Rh7; reference `Bg3 Bg5 Rg6 Bd2 Rg4+ Kf3 Nf6 Nf5 Rf1+ Ke2`. After the played move: `Rh7 e5 Rg1 Ne3 Rf7 Ne2 Rgg7 a5 b3 b4`.
- **53HukkXX, zero-based ply 79:** played Rxg3; reference `e4 Kxh6 Re7 g4 e3 Ra1 Rxc2 g5 d4 g6`. After the played move: `Rxg3 Rf7+ Kd6 Rxa7 Rxh3+ Kg6 Rg3+ Kh7 Rc8 Rf7`.
- **jmVmUcjT, zero-based ply 80:** played g3; reference `Ra6 Ke5 h4 Rc4 Ke2 Rc2+ Kd1 Rxg2 Rxb6 Ra2`. After the played move: `g3 h5 Ra6 Kf5 gxf4 Rxf4 Ke2 Ke5 Kd2 f5`.
- **YMHob41k, zero-based ply 81:** played h3; reference `Be6 Qxh4+ Kg8 Qh6 Bd7 Be5 Rh7 Qf4 Rh3 Bb8`. After the played move: `h3 Qh4+ Kg8 Be5 Rb7 Qh8+ Kf7 Qh7+ Ke8 Qh6`.
- **nViUqeab, zero-based ply 42:** played Qxe5; reference `Qg5+ Kh8 Qe7 Kg8 Ne4 Rh6 N4f6+ Rxf6 Qxf6 Rxb2`. After the played move: `Qxe5 f5 Nf6+ Kf7 Nxh7 Kg6 Qe7 Bb7 Nxd7 f4`.
- **m6B9i3bN, zero-based ply 73:** played f5; reference `b5 Kb1 Qe1+ Kb2 Qc3+ Kb1 g5 a4 bxa4 Kc1`. After the played move: `f5 a4 a5 Kb3 d3 cxd3 Qxd3+ Kb2 f4 Ka1`.
- **vMJ64VjF, zero-based ply 73:** played Bh4+; reference `b5 f5 b4 f6 Bxf6 h4 b3 Kf4 b2 Kf3`. After the played move: `Bh4+ Kf3 b5 f5 b4 f6 Bxf6 Kf4 b3 Kf5`.
- **35s0eXdt, zero-based ply 90:** played a5; reference `Rc6+ Kd7 b7 Bh7 Rb6 Bg8 b8=Q Bb3+ Kxb3 Ke7`. After the played move: `a5 Be6 a6 Ke5 a7 Bd5 Bxd5 Kf6 a8=Q Ke5`.
- **WwybQSjn, zero-based ply 22:** played Qh5+; reference `Nd3 Qc4 Qh5+ g6 Qe2 O-O-O dxe6 Nxe6 O-O-O Nd4`. After the played move: `Qh5+ g6 Qe2 Nh6 Ne4 Qa5+ b4 Qxb4+ c3 Qa5`.
- **ZiWwAAsF, zero-based ply 38:** played Qb2; reference `Qb2 d5 Qxc3+ d4 Qxd4+ Kg6 O-O-O f6 Bxa6 f4`. After the played move: `Qb2 d5 Qxc3+ d4 Qxd4+ Kg6 O-O-O f6 Bxa6 f4`.
- **1TrbIGaE, zero-based ply 12:** played dxc6; reference `dxc6 e4 cxb7 Bxb7 Ng1 Nf6 e3 b4 b3 Bd6`. After the played move: `dxc6 e4 cxb7 Bxb7 Ng1 Nf6 e3 b4 b3 Bd6`.
- **2PKiqrUW, zero-based ply 12:** played fxe5; reference `fxe5 Ng4 Nf3 c6 h3 Nh6 Bd3 cxd5 Nxd5 e6`. After the played move: `fxe5 Ng4 Nf3 c6 h3 Nh6 Bd3 cxd5 Nxd5 e6`.
