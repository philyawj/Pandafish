# Conservative dead-material search: reject after the standard screen

Baseline is tagged v0.23; cleanup PR #42 is already merged without a release.
The preceding compact search is rejected after its46% direct screen. Its real
8.162% fixed-work saving does not transfer to this separate rule correction.

## Mechanism and correctness

The [frozen plan](PLAN.md) and [implementation boundary](IMPLEMENTATION.md)
cover only bare kings, one minor in total, and bishops confined to one square
colour without knights. An allocation-free additive predicate returns zero in
main search/quiescence before cache or stand-pat, after mate/stalemate checks.
Direct evaluation agrees. Other evaluation, search/q eligibility, ordering,
TT/history, clocks and legal-move behavior remain v0.23. Root search returns a
legal move on already drawn boards. No complete dead-position solver, tablebase
or optional move-count claim policy is added.

All eight both-module normal/race/vet/build stages pass.11,969 material/square
comparisons agree with the existing map-based game material implementation;
these are independent representations, not an external rules oracle.596 saved
independent move/state vectors still pass. Positive/live controls, promoted
bishops, two-knight/opposite-bishop negatives, captures, promotions, EP, castling,
cache bypass, history restoration, cancellation and zero allocations pass.

Three old tests expected a lone knight's material/placement to count after its
position was already drawn. Original source/failures remain; live placement
fixtures now use equal pawns, and the queen-fork continuation expects KN/K to
draw. Two early policy tables had invalid opposite-turn controls; all original
source/rows remain marked invalid. Final25 vectors explicitly check both-turn
legality and preserve mate precedence. Move-count mismatch is documented but
left separate because ordinary50-move claims require a policy/action contract.

## Cost and changed-value review

Three alternating warm pairs pass the prospectively frozen5% overhead ceiling:
median ratio**1.006937**, or**0.694% overhead**, with1,145 complete depth-five
roots. Baseline43,118,858 nodes; candidate43,127,210. This is a correctness
change with a cost limit, not a speed claim.

[Witness review](WITNESS_PLAN.md) retains all51 changed rows, including three
score and two move changes. Every changed root has a directly observed dead
search descendant. Shadow replays retain each actor's original moves/scores/
nodes exactly.1,979 unique dead witness FENs agree with game adjudication.
One live opposite-bishop policy root changes6→5cp through a legal capture into
KB/K; the other16 live/count/mate controls stay exact. Static formulas are
unchanged. The original overbroad assertion that all live roots must keep their
search value, and its failure, are retained. A live ancestor may correctly
change value when a formerly mis-scored child becomes a draw.

## Existing actual-budget controls

| Suite | Mean finite reference loss, baseline → candidate (cp) | Gate |
| --- | ---: | --- |
| Public100ms | 45.896 → 45.896 | Pass |
| Public1000ms | 54.851 → 54.851 | Pass |
| Newer100ms | 30.908 → 30.832 | Pass |
| Newer1000ms | 24.879 → 24.947 | Pass |

No new losing-mate flags or overruns. The newer100ms white-turn mean worsens
slightly and its worst saved regression is5cp; newer1000ms worsens overall with
a9cp worst row. These stay inside unchanged tolerances. Surrogate passes do
not establish strength.

## Frozen fresh controls and pending decision

Seed2026112202/stream125 retains144 roots,72 per turn/all16 families. Explicit
exclusions cover historical/public/newer/prior-fresh, warm/oracle, saved live
analysis, previous compact-screen games and draw policy controls. Temporary Go
adapters and hashes avoid tracking duplicate trajectories. This cohort is
consumed regardless of the next result. [Collection review](COLLECTION_REVIEW.md)
preserves the wrapper's self-output identity failure and verifies all523 actual
original inputs unchanged; no query is repeated or cohort replaced.

Both fresh gates pass: mean finite loss is 37.136 → 37.098cp at 100ms and
35.976 → 35.976cp at 1000ms, with no new losing-mate flags or overruns.

## Sole standard screen and decision

The unchanged concurrent screen runs 50 games each, 100ms requests/50ms grace,
400-ply cap and the existing ten openings on the recorded Apple M2/AC host.
No deliberate heavy work overlaps. All 150 games finish and replay legally.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.23 | 22 / 6 / 22 | 0 | 0 | 50% |
| Stockfish Skill 3 | 28 / 3 / 19 | 0 | 0 | 59% |
| Pinned Sacrifice | 35 / 7 / 8 | 0 | 0 | 77% |

Candidate timing maxima are 97.5/98.3/99.1ms, with zero candidate attempts above
100ms. [The complete report](screen/README.md), [legal replay](games-audit.json),
raw summaries/PGNs and console remain preserved. All 87 frozen source files and
binary identities match before and after play. Candidate SHA:
`daf742221cc09e5124092989cdcf45b8893fc67955eb17e28eb59eab3b230440`.

**Reject for promotion:** 50% does not exceed the frozen 50% direct gate.
No unchanged retry or larger confirmation occurs. The material-rule correction
is functionally established, but no measured general strength gain follows.
External scores do not override the failed direct gate; separate historical
screens are not simultaneous controls or precise Elo estimates.

Restore v0.23 playing source before another mechanism. The next reassessment
returns to the higher-ranked contextual evaluation opportunity, beginning with
its bounded capacity/inference diagnosis rather than another coefficient fit.
All research evidence remains local and uncommitted; cleanup PR #42 remains
merged without a release. No live play or research publication is authorized.
