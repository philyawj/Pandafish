# Contextual pawn structure, separate from placement

Prospective sixth attempt, conditional on the verified-pass cost gate failing.
The older architecture review explicitly recommends pawn structure and promotion
paths; this does not originate from one live mistake. Two guessed activity/king
representations and the checking screen have failed on this branch. A full
learned contextual evaluator needs a separate corpus/inference project. Before
that larger investment, test a compact, original pawn-only relation model with
one frozen coefficient map and inexpensive bitboard operations. No fit or sweep.

Keep material values, existing placement, search, normal quiescence, table and
clocks at released v0.22. Add White-minus-Black tapered pawn relations:

- Extra pawns on a file: minus 12 middlegame / 20 endgame per surplus pawn.
- A pawn with no friendly pawn on either adjacent file: minus 10 / 16.
- A passed pawn has no opposing pawn ahead on its own or neighboring files.
  Bonuses by relative rank 0..7 are 0/0/5/12/25/45/75/0 middlegame and
  0/0/10/25/50/90/150/0 endgame. A directly occupied promotion-path square
  halves its bonus. This is a positional incentive, not a forced promotion.

Use file/fill masks and rank populations, no generated legal moves, attack maps,
heap state, pawn cache, root move requirements or term fitted to saved outcomes.
The passed-pawn blocker is immediate, not a proof of an unobstructed whole path.
Avoid combining rook-file, bishop-pair, king safety or mobility changes.

Pawn counts/structure appeared in the rejected older linear feature models;
the chess concepts are not novel and those failures remain relevant adverse
evidence. This standalone frozen map is a bounded engineering probe of their
cost/decision consequences on the current kernel, not a claim to complete the
older recommendation for independent-data calibration or nonlinear learning.
If it fails, another hand-weighted pawn map is not the automatic next attempt.

First require scalar-oracle comparison on seeded legal positions, both-color
mirrors/perspective, isolated/connected/doubled/pass/blocker controls, unchanged
material and mate/draw handling, restoration and full/race tests. Measure
component cost and frozen warm depth-five elapsed; reject a median slowdown
above 15% in three alternating unprofiled pairs before native tests. Changed
nodes/scores are expected and not evidence of strength.

Then require all four existing actual-budget suites with unchanged quality,
both-turn and safety gates. Use untouched seed 2026112063, 16 fixed families,
excluding all earlier v0.23 control roots, for both final budgets. Keep every
failure; do not retune or choose another seed after inspection. A surviving
candidate gets one standard concurrent 50-game-per-opponent screen. No larger
confirmation or publication follows automatically; check in before commit.
