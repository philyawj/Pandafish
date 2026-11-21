# Material-guarded reductions: substantial cost gain, reject quality gates

The general phase>=8 guard restores the original independent sparse-mate control.
Normal/race/vet/build checks pass. Three alternating 1,145-root nominal-depth-five
pairs have deterministic node ratio 0.500029 and median elapsed ratio 0.545518:
**about 50% fewer nodes and 45.448% less search time**. This is selective work,
not an equivalent full-width-depth speedup. Material/evaluation/quiescence,
cache identity/depth/replacement and clocks remain v0.22; no parallel code is active.

Original-tag full-depth review of all selected moves passes: 1,103 finite roots,
mean positive regret 1.014cp, maximum 273cp, zero mate-risk flags and zero value
contradictions. This supports the bounded mechanism but does not establish strength.

| Actual-budget suite | Baseline → candidate finite mean loss | Decision |
| --- | ---: | --- |
| Public 100ms | 45.8955 → 46.9403cp | Pass |
| Public 1000ms | 54.8507 → 52.1642cp | Pass |
| Newer 100ms | 30.6794 → 32.5573cp | Pass |
| Newer 1000ms | 24.2803 → 23.0455cp | Pass |
| Fresh 100ms | 50.3307 → 48.4788cp | **Fail**, one new reference mate flag |
| Fresh 1000ms | 37.7460 → 37.5000cp | **Fail**, White 33.0476→38.5079cp exceeds 38.0476cp limit |

No suite has timing overruns, illegal moves or query errors. Fresh controls use
seed 2026112119/stream119, 141 roots (71 White/70 Black) from all 16 families,
excluding prior controls and live survey boards. The failed fields stay false.

The [explicit scope reassessment](SCREEN_DECISION.md) preserves a confirmed
short-budget mate concern and the original no-exception failure. At one million
reference nodes, candidate Rfd1 reports mate -5 while baseline Rh1 remains finite
but badly losing. At matched depths both engines choose Rh1 at five and Rfd1 at
six. Thus that observed move is not uniquely attributable to reductions; its
equal-time disadvantage remains. The separate decision permits only the
previously unexecuted long-budget suite and requires it to pass before games.
It fails its per-turn gate, so **no strength screen runs**. No further exception,
favorable rerun or larger confirmation follows. There is no measured strength gain.

All source, raw rows/queries, source identities, original/corrected validation,
failed/stopped workflows, confirmation and matched-depth diagnostics remain.
The first validation's misplaced killer-audit reference change is corrected;
the unchanged independent mate assertion passes. Warm-plan preflight source/status
checks address the preceding prototype's sequencing error. Restore released code
before a separately frozen, more conservative frontier-protection revision.
