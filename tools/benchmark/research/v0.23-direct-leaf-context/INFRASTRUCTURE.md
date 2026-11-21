# Preserved pre-query input failure

The first label plan stops before creating a run directory or querying an
engine: the Homebrew Stockfish path is a symlink and the shared input hasher
initially accepts regular files only. Keep the original plan/error here.

Workflow inputs now hash both the resolved file contents and the declared
symlink/target identity. A retarget is detected even when the new target has
identical bytes. Source snapshots still reject symlinks. Focused controls cover
both properties. The protocol, selected roots, teacher budget and gates stay
fixed; no scientific result is retried by this pre-query infrastructure fix.

Label collection uses separate immutable diagnostic and remaining journal
segments. Completion reuses the diagnostic segment but does not modify it, so
workflow output hashes and any later resume remain valid. Failed queries are
preserved and cannot be implicitly re-queried.
