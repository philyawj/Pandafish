"""Independent python-chess legal moves, status, child FEN/hash and perft oracles."""
import hashlib,json,random
from pathlib import Path
import chess,chess.polyglot
r=Path(__file__).resolve().parent
if(r/'oracle-inputs.json').exists():raise SystemExit('Refuse overwrite')
base=json.loads((r/'inputs.json').read_text())['cases'];fens=[c['fen']for c in base]
randomizer=random.Random(170016);b=chess.Board()
for i in range(512):
 if b.is_game_over(claim_draw=True)or i%80==79:b=chess.Board()
 fens.append(b.fen(en_passant='fen'))
 moves=sorted(b.legal_moves,key=lambda m:m.uci());b.push(randomizer.choice(moves))
def perft(b,depth):
 if depth==0:return 1
 total=0
 for m in list(b.legal_moves):b.push(m);total+=perft(b,depth-1);b.pop()
 return total
cases=[]
for i,fen in enumerate(fens):
 b=chess.Board(fen);assert b.is_valid(),fen
 moves=[]
 for m in b.legal_moves:
  tactical=b.is_capture(m)or m.promotion is not None
  capture=b.is_capture(m)and not b.is_en_passant(m)
  ep=b.is_en_passant(m);castle=b.is_castling(m);check=b.gives_check(m)
  b.push(m)
  moves.append({'uci':m.uci(),'tactical':tactical,'capture':capture,'ep':ep,'castle':castle,'check':check,'child_fen':b.fen(en_passant='fen'),'child_hash':chess.polyglot.zobrist_hash(b)})
  b.pop()
 cases.append({'id':i,'fen':fen,'hash':chess.polyglot.zobrist_hash(b),'status':'mate'if b.is_checkmate()else'stalemate'if b.is_stalemate()else'normal','moves':moves,'perft2':perft(b,2)})
(r/'oracle-inputs.json').write_text(json.dumps({'method':'python-chess independent legality/tags, full FEN EP convention, Polyglot hash, mate/stalemate and depth2 perft. Original 84 FENs plus 512 seeded legal random positions, sampled without candidate output. Duplicates retained.','python_chess_version':chess.__version__,'cases':cases},indent=2)+'\n')
print('Independent oracle positions:',len(cases),'unique FENs:',len(set(fens)))
