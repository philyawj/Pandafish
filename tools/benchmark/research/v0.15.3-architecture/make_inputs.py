"""Freeze inputs before architecture probes. No result-based selection."""
import hashlib,json
from pathlib import Path
import chess,chess.pgn
root=Path.cwd();out=root/'tools/benchmark/research/v0.15.3-architecture'
source=root/'tools/benchmark/results/v0.15-conversion-dev-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/previous.pgn'
cases=[]
with source.open() as f:
 for index in range(24):
  g=chess.pgn.read_game(f);assert g is not None
  moves=list(g.mainline_moves());target=[16,32,48,64][index%4];ply=min(target,len(moves)-2);assert ply>=8
  b=g.board()
  for move in moves[:ply]:b.push(move)
  assert not b.is_game_over()
  cases.append({'id':f'neutral-{index+1:02d}-ply{ply}','group':'neutral','fen':b.fen(en_passant='fen'),'position':'position startpos moves '+' '.join(m.uci() for m in b.move_stack),'next_two':[m.uci() for m in moves[ply:ply+2]],'source_round':g.headers.get('Round'),'chosen_ply':ply,'target_ply':target})
cases += [{'id':'start','group':'control','fen':chess.STARTING_FEN,'position':'position startpos','next_two':['e2e4','e7e5']},{'id':'sparse-transpositions','group':'control','fen':'4k1n1/8/2n5/8/8/2N5/8/4K1N1 w - - 0 1','position':'position fen 4k1n1/8/2n5/8/8/2N5/8/4K1N1 w - - 0 1','next_two':[]},{'id':'known-repetition','group':'rule_control','position':'position fen 4k1n1/8/8/8/8/8/8/3QK1N1 w - - 0 1 moves g1f3 g8f6 f3g1 f6g8 g1f3 g8f6 f3g1','next_two':[]}]
for c in json.loads((root/'tools/benchmark/live/v0.15-live-20261113-074740/analysis/engine-replay/inputs.json').read_text()):cases.append({'id':c['game'].split('/')[-1],'group':'saved_error','fen':c['fen'],'position':c['position'],'next_two':[]})
pawn=json.loads((root/'tools/benchmark/research/v0.16-next-assessment/diagnosis.json').read_text())['records'][2]
cases.append({'id':'Dsq-pawn-ending','group':'saved_ending','fen':pawn['fen'],'position':pawn['position'],'next_two':[]})
for c in cases:
 s=c['position'].split();i=s.index('moves') if 'moves'in s else len(s);b=chess.Board() if s[1]=='startpos' else chess.Board(' '.join(s[2:i]))
 for m in s[i+1:]:b.push_uci(m)
 assert b.is_valid(),c['id'];c['fen']=b.fen(en_passant='fen')
(out/'inputs.json').write_text(json.dumps({'method':'Frozen before architecture probe results. First 24 games from prior conversion previous.pgn, one position per game at cycling plies 16/32/48/64, clamped to len-2 if needed; no selection by result or new-engine output. Plus start/sparse/repetition controls, all eight saved live errors and the diagnosed pawn ending. Public, correlated, not blind or population representative.','source':str(source.relative_to(root)),'source_sha256':hashlib.sha256(source.read_bytes()).hexdigest(),'cases':cases},indent=2)+'\n')
print('Frozen',len(cases),'positions')
