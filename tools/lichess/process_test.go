package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

func fakeStockfish() {
	position := "position startpos"
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if path := os.Getenv("PANDAFISH_TEST_SF_LOG"); path != "" {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				os.Exit(2)
			}
			fmt.Fprintln(f, line)
			f.Close()
		}
		switch {
		case line == "uci":
			fmt.Println("id name Stockfish test")
			for _, name := range []string{"Threads", "Hash", "MultiPV", "Skill Level", "UCI_LimitStrength", "Ponder"} {
				fmt.Println("option name " + name + " type spin default 1")
			}
			fmt.Println("uciok")
		case line == "isready":
			fmt.Println("readyok")
		case strings.HasPrefix(line, "position "):
			position = line
		case strings.HasPrefix(line, "go "):
			game, err := positionGame(position)
			if err != nil {
				os.Exit(3)
			}
			moves := game.ValidMoves()
			if len(moves) == 0 {
				os.Exit(4)
			}
			best := moves[0].String()
			score := 100
			if _, restricted, ok := strings.Cut(line, " searchmoves "); ok {
				best = restricted
				score = -100
			}
			if os.Getenv("PANDAFISH_TEST_SF_BOUND") == "1" && !strings.Contains(line, "searchmoves") && len(moves) > 1 {
				other := moves[1].String()
				fmt.Printf("info depth 7 nodes 50 time 1 score cp 20 pv %s\ninfo depth 8 nodes 100 time 2 score cp 30 lowerbound pv %s\nbestmove %s\n", other, best, best)
			} else {
				fmt.Printf("info depth 8 nodes 100 time 1 score cp %d pv %s\nbestmove %s\n", score, best, best)
			}
		}
	}
}
func processFixture(t *testing.T) (processOptions, string) {
	t.Helper()
	root := t.TempDir()
	session := filepath.Join(root, "session")
	if err := os.MkdirAll(filepath.Join(session, "games"), 0755); err != nil {
		t.Fatal(err)
	}
	fixturePGN(t, filepath.Join(session, "games"), "game0001", "1. e4 e5", "1-0")
	engine := []byte("recorded engine fixture")
	config := []byte("token: REDACTED\n")
	for name, data := range map[string][]byte{"pandafish": engine, "config.redacted.yml": config, "lichess-bot.log": []byte("2026-01-01T00:00:00Z INFO ordinary log\n")} {
		if err := os.WriteFile(filepath.Join(session, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now().UTC()
	m := manifest{Schema: 1, Status: "finished", Started: now, Finished: &now, EngineSHA256: digest(engine), ConfigSHA256: digest(config)}
	if err := writeJSON(filepath.Join(session, "manifest.json"), m); err != nil {
		t.Fatal(err)
	}
	ply := 0
	row := searchRecord{Game: "https://lichess.org/game0001", Ply: &ply, Process: "engine-test", Status: "completed", FEN: chess.StartingPosition().String(), Position: "position startpos", Go: "go movetime 10000", Metrics: &telemetry.Metrics{Depth: 2, Timed: true, HardMS: 10, ElapsedMS: 9}}
	b, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(session, "moves.jsonl"), append(b, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PANDAFISH_TEST_STOCKFISH", "1")
	log := filepath.Join(root, "stockfish-commands.txt")
	t.Setenv("PANDAFISH_TEST_SF_LOG", log)
	return processOptions{Session: session, Archive: filepath.Join(root, "archive"), Stockfish: binary, Bot: "Pandafish", Nodes: 100, Limit: 1}, log
}
func TestProcessArchiveResumeAndSettings(t *testing.T) {
	o, log := processFixture(t)
	before, err := inventory(o.Session, true)
	if err != nil {
		t.Fatal(err)
	}
	if err = runProcess(o); err != nil {
		t.Fatal(err)
	}
	after, err := inventory(o.Session, true)
	if err != nil {
		t.Fatal(err)
	}
	b1, _ := json.Marshal(before)
	b2, _ := json.Marshal(after)
	if string(b1) != string(b2) {
		t.Fatal("process changed raw source")
	}
	o.Limit = 0
	if err = runProcess(o); err != nil {
		t.Fatal(err)
	}
	commands, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	searches := strings.Count(string(commands), "go nodes")
	if err = runProcess(o); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(again), "go nodes") != searches {
		t.Fatal("resuming repeated cached analysis")
	}
	raw, err := os.ReadFile(filepath.Join(o.Archive, "analysis/stockfish-n100/summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var summary map[string]any
	if err = json.Unmarshal(raw, &summary); err != nil {
		t.Fatal(err)
	}
	if summary["processed_plies"] != float64(2) || summary["expected_plies"] != float64(2) || summary["bot_moves"] != float64(1) {
		t.Fatal(string(raw))
	}
	o.Bot = "Other"
	if err = runProcess(o); err == nil {
		t.Fatal("mixed settings into existing analysis")
	}
	if err = os.WriteFile(filepath.Join(o.Archive, "raw", "lichess-bot.log"), []byte("tampered"), 0644); err != nil {
		t.Fatal(err)
	}
	if err = archiveSession(o.Session, o.Archive); err == nil {
		t.Fatal("accepted modified raw archive")
	}
}
func TestArchiveRefusesCredentialAndSymlink(t *testing.T) {
	o, _ := processFixture(t)
	if err := os.WriteFile(filepath.Join(o.Session, "extra.yml"), []byte("token: real-private-value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := archiveSession(o.Session, o.Archive); err == nil {
		t.Fatal("archived unredacted config")
	}
	if _, err := os.Stat(o.Archive); !os.IsNotExist(err) {
		t.Fatal("left an archive after failed scan")
	}
	if err := os.Remove(filepath.Join(o.Session, "extra.yml")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("manifest.json", filepath.Join(o.Session, "link")); err != nil {
		t.Fatal(err)
	}
	if err := archiveSession(o.Session, o.Archive); err == nil {
		t.Fatal("accepted symlink")
	}
}
func TestSFParsingAndMateComparison(t *testing.T) {
	result, ok := parseSFInfo("info depth 12 score cp -42 nodes 1234 time 90 pv e2e4 e7e5")
	if !ok || result.Score.Value != -42 || result.Nodes != 1234 || len(result.PV) != 2 {
		t.Fatal(result)
	}
	if _, ok = parseSFInfo("info depth 12 score cp 42 lowerbound nodes 100 pv e2e4"); ok {
		t.Fatal("accepted bound as exact")
	}
	if _, ok = parseSFInfo("info depth 12 nodes 100"); ok {
		t.Fatal("accepted missing score")
	}
	for _, c := range []struct {
		best, played sfScore
		want         string
	}{{sfScore{"cp", 50}, sfScore{"cp", -300}, "loss_300cp_plus"}, {sfScore{"cp", 0}, sfScore{"cp", 10}, "search_disagreement"}, {sfScore{"mate", 3}, sfScore{"cp", 600}, "missed_forced_mate"}, {sfScore{"cp", 10}, sfScore{"mate", -2}, "allowed_forced_mate"}, {sfScore{"mate", -2}, sfScore{"mate", -1}, "mate_score_comparison"}} {
		loss, _, category := compareScores(c.best, c.played)
		if category != c.want {
			t.Fatalf("got %s want %s", category, c.want)
		}
		if c.best.Kind == "mate" || c.played.Kind == "mate" {
			if loss != nil {
				t.Fatal("invented cp loss for mate")
			}
		}
	}
	s := stockfish{lines: make(chan string), cancel: func() {}}
	if _, err := s.until("readyok", time.Millisecond); err == nil {
		t.Fatal("wait did not time out")
	}
}
func TestWrappedClientErrorClassification(t *testing.T) {
	path := filepath.Join(t.TempDir(), "client.log")
	text := "2026-01-01T00:00:00Z ERROR {'error': 'Bot played 100 matchmaking.py:129\n2026-01-01T00:00:00Z games today', 'ratelimit': {'key': 'bot.vsBot.day'},\n2026-01-01T00:00:00Z 'opponent_is_rate_limited': True, 'bot_is_rate_limited': False}\n2026-01-01T00:00:00Z INFO continuing normally\n"
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := classifyClientLog(path)
	if err != nil || len(got) != 1 || got[0].Kind != "opponent_daily_game_cap" {
		t.Fatal(got, err)
	}
}

func TestCloneArchiveWithoutExecutable(t *testing.T) {
	o, _ := processFixture(t)
	o.Limit = 0
	if err := runProcess(o); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(o.Archive, "raw", "pandafish")); err != nil {
		t.Fatal(err)
	}
	o.Session = filepath.Join(o.Archive, "raw")
	if err := runProcess(o); err != nil {
		t.Fatal("clone could not reprocess data without a compiled engine:", err)
	}
	if err := os.Remove(filepath.Join(o.Archive, "raw", "moves.jsonl")); err != nil {
		t.Fatal(err)
	}
	if err := archiveSession(o.Session, o.Archive); err == nil {
		t.Fatal("missing data file accepted")
	}
}
func TestRecordingHardwareAndAutomaticReport(t *testing.T) {
	o := recordingFixture(t)
	o.PrepareOnly = false
	o.Hardware = "Apple M2 MacBook Air"
	o.Workload = "browser and stream"
	t.Setenv("PANDAFISH_TEST_CLIENT", "1")
	t.Setenv("PANDAFISH_TEST_MODE", "success")
	if err := record(o); err != nil {
		t.Fatal(err)
	}
	m := readManifest(t, o.Session)
	if m.Hardware != o.Hardware || m.Workload != o.Workload {
		t.Fatal("hardware metadata not saved")
	}
	for _, name := range []string{"hardware.txt", "summary.md", "summary.json", "moves.jsonl"} {
		if _, err := os.Stat(filepath.Join(o.Session, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestBoundOnlyFinalChoiceIsRescored(t *testing.T) {
	o, _ := processFixture(t)
	t.Setenv("PANDAFISH_TEST_SF_BOUND", "1")
	sf, err := openStockfish(o.Stockfish)
	if err != nil {
		t.Fatal(err)
	}
	defer sf.close()
	result, err := sf.analyse("position startpos", "", 100)
	if err != nil {
		t.Fatal(err)
	}
	if !result.ReferenceRechecked || len(result.SelectionOutput) == 0 || result.PV[0] != result.Bestmove || result.Score.Value != -100 {
		t.Fatalf("unrelated PV score used for final move: %+v", result)
	}
}
