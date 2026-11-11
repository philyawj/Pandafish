# Pandafish

A bare-bones chess engine written in Go. Pandafish v0.6 speaks UCI over
stdin/stdout and uses three-ply alpha-beta search with bounded quiescence.
Standard chess only.

The engine code is in `main.go`, `search.go`, and `evaluation.go`.
[corentings/chess](https://github.com/corentings/chess) handles chess rules,
including castling, en passant, and promotions. The engine assigns values of
1 to pawns, 3 to knights and bishops, 5 to rooks, and 9 to queens. It picks a
legal move by looking ahead through its move, the opponent's reply, and its
next response, assuming each side chooses its best continuation. Alpha-beta
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
static evaluation is used even
if captures or checks remain, so long tactical sequences can still be missed.

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
printf 'uci\nisready\nposition startpos moves e2e4\ngo\nquit\n' | ./pandafish
```

The last response is `bestmove` followed by the selected legal black move.
When no legal move exists, the engine returns `bestmove 0000`.

The engine accepts `position startpos` or `position fen ...`, optionally followed
by `moves ...`. Normal `go` commands return after the three-ply main search
and bounded quiescence extension; clock/depth limits are ignored.
`searchmoves` restricts the choices. `go infinite` waits for `stop`;
`go ponder` waits for `stop` or `ponderhit`. Unknown commands and engine options
are ignored.

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

[Benchmark documentation](benchmarks/README.md) covers the validation workflow,
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
random tie breaking are unchanged. Time-aware iterative deepening remains
the next search foundation to add before increasing search work further.
