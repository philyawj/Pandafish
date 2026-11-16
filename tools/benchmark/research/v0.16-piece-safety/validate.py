"""Record offline candidate checks; run only after performance probes finish."""
import json,os,subprocess,time
from pathlib import Path
root=Path.cwd();r=root/'tools/benchmark/research/v0.16-piece-safety';env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache');results=[]
steps=[('focused',['go','test','-count=1','-run','TestPieceSafety|TestSafetySummary','-v','.'],{}),('tests',['go','test','-count=1','./...'],{}),('race',['go','test','-race','-count=1','./...'],{}),('vet',['go','vet','./...'],{}),('arm64',['go','build','-o','/tmp/pandafish-v016-arm64','.'],{'GOOS':'darwin','GOARCH':'arm64'}),('candidate-build',['go','build','-o','tools/benchmark/bin/pandafish-v0.16-piece-safety-dev','.'],{}),('wrapper',['/Users/philyawj/projects/lichess-bot/venv/bin/python',str(r/'wrapper-check.py')],{})]
for name,args,extra in steps:
 start=time.monotonic();p=subprocess.run(args,env=dict(env,**extra),text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 (r/(name+'-check.txt')).write_text(p.stdout);results.append({'name':name,'args':args,'extra_env':extra,'exit_code':p.returncode,'elapsed_s':time.monotonic()-start});(r/'candidate-checks.json').write_text(json.dumps(results,indent=2)+'\n');print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
