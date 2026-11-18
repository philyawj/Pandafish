# Reassess implementation cost after the larger-gain study

Seven implemented variants fail their applicable larger-gain objectives. Guarded
LMR, strict/certified quiescence caches and the stand-pat opportunity diagnostics
also fail their predefined large-work criteria before a playing implementation.
None supplies new native or game-strength evidence. Preserve all failures rather
than claim the 15% objective passed or keep expanding speculative complexity.

The monotonic-deadline variant remains a small **functional** improvement:
three cold ratios 0.93045/0.92912/0.93135, deep 0.92772 and newer 0.89227, with all
325 scores/moves/nodes exact and passing deadline/race/vet/build controls. It
fails the original larger 15% goal. However, its short deadline helper uses a
public Go API and preserve every time allowance/check site. That cost/risk is
different from adding per-ply mutable storage or a new cache/proof subsystem for
a similarly small gain. The 15% research target was introduced on this branch;
it is not a repository requirement or a user-specified minimum for an otherwise
worthwhile candidate.

Open a separate, prospective **small-scope release-candidate assessment** of the
exact same frozen clock code. Do not rerun or change the original cold results,
source, thresholds or failed field. Run three previously unexecuted warm pairs
with depth four plus the frozen 30s allowance, so they exercise production
deadline reads without reaching a time boundary. Require at least **5% median
search elapsed reduction**, all complete score/move/node identities and no
contradictory existing cohort. Freeze that target before warm observations.

Then require the same four actual-budget native safety/quality gates and one
standard 50-games-per-opponent concurrent screen versus tagged v0.21, Skill 3
and pinned Sacrifice. These have never run for this code; no failed match/native
observation is rerun or discarded. Do not tune the runtime code to those results.
If this smaller assessment fails, keep it inactive and resume a different
hypothesis. If it passes, present the modest functional gain and actual game
screen separately, with the entire unsuccessful larger-gain study included.
No precise Elo, broad strength claim, larger confirmation or publication follows
automatically. Check in before committing a viable candidate.

Larger research remains open. Current evidence makes a low-cost checking-move/
quiet-threat search interface or a genuinely different evaluation representation
more substantive future projects than repeatedly adding small ordering/cache
layers. Neither is implemented or claimed successful here. No opening book or
online dataset is needed to validate this clock change.
