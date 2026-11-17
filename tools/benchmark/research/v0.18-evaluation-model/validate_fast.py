"""One exact-cost revision and fast-base ablation, preserving the slow candidate."""
import hashlib,io,json,os,shutil,subprocess,tarfile,tempfile,time
from pathlib import Path
repo=Path.cwd();root=repo/'tools/benchmark/research/v0.18-evaluation-model';r=root/'fast-kernel';env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache');checks=[]
assert not(r/'checks.json').exists()
def run(name,cmd,extra=None):
 start=time.monotonic();p=subprocess.run([str(x)for x in cmd],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(r/(name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':[str(x)for x in cmd],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);assert p.returncode==0,p.stdout
run('race',['go','test','-race','-count=1','./...']);run('vet',['go','vet','./...'])
scratch=Path(tempfile.mkdtemp(prefix='pandafish-v018-fast-'));archive=subprocess.check_output(['git','archive','v0.17'])
for label in ['candidate','fast-base']:
 tree=scratch/label;tree.mkdir()
 with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(tree,filter='data')
 for p in repo.glob('*.go'):(tree/p.name).write_bytes(p.read_bytes())
 shutil.copytree(repo/'third_party',tree/'third_party',dirs_exist_ok=True)
 fixture=tree/'tools/benchmark/research/v0.18-evaluation-model';fixture.mkdir(parents=True,exist_ok=True);shutil.copyfile(root/'scalar-expectations.json',fixture/'scalar-expectations.json')
 (tree/'research_probe_test.go').write_bytes((root/'driver_test.go.txt').read_bytes())
 if label=='fast-base':
  p=tree/'evaluation.go';p.write_text(p.read_text().replace(' + modelAdjustment(position, perspective)',''))
  for name in ['main.go','main_test.go']:
   p=tree/name;p.write_text(p.read_text().replace('v0.18-evaluation-model-fast-dev','v0.18-fast-base-diagnostic'))
 binary=repo/'tools/benchmark/bin'/('pandafish-v0.18-evaluation-model-fast-dev'if label=='candidate'else'pandafish-v0.18-fast-base-diagnostic')
 run(label+'-build',['go','-C',tree,'build','-o',binary,'.'])
 for mode,inputs,depths in [('fixed','inputs.json','3,4')]+([('deep','deep-inputs.json','5,6')]if label=='candidate'else[]):
  run(label+'-'+mode,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','15m','-v'],{'ARCH_INPUT':str(root/inputs),'ARCH_OUTPUT':str(r/(label+'-'+mode+'.json')),'ARCH_MODE':'fixed','ARCH_DEPTHS':depths})
 run(label+'-ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/(label+'-ghi.json'))})
# Strong semantic controls distinguish cost from model changes.
for label,prefix in [('candidate','candidate'),('fast-base','baseline')]:
 old={(x['ID'],x['Depth']):x for x in json.loads((root/(prefix+'-fixed.json')).read_text())}
 for x in json.loads((r/(label+'-fixed.json')).read_text()):
  a=old[x['ID'],x['Depth']];assert x['Complete']and x['Score']==a['Score']and x['Nodes']==a['Nodes']and x['Best']==a['Best'],(label,x['ID'])
for x,a in zip(json.loads((r/'candidate-deep.json').read_text()),json.loads((root/'candidate-deep.json').read_text())):assert x['Complete']and x['Score']==a['Score']and x['Nodes']==a['Nodes']and x['Best']==a['Best'],x['ID']
(r/'worktrees.json').write_text(json.dumps({label:str(scratch/label)for label in ['candidate','fast-base']},indent=2)+'\n')
paths=list(repo.glob('*.go'))+[repo/'go.mod',repo/'go.sum',root/'scalar-expectations.json']+list((repo/'third_party/chess').glob('*'))
(r/'source-hashes.json').write_text(json.dumps({str(p.relative_to(repo)):{'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'bytes':p.stat().st_size}for p in paths if p.is_file()},indent=2)+'\n')
print('Exact cost revision and base-only ablation checks pass.')
