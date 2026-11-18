package chess

import "testing"

func TestSearchEvasionMaskChecksAndFallbacks(t *testing.T) {
	position := func(fen string) *Position {
		option, err := FEN(fen)
		if err != nil {
			t.Fatal(err)
		}
		return NewGame(option).Position()
	}
	for _, c := range []struct {
		fen     string
		squares []Square
	}{
		{"4r1k1/8/8/8/8/8/8/4K3 w - - 0 1", []Square{E2, E3, E4, E5, E6, E7, E8}},
		{"6k1/8/8/8/1b6/8/8/4K3 w - - 0 1", []Square{D2, C3, B4}},
		{"6k1/8/8/8/8/5n2/8/4K3 w - - 0 1", []Square{F3}},
		{"8/8/8/3k4/4P3/8/8/6K1 b - - 0 1", []Square{E4}},
		{"4r1k1/8/8/8/1b6/8/8/4K3 w - - 0 1", nil},
	} {
		p := position(c.fen)
		got, active := searchEvasionMask(p)
		want := bitboard(0)
		for _, s := range c.squares {
			want |= bbForSquare(s)
		}
		if !active || got != want {
			t.Fatal("checker mask", c.fen, got, want, active)
		}
	}
	p := position("4r1k1/8/8/8/8/8/8/4K3 w - - 0 1")
	p.inCheck = false
	if _, active := searchEvasionMask(p); active {
		t.Fatal("unchecked position uses mask")
	}
	p.inCheck = true
	p.board.bbWhiteKing = 0
	if _, active := searchEvasionMask(p); active {
		t.Fatal("missing king uses mask")
	}
	p.board.bbWhiteKing = bbForSquare(E1) | bbForSquare(D1)
	if _, active := searchEvasionMask(p); active {
		t.Fatal("nonunique king uses mask")
	}
	p = position("6k1/8/8/8/8/8/8/4K3 w - - 0 1")
	p.inCheck = true
	if _, active := searchEvasionMask(p); active {
		t.Fatal("inconsistent check flag uses mask")
	}
}
