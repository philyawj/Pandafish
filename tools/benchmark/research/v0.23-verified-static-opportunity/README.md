# Reject: stand-pat does not verify a static cutoff

The unchanged1,145-root/43,118,858-node tree matches released v0.22 exactly.
The net proxy is18.308% after charging1,077,982 verification nodes; both-turn
headroom passes. However the same179 false predictions and one concealed losing
mate remain, so the mandatory safety gate fails. No pruning is implemented.

The proposed verification is ineffective by construction: outside check,
quiescence may stand pat. Every initial static prerequisite already reaches the
verification threshold, so every q probe returns immediately after one node.
It cannot disprove the static cutoff. This design mistake should have been
identified before the run; the23-second diagnosis is preserved rather than
represented as meaningful tactical verification. No threshold adjustment can
repair that underlying issue. Close this static-pruning line, restore released
play and reassess. No native/fresh tests or games run, and no new cohort is used.

The source snapshot, full179 counterexamples and failed workflow remain. This
is the23rd unsuccessful substantive attempt (18 prototypes, one offline model,
four diagnoses), not a23rd strength screen. No publication/live play is authorized.
