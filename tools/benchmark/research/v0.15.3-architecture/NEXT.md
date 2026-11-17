# Best larger projects to try next

**First: a coordinated search-core project, beginning with the archived
pipeline prototype and then removing the full-legal-move work from quiet
capture search. Second: a coherent, data-calibrated evaluation/threat model.
Keep repetition correctness throughout.** These are ranked experiments, not
promised Elo gains or a fixed checklist. Reassess after each result.

The bishop-pair source stays archived. Its practical 57% direct screen is not
a demonstrated broad gain. The maintainer chose unchanged released play for
v0.15.3 and broader architectural research ahead of its optional confirmation.

## 1. Search core: useful depth with much less work

This is the strongest measured engineering opportunity. The present profile
spends 50.9% of sampled CPU under `orderedMoves` (generation plus sorting), and
the timed kernel allocates about 72.8 MB / 120,241 objects per search. Cold
fixed-depth nodes are 91.3% quiescence. It generates 6.61 million full legal
moves to visit 0.453 million eligible tactical/check-evasion moves at depth four.

A coordinated prototype already supplies a more substantial starting point
than another chess bonus: **37% fewer nodes at depth six** on the six frozen
deep cases, preserving all tested fixed-depth scores and selecting members of
the released best set. It combines a single-best root contract, PVS,
history-safe board hints/reversible fingerprints and tactical-only sorting.
The components have separate ablations. It is exploratory, has no strength
match, and intentionally fails the old all-root-ties API assertion. Its apparent
promise must not be promoted as an established strength gain.

### Stage A: turn the prototype into a reviewable candidate

- Start from tagged v0.15.3's released evaluator, with no bishop-pair or other
  new evaluation weights. Preserve the archived prototype exactly for comparison.
- Specify that production returns one proven best root move, rather than
  certifying every exact tie for random selection. Keep the diagnostic ability
  to score individual alternatives separately. Update contract tests explicitly:
  require the selected move to belong to the independently established best
  set, with correct score and mate distance. Do not weaken legality, interruption
  or mate-distance proofs.
- Keep board-only **hints** separate from history-qualified **scores**. Reset
  only irrelevant cache fingerprint history at irreversible moves, while
  preserving correct repetition counts and unwind/cancellation behavior.
- Harden the GHI tests with different real histories, warm-table searches,
  pawn/capture resets, castling-right changes, en passant, halfmove counters,
  repetition choices and mate-before-draw cases. Merely refusing to cache draws
  does not solve the demonstrated nonzero-score contamination.
- Validate full-depth score parity and move membership on a broader corpus;
  measure cost at 100 ms and at longer budgets, including the routine three-search load. Then run the standard 50 games
  each against the new released baseline, Skill 2 and pinned Sacrifice, all
  concurrent, saving failures. A positive screen leads to a separately chosen
  confirmation, not an automatic promotion.

**Stop or revise** if draw/cancellation behavior fails, equivalent searches
change scores unexpectedly, the root-contract trade produces material tactical
regressions, or equal-time gain disappears outside the selected cases.
Do not transplant the old PVS/LMR bundle or change only `==` to `>=` in the TT.

### Stage B: a search-oriented move and position interface

The tactical-sort prototype still generates every quiet legal move before
filtering. Its approximately 15% extra reported nodes at equal time shows
some avoidable overhead, not the full opportunity. The next larger change is
legal tactical generation and a staged move picker with reusable storage.

Evaluate two implementation routes before committing to a rewrite:

1. Extend/fork the pinned library with a search API: efficient `hasAnyLegalMove`,
   legal captures plus all promotions, full evasions in check, caller-owned
   move buffers and incremental board metadata. The existing library already
   has bitboards; the problem is its general-purpose traversal, legality tags,
   defensive copies and our call pattern, not the absence of bitboards.
2. If the API cannot eliminate the measured cost, introduce a compact internal
   search position with make/unmake, an undo stack and incremental hash/state,
   retaining the current library as UCI/game-rules frontend and test oracle.

Preserve non-capture promotions, underpromotions, en passant discovered checks,
castling legality, pins, double checks, stalemate detection and all check
evasions. A capture-only list being empty is not proof of stalemate. Never sort
`ValidMovesUnsafe()` in place. Keep the old evaluator/search policy while
measuring the backend, so throughput gains have an identifiable cause.

Require independent perft and randomized differential legal-move sets, board/
hash/make-unmake round trips, existing mate/repetition tests and unchanged
full-depth scores. Use a second oracle for special-rule tests rather than
assuming the old library is infallible. Preserve exact failed positions.

An engineering acceptance target such as **at least 1.5x useful search
throughput with markedly lower allocation** is a project goal, not a prediction.
If only a few percent survives, stop expanding the rewrite. Even a 2x speedup
would not imply a particular Elo gain or fix every wrong leaf model.

### Stage C: selective search on a sound, efficient foundation

Only after the above is measured, revisit a coherent move-ordering/selective
search design: PV/non-PV treatment, history ordering, aspiration, verified
reductions and checking-move protection. Main-node beta cutoffs already happen
on the first move about 78% of the time in this shallow suite, so history alone
is not established as the largest opportunity. Reductions require tactical and
zugzwang controls plus actual game validation.

The earlier PVS/LMR attempt was before v0.15's repetition fingerprint, so that
new guard cannot explain its original failure. The new prototype differs in
root contract and cache/move-pipeline coordination; it is not evidence that a
blind retry of the old settings will work.

## 2. Evaluation and threat modeling with data, not independent bonuses

Search speed alone is insufficient. In the selected Qd6 line, the released
quiet leaf scores +161 cp for Black while the finite-node referee reports
-662 cp. In the Kxg5 line a non-checked leaf scores +400 cp for Black but the
referee finds a mate in one for White. In the bad rook exchange, the engine
chooses a different opponent continuation because it undervalues the actual
pawn ending after `Rxd5 exd5`. These are different failure mechanisms.

Build a larger position-and-response corpus containing engine lines,
stronger opponent deviations, quiet threats, king attacks, exchanges and
endgames. Split by game/opening family **before** tuning; the 30 public cases
are regressions, not a training holdout. Measure feature cost using the new
search-state summaries, avoiding a fresh expensive board/attack reconstruction
at every leaf.

Choose and compare two serious model routes:

- **A compact tapered classical model:** coherent king shelter/exposure and
  attacker access; mobility and pinned/blocked pieces; pawn structure, promotion
  paths, king blockade and pawn breaks; coordinated phase-sensitive placement.
  Calibrate a bounded parameter set against independent data, with ablations,
  rather than guessing each term on one saved game. Earlier failed bundles
  reject those tested weights/designs, not every classical model.
- **A compact trained evaluator / NNUE route:** treat representation, incremental
  updates, data/label quality, training/validation, runtime throughput and
  provenance/licensing as one project. Do not drop a Stockfish net into an
  incompatible evaluator or reproduce its current parameter count by default.
  The [official training project](https://github.com/official-stockfish/nnue-pytorch)
  is a primary reference for the workflow, not evidence of a free gain here.

Separate static evaluation from threat-search coverage. A quiet mate-in-one
at the quiescence boundary suggests a threat-handling control, but the previous
naive quiet-check extension timed out. Any replacement must have a measured
budget, termination controls and holdout tactics, not unbounded extra checks.
Do not require any particular Stockfish move when several are sound.

**Acceptance:** out-of-sample decision-quality improvements across categories,
acceptable equal-time cost, retained tactical/rule controls, then controlled
strength screens and an independently planned confirmation. Lower supervised
loss or a few corrected moves do not establish playing strength.

## 3. Repetition-safe TT policy as architecture, not a magic size knob

The safe cache redesign belongs in Stage A and scales better at depth six.
Further score-sharing needs a reasoned GHI design and explicit exact-depth
versus minimum-depth semantics. Board-only scores are rejected by a reproduced
-902 cp versus draw contamination. A 32 MiB table did not help cold shallow
nodes. Neither is the next large bet on its own. See [CACHE.md](CACHE.md).

## Testing strategy for larger ambitions

Retain 50 games per opponent, three concurrent matchups, as triage. Design
future tests around a stated effect and frozen source, include both 100 ms and
longer diagnostic budgets, use a broader independent opening/position set,
and plan confirmations before running them. Preserve separate opponents and
all failed attempts. The research does not implement SPRT or justify adding
low-confidence feature Elo guesses.

After each candidate, reassess whether the limiting evidence is search work,
wrong leaf models or actual game outcomes. The best immediate task is Stage A,
with Stage B as the larger throughput project. Evaluation is the next major
strength workstream, not another 25 cp bonus while the engine continues to
spend most of its time generating and sorting irrelevant quiet moves.
