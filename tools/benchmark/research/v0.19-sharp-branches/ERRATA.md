# Preserved preflight finding

The first combined unittest invocation used the chess-client virtual environment,
which has python-chess but no NumPy. Boundary/query controls passed, while the
optimizer tests could not import NumPy. I retain preflight-tests.txt. I run the
boundary/query tests with that environment and the optimizer tests with the
existing NumPy-capable python3; no dependency install or data collection occurred.

The first collection stopped on a legal 25ms v0.18 fallback with no completed
scored depth. The logger wrongly required an exact score from timed student
rollouts. I preserve its source, checkpoint and traceback. Teacher node-limited
queries still require exact results; timed student queries require a valid move.
The resume manifest freezes the recorded query bytes and completed rollouts. I
replay that prefix (including the original fallback) without querying engines,
then append only unfinished work. No data/split/fit result was selected or rerun.

Before any fit or validation/holdout quality metric, I found a CP-pair target
scaling error in the declared experimental objective. Clipping the teacher
gap before subtracting the original base gap could ask the learner to shrink
an already-correct large material margin, or move a correction in the wrong
direction. FIT_REFINEMENT.json freezes the algebraic correction and examples.
I clip the residual correction instead, retaining the original plan and old
model source for provenance. No empirical fit/holdout result informed this
change; all other bounds, grid, ordinal constraints and gates stay fixed.
This affects experimental fitting only; v0.18 contains no learned residual.

The first final audit compared the declared commit with the annotated tag
object from `git rev-parse v0.18` and stopped. I preserve its source and
exit/diagnosis. I correct the audit to dereference `v0.18^{commit}`; this
changes no experimental source, observation, input or acceptance gate.
