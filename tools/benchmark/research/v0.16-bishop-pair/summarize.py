"""Summarize preserved fixed-depth, equal-time and finite-node referee records."""
import collections,json,statistics
from pathlib import Path
r=Path(__file__).resolve().parent
load=lambda n:json.loads((r/n).read_text())
a=load('baseline-fixed.json');b=load('candidate-fixed.json')
ca=load('clean-1-baseline.json')+load('clean-4-baseline.json');cb=load('clean-2-candidate.json')+load('clean-3-candidate.json')
d={'all_complete_legal':all(x['Complete'] and x['Legal'] for x in a+b+ca+cb),'fixed_best_sets_changed':sum(set(x['Moves'])!=set(y['Moves']) for x,y in zip(a,b)),'probes_per_engine':len(a)+len(ca),'clean_probes_per_engine':len(ca),'clean_depth_higher_equal_lower':dict(collections.Counter('higher' if y['Depth']>x['Depth'] else 'lower' if y['Depth']<x['Depth'] else 'equal' for x,y in zip(ca,cb))),'clean_node_ratio':sum(x['Nodes'] for x in cb)/sum(x['Nodes'] for x in ca)}
for label,rows in [('baseline',ca),('candidate',cb)]:d[label]={'nodes':sum(x['Nodes'] for x in rows),'mean_depth':statistics.mean(x['Depth'] for x in rows),'mean_elapsed_ms':statistics.mean(x['ElapsedMS'] for x in rows),'max_elapsed_ms':max(x['ElapsedMS'] for x in rows)}
sf={(x['game'],x['ply'],x['move']):x['score'] for x in load('referee-results.json')['results']};groups=collections.defaultdict(list);mates=[];changes=[]
for x,y in list(zip(a,b))+list(zip(ca,cb)):
 av=[sf[x['Game'],x['Ply'],m] for m in x['Moves']];bv=[sf[y['Game'],y['Ply'],m] for m in y['Moves']]
 if any('mate'in v for v in av+bv):mates.append({'game':x['Game'],'mode':x['Mode'],'depth':x['Depth'],'baseline':av,'candidate':bv});continue
 delta=statistics.mean(v['cp'] for v in bv)-statistics.mean(v['cp'] for v in av)
 groups[x['Mode'],x['Depth'] if x['Mode']=='fixed' else 0,x['Split']].append(delta)
 if x['Mode']=='fixed' and x['Moves']!=y['Moves']:changes.append({'game':x['Game'],'split':x['Split'],'depth':x['Depth'],'baseline_moves':x['Moves'],'candidate_moves':y['Moves'],'mean_referee_cp_delta':delta})
d['referee']={'method':'Every fixed best tie averaged, each timed selection scored; repeated cases and geometric mirrors are correlated, no Elo inference. Positive cp favors candidate. Mate transitions kept separate. Existing identical full-history root queries reused from piece-safety archive; newly selected moves get same 100000-node fresh-hash configuration.','groups':[{'mode':k[0],'depth':k[1],'split':k[2],'n':len(v),'mean_cp_delta':statistics.mean(v),'better_same_worse':[sum(x>0 for x in v),sum(x==0 for x in v),sum(x<0 for x in v)]}for k,v in groups.items()],'mates':mates,'fixed_changes':changes}
(r/'assessment.json').write_text(json.dumps(d,indent=2)+'\n')
print(json.dumps(d,indent=2))
