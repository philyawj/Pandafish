# Move representation opportunity: proceed to a bounded exact kernel

The native general Move is128bytes; origin/destination, promotion and tags fit
in6bytes. Board and Position sizes are192/80bytes. This is a concrete search
representation mismatch, not evidence that bitboards or scratch state are absent.

One unchanged warm allocation sample preserves all1,145 roots /43,118,858 nodes.
Normal512-KiB sampling and one end-of-test heap profile replace dense per-root
memory-stat reads. `alloc_space` attributes76.76MiB flat to library search
generation and63.02MiB to castleMoves, versus58.37MiB to tables. These sampled
owners are separate flat allocations; cumulative rows overlap. The sample also
includes PGN loading and report I/O. Allocation volume is not elapsed headroom.

The castle helper returns a two-full-Move backing array even when its occupied-
path conditions reject every castle. The additive search API can prefilter
impossible rights/occupancy before entering that exact existing rule routine.
Candidate move tagging can use only the six needed bytes, without copying the
full PGN/tree value on every pseudo-legal candidate and existence probe.

Proceed with one **compact candidate-generation kernel**, not a full public
Move/API or board rewrite. Full emitted/cached Move values retain metadata;
engine search policy stays unchanged. The prospective8% exact warm target and
all later correctness/decision gates remain necessary. This diagnosis passes
the representation/ownership opportunity condition, not a cost or strength gate.

The raw allocation profile, compact size result, Go source, profile log and
shared exactness report remain here. Full warm rows are temporary and hash-
identified by exactness.json; original exact rows remain in the v0.23 archive.
Reproduce with the existing warm harness's `-test.memprofile` and default
`-test.memprofilerate=524288`, then the standalone installed pprof analyzer's
`-alloc_space`/`-alloc_objects`. No instrumented elapsed ratio is used as gain.
Tagged baseline test hash remains the post-v0.23 audit's recorded identity.
