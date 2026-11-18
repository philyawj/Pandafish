# Selection: generate later main-search stages only when needed

The fixed 1,145-root depth-five diagnosis creates 132,855,339 main-search Move
records, including 122,425,191 ordinary quiet records. Existing table/tactical/
killer-stage cutoffs occur before ordinary quiet traversal at 3,305,521 nodes
and cover 97,606,034 quiet records (**79.73%**). Table-preferred first cutoffs
alone cover 41,173,522 full records (**30.99%**). Both exceed the preset half-
quiet-record opportunity criterion. These are deferrable-generation counts,
not predicted elapsed savings: preferred moves still need direct validation and
late nodes can require another pseudo-destination pass. All diagnostic results
remain original v0.21; no elapsed claim uses instrumentation.

Implement one staged main-search picker in the **existing exact order**:
validated captured TT/board hint; stable ordered legal captures/promotions;
validated legal ordinary-quiet killer 0/1; remaining ordinary quiet moves in
their original generation order. Skip duplicate signatures at later stages.
Capture the hint and killer priorities before searching, as in the eager path.
Never search killers before captures. Root and all quiescence ordering/eligibility
remain unchanged. Use an additive legal hint lookup and quiet-only append API
inside the pinned extension boundary; do not alter original upstream runtime.

No node, move, priority/tie rule, evaluation, horizon, PVS/TT/history policy or
clock change is included. If a preferred/tactical/killer stage cuts off, later
quiet generation never happens. If it does not, generate remaining quiet moves
without regenerating/tagging tactics. Keep active-ply move storage private and
never overwrite a stage still being iterated. All previous failed clock/history/
cache/proof code stays inactive; baseline is released v0.21.

Controls must compare every yielded whole Move against the eager ordered list
over cold/warm, seeded, special-move, check/evasion, stale/malformed/duplicate
hint and metadata cases. Check legality of every direct hint, quiet subset/cache
ownership, perft, parent/child isolation and zero-buffer-generation after a legal
first hint. Keep independent original-library and frozen move/tag/state vectors,
full race/vet/build and original upstream replay. Require all 325 original
score/move/node comparisons and 1,145-root persistent-TT sequence to stay exact.

Efficiency returns to the larger exact-kernel objective: **15% median elapsed
reduction in three alternating public depth-four pairs**, beneficial deep/newer
work and at least 15% median reduction in three warm depth-five iterative pairs.
The latter uses the actual main-search lifetime and frozen game trajectory.
Freeze drivers before measurement. The same four native safety/quality gates
and one standard three-concurrent 150-game screen follow. Preserve every failure;
no unchanged favorable rerun, larger confirmation, live game or publication is
automatic. Stop before committing a viable candidate.
