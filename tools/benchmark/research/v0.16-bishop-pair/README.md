## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# Archived v0.16 bishop-pair candidate

**Archived in v0.15.3, unpromoted.** I selected unchanged released
play and a broader architecture investigation. The original screen and its
then-current confirmation recommendation remain below as historical evidence.
The [new plan](../v0.15.3-architecture/NEXT.md) supersedes the immediate 200-game
confirmation recommendation; no such confirmation ran.

Original working hypothesis: a bounded interaction for retaining bishops of both square
colors can improve equal-material exchanges at negligible evaluation cost.
This follows a [fresh reassessment](../v0.16-next-assessment/README.md), which
explicitly deferred the previously second-ranked pawn-ending implementation
after new counterfactuals failed to support the obvious compact fixes.

Start from released v0.15.2, merge commit
`64e60edb2f7c210e942c9643453237647c4867b3`. The candidate is labeled
`v0.16-bishop-pair-dev`. Material values, placement formulas, main search,
quiescence, move ordering, table policy, repetition and clocks are unchanged.
No rook-file, mobility, pawn, king-safety or attack-map component is included.

## Implementation and controls

During the existing evaluation traversal, OR a bit for the occupied square
color of each bishop, separately per side. Give **one untuned 25 cp bonus**
when both bits are present. Two same-color promoted bishops receive no bonus;
three or more bishops receive at most one bonus. Both sides having a pair
cancel. Material remains 100/300/300/500/900 cp.

The weight was chosen once as a small quarter-pawn interaction and is not
fitted to the reference score or tuned after viewing regressions. It cannot
model blocked or trapped bishops, relative minor-piece strength in every pawn
structure, or all of the coordination in the selected exchange fixture.

[The original diagnosis](../v0.16-next-assessment/bishop-diagnosis.json) and
[fixture baseline](fixture-baseline.json) reproduce Bxc6 at depths three and
five. [The candidate](fixture-candidate.json) instead chooses a pair-preserving
alternative in unrestricted searches at both depths and in all three 100 ms
trials. This is a selected functional correction, not a strength measurement.
The reference's move is not required by the implementation or root filtering.

Focused tests cover single bishops, actual two-color pairs, promoted same-color
bishops, excess bishops, both-side cancellation, both perspectives, color/rank
and file mirrors, unchanged material, capturing a free queen and the
unrestricted smothered-mate queen-offer control. The source and tests are
preserved in [candidate-source.json](candidate-source.json).

## Frozen broader checks

[probe-inputs.json](probe-inputs.json) reuses all 30 public recent cases with
their original 15/15 split and all ten older v0.10 cases. It adds only the
pre-diagnosed synthetic exchange family and two geometric mirrors. These are
public regression cases, not blind or representative positions. There is no
retuning on the holdout. Every fixed-depth best tie and timed choice is retained.

Fixed depths are three and five. Timed probes use the production 100 ms budget,
95 ms internal hard deadline, cold tables and full reversible history. A balanced
baseline/candidate/candidate/baseline schedule uses three trials per position
per round. Static microbenchmarks use the same three boards and toolchain as
the rejected piece-safety test. No other engine task deliberately overlaps
performance probes; no thermal/frequency telemetry is collected.

The conditional game screen is the repository default: 25 color pairs each
against the tagged v0.15.2 executable, Stockfish 19 Skill 2 and Sacrifice v0.2.0
at `b528f0ba373910b72052d991700a990ca1cb07fd`, all three matchups concurrent.
A failed previous-release gate will be preserved. There is no automatic
confirmation, commit, push, PR or live client for this candidate.

## Pre-screen diagnostics

The fixed source and 25 cp weight were frozen in [preflight.json](preflight.json)
before the routine screen. All 688 main diagnostic probes complete legally:
86 fixed-depth and 258 timed probes per engine. This excludes the separately
preserved initial five-probe exchange smoke per engine.

| Balanced 100 ms probes | Tagged v0.15.2 | Candidate |
| --- | ---: | ---: |
| Total nodes | 4,174,263 | 4,105,648 |
| Mean completed depth | 3.783 | 3.740 |
| Mean elapsed ms | 92.913 | 93.194 |
| Maximum elapsed ms | 95.054 | 95.067 |

Aggregate candidate nodes fall 1.64%. Completed depth is higher/equal/lower in
3/241/14 aligned slots. Scores alter trees, so this is not a pure evaluator
throughput measurement. The microbenchmark medians in ns/op are 562.2/612.7
at start, 467.5/506.5 on the middlegame board and 309.5/326.5 in the pawn/rook
ending, baseline/candidate respectively. This is a small absolute overhead,
not zero cost. No thermal or CPU-frequency telemetry was collected.

[assessment.json](assessment.json) summarizes every result; [summarize.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz)
reproduces that summary. Raw fixed best sets and balanced timed choices stay
in this directory. The synthetic family changes from bishop exchange to knight
exchange in all three reflected boards at both fixed depths and in all timed
slots. The three reflections are one underlying example, not three independent
discoveries.

Every selected move and fixed best tie is compared with the same full-strength
Stockfish 19 100,000-node/fresh-hash setup, one thread and 64 MiB, with full
history. Identical historical queries are reused from the prior archived
referee run; new moves are analyzed with the same settings. This is separate
from Skill 2 match play. Stockfish is a finite-budget comparator and its single
move is not a required answer. Mate cases remain separate.

Mean candidate-minus-baseline referee cp differences:

| Diagnostic | Development | Public holdout | Older v0.10 | Synthetic family |
| --- | ---: | ---: | ---: | ---: |
| Depth 3, all best ties averaged | +6.6 | +7.5 | +5.9 | +235.7 |
| Depth 5, all best ties averaged | +18.5 | +0.3 | +68.3 | +235.7 |
| Clean 100 ms choices | -15.7 | +7.9 | +8.0 | +235.7 |

There are 15/15/9/3 cp boards in each fixed-depth row and 90/90/54/18 repeated
slots in the timed row. One older board is still reported as mate -4 for both
engines and is not converted to cp. Repetitions and mirrors are strongly
correlated. No statistical strength inference follows from these means.

Positive depth-five examples include a6 instead of the Nd4/O-O tie in
16G1f9rx, castling instead of Bxf3 in qaE8fqwI, and two older cases.
Serious regressions are retained too: O4wRsvjl returns to Kf7 at depth five,
worsening the referee average by 211 cp. BhrzqEje also worsens slightly. The
preflight summary's phrase about small regressions was too narrow: this
per-case review exposes a large king-choice regression. No weight or source
behavior is changed in response. The clearer controlled exchange mechanism,
small equal-time cost and generally neutral/positive fixed-depth comparisons
justify only the limited screen, not a promotion or strength claim.

A final test-only refinement adds a pawn recapture after the queen capture,
so the tactical control explicitly gives up the bishop pair to win a queen.
It does not change the playing source or the screened binary.

## Standard concurrent screen

One frozen 50-game screen per opponent, 150 games total, completed with three
matchups concurrent. Settings: 100 ms per move plus 50 ms protocol grace,
400-ply cap, the ten-position screening book (SHA-256
`80ceed68e4575b7a0d8649aa67efd35da4c8c46e480ad31b1ce12fefa94e89ca`),
color-swapped opening pairs, Stockfish Skill 2 / one thread / 16 MiB hash /
strength limiting off, and pinned Sacrifice with its verified Tal asset.
The machine is Apple M2, darwin/arm64, Go 1.25.0. Background system work was not
controlled and no thermal/utilization telemetry was collected.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.15.2 | 26 / 5 / 19 | 0 | 0 | 57% |
| Stockfish 19 Skill 2 | 21 / 2 / 27 | 0 | 0 | 44% |
| Sacrifice v0.2.0 b528f0ba | 7 / 11 / 32 | 0 | 0 | 25% |

The previous-release practical criterion passes: strictly over 50% with no
unfinished games or forfeits. It is not a significance test. All matches ran
to the frozen sample; there was no gate-based early cancellation, discarded
failure, favorable-score rerun or larger confirmation.

[Complete screen, exact identities/hashes/settings and every game](../../results/v0.16-bishop-pair-dev-vs-v0.15.2-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json),
[PGNs](../../results/v0.16-bishop-pair-dev-vs-v0.15.2-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/games.pgn)
and [screen command output](screen.txt) preserve the original candidate label.
The executable hash is
`f3872e955bfa820b2d7b56681e1d8367e491640dd305cb15341b83050cfe2b00`.
No release binary is relabeled as the development screen.

The candidate returned 8,577 moves across the screen. None exceeded the
requested 100 ms response budget. Maximum responses by matchup were 99.105,
97.771 and 97.913 ms; p95 values were 96.558, 96.782 and 96.667 ms. The external
engines often slightly exceed 100 ms within the explicit protocol grace;
there were no forfeits. These observations are not deadline guarantees.
[screen-timing.json](screen-timing.json) retains counts and timings.

## Uncertainty and recommendation

**Hold the frozen candidate for confirmation; do not promote it yet.** The
synthetic exchange behavior is a functional improvement, and the 57% direct
screen is a positive strength signal. It does not establish an overall gain.
The external results are weak; they do not independently compare this exact
candidate against a simultaneously screened baseline on those opponents.
Different earlier 50-game runs cannot establish a regression or improvement
by simply subtracting percentages. Do not pool three opponents into one score.

[sampling.json](sampling.json), generated by [sampling.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz), resamples
100,000 times with seed 160016. Approximate descriptive 95% score intervals:

| Opponent | Resample color pairs | Resample opening clusters |
| --- | --- | --- |
| v0.15.2 | 44% to 70% | 43.3% to 70.2% |
| Skill 2 | 32% to 56% | 32.3% to 54.5% |
| Sacrifice | 17% to 33% | 19% to 30% |

There are only 25 pairs and ten repeated opening clusters per opponent. These
intervals condition on the observed book/candidate and are not calibrated Elo
bounds, a sequential test or a correction for selecting features after earlier
experiments. The direct interval comfortably includes equality. The public
regression/fixture selection is another limitation.

**A separately planned 200-game confirmation per opponent is justified**, if
I want to pursue this source: the direct gate passes, the
mechanism is concrete, and equal-time cost is small. Keep the source/weight
frozen rather than tuning to the mixed screen. Use the public 50-position book,
three concurrent matchups and the same binaries/settings. Even 200 may remain
inconclusive for a small gain. Do not automatically add further games until
favorable. No confirmation has been started.

The known Kf7 regression is retained as a negative control and an item for
further diagnosis. It is not patched by this candidate. A pair bonus does not
understand mobility, trapped bishops, king exposure or every exchange. The
screen does not justify adding those terms to this experiment.

I need to review the results before any new commit, push or PR.
Current working source remains isolated on a local development branch; the
published v0.15.2 playing policy and active bot executable/config are untouched.
No live client or online game was started.

## Validation, preservation and replay

The screen summary records passing gofmt, `go test -race ./...`, vet and build.
[tests.txt](tests.txt) records the normal Go tests; [wrapper-check.txt](wrapper-check.txt)
records the offline 12-ply public-config wrapper, legal mate and history draw.
The named candidate was built explicitly for darwin/arm64 before the wrapper.
The stronger final tactical test passes in [controls-final.txt](controls-final.txt).
It changed tests only, not production code or the screened binary.

[candidate-source.json](candidate-source.json) is the frozen pre-screen overlay
and [candidate-source-final.json](candidate-source-final.json) adds only the
post-launch test refinement. [provenance.json](provenance.json) records both
hashes, exact baseline/opponents, unchanged policy-file checks and binary hashes.
No experimental piece-safety code is in this build.

Build/control replay of the final overlay passed in
[archive-smoke.txt](archive-smoke.txt). To reproduce in a new scratch directory,
without touching this record or running games:

```sh
python3 tools/benchmark/research/v0.16-bishop-pair/reproduce.py /tmp/pandafish-bishop-replay
```

Add `--checks-only` for the archive build/control smoke. The Go diagnostic
drivers are data files copied into scratch root packages. The original local
[run-probes.sh](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) and [referee.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) preserve the methods;
the latter requires the installed Stockfish and python-chess environment.
Random root ties and timing need not reproduce byte-for-byte.

Run the command below only if I approve the proposed confirmation. Verify the frozen candidate SHA above first; this is **not executed**:

```sh
go run ./tools/benchmark \
  --candidate tools/benchmark/bin/pandafish-v0.16-bishop-pair-dev \
  --version v0.16-bishop-pair-confirmation-dev \
  --previous tools/benchmark/bin/pandafish-v0.15.2 --previous-version v0.15.2 \
  --stockfish /opt/homebrew/bin/stockfish --stockfish-skill 2 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --pairs 100 --parallel-all \
  --openings tools/benchmark/openings-confirmation-50-v1.txt
```

This is a broader public confirmation, not blind validation. Screen and
confirmation stay separately labeled and are never selectively pooled.
