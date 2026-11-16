# Preserved early failures and corrections

- Initial Go build and trace test used the default Go cache, which the sandbox
  could not read/write. Repeated with `GOCACHE=/tmp/pandafish-v016-go-cache`.
  Go additionally warned about its module metadata stat-cache write; builds
  completed. The previous executable was subsequently rebuilt from the tag's
  scratch archive, not from candidate source.
- The first focused test compile called a nonexistent `Piece.FEN` method.
  Replaced it with an explicit FEN character mapping in the test helper.
- `controls-first.txt` preserves two bad expected results. The pinned attacker
  fixture also exposed the pinned knight to the white rook (+16); a black pawn
  now defends that knight. The purported opened-ray bishop was on the wrong
  diagonal. The corrected board uses f6-e5-d4, with the other targets defended.
  Production code was unchanged for these fixture corrections.
- `baseline-traces.json` used every historical position for its draw fingerprint.
  `baseline-traces-history-corrected.json` uses the production reversible-history
  helper. Both are retained; the corrected record is authoritative.
- The initial static microbenchmark overlapped the first baseline probe run.
  That first probe run is retained as a workload-contaminated diagnostic, not
  used as the final equal-time cost comparison. Clean paired probes follow.

- Release smoke completed successfully, but the first summary-reader looked in
  the parent results directory. The runner saves a nested matchup-label folder;
  the documented link and release-check record now point at that saved summary.
  No smoke game was rerun.
