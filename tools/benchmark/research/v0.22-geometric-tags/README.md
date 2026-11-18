# Rejected full-ray geometry attempt

This version retains exact move tags and all 325 fixed-depth score/move/node
comparisons, but fails the unchanged 15% performance gate. Three alternating
depth-four ratios are 0.99322, 0.98892 and 0.97859 (median **1.1% less time**).
Depth-six is 1.01397 (**1.4% slower**); newer depth-four is 0.96647. The timing
driver exits 1 and preserves `timing_gate_passed: false`.

The [selection](../v0.22-next-assessment/SELECTION_TAGS.md) describes the mechanism.
Full geometric king rays classify too many origins/destinations as potentially
relevant, retaining simulation work while adding context/branch overhead. This
is an inference from the conservative predicates and timings; no per-branch
causal percentage is claimed. Full race/vet/build, 100,650 pseudo-legal tag
comparisons, geometry/quiet/special controls, the independent 596-position oracle
and original upstream root-package race tests all pass. Original library methods
and all 21 pinned runtime files stay intact. No native decision probes, warm
candidate workload or strength games run after the failed timing gate.

The initial lookup-initialization defect caught before tests is preserved in
[ERRATA.md](ERRATA.md) and `search_tags.initial.go.txt`. `candidate-source.json`
contains the corrected three-file Go overlay actually tested, based on v0.21;
`source-hashes.json` records exact bytes. Restore through `tools/research replay`
in a new scratch directory. Drivers with unused warm/profile modes remain
source only, not evidence of an executed check.

Keep this version inactive. A separate revision will use current occupancy to
identify the first king-ray blockers and visible checking destinations, rather
than every aligned square. It must pass the same gate; no failed result is
relabelled or folded into another candidate's measurement.
