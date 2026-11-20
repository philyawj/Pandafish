# Lichess session report

Engine: Pandafish v0.22

Engine SHA-256: `b1f7f1067c1cacbe6fb7b88264286c59fe150d123cdd603c7f6909ead1eb813d`

| Measure | Value |
| --- | --- |
| Session status | finished |
| Saved games / unfinished | 19 / 0 |
| Engine processes | 20 |
| Searches / completed / cancelled / incomplete or illegal | 834 / 834 / 0 / 0 |
| Searches linked to PGN | 833 |
| Searches without final metrics | 0 |
| Legal fallback moves | 0 |
| Searches past engine hard budget | 681 |
| Maximum recorded response | 32212.466 ms |

Maximum hard-budget overrun: 8.420 ms. Response p95: 9997.771 ms.

| Completed depth | Searches |
| --- | --- |
| 6 | 2 |
| 7 | 107 |
| 8 | 402 |
| 9 | 237 |
| 10 | 53 |
| 11 | 18 |
| 12 | 12 |
| 13 | 3 |

| Stop reason | Searches |
| --- | --- |
| deadline | 681 |
| soft_limit | 153 |

Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores are sign-only through v0.14.4 and signed moves to mate from v0.15 onward. PV currently contains only the chosen root move.

## Recording warnings

- some searches have missing or ambiguous PGN linkage; see moves.jsonl
