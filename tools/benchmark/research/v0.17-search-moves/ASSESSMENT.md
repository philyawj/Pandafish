# Fresh choice after v0.16

I published v0.16 through [PR #31](https://github.com/philyawj/Pandafish/pull/31)
and [the release](https://github.com/philyawj/Pandafish/releases/tag/v0.16).
The tagged commit is `4ae65b14f69d95a1ad498baa759fec5e6bd0a1a6`. I approved
promotion from the standard 65% / 43% / 49% screen and deferred confirmation.
The original uncertainty, color split and every result remain preserved.

I reassessed the next two recommendations rather than assuming both should
be implemented. A fresh tagged-v0.16 profile uses the same 24 neutral-position,
100 ms timed kernel, with allocation-stack sampling disabled and no other
engine workload. It reports about **45.9% cumulative CPU in legal generation**,
**34.1% in move tagging**, **6.1% in static evaluation**, about **71.1 MB** and
**154,761 allocations per request**. Cumulative rows overlap; they are not
added. CPU share measures runtime cost, not the potential strength value of an
evaluator change; the choice also relies on the concrete safe backend route
and the positive v0.16 direct screen. This is a short, workload-specific profile, with visible GC/runtime
cost and uncontrolled desktop/thermal variance. Raw profiles and commands stay
in this directory. More objects than the older profile can reflect more
completed work at the same requested time; it is not a per-node regression claim.

| Candidate | Evidence now | Decision |
| --- | --- | --- |
| Specialized legal tactical generation and reusable move storage | v0.16 still builds/tag-checks quiet moves in capture search; generation remains the largest measured CPU cost. Prior search-efficiency work passed its direct strength screen. | Select one bounded additive-library-API experiment. |
| Coherent evaluation/threat model | Known quiet-mate, attack and pawn-ending weaknesses remain; Stockfish improvement is unproven. It needs independent response data and model/runtime decisions. | Retain as next major strength workstream; move it ahead if useful-depth gains stop addressing losses. |
| Full internal board/make-unmake rewrite | Existing pinned library already has bitboards; a narrow extension can filter destinations before tag/legality work and write into reusable storage. | Defer unless that API route fails its engineering gate. |
| More cache knobs, old LMR bundle, isolated bonuses or a playing book | No new result establishes a larger bottleneck or a missing prerequisite here. | Do not bundle. |

I chose an additive API in a local copy of pinned `corentings/chess` v2.6.1.
All existing runtime files and public behavior remain unchanged; source hashes
and the MIT license are retained. The additions generate only eligible legal
captures/promotions outside check, append full legal moves for main search and
check evasions, and test legal-move existence without allocating a list. Per-ply
move slices are reused without mutating the position's cached moves.

This is a search backend experiment, not evaluation work or another pruning
policy. I require identical fixed-depth node counts, scores and legal best-move
membership, plus legal-move/tag/status and child-FEN/hash equivalence. Empty
capture lists must never imply stalemate. Promotions, EP, castling, pins,
double checks, repetition, mate precedence and interruption remain controls.

I will run differential/perft checks against unchanged upstream code and
python-chess, the full tests/race checks, and equal-time fixed/timed measures.
A substantial reduction in work cost/allocation with unchanged decisions is
promising functional evidence; then the standard 50 games per opponent against
tagged v0.16, Skill 2 and pinned Sacrifice run concurrently. All failures stay.
No larger confirmation or live games run automatically. I will check in before
another commit, push or PR. No other candidate is implemented in parallel.
