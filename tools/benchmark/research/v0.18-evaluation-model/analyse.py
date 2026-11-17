"""Keep evaluator-induced score/tree changes separate from performance and strength."""
import collections,json,statistics
from pathlib import Path
r=Path(__file__).resolve().parent;report={}
for mode in ['fixed','deep']:
 a=json.loads((r/('baseline-'+mode+'.json')).read_text());b=json.loads((r/('candidate-'+mode+'.json')).read_text());index={(x['ID'],x['Depth']):x for x in a};groups=collections.defaultdict(lambda:{'cases':0,'baseline_nodes':0,'candidate_nodes':0,'baseline_ms':0,'candidate_ms':0,'changed_moves':0,'changed_scores':0});failures=[]
 for x in b:
  old=index[x['ID'],x['Depth']]
  if not(old['Complete']and x['Complete']and old['Legal']and x['Legal']and len(x['Best'])==1):failures.append({'baseline':old,'candidate':x})
  g=groups[str(x['Depth'])];g['cases']+=1;g['baseline_nodes']+=old['Nodes'];g['candidate_nodes']+=x['Nodes'];g['baseline_ms']+=old['MS'];g['candidate_ms']+=x['MS'];g['changed_moves']+=int(x['Best']!=old['Best']);g['changed_scores']+=int(x['Score']!=old['Score'])
 for g in groups.values():g['node_ratio']=g['candidate_nodes']/g['baseline_nodes'];g['time_ratio']=g['candidate_ms']/g['baseline_ms']
 report[mode]={'rows':len(b),'failures':failures,'depths':dict(groups),'note':'Different evaluation changes scores, decisions and pruning; node/time ratios are not fixed-work speedups.'}
rows=json.loads((r/'uci-performance.json').read_text())['records'];blocks={}
for load,budget in [('serial',100),('serial',1000),('three-search',100),('warm',100)]:
 for actor in ['baseline','candidate']:
  selected=[x for x in rows if x['actor']==actor and x['load']==load and x['budget_ms']==budget];good=[x for x in selected if x['completed']and not x['terminal']]
  blocks[f'{load}/{budget}/{actor}']={'requests':len(selected),'mean_depth':statistics.mean(x['depth']for x in good),'total_nodes':sum(x['nodes']for x in good),'max_response_ms':max(x['response_ms']for x in good),'terminal_nulls':sum(x['terminal']and x['move']=='0000'for x in selected),'failures':[x for x in selected if not x['completed']]}
report['uci']={'requests':len(rows),'blocks':blocks}
(r/'analysis.json').write_text(json.dumps(report,indent=2)+'\n')
print(json.dumps(report,indent=2))
assert len(rows)==2046 and not report['fixed']['failures']and not report['deep']['failures']and all(not b['failures']for b in blocks.values())
