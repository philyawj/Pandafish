# v0.17 legal tactical generation and move storage

I approved promotion of the screened candidate as v0.17 and authorized its
PR, merge and release. The broader confirmation remains deferred. The earlier
recommendation and development-state notes below are retained as the original
assessment, superseded by this explicit release decision. I also authorized
reassessing and starting one next candidate after release, with review before
its next PR. No live games are authorized.

**Recommendation: promote this candidate under the existing practical gate,
with provisional strength claims.**
I selected one bounded search-backend experiment after publishing v0.16 and
[reassessing the two next projects](ASSESSMENT.md). No new commit, PR or release
has been prepared for this candidate.

The hypothesis is that filtering tactical destinations before legal-move tagging,
testing move existence without a list, and reusing per-ply move storage provide
substantially more useful work at equal time. The evaluator, search policy,
quiescence eligibility/horizons, tables, history, mate scores and clocks are
unchanged. A local pinned-library root-package copy adds the search API; every
original runtime file and the MIT license is byte-preserved. [METHOD.md](METHOD.md)
describes the implementation, controls and limitations.

All existing tests, race suite, vet, native/ARM64 builds and offline wrapper pass.
The original upstream library tests pass under race. New fixtures/500 seeded
positions match the unchanged generator, including tactical order and state/hash.
The independent python-chess oracle passes all **596 positions**, every legal
move/tag/child FEN/hash and depth-two perft. Starting depth-four perft passes.

All **192 fixed-depth pairs retain exact node counts, scores and chosen moves**.
This measures lower cost for the same search, before equal-time decisions change.

| Probe | Positions | Same nodes, each actor | Baseline summed time | Candidate summed time |
| --- | ---: | ---: | ---: | ---: |
| Depth 4 | 84 | 1,010,100 | 4.891 s | 2.779 s |
| Depth 6 | 12 | 3,658,933 | 15.795 s | 9.579 s |

That is approximately 1.76x and 1.65x fixed-work throughput on these different
cohorts, with one run per case/actor. The same 60-request timed profile uses
**71.1 MB → 9.54 MB allocated/request**, and **154,761 → 73,499 allocations**.
Both use unchanged time management; they may complete different amounts of
work. These are strong functional/efficiency findings, not measured strength.

All **2,046 native UCI requests** complete as expected (2,018 nonterminal
searches and 28 terminal null responses). The requested budgets and unchanged
engine margins match.

| Equal-time block | Paired requests | Baseline mean depth | Candidate mean depth | Deeper / equal / shallower |
| --- | ---: | ---: | ---: | ---: |
| Serial, 100 ms | 252 | 4.702 | 5.115 | 104 / 148 / 0 |
| Serial, 1,000 ms | 84 | 6.405 | 6.738 | 28 / 56 / 0 |
| Three-search load, 100 ms | 504 | 4.631 | 5.073 | 224 / 279 / 1 |

The one shallower pair is `known-repetition`, trial 5, depth 8 versus 7.
It remains a negative timing control; its cause is not isolated. The other five
candidate repetitions complete depth 8. Its depth-3/4 score/node/move parity
also passes. I do not claim a per-position or universal speedup. Maximum
observed native responses stay within requested time plus the 50 ms protocol
allowance (about 125.3 ms baseline and 120.4 ms candidate under shared load).

The standard screen completed all three matchups concurrently, 50 games each,
100 ms requested moves and the unchanged 50 ms protocol allowance:

| Opponent | W / D / L | Unfinished | Forfeits | Completed score |
| --- | --- | ---: | ---: | ---: |
| Tagged v0.16 | 24 / 11 / 15 | 0 | 0 | 59% |
| Stockfish Skill 2 | 25 / 3 / 22 | 0 | 0 | 53% |
| Pinned Sacrifice | 22 / 9 / 19 | 0 | 0 | 53% |

The practical direct gate passes. All 150 PGNs pass independent legal-move,
result and checkmate checks. The run took about 10.7 minutes; maximum candidate
response was about 98.2 ms. The original labels, hashes, every game/loss and
[full settings](../../results/v0.17-search-moves-dev-vs-v0.16-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/summary.json)
remain intact. [screen-review.json](screen-review.json) includes opening/color
breakdowns and nominal paired intervals; [the result report](../../results/v0.17-search-moves-dev-vs-v0.16-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
links all PGNs. Every result and failure stays. No larger screen,
live client or further feature is started automatically. I will check in before
committing, pushing or preparing another PR.

[Complete source overlay](candidate-source.json), [source hashes](source-hashes.json),
[upstream hashes/license provenance](upstream-source.json), [code checks](checks.json),
[independent oracle](oracle-results.json), [profiles/checks](measurement-checks.json),
[fixed-depth analysis](analysis.json), [preserved initial failure](ERRATA.md).

## Interpretation and next step

I recommend promotion based on exact fixed-work equivalence, materially lower
cost/allocation, deeper native work and a positive practical screen with no
reliability failures. This recommendation uses the existing screening policy;
it does not turn a finite screen into a confirmed Elo gain.

The direct 59% is against the stronger, newly released v0.16 baseline. Its
nominal paired 95% t interval is about 47.5–70.5%, and repeated openings weaken
the independence assumption. The direct color split is **44% White / 74% Black**.
Both external 53% samples also have wide uncertainty. I cannot claim a reliably
above-50% rate against either opponent or advance the Stockfish target from
this screen. Historical v0.16 external scores are context, not simultaneous
baseline controls. Functional improvement is established within the controls;
measured strength evidence is promising and provisional.

A separately planned broader confirmation is justified if the next decision
requires stronger strength evidence. It is not automatic or a prerequisite
silently added after the practical gate passes. I have not run it. The main
implementation tradeoff is maintaining a local pinned root-package copy; hashes,
upstream tests, independent oracles and the retained license make the change
reviewable. The immutable position model remains; this is not a full board rewrite.

[Next-step reassessment](NEXT.md) now favors the coherent data-calibrated
evaluation/threat workstream over expanding into a full board rewrite. The
bounded backend route met its engineering target, while known quiet-threat and
endgame weaknesses remain. No third experiment has started. The candidate and
all new evidence remain local and uncommitted, pending review before a PR.

## Approved v0.17 release preparation

I approved promotion from the standard 50-game screen, deferred the proposed
larger confirmation, and authorized a separately reassessed next candidate
after release. [release-source-identity.json](release-source-identity.json)
checks that only the UCI name/test differ from the frozen source. The repeated
race/vet/build/native-ARM64/offline-wrapper checks pass.

The six-game integration smoke is preserved separately: two draws against
v0.16, two losses against Skill 2, and two wins against Sacrifice. All six games
are independently legal, with zero forfeits/unfinished games. Its strict
two-game previous-version gate returns exit 1 at 50%; this is integration-only
evidence and was not rerun for a favorable score. The original 50-game strength
screen is unchanged. [Release checks](release-checks.json),
[provenance](release-provenance.json) and [PGN check](release-pgn-check.json)
retain the distinction.
