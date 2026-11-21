# Shallow static-margin pruning: bounded opportunity test

The two checking-budget attempts are sealed failures: deeper coverage exceeds
cost; one-choice coverage improves fresh100 but fails fresh1000 Black. Restore
released v0.22. Do not repeat those binaries, extend their screens or adjust
their margins/seeds. Twenty-one substantive attempts are unsuccessful so far.

Reassess older alternatives: positive/signed history and aspiration saved only
about0–6% warm work on the older core; child reuse, tags, stages and mate bounds
also miss larger-gain targets. Another weight fit lacks new signal. Make/unmake
has a larger correctness boundary and allocation savings already underperform
elapsed estimates. Prefer a cheap selectivity diagnosis before another prototype.

Investigate reverse futility: at a shallow null-window main node, static value
well above beta may justify a selective fail-high without searching children.
This is a different mechanism from null-move probes and late quiet reductions.
The official [Stockfish19 search source](https://github.com/official-stockfish/Stockfish/blob/sf_19/src/search.cpp)
contains child-node futility pruning; its evaluator, margins and surrounding
search are different and do not establish a gain for Pandafish. No code or tuned
constants are imported from it.

Freeze the prospective boundary: main depth1–3, ply>=3, exactly one-centipawn
window, not checked, alpha/beta outside mate bands, total nonpawn phase>=8,
and released staticScore >=beta+150*depth. The150cp/ply safety margin is a
conservative untuned hypothesis, not a fitted optimum. Terminal/repetition/TT
checks retain precedence. Quiescence, material, clocks and cache rules stay v0.22.

First measure on the unchanged1,145-root depth-five trajectory. Do not actually
prune or count nested work twice. Count descendant work of outer qualifying
nodes, every qualification, false static predictions against original completed
node results, and any hidden forced losing mate. Require exact released values,
moves and nodes, >=20% outer descendant opportunity overall and >=15% each turn,
and no original forced losing mate concealed by a prospective cut. Instrumented
time is not speed evidence. Stop before a prototype if the opportunity/risk test
fails. No threshold sweep or fresh corpus is part of this diagnosis.

If promising, freeze a separate implementation/control/cost plan before actual
pruning. It must earn independent selected-move review, existing/fresh native
controls and the usual concurrent screen. No publication or live play; check in
before committing a recommended candidate. Reuse current Go tooling, with no
new general framework or expanded research campaign for this initial test.
