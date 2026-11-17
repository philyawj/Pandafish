"""Descriptive, conditional-on-book intervals; no significance gate or Elo claim."""
import json,random,statistics,sys
from pathlib import Path
summary=Path(sys.argv[1]);data=json.loads(summary.read_text());rng=random.Random(160016);runs=100000;out={'method':'100000 bootstrap resamples, deterministic seed 160016; color pairs and repeated opening clusters. Conditional on this public book and selected candidate; not calibrated rating intervals, SPRT, or correction for multiple experiments. Opponents remain separate.','summary':str(summary),'seed':160016,'resamples':runs,'opponents':[]}
def interval(rows):
 rows.sort();return [rows[int(.025*(len(rows)-1))]*100,rows[int(.975*(len(rows)-1))]*100]
for match in data['matches']:
 games=match['games'];scores=[]
 for g in games:
  result=g['result']
  if result=='1/2-1/2':score=.5
  elif result in ['1-0','0-1']:score=float((result=='1-0')==g['candidate_white'])
  else:score=None
  scores.append(score)
 row={'opponent':match['engine']['label'],'stats':match['stats']}
 if len(scores)%2 or any(x is None for x in scores):row['limitation']='Incomplete pair or unfinished game: no interval inferred.'
 else:
  pairs=[(scores[i]+scores[i+1])/2 for i in range(0,len(scores),2)];clusters={}
  for i,score in enumerate(pairs):clusters.setdefault(tuple(games[2*i]['opening']),[]).append(score)
  ps=[sum(rng.choices(pairs,k=len(pairs)))/len(pairs) for _ in range(runs)]
  cs=[];items=list(clusters.values())
  for _ in range(runs):
   selected=rng.choices(items,k=len(items));cs.append(sum(sum(v) for v in selected)/sum(len(v) for v in selected))
  row.update({'color_pairs':len(pairs),'opening_clusters':len(items),'paired_score_interval_95_percent':interval(ps),'opening_cluster_interval_95_percent':interval(cs)})
 out['opponents'].append(row)
Path(__file__).with_name('sampling.json').write_text(json.dumps(out,indent=2)+'\n')
print(json.dumps(out,indent=2))
