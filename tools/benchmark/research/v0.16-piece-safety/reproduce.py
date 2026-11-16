"""Rebuild and replay in a NEW output directory. No matches or online play."""
import json,os,subprocess,sys,time
from pathlib import Path
repo=Path.cwd();source=Path(__file__).resolve().parent
if len(sys.argv) not in [2,3] or len(sys.argv)==3 and sys.argv[2]!='--checks-only':raise SystemExit('Usage: python3 reproduce.py /absolute/new/output-dir [--checks-only] (from repo root)')
checks_only=len(sys.argv)==3
out=Path(sys.argv[1]).resolve();out.mkdir(parents=True,exist_ok=False)
overlay=json.loads((source/'candidate-source.json').read_text())
env=dict(os.environ,GOCACHE=os.environ.get('GOCACHE','/tmp/pandafish-v016-go-cache'))
checks=[]
def run(label,args,extra=None):
 start=time.monotonic();p=subprocess.run([str(x) for x in args],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 (out/(label+'.txt')).write_text(p.stdout);checks.append({'name':label,'args':[str(x) for x in args],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(out/'replay-checks.json').write_text(json.dumps(checks,indent=2)+'\n')
 if p.returncode:raise SystemExit(f'{label} failed, output saved in {out}')
 print(label,'PASS',flush=True)
archive=subprocess.check_output(['git','archive',overlay['base_tag']])
for label in ['baseline','candidate']:
 tree=out/label;tree.mkdir();subprocess.run(['tar','-x','-C',str(tree)],input=archive,check=True)
 if label=='candidate':
  for name,text in overlay['files'].items():(tree/name).write_text(text)
 for name in ['trace','probe','cost']+( ['leaf'] if label=='candidate' else []):
  (tree/('research_'+name+'_test.go')).write_text((source/(name+'_test.go.txt')).read_text())
 run(label+'-build',['go','-C',tree,'build','-o',out/('pandafish-'+label),'.'])
 run(label+'-tests',['go','-C',tree,'test','-run','TestPieceSafety|TestSafetySummary','-v','.'])
 if checks_only:continue
 run(label+'-original-probes',['go','-C',tree,'test','-count=1','-run','TestResearchProbe','-timeout','8m'],{'PROBE_INPUT':str(source/'probe-inputs.json'),'PROBE_OUTPUT':str(out/(label+'-probes.json'))})
 run(label+'-static-cost',['go','-C',tree,'test','-run','^$','-bench','BenchmarkResearchStatic','-benchtime=300ms','-count=3'])
if checks_only:
 print('Archive build/control smoke finished:',out);raise SystemExit(0)
run('trace',['go','-C',out/'baseline','test','-count=1','-run','TestResearchTrace','-timeout','8m'],{'TRACE_INPUT':str(source/'trace-inputs.json'),'TRACE_OUTPUT':str(out/'baseline-traces.json')})
for trial,label in enumerate(['baseline','candidate','candidate','baseline'],1):
 run(f'clean-{trial}-{label}',['go','-C',out/label,'test','-count=1','-run','TestResearchProbe','-timeout','3m'],{'PROBE_INPUT':str(source/'probe-inputs.json'),'PROBE_OUTPUT':str(out/f'clean-{trial}-{label}.json'),'PROBE_TIMED_ONLY':'1'})
for label in ['baseline','candidate']:
 run(label+'-live-fullroot',['go','-C',out/label,'test','-count=1','-run','TestResearchProbe','-timeout','8m'],{'PROBE_INPUT':str(source/'live-fullroot-inputs.json'),'PROBE_OUTPUT':str(out/(label+'-live-fullroot.json')),'PROBE_LIVE':'1','PROBE_FIXED_ONLY':'1'})
run('leaf-safety',['go','-C',out/'candidate','test','-count=1','-run','TestResearchLeafSafety'],{'LEAF_INPUT':str(out/'baseline-traces.json'),'LEAF_OUTPUT':str(out/'leaf-safety.json')})
print('Finished offline replay:',out)
