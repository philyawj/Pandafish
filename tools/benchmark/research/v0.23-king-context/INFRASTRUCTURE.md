# Preserved documentation-map preflight failure

The first fit workflow stops in its evidence audit with `evidence mismatch:
AGENTS.md`. No snapshot, training fit or new development metric runs. The user
authorized new research guidance; the legacy document migration map still
described the preceding AGENTS.md/tools guide bytes.

The Go `map-docs` helper records explicit before/after transitions for AGENTS.md,
tools/README.md and the authorized release-guidance addition. It updates only
current Markdown mappings, leaving original hashes, retired-source metadata
and raw observations intact. Focused race controls check those invariants and
reject raw-data/unregistered paths. The new transition record remains here.

`fit-plan-docmap.json` and `fit-run-docmap` are a new preflight run after this
bookkeeping repair; the original failed plan/run remain unchanged. This is not
a rerun of a fit or disappointing scientific outcome: neither had executed.
The repaired plan also freezes the migration map and affected guide identities.

That second preflight discovers one additional stale mapping for the live index,
whose content already matches HEAD from the earlier merged review. The bulk
audit now reports all original mismatches together rather than revealing them
one at a time. An explicit remaining transition records that HEAD match and a
release-guide paragraph-layout correction. Both failed preflights are retained;
the final ready plan starts before the first and only fit. No label or model
outcome is rerun during these documentation repairs.
