## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# v0.15.3: larger-gain architecture investigation

**Recommend a coordinated search-core project first, followed by a coherent,
data-calibrated evaluation/threat model.** The evidence supports a much larger
workstream than another independent small bonus. It does not establish a new
strength gain, identify one guaranteed Elo fix, or justify reverting draw
correctness. [NEXT.md](NEXT.md) gives the implementation order, acceptance
criteria, risks and stop conditions. [CACHE.md](CACHE.md) answers the cache
question in detail; [METHOD.md](METHOD.md) documents all probes and failures.

v0.15.3 is a **research release**. The bishop-pair candidate and every result
are archived, and I chose to keep released playing behavior unchanged. All
architecture prototypes are data outside the normal build. The release retains v0.15.2/v0.15 playing policy; public identity
and documentation change only. No live Lichess client or larger bishop-pair
confirmation starts.

## Why the earlier attempts did not establish the requested leap

The piece-safety experiment detects tested exposure but does not explain its
motivating quiet leaves and costs substantial equal-time work. Bishop pair
corrects a selected exchange and passes the direct practical screen at 57%,
but 44% against Skill 2 and 25% against Sacrifice do not demonstrate a broad
gain. Its direct paired interval is about 44% to 70%. These are failures to
establish the desired strength leap, not proof that every such chess concept
is bad. The original source, games, score uncertainty and negative controls
remain in [the bishop-pair archive](../v0.16-bishop-pair/README.md).

This investigation changes the experimental question from "which plausible
bonus next?" to "where does useful search work go, what information is lost,
and what mechanism could affect a large part of play?" It examines cache
correctness/reuse, CPU/allocation cost, root/search contracts, move generation,
native Sacrifice comparison, equal-time decision quality, selected leaf and
opponent-response errors, and a current released-policy game control.

## Scope and measured findings

The frozen suite contains 36 full-history inputs: 24 systematically selected
positions from older games, three controls, eight live mistakes and the
counterfactual pawn ending. The deep subset is chosen by input order before
its results. These are public, correlated diagnostics, not blind validation or
an Elo population. No prototype is adjusted to produce a named root move.

### 1. Cache restrictions are real; an unsafe rollback is wrong

The v0.15 history fingerprint restricts both score reuse and ordering hints.
Cold depth-four probes make 56,145 lookups with zero hits. An unbounded shadow
map finds 5,378 previously completed board/halfmove contexts, all with different
fingerprints. That identifies lost opportunities, not safe reusable scores.

| Diagnostic policy | Depth-four nodes / baseline | Depth-six nodes / baseline | Correctness/contract status |
| --- | ---: | ---: | --- |
| Released baseline | 100% | 100% | Reference |
| Table disabled | 100% | not probed | Deliberate cache-presence test failures |
| Board-only move hints | 97.3% | not probed | Root/internal and GHI controls pass |
| Reversible-context fingerprint | 97.8% | not probed | Root/internal and GHI controls pass |
| Both history-safe changes | 95.5% | 83.9% | All tested scores/best sets match; root/internal and GHI pass |
| Unsafe board-only scores | 94.3% | 83.4% | **Rejected: actual draw-context contamination** |
| `depth >= requested` alone | 100% | not probed | Changes exact-depth/window contract |
| 32 MiB table, 16x capacity | 100% | not probed | Exceeds intentional memory budget; no shallow node benefit |

Depth four covers 36 inputs; depth six covers six prespecified inputs, so the
columns are different cohorts. Node reductions are search-work measurements,
not strength or identical-cost throughput measurements. Warm and native timed
records complement these cold probes.

The unsafe cache reuses a -902 cp score in a history where a fresh search finds
a draw, despite the drawing continuation being legal. Original and safe
prototypes return zero. It is insufficient to refuse to cache zero scores:
the problematic cached score is nonzero. Preserve the guard's purpose while
redesigning identity and irrelevant-history handling. See [CACHE.md](CACHE.md).

The original table introduction, exact-depth reuse and v0.15 history correction
are separate changes. The older PVS/LMR attempt predates the latter. There is
no support for blaming its failure on a guard introduced afterward.

### 2. Generating and ordering moves is the largest measured CPU opportunity

A clean CPU profile of the released timed search kernel gives these cumulative
shares. Rows overlap and must not be added:

| Work | Sampled CPU share |
| --- | ---: |
| Quiescence, including its descendants | 63.2% |
| `orderedMoves`, including generation and sorting | 50.9% |
| Library `CalcMoves`, across all callers | 40.9% |
| Stable move sorting | 13.6% |
| Deadline/cancellation checks | 8.1% |
| Static evaluation | 6.3% |

The kernel rotates 24 neutral positions under normal 100 ms management. It
allocates about **72.8 MB and 120,241 objects per search** in that workload.
A separate allocation sample attributes about 41% of bytes directly to
`standardMoves`, 38% to `ValidMoves` defensive copies and 9% to castle-list
creation. These are approximate sampled ownership figures. They are not an
additional CPU percentage or a universal per-move allocation rate.

In the depth-four node audit, **589,377 of 645,606 nodes (91.3%) are quiescence**.
It generates 6,608,863 legal moves to examine 452,690 eligible tactical/evasion
moves, about 6.8%. Even the tactical-sort prototype still generates all legal
moves before filtering; it mainly removes irrelevant sorting/copy work. It
produces about 15% more reported nodes at equal requested time in the 100 ms
suite, with only a modest completed-depth increase. It is a lower-bound probe
of a deeper search-API opportunity, not a demonstrated 15% strength gain.

The library already uses bitboards and a mailbox. Its general-purpose API,
legality/check-tag simulation and defensive-copy semantics are costly for this
caller pattern. A search-oriented legal-capture/evasion interface, staged move
picker and reusable storage should be evaluated before a full internal-board
rewrite. A full rewrite has much higher special-rule correctness risk.

Raw profiles and reports include [clean CPU](native-cpu-clean.pb.gz),
[cumulative CPU report](cpu-clean-cumulative.txt),
[allocation ownership](alloc_space.txt) and [kernel benchmark](cpu-clean-profile.txt).
The first joint profile is preserved but not used for normal CPU percentages:
its dense allocation-stack sampling consumed substantial CPU itself.

### 3. Search contracts and interactions offer more than an isolated patch

Returning one proven best root move instead of certifying every exact tie saves
16.7% of nodes at depth four across the full suite, but only 2.6% at depth six
in the separate deep subset. PVS alone is not uniformly cheaper. A root move's
mere agreement with Stockfish is not the metric; all tested fixed-depth scores
and selected-move membership are checked.

The exploratory coordinated prototype combines the explicit one-best contract,
PVS, both safe-cache changes and tactical-only ordering:

| Depth | Inputs | Baseline nodes | Coordinated nodes | Reduction |
| --- | ---: | ---: | ---: | ---: |
| 3 | 36 | 122,595 | 112,474 | 8.3% |
| 4 | 36 | 645,606 | 523,194 | 19.0% |
| 5 | 6 | 485,501 | 307,732 | 36.6% |
| 6 | 6 | 2,405,637 | 1,517,750 | 36.9% |

All these searches complete legally with identical scores; every selected move
belongs to the released best set. The coordinated normal tests still flag the
changed all-ties API, specifically the requirement to return both equally long
forced-loss resistance moves. The selected move and mate-distance score remain
valid. That contract needs explicit design/test review before production use.
The prototype is not silently called correctness-complete because it is fast.

In its contemporaneous 100 ms block, neutral-case mean completed depth rises
from 3.67 to 4.08. A finite-node reference quality mean improves in that block,
but the suite is small, correlated and used to explore these designs. There
is no architecture-prototype game-strength screen. This is enough to justify
hardening and testing a larger search-core candidate, not enough to promote it.

### 4. Sacrifice comparison: features are not isolated explanations

The pinned Sacrifice code combines PVS, reductions, history ordering,
aspiration, depth-qualified TT bounds, killers and a richer evaluation of
mobility, shelter, activity, center and attacks. Tal game data affects root
ordering rather than positional evaluation. Its author background does not
identify a performance mechanism. The code uses the Rust chess library; it is
not a bare material-only reference.

Native cold probes give mean reported depths of 5.63 versus 4.01 at 100 ms and
7.83 versus 5.67 at 1,000 ms, Sacrifice versus Pandafish. Its reported node count
is about four times as high at 100 ms, **but this is not four times comparable
throughput**: node definitions differ, it counts depth-zero main calls again
when entering quiescence, reductions change tree meaning, and evaluation and
memory budgets differ. Nominal depth does not establish better moves either.

On the 21 neutral cases where both scores remain cp, the referee's mean
clamped move loss is 56.0 cp for Pandafish versus 78.2 for Sacrifice at 100 ms;
at 1,000 ms it is 40.9 versus 22.8. Three other neutral cases involve mate
scores and stay separate. These positions come from older Pandafish games and
can favor its familiar distribution. The mixed results prevent presenting
search depth, NPS or this small position suite as an engine rating.

Native warm requests use real UCI generations. Of 25 parent/next episodes, two
end in mate; their expected null terminal moves are not errors. Among 23
nonterminal warm/cold pairs, safe cache reaches higher depth in four and equal
depth in 19, with no lower cases. Baseline is higher in one/equal in 22;
Sacrifice is higher in two/equal in 20/lower in one. One warm pass is diagnostic,
not a stable strength estimate.

### 5. Wrong leaves and missed opponent responses are distinct problems

The 16 selected reconstructed leaf checks use full histories and the last
completed exact score from fresh 100,000-node reference searches. Illustrations:

| Diagnostic alternative | Engine quiet-leaf Black score | Referee Black score | Interpretation |
| --- | ---: | ---: | --- |
| Nd4 | -12 cp | -5 cp | This preferred-line leaf is reasonably assessed; investigate stronger opponent deviations/horizon. |
| Qd6 | +161 cp | -662 cp | Large leaf-model/threat disagreement. |
| Qe8 | +75 cp | -227 cp | Also underestimated disadvantage, but less severe than Qd6; root ranking is reversed by the leaf model. |
| Rd5, engine's Rc8 response line | +14 cp | +16 cp | Preferred-line agreement misses the actual stronger Rxd5/exd5 ending. |
| Kxg5 continuation | +400 cp | mate -1 | Non-checked leaf misses a quiet mating threat. |

The quiet mate is also independently verified by both the Go rules library
and python-chess in [quiet-mate-proof.json](quiet-mate-proof.json), rather than
being accepted solely from Stockfish. Capture quiescence returns -400 cp for
White there, while a one-ply full root search finds mate. This is a reconstructed
line illustrating a boundary failure, not a claim that the opponent executed
the mate in the live game.

The pawn-ending counterfactual separately remains strongly favorable for White;
splitting the doubled pawns does not remove the advantage. A generic passer
bonus could favor the wrong side. This calls for response-branch training and
king-blockade/pawn-break understanding, not another guessed pawn penalty.
The previous naive quiet-check extension timed out, so the mate example does
not justify unbounded additional checking search.

[leaf-referee.json](leaf-referee.json) preserves all 16 comparisons, including
counterexamples and disagreements. An individual Stockfish score is finite
and may change with budget. Main-PV leaf agreement is not proof that the
opponent model is right. A new corpus must include stronger response branches.

## Retrospective strength control, not a prototype match

One standard concurrent screen gives context for the released playing policy:

| Actual tested engine / opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.15.2 / tagged v0.14.4 | 21 / 14 / 15 | 0 | 0 | 56% |
| Tagged v0.15.2 / Stockfish Skill 2 | 21 / 2 / 27 | 0 | 0 | 44% |
| Tagged v0.15.2 / pinned Sacrifice | 7 / 13 / 30 | 0 | 0 | 27% |

Fifty games per opponent, 25 color pairs, 100 ms requested moves, 50 ms protocol
grace, 400-ply cap, the standard ten-position book, three concurrent matchups.
No other engine probes deliberately overlap the match. The machine reports AC
power; background work and temperature are not controlled. The actual old-tag
revision and binary hash are documented in [control-provenance.json](control-provenance.json),
including the mistaken descriptive revision label supplied to the command.
The tagged executable itself is correct. All original metadata are preserved.

[Complete settings, hashes and games](../../results/v0.15.2-architecture-control-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json)
and [all PGNs](../../results/v0.15.2-architecture-control-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/games.pgn)
remain separately labeled from the bishop-pair screen and release smoke.

The practical previous-release gate passes, but this is a small retrospective
screen. v0.15 changes mate-distance and repetition behavior as well as cache
context, so the 56% cannot be assigned to any one of them. It does not establish
a cache-induced strength regression or prove the guard has no cost.

The control's 44%/27% external scores resemble the bishop candidate's 44%/25%
in this profile. That weakens the earlier impression that the bishop bonus
itself explains the external shortfall compared with older historical runs.
Neither the two-point Sacrifice difference nor subtracting other old screens
establishes a regression. Opponents, repeated openings, random ties, workload
and finite samples matter. The original 57% bishop direct lead remains
inconclusive; I chose to archive the feature.

## Recommendation and what ships

1. **Harden the coordinated search-core prototype as the next playing
   experiment**, with explicit one-best semantics, safe cache identities and
   tactical move staging. Its deeper work savings make this a larger concrete
   candidate than a new standalone evaluation bonus. No prototype Elo claim.
2. **Remove the full-legal-generation/allocation bottleneck in tactical search.**
   Extend the library's search API or introduce a tested internal search
   position only if needed. Preserve all special moves, stalemate, checks and
   draw invariants. This is the largest measured CPU/allocation project.
3. **Develop a coherent data-calibrated evaluation/threat model**, using stronger
   opponent deviations as well as principal lines. Compare a compact tapered
   classical model with a compact trained/NNUE route after measuring inference
   and data costs. Neither a larger weight nor naive quiet checks is justified.

[NEXT.md](NEXT.md) details the stages, tests and stop conditions. Safe cache
repair is part of that architecture, not a rollback of repetition correctness
or a promise of a large standalone gain. The deadline margin, cache size,
author's age, language choice and one feature's presence do not establish the
strength explanation. The old PVS/LMR bundle is not automatically reintroduced.

The release archives every prototype, source overlay, unit failure, position
probe, profile and game result. It restores the tagged released evaluator;
material, placement, search, quiescence, TT policy, repetition and clocks are
unchanged. Only public release identity/greetings and documentation change.
No NNUE, new backend, cache prototype or bishop-pair bonus ships. The optional
bishop-pair 200-game confirmation is deferred, not silently run or made mandatory.

## Reproduce and inspect

- [METHOD.md](METHOD.md): input selection, counters, profiling caveats, native
  protocols, quality scoring, intentional versus substantive failures.
- [CACHE.md](CACHE.md): source history, lost reuse, safe designs and actual GHI
  contamination.
- [ERRATA.md](ERRATA.md): incorrect Bf3 prose (actual Qf3), main-PV reasoning,
  harness/profile mistakes and revision-label correction.
- [analysis.json](analysis.json), [analyse.py](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz): all fixed/deep,
  timing and per-case reference tables; raw data remain alongside them.
- [sources.json](sources.json): primary references and the inspected pinned
  Go/Rust library/source identities.
- [provenance.json](provenance.json), [release-checks.json](release-checks.json):
  source/binary identities, unchanged-policy checks and separate release validation.

The [reproduction driver](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz) creates new scratch trees and outputs,
without overwriting this archive or running matches or live games. See its
`--help` and the method notes. Source overlays apply to tagged v0.15.2; the
prototype code is data, not compiled by normal repository builds.

Primary context: the [GHI paper](https://cdn.aaai.org/AAAI/2004/AAAI04-102.pdf)
explains why history-sensitive score reuse requires care;
[pinned Stockfish search](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp)
and [pinned Sacrifice search](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs)
show coordinated search systems rather than isolated gain estimates.
[Fishtest mathematics](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-Mathematics.html)
and its [FAQ](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-FAQ.html)
provide statistical context; this runner does not implement their sequential
statistical guarantees. No numerical feature-Elo guesses are added or summed.

## Research-release validation

[release-checks.json](release-checks.json) records passing formatting, full
repository race tests, vet, build, explicit darwin/arm64 build, offline public-
config wrapper and archived-source replay. The three-way integration smoke
completes all six games without forfeits or unfinished games: 1/1/0 versus
v0.15.2, 1/0/1 versus Skill 2 and 0/1/1 versus Sacrifice. These two-game samples
are integration evidence only; no strength claim is made for the release.
The archived smoke summary preserves its actual development-build hash.
