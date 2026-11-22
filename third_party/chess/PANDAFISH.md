# Local search API extension

This local module originated from the runtime source of
`github.com/corentings/chess/v2` v2.6.1. The upstream MIT license is retained in
LICENSE. Original source hashes are recorded in
`tools/benchmark/research/v0.17-search-moves/upstream-source.json` at the repository
root. Subsequent search-specific extensions are described in the
[backend note](../../docs/CHESS_SEARCH_BACKEND.md).

The additive APIs provide caller-owned move storage, staged tactical/quiet
generation, evasion filtering, worker-owned child state, exact exchange helpers
and bounded protected queen-contact queries. General Position/Game/Board APIs
remain available. Root tests exclude this nested module; run its tests, race
checks and vet explicitly. Independent move/state vectors check the rules
boundary separately from engine search performance. v1.0 pauses further engine
development; historical research plans are not an active queue.
