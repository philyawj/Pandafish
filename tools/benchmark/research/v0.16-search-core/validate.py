"""Preserve every check output. Run once from the repository root."""
import hashlib, io, json, os, subprocess, tarfile, tempfile, time
from pathlib import Path
root=Path.cwd();r=root/'tools/benchmark/research/v0.16-search-core'
env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache')
checks=[]
def run(name,cmd,extra=None):
 output=r/(name+'.txt')
 if output.exists():raise RuntimeError('Refuse overwrite: '+str(output))
 start=time.monotonic()
 p=subprocess.run([str(x)for x in cmd],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 output.write_text(p.stdout)
 checks.append(dict(name=name,command=[str(x)for x in cmd],exit=p.returncode,elapsed_s=time.monotonic()-start))
 (r/'checks.json').write_text(json.dumps(checks,indent=2)+'\n')
 print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
 return p.stdout
# No timed performance blocks run while these checks execute.
formatted=run('format',['gofmt','-l',*[p for p in root.glob('*.go')]])
assert not formatted.strip(),formatted
run('race-tests',['go','test','-race','-count=1','./...'])
run('vet',['go','vet','./...'])
run('candidate-build',['go','build','-o',root/'tools/benchmark/bin/pandafish-v0.16-search-core-dev','.'])
run('arm64-build',['go','build','-o','/tmp/pandafish-v016-search-core-arm64','.'],{'GOOS':'darwin','GOARCH':'arm64'})
run('runner-build',['go','build','-o',root/'tools/benchmark/bin/benchmark','./tools/benchmark'])
scratch=Path(tempfile.mkdtemp(prefix='pandafish-v016-core-'))
archive=subprocess.check_output(['git','archive','v0.15.3'])
for label in ['baseline','candidate']:
 tree=scratch/label;tree.mkdir()
 with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(tree,filter='data')
 if label=='candidate':
  for p in root.glob('*.go'):(tree/p.name).write_bytes(p.read_bytes())
 (tree/'research_probe_test.go').write_bytes((r/'driver_test.go.txt').read_bytes())
 if label=='baseline':run('baseline-build',['go','-C',tree,'build','-o',root/'tools/benchmark/bin/pandafish-v0.15.3','.'])
# Freeze measured source, including every test; do not change runtime after this.
sources={p.name:{'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'bytes':p.stat().st_size}for p in sorted(root.glob('*.go'))}
(r/'source-hashes.json').write_text(json.dumps(sources,indent=2)+'\n')
(r/'worktrees.json').write_text(json.dumps({label:str(scratch/label)for label in ['baseline','candidate']},indent=2)+'\n')
run('offline-wrapper',['/Users/philyawj/projects/lichess-bot/venv/bin/python',r/'wrapper-check.py'])
for label in ['baseline','candidate']:
 tree=scratch/label
 for mode,inputs,depths in [('fixed','inputs.json','3,4'),('deep','deep-inputs.json','5,6')]:
  run(label+'-'+mode,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','15m','-v'],{'ARCH_INPUT':str(r/inputs),'ARCH_OUTPUT':str(r/(label+'-'+mode+'.json')),'ARCH_MODE':'fixed','ARCH_DEPTHS':depths})
 run(label+'-ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/(label+'-ghi.json'))})
print('All correctness/fixed-depth checks saved.',flush=True)
