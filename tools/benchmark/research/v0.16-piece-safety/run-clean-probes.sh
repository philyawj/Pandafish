#!/bin/sh
set -eu
# Scratch dirs are prepared from tag plus saved overlay; do not run concurrently
# with engine checks, games, or the referee assessment.
research="$PWD/tools/benchmark/research/v0.16-piece-safety"
baseline_dir=$(cat /tmp/pandafish-v016-baseline-dir)
candidate_dir=$(cat /tmp/pandafish-v016-candidate-dir)
export GOCACHE=/tmp/pandafish-v016-go-cache
for label in baseline candidate candidate baseline; do
 case "$label" in baseline) source_dir="$baseline_dir";;candidate) source_dir="$candidate_dir";;esac
 trial=${trial:-0};trial=$((trial+1))
 PROBE_TIMED_ONLY=1 PROBE_INPUT="$research/probe-inputs.json" PROBE_OUTPUT="$research/clean-$trial-$label.json" go -C "$source_dir" test -count=1 -run TestResearchProbe -v -timeout 3m > "$research/clean-$trial-$label-check.txt" 2>&1
done
for label in baseline candidate; do
 case "$label" in baseline) source_dir="$baseline_dir";;candidate) source_dir="$candidate_dir";;esac
 PROBE_FIXED_ONLY=1 PROBE_LIVE=1 PROBE_INPUT="$research/live-fullroot-inputs.json" PROBE_OUTPUT="$research/$label-live-fullroot.json" go -C "$source_dir" test -count=1 -run TestResearchProbe -v -timeout 8m > "$research/$label-live-fullroot-check.txt" 2>&1
done
LEAF_INPUT="$research/baseline-traces-history-corrected.json" LEAF_OUTPUT="$research/leaf-safety.json" go -C "$candidate_dir" test -count=1 -run TestResearchLeafSafety -v > "$research/leaf-check.txt" 2>&1
