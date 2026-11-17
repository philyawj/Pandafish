"""Finite-node scoring of actual selected moves; no candidate move is forced."""
import hashlib,json
from pathlib import Path
import chess,chess.engine
r=Path(__file__).resolve().parent
if(r/'decision-reference.json').exists():raise SystemExit('Refuse overwrite')
rows=[x for x in (json.loads((r/'uci-performance.json').read_text())['records']+json.loads((r/'ablation-uci.json').read_text())['records'])if x['load']=='serial'and x['legal']]
positions={}
for x in rows:positions.setdefault(x['position'],set()).add(x['move'])
def board(command):
 words=command.split();i=words.index('moves')if'moves'in words else len(words)
 b=chess.Board()if words[1]=='startpos'else chess.Board(' '.join(words[2:i]))
 for move in words[i+1:]:b.push_uci(move)
 return b
prior=json.loads((r.parent/'decision-reference.json').read_text());old={(x['position'],x['move']):x for x in prior['queries']}
report={'reused_query_count':0,'method':json.loads((r.parent/'DECISION_PLAN.json').read_text()),'stockfish_sha256':hashlib.sha256(Path('/opt/homebrew/bin/stockfish').read_bytes()).hexdigest(),'queries':[]}
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30)as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False});report['identity']=sf.id
 for index,(position,moves)in enumerate(positions.items()):
  b=board(position)
  for move in [None,*sorted(moves)]:
   if (position,move) in old:
    report['queries'].append(old[position,move]);report['reused_query_count']+=1;continue
   sf.configure({'Clear Hash':None});iterations=[]
   try:
    with sf.analysis(b,chess.engine.Limit(nodes=100000),root_moves=[chess.Move.from_uci(move)]if move else None,game=object())as stream:
     for info in stream:
      if'score'not in info:continue
      score=info['score'].pov(b.turn)
      iterations.append({'depth':info.get('depth'),'nodes':info.get('nodes'),'cp':score.score(),'mate':score.mate(),'lowerbound':info.get('lowerbound',False),'upperbound':info.get('upperbound',False),'pv':[m.uci()for m in info.get('pv',[])]})
    exact=[x for x in iterations if not x['lowerbound']and not x['upperbound']]
    report['queries'].append({'position':position,'move':move,'iterations':iterations,'last_exact':exact[-1]if exact else None});assert exact
   except Exception as error:
    report['queries'].append({'position':position,'move':move,'iterations':iterations,'error':repr(error)});(r/'decision-reference.json').write_text(json.dumps(report,indent=2)+'\n');raise
  (r/'decision-reference.json').write_text(json.dumps(report,indent=2)+'\n');print('Reference positions',index+1,'/',len(positions),flush=True)
queries={(x['position'],x['move']):x['last_exact']for x in report['queries']};out={};all_passed=True
for budget in [100,1000]:
 cases={};mates=[]
 for position in positions:
  relevant=[x for x in rows if x['position']==position and x['budget_ms']==budget]
  root=queries[position,None];finite=root['cp']is not None and all(queries[position,x['move']]['cp']is not None for x in relevant)
  if not finite:
   mates.append({'position':position,'unrestricted':root,'selected':[{'actor':x['actor'],'trial':x['trial'],'move':x['move'],'reference':queries[position,x['move']]}for x in relevant]});continue
  values={}
  for actor in ['baseline','candidate','fast-base']:
   actual=[x for x in relevant if x['actor']==actor];assert len(actual)==(3 if budget==100 else 1)
   losses=[min(1000,max(0,root['cp']-queries[position,x['move']]['cp']))for x in actual]
   values[actor]={'mean_loss_cp':sum(losses)/len(losses),'trial_losses':losses,'moves':[x['move']for x in actual]}
  cases[position]=values
 mean={a:sum(c[a]['mean_loss_cp']for c in cases.values())/len(cases)for a in ['baseline','candidate','fast-base']}
 severe={a:sum(c[a]['mean_loss_cp']>300 for c in cases.values())for a in ['baseline','candidate','fast-base']}
 passed=mean['candidate']<=mean['baseline']+max(5,.1*mean['baseline'])and severe['candidate']<=severe['baseline']+2
 out[str(budget)]={'finite_cp_positions':len(cases),'mean_loss_cp':mean,'severe_positions_over_300cp':severe,'gate_passed':passed,'cases':cases,'mate_cases':mates};all_passed=all_passed and passed
(r/'decision-quality.json').write_text(json.dumps({'budgets':out,'gate_passed':all_passed,'notes':'Finite-node diagnostic estimates, not unique acceptable moves or strength evidence; all mate-valued cases preserved separately. Weights remain frozen.'},indent=2)+'\n')
print(json.dumps({'gate_passed':all_passed,'budgets':{k:{f:v for f,v in x.items()if f not in ['cases','mate_cases']}for k,x in out.items()}},indent=2))
