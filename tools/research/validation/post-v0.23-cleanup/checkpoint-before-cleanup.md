# Current research checkpoint

- Branch `research/v0.23-candidates`, HEAD `8b0109d`. Baseline tagged v0.22,
  `1953f2af97f956c4499cbc5b61c3c1c698856740`. Continue through failures until a
  recommended candidate; stop **before commit**. No publication/live play.
- User reports excessive model usage after multiple five-hour runs, with53%
  weekly capacity remaining. Follow the bounded reassessment instructions in
  [the tooling flowchart](RESEARCH_TOOLING.md) and [workflow](RESEARCH_WORKFLOW.md).
  One next measurement; Go executes it. No speculative sweep/general framework.
  Live mistakes are regression controls, not the candidate-selection agenda.
- **23 unsuccessful substantive attempts:**18 prototypes, one offline model,
  four diagnoses. Two successful opportunity studies are separate. Consolidate
  the older assessment index/RESEARCH_V023 summary at a viable milestone; current
  attempt READMEs/raw states take precedence over stale pending rows.
- Three sole screens: quiet checks50/50/77%, root parallel50/47/80%, quiet scouts
 57/45/77% (direct/Skill3/Sacrifice). All450 legal/finished, no forfeits. The last
  direct pass does not offset its adverse external signal. No unchanged reruns.
- `v0.23-check-budget`: correctness passes, elapsed ratio1.363142 exceeds1.25.
  Sealed45files. Two timeouts, bad initial fixture and oracle memoization repair
  preserved; expanded race passes at319.928s. No native/games.
- `v0.23-check-budget-one`: controls, cost (elapsed0.714263/nodes0.624976),
  same-horizon regret and four existing suites pass. Fresh100 loss42.919→32.879
  cp improves23.394%; fresh1000 Black19.677→26.462 fails. No exception/games.
  Sealed103files. Seed2026112122/stream122 is consumed:142roots/71eachturn.
  Future fresh collection must exclude this cohort along with all earlier ones.
- `v0.23-static-cut-opportunity`: exact1,145roots/43,118,858nodes;20.808% outer
  descendant opportunity,179 false predictions including one hidden losing mate.
  Safety gate fails. Sealed12files. Do not activate unverified static pruning.
- Verified-static follow-up fails the same hidden-mate gate. Its proposed q
  verification is flawed: stand-pat guarantees immediate acceptance from the
  static prerequisite. All1,077,982 probes cost exactly one node. Sealed12files.
  The design mistake is explicit; this pruning line is closed. No native/games.
- **Current** `v0.23-state-kernel`: normal/race/vet pass in both modules, both
  builds and596 independent vectors pass. One coordinated state project combines
  archived child-storage reuse and virtual ordinary tags. Restore v0.22 first;
  no pruning/eval/order/q/table/clock changes. Earlier separate15% gates stay
  failed. New combined target:>=12% median warm saving, exact1,145values/moves/
  nodes; all six native/fresh gates still mandatory. Seed23/123 unconsumed.
- Active source: `search.go`,`iterative.go`,`search_position.go/test`, new chess
  SearchPosition, virtual tags and renamed simulated fallback in search_update.
  Nested tag oracle test and independent596-vector helper exercise actual
  scratch child FEN/hash/perft. Shared validation includes nested normal/race/vet.
  Three warm pairs pass: elapsed ratio0.858578 (14.142% saving); all1,145 scores,
  moves and43,118,858nodes exact. Snapshot118Go files saved. Running
  All four existing and both fresh suites pass without exceptions/new mate
  flags/overruns. Fresh seed23/123 yields148roots,74eachturn,16families;
  loss100ms42.325→40.184cp,1000ms30.630→28.197cp. Cohort is now consumed.
  **Sole standard screen complete/reviewed**, direct benchmark, console`screen-run.txt`,
  results`screen/v0.23-state-kernel-dev-vs-v0.22-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba`.
  All3concurrent/50each/100ms+50ms: direct20W11D19L=51%; Stockfish34W0D15L,
  one400-ply cap=69.4% completed/68–70% possible; Sacrifice34W8D8L=76%.
  All150 replay legally,149 finish, no forfeits, candidate max99.3ms/no100ms
  overruns. **Recommend moderate functional promotion; strength inconclusive.**
  Direct lead is one game. No large Elo claim/Skill increase; broader-opening
  confirmation is justified before such claims, not run automatically.
  Promotion, PR, merge and v0.23 release are now explicitly approved. Complete
  RELEASING.md with version-only validation/smokes and preserve this screen's
  original identity. Stop after publication; no next experiment or live play.
  Final source snapshot is byte-identical to the tested118-file snapshot;
  candidate binary SHA6926bf6ce9cd13fbac87e3c9cf2fe9e033ca26b9ba9f2c6495ee45fdc53fe745.
- Use `GOCACHE=/tmp/pandafish-v016-go-cache`. AppleM2/8CPUs/AC. No builds/tests/
  profiles overlap timing, rollout collection, native tests or games. Baseline
  warm test `/tmp/pandafish-warm-v022-shared.test`. Shared fail-closed validation,
  warm/native plans remain. Standard50each/all3concurrent/100ms+50ms/10openings.
  Call benchmark directly; Markdown is at results-root README, not nested *.md.
- Preserve all failures. Update EFFICIENCY_LOG at eventual PR checkpoint: Go
  avoids manual orchestration and stops unused stages, but no measured token
  saving/net payback exists and model supervision remains too costly.
