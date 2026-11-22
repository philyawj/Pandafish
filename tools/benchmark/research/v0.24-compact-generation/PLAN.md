# Compact candidate-generation kernel

Base playing source: v0.23. Earlier full pin-context and target-attack attempts
remain failed under their original8% objectives; source, controls and all warm
observations are preserved. No earlier candidate component is activated.

The sizeof/allocation diagnosis supports a narrow internal representation change:
use a6-byte search candidate through generation/tagging, constructing the full
128-byte public Move only for emitted moves or special-move fallback. Existing
cached Move values retain all metadata and order. Lookup/status/tactical/quiet
paths share that candidate representation. Precheck castling rights and occupied
paths before calling the exact existing castle rule routine, which otherwise
allocates a two-Move array even when no castle is possible.

Keep original upstream runtime files, public APIs, emitted tags/order, transition
hashes, evaluation, search/q/TT/repetition/mate policy and clocks unchanged.
No new attack algorithm, pin context, fitted term or pruning is included.

Correctness requires both-module normal/race/vet, focused full cached metadata,
malformed-rights, specials, zero-allocation preferred lookup and generator/oracle
controls, plus596 independent vectors. Three alternating warm pairs must retain
all1,145 moves/scores/nodes with median time ratio<=0.92 (>=8% saving). Failure
stops before native/fresh/game work; no unchanged favorable rerun or relaxed gate.

Only earned candidates receive four existing and two untouched fresh actual-
budget suites, then the standard concurrent50 games each vs v0.23, Skill3 and
pinned Sacrifice. Preserve every failure. No playing-candidate commit/publication
or live play is authorized. Use shared Go plans and evidence conventions.
