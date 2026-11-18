# Quiet-mate proof: archived, not promoted

The nonrecursive proof finds the 15 saved quiet mates and passes enumeration,
discovered-check, castling, mate-distance, repetition and cancellation controls.
The first depth-four suite costs 8.33% more elapsed time, with a changed tree
(991,089 versus 1,010,100 counted search nodes). This is fixed-depth cost, not
an identical-work speed comparison. Public 100ms/1000ms actual-choice aggregate
gates pass without new losing-mate reports or response overruns.

Fresh development controls contain 145 roots from 16 baseline/reference rollouts,
73 White to move and 72 Black to move. At 100ms, 130 fully finite cases give:

| Diagnostic | Baseline | Quiet proof v1 |
| --- | ---: | ---: |
| Aggregate mean loss, 25cp band and 1000 cap | 40.4154 | 42.3769 |
| White mean loss | 42.7692 | 39.7692 |
| Black mean loss | 38.0615 | 44.9846 |
| Severe positions, loss >300 | 4 | 5 |

The CLI's aggregate gate passes, but the additional predeclared fresh-turn gate
fails: Black exceeds baseline + max(5cp,10%). Overall acceptance therefore fails.
The absence of new mate-risk flags and deadline overruns does not override that
failure. No strength matches ran. No 1000ms fresh expansion or favorable retry
followed. These finite-reference diagnostics do not measure playing strength.

A direct en-passant control also caught a cold/warm API classification defect.
The corrected v2 excludes tactical EP and still finds EP mates through ordinary
quiescence. Full race tests, vet and ARM64 build pass. That repair does not target
the observed equal-time regression, so it does not justify rerunning the same
candidate for a better score. Both original and corrected source overlays,
the failing control, all successful controls and every native/reference query
are preserved. The original raw reports retain their original aggregate status.

Playing source is restored to tagged v0.19 before the next implementation.
The proof is a useful archived functional mechanism; its standalone promotion
is rejected. It may be reconsidered only as a separately declared future test
with a credible cost/selection change and new acceptance evidence.

## Reassessment

The new baseline CPU sample puts 42.08% cumulatively in move tagging and 23.08%
in board update, including 11.31% directly in repeated bitboard setters. These
percentages overlap and must not be added. The current board update loops over
all twelve piece bitboards even for an ordinary quiet move.

The next experiment will specialize **immutable board transitions** for search:
update only moving/captured/promoted/castling/en-passant pieces, keep original
public methods untouched as oracles, and require exact legal tags, child state,
hashes, scores, root moves and node counts. This is a stronger profile-based
target than guessed weights or a full make/unmake rewrite. Quiet-mate policy,
attack-ray changes and unused-tag skipping will not be bundled into it.
