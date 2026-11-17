## Publication note for v0.19

I preserve the original local experiment report below, including its stop/check-in
state. I subsequently received approval to publish these attempts together in
v0.19. I keep playing policy at v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

# Current handoff

I completed the user-authorized v0.18 promotion: PR #33 merged, GitHub release
published, tag points to `75f12bbdbeb0e38d148690ddc0f1d6405c5af3d3`.
The standard development screen was 57% vs v0.17, 73% vs Skill 2 and 51% vs
Sacrifice, all 150 legally replayed, zero forfeits/unfinished. The six-game
release integration smoke's 50% direct gate returned exit 1 and is preserved.
No larger confirmation or live client ran. The tagged baseline is built as
`tools/benchmark/bin/pandafish-v0.18`, with verified v0.18 UCI identity.

I then created local branch `v0.19-decision-evaluation`, wrote the short ranked
RESEARCH.md, froze PLAN.json and completed its top-choice experiment. No new
production Go code is active: all runtime source remains exactly tagged v0.18.
The entire new research directory is uncommitted/unpushed; check in before a PR.

The new objective uses independent generated families, full-history root
alternatives and exact depth-three released-search quiet-leaf witnesses.
Collection: 40 families, 211 sampled roots, declared splits129/39/43, 693
alternatives. All paths/FENs replay independently; all676 static vectors and
branch scalars match. One sampled terminal-root collector assertion was fixed,
with original source/query/trace and completed-label resume manifest retained.

The fit excludes25 incomplete/mate/nonstatic roots, leaving114/34/38. It uses302
training pairs and selects ridge0.01 on validation decision regret. Frozen
weights precede holdout inspection. Validation60.26→50.79cp improves, but
holdout54.02→54.74 worsens1.32%, so the ≥5% improvement gate FAILS. No model
activation, runtime screen or strength games follow. All weights/grid/data are
preserved. Do not retune these controls until favorable.

A stricter post-result audit fromply4 finds two early positions shared across
training/validation families16/19/21, and three early positions related to the
reserved confirmation suite. No sampled root duplicates cross splits, but
independence is weaker than intended; grouping fromply8 missed those related
families. Retain original splits/results without claiming blind validation or
resplitting after outcome. A future model needs wholly fresh families plus
canonical grouping/exclusion from the early family boundary. These research
observations do not alter the earlier frozen v0.18 release.

Read README.md, RESEARCH.md, PLAN.json, ERRATA.md, NEXT.md and final-audit.json.
The next recommendation is repair/test the data boundary before another learned
fit; if choosing a separate playing-code improvement, reprofile v0.18 legal
check/tag/attack work before selecting one exact backend optimization. Do not
automatically jump to a large board rewrite, another weight or NNUE solely
because this fit failed. No automatic next candidate or PR follows this report.
