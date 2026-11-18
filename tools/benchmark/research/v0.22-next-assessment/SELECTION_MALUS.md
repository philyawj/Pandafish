# Revision: distinguish unsuccessful quiet moves

The positive-only version preserves all 325 scores/root moves, but fails timing:
3.40% median public gain, approximately 9.5% deep and 10.0% newer. Keep its exact
source/results inactive as a separate [archive](../v0.22-quiet-history/README.md).

After a completed main-search **quiet** beta cutoff, reward that move as before
and penalize earlier searched quiet moves at the same node with the same depth-
squared magnitude. Ignore earlier captures/promotions/EP. A tactical cutoff
causes no history update. Use `h += bonus - h*abs(bonus)/16384`, keeping both signs
bounded; the reward cap and lifetime/indices stay unchanged. No negative update
comes from an incomplete child. This is a revision of ordering feedback, not
score penalties, selective pruning, offline training or a different evaluator.

Add relative-success, signed saturation and tactical/incomplete negative
controls. Retain the 15% median elapsed target, full score/best-set proof,
warm-game controls, all native gates and the standard concurrent screen. Every
production move remains eligible. No other failed code is bundled, and no larger
confirmation or publication follows automatically.
