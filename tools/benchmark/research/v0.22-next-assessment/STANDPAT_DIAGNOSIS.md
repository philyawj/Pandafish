# Diagnose moving an existing stand-pat cutoff before child creation

At an unchecked quiescence node, consider an ordinary non-checking capture or EP
with no promotion. Its child is not checked. If the exact child static score from
the parent's perspective is at most alpha, then the child's stand-pat score is
at least its beta (-alpha): its existing quiescence would immediately cut off.
Stalemate/repetition can override that score with zero, so require alpha>=0 too.
Both outcomes then prove this move cannot improve alpha. No guessed victim-gain
margin or changed evaluation is needed.

Diagnose before implementation using the original immutable child and exact
staticScore, leaving the engine's actual SEE/child search unchanged. Count
potential pre-SEE candidates, SEE calls they could avoid, actual child visits
that would be avoided and any contradiction of the implied bound. Reuse the
frozen 20-game depth-five iterative trajectory, with no timing claim from the
extra diagnostic board copies/evaluations.

Require at least 20% potentially avoided actual node visits and 30% of actual SEE
calls, with zero observed bound contradictions, before implementing a cheap exact
piece-square delta and conservative bound shortcut. If a shortcut skips work
before SEE, retain a valid fail-low upper bound by taking max(best, predicted
parent-perspective child static score, 0); it stays <=alpha. A full-window/in-
window result must remain exact. Checks, evasions, promotions and negative-alpha
cases remain ordinary search. Do not claim identical fail-soft numbers outside
the window if the valid bound becomes looser; require exact completed root scores
and independent best-set membership. The new static delta must reproduce all
original phase/rounding/king effects, including promoted phase and special moves.

This is an exact relocation of the existing stand-pat decision, not heuristic
delta/futility pruning or a new material margin. Its actual cost and game impact
remain unestablished until implementation, controls and timed/native screens.
