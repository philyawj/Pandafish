package main

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

func TestUCI(t *testing.T) {
	tests := []struct {
		name, commands, want string
	}{
		{"handshake", "\n  uci \nisready\nsetoption name Hash value 16\nunknown\nstop\nquit\nisready", "id name Pandafish v0.1\nid author philyawj\nuciok\nreadyok\n"},
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

func TestRandomLegalMove(t *testing.T) {
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
