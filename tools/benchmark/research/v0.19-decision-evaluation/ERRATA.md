## v0.19 tooling migration

This historical report retains its original measurements. Python/shell runners
were retired in v0.19; commands below describe the original experiment.
Current verification uses the [Go research tools](/tools/research/README.md).
Original scripts and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# Preserved collection correction

A generated rollout reached checkmate at a sampled ply. The finite-reference
query correctly returned mate zero and no bestmove. The collector incorrectly
required a legal bestmove even for terminal roots, and stopped. I preserve the
raw query, original source and complete traceback as `collect-initial.py.txt`
and `collection-run-initial.txt`. I retain the root as terminal/excluded rather
than treating it as a missing favorable label.

`resume-plan.json` freezes the completed root IDs and original membership hash.
I correct terminal handling and resume only unfinished branch labels, without
regenerating families, changing splits or repeating completed reference queries.
The model has not been fitted or activated.

The post-result stricter independence audit finds early transpositions that
the declared ply-eight grouping did not cover. Prefix strings can differ
while reaching the same position at plies four/five. Two positions connect
training/validation families16/19/21; three positions from two families relate
to reserved confirmation openings. The sampled roots are distinct, but the
family/confirmation-independence claim is weaker than intended. I preserve
that limitation, add synthetic grouping controls, and retain the failed fit
without resplitting/refitting. Future data must check canonical early positions
before collection/labels, with a fresh holdout.
