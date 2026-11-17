# My v0.19 research summary and next steps

I release v0.19 to preserve the completed evaluation investigation and return
all repository source to Go. **Playing policy remains v0.18.** I promote none
of the fitted evaluation models and make no new playing-strength or Elo claim.
The only engine change is its reported version. I selected v0.19 for this
substantial research/tooling increment.

## What I attempted, and when

Before releasing v0.18, I tried a compact fitted board evaluator and a faster
implementation of the same idea. Static fitting improved, but actual short-budget
decisions failed. I archived both in the v0.18 release, then promoted the separate
exact-score computation kernel. Its efficiency/strength evidence belongs to
[that report](../tools/benchmark/research/v0.18-exact-evaluation/README.md), not
to the rejected learned values.

After the v0.18 merge/release I made three further local attempts. They are
published together here for the first time. Their `v0.19-` directory names are
research labels, not claims that a v0.19 engine played those experiments.

| Attempt | What I changed | Result and decision |
| --- | --- | --- |
| [Decision-aligned fit](../tools/benchmark/research/v0.19-decision-evaluation/README.md) | Fit branch preferences at saved student leaves rather than static error alone; fresh rollout families and multiple acceptable moves. | Holdout surrogate regret 54.02 to 54.74, about 1.32% worse. I rejected it before activation. A stricter later audit also found early cross-fold transpositions. |
| [Independent-family retry](../tools/benchmark/research/v0.19-independent-families/README.md) | Block prior/reserved positions from ply four, audit all legal student witness paths before labels, and split whole transitive components. | Fresh surrogate regret improved 8.35%. Actual fresh finite-reference choices improved, but public 100ms loss worsened 55.27 to 68.04 and a fresh 1000ms choice acquired a losing-mate report. I rejected it and restored v0.18 source. |
| [Sharp-branch/both-turn retry](../tools/benchmark/research/v0.19-sharp-branches/README.md) | Retain complete mate-containing roots, separate mate ordering from finite-CP regression, expand negatives, balance turns and correct target scaling before fitting. | Validation improved, but fresh holdout penalty rose 44.04 to 49.99, **13.51% worse**; severe cases rose 3 to 6 and mate-risk roots 3 to 4. I rejected it before activation. |

I kept all failed runs, initial plumbing/test failures, raw queries, coefficients,
boundaries, source identities and diagnostics. I did not retune inspected
holdouts, rerun unfavorable games or silently weaken gates. No strength matches
ran for these three post-release fits; none reached the whole acceptance plan.
The second fit's finite-reference improvements are functional diagnostics, not
measured strength. Its mate flag describes survival under a finite reference
budget in an already bad position, not proof that the baseline could save it.

The third experiment's primary penalty includes an explicit capped mate-risk
convention. I do not compare its number directly with the earlier CP-only
surrogate losses. I report fully finite and per-turn results separately in its
archive. A reported finite score is not proof that a position is drawable.

## What I learned

My original narrow evaluation preference was not sufficient. Cleaner boundaries
and difficult labels were necessary repairs, but they did not make this bounded
12-feature, 24-coefficient model generalize safely. One fit passing frozen
student-leaf selection also failed actual search. The next learning experiment
needs a different representation/selection design, rather than more hand-tuning
of this same inspected sample.

The offline pair target originally clipped the teacher gap before subtracting
the baseline gap. That can ask for the wrong correction direction on large
already-correct margins. I froze the algebraic correction **before fitting the
third attempt**, tested it and preserved the old source. This was a research
objective error; released v0.18 never used these fitted weights. I cannot
attribute the third result to that fix alone because its data/design changed too.

Exact bounded-capacity probes show 15 training mate-pair preferences cannot be
expressed within the original coefficient box, even allowing integer rounding.
Four dominating bad branches span two distinct roots. These are frozen-witness
limits, not a proof that all classical or neural evaluation must fail.

A broader training-only scan provides a more concrete next mechanism:
**15 saved static records contain a legal quiet mate in one**, across 14 distinct
leaf FENs, 13 roots and 12 families. Original python-chess and unchanged tagged
Go independently verify them. Full-window v0.18 quiescence returns a finite
stand-pat score at every one. That establishes a local capture-only horizon
omission. It does not establish a forced root loss, a frequent real-game problem
or a future strength gain. Hypothetical opponent mates after a pass are only
threats and cannot justify a forced-loss score.

I have no new evidence that a prior cache improvement is blocking progress.
Search history/cache correctness remains covered by the earlier controls.
The successful v0.16–18 engineering changes and failed fitted values are
different outcomes. I have not established that I have found every weakness or
that an opening book or Lichess dataset is required. I can test the next concrete
mechanism using offline positions and legal proofs.

## My next ranked experiments

1. **Bounded exact quiet-mate proof at the quiescence boundary.** I would profile
   a nonrecursive legal checking-move/checkmate proof kernel. I would score only
   proven actual mates, preserve mate distance/history/deadlines and avoid
   arbitrary recursive quiet-check expansion. The previous naive extension timed
   out. I would require positive proofs, checks with legal escapes/captures/blocks,
   pinned moves, stalemate and special-move negatives, then cost and equal-time
   decision gates on fresh both-turn data. I would reject an unaffordable kernel
   even if it finds the saved mates. This is my top recommendation.
2. **Profile-driven reuse of legal/check/status information.** I would reprofile
   v0.19 before selecting one additive backend change. The local library boundary
   already delivered useful exact-tree efficiency gains. A new speed-only
   experiment must preserve scores, moves, node counts and rules; I would keep
   it separate from the changed-search mate experiment.
3. **Different evaluation representation plus actual-budget validation.** I
   would compare contextual threat features or a compact nonlinear model on
   representative training data, declare data/inference budgets, select with
   actual-budget validation and use a fresh untouched holdout. I would implement
   the new pipeline in Go. Larger caps or a neural label alone are not a plan.
4. **Selective search or make/unmake after a new measured case.** I would not
   repeat the failed broad pruning bundle or rewrite board state merely because
   the evaluation attempts failed. Profiling and focused correctness controls
   must justify either larger project.

I reassess after each experiment. Only a candidate that passes mechanism,
correctness, cost and actual-choice gates reaches the standard 50 games each
against the tagged baseline, Skill 2 and pinned Sacrifice, with all three
matchups concurrent. No larger confirmation runs automatically. I have not
implemented any of these next playing-code hypotheses in v0.19.

## My Go-only migration and preservation boundary

I removed 131 Python/shell source files, including archived `.py.txt` files,
and replaced the shell test helper with a Go subprocess. I added the shared
[Go research tools](../tools/research/README.md) for inventory verification,
archived Go-engine restoration, frozen integer-decision checks and mate proofs.
I explicitly retired old one-off collectors, NumPy trainers and external-oracle
drivers; I do not claim that replaying their outputs is an equivalent Go port.
All future repository programs/tests/research source must be Go, enforced by a
source-policy test and the short tracked [AGENTS.md](../AGENTS.md).

I preserve all 1,148 original research/document files in the release attachment
`v0.19-historical-research-evidence.tar.gz`, including scripts that had never been
committed. Its SHA-256 is
`ff77f6be062c46874387137c636741c81c3aaceda9eb256aef8a926dc4e403f5`.
The attachment contains its own per-file manifest. Before migration, all 1,023
existing research-inventory entries passed their original hashes. I keep those
inventories unchanged and explicitly map retired scripts/document rewrites in
[migration.json](../tools/research/migration.json). Raw results remain in Git;
original scripts also remain in older tags where previously published. The
release attachment preserves the original environments, not a current Python
dependency for Pandafish. No source language is hidden through GitHub attributes.

Markdown command examples and the separately installed optional lichess-bot
client retain their normal invocation instructions. I have not rewritten or
vendored that external client. Repository engine, developer tools, tests and
research programs are Go.

I link release validation, integration-smoke outcomes and remaining limits from
[the v0.19 validation record](../tools/research/v0.19-validation.md). Those checks
establish release/tooling integration and unchanged playing policy, not a new
strength gain.
