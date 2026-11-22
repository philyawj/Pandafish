# Research checkpoint: v1.0 milestone and development pause

**Active engine development is paused.** v1.0 is the stable-engine milestone:
a working tested UCI engine, reproducible benchmarks and preserved research,
clear documentation/known limitations, and stable scope. GoPandaFish and other
Panda Bots will use the engine on Lichess. Further work is paused to focus on
other efforts. Release publication does not start a bot or change private configs.

The authorized scope is documentation cleanup, version-only release validation,
commit/PR/merge and v1.0 publication. Stop after publication/status reporting.
No next experiment, larger confirmation, live session or deployment starts
automatically. Research resumes only through a deliberate scope decision.

## Promoted research boundary

v1.0 promotes [protected queen-contact defense](../tools/benchmark/research/v0.24-contact-threat/README.md),
researched against tagged v0.23 (`28948e7473251d950003ec1f966f0cf26369b817`).
A legal quiet protected queen-contact mate selects actual defense verification;
a hypothetical pass never supplies a loss score. All legal defenses are eligible,
and the declared reply must remain legal and mate after the actual defense.
One event per q path, two remaining ordinary plies, window-independent behavior,
no SEE omissions during defense. Other evaluation/search policies remain.

- Original87-file source SHA:
  `447af5b0094a0319c92a1eea9c685cf16789204993556a6182e585025426ad5c`.
- Original candidate binary SHA:
  `ce4d89de73f89831c905e8810deb3384386e0cf73fd838e66475d8b8c995f683`.
  UCI reported v0.23; original development labels/hashes stay unchanged.
- Both-module normal/race/vet/builds,596 independent vectors, actual-response/
  window/history/mate controls pass. Cost:5.180% median overhead,0.987% more nodes.
  Three warm pairs retain all changed rows:19 scores/11 choices/406 rows.
- All six native/fresh gates pass. Adverse119cp/5cp existing and67cp fresh rows
  remain. Fresh seed2026112205/stream127 consumes138 roots,69 per turn/all16
  families, with full resolved exclusions and hashed temporary adapters.
- Sole standard screen:24W/6D/20L (54%) vs v0.23,66% vs Skill3 and77% vs Sacrifice.
  All150 games finish and replay legally; no forfeits or candidate attempts above
  100ms.87 frozen source files/binary identities match before/after play.
  The small selected lead is not a reliable general-strength or Elo estimate.
- [v1.0 validation](../tools/research/v1.0-validation.md) records normalized
  source identity, version-only integration smokes and platform build limits.
  Release archives and checksums are GitHub assets; binaries stay outside Git.

## Preserved work and future starting point

[Investigation summary](RESEARCH_V024.md), [original audit](REPOSITORY_AUDIT_V023.md),
[historical ranking](NEXT_ENGINE_ASSESSMENT_V023.md), [version history](HISTORY.md)
and [storage policy](EVIDENCE_STORAGE.md). Cleanup PR #42 merged without a release
at `68ef54b08ef9b63d66332243d148ed2044c80e38`. Local `handoff.md` stays outside Git.
Existing data remains intact; new diagnostics use compact records/source hashes.

Full legality context fails its8% net objective; target attacks save1.294% and
partial compact generation2.166%, both below8%. Coordinated compact search saves
8.162% but fails its screen at46% direct /54% Skill3 /76% Sacrifice. Dead-material
correction fails at50% /59% /77%; move-count claims fail at48% /63% /70%.
Contextual capacity stops before fitting because a required fixed-witness pair
has exactly identical FENs. No board-only evaluator can separate those targets.
All original failures/source versions and inventories remain preserved.

Consumed recent cohorts: streams122/123 (seeds2026112122/2026112123;142/148 roots),
124/125 (2026112201/2026112202;144/144),126 (2026112204;129),127 (2026112205;138),
plus all prior/public/newer/live/warm/oracle controls and contact witnesses.
If research resumes, review resolved exclusions; a new seed alone is not
independence. Use the [workflow](RESEARCH_WORKFLOW.md) and
[Go tooling map](RESEARCH_TOOLING.md), with v1.0 as the released baseline.
Preserve failed gates and avoid broad fitting or variant sweeps.
