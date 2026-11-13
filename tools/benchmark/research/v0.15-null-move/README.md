# Guarded null-move pruning: experiment and findings

Published in **v0.14.3 as a research archive**. The release retains v0.14.2's
playing behavior; experimental search code and tests are archived data outside
the normal Go build.

The assessment below was made after publishing v0.14.2, before match results.
The positional-evaluation candidate is archived, not enabled. The previously
suggested pawn/king-only variant remains a reasonable narrower hypothesis, but
fresh reference comparison favors testing a larger search opportunity first.
No numerical strength gain is guaranteed.

## Current primary-source comparison

| Engine and pinned source | Approach relevant here | Implication for minimal Go Pandafish |
| --- | --- | --- |
| [Stockfish](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/search.cpp) | Null-move pruning with material, score, node-context guards and deeper verification; extensive other selective search. Its [evaluation](https://github.com/official-stockfish/Stockfish/blob/def601cc156f7d41f9ddea76d74515eddf232aaf/src/evaluate.cpp) uses neural networks. | Borrow the search question and safeguards, not its parameters or full architecture. NNUE is a much larger project. |
| [Sunfish](https://github.com/thomasahle/sunfish/blob/e396948acbd80cdcbad3dccfb07e878fba9df8f4/sunfish.py) | Compact zero-window bound search, positional tables, guarded null moves; current code also has deep null probes, reductions, repetition and mate-distance handling. | A minimal implementation can still benefit from pruning. Its current null policy is more nuanced than a universal pass cutoff. Do not describe it as just random play plus tables. |
| [Sacrifice](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/search.rs) | PVS, late quiet-move reductions, killers/history ordering, depth-qualified transposition bounds and mate-distance normalization. [Evaluation](https://github.com/piratesboys/Sacrifice-Engine-/blob/b528f0ba373910b72052d991700a990ca1cb07fd/src/evaluate.rs) combines placement, mobility, shelter, activity, center and attacks. No null-move mechanism found in this search file. | Its strength cannot be assigned to one feature. Our PVS/reduction experiment already failed its desired strength goal. Copying the complete feature set is not a controlled next step. |

The default-branch heads fetched for this review match these revisions. Sacrifice
remains the same pinned opponent; no silent benchmark upgrade. Source inspection,
not playing-strength measurement of these reference engines, supports the table.
This is original Go code, not copied implementations or evaluation tables.

## Ranked next hypotheses

1. **Guarded null-move pruning.** Pandafish currently searches real continuations
   before obtaining any such bound. A cheap pass probe can cut branches that
   remain adequate even after giving the opponent a free turn. This may save
   enough work to complete deeper useful iterations. Chosen for broad upside.
2. **Pawn/king-only evaluation.** Cheaper than the failed combined evaluator and
   encouraging on selected positions. Still unproven in games; static cost and
   overvalued passers/shelter remain concerns.
3. **Mate distance and draw-aware conversion.** Constant mate scores and absent
   repetition scoring can harm conversion. This requires careful score-table
   normalization and history context, not a cosmetic evaluation bonus.
4. **History-based quiet ordering.** Complements killers and table moves, but
   likely starts as an efficiency improvement rather than the large gain sought.
5. **Richer attack/development evaluation.** Potentially helpful but larger and
   more expensive; the last screen argues against bundling more guesses now.

## Fixed candidate policy

Keep released evaluation, quiescence, root ties, table layout, time allocation,
and real-move ordering. No PVS/LMR, databases, new UCI options, or new dependencies.
Before searching real moves at a nonterminal node, permit a pass probe only if:

- depth is at least three and the side is not in check;
- the side has a knight, bishop, rook or queen;
- the beta bound is outside the mate band and static evaluation reaches beta;
- the immediately preceding move was not a hypothetical pass.

Search the passed position at depth minus three, with a one-point window around
beta. Return only beta on a non-mate fail-high, not an exact or mate score. The
pass uses the existing chess library: turn changes, en passant clears, board and
castling rights remain, and the real parent is immutable. The pass subtree has
no access to the real transposition table or killer hints. Cancelled work cannot
become a completed score. The internal reference flag disables null pruning for
offline controls only and is never changed within a live engine/table.

Pawn-only exclusion does not eliminate all zugzwang: minor-piece endings can
still violate the assumption. This modest policy is a hypothesis, not a proof
of minimax equivalence. No high-depth verification mechanism is claimed.

## Validation and decision before scores

Run focused guards/state/bounds/cancellation tests, full tests, race checks, vet,
ARM64 build and offline wrapper, plus saved-position depth-four and 100 ms
controls. Preserve score differences rather than requiring a selective search
to be identical at every position. Never infer strength from node reduction.

Use 50 games per opponent against tagged v0.14.2, Stockfish Skill 2 and pinned
Sacrifice, with all three matchups concurrent. This is the standard routine
confirmed for this session. Review the previous-release criterion after play and
preserve all results. Historical 60%/34% external scores provide context, not a
simultaneous baseline control. The descriptive aim remains a substantial visible gain, not merely a narrow gate
pass. No automatic retries, weight sweep or longer matches to rescue a result.
Check in before any candidate commit or PR.

## Pre-match diagnostics

On the same 19 selected saved/synthetic positions used in the last experiment,
fixed-depth-four scores were unchanged in all cases. Nodes fell from 387,832 to
279,122 (28.0%). Separate 100 ms iterative probes reached higher depth in four
cases, the same in fifteen, and lower in none. Two timed scores differed. Total
visited nodes were 367,825 for the reference and 405,059 for the candidate;
visiting more nodes under a clock can coexist with less work at fixed depth.
These selected cases, random root ties, and one timed sample per policy/case
do not establish game strength. Raw rows and exact source are preserved in the
[results directory](../../results/v0.15-null-dev-vs-v0.14.2-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md). No policy or parameter was changed after these diagnostics.

## Completed screen: hold

[All results and interpretation](../../results/v0.15-null-dev-vs-v0.14.2-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md).
The candidate scored 56% against v0.14.2, 36% against Skill 2, and 23% against
Sacrifice, with no forfeits or unfinished games in 150 total games. The direct
criterion passed, but external evidence does not support promotion. The candidate
was not promoted. Next, diagnose saved-game failures and test
conversion gaps instead of assuming another textbook feature yields strength.

## Representative code

A null move is a search probe, never a move sent through UCI:

```go
probe := *worker
probe.table, probe.killers = nil, nil
score, complete := probe.searchNode(
    position.Update(nil), depth-3, -beta, 1-beta, false, ply+1, false,
)
worker.nodes = probe.nodes
```

Cancellation remains distinct from a score. A non-mate fail-high returns the
beta bound only. The final false argument prevents an immediately consecutive
pass; descendants after a real move can consider a pass again. The full source
and focused tests are in `candidate-source.json`, not the released engine.

## Replay the archived candidate

Run from the repository root:

```sh
repo_root="$PWD"
replay_root=$(mktemp -d)
record_dir="$repo_root/tools/benchmark/results/v0.15-null-dev-vs-v0.14.2-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba"
git archive v0.14.2 | tar -x -C "$replay_root"
cp tools/benchmark/research/v0.15-selective-search/restore-snapshot.go.txt /tmp/pandafish-restore-snapshot.go
go run /tmp/pandafish-restore-snapshot.go "$record_dir/candidate-source.json" "$replay_root"
(cd "$replay_root" && go test ./...)
(cd "$replay_root" && go build -o /tmp/pandafish-null-replay .)
(cd "$replay_root" && PANDAFISH_NULL_AUDIT=/tmp/null-audit.json go test -run '^TestNullAudit$' -count=1 .)
```

Use fresh tables/processes for controls. Compiler versions, scheduling and random
root ties can change binary hashes, timings and outcomes. The archived overlay
recovers the tested policy, not guaranteed identical match results.

## Research-release verification

[Release checks](research-release-checks.json) separately record restored-engine
normal/race tests, vet, ARM64 build, the offline wrapper, archived source replay,
and six short integration games. Those smoke scores do not establish strength.
The release restores search byte-for-byte from v0.14.2 and changes only its UCI
identity/handshake expectation to v0.14.3. All experimental labels and game hashes
remain unchanged; none are presented as measurements of the release binary.
