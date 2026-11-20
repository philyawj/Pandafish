# Verification of the combined review

Only research data, review documents, Go provenance and the live archive index
changed. Engine, search backend, timing policy, private configuration and live
recordings were not edited.

| Check | Result |
| --- | --- |
| Existing archive file fingerprints | 525 files hashed; unchanged on repeat audit |
| Raw inventory checks | 341 entries match SHA-256 and length |
| Original recording comparisons | Same 341 files match archived bytes |
| Cross-session game IDs | 157 occurrences / 157 unique IDs; no duplicates or conflicts |
| Complete PGN/analysis agreement | 15,012 plies; all reconstructed moves, FENs and bot-side identities match |
| Referee classifications | Recomputed from scores; source summary totals agree |
| Captured legal completed searches | 7,507 / 7,507 legal; all have metrics; zero saved fallbacks |
| Referee PV legality | Complete survey and selected confirmation PVs replay legally |
| Larger-budget completion | 32 / 32 selected cases saved with raw UCI output |
| Positive diagnostics | 24 / 24 selected CP errors persist at least 150 cp |
| Negative diagnostics | Both agreement controls remain equal; disagreement cases retained without forcing reference agreement |
| Mate interpretation controls | All four selected survey missed-mate cases still win by mate after played moves |
| Actual overlap / suspension | 20 cross-process interval pairs in two games; one wall/monotonic gap above one second |
| Recording/processing package tests | `go test ./tools/lichess` passes |
| Go provenance static checks | `gofmt` and `go vet` pass; preserved Stockfish adapter matches production bytes |
| Authored tracked whitespace | `git diff --check` passes |

The Go provenance program was formatted and compiled by `go run`; source data
was not reformatted. Its empty supplementary selection was not sent to Stockfish.
The existing three 100,000-node analyses were reused, not overwritten or rerun.
The completed confirmation was resumed for validation without issuing more
position queries; its saved input prefix, scores and frozen method match.
The [derived inventory](derived-inventory.json) records the resulting review
files independently of the original-session fingerprints.

These checks validate review correctness and recorded integration. There is no
new equal-time engine comparison, strength screen or live Elo estimate: no
engine implementation changed. Classification causes, final holdout behavior,
exact dirty client state and OS sleep telemetry remain unavailable.
