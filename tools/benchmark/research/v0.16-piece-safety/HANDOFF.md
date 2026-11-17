# Piece-safety result and research release

The narrow v0.16 piece-safety candidate is rejected and archived in
[candidate-source.json](candidate-source.json), over tagged v0.15.1. It is not
compiled by normal builds. v0.15.2 is a research patch that changes public
identity and documentation, retaining v0.15 playing behavior.

Read [README.md](README.md), [assessment.json](assessment.json),
[candidate-checks.json](candidate-checks.json) and
[early-failures.md](early-failures.md). All raw diagnostics and source drivers
are local to this directory. The final balanced timing schedule was ABBA;
there was no 50-game strength screen because the mechanism/cost gate failed.
Release smoke games are integration checks only and stay separate.

The motivating Nd4/Qd6 quiet leaves receive no safety correction. A simple
geometric queen penalty would even misread Nf6's absolute pin. Correct pins and
recaptures pass focused controls, but evaluation costs about three times as
much on the microbenchmarks; equal-time node count falls 18.2%. Decision changes
include real improvements and serious regressions. No Elo claim follows.

I authorized committing, merging and publishing this work as
v0.15.2, then asked for a **fresh choice of the next hypothesis**, rather than
following the previous five recommendations as a checklist. After publication,
start the next experiment from the tagged v0.15.2 source and executable. The
pawn ending after Rd5 is a leading diagnostic, but reassess all alternatives
against the new result before selecting implementation scope. Preserve rejected
findings and stop for my review before preparing that next PR.

Keep material values, search, quiescence, table policy and clocks fixed for a
new evaluation experiment. Run positive/negative controls, public regressions,
fixed-depth and equal-time checks first. If promising, use the standard 50 games
each against tagged v0.15.2, Stockfish Skill 2 and pinned Sacrifice, three
matchups concurrent. Keep all failures. Do not run larger confirmation or live
Lichess games automatically. No next feature is implemented in this release.
