# Exact removal of a redundant ring expression

The initial public-100 diagnostic passes: 41.164 cp versus 45.896 cp baseline,
severe cases one versus two, no new losing-mate reports or timing overruns.
Its source, binary and complete result remain the initial implementation.

Inspection also identifies a structurally zero expression: `ring & ownAttacks
&^ enemyAttacks` is empty because the defending king geometrically attacks its
entire own ring. Therefore this first representation has no active weak-ring
feature. Do not claim that the intended weak-square weight contributed to the
result. No new weak-square definition or coefficient is substituted after seeing
that result. Safe checking destinations still use defenders and double attacks.

Remove only that zero expression, preserving every scalar, and test the compact
implementation as the actual candidate before further selection. Original and
compact source/binary identities remain distinct. The old public-100 result is
diagnostic; run all four standard actual-choice suites for the compact build.
No parameter search or adverse-gate exception is authorized by this cleanup.
The initial component-cost sample overlapped brief build/snapshot work at
startup and is provisional; retain it and run an isolated cost sample.
