# Evidence-review findings

The independent Go choice reviewer initially needed three corrections before
its complete successful run. All occur in the review utility, outside playing
code and the timed query driver:

1. Its first compilation omitted the pinned `Game.Move` options argument. The
   call was corrected to the actual API with nil options.
2. Its first comparison assumed every historical quality report contained
   explicit turn gates. Quiet v1 predates those reporting fields; its aggregate
   report passes while the separately frozen Black-turn criterion fails. The
   reviewer preserves the original schema, recomputes both turn means, and
   compares the original reported aggregate status. The quiet candidate remains
   rejected overall; no observation or acceptance criterion was changed.
3. Exact FEN text comparison exposed the documented historical convention for
   uncapturable en-passant squares: the old input omits them while the Go history
   replay retains the double-push target. The reviewer normalizes EP only when
   no legal EP move exists, preserving all other FEN fields and actual legal EP.
   The first example was `broader-previous-04-ply65`, saved `-`, replayed `h3`.

The corrected reviewer then agrees with every completed suite across quiet v1,
the failed transition-only candidate and the legality-kernel revision. It
independently parses exact raw UCI scores, follows original cached observations,
validates legal histories and selected moves, and recomputes aggregate/turn
losses, severe cases, mate flags and response gates. The teacher is a finite
reference, not a unique acceptable-move oracle. Legacy cache records preserve
structured iterations rather than original raw UCI lines; that boundary remains
explicit. The saved review source permits reproducing these checks without
starting engines or games.

`source-review-initial.json` checks all five archived source snapshots and their
hash maps, the active frozen overlay, original upstream files, protected playing
policy files, candidate/baseline/opponent binaries and reference-cache hashes.
The word “initial” means this review preceded final PGNs and archive inventories;
it is not a superseded failed source result. Final inventories and ordinary Go
audit checks cover the completed archive separately.

The first Markdown review found an extra separator cell in the two manually
written matchup tables (five headers, six separator cells). The tables were
corrected without changing any game data or values. The failed document review
is retained as `docs-review-initial.json`; the subsequent result is separate.
