"""Illustrative leaf-score audit; not a representative error-rate estimate."""
import json
from pathlib import Path
import chess,chess.engine
r=Path(__file__).resolve().parent;inputs={x['game']:x for x in json.loads((r/'leaf-inputs.json').read_text())};out=[]
def board(command,pv):
 s=command.split();i=s.index('moves')if'moves'in s else len(s);b=chess.Board()if s[1]=='startpos'else chess.Board(' '.join(s[2:i]))
 for m in s[i+1:]+pv:b.push_uci(m)
 return b
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30)as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False})
 for trace in json.loads((r/'leaf-traces.json').read_text()):
  if trace['Kind']!='baseline_exact_research_pv' or not trace['Nodes']:continue
  c=inputs[trace['Game']];b=board(c['position'],trace['PV']);last=trace['Nodes'][-1];assert b.fen(en_passant='fen')==last['FEN']
  record={'game':trace['Game'],'forced_diagnostic_root':trace['Root'],'depth':trace['Depth'],'pv':trace['PV'],'leaf_fen':b.fen(en_passant='fen'),'baseline_black_cp':last['Static'],'in_check':b.is_check(),'terminal_no_legal_moves':not bool(list(b.legal_moves)),'iterations':[]};out.append(record)
  if record['terminal_no_legal_moves']:(r/'leaf-referee.json').write_text(json.dumps(out,indent=2)+'\n');continue
  sf.protocol.send_line('ucinewgame');sf.configure({'Clear Hash':None})
  with sf.analysis(b,chess.engine.Limit(nodes=100000),game=object())as stream:
   for info in stream:
    if 'score'not in info:continue
    v=info['score'].pov(chess.BLACK);record['iterations'].append({'black_score':{'mate':v.mate()}if v.is_mate()else{'cp':v.score()},'depth':info.get('depth'),'nodes':info.get('nodes'),'bound':bool(info.get('upperbound')or info.get('lowerbound')),'pv':[m.uci()for m in info.get('pv',[])]})
  exact=[x for x in record['iterations']if not x['bound']];record['last_exact']=exact[-1]if exact else None;(r/'leaf-referee.json').write_text(json.dumps(out,indent=2)+'\n')
print('Saved',len(out),'selected leaf comparisons')
