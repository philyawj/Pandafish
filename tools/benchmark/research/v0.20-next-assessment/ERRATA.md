# Preserved initial findings

The downloaded Go 1.25 toolchain lacks the pprof tool. CPU capture succeeded;
the installed Homebrew Go pprof tool rendered the saved profile. This changed
only profile presentation, not the engine compiler or timing results.

An added direct-checking en-passant mate control exposed a cold/warm API
disagreement in quiet-proof v1. The cold generator could classify that tactical
move as a quiet mate, while the cached-list path correctly excluded it. A
discovered-check EP mate did not expose the issue because the quiet geometry
filter already excluded that move; both controls remain in the final tests.
The direct-control failure, complete failing source and all v1 choice probes
remain archived. No v1 strength match ran.

V2 explicitly excludes en passant from the quiet proof. Ordinary tactical
quiescence still finds both EP mates. This is a local API/control repair, not
a material, clock, table or guessed-score change. V1 timed evidence remains
labelled as v1; final acceptance requires fresh output directories for v2.
The frozen development inputs can serve as regression controls after this
correctness repair; they are not relabelled as a new blind holdout.

## Reachable-history review

The first direct EP fixture retained a pawn on d7 despite declaring an EP target
of d6. It is accepted as a synthetic FEN by the rules parser, but cannot describe
the immediately preceding legal d7-d5 history. Its original failed log and source
are preserved; it is not presented as a reachable-history proof.

The separately saved `legal_ep_regression_test.go.txt` starts from
`1rrb4/1pkp4/8/2K1P3/8/7B/8/8 b - - 0 1`, legally plays d7d5 and reaches
`1rrb4/1pk5/8/2KpP3/8/7B/8/8 w - d6 0 2`. Original full legal enumeration proves
e5d6 is an EP checkmate. The cold quiet API still misclassifies that tactical mate
in v1 (expected failing log `quiet-v1-reachable-ep.txt`); corrected v2 passes
(`quiet-v2-reachable-ep.txt`). This confirms the API defect on a legal predecessor
without relying on the initial synthetic history. Both are scratch restorations
of the original overlays plus the same additional control.

V1's native evidence remains labelled v1. No fresh v2 timed acceptance or
strength screen follows: the independent Black-turn decision failure already
rejects this standalone hypothesis, and the classification repair does not
address that failure. The earlier statement that acceptance needs new v2 output
directories describes a necessary condition for any reconsideration, not a
claim that such a rerun occurred. No quiet proof is active in the presented
legality-kernel candidate.
