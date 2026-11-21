# Fresh decision controls, frozen before collection

If the frozen evaluation survives existing actual-choice suites, generate the
fixed 16-family rollout with seed 2026112023 and tagged v0.22 only. Candidate
code and outcomes do not choose source moves, positions or family count. Use
the preserved Go generator, otherwise unchanged from the v0.20 collector.
Sample adjacent White/Black roots at plies 12/13, 24/25, 36/37, 48/49 and 60/61.
Retain all attempted families, query failures, exclusions and histories.

Exclude canonical/color-mirrored old v0.19 blocked/forest states, public and
newer decision roots, every v0.22 live surveyed board and duplicate samples.
Require at least 80 retained roots and at least 30 for each turn. These controls
are unseen before the candidate is frozen, but share collection machinery and
finite referee estimates with older work; they are not an independent Elo test.
No fitting uses them. After inspection they are development evidence for future
attempts, not reusable hidden holdouts.

Run 100ms (three trials) and 1000ms (one), both turn gates, with fresh 100,000-node
reference queries and the existing max(10%,5cp) finite-loss tolerance, severe-error
and separate losing-mate controls. No extension of a disappointing batch or
favorable rerun. A new candidate after failure needs an explicitly frozen new
final-control seed rather than a claim that this inspected set remains hidden.
