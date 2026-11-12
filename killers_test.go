package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

func killerTestMove(t *testing.T, position *chess.Position, uci string) chess.Move {
	t.Helper()
	for _, move := range position.ValidMoves() {
		if move.String() == uci {
			return move
		}
	}
	t.Fatalf("missing legal fixture move %s", uci)
	return chess.Move{}
}

const killerOrderPosition = "r3k2r/1P3ppp/8/3q4/8/8/PP3PPP/R3K2R w KQkq - 0 1"

func TestKillerReplacementAndQuietOnly(t *testing.T) {
	position := evalFEN(t, killerOrderPosition)
	a, b := killerTestMove(t, position, "a2a3"), killerTestMove(t, position, "h2h3")
	table := new(killerTable)
	table.remember(3, a)
	table.remember(3, b)
	table.remember(3, a)
	want := [2]uint16{encodeMove(a), encodeMove(b)}
	if table[3] != want {
		t.Fatalf("replacement: %v, want %v", table[3], want)
	}
	table.remember(3, a) // Do not displace the second slot with a duplicate.
	for _, uci := range []string{"b7a8q", "b7b8n"} {
		table.remember(3, killerTestMove(t, position, uci))
	}
	ep := evalFEN(t, "k7/8/8/3pP3/8/8/8/7K w - d6 0 1")
	table.remember(3, killerTestMove(t, ep, "e5d6"))
	if table[3] != want {
		t.Fatal("capture, promotion, or en passant replaced a quiet hint")
	}
	before := *table
	table.remember(-1, a)
	table.remember(maxSearchDepth, a)
	if *table != before || table[2] != [2]uint16{} {
		t.Fatal("hint escaped its search ply")
	}
	(*killerTable)(nil).remember(0, a)
}

func TestKillerOrderingPreservesTacticsAndHashPriority(t *testing.T) {
	position := evalFEN(t, killerOrderPosition)
	moves := orderedMoves(position)
	original := slices.Clone(moves)
	table := new(killerTable)
	table.remember(0, killerTestMove(t, position, "a2a3"))
	table.remember(0, killerTestMove(t, position, "h2h3"))
	table.order(moves, 0)
	quiet := 0
	for quiet < len(original) && !isQuietMove(original[quiet]) {
		if original[quiet].String() != moves[quiet].String() {
			t.Fatal("killer displaced a capture or promotion")
		}
		quiet++
	}
	if quiet == 0 || moves[quiet].String() != "h2h3" || moves[quiet+1].String() != "a2a3" {
		t.Fatal("quiet hints were not placed after the tactical prefix")
	}
	// A stale quiet hint can now represent a capture or promotion. It must
	// not disturb tactical ordering, add a move, or change move metadata.
	stale := new(killerTable)
	stale[0] = [2]uint16{encodeMove(killerTestMove(t, position, "b7a8q")), 65535}
	unchanged := slices.Clone(original)
	stale.order(unchanged, 0)
	if !reflect.DeepEqual(original, unchanged) {
		t.Fatal("stale hints altered legal/tactical moves")
	}
	for _, ply := range []int{-1, maxSearchDepth} {
		table.order(unchanged, ply)
	}
	(*killerTable)(nil).order(unchanged, 0)
	if !reflect.DeepEqual(original, unchanged) {
		t.Fatal("disabled/out-of-range hints changed ordering")
	}
	preferMove(moves, encodeMove(killerTestMove(t, position, "e1g1")))
	if moves[0].String() != "e1g1" {
		t.Fatal("hash move did not retain first priority")
	}
	before, after := make([]uint16, len(original)), make([]uint16, len(moves))
	for i := range moves {
		before[i], after[i] = encodeMove(original[i]), encodeMove(moves[i])
	}
	slices.Sort(before)
	slices.Sort(after)
	if !slices.Equal(before, after) {
		t.Fatal("ordering changed the set of legal moves")
	}
}

func TestKillerScoresAndRootTies(t *testing.T) {
	for _, command := range []string{"startpos", "fen " + queenRaidPosition, "fen " + timeoutPosition,
		"fen 7k/8/5KQ1/8/8/8/8/1q6 w - - 0 1"} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		position := game.Position()
		moves := rootMoves(position, nil)
		baseline := searcher{table: newTable()}
		ordered := searcher{table: newTable(), killers: new(killerTable)}
		for depth := 1; depth <= 3; depth++ {
			wantMoves, wantScore, ok := baseline.root(position, slices.Clone(moves), depth)
			gotMoves, gotScore, complete := ordered.root(position, slices.Clone(moves), depth)
			slices.Sort(wantMoves)
			slices.Sort(gotMoves)
			if !ok || !complete || gotScore != wantScore || !slices.Equal(gotMoves, wantMoves) {
				t.Fatalf("%s depth %d: %v/%d, want %v/%d", command, depth, gotMoves, gotScore, wantMoves, wantScore)
			}
		}
		for _, window := range [][2]int{{-1, 0}, {0, 1}, {-searchInfinity, searchInfinity}} {
			want, _ := (&searcher{}).search(position, 2, -searchInfinity, searchInfinity, inCheck(position))
			got, complete := ordered.search(position, 2, window[0], window[1], inCheck(position))
			if !complete || (want <= window[0] && (got < want || got > window[0])) ||
				(want >= window[1] && (got > want || got < window[1])) ||
				(want > window[0] && want < window[1] && got != want) {
				t.Fatalf("window %v score %d, reference %d", window, got, want)
			}
		}
	}
}

func TestCancelledSearchDoesNotLearnKillers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := new(killerTable)
	worker := searcher{ctx: ctx, killers: table}
	if _, complete := worker.search(chess.StartingPosition(), 4, -searchInfinity, searchInfinity, false); complete {
		t.Fatal("cancelled search completed")
	}
	if *table != (killerTable{}) {
		t.Fatal("cancelled work added hints")
	}
}

func TestKillerNodeAudit(t *testing.T) {
	type nodeResult struct {
		Position string `json:"position"`
		Score    int    `json:"score"`
		Baseline uint64 `json:"baseline_nodes"`
		Ordered  uint64 `json:"ordered_nodes"`
	}
	rows := []nodeResult{}
	for _, command := range []string{"startpos", "fen " + timeoutPosition, "fen " + queenRaidPosition, "fen " + crowdedPosition} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		baseline := searcher{table: newTable()}
		ordered := searcher{table: newTable(), killers: new(killerTable)}
		want, score, ok := baseline.root(game.Position(), rootMoves(game.Position(), nil), 4)
		got, gotScore, complete := ordered.root(game.Position(), rootMoves(game.Position(), nil), 4)
		slices.Sort(want)
		slices.Sort(got)
		if !ok || !complete || score != gotScore || !slices.Equal(want, got) {
			t.Fatal("ordering changed exact root result")
		}
		if command == "startpos" && ordered.nodes >= baseline.nodes {
			t.Fatalf("startpos nodes %d, baseline %d", ordered.nodes, baseline.nodes)
		}
		t.Logf("%s: nodes %d -> %d; score %d", command, baseline.nodes, ordered.nodes, score)
		rows = append(rows, nodeResult{command, score, baseline.nodes, ordered.nodes})
	}
	// Optional artifact for the local experiment; routine tests write nothing.
	if path := os.Getenv("PANDAFISH_NODE_AUDIT"); path != "" {
		b, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, append(b, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
}
