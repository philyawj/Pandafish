import json
from pathlib import Path
import chess,chess.engine
r=Path(__file__).resolve().parent
b=chess.Board('r4bk1/pp5p/2n2n2/8/3N4/5BB1/PP5P/R5K1 w - - 0 1')
out=[]
for name,binary in [('baseline','tools/benchmark/bin/pandafish-v0.15.2'),('referee','/opt/homebrew/bin/stockfish')]:
 with chess.engine.SimpleEngine.popen_uci(binary,timeout=30) as engine:
  if name=='referee':engine.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False})
  for depth in [3,5]:
   info=engine.play(b,chess.engine.Limit(depth=depth),info=chess.engine.INFO_ALL,game=object());out.append({'engine':name,'fen':b.fen(),'requested_depth':depth,'move':info.move.uci(),'depth':info.info.get('depth'),'score':str(info.info.get('score')),'pv':[m.uci() for m in info.info.get('pv',[])]})
  if name=='referee':
   for move in ['f3c6','d4c6']:
    info=engine.analyse(b,chess.engine.Limit(nodes=500000),root_moves=[chess.Move.from_uci(move)],game=object());out.append({'engine':name,'fen':b.fen(),'restriction':move,'score':str(info['score']),'pv':[m.uci() for m in info.get('pv',[])]})
(r/'bishop-diagnosis.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
