# Search-ordering reassessment before implementation

The four exact-kernel variants preserve correctness but fail their 15% median
elapsed target; their source and failed gates remain inactive archives. Smaller
measured gains alone do not justify promoting another code path in this study.

The older architecture assessment lists history ordering, aspiration and verified
reductions after backend improvements, but warns that its shallow first-move
cutoff rate was already about 78%. Re-measure on v0.21 before treating history as
an established large opportunity. Instrument main beta cutoffs only, recording
first-move and quiet cutoff indices on the existing public/newer depth-four
cohorts. The instrumented baseline remains diagnostic scratch source and never
becomes a playing candidate. No timing claim uses instrumented observations.

If at least 5% of main cutoffs are quiet and occur after the first move on both
cohorts, test one bounded per-search quiet-history ordering mechanism. This is
an opportunity screen rather than proof the heuristic will select those moves.
Otherwise reconsider aspiration/another search hypothesis. No LMR, futility,
null move, new evaluation or lower clock-check frequency is bundled. Freeze the
actual ordering/update policy and correctness/native/screen gates before editing.
