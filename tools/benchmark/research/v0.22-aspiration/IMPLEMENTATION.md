# One aspiration probe per eligible iteration

The prior completed iteration supplies an ordering/value estimate, not a proof
at the next depth. At depths three and above, a finite previous score defines a
±50cp root window. The search still reaches the full requested depth and retains
all legal moves. Alpha-beta can stop disproven alternatives earlier within that
window. A result strictly inside the window is exact; a result on either edge
or outside is a bound and receives one original full-window retry.

The implementation adds `aspiration.go` and a window-capable root helper in
`iterative.go`. The existing `root` wrapper passes infinite bounds. With those
bounds, root child windows, strict best replacement and tie handling are the
original calculation: alpha tracks the incumbent and beta never cuts a legal
finite/mate score. Required direct-root comparisons check that score/move/node
behavior, separately from the intended iterative node changes.

Depths one/two and previously completed mate-band scores bypass the probe. A
finite-to-mate transition can fail the narrow window and retry. This can cost
extra work, especially on unstable tactical roots; no faster completion is
guaranteed at every position. There is no widening schedule, parameter fitting,
new aspiration-specific TT policy, guessed root preference or selective pruning.

Only an exact completed pass updates `searchResult`. If either phase stops,
the prior completed iteration remains the result, or the existing legal random
fallback remains if none completed. Both phases use the same worker's context,
deadline, TT, killer hints and repetition bookkeeping. Each root call pushes and
unwinds its history, and a retry begins from the restored root context. Partial
subtree bounds remain valid under the existing table rules. Time budgets and
cancellation check sites inside search remain unchanged; retries consume the
same remaining allowance.

Evaluation, main depth, PVS, quiescence eligibility/horizons/SEE, mate distance,
table capacity/replacement/draw identity, immutable updates, move generation and
stable ordering remain tagged v0.21. All failed storage/tag/clock/history versions
stay outside the playing build. Source overlays retain the native v0.21 greeting
until a reviewed promotion; development labels identify measured binaries.

Performance checks must invoke iterative deepening to exercise this feature.
Direct-root work remains an oracle and cannot substantiate aspiration speed.
Different tied best moves are acceptable when independent full-window value
proofs show optimal membership. Functional search-work reduction alone is not
measured playing strength; actual-budget decisions and the game screen are
separate gates. A later gradual-window or PV/selective-search change would be
a separate experiment, not an implied extension of this implementation.
