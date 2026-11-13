package main

import (
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

func TestQuiescence(t *testing.T) {
	for _, tt := range []struct {
		name, fen string
		want      int
	}{
		{"quiet position", "7k/8/8/8/8/8/P7/7K w - - 0 1", 1},
		{"white recaptures queen", "6k1/8/8/8/3q4/8/8/3R2K1 w - - 0 1", 5},
		{"black recaptures queen", "3r2k1/8/8/3Q4/8/8/8/6K1 b - - 0 1", 5},
		{"white avoids poisoned pawn", "3r2k1/8/8/3p4/8/8/8/3Q2K1 w - - 0 1", 3},
		{"black avoids poisoned pawn", "3q2k1/8/8/8/3P4/8/8/3R2K1 b - - 0 1", 3},
		{"capture and recapture", "r6k/q7/8/8/8/8/8/R6K w - - 0 1", -5},
		{"en passant", "k7/8/8/3pP3/8/8/8/7K w - d6 0 1", 1},
		{"quiet white promotion", "7k/P7/8/8/8/8/8/K7 w - - 0 1", 9},
		{"quiet black promotion", "k7/8/8/8/8/8/p7/7K b - - 0 1", 9},
		{"capture promotion", "r6k/1P6/8/8/8/8/8/7K w - - 0 1", 9},
		{"rook promotion avoids stalemate", "8/k1P5/2K5/8/8/8/8/8 w - - 0 1", 5},
		{"quiet king evasions lose queen", "k6r/p7/8/8/8/8/1b6/Q6K w - - 0 1", -9},
		{"only escape is quiet block", "7r/8/8/8/8/6k1/8/6BK w - - 0 1", -5},
		{"mate", "7k/6Q1/5K2/8/8/8/8/8 b - - 0 1", -mateScore},
		{"stalemate", "7k/5K2/6Q1/8/8/8/8/8 b - - 0 1", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			game, err := parsePosition(strings.Fields("fen " + tt.fen))
			if err != nil {
				t.Fatal(err)
			}
			before := game.FEN()
			want := tt.want
			if want != mateScore && want != -mateScore {
				want *= 100
			}
			if got := quiescence(game.Position(), quiescenceDepth, -searchInfinity, searchInfinity); (game.Position().Status() != chess.NoMethod && got != want) || (game.Position().Status() == chess.NoMethod && abs(got-want) > 90) {
				t.Fatalf("quiescence score = %d, want %d", got, tt.want)
			}
			if game.FEN() != before || len(game.Moves()) != 0 {
				t.Fatal("quiescence changed the root game")
			}
		})
	}
}

// Four tactical plies keep the dense exhaustive oracle bounded. Production
// search uses the larger safety allowance with alpha-beta cutoffs.
func TestQuiescenceMatchesExhaustiveOnTimeoutPositions(t *testing.T) {
	for _, fen := range []string{timeoutPosition, queenRaidPosition, crowdedPosition} {
		game, err := parsePosition(strings.Fields("fen " + fen))
		if err != nil {
			t.Fatal(err)
		}
		want := exhaustiveQuietScore(game.Position(), 4)
		if game.Position().Turn() == chess.Black {
			want = -want
		}
		if got := quiescence(game.Position(), 4, -searchInfinity, searchInfinity); got != want {
			t.Errorf("%s: quiescence score = %d, exhaustive score = %d", fen, got, want)
		}
	}
}

func TestQuiescenceCannotStandPatInCheck(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen k6r/p7/8/8/8/8/1b6/Q6K w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	// Static material is zero, but every legal escape lets the bishop take
	// our queen. A stand-pat cutoff at beta = 0 would hide that loss.
	if got := quiescence(game.Position(), quiescenceDepth, -searchInfinity, 0); got >= -800 {
		t.Fatalf("checked score = %d, want a lost queen", got)
	}
}

func TestQuiescenceDepthCapOutsideCheck(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 3r2k1/8/8/8/3Q4/8/8/6K1 b - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	if inCheck(game.Position()) {
		t.Fatal("fixture must not be checked")
	}
	// At the ordinary horizon we return static positional evaluation even if captures remain.
	if got := quiescence(game.Position(), 0, -searchInfinity, searchInfinity); got != evaluate(game.Position(), game.Position().Turn()) {
		t.Fatalf("capped score = %d, want static evaluation", got)
	}
}

func TestInCheck(t *testing.T) {
	for _, tt := range []struct {
		fen  string
		want bool
	}{
		{"7k/8/8/8/8/8/P7/7K w - - 0 1", false},
		{"k6r/8/8/8/8/8/8/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/5n2/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/6p1/7K w - - 0 1", true},
		{"k7/8/8/8/8/8/1b6/K7 w - - 0 1", true},
		{"7k/8/8/8/8/8/8/K6R b - - 0 1", true},
	} {
		game, err := parsePosition(strings.Fields("fen " + tt.fen))
		if err != nil {
			t.Fatal(err)
		}
		if got := inCheck(game.Position()); got != tt.want {
			t.Errorf("%s: check = %v, want %v", tt.fen, got, tt.want)
		}
	}
}

// Both queens give check at the ordinary quiescence horizon. Capturing the
// checker wins it, but the old cutoff returns the pre-capture material score.
func TestQuiescenceHorizonCheckEvasions(t *testing.T) {
	for _, fen := range []string{
		"6k1/8/8/8/3q4/8/8/3R2K1 w - - 0 1",
		"3r2k1/8/8/3Q4/8/8/8/6K1 b - - 0 1",
	} {
		game, err := parsePosition(strings.Fields("fen " + fen))
		if err != nil {
			t.Fatal(err)
		}
		position := game.Position()
		if !inCheck(position) {
			t.Fatal("fixture must be checked")
		}
		before := game.FEN()
		got := quiescence(position, 0, -searchInfinity, searchInfinity)
		if abs(got-500) > 90 {
			t.Fatalf("%s: horizon score %d, want the queen captured and a rook remaining", fen, got)
		}
		if game.FEN() != before {
			t.Fatal("search changed root")
		}
	}
}

func TestQuiescenceHorizonQuietEvasion(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 7r/8/8/8/8/6k1/8/6BK w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	position := game.Position()
	moves := position.ValidMoves()
	if len(moves) != 1 || moves[0].String() != "g1h2" {
		t.Fatal("fixture must force Bh2")
	}
	// The only evasion is a quiet bishop block. The score must describe the
	// resulting line, rather than treating the checked position as quiet.
	// Bh2 also gives check; Black can answer Kxh2, winning the bishop.
	want := exhaustiveQuietScore(position, 0)
	if got := quiescence(position, 0, -searchInfinity, searchInfinity); got != want || abs(got+500) > 90 {
		t.Fatalf("score %d, want forced quiet evasion score %d", got, want)
	}
}

func TestDeadlineDuringHorizonCheckEvasion(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 3r2k1/8/8/3Q4/8/8/8/6K1 b - - 0 1"))
	if err != nil {
		t.Fatal(err)
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
	if _, complete := worker.quiet(game.Position(), 0, -searchInfinity, searchInfinity, true); complete {
		t.Fatal("horizon evasion ignored its deadline")
	}
	if worker.nodes != 1 {
		t.Fatalf("visited %d nodes after deadline", worker.nodes)
	}
}

func TestQuiescenceHorizonCountercheck(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 3k4/8/8/8/3q4/8/8/3R2K1 w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	position := game.Position()
	var capture *chess.Position
	for _, move := range position.ValidMoves() {
		if move.String() == "d1d4" {
			capture = position.Update(&move)
		}
	}
	if capture == nil || !inCheck(capture) {
		t.Fatal("Rxd4 must give countercheck")
	}
	// Resolve the countercheck as well; evaluating immediately after Rxd4+
	// would leave the checked king's positional score unresolved.
	want := exhaustiveQuietScore(position, 0)
	worker := searcher{}
	got, complete := worker.quiet(position, 0, -searchInfinity, searchInfinity, true)
	if !complete || got != want || abs(got-500) > 90 || worker.nodes < 3 {
		t.Fatalf("countercheck: score %d, want %d, nodes %d, complete %v", got, want, worker.nodes, complete)
	}
}

func TestQuiescenceCheckEvasionSafetyCap(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 6k1/8/8/8/3q4/8/8/3R2K1 w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	worker := searcher{}
	position := game.Position()
	got, complete := worker.quiet(position, -checkEvasionDepth, -searchInfinity, searchInfinity, true)
	if !complete || got != staticScore(position, position.Turn()) || worker.nodes != 1 {
		t.Fatalf("safety cap: score %d, complete %v, nodes %d", got, complete, worker.nodes)
	}
	// Immediately before the cap, a forced evasion still receives one ply.
	got, complete = (&searcher{}).quiet(position, 1-checkEvasionDepth, -searchInfinity, searchInfinity, true)
	if !complete || abs(got-500) > 90 {
		t.Fatalf("last evasion: score %d, complete %v", got, complete)
	}
}

func TestQuiescenceResolvesCheckAfterCaptureAndRecapture(t *testing.T) {
	game, err := parsePosition(strings.Fields("fen 4k3/8/3q4/3p4/2B5/3R4/6K1/8 w - - 0 1"))
	if err != nil {
		t.Fatal(err)
	}
	// Either piece can take d5. Qxd5+ then reaches the ordinary two-ply
	// horizon, but the other white piece can take the checking queen.
	// Black must decline that exchange, so White can win the pawn.
	position := game.Position()
	got := quiescence(position, quiescenceDepth, -searchInfinity, searchInfinity)
	want := exhaustiveQuietScore(position, quiescenceDepth)
	if got != want || got <= staticScore(position, position.Turn())+50 {
		t.Fatalf("score %d, want safe pawn capture, reference %d, static %d", got, want, staticScore(position, position.Turn()))
	}
}

// The v0.12 two-ply cutoff rejected taking this pawn because it stopped before
// the final recapture. All moves in the example are ordinary non-checking captures.
const exchangeHorizonPosition = "4k3/8/3q4/3p4/2B5/3R3K/8/8 w - - 0 1"

func TestQuiescenceContinuesOrdinaryExchanges(t *testing.T) {
	for _, fen := range []string{
		exchangeHorizonPosition,
		"8/8/3r3k/2b5/3P4/3Q4/8/4K3 b - - 0 1",
	} {
		game, err := parsePosition(strings.Fields("fen " + fen))
		if err != nil {
			t.Fatal(err)
		}
		position := game.Position()
		if inCheck(position) {
			t.Fatal("fixture must not be checked")
		}
		before := game.FEN()
		short, oldComplete := (&searcher{unfiltered: true}).quiet(position, 2, -searchInfinity, searchInfinity, false)
		worker := searcher{}
		got, complete := worker.quiet(position, quiescenceDepth, -searchInfinity, searchInfinity, false)
		if !oldComplete || !complete || got <= short+50 || abs(got+100) > 90 {
			t.Fatalf("exchange score %d, old horizon %d, complete %v", got, short, complete)
		}
		if game.FEN() != before {
			t.Fatal("exchange changed root")
		}
	}
}

func TestExchangeFilterKeepsTacticalExceptions(t *testing.T) {
	// Full and filtered continuation must retain checking sacrifices, quiet
	// promotions, underpromotions, and every evasion, even when SEE is negative.
	for _, fen := range []string{
		"3r2k1/8/8/3p4/8/8/8/3Q2K1 w - - 0 1",
		"k6r/p7/8/8/8/8/1b6/Q6K w - - 0 1",
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
		"8/k1P5/2K5/8/8/8/8/8 w - - 0 1",
	} {
		position := evalFEN(t, fen)
		got, ok := (&searcher{}).quiet(position, 4, -searchInfinity, searchInfinity, inCheck(position))
		want, complete := (&searcher{unfiltered: true}).quiet(position, 4, -searchInfinity, searchInfinity, inCheck(position))
		if !ok || !complete || got != want {
			t.Fatalf("%s: filtered %d, full %d", fen, got, want)
		}
	}
}
