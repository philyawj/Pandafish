# Method and scope

I froze the scope in [PLAN.md](PLAN.md), then implemented the archived
coordinated prototype against tagged v0.15.3. The starting checkout was the
release commit; preexisting first-person documentation edits are listed in
[starting-state.txt](starting-state.txt) and remain intact. The initial experiment stopped before commit/push/PR; I subsequently approved
promotion and authorized the v0.16 release. The proposed larger confirmation
remains deferred. No live client is started.

## Implementation and root contract

- A completed nonempty root returns exactly one proven best move and its exact
  score. Later alternatives search against the incumbent; equality does not
  require proving the whole tie set. The selected move from the preceding
  completed iteration is searched first. Only the emergency unsearched fallback
  remains random. Tie behavior can change actual games even with identical
  fixed-depth scores; it is part of this candidate, not a strength-neutral claim.
- PVS scouts later main-search moves at full depth, re-searching strict
  in-window improvements. No LMR, null-move pruning, aspiration or added quiet
  checks is introduced. The re-search reuses the same immutable child position.
- A separate 65,536-entry board-hint table uses about 1 MiB, in addition to the
  original 2 MiB score table. Hints are matched to freshly generated legal
  moves, and cannot supply scores or bounds. Full keys reject slot collisions.
  `ucinewgame` clears both tables. The score table retains its exact-depth,
  halfmove, replacement and generation rules.
- Only a child actually reached by pawn move/capture (including EP/promotion)
  resets the irrelevant-history fingerprint. Repetition counts remain intact;
  both counts and the preceding fingerprint unwind after complete, cutoff and
  interrupted work. Arbitrary root/search entry points retain their history.
  Castling and lost castling rights stay conservatively history-qualified.
- Nonchecked quiescence copies eligible tactical moves from the library's
  read-only legal cache and stably sorts that copy. Checked positions still
  use every ordered legal evasion. Full legal generation, terminal detection,
  SEE eligibility, horizons and evaluator remain unchanged. This does not
  implement the deferred specialized-generator project.

The tests explicitly replace all-ties assertions with exact-score and
independent best-set membership checks, including an exhaustive mate-distance
proof. `fullWindowRoot` scores each alternative independently without the root
bound optimization or a score table. The existing exhaustive minimax/window,
legal exchange, special-move, horizon and cancellation tests remain. Real
warm-history controls run in both color directions and both cache-seeding
orders. A separate history contamination probe preserves the original -902 cp
versus draw distinction.

## Inputs and fixed-depth comparisons

[make_inputs.py](make_inputs.py) freezes [inputs.json](inputs.json) before
candidate measurements: the original 36 cases, 36 broader cases (first 12 games
per opponent from the released-policy control) and 12 explicit special-rule,
endgame and threat controls. Selection cycles through prescribed plies and
clamps to two plies before game end. All 84 full-history commands are unique.

The deep subset is the original six plus the first two new games per opponent,
chosen before measuring. Every case is searched cold at depths 3 and 4; the
12 deep cases also at depths 5 and 6. Each root has a ten-second independent
safety deadline; completeness is checked rather than treating interrupted
work as a score. The tagged baseline enumerates its exact best set; candidate
score must match and its single selected move must belong to that set.

Cold means fresh score/hint tables; the parsed root position is reused between
requested depths, so library legal-move caches may already be populated. Both
actors use the same driver.

These are public, correlated cases, partly from Pandafish games. They broaden
the architecture diagnostics but are not an unseen population, training split
or a proof of universal score equivalence. Raw fixed timings use a single run
per case/actor, so I use the native equal-time blocks for the timing gate.

## Native equal-time performance

[uci_performance.py](uci_performance.py) records full UCI responses, independent
python-chess legality, completed depth, node counts and actual response time.
Both engines receive the same requested budget and retain the same internal
soft/hard limits. Native UCI searches advance the actual cache generation.

- Serial: 84 cases × two actors × three trials at 100 ms. Actor and case orders
  reverse on alternate trials. One additional 1,000 ms trial per case/actor.
- Three-search load: three separate engine processes per actor, with two blocks
  in reversed actor order and three lanes/trials in each block. Each lane runs
  all 84 positions at 100 ms; there are six trials per actor. The collector's
  raw-output serialization overhead is included in this shared-host workload.
- Warm requests: every available actual next-two-move continuation is searched
  after its parent and separately cold, at 100 ms. Mate/stalemate requests with
  `bestmove 0000` are retained as expected terminal outcomes, not illegal moves.

I run correctness/build checks first, fixed probes next, timed blocks next,
and games only after those gates. No deliberate engine benchmark overlaps the
native timed blocks or strength screen. Background desktop activity, thermal
state and finite timing variance remain uncontrolled limitations. Reported
node definitions are unchanged between these two Pandafish actors; deeper
completion and node throughput still do not establish playing strength.

## Provenance and replay

[validate.py](validate.py) records full race tests, vet, native/ARM64 builds,
offline wrapper validation and all fixed probes. The baseline is rebuilt from
`git archive v0.15.3`, not the working tree. [source-hashes.json](source-hashes.json)
and [candidate-source.json](candidate-source.json) freeze every root Go source/
test and the complete candidate overlay. [provenance.json](provenance.json)
identifies baseline, candidate, Stockfish, Sacrifice and Tal asset hashes.
No opponent revision or playing asset is updated. The host reports eight
logical CPUs. [machine-power.txt](machine-power.txt) preserves the mixed power
report (AC Power with a discharging battery); I do not infer a fixed power or
thermal state from it.

Replay in a new directory from the repository root:

```sh
python3 tools/benchmark/research/v0.16-search-core/reproduce.py /tmp/pandafish-core-replay --fixed
```

This restores and tests source; it does not launch matches or online play.
Original outputs and failed commands stay alongside the successful checks;
[ERRATA.md](ERRATA.md) distinguishes test/harness mistakes from engine behavior.
