package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

// Independent material-only oracle: try every legal recapture on the target,
// including alternatives of equal value. Either side may decline to exchange.
func exchangeOracle(position *chess.Position, target chess.Square) int {
	best := 0
	for _, move := range position.ValidMoves() {
		if move.S2() != target || (!move.HasTag(chess.Capture) && !move.HasTag(chess.EnPassant)) {
			continue
		}
		gain := pieceValues[position.Board().Piece(target).Type()]
		if move.HasTag(chess.EnPassant) {
			gain = pieceValues[chess.Pawn]
		}
		if move.Promo() != chess.NoPieceType {
			gain += pieceValues[move.Promo()] - pieceValues[chess.Pawn]
		}
		best = max(best, gain-exchangeOracle(position.Update(&move), target))
	}
	return best
}

func TestStaticExchange(t *testing.T) {
	for _, tt := range []struct {
		name, fen, move string
		want            int
	}{
		{"free pawn", "6k1/8/8/3p4/8/8/8/3R2K1 w - - 0 1", "d1d5", 100},
		{"defended pawn", "3r3k/8/8/3p4/8/8/8/3Q2K1 w - - 0 1", "d1d5", -800},
		{"defended queen", "3r2k1/8/8/8/3q4/8/8/3R2K1 w - - 0 1", "d1d4", 400},
		{"revealed rook recapture", "3r2k1/8/8/3p4/8/8/3R4/3R2K1 w - - 0 1", "d2d5", 100},
		{"pinned bishop cannot recapture", "4k3/4b3/3p4/8/8/Q7/8/4R1K1 w - - 0 1", "a3d6", 100},
		{"pinned rook captures along pin", "4k3/4r3/4p3/8/8/1Q6/8/4R1K1 w - - 0 1", "b3e6", -300},
		{"king cannot recapture protected rook", "7k/7p/8/8/8/8/8/1B4KR w - - 0 1", "h1h7", 100},
		{"legal king recapture", "7k/7p/8/8/8/8/8/6KR w - - 0 1", "h1h7", -400},
		{"en passant rook recapture", "3r2k1/8/8/3pP3/8/8/8/6K1 w - d6 0 1", "e5d6", 0},
		{"en passant diagonal xray", "6k1/8/8/3pP3/8/8/7b/K7 w - d6 0 1", "e5d6", 0},
		{"capture promotion", "r6k/1P6/8/8/8/8/8/7K w - - 0 1", "b7a8q", 1300},
		{"underpromotion", "r6k/1P6/8/8/8/8/8/7K w - - 0 1", "b7a8n", 700},
		{"quiet promotion", "7k/P7/8/8/8/8/8/K7 w - - 0 1", "a7a8q", 800},
		{"recapturing pawn promotes", "k7/8/8/8/8/7R/6p1/K6r w - - 0 1", "h3h1", -800},
		{"knight defender", "6k1/8/5n2/3p4/8/8/8/3Q2K1 w - - 0 1", "d1d5", -800},
		{"black pawn defender", "6k1/8/8/4p3/3r4/8/8/3Q2K1 w - - 0 1", "d1d4", -400},
		{"white pawn defender", "3q2k1/8/8/3R4/4P3/8/8/6K1 b - - 0 1", "d8d5", -400},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(strings.Fields("fen " + tt.fen))
			if err != nil {
				t.Fatal(err)
			}
			position := game.Position()
			var selected chess.Move
			found := false
			for _, move := range position.ValidMoves() {
				if move.String() == tt.move {
					selected = move
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("fixture move %s is not legal", tt.move)
			}
			before := game.FEN()
			initial := captureValue(position, selected)
			if selected.Promo() != chess.NoPieceType {
				initial += pieceValues[selected.Promo()] - pieceValues[chess.Pawn]
			}
			oracle := initial - exchangeOracle(position.Update(&selected), selected.S2())
			if oracle != tt.want {
				t.Fatalf("fixture material oracle %d, want %d", oracle, tt.want)
			}
			if got := staticExchange(position, selected); got != tt.want {
				t.Fatalf("SEE %d, want %d", got, tt.want)
			}
			if game.FEN() != before || len(game.Moves()) != 0 {
				t.Fatal("SEE changed the root")
			}
		})
	}
}

func TestExchangeCancellation(t *testing.T) {
	position := evalFEN(t, "3r2k1/8/8/3p4/8/8/3R4/3R2K1 w - - 0 1")
	move := killerTestMove(t, position, "d2d5")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, complete := (&searcher{ctx: ctx}).exchange(position, move); complete {
		t.Fatal("cancelled exchange completed")
	}
	start := time.Unix(100, 0)
	calls := 0
	worker := searcher{deadline: start.Add(time.Second), now: func() time.Time {
		calls++
		if calls >= 3 {
			return start.Add(time.Second)
		}
		return start
	}}
	if _, complete := worker.exchange(position, move); complete {
		t.Fatal("exchange ignored deadline between recaptures")
	}
}

func TestExchangeFilterAgainstLegalCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/quiescence-validation.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Cases []struct {
			FEN string `json:"fen"`
		} `json:"cases"`
	}
	if err = json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, c := range fixture.Cases {
		position := evalFEN(t, c.FEN)
		if inCheck(position) {
			continue
		}
		for _, move := range position.ValidMoves() {
			if move.HasTag(chess.Check) || move.Promo() != chess.NoPieceType || isQuietMove(move) || captureValue(position, move) >= pieceValues[position.Board().Piece(move.S1()).Type()] {
				continue
			}
			checked++
			see := staticExchange(position, move)
			oracle := captureValue(position, move) - exchangeOracle(position.Update(&move), move.S2())
			if see < 0 && oracle >= 0 {
				t.Errorf("false-negative material filter: %s %s SEE %d oracle %d", c.FEN, move.String(), see, oracle)
			}
		}
	}
	if checked == 0 {
		t.Fatal("corpus exercised no filter candidates")
	}
	t.Logf("checked %d potential capture filters against full legal exchange oracle", checked)
}
