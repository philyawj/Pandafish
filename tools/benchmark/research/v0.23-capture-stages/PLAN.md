# Defer losing main-search captures without pruning

Restore v0.22 after the exact virtual-tag candidate saves only 4.543% elapsed.
Remaining small exact kernels have repeatedly missed the substantive gain goal;
the full parallel screen also fails. Reassess search work rather than bundle
those inactive changes or tune another inspected evaluation model.

The current main picker puts all captures before quiet killers/moves. MVV-LVA
can give a poisoned high-value target precedence despite a negative exact legal
exchange estimate. Use the released bitset SEE only for ordinary high-attacker/
lower-victim captures outside check; cheap/equal attackers, checks, promotions,
preferred moves and required evasions keep priority. Defer a negative estimate
until after all quiet moves, preserving original order within every stage.
Every deferred move remains eligible: this is ordering, never pruning or a
material score. Quiescence filtering, root order, tables/history and clocks stay
released. Per-ply deferred buffers must survive descendant/quiet-buffer reuse;
an interrupted exchange must propagate incomplete search, not a false score.

Compare complete ordered lists against a separately constructed eager oracle
using the frozen mailbox SEE, cover poisoned/unprotected captures, checking and
promotion exceptions, evasions, hints, cancellation and perft. Keep exact fixed-
depth scores/first root choices, deterministic candidate counts and the full
warm 1,145-root workload. Frozen cost gate: >=15% median elapsed saving over
three alternating pairs; an early >=15% node-saving gate stops later pairs.
No coefficient sweep, selection from a teacher root or retry of failed binaries.

If promising, use the four existing native suites and final seed
2026112116/stream116 excluding all consumed cohorts/live survey boards. The
standard screen remains 50 games each against v0.22, Skill 3 and pinned Sacrifice,
concurrent. No automatic larger confirmation or publication. Preserve adverse
outcomes and check in before committing a candidate that merits promotion.
