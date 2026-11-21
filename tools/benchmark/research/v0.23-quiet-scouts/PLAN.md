# Count quiet candidates rather than all moves

The frontier guard preserves exact depth-four controls but saves only 12.128%
nodes, below its unchanged target. It remains archived. Restore v0.22 before
this separate definition of late quiet search, using the material-protected
implementation only as an explicitly preserved foundation.

The original move index counts preferred moves, captures and promotions before
quiet generation. With three tactical moves, the first ordinary quiet can be
classified as late. Replace that ordinal with the count of previously searched
nonchecking, noncastling quiet candidates. Protect the first three such choices
regardless of preceding captures/promotions/checks. Castling remains full depth;
checking/tactical/preferred killer exceptions and phase>=8 remain. One-ply
reduction starts at parent depth three. Alpha-improving or mate-reporting probes
are verified at full depth before ordinary PVS re-search. No move is omitted.

This is a general ordering/eligibility definition, not a FEN or root-move fix.
It does not claim to repair the predecessor's pre-existing quiet-check horizon.
Keep material/evaluation/quiescence/cache/time policy and single-worker search.
Use guard, capture-before-quiet, verification, independent mate, context/race
and actual-budget controls. An explicitly unreduced killer-only audit retains
the original isolated-order contract; independent regret uses tagged v0.22.

Freeze >=20% node saving per warm pair and >=15% three-pair median elapsed
saving. Require all four existing native suites, a new both-turn final cohort
seed2026112121/stream121 with all consumed/live exclusions, and the same full-
depth selected-choice regret gate before the ordinary route to one concurrent
150-game screen. Preserve false gates and document any research-scope
reassessment before games. No automatic larger confirmation/publication;
check in before committing a viable recommendation.
