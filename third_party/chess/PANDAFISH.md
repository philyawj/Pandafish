# Local search API extension

This is the runtime source of `github.com/corentings/chess/v2` v2.6.1,
copied unchanged, plus `search_moves.go`. The upstream MIT license is retained
in LICENSE. Source hashes are recorded in
`tools/benchmark/research/v0.17-search-moves/upstream-source.json`.

The additional API provides caller-owned legal move storage, direct tactical
move generation and allocation-free legal-move existence checks. Existing
Position/Game/Board methods and state/hash behavior remain unchanged. The
original library tests and independent move/position oracles validate this
extension separately from engine search performance.
