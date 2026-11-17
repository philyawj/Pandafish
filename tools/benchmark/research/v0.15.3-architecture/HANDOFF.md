# Start here after v0.15.3

v0.15.3 is a **research patch with unchanged v0.15.2/v0.15 playing behavior**.
I authorized this PR/merge/release and explicitly chose to archive
the bishop-pair candidate rather than activate it. No cache, PVS, move-generator
or bishop-pair prototype is compiled into normal builds. No live bot or
200-game confirmation was started. Private configuration/active root executable
were not changed.

Read [README.md](README.md), [NEXT.md](NEXT.md), [CACHE.md](CACHE.md),
[METHOD.md](METHOD.md) and [ERRATA.md](ERRATA.md). The full report distinguishes
measured work, decision diagnostics, real/intentional test failures, old-version
control games and release smoke. Original labels/hashes stay unchanged.

## What changed our priorities

- The v0.15 repetition fingerprint causes zero score-cache hits in these cold
  shallow probes, but removing history causes an actual -902 cp/draw error.
  Safe board hints plus reversible-context identity save 4.5% of depth-four
  nodes on 36 inputs and 16.1% of depth-six nodes on six inputs. This is not
  a reason to revert draw correctness or buy a larger table.
- The clean CPU profile puts 50.9% under move generation/ordering, with about
  72.8 MB / 120k allocations per timed search. Quiescence is 91.3% of cold
  depth-four nodes and filters most of the fully generated quiet moves.
- A coordinated exploratory prototype saves about 37% of deep-case nodes,
  matching all tested scores and selecting members of baseline best sets.
  It intentionally changes the all-ties root contract and has no strength
  screen. Full contract review and broader validation are still required.
- Static/leaf blindness is real too: Qd6's leaf is +161 for Black versus -662
  reference cp; the Kxg5 leaf misses a quiet mate in one. Nd4's own PV leaf is
  reasonably scored, and Rd5's stronger opponent response is the bad pawn
  ending. Main-PV-only diagnosis is insufficient. PNfQoYP0's reference is Qf3,
  not the old prose's Bf3.

The retrospective released-policy screen completes all three 50-game matchups
concurrently: 56% against v0.14.4, 44% Skill 2, 27% Sacrifice. The old-version
comparison also changes mate/repetition rules, so it cannot assign strength
cost to the cache. A wrong descriptive previous-revision label is preserved
and corrected by the companion provenance; the executable itself is tag-built.

The original bishop screen is separately archived: 57% versus v0.15.2, 44%
Skill 2, 25% Sacrifice. It does not establish a broad leap or a definitive
external regression. Its optional confirmation was deferred by the broader
research request, not run in the background.

## Best next implementation

Start from **tagged v0.15.3**, building its previous executable from the tag.
Harden and test the coordinated search-core prototype first, then address
specialized legal tactical generation/buffer reuse as the larger throughput
project. Preserve history-dependent draws, mate distance, legal special moves,
termination and cancellation. Make the one-best root contract explicit; do not
silently weaken all-ties tests. Prototype overlays apply to v0.15.2 and are
replayable; the new released playing code differs only in identity.

The next major strength workstream is a coherent data-calibrated evaluation/
threat model on independent response-branch data, with classical versus compact
trained/NNUE routes evaluated for cost. Do not revive the rejected safety
formula, guess another 25 cp bonus, blindly retry old LMR parameters, or assume
a faster search fixes every wrong leaf.

Use fixed-depth parity/move-membership, broader tactical/rule controls and
native equal-time checks before the standard 50 games each against released
baseline, Skill 2 and pinned Sacrifice, three concurrent matchups. All results
stay. A larger confirmation is separately planned, never automatic. Reassess
after each result. No new playing-feature implementation follows this research
release automatically; check in with me on the next candidate.
