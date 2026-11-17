#!/bin/sh
set -eu
r="$PWD/tools/benchmark/research/v0.16-bishop-pair"
base=$(cat /tmp/pandafish-v0162-baseline-dir)
candidate=$(cat /tmp/pandafish-bishop-candidate-dir)
export GOCACHE=/tmp/pandafish-v016-go-cache
for label in baseline candidate; do
 case "$label" in baseline) tree="$base";;candidate) tree="$candidate";;esac
 PROBE_FIXED_ONLY=1 PROBE_INPUT="$r/probe-inputs.json" PROBE_OUTPUT="$r/$label-fixed.json" go -C "$tree" test -count=1 -run TestResearchProbe -timeout 8m -v > "$r/$label-fixed-check.txt" 2>&1
 go -C "$tree" test -run '^$' -bench BenchmarkResearchStatic -benchtime=300ms -count=3 > "$r/$label-cost.txt" 2>&1
done
trial=0
for label in baseline candidate candidate baseline; do
 case "$label" in baseline) tree="$base";;candidate) tree="$candidate";;esac
 trial=$((trial+1))
 PROBE_TIMED_ONLY=1 PROBE_INPUT="$r/probe-inputs.json" PROBE_OUTPUT="$r/clean-$trial-$label.json" go -C "$tree" test -count=1 -run TestResearchProbe -timeout 3m -v > "$r/clean-$trial-$label-check.txt" 2>&1
done
