"""Run exactly the routine 150-game screen once; preserve any failed decision."""
import datetime,hashlib,json,os,subprocess,time
from pathlib import Path
r=Path(__file__).resolve().parent;repo=r.parents[3]
if (r/'screen-command.json').exists():raise SystemExit('Refuse a second screen')
source=json.loads((r/'source-hashes.json').read_text())
for name,entry in source.items():assert hashlib.sha256((repo/name).read_bytes()).hexdigest()==entry['sha256'],name
p=json.loads((r/'provenance.json').read_text())
for label,path in [('candidate','tools/benchmark/bin/pandafish-v0.18-exact-evaluation-dev'),('baseline','tools/benchmark/bin/pandafish-v0.17'),('sacrifice','tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine'),('sacrifice_Tal','tools/benchmark/bin/sacrifice-source/data/Tal.pgn')]:assert hashlib.sha256((repo/path).read_bytes()).hexdigest()==p['binary_sha256'][label],label
assert hashlib.sha256(Path('/opt/homebrew/bin/stockfish').read_bytes()).hexdigest()==p['binary_sha256']['stockfish']
a=json.loads((r/'analysis.json').read_text());assert a['fixed']['groups']['4/all']['elapsed_ratio']<=.90 and a['deep']['groups']['6/all']['elapsed_ratio']<=.90;assert not a['fixed']['failures']and not a['deep']['failures'];assert a['uci']['requests']==2046
assert all(not b['failures']for b in a['uci']['blocks'].values())
command=[str(repo/'tools/benchmark/bin/benchmark'),'--candidate','tools/benchmark/bin/pandafish-v0.18-exact-evaluation-dev','--version','v0.18-exact-evaluation-dev','--revision',p['base_commit']+'-dirty','--previous','tools/benchmark/bin/pandafish-v0.17','--previous-version','v0.17','--previous-revision',p['base_commit'],'--stockfish','/opt/homebrew/bin/stockfish','--stockfish-skill','2','--sacrifice','tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine','--pairs','25','--parallel-all','--verbose','--machine','darwin/arm64; 8 logical CPUs; three concurrent matchups; no other engine probes or live client; power state captured in host.txt; no thermal telemetry']
record={'command':command,'started_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'gates':'race/vet/build/wrapper; all19200 table entries; 596 original-formula base-score equivalence cases; 192 paired exact node/score/move searches; 2046 valid UCI probes; frozen source and binary hashes verified'}
(r/'screen-command.json').write_text(json.dumps(record,indent=2)+'\n')
start=time.monotonic()
with (r/'screen-run.txt').open('w')as log:
 process=subprocess.run(command,cwd=repo,stdout=log,stderr=subprocess.STDOUT)
record.update(exit=process.returncode,elapsed_s=time.monotonic()-start,finished_utc=datetime.datetime.now(datetime.timezone.utc).isoformat())
(r/'screen-command.json').write_text(json.dumps(record,indent=2)+'\n')
print('Screen exit:',process.returncode,'all outputs preserved; no confirmation launched.',flush=True)
raise SystemExit(process.returncode)
