# Bitset computation of the existing exchange heuristic

The [plan](PLAN.md) selects repeated attacker scanning from the measured profile.
The generation foundation is the exact staged picker plus evasion prefilter;
their earlier rejected variants keep their own source and results. The new
change is confined to the exchange board and its attacker computation. Full
search, evaluation values, eligibility, table policy and clocks are unchanged.

The former local mailbox exchange board scans all 64 squares for each least
attacker and scans again to test king safety after a temporary capture. The new
board retains the same piece mailbox and king pointers, adding side/type bitsets.
It initializes them from the existing value-copy `EvaluationPieces` API and
iterates occupied squares. Missing king pointers retain the former A1 default;
nonunique king pointers retain the last ascending-square king. This preserves
the old heuristic even on the supplied malformed-king controls.

Pawn/knight/king origin masks directly identify their geometric attackers.
Precomputed rays intersect current occupancy; only the nearest occupied square
on each ray can supply a slider. Combined attacker bits are restricted to the
requested side. Geometry initializes from integer squares, is immutable and
occupies about 6.5 KiB. It is not a position/score cache or a new TT policy.

Least-attacker selection visits candidate bits in **ascending square order**,
exactly as the old mailbox scan. It keeps the original material costs, king-last
cost and strict lower-cost replacement. This preserves equal-value ties even
between bishops and knights. Temporary captures update mailbox, type/side sets
and the same king pointer before testing attacks on that king. Restoration
checks cover every field, including all new bitsets. Pins and newly opened
lines therefore remain legal-recapture checks, not geometric approximations.

The initial EP removal, promotion gain, recapturing-pawn queen promotion, bounded
gain array, side alternation and backward stop-or-exchange minimax remain the
original code. All calls to cancellation/deadline checks remain at the same
logical points. The old heuristic remains a least-valuable-legal-attacker
estimate; the rewrite does not make it a full tactical exchange search.

`reference_exchange_test.go` freezes the original v0.21 mailbox implementation
under test-only names. Production excludes that reference. No learned weights,
new capture filtering rule or root-specific exception is introduced. Twelve
frozen Go overlays include the generation foundation, new kernel and oracles.
All original upstream runtime files remain byte-identical.
