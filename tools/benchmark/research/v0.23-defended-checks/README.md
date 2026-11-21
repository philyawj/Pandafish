# Defender-aware checking-pressure result: reject

The first contextual king representation is not recommended for promotion.
It improves the old public decision sample but fails newer actual 100ms choices
for both root turns. No strength games ran, no weights were adjusted to rescue
the result and no fresh final controls were consumed. Playing code is restored
to v0.22 before the next hypothesis.

| Compact implementation, actual equal-time choices | v0.22 mean loss | Candidate mean loss | Result |
| --- | --- | --- | --- |
| Public 100ms, three trials | 45.687 cp | 41.164 cp | Pass |
| Public 1000ms, one trial | 54.851 cp | 50.119 cp | Pass |
| Newer 100ms, one trial | 30.679 cp | 36.115 cp | **Fail** |

The newer gate allows 35.679 cp. Black changes 25.431→30.446 and White
35.848→41.697, each failing its existing tolerance. Severe cases stay three;
there are no new losing-mate flags or timing overruns. The driver exits one
and does not run newer-1000, fresh controls or games. Preserve that negative
decision rather than treating public improvement as generalization.

The initial implementation's public-100 diagnostic also passes, 45.896→41.164.
Its [original source](candidate-source.json) remains separate from the
[compact source](compact/candidate-source.json). The
[exact compaction note](COMPACT_PLAN.md) records a structurally zero weak-ring
expression and its removal; no active weak-square feature is claimed. Safe
checking destinations use geometric defender and double-attack context, with
phase/resource-dependent pawn shelter. These remain approximations, not legal
checking moves or mate proofs.

Correctness controls pass blocker geometry, pawn edges, color mirrors,
perspective, bounds, terminal/material invariants, parent immutability and zero
allocations. The unchanged base score still matches its independent oracle.
The isolated component median rises about 53.01→371.9 ns. In the newer timed
sample, mean nodes fall 124,541→89,856 (27.9%) and completed depth falls
5.021→4.855. These are actual decision-workload costs, not a sevenfold whole-
engine slowdown. The provisional startup-overlap cost sample is retained but
does not support the isolated timing claim.

The outcome supports neither a generic king bonus nor the conclusion that
contextual evaluation is exhausted. Richer local features can improve inspected
old positions yet lose depth and newer decisions. No evidence attributes this
failure to transposition policy. A different representation needs more than a
coefficient tweak. Reassess the earlier checking-interface and exact backend
threads alongside the live-inspired evaluation route.

Reproduction uses `candidate-replay-source.json` with base tag v0.22 through
`go run ./tools/research replay --source PATH --out NEW-DIRECTORY --check`.
The compact [native checks](compact/native-checks.json) record exact commands,
identities and exits; all raw queries, cache provenance and referee results
remain in the three completed suite directories. No code is promoted from
this archive, and no publication is authorized by the ongoing investigation.

The original full snapshots are unchanged. `source-format-repair.json` records
the compatible overlay's omission of two unchanged nested module metadata
files, each verified against v0.22; runtime source is identical. The shared
replay format accepts those files from the base tag rather than an overlay.
