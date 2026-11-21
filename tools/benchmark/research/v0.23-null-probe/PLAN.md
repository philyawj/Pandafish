# Full-opponent-ply null probe with deeper verification

Select a substantive pre-live search direction after six rejected prototypes
and the failed capture-futility opportunity diagnosis. The older null policy's
failed screen stays failed. At its shallowest trigger it could stop immediately
in capture-only quiescence without a full opponent move. This is a concrete
horizon limitation, not a proven explanation of its old match losses.

The current branch's verified-pass attempt searches two full plies and verifies
every apparent cutoff. It saves only 9.89% nodes and no elapsed time. Its failure
motivates testing a different explicitly selective cost/risk boundary, not calling
the failed verification design successful or weakening its recorded gates.

Freeze this new policy on released v0.22:

- Main nodes at depth at least three, unchecked, with finite non-mate beta and
  static score at least beta+50cp. Own non-pawn material at least 900cp, at least
  one own rook/queen, total phase at least eight. Skip current-board repetition
  count above one, halfmove clock >=80, either side's advanced pawns within two
  promotion steps, malformed kings and an already attacked opposing king.
- Hypothetical pass depth is max(1, depth-3), so **every pass probe searches a
  full opponent ply**. Use released capture-only quiescence, with all existing
  check evasions/promotions. No checking-generator code from the failed normal
  quiet-check attempt is active. Skip a terminal pass state.
- The pass worker has no real repetition prefix, shared TT/hints/killers or
  nested null pruning. Share cancellation/deadline; add every visited probe node
  to the real worker. A mate-valued result cannot cause a cutoff.
- At nominal depth six or more, also require real-move verification at depth
  minus two, isolated from shared tables/pruning but carrying exact cloned real
  repetition state without double-entering the current board. Count its nodes.
- Return beta only on a completed non-mate fail-high. At depths three through
  five there is deliberately no real verification: this adds tactical/zugzwang
  risk. The material/pawn/repetition guards and opponent-ply coverage reduce
  specific risks without proving full-width equivalence. No precise mate score
  or hypothetical pass-history score is stored in the real table.

Material/evaluation, ordinary main ordering/root eligibility, quiescence policy,
TT identity/capacity/lifetime and clocks remain unchanged. Nominal depth and
heuristic bounds must not be described as an exact full-width proof.

Require focused both-color quiet-mate refutations, guards, cancellation/draw
priority, parent/history/table restoration, full normal/race tests and a clean
source replay. Frozen warm criteria remain >=20% fewer actual nodes and >=15%
median elapsed saving. Shared `warm-plan` stops immediately after an early
deterministic node-gate failure; no extra timing pairs can rescue it.

If it survives, run all four original actual-budget gates, then fixed seed
2026112083/stream83 final controls, excluding public/newer, both consumed v0.23
control pools and every v0.22 live survey board. Only then run the standard
concurrent 50 games per opponent. No tuning after those observations, unchanged
favorable rerun, automatic larger confirmation or publication.

Reassessed alternatives: a genuinely new contextual/learned representation with
independent data/inference budgets is the next larger project if this fails;
another hand-weighted activity/pawn map is not. Remaining exact state/attack work
needs measured warm headroom beyond the already rejected 9.14% reuse result.
The ranking remains conditional after each result, not a fixed checklist.
