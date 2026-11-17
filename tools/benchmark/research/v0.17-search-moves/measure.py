"""Run independent legality oracle, then isolated candidate CPU/allocation profile."""
import json,os,subprocess,time
from pathlib import Path
r=Path(__file__).resolve().parent;trees=json.loads((r/'worktrees.json').read_text());tree=Path(trees['candidate']);env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache')
(tree/'research_oracle_test.go').write_bytes((r/'oracle_test.go.txt').read_bytes())
checks=[]
def run(name,cmd,extra=None):
 assert not(r/(name+'.txt')).exists()
 start=time.monotonic();p=subprocess.run([str(x)for x in cmd],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 (r/(name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':[str(x)for x in cmd],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'measurement-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True);assert p.returncode==0,p.stdout
run('oracle',['go','-C',tree,'test','-count=1','-run','^TestSearchMoveOracle$','-timeout','5m','-v'],{'ORACLE_INPUT':str(r/'oracle-inputs.json'),'ORACLE_OUTPUT':str(r/'oracle-results.json')})
run('candidate-profile',['go','-C',tree,'test','-run','^$','-bench','^BenchmarkArchitectureProfile$','-benchtime','60x','-cpuprofile',r/'candidate-cpu.pb.gz','-o',tree/'profile.test'],{'ARCH_INPUT':str(r/'inputs.json'),'ARCH_NO_MEM_PROFILE':'1'})
for label,extra in [('flat',[]),('cumulative',['-cum'])]:run('candidate-cpu-'+label,['go','tool','pprof','-top',*extra,tree/'profile.test',r/'candidate-cpu.pb.gz'])
