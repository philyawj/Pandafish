# Are 50 games per opponent enough?

**Keep 50 games against each of three opponents as fast triage. Do not treat a
small lead in that screen as proof of improvement.** For the smaller gains
now being attempted, I recommend a separately planned **200-game confirmation
per opponent** for a promising strength candidate, still with three matchups
concurrent. Even 200 may be inconclusive for a 20 to 50 Elo change.

This is a recommendation for the next workflow decision. The program's default
remains 50 games per opponent, three matchups concurrent, as previously agreed.
No large new match or automatic escalation ran for this research release.
Correctness or tooling changes can still ship on appropriate functional
evidence without claiming a measured strength gain.

## Evidence from the last candidate

The original v0.15 candidate scored 20 wins, 13 draws and 17 losses against
v0.14.4: **53% in 50 games**, or 1.5 points more than half. Converting the point
score through the conventional Elo formula gives about +21 Elo, but that
point conversion alone says very little about the true difference.

The saved [sampling diagnostic](analysis/benchmark-sampling.json) resamples
the actual original results 100,000 times with a recorded deterministic seed:

| Method | Approximate 95% score interval |
| --- | --- |
| Independent-game normal approximation | 41.0% to 65.0% |
| Resample the 25 complete color pairs | 40.0% to 66.0% |
| Resample the ten observed opening clusters | 39.3% to 68.0% |

The paired interval corresponds to roughly -70 to +115 conventional Elo.
These are descriptive, conditional-on-the-observed-book intervals, not a
certified rating interval or SPRT verdict. There are only ten opening clusters;
the calculation does not model unseen openings, engine changes selected after
many attempts, or host/workload shifts. It illustrates why a 53% screen is
inconclusive rather than demonstrating that the feature is bad.

Fifty games against each of three different opponents are not 150 observations
of the same head-to-head strength parameter. Each opponent supplies a separate
comparison with its own draw rate and style. Look for consistency, but do not
pool all three into a single improvement percentage.

## How much more evidence?

For scale only, an independent-game normal model with the observed 26% draw
rate, two-sided 5% error and 80% power gives these approximate requirements:

| True conventional Elo difference | Expected score | Illustrative games |
| --- | --- | --- |
| +10 | 51.44% | 7,016 |
| +25 | 53.59% | 1,126 |
| +50 | 57.15% | 285 |
| +100 | 64.01% | 75 |

The calculation uses `p = 1 / (1 + 10^(-Elo/400))` and
`n = (1.96 + 0.842)^2 * (1 - drawRate) / (4 * (p - 0.5)^2)`.
It is an illustrative approximation, not a guaranteed stopping rule. Color-pair
covariance, book repetition and multiple trials change the requirements.
More games do not fix a biased book or inaccurate adjudication.

The [Fishtest mathematics documentation](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-Mathematics.html)
describes paired/pentanomial modeling and sequential tests. Its
[FAQ](https://official-stockfish.github.io/docs/fishtest-wiki/Fishtest-FAQ.html)
describes staged testing and warns about repeated attempts. Our fixed-size
runner does not implement GSPRT or give those statistical guarantees.

## Practical proposal

1. First require focused correctness and mechanism tests. For piece safety,
   show which quiet leaves change and why, and measure cost at equal time.
2. Run the existing 50-game routine against the previous release, Stockfish
   Skill 2 and pinned Sacrifice, three matchups concurrent. Retain failures,
   unfinished games, forfeits, clocks and hashes. A narrow lead is a reason to
   investigate, not a strength verdict.
3. For a promising strength candidate, freeze its source and choose the next
   sample before running. Use 200 fresh games per opponent with the broader
   50-position confirmation book. Do not overwrite or relabel the screen.
   Freeze baseline, opponent binaries, budgets and concurrency too.
4. For a close result, report it as inconclusive. A preplanned roughly 300-game
   test is more appropriate for a hoped-for 50 Elo difference; around 1,000 or
   more may be needed for smaller differences. Do not keep adding games until
   a desired score appears, or require all three noisy opponent percentages to
   rise. A working rule fix and an Elo claim are different release rationales.

Example confirmation after the default screen, if this proposed stage is chosen:

```sh
go run ./tools/benchmark \
  --candidate tools/benchmark/bin/pandafish-candidate \
  --version v0.16-piece-safety-confirmation-dev \
  --previous tools/benchmark/bin/pandafish-v0.15.1 \
  --previous-version v0.15.1 \
  --stockfish /absolute/path/to/stockfish --stockfish-skill 2 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --pairs 100 --parallel-all \
  --openings tools/benchmark/openings-confirmation-50-v1.txt
```

This is 100 color pairs, 200 games per opponent. The existing 50-opening book
is cycled twice, so this is broader but still contains repeated openings.
It is a public confirmation book, not a hidden holdout. The screen and
confirmation must remain separately labeled; do not selectively pool them.
See the [benchmark setup](../../README.md) for pinned opponent installation.

Keep three-way concurrency for development on the demonstrated host capacity.
It saves time and reflects the recorded calibration decision. Use serial
execution for a planned sensitivity check when investigating small gains or
host contention, not automatically. Parallel play is not free of variance,
but changing the workload between candidates can be worse than holding it
consistent. Live games provide failure discovery; controlled local matches
remain the strength-validation loop.
