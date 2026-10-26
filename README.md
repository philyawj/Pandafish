# Pandafish

A bare-bones Go chess engine. v0.1 chooses uniformly at random from the legal
moves and speaks UCI over stdin/stdout. Standard chess only.

All engine code is in `main.go`.
[corentings/chess](https://github.com/corentings/chess) handles chess rules,
including castling, en passant, and promotions. There is no evaluation or search.

Build with Go 1.25 or newer:

```sh
go build -o pandafish .
go test ./...
```

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

To connect through [lichess-bot](https://github.com/lichess-bot-devs/lichess-bot),
start with its default config and change these fields (preserve the other settings):

```yaml
engine:
  dir: "/absolute/path/to/Pandafish"
  name: "pandafish"
  protocol: "uci"
  ponder: false
  uci_options: {}

challenge:
  variants: [standard]
  modes: [casual]

matchmaking:
  challenge_mode: casual
```

Build the binary for the machine running lichess-bot. Keep opening books, online
moves, and tablebases disabled so every move comes from Pandafish. Remove the
default UCI options such as Hash and Threads; this engine advertises no options.
Use casual games for this random-move version, as directed by the
[lichess-bot default config](https://github.com/lichess-bot-devs/lichess-bot/blob/master/config.yml.default)
for engines that do not try to win.
