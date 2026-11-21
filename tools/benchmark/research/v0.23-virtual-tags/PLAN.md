# Exact ordinary-move tags without temporary Board simulation

After the parallel screen and root-scout cost failures, restore tagged v0.22.
Current profiles retain material move-tag cost; the successful v0.22 exchange
revision shows that computing the same predicate from compact bitsets can be
more effective than reducing allocation or adding conservative geometric gates.
Earlier full-ray/occupied-ray precheck attempts stay failed and inactive.

For ordinary moves only, compute resulting occupancy and attacker masks directly
from original piece sets plus moved/captured-square deltas. Preserve the exact
own-king legality and opponent-check predicates and cached king-square behavior.
Avoid a temporary full Board copy, mailbox writes and convenience recomputation.
Promotions, en passant and castling retain the original full simulation path.
No generation order/eligibility, child-state, evaluation, search, table or clock
change. The original library tagger and frozen independent vectors remain oracles.

Require pseudo-legal tag comparisons, both existence/full modes, special/invalid
FEN controls, original independent move/state/perft vectors, parent immutability,
normal/race checks and exact warm moves/scores/nodes. Three alternating depth-five
pairs require at least 15% median elapsed saving; no node-saving gate applies to
an exact-tree implementation. If cost fails, archive without native/games.

If cost passes, collect all four existing native suites. Use final seed
2026112115/stream115, excluding quiet-check/phase/root-parallel cohorts plus old
controls and surveyed live boards. A failed surrogate remains failed; an exact
efficiency mechanism may receive one separately documented research screen only
after matched-depth review of all changed choices, with no legality/timing/new
losing-mate issue. This possible exception is declared before native outcomes.
Keep the standard 50 games/opponent, three concurrent, direct >50% criterion and
full review of both external rows. No automatic larger confirmation or publication.
