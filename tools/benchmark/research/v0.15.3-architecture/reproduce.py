"""Replay archived diagnostic overlays in a NEW directory. No matches or live play."""
import argparse,io,json,os,subprocess,tarfile,time
from pathlib import Path
parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('output',type=Path);parser.add_argument('--smoke',action='store_true',help='baseline/safe-cache/coordinated build, GHI and six-input shallow checks only');parser.add_argument('--variants',help='comma-separated labels instead of all overlays');args=parser.parse_args()
repo=Path.cwd();source=Path(__file__).resolve().parent;out=args.output.resolve();out.mkdir(parents=True,exist_ok=False);env=dict(os.environ,GOCACHE=os.environ.get('GOCACHE','/tmp/pandafish-v016-go-cache'));checks=[]
names=subprocess.check_output(['git','ls-tree','--name-only','v0.15.2'],text=True).splitlines();paths=[n for n in names if n.endswith('.go')or n in['go.mod','go.sum','internal','testdata']];archive=subprocess.check_output(['git','archive','v0.15.2','--',*paths])
labels=args.variants.split(',')if args.variants else(['baseline','safe-cache','coordinated-kernel']if args.smoke else[p.stem.removeprefix('overlay-')for p in sorted(source.glob('overlay-*.json'))])
input_path=source/('deep-inputs.json'if args.smoke else'inputs.json')
def run(label,step,command,extra=None):
 start=time.monotonic();p=subprocess.run([str(x)for x in command],env=dict(env,**(extra or {})),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(out/f'{label}-{step}.txt').write_text(p.stdout);checks.append({'variant':label,'step':step,'command':[str(x)for x in command],'exit':p.returncode,'elapsed_s':time.monotonic()-start});(out/'checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(label,step,p.returncode,flush=True);return p.returncode
for label in labels:
 overlay=json.loads((source/f'overlay-{label}.json').read_text());tree=out/label;tree.mkdir()
 with tarfile.open(fileobj=io.BytesIO(archive))as tar:tar.extractall(tree,filter='data')
 for name,text in overlay['files'].items():(tree/name).write_text(text)
 driver='audit_test.go.txt'if label=='instrumented-baseline'else'driver_test.go.txt';(tree/'research_replay_test.go').write_text((source/driver).read_text())
 if run(label,'build',['go','-C',tree,'build','-o',out/('pandafish-'+label),'.']):continue
 if label=='instrumented-baseline':run(label,'audit',['go','-C',tree,'test','-count=1','-run','^TestArchitectureAudit$','-timeout','12m','-v'],{'ARCH_INPUT':str(input_path),'ARCH_OUTPUT':str(out/(label+'-audit.json'))});continue
 if not args.smoke:run(label,'root-tests',['go','-C',tree,'test','-count=1','./...'])
 run(label,'ghi',['go','-C',tree,'test','-count=1','-run','^TestArchitectureGHI$','-v'],{'ARCH_INPUT':str(input_path),'ARCH_OUTPUT':str(out/(label+'-ghi.json'))})
 run(label,'fixed',['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','12m','-v'],{'ARCH_INPUT':str(input_path),'ARCH_OUTPUT':str(out/(label+'-fixed.json')),'ARCH_MODE':'fixed','ARCH_DEPTHS':'3,4'})
 if not args.smoke:run(label,'warm',['go','-C',tree,'test','-count=1','-run','^TestArchitectureProbe$','-timeout','12m','-v'],{'ARCH_INPUT':str(source/'inputs.json'),'ARCH_OUTPUT':str(out/(label+'-warm.json')),'ARCH_MODE':'warm'})
print('Replay outputs preserved:',out)
print('Some full-mode diagnostic tests deliberately fail changed contracts or unsafe policies; inspect checks, CACHE.md and METHOD.md. No game suite was run.')
