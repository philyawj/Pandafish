# Combined v0.22 Lichess review

The three completed recordings contain **157 distinct games against 119 BOT
accounts: 61 wins, 11 draws, 80 losses and five unfinished games**. All 24
selected substantial CP mistakes persist at a fivefold referee budget. King
exposure, development and piece activity deserve the next decision-level
diagnosis; the previous mate-proof-first ranking should be revised.

This is an offline evidence review of released v0.22. No engine behavior changed,
new candidate was tested, live game was started, or commit, PR or release was
prepared. Further work requires a separate decision.

## Recordings and workloads

Every source folder remains intact. Session names use local Pacific time;
manifest intervals below use UTC. The preparation-only
`local-games/v0.22-setup-check-20261118` is excluded.

| Original session path | Preserved archive | UTC interval | Recorded workload | W / D / L / unfinished |
| --- | --- | --- | --- | --- |
| `local-games/v0.22-live-20261118-080529` | [20261118-080529](../v0.22-live-20261118-080529/OPERATIONS.md) | Nov 18 16:05:29–20:00:00 | M2 MacBook Air, terminal | 5 / 2 / 12 / 0 |
| `local-games/v0.22-live-20261118-120005` | [20261118-120005](../v0.22-live-20261118-120005/OPERATIONS.md) | Nov 18 20:00:05–Nov 19 23:04:36 | M2 MacBook Air, terminal | 42 / 6 / 52 / 5 |
| `local-games/v0.22-live-20261119-151427` | [20261119-151427](../v0.22-live-20261119-151427/OPERATIONS.md) | Nov 19 23:14:27–Nov 20 05:19:19 | M2 MacBook Air, terminal | 14 / 3 / 16 / 0 |

All manifests identify clean source revision
`1953f2af97f956c4499cbc5b61c3c1c698856740`, UCI `Pandafish v0.22`, and the same
captured executable SHA-256
`b1f7f1067c1cacbe6fb7b88264286c59fe150d123cdd603c7f6909ead1eb813d`.
No other Pandafish version is pooled. Opponent executable versions are not
recorded; account names are identifiers, not version guarantees.

The recorded workload is a declaration, not measured system utilization.
No other foreground workload was supplied for this review. The reported laptop
closure is documented below. The redacted configurations differ only in engine
and PGN recording paths; no intentional strength or time-setting change appears.
Incoming clocks varied because the unchanged configuration accepts several
time controls. Matchmaking requested rated standard 180+2, concurrency one,
no pondering, no book or online move sources, no automatic draw offers or
resignation, and 1,000 ms client move overhead. The external client revision is
the same in all sessions, but its working tree was dirty and its exact patch
was not captured.

## Results separated by clock and opponent type

All 157 PGNs label the opponent `BOT`; there is no human-opponent sample.
Each row below is therefore v0.22 versus BOT accounts at the stated clock.
Score is `(wins + draws/2) / finished games`. Unfinished games remain separate.

| Clock, seconds + increment | W | D | L | Unfinished | Finished score |
| --- | --- | --- | --- | --- | --- |
| 60+1 | 2 | 1 | 2 | 0 | 50.00% |
| 60+2 | 0 | 0 | 0 | 1 | — |
| 60+5 | 1 | 0 | 0 | 0 | 100.00% |
| 180+0 | 2 | 0 | 0 | 0 | 100.00% |
| 180+1 | 2 | 0 | 4 | 0 | 33.33% |
| **180+2** | **41** | **7** | **69** | **4** | **38.03%** |
| 180+3 | 1 | 0 | 0 | 0 | 100.00% |
| 300+0 | 0 | 0 | 1 | 0 | 0.00% |
| 300+2 | 0 | 0 | 1 | 0 | 0.00% |
| 300+3 | 4 | 0 | 0 | 0 | 100.00% |
| 300+5 | 2 | 0 | 0 | 0 | 100.00% |
| 420+2 | 0 | 1 | 0 | 0 | 50.00% |
| 480+2 | 1 | 0 | 0 | 0 | 100.00% |
| 600+2 | 0 | 0 | 1 | 0 | 0.00% |
| 600+3 | 2 | 0 | 0 | 0 | 100.00% |
| 600+5 | 2 | 1 | 0 | 0 | 83.33% |
| 900+3 | 0 | 1 | 0 | 0 | 50.00% |
| 900+5 | 1 | 0 | 1 | 0 | 50.00% |
| 900+10 | 0 | 0 | 1 | 0 | 0.00% |
| **All clocks, descriptive only** | **61** | **11** | **80** | **5** | **43.75%** |

The [machine-readable audit](audit.json) also separates session, color, event,
clock and opponent type. There are 150 rated and seven casual game records;
the three rated bullet records include one unfinished game. Opponents repeat,
and matchmaking is selected rather than randomized. The mixed-clock total,
session-to-session changes, account ratings and comparison with older live
sessions do **not** isolate an Elo gain. Functional evidence here is legal
integration and reproducible decision errors; measured strength remains the
separate controlled v0.22 release evidence.

## Shutdown, reconnects and actual concurrency

- **Expected laptop closure:** [2byZKYLT](https://lichess.org/2byZKYLT) is the
  only own time-forfeit loss. It ended after `1.e4 Nc6 2.Nf3`. The reported
  closure explains this interruption; it is not evidence of an engine timeout
  defect. Trace wall timestamps span 19:44:05.077–19:59:45.264 UTC, while the
  process's monotonic response duration is 7,461.685 ms. The roughly 933-second
  difference supports host suspension. No OS sleep log was collected. The loss
  stays in every unadjusted table. Excluding only this reported interruption
  gives 61/11/79 over 151 finished games, 44.04%; the 180+2 counterpart is
  41/7/68 over 116, 38.36%. These are operational sensitivity views.
- **Reconnects and matchmaking errors:** 19 control-stream reconnect markers,
  63 HTTP rate-limit error blocks, 61 opponent daily-cap blocks, and four rating
  mismatch blocks. These are grouped incidents, not traceback-line counts.
  They do not establish 19 lost games or engine crashes. Detailed incident
  records remain in each source archive.
- **Observed overlapping games:** session two ran
  [2WGLvZUk](https://lichess.org/2WGLvZUk) and
  [vMJ64VjF](https://lichess.org/vMJ64VjF), both against Cloudia-Chess,
  simultaneously. Trace timestamps show 20 overlapping search-interval pairs
  between two engine processes around Nov 19 17:30–17:32 UTC. The client log
  explicitly reports process count two. Both games were wins. Keep both distinct
  IDs and flag the shared workload; the concurrency-one configuration did not
  enforce exclusively serial play. A reconnect cause is not established.
- **Unfinished records:** `bbUUoSNG`, `7zlgD0un`, `ZsXkGPEs`, `bB4Gj7g3`,
  `DijN8b0i` remain preserved as abandoned/unfinished, with zero or one ply.
  The opponent's time-forfeit in `iJ00aKHL` remains an ordinary recorded win,
  distinct from the own host interruption.

All 7,507 completed searches have metrics and legal returned moves, with no
fallback, cancelled or incomplete search in the saved reports. Maximum engine
hard-budget overrun was 8.421 ms; maximum response-minus-search overhead was
3.508 ms. Tiny positive overruns in the raw counters are not clock forfeits.
171 requests arrived with the bot's own clock below ten seconds. Longer 24–37
second responses on slower clocks match their allocated budgets.

There are 7,429 uniquely linked searches, 75 ambiguous searches and three
missing links. Ambiguity is retained rather than resolved by guessing a game
from a shared opening. Consequently some per-position live metrics are absent
even though every played ply has referee analysis. The three missing links
include the unsubmitted move from the known closure. Reported aggregate search
elapsed time is 40,478.664 seconds, not CPU utilization and not suspension-
inclusive wall time. The audit checks actual trace wall timestamps separately.

## Reading the findings

- [Review and revised priorities](REVIEW.md): confirmed errors, negative
  controls, mechanism limits and ranked next work.
- [All 32 confirmations](CONFIRMED.md): scores, clocks, live depths and legal
  continuation samples, including the cases whose interpretation changed.
- [Processing and preservation](PROCESS.md): deduplication, referee method,
  selection, exact provenance and reproduction.
- [Verification](TESTING.md) and [handoff](HANDOFF.md): completed checks and the
  explicit boundary before any further experiment.
