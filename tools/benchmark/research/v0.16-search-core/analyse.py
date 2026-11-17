"""Summarize all probes without dropping failures, ties or terminal requests."""
import collections, json, statistics
from pathlib import Path
r=Path(__file__).resolve().parent
result={}
for mode in ['fixed','deep']:
 a=json.loads((r/f'baseline-{mode}.json').read_text());b=json.loads((r/f'candidate-{mode}.json').read_text())
 index={(x['ID'],x['Depth']):x for x in a};assert len(index)==len(a)
 groups=collections.defaultdict(lambda:{'cases':0,'baseline_nodes':0,'candidate_nodes':0,'baseline_ms':0,'candidate_ms':0})
 failures=[]
 for row in b:
  old=index[(row['ID'],row['Depth'])]
  if not (row['Complete']and old['Complete']and row['Legal']and old['Legal']and row['Score']==old['Score']and len(row['Best'])==1 and row['Best'][0]in old['Best']):failures.append({'baseline':old,'candidate':row})
  for label in ['all',row['Group']]:
   g=groups[str(row['Depth'])+'/'+label];g['cases']+=1;g['baseline_nodes']+=old['Nodes'];g['candidate_nodes']+=row['Nodes'];g['baseline_ms']+=old['MS'];g['candidate_ms']+=row['MS']
 for g in groups.values():
  g['node_ratio']=g['candidate_nodes']/g['baseline_nodes'];g['elapsed_ratio']=g['candidate_ms']/g['baseline_ms']
 result[mode]={'rows':len(b),'failures':failures,'groups':dict(groups)}
p=r/'uci-performance.json'
if p.exists():
 rows=json.loads(p.read_text())['records'];groups=collections.defaultdict(list)
 for x in rows:groups[(x['load'],x['budget_ms'],x['actor'])].append(x)
 blocks={}
 for key,group in groups.items():
  valid=[x for x in group if x['completed']and not x['terminal']]
  response=sorted(x['response_ms']for x in valid)
  blocks['/'.join(map(str,key))]={'requests':len(group),'nonterminal_completed':len(valid),'terminal_nulls':sum(x.get('terminal',False)and x.get('move')=='0000'for x in group),'failures':[x for x in group if not x['completed']], 'mean_depth':statistics.mean(x['depth']for x in valid)if valid else None,'total_nodes':sum(x['nodes']for x in valid),'mean_response_ms':statistics.mean(response)if response else None,'p95_response_ms':response[min(len(response)-1,int(.95*len(response)))]if response else None,'max_response_ms':max(response)if response else None}
 pairs={}
 for load,budget in [('serial',100),('serial',1000),('three-search',100)]:
  d={ (x['actor'],x['id'],x['trial']):x for x in rows if x['load']==load and x['budget_ms']==budget }
  gains=[];by_group=collections.defaultdict(list)
  for (actor,id,trial),new in d.items():
   if actor!='candidate':continue
   old=d.get(('baseline',id,trial))
   if old and new.get('completed')and old.get('completed')and not new.get('terminal'):
    gain=new['depth']-old['depth'];gains.append(gain);by_group[new['group']].append(gain)
  pairs[f'{load}/{budget}']={'pairs':len(gains),'deeper':sum(x>0 for x in gains),'equal':sum(x==0 for x in gains),'shallower':sum(x<0 for x in gains),'mean_depth_gain':statistics.mean(gains)if gains else None,'by_group':{k:{'pairs':len(v),'mean_depth_gain':statistics.mean(v)}for k,v in by_group.items()}}
 warm={}
 for actor in ['baseline','candidate']:
  d={x['id']:x for x in rows if x['load']=='warm'and x['actor']==actor};gains=[];nodes=[0,0]
  for id,new in d.items():
   if not id.endswith('/warm-next'):continue
   old=d[id.removesuffix('/warm-next')+'/cold-next']
   if new['completed']and old['completed']and not new['terminal']:
    gains.append(new['depth']-old['depth']);nodes[0]+=new['nodes'];nodes[1]+=old['nodes']
  warm[actor]={'nonterminal_pairs':len(gains),'warm_deeper':sum(x>0 for x in gains),'equal':sum(x==0 for x in gains),'warm_shallower':sum(x<0 for x in gains),'warm_nodes':nodes[0],'cold_nodes':nodes[1]}
 result['uci']={'requests':len(rows),'blocks':blocks,'paired_depths':pairs,'warm':warm}
(r/'analysis.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({'fixed_rows':result['fixed']['rows'],'fixed_failures':len(result['fixed']['failures']),'deep_rows':result['deep']['rows'],'deep_failures':len(result['deep']['failures']),'depth4':result['fixed']['groups']['4/all'],'depth6':result['deep']['groups']['6/all'],'uci_pairs':result.get('uci',{}).get('paired_depths')},indent=2))
if result['fixed']['failures']or result['deep']['failures']:raise SystemExit('Parity/membership failures: inspect all saved rows')
