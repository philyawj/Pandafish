# Research after released v0.21

**v0.22 promotes the [revised exact kernel](../tools/benchmark/research/v0.22-bitboard-exchange/RESULTS.md)
after the approved review checkpoint.** Its standard concurrent screen scores 64% against tagged v0.21,
71% against Stockfish Skill 3 and 70% against pinned Sacrifice, with all 150 games
legal, finished and free of forfeits. The original
[candidate checkpoint](../tools/research/v0.22-checkpoint.md) retains its
pre-approval wording; [release validation](../tools/research/v0.22-validation.md)
records the subsequent authorized promotion. No next candidate is started.

Fixed-work search time falls 35.35% in the cold median and 37.84% in a new warm
tagged-baseline pair, with identical moves, scores and nodes. The new exchange
kernel adds 14.65% warm savings over the first evasion prototype. These functional
results are distinct from the small game screen's measured strength evidence.
The public 100ms teacher-loss gate still fails; the other three suites pass.
The [documented pre-game exception](../tools/benchmark/research/v0.22-bitboard-exchange/SCREEN_DECISION.md)
keeps that failure explicit. No root is forced to a teacher move.

Released [v0.21](https://github.com/philyawj/Pandafish/releases/tag/v0.21) and
[PR #38](https://github.com/philyawj/Pandafish/pull/38) complete the preceding
precomputed-ordering promotion. The next investigation starts from that tagged
source. It reassesses candidates after each result, preserves unsuccessful
attempts and stops before committing a viable candidate for review.

The [assessment and complete attempt table](../tools/benchmark/research/v0.22-next-assessment/README.md)
link plans, diagnoses, source, measurements, adverse controls and decisions.
Allocation volume alone did not predict elapsed benefit: per-ply child storage
removed roughly 79% of warm allocated bytes but saved only about 5% of time.
The geometric-tag variants, positive/signed history and aspiration also failed
their declared larger-gain targets. Guarded reductions, certified quiescence
reuse and stand-pat relocation failed opportunity measurements before playing
implementations. This is measured evidence against these particular mechanisms
and workloads, not evidence that broader search or evaluation work is exhausted.

A separately documented small-scope clock reassessment saved 6.69% of timed warm
search and passed all native suites. Its one standard screen scored 46% against
v0.21, 47% against Skill 3 and 71% against Sacrifice. The direct gate failed.
Its code and all games remain archived; the clock helper is inactive again.
The original larger-gain failure and the cache-argument/reporting failure remain
visible. No favorable rerun or larger confirmation followed.

The subsequent [staged move-generation candidate](../tools/benchmark/research/v0.22-staged-picker/README.md)
targets work that is never used after an early cutoff. It preserves the eager
move order and all fixed-work moves, scores and nodes while saving 15.76% cold
and 19.81% warm median search time. Its original public 100ms teacher-loss gate
fails narrowly; the other three suites pass. Matched-depth inspection shows the
largest regression is a pre-existing deeper-search choice, not changed ordering.
The [explicit research-screen exception](../tools/benchmark/research/v0.22-staged-picker/SCREEN_DECISION.md)
keeps that failed gate false and documents the decision before one standard
concurrent game screen. Efficiency, correctness, surrogate decisions and measured
game outcomes must remain separate conclusions.

That first staged screen subsequently finishes at 47% direct, 68% against
Skill 3 and 62% against Sacrifice, with no unfinished games or forfeits. The
direct gate fails; that binary is archived without promotion. The investigation
continues with measured diagnoses, not a rerun of the same candidate. Lazy
within-stage validation fails its opportunity gate. An exact check-evasion mask
then reaches a separately documented local cost prototype; its original
equal-count opportunity failure and revised expensive-tag rationale remain
explicit. The prototype must add 10% warm savings before further native/games.

The first evasion version subsequently misses its extra warm target, 9.11%
versus 10%, and reaches no native/game screen. Profiling finds too little context-
helper cost to justify the proposed deferral tweak. The next revision replaces
the exchange evaluator's repeated square scans with exact bitsets, preserving
legal recaptures, square ties and every original score. It adds a substantive
14.65% over the faster evasion foundation before earning a separate native/game
decision. All earlier failures stay intact; this is not an unchanged-code retry.

The [full result](../tools/benchmark/research/v0.22-bitboard-exchange/README.md)
links hypotheses, implementation, controls, source/binary identities, timings,
all native reports, the failed gate, PGNs and independent reviews. The
[next-step ranking](../tools/benchmark/research/v0.22-bitboard-exchange/NEXT.md)
prioritizes a genuinely cheaper checking-move/tactical-horizon design, then a
different contextual evaluation project, with fresh diagnosis before either.
The old quiet proof, fitted values and small isolated bonuses are not reactivated.
A broader-opening confirmation is justified before a large general strength
claim or advancing the target; it has not run and is not automatic.

The current candidate combines exact staged generation, evasion prefiltering
and the new exchange computation. Material, search eligibility/order, horizons,
score-cache/history policy and clocks remain v0.21. There is no measured evidence
here that a prior cache improvement blocks progress; this does not prove every
fundamental opportunity has been identified. Skill 3 remains the standard.
All repository engine, tests, tools and archived research source remain Go.
No live game, opening book or online dataset is part of this investigation.
The evidence does not support weakening history/cache identity or forcing saved
root moves to match a teacher.
