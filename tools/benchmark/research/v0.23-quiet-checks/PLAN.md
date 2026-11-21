# Checking-only generation with one bounded quiescence layer

Choose the pre-live checking-interface thread after two documented rejections:
the live-inspired contextual evaluator fails newer 100ms decisions, and current-
kernel child reuse saves 9.14% warm versus its frozen 15% target. Neither is
active. No old failed quiet-mate proof, evaluator or storage layer is reenabled.

Implement an additive library API for legal noncapture, nonpromotion checking
moves. Prefilter direct checking destinations and discovered-check blockers
before full legality/tag computation; include checking castles. The final legal
tagger remains authoritative. Promotions and EP/captures stay in the existing
tactical generator. Malformed kings or a pre-existing attack on the enemy king
use a conservative unfiltered fallback, rather than silently losing moves.

First prove this API equals the original complete legal list filtered by quiet
check on frozen independent-rule positions, special discoveries/castles/pins,
both turns, checked positions and seeded legal walks. Cover cold/cached calls,
caller-owned buffers and immutable parent state. Benchmark versus generating
all legal quiets then filtering; require at least 50% aggregate generator saving
on the fixed public/newer board set before activating the search change.

At an unchecked main-search horizon only (initial quiescence depth sixteen),
search the additional quiet checks after the existing tactical moves. Every
required reply uses normal legal check evasions and normal capture continuation.
No further quiet checking layer is added after a reply or capture. Stand-pat
and its valid fail-high remain unchanged; scores only rise over that baseline
at the initial horizon. This is selective tactical expansion, not a complete
mate proof or a solution for every opponent threat. Search termination,
repetition context, mate distance, material, static evaluation, TT policy,
root tie rules and time management remain unchanged.

Use immediate quiet mates and defensive captures/blocks/escapes as positive and
negative controls for both colors; retain legal quiet promotions, discovered
checks and history/cancellation cases. No production root is forced. Assess
actual equal-time decisions on all four existing suites with their unchanged
tolerances and losing-mate guards. Freeze this one-layer trigger before probes;
do not widen it after a failed outcome. A surviving candidate receives fresh
seeded final controls and one standard concurrent 50-game-per-opponent screen.
Preserve all failures, stop before committing a viable candidate, and do not
automatically run a larger confirmation.
