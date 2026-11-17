"""Offline fresh diagnosis before selecting the next implementation."""
import json
from pathlib import Path
import chess,chess.engine
root=Path.cwd();out=root/'tools/benchmark/research/v0.16-next-assessment'
inputs=json.loads((root/'tools/benchmark/live/v0.15-live-20261113-074740/analysis/engine-replay/inputs.json').read_text())
records=[]
def board_for(c):
 s=c['position'].split();i=s.index('moves');b=chess.Board(' '.join(s[2:i]));
 for m in s[i+1:]:b.push_uci(m)
 return b
def snap(engine,b,nodes):
 engine.protocol.send_line('ucinewgame');engine.configure({'Clear Hash':None});d=engine.analyse(b,chess.engine.Limit(nodes=nodes),game=object());score=d['score'].pov(chess.WHITE)
 return {'white_score':{'mate':score.mate()} if score.is_mate() else {'cp':score.score()},'depth':d.get('depth'),'nodes':d.get('nodes'),'pv':[m.uci() for m in d.get('pv',[])]}
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30) as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False})
 for game,lines in [('DsqFd5Kj',[[],['d6d5'],['d6d5','c5d5','e6d5'],['d6d1']]),('PNfQoYP0',[[],['h2h4'],['e3f3']])]:
  c=next(x for x in inputs if x['game'].endswith(game))
  for line in lines:
   b=board_for(c)
   for m in line:b.push_uci(m)
   r={'game':game,'line':line,'fen':b.fen(),'position':'position startpos moves '+' '.join(m.uci() for m in b.move_stack),'stockfish':snap(sf,b,500000)}
   if game=='DsqFd5Kj' and len(line)==3:
    p=b.copy(stack=False);p.turn=not p.turn;r['opposite_turn']={'fen':p.fen(),'stockfish':snap(sf,p,500000)}
   records.append(r);(out/'diagnosis.json').write_text(json.dumps({'method':'Fresh 500000-node full-strength Stockfish comparison, white perspective, fresh hash, full history. Alternate turn is synthetic and has no historical repetition. Finite-node comparator, not proof or tuning target.','records':records},indent=2)+'\n')
print('Saved',len(records),'diagnostic boards')
