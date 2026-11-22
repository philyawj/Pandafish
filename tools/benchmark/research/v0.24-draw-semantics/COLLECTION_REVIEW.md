# Collection identity review

The collector completes all16 families and retains144 roots,72 per turn, using
the single frozen seed2026112202/stream125. No collection query is repeated.
The wrapper records `identity-failed` because its historical exclusion glob
also becomes a match for the newly written cohort's own `inputs.json` after the
step. That output did not exist when exclusions were resolved. The original
plan, log and failed state remain intact; they are not relabelled passed.

The explicit Go review verifies every one of523 original input identities,
including the teacher symlink identity. All are byte-identical. The sole new
glob match is exactly the expected collection output. It verifies cohort/family
coverage and pins the output hash. The initial verifier treated the wrapper's
reserved symlink identity as a filename and failed; its source is retained,
and the corrected verifier implements the wrapper's link-plus-resolved-path
identity. This repair does not collect data or change any candidate/gate.

The source/data identity check is therefore satisfied by the explicit review,
while the original wrapper's self-output error stays visible. Fresh native work
uses the completed output pinned by its own shared plan. Future collection
plans should enumerate pre-existing resolved exclusion files rather than use
an input glob that can include the planned output. No scientific threshold,
cohort or seed changes, and no favorable retry, occur.
