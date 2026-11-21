# Guarded pass probe with real-move verification

Reassess a substantive pre-live search thread after four rejected attempts.
The old v0.14.3 null experiment remains failed at 56% direct / 36% Skill 2 /
23% Sacrifice. Its minimum-depth-three, depth-minus-three pass could enter
capture-only quiescence without even one full opponent move. That is a design
limitation, not a claim that it has been proven to cause those old games.

This is a new selective policy on v0.22, not an exact-kernel optimization:

- Eligible nodes have depth at least four, are not checked, have finite non-mate
  beta, static score at least beta, at least 500cp own non-pawn material and
  total phase at least eight. Skip advanced pawns within two promotion steps,
  near-draw halfmove clocks, malformed kings and a prechecked opposing king.
- Probe the hypothetical pass with **two full opponent/reply plies**, followed
  by the already validated single checking layer inside this probe only. Normal
  quiescence remains released v0.22. The probe has no real-game repetition
  prefix, score table, hints or nested null pruning. Skip a terminal pass state.
- A non-mate fail-high also requires a **real-position verification at depth
  minus two**, without null pruning or shared score/hint tables, but with an
  exact copy of actual repetition counts/context. Avoid double-entering the
  already pushed real node. Both searches share the parent's cancellation and
  deadline. Count all their work in reported nodes.
- Return beta only after both non-mate fail-highs, never a mate/exact score from
  the hypothetical pass. No bound from either subsearch enters the real table.
  Existing real draw checks, cache identities and mate normalization stay intact.

Verification and material guards reduce risk but do not prove full-width minimax
equivalence or eliminate every zugzwang. Document nominal/selective depth and
heuristic bounds honestly. Do not change material/evaluation, root restrictions,
ordinary move ordering, clocks, table capacity or draw rules.

First require guards, both-color quiet-mate refutations, real-history and parent
restoration, cancellation, no nested pruning/cache pollution and full/race/rule
controls. On frozen warm depth-five work, require at least 20% fewer actual
visited nodes including probes and at least 15% median elapsed saving in three
alternating unprofiled pairs. Scores/moves may differ by design; every result
must be complete and legal. A cost failure is archived before native/games.

A surviving policy receives all four original actual-choice suites with their
unchanged gates, then a new fixed 16-family final control (seed 2026112053,
blocking all consumed v0.23 control roots), then one standard concurrent 50-game
screen per opponent. No numerical strength claim follows from saved nodes.
Preserve every failure and stop before committing a candidate worth promoting.
