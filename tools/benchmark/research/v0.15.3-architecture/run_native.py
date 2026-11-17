"""Corrected validation, native profile, and serial UCI measurements."""
import json,os,subprocess,time
from pathlib import Path
r=Path(__file__).resolve().parent;env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',ARCH_INPUT=str(r/'inputs.json'));manifest=json.loads((r/'worktrees.json').read_text());checks=[]
def run(label,args,extra=None):
 start=time.monotonic();p=subprocess.run(args,env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(r/(label+'.txt')).write_text(p.stdout);checks.append({'step':label,'command':args,'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'native-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(label,p.returncode,flush=True);return p.returncode
for row in manifest:run(row['label']+'-tests-corrected',['go','-C',row['tree'],'test','-count=1','./...'])
base=next(x['tree']for x in manifest if x['label']=='baseline')
run('native-profile',['go','-C',base,'test','-run','^$','-bench','^BenchmarkArchitectureProfile$','-benchtime=10s','-memprofilerate=4096','-cpuprofile',str(r/'native-cpu.pb.gz'),'-memprofile',str(r/'native-memory.pb.gz'),'-timeout','3m'])
for label,args in [('cpu-flat',['-top','-nodecount=45']),('cpu-cumulative',['-top','-cum','-nodecount=45'])]:run(label,['go','tool','pprof',*args,str(r/'native-cpu.pb.gz')])
for sample in ['alloc_space','alloc_objects']:run(sample,['go','tool','pprof','-sample_index='+sample,'-top','-nodecount=40',str(r/'native-memory.pb.gz')])
for mode in ['100ms','1000ms','warm']:
 if run('uci-'+mode+'-check',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'uci_probes.py'),mode]):break
