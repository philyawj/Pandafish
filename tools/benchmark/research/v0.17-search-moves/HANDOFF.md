# Start here: v0.17 search-backend candidate awaiting review

I approved promotion of the screened candidate as v0.17 and authorized its
PR, merge and release. The broader confirmation remains deferred. The earlier
recommendation and development-state notes below are retained as the original
assessment, superseded by this explicit release decision. I also authorized
reassessing and starting one next candidate after release, with review before
its next PR. No live games are authorized.

I published v0.16 through [PR #31](https://github.com/philyawj/Pandafish/pull/31)
and [the release](https://github.com/philyawj/Pandafish/releases/tag/v0.16),
then reassessed and implemented this separately authorized candidate. The
v0.16 confirmation was deferred. All current candidate changes are on local
`v0.17-search-moves`, uncommitted/unpushed, with no second PR prepared.

Read README.md, ASSESSMENT.md, METHOD.md, NEXT.md and ERRATA.md. Start from tag
v0.16 (`4ae65b14f69d95a1ad498baa759fec5e6bd0a1a6`) for restoration. The frozen
candidate overlay includes a local replacement of pinned corentings/chess
v2.6.1's root-package runtime and MIT license. Every original runtime file is
byte-identical; only the new search API is additive. Existing APIs/board/hash
and parser behavior are unchanged. The engine uses direct legal tactical
moves, no-list move existence/status, and reusable per-ply move slices.

Evaluation/material, search/PVS/root policy, capture eligibility, quiescence
horizons, score/hint tables, killers, mate/repetition behavior and clocks remain
v0.16. No bishop-pair, piece-safety, new pruning, opening book, trained model or
internal make/unmake rewrite is included.

Validation passes: full race/vet/build/native-ARM64/offline-wrapper checks;
original upstream tests under race; 500 seeded library comparisons and special
controls; starting depth-four perft; 596 independent python-chess move/tag,
child-FEN/hash/status and depth-two-perft oracles; all 192 paired depth searches
with identical nodes/scores/chosen moves; 2,046 valid native UCI requests.
The 28 terminal nulls are expected. One known-repetition timing pair is shallower
and retained. Fixed-work time falls 39–43%; timed-kernel allocation is about
71.1 MB versus 9.54 MB/request. These are functional/efficiency findings.

The standard screen completes all three 50-game matchups concurrently:

- v0.16: 24/11/15, **59%**; practical gate passes.
- Stockfish Skill 2: 25/3/22, **53%**.
- Pinned Sacrifice: 22/9/19, **53%**.

Zero forfeits/unfinished games, all 150 PGNs independently legal. The direct
lead is 44% White / 74% Black, and its nominal paired interval includes equality.
The external samples do not establish reliably above-50% strength. No Elo claim.
All raw data, initial formatting-scope failure, negative timing control and
source/license/binary provenance are retained. The source overlay was replayed
in a fresh directory; use reproduce.py without games for restoration.

I recommend promotion under the existing practical gate, pending review before
another commit, push or PR. A broader confirmation is separately justified for
stronger strength evidence, never automatic. No further implementation is
started. The next likely project is coherent data-calibrated evaluation/threat
modeling, with fresh response data and runtime controls, rather than extending
the full-board rewrite. Reassess when the review/result changes the evidence.

The released tag, active bot executable/private config and live games remain
untouched by this experiment. Named candidate/baseline binaries stay ignored
under tools/benchmark/bin. Original development labels/hashes are immutable;
do not relabel the screen as a released v0.17 binary.
