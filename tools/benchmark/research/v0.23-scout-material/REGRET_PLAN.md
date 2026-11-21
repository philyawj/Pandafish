# Independent selected-move full-depth review

Before a strength screen, replay every depth-five candidate choice from the
frozen 1,145-root warm workload using tagged, unreduced v0.22 in a separate
test executable. Reconstruct the same game history through the shared harness.
Use a fresh table and full root window for that one selected move; only the
diagnostic reference is restricted. Production root searches remain unrestricted.

Compare its value with the corresponding completed unreduced baseline optimum.
Any co-best move is accepted; report finite positive regret and excess beyond
25cp separately. Freeze review gate before queries: mean finite positive regret
<=25cp, no selected new losing-mate result when the baseline optimum avoids it,
no missed baseline winning mate, no unsupported selective mate claim and no
selected reference value above the known optimum. Preserve every row and adverse
case. This is a fixed-depth diagnostic, not a measured Elo or equal-time result.
The actual-budget six suites remain separate mandatory checks.
