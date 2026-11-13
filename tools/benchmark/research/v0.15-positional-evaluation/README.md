# Board-aware evaluation: experiment and findings

Published in v0.14.2 as a research archive. The release retains v0.14.1's
playing behavior; experimental evaluation is archived data, not active code.

This candidate starts from v0.14.1 and the
[archived selective-search findings](../v0.15-selective-search/README.md).
The unpromoted PVS/LMR policy is not enabled. The hypothesis is that richer
board-derived judgments can improve quiet choices even when extra nominal
search depth did not. No strength gain is assumed before the game screen.

The coherent scope is pawn structure, real king shelter, and restrained minor
activity. Material, original placement formulas, alpha-beta, quiescence,
transposition identity/layout, root ties, and time allocation are unchanged.
No opening database, learned weights, bishop-pair bonus, queen-mobility reward,
repetition logic, new UCI option, or incremental evaluation cache is added.

## Fixed initial coefficients

Original coefficients below were selected before reviewing match scores.
There is no weight sweep or automatic favorable-result retry.

| Feature | Middlegame cp | Endgame cp | Definition |
| --- | --- | --- | --- |
| Doubled pawn | -12 per extra | -18 per extra | File count minus one |
| Isolated pawn | -10 per pawn | -12 per pawn | No friendly pawn on either adjacent file |
| Passed pawn | 2 x rank squared | 4 x rank squared | Relative rank 1 on the starting pawn rank; no enemy pawn ahead on own/adjacent files |
| Blocked passer | Half bonus | Half bonus | The next forward square is occupied |
| Rear pawn on a friendly occupied file | No passer bonus | No passer bonus | Friendly pawn ahead on the same file avoids doubled passer credit |
| King shield one rank ahead | +14 per file | 0 | Own/adjacent king files; one pawn credit per file |
| King shield two ranks ahead | +7 per file | 0 | Used only when no closer shield pawn exists on that file |
| Semi-open file around king | -12 per file | 0 | No friendly pawn on own/adjacent file, enemy rook/queen remains |
| Fully open file around king | Additional -8 | 0 | Neither side has a pawn on that file, enemy rook/queen remains |
| Knight activity | +4 per destination | +4 per destination | Exclude friendly occupancy and enemy pawn attacks |
| Bishop activity | +3 per destination | +3 per destination | Stop rays at blockers; same exclusions |

Relative rank is capped at six for the bonus. Pawn and shelter terms use the
existing material phase: knight/bishop 1, rook 2, queen 4; initial/capped total
24. Shelter fades to zero in pawn-only endings, while existing endgame king
centralization remains. Doubled and isolated penalties can overlap; they are
kept explicit in diagnostics rather than disguised as one feature.

Mobility is pseudo-legal attack/destination availability, not legal or fully
safe mobility. Pins, attacks by other piece types, and king legality are not
modeled by this small preference. Capturing a hanging piece still belongs to
search. Queen activity is deliberately not rewarded with a broad mobility term.
All inputs are board-derived, so no new history-sensitive cache context arises.

## Cost and controls

A local board summary holds occupancy/pawn files, pawn attacks, king squares,
heavy-piece counts, and phase. It is built once per static evaluation. Feature
masks are internal diagnostic controls fixed per compiled binary, not UCI
options. Fresh engine processes/tables are used for ablations.

Unit checks cover disabled-feature parity with the released score on 19 varied
positions, material/component accounting, hand-counted pawn structures for both
colors, real shelter, phase limits, mobility and pawn attacks against independent
library attack maps, symmetry, terminal dominance, root immutability, and zero
heap allocations. The old en-passant test assumed the resulting pawn was worth
about 100 cp in total. It now checks material remains 100 cp and quiescence
returns the correct positional leaf, including the new advanced-passer value,
for both colors. This is a scoring-contract update, not a changed chess rule.

The isolated benchmark measures static-call overhead separately from whole
search. A higher cost is acceptable only if the complete evaluation earns it;
node counts alone cannot establish that. Compare baseline, pawn/king-only, and
full-feature position probes before the routine games. Do not tune coefficients
from the selected diagnostic or the match scores.

## Screen and decision

Use `v0.15-evaluation-dev` as a distinct label from the archived failed search
run. Run 50 games each against immutable v0.14.1, Stockfish Skill 2, and pinned
Sacrifice, all three matchups concurrent with the standard clocks/openings.
The reference external profile is 60% / 34%, using the playing algorithm that
v0.14.1 retains. The failed search candidate is not the control.

Before reviewing scores, the descriptive visible-gain goal remains at least
+10 percentage points on one external opponent without a decline greater than
five on the other. This is not statistical significance or an extra automatic
runner gate. Review the direct previous-version result and reliability too.
Preserve flat or negative evidence; no automatic longer screen or repeated run
until favorable. Code changes, configs, binaries, and results remain local;
check in with a recommendation before preparing a PR or release.

## Completed screen: hold

The [full record](../../results/v0.15-evaluation-dev-vs-v0.14.1-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba/README.md)
contains the tested source, diagnostics, all 150 games, timing, and checks.
The candidate scored 57% against v0.14.1, 45% against Stockfish Skill 2, and 39%
against pinned Sacrifice, with no forfeits or unfinished games. It passed the
previous-release criterion but missed the intended substantial external gain.
Selected-position improvements were insufficient evidence to promote. Keep the
active v0.14.1 engine and investigate a smaller pawn/king-only hypothesis next;
that ablation has not yet earned promotion through games.

## Replay the archived candidate

From the repository root, restore the exact source overlay into its recorded
base checkout. The shared Go helper validates paths before writing files.

```sh
repo_root="$PWD"
replay_root=$(mktemp -d)
record_dir="$repo_root/tools/benchmark/results/v0.15-evaluation-dev-vs-v0.14.1-and-stockfish-2-and-sacrifice-v0.2.0-b528f0ba"
git archive v0.14.1 | tar -x -C "$replay_root"
cp tools/benchmark/research/v0.15-selective-search/restore-snapshot.go.txt /tmp/pandafish-restore-snapshot.go
go run /tmp/pandafish-restore-snapshot.go "$record_dir/candidate-source.json" "$replay_root"
(cd "$replay_root" && go test ./...)
(cd "$replay_root" && go build -o /tmp/pandafish-evaluation-replay .)
```

For the pawn/king-only diagnostic, apply `pawn-king-source.json` to a separate
restored copy before building it. Build each policy separately with fresh engine
processes. Source recovery does not promise identical binary hashes or game
outcomes across compilers, random ties, or scheduling.

## Research-release verification

The release restores evaluation and tests from v0.14.1. Only engine identity and
its handshake expectation change to v0.14.2. See
[release checks](research-release-checks.json) for normal/race tests, vet, ARM64
build, offline wrapper, replay verification, and short integration games. The
experimental scores above belong to the archived development binary, not this
research release.

## Representative code and what it teaches

The archived implementation adds board context rather than another placement
formula. For example, a passed pawn has no opposing pawn ahead on its own or
adjacent files, while only the foremost friendly pawn receives the bonus:

```go
passed := board.pawns[side][file]&ahead == 0
for adjacent := max(0, file-1); adjacent <= min(7, file+1); adjacent++ {
    passed = passed && board.pawns[side.Other()][adjacent]&ahead == 0
}
```

Its initial advancement bonus was deliberately simple:

```go
mg, eg := 2*advance*advance, 4*advance*advance
if nextRank >= 0 && nextRank < 8 && board.pieces[nextRank*8+file] != chess.NoPiece {
    mg, eg = mg/2, eg/2
}
```

That recognizes a blocker, but does not establish that the pawn is supported,
safe, or able to promote. Similarly, the shelter score counts actual nearby
pawns but does not trace enemy attacks. Minor activity excludes enemy pawn
attacks but ignores pins and other attackers. These are testable approximations,
not evidence that their particular weights improve a complete engine.

The promising diagnostic and disappointing external match combination is the
main lesson. Correct feature definitions and attractive position scores do not
settle the tradeoff between evaluation incentives and timed search cost. The
next assessment should compare reference engines afresh, isolate a clear gap,
and state what result would justify retaining it before another game screen.
