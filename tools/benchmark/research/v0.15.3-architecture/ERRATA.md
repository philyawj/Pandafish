# Corrections and interpretation limits found during this investigation

- The old live-review prose calls PNfQoYP0's `e3f3` reference **Bf3**. The saved
  FEN has a white queen on e3: the move is **Qf3**. Raw FEN/history/UCI data are
  correct and unchanged. This is a queen-versus-rooks/minor-piece position,
  not direct evidence for a pure pawn-ending evaluator.
- An evaluation on the engine's chosen PV leaf need not expose the cause of a
  bad decision. Nd4's reconstructed leaf is close to the referee's score.
  Rd5's reconstructed rook leaf is also close, but its stronger opponent
  exchange leads to a badly assessed pawn ending. Probe opponent deviations
  as well as the engine's own preferred line.
- The v0.15 history fingerprint cannot explain the earlier pre-v0.15 PVS/LMR
  screen. The cache introduction, exact-depth semantics and repetition-context
  change happened at different times.
- The first architecture normal-test pass accidentally ran diagnostic tests
  without output paths. Logs are preserved; corrected drivers skip them unless
  explicitly requested. This was a harness error, not a released engine fault.
- The first joint CPU/allocation profile included substantial allocation-stack
  profiling overhead. A clean CPU follow-up disables allocation sampling and
  is used for CPU claims. Cumulative rows overlap.
- Twelve warm UCI requests are terminal mate boards after the two recorded
  continuation moves. Their null move is expected; raw `legal:false` is not
  interpreted as an illegal engine move.
- The retrospective match command supplied a wrong descriptive previous-
  revision abbreviation. The actual executable was built from tag v0.14.4 at
  `3f3ca331cbe8aac9c7c2194b9314b33b3da34a81`; its SHA and the original mistaken
  label are retained in [control-provenance.json](control-provenance.json).
  No playing result is changed or rerun to repair this metadata mistake.

- The final staged whitespace check caught a trailing-space Python blank line
  and the intentional prefix/tab layout of a saved Git patch. The script is
  cleaned; the patch is renamed `.patch` without changing bytes and narrowly
  exempted from whitespace checks with the raw check output. No playing or
  diagnostic algorithm changes.
