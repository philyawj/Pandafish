package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"slices"
	"strings"

	"github.com/corentings/chess/v2"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer) error {
	game := chess.NewGame()
	pending := ""
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "uci":
			fmt.Fprintln(out, "id name Pandafish v0.1\nid author philyawj\nuciok")
		case "isready":
			fmt.Fprintln(out, "readyok")
		case "ucinewgame":
			game, pending = chess.NewGame(), ""
		case "position":
			next, err := parsePosition(args[1:])
			if err != nil {
				fmt.Fprintln(out, "info string invalid position:", err)
				continue
			}
			game, pending = next, ""
		case "go":
			pending = randomMove(game, args[1:])
			// Infinite searches and pondering wait for the GUI to release the move.
			if !slices.Contains(args, "infinite") && !slices.Contains(args, "ponder") {
				fmt.Fprintln(out, "bestmove", pending)
				pending = ""
			}
		case "stop", "ponderhit":
			if pending != "" {
				fmt.Fprintln(out, "bestmove", pending)
				pending = ""
			}
		case "quit":
			return nil
		}
	}
	return scanner.Err()
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

func randomMove(game *chess.Game, args []string) string {
	var candidates []string
	searchMoves := slices.Index(args, "searchmoves")
	for _, move := range game.ValidMoves() {
		uci := (chess.UCINotation{}).Encode(game.Position(), &move)
		if searchMoves < 0 || slices.Contains(args[searchMoves+1:], uci) {
			candidates = append(candidates, uci)
		}
	}
	if len(candidates) == 0 {
		return "0000"
	}
	return candidates[rand.IntN(len(candidates))]
}
