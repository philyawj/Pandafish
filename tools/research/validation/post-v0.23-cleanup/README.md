# Post-v0.23 repository cleanup

This is a tooling/documentation increment without a release or engine-policy
change. Existing evidence remains in place under the
[storage policy](../../../../docs/EVIDENCE_STORAGE.md). New diagnostics retain
canonical observations, compact summaries and sparse adverse witnesses; old
producer/reader schemas are unchanged.

Implemented cleanup:

- Publish the repository audit, ranked assessment and their unchanged diagnostic
  evidence. Their original observations concern tagged v0.23.
- Replace the stale current checkpoint, advance active baseline examples to
  v0.23, add a current backend preface and distinguish historical guidance.
- Add root/nested chess normal/race/vet to the research validation builder before
  builds. Benchmark validation now checks the nested module too, respecting its
  race flag. Failure-injection tests cover nested test/vet stopping before build.
- Repair the old quiet-check screen link through a small navigation file, leaving
  its frozen candidate README, failed direct gate and inventory unchanged.
- Document deferral of compression, external archives, format migrations, a
  cohort ledger and shared UCI source. No data or source tree is reorganized.

[The validation run](validation-run/state.json) passes all eight steps: root
normal/race/vet, nested normal/race/vet, warm-test build and engine build. Root
UCI tests exercise subprocess protocol/search behavior. This increment changes
validation scheduling rather than the engine or match protocol, so no strength
screen, live client or additional game smoke is run. Original v0.23 matches and
failed release smokes retain their identities.

[Document validation](document-validation.json) checks 192 local inline links,
13 tables and all eight engine candidates. New policy/checkpoint wording and
Go formatting/whitespace checks pass. Original first-person historical reports
remain historical; current additions use factual engineering prose.

[The document-map transition](document-map-transition.json) updates four mapped
guides while preserving original hashes, retired-source metadata and raw data.
[The report amendment](report-amendment.json) verifies the exact saved pre-edit
audit/assessment against their original document-validation hashes. Only local
handoff linkage and subsequent-decision context change; the original validation
record stays intact. The old checkpoint is retained as a frozen copy here.
Copied historical documents retain their original relative-link context.

The local `handoff.md`, private configuration, executables and caches are excluded
from the staged publication. The next legality/attack-context research starts on
a separate branch after this PR is merged. No playing candidate is committed or
published by this cleanup, and no tag/release is created.
