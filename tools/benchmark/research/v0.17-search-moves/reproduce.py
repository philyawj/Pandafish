"""Restore the frozen v0.17 search-backend candidate in a NEW directory, no games."""
import argparse,io,json,os,subprocess,tarfile
from pathlib import Path
parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('output',type=Path);args=parser.parse_args()
r=Path(__file__).resolve().parent;out=args.output.resolve();out.mkdir(parents=True,exist_ok=False)
archive=subprocess.check_output(['git','archive','v0.16'])
with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(out,filter='data')
for name,text in json.loads((r/'candidate-source.json').read_text())['files'].items():
 p=out/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(text)
checks=[]
for name,command in [('tests',['go','-C',str(out),'test','-count=1','./...']),('build',['go','-C',str(out),'build','-o',str(out/'pandafish-candidate'),'.'])]:
 p=subprocess.run(command,env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache'),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(out/(name+'.txt')).write_text(p.stdout);checks.append({'name':name,'command':command,'exit':p.returncode});(out/'replay-checks.json').write_text(json.dumps(checks,indent=2)+'\n');print(name,p.returncode,flush=True)
 if p.returncode:raise SystemExit(p.returncode)
print('Restored source and all replay outputs:',out)
