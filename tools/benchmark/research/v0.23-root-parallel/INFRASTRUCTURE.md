# Retained infrastructure failures

Initial source replay runs all package tests successfully but exits with a Go
cache-trim permission failure. Its log remains. A new replay directory with the
writable declared GOCACHE passes; no engine source changes between replays.

Initial fresh collection stops before queries because a completed analysis
record has no moves. The matching archived PGN for ZsXkGPEs is an abandoned,
zero-ply game. The exclusion reader now records an empty completed game as zero
positions, while still rejecting incomplete empty analysis and malformed inputs.
Regression and race controls cover both cases. No surveyed board is omitted.

The original plan, failed workflow and method-only output remain unchanged.
The corrected collection uses a new plan/run/output directory, the same frozen
seed 2026112113/stream113 and the same exclusions. This is pre-query input repair;
no candidate result or generated cohort has been inspected or resampled.
Production source, binary and passed cost/native observations are unchanged.
