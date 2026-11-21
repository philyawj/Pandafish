# Control contract and retained failures

The first run fails the old every-capture-before-quiet order assertion, as
expected for this declared ordering experiment. The reference is updated to
construct the full eager list first and partition losing captures with the
separate frozen mailbox exchange oracle. It does not call the new live picker
or bitset classifier. All moves and original within-stage order remain required.
Checking captures, promotions, hints and required evasions keep their priority.

The first storage implementation also causes a real regression: the existing
preferred-only zero-allocation control reports one allocation. Direct writes
through the picker to worker-owned deferred storage make the compiler retain
the worker on the heap. Isolating append/reset in a searcher helper restores
zero allocations on the unchanged control and allocates deferred storage only
when a losing capture is actually deferred. The initial failing logs and escape
analysis remain; the control is not weakened. An intermediate unused test import
compile failure is also retained.

The corrected normal suite passes. Poisoned/unprotected/checking-capture controls,
preferred-move override, cancellation propagation, perft and whole-list comparisons
cover the changed contract. No root is forced to a reference-engine move.
