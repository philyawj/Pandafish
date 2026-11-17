"""Restore frozen candidate from v0.15.3 in a NEW directory. No matches/live play."""
import argparse, io, json, os, subprocess, tarfile
from pathlib import Path
parser=argparse.ArgumentParser(description=__doc__)
parser.add_argument('output',type=Path)
parser.add_argument('--fixed',action='store_true',help='also run the 84-case depth 3/4 diagnostic')
args=parser.parse_args();out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
r=Path(__file__).resolve().parent
archive=subprocess.check_output(['git','archive','v0.15.3'])
with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(out,filter='data')
for name,source in json.loads((r/'candidate-source.json').read_text())['files'].items():(out/name).write_text(source)
env=dict(os.environ,GOCACHE=os.environ.get('GOCACHE','/tmp/pandafish-v016-go-cache'))
checks=[]
def run(name,cmd,extra=None):
 p=subprocess.run([str(x)for x in cmd],stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,env=dict(env,**(extra or {})))
 (out/(name+'.txt')).write_text(p.stdout);checks.append({'step':name,'command':[str(x)for x in cmd],'exit':p.returncode})
 (out/'replay-checks.json').write_text(json.dumps(checks,indent=2)+'\n')
 print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
run('tests',['go','-C',out,'test','-count=1','./...'])
run('build',['go','-C',out,'build','-o',out/'pandafish-candidate','.'])
if args.fixed:
 (out/'research_probe_test.go').write_bytes((r/'driver_test.go.txt').read_bytes())
 run('fixed',['go','-C',out,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','15m','-v'],{'ARCH_INPUT':str(r/'inputs.json'),'ARCH_OUTPUT':str(out/'fixed.json'),'ARCH_MODE':'fixed','ARCH_DEPTHS':'3,4'})
print('Replay and all outputs:',out)
