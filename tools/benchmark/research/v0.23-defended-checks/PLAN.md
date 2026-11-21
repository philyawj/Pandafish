# Defender-aware checking-pressure hypothesis

Prospective first v0.23 candidate, based on the new depth/leaf diagnosis.
Keep material values, original placement, all search eligibility/horizons,
transposition policy, time allocation and tie rules unchanged.

The old failed king model counted weighted attacks into a king ring and
coordination, without defenders or safe checking destinations. This design
instead builds geometric attack unions and double-attack sets, estimates
defender-aware checking destinations for knights/sliders, and conditions local
pawn shelter on enemy attacking resources. These are evaluation hints, not
legal checking moves or mate proofs. Pins and movement-induced occupancy changes
are explicit approximation limits. No fitted weights or old proof is activated.

Freeze the initial hand-designed constants before native/reference results:
checking destinations (up to four each) weight knight 24, bishop 16, rook 24,
queen 12; weak ring squares weight four; ring attack units N/B two, R three,
Q four; coordination multiplier at most four. Missing near-file shelter costs
20, a one-rank-advanced shelter six, two ranks twelve; enemy material resource
weights Q four, R two, minors one, capped at twelve. King danger is capped at
600 cp and tapered by the existing phase. Queenless danger is quartered.

Require independent blocker/attack geometry, symmetry, terminal/material and
hot-path allocation controls. Record actual cost rather than inferring strength
from static preferences. Existing kernel exactness assertions must continue
checking the unchanged base separately from the new complete evaluation.

Run the public and newer actual-choice suites at 100ms and 1000ms, with the
existing max(10%,5cp) tolerance and separate losing-mate guards, balanced-turn
checks for newer roots. A clear adverse gate rejects this initial design before
games; preserve it without coefficient rescue. A surviving frozen design then
receives new independently seeded decision controls and one standard concurrent
50-game screen versus tagged v0.22, Skill 3 and pinned Sacrifice. Fresh controls
are evaluated once and are development evidence afterward. Preserve every
failure, source snapshot and result, then reassess if rejected. No publication
or larger confirmation runs automatically.
