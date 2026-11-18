# Reassessment after the rank-order screen

Recommend promoting the [rank-order candidate](../v0.21-rank-order/RESULTS.md)
under the established practical gate, pending the requested review before a
commit. Exact measured trees, 18.9% median fixed-work time reduction, safe native
controls and the clean 55%/50%/71% screen support that recommendation. The 55%
direct lead remains an uncertain small-sample result. No precise Elo gain or
automatic Skill 4 advancement is supported.

The initial plan prioritized remaining legal-state work. Fresh profiles showed
stable ordering consuming about 21% CPU and made precomputed small sort records
a stronger first choice. That choice clears its gates on its first measured
implementation, so no further feature is bundled to chase a larger game score.
The next hypothesis should be selected after this candidate's review rather
than worked through as a predetermined checklist.

1. **Warm-game allocation/copy profile before choosing the next larger kernel.**
   Immutable child boards/positions account for substantial selected allocation,
   while the fixed-depth driver also allocates a fresh TT per root. Production
   reuses that table within a game. Measure a warm realistic workload before
   estimating the opportunity. Coallocated immutable child state or caller-owned
   per-ply destinations may reduce allocations with a smaller correctness burden
   than make/unmake, but ownership/history/cancellation controls and a meaningful
   cost target must precede any claim. Neither is implemented here.
2. **Remaining attack/tag work or exact exchange computation.** Candidate CPU
   sampling still exposes temporary move tagging and the mailbox exchange
   estimator. Choose one measured mechanism. Preserve full legality, pins,
   special moves and deterministic least-attacker ties, with independent/original
   comparison oracles. Do not bundle different SEE thresholds, quiet checks or
   pruning into an exact implementation change.
3. **A different evaluation representation remains separate research.** Prior
   fitted values failed safe generalization/actual-budget gates. A contextual or
   nonlinear Go pipeline needs a new representation, runtime/data budget,
   actual-search selection and unused final families. The easier reference
   result and Sacrifice sample do not justify guessed weights or reactivation of
   an archived quiet-mate proof. Online datasets/books remain optional future
   data choices; this implementation needs neither.

The ladder's three timeouts stay a separate unresolved reliability finding.
This candidate's 1,252 native requests and 150-game screen have no response
overruns or forfeits. That establishes these observations, not an absolute future
scheduling guarantee or a diagnosis of the earlier outliers. Preserve those
failures; investigate them deliberately if another reliability study is needed.

The optional broader-opening confirmation is justified before a large/general
strength claim, but has not run and is not a new mandatory release gate. No
additional engine experiment, commit, push, PR, release or live game follows
automatically from this checkpoint.
