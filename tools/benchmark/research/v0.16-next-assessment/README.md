# Fresh candidate choice after v0.15.2

**Select a single, low-cost bishop-pair interaction for the next experiment.**
This is a fresh selection after rejecting the attacked/undefended-piece term,
not progression through the old five-item list. No new gain or Elo estimate is
assigned. The implementation and validation live in
[the bishop-pair record](../v0.16-bishop-pair/README.md).

The v0.15.2 research release is published at
[GitHub](https://github.com/philyawj/Pandafish/releases/tag/v0.15.2), through
[PR #29](https://github.com/philyawj/Pandafish/pull/29), merge
`64e60edb2f7c210e942c9643453237647c4867b3`. It preserves playing behavior.
The next baseline is built from that tag, not from experimental source.

## Evidence that changed the choice

The rejected term passed its unit controls but corrected neither motivating
quiet leaf. It cost 18.2% of nodes in balanced 100 ms probes and had mixed or
negative saved-decision effects. Functional detection is not enough; the next
feature should have a demonstrable positional interaction and a much smaller
per-leaf cost. Reweighting that term is not the next experiment.

The first fresh diagnosis still examined the proposed pawn ending, rather than
assuming it was the automatic second item. [diagnosis.json](diagnosis.json)
records seven new 500,000-node reference boards, with full history where
applicable. After `Rd5 Rxd5 exd5`, equal material has Black's central passer
and White's kingside majority. White's reference advantage remains substantial
with either side to move. A generic passed-pawn bonus risks reinforcing Black's
bad exchange. PNfQoYP0 is still a queen/rook/minor-piece position; it is not
direct evidence for a pure pawn-ending formula.

A second frozen counterfactual set asks whether Black's doubled g-pawns explain
that ending. [pawn-variants.json](pawn-variants.json) retains every exact
iteration, excluding incomplete bound updates. With equal material, the last
exact white scores are +386 cp in the original ending and +386, +209, +296,
+230 after four ways of splitting the g-pawns. Moving the central passer to an
extra a-pawn still gives +322. Finite-node differences are not proofs, but they
contradict the simple story that one doubled-pawn penalty explains the result.
The continuation includes king access, blockading and pawn breaks. A compact
correct model has not yet been established. No pawn feature or weight was
implemented or fitted to Rd5.

Next, an independently constructed equal-material exchange fixture tests the
missing distinction between B+B+N and B+N+N. In
`r4bk1/pp5p/2n2n2/8/3N4/5BB1/PP5P/R5K1 w - - 0 1`, the released engine gives
up its two-color bishop pair with Bxc6 at depths three and five. The knight
capture Nxc6 preserves both bishops. Fresh 500,000-node restricted reference
checks score Bxc6 +50 and Nxc6 +304 cp; the unrestricted reference also selects
other moves at shallow depth, including Ne6. We do not treat Nxc6 as the only
acceptable move or attribute the entire reference difference to bishop pair.
[bishop-diagnosis.json](bishop-diagnosis.json) preserves the moves, scores and
continuations. The fixture also involves piece placement and rook activity.

A fixed 25 cp pair bonus, chosen once without tuning the public holdout,
changes unrestricted engine searches to a pair-preserving alternative at
depths three and five and at 100 ms. The property is explicit and testable:
access to both bishop square colors is worth a small interaction score while
material values remain identical. It can be computed inside the existing
64-square traversal, without attack maps or new search work.

## Reassessed priorities

| Priority | Candidate | Assessment after the result |
| --- | --- | --- |
| 1 | Bishop pair alone | Concrete controlled exchange, clear bounded interaction, almost no extra board work. Screen it; a general strength gain is unknown. Do not bundle rook-file scoring. |
| 2 | Pawn-ending diagnosis | Concrete serious failure, but simple passer/doubled stories do not explain the counterfactuals. Diagnose blockade, king access and breaks before implementing. |
| 3 | Phase-sensitive placement | Cheap and plausible, but needs a stronger controlled decision example before changing broad incentives. |
| 4 | King exposure | Real cases, but the earlier king-danger screen and unstable home-rank probe require a new mechanism, not another generic penalty. |
| 5 | Revised piece safety | Some positive cases survive, but the tested narrow term failed mechanism/cost checks. A future revision needs specific quiet threats and lower cost first. |

This ranking can change again after the bishop-pair result. Rook-file scoring
remains separate and unselected. A failed result will be preserved rather than
silently stacked into the next evaluator.

## Validation plan fixed before screening

Keep material, placement, search, quiescence, table policy and clocks unchanged.
Test incomplete/same-color promoted pairs, both colors/reflections, multiple
bishops, equal material, immediate tactics and queen sacrifices. Preserve all
40 earlier public regression boards and the same 15/15 split; add the synthetic
exchange family with color/file mirrors. Do not tune weights on those results.

Measure unrestricted fixed-depth best sets, balanced repeated 100 ms choices
and static cost. If the mechanism survives, run the standard 50 games each
against tagged v0.15.2, Skill 2 and pinned Sacrifice, three matchups concurrent.
Preserve every failed gate and game. No automatic larger confirmation or live
games. Stop before committing, pushing or preparing the next PR for maintainer
review. The v0.15.2 publication authorization does not publish this candidate.

## Reassess again after the bishop-pair screen

The bishop-pair screen passes the practical previous-release gate at 57%, but
44% against Skill 2 and 25% against Sacrifice do not establish broad strength.
All 150 games finish without forfeits. The direct paired interval is roughly
44% to 70%. The [complete result](../v0.16-bishop-pair/README.md) therefore
recommends holding this source for a separately approved, preplanned 200-game
confirmation per opponent. No larger sample or next PR is started.

The immediate priority is resolving this candidate's evidence while keeping
its weight frozen, rather than stacking another evaluation term. If choosing
a subsequent mechanism after that decision, the diagnostic queue changes:

| Priority | Next work | Why |
| --- | --- | --- |
| 1 | Confirm or retire the frozen bishop-pair candidate | Concrete functional correction and 57% direct lead justify confirmation, while uncertainty and external scores prevent promotion now. |
| 2 | Trace king-exposure leaves | Kf7 is a serious new negative control for preserving a pair. Explain the searched leaf and defensive alternatives before another king term. |
| 3 | Diagnose pawn-ending blockade and breaks | The original ending remains bad and simple doubled/passer counterfactuals are insufficient. No ready compact model yet. |
| 4 | Phase-sensitive placement | Cheap, but still needs a focused causal contrast before implementation. |
| 5 | Revisit piece safety or rook files only with new evidence | Do not revive the rejected attack-map cost or add another absent feature simply because it is on a list. |

This is a decision queue, not an instruction to implement five features in
order. No new implementation follows the bishop-pair result automatically.


## Subsequent maintainer decision

The maintainer chose archive-only treatment in v0.15.3 and a broader investigation
of larger opportunities before a bishop-pair confirmation. The original choices
and results above remain historical evidence. Read the
[new measured architecture plan](../v0.15.3-architecture/NEXT.md) for the current
priority; released playing behavior remains unchanged.
