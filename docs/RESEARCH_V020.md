# Next-candidate investigation after v0.19

**Recommendation: promote the search-legality kernel after review.**
The strongest supported direction is **reducing repeated legal-move work while
preserving the search tree**. The current candidate specializes immutable
transitions, reuses computed attack rays and skips check tags that an existence
query discards. Selected fixed-work elapsed time falls about 39% versus tagged
v0.19, with exact scores, moves and nodes. Its standard concurrent screen scores
65% against tagged v0.19, 67% against Skill 2 and 58% against pinned Sacrifice,
with zero forfeits or unfinished games in all 150 games. Strength and limits
are assessed separately in the [candidate results](../tools/benchmark/research/v0.20-search-legality/RESULTS.md).
The candidate checkpoint preceded publication. Promotion to v0.20 was then
approved, with release checks recorded [separately](../tools/research/v0.20-validation.md).
Archived checkpoint notes retain their original pre-approval wording. The later
tagged Skill 2–10 benchmark is separate from candidate selection and does not
automatically change the default reference level.

## Fresh assessment rather than a fixed checklist

The [v0.19 research](RESEARCH_V019.md) ranked a quiet-mate proof ahead of backend
reuse and a different evaluation model. That ranking was reasonable: 15 saved
static records contained legally verified quiet mates missed by capture-only
quiescence. Previous compact fitted evaluators repeatedly failed generalization
or actual-budget checks; another coefficient adjustment had less support.
Earlier v0.16–18 exact-computation changes had produced practical gains, but
their success did not establish that every backend optimization would help.

The older [architecture ranking](../tools/benchmark/research/v0.15.3-architecture/README.md#recommendation-and-what-ships)
put coordinated search, tactical generation and calibrated evaluation first.
The first two became v0.16/v0.17; v0.18 separately accelerated the original
evaluation after fitted values failed. The
[post-v0.18 note](../tools/benchmark/research/v0.18-exact-evaluation/NEXT.md)
then prioritized decision-aligned learning, which the v0.19 attempts tested and
rejected. The [backend note](CHESS_SEARCH_BACKEND.md)
also identified repeated tags and shared legal/status information as possible
future uses, conditional on profiling. The current measured attack/tag kernel
therefore has stronger support than either mechanically continuing the old top
three or launching an unprofiled board rewrite.

Three isolated steps followed. Each had a written plan before measurements;
failed outcomes changed the next choice.

| Attempt | Evidence | Decision |
| --- | --- | --- |
| [Nonrecursive quiet-mate proof](../tools/benchmark/research/v0.20-next-assessment/QUIET_RESULT.md) | Finds all 15 saved proofs, with about 8.3% additional fixed-depth cost. Public controls pass, but newer Black-to-move 100ms loss rises 38.06 to 44.98, beyond the frozen tolerance. An EP quiet/tactical classification defect is separately corrected. | Rejected before strength games; v1 and corrected v2 source are archived. No favorable retry or fresh 1000ms expansion. |
| [Affected-piece transitions](../tools/benchmark/research/v0.20-search-transitions/RESULT.md) | Exact trees and about 21.1% less fixed-work elapsed time; all native gates pass. Concurrent 50-game screens score 45% versus v0.19, 61% versus Skill 2 and 57% versus Sacrifice. | Rejected for promotion because the direct previous-release criterion fails. All 150 games and runner failure retained. |
| [Legality-kernel revision](../tools/benchmark/research/v0.20-search-legality/README.md) | A new profile identifies repeated attack rays, generic piece selection and unused opponent-check tags. The revision produces 38.7% total elapsed reduction and another 22.3% versus the transition revision, above its separately frozen extra-gain gate. Native gates pass; concurrent screens score 65% / 67% / 58%. | Recommend promotion after review; no commit or larger confirmation is automatic. |

The transition-only failure is not erased by reusing its functionally verified
code. The later revision changes a measured hot path and exceeds a new cost
threshold before another screen. It is not an unchanged-code rerun for a more
favorable outcome. Selecting after failed attempts still introduces selection
uncertainty; the final screen cannot be treated as an isolated significance test.

## Why this implementation

The initial v0.19 profile places 42.08% cumulative CPU in move tagging and 23.08%
in board update, including 11.31% directly in repeated bitboard setters. These
overlap and must not be added. A normal move scans all twelve piece sets through
the general update. Search instead needs only the moving, captured, promoted,
en-passant and castling-rook sets. Updating those directly removes repeated work
while retaining immutable child boards and original state/hash helpers.

The transition-only profile then places 22.19% cumulative CPU in the original
attack predicate and 11.93% in generic piece construction. The predicate repeats
ray calculations and piece selection; a search-only equivalent computes each
ray once and addresses each side's piece sets directly. A legal-move-existence
query also discards the move, so it needs own-king legality without classifying
whether the opponent is checked. Full move lists still retain every check tag.

This is a compact extension of the authored v0.17 library boundary. It avoids
the larger correctness burden of reversible board mutation. Original upstream
runtime, metadata and license remain unchanged; new helpers have to remain
synchronized with any future upstream upgrade. An independent Go mailbox attack
oracle, original transition oracles and exact-tree controls cover this cost.
There is no board attack cache, stale-tag assumption or shared mutable state.
The [implementation note](../tools/benchmark/research/v0.20-search-legality/IMPLEMENTATION.md)
describes the boundary in detail.

## What the evidence establishes

Correctness coverage includes 37,795 legal child transitions and 128,000 attack
predicates, existing perft/special-move/cancellation/history controls and full
root race/vet/build checks. All 325 distinct fixed-depth position/depth
comparisons retain nodes, score and selected move. Three alternating depth-four
pairs give a median elapsed ratio of 0.61256 versus v0.19. Depth-three,
depth-six and newer development roots also improve. These are selected-workload
efficiency results, not a universal throughput or Elo guarantee.

Four native equal-time suites pass the unchanged finite-reference, mate-risk,
severity and response gates. Newer positions also pass separately for each root
turn. Close alternatives to the reference move are accepted, and production
search is never root-restricted. The 145 newer roots were frozen before the
quiet-proof probe; their reuse makes them inspected regression controls, not a
blind holdout. Their 16 rollout families and timed actor trials are correlated.
Finite-budget teacher scores do not prove a position was savable, and absence
of a new mate flag does not prove tactical safety.

The game screen is the separate practical strength check, using tagged v0.19,
Skill 2 and pinned Sacrifice concurrently under the standard settings. Absolute
external scores across different candidate runs are not controlled external
baseline comparisons. Repeated openings, scheduling and modest samples limit
precision. No automatic larger confirmation or live game is authorized.

## Next priorities

1. **Review this kernel before another feature.** The 65% direct score and both
   external screens support promotion under the established practical gate.
   Complete legal replay, source identity and evidence audits before committing.
   A broader-opening confirmation is justified before claiming a large/general
   strength gain or advancing the target, especially after candidate selection.
   The optional 100 games per opponent over 50 openings is suitable for that
   question; it has not run and is not a new mandatory promotion gate.
2. **Profile the resulting engine before selecting further backend work.**
   Reassess remaining tag/status, copying/allocation and exchange-estimate cost.
   Prefer one exact-tree specialization with a meaningful cost target. Make/unmake
   becomes justified only if immutable child ownership/copying is a demonstrated
   dominant remaining cost and the projected gain exceeds its risk. Prior
   percentages do not identify the new bottleneck after this change.
3. **Reconsider tactical horizon coverage only with a new design.** The legal
   quiet-mate omission is real, but the standalone proof failed actual-budget
   balance. A later experiment needs a credible lower-cost or selective trigger,
   proof/negative controls and new both-turn actual-budget evidence. Do not simply
   reactivate corrected v2 or broaden quiet checks because saved mates exist.
4. **A different evaluation representation remains a larger research project.**
   Use contextual threats or a compact nonlinear representation with Go-only
   collection/inference, a declared inference budget, actual-budget validation
   and untouched final families. Repeating the inspected 24-coefficient model
   or raising caps has weaker support. No opening book or online dataset is
   required for the current mechanism; their potential future usefulness has
   not been ruled out.

No evidence here points to a prior cache change blocking progress. Nor does it
prove that all fundamental opportunities have been identified. The useful lesson
is narrower: directly verified tactical omissions can still lose at equal time,
and genuine speedups can still fail a small game screen. Both functional and
playing evidence are needed. The next ranking follows measured results rather
than exhaustion of a predetermined list.

All source, failed controls, exact identities, native/reference observations and
games are retained in the three linked archives. The
[testing guide](../tools/benchmark/research/v0.20-search-legality/TESTING.md) and
[Go tooling guide](../tools/research/README.md) describe reproduction and limits.
No Python or shell source, active learned model, live session or release/version
change is introduced by this candidate.
