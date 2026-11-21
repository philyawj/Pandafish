# Scope reassessment before remaining controls or games

The original plan's all-six-gates/no-surrogate-exception condition is **not
satisfied**. Fresh 100ms remains failed on one new losing-mate report,
proof-fresh-12-ply36. The node/time/full-depth review successes do not pass that
gate. Original false fields, the stopped workflow and all raw records remain.
This document explicitly revises the evidence-gathering scope before any game;
it is not a claim that the original protocol passed or promotion is justified.

All four existing suites pass. Fresh 100ms finite mean improves 50.3307→48.4788cp,
both turns improve, severe errors stay three and no timing/legality failure
occurs. Nevertheless candidate Rfd1 receives mate -5 at the independently fixed
one-million-node confirmation, while baseline Rh1 remains finite (-854cp).
Unrestricted reference remains finite (-873cp). The finite alternatives are
not proved saveable, and the surviving candidate mate report is not dismissed.
The confirmation queries the unrestricted root and the complete union of both
actors' moves; no timed actor run or old gate is repeated or overwritten.

Independent native matched-depth diagnosis gives:

| Depth | Tagged v0.22 | Candidate |
| --- | --- | --- |
| 5 | Rh1, -18cp | Rh1, -18cp |
| 6 | Rfd1, -250cp | Rfd1, -231cp |

The bad choice is also preferred by unreduced search at depth six. This rules
out attributing this observed choice solely to a reduced branch; it does not
eliminate its equal-time cost or prove all pruning decisions safe. The forcing
reference line contains further quiet checks beyond the capture-only horizon,
a known search omission that additional nominal depth does not fix monotonically.
No particular production root move or FEN exception is added.

The mechanism has stronger evidence than a small timing variation: 45.448%
median elapsed saving and approximately 50% fewer nodes over three pairs, all
1,145 results legal/deterministic, original independent mate/race controls
passing, and full-depth selected-move regret averaging 1.014cp over 1,103 finite
roots, worst 273cp, with zero mate-risk flags/value contradictions in that review.
These are selective-search diagnostics, not equivalent full-width depth or Elo.

That combination justifies completing the **previously unexecuted fresh 1000ms
suite once**, with unchanged inputs/binary/settings and only matching completed
reference data reused. It must pass its original gate; no further exception is
planned if it fails. If it passes, run **one standard research screen**, 50 games
each against tagged v0.22, Skill 3 and pinned Sacrifice, all three concurrent.
Keep the unchanged direct >50% conservative criterion and zero-forfeit review.
Preserve both external rows and every raw game, then legally audit all games.

The exception is a post-result research judgment and adds selection uncertainty.
A favorable screen cannot erase the failed surrogate or establish a precise Elo
gain. Review the confirmed weakness, all failed predecessors and resources with
the complete game outcomes before any recommendation. No second unchanged
screen, larger confirmation, live game, commit or publication follows automatically.
