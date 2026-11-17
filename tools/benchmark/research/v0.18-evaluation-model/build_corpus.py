"""Frozen offline corpus: played positions plus stronger reference replies."""
import hashlib,json
from pathlib import Path
import chess,chess.engine,chess.pgn
r=Path(__file__).resolve().parent;repo=r.parents[3];plan=json.loads((r/'PLAN.json').read_text());source=repo/plan['source_directory']
roots=[];sources={}
for opponent in ['previous','stockfish','sacrifice']:
 p=source/(opponent+'.pgn');sources[p.name]=hashlib.sha256(p.read_bytes()).hexdigest()
 with p.open()as f:
  for game_index in range(20):
   game=chess.pgn.read_game(f);assert game and not game.errors
   moves=list(game.mainline_moves());family=' '.join(m.uci()for m in moves[:4]);seen=set()
   for target in [12,24,36,48,64,80]:
    ply=min(target,len(moves)-2);b=game.board()
    for m in moves[:ply]:b.push(m)
    fen=b.fen(en_passant='fen')
    if fen in seen or not any(b.legal_moves):continue
    seen.add(fen);roots.append({'id':f'{opponent}-{game_index+1:02}-p{ply}','game':f'{opponent}-{game_index+1:02}','family':family,'target_ply':target,'chosen_ply':ply,'fen':fen,'position':'position startpos moves '+' '.join(m.uci()for m in moves[:ply]),'checked':b.is_check(),'claimable':b.can_claim_draw(),'halfmove':b.halfmove_clock})
families=sorted({x['family']for x in roots});assert len(families)==10
splits={family:'train'if i<6 else'validation'if i<8 else'holdout'for i,family in enumerate(families)}
queries={};queryfile=r/'teacher-queries.json'
if queryfile.exists():
 old=json.loads(queryfile.read_text());assert old['plan_sha256']==hashlib.sha256((r/'PLAN.json').read_bytes()).hexdigest();queries=old['queries']
report={'plan_sha256':hashlib.sha256((r/'PLAN.json').read_bytes()).hexdigest(),'source_sha256':sources,'stockfish_sha256':hashlib.sha256(Path('/opt/homebrew/bin/stockfish').read_bytes()).hexdigest(),'split_families':splits,'queries':queries}
def save():queryfile.write_text(json.dumps(report,indent=2)+'\n')
with chess.engine.SimpleEngine.popen_uci('/opt/homebrew/bin/stockfish',timeout=30)as sf:
 sf.configure({'Threads':1,'Hash':64,'Skill Level':20,'UCI_LimitStrength':False});report['identity']=sf.id
 def query(fen):
  if fen in queries:return queries[fen]
  b=chess.Board(fen);rows=[];sf.configure({'Clear Hash':None})
  try:
   with sf.analysis(b,chess.engine.Limit(nodes=50000),game=object())as stream:
    for info in stream:
     if 'score'not in info:continue
     score=info['score'].pov(chess.WHITE)
     rows.append({'depth':info.get('depth'),'nodes':info.get('nodes'),'cp':score.score(),'mate':score.mate(),'lowerbound':info.get('lowerbound',False),'upperbound':info.get('upperbound',False),'pv':[m.uci()for m in info.get('pv',[])]})
   exact=[x for x in rows if not x['lowerbound']and not x['upperbound']]
   result={'iterations':rows,'last_exact':exact[-1]if exact else None};queries[fen]=result;save();assert result['last_exact']
   return result
  except Exception as error:
   queries[fen]={'iterations':rows,'error':repr(error)};save();raise
 cases=[]
 for index,c in enumerate(roots):
  q=query(c['fen']);root=dict(c,kind='played',split=splits[c['family']],label=q['last_exact']);cases.append(root)
  pv=q['last_exact']['pv']
  if pv:
   b=chess.Board(c['fen']);move=chess.Move.from_uci(pv[0]);assert move in b.legal_moves;b.push(move)
   if any(b.legal_moves):
    fen=b.fen(en_passant='fen');q2=query(fen)
    cases.append(dict(c,id=c['id']+'/reference-reply',kind='reference-reply',split=splits[c['family']],fen=fen,position='position fen '+fen,checked=b.is_check(),claimable=c['claimable']or b.can_claim_draw(),halfmove=b.halfmove_clock,label=q2['last_exact']))
  if index%25==0:print('Roots',index+1,'/',len(roots),'queries',len(queries),flush=True)
  (r/'corpus.json').write_text(json.dumps({'sources':sources,'split_families':splits,'method':plan,'cases':cases},indent=2)+'\n')
print('Corpus complete:',len(cases),'rows,',len(queries),'distinct reference FENs',flush=True)
