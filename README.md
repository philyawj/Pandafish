# Pandafish

Pandafish v1.0 is a standard-chess engine written entirely in Go. It speaks
[UCI](docs/ENGINE.md) and chooses moves with time-aware iterative deepening,
alpha-beta search and bounded quiescence. A simple material/placement evaluator
and targeted mating-threat verification guide its decisions.

## Download v1.0

Download the archive for the operating system and processor from the
[v1.0 release](https://github.com/philyawj/Pandafish/releases/tag/v1.0).

| Platform | Processor | Download |
| --- | --- | --- |
| macOS | Apple silicon (M1 or newer) | [darwin-arm64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-darwin-arm64.tar.gz) |
| macOS | Intel | [darwin-amd64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-darwin-amd64.tar.gz) |
| Linux | x64 | [linux-amd64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-linux-amd64.tar.gz) |
| Linux | ARM64 | [linux-arm64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-linux-arm64.tar.gz) |
| Windows | x64 | [windows-amd64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-windows-amd64.zip) |
| Windows | ARM64 | [windows-arm64](https://github.com/philyawj/Pandafish/releases/download/v1.0/Pandafish-v1.0-windows-arm64.zip) |

Extract the archive. On macOS/Linux, open a terminal in the extracted directory
and run `./pandafish`; on Windows, run `.\pandafish.exe` from PowerShell.
The engine is a command-line UCI program: use a chess GUI or lichess-bot for a
board interface. [Download instructions](docs/DOWNLOADS.md) cover processor
selection, checksums and platform details. Archives include both required
license notices. Only macOS ARM64 is runtime-tested on the release host;
other targets are cross-built.

## Build from source

Go 1.25 or newer is required:

```sh
git clone https://github.com/philyawj/Pandafish.git
cd Pandafish
go build -o pandafish .
go test ./...
```

On Windows, use `go build -o pandafish.exe .`. Build on the target machine or use
the matching release archive. Contributor checks include the separate chess
module; see [CONTRIBUTING.md](CONTRIBUTING.md).

## Try the engine

Start the executable, then enter:

```text
uci
isready
position startpos moves e2e4
go movetime 100
```

Wait for `bestmove`, then enter `quit`. Closing stdin or sending `quit`
immediately after `go` cancels the search. The engine returns `bestmove 0000`
when no legal move exists. [Engine documentation](docs/ENGINE.md) describes
supported UCI commands, timing, search behavior and draw-rule limitations.

## Play on Lichess

Use the external lichess-bot client with a BOT account and a private API token.
Follow the [Lichess setup guide](docs/LICHESS.md) and the supplied
[example configuration](lichess-bot-config-example.yml). Credentials and private
configuration stay outside Git.

## Documentation and evidence

| Guide | Contents |
| --- | --- |
| [Engine behavior](docs/ENGINE.md) | Search, evaluation, UCI, timing and known limits |
| [Lichess setup](docs/LICHESS.md) | Client installation, configuration and launch instructions |
| [Benchmarks](tools/benchmark/README.md) | Reproducible comparisons, results and uncertainty |
| [Version history](docs/HISTORY.md) | Earlier milestones and original research links |
| [v1.0 investigation](docs/RESEARCH_V024.md) | Accepted tactical change and preserved unsuccessful attempts |
| [Developer tools](tools/README.md) | Benchmarking, recording and evidence tools |
| [Evidence storage](docs/EVIDENCE_STORAGE.md) | Existing archives and leaner future records |
| [Contributing](CONTRIBUTING.md) | Go-only source rules and validation workflow |

## Milestones and plans

**v1.0 marks a stable, usable engine and a pause in active development.**
The milestone recognizes:

- A working, tested UCI chess engine.
- Reproducible benchmarks and preserved research.
- Clear documentation and known limitations.
- A stable scope worth preserving.

Development is reaching diminishing returns within the current search and
evaluation design. v1.0 is being released as the stability milestone for this
work. GoPandaFish and other Panda Bots will use the engine to play on Lichess.
Future engine development is currently on pause while effort focuses on other
projects. Any return to development will begin with a deliberate scope review;
no next experiment or live session starts automatically.

The release includes bounded queen-contact threat defense. Its standard screen
passes the project's practical gate, but the small lead does not establish a
reliable general strength or Elo gain. Full evidence and limitations remain in
the [research report](tools/benchmark/research/v0.24-contact-threat/README.md).
Earlier milestones are preserved in [version history](docs/HISTORY.md).

Pandafish uses the [MIT License](LICENSE). The chess library's license is in
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md). Both notices accompany releases.
