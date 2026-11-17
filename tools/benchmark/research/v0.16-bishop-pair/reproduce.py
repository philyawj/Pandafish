"""Offline reconstruction in a new directory; never runs matches or live play."""
import json,os,subprocess,sys,time
from pathlib import Path
repo=Path.cwd();r=Path(__file__).resolve().parent
if len(sys.argv) not in [2,3] or len(sys.argv)==3 and sys.argv[2]!='--checks-only':raise SystemExit('Usage: python3 reproduce.py /absolute/new/output-dir [--checks-only], from repo root')
checks_only=len(sys.argv)==3
out=Path(sys.argv[1]).resolve();out.mkdir(parents=True,exist_ok=False)
p=r/'candidate-source-final.json';overlay=json.loads((p if p.exists() else r/'candidate-source.json').read_text());archive=subprocess.check_output(['git','archive',overlay['base_tag']]);env=dict(os.environ,GOCACHE=os.environ.get('GOCACHE','/tmp/pandafish-v016-go-cache'));checks=[]
def run(name,args,extra=None):
 started=time.monotonic();p=subprocess.run([str(a) for a in args],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(out/(name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':[str(a) for a in args],'exit':p.returncode,'elapsed_s':time.monotonic()-started});(out/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit('Replay failed; output preserved in '+str(out))
for label in ['baseline','candidate']:
 tree=out/label;tree.mkdir();subprocess.run(['tar','-x','-C',str(tree)],input=archive,check=True)
 if label=='candidate':
  for name,text in overlay['files'].items():(tree/name).write_text(text)
 for name in ['probe','cost']:(tree/('research_'+name+'_test.go')).write_text((r/(name+'_test.go.txt')).read_text())
 run(label+'-build',['go','-C',tree,'build','-o',out/('pandafish-'+label),'.'])
 run(label+'-controls',['go','-C',tree,'test','-run','TestBishopPair|TestPairMaterial|TestPlacementPreferences|TestEvaluationColorSymmetry','-v','.'])
 if checks_only:continue
 run(label+'-fixed',['go','-C',tree,'test','-count=1','-run','TestResearchProbe','-timeout','8m'],{'PROBE_INPUT':str(r/'probe-inputs.json'),'PROBE_OUTPUT':str(out/(label+'-fixed.json')),'PROBE_FIXED_ONLY':'1'})
 run(label+'-static-cost',['go','-C',tree,'test','-run','^$','-bench','BenchmarkResearchStatic','-benchtime=300ms','-count=3'])
if checks_only:
 print('Archive build and controls passed:',out);raise SystemExit(0)
for trial,label in enumerate(['baseline','candidate','candidate','baseline'],1):
 run(f'clean-{trial}-{label}',['go','-C',out/label,'test','-count=1','-run','TestResearchProbe','-timeout','3m'],{'PROBE_INPUT':str(r/'probe-inputs.json'),'PROBE_OUTPUT':str(out/f'clean-{trial}-{label}.json'),'PROBE_TIMED_ONLY':'1'})
print('Replay finished without matches or live games:',out)
