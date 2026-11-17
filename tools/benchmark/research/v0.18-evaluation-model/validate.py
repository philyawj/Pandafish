"""Correctness and fixed-depth backend comparison against tagged v0.16."""
import hashlib,io,json,os,shutil,subprocess,tarfile,tempfile,time
from pathlib import Path
repo=Path.cwd();r=repo/'tools/benchmark/research/v0.18-evaluation-model';env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache');checks=[]
def run(name,cmd,extra=None):
 output=r/(name+'.txt');assert not output.exists(),name
 start=time.monotonic();p=subprocess.run([str(x)for x in cmd],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 output.write_text(p.stdout);checks.append({'name':name,'command':[str(x)for x in cmd],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
 return p.stdout
assert not(r/'checks.json').exists()
files=subprocess.check_output(['rg','--files','-g','*.go','-g','!tools/benchmark/bin/**','-g','!third_party/**'],text=True).splitlines()
files.extend(['third_party/chess/search_moves.go','third_party/chess/eval_features.go'])
assert not run('format',['gofmt','-l',*files]).strip()
run('race',['go','test','-race','-count=1','./...'])
run('vet',['go','vet','./...'])
run('candidate-build',['go','build','-o',repo/'tools/benchmark/bin/pandafish-v0.18-evaluation-model-dev','.'])
run('arm64',['go','build','-o','/tmp/pandafish-v017-search-moves-arm64','.'],{'GOOS':'darwin','GOARCH':'arm64'})
run('wrapper',['/Users/philyawj/projects/lichess-bot/venv/bin/python',r/'wrapper-check.py'])
# Replay the original library's tests/fixtures against our additive API source.
scratch=Path(tempfile.mkdtemp(prefix='pandafish-v017-check-'));upstream=Path('/Users/philyawj/go/pkg/mod/github.com/corentings/chess/v2@v2.6.1')
lib=scratch/'library';shutil.copytree(repo/'third_party/chess',lib)
for p in upstream.glob('*_test.go'):shutil.copyfile(p,lib/p.name)
shutil.copytree(upstream/'fixtures',lib/'fixtures')
run('upstream-library-tests',['go','-C',lib,'test','-race','-count=1','./...'])
# All copied runtime files are unchanged. Only search_moves.go is additive.
for name,entry in json.loads((r/'upstream-source.json').read_text())['files'].items():assert hashlib.sha256((repo/'third_party/chess'/name).read_bytes()).hexdigest()==entry['sha256'],name
archive=subprocess.check_output(['git','archive','v0.17'])
for label in ['baseline','candidate']:
 tree=scratch/label;tree.mkdir()
 with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(tree,filter='data')
 if label=='candidate':
  for p in repo.glob('*.go'):(tree/p.name).write_bytes(p.read_bytes())
  for name in ['go.mod','go.sum']:(tree/name).write_bytes((repo/name).read_bytes())
  shutil.copytree(repo/'third_party',tree/'third_party',dirs_exist_ok=True)
  target=tree/'tools/benchmark/research/v0.18-evaluation-model';target.mkdir(parents=True,exist_ok=True);shutil.copyfile(r/'scalar-expectations.json',target/'scalar-expectations.json')
 (tree/'research_probe_test.go').write_bytes((r/'driver_test.go.txt').read_bytes())
 for mode,inputs,depths in [('fixed','inputs.json','3,4'),('deep','deep-inputs.json','5,6')]:
  run(label+'-'+mode,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','15m','-v'],{'ARCH_INPUT':str(r/inputs),'ARCH_OUTPUT':str(r/(label+'-'+mode+'.json')),'ARCH_MODE':'fixed','ARCH_DEPTHS':depths})
 run(label+'-ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/(label+'-ghi.json'))})
(r/'worktrees.json').write_text(json.dumps({name:str(scratch/name)for name in ['baseline','candidate','library']},indent=2)+'\n')
paths=list(repo.glob('*.go'))+[repo/'go.mod',repo/'go.sum']+list((repo/'third_party/chess').glob('*'))
hashes={str(p.relative_to(repo)):{'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'bytes':p.stat().st_size}for p in paths if p.is_file()}
(r/'source-hashes.json').write_text(json.dumps(hashes,indent=2)+'\n')
print('Correctness and fixed-depth comparisons complete.',flush=True)
