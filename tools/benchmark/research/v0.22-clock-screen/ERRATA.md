# Reference-cache argument failure after the first native block

The initial Go native driver supplies the v0.21 `reference-cache.json` metadata
manifest as `--cache`. The choices tool accepts a complete prior run directory
or a flat legacy query cache with an engine hash; the metadata manifest contains
a journal hash/source instead. All 84 public cases / 504 actor queries complete
first, then reference loading exits with `reference cache engine hash differs`.
The pinned Stockfish binary still matches its original SHA-256. This is a driver
argument failure, not evidence of an engine/hash change or a failed quality gate.

Preserve the initial driver, native log/check record and every native observation.
A Go reference-only continuation calls the existing private reference function
on those saved observations with the correct prior run directory. It checks
teacher/budget/case identities and proves native.json, native.jsonl and method.json
are byte-identical before/after. It starts no Pandafish process and reruns no
engine choice. The reference output and repair status remain separately saved.

Only if that original block's actual quality report passes, continue the three
previously unexecuted suites with corrected cache-directory arguments. Keep the
initial failed check file unchanged; continuation checks have a separate file.
No favorable replacement trial or quality retry is authorized.
