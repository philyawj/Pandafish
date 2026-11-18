# Horizon certificates as a separate quiescence-cache hypothesis

The exact-remaining-depth simulated cache has 4,023,265 applicable hits across
38,446,985 quiescence visits, but avoids only 2,249,602 non-overlapping descendants.
Its complete warm node total/fraction and failed opportunity gate remain in the
strict diagnostic review. Do not remove depth or draw guards to inflate reuse.

A different, provable horizon condition may permit reuse across remaining-depth
values. Diagnose this separately before implementing a playing cache. A completed
quiescence proof that never touches a horizon cap can remain valid for any
remaining depth sufficient to traverse the visited proof tree. For each visited
nonterminal non-draw node at relative ply r, require remaining depth at least r+1
outside check and r-3 in check (the unchanged four-ply evasion cap). Terminal/draw
leaves impose no horizon requirement. If any visited node stops at a horizon,
mark the entire conservative certificate invalid. Otherwise propagate the maximum
required depth through recursive calls. Bound applicability, mate normalization,
check/unfiltered mode, halfmove and complete repetition context remain mandatory.

This is not an ordinary deeper-entry shortcut. A stable certificate certifies
that the visited bound proof never depends on the omitted horizon; it is accepted
only above its recorded minimum. Capped work still requires exact depth. Using
all visited branches is conservative: even an earlier capped branch invalidates
a later cutoff's certificate. No selective-search result or evaluation change
is introduced.

Repeat only the instrumented opportunity diagnosis for this different mechanism,
still recomputing every result. Record exact/cross-depth hits, invalid caps and
observed bound contradictions against recomputed results. Require zero such
contradictions and the same >20% non-overlapping descendant opportunity before
implementation. An instrumented potential is not a cache speed claim. Future
code would additionally need independent proof/cap/cold-warm/depth/history/mate/
cancellation tests, production timed performance and the standard native/screen
gates. Do not claim the strict cache passed or rerun it for favorable evidence.
