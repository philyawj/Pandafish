"""Resume candidate setup after the preserved scratch-copy harness failure."""
import hashlib,json,os,shutil,subprocess,time
from pathlib import Path
repo=Path.cwd();r=repo/'tools/benchmark/research/v0.18-evaluation-model';checks=json.loads((r/'checks.json').read_text());baseline=Path(next(x['command'][2]for x in checks if x['name']=='baseline-fixed'));scratch=baseline.parent;tree=scratch/'candidate'
shutil.copytree(repo/'third_party',tree/'third_party',dirs_exist_ok=True)
for p in repo.glob('*.go'):(tree/p.name).write_bytes(p.read_bytes())
target=tree/'tools/benchmark/research/v0.18-evaluation-model';target.mkdir(parents=True,exist_ok=True);shutil.copyfile(r/'scalar-expectations.json',target/'scalar-expectations.json')
(tree/'research_probe_test.go').write_bytes((r/'driver_test.go.txt').read_bytes())
def run(name,cmd,extra):
 assert not(r/(name+'.txt')).exists(),name
 start=time.monotonic();p=subprocess.run([str(x)for x in cmd],env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',**extra),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(r/(name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':[str(x)for x in cmd],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);assert p.returncode==0
for mode,inputs,depths in [('fixed','inputs.json','3,4'),('deep','deep-inputs.json','5,6')]:run('candidate-'+mode,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','15m','-v'],{'ARCH_INPUT':str(r/inputs),'ARCH_OUTPUT':str(r/('candidate-'+mode+'.json')),'ARCH_MODE':'fixed','ARCH_DEPTHS':depths})
run('candidate-ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/'candidate-ghi.json')})
(r/'worktrees.json').write_text(json.dumps({name:str(scratch/name)for name in ['baseline','candidate','library']},indent=2)+'\n')
paths=list(repo.glob('*.go'))+[repo/'go.mod',repo/'go.sum',r/'scalar-expectations.json']+list((repo/'third_party/chess').glob('*'))
(r/'source-hashes.json').write_text(json.dumps({str(p.relative_to(repo)):{'sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'bytes':p.stat().st_size}for p in paths if p.is_file()},indent=2)+'\n')
print('Candidate validation resumed; all earlier measurements retained.')
