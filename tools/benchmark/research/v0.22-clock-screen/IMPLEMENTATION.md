# Same deadline, cheaper clock observation

`searcher.stopped` retains its context and zero-deadline guards, then delegates
the deadline predicate to one short helper. Injected clocks keep the original
`!now().Before(deadline)` boundary. Ordinary deadlines use
`time.Until(deadline) <= 0`. For production Time values derived from time.Now,
Go 1.25's public API reads only the embedded monotonic clock. Time.Before already
compares that same monotonic component; the discarded wall-clock read is unused
by the comparison. Wall-only deadlines use the standard Time API fallback.

Every check site, allowance, soft/hard margin, context cancellation, root choice,
search horizon, quiescence rule, TT policy, repetition identity and evaluation
remains v0.21. The change introduces no polling interval, cached clock, timer flag
or additional thinking time. Soft iteration checks and final metrics retain the
existing clock helper. Depth-only searches without a deadline do not benefit.

The two-file snapshot is byte-identical to the earlier clock research source.
The rebuilt development executable is also byte-identical, SHA-256
`ee73a8f0252a7d873888bec56e9ea0a56348759b234f1ecfcc47f496ca41e98a`.
Its UCI greeting remains Pandafish v0.21. A new development label identifies the
separate screen objective; it does not imply a runtime revision or a release.

The original 15% gate remains failed in the earlier archive and in the identical
copied `timing-analysis.json`. This assessment follows an explicit prospective
[scope reassessment](../v0.22-next-assessment/REASSESSMENT_SMALL_SCOPE.md), with a
5% target for timed warm work before native/game checks. No original cold trial
is rerun or replaced. The lower percentage reflects a different implementation-
cost objective, not success of the original larger-gain hypothesis.

Warm pairs use the original frozen 20 games and 1,145 candidate-turn roots,
production iterative deepening, one TT per game and fresh workers per root, at
depth four plus a 30s movetime allowance. This exercises deadline reads without
reaching a timer boundary; untimed depth-only pairs omit the mechanism. Both
actors use the unchanged budget formula, root trajectory and history. The
original 100ms baseline profiles remain descriptive background, not paired
candidate efficiency evidence. Per-root MemStats/output work is outside the
reported search interval and remains instrumentation overhead.

The native four-suite gates still use actual 100/1000ms budgets and three-search
load, with full-strength Stockfish labels and unrestricted production roots.
The routine game opponent remains Skill 3. Functional timing, native decisions
and one exploratory game screen must be reported separately; none establishes
a precise Elo gain, absolute scheduling guarantee or automatic promotion.
