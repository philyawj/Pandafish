# Initial control fixture correction

The first control run fails only the far-past wall-date case. Its fixture uses
0001-01-01 00:00:00 UTC, which is Go's zero Time and deliberately means an
unlimited search in both baseline and candidate. The expected stopped result was
incorrect. Preserve `deadline_test.initial.go.txt` and the complete failed
`clock-controls-initial.txt`; replace that fixture with year 2, which is still
beyond Duration's negative range but is not the unlimited sentinel. Runtime code
is unchanged by the correction. Ordinary, injected boundary, cancellation and
recursive deadline controls already pass in the initial run.
