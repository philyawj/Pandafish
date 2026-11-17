"""Native fast-base-only control after the main timing blocks, no overlap."""
import hashlib,json,time
from pathlib import Path
# Reuse definitions only; do not execute the main timing schedule.
source=(Path(__file__).resolve().parent/'uci_performance.py').read_text();exec(source[:source.index('outpath=')])
path=repo/'tools/benchmark/bin/pandafish-v0.18-fast-base-diagnostic';output=r/'ablation-uci.json'
assert not output.exists()
records=[];report={'actor':'fast-base','binary_sha256':hashlib.sha256(path.read_bytes()).hexdigest(),'method':'Residual disabled; same exact fast base kernel, fresh ucinewgame, 100ms three trials and 1000ms one trial, no other engine probes. Diagnostic ablation, not a strength screen.','records':records}
for budget,repeats in [(100,3),(1000,1)]:
 for trial in range(repeats):
  engine=Engine(path);report['handshake']=engine.handshake
  try:
   for c in(cases if trial%2==0 else cases[::-1]):
    row={'actor':'fast-base','load':'serial','budget_ms':budget,'trial':trial,'id':c['id'],'group':c['group'],'position':c['position']}
    try:
     engine.send('ucinewgame');engine.send('isready');engine.wait('readyok',10)
     start=time.monotonic();engine.send(c['position']);engine.send(f'go movetime {budget}');lines=engine.wait('bestmove',budget/1000+5);elapsed=(time.monotonic()-start)*1000
     b=board(c['position']);move=lines[-1].split()[1];terminal=not any(b.legal_moves);legal=chess.Move.from_uci(move)in b.legal_moves
     m=[json.loads(line.removeprefix('info string pandafish '))for line in lines if line.startswith('info string pandafish ')][-1]
     row.update(output=lines,response_ms=elapsed,move=move,terminal=terminal,legal=legal,valid=legal or(terminal and move=='0000'),depth=m['depth'],nodes=m['nodes'],metrics=m);row['completed']=row['valid']and(row['depth']>0 or terminal)
    except Exception as error:row.update(error=repr(error),completed=False,valid=False)
    records.append(row);output.write_text(json.dumps(report,indent=2)+'\n')
    assert row['completed'],row
  finally:engine.close()
  print('Fast-base-only',budget,'trial',trial,'complete',flush=True)
assert len(records)==336
