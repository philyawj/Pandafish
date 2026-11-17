"""Bounded joint calibration; freeze weights before unlocking holdout metrics."""
import hashlib,json
from pathlib import Path
import numpy as np
r=Path(__file__).resolve().parent
assert not(r/'model-weights.json').exists(),'Refuse to retune the frozen candidate'
assert not json.loads((r/'feature-oracle-results.json').read_text())['failures']
plan=json.loads((r/'CALIBRATION_PLAN.json').read_text())
cases=json.loads((r/'corpus.json').read_text())['cases'];features={x['ID']:x for x in json.loads((r/'corpus-features.json').read_text())}
rows=[];excluded=[];seen={}
for c in cases:
 f=features[c['id']];reason=[];label=c['label']['cp'];key=' '.join(c['fen'].split()[:4])
 if key in seen:reason.append('duplicate-of:'+seen[key])
 else:seen[key]=c['id']
 if label is None:reason.append('mate-label')
 if c['checked']:reason.append('checked')
 if c['claimable']:reason.append('claimable-draw')
 if c['halfmove']>60:reason.append('high-halfmove')
 if abs(f['Base'])>1500 or(label is not None and abs(label)>1500):reason.append('outside-cp-range')
 if abs(f['Quiet']-f['Base'])>50:reason.append('unresolved-tactics')
 if reason:excluded.append({'id':c['id'],'reasons':reason});continue
 phase=f['Features']['Phase']/24;v=np.array(f['Features']['Values'],dtype=float)
 rows.append({'id':c['id'],'family':c['family'],'split':c['split'],'base':f['Base'],'teacher':label,'x':np.r_[v*phase,v*(1-phase)].tolist()})
(r/'fit-rows.json').write_text(json.dumps({'eligible':rows,'excluded':excluded},indent=2)+'\n')
X=np.array([x['x']for x in rows]);base=np.array([x['base']for x in rows]);teacher=np.array([x['teacher']for x in rows]);y=np.clip(teacher-base,-400,400)
indices={name:np.array([i for i,x in enumerate(rows)if x['split']==name],dtype=int)for name in ['train','validation','holdout']}
assert all(len(v)>=15 for v in indices.values()),{k:len(v)for k,v in indices.items()}
bounds=np.array(plan['middle_bounds']+plan['end_bounds']);train=indices['train'];scale=np.maximum(1,np.sqrt(np.mean(X[train]**2,axis=0)));Z=X[train]/scale
G=Z.T@Z;rhs=Z.T@y[train]
def metrics(weights,index):
 predicted=base[index]+np.trunc(X[index]@weights)
 error=np.abs(np.clip(predicted,-1000,1000)-np.clip(teacher[index],-1000,1000))
 return {'rows':len(index),'clipped_mae_cp':float(np.mean(error)),'raw_catastrophes_over_300cp':int(np.sum(np.abs(predicted-teacher[index])>300))}
fits=[]
for regularization in plan['regularization']:
 gram=G+len(train)*regularization*np.eye(len(bounds));beta=np.zeros(len(bounds));iteration=0
 for iteration in range(1000):
  before=beta.copy()
  for j in range(len(beta)):
   beta[j]=np.clip((rhs[j]-gram[j]@beta+gram[j,j]*beta[j])/gram[j,j],0,bounds[j]*scale[j])
  if np.max(np.abs(beta-before))<1e-7:break
 weights=np.clip(np.rint(beta/scale),0,bounds).astype(int)
 fits.append({'regularization':regularization,'iterations':iteration+1,'weights':weights.tolist(),'train':metrics(weights,train),'validation':metrics(weights,indices['validation'])})
best=min(fits,key=lambda f:f['validation']['clipped_mae_cp'])
frozen={'base_tag':'v0.17','feature_count':12,'feature_names':['knight-safe-mobility','bishop-safe-mobility','rook-safe-mobility','queen-safe-mobility','king-zone-pressure','coordinated-attackers-squared','pawn-shelter','passed-pawn-progress','passed-pawn-king-access','negative-isolated-pawns','negative-doubled-pawns','negative-blocked-pawns'],'selected_on':'validation only, after integer rounding','regularization':best['regularization'],'middle':best['weights'][:12],'ending':best['weights'][12:],'source_corpus_sha256':hashlib.sha256((r/'corpus.json').read_bytes()).hexdigest(),'calibration_plan_sha256':hashlib.sha256((r/'CALIBRATION_PLAN.json').read_bytes()).hexdigest()}
(r/'model-weights.json').write_text(json.dumps(frozen,indent=2)+'\n')
# This is the first evaluation of holdout quality. Do not retune from it.
w=np.array(best['weights']);zero=np.zeros(len(bounds));quality={name:{'baseline':metrics(zero,index),'candidate':metrics(w,index)}for name,index in indices.items()}
for value in quality.values():value['relative_mae_improvement']=1-value['candidate']['clipped_mae_cp']/value['baseline']['clipped_mae_cp']
h=quality['holdout'];passed=h['relative_mae_improvement']>=.05 and h['candidate']['raw_catastrophes_over_300cp']<=1.1*h['baseline']['raw_catastrophes_over_300cp']
report={'plan':plan,'numpy_version':np.__version__,'raw_rows':len(cases),'eligible_rows':len(rows),'excluded_rows':len(excluded),'validation_grid':fits,'selected':frozen,'quality':quality,'holdout_gate_passed':bool(passed),'notes':'Finite-node teacher, public correlated source games, few held-out opening families. Static label quality is a prerequisite, not a playing-strength claim. Weights were frozen before holdout metrics; no holdout retuning.'}
(r/'calibration-results.json').write_text(json.dumps(report,indent=2)+'\n')
print(json.dumps({'rows':{k:len(v)for k,v in indices.items()},'weights':frozen,'quality':quality,'holdout_gate_passed':bool(passed)},indent=2))
