# Execution environment

Apple M2, eight logical CPUs, macOS/darwin-arm64, Go1.25 toolchain and
`GOCACHE=/tmp/pandafish-v016-go-cache`. Before warm timing, `pmset -g batt`
reports AC power,100% charged. `uptime` after validation reports load averages
3.87/3.29/3.36. Temperature and unrelated process utilization are unmeasured;
process listing is unavailable under the sandbox. No deliberate builds, tests,
profiles or competing research measurements overlap warm/native/game timing.
Three paired elapsed ratios describe this workload, not a universal cost bound.
