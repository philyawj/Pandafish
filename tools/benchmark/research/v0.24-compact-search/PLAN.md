# Coordinated compact search-move representation

Baseline playing policy is tagged v0.23. The temporary compact-generation
foundation saves2.166% but fails its8% objective; all original source, controls,
cost rows and gate flags stay preserved. This project extends that foundation
through the search consumer rather than promoting it under a smaller target.

Introduce a6-byte search-only move value and one shared generator with full-
metadata or compact emission. Keep existing public Move/Position APIs and
cached metadata. Worker-owned compact per-ply buffers, stable ordering, staged
picker and internal descent use the small value. Root/frontend and diagnostic
compatibility interfaces can retain full Move values. Compact child update
converts at the existing state/hash boundary, avoiding a new board rules system.

Original upstream runtime files remain byte-identical. All generated signatures,
tags, order, state/hash, terminal status, search/evaluation/SEE policy, TT/history,
mate/repetition and clocks stay v0.23. No target-attack or full pin context is
included. Small castling prechecks remain the exact foundation; the original
castle routine still decides legality/tags.

Require full/compact differential lists, cached metadata and round-trip controls,
special/malformed moves, stable order, active-ply ownership, preferred-only zero
allocation, both-module normal/race/vet, builds and596 independent vectors.
Three alternating warm pairs vs the tagged baseline must retain every1,145
move/score and43,118,858 nodes, with median ratio<=0.92 (>=8% saving). No adjusted
threshold or unchanged favorable retry. Instrumented/alloc samples are separate.

Only a passed exact cost gate permits four existing actual-budget suites,
followed by prospectively frozen untouched fresh both-budget/both-turn controls
and the standard concurrent50 games each vs v0.23, Skill3 and pinned Sacrifice.
Stop at a recommended candidate before any research commit/publication. No live
play or new release is authorized. Keep deterministic execution in shared Go
tools and preserve every failure with compact evidence where readers permit it.
