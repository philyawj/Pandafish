# Preserved validation findings

The upstream replay helper was launched with the original whole-module command
before its scope correction was applied. As in the transition-only attempt, the
vendored root runtime tests passed while the unvendored rendering package could
not acquire its uncached SVG dependency under sandbox cache restrictions. Both
the failed whole-module log and the explicit root-only result are retained.
No engine code, original runtime or dependency version was changed for this.

The Go reference-query driver now journals every query before validating its
score. An added Go subprocess control verifies that a scoreless reference error
retains its complete output. All earlier observed native/reference queries were
successful; their records and accepted metrics remain unchanged. This repair
strengthens failure preservation and does not affect the playing engine.
