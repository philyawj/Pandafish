# Phase-sensitive piece activity: reject

The frozen coordinate/phase map is not recommended for promotion. Material,
the king formula, precomputed table dimensions and hot static-score loop remain
unchanged, and all normal/invariant tests pass. Cheap inference alone does not
make these values generalize.

| Actual-choice suite | v0.22 mean loss | Candidate | Result |
| --- | --- | --- | --- |
| Public 100ms | 45.896 cp | 43.119 cp | Pass |
| Public 1000ms | 54.851 cp | 51.478 cp | Pass |
| Newer 100ms | 31.511 cp | 33.420 cp | Pass |
| Newer 1000ms | 24.242 cp | 25.947 cp | **Fail Black bucket** |
| New independent 100ms | 37.301 cp | 40.979 cp | **Fail stricter independent objective** |

The newer Black 1000ms bucket is 27.092 versus 21.646 baseline, above its
26.646 limit. Its false field and original driver exit remain. A prospective
[conditional fresh review](CONDITIONAL_FRESH_REVIEW.md) allows one new fixed
151-root, 16-family control after this narrow older miss, requiring at least 5%
mean improvement at each budget and no additional severe cases. The map remains
unchanged. At 100ms it instead worsens about 9.9%, although the original generic
tolerance still passes and severe counts stay six. The stronger condition fails,
the driver exits one and neither fresh-1000 nor games run. No alternate seed,
coefficient adjustment or hidden favorable retry follows.

All old and fresh safety reports retain zero new losing-mate flags and timing
overruns. Hand-checked endpoints, phase blending, symmetry, material units,
precomputed-table exactness, no allocation and root immutability pass. The
largest older Black deterioration, Bxe5 in `proof-fresh-10-ply37`, stays in the
raw data. Source, binary and every native/reference observation are preserved.

The result rejects this particular original table map, not every inexpensive
evaluation representation. Coordinate/phase preferences still cannot identify
actual blockers, shelter, support or safe activity. Together with the rejected
king representation and checking screen, it argues against another immediate
coefficient guess. Restore released v0.22 and reassess selective search with
explicit tactical/draw safeguards and a measured work opportunity.
