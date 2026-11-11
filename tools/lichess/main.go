// Command lichess records local lichess-bot sessions and reports on saved evidence.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "lichess:", err)
		os.Exit(1)
	}
}

func execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("use record or report; see tools/lichess/README.md")
	}
	switch args[0] {
	case "record":
		flags := flag.NewFlagSet("record", flag.ContinueOnError)
		var o recordOptions
		flags.StringVar(&o.Client, "client", "", "lichess-bot checkout directory (required)")
		flags.StringVar(&o.Config, "config", "config.yml", "config path, relative to client checkout")
		flags.StringVar(&o.Engine, "engine", "", "native UCI binary override; otherwise use config")
		flags.StringVar(&o.Session, "session", "", "new session directory (required; use local-games/NAME)")
		flags.StringVar(&o.Python, "python", "", "Python executable; default CLIENT/venv/bin/python")
		flags.BoolVar(&o.PrepareOnly, "prepare-only", false, "prepare a session without starting the client or connecting online")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("unexpected arguments")
		}
		return record(o)
	case "report":
		flags := flag.NewFlagSet("report", flag.ContinueOnError)
		session := flags.String("session", "", "saved session directory (required)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("unexpected arguments")
		}
		return report(*session)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
