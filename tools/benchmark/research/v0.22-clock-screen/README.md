# Rejected clock strength screen

Decision: **do not promote**. The [results](RESULTS.md) distinguish the 6.69%
timed-warm functional gain, passing actual-budget native controls and the failed
46% previous-release screen. Skill 3 scores 47%, Sacrifice 71%; all 150 games
finish legally with zero forfeits. No larger confirmation or unchanged retry
runs. Source is archived and restored before the next experiment.

[Plan](PLAN.md), [implementation](IMPLEMENTATION.md), [testing](TESTING.md) and
[cache-argument errata](ERRATA.md) retain the complete scope and correction. The
two-file source and binary are identical to the earlier clock attempt, whose
larger 15% goal stays failed. The separately declared smaller functional gate
passes; it does not turn either failed objective or the game gate into a pass.

`candidate-source.json` restores this exact code over v0.21 through
`tools/research replay --source ... --out NEW --check`. All raw warm/native/game
data, the initial cache-argument failure, reference-only repair and independent
reviews remain. Native sources/reference caches keep their original identities;
no development result is relabelled as a released version.
