# Move-count draw claims: one isolated search contract

Begin from tagged v0.23, with all previous failures archived and no material or
contextual evaluator reactivated. The diagnosed gap remains: game adjudication
recognizes move-count draws while search ignores them. Existing score entries
already compare halfmove clocks; preserve and test that protection.

Model a valid 50-move claim as an available zero-valued action, rather than
forcing every claimable winning position to zero. A current halfmove clock of
at least100 permits the claim; at99 a legal non-pawn noncapture intended move
permits it. At150 termination is automatic, after checkmate precedence. Captures
and pawn moves reset the clock. Existing threefold policy stays unchanged.
Root search must associate an intended-move claim with a qualifying allowed
legal move, and retain a winning mate over an optional claim. UCI has no claim
verb: the score models claim availability; the game client handles the claim.
The benchmark already claims automatically when eligible. No live-client
configuration or live games change.

[FIDE articles9.3 and9.6.2](https://handbook.fide.com/chapter/e012023) define the
claim/automatic distinction and mate precedence. Search/quiescence keep their
existing finite horizons; positive heuristic scores remain estimates, not proof
that a claimable game is won. This does not implement every optional repetition
choice or solve the rejected material-rule candidate.

Before timing, require boundary truth tables (98/99/100/149/150), qualifying
intended moves versus pawn/capture-only evasions, current claims while checked,
mate precedence, root searchmoves restrictions, irreversible resets, both-turn
symmetry, cancellation/history restoration and warm score-cache clock contexts.
Use an exhaustive small-tree reference with explicit claim actions for focused
main-search value controls. Normal/race/vet in both modules, builds and596
independent move/state vectors must pass.

Freeze a three-pair warm cost ceiling of5% (median ratio<=1.05), identifying every
changed value/move and retaining all adverse rows. This is a correctness change,
not an exact speed rewrite. Require unchanged existing/fresh actual-budget gates
before one standard concurrent50-game-per-opponent screen. No unchanged retry
or threshold changes. No research commit/publication, release or live play.
