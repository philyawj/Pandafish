# Validation boundary

The [source snapshot](candidate-source.json) and [hashes](source-hashes.json)
contain only `search.go` and `deadline_test.go` over released v0.21. The original
[clock controls](../v0.22-monotonic-deadline/clock-controls-corrected.txt), full
race/vet/build logs and 325 cold score/move/node pairs test exactly those bytes.
`build-identity.json` proves the new development executable is byte-identical
to the original measured clock binary. No runtime adjustment follows the earlier
15% failure. The original fixture failure and Go source remain archived.

Controls cover zero/unlimited time, monotonic past/future, wall-only past/future,
far wall dates, injected before/equal/after boundaries, exact call counts and
canceled-context short circuiting. Existing recursive deadlines, last-completed
iteration, fallback, UCI stop/ponder/position, mate/repetition/GHI, perft and library
oracles remain active. No local-library file changes, so this scope introduces
no new move/transition/attack implementation requiring a new chess oracle.

Three alternating baseline/candidate warm pairs exercise real deadline reads at
depth four with the frozen 30s allowance. One TT persists per recorded game;
workers/history buffers are per go as in production. The 20-game / 1,145-root
workload is identical to the original assessment. Every score, move and node
count agrees across actors and repeated pairs; each visits 9,874,565 nodes.
Median search ratio is 0.93306 (6.69% less time), passing the prospectively
declared smaller 5% gate. No single pair or game is selected for its outcome.
Raw parsing, allocation, search intervals, commands, environments and inputs are
retained. The generic `membership_proof_pending` warm field is a checkpoint;
identical baseline choices/values give direct membership and final independent
row review resolves that checkpoint separately rather than changing raw output.

Four native suites use the existing 84 public and 145 newer inspected roots:
three public 100ms trials and one 1000ms trial, plus one newer trial per budget.
Three searches run at once, and all native work finishes before teacher queries.
The fixed teacher is Stockfish 19 at full Skill 20, 100,000 nodes, Threads 1,
Hash 64 MiB, MultiPV 1. Matching v0.21 reference caches retain original observations
and hashes. Newer inputs are regression controls, not a blind holdout.

The unchanged finite paired loss is the trial-mean clamped
`referenceCP-selectedCP-25`, in [0,1000]. Mean loss must remain within baseline
plus max(5cp,10%), severe (>300) cases within baseline plus two, with no new
losing-mate flag under the established finite/nonlosing-root convention and no
response beyond allowance plus 50ms. Require each newer turn's mean gate too.
Production roots are unrestricted; forced selected-move teacher queries only
label already chosen moves. Close alternatives are acceptable. Finite labels do
not prove a position was savable or that the teacher's first choice was unique.

Only after those gates pass, run the one routine 50-game screen against tagged
v0.21, Skill 3 and pinned Sacrifice, all three concurrent. Keep the original ten
openings, 100ms allowance, 50ms grace and 400-ply cap. Preserve all results,
timeouts, unfinished games, commands, identities and raw PGNs; legally replay
the complete screen and independently check settings/statistics/timing counts.
No CPU-heavy validation overlaps native/timing/game work. No automatic larger
confirmation or live game is authorized. A viable result stops before committing.
