package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPreviousGate(t *testing.T) {
	for _, c := range []struct {
		s  stats
		ok bool
	}{{stats{Wins: 60, Draws: 10, Losses: 30}, true}, {stats{Wins: 40, Draws: 20, Losses: 40}, false}, {stats{Wins: 50, Draws: 10, Losses: 30, Unfinished: 10}, true}, {stats{Wins: 40, Draws: 20, Losses: 30, Unfinished: 10}, false}, {stats{Wins: 60, Losses: 40, Forfeits: 1}, false}, {stats{}, false}} {
		if favorablePrevious(c.s) != c.ok {
			t.Fatalf("gate %+v", c.s)
		}
	}
}
func TestPreviousBeforeParallelExternal(t *testing.T) {
	roles := []opponent{{Role: "previous"}, {Role: "stockfish"}, {Role: "sacrifice"}}
	started := make(chan string, 3)
	release := make(chan struct{})
	done := make(chan error, 1)
	play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
		started <- o.Role
		if o.Role != "previous" {
			<-release
		}
		return match{Engine: identity{Label: o.Role}, Stats: stats{Wins: 60, Losses: 40}}, nil
	}
	go func() {
		saved := 0
		status := ""
		matches, err := playScheduleWith("", identity{}, roles, nil, settings{ExternalConcurrency: 2}, "", time.Now(), false, func(ms []match) error {
			saved++
			if ms[0].Engine.Label != "previous" {
				return fmt.Errorf("previous checkpoint missing")
			}
			return nil
		}, func(s string) { status = s }, play)
		if err == nil && (len(matches) != 3 || saved != 3 || !strings.HasPrefix(status, "passed")) {
			err = fmt.Errorf("bad completion %d/%d/%s", len(matches), saved, status)
		}
		done <- err
	}()
	for i := 0; i < 3; i++ {
		select {
		case role := <-started:
			if i == 0 && role != "previous" {
				t.Fatal("external ran first")
			}
		case <-time.After(2 * time.Second):
			close(release)
			t.Fatal("external matchups did not overlap")
		}
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
func TestUnfavorablePreviousStopsExternal(t *testing.T) {
	calls := 0
	gate := ""
	play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
		calls++
		return match{Stats: stats{Wins: 40, Draws: 20, Losses: 40}}, nil
	}
	saved := 0
	_, err := playScheduleWith("", identity{}, []opponent{{Role: "previous"}, {Role: "stockfish"}, {Role: "sacrifice"}}, nil, settings{ExternalConcurrency: 2}, "", time.Now(), false, func(ms []match) error { saved = len(ms); return nil }, func(s string) { gate = s }, play)
	if err == nil || calls != 1 || saved != 1 || !strings.HasPrefix(gate, "blocked") {
		t.Fatal("weak previous screen did not block external games")
	}
}
func TestSequentialExternalOrdering(t *testing.T) {
	order := make(chan string, 2)
	play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
		order <- o.Role
		return match{}, nil
	}
	_, err := playScheduleWith("", identity{}, []opponent{{Role: "stockfish"}, {Role: "sacrifice"}}, nil, settings{ExternalConcurrency: 1}, "", time.Now(), false, func([]match) error { return nil }, func(string) {}, play)
	if err != nil || <-order != "stockfish" || <-order != "sacrifice" {
		t.Fatal("sequential fallback order changed")
	}
}
func TestChangedSacrificeAssetStopsMatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Tal.pgn")
	if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}
	hash, err := digest(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = runMatch("", identity{}, opponent{Role: "sacrifice", AssetPath: path, AssetSHA: hash}, nil, settings{}, t.TempDir(), time.Now(), false)
	if err == nil || !strings.Contains(err.Error(), "asset changed") {
		t.Fatal("asset mutation not detected")
	}
}

func TestReportTablesAndOpponentPGNLinks(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "candidate")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	s := summary{Schema: 2, Complete: true, Settings: settings{ExternalConcurrency: 2}, Matches: []match{{Engine: identity{Label: "Sacrifice"}, PGN: "sacrifice.pgn"}}}
	if err := saveSummary(dir, s); err != nil {
		t.Fatal(err)
	}
	if err := report(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "[PGN](candidate/sacrifice.pgn)") {
		t.Fatal("report did not link to opponent-specific games")
	}
	columns := 0
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "|") {
			columns = 0
			continue
		}
		if columns == 0 {
			columns = strings.Count(line, "|")
		} else if strings.Count(line, "|") != columns {
			t.Fatalf("inconsistent Markdown table columns: %s", line)
		}
	}
}

func TestParallelFailureReapsAndSavesOtherMatch(t *testing.T) {
	started := make(chan string, 2)
	release := make(chan struct{})
	done := make(chan error, 1)
	play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
		started <- o.Role
		if o.Role == "stockfish" {
			return match{}, fmt.Errorf("engine exited")
		}
		<-release
		return match{Engine: identity{Label: o.Role}}, nil
	}
	go func() {
		saved := false
		matches, err := playScheduleWith("", identity{}, []opponent{{Role: "stockfish"}, {Role: "sacrifice"}}, nil, settings{ExternalConcurrency: 2}, "", time.Now(), false, func(ms []match) error {
			saved = len(ms) == 1 && ms[0].Engine.Label == "sacrifice"
			return nil
		}, func(string) {}, play)
		if err == nil || !saved || len(matches) != 1 {
			done <- fmt.Errorf("lost successful matchup after peer failure")
			return
		}
		done <- nil
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			close(release)
			t.Fatal("missing parallel job")
		}
	}
	select {
	case <-done:
		close(release)
		t.Fatal("scheduler returned with a live writer")
	default:
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestValidationRequiresRoutineOpponents(t *testing.T) {
	for _, c := range []struct {
		stage, previous, sacrifice string
		validation, valid          bool
	}{
		{"all", "previous", "sacrifice", true, true},
		{"all", "", "sacrifice", true, false},
		{"all", "previous", "", true, false},
		{"previous", "previous", "", true, true},
		{"previous", "", "", false, false},
		{"external", "", "sacrifice", true, true},
		{"external", "", "", true, false},
		{"all", "", "", false, true}, // Historical direct Stockfish-only runs.
		{"unknown", "previous", "sacrifice", true, false},
	} {
		if got := validateSchedule(c.stage, c.previous, c.sacrifice, c.validation); (got == nil) != c.valid {
			t.Fatalf("schedule %+v: %v", c, got)
		}
	}
}

func TestAllMatchupsOverlapAndReviewPrevious(t *testing.T) {
	for _, favorable := range []bool{true, false} {
		t.Run(fmt.Sprint(favorable), func(t *testing.T) {
			roles := []opponent{{Role: "previous"}, {Role: "stockfish"}, {Role: "sacrifice"}}
			started := make(chan string, 3)
			release := make(chan struct{})
			t.Cleanup(func() { close(release) })
			type result struct {
				matches []match
				err     error
				status  string
				saves   int
			}
			done := make(chan result, 1)
			play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
				started <- o.Role
				<-release
				s := stats{Wins: 35, Draws: 10, Losses: 5}
				if !favorable && o.Role == "previous" {
					s = stats{Wins: 10, Draws: 20, Losses: 20}
				}
				return match{Engine: identity{Label: o.Role}, Stats: s}, nil
			}
			go func() {
				status, saves := "", 0
				ms, err := playScheduleWith("", identity{}, roles, nil, settings{MatchupConcurrency: 3, ExternalConcurrency: 2}, "", time.Now(), false, func(ms []match) error {
					saves++
					return nil
				}, func(s string) { status = s }, play)
				done <- result{ms, err, status, saves}
			}()
			seen := map[string]bool{}
			for range roles {
				select {
				case role := <-started:
					seen[role] = true
				case <-time.After(2 * time.Second):
					t.Fatal("all three matchups did not overlap")
				}
			}
			if len(seen) != 3 {
				t.Fatal("duplicate/missing matchups", seen)
			}
			// Send one token per worker instead of closing, so cleanup can safely
			// release blocked workers even if an assertion above fails.
			for range roles {
				release <- struct{}{}
			}
			r := <-done
			if len(r.matches) != 3 || r.saves != 3 || (r.err == nil) != favorable {
				t.Fatalf("completion: %+v", r)
			}
			for i, o := range roles {
				if r.matches[i].Engine.Label != o.Role {
					t.Fatal("completion order changed report order")
				}
			}
			prefix := "passed"
			if !favorable {
				prefix = "blocked"
			}
			if !strings.HasPrefix(r.status, prefix) {
				t.Fatal("wrong previous-release decision", r.status)
			}
		})
	}
}

func TestAllMatchupFailureReapsAndSavesPeers(t *testing.T) {
	started := make(chan string, 3)
	release := make(chan struct{})
	defer close(release)
	done := make(chan error, 1)
	play := func(_ string, _ identity, o opponent, _ [][]string, _ settings, _ string, _ time.Time, _ bool) (match, error) {
		started <- o.Role
		if o.Role == "stockfish" {
			return match{}, fmt.Errorf("engine exited")
		}
		<-release
		return match{Engine: identity{Label: o.Role}, Stats: stats{Wins: 50}}, nil
	}
	go func() {
		saved := 0
		ms, err := playScheduleWith("", identity{}, []opponent{{Role: "previous"}, {Role: "stockfish"}, {Role: "sacrifice"}}, nil, settings{MatchupConcurrency: 3}, "", time.Now(), false, func(ms []match) error { saved = len(ms); return nil }, func(string) {}, play)
		if err == nil || len(ms) != 2 || saved != 2 || ms[0].Engine.Label != "previous" || ms[1].Engine.Label != "sacrifice" {
			done <- fmt.Errorf("peer results lost: %v %d %d", err, len(ms), saved)
			return
		}
		done <- nil
	}()
	for range 3 {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatal("missing job")
		}
	}
	select {
	case <-done:
		t.Fatal("returned while peer writers were active")
	default:
	}
	release <- struct{}{}
	release <- struct{}{}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestReportRecordsAllMatchupConcurrency(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "parallel-all")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := saveSummary(dir, summary{Schema: 2, Complete: true, Settings: settings{MatchupConcurrency: 3, ExternalConcurrency: 2}}); err != nil {
		t.Fatal(err)
	}
	if err := report(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil || !strings.Contains(string(data), "3 (all)") {
		t.Fatal("all-matchup concurrency missing from report", err)
	}
}

func TestScheduleConcurrencyDefaultsAndOverrides(t *testing.T) {
	for _, tt := range []struct {
		name, stage                                    string
		previous, sacrifice, parallel, all, sequential bool
		cpus, wantAll, wantExternal                    int
		wantError                                      bool
	}{
		{"routine defaults", "all", true, true, false, false, false, 8, 3, 2, false},
		{"four core host", "all", true, true, false, false, false, 4, 3, 2, false},
		{"small host fallback", "all", true, true, false, false, false, 2, 0, 1, false},
		{"explicit serial", "all", true, true, false, false, true, 8, 0, 1, false},
		{"external-only legacy mode", "all", true, true, true, false, false, 8, 0, 2, false},
		{"explicit all", "all", true, true, false, true, false, 8, 3, 2, false},
		{"previous stage", "previous", true, true, false, false, false, 8, 0, 1, false},
		{"no sacrifice", "all", true, false, false, false, false, 8, 0, 1, false},
		{"no previous", "all", false, true, false, false, false, 8, 0, 1, false},
		{"explicit all needs capacity", "all", true, true, false, true, false, 2, 0, 0, true},
		{"explicit all needs full stage", "external", true, true, false, true, false, 8, 0, 0, true},
		{"explicit all needs previous", "all", false, true, false, true, false, 8, 0, 0, true},
		{"explicit all needs sacrifice", "all", true, false, false, true, false, 8, 0, 0, true},
		{"conflicting parallel modes", "all", true, true, true, true, false, 8, 0, 0, true},
		{"conflicting serial modes", "all", true, true, false, true, true, 8, 0, 0, true},
		{"conflicting legacy mode", "all", true, true, true, false, true, 8, 0, 0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			all, external, err := scheduleConcurrency(tt.stage, tt.previous, tt.sacrifice, tt.parallel, tt.all, tt.sequential, tt.cpus)
			if (err != nil) != tt.wantError || all != tt.wantAll || external != tt.wantExternal {
				t.Fatalf("got %d/%d/%v, want %d/%d/error=%v", all, external, err, tt.wantAll, tt.wantExternal, tt.wantError)
			}
		})
	}
}
