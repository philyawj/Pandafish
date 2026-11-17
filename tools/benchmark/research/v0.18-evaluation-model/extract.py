"""Extract the frozen prototype features in a fresh v0.17-based scratch package."""
import io,json,os,shutil,subprocess,tarfile,tempfile,time
from pathlib import Path
r=Path(__file__).resolve().parent;repo=r.parents[3];tree=Path(tempfile.mkdtemp(prefix='pandafish-v018-features-'))
with tarfile.open(fileobj=io.BytesIO(subprocess.check_output(['git','archive','v0.17'],cwd=repo)))as t:t.extractall(tree,filter='data')
for p in repo.glob('*.go'):(tree/p.name).write_bytes(p.read_bytes())
shutil.copyfile(repo/'third_party/chess/eval_features.go',tree/'third_party/chess/eval_features.go')
base=subprocess.check_output(['git','show','v0.17:evaluation.go'],cwd=repo,text=True).replace('func staticScore(', 'func baseStaticScore(')
base+='\nfunc staticScore(position *chess.Position, perspective chess.Color) int { return baseStaticScore(position, perspective) }\n'
(tree/'evaluation.go').write_text(base)

(tree/'research_extract_test.go').write_bytes((r/'extract_test.go.txt').read_bytes())
checks=[]
for label,source in [('oracle','feature-oracle-input.json'),('corpus','corpus.json')]:
 cmd=['go','-C',str(tree),'test','-count=1','-run','^TestExportEvaluationFeatures$','-v'];start=time.monotonic()
 p=subprocess.run(cmd,env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',MODEL_INPUT=str(r/source),MODEL_OUTPUT=str(r/(label+'-features.json'))),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 (r/(label+'-extract.txt')).write_text(p.stdout);checks.append({'label':label,'command':cmd,'exit':p.returncode,'elapsed_s':time.monotonic()-start});(r/'extract-checks.json').write_text(json.dumps(checks,indent=2)+'\n');assert p.returncode==0,p.stdout
expected={x['id']:x for x in json.loads((r/'feature-oracle-input.json').read_text())['cases']}
errors=[]
for row in json.loads((r/'oracle-features.json').read_text()):
 if row['Features']!=expected[row['ID']]['expected']:errors.append({'id':row['ID'],'actual':row['Features'],'expected':expected[row['ID']]['expected']})
(r/'feature-oracle-results.json').write_text(json.dumps({'positions':len(expected),'passed':len(expected)-len(errors),'failures':errors},indent=2)+'\n');assert not errors,errors[:3]
(r/'worktree.json').write_text(json.dumps({'tree':str(tree)},indent=2)+'\n')
print('Independent features pass:',len(expected),'corpus rows:',len(json.loads((r/'corpus-features.json').read_text())))
