# Transition-only result: promotion rejected

The specialization passes the predeclared correctness, performance and native
choice gates. It preserves 325 distinct position/depth tree comparisons and
reduces fixed-work elapsed time about 21%. The 37,795 transition controls,
race/vet/build checks and vendored upstream runtime replay pass.

The standard concurrent screen completed 50 games per opponent:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.19 | 19 / 7 / 24 | 0 | 0 | 45% |
| Stockfish Skill 2 | 29 / 3 / 18 | 0 | 0 | 61% |
| Pinned Sacrifice | 20 / 17 / 13 | 0 | 0 | 57% |

The previous-release criterion fails, so promotion is rejected. The runner's
exit 1 and every completed game remain preserved. The fixed-tree efficiency
gain is functional evidence, not established playing-strength improvement.
Small repeated-opening screens are noisy, and the external scores are not
simultaneous external baseline controls. No larger confirmation or favorable
rerun follows. The original candidate binary, complete source overlay, all
timings and native/reference records remain archived.

The next revision will follow a new profile of this kernel. Existence checks
need own-king legality but discard opponent-check tags. The original attack
predicate also repeatedly computes rays and piece lookups. A coherent search-
legality refinement can remove that unused/repeated work without changing
search values or eligibility. It must produce at least 10% additional fixed-
work time reduction versus this revision before another strength screen is
justified. No opening, clock, move preference or value will be tuned against
these games, and quiet-mate search remains absent.
