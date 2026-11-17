"""Audit the new candidate evidence without changing the released archives."""
import ast,hashlib,json,re,subprocess
from pathlib import Path
r=Path(__file__).resolve().parent;repo=r.parents[3];checks={}
for p in r.glob('*.json'):json.loads(p.read_text())
for p in r.glob('*.py'):ast.parse(p.read_text(),filename=str(p))
checks['json_and_python_syntax']=True
def screened(name):
 text=(repo/name).read_text()
 if name in ['main.go','main_test.go'] and 'Pandafish v0.17-search-moves-dev' not in text:
  text=text.replace('Pandafish v0.17','Pandafish v0.17-search-moves-dev')
 return text.encode()
for name,entry in json.loads((r/'source-hashes.json').read_text()).items():assert hashlib.sha256(screened(name)).hexdigest()==entry['sha256'],name
for name,text in json.loads((r/'candidate-source.json').read_text())['files'].items():assert screened(name).decode()==text,name
checks['screened_candidate_matches_except_release_identity']=True
for name,entry in json.loads((r/'upstream-source.json').read_text())['files'].items():assert hashlib.sha256((repo/'third_party/chess'/name).read_bytes()).hexdigest()==entry['sha256'],name
checks['original_upstream_runtime_and_license_byte_preserved']=True
for name,unchanged in json.loads((r/'unchanged-policy.json').read_text()).items():assert unchanged and (repo/name).read_bytes()==subprocess.check_output(['git','show','v0.16:'+name]),name
checks['evaluation_search_policy_and_clocks_unchanged']=True
prior=0
for directory in ['v0.15.3-architecture','v0.16-bishop-pair','v0.16-next-assessment','v0.16-piece-safety','v0.16-search-core']:
 root=r.parent/directory
 for name,entry in json.loads((root/'inventory.json').read_text()).items():
  data=(root/name).read_bytes();assert len(data)==entry['bytes']and hashlib.sha256(data).hexdigest()==entry['sha256'],(directory,name);prior+=1
checks['prior_archive_entries_verified']=prior
analysis=json.loads((r/'analysis.json').read_text());assert analysis['fixed']['rows']==168 and analysis['deep']['rows']==24
assert not analysis['fixed']['failures']and not analysis['deep']['failures'];assert analysis['uci']['requests']==2046
assert all(not b['failures']for b in analysis['uci']['blocks'].values())
assert json.loads((r/'oracle-results.json').read_text())=={'failures':[],'passed':596,'positions':596}
checks['node_score_move_and_independent_oracle_gates']=True
review=json.loads((r/'screen-review.json').read_text());assert sum(x['pgn_legal_games']for x in review['matches'])==150
checks['all_150_pgns_independently_replayed']=True
assert json.loads((r/'archive-replay-checks.json').read_text())['exit']==0
checks['source_overlay_replay']=True
reports=list(r.glob('*.md'))+[repo/'tools/benchmark/README.md',repo/'tools/benchmark/results/v0.17-search-moves-dev-vs-v0.16-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md']
for p in reports:
 for target in re.findall(r'\]\(([^)]+)\)',p.read_text()):
  if '://'in target or target.startswith('#'):continue
  assert(p.parent/target.split('#')[0]).exists(),(p,target)
 for i,line in enumerate(p.read_text().splitlines(),1):assert line==line.rstrip(),(p,i)
checks['links_and_prose_whitespace']=True
p=subprocess.run(['git','diff','--check'],cwd=repo,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);assert p.returncode==0,p.stdout
checks['git_diff_check']=True
(r/'final-audit.json').write_text(json.dumps(checks,indent=2)+'\n')
entries={}
for p in sorted(r.rglob('*')):
 if not p.is_file()or p.name=='inventory.json'or'__pycache__'in p.parts:continue
 data=p.read_bytes();entries[p.relative_to(r).as_posix()]={'sha256':hashlib.sha256(data).hexdigest(),'bytes':len(data)}
(r/'inventory.json').write_text(json.dumps(entries,indent=2)+'\n')
print('PASS:',len(entries),'new archive files;',prior,'prior entries intact; source, license, raw results, replay and links verified.')
