# Current v0.18 evaluation revision

I completed the user-authorized follow-up on branch `v0.18-evaluation-model`.
The current source is the isolated **exact v0.17 evaluation kernel**, with UCI
identity `Pandafish v0.18` for release; the screened candidate identity remains
`Pandafish v0.18-exact-evaluation-dev` in all original records. No learned residual or feature
extractor is in the current build. The old slow/fast fitted variants remain
fully reproducible in `../v0.18-evaluation-model/`; their historical handoff
predates this continuation. No old inventory entry was changed.

I recommend promotion of this kernel under the practical gate, not promotion
of the rejected model. I subsequently approved PR/merge/release as v0.18 and reassessment/start of
one next candidate after publication. The next candidate needs a check-in
before its PR. The prior baseline v0.17 remains tagged at
`76e73c8716d2ddcc826bdf349ff1c7a7ac71ecf3`; active bot/config are untouched.
No live games or automatic larger confirmation are authorized by this run.

Read README.md, PLAN.json, NEXT.md, LEARNING_REVISION.md and ERRATA.md. The new
kernel uses an additive value-copy bitboard getter and a 38,400-byte table of
original integer formulas. All 19,200 entries and 596 base-score positions/both
perspectives match the old routine; new stress/ownership/allocation controls
pass. Original v0.17 root tests are restored, changing only the identity test.
Full race/vet/native/ARM64/offline wrapper/history controls pass.

All 192 fixed-depth node/score/move pairs match. Depth-four time falls
2.785→2.450 s and depth-six 9.580→8.150 s (~12%/~15%). All 2,046 native
requests are valid, including 28 expected terminal nulls; equal-time nodes
increase roughly 16–18%. Three serial and one warm cross-actor depth regressions
remain in shallower-timing-pairs.json. Their cause is not isolated.

The frozen candidate completed the standard 50 games per opponent with three
matchups concurrent at 100 ms plus 50 ms wall allowance:

- Tagged v0.17: 24/9/17, 57%.
- Stockfish Skill 2: 36/1/13, 73%.
- Pinned Sacrifice: 21/9/20, 51%.

All 150 PGNs independently replay legally; zero forfeits/unfinished. Screen
exit 0. The direct practical gate passes, but its nominal pair interval is
42.25–71.75%. Its first two openings yield 12 wins in 12 games; equal opening
weighting descriptively gives 52.08%, not a replacement official score/gate.
Direct colors are 58% White / 56% Black. No causal external-score delta, Elo
or reliably large strength claim follows. A separate broader 100-game-per-
opponent/50-opening confirmation is justified before broad claims or raising
Skill 2; it has not run and is optional under the current policy.

The next larger recommendation is fresh independent response/leaf data and a
student-budget decision/ranking objective, with severe-error controls. Do not
retune the old fitted weights on inspected failures. Their bounds and static
MAE selection are diagnostic clues, not a coefficient fix. A compact neural
representation is still an alternative; board rewrite/pruning remain separate.

candidate-source.json plus reproduce.py restores the current overlay on tagged
v0.17 in a new directory, then tests/builds without games. This clean replay
passes. final-audit.json and inventory.json preserve final provenance and all
699 earlier inventory entries. The backend explanation requested earlier is
`docs/CHESS_SEARCH_BACKEND.md`, linked from the main README; PR #32 already
clarifies that v0.16 is before-state and v0.17 introduces the local backend.
Retain that note and both experiment archives through any eventual PR.

Release validation passes race/vet/native/ARM64/offline wrapper. The identity-
only release source proof and six-game integration smoke are preserved. The
smoke returns exit 1 because its direct two games split 1/0/1 (50%); Skill 2
also splits 1/0/1 and Sacrifice is 2/0/0. All six independently replay legally,
zero forfeits/unfinished; no favorable rerun. See release-provenance.json.
The release target is https://github.com/philyawj/Pandafish/releases/tag/v0.18 .
