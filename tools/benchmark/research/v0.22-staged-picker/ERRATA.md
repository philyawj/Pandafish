# Supplied castling-right compatibility control

Initial standard/seeded controls pass 406,158 lookup/subset cases, 10,240 whole
ordered lists, preferred-only zero-allocation/deferred-buffer controls and perft.
Adding inconsistent-FEN castling rights exposes a compatibility difference: the
pinned generator trusts those rights and can emit moves from an empty king
origin, or duplicate signatures when another piece occupies that origin. The
new direct lookup initially rejects them and signature deduplication changes
the eager list. Preserve the failed `compatibility-initial.txt` and complete
five-file `initial-source.json` / hashes before correction.

This experiment does not fix that pre-existing malformed-position behavior.
Direct lookup preserves the pinned castle-list path before ordinary piece
validation. The picker uses the original eager path when supplied current-side
castling rights lack the king on its starting square, preserving duplicates and
metadata. Standard legal positions retain the staged path. Byte-indexed rights
checks avoid introducing per-node string conversion/uppercase allocation.

Re-run all standard/seeded and added cold/warm inconsistent-FEN controls before
timing. This is an input-compatibility correction in an unreleased prototype;
it does not change original upstream rules or bundle a castling fix.

The first independent warm reviewer assumes both scratch modules name the frozen
helper `archived_proof_test.go`; the baseline actually names it
`warm_profile_test.go`. It stops with a missing-file error after checking the
saved pairs, before producing a review. Preserve `review_warm.initial.go.txt`.
The corrected reviewer checks each actual path against the same frozen helper
hash. No helper, observation, timing or actor query changes.
