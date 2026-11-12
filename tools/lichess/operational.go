package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/corentings/chess/v2"
)

type gameRecord struct {
	Game          string `json:"game"`
	White         string `json:"white"`
	Black         string `json:"black"`
	BotColor      string `json:"bot_color"`
	Opponent      string `json:"opponent"`
	OpponentTitle string `json:"opponent_title"`
	Result        string `json:"result"`
	Outcome       string `json:"outcome"`
	Termination   string `json:"termination"`
	Event         string `json:"event"`
	TimeControl   string `json:"time_control"`
	Start         string `json:"start_utc"`
	BotRating     string `json:"bot_rating"`
	RatingChange  string `json:"rating_change"`
	Plies         int    `json:"plies"`
}
type clientIncident struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}

var logLevel = regexp.MustCompile(`\b(INFO|ERROR|WARNING|DEBUG|CRITICAL)\b`)

func classifyClientLog(path string) ([]clientIncident, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	var incidents []clientIncident
	block := ""
	errorBlock := false
	flush := func() {
		if !errorBlock {
			return
		}
		text := strings.Join(strings.Fields(block), " ")
		kind := "other_client_error"
		lower := strings.ToLower(text)
		if strings.Contains(lower, "bot.vsbot.day") && strings.Contains(lower, "'opponent_is_rate_limited': true") {
			kind = "opponent_daily_game_cap"
		} else if strings.Contains(lower, "bot.vsbot.day") && strings.Contains(lower, "'bot_is_rate_limited': true") {
			kind = "own_daily_game_cap"
		} else if strings.Contains(lower, "100 games against other bots") {
			kind = "daily_game_cap_unspecified"
		} else if strings.Contains(lower, "429") || strings.Contains(lower, "too many requests") {
			kind = "http_rate_limit"
		} else if strings.Contains(lower, "engine") {
			kind = "engine_related_error"
		}
		incidents = append(incidents, clientIncident{kind, text})
	}
	for scanner.Scan() {
		line := scanner.Text()
		_, body, ok := strings.Cut(line, " ")
		if !ok {
			body = line
		}
		// The client prints control-stream exceptions directly, without a log level.
		// Count the reconnect marker once, not each wrapped traceback line.
		if strings.HasPrefix(strings.TrimSpace(body), "Control stream error, reconnecting:") {
			flush()
			errorBlock = false
			block = ""
			incidents = append(incidents, clientIncident{"control_stream_reconnect", strings.TrimSpace(body)})
			continue
		}
		if level := logLevel.FindString(body); level != "" {
			flush()
			errorBlock = level == "ERROR" || level == "CRITICAL"
			block = body
		} else if errorBlock {
			block += " " + body
		}
	}
	flush()
	return incidents, scanner.Err()
}
func operationalReport(raw, output, bot string) error {
	data, err := os.ReadFile(filepath.Join(raw, "manifest.json"))
	if err != nil {
		return err
	}
	var m manifest
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}
	paths, err := filepath.Glob(filepath.Join(raw, "games", "*.pgn"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	var games []gameRecord
	outcomes := map[string]int{}
	controls := map[string]int{}
	terminations := map[string]int{}
	opponents := map[string]bool{}
	for _, path := range paths {
		pgn, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		opt, e := chess.PGN(strings.NewReader(string(pgn)))
		if e != nil {
			return e
		}
		g := chess.NewGame(opt)
		color := "White"
		if !strings.EqualFold(g.GetTagPair(color), bot) {
			color = "Black"
		}
		if !strings.EqualFold(g.GetTagPair(color), bot) {
			return fmt.Errorf("bot %q missing from %s", bot, filepath.Base(path))
		}
		other := "Black"
		winResult := "1-0"
		if color == "Black" {
			other = "White"
			winResult = "0-1"
		}
		result := g.GetTagPair("Result")
		outcome := "unfinished"
		switch result {
		case "1/2-1/2":
			outcome = "draw"
		case "1-0", "0-1":
			outcome = "loss"
			if result == winResult {
				outcome = "win"
			}
		}
		row := gameRecord{Game: g.GetTagPair("Site"), White: g.GetTagPair("White"), Black: g.GetTagPair("Black"), BotColor: color, Opponent: g.GetTagPair(other), OpponentTitle: g.GetTagPair(other + "Title"), Result: result, Outcome: outcome, Termination: g.GetTagPair("Termination"), Event: g.GetTagPair("Event"), TimeControl: g.GetTagPair("TimeControl"), Start: g.GetTagPair("UTCDate") + " " + g.GetTagPair("UTCTime"), BotRating: g.GetTagPair(color + "Elo"), RatingChange: g.GetTagPair(color + "RatingDiff"), Plies: len(g.Moves())}
		games = append(games, row)
		outcomes[outcome]++
		controls[row.TimeControl]++
		terminations[row.Termination+" / "+outcome]++
		opponents[row.Opponent] = true
	}
	sort.Slice(games, func(i, j int) bool { return games[i].Start < games[j].Start })
	capture, err := loadCapture(filepath.Join(raw, "moves.jsonl"))
	if err != nil {
		return err
	}
	var longest *searchRecord
	fixed := map[string]int{}
	clockCases, lowClockCases := 0, 0
	maxOverrun, maxProtocolOverhead := 0.0, 0.0
	totalSearchMS := 0.0
	for _, r := range capture.All {
		if longest == nil || r.ResponseMS > longest.ResponseMS {
			longest = r
		}
		if r.Metrics == nil {
			continue
		}
		if r.Metrics.Timed {
			maxOverrun = max(maxOverrun, r.Metrics.ElapsedMS-r.Metrics.HardMS)
		}
		maxProtocolOverhead = max(maxProtocolOverhead, r.ResponseMS-r.Metrics.ElapsedMS)
		totalSearchMS += r.Metrics.ElapsedMS
		fields := strings.Fields(r.Go)
		if strings.Contains(r.Go, "movetime") {
			fixed[r.Go]++
		} else {
			clockCases++
			clockKey := "wtime"
			fenFields := strings.Fields(r.FEN)
			if len(fenFields) < 2 {
				return fmt.Errorf("completed search has no valid FEN")
			}
			if fenFields[1] == "b" {
				clockKey = "btime"
			}
			for i, key := range fields {
				if key == clockKey && i+1 < len(fields) {
					v, _ := strconv.ParseInt(fields[i+1], 10, 64)
					if v < 10000 {
						lowClockCases++
					}
				}
			}
		}
	}
	incidents, err := classifyClientLog(filepath.Join(raw, "lichess-bot.log"))
	if err != nil {
		return err
	}
	incidentCounts := map[string]int{}
	for _, i := range incidents {
		incidentCounts[i.Kind]++
	}
	linked, repeated := 0, 0
	for _, attempts := range capture.ByPly {
		linked += len(attempts)
		if len(attempts) > 1 {
			repeated++
		}
	}
	summary := struct {
		CompletedSearches       int            `json:"completed_searches"`
		LinkedSearches          int            `json:"linked_searches"`
		UniqueLinkedPositions   int            `json:"unique_linked_positions"`
		RepeatedLinkedPositions int            `json:"repeated_linked_positions"`
		UnlinkedSearches        int            `json:"unlinked_completed_searches"`
		Games                   []gameRecord   `json:"games"`
		Outcomes                map[string]int `json:"outcomes"`
		Controls                map[string]int `json:"time_controls"`
		Terminations            map[string]int `json:"terminations"`
		UniqueOpponents         int            `json:"unique_opponents"`
		FixedRequests           map[string]int `json:"fixed_requests"`
		ClockRequests           int            `json:"clock_requests"`
		Below10s                int            `json:"requests_with_own_clock_below_10s"`
		MaxOverrunMS            float64        `json:"max_hard_overrun_ms"`
		MaxOutputOverheadMS     float64        `json:"max_response_minus_search_ms"`
		SearchWallSeconds       float64        `json:"total_search_wall_seconds"`
		ClientErrors            map[string]int `json:"client_error_blocks"`
		Longest                 *searchRecord  `json:"longest_search"`
	}{len(capture.All), linked, len(capture.ByPly), repeated, len(capture.All) - linked, games, outcomes, controls, terminations, len(opponents), fixed, clockCases, lowClockCases, maxOverrun, maxProtocolOverhead, totalSearchMS / 1000, incidentCounts, longest}
	if err = writeJSON(filepath.Join(output, "operations.json"), summary); err != nil {
		return err
	}
	if err = writeJSON(filepath.Join(output, "client-incidents.json"), incidents); err != nil {
		return err
	}
	duration := time.Duration(0)
	if m.Finished != nil {
		duration = m.Finished.Sub(m.Started)
	}
	text := fmt.Sprintf("# Live session operations\n\nEngine: %s; source revision `%s`; clean engine build: %t.\n\nSession duration: %s. Platform: %s.\n\n%d games against %d distinct opponents: %d wins, %d draws, %d losses, %d unfinished. These are selected live opponents and mixed clocks, not a controlled strength estimate.\n\n| Time control | Games |\n| --- | --- |\n", m.EngineID, m.EngineRevision, m.EngineDirty == "false", duration.Round(time.Second), m.Platform, len(games), len(opponents), outcomes["win"], outcomes["draw"], outcomes["loss"], outcomes["unfinished"])
	keys := []string{}
	for k := range controls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		text += fmt.Sprintf("| %s | %d |\n", k, controls[k])
	}
	text += fmt.Sprintf("\nMaximum hard-budget overrun: %.6f ms. Maximum response minus search time: %.3f ms. %d clock-based requests had less than 10 seconds left on the moving side's supplied clock. Search wall time totals %.1f seconds; this is not measured CPU utilization.\n\nFixed requests: `%v`. Clock requests: %d.\n", maxOverrun, maxProtocolOverhead, lowClockCases, totalSearchMS/1000, fixed, clockCases)
	text += fmt.Sprintf("\nCompleted searches: %d; linked searches: %d; unique linked positions: %d; positions with repeated searches: %d; unlinked completed searches: %d. Timing totals include all completed searches, including repeated and unlinked attempts.\n", len(capture.All), linked, len(capture.ByPly), repeated, len(capture.All)-linked)
	if longest != nil && longest.Metrics != nil && longest.Ply == nil {
		text += fmt.Sprintf("\nLongest response: %.3f ms in unlinked process `%s`. Request: `%s`; allocated hard budget %.3f ms.\n", longest.ResponseMS, longest.Process, longest.Go, longest.Metrics.HardMS)
	}
	if longest != nil && longest.Metrics != nil && longest.Ply != nil {
		text += fmt.Sprintf("\nLongest response: %.3f ms at [%s](%s), zero-based ply %d. Request: `%s`; allocated hard budget %.3f ms. Compare the request and budget before calling a long search a stall.\n", longest.ResponseMS, longest.Game, longest.Game, *longest.Ply, longest.Go, longest.Metrics.HardMS)
	}
	text += "\n## Client errors\n\nCounts refer to grouped ERROR/CRITICAL records plus explicit unlabelled control-stream reconnect markers, not keyword hits. Other severity levels and individual wrapped traceback lines are not counted as separate failures. Opponent daily-cap responses describe the challenged account, not an engine crash.\n\n"
	keys = nil
	for k := range incidentCounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		text += fmt.Sprintf("- %s: %d\n", k, incidentCounts[k])
	}
	text += "\n## Games\n\n| Game | Bot color | Opponent | Clock | Result | Termination |\n| --- | --- | --- | --- | --- | --- |\n"
	for _, g := range games {
		text += fmt.Sprintf("| [game](%s) | %s | %s | %s | %s | %s |\n", g.Game, g.BotColor, g.Opponent, g.TimeControl, g.Outcome, g.Termination)
	}
	return os.WriteFile(filepath.Join(output, "OPERATIONS.md"), []byte(text), 0644)
}
