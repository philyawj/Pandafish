# Recommended candidate checkpoint

[PR #37](https://github.com/philyawj/Pandafish/pull/37) is merged and Skill 3 is
the routine benchmark default. That data/tooling PR publishes every tagged
v0.20 Skill 2–10 matchup and its failures without another release. The current
research branch starts from its merge commit
`6bfb4d565c2adf2a2ad59216a462b0f4f0c01503`.

Recommend promotion of this rank-order candidate after review. Its first
implementation passes all declared gates: exact scores/moves/nodes on 325
distinct position/depth combinations; 18.9% median fixed-work elapsed reduction;
four native decision suites including both newer root turns; and the standard
three-concurrent 50-game screen. Results are 23/9/18 versus v0.20 (55%), 23/4/23
versus Skill 3 (50%), and 28/15/7 versus Sacrifice (71%). All 150 games finish
with zero forfeits. The direct lead is modest; no precise/general strength
gain follows from this sample alone.

Start with [results](RESULTS.md), [implementation](IMPLEMENTATION.md),
[testing](TESTING.md), the [fresh selection](../v0.21-next-assessment/SELECTION.md)
and [reassessment](../v0.21-next-assessment/NEXT.md). Every library file, value,
priority/tie rule, move eligibility, search/TT policy and clock remains unchanged.
Full race/vet/build controls, independent choice/screen accounting, source/binary
hash review, all 150 legal PGN replays, scratch restoration/tests/build,
Linux/ARM64 build and source-policy checks pass. Raw queries, profiles, commands,
timings and original labels remain preserved.

The measured label is `v0.21-rank-order-dev`, native identity `Pandafish v0.20`,
binary hash `d9f2aee7c30832e00690e1c24f1490e27b7dcfe0dda6810858d71dd4bb232f1f`.
Three Go overlay files and hashes reproduce the candidate from the declared
base. No experimental code is committed, pushed or prepared as a PR. No version
bump, additional release, larger confirmation or live game has occurred.

The next action is the requested review before a commit. A broader-opening
confirmation is justified before a large/general strength claim, but is not a
new mandatory promotion gate or an automatic next run. Future hypotheses must
be reassessed from the current profile; no further feature is bundled into
this successful candidate. Failed historical attempts remain archived/inactive;
no new failed runtime variant was discarded in this cycle.
