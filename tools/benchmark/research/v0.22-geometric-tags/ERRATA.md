# Pre-test initialization correction

Source inspection catches a lookup-table dependency in the first geometry
initializer before any candidate test, timing or native build is run. Its
`bbForSquare` calls read `bbSquares`, populated by an original package init
function that runs after variable initializers. The initial source is preserved
as `search_tags.initial.go.txt`; no measurements are claimed for it.

Replace those calls with integer shifts using the pinned A1-at-bit-63 mapping.
Add explicit geometry checks against independently computed alignment and live
lookup masks after initialization. This repeats the class of initialization
hazard previously documented for v0.18 and reinforces the need for direct table
controls. Preserve the caught defect; only the corrected source reaches tests.
