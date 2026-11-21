# Child-state reuse on v0.22: reject the larger-gain attempt

All 1,145 warm roots in each of three alternating pairs retain identical moves,
scores, depths and nodes. Full normal/race tests pass, including the archived
immutable-transition, ownership, aliasing, perft and cancellation controls.
The scratch baseline's 71 Go files match tagged v0.22 byte-for-byte.

| Warm depth-five pair | Candidate / baseline search time | Saving |
| --- | --- | --- |
| 1 | 0.919631 | 8.04% |
| 2, reverse order | 0.906810 | 9.32% |
| 3 | 0.908627 | 9.14% |

The median **9.137%** elapsed saving misses the prospectively frozen **15%**
warm target. The driver exits one with `timing_gate_passed:false`; no cold
325-case suite, native choice tests, fresh final controls or games follow.
Source and the prepared but unexecuted cold driver are retained. Restore all
playing files to v0.22 before the next hypothesis.

This is a real engineering saving, larger than the old v0.21 warm result of
4.65%, but insufficient for the stated larger-gain objective. Neither failed
attempt is relabeled as passed. Profiles show allocation volume and overlapping
runtime costs, not a promise that removing bytes removes proportional elapsed
time. Reuse remains a possible independently scoped memory/efficiency project;
it is not bundled into a later candidate or tested repeatedly for a preferred
strength score.

The [source overlay](candidate-source.json), exact commands, every paired row,
baseline source verification, tests/race output and analysis remain available.
Go reproduction uses `candidate-replay-source.json` with the normal research
`replay` command and base v0.22. The original snapshot remains intact;
`source-format-repair.json` verifies that omitted nested module metadata is
unchanged from that base, with no runtime source change.
The warm driver uses compiled test binaries without CPU profiling or per-root
MemStats. The earlier profile remains distinct in the assessment directory.
No measured playing-strength claim, commit or publication follows.
