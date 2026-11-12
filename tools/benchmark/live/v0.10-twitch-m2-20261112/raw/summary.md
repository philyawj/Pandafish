# Lichess session report

Engine: Pandafish v0.10

Engine SHA-256: `8a5829be45932c8d6923b5a90fa2d151e7d74f9386564c55e501c059c4e7b1a9`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 14 / 1 |
| Engine processes | 18 |
| Searches / completed / cancelled / incomplete or illegal | 550 / 550 / 0 / 0 |
| Searches linked to PGN | 510 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 440 |
| Maximum recorded response | 9995.400 ms |

Maximum hard-budget overrun: 0.055 ms. Response p95: 7273.105 ms.

| Completed depth | Searches |
| --- | --- |
| 5 | 53 |
| 6 | 207 |
| 7 | 170 |
| 8 | 73 |
| 9 | 23 |
| 10 | 11 |
| 11 | 10 |
| 12 | 1 |
| 13 | 2 |

| Stop reason | Searches |
| --- | --- |
| deadline | 440 |
| soft_limit | 110 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores give a sign without distance. PV currently contains only the chosen root move.

## Recording warnings

- some searches have missing or ambiguous PGN linkage; see moves.jsonl
