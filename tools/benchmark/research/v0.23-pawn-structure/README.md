# Contextual pawn map: reject on actual-budget decisions

The [frozen map](PLAN.md) adds isolated/doubled/pass relations and an immediate
blocker discount to original placement. Material, search, quiescence, TT and
clocks remain v0.22. This is an original hand-designed standalone map, not a
newly calibrated contextual model. Earlier fitted pawn features remain failed.

Mailbox-oracle checks, color/perspective symmetry, parent restoration, material/
terminal controls, normal/race/vet/Linux builds and a clean replay pass. Original
control expectation mistakes and provisional timing overlap remain in
[control notes](CONTROL_NOTES.md). The first kernel costs about 60.68ns with no
allocation and slows the 1,145-root warm workload by a median 13.60%.

Before any quality observation, an [exact compaction](COMPACT_PLAN.md) changes
file/rank loops to bitboard fills and sparse passer traversal. Component cost is
16.54ns, no allocations; three new warm pairs show 3.95% median overhead.
All 6,870 row comparisons retain initial-model moves, scores, depths and nodes.
No coefficient, pawn definition or phase rounding changes.

The first public 100ms suite fails: mean finite loss **45.706→57.881 cp**, above
the original baseline-plus-5cp limit. Black 107.544→133.316, White 21.229→28.021.
Severe cases rise 2→3; no new losing mates or timing overruns occur. All 504
native queries and raw reference records are retained. The driver exits one.
No long/newer suite, fresh seed collection or games run; the prepared seed
2026112063 generator is unexecuted. No favorable coefficient/seed retry follows.

This map changes board-context valuation correctly, but that functional behavior
does not establish better decisions or strength. Its adverse result adds to the
case against immediately guessing another classical weight map. Restore v0.22
and measure a different search opportunity before implementing it.

Final full source SHA256:
`caa66531c365d4653d8a598c4b82ed02b958c06002519a40838d86d20e61dce9`.
Native binary SHA256:
`0740a69e3a49f2c26d954738c51f1bbf30cdf31eddffe058c312f3334cf6ee78`.
The initial kernel, compatible replay manifests, exact comparison and every
cost/native failure remain separate records.
