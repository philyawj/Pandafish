package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/corentings/chess/v2"
)

type processOptions struct {
	Session, Archive, Stockfish, Bot, Focus string
	Profile                                 string
	Nodes                                   int64
	Limit                                   int
}
type analysisMethod struct {
	Schema          int    `json:"schema"`
	Engine          string `json:"engine"`
	EngineSHA256    string `json:"engine_sha256"`
	InventorySHA256 string `json:"inventory_sha256"`
	Nodes           int64  `json:"nodes_per_search"`
	Threads         int    `json:"threads"`
	HashMB          int    `json:"hash_mb"`
	MultiPV         int    `json:"multipv"`
	Skill           int    `json:"skill"`
	Bot             string `json:"bot"`
	FocusSHA256     string `json:"focus_sha256,omitempty"`
	Perspective     string `json:"score_perspective"`
	Comparison      string `json:"comparison"`
}
type analysedMove struct {
	Game           string          `json:"game"`
	Ply            int             `json:"ply"`
	Side           string          `json:"side"`
	IsBot          bool            `json:"is_bot"`
	FEN            string          `json:"fen"`
	Position       string          `json:"position"`
	Played         string          `json:"played"`
	PlayedSAN      string          `json:"played_san"`
	Best           sfResult        `json:"best"`
	PlayedResult   sfResult        `json:"played_result"`
	Reused         bool            `json:"played_search_reused"`
	CPLoss         *int            `json:"cp_loss,omitempty"`
	RawCPDelta     *int            `json:"raw_cp_delta,omitempty"`
	Classification string          `json:"classification"`
	Capture        *searchRecord   `json:"capture,omitempty"`
	Captures       []*searchRecord `json:"captures,omitempty"` // Multiple attempts, without guessing which was submitted.
}
type gameAnalysis struct {
	Game      string         `json:"game"`
	PGNSHA256 string         `json:"pgn_sha256"`
	Complete  bool           `json:"complete"`
	Moves     []analysedMove `json:"moves"`
}
type candidate struct {
	Game           string `json:"game"`
	Ply            int    `json:"ply"`
	FEN            string `json:"fen"`
	Position       string `json:"position"`
	Played         string `json:"played"`
	Best           string `json:"best"`
	Classification string `json:"classification"`
	CPLoss         *int   `json:"cp_loss,omitempty"`
}

func compareScores(best, played sfScore) (*int, *int, string) {
	if best.Kind == "cp" && played.Kind == "cp" {
		raw := best.Value - played.Value
		loss := max(0, raw)
		category := "small_or_equal"
		if raw < 0 {
			category = "search_disagreement"
		} else if loss >= 300 {
			category = "loss_300cp_plus"
		} else if loss >= 100 {
			category = "loss_100cp_plus"
		}
		return &loss, &raw, category
	}
	if best.Kind == "mate" && best.Value > 0 && !(played.Kind == "mate" && played.Value > 0) {
		return nil, nil, "missed_forced_mate"
	}
	if played.Kind == "mate" && played.Value <= 0 && !(best.Kind == "mate" && best.Value <= 0) {
		return nil, nil, "allowed_forced_mate"
	}
	return nil, nil, "mate_score_comparison"
}

type captureIndex struct {
	All   []*searchRecord
	ByPly map[string][]*searchRecord
}

func (index *captureIndex) attach(row *analysedMove, key string) {
	attempts := index.ByPly[key]
	if len(attempts) == 1 {
		row.Capture = attempts[0]
	} else if len(attempts) > 1 {
		row.Captures = attempts
	}
}

func loadCapture(path string) (*captureIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 4*1024*1024)
	rows := &captureIndex{ByPly: map[string][]*searchRecord{}}
	for scanner.Scan() {
		r := new(searchRecord)
		if err = json.Unmarshal(scanner.Bytes(), r); err != nil {
			return nil, err
		}
		if r.Status != "completed" {
			continue
		}
		rows.All = append(rows.All, r)
		if r.Game != "" && r.Ply != nil {
			key := fmt.Sprintf("%s/%d", r.Game, *r.Ply)
			rows.ByPly[key] = append(rows.ByPly[key], r)
		}
	}
	return rows, scanner.Err()
}
func runProcess(o processOptions) error {
	if o.Session == "" || o.Archive == "" || o.Bot == "" || o.Nodes < 1 || o.Nodes > 100000000 || o.Limit < 0 {
		return fmt.Errorf("process requires --session, --archive, --bot and nodes 1..100000000")
	}
	if err := archiveSession(o.Session, o.Archive); err != nil {
		return err
	}
	raw := filepath.Join(o.Archive, "raw")
	enginePath, err := exec.LookPath(o.Stockfish)
	if err != nil {
		return err
	}
	engineBytes, err := os.ReadFile(enginePath)
	if err != nil {
		return err
	}
	sf, err := openStockfish(enginePath)
	if err != nil {
		return err
	}
	defer sf.close()
	indexBytes, err := os.ReadFile(filepath.Join(o.Archive, "inventory.json"))
	if err != nil {
		return err
	}
	method := analysisMethod{Schema: 1, Engine: sf.name, EngineSHA256: digest(engineBytes), InventorySHA256: digest(indexBytes), Nodes: o.Nodes, Threads: 1, HashMB: 64, MultiPV: 1, Skill: 20, Bot: o.Bot, Perspective: "root mover; positive favors the player to move; cp and mate remain separate", Comparison: "unrestricted root versus played-move restricted root at the same node budget; fresh hash and ucinewgame per search"}
	focus := map[string]bool{}
	if o.Focus != "" {
		data, e := os.ReadFile(o.Focus)
		if e != nil {
			return e
		}
		var candidates []candidate
		if e = json.Unmarshal(data, &candidates); e != nil {
			return e
		}
		if len(candidates) == 0 {
			return fmt.Errorf("empty focus list")
		}
		for _, c := range candidates {
			focus[fmt.Sprintf("%s/%d", c.Game, c.Ply)] = true
		}
		method.FocusSHA256 = digest(data)
	}
	profile := fmt.Sprintf("stockfish-n%d", o.Nodes)
	if o.Focus != "" {
		profile += "-focus"
	}
	if o.Profile != "" {
		if filepath.Base(o.Profile) != o.Profile || o.Profile == "." || o.Profile == ".." {
			return fmt.Errorf("profile must be a directory name")
		}
		profile = o.Profile
	}
	output := filepath.Join(o.Archive, "analysis", profile)
	if err = os.MkdirAll(output, 0755); err != nil {
		return err
	}
	methodPath := filepath.Join(output, "method.json")
	if existing, e := os.ReadFile(methodPath); e == nil {
		var previous analysisMethod
		if json.Unmarshal(existing, &previous) != nil || !reflect.DeepEqual(previous, method) {
			return fmt.Errorf("analysis settings/input changed; use a different node profile/archive instead of overwriting")
		}
	} else if !os.IsNotExist(e) {
		return e
	}
	if err = writeJSON(methodPath, method); err != nil {
		return err
	}
	if err := operationalReport(raw, o.Archive, o.Bot); err != nil {
		return err
	}
	capture, err := loadCapture(filepath.Join(raw, "moves.jsonl"))
	if err != nil {
		return err
	}
	paths, err := filepath.Glob(filepath.Join(raw, "games", "*.pgn"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	processed, expected := 0, 0
	var reports []gameAnalysis
	started := time.Now()
	for _, path := range paths {
		pgn, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		option, err := chess.PGN(strings.NewReader(string(pgn)))
		if err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(path), err)
		}
		game := chess.NewGame(option)
		site := game.GetTagPair("Site")
		id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		target := filepath.Join(output, id+".json")
		result := gameAnalysis{Game: site, PGNSHA256: digest(pgn), Moves: []analysedMove{}}
		if cached, e := os.ReadFile(target); e == nil {
			if json.Unmarshal(cached, &result) != nil || result.PGNSHA256 != digest(pgn) || result.Game != site {
				return fmt.Errorf("invalid/stale checkpoint %s", id)
			}
		} else if !os.IsNotExist(e) {
			return e
		}
		byPly := map[int]analysedMove{}
		for _, m := range result.Moves {
			if _, exists := byPly[m.Ply]; exists {
				return fmt.Errorf("duplicate checkpoint ply")
			}
			byPly[m.Ply] = m
		}
		positions, moves := game.Positions(), game.Moves()
		history := []string{}
		result.Moves = nil
		gameExpected := 0
		for ply, move := range moves {
			key := fmt.Sprintf("%s/%d", site, ply)
			wanted := len(focus) == 0 || focus[key]
			if wanted {
				expected++
				gameExpected++
				position := "position fen " + positions[0].String()
				if len(history) > 0 {
					position += " moves " + strings.Join(history, " ")
				}
				if cached, exists := byPly[ply]; exists {
					if cached.Position != position || cached.Played != move.String() || cached.Best.Score.Kind == "" || cached.PlayedResult.Score.Kind == "" {
						return fmt.Errorf("checkpoint position differs in %s ply %d", id, ply)
					}
					// Preserve the original unrestricted output while repairing a bound-only final choice.
					if len(cached.Best.PV) == 0 || cached.Best.PV[0] != cached.Best.Bestmove {
						repaired, e := sf.analyse(position, cached.Best.Bestmove, o.Nodes)
						if e != nil {
							return e
						}
						repaired.SelectionOutput = cached.Best.Output
						repaired.ReferenceRechecked = true
						cached.Best = repaired
						if cached.Reused {
							cached.PlayedResult = repaired
						}
					}
					if len(cached.PlayedResult.PV) == 0 || cached.PlayedResult.PV[0] != cached.PlayedResult.Bestmove {
						return fmt.Errorf("invalid played-move checkpoint")
					}
					cached.Capture, cached.Captures = nil, nil
					capture.attach(&cached, key)
					cached.CPLoss, cached.RawCPDelta, cached.Classification = compareScores(cached.Best.Score, cached.PlayedResult.Score)
					result.Moves = append(result.Moves, cached)
					processed++
				} else if o.Limit == 0 || processed < o.Limit {
					best, e := sf.analyse(position, "", o.Nodes)
					if e != nil {
						return e
					}
					// Validate Stockfish's reference move before saving a comparison.
					ref := chess.NewGame()
					fen, e := chess.FEN(positions[ply].String())
					if e != nil {
						return e
					}
					ref = chess.NewGame(fen)
					if e = ref.PushNotationMove(best.Bestmove, chess.UCINotation{}, nil); e != nil {
						return fmt.Errorf("illegal Stockfish reference move")
					}
					played, reused := best, best.Bestmove == move.String()
					if !reused {
						played, e = sf.analyse(position, move.String(), o.Nodes)
						if e != nil {
							return e
						}
					}
					cp, delta, category := compareScores(best.Score, played.Score)
					side := "White"
					if positions[ply].Turn() == chess.Black {
						side = "Black"
					}
					row := analysedMove{Game: site, Ply: ply, Side: side, IsBot: strings.EqualFold(game.GetTagPair(side), o.Bot), FEN: positions[ply].String(), Position: position, Played: move.String(), PlayedSAN: chess.AlgebraicNotation{}.Encode(positions[ply], move), Best: best, PlayedResult: played, Reused: reused, CPLoss: cp, RawCPDelta: delta, Classification: category}
					capture.attach(&row, key)
					result.Moves = append(result.Moves, row)
					processed++
					result.Complete = false
					if e = writeJSON(target, result); e != nil {
						return e
					}
				}
			}
			history = append(history, move.String())
		}
		result.Complete = len(result.Moves) == gameExpected
		if err = writeJSON(target, result); err != nil {
			return err
		}
		reports = append(reports, result)
		fmt.Printf("%s: %d/%d positions; session %d evaluated; %.1fs elapsed\n", id, len(result.Moves), gameExpected, processed, time.Since(started).Seconds())
	}
	if len(focus) > 0 && expected != len(focus) {
		return fmt.Errorf("focus entries not all found in PGNs")
	}
	return analysisSummary(output, method, reports, processed, expected)
}
func analysisSummary(output string, method analysisMethod, reports []gameAnalysis, processed, expected int) error {
	counts := map[string]int{}
	var own []analysedMove
	losses := []int{}
	var candidates []candidate
	for _, report := range reports {
		for _, r := range report.Moves {
			if !r.IsBot {
				continue
			}
			own = append(own, r)
			counts[r.Classification]++
			if r.CPLoss != nil {
				losses = append(losses, *r.CPLoss)
			}
			if r.Classification == "missed_forced_mate" || r.Classification == "allowed_forced_mate" || (r.CPLoss != nil && *r.CPLoss >= 100) {
				candidates = append(candidates, candidate{r.Game, r.Ply, r.FEN, r.Position, r.Played, r.Best.Bestmove, r.Classification, r.CPLoss})
			}
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if (a.CPLoss == nil) != (b.CPLoss == nil) {
			return a.CPLoss == nil
		}
		if a.CPLoss != nil && b.CPLoss != nil && *a.CPLoss != *b.CPLoss {
			return *a.CPLoss > *b.CPLoss
		}
		if a.Game != b.Game {
			return a.Game < b.Game
		}
		return a.Ply < b.Ply
	})
	if err := writeJSON(filepath.Join(output, "candidates.json"), candidates); err != nil {
		return err
	}
	// A deterministic small follow-up set: at most one position per game, up to 20.
	selected := []candidate{}
	seen := map[string]bool{}
	// Prioritize the first substantial loss of a roughly playable position.
	chronological := append([]analysedMove(nil), own...)
	sort.Slice(chronological, func(i, j int) bool {
		if chronological[i].Game != chronological[j].Game {
			return chronological[i].Game < chronological[j].Game
		}
		return chronological[i].Ply < chronological[j].Ply
	})
	for _, r := range chronological {
		if seen[r.Game] || r.Best.Score.Kind != "cp" || r.Best.Score.Value < -100 {
			continue
		}
		if !((r.CPLoss != nil && *r.CPLoss >= 200 && r.PlayedResult.Score.Value <= -100) || r.Classification == "allowed_forced_mate") {
			continue
		}
		selected = append(selected, candidate{r.Game, r.Ply, r.FEN, r.Position, r.Played, r.Best.Bestmove, r.Classification, r.CPLoss})
		seen[r.Game] = true
		if len(selected) == 20 {
			break
		}
	}
	for _, c := range candidates {
		if len(selected) >= 20 {
			break
		}
		if !seen[c.Game] {
			selected = append(selected, c)
			seen[c.Game] = true
		}
	}

	if err := writeJSON(filepath.Join(output, "confirmation.json"), selected); err != nil {
		return err
	}
	summary := struct {
		Method           analysisMethod `json:"method"`
		Processed        int            `json:"processed_plies"`
		Expected         int            `json:"expected_plies"`
		BotMoves         int            `json:"bot_moves"`
		Classifications  map[string]int `json:"classifications"`
		MeanCPLossCPOnly float64        `json:"mean_cp_loss_cp_only"`
		CPCases          int            `json:"cp_cases"`
		Candidates       int            `json:"candidate_positions"`
	}{Method: method, Processed: processed, Expected: expected, BotMoves: len(own), Classifications: counts, CPCases: len(losses), Candidates: len(candidates)}
	for _, v := range losses {
		summary.MeanCPLossCPOnly += float64(v)
	}
	if len(losses) > 0 {
		summary.MeanCPLossCPOnly /= float64(len(losses))
	}
	if method.FocusSHA256 != "" {
		if err := writeRegressions(output, method, own); err != nil {
			return err
		}
	}
	if err := writeJSON(filepath.Join(output, "summary.json"), summary); err != nil {
		return err
	}
	text := fmt.Sprintf("# Stockfish analysis\n\n%s; %d nodes/search; one thread; 64 MiB hash; Skill Level 20; MultiPV 1.\n\nProcessed %d/%d played plies, including %d by %s.\n\n| Classification | Bot moves |\n| --- | --- |\n", method.Engine, method.Nodes, processed, expected, len(own), method.Bot)
	keys := []string{}
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		text += fmt.Sprintf("| %s | %d |\n", k, counts[k])
	}
	text += fmt.Sprintf("\nMean centipawn loss among the %d cp-to-cp comparisons: %.1f cp. Mate transitions are separate, not assigned artificial centipawns. These finite-budget comparisons are estimates, not Lichess accuracy or an Elo estimate. Negative raw deltas indicate search disagreement and remain recorded.\n\nFull per-game JSON includes root history, scores, PVs and raw Stockfish output. Scores favor the root mover. Each independent query starts with an empty hash and ucinewgame. Played moves use a restricted-root search at the same node budget unless the unrestricted best move already matches.\n\n%d candidate positions are saved in candidates.json. confirmation.json selects at most one per game, up to 20, for a larger-budget check. It first prefers the earliest substantial deterioration from a roughly playable position (reference at least -100 cp, loss at least 200 cp with the played result at most -100 cp, or newly allowed mate), then fills remaining slots from the ranked candidates. Selection is based on this survey and is not a random sample.\n", len(losses), summary.MeanCPLossCPOnly, len(candidates))
	if processed != expected {
		text += "\nThis run is partial. Repeating the same command resumes checkpoints.\n"
	}
	if err := os.WriteFile(filepath.Join(output, "README.md"), []byte(text), 0644); err != nil {
		return err
	}
	fmt.Printf("Analysis summary: %s (%s positions)\n", filepath.Join(output, "README.md"), strconv.Itoa(processed))
	return nil
}

type regressionCase struct {
	Game        string  `json:"game"`
	Ply         int     `json:"ply"`
	Split       string  `json:"split"`
	FEN         string  `json:"fen"`
	Position    string  `json:"position"`
	Played      string  `json:"played"`
	Reference   string  `json:"reference"`
	BestScore   sfScore `json:"best_score"`
	PlayedScore sfScore `json:"played_score"`
	CPLoss      int     `json:"estimated_cp_loss"`
}

func writeRegressions(output string, method analysisMethod, rows []analysedMove) error {
	cases := []regressionCase{}
	for _, r := range rows {
		if r.CPLoss == nil || *r.CPLoss < 100 || r.Best.Score.Kind != "cp" || r.Best.Score.Value < -150 || r.Best.Score.Value > 200 {
			continue
		}
		cases = append(cases, regressionCase{Game: r.Game, Ply: r.Ply, FEN: r.FEN, Position: r.Position, Played: r.Played, Reference: r.Best.Bestmove, BestScore: r.Best.Score, PlayedScore: r.PlayedResult.Score, CPLoss: *r.CPLoss})
	}
	sort.Slice(cases, func(i, j int) bool {
		if cases[i].Game != cases[j].Game {
			return cases[i].Game < cases[j].Game
		}
		return cases[i].Ply < cases[j].Ply
	})
	for i := range cases {
		cases[i].Split = "development"
		if i%2 == 1 {
			cases[i].Split = "holdout"
		}
	}
	return writeJSON(filepath.Join(output, "regression-candidates.json"), struct {
		Method analysisMethod   `json:"method"`
		Note   string           `json:"note"`
		Cases  []regressionCase `json:"cases"`
	}{method, "Public reference positions; the recorded engine is not expected to pass improvement cases retroactively. Reference moves can have equally good alternatives. The deterministic split is for later engine development; it is not a hidden or statistically independent test set.", cases})
}
