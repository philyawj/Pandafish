"""Freeze inputs before candidate measurements; no selection by engine output."""
import hashlib, json
from pathlib import Path
import chess, chess.pgn

root = Path(__file__).resolve().parent
repo = root.parents[3]
old = json.loads((root.parent/'v0.15.3-architecture/inputs.json').read_text())
deep_ids = {c['id'] for c in json.loads((root.parent/'v0.15.3-architecture/deep-inputs.json').read_text())['cases']}
cases = [dict(c, deep=c['id'] in deep_ids) for c in old['cases']]
control = repo/'tools/benchmark/results/v0.15.2-architecture-control-vs-v0.14.4-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba'
sources = {}
for opponent in ['previous','stockfish','sacrifice']:
 path = control/(opponent+'.pgn')
 if not path.exists():
  raise SystemExit('Missing frozen source: '+str(path))
 sources[str(path.relative_to(repo))] = hashlib.sha256(path.read_bytes()).hexdigest()
 with path.open() as f:
  for i in range(12):
   g = chess.pgn.read_game(f)
   assert g is not None and not g.errors
   moves = list(g.mainline_moves())
   target = [12,25,40,65,90,120][i%6]
   ply = min(target,len(moves)-2)
   b = g.board()
   for m in moves[:ply]: b.push(m)
   assert not b.is_game_over(claim_draw=True)
   cases.append({'id':f'broader-{opponent}-{i+1:02}-ply{ply}', 'group':'broader-'+opponent,
    'position':'position startpos moves '+' '.join(m.uci() for m in moves[:ply]),
    'fen':b.fen(), 'next_two':[m.uci() for m in moves[ply:ply+2]],
    'source_round':g.headers.get('Round'), 'chosen_ply':ply,'target_ply':target,'deep':i<2})
fixtures = [
 ('castling','r3k2r/ppp2ppp/2n5/3pp3/3PP3/2N5/PPP2PPP/R3K2R w KQkq - 0 1'),
 ('castling-black','r3k2r/ppp2ppp/2n5/3pp3/3PP3/2N5/PPP2PPP/R3K2R b KQkq - 17 1'),
 ('ep','4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1'),
 ('ep-pinned','k3r3/8/8/3pP3/8/8/8/4K3 w - d6 0 1'),
 ('quiet-promotion','7k/P7/8/8/8/8/8/4K3 w - - 0 1'),
 ('capture-promotion','r6k/1P6/8/8/8/8/8/7K w - - 0 1'),
 ('underpromotion','8/k1P5/2K5/8/8/8/8/8 w - - 0 1'),
 ('quiet-evasion','7r/8/8/8/8/6k1/8/6BK w - - 0 1'),
 ('mate-before-50','7k/5Q2/5K2/8/8/8/8/8 w - - 99 1'),
 ('resist-mate','7k/8/5Q2/5K2/8/8/8/8 b - - 0 1'),
 ('pawn-ending','8/8/4k3/3p4/3P4/4K3/8/8 w - - 34 1'),
 ('quiet-mate-threat','r6r/pppbn1Q1/2nqp3/5p2/3P1P1k/2B4P/PPP3P1/R4RK1 w - - 1 19'),
]
for name,fen in fixtures:
 b=chess.Board(fen);assert b.is_valid() and len(list(b.legal_moves))>0
 cases.append({'id':'special-'+name,'group':'special','position':'position fen '+fen,'fen':fen,'next_two':[],'deep':False})
for c in cases:
 assert chess.Board(c['fen']).is_valid(), c['id']
output={'method':'36 original diagnostics + first 12 games each from the frozen released-policy control previous/Stockfish/Sacrifice PGNs, at cycling plies 12/25/40/65/90/120 clamped to len-2 + 12 explicit rule/threat controls. No candidate-result selection. Deep subset: original six + first two new games per opponent. Public and correlated, not blind training/validation data.','sources_sha256':sources,'cases':cases}
for name, selected in [('inputs.json',cases),('deep-inputs.json',[c for c in cases if c['deep']])]:
 out=root/name
 if out.exists():raise SystemExit('Refuse overwrite: '+str(out))
 out.write_text(json.dumps(dict(output,cases=selected),indent=2)+'\n')
print('Frozen inputs:',len(cases),'deep:',sum(c['deep'] for c in cases),'unique commands:',len({c['position']for c in cases}))
