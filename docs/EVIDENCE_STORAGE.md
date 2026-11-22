# Evidence storage and leaner research records

## Current decision

Keep existing evidence in its current repository paths. No cold-data compression,
external storage migration, deletion or history rewrite is scheduled. The
[post-v0.23 audit](REPOSITORY_AUDIT_V023.md) measured 1.375 GB of tracked checkout
content but only about 187 MiB of local Git objects. Checkout convenience and
compressed repository storage are different concerns; neither measurement alone
predicts remote clone-transfer size.

The archive pilot and under-200-MB checkout target in that audit are deferred
options, superseded by this decision. Existing raw data, failed attempts, hashes,
inventories and source identities remain intact. A future migration requires
its own review and a verified recovery path before any removal.

## Effective practice for new work

- Use existing Go plans, the workflow runner, comparison tools and inventories.
  Add a one-off diagnostic only when an existing command cannot answer its
  concrete question; preserve its Go source and frozen method.
- Save one canonical raw observation stream for a new diagnostic, with stable
  case IDs, source/input hashes and declared settings. Generate compact aggregate
  findings and retain exact adverse witnesses. Do not save another full copy of
  the same observations merely for a different report layout.
- Use immutable input/reference identities and paths instead of copying consumed
  corpora, complete cached queries or source trees into every new directory.
  References must remain resolvable; record whether external bytes were actually
  verified. Source snapshots still use the supported replayable Go format.
- Let Go derive counts, ratios, gates and summaries from the raw records.
  Prefer compact summary JSON and Markdown to per-position expanded JSON when
  no consumer requires the expanded representation. Diagnostic instrumentation
  must not produce millions of routine trace rows when aggregate counters plus
  deterministic inputs and sparse contradiction records answer the hypothesis.
- Preserve failures and interrupted runs, including partial observations and
  numerical/tooling errors. Fewer redundant representations must not mean fewer
  unfavorable results or less provenance. Never rerun engines to regenerate a
  missing presentation file when saved raw observations suffice.
- Freeze the protocol, source and gates, then seal completed evidence once.
  Later document corrections use explicit mappings or amendment records;
  inventories must not be rebuilt to conceal changed observations.

## Existing readers and a future format change

Current `choices` writes both `native.jsonl` and sorted `native.json`, and current
quality/depth-review readers depend on those files. Those established outputs
remain supported and retained. This policy does not authorize deleting either
one or silently introducing compressed files that readers cannot load.

When that producer/reader contract next needs work, prefer a narrow compatible
change: read a canonical journal, reconstruct the sorted view deterministically,
and write only compact final summaries by default for the new schema. Validate
case/trial equality, all quality decisions, interruption recovery and historical
schema compatibility without rerunning native searches. Keep old evidence bytes
unchanged. Apply a similar rule to duplicated final quality/checkpoint bodies
and embedded reference-cache originals only after their consumers are checked.

Lossless in-repository compression remains an optional later step for cold data.
Retain original decompressed hashes, avoid oversized combined archives, and
provide explicit extraction/read support before adopting it. Compressing the
latest checkout does not erase historical blobs and may not reduce Git's existing
compressed storage proportionally.

## Monitoring and reconsideration

After a substantial research batch, record `git count-objects -vH`, tracked
checkout size, the largest added files and the new evidence volume in its compact
summary. Record loose and packed object sizes separately; local unreachable
objects and packing can affect comparisons. Reuse the audit's Go sizing method
only when a detailed breakdown is needed, rather than repeatedly recompressing
the whole collection.

Reconsider storage policy when checkout bulk becomes inconvenient, a batch adds
unexpectedly large raw evidence, or compressed repository storage approaches
1 GB. That is a planning trigger, not a hard GitHub limit or a forecast. Review
files approaching 50 MiB individually. GitHub's current guidance is linked from
the audit; check it again before designing a migration around a service limit.

No Git-size emergency or engine-strength benefit is claimed for this policy.
The immediate aim is slower evidence growth with the same scientific record.
