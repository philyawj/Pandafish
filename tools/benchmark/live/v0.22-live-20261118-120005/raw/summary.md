# Lichess session report

Engine: Pandafish v0.22

Engine SHA-256: `b1f7f1067c1cacbe6fb7b88264286c59fe150d123cdd603c7f6909ead1eb813d`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 105 / 5 |
| Engine processes | 106 |
| Searches / completed / cancelled / incomplete or illegal | 5160 / 5160 / 0 / 0 |
| Searches linked to PGN | 5084 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 4130 |
| Maximum recorded response | 37461.449 ms |

Maximum hard-budget overrun: 0.033 ms. Response p95: 12559.799 ms.

| Completed depth | Searches |
| --- | --- |
| 6 | 15 |
| 7 | 578 |
| 8 | 1872 |
| 9 | 1319 |
| 10 | 752 |
| 11 | 439 |
| 12 | 58 |
| 13 | 58 |
| 14 | 44 |
| 15 | 17 |
| 16 | 6 |
| 17 | 1 |
| 64 | 1 |

| Stop reason | Searches |
| --- | --- |
| deadline | 4130 |
| depth | 1 |
| soft_limit | 1029 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores are sign-only through v0.14.4 and signed moves to mate from v0.15 onward. PV currently contains only the chosen root move.

## Recording warnings

- some searches have missing or ambiguous PGN linkage; see moves.jsonl
