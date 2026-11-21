package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

type candidateProbeCase struct {
	ID, Position, FEN string
	Deep              bool
}

func candidateProbeCases(t testing.TB) []candidateProbeCase {
	t.Helper()
	b, err := os.ReadFile(os.Getenv("CANDIDATE_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var input struct{ Cases []candidateProbeCase }
	if err := json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}
	return input.Cases
}

func candidateProbeGame(t testing.TB, c candidateProbeCase) *chess.Game {
	t.Helper()
	text := strings.TrimPrefix(c.Position, "position ")
	if text == "" {
		text = "fen " + c.FEN
	}
	g, err := parsePosition(strings.Fields(text))
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestCandidateFixedWork(t *testing.T) {
	output := os.Getenv("CANDIDATE_OUTPUT")
	if output == "" {
		t.Skip("explicit diagnostic input/output required")
	}
	depth := 4
	if raw := os.Getenv("CANDIDATE_DEPTH"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			t.Fatal(err)
		}
		depth = n
	}
	type result struct {
		ID       string
		Depth    int
		Nodes    uint64
		Score    int
		Best     []string
		MS       float64
		Complete bool
	}
	rows := []result{}
	for _, c := range candidateProbeCases(t) {
		if depth > 4 && !c.Deep {
			continue
		}
		g := candidateProbeGame(t, c)
		p := g.Position()
		w := searcher{table: newTable(), killers: new(killerTable), deadline: time.Now().Add(30 * time.Second)}
		for _, k := range repetitionHistory(g) {
			w.pushKey(k)
		}
		moves := rootMoves(p, nil)
		start := time.Now()
		best, score, complete := w.root(p, moves, depth)
		rows = append(rows, result{c.ID, depth, w.nodes, score, best, float64(time.Since(start)) / float64(time.Millisecond), complete})
		data, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(output, append(data, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
		if !complete {
			t.Fatal("incomplete fixed work", c.ID)
		}
	}
}

func BenchmarkCandidateProfile(b *testing.B) {
	if os.Getenv("CANDIDATE_INPUT") == "" {
		b.Skip("explicit input required")
	}
	cases := candidateProbeCases(b)
	games := make([]*chess.Game, len(cases))
	for i, c := range cases {
		games[i] = candidateProbeGame(b, c)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := games[i%len(games)]
		w := searcher{table: newTable(), killers: new(killerTable)}
		for _, k := range repetitionHistory(g) {
			w.pushKey(k)
		}
		_, _, ok := w.root(g.Position(), rootMoves(g.Position(), nil), 4)
		if !ok {
			b.Fatal("interrupted")
		}
	}
}
