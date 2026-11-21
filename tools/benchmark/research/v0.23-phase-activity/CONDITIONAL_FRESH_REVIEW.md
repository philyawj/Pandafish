# Prospective independent review of the narrow native failure

The original native prerequisite is **not satisfied**. Newer Black 1000ms mean
loss is 27.0923 cp versus baseline 21.6462 and a 26.6462 limit: 0.4462 cp beyond
the tolerance. That report remains `passed:false`, and the original driver
exits one. No weight, reference score, threshold or root restriction is changed.

The other three suites pass. Public means improve 45.896→43.119 at 100ms and
54.851→51.478 at 1000ms; newer 100ms changes 31.511→33.420 within both-turn
tolerances. Newer 1000ms overall changes 24.242→25.947, with White improving
26.761→24.836. Severe counts and losing-mate flags do not worsen; there are no
timing overruns. The largest Black deterioration is the actual depth-six Bxe5
choice in `proof-fresh-10-ply37`, 235cp worse by the existing loss surrogate.
It is retained as an adverse example, not forced to Ba6 or removed from scoring.

This differs from the rejected expensive evaluator's failed standard-budget
controls: all primary 100ms gates pass here, and hot evaluation work is unchanged.
The coordinate/phase map is small, symmetric and frozen. The previous checking
experiment also demonstrates that passing every surrogate does not guarantee a
match gain. A narrow finite-reference bucket miss is therefore not treated as
conclusive playing-strength evidence in either direction.

Before inspecting new outcomes, allow the already planned fixed 16-family seed
2026112043 control, blocking the consumed quiet-check control boards too. It
must pass **both existing turn/safety gates and at least 5% lower aggregate
mean finite loss at each of 100ms and 1000ms**, with no increase in severe cases.
No additional batch, coefficient rescue or favorable rerun is permitted. If
that stricter independent control fails, archive the map without games. If it
passes, document a separate bounded screen decision before one standard game
screen, keeping the original native failure explicit. This is a prospective
exception for evidence gathering, not a retroactive native pass or promotion.
