"""Verify final source, earlier evidence, result bytes and authored links."""
import ast
import hashlib
import json
import re
import subprocess
from pathlib import Path

research = Path(__file__).resolve().parent
repo = research.parents[3]
issues = []

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

source = json.loads((research / 'source-hashes.json').read_text())
for name, entry in source.items():
    data = (repo / name).read_bytes()
    if name in ['main.go', 'main_test.go'] and b'Pandafish v0.18-exact-evaluation-dev' not in data:
        data = data.replace(b'Pandafish v0.18', b'Pandafish v0.18-exact-evaluation-dev')
    if hashlib.sha256(data).hexdigest() != entry['sha256']:
        issues.append('changed frozen source: ' + name)
old_archives = {}
for inventory in (repo / 'tools/benchmark/research').glob('*/inventory.json'):
    if inventory.parent == research:
        continue
    entries = json.loads(inventory.read_text())
    bad = [name for name, value in entries.items()
           if digest(inventory.parent / name) != value['sha256']]
    old_archives[str(inventory.parent.relative_to(repo))] = {
        'files': len(entries), 'mismatches': bad}
    issues.extend(str(inventory) + ': ' + name for name in bad)
for path in research.rglob('*.json'):
    json.loads(path.read_text())
for path in research.rglob('*.py'):
    ast.parse(path.read_text())
result = repo / 'tools/benchmark/results/v0.18-exact-evaluation-dev-vs-v0.17-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba'
markdown = [*research.glob('*.md'), result / 'README.md',
            repo / 'docs/CHESS_SEARCH_BACKEND.md', repo / 'README.md',
            repo / 'tools/benchmark/README.md']
for path in markdown:
    for target in re.findall(r'\]\(([^)]+)\)', path.read_text()):
        if '://' in target or target.startswith('#'):
            continue
        # This report is created at the end of this same audit.
        if (path.parent / target.split('#')[0]).resolve() == research / 'final-audit.json':
            continue
        if not (path.parent / target.split('#')[0]).exists():
            issues.append('missing link: ' + str(path) + ' ' + target)
    lines = path.read_text().splitlines()
    for i, line in enumerate(lines):
        if re.fullmatch(r'\|[ :|\-]+\|', line):
            if i == 0 or line.count('|') != lines[i - 1].count('|'):
                issues.append('table column mismatch: ' + str(path) + ':' + str(i + 1))
    if path.parent == research or path.parent.name == 'docs':
        if re.search(r'\bmaintainer\b', path.read_text(), re.I):
            issues.append('unwanted wording: ' + str(path))
provenance = json.loads((research / 'provenance.json').read_text())
paths = {
    'candidate': repo / 'tools/benchmark/bin/pandafish-v0.18-exact-evaluation-dev',
    'baseline': repo / 'tools/benchmark/bin/pandafish-v0.17',
    'stockfish': Path('/opt/homebrew/bin/stockfish'),
    'sacrifice': repo / 'tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine',
    'sacrifice_Tal': repo / 'tools/benchmark/bin/sacrifice-source/data/Tal.pgn',
}
for actor, path in paths.items():
    if digest(path) != provenance['binary_sha256'][actor]:
        issues.append('changed binary/asset: ' + actor)
model = research.parent / 'v0.18-evaluation-model'
weights_unchanged = digest(model / 'model-weights.json') == json.loads(
    (model / 'fast-kernel/provenance.json').read_text())['frozen_weights_sha256']
assert weights_unchanged
summary = json.loads((result / 'summary.json').read_text())
assert summary['complete'] and summary['settings']['matchup_concurrency'] == 3
assert summary['settings']['pairs'] == 25
assert all(m['stats']['forfeits'] == 0 and m['stats']['unfinished'] == 0
           and len(m['games']) == 50 for m in summary['matches'])
assert subprocess.check_output(['git', '-C', str(repo), 'rev-parse', 'v0.17^{}'],
                               text=True).strip() == provenance['base_commit']
report = {
    'current_source_files': len(source), 'frozen_model_weights_unchanged': weights_unchanged,
    'prior_archives': old_archives,
    'prior_inventory_entries': sum(x['files'] for x in old_archives.values()),
    'result_files': {str(p.relative_to(repo)): {'sha256': digest(p), 'bytes': p.stat().st_size}
                     for p in sorted(result.iterdir()) if p.is_file()},
    'screened_source_unchanged_except_release_identity_and_all_three_matchups_complete': True,
    'issues': issues,
}
(research / 'final-audit.json').write_text(json.dumps(report, indent=2) + '\n')
print(json.dumps({key: report[key] for key in
                  ['current_source_files', 'prior_inventory_entries', 'issues']}, indent=2))
assert not issues, issues
