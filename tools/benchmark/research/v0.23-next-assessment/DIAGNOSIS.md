# Baseline depth and leaf diagnosis

All six fresh-table, full-history v0.22 searches reproduce the live move and
score at the recorded completed depth. Every requested next depth also completes;
no 30-second watchdog is reached. These are fixed-depth diagnostics, not timed
comparisons or exact reproduction of the live warm table.

| Game | Live depth / choice | Next depth / choice |
| --- | --- | --- |
| gNQRBta1 | 8 / Kf8, −16 | 9 / Nxe4, −8 |
| IK3YVmhF | 8 / Nd2, −11 | 9 / Nd2, −10 |
| PV4oUeoe | 7 / Qf3, +25 | 8 / Ng4, +4 |
| 6Nup5qBd | 7 / f6, +159 | 8 / Nh6, +132 |
| Mido8Dw4 | 8 / Kd2, −4 | 9 / Kd2, −1 |
| u4xtX1iv | 7 / Ne1, +92 | 8 / Rc2, +89 |

Four choices change and two persist. Changed moves are not assumed to be good:
the existing confirmation labels do not evaluate these new alternatives.
All six next-depth moves differ from the
teacher reference, but teacher agreement is not a correctness condition.

Four separately restricted diagnostic branches recover legal exact-score
stand-pat witnesses through unchanged production searches. They do not restrict
the playing engine. At depth eight, Kf8 and O-O both score −16; their witness
leaves differ by a rook on h8 versus e8, but material and coarse placement cannot
distinguish that activity. The Kf8 root is therefore not uniquely preferred.
Nd2 scores −11 while the teacher's Ne5 branch scores −40. Both witnesses have
equal material, and the original placement difference is exactly the branch
scalar. Every recovered leaf scalar agrees with root parity; no unverified PV
is presented as unique or forced. Source, paths, FENs, nodes and all runs remain
in this directory.

This establishes a coarse evaluation limitation, not that all six mistakes
share one cause or that another generic king penalty will succeed. The old
attacker-count-only term and fitted 12-feature models already failed. A new
defender-aware checking-pressure representation has a different mechanism;
actual-budget controls must determine whether it is useful. The Nd2 witness
also illustrates a horizon alternative: its selected line does not follow the
referee's Bxh2+ attack. Neither representation nor checking search is proven
as the complete remedy.
