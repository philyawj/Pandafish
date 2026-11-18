# Candidate review checkpoint

Recommend promotion of the search-legality kernel, pending approval before any
commit. No commit, push, PR, version bump, release or live Lichess game has been
performed for this investigation. The active branch remains `main` at tagged
v0.19 with local changes. Development labels are preserved, not released names.

Start with [results](RESULTS.md), [implementation](IMPLEMENTATION.md),
[testing](TESTING.md) and the [combined assessment](../../../../docs/RESEARCH_V020.md).
The median selected fixed-work elapsed reduction is 38.7%, with all 325 distinct
score/move/node comparisons exact. All four native quality suites pass, including
both newer root turns. The routine concurrent 50-game screens score:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| Tagged v0.19 | 28 / 9 / 13 | 0 | 0 | 65% |
| Stockfish Skill 2 | 33 / 1 / 16 | 0 | 0 | 67% |
| Pinned Sacrifice | 21 / 16 / 13 | 0 | 0 | 58% |

Full root race/vet/tests, independent mailbox attacks, transition/perft controls,
native and Linux/ARM64 builds, restored-candidate tests/build and all 150 legal
PGN replays pass. Original upstream runtime/module/license bytes are unchanged.
`candidate-source.json` and `source-hashes.json` identify the tested Go overlay;
the active source and binary hashes still agree. The unused upstream rendering
dependency failure remains in the original log, alongside the scoped root pass.

Two preceding failures remain preserved. The
[quiet-mate proof](../v0.20-next-assessment/QUIET_RESULT.md) fails its fresh
Black-turn gate; v1's EP defect is reproduced from a legal predecessor and repaired
in archived v2. **No quiet proof is active.** The
[transition-only revision](../v0.20-search-transitions/RESULT.md) gains about 21%
fixed-work time but fails its direct screen at 45%. Its entire source and all
150 games remain. The current revision adds measured attack/tag savings under a
new extra-gain threshold; earlier failures are not relabelled as successful tests.

Fixed-tree efficiency is firmer evidence than the exact size/generalization of
the strength gain. Repeated openings, scheduling and selection after failed
attempts remain limitations. Prior external scores are not simultaneous baseline
controls. A broader-opening confirmation is justified before a large/general
strength claim or target advancement; the optional 100 games/opponent over 50
openings has not run and is not a new mandatory promotion gate.

The next action is the requested review checkpoint. After approval, follow
[RELEASING.md](../../../../RELEASING.md) for any publication, preserving these
tested development identities and distinguishing future release smoke checks.
Do not rerun games for a preferred score. Reprofile the resulting engine before
selecting another backend change; quiet-mate reconsideration or a new evaluation
representation needs a separately declared design and new evidence. All future
programs, controls and research drivers remain Go.
