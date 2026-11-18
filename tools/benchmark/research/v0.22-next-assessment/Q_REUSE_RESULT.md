# Quiescence reuse opportunities rejected before implementation

Both diagnostics complete the same 1,145 roots / 43,118,858 nodes. The strict
remaining-depth simulation has 4,023,265 applicable hits and 2,249,602 potentially
avoidable descendants (5.22%). The independently scoped horizon-certificate
simulation raises those to 5,651,188 / 3,151,843 (7.31%), including 1,627,923 cross-
depth hits. It observes 521 horizon caps and no bound contradiction. Raw traces,
proof-tracker source, counters and both failed >20% gates remain preserved.

Neither result justifies adding the extra score table or proof machinery under
the stated opportunity criterion. No playing cache is implemented. The certificate
diagnosis is not an unsafe depth>= shortcut and does not establish a global proof
from finite observations; its mathematical condition and limits remain in the
[separate note](Q_CERTIFICATE_DIAGNOSIS.md). Hypothetical descendant stores/bound
returns differ from an actual cache, so these node potentials are not measured
elapsed savings. No instrumented timing or allocation is promoted as efficiency.

Low duplicate-work coverage, rather than a demonstrated cache defect, limits
this route. Do not remove draw context, mode or horizon rules to manufacture a
larger number. Reassess work that can be avoided before creating a child, such
as a guaranteed existing quiescence stand-pat cutoff, rather than caching a
small share of already short leaves.
