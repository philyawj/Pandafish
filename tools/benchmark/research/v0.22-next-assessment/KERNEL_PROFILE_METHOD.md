# Post-screen kernel profile commands and identity

The staged and evasion profiles are single diagnostic observations, not paired
efficiency evidence. Their workspaces replay the respective frozen candidate
JSON over v0.21 and add `warm_profile_test.go.txt` under a test filename. The
staged workspace is `/tmp/pandafish-v022-picker-profile`; the evasion workspace
is `/tmp/pandafish-v022-evasion-profile`. Engine files match the phase's frozen
source. Extra gated reporting tests in a scratch module do not execute under
the selected test name.

For each phase, run from that scratch module with absolute repository paths:

```sh
GOCACHE=/tmp/pandafish-v016-go-cache \
WARM_INPUT=/Users/philyawj/projects/Pandafish/tools/benchmark/research/v0.21-rank-order/screen/v0.21-rank-order-dev-vs-v0.20-and-stockfish-3-and-sacrifice-v0.2.0-b528f0ba \
WARM_OUTPUT=/absolute/path/to/new-profile-rows.json \
WARM_MODE=pipeline WARM_DEPTH=5 \
go test -count=1 -timeout=15m -run '^TestWarmPipelineProfile$' -v \
  -cpuprofile /absolute/path/to/new-profile.cpu
```

The original phase-specific output names are `staged-profile-warm-depth5.json`
and `evasion-profile-warm-depth5.json`, with `staged-warm-depth5.cpu` and
`evasion-warm-depth5.cpu` in this assessment directory. Matching text logs retain
all 20 game checkpoints. The same frozen input hashes, 1,145 roots and depth-five
results must agree with the original warm baseline. The lazy and evasion
diagnostic replays use their separately frozen reporting helper and
`Q_DIAGNOSTIC_OUTPUT`; production scoring remains unrestricted and unchanged.

The Go 1.25 toolchain lacks its pprof executable in this environment. Rendering
uses the installed local Go tool without changing the engine compiler:

```sh
GOTOOLCHAIN=local /opt/homebrew/opt/go/libexec/bin/go tool pprof \
  -top -cum -nodefraction=0 -nodecount=300 \
  /tmp/pandafish-v022-evasion-profile/Pandafish.test \
  tools/benchmark/research/v0.22-next-assessment/evasion-warm-depth5.cpu
```

Top/cumulative/all-node reports stay beside the raw profiles. Setup, per-root
MemStats, GC, cancellation timers and output add runtime overhead; percentages
overlap and cannot be added or converted directly to expected wall-time gains.
No heavy validation or second profile overlaps these observations. Independent
paired fixed-work measurements, not these profiles, decide elapsed cost.
