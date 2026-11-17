# Start here: v0.16 coordinated search core

I approved promoting Stage A from the v0.15.3 architecture plan as v0.16 and
authorized the PR, merge and release. The proposed larger confirmation is
deferred. First-person documentation edits are included in this release.
The active bot executable/config and live games remain unchanged.

Read [README.md](README.md), [METHOD.md](METHOD.md), [NEXT.md](NEXT.md) and
[ERRATA.md](ERRATA.md). The source starts at tagged v0.15.3; the original
architecture prototype is still archived unchanged. The candidate has a
separate development UCI identity and named ignored executable:
`tools/benchmark/bin/pandafish-v0.16-search-core-dev`.

Completed roots return one proven best move, not every tie. PVS uses no depth
reductions. Legal board-only hints are separate from history-qualified scores.
Irreversible child descent trims only unreachable fingerprint history and
restores state on unwind. Quiescence sorts only copied eligible tactics but
still generates full legal moves. No evaluator/material, SEE, horizon, clock,
score-table capacity/replacement/exact-depth or killer-policy change is bundled.
The hint table adds about 1 MiB. No bishop-pair bonus is active.

Validation passes: full race suite, vet, native/ARM64 builds, offline wrapper,
archive-source replay, 192 paired exact-score/move-membership searches, and
2,046 valid native UCI requests. The 28 terminal null responses are expected.
All raw data, tests and initial harness/test-authoring failures are preserved.
The 84-position corpus and 12-position deep subset were frozen before results.
Depth-four work falls 14.8%; depth-six work falls 35.2%. Three-search equal-time
paired depths are deeper/equal/shallower in 150/354/0 cases. These are functional
and efficiency findings, separate from game strength.

The standard concurrent 50-game screen completes in about 11 minutes:

- v0.15.3: 27/11/12, **65%**; practical gate passes.
- Stockfish Skill 2: 21/1/28, **43%**.
- Pinned Sacrifice: 16/17/17, **49%**.

All 150 games complete, zero forfeits/unfinished; independent python-chess replay
validates every PGN move/result and checkmate termination. The direct lead is
50% White / 80% Black, which warrants broader color/opening checking. The
Sacrifice result is encouraging relative to the historical 27% released-policy
sample, but this is not a contemporaneous external control. No Elo claim.

My earlier recommendation was a frozen 200-game-per-opponent confirmation,
which remains proposed in NEXT.md. I chose promotion from the standard screen
instead, retaining all limitations. No larger run was performed. I also
explicitly authorized reassessing the next two projects and starting one new
candidate after the release. Reassess current v0.16 performance before choosing
specialized move generation/storage versus evaluation/threat modeling or a
better-supported alternative. Keep that next experiment separate, preserve all
results, and check in before another PR. No live games are authorized.

Provenance, exact commands, source/binary hashes, raw failures and replay tools
are linked in the main report. The candidate overlay restores from v0.15.3 in a
new directory. Preserve the original game labels/hashes; do not relabel this
screen as a released v0.16 binary. If code changes, this screen no longer
validates the changed candidate.
