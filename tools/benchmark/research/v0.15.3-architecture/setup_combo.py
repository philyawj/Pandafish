"""Exploratory interaction probe after individual ablations, not a strength candidate."""
import io,json,subprocess,tarfile,tempfile
from pathlib import Path
r=Path(__file__).resolve().parent;repo=Path.cwd();names=subprocess.check_output(['git','ls-tree','--name-only','v0.15.2'],text=True).splitlines();base={n:subprocess.check_output(['git','show','v0.15.2:'+n],text=True)for n in names if n.endswith('.go')};files=dict(base);files.update(json.loads((r/'overlay-pvs-safe-cache.json').read_text())['files'])
s=files['iterative.go'].replace('1-bestScore','-bestScore').replace('if score == bestScore {','if score == bestScore && len(best)==0 {');files['iterative.go']=s
helper=json.loads((r/'overlay-tactical-ordering.json').read_text())['files']['search.go'];helper='func orderedTacticalMoves'+helper.split('func orderedTacticalMoves',1)[1]
s=files['search.go'];needle='for _, move := range orderedMoves(position) {';assert needle in s;s=s.replace(needle,'var moves []chess.Move\n if checked {moves=orderedMoves(position)} else {moves=orderedTacticalMoves(position)}\n for _, move := range moves {',1);s+='\n'+helper;files['search.go']=s
for name in ['main.go','main_test.go']:files[name]=files[name].replace('v0.15.3-diag-pvs-safe-cache','v0.15.3-diag-coordinated-kernel')
paths=[n for n in names if n.endswith('.go')or n in['go.mod','go.sum','internal','testdata']];tree=Path(tempfile.mkdtemp(prefix='pandafish-coordinated-kernel-'))
with tarfile.open(fileobj=io.BytesIO(subprocess.check_output(['git','archive','v0.15.2','--',*paths])))as tar:tar.extractall(tree,filter='data')
for name,text in files.items():(tree/name).write_text(text)
subprocess.run(['gofmt','-w',*[str(tree/n)for n in files]],check=True)
changed={name:(tree/name).read_text()for name in files if base.get(name)!=(tree/name).read_text()}
(r/'overlay-coordinated-kernel.json').write_text(json.dumps({'base_tag':'v0.15.2','variant':'coordinated-kernel','selection_note':'Exploratory combination selected after measuring individual components; no independent strength validation. Exact fixed-depth score retained, root tie set intentionally becomes one best move.','files':changed},indent=2)+'\n')
(tree/'research_driver_test.go').write_text((r/'driver_test.go.txt').read_text());manifest=json.loads((r/'worktrees.json').read_text());manifest.append({'label':'coordinated-kernel','tree':str(tree)});(r/'worktrees.json').write_text(json.dumps(manifest,indent=2)+'\n');print(tree)
