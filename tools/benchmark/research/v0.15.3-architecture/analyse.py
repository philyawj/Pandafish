"""Generate compact, auditable research tables; no Elo inference."""
import collections,json,statistics
from pathlib import Path
r=Path(__file__).resolve().parent
load=lambda n:json.loads((r/n).read_text())
base=load('baseline-fixed.json');out={'fixed':[],'deep':[],'audit':{},'timed':[],'quality':[],'warm_native':[]}
for row in load('worktrees.json'):
 label=row['label'];p=r/(label+'-fixed.json')
 if not p.exists():continue
 data=load(p.name)
 for depth in [3,4]:
  pairs=[(a,b)for a,b in zip(base,data)if a['Depth']==depth];assert len(pairs)==36
  out['fixed'].append({'variant':label,'depth':depth,'completed_legal':sum(b['Complete']and b['Legal']for a,b in pairs),'nodes':sum(b['Nodes']for a,b in pairs),'node_ratio':sum(b['Nodes']for a,b in pairs)/sum(a['Nodes']for a,b in pairs),'score_differences':sum(a['Score']!=b['Score']for a,b in pairs),'best_set_differences':sum(set(a['Best'])!=set(b['Best'])for a,b in pairs),'selected_outside_baseline_best':sum(not set(b['Best']).issubset(a['Best'])for a,b in pairs)})
if (r/'baseline-deep.json').exists():
 db=load('baseline-deep.json')
 for label in ['baseline','safe-cache','unsafe-board-scores','root-single-best','tactical-ordering','coordinated-kernel']:
  p=r/(label+'-deep.json')
  if not p.exists():continue
  data=load(p.name)
  for depth in [5,6]:
   pairs=[(a,b)for a,b in zip(db,data)if a['Depth']==depth];out['deep'].append({'variant':label,'depth':depth,'completed_legal':sum(b['Complete']and b['Legal']for a,b in pairs),'nodes':sum(b['Nodes']for a,b in pairs),'node_ratio':sum(b['Nodes']for a,b in pairs)/sum(a['Nodes']for a,b in pairs),'score_differences':sum(a['Score']!=b['Score']for a,b in pairs),'best_set_differences':sum(set(a['Best'])!=set(b['Best'])for a,b in pairs),'selected_outside_baseline_best':sum(not set(b['Best']).issubset(a['Best'])for a,b in pairs)})
for depth in [3,4]:
 rows=[x for x in load('cache-audit.json')if x['depth']==depth];a=collections.Counter()
 for x in rows:a.update(x['audit'])
 out['audit'][str(depth)]=dict(a,unique_board_contexts=sum(x['unique_board_contexts']for x in rows),unique_history_contexts=sum(x['unique_history_contexts']for x in rows))
sf={(x['position'],x['move']):x['last_exact']['score']for x in load('referee-results.json')['results']if x['last_exact']}
for filename in ['uci-100ms.json','uci-1000ms.json','uci-warm.json','uci-historical.json','uci-historical-fixed.json','uci-combo.json']:
 p=r/filename
 if not p.exists():continue
 data=load(filename);by=collections.defaultdict(list)
 for x in data['records']:by[x['variant'],x['group']].append(x)
 for (label,group),rows in by.items():
  nonterminal=[x for x in rows if x.get('depth')is not None and x.get('legal')];elapsed=[x['response_ms']for x in rows if'response_ms'in x]
  out['timed'].append({'file':filename,'variant':label,'group':group,'requests':len(rows),'completed_legal':sum(x.get('completed',False)for x in rows),'mean_depth':statistics.mean(x['depth']for x in nonterminal)if nonterminal else None,'mean_reported_nodes':statistics.mean(x['nodes']for x in nonterminal if x.get('nodes')is not None)if nonterminal else None,'mean_response_ms':statistics.mean(elapsed)if elapsed else None,'max_response_ms':max(elapsed)if elapsed else None})
  cases=collections.defaultdict(list);mates=collections.Counter();negative=0
  for x in rows:
   if not x.get('legal'):continue
   origin=sf.get((x['position'],None));chosen=sf.get((x['position'],x['move']))
   if not origin or not chosen:continue
   if 'cp'in origin and 'cp'in chosen:
    delta=origin['cp']-chosen['cp'];negative+=int(delta<0);cases[x['id']].append({'loss':max(0,delta),'raw':delta,'move':x['move'],'ref':origin['cp'],'chosen':chosen['cp']})
   else:
    if chosen.get('mate',0)<0 and not origin.get('mate',0)<0:mates['allowed_losing_mate']+=1
    elif origin.get('mate',0)>0 and not chosen.get('mate',0)>0:mates['missed_winning_mate']+=1
    elif 'mate'in origin and 'mate'in chosen:mates['both_mate']+=1
    else:mates['other_mate_transition']+=1
  case_rows=[{'id':k,'mean_clamped_cp_loss':statistics.mean(v['loss']for v in vals),'mean_raw_cp_loss':statistics.mean(v['raw']for v in vals),'observations':len(vals),'choices':vals}for k,vals in cases.items()]
  out['quality'].append({'file':filename,'variant':label,'group':group,'cp_cases':len(case_rows),'mean_case_cp_loss':statistics.mean(x['mean_clamped_cp_loss']for x in case_rows)if case_rows else None,'median_case_cp_loss':statistics.median(x['mean_clamped_cp_loss']for x in case_rows)if case_rows else None,'case_mean_loss_at_least_200':sum(x['mean_clamped_cp_loss']>=200 for x in case_rows),'negative_comparisons':negative,'mate_observations':dict(mates),'cases':case_rows})
 if filename=='uci-warm.json':
  for label in data['identities']:
   warm={x['id'].split('/')[0]:x for x in data['records']if x['variant']==label and x['id'].endswith('/warm-next')and x.get('completed')}
   cold={x['id'].split('/')[0]:x for x in data['records']if x['variant']==label and x['id'].endswith('/cold-next')and x.get('completed')}
   pairs=[(warm[k],cold[k])for k in warm if k in cold];out['warm_native'].append({'variant':label,'nonterminal_pairs':len(pairs),'warm_depth_higher_equal_lower':dict(collections.Counter('higher'if a['depth']>b['depth']else'lower'if a['depth']<b['depth']else'equal'for a,b in pairs)),'warm_mean_depth':statistics.mean(a['depth']for a,b in pairs),'cold_mean_depth':statistics.mean(b['depth']for a,b in pairs)})
(r/'analysis.json').write_text(json.dumps(out,indent=2)+'\n')
print('Saved analysis.json')
