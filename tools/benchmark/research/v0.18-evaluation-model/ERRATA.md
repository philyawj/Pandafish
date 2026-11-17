# Preserved prototype findings

The first two independent feature extractions fail 320 of 360 cases, while
phase/mobility/pressure and symmetry mostly agree. The new pawn masks were
initialized as package variables via `bbForSquare`, but that upstream lookup
is populated later in an `init()` function. The masks were therefore zero.
The original helper and every failed extraction/vector are retained as
`feature-prototype-initial.go.txt` and `*-initial*` artifacts. No weights were
fit and no production evaluator was activated from those vectors.

I corrected initialization to use direct square bit shifts, independent of
upstream init order, and strengthened exact passed-pawn and asymmetric-shelter
controls. The independent oracle is rerun before any calibration. The initial
cost sample is also retained and superseded by a corrected cost sample because
the feature implementation changed.

Before fitting any weights or inspecting validation/holdout quality, I refined
eligibility to quiescent-stable rows (baseline static and full-window v0.17
quiescence within 50 cp), to align this positional model with its runtime use.
`CALIBRATION_PLAN.json` records that decision separately; the original corpus
plan/hash and all played/reference-reply labels remain unchanged. Earlier
feature exports are retained, and the exporter now derives baseline evaluation
from the tag even if a future working tree activates the candidate.

The first active-model full test run flags five old material-approximation
assertions (queen capture, forced quiet evasions and an extended exchange).
Their legal/tactical outcomes and exact independent minimax scores are intact;
positional residuals exceed the old +/-90 cp tolerance. I replaced those
nonterminal approximations with exact unpruned, unordered minimax comparisons,
retaining terminal values, tactical capture identity and exact material units.
An old comment incorrectly named Kxh2 after Bh2; the legal capture is Rxh2+.
Original failure output is retained in `tests-active-initial.txt`.

A diagnostic-message edit also changed an unrelated check-test reference from
`tt.want` to undefined `want`. The resulting compilation errors are preserved
in `tests-scalar.txt` and `cost-fitted.txt`; the check assertion is restored.
This was a test-authoring error, not a playing-code failure. Exact integer
quality was recomputed for the frozen weights without retuning; a few
sub-centipawn aggregate differences from floating truncation do not change
the held-out gate result.

The first scratch-candidate setup stopped because tagged v0.17 already contains
`third_party`, while the prior backend harness expected that directory to be
new. The raw traceback is in `validation-run-initial.txt`; all completed code
checks and baseline probes remain intact. I changed the scratch copy to allow
existing directories and resumed only the unfinished candidate probes, without
rerunning or replacing earlier measurements.
