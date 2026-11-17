"""Independent feature geometry, including color/rank symmetry, via python-chess."""
import json
from pathlib import Path
import chess
r=Path(__file__).resolve().parent
fens=[c['fen']for c in json.loads((r.parent/'v0.17-search-moves/oracle-inputs.json').read_text())['cases'][:180]]
def features(b):
 values=[[0]*12 for _ in range(2)];colors=[chess.WHITE,chess.BLACK];zones=[];pawn_attacks=[];phase=0
 for c in colors:
  k=b.king(c);zones.append(set(b.attacks(k))|{k})
  pawn_attacks.append(set().union(*(set(b.attacks(p))for p in b.pieces(chess.PAWN,c))))
 for i,c in enumerate(colors):
  king=b.king(c);file=chess.square_file(king);rank=chess.square_rank(king);step=1 if c else -1
  values[i][6]=sum(1 for p in b.pieces(chess.PAWN,c)if abs(chess.square_file(p)-file)<=1 and (chess.square_rank(p)-rank)*step in [1,2])
  values[i][4]=len(pawn_attacks[i]&zones[1-i]);attackers=int(values[i][4]>0)
  for j,piece in enumerate([chess.KNIGHT,chess.BISHOP,chess.ROOK,chess.QUEEN]):
   for sq in b.pieces(piece,c):
    phase+=[1,1,2,4][j];attacks=set(b.attacks(sq))
    values[i][j]+=len(attacks-set(chess.SquareSet(b.occupied_co[c]))-pawn_attacks[1-i])
    hit=len(attacks&zones[1-i]);values[i][4]+=hit;attackers+=int(hit>0)
  values[i][5]=attackers*attackers
  files=[sum(chess.square_file(p)==f for p in b.pieces(chess.PAWN,c))for f in range(8)]
  values[i][10]=-sum(max(0,n-1)for n in files)
  for p in b.pieces(chess.PAWN,c):
   f=chess.square_file(p);rank=chess.square_rank(p);relative=rank if c else 7-rank;next=p+(8 if c else -8)
   values[i][9]-=int((f==0 or files[f-1]==0)and(f==7 or files[f+1]==0))
   values[i][11]-=int(0<=next<64 and b.piece_at(next)is not None)
   passed=not any(abs(chess.square_file(other)-f)<=1 and (chess.square_rank(other)-rank)*(1 if c else -1)>0 for other in b.pieces(chess.PAWN,not c))
   if passed:
    values[i][7]+=relative*relative
    if 0<=next<64:
     def distance(a,d):return max(abs(chess.square_file(a)-chess.square_file(d)),abs(chess.square_rank(a)-chess.square_rank(d)))
     values[i][8]+=max(0,relative-2)*(distance(b.king(not c),next)-distance(b.king(c),next))
 return {'Values':[a-b for a,b in zip(values[0],values[1])],'Phase':min(24,phase)}
cases=[]
for i,fen in enumerate(fens):
 b=chess.Board(fen)
 for label,board in [('original',b),('mirror',b.mirror())]:cases.append({'id':f'{i}-{label}','fen':board.fen(en_passant='fen'),'expected':features(board)})
(r/'feature-oracle-input.json').write_text(json.dumps({'cases':cases},indent=2)+'\n');print('Independent feature cases',len(cases))
