# Decision after the coordinated search-core screen

I subsequently approved promoting the screened code as v0.16, authorized its
PR/merge/release, and deferred the proposed larger confirmation. The original
recommendation and proposed protocol below remain historical evidence. I also
authorized reassessing and starting one next candidate after release. That
new experiment will have its own source, gates and review before another PR.

**I recommend advancing the frozen candidate to confirmation and holding release
promotion until that review.** I do not recommend rejecting the mechanism or
changing its weights/search settings in response to this screen.

The functional result is clear within the tested scope: all 192 paired fixed-
depth searches retain scores and best-set membership; real draw-history,
mate-distance and cancellation controls pass; native equal-time depth improves
under serial and three-search loads. The game result is a promising direct
65% screen, encouraging Sacrifice evidence and no clear Stockfish improvement.
That is enough to justify confirmation, not enough to declare a general leap.

The direct score is 50% as White and 80% as Black over 25 games per color.
Across the ten repeated opening positions, seven score above 50%, two at 50%
and one below. The lead is therefore not confined to one opening, but its
color concentration strengthens the case for a broader, separately reviewed
confirmation. Repeated openings and finite timing variance limit the nominal
paired interval. The earlier 27% Sacrifice control is historical context, not
an external baseline run concurrently with this candidate.

## Proposed confirmation — not executed

- Freeze the exact candidate source/binary hashes in `provenance.json` and
  `source-hashes.json`. Do not tune between opponents or after intermediate
  results. The already-tested source overlay is in `candidate-source.json`.
- Play 200 games per opponent (`--pairs 100`) over the existing public
  `openings-confirmation-50-v1.txt`: 50 positions, both colors, two cycles.
  This broadens the opening set as well as increasing the count. It remains
  public and partially overlapping with diagnostics, not a blind holdout.
- Retain tagged v0.15.3, Stockfish Skill 2, the same pinned Sacrifice/Tal asset,
  100 ms requested time, 50 ms protocol allowance and three concurrent matchups.
  Capture source/binary/asset hashes and machine/power notes again. Run no
  other engine probes or live client concurrently.
- Preserve every game, failure, timeout, gate decision and color/opening
  breakdown. Require a positive direct result with zero forfeits and inspect
  uncertainty and variation rather than accepting any barely-positive score
  as a broad gain. Review each external opponent separately; do not add their
  scores or compare historical samples as if they were paired controls.
- If reliability/correctness fails, investigate the exact position/request.
  If the direct advantage disappears or the broader set exposes a material
  regression, hold promotion and diagnose it; do not rerun until favorable.
  Even 200 games per opponent may leave a modest effect unresolved.

This is a proposed protocol awaiting a separate decision. No confirmation,
commit, push, PR, release or live game follows automatically.

## Reassessed implementation queue

I keep specialized tactical generation/reusable storage as the next likely
engineering experiment if confirmation supports this candidate. The combined
result supports the search-efficiency direction, so I do not need another
broad investigation before deciding that project. The larger profile target
remains: this candidate still generates the full legal list before filtering.

I retain coherent evaluation/threat modeling as the next major strength
workstream. The evaluator and diagnosed quiet-threat/endgame weaknesses are
unchanged. It should move ahead of the backend project if confirmation or
focused loss diagnostics show that the extra useful depth is not addressing
the actual failures. I have not demonstrated that causal claim here.

I will choose one next experiment after reviewing the confirmation decision,
not implement all three as a checklist. No bishop-pair or rejected piece-safety
term is reintroduced, and no old PVS/LMR bundle is silently restored.
