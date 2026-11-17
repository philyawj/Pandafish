"""Serial native/diagnostic UCI probes with raw output and independent watchdogs."""
import hashlib,json,os,queue,subprocess,sys,threading,time
from pathlib import Path
import chess
repo=Path.cwd();r=Path(__file__).resolve().parent;cases=json.loads((r/'inputs.json').read_text())['cases']
class Engine:
 def __init__(self,path):
  self.path=path;self.process=subprocess.Popen([path],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,bufsize=1);self.lines=queue.Queue()
  def read():
   for line in self.process.stdout:self.lines.put(line.strip())
   self.lines.put(None)
  threading.Thread(target=read,daemon=True).start();self.send('uci');self.handshake=self.wait('uciok',10);self.send('ucinewgame');self.send('isready');self.wait('readyok',10)
 def send(self,s):self.process.stdin.write(s+'\n');self.process.stdin.flush()
 def wait(self,prefix,timeout):
  deadline=time.monotonic()+timeout;out=[]
  while True:
   line=self.lines.get(timeout=max(.001,deadline-time.monotonic()))
   if line is None:raise RuntimeError('engine closed stdout')
   out.append(line)
   if line==prefix or line.startswith(prefix+' '):return out
   if time.monotonic()>deadline:raise TimeoutError('missing '+prefix)
 def close(self):
  try:self.send('quit');self.process.wait(timeout=3)
  except Exception:self.process.kill();self.process.wait()
def board(c):
 s=c['position'].split();i=s.index('moves')if'moves'in s else len(s);b=chess.Board()if s[1]=='startpos'else chess.Board(' '.join(s[2:i]))
 for m in s[i+1:]:b.push_uci(m)
 return b
manifest={x['label']:str(repo/'tools/benchmark/bin'/('architecture-'+x['label'])) for x in json.loads((r/'worktrees.json').read_text()) if x['label']!='instrumented-baseline'}
# Native baseline is the exact tagged executable, not the profiled test tree.
manifest['baseline']=str(repo/'tools/benchmark/bin/pandafish-v0.15.2');manifest['v0.14.4']=str(repo/'tools/benchmark/bin/pandafish-v0.14.4-architecture-control');manifest['sacrifice']=str(repo/'tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine')
mode=sys.argv[1] if len(sys.argv)>1 else '100ms'
if mode=='100ms':labels=['baseline','sacrifice','no-tt','board-hints','trim-history','safe-cache','depth-atleast','large-tt','root-single-best','pvs','pvs-safe-cache','tactical-ordering'];budgets=[100];repeats=3
elif mode=='1000ms':labels=['sacrifice','baseline','safe-cache'];budgets=[1000];repeats=1
elif mode=='combo':labels=['baseline','coordinated-kernel'];budgets=[100];repeats=3
elif mode=='historical':labels=['v0.14.4','baseline'];budgets=[100];repeats=3
elif mode=='historical-fixed':labels=['v0.14.4','baseline'];budgets=[0];repeats=1
elif mode=='warm':
 labels=['baseline','sacrifice','safe-cache'];budgets=[100];repeats=1;expanded=[]
 for c in cases:
  if len(c['next_two'])!=2:continue
  parent=dict(c,id=c['id']+'/parent');position=c['position']+('' if ' moves 'in c['position'] else ' moves')+' '+' '.join(c['next_two'])
  expanded += [parent,dict(c,id=c['id']+'/warm-next',position=position,warm_continue=True),dict(c,id=c['id']+'/cold-next',position=position)]
 cases=expanded
else:raise SystemExit('mode must be 100ms or 1000ms')
# Case order and engine order alternate by trial to temper simple order effects.
records=[];metadata={'mode':mode,'method':'One engine searching at a time; native tagged baseline and pinned Sacrifice, diagnostic variants isolated. Fresh ucinewgame each position; actual full history; equal requested budget, unchanged engine margins (Pandafish hard95/995, Sacrifice nominal100/1000). Depth and node definitions differ across engines. No game-strength claim.','identities':{},'records':records}
for trial in range(repeats):
 order=labels if trial%2==0 else labels[::-1]
 for label in order:
  engine=Engine(manifest[label]);metadata['identities'][label]={'path':manifest[label],'sha256':hashlib.sha256(Path(manifest[label]).read_bytes()).hexdigest(),'handshake':engine.handshake}
  try:
   for c in (cases if trial%2==0 else cases[::-1]):
    for budget in budgets:
     rec={'variant':label,'id':c['id'],'group':c['group'],'trial':trial,'requested_ms':budget,'position':c['position']};records.append(rec)
     try:

      if not c.get('warm_continue'):engine.send('ucinewgame')
      engine.send('isready');engine.wait('readyok',10);start=time.monotonic();engine.send(c['position']);engine.send('go depth 4' if mode=='historical-fixed' else f'go movetime {budget}');lines=engine.wait('bestmove',20 if mode=='historical-fixed' else max(5,budget/1000+3));rec['response_ms']=(time.monotonic()-start)*1000;rec['output']=lines;move=lines[-1].split()[1];rec['move']=move;rec['legal']=chess.Move.from_uci(move)in board(c).legal_moves
      infos=[s.split()for s in lines if s.startswith('info depth ')];info=infos[-1]if infos else[]
      rec['depth']=int(info[info.index('depth')+1])if'depth'in info else None;rec['nodes']=int(info[info.index('nodes')+1])if'nodes'in info else None
      rec['completed']=rec['legal']and rec['depth']is not None
     except Exception as e:rec['error']=repr(e);rec['completed']=False;rec['legal']=False;raise
     finally:(r/f'uci-{mode}.json').write_text(json.dumps(metadata,indent=2)+'\n')
  finally:engine.close()
  print(mode,trial,label,'complete',flush=True)
