## Publication note for v0.19

I preserve the original local experiment report below, including its stop/check-in
state. I subsequently received approval to publish these attempts together in
v0.19. I keep playing policy at v0.18; the [release summary](/docs/RESEARCH_V019.md)
and [Go tooling guide](/tools/research/README.md) describe current status and
replace historical script commands. Original source remains in the hashed
release attachment.

## v0.19 tooling migration

I preserve this historical report and its original measurements. I retired the
Python/shell runners in v0.19; commands below describe the original experiment.
I now use the [Go research tools](/tools/research/README.md). Original scripts
and unedited evidence remain in the [hashed release archive](https://github.com/philyawj/Pandafish/releases/download/v0.19/v0.19-historical-research-evidence.tar.gz).

# Local handoff after the third post-v0.18 fit

I am on `v0.19-decision-evaluation`, based on released v0.18
`75f12bbdbeb0e38d148690ddc0f1d6405c5af3d3`. I have no new commits. The playing
Go source and module/local runtime files match the tag. All three post-release
research attempts remain local. I must check in before commit, push or PR;
I have no authorization for a new release or live Lichess games.

I completed the authorized sharp-branch/both-turn retry. The pre-fit algebraic
objective correction is frozen separately from the original plan. The selected
model fails fresh holdout before activation: penalty 44.0435 to 49.9918,
severe roots 3 to 6, mate-risk roots 3 to 4. I ran no actual candidate engine,
equal-time performance test or strength games for this model. I will not refit
it or reuse this holdout as unseen evidence.

I recommend rejecting the weights, keeping the data/objective repairs and
testing a separate bounded quiet-mate proof next. The training-only full scan
finds 15 static records with legal quiet mates; unchanged v0.18 Go and
python-chess independently verify them. This is a functional omission and
prospective hypothesis, not measured strength. See [NEXT.md](NEXT.md).

## My preserved evidence

- `PLAN.json`, `FIT_REFINEMENT.json`, original objective/source and
  `ERRATA.md`: exact chronology, controls and pre-result correction.
- `generation*.json*`/logs: all attempts, recorded prefix, failure/resume and
  80 completed rollout families. `blocked.json`: 140,130 states/16 inputs.
- `forest-leaves.jsonl`/`.json`: every legal alternative before branch labels;
  29,357 records. Feature oracle: 28,890 vectors/scalars, zero failures.
- Membership/audits: 78 components, raw split 504/186/186, zero overlap.
- `corpus.json`, `label-queries.jsonl`, `leaves.json`, `fit-data.json`:
  complete sharp/terminal/mate alternatives, no root exclusions.
- `model-weights.json`, `fit-results.json`, `fit-run.txt`: one frozen model,
  every validation grid fit and failed holdout. No candidate binary exists.
- `capacity-diagnosis.json`: initial exploratory output retained.
  `diagnose_capacity.py` / `capacity-diagnosis-robust.json`: exact rational
  continuous bounds plus a conservative 2cp integer margin; four dominating
  bad branches span two distinct training roots, not four roots.
- `quiet-threat-diagnosis.json`: initial six-branch probe, no mate found.
  `diagnose_quiet_threats.py` / `training-quiet-threat-scan.json`: all 2,557
  selected learnable training witnesses, 15 actual quiet-mate records and 54
  hypothetical pass threats. `quiet-proof.txt`: independent tagged-Go proof.
- `final_audit.py`, `final-audit.json`, `provenance.json`, `inventory.json`:
  source, decisions, prefix and previous-archive integrity. Python caches are
  disposable and excluded.

## My reproduction commands

I use the existing chess-client virtual environment for python-chess, and
system `python3` for NumPy. Both are Python 3.12.4. I installed nothing.

```sh
cd /Users/philyawj/projects/Pandafish/tools/benchmark/research/v0.19-sharp-branches
/Users/philyawj/projects/lichess-bot/venv/bin/python -m unittest -v test_boundary test_query_controls
OPENBLAS_NUM_THREADS=1 python3 -m unittest -v test_model
python3 final_audit.py
```

The audit does not refit, contact engines or change frozen observations. The
diagnostic scripts regenerate only their derived capacity/training-scan files;
I would run them in a copied archive to preserve the recorded artifacts.

For the independent Go mate proof I would copy this directory to scratch,
remove only its derived `quiet-mate-go-checks.json`, and run
`python3 prove_quiet_mates.py`. It creates a tagged-v0.18 scratch tree, adds
only the archived test, uses `GOCACHE=/tmp/pandafish-v016-go-cache`, and preserves
the format/test exits and proof log. It makes no playing-source changes.

`fit.py` refuses to overwrite frozen weights. An exact fit reproduction would
use a fresh copied archive with only derived weights/results removed and the
recorded data unchanged; it would not create independent confirmation. I do
not rerun live collection or pick new coefficients after inspecting holdout.
