# One-ply late-child scout reductions

The separately archived unchanged-tree study passes: eligible outer scouts cover
66.436% of 43,118,858 nodes, 65.957–66.832% by turn. The old narrow-parent
guarded-LMR study remains rejected; this policy includes late child scouts of
wide-window parents and therefore carries a different selective risk.

Start from uninstrumented tagged v0.22. On the fourth/later main-search move,
parent depth >=3 and ply >=1, use one less ply only for an ordinary nonchecking
quiet child's existing one-centipawn PVS scout. Protect checked parents, killer
hints, captures/promotions and mate-band windows. Parent window width does not
exclude an individually narrow child scout. Do not reduce the first three moves.
If a reduced probe improves alpha **or reports a mate-band value**, repeat its
null-window probe at the original full depth. Then retain the original full-
window PVS re-search for a full-depth in-window improvement. No move is omitted.

This changes selective search and fixed-depth scores may change. A reduced fail-
low can overlook a useful quiet move; full-depth verification of improvements
does not eliminate that risk. Preserve actual history, terminal/mate rules,
quiescence, material/evaluation, score-key/depth/replacement policy and clocks.
Test-only unreduced references use fresh/nil tables; production never mixes
algorithms in one table. No parallel, root-scout, capture-stage or clock code.

Require focused guard, full-depth verification, mate, draw/history, legality,
context restoration and cancellation controls plus normal/race tests. Use all
1,145 warm roots for three alternating depth-five pairs. Frozen targets: >=20%
node saving after each pair and >=15% median search elapsed saving. Values are
not required to match the unreduced baseline; candidate results/counts must be
deterministic and every result legal and complete. Keep all adverse observations.

A surviving candidate must pass all four existing native suites and new both-
turn final controls, seed 2026112118/stream118, excluding all consumed cohorts
and surveyed live boards. No exact-efficiency surrogate exception applies to
this selective policy. Independently review its selected-move full-depth regret
before any game screen; accept co-best alternatives, never force a root move.
Then one standard 50-game-per-opponent screen, all three concurrent, with the
unchanged direct criterion and complete external review. No automatic larger
confirmation or publication; stop before commit for a viable recommendation.
