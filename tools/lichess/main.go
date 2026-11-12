// Command lichess records local lichess-bot sessions and reports on saved evidence.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "lichess:", err)
		os.Exit(1)
	}
}

func execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("use record, report, or process; see tools/lichess/README.md")
	}
	switch args[0] {
	case "record":
		flags := flag.NewFlagSet("record", flag.ContinueOnError)
		var o recordOptions
		flags.StringVar(&o.Client, "client", "", "lichess-bot checkout directory (required)")
		flags.StringVar(&o.Config, "config", "config.yml", "config path, relative to client checkout")
		flags.StringVar(&o.Engine, "engine", "", "native UCI binary override; otherwise use config")
		flags.StringVar(&o.Session, "session", "", "new session directory; default local-games/session-TIMESTAMP")
		flags.StringVar(&o.Python, "python", "", "Python executable; default CLIENT/venv/bin/python")
		flags.StringVar(&o.Hardware, "hardware", "", "hardware description saved with the session")
		flags.StringVar(&o.Workload, "workload", "", "other workload, such as streaming and browser")
		flags.BoolVar(&o.PrepareOnly, "prepare-only", false, "prepare a session without starting the client or connecting online")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("unexpected arguments")
		}
		if o.Session == "" {
			o.Session = filepath.Join("local-games", "session-"+time.Now().Format("20060102-150405"))
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
	case "process":
		flags := flag.NewFlagSet("process", flag.ContinueOnError)
		var o processOptions
		flags.StringVar(&o.Session, "session", "", "finished session directory")
		flags.StringVar(&o.Archive, "archive", "", "new or matching archive directory")
		flags.StringVar(&o.Stockfish, "stockfish", "stockfish", "Stockfish executable")
		flags.StringVar(&o.Bot, "bot", "GoPandaFish", "bot account name in PGNs")
		flags.Int64Var(&o.Nodes, "nodes", 100000, "nodes per independent Stockfish search")
		flags.IntVar(&o.Limit, "limit", 0, "preview at most this many plies; zero resumes all")
		flags.StringVar(&o.Profile, "profile", "", "analysis directory name override for a separate method/version")
		flags.StringVar(&o.Focus, "focus", "", "JSON candidate list for a selected-position follow-up")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return fmt.Errorf("unexpected arguments")
		}
		return runProcess(o)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
