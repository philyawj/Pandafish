"""Score every distinct selected UCI move, preserving completed exact iterations."""
import hashlib,json
from pathlib import Path
import chess,chess.engine
r=Path(__file__).resolve().parent
inputs=[r/'uci-100ms.json',r/'uci-1000ms.json',r/'uci-warm.json']
for name in ['uci-historical.json','uci-historical-fixed.json','uci-combo.json']:
 if (r/name).exists():inputs.append(r/name)
queries={};ids={}
for p in inputs:
 for row in json.loads(p.read_text())['records']:
  if not row.get('legal'):continue
  position=row['position'];queries[position,None]=True;queries[position,row['move']]=True;ids.setdefault(position,set()).add(row['id'])
def board(command):
 s=command.split();i=s.index('moves')if'moves'in s else len(s);b=chess.Board()if s[1]=='startpos'else chess.Board(' '.join(s[2:i]))
 for move in s[i+1:]:b.push_uci(move)
 return b
prior={}
if (r/'referee-results.json').exists():
 for old in json.loads((r/'referee-results.json').read_text())['results']:prior[old['position'],old['move']]=old
out={'reused_query_count':0,'method':'Every distinct legal selected move plus an unrestricted comparator for each full-history position. Fresh hash, one thread, Hash64, Skill20, strength limiting off, 100000 nodes each. Keep every score-bearing iteration; use last completed exact iteration, not later bound updates. Finite-node comparator, not ground truth or an Elo test.','nodes_per_query':100000,'stockfish_sha256':hashlib.sha256(Path('/opt/homebrew/bin/stockfish').read_bytes()).hexdigest(),'results':[]}
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30) as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False});out['identity']=sf.id
 for position in dict.fromkeys(p for p,m in queries):
  b=board(position)
  for move in [m for p,m in queries if p==position]:
   if (position,move) in prior:
    out['results'].append(prior[position,move]);out['reused_query_count']+=1;continue
   sf.protocol.send_line('ucinewgame');sf.configure({'Clear Hash':None});iterations=[]
   with sf.analysis(b,chess.engine.Limit(nodes=100000),root_moves=[chess.Move.from_uci(move)]if move else None,game=object())as stream:
    for info in stream:
     if 'score'not in info:continue
     s=info['score'].pov(b.turn);iterations.append({'score':{'mate':s.mate()}if s.is_mate()else{'cp':s.score()},'depth':info.get('depth'),'nodes':info.get('nodes'),'upperbound':info.get('upperbound',False),'lowerbound':info.get('lowerbound',False),'pv':[m.uci()for m in info.get('pv',[])]})
   exact=[i for i in iterations if not i['upperbound']and not i['lowerbound']];out['results'].append({'position':position,'ids':sorted(ids[position]),'move':move,'iterations':iterations,'last_exact':exact[-1]if exact else None})
   (r/'referee-results.json').write_text(json.dumps(out,indent=2)+'\n')
  print('position',len(out['results']),'/ total queries',len(queries),flush=True)

(r/'referee-results.json').write_text(json.dumps(out,indent=2)+'\n')
