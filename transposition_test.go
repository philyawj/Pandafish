package main

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

func TestTableIdentityAndReplacement(t *testing.T) {
	table := newTable()
	if got := uintptr(len(table.entries)) * unsafe.Sizeof(tableEntry{}); got > 2<<20 {
		t.Fatalf("table exceeds 2 MiB: %d", got)
	}
	position := chess.StartingPosition()
	table.store(position, 3, 25, -100, 100, 0)
	entry, ok := table.probe(position)
	if !ok || entry.bound != exactBound || entry.score != 25 {
		t.Fatal("missing exact entry")
	}
	for _, fen := range []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
	} {
		if _, ok := table.probe(evalFEN(t, fen)); ok {
			t.Fatal("reused score across distinct position context")
		}
	}
	// Force two different full keys into the same slot.
	table.entries = make([]tableEntry, 1)
	table.store(position, 3, 25, -100, 100, 0)
	other := evalFEN(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	if _, ok := table.probe(other); ok {
		t.Fatal("index collision accepted as a hit")
	}
	table.store(other, 1, 3, -100, 100, 0)
	if _, ok := table.probe(position); !ok {
		t.Fatal("shallower collision replaced current deeper work")
	}
	table.generation++
	table.store(other, 1, 3, -100, 100, 0)
	if _, ok := table.probe(position); ok {
		t.Fatal("old collision still hit after replacement")
	}
	table.clear()
	if _, ok := table.probe(other); ok {
		t.Fatal("clear retained an entry")
	}
}

func TestTableEnPassantAndTransposedPositions(t *testing.T) {
	table := newTable()
	withEP := evalFEN(t, "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1")
	withoutEP := evalFEN(t, "4k3/8/8/3pP3/8/8/8/4K3 w - - 0 1")
	table.store(withEP, 2, 100, -200, 200, 0)
	if _, ok := table.probe(withoutEP); ok {
		t.Fatal("en passant rights ignored")
	}
	a, err := parsePosition(strings.Fields("startpos moves g1f3 g8f6 b1c3 b8c6"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := parsePosition(strings.Fields("startpos moves b1c3 b8c6 g1f3 g8f6"))
	if err != nil {
		t.Fatal(err)
	}
	table.store(a.Position(), 2, 5, -100, 100, 0)
	if entry, ok := table.probe(b.Position()); !ok || entry.score != 5 {
		t.Fatal("transposed move orders did not reuse the position")
	}
}

func TestCachedSearchMatchesUncachedWindows(t *testing.T) {
	for _, command := range []string{"startpos", "fen 4k3/1q6/8/5N2/8/8/8/K7 w - - 0 1", "fen k7/8/8/8/5n2/8/1Q6/4K3 b - - 0 1", "fen 7k/8/5KQ1/8/8/8/8/8 w - - 0 1"} {
		game, err := parsePosition(strings.Fields(command))
		if err != nil {
			t.Fatal(err)
		}
		p := game.Position()
		table := newTable()
		worker := searcher{table: table}
		for depth := 3; depth >= 1; depth-- {
			want := search(p, depth, -searchInfinity, searchInfinity)
			for _, window := range [][2]int{{want - 1, want}, {want, want + 1}, {-searchInfinity, searchInfinity}, {-searchInfinity, searchInfinity}} {
				got, complete := worker.search(p, depth, window[0], window[1], inCheck(p))
				if !complete || (want <= window[0] && (got < want || got > window[0])) || (want >= window[1] && (got > want || got < window[1])) || (window[0] < want && want < window[1] && got != want) {
					t.Fatalf("%s depth %d window %v: cached %d, want %d", command, depth, window, got, want)
				}
			}
		}
	}
}

func TestCacheReuseReducesNodesAndPreservesRootTies(t *testing.T) {
	p := chess.StartingPosition()
	moves := rootMoves(p, nil)
	baseline := searcher{}
	best, score, _ := baseline.root(p, moves, 3)
	worker := searcher{table: newTable()}
	for i := 0; i < 2; i++ {
		worker.nodes = 0
		got, value, complete := worker.root(p, moves, 3)
		slices.Sort(got)
		slices.Sort(best)
		if !complete || value != score || !slices.Equal(got, best) {
			t.Fatal("cache changed root ties or score")
		}
		t.Logf("root nodes: uncached=%d cached_pass_%d=%d", baseline.nodes, i+1, worker.nodes)
		if i == 1 && worker.nodes >= baseline.nodes {
			t.Fatalf("warm cache nodes %d, uncached %d", worker.nodes, baseline.nodes)
		}
	}
	// An illegal hint cannot introduce a move; promotions retain their identity.
	before := append([]chess.Move(nil), moves...)
	preferMove(moves, 65535)
	for i := range moves {
		if moves[i].String() != before[i].String() {
			t.Fatal("invalid hint changed moves")
		}
	}
	promo := evalFEN(t, "7k/P7/8/8/8/8/8/4K3 w - - 0 1")
	seen := map[uint16]bool{}
	for _, move := range promo.ValidMoves() {
		key := encodeMove(move)
		if seen[key] {
			t.Fatal("move encoding collision")
		}
		seen[key] = true
	}
}

func TestInterruptedNodeIsNotCached(t *testing.T) {
	p := chess.StartingPosition()
	table := newTable()
	start := time.Unix(100, 0)
	calls := 0
	worker := searcher{table: table, deadline: start.Add(time.Second), now: func() time.Time {
		calls++
		if calls > 12 {
			return start.Add(time.Second)
		}
		return start
	}}
	if _, complete := worker.search(p, 3, -searchInfinity, searchInfinity, false); complete {
		t.Fatal("expected interruption")
	}
	if _, ok := table.probe(p); ok {
		t.Fatal("incomplete parent cached")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	worker = searcher{ctx: ctx, table: table}
	if _, complete := worker.search(p, 3, -searchInfinity, searchInfinity, false); complete {
		t.Fatal("cancelled search completed")
	}
	if _, ok := table.probe(p); ok {
		t.Fatal("cancelled parent cached")
	}
}

func TestTablePersistsAcrossSearchesAndRestrictedRootIsSafe(t *testing.T) {
	p := chess.StartingPosition()
	table := newTable()
	for _, restrict := range [][]string{{"a2a3"}, nil, {"e2e4"}, nil} {
		limits := goLimits{depth: 3, searchMoves: restrict}
		result := iterativeSearchWithTable(context.Background(), p, limits, timeBudget{}, time.Now(), table)
		best, score, _ := (&searcher{}).root(p, rootMoves(p, restrict), 3)
		if result.score != score || !slices.Contains(best, result.move) {
			t.Fatal("restricted search poisoned cache")
		}
	}
	if table.generation != 4 {
		t.Fatal("search generation not advanced")
	}
}

func TestUCICacheLifecycle(t *testing.T) {
	session := newUCISession(t)
	probe := func() uint64 {
		session.send("position startpos")
		session.send("go depth 3")
		var metrics telemetry.Metrics
		for {
			select {
			case line := <-session.lines:
				if strings.HasPrefix(line, telemetry.Prefix) {
					if err := json.Unmarshal([]byte(strings.TrimPrefix(line, telemetry.Prefix)), &metrics); err != nil {
						t.Fatal(err)
					}
				}
				if strings.HasPrefix(line, "bestmove ") {
					if metrics.Depth != 3 || metrics.Nodes == 0 {
						t.Fatalf("invalid completed metrics: %+v", metrics)
					}
					return metrics.Nodes
				}
			case <-time.After(3 * time.Second):
				t.Fatal("UCI cache probe timed out")
			}
		}
	}
	cold, warm := probe(), probe()
	session.send("ucinewgame")
	fresh := probe()
	if warm >= cold || fresh <= warm {
		t.Fatalf("cache lifecycle nodes: first %d, reused %d, new game %d", cold, warm, fresh)
	}
	t.Logf("UCI nodes: first=%d reused=%d after_new_game=%d", cold, warm, fresh)
}
