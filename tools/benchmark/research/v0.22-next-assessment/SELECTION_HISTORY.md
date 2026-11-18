# Selection: bounded quiet-cutoff history ordering

The instrumented v0.21 diagnosis passes the predeclared opportunity criterion.
Public depth-four roots have 94,770 main cutoffs, including 13,184 quiet cutoffs
after the first move (**13.91%**); newer roots have 295,696 / 43,222 (**14.62%**).
After the second move, those counts remain 8,687 / 31,566. First-move cutoffs
already account for 81.25% / 80.54%, so history is an opportunity, not an assured
large gain. Raw per-root results, counter outputs and instrumented source are
retained. The two diagnosis processes run concurrently; no elapsed claim uses
their instrumented timings.

Implement only per-search quiet history. Retain a lazy 2×64×64 int16 table indexed
by moving side/from/to. A completed main-search quiet beta cutoff receives a
depth-squared bonus capped at 2,048 and the bounded update
`h += bonus - h*bonus/16384`. No negative updates or parameter fitting in this
first attempt. Share within iterative deepening; discard at every new search,
like existing killer hints. This history is distinct from repetition history.

In main search, stably order only the quiet suffix after captures/promotions,
using existing small rank/permutation scratch. Then apply existing killers and
TT preference in their current precedence. Never add/remove a move. Root
ordering, all quiescence ordering/eligibility, evaluation, horizons, PVS windows,
TT/history identities and clocks stay v0.21. No earlier failed optimization,
aspiration, LMR, null move or guessed positional term is bundled.

Ordering changes expected node counts and may select a different exact tied
best move. Replace the exact-tree equality requirement with equal completed
scores and independent tagged-baseline full-window selected-move membership on
all 325 public/deep/newer comparisons and the persistent-TT sequence. Do not
force the baseline's root move. Retain the **15% median elapsed target** on three
alternating public depth-four pairs, beneficial deep/newer and warm work, plus
node accounting. Require bounds/reset/precedence/tie/ownership controls and full
race/vet/build. A zero-history control must preserve original stable ordering.
The same four actual-budget native gates and one three-concurrent 50-games-per-
opponent screen follow only after those gates pass. Preserve failures and stop
before committing a viable candidate. This is full-width ordering research, not
selective-search pruning or an evaluation change.
