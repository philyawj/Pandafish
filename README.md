# Pandafish

A bare-bones Go chess engine written in Go. Pandafish v0.4 speaks UCI over
stdin/stdout and looks ahead one opponent reply using depth-two minimax.
Standard chess only.

The engine code is in `main.go` and `evaluation.go`.
[corentings/chess](https://github.com/corentings/chess) handles chess rules,
including castling, en passant, and promotions. The engine assigns values of
1 to pawns, 3 to knights and bishops, 5 to rooks, and 9 to queens. It picks a
legal move by scoring every possible opponent reply from the moving side's
perspective. It assumes the opponent chooses the lowest score, then chooses
the move with the highest such worst-case score. Equal scores are still
selected randomly. Checkmate scores above any material advantage, being
checkmated scores below any material disadvantage, and stalemate scores zero.
Other draw rules and positional features have no separate search evaluation.

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
by `moves ...`. Normal `go` commands return after the fixed depth-two search;
clock/depth limits are ignored. `searchmoves` restricts
the choices. `go infinite` waits for `stop`;
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

[Benchmark results and progress](benchmarks/README.md) track each version
against Stockfish Skill Level 1 and the previous Pandafish release. The local
runner saves game records, settings, and a results table separately from the
engine. These matches measure development progress, not a Lichess rating.
Use the [development candidate workflow](benchmarks/README.md#test-a-development-candidate)
to test changes before tagging a release.
The long-term goal is to progress through Stockfish Skill Levels 1 to 20,
preserving a confirmed result above 50% before advancing to the next level.
The [benchmark decision process](benchmarks/README.md#deciding-whether-a-change-helps)
explains routine 100-game matches, optional longer comparisons, and the limits
of those scores.
The [release workflow](RELEASING.md) covers preparing a version, reviewing and
merging its PR, and publishing the release from the merged commit.

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
random tie breaking are unchanged. Focused position tests and the benchmark
workflow verified the development candidate against v0.2 and Stockfish. The
saved 100-game matchups scored 40 wins and 60 draws against v0.2, and 100 losses
against Stockfish Skill Level 1, with no forfeits or unfinished games.
The [development benchmark records](benchmarks/results/README.md) retain their
original candidate labels and binary hashes.

An earlier v0.3 screening run recorded Pandafish's first local benchmark draw
against Stockfish 19 at Skill Level 1: a stalemate with Pandafish playing White
in round 25. The milestone game is preserved separately from routine results.
[Download the game PGN](benchmarks/milestones/first-stockfish-level-1-draw.pgn).

**v0.4: Know when the game is over (complete).** The evaluator recognizes
checkmate and stalemate, including an opponent's immediate mating reply.
Depth-two search, material values, and random tie breaking are unchanged.
Repetition and fifty-move draw detection inside the search remain future work.

The development candidate scored 44 wins, 39 draws, and 17 losses against v0.3
in 100 games, a 63.5% match score. This is a measured improvement over v0.3
under the recorded conditions. Focused tests also verify that both colors
choose mate over material, avoid immediate mating replies, and avoid stalemate
when ahead. The sample supports this early increment; it does not establish
a precise Elo gain or stronger play against every opponent or time control.
The candidate lost all 100 games against Stockfish Level 1. Both matchups had
no forfeits or unfinished games. [Results and conditions](benchmarks/results/README.md)
retain their original development labels and hashes.
