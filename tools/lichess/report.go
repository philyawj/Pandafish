package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

type searchRecord struct {
	Process     string             `json:"process"`
	GameNumber  int                `json:"game_number"`
	Search      int                `json:"search"`
	Time        time.Time          `json:"time"`
	Position    string             `json:"position"`
	Go          string             `json:"go"`
	FEN         string             `json:"fen,omitempty"`
	Bestmove    string             `json:"bestmove,omitempty"`
	ResponseMS  float64            `json:"response_ms"`
	Status      string             `json:"status"`
	Metrics     *telemetry.Metrics `json:"metrics,omitempty"`
	Game        string             `json:"game,omitempty"`
	Ply         *int               `json:"ply,omitempty"`
	Correlation string             `json:"correlation"`
	Candidates  []string           `json:"candidates,omitempty"`
	startedNS   int64
}

type savedGame struct {
	id, result string
	moves      []string
	positions  []string
}

type sessionReport struct {
	Schema             int            `json:"schema"`
	SessionStatus      string         `json:"session_status"`
	Processes          int            `json:"processes"`
	Searches           int            `json:"searches"`
	Completed          int            `json:"completed"`
	Cancelled          int            `json:"cancelled"`
	Incomplete         int            `json:"incomplete"`
	Matched            int            `json:"matched"`
	MissingMetrics     int            `json:"missing_metrics"`
	Fallbacks          int            `json:"fallbacks"`
	PastHardBudget     int            `json:"past_hard_budget"`
	MaxResponseMS      float64        `json:"max_response_ms"`
	Games              int            `json:"games"`
	UnfinishedGames    int            `json:"unfinished_games"`
	MaxBudgetOverrunMS float64        `json:"max_budget_overrun_ms"`
	ResponseP95MS      float64        `json:"response_p95_ms"`
	DepthCounts        map[int]int    `json:"depth_counts"`
	StopReasons        map[string]int `json:"stop_reasons"`
	Warnings           []string       `json:"warnings"`
}

func readTrace(path string) (rows []*searchRecord, warnings []string) {
	f, err := os.Open(path)
	if err != nil {
		return nil, []string{"cannot read trace " + filepath.Base(path)}
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 4*1024*1024)
	position := "position startpos"
	gameNumber, number := 0, 0
	var pending *searchRecord
	var sequence uint64
	var elapsed int64
	exited := false
	abandon := func(status string) {
		if pending != nil {
			pending.Status = status
			pending = nil
		}
	}
	for scanner.Scan() {
		var event telemetry.Event
		if json.Unmarshal(scanner.Bytes(), &event) != nil || event.Schema != 1 || event.Sequence != sequence+1 || event.ElapsedNS < elapsed {
			warnings = append(warnings, "invalid or truncated trace: "+filepath.Base(path))
			break
		}
		sequence, elapsed = event.Sequence, event.ElapsedNS
		switch event.Kind {
		case "in":
			args := strings.Fields(event.Line)
			if len(args) == 0 {
				continue
			}
			switch args[0] {
			case "position":
				if _, e := positionGame(event.Line); e != nil {
					warnings = append(warnings, "invalid position request in "+filepath.Base(path))
					continue
				}
				abandon("cancelled")
				position = event.Line
			case "ucinewgame":
				abandon("cancelled")
				gameNumber++
				position = "position startpos"
			case "go":
				abandon("cancelled")
				number++
				pending = &searchRecord{Process: filepath.Base(path), GameNumber: gameNumber, Search: number, Time: event.Time, Position: position, Go: event.Line, Status: "incomplete", Correlation: "missing", startedNS: event.ElapsedNS}
				if g, e := positionGame(position); e == nil {
					pending.FEN = g.FEN()
				} else {
					warnings = append(warnings, fmt.Sprintf("invalid position in %s search %d", filepath.Base(path), number))
				}
				rows = append(rows, pending)
			case "quit":
				abandon("cancelled")
			}
		case "out":
			for _, line := range strings.Split(strings.TrimSpace(event.Line), "\n") {
				if pending == nil {
					continue
				}
				if strings.HasPrefix(line, "info string invalid go:") {
					pending.Status = "invalid_request"
					pending = nil
					warnings = append(warnings, "invalid go request in "+filepath.Base(path))
					continue
				}
				if strings.HasPrefix(line, telemetry.Prefix) {
					var m telemetry.Metrics
					if json.Unmarshal([]byte(strings.TrimPrefix(line, telemetry.Prefix)), &m) == nil {
						pending.Metrics = &m
					} else {
						warnings = append(warnings, "invalid metrics in "+filepath.Base(path))
					}
				}
				if strings.HasPrefix(line, "bestmove ") {
					fields := strings.Fields(line)
					if len(fields) < 2 {
						continue
					}
					pending.Bestmove = fields[1]
					pending.Status = "completed"
					pending.ResponseMS = float64(event.ElapsedNS-pending.startedNS) / float64(time.Millisecond)
					if g, e := positionGame(pending.Position); e == nil {
						if pending.Bestmove == "0000" {
							if len(g.ValidMoves()) > 0 {
								e = fmt.Errorf("null move with legal moves")
							}
						} else {
							e = g.PushNotationMove(pending.Bestmove, chess.UCINotation{}, nil)
						}
						if e != nil {
							pending.Status = "illegal_move"
							warnings = append(warnings, "illegal bestmove in "+filepath.Base(path))
						}
					}
					pending = nil
				}
			}
		case "exit":
			if event.Line == "error" {
				warnings = append(warnings, "engine exited with an error: "+filepath.Base(path))
			}
			exited = true
			abandon("cancelled")
		}
	}
	if scanner.Err() != nil {
		warnings = append(warnings, "trace read failed: "+filepath.Base(path))
	}
	if !exited {
		warnings = append(warnings, "missing engine exit event: "+filepath.Base(path))
	}
	return rows, warnings
}

func positionGame(command string) (*chess.Game, error) {
	args := strings.Fields(command)
	if len(args) < 2 || args[0] != "position" {
		return nil, fmt.Errorf("invalid position")
	}
	args = args[1:]
	game := chess.NewGame()
	switch {
	case args[0] == "startpos":
		args = args[1:]
	case args[0] == "fen" && len(args) >= 7:
		fen, err := chess.FEN(strings.Join(args[1:7], " "))
		if err != nil {
			return nil, err
		}
		game = chess.NewGame(fen)
		args = args[7:]
	default:
		return nil, fmt.Errorf("invalid position")
	}
	if len(args) > 0 {
		if args[0] != "moves" {
			return nil, fmt.Errorf("invalid move list")
		}
		for _, move := range args[1:] {
			if err := game.PushNotationMove(move, chess.UCINotation{}, nil); err != nil {
				return nil, err
			}
		}
	}
	return game, nil
}

func loadGames(dir string) ([]savedGame, []string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []string{"cannot read games directory"}
	}
	byID := map[string]savedGame{}
	var warnings []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".pgn") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			warnings = append(warnings, "cannot read PGN "+entry.Name())
			continue
		}
		option, err := chess.PGN(f)
		f.Close()
		if err != nil {
			warnings = append(warnings, "invalid or truncated PGN "+entry.Name())
			continue
		}
		g := chess.NewGame(option)
		id := g.GetTagPair("Site")
		if id == "" || id == "?" {
			warnings = append(warnings, "PGN without Site: "+entry.Name())
			continue
		}
		saved := savedGame{id: id, result: g.GetTagPair("Result")}
		for _, m := range g.Moves() {
			saved.moves = append(saved.moves, m.String())
		}
		for _, p := range g.Positions() {
			saved.positions = append(saved.positions, p.String())
		}
		if previous, ok := byID[id]; ok {
			if len(previous.moves) > len(saved.moves) {
				continue
			}
			if len(previous.moves) == len(saved.moves) && previous.result != "*" {
				continue
			}
		}
		byID[id] = saved
	}
	games := make([]savedGame, 0, len(byID))
	for _, g := range byID {
		games = append(games, g)
	}
	sort.Slice(games, func(i, j int) bool { return games[i].id < games[j].id })
	return games, warnings
}

// Require the recorded history, root FEN and returned move, not just a board match.
func matchingPly(row *searchRecord, g savedGame) (int, bool) {
	request, err := positionGame(row.Position)
	if err != nil {
		return 0, false
	}
	history := request.Moves()
	positions := request.Positions()
	if len(positions) == 0 {
		return 0, false
	}
	for base, p := range g.positions {
		ply := base + len(history)
		if p != positions[0].String() || ply >= len(g.positions) || g.positions[ply] != row.FEN {
			continue
		}
		ok := true
		for i, m := range history {
			if base+i >= len(g.moves) || g.moves[base+i] != m.String() {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if row.Bestmove != "" && row.Bestmove != "0000" && (ply >= len(g.moves) || g.moves[ply] != row.Bestmove) {
			continue
		}
		return ply, true
	}
	return 0, false
}

func correlate(rows []*searchRecord, games []savedGame) {
	groups := map[string][]*searchRecord{}
	for _, r := range rows {
		key := fmt.Sprintf("%s/%d", r.Process, r.GameNumber)
		groups[key] = append(groups[key], r)
	}
	for _, group := range groups {
		var candidates []savedGame
		for _, g := range games {
			ok := true
			for _, r := range group {
				if _, match := matchingPly(r, g); !match {
					ok = false
					break
				}
			}
			if ok {
				candidates = append(candidates, g)
			}
		}
		for _, r := range group {
			r.Correlation = "missing"
			r.Candidates = nil
			r.Game = ""
			r.Ply = nil
			for _, g := range candidates {
				r.Candidates = append(r.Candidates, g.id)
			}
			if len(candidates) > 1 {
				r.Correlation = "ambiguous"
			}
			if len(candidates) == 1 {
				r.Correlation = "matched"
				r.Game = candidates[0].id
				ply, _ := matchingPly(r, candidates[0])
				r.Ply = &ply
			}
		}
	}
}

func report(session string) error {
	if session == "" {
		return fmt.Errorf("report requires --session")
	}
	raw, err := os.ReadFile(filepath.Join(session, "manifest.json"))
	if err != nil {
		return err
	}
	var m manifest
	if err = json.Unmarshal(raw, &m); err != nil || m.Schema != 1 {
		return fmt.Errorf("invalid session manifest")
	}
	result := sessionReport{Schema: 1, SessionStatus: m.Status, Warnings: []string{}, DepthCounts: map[int]int{}, StopReasons: map[string]int{}}
	for _, artifact := range []struct{ name, expected string }{{"pandafish", m.EngineSHA256}, {"config.redacted.yml", m.ConfigSHA256}} {
		data, e := os.ReadFile(filepath.Join(session, artifact.name))
		if e != nil {
			result.Warnings = append(result.Warnings, "cannot verify saved "+artifact.name)
			continue
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != artifact.expected {
			result.Warnings = append(result.Warnings, "saved "+artifact.name+" SHA-256 does not match manifest")
		}
	}
	if m.Status != "finished" {
		result.Warnings = append(result.Warnings, "session is not marked finished; this report may be partial")
	}
	traces, err := filepath.Glob(filepath.Join(session, "traces", "engine-*.jsonl"))
	if err != nil {
		return err
	}
	result.Processes = len(traces)
	var rows []*searchRecord
	for _, path := range traces {
		r, w := readTrace(path)
		rows = append(rows, r...)
		result.Warnings = append(result.Warnings, w...)
	}
	games, warnings := loadGames(filepath.Join(session, "games"))
	result.Warnings = append(result.Warnings, warnings...)
	result.Games = len(games)
	for _, g := range games {
		if g.result == "*" || g.result == "" {
			result.UnfinishedGames++
		}
	}
	correlate(rows, games)
	var jsonl strings.Builder
	var responses []float64
	for _, r := range rows {
		result.Searches++
		switch r.Status {
		case "completed":
			result.Completed++
			responses = append(responses, r.ResponseMS)
		case "cancelled":
			result.Cancelled++
		default:
			result.Incomplete++
		}
		if r.Correlation == "matched" {
			result.Matched++
		}
		if r.Metrics == nil {
			result.MissingMetrics++
		} else {
			result.DepthCounts[r.Metrics.Depth]++
			result.StopReasons[r.Metrics.Reason]++
			if r.Metrics.Fallback {
				result.Fallbacks++
			}
			if r.Metrics.Timed && r.Metrics.ElapsedMS > r.Metrics.HardMS {
				result.PastHardBudget++
				result.MaxBudgetOverrunMS = max(result.MaxBudgetOverrunMS, r.Metrics.ElapsedMS-r.Metrics.HardMS)
			}
		}
		if r.ResponseMS > result.MaxResponseMS {
			result.MaxResponseMS = r.ResponseMS
		}
		data, e := json.Marshal(r)
		if e != nil {
			return e
		}
		jsonl.Write(data)
		jsonl.WriteByte('\n')
	}
	sort.Float64s(responses)
	if len(responses) > 0 {
		result.ResponseP95MS = responses[(95*len(responses)+99)/100-1]
	}
	if result.Searches == 0 {
		result.Warnings = append(result.Warnings, "no recorded searches")
	}
	if result.Matched < result.Searches {
		result.Warnings = append(result.Warnings, "some searches have missing or ambiguous PGN linkage; see moves.jsonl")
	}
	if result.MissingMetrics > 0 {
		result.Warnings = append(result.Warnings, "some searches have no final engine metrics")
	}
	if result.Incomplete > 0 {
		result.Warnings = append(result.Warnings, "incomplete or illegal searches require inspection")
	}
	if err = os.WriteFile(filepath.Join(session, "moves.jsonl"), []byte(jsonl.String()), 0644); err != nil {
		return err
	}
	if err = writeJSON(filepath.Join(session, "summary.json"), result); err != nil {
		return err
	}
	text := fmt.Sprintf("# Lichess session report\n\nEngine: %s\n\nEngine SHA-256: `%s`\n\n| Measure | Value |\n| --- | --- |\n| Session status | %s |\n| Saved games / unfinished | %d / %d |\n| Engine processes | %d |\n| Searches / completed / cancelled / incomplete or illegal | %d / %d / %d / %d |\n| Searches linked to PGN | %d |\n| Searches without final metrics | %d |\n| Legal fallback moves | %d |\n| Searches past engine hard budget | %d |\n| Maximum recorded response | %.3f ms |\n\n", m.EngineID, m.EngineSHA256, result.SessionStatus, result.Games, result.UnfinishedGames, result.Processes, result.Searches, result.Completed, result.Cancelled, result.Incomplete, result.Matched, result.MissingMetrics, result.Fallbacks, result.PastHardBudget, result.MaxResponseMS)
	text += fmt.Sprintf("Maximum hard-budget overrun: %.3f ms. Response p95: %.3f ms.\n\n", result.MaxBudgetOverrunMS, result.ResponseP95MS)
	text += "| Completed depth | Searches |\n| --- | --- |\n"
	var depths []int
	for depth := range result.DepthCounts {
		depths = append(depths, depth)
	}
	sort.Ints(depths)
	for _, depth := range depths {
		text += fmt.Sprintf("| %d | %d |\n", depth, result.DepthCounts[depth])
	}
	text += "\n| Stop reason | Searches |\n| --- | --- |\n"
	var reasons []string
	for reason := range result.StopReasons {
		reasons = append(reasons, reason)
	}
	sort.Strings(reasons)
	for _, reason := range reasons {
		text += fmt.Sprintf("| %s | %d |\n", reason, result.StopReasons[reason])
	}
	text += "\n"
	text += "Engine elapsed time and response time are separate measurements. Small hard-budget overruns can occur while an in-progress node finishes; inspect their size in `moves.jsonl`. Client submission and network latency are not measured. Scores use the root side's perspective; mate scores give a sign without distance. PV currently contains only the chosen root move.\n\n"
	if len(result.Warnings) > 0 {
		text += "## Recording warnings\n\n"
		for _, w := range result.Warnings {
			text += "- " + w + "\n"
		}
	} else {
		text += "No recording warnings.\n"
	}
	if err = os.WriteFile(filepath.Join(session, "summary.md"), []byte(text), 0644); err != nil {
		return err
	}
	fmt.Printf("Report: %s\nSearches: %d; completed: %d; linked: %d; warnings: %d\n", filepath.Join(session, "summary.md"), result.Searches, result.Completed, result.Matched, len(result.Warnings))
	return nil
}
