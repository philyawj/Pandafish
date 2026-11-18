# Search transition specialization plan

Baseline is released v0.19. The preceding quiet-proof candidate failed its
fresh-turn gate and is archived separately. Its runtime/search policy is absent.

Hypothesis: the immutable library transition needlessly scans and rewrites all
twelve bitboards for each move, both while tagging candidate legality and while
creating a searched child. A search-only transition that updates affected
piece sets directly should remove meaningful cost without changing any value,
move eligibility/order, hash, history, table, quiescence or time policy.

Scope: additive board/position helpers; use them in the authored search generator
and three engine child-update sites. Original pinned runtime methods and license
stay byte-identical. Retain immutable copies, original attack functions and
complete check tagging. Do not bundle quiet-mate search, make/unmake, attack-ray
reuse, score changes or selective pruning.

Correctness: compare generated legal moves/tags and immutable children against
the unchanged original methods over special-move fixtures and seeded legal
traversals. Check FEN, all piece bitboards/mailbox, hashes, check/status, cached
and uncached operation, source immutability, repetition and perft. Run root and
nested-module tests, race checks, vet and native/ARM64 builds. All frozen public
84-position depth-three/four and deep-subset depth-six probes must retain exact
scores, chosen moves and node counts. Add newly generated roots as regression
controls, without claiming those inspected development cases as a blind holdout.

Performance: after correctness, retain initial timings and run three alternating
baseline/candidate timing pairs at depth four. Require at least 10% lower total
elapsed time on their median ratios, without a contradictory deep-suite result.
Keep all outliers, rather than rerunning to obtain a preferred ratio. Then run
native equal-time choice/legality/deadline controls at 100ms and 1000ms using the
same declared finite-reference and survival criteria; fresh-turn tolerances
remain required. An exact-tree gain is functional efficiency, not strength.

If those gates pass, run the standard 50 games each versus tagged v0.19,
Stockfish Skill 2 and pinned Sacrifice, all three matchups concurrent. Preserve
every result; use the established practical previous-release criterion and
review both external matchups. No automatic larger confirmation, live games,
commit, push or PR. Check in before committing a candidate worth presenting.
