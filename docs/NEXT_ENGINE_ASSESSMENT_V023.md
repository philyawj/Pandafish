# Ranked engine assessment after v0.23

Assessment date: 2026-11-21. Released baseline:
`28948e7473251d950003ec1f966f0cf26369b817` (v0.23).

Subsequent repository decision: [existing evidence stays in place](EVIDENCE_STORAGE.md).
The archive pilot referenced in this original assessment is deferred. The next
engine project remains the bounded legality/attack-context diagnosis; cleanup
publication does not publish a playing candidate or another release.

**Recommended first project: shared legality and attack context for the search
node, beginning with a bounded opportunity diagnosis.** Current sampling and
source inspection support a coherent next backend project, while the simpler
allocation/tagging improvements have already shipped. A contextual evaluator
has greater potential to change playing decisions, but a lower probability of
near-term success and much higher data/validation cost. Neither has a measured
strength gain yet.

This is an assessment of alternatives, not an implementation checklist. No
candidate, fitting campaign, strength screen or live game was started. The
[repository audit](REPOSITORY_AUDIT_V023.md) covers tooling, storage and guidance
remedies separately. Its cold-evidence pilot is useful repository work, not a
playing-strength candidate.

## What the successful architecture chain actually established

The [v0.15.3 investigation and original ranking](../tools/benchmark/research/v0.15.3-architecture/NEXT.md)
identified substantial repeated work, an unnecessarily strong all-root-ties
contract, safe versus unsafe history reuse, and distinct leaf-model/threat
failures. It did not prescribe a universal board rewrite. Subsequent outcomes
support that distinction:

| Project | Substantive mechanism and measured result | Playing-strength limitation |
| --- | --- | --- |
| [v0.16 search core](../tools/benchmark/research/v0.16-search-core/README.md) | One-best root contract, PVS, separate board hints and safe irreversible-history identity; 35.2% fewer depth-six nodes in the selected 12 roots. | 65% direct /43% Skill 2 /49% Sacrifice screen; color/opening concentration and no larger confirmation. Components have no individual strength attribution. |
| [v0.17 search backend](../tools/benchmark/research/v0.17-search-moves/README.md) | Tactical destination generation, move-existence API and per-ply buffers; selected identical-work throughput about 1.65–1.76x. | 59% direct /53% /53%; a native shallower case and uncertainty retained. |
| [v0.18 exact evaluation](../tools/benchmark/research/v0.18-exact-evaluation/README.md) | Compute the same material/placement formula from bitboards and tables; selected fixed-work time about 12–15% lower. | 57% /73% /51%; direct gain concentrated by opening. The rejected learned values did not ship. |
| [v0.20 legality kernel](../tools/benchmark/research/v0.20-search-legality/README.md) | Specialize transitions and avoid discarded legality/tag work; 38.7% median selected fixed-work reduction, 325 exact comparisons. | 65% /67% /58%; transition-only predecessor failed its screen. |
| [v0.21 stable rank order](../tools/benchmark/research/v0.21-rank-order/README.md) | Compute priorities once and sort small records, preserving stable ties; 18.9% lower selected depth-four time. | 55% /50% Skill 3 /71%; modest direct evidence. |
| [v0.22 generation/exchange kernel](../tools/benchmark/research/v0.22-bitboard-exchange/README.md) | Defer unused generation, prefilter evasions and use exact bitset exchange; 37.84% total warm saving in one tagged pair, plus three additional-gain pairs. | 64% /71% /70%; public 100ms teacher-loss gate remains failed under its explicit pre-screen exception. |
| [v0.23 state kernel](../tools/benchmark/research/v0.23-state-kernel/README.md) | Worker-owned child storage plus virtual ordinary tags; 14.142% less warm time, exact 1,145 roots /43,118,858 nodes in three pairs. | 51% /69.4% completed /76%, one capped game. One-game direct lead establishes no meaningful strength or Elo gain. |

Skill 2 rows and Skill 3 rows are different opponent profiles. Separate releases'
external results are not simultaneous controls, and the efficiency ratios use
different evolving baselines/cohorts. They must not be multiplied into a
historical strength estimate. v0.19 was a research/tooling release with v0.18
playing policy, not an omitted engine gain.

The successful work removed costs repeatedly paid by the real search. Some
projects were focused, but acted on large call counts. A large rewrite does not
automatically offer more headroom. Smaller exact components that missed an
ambitious size target could still have value: the separate v0.23 storage and
virtual-tag attempts saved 9.137% and 4.543%, and the later prospectively scoped
combined project earned promotion. Their original gates remain failed.

## What the failed work constrains

[v0.19 evaluation history](RESEARCH_V019.md) rejects the tested low-dimensional
linear representations, targets and selection protocols. Cleaner splits and
lower offline regret were insufficient: one activated fit failed short-budget
choices and added a losing-mate report. The later
[450-coefficient piece–king model](../tools/benchmark/research/v0.23-king-context/README.md)
also worsened validation and historical holdout before runtime activation.
The [direct-leaf target study](../tools/benchmark/research/v0.23-direct-leaf-context/README.md)
missed its opportunity gates. Another fit on those inspected records is not
the recommended project.

The [v0.23 attempt index](../tools/benchmark/research/v0.23-next-assessment/README.md)
preserves 23 unsuccessful substantive attempts: 18 engine prototypes, one offline
model and four diagnoses. Only four full screens ran across that investigation,
including the promoted kernel. Cost failures and failed decision surrogates
are not 23 strength defeats. Important counterevidence includes:

- Quiet-check generation was about 78% faster than its old route, but the
  changed-horizon candidate scored 50% direct and failed promotion.
- Root parallelism saved 25.86% wall time using extra CPU, yet scored
  50% direct /47% Skill 3 /80% Sacrifice after a failed fresh control.
- Quiet-count reductions saved 43.128% nominal-depth time and passed all six
  decision gates; their 57% direct /45% Skill 3 /77% screen still did not merit
  promotion. Nominal depth and finite teacher loss were insufficient selectors.
- A narrower check-budget combination passed cost and fresh 100ms, but failed
  fresh long-budget Black. More tactical visibility did not establish balance.
- Static-margin pruning had 20.808% outer-node opportunity but a hidden losing
  mate. The attempted quiescence verification was logically ineffective because
  stand-pat guaranteed acceptance. That line remains closed.
- Aspiration, history, root scouts, capture staging, mate bounds and strict
  quiescence reuse found limited benefit under their tested designs. Familiarity
  with those techniques is not new evidence for another sweep.

These results constrain proposals; they do not prove a playing-strength ceiling.
Live mistakes provide regression controls and operational facts. They do not
determine this ranking, prove one common cause, or make a book/dataset mandatory.

## Current mechanics and fresh baseline evidence

Source inspection of [search](../search.go), [iteration](../iterative.go),
[tables](../transposition.go), [evaluation](../evaluation.go),
[move picking](../main_move_picker.go) and [the local backend](../third_party/chess/search_moves.go)
establishes these boundaries:

- Bitboards, mailbox, incremental Zobrist updates, PVS, iterative deepening,
  tactical generation, deferred quiet generation, killers, board hints,
  history-qualified scores and scratch child storage already exist.
- Each main/q node separately asks `SearchStatus`, which starts an existence
  search, before normal generation or a stand-pat cutoff. Ordinary legality
  and opponent-check classification still query attacks for each candidate.
  Evasion masks already optimize checked nodes; they are not a general pin map.
- Child storage avoids allocation but copies the whole board. Evaluation
  recounts phase and traverses occupied pieces at each q node; it does not
  maintain incremental material/placement totals. Hash updates are already
  incremental and should not be presented as a missing feature.
- The evaluator is material plus small additive placement formulas. Except for
  phase and occupancy, it cannot represent the relationship between a piece,
  its defenders, mobility, pawn blockade and king shelter. More time searching
  this leaf model can expose another bad horizon.
- Quiescence has captures/promotions and complete check evasions, with SEE
  filtering and finite caps. It has no ordinary quiet-threat layer; at the last
  checked safety cap it returns static evaluation. Proven quiet-mate omissions
  remain genuine, but their population frequency is unmeasured.
- The direct-mapped score table uses exact requested depth and repetition-count
  context; deeper entries can order moves without supplying scores. Q scores
  are uncached. No automatic inference of safely reusable cross-history scores
  follows from the existence of transpositions.
- Search handles mate/stalemate and threefold repetition. It does not separately
  score move-count draws or dead-material draws, although the game/benchmark
  layer has additional draw handling. Time allocation uses a fixed horizon and
  soft fraction, without a demonstrated instability-based allocation policy.

The unchanged v0.23 CPU sample uses the established warm trajectory, not selected
Lichess losses. All 1,145 outputs and nodes match the archived candidate.
[Method and raw profile](../tools/benchmark/research/post-v0.23-audit/README.md)
record 16.97s sampled CPU over an 18.96s profile duration:

| Function/subtree | Sampled cumulative CPU | Interpretation |
| --- | ---: | --- |
| Library `Position.searchMoves` | 27.52% | Generation, filtering and tags; largest specific backend subtree. |
| Ordinary/full move tagging | 11.79% | Inside generation; do not add to its 27.52%. |
| `SearchStatus` | 7.07% | Existence traversal, also inside generation. |
| `SearchPosition.Update` | 12.37% | Child-copy/transition/hash work; incremental hash is a subset. |
| Static evaluation | 8.84% | Exact computation and potential future inference budget. |
| Exchange | 8.66% | Already uses the v0.22 bitset kernel. |
| Quiescence and descendants | 56.92% | Encloses several rows; not a separate optimization budget. |

Runtime `madvise` accounts for 8.31% flat and generic copies for 7.31% flat;
their origin/recoverability is not isolated. Eliminating allocations cannot be
assumed to recover all runtime samples. This single profile gives priorities,
not a stable per-position fraction or an elapsed saving. No heavy validation
overlapped it; background system work remained uncontrolled.

## Ranking and cost scale

Recommended impact means expected useful result relative to implementation,
validation, compute and model-supervision costs. It is not the largest
hypothetical Elo gain. Confidence distinguishes a real limitation from confidence
that a new implementation will improve games. Costs below are relative: small
is a focused bounded change, medium crosses engine/backend or protocol controls,
large needs a new representation or data/search design. They are planning
judgments, not measured labor estimates.

| Rank | Coherent project | Likely useful benefit | Confidence of useful result | Engineering / validation / compute / model cost |
| --- | --- | --- | --- | --- |
| 1 | Shared legality and attack context | Moderate exact-throughput opportunity; potential downstream tactical interface | Medium | Medium / high / low diagnostic then normal gates / medium |
| 2 | Compact contextual nonlinear evaluator | Largest broad decision-quality potential; unproven | Low–medium | Large / high / high if diagnosis earns new data / high but capped |
| 3 | Threat-aware tactical frontier | Correct some known horizon omissions without blanket check expansion | Low–medium | Large / very high / medium / high |
| 4 | Draw-aware endgame search semantics | Concrete correctness and conversion consistency | High for fixing semantics; unknown population strength | Medium / high / low / medium |
| 5 | Repetition-safe cache contract and reuse | Potential work reduction at deeper/warm searches | Low–medium | Large / very high / low diagnostic, medium later / high |
| 6 | Compact incremental search state/evaluation | Lower copying and leaf cost; support richer inference | Medium for mechanism, low for size | Large / very high / low diagnostic, medium later / medium–high |
| 7 | Confidence-guided selective depth allocation | Large possible node saving with severe tactical risk | Low | Large / very high / high / high |
| 8 | Cooperative root scheduling with explicit CPU policy | Lower latency on spare cores | Low for strength benefit | Large / very high / medium–high / medium–high |

### 1. Shared legality and attack context

**Mechanism.** Compute useful king/checker/pin information once for a node and
reuse it across existence, preferred-move lookup and staged generation. Ordinary
nonking moves outside pin/check constraints can avoid repeated own-king attack
classification; king moves, EP, promotions and castling retain full special
validation. Preserve opponent-check tags and the exact stable emitted order.
Integrate status with available legal-move evidence without confusing an empty
tactical list with stalemate.

**Why this differs.** v0.23 virtual tags still recompute attacks after each
candidate. The old evasion mask only handles checked nodes. The failed within-
stage lazy study measured unvisited records, not redundant legality reasoning
for visited ordinary moves. The earlier 0.54%-CPU evasion-context helper cannot
justify optimizing that helper again; this proposal targets the broader current
generation/status work.

**Benefit and risk.** Most credible near-term exact-efficiency project, with
possible reusable information for later threat evaluation. No strength magnitude
is established. Pins, discovered checks, king x-rays and stale scratch contexts
are high-risk boundaries; the full legal rules remain authoritative.

**First bounded diagnostic.** One diagnostic on the same frozen 1,145-root tree:
count non-overlapping candidate/tag/status work that a shared context could
avoid, separating checked nodes, king/special moves and context construction.
Check every proposed ordinary eligibility/tag decision against both existing
transition oracles and the 596 independent vectors. Keep diagnostic elapsed
time out of efficiency claims. Budget one study, no context variants or screens.

**Prospective continue/stop rule.** Continue only with zero legal/tag/order
contradictions and a cost model plausibly supporting at least 8% whole-work
saving after context overhead. This is a proposed project-specific threshold,
not a measured gain or universal promotion gate. Stop if benefit is mainly the
already cheap evasion helper or requires weakening legality. A later exact
implementation would need both-module checks, perft/special controls, three
alternating exact warm pairs and all actual-budget gates before games.

### 2. Compact contextual nonlinear evaluation

**Mechanism.** Represent interactions that the additive evaluator cannot express:
king-relative piece configurations and pawn/piece relations, using a small
quantized nonlinear model with an incremental sparse accumulator. Keep current
material/placement as an explicit control. A compact jointly calibrated classical
interaction model is a cheaper alternative if the initial capacity study shows
that a neural representation is unnecessary; do not implement both blindly.

**Why this differs.** It changes representational capacity and inference design,
not coefficients in the rejected 12-feature or 450-coefficient linear models.
Past cleaned splits, sharper labels and target repairs did not establish
generalization. The direct-leaf study does not justify relabeling the old corpus
at scale. Historical holdouts are inspected development material.

**Benefit and dependencies.** Broadest potential strength workstream because
leaf preferences change across the entire tree, including endgames and quiet
defense. Improvement is uncertain and slower inference can erase it. A new
Go-only training/inference pipeline, licensed data/model provenance, robust
quantization and board-update parity would be substantial work. No external
book or online dataset is required by current evidence.

The official [NNUE design documentation](https://official-stockfish.github.io/docs/nnue-pytorch-wiki/docs/nnue.html)
describes sparse accumulators and move updates. That is an implementation
principle, not a compatible Stockfish net or evidence that its size is suitable
here. Stockfish's evaluator, feature layout and surrounding search differ.

**First bounded diagnostic.** Freeze one representation, parameter/memory cap,
and inference-cost envelope. On training-only historical counterexample groups,
test whether its features can separate the failed decisive preferences, with
color/mirror and exact full-versus-incremental inference controls. Use synthetic
weights for the cost test; do not collect new labels or train a large model.

**Continue/stop.** Continue to a separately budgeted data plan only if the chosen
representation removes the declared capacity contradictions and fits an initial
10% whole-search cost envelope in a scratch diagnostic. That budget is a design
choice, not a guarantee that 10% is acceptable for play. Stop if it merely fits
inspected examples, requires an unbounded collection campaign, or loses update
parity. Later success requires family/transposition-separated untouched data,
decision-level and mate controls for both turns, actual equal-time improvement
and a standard screen. Lower supervised error alone cannot pass.

### 3. Threat-aware tactical frontier

**Mechanism.** Give the main/q boundary an explicit tactical contract: detect a
limited class of legal quiet mating/checking threats and search their actual
defensive responses under a frozen local budget. A pass-derived threat is never
a forced loss score. Required evasions stay complete; beta verification must
search evidence capable of refuting the claim instead of accepting stand-pat.

**History and difference.** Independent v0.19 training proofs establish quiet
mate omissions, but v0.20 standalone proof failed actual-budget balance. v0.23
quiet checks and two check-budget designs did not earn promotion. This proposal
differs only if it conditions on a validated threat/response mechanism and
replaces a specified frontier contract. Another first-N-checks layer would be a
retry of the failed family and should stop at design review.

**Benefit/risk/cost.** Could improve tactical reliability that exact speedups
cannot fix. It changes search values and TT/q horizon compatibility; false
threats, missed defenses and tactical explosion are major risks. Shared cheap
attack context from project 1 may help, but is not sufficient. Engineering and
model analysis costs are high; do not combine it with evaluator changes.

**First diagnostic and gate.** Audit a fixed training-only set of existing quiet
mate positives plus legal escapes, blocks, captures, EP, promotions, stalemate
and hypothetical-pass negatives. Before code activation, specify what finite
certificate justifies extra search and count its frequency/cost on neutral
baseline leaves. Continue only with zero false forced-mate classifications and
a projected bounded overhead at or below 10%; stop if the proposal recreates
blanket checking search or cannot state complete defense coverage. Later compare
against a full-width same-horizon reference, then both-budget/both-turn native
controls. Preserving named root moves is not the goal.

### 4. Draw-aware endgame search semantics

**Mechanism.** Align internal terminal/claim evaluation with the game interface:
distinguish claimable draws from automatic termination, preserve checkmate
precedence, and implement a conservative proved subset of dead-position tests.
Use the halfmove clock and complete supplied reversible history. This can avoid
optimizing a material advantage that cannot become a game win.

**History and difference.** The README explicitly documents missing move-count
draw evaluation. This is not the rejected generic passer penalty, endgame king
bonus or conversion heuristic. The benchmark currently claims eligible
threefold/fifty-move draws outside the search, so the engine and referee can
reason about different outcomes. The capped v0.23 game alone does not establish
that this gap caused it or that implementing a draw rule would have won it.

**Benefit/risk/cost.** Clear correctness/semantic value and potentially useful
endgame choices; frequency and strength impact unknown. A claim is an available
action, not necessarily a forced zero in a winning position. Do not declare all
small-material configurations dead. Changing draw semantics also changes score
cache validity and regression expectations.

**First diagnostic and gate.** Reconcile the current game-layer policy and
intended UCI-client behavior, then build a small truth table around halfmove
boundaries, pawn/capture resets, mate precedence, repetition histories and
unambiguous dead-material cases. Establish explicit rule expectations before
changing play. Continue if a policy mismatch is reproduced and can be corrected
without altering non-draw positions; stop broadening the material detector if
its exact claim cannot be proved. No tablebase download is needed. A later
implementation needs updated oracle semantics, GHI/history tests and narrowly
identified changed-value controls before ordinary decision/game gates.

### 5. Repetition-safe transposition/search contract

**Mechanism.** Redesign which search facts can be safely reused, keeping move
hints, exact-depth values, minimum-depth bounds, history obligations and q-horizon
identity explicit. Start with shadow accounting rather than accepting more
scores. A proof-carrying or conservatively history-qualified reuse scheme is a
project; changing `depth == requested` to `>=` is not a correctness argument.

**History and difference.** The architecture study's unsafe board-only scores
actually reused −902cp in a drawable history. The safe split and irreversible
fingerprint already shipped in v0.16. A larger table alone had no useful shallow
node benefit. Strict q reuse later exposed only 5.2–7.3% potential node saving.
This proposal earns a new attempt only if v0.23's current deeper/warm workload
reveals a different substantial reuse boundary.

**Benefit/risk/cost.** Potentially broad search-work benefit, especially as useful
depth grows, but unknown headroom. History contamination, bound misuse and mate
normalization make validation very expensive. No cache-memory increase is assumed.
Projects 3, 4 and 7 would change cache semantics and cannot be bundled into the
diagnostic. Current Stockfish search uses coordinated table, rule-count and
search policies; copying its lookup predicate is not a safe transplant.
[Stockfish search source](https://github.com/official-stockfish/Stockfish/blob/master/src/search.cpp)

**First diagnostic and gate.** One unchanged-tree shadow ledger counts matching
boards with rejected depth/history identities and non-overlapping descendant
work, retaining witness histories for every proposed safe category. Continue
only if a specified sound subset covers at least 10% of current work and all
history/halfmove/mate witnesses agree. Stop on a contamination witness or if
benefit requires arbitrary deeper values in an exact-depth API. Keep unsafe
categories as observations rather than implementing them.

### 6. Compact incremental state and evaluation

**Mechanism.** Retain search-owned state with small undo records or compact
copy-on-ply fields and incremental phase/material/placement totals. Preserve
the current integer rounding exactly, including captures/promotions that change
phase. Evaluate the smaller state route before a full rules/backend replacement.

**History and difference.** Board-copy allocation is already removed in v0.23;
reusing the same scratch cells again is not new. Current copies, transition/hash
work and leaf scans remain. Earlier allocation savings repeatedly overpredicted
elapsed savings, and `runtime.duffcopy` includes more than board copying.
Project 1 reduces generation logic; this project changes state movement and
evaluation maintenance. Their overlapping gains cannot be added in advance.

**Benefit/risk/cost.** Moderate plausible exact-throughput gain and an enabling
path for project 2, with high special-move, unwind, cancellation and ownership
risk. A broad make/unmake rewrite would need full round-trip/hash/perft proofs.
Rounding means a naive single tapered sum may change released scores.

**First diagnostic and gate.** Attribute the existing profile's copy stacks and
count phase-changing versus ordinary updates on the frozen trajectory. Design
one exact incremental-value representation, without replacing board rules.
Continue only if affected non-overlapping work supports at least 10% net saving
after update/undo overhead and exact rounding is feasible. Otherwise retain a
small enabling representation only under an explicitly smaller future objective;
do not turn the current gate into a pass. Actual promotion would require source,
state/hash restoration and full-tree equivalence in both modules.

### 7. Confidence-guided selective depth allocation

**Mechanism.** Allocate depth using a frozen measure of move-order or threat
confidence, protecting forcing continuations and verifying reduced improvements
at full depth. This needs trustworthy rank/uncertainty information, not merely
the late position of a quiet move. It is a changed-search policy.

**History and difference.** Quiet-count LMR already achieved large time savings
and passed all six surrogates, yet failed the external strength assessment.
Material guards and frontier protection also failed different gates. A new
reduction formula or margin sweep is not justified. The old static-pruning
verification remains closed. A materially different confidence signal, possibly
from projects 1–3, is a prerequisite.

**Benefit/risk/cost.** Potentially large tree savings; low present confidence
that they buy strength rather than deeper bad decisions. Data and tactical
validation costs are high. The model-selection process itself is a significant
cost and source of overfitting.

**First diagnostic and gate.** On training-only complete full-depth records,
test whether one declared confidence signal distinguishes costly late branches
from decisive counterexamples across both turns and opening families. Continue
only with meaningful non-overlapping work coverage (proposed minimum 20%) and
zero missed forced-mate counterexamples in the frozen diagnostic. This does not
prove safety; later unreduced regret, fresh controls and external games remain
mandatory. Stop if the signal is merely move count or tuned to the inspected
failed cohort.

### 8. Cooperative root scheduling with explicit CPU policy

**Mechanism.** Explore cooperative bound sharing, cancellation of superseded
work and an explicit thread/CPU contract. Measure total work and process count
alongside latency. Default single-worker behavior should remain separately
testable; live client configuration is outside this project.

**History and difference.** The existing two-worker candidate saved 25.86% time
but required extra CPU, had a failed fresh gate and failed its direct screen.
An unchanged port onto v0.23 or adding a Threads option would not be a new
strength hypothesis. A new study must show that stale-window/duplicate work,
not insufficient cores alone, is the important residual scheduling loss.

**Benefit/risk/cost.** Useful latency on genuinely spare cores, but poor present
evidence for higher strength under the standard three-matchup load. Shared-table
races, nondeterministic ties, cancellation/joins and oversubscription raise
validation cost. This ranks last because an earlier implementation already
tested the general idea with mixed games.

**First diagnostic and gate.** Analyze the preserved parallel records and
opportunity proxy first. Continue only if avoidable stale/duplicate work is
large enough to support a prospective 15% additional latency improvement over
the archived parallel design without increased core budget. If those records
cannot establish it, return an open question rather than start games or assume
the serial proxy is real parallel speed. Any later measurements require isolated
serial and standard concurrent-load controls, aggregate CPU accounting, race
and cancellation tests, and the ordinary decision/screen gates.

## Other options and the recommended decision

Opening books, online move data, tablebases and adaptive time allocation are
legitimate possibilities, but none has a stronger measured immediate case here.
A book mainly affects opening decisions and may conceal repeated-opening
weaknesses. Tablebases add exact restricted-material knowledge with a new data
dependency; project 4 can first fix semantics without them. Adaptive allocation
needs a real-clock instability opportunity study because the routine screen
uses fixed `movetime`; a faster clock check alone already failed its screen.
No general prohibition or engine ceiling is inferred from their present rank.

Before a future engine implementation, resolve the shared-validation omission
and refresh the current checkpoint. Then approve **one bounded project-1
opportunity diagnostic**, with the frozen eligibility, exactness and net-headroom
criteria above. It ranks ahead of a full state rewrite because current generation
work is directly measured and the existing backend can be extended. It ranks
ahead of another evaluator/search attempt because it has a narrower correctness
contract and no new labels, fitting or horizon policy. Project 2 remains the
strongest larger strength workstream if project 1 lacks useful headroom.

The standard later screen remains 50 games each versus tagged v0.23, Skill 3 and
pinned Sacrifice, all three concurrent, after earned correctness/cost/native
gates. Its ten-opening sample is triage, not a reliable small-Elo estimator.
Broader independent openings and a prospectively planned confirmation are
appropriate when stronger strength evidence is required. Do not retry an unchanged
screen or infer calibrated statistical error rates from the practical gate.
Fishtest's paired/sequential methodology is a separate system, not a guarantee
implemented by this runner.
[Fishtest statistical methods](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-Mathematics.html)

All diagnostic thresholds above are prospective proposals for review. None has
been measured or passed by a new candidate. Historical failures, original cost
targets and adverse games remain intact. The next action is review of the audit,
storage pilot and first diagnostic; candidate implementation and publication
remain separate decisions.
