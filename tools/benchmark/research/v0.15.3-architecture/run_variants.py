"""Serial compiler/check/probe driver. Failed diagnostics are preserved."""
import json,os,subprocess,time
from pathlib import Path
r=Path(__file__).resolve().parent;repo=Path.cwd();env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',ARCH_INPUT=str(r/'inputs.json'));checks=[]
manifest=json.loads((r/'worktrees.json').read_text())
def run(label,step,args,extra=None):
 started=time.monotonic();p=subprocess.run(args,env=dict(env,**(extra or {})),text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);(r/f'{label}-{step}.txt').write_text(p.stdout);checks.append({'variant':label,'step':step,'command':args,'exit':p.returncode,'elapsed_s':time.monotonic()-started});(r/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(label,step,p.returncode,flush=True);return p.returncode
for row in manifest:
 label,tree=row['label'],row['tree'];binary=str(repo/'tools/benchmark/bin'/('architecture-'+label))
 if run(label,'build',['go','-C',tree,'build','-o',binary,'.']):continue
 # Tests may intentionally reject changed score/tie semantics; never hide them.
 if label!='instrumented-baseline':
  run(label,'tests',['go','-C',tree,'test','-count=1','./...'])
  run(label,'ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_OUTPUT':str(r/f'{label}-ghi.json')})
  for mode in ['fixed','warm']:
   run(label,mode,['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','12m','-v'],{'ARCH_MODE':mode,'ARCH_OUTPUT':str(r/f'{label}-{mode}.json')})
 else:run(label,'audit',['go','-C',tree,'test','-count=1','-run','^TestArchitectureAudit$','-timeout','12m','-v'],{'ARCH_OUTPUT':str(r/'cache-audit.json')})
