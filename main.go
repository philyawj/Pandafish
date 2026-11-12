package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

func main() {
	var trace *telemetry.Recorder
	if dir := os.Getenv(telemetry.Environment); dir != "" {
		var err error
		trace, err = telemetry.Open(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "capture:", err)
			os.Exit(1)
		}
	}
	defer trace.Close()
	if err := runEngine(os.Stdin, telemetry.Output{Destination: os.Stdout, Recorder: trace}, trace); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type inputLine struct {
	text string
	err  error
	eof  bool
}

type activeSearch struct {
	cancel   context.CancelFunc
	results  chan searchResult
	result   searchResult
	finished bool
	limits   goLimits
}

func run(in io.Reader, out io.Writer) error { return runEngine(in, out, nil) }

func runEngine(in io.Reader, out io.Writer, trace *telemetry.Recorder) (err error) {
	defer func() {
		status := "ok"
		if err != nil {
			status = "error"
		}
		if e := trace.Write("exit", status); err == nil {
			err = e
		}
	}()
	// Only this loop owns output and protocol state. The search worker owns its
	// position until joined; a replacement cannot race old position caches.
	lines, done := make(chan inputLine), make(chan struct{})
	defer close(done)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			select {
			case lines <- inputLine{text: scanner.Text()}:
			case <-done:
				return
			}
		}
		select {
		case lines <- inputLine{eof: true, err: scanner.Err()}:
		case <-done:
		}
	}()
	game := chess.NewGame()
	table := newTable()
	var active *activeSearch
	finish := func(emit bool, reason string) {
		if active == nil {
			return
		}
		active.cancel()
		if !active.finished {
			active.result = <-active.results
		}
		if emit {
			if reason != "" {
				active.result.metrics.Reason = reason
			}
			reportSearch(out, active.result)
			fmt.Fprintln(out, "bestmove", active.result.move)
		}
		active = nil
	}
	defer func() { finish(false, "") }()
	start := func(limits goLimits) {
		ctx, cancel := context.WithCancel(context.Background())
		results := make(chan searchResult, 1)
		budget, _ := limits.budget(game.Position().Turn()) // Validated before start.
		active = &activeSearch{cancel: cancel, results: results, limits: limits}
		position, started := game.Position(), time.Now()
		go func() { results <- iterativeSearchWithTable(ctx, position, limits, budget, started, table) }()
	}
	for {
		var results <-chan searchResult
		if active != nil && !active.finished {
			results = active.results
		}
		select {
		case result := <-results:
			active.result, active.finished = result, true
			if !active.limits.infinite && !active.limits.ponder {
				finish(true, "")
			}
		case line := <-lines:
			if line.eof {
				return line.err
			}
			_ = trace.Write("in", line.text) // Recording failures surface on exit; play continues.
			args := strings.Fields(line.text)
			if len(args) == 0 {
				continue
			}
			switch args[0] {
			case "uci":
				fmt.Fprintln(out, "id name Pandafish v0.10-dev\nid author philyawj\nuciok")
			case "isready":
				fmt.Fprintln(out, "readyok")
			case "ucinewgame":
				finish(false, "")
				game = chess.NewGame()
				table.clear()
			case "position":
				next, err := parsePosition(args[1:])
				if err != nil {
					fmt.Fprintln(out, "info string invalid position:", err)
					continue
				}
				finish(false, "")
				game = next
			case "go":
				limits, err := parseGo(args[1:])
				if err == nil {
					_, err = limits.budget(game.Position().Turn())
				}
				if err != nil {
					fmt.Fprintln(out, "info string invalid go:", err)
					continue
				}
				finish(false, "")
				start(limits)
			case "stop":
				finish(true, "stop")
			case "ponderhit":
				if active != nil && active.limits.ponder {
					limits := active.limits
					limits.ponder = false
					finish(false, "")
					if _, err := limits.budget(game.Position().Turn()); err != nil {
						fmt.Fprintln(out, "info string invalid go:", err)
					} else {
						start(limits)
					}
				}
			case "quit":
				return nil
			}
		}
	}
}

func parsePosition(args []string) (*chess.Game, error) {
	game := chess.NewGame()
	switch {
	case len(args) >= 1 && args[0] == "startpos":
		args = args[1:]
	case len(args) >= 7 && args[0] == "fen":
		fen, err := chess.FEN(strings.Join(args[1:7], " "))
		if err != nil {
			return nil, err
		}
		game = chess.NewGame(fen)
		args = args[7:]
	default:
		return nil, fmt.Errorf("expected startpos or fen with six fields")
	}
	if len(args) > 0 {
		if args[0] != "moves" {
			return nil, fmt.Errorf("expected moves")
		}
		for _, move := range args[1:] {
			if err := game.PushNotationMove(move, chess.UCINotation{}, nil); err != nil {
				return nil, err
			}
		}
	}
	return game, nil
}
