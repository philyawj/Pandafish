# Full-opponent-ply null probe: reject on first cost pair

The [frozen policy](PLAN.md) protects the shallowest pass with a full opponent
move, adds material/pawn/repetition guards and verifies apparent cutoffs at
depth six or more. It is selectively heuristic, not full-width equivalence.
Ordinary quiescence, evaluation, table identity/capacity and clocks stay v0.22.

Full normal/race tests, focused both-color quiet-mate refutations, real-history
restoration, exact probe/deeper-verification node accounting, cancellation inside
the probe, vet and a clean source replay pass. The shared first warm pair covers
all 1,145 complete legal roots. Baseline 43,118,858 nodes becomes 40,958,134:
**5.01% fewer**, below the frozen 20% requirement. The one elapsed ratio is
0.958118 (4.19% faster), not a three-pair median performance finding.

The shared runner automatically stops at the early node gate, retaining its
false report and every row/log. Two further timing pairs, native controls, fresh
seed collection and games never run. No margin, guard or reduction is retuned.
The sparse gain does not justify this policy's tactical/zugzwang risk or a
promotion recommendation. Restore released playing files and reassess the
contextual-model project using preserved labels where suitable.

Full source SHA256:
`e4d9d85cef9b93d76b998131cd18d35824918c77a9b653929d7fe95647165964`.
The shared optional harness is research-only; source and complete controls are
retained in the snapshot and `provenance`. No strength conclusion follows.
