# Selection: one bounded aspiration probe with full-window recovery

The tagged depth-five diagnosis passes the fixed opportunity criterion: adjacent
finite depth-three-through-five changes stay strictly inside ±50cp in 204/219
public transitions (93.15%) and 370/416 newer transitions (88.94%). Mate-band
transitions remain separately counted (33/19). Both turns are retained. These
instrumented, concurrently collected traces establish stability on selected roots,
not timing or strength. The 50cp band is not tuned after observing these data.

Start from restored v0.21. At depth three and later, when the previous completed
score is outside the mate band, search the root at full requested depth with
alpha/beta equal to previous score ±50cp. Accept only a complete result strictly
inside the window. On either bound failure (including equality), retry once using
the original full root window. No gradual widening, depth reduction, omitted move,
new root preference, evaluation term, history table or failed kernel is included.
Prior mate scores and depths one/two use the full root search. Unfinished probes
or retries never replace a previously completed iteration.

Factor the root into a window-capable helper with the old full root as its
unchanged public/internal wrapper. At infinite bounds, preserve old child windows,
score/tie behavior and nodes. Other alpha-beta/PVS/quiescence/TT/history and clocks
remain intact. Root bounds are private to the current pass; only exact completed
results become the production result. Preserve mate distance and history unwind.

Add direct in-window, fail-low/high equality/recovery, mate/depth bypass,
complete-score/best-set, canceled-probe/retry and last-completed-result controls.
Require race/vet/build and original root/node controls. Use a frozen Go test helper
that invokes the actual production iterative entry point in both actors: public
depth three/four, 12 deep depth-six roots, newer depth-four, plus public/newer
depth-five controls. Node counts may change; scores and independent tagged-baseline
selected-best membership must agree. No unique baseline/reference root is forced.

Efficiency follows the prospective iterative-workload principle: require no
aggregate regression on public depth-five median, deep depth-six and newer
depth-four/five controls, then **at least 15% median elapsed reduction on three
alternating 1,145-root warm iterative game-sequence pairs**. Freeze all helper/
driver changes before measuring; preserve cold/full-root and every warm result.
The actual-budget four native gates and standard three-concurrent 150-game screen
remain mandatory. No larger confirmation or publication follows automatically.
