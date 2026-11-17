# Preserved active-model test findings

The first normal test run passes the 13,495 independent scalar/perspective
fixtures but flags two old material-approximation assertions: Black's queen
capture scores308cp versus about400, and en passant scores208cp versus about100.
The learned positional residual changes those nonterminal values. I retain the
original output and test contracts as active-tests-initial.txt and
initial-test-contracts.json.

I compare the nonterminal tactical scores to independent unpruned/unordered
minimax instead of a material-only tolerance. Terminal mate/draw values remain
exact. I retain the unique safe queen capture, its literal400cp material balance,
and an explicit en-passant gain/material control. No learned weight or search
policy changes to satisfy those checks.
