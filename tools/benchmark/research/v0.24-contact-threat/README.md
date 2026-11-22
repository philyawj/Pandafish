# Protected queen-contact threat: recommended candidate

The [frozen plan](PLAN.md) selects a specific legal quiet queen-contact mating
threat, then searches actual defenses. The checking destination must be adjacent
to the defender king and protected by an attacker pawn, knight or king. A pass
selects a threat only; it never supplies a loss score. Every actual defense must
leave the declared mating move legal and mate before its mate score is used.
[Implementation](IMPLEMENTATION.md) and [process](PROCESS.md) define scope.

The activated frontier applies identically in full/narrow windows, disables
stand-pat during defense, searches all legal replies without SEE omissions and
uses one verification event per q path. Two remaining ordinary plies are required.
Capture, countercheck, king movement, protection removal and immediate mate can
refute the threat. Unresolved children retain the original q heuristic. This is
a bounded tactical contract, not a complete quiet-threat solver. No rejected
compact, material, move-count or learned evaluator is reactivated.

## Correctness and realistic cost

Training diagnosis covers 11 distinct saved threats with zero false certificates
across 148 cold/cached comparisons. All 596 independent move/state vectors pass,
including 17 generated contact moves. Both modules pass normal/race/vet/builds.
A complete legal-response reference and public immutable child transitions check
forced-mate positives, capture/countercheck refutations, mirrored colors,
full/narrow windows, stalemate, budgets, repetition/history and cancellation.
The initial scratch harness compile failure/source remain preserved.

Recognition-only shadow has 0.883% overhead and exact baseline moves/scores/nodes.
The separate instrumented ledger records 53,508 certified occurrences at
25,813,637 stand-pat cutoffs, with 348,626 geometric guards. Those occurrences
are not independent mistakes or forced root losses. Ledger time is not a cost
claim; its trajectory matches the original 1,145 roots and 43,118,858 nodes.

Three activated warm pairs pass the prospective 10% elapsed/1.25 early node
ceilings: **5.180% median overhead**, **0.987% more search nodes**. All 1,145 roots
complete; 19 scores, 11 choices and 406 rows change deterministically. Every
changed row remains in `candidate-warm/changes.json`. Elapsed includes direct
certificate queries, which are not recursively counted search nodes.

## Existing actual-budget controls

| Suite | Mean finite reference loss, baseline → candidate (cp) | Gate |
| --- | ---: | --- |
| Public 100ms | 45.896 → 45.896 | Pass |
| Public 1000ms | 54.851 → 54.851 | Pass |
| Newer 100ms | 30.794 → 31.038 | Pass |
| Newer 1000ms | 24.212 → 23.515 | Pass |

No new losing-mate flags or timing overruns. The newer short-budget suite worsens
for both turns, with retained 119cp Black and 5cp White adverse rows. Its longer
budget improves overall through Black. These small surrogate results do not
establish general strength or explain every changed choice.

## Frozen fresh controls

The single seed 2026112205/stream127 cohort retains 138 roots, 69 per turn,
across all16 families. Collection passes with561 pinned pre-existing paths.
Exclusions cover historical/public/newer/prior-fresh, warm/oracle/live, all three
previous complete screens and22,583 contact witness/control cases via temporary
hashed adapters. No query/cohort retry or duplicate tracked trajectory.
This cohort is consumed. Both fresh gates pass: mean finite loss improves
23.279→22.852cp at100ms and stays23.221cp at1000ms. No new losing-mate flags or
overruns; a67cp White short-budget regression remains recorded.

## Single standard screen and recommendation

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.23 | 24 / 6 / 20 | 0 | 0 | 54% |
| Stockfish Skill 3 | 33 / 0 / 17 | 0 | 0 | 66% |
| Pinned Sacrifice | 35 / 7 / 8 | 0 | 0 | 77% |

The original greater-than-50% direct gate passes. All150 games finish and replay
legally; no forfeits occur. Candidate timing maxima are97.6/99.1/97.3ms, with
zero candidate attempts above100ms. Settings remain100ms requests/50ms grace,
400-ply cap, unchanged ten openings, three concurrent matchups on AppleM2/AC.
No deliberate heavy work overlaps. [Complete report](screen/README.md),
[legal replay](games-audit.json) and all raw games/results are preserved.

All87 source files and binary identities match before and after play. Source SHA:
`447af5b0094a0319c92a1eea9c685cf16789204993556a6182e585025426ad5c`;
candidate binary SHA:
`ce4d89de73f89831c905e8810deb3384386e0cf73fd838e66475d8b8c995f683`.
The UCI name remains v0.23; development labels/hashes identify the tested build.

**Recommend promotion review as a bounded tactical-reliability improvement.**
Its local legal-defense behavior is established and all predeclared correctness,
cost, native/fresh and practical game gates pass without exceptions. The54%
direct result is a narrow small-sample lead, selected after preserved unsuccessful
attempts; it does not establish a statistically reliable general gain, precise
Elo or causal attribution of particular wins. Separate historical external
screens are not simultaneous baseline controls. The5.18% cost and adverse native
rows remain limitations. A separately planned broader confirmation would be
needed for a stronger strength claim; it does not start automatically.

Research source/evidence remain uncommitted for review. Cleanup PR #42 remains
merged without a release, and `handoff.md` is absent from that commit. No research
publication, release or live play follows without separate authorization.
