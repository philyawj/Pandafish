# v0.13 selective-quiescence candidate

The maintainer reviewed and approved this candidate for v0.13. Original
development labels and binary hashes are preserved; these games tested
v0.13-dev, rather than the subsequently tagged release binary.
The selected experiment follows the [post-v0.12 comparison and three recommendations](../../research/v0.13-candidates/README.md).

## Change and limits

Continue ordinary capture/promotion quiescence for up to 16 plies instead of
two. Keep the existing four extra checked-node evasion plies and final static
safety cutoff. Outside check, use static exchange evaluation (SEE) to skip
negative ordinary non-checking captures. Checking captures, promotions, and
required evasions remain eligible. Equal/cheaper attackers bypass SEE.

SEE follows legal least-valuable recaptures on the destination, accounting for
pins, exposed rays, king safety, en passant, and promotions. It is a material
heuristic: equal-valued attacker alternatives and compensation elsewhere can
matter. Skipping a sacrifice can miss a tactic. This is selective search, not
full-width minimax equivalence. The safety bounds can still stop unresolved
sequences. Evaluation, main search, TT rules, killers, clocks, and exact ties
among the searched root scores are unchanged. No new UCI option or dependency.

## Reproduction

Apple M2 MacBook Air; all games and matchups serial; no live bot. The two
[power observations](power-observations.json) reported AC power at start/end,
but the battery fell from 99% to 79% and reported discharging. These are not
continuous CPU, thermal, or power measurements. Scheduling and operating
conditions limit comparisons with earlier screens.

Use the unchanged ten-opening screening suite, 50 color pairs per opponent,
100 ms requested per move, 50 ms protocol allowance, and 400-ply cap. First
play v0.12 and require zero forfeits and a score strictly above 50% counting
unfinished games as losses. Then play Skill 1 and pinned Sacrifice serially.
No retries, longer runs, or Skill 2 extension were used.

```sh
caffeinate -i go run ./tools/benchmark --validate --race --version v0.13-dev \
  --previous tools/benchmark/bin/pandafish-v0.12 --previous-version v0.12 \
  --previous-revision a7f6101819697cc0e756ce233a8468c9d41183e7 \
  --sacrifice tools/benchmark/bin/sacrifice-source/target/release/Sacrifice-Engine \
  --revision a7f6101819697cc0e756ce233a8468c9d41183e7-dirty \
  --machine 'Apple M2 MacBook Air; AC power reported at start; all matchups serial; no live bot' \
  --verbose
```

The [exact candidate source snapshot](candidate-source.json) is relative to
tag v0.12 at the recorded base revision. The match and selected-position audit
used the identical candidate SHA-256:
`394243a8bcb0368f5288363785f8fff9226384414a07bcbf25ec79591b20d63f`.
[Summary JSON](summary.json) records opponent/runner hashes and checks. The
Sacrifice pin is v0.2.0, b528f0ba373910b72052d991700a990ca1cb07fd.

After matches, a [comment-only clarification](post-match-comments.patch)
qualified the equal/cheaper-attacker shortcut for recapture promotions and
specified checking captures. Apply the zero-context patch with `git apply --unidiff-zero`. It changes no
executable search statements or
parameters. The candidate wrapper build before release preparation has SHA-256
`be5f70adefeea8965f3286ce52bae0e2b5c5092c8f627c90105afdf4ba7e33fa`.
No algorithm was retuned after inspecting any diagnostic or match results.

## Verification and diagnostic

Formatting, `go test -race ./...`, `go vet ./...`, and the candidate build passed
before matches. Ordinary `go test ./...`, vet, and whitespace checks also passed
after the comment clarification. The ARM64 build and offline lichess-bot
config/UCI wrapper returned 12 legal plies with identity Pandafish v0.13-dev.
The active v0.12 executable was restored byte-for-byte after each wrapper check.
Neither private bot config changed; the public example's development greetings
were updated and validated with temporary resolved engine paths and a placeholder
token. No live listener, challenges, or online games were started.

[The selected-position audit](../../research/v0.13-candidates/selective-quiescence-audit/README.md)
preserves focused test details, the independent legal-material oracle,
unfiltered continuation ablation, all 78 legal probes, Stockfish referee
outputs, source, hashes, and compact aggregates. The known exchange fixture is
fixed by SEE even at two plies; it does not isolate a benefit from the new
16-ply bound. Fixed-time choices were mixed. Against broader unfiltered
continuation, filtering reduced depth-three nodes 27.7% on the selected corpus;
against v0.12, fewer nodes did not mean less elapsed time. Those diagnostics
are distinct from these game results.

## Results

| Opponent | W/D/L | Unfinished | Forfeits | Completed score |
| --- | --- | --- | --- | --- |
| v0.12 | 46/43/11 | 0 | 0 | 67.5% |
| Stockfish Skill 1 | 58/11/31 | 0 | 0 | 63.5% |
| Sacrifice v0.2.0 b528f0ba | 24/28/48 | 0 | 0 | 38% |

All 300 games completed. The previous-release gate passed. Each forfeits cell
counts either side; there were none. Draws score half a point.

| Matchup | Candidate responses | Above requested 100 ms | Maximum response ms |
| --- | --- | --- | --- |
| v0.12 | 6,012 | 6 | 136.474 |
| Stockfish Skill 1 | 4,812 | 7 | 129.334 |
| Sacrifice | 4,674 | 8 | 131.088 |

All 15,498 candidate responses stayed inside the 150-ms deadline including
protocol allowance. Twenty-one exceeded the requested 100 ms. The maximum
is an observed response, not an engine time guarantee or a measured cause.
Opponents also exceeded the request at times; complete timings remain in JSON.

## Interpretation and recommendation

This is a promising local improvement. The 67.5% direct score against v0.12 is
more substantial than the last two previous-release screens. Relative to the
separate v0.12 development screen, Skill 1 rises from 57.5% to 63.5% (+6 points),
and Sacrifice from 17.5% to 38% (+20.5 points).

Before reviewing match scores, we described a visible external gain as at least
10 points on one reference without a decline greater than 5 points on the
other. This screen meets that descriptive target. It is not a significance
test: these are separate 100-game runs with repeated openings, random ties,
and varying operating conditions. Do not infer a precise Elo gain, guaranteed
long-control benefit, or reliable Skill 1 win rate above 50% from this sample.
Sacrifice still wins substantially more than it loses against the candidate.

The maintainer approved proceeding to PR, merge, and v0.13 release. The intended
mechanism has regression coverage, the routine gate passed, reliability checks
passed, and both external scores improved. This supports a stronger claim than
efficiency alone, while retaining the limits above. Keep Stockfish Skill 1 as
the routine target; any ladder change or longer confirmation is a separate
decision. Selective main search and richer classical evaluation remain the
next two research options, not part of this increment.

## Saved records

- [Compact results](compact-results.json) and [complete summary/timings](summary.json)
- [All games](games.pgn), [v0.12](previous.pgn), [Stockfish](stockfish.pgn), [Sacrifice](sacrifice.pgn)
- [Exact match source](candidate-source.json) and [post-match comment patch](post-match-comments.patch)
- [Position audit, sources, and raw records](../../research/v0.13-candidates/selective-quiescence-audit/README.md)


## Release preparation

The maintainer approved commit, push, PR, merge, and release as v0.13. UCI
identity, its test, greetings, current README, and future-baseline instructions
were updated. The benchmarked search policy and original run labels/hashes
remain unchanged. Format checks, full race-enabled tests, vet, and build passed
again. Two-game integration smokes against each opponent completed with no
forfeits or unfinished games; those small samples validate integration rather
than strength and remain under `/tmp/pandafish-v013-release-smoke`.
The release ARM64 build passed 12 offline lichess-bot wrapper plies. Both private
config copies received version-only greeting updates and passed config loading;
credentials and challenge policies were preserved. No live games were started.
