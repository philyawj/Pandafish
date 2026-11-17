"""Create isolated, replayable diagnostic variants from v0.15.2. Never edits root Go."""
import io,json,subprocess,tarfile,tempfile
from pathlib import Path
repo=Path.cwd();r=repo/'tools/benchmark/research/v0.15.3-architecture'
names=subprocess.check_output(['git','ls-tree','--name-only','v0.15.2'],text=True).splitlines();paths=[x for x in names if x.endswith('.go') or x in ['go.mod','go.sum','internal','testdata']]
archive=subprocess.check_output(['git','archive','v0.15.2','--',*paths]);base={n:subprocess.check_output(['git','show','v0.15.2:'+n],text=True) for n in names if n.endswith('.go')}
configs={}
configs['baseline']=dict(base)
c=dict(base);c['transposition.go']=c['transposition.go'].replace('if table == nil {','if true {');configs['no-tt']=c
hint='''
type boardHint struct { key uint64; move uint16 }
func (t *transpositionTable) boardHint(p *chess.Position) uint16 {
 if t==nil{return 0};h:=p.ZobristHash();e:=t.hints[h%uint64(len(t.hints))];if e.key==h{return e.move};return 0
}
func (t *transpositionTable) storeHint(p *chess.Position, move uint16) {
 if t==nil{return};h:=p.ZobristHash();t.hints[h%uint64(len(t.hints))]=boardHint{h,move}
}
'''
def hints(src):
 c=dict(src);s=c['transposition.go'].replace('entries    []tableEntry','hints []boardHint\n\tentries    []tableEntry').replace('entries: make([]tableEntry, tableEntries)','entries: make([]tableEntry, tableEntries), hints: make([]boardHint,tableEntries)');s=s.replace('clear(table.entries)','clear(table.entries)\n clear(table.hints)');s+=hint;c['transposition.go']=s
 s=c['search.go'].replace('preferMove(moves, entry.move)\n\t}', 'preferMove(moves, entry.move)\n\t} else { preferMove(moves, worker.table.boardHint(position))\n\t}')
 needle='worker.table.store(position, depth, scoreToTable(best, ply), scoreToTable(originalAlpha, ply), scoreToTable(beta, ply), bestMove, worker.historyKey)';assert needle in s;s=s.replace(needle,needle+'\n worker.table.storeHint(position,bestMove)');c['search.go']=s;return c
configs['board-hints']=hints(base)
trimhelper='''package main
import "github.com/corentings/chess/v2"
// Only a child reached by an irreversible pawn/capture move resets history.
// Direct entry points (including FEN and repeated-position tests) do not reset.
// Keep counts for draw detection; reset only the cache fingerprint. Earlier
// positions cannot recur after a pawn move/capture, so their counts cannot
// affect descendant repetition results. Restore the prior key on return.
func (w *searcher) childSearch(p *chess.Position,d,a,b int,c bool,ply int)(int,bool){
 if p.HalfMoveClock()!=0{return w.searchAtPly(p,d,a,b,c,ply)}
 key:=w.historyKey;w.historyKey=0
 defer func(){w.historyKey=key}()
 return w.searchAtPly(p,d,a,b,c,ply)
}
func (w *searcher) childQuiet(p *chess.Position,d,a,b int,c bool,ply int)(int,bool){
 if p.HalfMoveClock()!=0{return w.quietAtPly(p,d,a,b,c,ply)}
 key:=w.historyKey;w.historyKey=0
 defer func(){w.historyKey=key}()
 return w.quietAtPly(p,d,a,b,c,ply)
}
'''
def trim(src):
 c=dict(src);c['research_context.go']=trimhelper
 for f in ['search.go','iterative.go']:c[f]=c[f].replace('worker.searchAtPly(position.Update(&move),','worker.childSearch(position.Update(&move),').replace('worker.quietAtPly(position.Update(&move),','worker.childQuiet(position.Update(&move),')
 return c
configs['trim-history']=trim(base);configs['safe-cache']=trim(hints(base))
c=dict(base);c['transposition.go']=c['transposition.go'].replace('if len(history) > 0 {\n\t\tkey ^= history[0]\n\t}', 'if false { key ^= 0 }');configs['unsafe-board-scores']=c
c=dict(base);c['search.go']=c['search.go'].replace('int(entry.depth) == depth','int(entry.depth) >= depth');configs['depth-atleast']=c
c=dict(base);c['transposition.go']=c['transposition.go'].replace('tableEntries = 1 << 16','tableEntries = 1 << 20');configs['large-tt']=c
c=dict(base);c['iterative.go']=c['iterative.go'].replace('1-bestScore','-bestScore').replace('if score == bestScore {','if score == bestScore && len(best)==0 {');configs['root-single-best']=c
old='''score, complete := worker.searchAtPly(position.Update(&move), depth-1, -beta, -alpha, move.HasTag(chess.Check), ply+1)'''
new='''var score int
 var complete bool
 if moveIndex==0 { score,complete = worker.searchAtPly(position.Update(&move),depth-1,-beta,-alpha,move.HasTag(chess.Check),ply+1) } else {
 score,complete = worker.searchAtPly(position.Update(&move),depth-1,-alpha-1,-alpha,move.HasTag(chess.Check),ply+1)
 if complete && -score>alpha && -score<beta { score,complete = worker.searchAtPly(position.Update(&move),depth-1,-beta,-alpha,move.HasTag(chess.Check),ply+1) }
 }'''
def pvs(src):
 c=dict(src);s=c['search.go'];target=old.replace('searchAtPly','childSearch') if 'research_context.go'in c else old
 replacement=new.replace('searchAtPly','childSearch') if 'research_context.go'in c else new
 assert target in s;s=s.replace(target,replacement,1).replace('for _, move := range moves {','for moveIndex, move := range moves {',1);c['search.go']=s;return c
configs['pvs']=pvs(base);configs['pvs-safe-cache']=pvs(configs['safe-cache'])
c=dict(base);s=c['search.go'];start=s.index('func orderedMoves(');tactical=s[start:].replace('func orderedMoves(', 'func orderedTacticalMoves(',1).replace('moves := position.ValidMoves()','''var moves []chess.Move
 for _,m:=range position.ValidMovesUnsafe(){if m.HasTag(chess.Capture)||m.HasTag(chess.EnPassant)||m.Promo()!=chess.NoPieceType{moves=append(moves,m)}}''',1)
needle='for _, move := range orderedMoves(position) {';assert needle in s;s=s.replace(needle,'moves := orderedMoves(position)\n if !checked { moves = orderedTacticalMoves(position) }\n for _, move := range moves {',1)
# Avoid generating/sorting the full list before the tactical-only branch.
s=s.replace('moves := orderedMoves(position)\n if !checked { moves = orderedTacticalMoves(position) }', 'var moves []chess.Move\n if checked { moves=orderedMoves(position) } else { moves=orderedTacticalMoves(position) }',1)
s+='\n'+tactical;c['search.go']=s;configs['tactical-ordering']=c

# Diagnostic counters are added only to one tree, never to timed executables.
c=dict(base);s=c['search.go'];s=s.replace('type searcher struct {','type searcher struct {\n audit *researchAudit')
s=s.replace('worker.nodes++','worker.nodes++\n if worker.audit!=nil { worker.audit.MainNodes++ }',1)
needle='entry, hit := worker.table.probe(position, worker.historyKey)';s=s.replace(needle,needle+'\n if worker.audit!=nil {worker.audit.probe(position,depth,alpha,beta,ply,worker.historyKey,entry,hit)}')
s=s.replace('return score, !worker.stopped()','if worker.audit!=nil {worker.audit.Cutoffs++}; return score, !worker.stopped()',1)
s=s.replace('for _, move := range moves {','for moveIndex, move := range moves {',1)
s=s.replace('worker.killers.remember(ply, move)','if worker.audit!=nil {worker.audit.BetaCutoffs++;worker.audit.CutoffOrdinal+=uint64(moveIndex+1);if moveIndex==0{worker.audit.FirstMoveCutoffs++}}\n worker.killers.remember(ply, move)',1)
needle='worker.table.store(position, depth, scoreToTable(best, ply), scoreToTable(originalAlpha, ply), scoreToTable(beta, ply), bestMove, worker.historyKey)';s=s.replace(needle,needle+'\n if worker.audit!=nil {worker.audit.store(position,depth,scoreToTable(best,ply),scoreToTable(originalAlpha,ply),scoreToTable(beta,ply),worker.historyKey)}')
# Quiescence begins at the second worker.nodes increment.
q=s.index('func (worker *searcher) quietAtPly');s=s[:q]+s[q:].replace('worker.nodes++','worker.nodes++\n if worker.audit!=nil {worker.audit.QNodes++}',1).replace('static := staticScore(position, position.Turn())','if worker.audit!=nil {worker.audit.Evals++}; static := staticScore(position, position.Turn())',1).replace('for _, move := range orderedMoves(position) {','moves:=orderedMoves(position)\n if worker.audit!=nil {worker.audit.QGenerated+=uint64(len(moves))};for _, move := range moves {',1).replace('// Equal-or-cheaper attackers','if worker.audit!=nil {worker.audit.QEligible++}\n // Equal-or-cheaper attackers',1)
c['search.go']=s;c['research_audit.go']=(r/'audit.go.txt').read_text();configs['instrumented-baseline']=c
manifest=[];work=Path(tempfile.mkdtemp(prefix='pandafish-architecture-'))
for label,files in configs.items():
 tree=work/label;tree.mkdir()
 with tarfile.open(fileobj=io.BytesIO(archive)) as tar:tar.extractall(tree,filter='data')
 changed={}
 for name,text in files.items():
  if name in ['main.go','main_test.go'] and label!='baseline':text=text.replace('Pandafish v0.15.2','Pandafish v0.15.3-diag-'+label)
  (tree/name).write_text(text)
  if base.get(name)!=text:changed[name]=text
 subprocess.run(['gofmt','-w',*[str(tree/n) for n in files]],check=True)
 for name in changed:changed[name]=(tree/name).read_text()
 if label!='instrumented-baseline':(tree/'research_driver_test.go').write_text((r/'driver_test.go.txt').read_text())
 else:(tree/'research_audit_test.go').write_text((r/'audit_test.go.txt').read_text())
 (r/('overlay-'+label+'.json')).write_text(json.dumps({'base_tag':'v0.15.2','variant':label,'files':changed},indent=2)+'\n');manifest.append({'label':label,'tree':str(tree)})
(r/'worktrees.json').write_text(json.dumps(manifest,indent=2)+'\n');print(work)
