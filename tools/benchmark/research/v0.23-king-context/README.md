# Learned piece–king model: reject before runtime activation

The different 450-coefficient representation reuses the balanced sharp corpus,
with all 4,412 learnable FEN/base links verified against original evaluation.
Feature orientation, analytic optimum, gradient/projection, rounding, mate-order
and race controls pass. Old validation/holdout are inspected development data,
not newly unseen evidence. All source and prior labels remain unchanged.

The initial solver fails its numerical gate before development evaluation.
[The numerical refinement](NUMERICAL_REFINEMENT.md) preserves that failure and
keeps every model/data/quality parameter fixed. The first converged solution
then finishes in 177 steps. It uses 5,265 finite and 152 ordinal constraints,
including 28 constant constraints; it does not drop whole mate-containing roots.

| Inspected split | Baseline mean penalty | Model mean penalty | Severe baseline→model |
| --- | ---: | ---: | ---: |
| Training, 494 roots | 51.150 | 43.004 | 18→14 |
| Validation, 182 roots | 36.095 | 37.462 | 3→3 |
| Historical holdout, 184 roots | 44.043 | 56.179 | 3→6 |

The historical holdout also adds the avoidable-loss-reported choice
`sharp-028-ply36`. Both development gates fail. The source/weights/complete
metrics remain in `fit-numerical`; the shared runner stops automatically.
No runtime context table, cost/native study, fresh seed or strength game runs.
The production engine remains released v0.22. No ridge, cap, subset or feature
bucket is tuned against this result.

This rejects the fitted model, not every contextual evaluator. The supervised
targets are reference values of root alternatives, while features describe a
weak student's later witness board. Those are not the same value problem.
Their mismatch is a candidate explanation, not an established cause. The next
assessment first measures that mismatch using fixed training cases before
authorizing a larger direct-leaf labeling effort. Previous static-fit failures
and finite-label/partial-alternative limitations remain relevant.

[Infrastructure notes](INFRASTRUCTURE.md) retain both document-map preflight
failures and explicit repairs. The complete source snapshots also include the
Go toolbox, so the new fitter is replayable without retired scripts.
