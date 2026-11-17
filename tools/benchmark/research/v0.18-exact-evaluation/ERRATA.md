# Preserved checks and corrections

I invoked `analyse.py` before the sequential validation process had written
`candidate-deep.json`. It returned exit 1 with `FileNotFoundError` for that
not-yet-produced file. This was an ordering error in my analysis invocation,
not a failed engine probe. I retain the partial probe outputs and let the
original validation finish; I do not rerun its measurements. The completed
analysis is run only after the required files exist.

The first final-audit invocation checked the README's link to its own output
`final-audit.json` before creating that file. The bootstrap link check failed,
with the original report/trace retained as `final-audit-initial.*`. I corrected
the audit to defer its known generated-output link. No source, probe, model
weight or game result changed. The completed audit remains a separate output.
