# Next-candidate assessment after Skill 3 adoption

The baseline playing engine is released v0.20, tag commit
`c90ca44d71ce4c3897f43402976f4ca39ee1dc73`. The working research base is merged
Skill 3 tooling/data PR #37, commit
`6bfb4d565c2adf2a2ad59216a462b0f4f0c01503`; its engine/library policy is identical.
No new release accompanies that adoption. Future strength screens use tagged
v0.20, **Stockfish Skill 3** and pinned Sacrifice, 50 games each, all concurrent.

## Fresh priorities before profiling or implementation

1. **One profile-supported reduction in temporary legal-state work.** v0.20
   removed repeated transitions/rays/tags, gained 38.7% selected fixed-work time
   and passed its direct screen at 65%. Temporary move-tag validation still
   constructs board state for attack predicates; immutable searched children
   serve a different purpose. Reprofile the resulting engine before selecting
   a precise change. Attack-only state or affected convenience-state updates
   are candidates if their measured cost remains substantial.
2. **A different evaluation representation with actual-budget selection.** Prior
   fitted values repeatedly failed generalization/safety. A new contextual or
   nonlinear Go pipeline needs representative data, an inference budget and a
   truly unused final cohort. More tuning of inspected coefficients has weaker
   support than a measured exact-computation target.
3. **Reconsider tactical-horizon proofs only under a new cost/trigger design.**
   Saved legal quiet mates establish an omission; the standalone proof failed
   fresh Black-turn choices. Corrected EP classification does not justify
   reactivation. Search selectivity or make/unmake likewise needs a stronger
   measured case and separate correctness boundaries.

The ladder's clean 62%/48% at Skills 2/3 motivates the new reference, not a
specific evaluation term. Its three timeouts (two Stockfish, one Pandafish) and
one unfinished game remain preserved. Their causes are unresolved; no clock or
scheduling change is bundled into a speed-only attempt. Native response gates
and all failures remain mandatory evidence.

## Selection and acceptance protocol

Capture CPU/allocation profiles of the tagged engine's existing Go fixed-depth
driver on the 84 public roots and the 145 inspected development roots. Preserve
all commands/profiles/logs. The installed local Go pprof tool may render the
profile; the engine compiler remains Go 1.25. Profiles are selected workloads,
not playing strength. Freeze the chosen mechanism before editing runtime code.

An exact-computation candidate must preserve original upstream runtime/license
bytes, immutable ownership, rules, complete returned tags, move order, hashes,
mate/repetition and cancellation. Run differential legal-transition/status/
special-move/perft controls and independent mailbox attack checks. Preserve
scores, chosen moves and nodes on all 325 distinct public/deep/newer fixed-depth
combinations. Run three alternating depth-four timing pairs; require at least
**15% median elapsed reduction** versus v0.20 and no contradictory deeper result
before native probes. A smaller gain needs reassessment, not an automatic screen.
Material, values, search eligibility, quiescence, TT policy and clocks remain
unchanged in this first exact-kernel experiment.

Surviving native checks use the frozen public 84 inputs at three 100ms trials
per actor and one 1000ms trial, plus the inspected 145 newer development inputs
at one trial per actor/budget. Use three concurrent native searches; finish
native work before reference queries. All legal production root moves remain
eligible. Full-strength pinned Stockfish labels selected moves at 100,000 nodes,
Threads 1, Hash 64, MultiPV 1; existing matching caches may be reused with hashes.
The game opponent is separately Skill 3. Close alternatives are acceptable.

For fully finite roots, retain the per-position trial mean of
`clamp(root_reference_CP-selected_reference_CP-25,0,1000)`. Require aggregate
mean no worse than baseline plus max(5cp,10%), severe cases (>300) no more than
baseline plus two, no newly selected losing-mate flag under the existing
finite/nonlosing-root convention, and no response past allowance plus 50ms.
Each newer root turn must also meet the mean tolerance. These reused positions
are regression controls, not a blind holdout or proof that finite scores are
salvageable. No root move is forced and no reference move is uniquely required.

Only a surviving mechanism reaches the standard 150-game concurrent screen.
Require the established previous-release practical gate (>50% counting unfinished
as losses, zero forfeits), then review both external matchups and all evidence.
Preserve every failure, source revision and raw game. No unchanged favorable
retry or automatic larger confirmation. Archive rejected Go code outside active
play and reassess the next hypothesis; continue autonomously until a candidate
warrants presentation. Check in before committing that candidate. No live games,
commit, push, PR, tag or release is authorized for the research phase.
