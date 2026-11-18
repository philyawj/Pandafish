# Exact exchange and integration controls

The unchanged original material/legal exchange fixtures cover free/defended
captures, pins, legal/illegal king recaptures, opened rook/diagonal rays, EP,
quiet/capture underpromotion and recapturing-pawn promotion. Existing independent
full legal-recapture oracle/filter controls and deadline/cancellation tests pass.

New seeded and focused controls compare **46,794 full exchange scores** against
the frozen v0.21 mailbox implementation, including ordinary quiet moves beyond
the runtime filter's capture-only use. **33,152 attacker sets**, attacked flags,
least legal attacker choices and complete parent restorations agree. Constructor
mailbox/king pointers agree on every checked board. Pins, king safety, EP,
promotion, supplied missing/nonunique kings and special castling rights are
included. Equality is checked before score thresholds or root choices.

All earlier generation/ordered-Move controls stay active. Root race/vet/build,
pinned upstream root-package race and 596 independent frozen rule vectors pass.
The oracle invocation uses the actual frozen `oracle-inputs.json`; the prior
evasion phase's wrong-path failure stays in that archive and is not silently
rewritten. Source is frozen before the fixed-work drivers run.

Cold comparisons use the original three alternating public depth-four pairs,
public depth three, twelve depth-six roots and newer depth-four roots: 325
distinct position/depth combinations. The 15% total median and deeper/newer
requirements remain. Three warm depth-five pairs compare with the **first
evasion revision**, requiring 10% extra median elapsed savings and identical
scores/moves/nodes over the frozen 20-game / 1,145-root persistent-table sequence.
A separate new warm pair compares total elapsed time directly with tagged v0.21.
Saved drivers refuse existing result paths. Per-root parsing/MemStats/output
overhead is separate from search-time totals; profiles do not set efficiency.

Native budgets and games run separately from timings and heavy validation.
All actor/reference errors, failed gates, identities and PGNs remain. The four
native reports must be reviewed explicitly, not only by command exit status.
Teacher and close-alternative conventions remain unchanged. No screen, favorable
unchanged rerun, larger confirmation, live game or publication is automatic.
