# Processing, selection and preservation

## Existing evidence

All work was offline. The original recordings and three processed archives
were read without regenerating their reports or modifying their contents.
The combined review is a new sibling directory, not a replacement session.
No Lichess API or external analysis service was queried.

The [source fingerprints](source-fingerprints.json) cover **525 existing archive
files**, including raw records, inventories, analysis, incident lists and reports.
Every one of the **341 raw inventory entries** was checked against SHA-256 and
length and compared with its original `local-games/` counterpart. Executables
remain local/ignored but are included in the integrity check. All source
fingerprints were checked again after the review. No raw record, unsuccessful
result, ambiguous link or original selection was deleted or rewritten.

Session manifests carry different configuration hashes because the recording
paths differ. A textual diff of the redacted configurations shows only `engine.dir`
and `pgn_directory` changes. All engine hashes, source revisions, UCI versions,
client revisions and declared hardware/workload match. The exact dirty external
client patch is unavailable; its revision alone cannot reproduce that state.

## Deduplication and legality

PGNs were parsed independently of the saved operational totals. A game key is
the eight-character Lichess ID from its Site tag. The 157 PGN occurrences have
157 distinct IDs. Repeated opponents and repeated openings are retained;
Cloudia-Chess's simultaneous games are different IDs, not duplicates.

The audit checks overlapping duplicate histories for consistency and would keep
the longest compatible record with all source paths; no duplicate-resolution
case was needed. Conflicting metadata or move prefixes would stop the audit
rather than silently select a result. Original paths are listed in [games.json](games.json).

Every per-game analysis file is complete, matches its PGN SHA-256 and covers
exactly the recorded plies. Each analyzed move and FEN was compared with the
PGN reconstruction. Referee PVs were replayed legally, including complete
500,000-node confirmation PVs. Bot-side classifications were recomputed.
All 7,507 completed captured best moves were separately validated as legal.

Search overlap uses actual `go` and `bestmove` trace wall timestamps. It does
not infer real wall duration from monotonic engine elapsed time: host suspension
produced a large discrepancy in the known time-forfeit. Twenty cross-process
overlap pairs belong to the two concurrent Cloudia-Chess games. One
wall-minus-monotonic gap exceeds one second, in the reported closure search.
The detailed intervals and clock gap remain in [audit.json](audit.json).

## Referee method

The complete survey is reused from each source archive's
`analysis/stockfish-n100000/`. All three methods specify full-strength
**Stockfish 19**, SHA-256
`dc2f18c34ae962dff591b66147d220ec06e61d756b93d8a4f5e04fd8e55c251f`,
100,000 nodes per query, one thread, 64 MiB hash, Skill 20 and MultiPV 1.
Played-root and unrestricted queries preserve the full supplied history and
start with fresh hash and `ucinewgame`. Scores favor the root mover.
CP and mate comparisons remain separate. Negative CP deltas remain recorded
as search disagreement, with only the derived CP loss clamped to zero.

The combined selected-position follow-up uses the same executable/settings
at **500,000 nodes per query**, through the preserved production Stockfish
adapter. That adapter rechecks a selected reference when final bestmove and
the score PV disagree. A played-root query is reused only when the reference
already selects the played move. It never forces a Pandafish root choice.
There was no Pandafish replay, policy modification or heuristic experiment.

## Selection before larger-budget queries

The [primary selection](selection.json) contains 32 distinct games:

- Eight cases in each of opening (zero-based ply below 20), middlegame
  (20–59) and late (60 or later). Require survey reference CP in [−100,+300]
  and loss at least 150 cp, or a newly allowed mate from that reference band.
  Within each group, sort by earliest ply, larger finite loss, then game ID.
  Select at most one case per game globally, taking the phase groups in that
  order. Phase names are reproducible ply buckets, not a material classifier.
- Four winning-mate survey flags, ordered by earliest ply and game ID,
  excluding already selected games.
- Two negative-delta controls with raw delta at most −100 cp and reference
  at least −100 cp, then two played/reference-agreement controls with reference
  at least −100 cp and ply at least 12. Use the same earliest-ply ordering
  and global one-game limit.

The initial selection draft admitted already-winning reference CP above +300.
It was narrowed **before any larger-budget query** to emphasize contested
decisions. That unqueried draft is retained as
[selection-draft-before-winning-filter.json](selection-draft-before-winning-filter.json).
The final selection SHA-256 is bound into
[confirmation-method.json](confirmation-method.json). No new positions were
added after inspecting its results. Every selected result is preserved in
[confirmation-n500000.json](confirmation-n500000.json).

A separate playable newly-allowed-mate filter was checked against all survey
bot moves, independently of selection order. It finds no reference CP at least
−100. Its saved [empty selection](selection-mate-defence.json) causes no queries.
[survey-mate-transitions.json](survey-mate-transitions.json) preserves all 61
survey mate-transition flags and scores, so the empty result is inspectable.
Mate-defence commands are present for reproducibility but were not run.

## Go provenance and reproduction

[review.go.txt](provenance/review.go.txt) is Go source, with its helper
[stockfish.go.txt](provenance/stockfish.go.txt) copied unchanged from
`tools/lichess/stockfish.go`. The `.go.txt` extension isolates research entry
points from normal package discovery. No other-language program was added.
From the repository root, with all original sessions and source archives present:

```sh
mkdir -p /tmp/pandafish-v022-review
cp tools/benchmark/live/v0.22-campaign-review-20261120/provenance/review.go.txt /tmp/pandafish-v022-review/review.go
cp tools/benchmark/live/v0.22-campaign-review-20261120/provenance/stockfish.go.txt /tmp/pandafish-v022-review/stockfish.go
GOCACHE=/tmp/pandafish-v016-go-cache go run /tmp/pandafish-v022-review/*.go
GOCACHE=/tmp/pandafish-v016-go-cache go run /tmp/pandafish-v022-review/*.go confirm
GOCACHE=/tmp/pandafish-v016-go-cache go run /tmp/pandafish-v022-review/*.go summarise
```

The audit refuses changed source fingerprints. Confirmation checkpoints after
each case; an existing completed file resumes without repeating queries and
checks the saved selection prefix and frozen method. Keep
its method and input hashes together. For a different selection, executable or
budget, use a separate output profile rather than overwriting this evidence.
The standard raw archive checks also work without the ignored executable;
this stricter local reproduction additionally requires every original file.

After all review files are final, `go run /tmp/pandafish-v022-review/*.go inventory`
writes their separate [derived inventory](derived-inventory.json), excluding
that inventory itself; `verify-inventory` checks it without rewriting it.
