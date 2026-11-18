# Search-legality kernel candidate

The candidate specializes immutable board transitions and removes repeated or
discarded attack/tag work. Read the frozen [plan](PLAN.md), detailed
[implementation](IMPLEMENTATION.md), [testing record](TESTING.md),
[results](RESULTS.md), and [preserved findings](ERRATA.md).
The [handoff](HANDOFF.md) states the exact review checkpoint and outstanding action.
The [combined research assessment](../../../../docs/RESEARCH_V020.md) records
the rejected predecessors and next priorities.

The base is tagged v0.19, commit
`e7bf00fe5b158e3d58ca77e83385ed918f839717`. Search depth/eligibility, material,
positional values, transposition policy, repetition, ordering, quiescence and
time management remain unchanged. Original upstream runtime methods, module
metadata and license remain unchanged. There is no opening book or online data.
The quiet-mate experiment is archived and absent from the active engine.

The development label is `v0.20-search-legality-dev`, with native UCI name still
`Pandafish v0.19`. Tested binary SHA-256:
`1f40c45929750169a68af3101ff364821218a236b8e4f5229ad0789ca81861b4`.
`candidate-source.json` and `source-hashes.json` freeze the complete Go overlay
used to build it. This label is not a published version or a test of a future tag.

The standard concurrent 50-game-per-opponent screen scores **65% against tagged
v0.19, 67% against Skill 2 and 58% against pinned Sacrifice**, with no forfeits or
unfinished games. Promotion is recommended under the established practical
gate, pending review. The exact-tree speed finding is firmer than the size or
generality of the playing-strength improvement.

The initial predecessor's [failed screen](../v0.20-search-transitions/RESULT.md)
remains part of the evidence. A further 22.3% reduction in normalized fixed-work
time justified testing this substantive revision. The median total reduction
versus v0.19 is 38.7%, with all 325 distinct score/move/node comparisons exact.
Correctness and speed alone do not establish playing strength.

All raw profiles, initial timings, alternating comparisons, native choices,
teacher queries, failures and games are retained. The inventory records their
bytes and hashes. No larger confirmation, live game, commit, push or PR runs
automatically; the candidate must be reviewed before a commit.
