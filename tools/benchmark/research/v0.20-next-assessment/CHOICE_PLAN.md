# Quiet-proof actual-choice validation plan

This plan precedes new native/reference quality results. Production root choices
remain unrestricted. Baseline: tagged v0.19. Candidate: the frozen quiet-proof
source and named binary. The teacher is pinned Stockfish, full strength,
Threads 1, Hash 64 MiB, fresh hash per distinct 100,000-node query. Existing
matching public reference records may be reused with source/binary hashes.

Public controls are the existing 84 positions. Each actor gets three 100ms
trials and one 1000ms trial under three-search load. Actor order alternates;
reference work starts only after each timed native block finishes. Every raw
UCI line, telemetry record, selected move and failure is retained. Terminal
null moves are controls, not errors. Native response must stay within allowance
plus the unchanged 50ms protocol grace.

For fully finite cases, the position loss is the trial mean of
`clamp(unrestricted_reference_CP - selected_move_CP - 25, 0, 1000)`. Close
alternatives are acceptable. Mate cases remain separate; a finite score may
still represent a lost position. At both budgets each suite must have mean loss
no greater than baseline plus max(5cp,10%), and severe positions (>300 loss)
no greater than baseline plus two. A newly selected losing-mate report fails
when the unrestricted reference is finite/nonlosing and no baseline trial has
a losing-mate report. This is a finite-budget survival diagnostic, not a proof
that the original position was salvageable.

After public controls pass, one fresh batch will use 16 offline rollout
families, seed 2026111728, full-strength 5,000-node MultiPV8 exploration for the
first eight plies, then tagged v0.19 at 25ms versus 5,000-node Stockfish with
student color alternating by family. Sample paired plies 12/13,24/25,36/37,
48/49,60/61 through ply 64 or game end. Require at least 80 live roots and
30 per root turn. Reject sampled states matching prior blocked positions,
any reconstructed sharp-batch witness-path state or another sample's color/rank
mirror; ignore counters and unusable en passant. Retain every attempted family
and excluded sample; do not extend the batch to improve coverage or quality.
Freeze every history/input before candidate probing. These
are new development controls, not statistically independent strength games or
a fitted-model holdout. No old weights are fitted or selected. Timed engine
choices can vary; preserve exact observed queries rather than rerolling them.
The same aggregate gates apply, plus each fresh root turn's mean tolerance.

Only after all gates pass will the standard 50 games per opponent run against
tagged v0.19, Skill 2 and pinned Sacrifice, all three matchups concurrent.
Retain every game regardless of outcome. No automatic larger confirmation,
live game, commit, push or PR follows. A failed candidate is archived and the
next mechanism reassessed under the current authorization.
