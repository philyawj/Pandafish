# Running Pandafish through lichess-bot

Pandafish v1.0 will power GoPandaFish and other Panda Bots on Lichess.
Engine development is paused while effort focuses on other projects.
These instructions prepare a local client; publishing a release does not start
a bot or change a private configuration.

To play through [lichess-bot](https://github.com/lichess-bot-devs/lichess-bot),
install that client and its Python dependencies following its setup instructions.
Use a Lichess BOT account and an
[API token with the Play bot moves permission](https://lichess.org/account/oauth/token/create?scopes[]=bot:play).
Pandafish itself needs no credentials; lichess-bot handles the connection.

Copy the supplied [example config](../lichess-bot-config-example.yml) from this
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
