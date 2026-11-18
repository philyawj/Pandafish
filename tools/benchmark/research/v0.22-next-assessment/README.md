# Assessment after released v0.21

**Current recommendation: promote the [revised exact kernel](../v0.22-bitboard-exchange/RESULTS.md)
after review.** The clean standard screen is 64% / 71% / 70%, with no forfeits or
unfinished games; all 150 games replay legally. Fixed-work cold/warm savings are
35.35% / 37.84% total, with a separately measured 14.65% warm gain over the first
evasion version. The public 100ms surrogate still fails and is disclosed with
its explicit pre-game exception. No commit, PR, larger confirmation or live
game follows automatically. [The short project report](../../../../docs/RESEARCH_V022.md)
and [next ranking](../v0.22-bitboard-exchange/NEXT.md) separate those conclusions.

[PR #38](https://github.com/philyawj/Pandafish/pull/38) is merged and
[v0.21](https://github.com/philyawj/Pandafish/releases/tag/v0.21) is published at
`494604686a7afbda256da1625aa907baa6060084`. This branch starts from that tag.
Skill 3 remains the routine reference; no live games run. All research source,
including frozen/instrumented source, is Go.

The [initial plan](PLAN.md) reassesses earlier notes through realistic warm-game
profiles before implementation. Both modes complete the same 1,145 recorded roots
across 20 games and all ten openings. Position reconstruction accounts for about
0.53% of the measured pipeline. Child-state allocation dominates bytes, while
tagging, clock observation and runtime work dominate larger CPU portions.
Allocation dominance does not by itself predict elapsed savings. Instrumented
profiles/timings include setup/measurement overhead; actual paired tests decide
efficiency. Original raw profiles, identities, helper source and failures remain.

| Attempt | Observation | Decision under its stated objective |
| --- | --- | --- |
| [Per-ply child state](../v0.22-child-storage/README.md) | 3.3% cold median / 4.65% warm gain; warm bytes fall 79.3%. | Reject larger-gain implementation; preserve correct inactive code. |
| [Full-ray tag geometry](../v0.22-geometric-tags/README.md) | 1.1% median gain; deep subset 1.4% slower. | Reject; preserve caught initialization defect and corrected code. |
| [Occupied-ray tags](../v0.22-visible-tags/README.md) | 3.35% median gain; deep effectively flat. | Reject. |
| [Monotonic deadline, original scope](../v0.22-monotonic-deadline/README.md) | 6.95% cold median; deep/newer gain 7.2%/10.8%. | Fails original 15% larger-gain target; later separate small-scope assessment below. |
| [Positive-only history](../v0.22-quiet-history/README.md) | 3.4% cold median; deep/newer gain about 9–10%. | Reject; no native/game screen. |
| [Signed history](../v0.22-history-malus/README.md) | 5.6% fewer warm nodes; median elapsed 0.2% slower. | Reject, including after prospective workload correction. |
| [Aspiration](../v0.22-aspiration/README.md) | All 554 iterative scores/moves agree; 2.1% warm median gain. | Reject larger-gain attempt. |
| [Guarded LMR diagnosis](LMR_RESULT.md) | Non-overlapping eligible work: 27.4% deep public / 6.0% newer. | Fails preset opportunity; no playing reduction code. |
| [Strict/certified quiescence reuse](Q_REUSE_RESULT.md) | Potential node savings 5.2% / 7.3%; no observed bound contradiction. | Fails opportunity; no playing cache/proof code. |
| [Stand-pat relocation diagnosis](STANDPAT_RESULT.md) | Potential node savings 5.6%, SEE calls 38.4%. | Fails paired opportunity; no playing shortcut. |
| [Clock small-scope reassessment](../v0.22-clock-screen/README.md) | 6.69% timed warm gain; 46% / 47% / 71% screen. | Reject direct gate; original larger-target failure also stays failed. |
| [First staged picker](../v0.22-staged-picker/RESULTS.md) | 19.81% warm gain; 47% / 68% / 62% screen. | Reject direct gate; retain native failure and screen exception. |
| [Lazy within-stage validation](LAZY_LEGAL_DIAGNOSIS.md) | 25.26% unreached records versus 30% criterion. | No playing implementation. |
| [First evasion mask](../v0.22-evasion-mask/README.md) | 24.34% cold total; 9.11% warm extra versus 10% target. | Reject before native/games; original opportunity failure and scope review stay explicit. |
| [Evasion context deferral](EVASION_CONTEXT_REVIEW.md) | Entire helper only 0.54% sampled CPU. | Insufficient larger-gain headroom; no playing revision. |
| [Bitset exchange revision](../v0.22-bitboard-exchange/RESULTS.md) | 35.35% cold / 37.84% warm total; 14.65% extra over evasion; 64% / 71% / 70% screen. | Recommend review/promotion, with failed public surrogate disclosed. |

The [history workload correction](HISTORY_METHOD_CORRECTION.md) is frozen before
signed-history measurements; its failure remains a failure after using the more
appropriate iterative lifetime. Aspiration also uses the actual iterative entry
point. No new native or strength games are spent on those failed larger-work
attempts. Bound/identity/horizon safeguards are not relaxed to manufacture an
opportunity. All original thresholds and adverse diagnostics remain visible.

The [explicit later reassessment](REASSESSMENT_SMALL_SCOPE.md) chooses the exact
same short clock helper for a different, smaller implementation-cost objective.
It does not claim the original 15% target passed. The
[small-scope assessment](../v0.22-clock-screen/PLAN.md) requires previously
unexecuted timed warm pairs, all native gates and one routine concurrent screen.
Its three warm pairs pass the prospectively declared 5% gate with 6.69% median
gain and all scores/moves/nodes exact. The rebuilt binary is byte-identical to
the original clock artifact. All 1,252 native requests pass, but its one standard
screen scores 46% against v0.21, 47% against Skill 3 and 71% against Sacrifice.
All 150 games replay legally with no forfeits/unfinished games. The direct gate
fails; the clock code is archived again. The reference-cache argument failure
and reference-only repair remain explicit; no native actor query is rerun.
No precise Elo claim or automatic larger confirmation follows.

This study provides evidence against several attractive larger-gain assumptions,
not a general claim that engine development is complete. Main late-quiet work is
small at these measured depths, quiescence duplicates cover limited descendant
work, and the existing coordinated root/PVS/table pipeline leaves little for
simple aspiration/history layers. A future low-cost checking-move/quiet-threat
interface or different evaluation representation is a more substantive research
project; neither is implemented here. No book or online dataset is required for
the measured clock change. The failed game screen triggers another reassessment:
defer legal move generation until preferred/tactical stages fail to refute a node.
This previously suggested staged-picker route is different from changing sort
record size, cache policy or search breadth. Measure its opportunity and preserve
the original complete move order before any implementation claim.

The [first staged picker](../v0.22-staged-picker/README.md) achieves 15.76% cold
and 19.81% warm median savings with all fixed-work results exact. One public
100ms teacher-loss gate fails; the remaining suites pass. An explicit documented
exception authorizes one research screen after matched-depth diagnosis. It
scores 47% / 68% / 62%, with all games legal and no forfeits/unfinished games.
The direct gate fails, so that version is rejected. Its strong Skill 3 row and
real efficiency gain do not silently satisfy the promotion requirement.

The next [lazy-validation diagnosis](LAZY_LEGAL_DIAGNOSIS.md) finds only 25.26%
unreached legal records within materialized stages, below its 30% opportunity.
No playing lazy iterator follows. A distinct [evasion-mask diagnosis](EVASION_DIAGNOSIS.md)
also fails its original equal-count criterion (38.01% versus 50%), with zero
legal rejections and exact warm results. The
[scope review](EVASION_SCOPE_REVIEW.md) identifies that full two-king tags have
51.66% avoidable calls while cheaper existence checks have little opportunity.
It transparently authorizes one implementation/cost probe, preserving the
original failure and the unchanged 10% additional warm timing gate. The
[generation-kernel revision](../v0.22-evasion-mask/PLAN.md) retains the exact
staged base plus one geometric check-evasion prefilter. No old failed clock,
history, storage, aspiration or quiescence-reuse code is activated.

The first evasion revision passes its cold and correctness controls but misses
the unchanged warm extra-gain target: 9.11% versus 10%. It is
[archived before native/games](../v0.22-evasion-mask/README.md). A
[context-cost profile](EVASION_CONTEXT_REVIEW.md) puts its helper at only 0.54%
sampled CPU, too little larger-gain headroom for the suggested deferral tweak;
no such playing revision is implemented. The new
[exchange-kernel revision](../v0.22-bitboard-exchange/PLAN.md) instead replaces
repeated 64-square attacker scans with exact local bitsets. It retains the
verified generation foundation, keeps all prior failures explicit and requires
10% further warm saving over the already faster evasion version. All heuristic
semantics, square ties, legal recaptures and cancellation points must agree with
the frozen mailbox oracle before timings or playing evidence.
