# Pandafish

A bare-bones Go chess engine. v0.1 chooses uniformly at random from the legal
moves and speaks UCI over stdin/stdout. Standard chess only.

All engine code is in `main.go`.
[corentings/chess](https://github.com/corentings/chess) handles chess rules,
including castling, en passant, and promotions. There is no evaluation or search.

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

The last response is `bestmove` followed by a random legal black move. When no
legal move exists, the engine returns `bestmove 0000`.

The engine accepts `position startpos` or `position fen ...`, optionally followed
by `moves ...`. Normal `go` commands return immediately; clock/depth limits are
ignored. `searchmoves` restricts the choices. `go infinite` waits for `stop`;
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
move comes from Pandafish. Keep `uci_options` empty: v0.1 has no Hash, Threads, or
other configurable engine options. Casual play follows the guidance in the
[lichess-bot default config](https://github.com/lichess-bot-devs/lichess-bot/blob/master/config.yml.default)
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

## Bots using Pandafish

Have you run a Lichess bot with any version of Pandafish? Add it here with a
pull request, or open an issue with your bot's profile link and the version you
used. Past and current users are welcome.

| Bot | Pandafish versions used |
| --- | --- |
| [GoPandaFish](https://lichess.org/@/GoPandaFish) | v0.1 |

## Milestones and plans

**v0.1: Get games working (complete).** The first step was to build a minimal
Go engine that speaks UCI and can complete games on Lichess through lichess-bot.
Random legal moves were enough to prove that connection worked. GoPandaFish's
first successful games included:

- [Casual correspondence / unlimited game](https://lichess.org/p3w5nk1V)
- [Casual timed game: 3+0 blitz](https://lichess.org/wb007ifr)

Both games reached checkmate. This milestone was about getting the engine
playing complete games, regardless of the result.

**Next: Choose moves for a reason.** Add chess logic so Pandafish evaluates its
choices and selects moves beyond uniform randomness. The next step is to build
playing strength on top of the working UCI and lichess-bot integration, keeping
the engine small and easy to evolve.
