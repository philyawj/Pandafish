# Revision: occupied king rays

The [full-ray attempt](../v0.22-geometric-tags/README.md) fails timing before
native or strength screens. Preserve its corrected source, caught initializer
defect, complete controls and all slower/faster timing observations.

Refine the same conservative mechanism using current occupancy. Only the first
occupied square on a king ray can expose that king when moved; a piece behind
another blocker cannot expose a slider by leaving. EP/castling remain special
simulation cases. Cache the visible orthogonal/diagonal rays once per king and
intersect origin-risk masks with the moving side's pieces. Use those visible
rays for sliding-piece checking destinations too: if an earlier blocker remains,
a destination behind it cannot check. Moving the blocker itself is already
covered by its origin-risk mask. Pawn/knight/king direct-check masks and promotion
type handling remain as before. Already-attacked opponent kings and checked
parents retain conservative tests; all simulations use the unchanged primitives.

Add shielded-origin controls to the prior complete tag vectors. Keep the exact
15% elapsed gate and all subsequent controls unchanged. No storage reuse or
other failed code is bundled. If it fails, archive and reassess another mechanism.
