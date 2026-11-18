# Deferred main-search move generation over v0.21

**Decision: archive without promotion.** The [complete screen](RESULTS.md)
scores 47% against v0.21, 68% against Skill 3 and 62% against Sacrifice. All 150
games complete without forfeits. The failed direct gate and the earlier native
surrogate failure remain preserved despite the exact-work efficiency gain.

This candidate implements an exact staged main-search move picker. The previous
pipeline constructs the full legal list before trying a preferred move; the new
pipeline validates that hint directly and constructs tactics, quiet killers and
ordinary quiet moves only when each stage is reached. Full move order is
preserved. All earlier failed experimental code is inactive.

The [post-v0.21 assessment](../v0.22-next-assessment/README.md) records the changing
ranking and rejected attempts. The [opportunity diagnosis](../v0.22-next-assessment/PICKER_DIAGNOSIS.md)
finds that roughly 80% of generated main-search quiet records could be deferred
past early cutoffs. This identifies avoidable work without assuming a strength
gain. [Implementation details](IMPLEMENTATION.md) describe ordering, ownership,
direct validation, special moves and the malformed-castling compatibility path.

The source baseline is released v0.21 at
`494604686a7afbda256da1625aa907baa6060084`. The tested development label is
`v0.22-staged-picker-dev`; its UCI identity remains v0.21. Frozen
`candidate-source.json` and `source-hashes.json` contain all five Go overlays.
It is an uncommitted candidate, not a released v0.22 binary.
Its tested executable SHA-256 is
`e96f2c19001506d7d51d3d607f11abd62c53913f16eee688f0bae80ad306ce74`.

## Correctness and efficiency

[Controls and methods](TESTING.md) include 406,674 direct lookup/subset cases,
10,300 exact whole-Move ordered lists, per-ply ownership, preferred-only zero
allocation, startpos depth-four perft, 596 independent frozen rule vectors and
the pinned upstream root-package race suite. Existing transition, UCI,
cancellation, mate and history controls remain active. Full race, vet and build
checks pass. The initial malformed-rights compatibility failure and its complete
source are preserved in [ERRATA.md](ERRATA.md).

| Fixed-work cohort | Search-time reduction | Search result |
| --- | --- | --- |
| Public depth four, three alternating pairs | 15.76% median | Identical moves, scores and nodes. |
| Public depth three | 7.29% | Identical moves, scores and nodes. |
| Twelve deeper roots, depth six | 23.72% | Identical moves, scores and nodes. |
| Newer roots, depth four | 14.73% | Identical moves, scores and nodes. |
| Twenty-game warm trajectory, depth five, three pairs | 19.81% median | All 1,145 roots and 43,118,858 nodes per actor agree. |

All 325 distinct cold comparisons agree. Warm pairs retain one table per game,
reconstruct the recorded history and preserve the same legal trajectory. The
independent `warm-review.json` checks all 3,435 paired observations, input/helper
hashes, repeated actor results and baseline-best membership through equality.
Allocation falls only about 1.1%; the main measured benefit is avoiding generation
and ordering work. These are selected workload efficiency results, not Elo.

## Equal-time native decisions

All 1,252 original actor queries complete. The independent `choice-review.json`
recomputes saved legal choices, finite losses, mate flags and timings. The fixed
teacher is Stockfish 19, full Skill 20, one thread, Hash 64, MultiPV 1 and
100,000 nodes. This diagnostic teacher differs from the Skill 3 game opponent.
No production root is restricted to a reference move.

| Suite | Finite roots | Mean loss: baseline → candidate (cp) | Severe cases: baseline → candidate | Original gate |
| --- | --- | --- | --- | --- |
| Public 100ms, three trials | 67 | 40.0448 → 45.9552 | 1 → 2 | **Fail**, limit 45.0448. |
| Public 1000ms | 67 | 52.5821 → 54.8507 | 3 → 3 | Pass. |
| Newer 100ms | 131 | 35.0840 → 31.7634 | 3 → 3 | Pass, both turns. |
| Newer 1000ms | 132 | 30.0152 → 29.8864 | 2 → 2 | Pass, both turns. |

There are no new losing-mate flags, timing overruns or query errors. The failed
public 100ms result is not rerun or marked passed. Its Black-only diagnostic
tolerance also fails; that turn gate is not required by the original public
protocol. The largest loss change is the baseline's own depth-five Qd4 choice
on HhPI3vCx, reached by the faster engine. All five changed roots, including the
improving choice, agree exactly between actors at every observed depth (ten
unrestricted comparisons). This explains those observations without excusing
the worse choice or proving all future behavior equivalent.

[The diagnosis](NATIVE_DIAGNOSIS.md) and
[explicit screen decision](SCREEN_DECISION.md) preserve the failed prerequisite
and explain the bounded research exception before games. A game screen cannot
retroactively make the original native gate pass. No threshold, root-specific
behavior or actor observation is changed.

## Reproduction and evidence

Restore the candidate into a new directory using the Go research replay tool:

```sh
go run ./tools/research replay \
  --source tools/benchmark/research/v0.22-staged-picker/candidate-source.json \
  --out /tmp/pandafish-staged-picker-replay --check
```

`timing_driver.go.txt`, `warm_pairs_driver.go.txt`, the assessment's native
drivers and `screen_driver.go.txt` preserve commands and settings. Their saved
records are the original observations; fresh timings/games will differ and must
use new output paths. Warm profiling requires the separately archived Go helper
and recorded v0.21 game input, not only the five runtime overlays. Native teacher
cache arguments are prior run directories. The independent native/screen review
programs are the existing Go reviewers referenced by the assessment.

No opening book, online dataset or live Lichess play is used. Earlier failed
code remains archived separately and is absent from the current executable.
No automatic larger confirmation, commit, PR or release follows this research.

The same roots and opening families have informed earlier development; they are
not a fresh holdout. A single 50-game row has wide uncertainty, and the ten
repeated opening pairs do not provide 50 independent position families. Three
concurrent matchups share host scheduling and can change depth boundaries. The
matching-depth evidence supports functional equivalence on tested work, not a
universal proof of identical elapsed behavior or a fix for existing horizon
mistakes. The failed public surrogate remains a concrete limitation.
