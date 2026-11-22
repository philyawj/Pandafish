# Draw semantics: bounded policy diagnosis

Tagged v0.23 playing source is restored after the rejected compact-search
screen. No speed or playing-strength result from that candidate is transferred.

Hypothesis: search assigns nonzero values to a conservative subset of positions
that the existing game interface automatically ends for dead material. First
record an explicit truth table with both root turns: bare kings, one minor,
same-colour bishops, and live-material negative controls. Compare game outcome,
position terminal status, static value, main search and quiescence. Also record
99/100/149/150 halfmove boundaries and checkmate precedence without changing
move-count policy. No new teacher queries, datasets or games in this diagnosis.

Continue with one isolated dead-material correction only if the mismatch is
reproduced and the accepted subset has a direct no-possible-mate argument.
Do not combine claim policy or a general dead-position solver. A fifty-move
claim is optional under ordinary over-the-board rules; the benchmark claims
eligible draws automatically. UCI has no explicit claim action. Keep that
separate policy question open rather than silently forcing all100-halfmove
positions to zero.

If implemented, require normal/race/vet in both modules, sound material negative
controls, special-move transitions, search/q/root consistency, cache/history and
cancel checks. Freeze a three-pair warm overhead ceiling of5% (median<=1.05),
not an efficiency claim. Every changed fixed-depth row must be identified;
non-draw controls must remain exact. Require all ordinary native/fresh gates
before one unchanged standard screen. Functional rule correction alone is not
evidence of increased playing strength. Preserve failures, restore baseline
between mechanisms and stop before publishing research source.

Rule expectations use [FIDE Laws5.2.2,9.3 and9.6](https://handbook.fide.com/chapter/e012023):
dead positions end immediately;50 moves permit a claim;75 moves end the game
automatically, with checkmate precedence. The conservative material subset
does not attempt every blocked or otherwise dead position.
