# Explicit wall-time and value protocol

The three alternating 1,145-root depth-five pairs use two candidate workers
versus the tagged single-worker baseline. This intentionally spends additional
CPU under the same wall-clock/depth workload. Record actual node ratios; do not
require scheduling-deterministic counts or substitute this profile into any old
failed deterministic-tree report.

The shared `--parallel` profile requires exact baseline values, complete depth,
trajectory identity and a legal move in every pair. Candidate tied moves and
nodes may vary; baseline results must remain deterministic. Median search
elapsed ratio must be <=0.85. The original exact-tree mode remains unchanged.

After all timed pairs, a separate serial full-window pass checks each distinct
recorded selected move on every root with its actual repetition history. Its
value must equal the recorded optimum. This reference may restrict a selected
move to verify it; production roots remain unrestricted. Verification runs after
timing, not between measured roots, and its durations are not speed evidence.
Source/race/cancellation/worker-bound tests and later native/fresh/game controls
remain mandatory. No strength claim follows from this cost gate alone.
