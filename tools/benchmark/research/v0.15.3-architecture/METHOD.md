## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# Methods, records and reproducibility

## Frozen position suite

[make_inputs.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) freezes [inputs.json](inputs.json) before probes:
24 positions from the first 24 games of the older conversion previous-version
match, at cycling plies 16/32/48/64, clamped to two plies before game end if
necessary. No selection by result or new-engine output. Add the starting board,
a sparse transposition control, the known repetition control, all eight live
error cases and the diagnosed pawn ending: 36 inputs in total.

The games and errors are public and correlated. This is neutral with respect
to the new architecture results, not blind validation or a representative Elo
sample. The older game distribution can favor Pandafish's familiar positions.
The six deeper inputs are fixed by input order: first four neutral cases,
starting board and sparse control, at depths five and six. No favorable subset
is selected after seeing savings.

## Isolated variants

[setup.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) creates small scratch copies of tagged v0.15.2 root Go,
module, internal and testdata files. Each [overlay JSON](overlay-safe-cache.json)
records exactly what differs from that tag. No diagnostic source is compiled
by the repository's normal build. No bishop-pair evaluation is in these trees.

The variants are baseline, no table, board-only hints, reversible fingerprint,
both safe-cache changes, unsafe board-only scores, minimum-depth score reuse,
16x table capacity, one best root move, PVS, PVS plus safe cache, tactical-only
quiescence sorting, and an instrumented baseline. All identities/hashes are
preserved in the UCI result files.

The coordinated-kernel prototype is an explicitly exploratory combination,
chosen after individual ablations: one-best root, PVS, safe cache and tactical
ordering. [setup_combo.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) records its construction. This is
not an independent confirmation, tuned weight or validated strength candidate.

## Fixed-depth and correctness checks

[run_variants.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) builds each scratch tree, runs root/internal
tests and the adversarial graph-history control, then depth-three/four searches
with full history and an independent ten-second per-search limit. Every best
root tie is retained except in variants explicitly changing the tie contract.
Completed scores and membership in the baseline best set are checked separately.
Incomplete placeholders are not interpreted as evaluations.

Warm fixed-depth probes repeat the root and search the recorded next position
with warm and fresh tables. This toy harness keeps table generation fixed;
native UCI probes separately test actual per-request generation semantics.
[run_followup.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) adds the prespecified depth-five/six probes.

The first normal test invocation accidentally ran diagnostic functions without
an output path. Those driver failures are retained in `*-tests.txt` and
[checks.json](checks.json). The driver now skips diagnostics unless explicitly
requested. Corrected root/internal tests are in `*-tests-corrected.txt` and
[native-checks.json](native-checks.json). The archived v1 driver permits
reproducing the harness mistake. No playing code was repaired for this error.

Expected or real failures remain distinct:

- No-table tests fail assertions that a functioning cache exists.
- Unsafe board-only scores fail draw-history isolation and the actual warm GHI
  control. This is a substantive correctness rejection.
- Minimum-depth reuse fails exact-depth score/window expectations. It requires
  a deliberately changed search contract, not silently weakened assertions.
- The large table fails the existing 2 MiB resource limit.
- One-best root and coordinated variants fail the test requiring *all* tied
  longest-resistance mate moves. Their selected moves and mate-distance scores
  remain correct in the independent proof; the API contract changed.

Safe hints, reversible fingerprint, their combination, PVS, PVS plus safe cache
and tactical sorting pass all root/internal tests. The release is separately
validated in the full repository after restoring tagged playing behavior.

## CPU and allocation profiling

A native baseline test kernel rotates the 24 neutral positions, clears a
preallocated table, and invokes the production iterative search with normal
100 ms time management and cancellation context. It includes deadline checks,
allocations and garbage collection. It excludes repeatedly allocating a new
2 MiB table per move. Go's benchmark calibration also runs, so profile duration
exceeds the final reported benchmark duration.

The first joint profile used allocation sampling at 4 KiB. It visibly spent
17.6% of CPU collecting allocation stack traces, so it is **not** the primary
normal-CPU estimate. Its allocation ownership remains useful. Both raw profiles
and text reports are retained. The follow-up CPU profile disables allocation
stack sampling via the benchmark-only runtime setting; garbage collection and
allocation counters remain enabled. Use `native-cpu-clean.pb.gz` and
`cpu-clean-*.txt` for CPU claims, and the separately sampled memory profile for
approximate allocation ownership. No such runtime setting is in a playing binary.

Cumulative profile rows overlap and must not be added. The 50.9% `orderedMoves`
row includes generation and sorting. The 40.9% `CalcMoves` row is nested work,
not another independent 40.9% cost. Bytes/op and allocations/op describe this
selected timed workload, not every chess position or a universal GC overhead.

## Native equal-time probes

[uci_probes.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) runs one searching process at a time, records raw
UCI output and legal moves, and has independent watchdogs. Startup and readiness
are outside the request timer; full-history `position` parsing is inside.
Fresh `ucinewgame` resets each cold request. Three trials use reversed engine
and case order in the middle trial; this tempers order effects without making
tiny differences statistically robust.

- 100 ms: 36 positions, three trials, native released baseline, pinned Sacrifice
  and ten uninstrumented diagnostic variants.
- 1,000 ms: 36 positions, one trial, baseline/Sacrifice/safe-cache.
- Warm: 25 parent/next/cold-next episodes, actual UCI cache generations, three
  engines. Two episodes end in mate after the recorded two moves. Their null
  best moves are expected terminal results, not illegal moves or fallbacks.
- Historical: tagged v0.14.4 versus v0.15.2, three 100 ms trials, plus one
  iterative depth-four request on each position. These versions differ in
  mate/draw correctness too; this is not a one-variable cache strength test.
- Coordinated: a new baseline/combo block, three 100 ms trials. Baseline is
  repeated as a contemporaneous control, not rerun to obtain a favorable score.

Requested time budgets match, but engine margins remain unchanged: Pandafish
internally uses 95/995 ms; Sacrifice nominally uses 100/1,000 ms. Report actual
response times rather than treating budgets as identical consumed CPU time.
Sacrifice uses a larger table and a different search/evaluation policy. Its
node counter counts depth-zero main calls before entering quiescence, whereas
Pandafish counts that transition only in quiescence. Nominal depths also differ
under LMR. Reported depth/NPS therefore do not directly measure relative strength.

## Move-quality and leaf diagnostics

[referee.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) checks every distinct selected legal move and an
unrestricted comparator for each full-history position using Stockfish 19,
Skill 20, limiting off, one thread, 64 MiB and 100,000 nodes per fresh-hash query.
All score-bearing iterations are retained; the last completed exact iteration
is used instead of a later incomplete bound update. Identical queries are
reused when the coordinated block is added; the earlier file is retained.

[analyse.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) retains raw and clamped cp differences, negative
reference disagreements, mate transitions and per-case results. It averages
repeated choices within a case before group means. Different mate/cp mixtures
must not be pooled into one cp or Elo statistic. Equal-budget move-quality
probes are diagnostics, not games or a calibrated strength predictor.

[leaf_trace_test.go.txt](leaf_trace_test.go.txt) reconstructs one maximizing
line per saved diagnostic alternative from the released evaluator, including
capture quiescence. It uses exact re-searches, deterministic first ties and a
separate watchdog, not a warm live PV reconstruction. [leaf_referee.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz)
checks the resulting 16 leaf positions with full histories and the same finite
node referee. An apparently reasonable leaf on the engine's own line does not
validate its evaluation of a stronger opponent response. The counterfactual
pawn ending is a separate example of that distinction.

## Games and release separation

The retrospective control is the standard 50 games each, all matchups
concurrent, using the exact tagged v0.15.2 binary against pre-v0.15 v0.14.4,
Skill 2 and pinned Sacrifice. It is not a larger bishop-pair confirmation or a
strength test of any architecture prototype. The supplied previous-revision
label was accidentally wrong; [control-provenance.json](control-provenance.json)
records the actual tagged source and executable hash. Original metadata and
results remain unchanged rather than being silently relabeled or rerun.

The original bishop-pair 150 games are retained separately. v0.15.3 restores
v0.15.2 playing behavior and changes public identity only. Release smoke games
validate integration, not strength. No live client or automatic 200-game
confirmation is launched.


The reused leaf driver also emits 16 empty optional reference-continuation
records because the eight-case replay input contains root alternatives rather
than those PV lists. Only the 16 nonempty engine traces are assessed; no
score is inferred for empty optional records. The initial reference-extension
driver wrote incrementally only on new queries. A final-write correction records
the all-reused-query metadata as well; old reference results are retained and
no identical query is re-searched for a preferred score.
