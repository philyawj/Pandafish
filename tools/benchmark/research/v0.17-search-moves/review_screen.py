"""Review the completed screen, including independent PGN legality and all losses."""
import collections, hashlib, json, math, statistics
from pathlib import Path
import chess, chess.pgn
r=Path(__file__).resolve().parent;repo=r.parents[3]
dir=repo/'tools/benchmark/results/v0.17-search-moves-dev-vs-v0.16-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba'
s=json.loads((dir/'summary.json').read_text());assert s['complete']and s['planned_matches']==3
assert s['settings']['matchup_concurrency']==3 and s['settings']['pairs']==25
assert s['settings']['move_ms']==100 and s['settings']['timeout_grace_ms']==50
assert s['settings']['openings_sha256']=='80ceed68e4575b7a0d8649aa67efd35da4c8c46e480ad31b1ce12fefa94e89ca'
p=json.loads((r/'provenance.json').read_text());assert s['candidate']['sha256']==p['binary_sha256']['candidate']
assert s['candidate']['source_revision']==p['base_commit']+'-dirty'
review={'result_directory':str(dir.relative_to(repo)),'previous_gate':s['previous_gate'],'settings':s['settings'],'candidate':s['candidate'],'matches':[]}
for m in s['matches']:
 games=m['games'];assert len(games)==50
 stats=m['stats'];rows=[];points=[];opening=collections.defaultdict(list);by_color=collections.defaultdict(list)
 with (dir/m['pgn_file']).open()as f:
  for index,meta in enumerate(games):
   g=chess.pgn.read_game(f);assert g is not None and not g.errors
   board=g.board();plies=0
   for move in g.mainline_moves():
    assert move in board.legal_moves,(m['engine']['label'],index,move)
    board.push(move);plies+=1
   result=g.headers['Result'];assert result==meta['result']
   point=.5 if result=='1/2-1/2'else 1.0 if result==('1-0'if meta['candidate_white']else'0-1')else 0.0
   if result=='*':point=0.0
   if meta['termination']=='Checkmate':assert board.is_checkmate()
   points.append(point);opening[' '.join(meta['opening'])].append(point);by_color['white'if meta['candidate_white']else'black'].append(point)
   rows.append({'game':index+1,'candidate_white':meta['candidate_white'],'result':result,'points':point,'termination':meta['termination'],'plies':plies,'final_fen':board.fen()})
  assert chess.pgn.read_game(f)is None
 pairs=[]
 for i in range(0,50,2):
  assert games[i]['opening']==games[i+1]['opening']and games[i]['candidate_white']!=games[i+1]['candidate_white']
  pairs.append((points[i]+points[i+1])/2)
 mean=statistics.mean(pairs);margin=2.0639*statistics.stdev(pairs)/math.sqrt(len(pairs))
 expected=(stats['wins']+.5*stats['draws'])/50
 assert abs(mean-expected)<1e-10
 responses=[x['elapsed_ms']for game in games for x in game['response_times']if x['candidate']]
 review['matches'].append({'engine':m['engine'],'stats':stats,'pgn_sha256':hashlib.sha256((dir/m['pgn_file']).read_bytes()).hexdigest(),'pgn_legal_games':len(rows),'color_points':{k:{'games':len(v),'points':sum(v),'score_percent':100*statistics.mean(v)}for k,v in by_color.items()},'opening_points':{k:{'games':len(v),'points':sum(v),'score_percent':100*statistics.mean(v)}for k,v in opening.items()},'paired_scores':pairs,'approx_pair_t95_percent':[100*max(0,mean-margin),100*min(1,mean+margin)],'interval_caveat':'Illustrative t interval over 25 color pairs; repeated public openings and shared-host effects violate full independence. Not an Elo estimate or sequential significance test.','candidate_response_count':len(responses),'candidate_response_max_ms':max(responses),'candidate_response_mean_ms':statistics.mean(responses),'games':rows})
review['notes']='All results retained. Historical external scores are context, not a concurrent released-baseline external control or a causal attribution to any individual component.'
(r/'screen-review.json').write_text(json.dumps(review,indent=2)+'\n')
for m in review['matches']:print(m['engine']['label'],m['stats'],'pair interval',m['approx_pair_t95_percent'],'legal PGNs',m['pgn_legal_games'])
