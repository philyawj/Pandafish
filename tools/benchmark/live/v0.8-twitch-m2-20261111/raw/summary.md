# Lichess session report

Engine: Pandafish v0.8

Engine SHA-256: `1baf9468c3634c565964b4cb0485fe65ff6b2d10fafa86ce050169a430723302`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 38 / 0 |
| Engine processes | 39 |
| Searches / completed / cancelled / incomplete or illegal | 1531 / 1531 / 0 / 0 |
| Searches linked to PGN | 1531 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 1263 |
| Maximum recorded response | 36772.187 ms |

Maximum hard-budget overrun: 0.104 ms. Response p95: 9995.221 ms.

| Completed depth | Searches |
| --- | --- |
| 4 | 21 |
| 5 | 382 |
| 6 | 793 |
| 7 | 264 |
| 8 | 64 |
| 9 | 6 |
| 13 | 1 |

| Stop reason | Searches |
| --- | --- |
| deadline | 1263 |
| soft_limit | 268 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores give a sign without distance. PV currently contains only the chosen root move.

No recording warnings.
