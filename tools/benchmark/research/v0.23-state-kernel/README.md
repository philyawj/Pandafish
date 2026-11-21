# Coordinated exact state kernel

The [frozen plan](PLAN.md) selects one state-handling project after the failed
selective-search/evaluation attempts. Integrate caller-owned per-ply child
Position/Board storage with ordinary-move tagging from virtual occupancy and
attacker sets. The separate prototypes remain failed under their original15%
size targets (9.137% and4.543% saving); neither is silently promoted. This combined
candidate prospectively targets12% after reassessing the cost of rejecting
moderate exact improvements before actual-budget evidence.

[Child storage](../../../../search_position.go) belongs to one worker and ply.
A child remains valid through recursive work and PVS re-search; descendants use
different cells. A later sibling invalidates the old scratch result. Parent/
destination aliasing panics before mutation, legal caches reset on reuse and
outer slice growth cannot move active cells. Immutable general APIs remain.
[Virtual tags](../../../../third_party/chess/search_virtual_tags.go) classify
ordinary moves without copying/updating a temporary Board. Promotions, en
passant and castling retain the original simulated path. No shared global pool,
make/unmake rewrite or unrelated library rules are introduced.

Material/evaluation, move eligibility/order, main/PVS/root/q policy, table
identity/depth/replacement, repetition/mate rules and clocks remain v0.22.
There are no archived reductions, extra quiet checks, static pruning, parallel
workers or fitted terms in this candidate. Equal-time choices may differ when
the same search completes more work.

Normal/race/vet pass in the root and nested chess module, as do test/engine
builds. Controls cover ownership, stale caches, aliases, worker isolation, parent
stability, history/cancellation, castling/EP/promotions, perft and140,778 tag
comparisons with both transition oracles. All596 frozen independent vectors
also pass, exercising actual scratch child FEN/hash and depth-two perft.

Three alternating warm pairs retain all1,145 roots from20 games/all10 openings,
persistent per-game tables and completed depth five. Every released score,
selected move and all43,118,858 nodes agree exactly; actor results repeat.
Median elapsed ratio is **0.858578**, or **14.142% less time**. The12% cost gate
passes. This establishes a moderate functional efficiency gain on the selected
workload, not an Elo gain or universal speed guarantee.

| Equal-time suite | Mean finite reference loss, v0.22 → candidate (cp) | Gate |
| --- | ---: | --- |
| public-100 | 45.8955 → 45.8955 | Pass |
| public-1000 | 54.8507 → 54.8507 | Pass |
| newer-100 | 30.6794 → 30.8321 | Pass |
| newer-1000 | 24.2424 → 24.1667 | Pass |
| fresh-100 | 42.3255 → 40.1837 | Pass |
| fresh-1000 | 30.6299 → 28.1969 | Pass |

All six unchanged quality gates pass, including both-turn requirements for
newer/fresh suites. No exception, new losing-mate flag or timing overrun occurs.
Fresh severe errors remain5→5 at100ms and2→2 at1000ms. Aggregate passes do not
mean every move improves: fresh short-budget White worsens49.1719→50.6719cp,
and fresh long-budget Black worsens20.5938→21.6406cp. All adverse rows remain.

Fresh seed2026112123/stream123 retains148 roots,74 per turn, all16 families.
Canonical exclusions cover every consumed public/newer/fresh cohort and saved
live survey boards, including the one-choice budget's142 roots. Collection uses
only the released student/pinned teacher, not candidate-selected samples.

The sole standard concurrent screen completes:50 each against taggedv0.22,
Stockfish Skill3 and pinned Sacrifice;100ms requests/50ms grace,400-ply cap,
10 openings, AppleM2/8CPUs/AC, one search thread per engine. No build/validation
work overlaps matches; external workload is unmeasured. The candidate still
identifies as v0.22 in UCI; the development benchmark label and binary/source
hashes distinguish it.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.22 | 20 / 11 / 19 | 0 | 0 | 51% |
| Stockfish Skill3 | 34 / 0 / 15 | 1 | 0 | 69.4% |
| Pinned Sacrifice | 34 / 8 / 8 | 0 | 0 | 76% |

The Stockfish unfinished record is game38 (zero-based index37), candidate Black,
Queen's Pawn opening,400-ply move cap. Its possible50-game score is68–70% when
the unfinished game counts as a loss through a win. It is not a timeout or
shutdown. All150 records replay legally and match reported results/lengths;
149 games finish, zero forfeits. Candidate response maxima are98.4/99.3/97.7ms
across the three rows; zero candidate attempts exceed100ms. Stockfish/Sacrifice
have1669/2876 attempts over the request but none beyond the150ms allowance.
Opponent binary/options/asset hashes and openings match the pinned profile.

**Recommend promotion as a moderate exact efficiency improvement.** The practical
direct gate passes, all six choice gates pass without exceptions, and external
rows avoid the earlier selective candidate's large adverse signal. However the
direct lead is only one game, and repeated openings, many investigated candidates,
timing variability and one capped game limit strength inference. The earlier
v0.22 screen's71% Stockfish/70% Sacrifice rows are separate small samples, not
paired proof that69.4%/76% improves or preserves true strength. No meaningful
Elo or large playing-strength jump is established.

A broader-opening confirmation is justified before a strength-gain claim or
target increase; it is not automatically authorized and has not run. The14.142%
same-work saving is independently demonstrated by the warm protocol. Stop at
the requested checkpoint before commit/push/PR/release. No further engine
experiment, unchanged rerun or live session starts as part of this recommendation.

See the [runner report](screen/README.md), [legal replay](games-audit.json),
[settings, identities and raw records](screen/v0.23-state-kernel-dev-vs-v0.22-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba/summary.json),
and [complete research summary](../../../../docs/RESEARCH_V023.md).

The118-file replayable source/module snapshot is `candidate-source.json`;
all commands, gates, raw data and failures remain. Compiled binaries are ignored.
