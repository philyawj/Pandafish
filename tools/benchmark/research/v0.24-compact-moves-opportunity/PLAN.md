# Compact search-move representation: bounded opportunity diagnosis

Restore tagged v0.23 playing source first. Preserve both previous failed
legality/attack studies and their original gates. No playing policy or code is
changed by this diagnosis.

The general library Move contains PGN/tree metadata, while the search uses
origin/destination, promotion and tags. Measure its native size against that
minimal data record and sample allocation ownership on the same frozen1,145-root
warm trajectory. The earlier CPU sample's7.31% generic-copy and8.31% madvise
shares are overlapping or unassigned costs, not a promised gain. Do not sum
them into a speed estimate.

One unchanged-baseline allocation sample, using normal sampling and an end-of-run
heap profile rather than per-root MemStats. Verify whole-tree exactness. Read
alloc_space/alloc_objects to distinguish move-buffer ownership from tables,
positions and PGN parsing. Full warm rows remain temporary; use existing saved
exact rows and hashes. Preserve the Go diagnostic, raw profile and compact
reports. No matches, teacher queries or fitting.

Continue only if the compact record is materially smaller and move-buffer work
is a substantial allocation/copy owner in the real search, rather than only
frontend PGN data. This does not establish elapsed headroom: a later exact
implementation still needs its own frozen >=8% three-pair warm gate, all rules/
state checks and actual-budget gates. Stop if most identified cost belongs to
unrelated tables/board copies or if compact metadata cannot preserve semantics.
