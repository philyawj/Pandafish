# Precomputed stable move ordering: candidate for promotion

**Recommend promotion after review.** Calculate each existing move-ordering
priority once, stable-sort small index/priority records, and apply the exact
permutation to the original moves. Selected depth-four elapsed time falls 18.9%
versus v0.20, with all 325 distinct score/move/node comparisons exact. The
standard concurrent screen scores 55% against v0.20, 50% against Skill 3 and 71%
against pinned Sacrifice, with no forfeits or unfinished games.

Read the [plan](PLAN.md), [implementation](IMPLEMENTATION.md),
[testing and reproduction](TESTING.md), [results and limitations](RESULTS.md),
and [reassessment](../v0.21-next-assessment/NEXT.md). The
[initial selection](../v0.21-next-assessment/SELECTION.md) explains why fresh
profiles moved ordering ahead of more temporary board-state work.

The engine/library base is released v0.20. Replay uses merged Skill 3 tooling
commit `6bfb4d565c2adf2a2ad59216a462b0f4f0c01503`, whose engine policy is v0.20.
`candidate-source.json` contains only `search.go`, `search_moves.go` and the new
Go rank-order controls. Values, tie ordering, eligibility, quiescence, TT policy,
repetition, clocks and every library file remain unchanged.

The original measured label is `v0.21-rank-order-dev`; native UCI identity still
says `Pandafish v0.20`. Its SHA-256 is
`d9f2aee7c30832e00690e1c24f1490e27b7dcfe0dda6810858d71dd4bb232f1f`.
This is not a published v0.21 or a test of a future renamed/tagged executable.
The 55% direct lead passes the practical gate but remains a modest, uncertain
strength sample. Efficiency is established more firmly than the generality or
size of a playing-strength gain.

All initial profiles, controls, native/reference queries, timing pairs, commands,
hashes and games are preserved. This first implementation clears every declared
gate; no failed variant was discarded, favorable rerun made or larger confirmation
started. No experimental commit, push, PR, release or live game follows
automatically. The current branch is `research/v0.21-next-candidate`; the next
action is the requested candidate review before a commit.
