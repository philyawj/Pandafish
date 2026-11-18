# Exact geometric check-evasion prefilter

The [plan](PLAN.md) adds one new generation optimization to the functionally
verified but unpromoted staged picker. The first picker's source, failed native
gate and 47% direct screen remain unchanged in its archive. This revision needs
at least 10% further warm elapsed gain before spending another game screen.

When an uncached checked position requests moves or legal-move existence, the
generator computes actual geometric checkers from its existing piece bitboards.
Orthogonal/diagonal attack rays identify visible rook/bishop/queen checkers;
fixed king/knight masks and pawn-origin masks identify direct checkers. If exactly
one checker exists, a precomputed line admits destinations that capture it or
interpose. For a non-sliding checker, the line contributes no extra destinations.
If multiple checkers exist, ordinary nonking destinations have an empty mask.

King candidates are never filtered by this mask. Pawn destinations at the
supplied EP square also bypass it: EP removes a piece away from the destination
and can capture a checking pawn or block a line. Every retained move still passes
the existing full move-tag/own-king legality predicate. The filter proves only
that excluded ordinary destinations cannot evade the current check. It does not
replace pin, king-escape, EP or opponent-check validation.

Unchecked positions, missing/nonunique own kings and a check flag without a
matching geometric checker use the original candidate set. Complete cached
legal lists retain their original metadata path. Castling keeps the original
checked-position rule and malformed-right compatibility path from the first
picker. Original upstream source, move order, hash/transition semantics and
search/TT/clock policy do not change.

Line and pawn-origin arrays initialize from integer square shifts, independent
of the upstream square-bit array's later `init()`. The full line table occupies
32 KiB; pawn origins add 1.5 KiB. They are immutable geometry, not a position or
score cache. The private helper adds no exported library API. The diagnostic's
slower mailbox/ray implementation remains separate archived source; counters
and its exported reporting map are absent from the playing candidate.

Eight frozen Go overlays include the original picker, the modified authored
search generator, the private geometry helper and focused root/nested tests.
All 21 original upstream runtime files stay byte-identical. The extra complexity
is checker geometry and its conservative exceptions, with original legal methods
and independent frozen vectors retained as oracles.
