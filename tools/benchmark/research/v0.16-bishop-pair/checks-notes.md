# Check and preservation notes

- The Go compiler warned that the sandbox denied a module metadata stat-cache
  write under `/Users/philyawj/go/pkg/mod/cache/download/.../v0.15.2.info...tmp`.
  The explicit ARM64 build returned exit 0 and produced the candidate. All
  normal/race tests, vet, production build and offline wrapper passed.
- `controls.txt` predates the added mirror test and stronger pawn recapture in
  the tactical control. `controls-final.txt` is authoritative for the final
  tests. There is no production source change after the match source freeze.
- The original preflight wording described regressions too narrowly. Full
  per-case review records the 211 cp Kf7 regression explicitly. The record is
  preserved and no weight is tuned to it.
- Exactly one routine 150-game screen ran. No failed match or gate is hidden;
  all three matchups complete, and the practical previous-release gate passes.
  This does not establish a strength gain.
