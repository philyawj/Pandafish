# Repaired linkage, derived report

This directory contains the output of the v0.15.1 reporting tool run against
a temporary copy of the original raw session. The archived original
[summary](../../raw/summary.md) and [moves](../../raw/moves.jsonl) remain intact.

| Measure | Original | Repaired |
| --- | --- | --- |
| Completed searches | 12,091 | 12,091 |
| Uniquely linked searches | 11,615 | 12,021 |
| Unresolved searches | 476 | 70 |
| Illegal searches / missing final metrics / fallbacks | 0 / 0 / 0 | 0 / 0 / 0 |

The 406 recovered links agree with the independent conservative audit. The
matching rules still require recorded history, root FEN and the actual next
PGN move. An unmatched final request no longer erases earlier valid links.
Consistent later evidence may resolve common opening prefixes; contradictory
evidence does not justify guessing a game.

Sixty unresolved requests belong to two identical 30-move games that cannot
be distinguished by these histories. Nine others have roots matching saved
final positions but returned moves absent from the PGN. One more remains
unmatched. They are retained with missing/ambiguous status, not silently
attributed by date or presumed submission. There are 12,021 uniquely covered
bot moves and 61 of the 12,082 PGN bot moves without unique captured metadata.

Stockfish's full survey evaluates actual PGN moves independently of this
linkage. Its evaluations are unchanged; this repair improves capture
attribution and future reports, not chess strength.

See [PROCESS.md](../../PROCESS.md) for the command and tests.
