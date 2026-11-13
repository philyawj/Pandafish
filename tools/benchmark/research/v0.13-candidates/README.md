# Choosing the next strength experiment after v0.12

Reviewed November 12, 2026, after releasing v0.12. The original review preceded engine changes. The horizon diagnostic below
describes the released baseline; the selected implementation and its separate
audit are recorded in [the local experiment](selective-quiescence-audit/README.md).

## What the last two releases actually showed

v0.11 addressed check evasions at the quiescence boundary. v0.12 reduced work
through quiet cutoff hints. Both passed the previous-release screen, but their
external scores are nearly unchanged:

| Development candidate | Previous-release score | Stockfish Skill 1 | Pinned Sacrifice |
| --- | --- | --- | --- |
| v0.11-dev | 55.5% | 59.5% | 18% |
| v0.12-dev | 58.5% | 57.5% | 17.5% |

These are separate 100-game screens with randomized choices, repeated openings,
and differing power histories. They do not establish an Elo difference between
versions. The v0.12 node audit improved four depth-four cases by 35.8% in total,
but all six saved 100-ms probes still completed the same depth. Less work did
not automatically produce another fully completed ply at those positions.

The earlier recommendations favored small, easily isolated engineering steps.
That was useful groundwork, but it was not enough evidence to expect a large
strength jump. Having room to improve is different from knowing which change
will exploit it. The next experiment should target a substantial tactical or
search limitation and show its intended effect before another release claim.

## Reference-engine comparison

The current default-branch commits checked through GitHub were unchanged from
the earlier source comparison: Sacrifice b528f0ba, Sunfish e396948a, and Stockfish
def601cc. Source inspection identifies differences, not their individual Elo
contributions. Sacrifice's claimed Stockfish Level 5 result was not reproduced
under our benchmark settings.

| Engine | Tactical leaf handling | Main search | Evaluation |
| --- | --- | --- | --- |
| Pandafish v0.12 | Two ordinary capture/promotion plies; extra checked-node evasions; final static safety cutoff | Iterative alpha-beta, equal-depth TT reuse, tactical and killer ordering; no PVS, LMR, or null-move pruning | Material and small placement formulas; phase-blended king placement; no explicit mobility, pawn structure, or pawn-shield score |
| [Sacrifice](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs) | Continues capture/promotion quiescence, with all evasions in check and a global 128-ply safety boundary | PVS, aspiration windows, LMR with full-depth re-search, killer/history ordering, deeper TT reuse | [Classical placement, mobility, king shelter, center, and attack terms](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/evaluate.rs); Tal PGNs affect ordering rather than static evaluation |
| [Sunfish](https://github.com/thomasahle/sunfish/blob/e396948acbd80cdcbad3dccfb07e878fba9df8f4/sunfish.py) | Selective tactical continuation with move-value thresholds and futility limits, rather than our ordinary two-ply cutoff | Repeated narrow bound searches, cached moves/scores, guarded null searches and selective reductions | Piece-square scores with incremental updates and differing king preferences by material conditions |
| [Stockfish](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp) | Capture/evasion quiescence with SEE and futility-based selectivity, draw checks, and a global ply safeguard | PV search, extensive ordering/history, LMR, verified null searches, and other selective pruning | [NNUE evaluation](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/evaluate.cpp) and scaling |

The practical lesson is broader tactical stabilization and selective tree
coverage, backed by useful evaluation. Another ordering hint alone does not
reproduce those systems. Sunfish's king-capture representation and bound-search
driver are not drop-in replacements for Pandafish's legal-move alpha-beta.
Sacrifice is the clearest intermediate classical example. Its combined feature
set cannot tell us which single addition accounts for its strength.

## A fresh diagnostic of the current tactical boundary

An isolated v0.12 source copy counted ten cold-table depth-four searches:
start position, three existing timeout positions, and the six confirmed stream
positions. The production engine was not modified. Results:

- 262,003 total visited nodes; 237,847 were quiescence nodes (90.8%).
- 39,062 non-checked leaves stopped at the normal two-ply quiescence cap.
- 35,552 of those leaves still had a legal capture or promotion (91.0%).
- No ultimate checked-node safety cutoff occurred in this small sample.

Remaining tactical moves are not necessarily good moves. This is not a claim
that 91% of leaf scores are wrong, nor a representative sample of all games.
It shows how often the cap can stop before exchanges settle, and why blindly
increasing it could be expensive: quiescence already dominates node visits.

A separate legal, non-check exchange fixture makes the limitation concrete:

```text
4k3/8/3q4/3p4/2B5/3R3K/8/8 w - - 0 1
```

After Rxd5, Black's Qxd5 can be answered by Bxd5. No check is involved at the
ordinary cutoff. The two-ply quiescence score is -213 cp; a four-ply diagnostic
scores -85 cp because the queen recapture does not refute taking the pawn.
The difference is 128 cp. This is a leaf-search example, not a claim that the
full engine cannot see the sequence from every shallower root or that four
plies solve quiescence generally. v0.11's checked-node extension does not
address this non-check sequence.

[Raw diagnostic](horizon-audit.json) and [counter patch/test source](horizon-audit-source.json)
are preserved. To reproduce, unpack tag v0.12 into a temporary checkout, apply
the saved search-counter patch, write its test source to
horizon_investigation_test.go, adjust the explicit fixture/output paths, and
run `go test -run '^TestHorizonInvestigation$' -count=1 .`. The test verifies
legal moves and that the relevant root/cutoff positions are not checked.
It adds counters and invokes the existing helper at two diagnostic depths;
it does not implement a new engine policy or run online games.

## Ranked candidates

Candidate 1 was selected by the maintainer for local implementation and
validation. Candidates 2 and 3 remain future alternatives. The original review
and diagnostic below describe the released v0.12 baseline.

### 1. Selective quiescence that resolves ordinary exchanges

**Preferred next experiment.** Replace the ordinary two-ply cutoff with bounded,
selective tactical continuation. Use static exchange evaluation (SEE) to judge
a capture's likely material outcome on its destination square and control
which captures deserve deeper work. Retain all required evasions and handle
checks, promotions, en passant, pinned attackers, and king legality carefully.
Keep a final safety bound and existing cancellation/deadline checks.

This changes the quality of the leaf scores across a large portion of the
search, rather than merely changing how quickly the same shallow leaves are
reached. Sacrifice's continuing exchanges and Stockfish's selective quiescence
are useful references. The verified non-check fixture gives a concrete failure
to fix. This is substantially broader than v0.11's checked-node exception.

Risk: unselective continuation can explode the tree and reduce main-search
depth; an incorrect SEE filter can discard tactical sacrifices. Do not simply
raise the cap or copy a negative-SEE cutoff blindly. Compare bounded broader
continuation alone against the same continuation with safe capture filtering
as internal ablations, keeping the evaluator and main-search policy fixed.
Test more unseen exchange positions, not only this fixture, and measure total
work, completed depth, response tails, and match scores. Do not release every
intermediate parameter adjustment as a separate strength milestone.

### 2. Selective main search: PVS followed by conservative LMR

**Largest search-depth upside, with more tactical risk.** PVS probes later
moves with narrow windows and re-searches when needed; LMR initially searches
late, low-priority quiet moves less deeply and restores depth for promising
results. Sacrifice implements this combination. Stockfish uses much richer
versions, and Sunfish also uses selective reductions in a different driver.

This can change the effective branching factor more than a constant-factor
ordering improvement. Begin with PVS at internal nodes, keeping root tie
handling exact. Verify its scores against the unreduced reference. Then test a
small LMR rule separately, excluding root moves, early moves, checks/evasions,
captures, promotions, and strong hints. A reduced fail-high needs full-depth
verification. Null-move pruning is a later alternative, not another change to
bundle into the same initial experiment.

LMR deliberately changes search coverage and can miss a quiet tactical move
whose reduced search fails low. Full-depth minimax parity is no longer a general
guarantee. Use explicit enabled/disabled comparisons and tactical regressions;
a higher reported depth alone is not proof of stronger play. If the goal is
maximum upside and more risk is acceptable, this is the alternative to choose
before candidate 3. Stabilized tactical leaves make this easier to assess.

### 3. A stronger classical positional evaluator

**The move-quality alternative.** Build on our small placement formulas with a
coherent board-aware evaluation experiment: phase-aware piece-square scoring,
both sides' mobility, pawn structure, and actual pawn shelter around the king.
Use feature ablations rather than adding arbitrary bonuses together. Start
with a specific observed weakness, such as missing pawn shelter, and validate
against broader quiet-position judgments as well as tactical safety.

Sacrifice provides a readable classical example; Sunfish shows that tuned
piece-square preferences and incremental scoring can stay compact. Stockfish's
NNUE shows how different the mature evaluation system is, but does not establish
that Pandafish should adopt neural evaluation now. Our current king shelter
preference depends on file/rank, not whether protective pawns actually exist.

Risk: plausible handcrafted terms can conflict, reward premature attacks, or
overfit the stream cases. Use original code/weights and a broader held-out
position corpus. Our saved Lichess games can supply examples; larger database
sampling or coefficient fitting fits later, rather than being required before
this increment. A database opening book would measure a different benefit.

## How to aim for a meaningful gain

Choose one workstream before changing code. Candidate 1 is the best-supported
next target here; candidate 2 is the higher-risk choice for a larger search-depth
jump. No source comparison or diagnostic supports promising a numerical Elo
gain. The evidence identifies opportunities, not an inevitable breakthrough.

Require a demonstrated intended effect on multiple unseen positions and useful
fixed-time behavior before spending the routine matches. Then use the normal
100-game previous-release gate followed by serial Skill 1 and pinned Sacrifice.
Compare with v0.12. Predeclare what would count as a meaningful external result,
and keep negative or flat results instead of rerunning to obtain a favorable
score. A small head-to-head lead supports a modest claim; flat external scores
should remain described as flat. A specific unresolved decision can justify a
larger confirmation, rather than making longer runs automatic.

The next release should be evaluated as a substantial tactical/search experiment,
not assumed to be a big strength jump because its mechanism is standard. Keep
the minimal Go alpha-beta architecture and validated clocks. Another live
session can follow meaningful local progress; it is not needed to discover or
reproduce the exchange-cutoff limitation above.

Before reviewing match scores, we chose a visible external improvement of at least
10 percentage points over the v0.12 screen on one reference opponent, without
a decline greater than 5 points on the other, as a useful descriptive target.
This is not a statistical test or an additional automatic release gate. The
normal previous-version gate still applies; all scores will be reported even
if the external result is flat or worse. No automatic reruns are planned.
