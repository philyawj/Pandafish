import json,os,subprocess,time
from pathlib import Path
repo=Path.cwd();r=Path(__file__).resolve().parent;m=json.loads((r/'worktrees.json').read_text());env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',ARCH_INPUT=str(r/'inputs.json'));checks=[]
def run(name,args,extra=None):
 start=time.monotonic();p=subprocess.run(args,env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(r/(name+'.txt')).write_text(p.stdout);checks.append({'step':name,'command':args,'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'combo-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);return p.returncode
label='coordinated-kernel';tree=next(x['tree']for x in m if x['label']==label)
run('coordinated-build',['go','-C',tree,'build','-o',str(repo/'tools/benchmark/bin/architecture-coordinated-kernel'),'.'])
run('coordinated-tests',['go','-C',tree,'test','-count=1','./...'])
run('coordinated-ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/'coordinated-kernel-ghi.json')})
for name,inp,depths in [('fixed','inputs.json','3,4'),('deep','deep-inputs.json','5,6')]:run('coordinated-'+name,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','8m','-v'],{'ARCH_MODE':'fixed','ARCH_DEPTHS':depths,'ARCH_INPUT':str(r/inp),'ARCH_OUTPUT':str(r/('coordinated-kernel-'+name+'.json'))})
run('uci-combo-check',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'uci_probes.py'),'combo'])
base=next(x['tree']for x in m if x['label']=='baseline');run('leaf-trace-check',['go','-C',base,'test','-count=1','-run','^TestResearchTrace$','-timeout','10m','-v'],{'TRACE_INPUT':str(r/'leaf-inputs.json'),'TRACE_OUTPUT':str(r/'leaf-traces.json')})
if (r/'referee-results.json').exists():(r/'referee-before-combo.json').write_bytes((r/'referee-results.json').read_bytes())
run('referee-combo-check',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'referee.py')])
run('analysis-check',['python3',str(r/'analyse.py')])
