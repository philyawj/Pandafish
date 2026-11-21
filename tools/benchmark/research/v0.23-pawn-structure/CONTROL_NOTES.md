# Preserved control failures and measurement limits

The first full test run fails two expectations. The manually calculated doubled
pawn value omitted that the rear pawn is immediately blocked by its own front
pawn. Correct the expected values from 5/23 to -1/10; the independently written
mailbox oracle already agrees with the bitboard implementation. The original
en-passant smoke case allows only 90cp departure from a material-only value;
the resulting d6 passer has contextual value. Require its exact legal EP result
under the scalar base-plus-pawn oracle instead of increasing the old tolerance.
All original failure output remains in `initial-tests.txt`. No model coefficient
or eligibility changed. The original base-table oracle remains intact; subtract
the separately checked scalar pawn correction to verify unchanged placement.

The first component measurement was mistakenly started before the race command
had returned completion. Keep it as `component-cost-overlapped.txt`, provisional
and unused for acceptance. The separate isolated component run and all warm
pairs start only after race validation has finished. No heavy agent workload
overlaps native choices or match play. External system activity is unmeasured.
