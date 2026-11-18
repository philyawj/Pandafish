# v0.20 versus Stockfish Skills 2–10

50 games per skill; 450 planned. Tagged v0.20; 100ms move allowance, 50ms protocol grace, 400-ply cap, ten paired openings. Three levels per batch (2/5/8, 3/6/9, 4/7/10); individual profiles record one matchup and method.json records the outer concurrency. All scores are from Pandafish's perspective.

| Skill | W / D / L | Unfinished | Forfeits | Completed score | Possible score | Records |
| --- | --- | --- | --- | --- | --- | --- |
| 2 | 29 / 4 / 17 | 0 | 0 | 62.0% | 62.0%–62.0% | [JSON](runs/skill-2/v0.20-vs-stockfish-2/summary.json), [PGN](runs/skill-2/v0.20-vs-stockfish-2/stockfish.pgn) |
| 3 | 22 / 4 / 24 | 0 | 0 | 48.0% | 48.0%–48.0% | [JSON](runs/skill-3/v0.20-vs-stockfish-3/summary.json), [PGN](runs/skill-3/v0.20-vs-stockfish-3/stockfish.pgn) |
| 4 | 16 / 0 / 34 | 0 | 0 | 32.0% | 32.0%–32.0% | [JSON](runs/skill-4/v0.20-vs-stockfish-4/summary.json), [PGN](runs/skill-4/v0.20-vs-stockfish-4/stockfish.pgn) |
| 5 | 14 / 6 / 30 | 0 | 0 | 34.0% | 34.0%–34.0% | [JSON](runs/skill-5/v0.20-vs-stockfish-5/summary.json), [PGN](runs/skill-5/v0.20-vs-stockfish-5/stockfish.pgn) |
| 6 | 9 / 0 / 41 | 0 | 1 | 18.0% | 18.0%–18.0% | [JSON](runs/skill-6/v0.20-vs-stockfish-6/summary.json), [PGN](runs/skill-6/v0.20-vs-stockfish-6/stockfish.pgn) |
| 7 | 4 / 2 / 43 | 1 | 0 | 10.2% | 10.0%–12.0% | [JSON](runs/skill-7/v0.20-vs-stockfish-7/summary.json), [PGN](runs/skill-7/v0.20-vs-stockfish-7/stockfish.pgn) |
| 8 | 0 / 2 / 48 | 0 | 0 | 2.0% | 2.0%–2.0% | [JSON](runs/skill-8/v0.20-vs-stockfish-8/summary.json), [PGN](runs/skill-8/v0.20-vs-stockfish-8/stockfish.pgn) |
| 9 | 1 / 0 / 49 | 0 | 2 | 2.0% | 2.0%–2.0% | [JSON](runs/skill-9/v0.20-vs-stockfish-9/summary.json), [PGN](runs/skill-9/v0.20-vs-stockfish-9/stockfish.pgn) |
| 10 | 0 / 1 / 49 | 0 | 0 | 1.0% | 1.0%–1.0% | [JSON](runs/skill-10/v0.20-vs-stockfish-10/summary.json), [PGN](runs/skill-10/v0.20-vs-stockfish-10/stockfish.pgn) |

The possible-score range treats unfinished games as all losses through all wins; it is not a confidence interval. Small repeated-opening screens and batch scheduling limit precision; adjacent skill samples need not be monotonic. No Elo equivalence, automatic target change or longer confirmation is implied. All raw runs, driver, method, failures and PGN checks are preserved. See [the frozen plan](PLAN.md).

**Recommend Skill 3 for the next routine reference**, pending review. See the
[assessment](ASSESSMENT.md) for the rationale, timeouts, sensitivity bounds and
limits. There are three timeout forfeits: Stockfish once at Skill 6; Stockfish
and Pandafish once each at Skill 9. Skill 7 has one unfinished 400-ply game.
All official scores remain unchanged. All 450 saved games replay legally;
data integrity passes while the reliability gate remains failed. The
[reporting findings](REPORTING_NOTES.md) preserve the collector's final assertion
failure and its separate report-only correction. No games were rerun.
