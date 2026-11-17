"""Profile tagged v0.16 before choosing the next implementation."""
import hashlib,io,json,os,subprocess,tarfile,tempfile,time
from pathlib import Path
r=Path(__file__).resolve().parent;repo=r.parents[3];out=Path(tempfile.mkdtemp(prefix='pandafish-v017-profile-'))
archive=subprocess.check_output(['git','archive','v0.16'],cwd=repo)
with tarfile.open(fileobj=io.BytesIO(archive))as t:t.extractall(out,filter='data')
(out/'research_profile_test.go').write_bytes((r.parent/'v0.16-search-core/driver_test.go.txt').read_bytes())
env=dict(os.environ,GOCACHE='/tmp/pandafish-v016-go-cache',ARCH_INPUT=str(r.parent/'v0.16-search-core/inputs.json'),ARCH_NO_MEM_PROFILE='1')
cmd=['go','-C',str(out),'test','-run','^$','-bench','^BenchmarkArchitectureProfile$','-benchtime','60x','-cpuprofile',str(r/'baseline-cpu.pb.gz'),'-o',str(out/'profile.test')]
start=time.monotonic();p=subprocess.run(cmd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
(r/'baseline-profile.txt').write_text(p.stdout)
(r/'profile-command.json').write_text(json.dumps({'base_tag':'v0.16','base_commit':subprocess.check_output(['git','rev-parse','v0.16^{commit}'],cwd=repo,text=True).strip(),'command':cmd,'exit':p.returncode,'elapsed_s':time.monotonic()-start,'tree':str(out),'cpu_allocation_stack_sampling_disabled':True,'collector':'same 24 neutral-position 100ms kernel as v0.15.3 architecture report; no other engine workloads run concurrently'},indent=2)+'\n')
assert p.returncode==0,p.stdout
for label,extra in [('flat',[]),('cumulative',['-cum'])]:
 cmd=['go','tool','pprof','-top',*extra,str(out/'profile.test'),str(r/'baseline-cpu.pb.gz')]
 p=subprocess.run(cmd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);(r/('baseline-cpu-'+label+'.txt')).write_text(p.stdout);assert p.returncode==0
print((r/'baseline-profile.txt').read_text());print((r/'baseline-cpu-cumulative.txt').read_text()[:5500])
