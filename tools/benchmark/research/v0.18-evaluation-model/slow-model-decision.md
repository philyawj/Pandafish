# Slow-model decision

The predeclared 100 ms decision-quality gate fails: mean clamped reference
move loss 55.57 cp baseline versus 64.10 candidate; severe cases 1 versus 3.
At 1,000 ms the candidate improves 68.0 to 56.33 cp, with severe cases 3/3.
All raw data, 234 reference queries and the frozen source are retained.
No standard games were run for this version and no weights were changed.

I am testing one separately labeled inference-cost revision with the exact
same weights/values: precompute the original material/placement functions and
iterate the existing piece bitboards. Exact score parity is required. This
tests whether cheaper inference can preserve the model's longer-budget
benefit at the standard budget. It does not rerun the identical candidate or
select new weights against these outcomes. The same decision gates apply.
