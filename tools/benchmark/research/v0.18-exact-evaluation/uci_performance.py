"""Balanced native UCI probes: equal time, serial/three-search load and warm play.
Run once after validate.py. No games or network services are started.
"""
import concurrent.futures, hashlib, json, os, queue, subprocess, threading, time
from pathlib import Path
import chess
r=Path(__file__).resolve().parent;repo=Path.cwd()
cases=json.loads((r/'inputs.json').read_text())['cases']
paths={'baseline':repo/'tools/benchmark/bin/pandafish-v0.17','candidate':repo/'tools/benchmark/bin/pandafish-v0.18-exact-evaluation-dev'}
class Engine:
 def __init__(self,path):
  self.p=subprocess.Popen([str(path)],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,bufsize=1)
  self.q=queue.Queue()
  def reader():
   for line in self.p.stdout:self.q.put(line.rstrip())
   self.q.put(None)
  threading.Thread(target=reader,daemon=True).start()
  self.send('uci');self.handshake=self.wait('uciok',10)
 def send(self,s):self.p.stdin.write(s+'\n');self.p.stdin.flush()
 def wait(self,prefix,seconds):
  out=[];until=time.monotonic()+seconds
  while True:
   line=self.q.get(timeout=max(.001,until-time.monotonic()))
   if line is None:raise RuntimeError('unexpected engine exit')
   out.append(line)
   if line==prefix or line.startswith(prefix+' '):return out
 def close(self):
  try:self.send('quit');self.p.wait(timeout=3)
  except Exception:self.p.kill();self.p.wait()
def board(command):
 words=command.split();index=words.index('moves')if'moves'in words else len(words)
 b=chess.Board()if words[1]=='startpos'else chess.Board(' '.join(words[2:index]))
 for move in words[index+1:]:b.push_uci(move)
 return b
outpath=r/'uci-performance.json'
if outpath.exists():raise SystemExit('Refuse overwrite')
records=[];lock=threading.Lock()
report={'method':'Native UCI. Same requested budgets and unchanged hard/soft margins. Serial 100 ms trials alternate actor/case order; serial 1000 ms is one trial. Three-search load runs three independent processes, two actor blocks with reversed order, three lanes/trials per block. Warm episodes use real generation increments; terminal null moves are expected. No strength claim.','logical_cpus':os.cpu_count(),'identities':{k:{'path':str(p),'sha256':hashlib.sha256(p.read_bytes()).hexdigest()}for k,p in paths.items()},'records':records}
def save(rec):
 with lock:
  records.append(rec);outpath.write_text(json.dumps(report,indent=2)+'\n')
def lane(label,trial,budget,load,inputs):
 engine=Engine(paths[label]);count=0
 try:
  with lock:report['identities'][label]['handshake']=engine.handshake
  for c in inputs:
   rec={'actor':label,'trial':trial,'budget_ms':budget,'load':load,'id':c['id'],'group':c['group'],'position':c['position']}
   try:
    if not c.get('warm'):engine.send('ucinewgame')
    engine.send('isready');engine.wait('readyok',10)
    started=time.monotonic();engine.send(c['position']);engine.send(f'go movetime {budget}')
    lines=engine.wait('bestmove',budget/1000+5)
    rec.update(response_ms=(time.monotonic()-started)*1000,output=lines,move=lines[-1].split()[1])
    b=board(c['position']);rec['terminal']=not any(b.legal_moves)
    rec['legal']=chess.Move.from_uci(rec['move'])in b.legal_moves
    rec['valid']=rec['legal']or(rec['terminal']and rec['move']=='0000')
    metrics=[json.loads(line.removeprefix('info string pandafish '))for line in lines if line.startswith('info string pandafish ')]
    if not metrics:
     # Preserve and parse the actual telemetry prefix from regular UCI info if
     # the engine's structured prefix changes; never invent a completed depth.
     info=[line.split()for line in lines if line.startswith('info depth ')]
     last=info[-1]if info else[]
     rec['depth']=int(last[last.index('depth')+1])if'depth'in last else 0
     rec['nodes']=int(last[last.index('nodes')+1])if'nodes'in last else 0
    else:rec.update(depth=metrics[-1]['depth'],nodes=metrics[-1]['nodes'],metrics=metrics[-1])
    rec['completed']=rec['valid']and(rec['depth']>0 or rec['terminal'])
   except Exception as error:rec.update(error=repr(error),valid=False,completed=False)
   save(rec);count+=1
   if 'error'in rec:break
 finally:engine.close()
 print(load,trial,label,budget,count,flush=True)
for trial in range(3):
 for label in (['baseline','candidate']if trial%2==0 else['candidate','baseline']):
  lane(label,trial,100,'serial',cases if trial%2==0 else cases[::-1])
for label in ['candidate','baseline']:lane(label,0,1000,'serial',cases)
for block in range(2):
 for label in (['baseline','candidate']if block==0 else['candidate','baseline']):
  with concurrent.futures.ThreadPoolExecutor(max_workers=3)as pool:
   futures=[pool.submit(lane,label,block*3+i,100,'three-search',cases if i%2==0 else cases[::-1])for i in range(3)]
   for f in futures:f.result()
episodes=[]
for c in cases:
 if len(c.get('next_two',[]))!=2:continue
 position=c['position']+(''if' moves 'in c['position']else' moves')+' '+' '.join(c['next_two'])
 episodes.extend([dict(c,id=c['id']+'/parent'),dict(c,id=c['id']+'/warm-next',position=position,warm=True),dict(c,id=c['id']+'/cold-next',position=position)])
for label in ['baseline','candidate']:lane(label,0,100,'warm',episodes)
print('All native UCI probes preserved:',len(records),flush=True)
if any(not row['completed']for row in records):raise SystemExit('One or more probes failed; inspect raw records')
