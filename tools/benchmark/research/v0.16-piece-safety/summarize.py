import collections,json,statistics
from pathlib import Path
r=Path(__file__).resolve().parent
out={}
def load(n):return json.loads((r/n).read_text())
a=load('baseline-probes.json');b=load('candidate-probes.json')
out['initial']={'note':'First baseline partly overlaps static microbenchmark; fixed scores valid, elapsed not final cost evidence.', 'records_per_engine':len(a),'all_complete_legal':all(x['Complete'] and x['Legal'] for x in a+b)}
for mode in ['fixed','100ms']:
 pairs=[(x,y) for x,y in zip(a,b) if x['Mode']==mode]
 out['initial'][mode]={'best_sets_changed':sum(set(x['Moves'])!=set(y['Moves']) for x,y in pairs),'baseline_nodes':sum(x['Nodes'] for x,y in pairs),'candidate_nodes':sum(y['Nodes'] for x,y in pairs)}
if (r/'clean-4-baseline.json').exists():
 clean_a=load('clean-1-baseline.json')+load('clean-4-baseline.json');clean_b=load('clean-2-candidate.json')+load('clean-3-candidate.json')
 out['clean_100ms']={'schedule':'baseline,candidate,candidate,baseline; three repeats per position per round; no parallel engine task','probes_per_engine':len(clean_a),'all_complete_legal':all(x['Complete'] and x['Legal'] for x in clean_a+clean_b)}
 for label,records in [('baseline',clean_a),('candidate',clean_b)]:
  out['clean_100ms'][label]={'nodes':sum(x['Nodes'] for x in records),'mean_depth':statistics.mean(x['Depth'] for x in records),'mean_elapsed_ms':statistics.mean(x['ElapsedMS'] for x in records),'max_elapsed_ms':max(x['ElapsedMS'] for x in records)}
 pairs=list(zip(clean_a,clean_b));out['clean_100ms']['depth_higher_equal_lower']=dict(collections.Counter('higher' if y['Depth']>x['Depth'] else 'lower' if y['Depth']<x['Depth'] else 'equal' for x,y in pairs))
 # Time-aligned position/trial slots; random root tie selection can differ.
 out['clean_100ms']['node_ratio']=sum(x['Nodes'] for x in clean_b)/sum(x['Nodes'] for x in clean_a)
if (r/'referee-results.json').exists():
 sf={(x['game'],x['ply'],x['move']):x['score'] for x in load('referee-results.json')['results']};groups=collections.defaultdict(list);mates=[]
 for x,y in zip(a,b):
  av=[sf[x['Game'],x['Ply'],m] for m in x['Moves']];bv=[sf[y['Game'],y['Ply'],m] for m in y['Moves']]
  if any('mate' in v for v in av+bv):mates.append({'case':x['Game'],'mode':x['Mode'],'depth':x['Depth'],'baseline':av,'candidate':bv});continue
  delta=statistics.mean(v['cp'] for v in bv)-statistics.mean(v['cp'] for v in av)
  groups[x['Mode'],x['Depth'] if x['Mode']=='fixed' else 0,x['Split']].append(delta)
 out['initial_referee_comparison']={'method':'Mean referee cp over every tied best move; repeated timed choices correlated; positive favors candidate. Selected errors, no Elo inference. Mate excluded and separately saved.','groups':[{'mode':k[0],'depth':k[1],'split':k[2],'n':len(v),'mean_cp_delta':statistics.mean(v),'better_same_worse':[sum(x>0 for x in v),sum(x==0 for x in v),sum(x<0 for x in v)]}for k,v in groups.items()],'mates':mates}
if (r/'referee-clean-results.json').exists():
 sf={(x['game'],x['ply'],x['move']):x['score'] for x in load('referee-clean-results.json')['results']};groups=collections.defaultdict(list);mates=[]
 for x,y in zip(clean_a,clean_b):
  av=sf[x['Game'],x['Ply'],x['Moves'][0]];bv=sf[y['Game'],y['Ply'],y['Moves'][0]]
  if 'cp' in av and 'cp' in bv:groups[x['Split']].append(bv['cp']-av['cp'])
  else:mates.append({'game':x['Game'],'baseline':av,'candidate':bv})
 out['clean_referee_comparison']={'method':'Position/trial slots from balanced rounds; repeated same-case scores highly correlated, not 240 independent samples. Mate cases separate.', 'groups':[{'split':k,'n':len(v),'mean_cp_delta':statistics.mean(v),'better_same_worse':[sum(x>0 for x in v),sum(x==0 for x in v),sum(x<0 for x in v)]}for k,v in groups.items()],'mates':mates}
if (r/'leaf-safety.json').exists():
 leaves=load('leaf-safety.json');out['leaf_last']=[x for i,x in enumerate(leaves) if i+1==len(leaves) or (x['game'],x['kind'],x['root'])!=(leaves[i+1]['game'],leaves[i+1]['kind'],leaves[i+1]['root'])]
(r/'assessment.json').write_text(json.dumps(out,indent=2)+'\n')
print(json.dumps(out,indent=2))
