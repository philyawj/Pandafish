# Lichess session report

Engine: Pandafish v0.22

Engine SHA-256: `b1f7f1067c1cacbe6fb7b88264286c59fe150d123cdd603c7f6909ead1eb813d`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 33 / 0 |
| Engine processes | 34 |
| Searches / completed / cancelled / incomplete or illegal | 1513 / 1513 / 0 / 0 |
| Searches linked to PGN | 1512 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 1217 |
| Maximum recorded response | 23711.444 ms |

Maximum hard-budget overrun: 0.030 ms. Response p95: 9995.186 ms.

| Completed depth | Searches |
| --- | --- |
| 6 | 1 |
| 7 | 229 |
| 8 | 686 |
| 9 | 395 |
| 10 | 140 |
| 11 | 37 |
| 12 | 11 |
| 13 | 7 |
| 14 | 4 |
| 15 | 2 |
| 64 | 1 |

| Stop reason | Searches |
| --- | --- |
| deadline | 1217 |
| depth | 1 |
| soft_limit | 295 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores are sign-only through v0.14.4 and signed moves to mate from v0.15 onward. PV currently contains only the chosen root move.

## Recording warnings

- some searches have missing or ambiguous PGN linkage; see moves.jsonl
