# Numerical failure before any development result

The first optimizer reaches its frozen 1,500-step ceiling and fails convergence.
`fit/model.json`, the original tool source snapshot and runner failure remain
unchanged. `fit/summary.json` records `development_evaluated:false`: neither
inspected development split has been scored, and no runtime model is active.

The objective is about 8.27 million. Curvature stays 512 through step 1,000,
then grows to 68,719,476,736 while the objective changes by only thousandths.
The absolute 1e-8 tangent-test allowance is too close to summed-loss round-off;
inertial steps then persist as the step size collapses. This is numerical
diagnosis, not evidence that the board representation succeeds or fails.

Keep every datum, feature, objective, ridge, bound, quantization rule, 1,500-step
ceiling and development/safety threshold fixed. Refine only numerical execution:
use a 1e-12 relative round-off allowance and momentum based on the known ridge
strong-convexity bound. Still require the 1e-7 projected-stationarity tolerance;
scale the projected step by curvature/ridge so inflated curvature cannot itself
declare convergence. Analytic constrained optima, finite-difference gradients,
projection/rounding bounds, zero residuals and a large-constant-loss analogue
remain pre-fit controls. A new source snapshot/output directory distinguishes
the numerical recovery from the preserved initial failure.

This is not a ridge/cap/grid change or a retry after observing development
quality. The first converged solution alone may reach those gates. A subsequent
quality failure remains a failure, without another fit to its outcomes.
