import json,os,subprocess,time
from pathlib import Path
r=Path(__file__).resolve().parent;env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',ARCH_INPUT=str(r/'inputs.json'));m=json.loads((r/'worktrees.json').read_text());checks=[]
def run(name,args,extra=None):
 start=time.monotonic();p=subprocess.run(args,env=dict(env,**(extra or {})),text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);(r/(name+'.txt')).write_text(p.stdout);checks.append({'step':name,'command':args,'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'followup-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);return p.returncode
base=next(x['tree']for x in m if x['label']=='baseline')
run('cpu-clean-profile',['go','-C',base,'test','-run','^$','-bench','^BenchmarkArchitectureProfile$','-benchtime=10s','-cpuprofile',str(r/'native-cpu-clean.pb.gz'),'-timeout','3m'],{'ARCH_NO_MEM_PROFILE':'1'})
for name,flags in [('cpu-clean-flat',['-top']),('cpu-clean-cumulative',['-top','-cum'])]:run(name,['go','tool','pprof',*flags,'-nodecount=50',str(r/'native-cpu-clean.pb.gz')])
for mode in ['historical','historical-fixed']:run('uci-'+mode+'-check',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'uci_probes.py'),mode])
for label in ['baseline','safe-cache','unsafe-board-scores','root-single-best','tactical-ordering']:
 tree=next(x['tree']for x in m if x['label']==label);run(label+'-deep-check',['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','8m','-v'],{'ARCH_INPUT':str(r/'deep-inputs.json'),'ARCH_MODE':'fixed','ARCH_DEPTHS':'5,6','ARCH_OUTPUT':str(r/(label+'-deep.json'))})
run('referee-check',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'referee.py')])
