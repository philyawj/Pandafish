"""Verify source/data consistency and local links; refresh this archive inventory."""
import ast,hashlib,json,re,subprocess
from pathlib import Path
r=Path(__file__).resolve().parent;repo=r.parents[3];checks={}
for p in r.glob('*.json'):json.loads(p.read_text())
checks['json_parse']=True
for p in r.glob('*.py'):ast.parse(p.read_text(),filename=str(p))
checks['python_syntax']=True
source=json.loads((r/'source-hashes.json').read_text())
def screened_bytes(name):
 text=(repo/name).read_text()
 if name in ['main.go','main_test.go'] and 'Pandafish v0.16-search-core-dev' not in text:
  text=text.replace('Pandafish v0.16','Pandafish v0.16-search-core-dev')
 return text.encode()
for name,entry in source.items():assert hashlib.sha256(screened_bytes(name)).hexdigest()==entry['sha256'],name
for name,text in json.loads((r/'candidate-source.json').read_text())['files'].items():assert screened_bytes(name).decode()==text,name
checks['frozen_source_matches_working_tree_and_overlay_except_release_identity']=True
prior=0
for name in ['v0.15.3-architecture','v0.16-bishop-pair','v0.16-next-assessment','v0.16-piece-safety']:
 root=r.parent/name
 for key,entry in json.loads((root/'inventory.json').read_text()).items():
  raw=(root/key).read_bytes();assert hashlib.sha256(raw).hexdigest()==entry['sha256'],(name,key);assert len(raw)==entry['bytes'];prior+=1
checks['previous_archive_entries_verified']=prior
reports = list(r.glob('*.md')) + [repo/'tools/benchmark/README.md', repo/'tools/benchmark/results/v0.16-search-core-dev-vs-v0.15.3-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md']
for p in reports:
 for target in re.findall(r'\]\(([^)]+)\)',p.read_text()):
  if '://'in target or target.startswith('#'):continue
  target=target.split('#')[0]
  assert (p.parent/target).exists(),(p,target)
checks['local_document_links']=True
fixed=json.loads((r/'analysis.json').read_text());assert fixed['fixed']['rows']==168 and fixed['deep']['rows']==24
assert not fixed['fixed']['failures']and not fixed['deep']['failures']
assert fixed['uci']['requests']==2046 and all(not x['failures']for x in fixed['uci']['blocks'].values())
checks['fixed_depth_and_native_uci_gates']=True
review=json.loads((r/'screen-review.json').read_text())
assert sum(x['pgn_legal_games']for x in review['matches'])==150
assert all(x['stats']['forfeits']==0 and x['stats']['unfinished']==0 for x in review['matches'])
checks['pgn_legality_and_reliability']=True
assert json.loads((r/'archive-replay-checks.json').read_text())['exit']==0
checks['archived_source_replay']=True
p=subprocess.run(['git','diff','--check'],cwd=repo,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
checks['git_diff_check']={'exit':p.returncode,'output':p.stdout};assert p.returncode==0,p.stdout
# Check newly created prose/scripts too, since git diff omits untracked files.
for p in list(r.glob('*.py'))+reports+[repo/'search_context.go',repo/'search_core_test.go']:
 for i,line in enumerate(p.read_text().splitlines(),1):assert line==line.rstrip(),(p,i)
checks['new_source_and_prose_whitespace']=True
(r/'final-audit.json').write_text(json.dumps(checks,indent=2)+'\n')
entries={}
for p in sorted(r.rglob('*')):
 if not p.is_file()or p.name=='inventory.json'or '__pycache__'in p.parts:continue
 raw=p.read_bytes();entries[p.relative_to(r).as_posix()]={'sha256':hashlib.sha256(raw).hexdigest(),'bytes':len(raw)}
(r/'inventory.json').write_text(json.dumps(entries,indent=2)+'\n')
print('PASS:',len(entries),'archive files;',prior,'prior archive checksums; frozen source, PGNs, replay, syntax, links and whitespace.')
