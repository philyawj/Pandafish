# Precomputed stable ordering experiment

This implements the [profile selection](../v0.21-next-assessment/SELECTION.md)
under the [frozen assessment protocol](../v0.21-next-assessment/PLAN.md).
Playing baseline: tagged v0.20. Research/replay base: merged tooling commit
`6bfb4d565c2adf2a2ad59216a462b0f4f0c01503`, whose engine policy is v0.20 and
benchmark default is Skill 3. No candidate commit or release is authorized.

Calculate the existing priority formula once per move. Stable-sort 16-byte
original-index/priority records, then apply their permutation with one temporary
128-byte Move per cycle. Already descending keys need no permutation. Per-ply
scratch storage remains private to the searcher. The priority formula, stable
ties, move eligibility/order, values, library state, search, TT and clocks remain
unchanged. Original ordering helpers remain comparison oracles.

Require the declared exact-tree controls and at least 15% median depth-four
elapsed reduction before native choice validation. A failed mechanism/cost or
native gate is archived, not promoted. A surviving candidate reaches 50 games
each versus tagged v0.20, Skill 3 and pinned Sacrifice, all three concurrent.
Retain every result and preserve runtime/source/binary labels. Reassess after
each result; no favorable unchanged retry or automatic larger confirmation.
