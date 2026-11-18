# Review before the evasion cost prototype

The original opportunity criterion fails: 77,273,086 of 203,289,758 counted tag
validations are avoidable (**38.01%**, below 50%). All 1,145 completed roots,
moves, scores and nodes agree; no legal candidate is excluded by the diagnostic.
Preserve this failed criterion in `evasion-opportunity-review.json`.

The count treats an existence-only validation and a full two-king classification
as one equal unit. That does not model their actual work. Among full tags,
75,752,084 of 146,644,450 are avoidable (**51.66%**); existence checks already
find king escapes early and offer little additional opportunity. The frozen
staged profile attributes 16.26% cumulative sampled CPU to tagging and 24.88%
to generation. These overlap and include setup/measurement runtime overhead;
they are not percentages that can be added or a reliable speedup prediction.

This post-diagnostic review authorizes **one local implementation/cost probe**
using cheap bitboard checker geometry and precomputed interposition masks.
It does not claim the original equal-count opportunity passed. The original
10% extra warm elapsed gate over the rejected staged picker and 15% total gate
over v0.21 remain unchanged. A failed cost probe is archived before native or
game work. This is a transparent revision of the opportunity model, not permission
to buy another screen with a trivial edit or a favorable rerun.

Retain the functionally exact staged picker as the base of this genuinely new
generation revision, with the old 47% screen explicitly rejected. Only the
evasion prefilter is added. Use original upstream/full legal methods, full ordered
Move/metadata controls, independent frozen vectors, special checks/EP/promotions,
missing/nonunique kings, ownership and race checks. All retained moves still
undergo full legality. No heuristic, score, cache, horizon or clock change is
included. Any native exception requires a separate explicit review before games.
