# Candidate assessment after v0.19

The released v0.19 engine retains v0.18 playing policy. No rejected fitted
weights remain active. Three later learning attempts failed either held-out
surrogate selection or actual equal-time safety. Another adjustment to those
24 coefficients has a weaker basis than a directly verified tactical omission.

| Rank | Candidate | Evidence and main risk |
| --- | --- | --- |
| 1 | Nonrecursive exact quiet-mate proof | Fifteen training witness records contain verified legal quiet mates that capture-only quiescence omits. Direct legal proofs avoid fitted score guesses. Per-node overhead could outweigh the benefit. |
| 2 | Purpose-specific legal/check tagging | Earlier backend changes delivered exact-tree speedups. Existence checks currently compute unused opponent-check tags, and attack tests repeat ray calculations. A fresh profile must establish remaining cost. |
| 3 | Different evaluation representation and actual-budget validation | Bounded feature capacity and generalization remain inadequate. This requires a new Go data/model design, rather than more inspected-holdout tuning. |
| 4 | Isolated selective search or make/unmake | Prior broad pruning failed; no new evidence presently outweighs the correctness burden of a board-state rewrite. |

The first implementation will test rank 1. Candidate quiet checks will be
filtered with board geometry, then validated and proved by legal-reply
existence. No recursive quiet-check expansion or opponent-pass loss score is
allowed. Proof runs only outside check, above the existing depth cap, after a
valid stand-pat cutoff has failed. Material, evaluation, table and clock policies
remain unchanged. Mate distance, repetition and cancellation must remain correct.

Initial acceptance requires all saved positive proofs, negative/legal-state
controls, agreement with full legal enumeration over seeded legal positions,
and fixed-work overhead no greater than 15% on the frozen 84-position public
suite. Larger overhead rejects this implementation before expensive strength
testing. This cost threshold is declared before its measurements. A surviving
candidate still requires native equal-time decisions, fresh controls and the
standard three concurrent 50-game matchups. Functional proof coverage is not
measured strength. No larger confirmation runs automatically.

If the proof implementation fails its cost/mechanism gate, its exact Go source,
tests and every result will be archived and playing source restored. The next
choice will be reassessed using the result and a current profile, with rank 2
as the leading separate alternative. Failed experiments do not require another
approval cycle. No commit, push, PR or live Lichess game is authorized here.
