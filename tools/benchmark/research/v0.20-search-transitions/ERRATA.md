# Validation environment findings

The first upstream replay ran the entire original module, including rendering
packages outside the vendored runtime. Its root runtime race suite passed, but
the unvendored image package could not acquire the uncached SVG dependency under
the sandbox's module-cache restrictions. The complete failed log is preserved.

The scoped root-runtime replay matches the actual vendored boundary and runs
without that unused rendering dependency. No dependency version, original
runtime source or engine behavior was changed to work around the environment.
