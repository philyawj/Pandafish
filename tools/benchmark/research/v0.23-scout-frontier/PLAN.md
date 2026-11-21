# Protect the last two main-search plies

The phase-guarded predecessor has a substantial nominal-depth cost gain and low
mean full-depth selected regret, but fails both fresh gates. Its one reported
mate concern survives a one-million-node check; both algorithms select that bad
move at the same deeper depth, so no unique reduction cause is asserted. The
remaining long-budget suite then fails its White tolerance. Preserve the
conditional scope decision and all false gates; no matches run for that binary.

Restore v0.22 before this separately frozen revision. Keep the general low-
material protection and full-depth verification, but permit a reduced scout
only at parent depth >=4, protecting both final main-search plies. At root depth
four the entire search is unreduced. This addresses the broad tactical risk of
skipping a quiet reply immediately above quiescence; it is not a root/FEN
exception or a claim to solve the predecessor's horizon-sensitive mate case.
No coefficient fitting, new evaluator, quiet-check expansion, cache change,
parallel search or old small kernel is bundled.

Require guard/verification/context/race controls, the unchanged independent mate
assertion and exact selected values/co-best membership through root depth four.
Keep the same 20% early node-saving and 15% three-pair median elapsed targets.
Validate through the shared Go validation plan and source-pinned warm preflight.
All four existing and both fresh gates must pass for the ordinary route; a
failure is documented before any separately considered scope reassessment.
Final seed 2026112120/stream120 excludes the now-consumed phase-guarded cohort
as well as earlier controls/live boards. Reuse the original-tag full-depth
regret review before one standard concurrent screen. No automatic larger
confirmation/publication, and check in before committing a viable candidate.
