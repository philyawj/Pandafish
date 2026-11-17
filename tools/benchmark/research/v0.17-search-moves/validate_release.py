"""Validate identity-only release preparation without another strength screen."""
import hashlib,json,os,subprocess,time
from pathlib import Path
repo=Path.cwd();r=repo/'tools/benchmark/research/v0.17-search-moves';env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache')
if(r/'release-checks.json').exists():raise SystemExit('Refuse overwrite')
checks=[]
def run(name,cmd,extra=None):
 start=time.monotonic();p=subprocess.run([str(x)for x in cmd],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 (r/('release-'+name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':[str(x)for x in cmd],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'release-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);return p
sources={}
for name,entry in json.loads((r/'source-hashes.json').read_text()).items():
 data=(repo/name).read_bytes();normalized=data
 if name in ['main.go','main_test.go']:normalized=data.replace(b'Pandafish v0.17\\n',b'Pandafish v0.17-search-moves-dev\\n')
 assert hashlib.sha256(normalized).hexdigest()==entry['sha256'],name
 sources[name]={'sha256':hashlib.sha256(data).hexdigest(),'screened_sha256':entry['sha256'],'identity_only_change':data!=normalized}
(r/'release-source-identity.json').write_text(json.dumps(sources,indent=2)+'\n')
files=subprocess.check_output(['rg','--files','-g','*.go','-g','!tools/benchmark/bin/**','-g','!third_party/**'],text=True).splitlines()
files.append('third_party/chess/search_moves.go')
p=run('format',['gofmt','-l',*files]);assert p.returncode==0 and not p.stdout.strip()
for name,cmd,extra in [
 ('race',['go','test','-race','-count=1','./...'],None),('vet',['go','vet','./...'],None),
 ('build',['go','build','-o','tools/benchmark/bin/pandafish-v0.17-release','.'],None),
 ('arm64',['go','build','-o','/tmp/pandafish-v016-release-arm64','.'],{'GOOS':'darwin','GOARCH':'arm64'}),
 ('wrapper',['/Users/philyawj/projects/lichess-bot/venv/bin/python',r/'release-wrapper.py'],None)]:
 assert run(name,cmd,extra).returncode==0,name
cmd=['tools/benchmark/bin/benchmark','--candidate','tools/benchmark/bin/pandafish-v0.17-release','--version','v0.17-release-smoke','--revision','4ae65b14f69d95a1ad498baa759fec5e6bd0a1a6-dirty','--previous','tools/benchmark/bin/pandafish-v0.16','--previous-version','v0.16','--previous-revision','4ae65b14f69d95a1ad498baa759fec5e6bd0a1a6','--stockfish','/opt/homebrew/bin/stockfish','--sacrifice','tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine','--pairs','1','--parallel-all','--results',str(r/'release-smoke')]
# The six-game direct gate may fail by sampling; preserve it and inspect integration.
p=run('smoke',cmd);assert p.returncode in [0,1]
paths=list((r/'release-smoke').glob('*/summary.json'));assert len(paths)==1
j=json.loads(paths[0].read_text());assert j['complete']and j['settings']['matchup_concurrency']==3
assert sum(len(m['games'])for m in j['matches'])==6
assert all(m['stats']['forfeits']==0 and m['stats']['unfinished']==0 for m in j['matches'])
(r/'release-provenance.json').write_text(json.dumps({'released_version':'v0.17','screened_candidate':'v0.17-search-moves-dev','only_Go_differences':'UCI identity and handshake expectation','release_binary_sha256':hashlib.sha256((repo/'tools/benchmark/bin/pandafish-v0.17-release').read_bytes()).hexdigest(),'smoke_summary':str(paths[0].relative_to(r)),'smoke_gate':j['previous_gate'],'smoke_is_integration_only':True,'decision':'I explicitly approved promotion from the standard 50-game screen, deferred the proposed larger confirmation, and authorized PR/merge/release plus reassessment of one next candidate.'},indent=2)+'\n')
print('Release validation complete; original screen untouched.',flush=True)
