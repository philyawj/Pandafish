"""Counterfactual pawn redundancy diagnosis; no evaluation changes yet."""
import json
from pathlib import Path
import chess,chess.engine
r=Path(__file__).resolve().parent
d=json.loads((r/'diagnosis.json').read_text());original=next(x for x in d['records'] if x['game']=='DsqFd5Kj' and len(x['line'])==3)
variants=[]
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30) as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False})
 for label,changes in [('original',[]),('split g7 to f7',[('g7','f7')]),('split g7 to h7',[('g7','h7')]),('split g6 to f6',[('g6','f6')]),('split g6 to h6',[('g6','h6')]),('remove central passer for isolated extra a pawn',[('d5','a6')])]:
  b=chess.Board(original['fen'])
  for a,z in changes:b.set_piece_at(chess.parse_square(z),b.remove_piece_at(chess.parse_square(a)))
  assert b.is_valid()
  exact=[]
  with sf.analysis(b,chess.engine.Limit(nodes=500000),game=object()) as analysis:
   for info in analysis:
    if 'score' in info and not info.get('lowerbound') and not info.get('upperbound'):
     s=info['score'].pov(chess.WHITE);exact.append({'white_score':{'mate':s.mate()} if s.is_mate() else {'cp':s.score()},'depth':info.get('depth'),'nodes':info.get('nodes'),'pv':[m.uci() for m in info.get('pv',[])]})
  variants.append({'label':label,'fen':b.fen(),'last_exact':exact[-1],'exact_iterations':exact});(r/'pawn-variants.json').write_text(json.dumps({'method':'Six predetermined same-material synthetic counterfactuals before choosing feature/weights. Fresh hash, 500000 nodes, Skill20, one thread, Hash64. Last completed exact iteration retained with all exact iterations. No historical repetition for edited boards.','variants':variants},indent=2)+'\n')
print([(x['label'],x['last_exact']['white_score']) for x in variants])
