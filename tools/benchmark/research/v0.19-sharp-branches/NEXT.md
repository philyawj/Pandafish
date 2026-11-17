## Publication note for v0.19

The original local experiment report below retains its stop/check-in state.
These attempts were subsequently approved for publication together in v0.19.
Playing policy remains v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# My reassessment after the sharp-branch retry

I reject the selected fit and keep released v0.18. Retaining difficult labels,
balancing turns and fixing the offline target direction did not produce a
generalizing candidate. I will not keep trying this same bounded 24-weight fit
on an inspected holdout. These priorities supersede the previous ordering;
I will reassess after each result rather than complete a fixed checklist.

1. **Test a bounded exact quiet-mate proof at the quiescence boundary.** My new
   training scan finds 15 static records with legal quiet mates in one, verified
   independently by python-chess and unchanged tagged Go code. This extends the
   older single reconstructed quiet-mate example with fresh training evidence.
   I would first profile a nonrecursive proof kernel: generate candidate quiet
   checks and prove checkmate by legal-response existence. I would return a mate
   score only for an actual legal mate, preserving mate distance and history.
   I would not recursively search arbitrary quiet checks or turn hypothetical
   opponent-pass threats into loss scores. The earlier recursive approach's
   timeout is a reason for a strict bounded-cost design.

   I would predeclare positive controls from these proofs and fresh discovery
   data, plus negative checks that have a capture, block or king escape, pinned
   pseudo-checks, stalemate, castling/en-passant/promotion edge cases, repetition
   and depth/deadline caps. I would require independent legal proof, ordinary
   position overhead and deadline tests, equal-time actual-choice controls and
   fresh both-turn holdout gates. A naive full legal list at every q-node may
   cost more than it helps; I would measure and reject that implementation if
   needed. Only a promising, correct, affordable candidate would proceed to the
   standard three concurrent 50-game matchups. No larger confirmation follows
   automatically. This is a proposal, not a tested improvement.
2. **Reprofile legal tagging/check/status reuse on released v0.18.** I retain
   this as the leading separate engineering project. v0.17/18 already produced
   measured efficiency gains; the local backend offers a bounded extension
   point for sharing repeated work. I need a current profile before selecting
   one change. A quiet-check proof kernel may expose useful primitives, but I
   would keep an exact-tree speed experiment separate from changed-search play.
3. **Replace the evaluation representation/selection design before another
   learned retry.** The bounded feature space cannot express some training
   preferences, and the earlier frozen-leaf pass did not guarantee actual-search
   safety. I would compare contextual king-threat features or a small nonlinear
   representation using representative training data, with inference and data
   budgets declared first. I would select on actual-budget validation searches,
   then use a fresh untouched holdout. A neural model is not an automatic fix;
   increasing caps against these known failures is not a new hypothesis.
4. **Revisit selective search or make/unmake only with a new measured case.**
   Previous broad search changes failed, and a board rewrite adds large undo,
   legality and hash obligations. I have no new evidence that either should
   displace the concrete quiet-mate omission and profile-driven backend work.

I do not need an opening book or Lichess dataset for the top proof experiment:
the mechanism and its counterexamples are reproducible offline. A future broad
learner may need a more representative source/distribution than these deliberately
disjoint rollouts. The current results cannot settle that data question or rule
out other fundamental weaknesses. I have not established a precise Elo ceiling
or promised a large gain from any one item.

I stop this authorized retry at its failed gate and preserve its result. I will
check in before a new PR, commit or push, and before treating a different research
project as the next active implementation.
