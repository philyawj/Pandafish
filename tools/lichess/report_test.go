package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/philyawj/Pandafish/internal/telemetry"
)

func fixtureTrace(t *testing.T, dir string, lines []telemetry.Event, tail string) string {
	t.Helper()
	path := filepath.Join(dir, "engine-test.jsonl")
	var b strings.Builder
	for i, e := range lines {
		e.Schema = 1
		e.Sequence = uint64(i + 1)
		e.ElapsedNS = int64(i) * int64(time.Millisecond)
		e.Time = time.Unix(100, int64(i)*int64(time.Millisecond)).UTC()
		raw, err := json.Marshal(e)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(raw)
		b.WriteByte('\n')
	}
	b.WriteString(tail)
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
func fixturePGN(t *testing.T, dir, id, moves, result string) {
	t.Helper()
	text := "[Event \"Offline test\"]\n[Site \"https://lichess.org/" + id + "\"]\n[White \"Pandafish\"]\n[Black \"Other\"]\n[Result \"" + result + "\"]\n\n" + moves + " " + result + "\n"
	if err := os.WriteFile(filepath.Join(dir, id+".pgn"), []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
}
func TestReportRecordedSearchAndRepeatability(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"games", "traces"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0755); err != nil {
			t.Fatal(err)
		}
	}
	data := []byte("offline fixture")
	sum := sha256.Sum256(data)
	for _, name := range []string{"pandafish", "config.redacted.yml"} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := writeJSON(filepath.Join(dir, "manifest.json"), manifest{Schema: 1, Status: "finished", EngineID: "test", EngineSHA256: hex.EncodeToString(sum[:]), ConfigSHA256: hex.EncodeToString(sum[:])}); err != nil {
		t.Fatal(err)
	}
	metrics := `{"depth":2,"nodes":125,"elapsed_ms":9,"timed":true,"soft_ms":7,"hard_ms":10,"reason":"soft_limit","fallback":false,"score_kind":"cp","score":100,"pv":["e7e5"]}`
	path := fixtureTrace(t, filepath.Join(dir, "traces"), []telemetry.Event{
		{Kind: "start"}, {Kind: "in", Line: "ucinewgame"}, {Kind: "in", Line: "position startpos moves e2e4"}, {Kind: "in", Line: "go wtime 1000 btime 900 winc 1 binc 2"},
		{Kind: "out", Line: telemetry.Prefix + metrics + "\n"}, {Kind: "out", Line: "bestmove e7e5\n"}, {Kind: "exit"},
	}, "")
	rows, warnings := readTrace(path)
	if len(warnings) != 0 || len(rows) != 1 {
		t.Fatalf("rows %d, warnings %v", len(rows), warnings)
	}
	r := rows[0]
	if r.Bestmove != "e7e5" || r.ResponseMS != 2 || r.Metrics.Score != 100 || r.Metrics.Depth != 2 || r.Status != "completed" || !strings.Contains(r.Go, "btime 900") {
		t.Fatalf("record %+v", r)
	}
	fixturePGN(t, filepath.Join(dir, "games"), "abcd1234", "1. e4 e5 2. Nf3", "1-0")
	if err := report(dir); err != nil {
		t.Fatal(err)
	}
	var summary sessionReport
	raw, err := os.ReadFile(filepath.Join(dir, "summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &summary); err != nil {
		t.Fatal(err)
	}
	if summary.Matched != 1 || summary.Completed != 1 || len(summary.Warnings) != 0 || summary.ResponseP95MS != 2 || summary.DepthCounts[2] != 1 || summary.StopReasons["soft_limit"] != 1 {
		t.Fatalf("summary %+v", summary)
	}
	before, err := os.ReadFile(filepath.Join(dir, "moves.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if err := report(dir); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(dir, "moves.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("report is not repeatable")
	}
	var got searchRecord
	if err = json.Unmarshal([]byte(strings.TrimSpace(string(after))), &got); err != nil {
		t.Fatal(err)
	}
	if got.Game != "https://lichess.org/abcd1234" || got.Ply == nil || *got.Ply != 1 {
		t.Fatalf("link %+v", got)
	}
}
func TestCorrelationAmbiguityHistoryAndProcessEvidence(t *testing.T) {
	dir := t.TempDir()
	fixturePGN(t, dir, "game0001", "1. e4 e5 2. Nf3 Nc6", "1-0")
	fixturePGN(t, dir, "game0002", "1. e4 c5", "*")
	games, warnings := loadGames(dir)
	if len(warnings) != 0 {
		t.Fatal(warnings)
	}
	makeRow := func(position, best string) *searchRecord {
		g, err := positionGame(position)
		if err != nil {
			t.Fatal(err)
		}
		return &searchRecord{Process: "one", Position: position, FEN: g.FEN(), Bestmove: best}
	}
	first := makeRow("position startpos", "e2e4")
	correlate([]*searchRecord{first}, games)
	if first.Correlation != "ambiguous" || len(first.Candidates) != 2 {
		t.Fatalf("guessed ambiguous game: %+v", first)
	}
	later := makeRow("position startpos moves e2e4 e7e5", "g1f3")
	correlate([]*searchRecord{first, later}, games)
	if first.Correlation != "matched" || later.Correlation != "matched" || !strings.HasSuffix(first.Game, "game0001") {
		t.Fatalf("process evidence did not resolve linkage: %+v", first)
	}
	missing := makeRow("position startpos moves d2d4", "d7d5")
	missing.Process = "two"
	correlate([]*searchRecord{missing}, games)
	if missing.Correlation != "missing" {
		t.Fatal("invented missing match")
	}
}
func TestPartialCancellationIllegalMoveAndMalformedTrace(t *testing.T) {
	path := fixtureTrace(t, t.TempDir(), []telemetry.Event{
		{Kind: "start"}, {Kind: "in", Line: "go infinite"}, {Kind: "in", Line: "position startpos moves e2e4"},
		{Kind: "in", Line: "go depth 1"}, {Kind: "out", Line: "bestmove e2e5"}, {Kind: "in", Line: "go movetime 5"},
	}, `{"schema":1`)
	rows, warnings := readTrace(path)
	if len(rows) != 3 || rows[0].Status != "cancelled" || rows[1].Status != "illegal_move" || rows[2].Status != "incomplete" || len(warnings) != 3 {
		t.Fatalf("rows %+v, warnings %v", rows, warnings)
	}
}
func TestPGNMissingAndTruncated(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "broken.pgn"), []byte("[Site \"unfinished"), 0644); err != nil {
		t.Fatal(err)
	}
	games, warnings := loadGames(dir)
	if len(games) != 0 || len(warnings) != 1 {
		t.Fatalf("games %d warnings %v", len(games), warnings)
	}
	if err := report(t.TempDir()); err == nil {
		t.Fatal("accepted missing manifest")
	}
}
func TestUCIPositionFENHistory(t *testing.T) {
	command := "position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves e2e4 e7e5"
	got, err := positionGame(command)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := positionGame("position startpos moves e2e4 e7e5")
	if err != nil {
		t.Fatal(err)
	}
	if got.FEN() != expected.FEN() {
		t.Fatal("FEN history lost")
	}
	for _, command := range []string{"position", "position fen bad w - - 0 1", "position startpos moves e2e5", "position startpos unexpected"} {
		if _, err := positionGame(command); err == nil {
			t.Fatalf("accepted %q", command)
		}
	}
}

func TestDuplicateEventIsFlagged(t *testing.T) {
	path := fixtureTrace(t, t.TempDir(), []telemetry.Event{{Kind: "start"}, {Kind: "in", Line: "go depth 1"}}, "")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.NewEncoder(f).Encode(telemetry.Event{Schema: 1, Sequence: 2, Kind: "out", Line: "bestmove e2e4"}); err != nil {
		t.Fatal(err)
	}
	f.Close()
	rows, warnings := readTrace(path)
	if len(rows) != 1 || rows[0].Status != "incomplete" || len(warnings) != 2 {
		t.Fatalf("rows %+v warnings %v", rows, warnings)
	}
}
