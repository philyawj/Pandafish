# Quiet-count scout reductions: direct pass, revise before promotion

This candidate follows the failed initial/material/frontier reduction versions.
The [plan](PLAN.md) distinguishes the new eligibility definition: count earlier
ordinary quiet choices, rather than every move before them. Captures, promotions
and checks cannot make the first ordinary quiet appear late. Protect the first
three ordinary quiets, castling, checks, tactical moves, killer hints, mate-band
windows and positions with total nonpawn phase below eight.

For an eligible later main-search child at parent depth >=3, reduce only its
PVS null-window probe by one ply. Any alpha improvement or mate-band result is
verified at full depth before the existing full-window PVS re-search. Every
legal move remains eligible. A reduced finite fail-low can still overlook a
useful move; completed nominal depth does not prove full-width optimality.
Material/evaluation, quiescence, table identity/depth/replacement and clock
formulas remain v0.22. Single-worker production has no archived parallel code.

Normal/race/vet/build controls pass through the shared validation workflow.
Focused controls cover earlier captures before quiets, protected moves, positive
verification, sparse material, parent/history restoration and cancellation.
The original independent mate proof passes. The killer-only audit isolates
unreduced ordering in both reference actors; its exact assertion remains.

Three alternating warm pairs retain all 1,145 game-history roots, all ten
openings and complete nominal depth five. Candidate values/nodes are
deterministic. Node ratio is **0.528777**, or **47.122% fewer nodes**; median search
elapsed ratio is **0.568719**, or **43.128% less time**. This is reduced selective
work, not an equivalent full-width-tree speedup. Original-tag full-depth review
of every selected choice reports mean finite regret **0.996cp** over 1,103 roots,
maximum **273cp**, with zero mate-risk flags/value contradictions. Raw adverse
cases remain. The fixed warm cohort is inspected development evidence.

| Equal-time suite | Mean finite loss, v0.22 → candidate | Gate |
| --- | ---: | --- |
| Public 100ms, three trials | 45.8955 → 46.4328cp | Pass |
| Public 1000ms | 54.8507 → 55.7463cp | Pass |
| Newer 100ms | 30.6794 → 32.5649cp | Pass |
| Newer 1000ms | 24.2803 → 23.0833cp | Pass |
| Fresh 100ms, three trials | 38.2689 → 30.7311cp | Pass |
| Fresh 1000ms | 34.2269 → 34.5462cp | Pass |

All six original gates pass, including both-turn requirements for newer/fresh
inputs. No exception is used. The fresh cohort uses frozen seed2026112121/
stream121: 130 retained roots, 65 per turn, from all 16 families, excluding all
consumed controls and live survey boards. No new losing-mate flags or timing
overruns occur. Means do not all improve: the fresh long-budget severe count
rises 1→2, and its largest regression is proof-fresh-07-ply13, 98→438cp loss.
These limitations remain visible despite the passed aggregate gate.

The one standard concurrent screen completes:

| Opponent | W / D / L | Score |
| --- | --- | ---: |
| Tagged v0.22 | 22 / 13 / 15 | **57%**, practical direct gate passes |
| Stockfish Skill 3 | 22 / 1 / 27 | **45%** |
| Pinned Sacrifice | 32 / 13 / 5 | **77%** |

All 150 games finish legally without forfeits. Candidate response maxima stay
below 98ms; no candidate attempt exceeds the requested 100ms. Opponent hashes,
Skill 3, one-thread/16MiB settings, openings and 100ms/50ms profile match the
prior v0.22 screen, which recorded 71% against Stockfish. Separate small screens
do not isolate an Elo change, but that adverse external signal is too large to
ignore. Direct pass and real efficiency do not establish broad playing gain.
**Recommend revision, not promotion**; do not repeat this unchanged screen.

The Go benchmark exits successfully and saves its Markdown report at the
screen root. The wrapper's erroneous nested *.md output expectation then fails.
That infrastructure failure stays recorded; no game is rerun. The legal-game
audit is executed separately and verifies all 150 records. Future screen plans
must expect the report at its actual root path.

Preserve this promising selective foundation for a separately frozen tactical-
horizon revision. No automatic larger confirmation, live game or publication.
