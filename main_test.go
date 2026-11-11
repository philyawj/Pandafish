package main

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

func TestUCI(t *testing.T) {
	tests := []struct {
		name, commands, want string
	}{
		{"handshake", "\n  uci \nisready\nsetoption name Hash value 16\nunknown\nstop\nquit\nisready", "id name Pandafish v0.5\nid author philyawj\nuciok\nreadyok\n"},
		{"history", "position startpos moves e2e4 e7e5\ngo wtime 1000 btime 1000 searchmoves g1f3", "bestmove g1f3\n"},
		{"no implicit move", "position startpos\ngo searchmoves e2e4\ngo searchmoves e2e4", "bestmove e2e4\nbestmove e2e4\n"},
		{"new game", "position startpos moves e2e4\nucinewgame\ngo searchmoves d2d4", "bestmove d2d4\n"},
		{"castling", "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1\ngo searchmoves e1g1", "bestmove e1g1\n"},
		{"castling history", "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1g1 e8c8\ngo searchmoves f1e1", "bestmove f1e1\n"},
		{"en passant", "position fen k7/8/8/3pP3/8/8/8/7K w - d6 0 1\ngo searchmoves e5d6", "bestmove e5d6\n"},
		{"en passant history", "position fen k7/8/8/3pP3/8/8/8/7K w - d6 0 1 moves e5d6\ngo searchmoves a8b8", "bestmove a8b8\n"},
		{"promotion", "position fen 7k/P7/8/8/8/8/8/K7 w - - 0 1\ngo searchmoves a7a8n", "bestmove a7a8n\n"},
		{"promotion history", "position fen 7k/P7/8/8/8/8/8/K7 w - - 0 1 moves a7a8n\ngo searchmoves h8g8", "bestmove h8g8\n"},
		{"checkmate", "position fen 7k/6Q1/5K2/8/8/8/8/8 b - - 0 1\ngo", "bestmove 0000\n"},
		{"stalemate", "position fen 7k/5K2/6Q1/8/8/8/8/8 b - - 0 1\ngo", "bestmove 0000\n"},
		{"king safety", "position fen 4r2k/8/8/8/8/8/4R3/4K3 w - - 0 1\ngo searchmoves e2d2 e2e8", "bestmove e2e8\n"},
		{"infinite", "go infinite searchmoves e2e4\nisready\nstop\nstop", "readyok\nbestmove e2e4\n"},
		{"ponder", "go ponder searchmoves d2d4\nisready\nponderhit\nstop", "readyok\nbestmove d2d4\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := run(strings.NewReader(tt.commands), &out); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLegalMoveSelection(t *testing.T) {
	legal := strings.Fields("a2a3 a2a4 b2b3 b2b4 c2c3 c2c4 d2d3 d2d4 e2e3 e2e4 f2f3 f2f4 g2g3 g2g4 h2h3 h2h4 b1a3 b1c3 g1f3 g1h3")
	var out bytes.Buffer
	if err := run(strings.NewReader(strings.Repeat("go\n", 100)), &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 100 {
		t.Fatalf("got %d responses, want 100", len(lines))
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "bestmove ") || !slices.Contains(legal, strings.TrimPrefix(line, "bestmove ")) {
			t.Fatalf("invalid response: %q", line)
		}
	}
}

func TestInvalidPosition(t *testing.T) {
	for _, command := range []string{"", "fen", "fen invalid w - - 0 1", "startpos moves e2e5", "startpos unexpected"} {
		if _, err := parsePosition(strings.Fields(command)); err == nil {
			t.Errorf("accepted invalid position %q", command)
		}
	}
}

func TestEvaluateMaterial(t *testing.T) {
	fen, err := chess.FEN("7k/8/8/8/8/8/8/KQ6 w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}
	game := chess.NewGame(fen)
	if got := evaluate(game.Position(), chess.White); got != 9 {
		t.Fatalf("white score = %d, want 9", got)
	}
	if got := evaluate(game.Position(), chess.Black); got != -9 {
		t.Fatalf("black score = %d, want -9", got)
	}
}
func TestBestMoveCapturesQueen(t *testing.T) {
	fen, err := chess.FEN("7k/q7/8/8/8/8/8/R6K w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}
	game := chess.NewGame(fen)
	if got := bestMove(game, nil); got != "a1a7" {
		t.Fatalf("move = %q, want a1a7", got)
	}
}

func TestBestMoveConsidersOpponentReply(t *testing.T) {
	for _, tt := range []struct{ name, fen, moves, want string }{
		{"white avoids recapture", "3r2k1/8/8/3p4/8/8/8/3Q2K1 w - - 0 1", "d1d5 d1e2", "d1e2"},
		{"black avoids recapture", "3q2k1/8/8/8/3P4/8/8/3R2K1 b - - 0 1", "d8d4 d8e7", "d8e7"},
		{"white saves threatened queen", "6k1/8/8/8/8/2p5/1Q6/6K1 w - - 0 1", "g1h1 b2a2", "b2a2"},
		{"black saves threatened queen", "6k1/1q6/2P5/8/8/8/8/6K1 b - - 0 1", "g8h8 b7a7", "b7a7"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(append([]string{"fen"}, strings.Fields(tt.fen)...))
			if err != nil {
				t.Fatal(err)
			}
			before := game.FEN()
			if got := bestMove(game, strings.Fields("searchmoves "+tt.moves)); got != tt.want {
				t.Fatalf("move = %q, want %q", got, tt.want)
			}
			if game.FEN() != before || len(game.Moves()) != 0 {
				t.Fatal("search changed the root game")
			}
		})
	}
}

func TestSearchFindsBestCapture(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen r6k/7P/8/8/8/8/Q7/7K b - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	// Taking our queen leaves our pawn (1) against their rook (5).
	// Taking only the pawn would leave a much better material balance for us.
	if got := search(game.Position(), 1, -searchInfinity, searchInfinity); got != 4 {
		t.Fatalf("black search score = %d, want 4", got)
	}
}

func TestSearchWithoutRepliesUsesTerminalScore(t *testing.T) {
	for _, tt := range []struct {
		fen  string
		want int
	}{
		{"7k/6Q1/5K2/8/8/8/8/8 b - - 0 1", -mateScore},
		{"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1", 0},
	} {
		game, err := parsePosition(append([]string{"fen"}, strings.Fields(tt.fen)...))
		if err != nil {
			t.Fatal(err)
		}
		if got := search(game.Position(), searchDepth, -searchInfinity, searchInfinity); got != tt.want {
			t.Fatalf("terminal score = %d, want %d", got, tt.want)
		}
	}
}

func TestEvaluateTerminalPositions(t *testing.T) {
	for _, tt := range []struct {
		name, fen  string
		whiteScore int
	}{
		{"white mates", "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1", mateScore},
		{"black mates", "8/8/8/8/8/5k2/6q1/7K w - - 0 1", -mateScore},
		{"black stalemated", "7k/5K2/6Q1/8/8/8/8/8 b - - 0 1", 0},
		{"white stalemated", "8/8/8/8/8/6q1/5k2/7K w - - 0 1", 0},
		{"check with an escape", "k6r/8/8/8/8/8/8/7K w - - 0 1", -5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(append([]string{"fen"}, strings.Fields(tt.fen)...))
			if err != nil {
				t.Fatal(err)
			}
			for _, side := range []chess.Color{chess.White, chess.Black} {
				want := tt.whiteScore
				if side == chess.Black {
					want = -want
				}
				if got := evaluate(game.Position(), side); got != want {
					t.Fatalf("score for %v = %d, want %d", side, got, want)
				}
			}
		})
	}
}

func TestBestMoveRecognizesGameEndings(t *testing.T) {
	for _, tt := range []struct{ name, fen, moves, want string }{
		{"white mates instead of taking queen", "7k/8/5KQ1/8/8/8/8/1q6 w - - 0 1", "g6g7 g6b1", "g6g7"},
		{"black mates instead of taking queen", "1Q6/8/8/8/8/5kq1/8/7K b - - 0 1", "g3g2 g3b8", "g3g2"},
		{"white avoids mate after taking rook", "8/8/r7/8/8/5kq1/Q7/7K w - - 0 1", "a2a6 a2b2", "a2b2"},
		{"black avoids mate after taking rook", "7k/q7/5KQ1/8/8/R7/8/8 b - - 0 1", "a7a3 a7b7", "a7b7"},
		{"white avoids stalemate while ahead", "7k/5K2/8/6Q1/8/8/8/8 w - - 0 1", "g5g6 g5g4", "g5g4"},
		{"black avoids stalemate while ahead", "8/8/8/8/6q1/8/5k2/7K b - - 0 1", "g4g3 g4g5", "g4g5"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(append([]string{"fen"}, strings.Fields(tt.fen)...))
			if err != nil {
				t.Fatal(err)
			}
			// Validate both choices so an illegal alternative cannot hide a bug.
			for _, move := range strings.Fields(tt.moves) {
				copy, err := parsePosition(append(append([]string{"fen"}, strings.Fields(tt.fen)...), "moves", move))
				if err != nil || copy == nil {
					t.Fatalf("invalid fixture move %s: %v", move, err)
				}
			}
			if got := bestMove(game, strings.Fields("searchmoves "+tt.moves)); got != tt.want {
				t.Fatalf("move = %s, want %s", got, tt.want)
			}
		})
	}
}
