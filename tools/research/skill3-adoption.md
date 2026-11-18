# Skill 3 reference adoption after v0.20

Explicit review of the tagged v0.20 ladder selects **Stockfish Skill 3 as the
new routine reference**. The benchmark CLI defaults to 3; explicit
`--stockfish-skill 2` retains the old reference. Routine comparisons remain 50
games each against the tagged previous release, Stockfish and pinned Sacrifice,
with all three matchups concurrent. No engine release accompanies this change.

The [complete ladder](../benchmark/results/v0.20-stockfish-ladder/README.md)
preserves every Skill 2–10 matchup, all 450 game objects, options, identities,
timings, PGNs and failures. The clean Skill 2/3 scores are 62%/48%, making Skill 3
a balanced next reference. Skills 4/5 score 32%/34%; their small inversion does
not establish a difficulty reversal. The [analysis](../benchmark/results/v0.20-stockfish-ladder/ASSESSMENT.md)
states the uncertainty and unresolved timeout causes.

One game at Skill 7 is unfinished at the 400-ply cap. Three timeout forfeits
remain: Stockfish at Skill 6, and Stockfish/Pandafish once each at Skill 9. Their
official results are kept, and affected rows are qualified. The adoption follows
review of the clean lower-level evidence; the entire ladder retains a failed
reliability status. No failed game was rerun, no sample was extended and no
upper-level result was used as a clean strength estimate.

The original collector's reporting assertion failure, initial diagnostic build
failure and report-only correction remain in the archive. All 450 saved games
pass legal PGN replay. The postcheck verifies the frozen plan/driver/binaries,
tag, color/opening balance, maximum concurrency and combined-PGN bytes. Its
inventory freezes 82 pre-adoption evidence files. Archived pending-review/local-
checkpoint wording describes that earlier stage; this document records the
subsequent authorization to merge the data and adopt Skill 3.

## Validation and version boundary

The implementation changes only the benchmark's default flag and workflow
documentation. Engine/library source remains byte-identical to tagged v0.20:
[65-file policy check](validation/skill3-adoption/engine-policy.json).
Benchmark [race tests](validation/skill3-adoption/race.txt) and
[vet](validation/skill3-adoption/vet.txt) pass. The short
[default-option smoke](validation/skill3-adoption/smoke/README.md) omits the skill
flag, so its saved Stockfish options verify that the new default is actually 3.
This validates integration, not playing strength. Every smoke result is kept.
The four smoke games complete without forfeits or unfinished games: 1/0/1
against Skill 3 and 2/0/0 against Sacrifice. The command exits 0, and the
independent PGN check legally replays all four saved games.
Full-module tests, vet and an engine build pass as well; their logs remain in
the same validation directory.

Current documentation mappings are updated explicitly while original archive
hashes and observations remain unchanged. Source-policy, archive, factual prose,
table and local-link checks supplement the focused validation. All new source
and diagnostic programs are Go.

Tagged v0.20 still contains its original benchmark default of 2. The post-release
PR changes the tool on `main`; it does not retag v0.20, change UCI identity or
publish another release. The ladder's explicit per-level flags and frozen tagged
runner remain reproducible. Future candidates use the current Skill 3 runner
and compare their engine against tagged v0.20. Larger confirmations and live
games remain separate deliberate decisions.
