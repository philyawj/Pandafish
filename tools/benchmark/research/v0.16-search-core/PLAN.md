# Coordinated search-core candidate: frozen scope and gates

I am testing Stage A from the v0.15.3 architecture plan. The hypothesis is that
coordinating the root contract, PVS, safe cache reuse and tactical sorting buys
more useful search depth at equal time without changing the evaluator or legal
search tree. One completed root move must belong to the independently proven
best set, with exact score and mate distance. Enumerating all tied root moves
is no longer a production promise; full-window alternative scoring remains a
diagnostic.

I keep material, static evaluation, capture eligibility/SEE, quiescence horizons,
killers, score-table size/replacement/exact-depth policy, time allocation and UCI
lifecycle unchanged. The separate board-hint table adds at most 1 MiB. Only its
ordering crosses history contexts; scores retain draw-sensitive identity.
Fingerprint resets occur on actual pawn/capture child descent, never arbitrary
entry points, and unwind on both completion and interruption. Full legal move
generation remains; specialized generation/reusable search state is deferred.

Before any strength screen:

1. Run existing rule, mate, cancellation, exhaustive minimax and cache tests;
   explicitly test the one-best contract against full-window alternatives.
2. Test real warm-history contamination in both directions, irreversible reset
   and unwind, EP/castling distinctions, hints/collisions, tactical eligibility,
   special promotions, quiet evasions and stalemate. Keep existing exact-depth
   score/window tests. Preserve failures and fixes.
3. Compare tagged v0.15.3 and candidate fixed-depth scores/best-set membership on
   the old 36 diagnostics plus a prespecified broader set selected without
   candidate results. Include deeper probes, actual warm successive UCI requests,
   and equal-time 100/1000 ms blocks, serial and three-search load.
4. Require full race tests, vet, build and offline wrapper validation. Record
   binary/source hashes, full commands and all raw outputs.

If correctness and equal-time evidence are promising, I will run exactly 50
games each against tagged v0.15.3, Stockfish Skill 2 and pinned Sacrifice, with
all three matchups concurrent. I will preserve all outcomes, including a failed
direct gate, and distinguish functional/work improvements from measured
strength. No automatic larger confirmation, live games, commit, push or PR.
I will reassess the ranking after reviewing this result.
