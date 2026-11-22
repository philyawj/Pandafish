# Contextual representation: stop at the capacity gate

The frozen scratch representation uses 655,488 parameter bytes and 160-byte
state cells. Full/incremental accumulators and outputs agree across 54,191
transitions from 602 positions, including the 596 independent-vector roots,
2,334 captures, 33 castles, five EP moves and 32 promotions. Both color/rank
and file mirrors agree. This is a diagnostic parity result, not a generalization
or playing-strength result. No production engine source is changed.

The required capacity certificate fails before output-head fitting. The pair
`sharp-001-ply37`, good `b8a7` versus bad `b8a8`, has identical activations.
[The collision review](collision-review.json) establishes a stronger fact:
both saved legal witness paths arrive at exactly the same six-field FEN and
base score -205. Therefore no evaluator depending only on that position can
separate these two fixed-witness targets. Changing width or random seed cannot
resolve that contradiction. The old feature-capacity observation was not
sufficient to attribute this pair to inadequate board representation.

The frozen gate remains failed. No seed/width/subset/target is changed, and
no cost, new labels, model training, native queries or games follow. The other
15 strict constraints do not receive a completed feasibility certificate.
The failure does not reject neural or contextual evaluation in general, nor
prove that changed evaluation could not alter the original search paths.
A future learning project needs a target/selection contract that does not demand
different values for identical inputs. Another fit on these fixed witness
preferences is not justified by this diagnostic.

[Original plan](PLAN.md), test log, parity/collision records, Go sources and the
minimal tagged-source overlay preserve the attempt. The latter replays the
four scratch Go files over v0.23; current released playing source stays exact.
Full race/vet/build timing gates were not reached after this early failure.
The next bounded mechanism resolves move-count claim semantics and verifies
cache/clock boundaries, separate from material-rule or model retries. Cleanup
PR #42 remains merged with no release; all research remains uncommitted.
