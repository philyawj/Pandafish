# Lichess session report

Engine: Pandafish v0.15

Engine SHA-256: `26f9af01670cfcb1bae9c261adf731cf7e2726034b717323130f8d0422abfdca`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 269 / 4 |
| Engine processes | 270 |
| Searches / completed / cancelled / incomplete or illegal | 12091 / 12091 / 0 / 0 |
| Searches linked to PGN | 12021 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 9675 |
| Maximum recorded response | 37462.405 ms |

Maximum hard-budget overrun: 1.132 ms. Response p95: 12211.459 ms.

| Completed depth | Searches |
| --- | --- |
| 4 | 17 |
| 5 | 751 |
| 6 | 4480 |
| 7 | 4570 |
| 8 | 1489 |
| 9 | 494 |
| 10 | 182 |
| 11 | 75 |
| 12 | 25 |
| 13 | 4 |
| 14 | 1 |
| 15 | 1 |
| 64 | 2 |

| Stop reason | Searches |
| --- | --- |
| deadline | 9675 |
| depth | 2 |
| soft_limit | 2414 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores are sign-only through v0.14.4 and signed moves to mate from v0.15 onward. PV currently contains only the chosen root move.

## Recording warnings

- some searches have missing or ambiguous PGN linkage; see moves.jsonl
