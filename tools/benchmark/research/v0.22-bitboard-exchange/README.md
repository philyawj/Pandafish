# Exact search-kernel revision over v0.21

**Recommendation: promote after review.** The [validated standard screen](RESULTS.md)
scores **64% against v0.21, 71% against Skill 3 and 70% against Sacrifice**,
with no unfinished games or forfeits. All 150 games replay legally. The public
100ms surrogate still fails and remains disclosed. No commit or PR is prepared.

This candidate defers unused legal-move work and replaces repeated exchange-board
square scans with exact local bitsets. It preserves the original complete move
order, exchange heuristic, material/placement values, search/TT/history policy
and clocks. The practical screen supports promotion; efficiency is a separate
established result on the selected workloads, without a precise Elo claim.

The [assessment](../v0.22-next-assessment/README.md) records the changing ranking,
failed attempts and new diagnoses. The generation foundation reuses the exact
first staged picker and evasion prefilter. Their original outcomes remain failed:
47% direct for the first picker; 9.11% extra warm gain versus a 10% target for the
first evasion version. The new exchange kernel must add another 10% warm saving
over that already faster foundation before any new playing test. No unchanged
binary is rerun for a favorable score.

The baseline is released v0.21 at
`494604686a7afbda256da1625aa907baa6060084`. The tested development label is
`v0.22-bitboard-exchange-dev`; UCI identity remains v0.21. Candidate executable
SHA-256 is `c5486656b922db819cee1fa5e2851e9ad1b744806cb9d84a301df6a64d657505`.
Twelve frozen Go overlays and exact hashes are in `candidate-source.json` and
`source-hashes.json`. It is uncommitted research, not a released v0.22 binary.

## Implementation and correctness

[The implementation note](IMPLEMENTATION.md) describes local side/type bitsets,
nearest slider blockers, king-safety simulation/restoration, ascending-square
tie preservation and the unchanged exchange minimax. Earlier generation notes
explain the [staged picker](../v0.22-staged-picker/IMPLEMENTATION.md) and
[evasion mask](../v0.22-evasion-mask/IMPLEMENTATION.md). No clock, history,
aspiration, quiescence-reuse, child-storage or learned-evaluation prototype is
activated. Original upstream runtime files and licenses remain unchanged.

[Controls](TESTING.md) include 46,794 full exchange-score comparisons with the
frozen mailbox implementation, 33,152 attacker-set/least-attacker/restoration
comparisons, independent legal-recapture fixtures, special moves, pins and
cancellation. All generation controls remain active: 406,674 lookup/subset cases,
10,300 whole-Move ordered lists, 6,036 extra evasion comparisons, active-ply
isolation, first-hint zero allocation and perft. Root and pinned upstream races,
vet/build and all 596 independent frozen rule vectors pass. All 21 original
upstream runtime files remain byte-identical.

## Fixed-work efficiency

| Cohort | Search-time reduction | Result |
| --- | --- | --- |
| Public depth four, three alternating pairs vs v0.21 | 35.35% median | All moves, scores and nodes agree. |
| Public depth three vs v0.21 | 25.71% | All agree. |
| Twelve depth-six roots vs v0.21 | 36.94% | All agree. |
| Newer depth-four roots vs v0.21 | 38.19% | All agree. |
| Warm depth five, three pairs vs first evasion revision | 14.65% additional median | All 1,145 roots / 43,118,858 nodes per actor agree and repeat. |
| New warm depth-five pair vs tagged v0.21 | 37.84% total | Same roots, moves, scores and 43,118,858 nodes. |

All 325 distinct cold combinations agree. Warm work retains one table per game,
reconstructs recorded histories and follows the same original 20-game trajectory.
The independent `../v0.22-next-assessment/kernel-pairs-review.json` recomputes
ratios, node accounting, all phase pass/failure flags, input/helper hashes,
best-move membership through equality and repeated results. Generic raw pending-
membership fields are earlier checkpoints, resolved by that separate review.
The total tagged comparison is one fresh pair; the additional-gain median has
three pairs. Neither is a universal throughput bound or an Elo measurement.

Allocated bytes barely change versus the evasion foundation (about 12.5 GB
cumulative in each warm run). The mechanism saves computation rather than peak
memory. The independent source review also checks the tagged/staged/evasion/
exchange benchmark workspaces against their declared frozen files. Profiles
include setup/measurement overhead and are not used as elapsed-gain evidence.

## Actual-budget native choices

Independent `choice-review.json` recomputes all original raw actor/reference
observations. Teacher: pinned Stockfish 19, full Skill 20, one thread, Hash 64,
MultiPV 1, 100,000 nodes. Production roots stay unrestricted. Reference scores
use the existing close-alternative band and finite-loss/mate conventions; the
teacher is not assumed to supply the only acceptable move.

| Suite | Finite roots | Mean loss baseline → candidate (cp) | Severe cases baseline → candidate | Original gate |
| --- | --- | --- | --- | --- |
| Public 100ms, three trials | 67 | 40.0746 → 45.8955 | 1 → 2 | **Fail**, limit 45.0746. |
| Public 1000ms | 67 | 52.5821 → 54.8507 | 3 → 3 | Pass. |
| Newer 100ms | 131 | 34.9313 → 30.6794 | 3 → 3 | Pass, both turns. |
| Newer 1000ms | 132 | 30.2500 → 24.2803 | 2 → 1 | Pass, both turns. |

All 1,252 actor queries complete with no errors, illegal moves, overruns or new
losing-mate flags under the existing convention. Public 100ms and its Black-only
diagnostic mean remain failed; that public turn requirement was not originally
mandatory. No failed actor request is rerun. [The diagnosis](NATIVE_DIAGNOSIS.md)
retains all six changed roots, including improvements. Twelve unrestricted
matched-depth comparisons have identical moves, scores and nodes. The dominant
loss is the same depth-five Qd4 choice in both actors, reached at different timed
depths. It remains a worse teacher-labelled choice and a real limitation.

The [explicit pre-game exception](SCREEN_DECISION.md) authorizes one bounded
research screen after reviewing those failures, the other budgets, exactness and
the larger efficiency revision. It does not pass the original gate, waive a
mate/timing/legality failure or authorize promotion. The prior preparer compile
error and missing dependent file are separately preserved in [ERRATA.md](ERRATA.md).

## Reproduction and limits

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.22-bitboard-exchange/candidate-source.json \
  --out /tmp/pandafish-exchange-kernel-replay --check
```

Go timing/warm/native/screen drivers preserve the original commands and output
paths; reproduction uses new paths and can produce different elapsed/game
results. The separate warm helper and original recorded-game inputs are required
for those trajectories. Stockfish reference-cache arguments are prior run
directories, not their metadata files. Binaries are ignored; their identities
remain in methods and game summaries, and source restores from tagged v0.21.

The same roots/openings have informed previous development. They are inspected
regression controls, not a fresh holdout. Selection after failed attempts adds
uncertainty. Repeated opening families and concurrent host scheduling limit the
precision of a 50-game row. Functional equivalence on the tested work does not
prove every future position or establish a general strength gain. Finite teacher
labels do not prove a position was savable; absent new mate flags do not prove
tactical safety. No opening book, online dataset or live game is used. No larger
confirmation, commit, PR or release follows automatically.
