"""Offline client wrapper check with the public placeholder-token config."""
import sys,tempfile,logging
from pathlib import Path
sys.path.insert(0,'/Users/philyawj/projects/lichess-bot')
import chess,chess.engine,yaml
from lib.config import load_config
from lib.engine_wrapper import create_engine
logging.getLogger().setLevel(logging.ERROR)
root=Path.cwd()
binary=root/'tools/benchmark/bin/pandafish-candidate'
config=yaml.safe_load((root/'lichess-bot-config-example.yml').read_text())
assert config['token']=='YOUR_LICHESS_BOT_API_TOKEN'
config['engine'].update(dir=str(binary.parent),name=binary.name,working_dir=str(root))
assert config['engine']['ponder'] is False
with tempfile.NamedTemporaryFile(mode='w',suffix='.yml') as f:
 yaml.safe_dump(config,f);f.flush();loaded=load_config(f.name)
 with create_engine(loaded) as wrapper:
  assert wrapper.engine.id['name']=='Pandafish v0.15.2'
  board=chess.Board()
  for _ in range(12):
   result=wrapper.engine.play(board,chess.engine.Limit(time=.1),ponder=False)
   assert result.move in board.legal_moves;board.push(result.move)
  mate=chess.Board('7k/5Q2/5K2/8/8/8/8/8 w - - 0 1')
  result=wrapper.engine.play(mate,chess.engine.Limit(depth=3),info=chess.engine.INFO_SCORE)
  assert result.move in mate.legal_moves
  mate.push(result.move);assert mate.is_checkmate()
  draw=chess.Board('4k1n1/8/8/8/8/8/8/3QK1N1 w - - 0 1')
  for move in 'g1f3 g8f6 f3g1 f6g8 g1f3 g8f6 f3g1'.split():draw.push_uci(move)
  result=wrapper.engine.play(draw,chess.engine.Limit(depth=2),info=chess.engine.INFO_SCORE)
  assert result.move.uci()=='f6g8' and result.info['score'].pov(chess.BLACK).score()==0
print('PASS: public placeholder-token config; candidate identity; 12 legal plies; mate; history-aware draw. No live client or online game started. Active binary/config untouched.')
