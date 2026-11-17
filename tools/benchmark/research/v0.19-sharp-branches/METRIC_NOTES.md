# What the mate-aware objective means

I keep finite score regression separate from mate ordering. I do not turn
`mate -4` into a numerical pawn value or change the engine's exact terminal
mate/draw scores. The new ordinal term is a soft 200 cp score-gap target between
an alternative with a reported losing mate and one without that report.

A finite score can still describe a lost position. Here “avoidable mate risk”
means that another labelled alternative has no losing mate found at the fixed
reference budget; it does not prove the position can be saved. I use this as
a conservative survival diagnostic, not a game-theoretic outcome claim.

The primary selection loss averages capped decision penalties over roots with
at least two alternatives. Finite alternatives use the same 25 cp uncertainty
band and 1,000-point cap as before. A selected losing-mate alternative gets a
1,000-point diagnostic penalty when a non-losing-mate-reported alternative
exists. This penalty is an explicit selection convention, not its training
CP label. I report mate-risk counts and CP-only results separately.

If all labelled alternatives lose by mate, there is no avoidable-mate penalty
within that subset. Winning mates remain acceptable, but I do not require the
reference's particular winning line when finite alternatives might also be
sound. Positive mate labels and distances remain raw for later review.

Ties cannot hide an unsafe choice: a predicted-best tie that contains a losing
mate still counts as a risky root. The unchanged model remains a valid selection.
Validation and holdout reject newly risky roots relative to that baseline.
Terminal roots and single-choice roots remain in the data/control records,
outside the decision-mean denominator.

A terminal or repetition witness keeps its exact search score and has zero
learnable features. A checked-cap witness retains the features the unchanged
engine actually evaluates at that cap. Thus I change supervision and selection,
not search eligibility, mate handling, material, transposition policy or clocks.

These are frozen-student-leaf surrogates. Any selected nonzero model still needs
unrestricted equal-time search checks, both-turn gates, mate controls and the
standard concurrent strength screen before a promotion recommendation. Lower
surrogate penalty is not playing-strength evidence.
