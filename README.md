# Pandafish

A bare-bones chess engine written in Go. Pandafish v0.7 speaks UCI over
stdin/stdout and uses time-aware iterative deepening with alpha-beta pruning
and bounded quiescence. Standard chess only.

The engine code is in `main.go`, `limits.go`, `iterative.go`, `search.go`, and
`evaluation.go`.
[corentings/chess](https://github.com/corentings/chess) handles chess rules,
including castling, en passant, and promotions. The engine assigns values of
1 to pawns, 3 to knights and bishops, 5 to rooks, and 9 to queens. It picks a
legal move by completing successively deeper searches, assuming each side
chooses its best continuation. It keeps the last fully completed result when
time runs out, with a legal fallback if no depth completes. Alpha-beta
pruning skips branches that cannot improve a search result, with captures
searched first to help pruning. Equal best moves are selected randomly;
root pruning keeps their scores exact while discarding inferior choices.
Checkmate scores above any material advantage, being
checkmated scores below any material disadvantage, and stalemate scores zero.
Other draw rules and positional features have no separate search evaluation.

Quiescence means quietness. At the normal search boundary, it examines
captures and promotions, and all legal escapes when the king is in check.
Away from check it can keep the current material score instead of making a
capture. This helps avoid counting a captured piece before considering the
recapture. The extension is capped at two additional plies. At that cap,
static evaluation is used even if captures or checks remain, so long tactical
sequences can still be missed.

Clone this repository and build with Go 1.25 or newer:

```sh
git clone https://github.com/philyawj/Pandafish.git
cd Pandafish
go build -o pandafish .
go test ./...
```

Build on the machine that will run the engine. On Apple silicon, this produces
a native ARM64 executable. On Windows, use `go build -o pandafish.exe .`.

Try it:

```sh
./pandafish
```

Then enter:

```text
uci
isready
position startpos moves e2e4
go movetime 100
```

Wait for `bestmove` followed by the selected legal black move, then enter
`quit`. Sending `quit` or closing stdin immediately after `go` cancels the
search. When no legal move exists, the engine returns `bestmove 0000`.

The engine accepts `position startpos` or `position fen ...`, optionally
followed by `moves ...`. Search accepts `movetime`, `wtime`, `btime`, `winc`,
`binc`, `movestogo`, and `depth` (1 through 64). Times are milliseconds.
`searchmoves` restricts root choices and must come last. A depth-only request
has no time deadline; time and depth limits otherwise apply together.

Bare `go` uses a 100 ms allowance. Explicit `movetime` takes precedence over
clock allocation, capped to the usable side-to-move clock when supplied.
With game clocks, reserve the smaller of 10 ms or 10% of the remaining time, divide the
usable clock by `movestogo` (30 if absent), and add 75% of the increment.
Spending is capped to half the usable clock, so uncredited increment cannot
exhaust the current clock.

Inside the resulting allowance, reserve the smaller of 5 ms or 10% for
scheduling and communication. The search stops starting new iterations at
75% of that hard budget and checks cancellation inside both main search and
quiescence. Scheduling can still affect observed response time; the engine
makes no absolute real-time guarantee.

`isready` remains responsive during search. `stop` returns the last completed
result or legal fallback. Position changes, new games, and replacement
searches cancel old work. `go infinite` holds its result until `stop`.
`go ponder` also waits; `ponderhit` restarts a fresh timed search using its
original limits. Pondering remains disabled in the example configuration.
Invalid recognized search limits produce an `info string` error; unknown
commands and unsupported engine options are ignored.

To play through [lichess-bot](https://github.com/lichess-bot-devs/lichess-bot),
install that client and its Python dependencies following its setup instructions.
Use a Lichess BOT account and an
[API token with the Play bot moves permission](https://lichess.org/account/oauth/token/create?scopes[]=bot:play).
Pandafish itself needs no credentials; lichess-bot handles the connection.

Copy the supplied [example config](lichess-bot-config-example.yml) from this
repository's root:

```sh
cp lichess-bot-config-example.yml lichess-bot-config.yml
```

Edit **only your local copy**: replace `YOUR_LICHESS_BOT_API_TOKEN` and review the
three numbered engine lines: `dir` (binary directory), `name` (binary filename),
and `working_dir` (normally the same directory). Use absolute paths.
`lichess-bot-config.yml` is gitignored; the example must always contain placeholders.

The example accepts incoming **casual standard chess challenges from humans**,
with one active game at a time. It accepts all BOT-supported time controls,
including correspondence and unlimited games. Rated games and bot challengers
are declined; automatic outgoing challenges and pondering are disabled.
Lichess does not allow BOT accounts to play UltraBullet.

From your lichess-bot checkout, run the following on macOS/Linux, replacing the
config path with your own. This assumes its Python environment is named `venv`:

```sh
env -u LICHESS_BOT_TOKEN ./venv/bin/python lichess-bot.py --config /absolute/path/to/Pandafish/lichess-bot-config.yml
```

The `env` prefix prevents another account's environment token from overriding
the config. On Windows, use `venv\Scripts\python.exe` and clear any existing
`LICHESS_BOT_TOKEN` environment variable first. Stop the client with Ctrl+C.

The example keeps opening books, online moves, and tablebases disabled so every
move comes from Pandafish. Keep `uci_options` empty: the engine has no Hash,
Threads, or other configurable engine options. Casual play follows the guidance
in the [lichess-bot default config](https://github.com/lichess-bot-devs/lichess-bot/blob/master/config.yml.default)
for engines that do not try to win.

Contributions are welcome. Keep changes focused and the engine simple. Before
opening a pull request, format changed Go files and run:

```sh
go test ./...
go vet ./...
go build -o pandafish .
git diff --check
```

For engine changes, also check that lichess-bot can initialize the rebuilt
executable and obtain legal moves through its UCI wrapper. Its built-in config
validation checks configuration and paths; a move-sequence smoke test checks
engine integration. These local checks do not require a live Lichess game.
Include the checks you ran in your PR. Never commit tokens, private configs,
or compiled binaries. Before committing, inspect `git diff --cached`.

Pandafish is licensed under the [MIT License](LICENSE). The chess library's
license is reproduced in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
Include both notices when distributing compiled Pandafish binaries.

## Local benchmarks

[Developer tools](tools/README.md) describes the shared tools layout and the
planned Lichess session recorder. The recorder has not been implemented.

[Benchmark documentation](tools/benchmark/README.md) covers the validation workflow,
playing-strength comparisons, results, and their interpretation.
See the [release workflow](RELEASING.md) for preparing and publishing a version.

## Milestones and plans

**v0.1: Get games working (complete).** This release built a minimal Go engine
that speaks UCI and can complete games on Lichess through lichess-bot. Random
legal moves were enough to prove that connection worked. GoPandaFish's first
successful games included:

- [Casual correspondence / unlimited game](https://lichess.org/p3w5nk1V)
- [Casual timed game: 3+0 blitz](https://lichess.org/wb007ifr)

Both games reached checkmate. This milestone was about getting the engine
playing complete games, regardless of the result.

**v0.2: Choose moves for a reason (complete).** Pandafish now evaluates the
immediate material balance after each legal move, chooses the move with the
best score from the moving side's perspective, and breaks ties randomly. This
is a one-ply material evaluation. It does not consider the opponent's reply,
checkmate, or positional features.

**v0.3: Look ahead one reply (complete).** Depth-two minimax now accounts
for the opponent's best immediate material response. The v0.2 piece values and
random tie breaking are unchanged.

**v0.4: Know when the game is over (complete).** The evaluator recognizes
checkmate and stalemate, including an opponent's immediate mating reply.
Depth-two search, material values, and random tie breaking are unchanged.
Repetition and fifty-move draw detection inside the search remain future work.

**v0.5: Build a deeper search (complete).** Replace the fixed two-ply
calculation with recursive minimax using alpha-beta pruning at three plies.
The v0.4 material evaluator, mate/stalemate handling, and random tie breaking
remain unchanged. This release retained fixed-depth search and ignored
clock/depth limits.

**v0.6: Finish the exchange before counting the pieces (complete).**
Add bounded quiescence search for captures, promotions, and check evasions
beyond the three-ply boundary. Material values, mate/stalemate scoring, and
random tie breaking are unchanged. This release still used a fixed main
search depth and ignored clock/depth limits.

**v0.7: Know when to stop thinking (complete).** Honor UCI time and
depth limits with conservative clock allocation, iterative deepening, and
interruptible search. Keep the last completed iteration and a legal fallback.
Material scoring, mate/stalemate handling, and the two-ply quiescence cap are
unchanged. Live-game logging and metric capture are planned for v0.7.1.
