package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
)

// The test executable doubles as a controllable UCI subprocess. Tests therefore
// exercise the real pipe, timeout and process cleanup code without Stockfish.
func TestMain(m *testing.M) {
	if mode := os.Getenv("PANDA_BENCH_TEST_ENGINE"); mode != "" {
		var history []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if trace := os.Getenv("PANDA_BENCH_TRACE"); trace != "" {
				f, err := os.OpenFile(trace, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					os.Exit(5)
				}
				fmt.Fprintln(f, line)
				f.Close()
			}
			switch {
			case line == "uci":
				fmt.Println("id name Test engine")
				if mode == "options" {
					for key := range stockfishOptions(1) {
						fmt.Printf("option name %s type string\n", key)
					}
				}
				fmt.Println("uciok")
			case line == "isready":
				fmt.Println("readyok")
				if mode == "block-input" {
					time.Sleep(time.Hour)
				}
			case strings.HasPrefix(line, "position startpos moves "):
				history = strings.Fields(strings.TrimPrefix(line, "position startpos moves "))
			case strings.HasPrefix(line, "go "):
				if mode == "slow" {
					time.Sleep(80 * time.Millisecond)
				}
				switch mode {
				case "hang":
					time.Sleep(time.Hour)
				case "crash":
					os.Exit(3)
				case "illegal":
					fmt.Println("bestmove e1e8")
				case "missing":
					fmt.Println("bestmove")
				case "mate":
					fmt.Println("bestmove", []string{"f2f3", "e7e5", "g2g4", "d8h4"}[len(history)])
				case "repeat":
					fmt.Println("bestmove", []string{"g1f3", "g8f6", "f3g1", "f6g8"}[len(history)%4])
				default:
					game, err := newGame(history)
					if err != nil {
						os.Exit(4)
					}
					moves := game.ValidMoves()
					fmt.Println("bestmove", (chess.UCINotation{}).Encode(game.Position(), &moves[0]))
				}
			}
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func helperPath(t *testing.T, mode string) string {
	t.Helper()
	t.Setenv("PANDA_BENCH_TEST_ENGINE", mode)
	path, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestUCIFailures(t *testing.T) {
	for _, tt := range []struct{ mode, want string }{{"hang", "timeout"}, {"crash", "exited"}, {"missing", "missing"}} {
		t.Run(tt.mode, func(t *testing.T) {
			e, err := startEngine(helperPath(t, tt.mode), nil)
			if err != nil {
				t.Fatal(err)
			}
			defer e.close()
			_, _, err = e.move([]string{"e2e4"}, 10*time.Millisecond, 20*time.Millisecond)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v, want %s failure", err, tt.want)
			}
		})
	}
}

func TestUnsupportedOptionRejected(t *testing.T) {
	_, err := startEngine(helperPath(t, "normal"), map[string]string{"Skill Level": "1"})
	if err == nil {
		t.Fatal("engine without Skill Level was accepted as configurable")
	}
}

func TestIllegalMoveForfeitsForBothColors(t *testing.T) {
	for _, white := range []bool{true, false} {
		path := helperPath(t, "illegal")
		game, record, err := playGame(path, path, nil, nil, white, settings{MoveMS: 10, GraceMS: 100, MaxPlies: 10})
		if err != nil {
			t.Fatal(err)
		}
		if record.Result != "0-1" || !strings.HasPrefix(record.Termination, "forfeit:") {
			t.Fatalf("unexpected result: %+v", record)
		}
		stats := summarize([]gameRecord{record})
		if white && stats.Losses != 1 || !white && stats.Wins != 1 {
			t.Fatalf("wrong color perspective: %+v", stats)
		}
		// A saved forfeit must carry the same result in the PGN move text.
		if !strings.HasSuffix(game.String(), record.Result) {
			t.Fatal("PGN result differs")
		}
	}
}

func TestMoveCapKeepsUnfinishedAndPreservesOpening(t *testing.T) {
	path := helperPath(t, "normal")
	opening := []string{"e2e4", "e7e5"}
	for _, white := range []bool{true, false} {
		game, record, err := playGame(path, path, nil, opening, white, settings{MoveMS: 10, GraceMS: 200, MaxPlies: 6})
		if err != nil {
			t.Fatal(err)
		}
		if record.Result != "*" || record.Plies != 6 || summarize([]gameRecord{record}).Score != nil {
			t.Fatalf("cap scored as completed: %+v", record)
		}
		if _, err := chess.PGN(strings.NewReader(game.String())); err != nil {
			t.Fatalf("invalid saved PGN: %v", err)
		}
		if record.Opening[0] != "e2e4" || record.Opening[1] != "e7e5" {
			t.Fatal("opening changed")
		}
	}
}

func TestScoring(t *testing.T) {
	games := []gameRecord{
		{Result: "1-0", CandidateWhite: true}, {Result: "0-1", CandidateWhite: false},
		{Result: "1-0", CandidateWhite: false}, {Result: "0-1", CandidateWhite: true, Termination: "forfeit: timeout"},
		{Result: "1/2-1/2"}, {Result: "*"},
	}
	s := summarize(games)
	if s.Wins != 2 || s.Draws != 1 || s.Losses != 2 || s.Unfinished != 1 || s.Forfeits != 1 || s.Score == nil || *s.Score != 50 {
		t.Fatalf("unexpected score: %+v", s)
	}
}

func TestTerminalResultTakesPrecedenceOverMoveCap(t *testing.T) {
	for _, tt := range []struct {
		mode           string
		cap            int
		result, method string
	}{
		{"mate", 4, "0-1", "Checkmate"},
		{"repeat", 8, "1/2-1/2", "ThreefoldRepetition"},
	} {
		t.Run(tt.mode, func(t *testing.T) {
			path := helperPath(t, tt.mode)
			_, record, err := playGame(path, path, nil, nil, true, settings{MoveMS: 10, GraceMS: 200, MaxPlies: tt.cap})
			if err != nil {
				t.Fatal(err)
			}
			if record.Result != tt.result || record.Termination != tt.method {
				t.Fatalf("terminal result lost at cap: %+v", record)
			}
		})
	}
}

func TestOpenings(t *testing.T) {
	openings, err := loadOpenings("openings-screening-10-v1.txt")
	if err != nil || len(openings) != 10 {
		t.Fatalf("openings: %v, %v", openings, err)
	}
}

func TestInvalidOpeningsRejected(t *testing.T) {
	for _, moves := range []string{"# empty suite", "e2e5", "f2f3 e7e5 g2g4 d8h4"} {
		path := filepath.Join(t.TempDir(), "openings.txt")
		if err := os.WriteFile(path, []byte(moves), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err := loadOpenings(path); err == nil {
			t.Fatalf("accepted invalid opening %q", moves)
		}
	}
}

func TestReportUsesSavedHistoryAndSkipsPartialRuns(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "partial"), 0755); err != nil {
		t.Fatal(err)
	}
	for i, label := range []string{"v0.1", "v0.2"} {
		dir := filepath.Join(root, fmt.Sprint(i))
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}
		s := summary{Schema: 2, Complete: true, Created: fmt.Sprint(i), Candidate: identity{Label: label}, Matches: []match{{Engine: identity{Label: "Stockfish Skill 1"}, Stats: summarize([]gameRecord{{Result: "0-1", CandidateWhite: true}})}}}
		data, _ := json.Marshal(s)
		if err := os.WriteFile(filepath.Join(dir, "summary.json"), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := report(root); err != nil {
		t.Fatal(err)
	}
	table, _ := os.ReadFile(filepath.Join(root, "README.md"))
	if strings.Count(string(table), "| v0.") != 2 || strings.Contains(string(table), "[JSON](partial/") {
		t.Fatal("report dropped history or included incomplete run")
	}
	if strings.Index(string(table), "v0.1") > strings.Index(string(table), "v0.2") {
		t.Fatal("history is not chronological")
	}
	if err := report(root); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil || string(table) != string(again) {
		t.Fatal("report regeneration changed the saved results")
	}
}

func TestRunDirectoriesPreservePreviousResults(t *testing.T) {
	root := t.TempDir()
	first, err := createRunDir(root, "v0.2", "v0.1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(first) != "v0.2-vs-v0.1-and-stockfish-1" {
		t.Fatalf("unexpected name: %s", first)
	}
	marker := filepath.Join(first, "summary.json")
	if err := os.WriteFile(marker, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, suffix := range []string{"-run-2", "-run-3"} {
		dir, err := createRunDir(root, "v0.2", "v0.1", 1)
		if err != nil {
			t.Fatal(err)
		}
		if dir != first+suffix {
			t.Fatalf("unexpected repeat name: %s", dir)
		}
	}
	data, err := os.ReadFile(marker)
	if err != nil || string(data) != "original" {
		t.Fatal("previous results were overwritten")
	}
	baseline, err := createRunDir(root, "v0.1", "", 1)
	if err != nil || filepath.Base(baseline) != "v0.1-vs-stockfish-1" {
		t.Fatalf("unexpected baseline: %s (%v)", baseline, err)
	}
}

func TestVersionLabelsStayInsideResultsDirectory(t *testing.T) {
	root := t.TempDir()
	dir, err := createRunDir(root, "../../test version", "../older", 1)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(dir) != root {
		t.Fatalf("label escaped results directory: %s", dir)
	}
}

func TestConfiguredSkillIsSentAndNamed(t *testing.T) {
	trace := filepath.Join(t.TempDir(), "commands.txt")
	t.Setenv("PANDA_BENCH_TRACE", trace)
	e, err := startEngine(helperPath(t, "options"), stockfishOptions(2))
	if err != nil {
		t.Fatal(err)
	}
	e.close()
	commands, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"setoption name Skill Level value 2", "setoption name UCI_LimitStrength value false", "setoption name Threads value 1", "setoption name Hash value 16"} {
		if !strings.Contains(string(commands), command+"\n") {
			t.Fatalf("missing command %q", command)
		}
	}
	dir, err := createRunDir(t.TempDir(), "candidate", "previous", 2)
	if err != nil || !strings.HasSuffix(dir, "stockfish-2") {
		t.Fatalf("skill not reflected in folder: %s, %v", dir, err)
	}
}

func TestSnapshotSurvivesOriginalRebuild(t *testing.T) {
	root := t.TempDir()
	original := filepath.Join(root, "engine")
	if err := os.WriteFile(original, []byte("original executable"), 0700); err != nil {
		t.Fatal(err)
	}
	snapshot, err := snapshotEngine(original, root, "snapshot")
	if err != nil {
		t.Fatal(err)
	}
	before, err := digest(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(original, []byte("rebuilt executable"), 0700); err != nil {
		t.Fatal(err)
	}
	after, err := digest(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(snapshot)
	if err != nil || before != after || string(data) != "original executable" {
		t.Fatal("snapshot changed with original")
	}
}

func TestWholeMoveDeadlineIncludesBlockedWrites(t *testing.T) {
	e, err := startEngine(helperPath(t, "block-input"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer e.close()
	started := time.Now()
	_, elapsed, err := e.move([]string{strings.Repeat("e2e4 ", 100000)}, 20*time.Millisecond, 20*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("blocked write not bounded: %v", err)
	}
	if elapsed < 40*time.Millisecond || time.Since(started) > 2*time.Second {
		t.Fatalf("unexpected deadline timing: %v", elapsed)
	}
}

func TestSlowEngineForfeitsAndRecordsTiming(t *testing.T) {
	path := helperPath(t, "slow")
	_, record, err := playGame(path, path, nil, nil, true, settings{MoveMS: 20, GraceMS: 20, MaxPlies: 10})
	if err != nil {
		t.Fatal(err)
	}
	if record.Result != "0-1" || !strings.Contains(record.Termination, "timeout") {
		t.Fatalf("late engine accepted: %+v", record)
	}
	if len(record.Timings) != 1 || !record.Timings[0].Candidate || record.Timings[0].ElapsedMS < 40 {
		t.Fatalf("missing failed attempt timing: %+v", record.Timings)
	}
}

func TestPossibleScoresIncludeAllUnfinishedOutcomes(t *testing.T) {
	for _, tt := range []struct {
		games        []gameRecord
		lower, upper float64
	}{
		{[]gameRecord{{Result: "*"}}, 0, 100},
		{[]gameRecord{{Result: "1-0", CandidateWhite: true}, {Result: "1/2-1/2"}, {Result: "*"}}, 50, 100 * 2.5 / 3},
		{[]gameRecord{{Result: "0-1", CandidateWhite: true}}, 0, 0},
	} {
		s := summarize(tt.games)
		if s.ScoreLower == nil || s.ScoreUpper == nil || math.Abs(*s.ScoreLower-tt.lower) > 1e-9 || math.Abs(*s.ScoreUpper-tt.upper) > 1e-9 {
			t.Fatalf("wrong possible score: %+v", s)
		}
	}
	if s := summarize(nil); s.ScoreLower != nil || s.ScoreUpper != nil {
		t.Fatal("empty match has score bounds")
	}
}

func TestCheckpointPreservesCompletedMatchOnFailedUpdate(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "run")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	s := summary{Schema: 2, PlannedMatches: 2, Candidate: identity{Label: "candidate"}, Matches: []match{
		{Engine: identity{Label: "previous"}, Games: []gameRecord{{Result: "*"}}},
	}}
	if err := saveSummary(dir, s); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "summary.json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := report(root); err != nil {
		t.Fatal(err)
	}
	table, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil || !strings.Contains(string(table), "partial (1/2 matchups)") || !strings.Contains(string(table), "0.0% to 100.0%") {
		t.Fatal("partial checkpoint not clearly reported")
	}
	value := math.NaN()
	s.Matches[0].Stats.Score = &value
	if err := saveSummary(dir, s); err == nil {
		t.Fatal("invalid summary update accepted")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed update destroyed checkpoint")
	}
	var saved summary
	if err := json.Unmarshal(after, &saved); err != nil || saved.Complete || len(saved.Matches) != 1 {
		t.Fatal("completed matchup was not preserved in partial run")
	}
}

func TestConfirmationOpeningsAreDistinctAndLegal(t *testing.T) {
	openings, err := loadOpenings("openings-confirmation-50-v1.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(openings) != 50 {
		t.Fatalf("got %d openings, want 50", len(openings))
	}
	seen := map[string]bool{}
	for _, opening := range openings {
		game, err := newGame(opening)
		if err != nil {
			t.Fatal(err)
		}
		position := strings.Join(strings.Fields(game.FEN())[:4], " ")
		if seen[position] {
			t.Fatalf("duplicate position: %v", opening)
		}
		seen[position] = true
	}
}
