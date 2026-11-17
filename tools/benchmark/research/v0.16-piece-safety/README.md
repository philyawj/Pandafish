# v0.16 piece-safety experiment

**Recommendation: reject this candidate.** Preserve the experiment as research;
keep v0.15.1 playing behavior. The compact exposure term passes its functional
controls but does not explain the motivating quiet leaves, costs search depth
at equal time, and produces mixed, often worse, saved-position decisions. This
is not a measured engine-strength improvement.

No routine match screen was started: the prerequisite mechanism/cost assessment
was not promising. Consequently there are no new results against tagged
v0.15.1, Stockfish Skill 2 or pinned Sacrifice, and no previous-release gate
claim. No larger confirmation is justified for this source. No live Lichess client
was started. I subsequently authorized publishing these findings
as the v0.15.2 research release, with released playing behavior unchanged.

## Scope and frozen inputs

Start: tag v0.15.1, commit
`42f4510` (full revision in [candidate-source.json](candidate-source.json)).
The previous executable was rebuilt from a `git archive v0.15.1` scratch copy.
The root executable and private bot configuration were untouched. Only an
isolated ignored candidate executable was built.

The hypothesis was that quiet leaves undervalue the burden of defending pieces
attacked by pawns or left without a usable defender. Before implementing or
choosing weights, [cases.json](cases.json) froze the 30 public live confirmations
by sorted game URL, alternating 15 development and 15 public regression
holdout cases. Both headline cases land in development. These cases were
already published; the holdout is not blind. No weights were tuned after
looking at either split. [probe-inputs.json](probe-inputs.json) adds all ten older
v0.10 confirmations, retaining their histories. Selected errors are not a
representative playing sample.

Material values, placement formulas, search, quiescence, move ordering,
transposition-table policy, repetition and time management were unchanged.
The candidate changed only the static score and diagnostic UCI identity, with
focused tests. No other ranked evaluation recommendation was implemented.

## Diagnose before scoring

[trace-inputs.json](trace-inputs.json) holds both full histories and the saved
500,000-node referee continuations. The baseline diagnostic re-searches exact
best replies after each forced diagnostic root alternative, following existing
capture-quiescence eligibility through its first maximizing quiet leaf. It
records geometric attacker/defender origins and a legal attacker check for
every threatened target. Re-searches with exact windows and deterministic
first ties reconstruct one legal maximizing line, not a warm live PV or all
equivalent continuations. Production search code is not instrumented or changed.

Authoritative output:
[baseline-traces-history-corrected.json](baseline-traces-history-corrected.json).
The initial all-history fingerprint trace is preserved separately, as explained
in [early-failures.md](early-failures.md).

- **16G1f9rx, Nd4, depth 7:** reconstructed line
  `Nd4 Nxd4 exd4 Bg5 O-O Nd2 Be5 Bxf6 Qxf6`, score -12 cp for Black.
  Its settled leaf has no clearly undefended or pawn-attacked non-pawn target.
  The referee instead uses `Nd4 Nxd4 exd4 e5 Bxe5 Re1`, an intervening pawn
  push and a quiet rook move. The bishop on e5 becomes genuinely loose after
  Re1; that fact is real, but quiescence already handles Rxe5 once it can
  capture. An exposure penalty on that later board does not establish that
  the earlier search will choose the missing forcing continuation.
- **HhPI3vCx, Qd6, depth 6:** reconstructed line
  `Qd6 Rae1 Qc6 Bxf6 Kxf6 Ne4+ Ke7`, score +161 cp for Black.
  The knight on e4 is defended by Re1; there is no qualifying loose target
  at the settled leaf. The referee's `Qd6 Rad1 Qc6 Rfe1 ...` is different
  counterplay, involving queen harassment and Nd5. The apparent attack of
  Nf6 on White's undefended Qh5 is illegal because the knight is pinned to
  Ke7 by Bg5. A geometric-only penalty would misdiagnose the queen.

The numerical feature overlay on every traced board is in
[leaf-safety.json](leaf-safety.json). Neither headline baseline played-line
quiet leaf receives a safety correction. Detecting a later loose bishop or
an immediate rook attack on a queen is not a demonstrated explanation for
the original choice. Qd6's queen is defended when Rd1 attacks it; cheaper
non-pawn attackers on defended targets are outside this first narrow term.

## Candidate and controls

[candidate-source.json](candidate-source.json) is a complete overlay on the
tag, including tests and diagnostic identity. Experimental source remains data
outside the normal build after rejection.

The evaluator constructs occupancy and per-side origin masks for pawn, knight,
king and blocked sliding attacks. Absolute pins nominate suspect defenders.
Only pawn-attacked targets or targets lacking an unpinned non-king geometric
defender proceed to one-capture verification. A capture must leave the
attacking king safe. Where needed, a legal recapture check on the changed
board handles absolute pins, king recapture safety and newly opened rays.
It reuses the released exchange-board legality helper, without running a full
SEE sequence or changing the quiescence filter.

Original, untuned weights in cp:

| Target | No legal recapture | Legal enemy pawn attack |
| --- | ---: | ---: |
| Pawn | 8 | 8 if undefended, otherwise 0 |
| Knight / bishop | 16 | 24 |
| Rook | 24 | 36 |
| Queen | 32 | 48 |

Use the maximum applicable penalty per target, never the sum; kings are not
penalized. This adds exposure to the existing material/placement score. It does
not deduct a piece's material value a second time. Pawn attacks still receive
a small penalty when recapture is possible because the pawn is cheaper.

[controls-first.txt](controls-first.txt) preserves initial fixture failures;
[controls-final.txt](controls-final.txt) and [focused-check.txt](focused-check.txt)
record corrected and final controls. Each of 12 attack/defence boards is checked
as given, reflected horizontally, and reflected vertically with colors swapped.
Controls cover undefended targets, safe and defended targets, pawn attacks on
three piece classes, blocked rays, pinned attackers/defenders, legal and illegal
king recaptures, and newly opened defender rays. Every origin/target attack bit
is compared with the existing geometric helper on three varied boards.
Quiescence must still resolve an ordinary loose-piece capture; unrestricted
root searches must preserve immediate mate and the smothered-mate queen-offer
control `Qg8+ Rxg8 Nf7#`. No saved root move is prescribed as the only answer.

Limitations are deliberate: unpinned geometric defenders allow a fast skip and
may be unusable because of relative pins, overloads or unrelated checks. This
conservatively misses exposure. The term does not model forks, quiet future
attacks, exchange sequences, defensive tempo, mobility, or defended queens
attacked by cheaper non-pawns. It can penalize temporary exposure that supports
a stronger forcing plan, and it can break otherwise acceptable root ties in
the wrong direction. These limits matter to the observed failures.

## Measurements

All probes ran locally on Apple M2 / darwin arm64 with the recorded Go toolchain.
No CPU frequency, temperature or utilization telemetry was collected; ordinary
background system work was not controlled. No other deliberately launched
engine task overlapped the final balanced performance probes.

Full-root fixed-depth searches at depths 3 and 5 save **every tied best move**,
not a random representative. Timed searches use the normal 100 ms allowance,
95 ms internal hard budget, cold table, full reversible history, and preserve
random tie selection. Both use released search functions. Every one of the
initial 400 probes and final 480 timed probes completed with legal choices.
The deeper diagnostic has separate watchdog failures below.

Raw initial records: [baseline-probes.json](baseline-probes.json),
[candidate-probes.json](candidate-probes.json). The first baseline run briefly
overlapped a static microbenchmark. Its fixed scores remain useful, but its
elapsed times are not final cost evidence. Nothing was overwritten.

Clean equal-time schedule: baseline, candidate, candidate, baseline; each round
covers 40 positions with three trials, for 240 probes per engine. Records:
[baseline round 1](clean-1-baseline.json), [candidate round 2](clean-2-candidate.json),
[candidate round 3](clean-3-candidate.json), [baseline round 4](clean-4-baseline.json).

| Clean 100 ms probes | Tagged baseline | Candidate |
| --- | ---: | ---: |
| Total nodes | 3,673,443 | 3,006,622 |
| Mean completed depth | 3.725 | 3.475 |
| Mean elapsed ms | 92.960 | 93.515 |
| Maximum elapsed ms | 95.053 | 95.079 |

Candidate node count falls **18.2%**. Position/trial-aligned completed depth is
higher in two, equal in 176 and lower in 62 probes. Node count is not an exact
throughput ablation because scores change the tree, but the static
microbenchmark independently confirms the additional evaluator cost:

| Static score, median ns/op, three runs | Baseline | Candidate |
| --- | ---: | ---: |
| Start | 585.4 | 1,685 |
| HhPI3vCx middlegame | 489.0 | 1,616 |
| Rook ending | 325.9 | 1,213 |

[Baseline microbenchmark](baseline-cost.txt), [candidate microbenchmark](candidate-cost.txt).
These are selected-board timings, not universal speed or deadline guarantees.

## Decision quality, including improvements

The saved referee is a comparator, not a unique move oracle. Fresh Stockfish 19
queries use Skill 20, limiting off, one thread, 64 MiB, equal 100,000-node limits,
fresh hash and full game history for unrestricted search and each selected move.
This is diagnostic referee configuration, not the Skill 2 match opponent.
Every tied fixed-depth best and timed choice from either engine is assessed;
there is no requirement to choose the old saved reference. Scores are from the
root side's perspective, mate reports are separate, and finite-node searches
can disagree. Raw PVs, scores and settings are in
[referee-results.json](referee-results.json) and
[referee-clean-results.json](referee-clean-results.json). The clean assessment
reuses identical existing queries rather than selectively rerunning them.

Mean candidate-minus-baseline referee score, cp:

| Diagnostic | Development | Public holdout | Older v0.10 |
| --- | ---: | ---: | ---: |
| Depth 3, all best ties averaged per case | +7.7 | -13.1 | +19.9 |
| Depth 5, all best ties averaged per case | -49.4 | -8.3 | -0.6 |
| Clean 100 ms choices | +21.5 | -27.6 | +84.1 |

The fixed-depth rows use 15, 15 and nine cp cases respectively. The clean timed
row uses 90, 90 and 54 cp observations respectively, repeated over those same
cases. The tenth old case, djj7sLqm, is already a referee-reported losing mate;
both engines' choices retain mate -4. It is preserved, not converted into cp
or counted as a successful correction. Repeated timed choices are strongly
correlated and are not independent evidence of that many improvements.

At depth five the development set has **0 better / 8 unchanged / 7 worse**
case averages. The holdout has **2 / 11 / 2**; the old cp cases have **1 / 5 / 3**.
For the clean timed slots the counts are **18 / 55 / 17**, **1 / 79 / 10** and
**22 / 32 / 0**, respectively. Counts and averages answer different questions;
large changes on a few selected positions drive the means.

Examples worth retaining:

- Positive: FtRp2MlL changes from Kb8 to Qxe4 at depth five, referee -270 to
  -71 cp, an improvement of 199 cp. CBtJ5xEy changes b5 to Nd6, -193 to -92 cp.
  The old glYlALQg gains 255 cp averaged over its candidate best set. Some old
  timed mistakes improve substantially. These are useful functional examples,
  not evidence that every new choice is good.
- Negative: in 16G1f9rx at depth five the baseline ties Nd4 and O-O; the
  candidate keeps only Nd4. The referee's scores are -461 and -38 cp, so the
  average over the best set deteriorates by 211.5 cp. This illustrates why
  removing a good tie is a regression even when both engines sometimes play
  the same move.
- Negative: BhrzqEje changes its depth-five best set from Qd3/a4 to b4,
  worsening its referee average by 273 cp. PNfQoYP0's broad baseline tie set
  shrinks to b4/h4, with a worse average by about 59 cp. This is not a
  demonstrated pawn-ending improvement, nor a reason to tune this term on
  that public case.

Unrestricted deeper diagnostics:

| Case / depth | Baseline best / cp | Candidate best / cp |
| --- | --- | --- |
| 16G1f9rx / 7 | Nd4 / -12 | Nd4 / +2 |
| 16G1f9rx / 8 | incomplete at 30 s | incomplete at 30 s |
| HhPI3vCx / 6 | Qd6 / +161 | Qd6 / +153 |
| HhPI3vCx / 7 | Qd6 / +159 | Qd6 / +159 |

All completed deeper choices are legal. The two depth-eight searches hit the
independent diagnostic worker deadline and have no completed score or best
move. Their JSON's placeholder zero and `Legal:false` mean unavailable output,
not an evaluated draw or an illegal move. They are retained in
[baseline-live-fullroot.json](baseline-live-fullroot.json) and
[candidate-live-fullroot.json](candidate-live-fullroot.json). No deadline was
extended to obtain a preferred answer. Cold single-depth full-root searches
can cost much more than the earlier restricted or iterative live replay.

## Validation and preservation

[candidate-checks.json](candidate-checks.json) and adjacent full check logs record
focused tests, the complete Go tests, race tests, vet, darwin/arm64 build and an
offline 12-ply lichess-bot wrapper check. The wrapper loads the public
placeholder-token example into a temporary config and checks candidate UCI
identity, legal moves, mate and full-history draw handling. It never starts a
live bot. A true smothered-mate queen-offer control was added after performance
probes; only test source changed at that point.

[provenance.json](provenance.json) records toolchain, executable hashes and
source identities. Sacrifice was locally verified at
`b528f0ba373910b72052d991700a990ca1cb07fd`; no game against it was launched.
[assessment.json](assessment.json) is generated by [summarize.py](summarize.py).
All initial outputs, corrected outputs, watchdog failures and failed fixture
expectations remain separate. There was one set of fixed-depth decisions and
one planned balanced timing repeat to resolve the explicitly identified
workload contamination, not repeated tuning until favorable.

After archiving the exact candidate, all root Go source and tests are restored
to tagged v0.15.1. [restored-checks.json](restored-checks.json) records normal-build
verification. The subsequent v0.15.2 research release changes public release identity and
example greetings only. Active engine configuration is untouched.

## What we learned and the next hypothesis

This experiment establishes that the term detects clear exposure with bounded
weights, honors the tested pins/recaptures, and preserves the tested forcing
sacrifices. Those are functional properties. It does **not** demonstrate a
strength gain. There are no new games or Elo measurements. Neither headline
quiet leaf is explained, the broad fixed-depth and public timed checks are
mixed or unfavorable, and the equal-time depth loss is substantial.

Reject this implementation and do not promote it or run a larger confirmation.
The promising isolated corrections can serve as future public controls, but
they do not offset the failure of the original mechanism/cost case. This does
not disprove every possible piece-safety evaluator. A revision would first
need a concrete missed quiet-leaf threat, cheaper-non-pawn/overload controls,
and a lower-cost implementation, before choosing new weights. Simply raising
these penalties or forcing Nd4/Qd6 alternatives is not justified.

**Leading next investigation: narrow pawn endings and passers**, starting with the actual
`Rd5 cxd5 exd5` pawn ending in DsqFd5Kj and comparing several legal alternatives
to Rd1. Trace king access, blocked paths and races before adding a bonus;
keep PNfQoYP0 and mirrored/side-to-move race controls separate. This is a provisional recommendation because it has a concrete resulting
position to explain, not an instruction to follow the old ranking in order.
I requested a fresh reassessment after publication. It does not validate the old bundled pawn scores. King exposure
remains a later separate experiment. No next feature is implemented here.

## Reproduce offline without overwriting this record

Run from the repository root:

```sh
python3 tools/benchmark/research/v0.16-piece-safety/reproduce.py /tmp/pandafish-piece-safety-replay
```

This creates fresh baseline/candidate source trees from the tag and overlay,
builds distinct binaries, and runs focused tests, original fixed-depth/timed
probes, corrected baseline traces, clean balanced timed probes, deeper full-root
probes and the leaf-feature overlay. Outputs go into the new directory. It does
not run matches or contact Lichess. Run it without other engine tasks for cost
comparability; repeated random ties and timing need not match byte-for-byte.
The deeper deadline failures remain valid possible outcomes.

The diagnostic Go drivers are `.go.txt` data, copied into scratch root packages
so they can call existing private engine functions. The original run's
[run-clean-probes.sh](run-clean-probes.sh), [validate.py](validate.py),
[referee.py](referee.py) and [wrapper-check.py](wrapper-check.py) preserve the
executed method and local paths. Referee replay requires the local Stockfish
executable and python-chess environment, with output/cache paths chosen anew;
these are not dependencies of the released engine.


## v0.15.2 release integration

The research patch restores the playing evaluator and changes public identity
only. [Release checks](release-checks.json) record formatting, full race tests,
vet, build, the offline wrapper, and a concurrent six-game smoke. The smoke
finished 1/1/0 against v0.15.1, 1/0/1 against Skill 2 and 0/2/0 against pinned
Sacrifice, all without forfeits or unfinished games. These tiny integration
samples provide no strength claim. The original candidate did not run matches.
See the [smoke summary](release-smoke/v0.15.2-release-smoke-vs-v0.15.1-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json)
and adjacent PGNs; original labels and hashes are preserved.
