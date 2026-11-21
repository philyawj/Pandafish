# Protect low-material search in the child-scout revision

The initial child-scout reduction policy fails its independent KQK delay-mate
control and remains archived, including interrupted timing. Reassess the failure
before a new binary. Low-material quiet moves can be essential to forced mates,
zugzwang and conversion; a reduced finite fail-low can hide that continuation,
so mate-band verification alone is insufficient.

Add one general scope guard: require total nonpawn phase >=8 under the released
phase weights (start 24), otherwise search every scout at full depth. This is a
structural low-material protection, not a FEN/move exception or tuned coefficient.
Keep the existing depth/move/check/killer/mate-window guards, one-ply reduction,
full-depth verification of alpha improvements/mate results and ordinary PVS
re-search. No other engine mechanism is introduced. The earlier KQK proof keeps
its exact original assertion. The killer-only audit uses an explicitly unreduced
reference in both actors to preserve that isolated ordering test's original scope.

The initial partial timing is exploratory only. Freeze new source/binaries and
require new normal/race/mate/context controls before any measurement. Use a
fail-closed Go validation workflow; the warm plan must pin its passed state and
matching source identities. No multiline test/build/timing chain can hide a test
failure. This tooling correction does not alter the engine hypothesis.

Targets remain >=20% node saving in each warm pair and >=15% median elapsed
saving over three alternating full-game pairs. Fixed-depth values may change;
legal completed results and deterministic candidate values/nodes are required.
All four existing native and both new final suites must pass; no surrogate
exception is permitted for this selective policy. Seed 2026112119/stream119,
with consumed-control/live exclusions. Independently inspect selected-move
full-depth regret before one standard concurrent 150-game screen. Preserve
all results, do not run larger confirmation automatically and check in before
committing a viable candidate. No publication or live play is authorized.
