## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# v0.16 coordinated search-core experiment

**Approved for v0.16.** I chose promotion from the standard screen and
authorized a PR, merge and release. The earlier recommendation for a larger
confirmation is preserved below; I deferred that run. The release changes only
public identity and documentation relative to the screened playing code. I tested
one coordinated search-core candidate from tagged v0.15.3. The original development label, hashes and every result remain unchanged.
The active bot/config and live games are outside this release workflow.

The hypothesis is that one proven best root move, full-depth PVS, draw-safe
cache reuse and tactical-only sorting together buy useful search depth at equal
time. This is Stage A of the [architecture plan](../v0.15.3-architecture/NEXT.md).
I have not added specialized tactical generation, a board rewrite, evaluation
weights, the bishop-pair bonus or another selective-search feature.

## What changed

Completed searches return one best move rather than enumerating all exact ties.
Every tested score still matches the released engine and the selected move
belongs to its independently scored best set. PVS probes later moves with a
null window and re-searches improvements without reducing depth. A separate
1 MiB table provides legal board-only move hints; score reuse still respects
draw history. Only unreachable pre-pawn/capture history is omitted from the
cache fingerprint during child descent. Nonchecked quiescence sorts a copied
tactical subset; full legal generation remains unchanged.

[PLAN.md](PLAN.md) records the frozen scope and gates. [METHOD.md](METHOD.md)
describes the contracts, controls, corpus, timing schedule and limitations.
[candidate-source.json](candidate-source.json) preserves the source/tests as an
overlay on tagged v0.15.3, and [reproduce.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) restores it in a new
directory without matches or live play.

## Correctness and fixed-depth work

The full race-enabled repository tests, vet, native/ARM64 builds and offline
public-config lichess-bot wrapper pass. The wrapper checks 12 legal plies, mate
and a history-dependent drawing move. Existing exhaustive minimax/window,
mate-distance, promotion/EP/evasion/stalemate and cancellation controls remain.
New tests cover one-best membership across move orders, warm history in both
colors/seeding orders, irreversible descent/unwind, castling conservatism,
legal/colliding hints and preservation of the library's move cache.

Across all 84 frozen positions at depths 3 and 4, plus 12 positions at depths
5 and 6, **all 192 paired searches complete legally, with identical scores and
valid best-set membership**. No particular Stockfish move is required.

| Cold probe | Positions | Baseline nodes | Candidate nodes | Work saved |
| --- | ---: | ---: | ---: | ---: |
| Depth 4 | 84 | 1,185,507 | 1,010,100 | 14.8% |
| Depth 6 | 12 | 5,650,566 | 3,658,933 | 35.2% |

The single-run summed elapsed time is approximately 30% lower at depth four
and 46% lower at depth six. These are diagnostic timings on different cohorts,
not confidence intervals or playing-strength gains. Savings vary by position.
The warm GHI control retains the actual -902 cp without relevant history and
0 with the drawing history, including a table seeded from the losing context.

[checks.json](checks.json) preserves every validation command/exit/time;
[analysis.json](analysis.json) summarizes raw probes. [ERRATA.md](ERRATA.md)
records the initial test compilation and input-filename mistakes alongside
their original failure logs. [unchanged-policy.json](unchanged-policy.json)
checks that evaluator/material, SEE, clocks, killers, repetition counting,
mate normalization and reporting remain byte-identical to the tag.

## Equal-time performance and strength

All **2,046 native UCI requests** complete as expected: 2,018 nonterminal
legal completed searches and 28 expected terminal null moves. The same requested
budgets and unchanged clock margins apply to both actors.

| Native block | Paired cases/trials | Baseline mean depth | Candidate mean depth | Deeper / equal / shallower |
| --- | ---: | ---: | ---: | ---: |
| Serial, 100 ms | 252 | 4.357 | 4.687 | 80 / 172 / 0 |
| Serial, 1,000 ms | 84 | 5.976 | 6.405 | 33 / 51 / 0 |
| Three-search load, 100 ms | 504 | 4.310 | 4.619 | 150 / 354 / 0 |

Maximum observed responses are 96.5 ms serial at 100 ms, 996.7 ms at 1,000 ms,
and 116.7 ms under three-search load. The last includes shared-host scheduling
and collection overhead; all stay within the routine requested-time-plus-50-ms
allowance. No claim of a universal timing bound follows.

For the 54 nonterminal warm/cold continuation pairs per actor, the baseline
finishes deeper/equal/shallower warm in 2/50/2 cases, and candidate in 2/51/1.
These small warm gains do not establish a large persistent-cache effect.

The [standard screen](../../results/v0.16-search-core-dev-vs-v0.15.3-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
completed all three 50-game matchups concurrently, at 100 ms per move plus the
unchanged 50 ms protocol allowance. The practical previous-release gate passed.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | ---: | ---: | ---: |
| Tagged v0.15.3 | 27 / 11 / 12 | 0 | 0 | 65% |
| Stockfish Skill 2 | 21 / 1 / 28 | 0 | 0 | 43% |
| Pinned Sacrifice | 16 / 17 / 17 | 0 | 0 | 49% |

I independently replayed every move in all 150 PGNs using python-chess and
checked the saved results; all pass. [screen-review.json](screen-review.json)
retains every game, all losses, opening/color breakdowns, timings and illustrative
paired uncertainty. [screen-command.json](screen-command.json) records the exact
command, frozen-hash/gate checks, duration and exit status.

The **direct 65% screen is promising measured strength evidence**, alongside a
separately established functional/efficiency improvement. It is not an Elo
estimate or proof that the population score is 65%. Its illustrative pair-based
95% t interval is about 53.5–76.5%; this treats 25 pairs as independent even though
only ten public openings are repeated, so it can understate uncertainty.
The direct score is 50% as White and 80% as Black; that color concentration
needs checking across broader openings before a release decision.

Sacrifice's 49% is encouraging compared with the earlier released-policy control
at 27%, but that older sample was not a simultaneous external control. I cannot
assign the 22-point difference causally to this candidate. Stockfish's 43% is
similar to the old 44% control; this screen does not demonstrate a broad leap
against every opponent. Tie selection, finite sampling and shared-host timing
remain explanations to test. No component gets a separate strength attribution.
No larger confirmation has run.

## Limitations and decision

I have a positive strength screen, not a confirmed general strength gain.
Public/correlated diagnostics, repeated openings, a small deep cohort, finite
timing variance and changed tie selection limit inference. Faster search retains
the same leaf evaluator and known threat blindness. The hint table adds 1 MiB. A count fingerprint remains probabilistic;
the controls exercise specific draw histories rather than proving every case.

**I recommend confirmation, not rejection or more code tuning now.** Correctness,
work savings, native completed depth and the direct strength screen agree well
enough to justify a larger independent check. I would freeze this exact source,
use 200 games per opponent over the existing 50-position confirmation suite,
retain three-way concurrency and the same opponents/time controls, and review
all results before any release decision. This is a proposal, not an executed
or automatically authorized run. Details and stop conditions are in [NEXT.md](NEXT.md).

The next implementation priority remains specialized tactical generation and
reusable search storage **if confirmation supports this candidate**. The result
supports reducing wasted search work; it does not justify another broad research
reset. The unchanged evaluator/threat model remains the next major strength
workstream, and should move ahead of that backend project if confirmation or
focused loss diagnostics point there. Neither project starts automatically.
I will check in before committing, pushing or preparing a PR.

## v0.16 promotion

I explicitly approved publishing this screened playing code as v0.16 and
proceeding to a separately reassessed next candidate. The proposed 200-game
confirmation remains deferred, and the original recommendation/uncertainty
above is preserved. [release-source-identity.json](release-source-identity.json)
verifies that only the UCI name and its test change relative to the frozen Go
sources. [release-checks.json](release-checks.json) and
[release-provenance.json](release-provenance.json) preserve the repeated release
checks and separate six-game integration smoke. The development source overlay,
original screen, timing results and all failed attempts retain their identities.
