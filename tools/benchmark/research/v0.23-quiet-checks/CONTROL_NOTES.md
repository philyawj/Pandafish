# Generator and initial search controls

The independent 596-position move/tag set and 2,507 explicit/seeded controls
match complete legal quiet-check lists, including cold/cached calls, order,
caller prefixes and parent state. The seed walk supplies 3,479 quiet checking
occurrences across 3,103 positions. Discovered checks, checking castles, pinned
checks, EP/promotion exclusions and an already checked opposing king are covered.

The initial API median is 349.9ns versus 1461ns for all-quiets-and-filter, a
76.05% component saving. An exact post-castle rook/discovery prefilter avoids
ordinary castle construction. It preserves all controls and measures 325.4ns
versus 1472ns, a 77.89% saving; allocated bytes fall from 92 to 10 per call in
this board mixture. Residual special-case allocations remain; no allocation-
free claim is made. These are generator measurements, not total search savings.

The first complete search test run retains two failures: the old capture-only
exhaustive oracle expects −582 cp on mirrored knight/queen forks, while the new
checking layer obtains +280 cp. The reference's complete-legal-list minimax is
updated to include exactly the newly specified initial quiet checks, without
using the new generator or alpha-beta pruning. Other depths retain the old
capture/promotion policy. The original failed output remains `search-tests.txt`;
subsequent output has a separate name. No production root is forced to the fork
or a teacher move. This is a changed horizon, not a fixed-tree-exact claim.

The first replay is refused because the full snapshot contains unchanged nested
`go.mod`/`go.sum` metadata, outside the shared overlay format. The full snapshot
is preserved. A compatible overlay omits only those two files after bytewise
comparison with v0.22, and clean-overlay tests/build pass. Original failure,
source-format proof and replay output remain separate; no runtime code changed.

The fresh collector's first compilation from `/tmp` is rejected by Go's internal-
package import boundary. The same sources are placed under ignored
`tmp/v023-generator` inside this module before collection. The failed log remains;
no rollout or candidate query occurred before that location correction.
