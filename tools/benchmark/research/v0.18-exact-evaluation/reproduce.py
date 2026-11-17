"""Restore a frozen exact-evaluation candidate over tagged v0.17 in a NEW directory.

This runs source tests and a build only; it starts no matches or live games.
"""
import argparse
import io
import json
import os
from pathlib import Path
import subprocess
import tarfile

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("output", type=Path)
args = parser.parse_args()
research = Path(__file__).resolve().parent
repo = research.parents[3]
snapshot = "candidate-source.json"
source = json.loads((research / snapshot).read_text())
out = args.output.resolve()
out.mkdir(parents=True, exist_ok=False)
archive = subprocess.check_output(["git", "-C", str(repo), "archive", source["base_tag"]])
with tarfile.open(fileobj=io.BytesIO(archive)) as archive_file:
    archive_file.extractall(out, filter="data")
for name, content in source["files"].items():
    target = out / name
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content)
checks = []
for name, command in (
    ("tests", ["go", "test", "-count=1", "./..."]),
    ("build", ["go", "build", "-o", "pandafish-replayed", "."]),
):
    result = subprocess.run(command, cwd=out,
                            env=dict(os.environ, GOCACHE="/tmp/pandafish-v016-go-cache"),
                            stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
    (out / (name + ".txt")).write_text(result.stdout)
    checks.append({"name": name, "command": command, "exit": result.returncode})
    (out / "replay-checks.json").write_text(json.dumps(checks, indent=2) + "\n")
    print(name, result.returncode, flush=True)
    if result.returncode:
        raise SystemExit(result.returncode)
print("Restored source and replay outputs:", out)
