# Decision: archive the staged-picker candidate

The first exact staged picker fails the direct previous-release criterion.
Preserve its functional success and all adverse playing evidence. Do not promote
or repeat the unchanged binary for a more favorable score.

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.21 | 17 / 13 / 20 | 0 | 0 | 47% |
| Stockfish Skill 3 | 34 / 0 / 16 | 0 | 0 | 68% |
| Pinned Sacrifice | 23 / 16 / 11 | 0 | 0 | 62% |

All three matchups run concurrently, 50 games each, 100ms plus 50ms grace,
400-ply cap and ten unchanged openings. Runtime is 605.63 seconds. The runner
returns exit 1 for the direct gate, after retaining all three rows. The tested
binary is `v0.22-staged-picker-dev`, SHA-256
`e96f2c19001506d7d51d3d607f11abd62c53913f16eee688f0bae80ad306ce74`.
The [original summary and PGNs](screen/v0.22-staged-picker-dev-vs-v0.21-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba/summary.json),
`screen-command.json`, `screen-status.json` and complete logs retain identities,
options, source/base revision, timing policy and host details.

Independent replay validates all 150 games. `screen-review.json` recomputes the
scores, color/opening pairs, game lengths, response accounting and failed direct
gate. `screen-strata.json` retains all opening families. The direct score splits
34% White / 60% Black; color-pair outcomes are six won, twelve drawn and seven
lost. Skill 3 splits 60% / 76%, Sacrifice 58% / 66%. This is limited evidence from
repeated opening families, not a precise Elo estimate. The 68% external row is
promising but does not erase the direct failure; comparing it with a different
run's absolute external score is not a paired baseline comparison.

Correctness controls, exact cold/warm search results and the 19.81% warm time
saving remain valid selected-workload findings. Public 100ms native loss still
fails; the other three native suites pass. The explicit pre-game exception in
[SCREEN_DECISION.md](SCREEN_DECISION.md) remains visible, without a threshold
change or actor-query retry. Stronger fixed-work efficiency did not establish a
promotion-worthy playing improvement in this screen.

Scratch replay/tests/build, source identity, current Go-only policy, Linux ARM64
build and the independent reviews pass after the screen. The phase-specific
source proof is `../v0.22-next-assessment/final-source-review.json`; it describes
the active five-file boundary at this checkpoint, before later research changes.
No larger confirmation or live games run. Source/tests stay archived and the
next mechanism requires a fresh diagnosis rather than an unchanged rerun.
