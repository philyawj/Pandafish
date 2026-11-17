# Preserved development results

The first validation stopped at formatting: two copied upstream files,
`board.go` and `pgn.go`, are not canonical gofmt output. They remain byte-exact
copies of the pinned dependency. I narrowed authored-source formatting checks
to the root module and the new search API file, matching the benchmark runner's
existing module-based scope. The original output/exit records are retained in
`format-initial.txt`, `checks-initial.json` and `validation-run-initial.txt`.
No playing code or measured result was changed to address this check.

The release integration run returns exit 1 because its two v0.16 games both
draw (50%, below the strict >50% gate). All six integration games complete
legally with zero forfeits/unfinished games. I preserved the failed tiny gate
and every result; it does not replace the passed 50-game candidate screen.

The first staged whitespace check flags the original upstream LICENSE's extra
blank EOF line. I retained the verbatim license and its upstream hash, and
added a whitespace-diagnostic exemption only for that license file. The raw
failed check and exit are preserved as `precommit-whitespace-initial.*`.
