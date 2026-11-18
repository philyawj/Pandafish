package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/philyawj/Pandafish/internal/telemetry"
)

type probeScore struct {
	CP      *int     `json:"cp"`
	Mate    *int     `json:"mate"`
	Depth   int      `json:"depth"`
	Nodes   uint64   `json:"nodes"`
	PV      []string `json:"pv"`
	Lower   bool     `json:"lowerbound"`
	Upper   bool     `json:"upperbound"`
	MultiPV int      `json:"multipv"`
}
type probeQuery struct {
	Position string             `json:"position"`
	Command  string             `json:"command"`
	Move     string             `json:"move"`
	MS       float64            `json:"response_ms"`
	Lines    []string           `json:"lines"`
	Score    probeScore         `json:"last_exact"`
	Metrics  *telemetry.Metrics `json:"metrics,omitempty"`
	Error    string             `json:"error,omitempty"`
}
type probeCase struct{ ID, Position, FEN string }
type probeTrial struct {
	Actor string     `json:"actor"`
	Trial int        `json:"trial"`
	Query probeQuery `json:"query"`
}
type probeNative struct {
	Case   probeCase    `json:"case"`
	Trials []probeTrial `json:"trials"`
	Error  string       `json:"error,omitempty"`
}

func parseProbeScore(line string) (probeScore, bool) {
	w := strings.Fields(line)
	s := probeScore{MultiPV: 1}
	if len(w) < 2 || w[0] != "info" || w[1] == "string" {
		return s, false
	}
	for i := 1; i < len(w); i++ {
		switch w[i] {
		case "depth":
			if i+1 < len(w) {
				s.Depth, _ = strconv.Atoi(w[i+1])
				i++
			}
		case "nodes":
			if i+1 < len(w) {
				s.Nodes, _ = strconv.ParseUint(w[i+1], 10, 64)
				i++
			}
		case "multipv":
			if i+1 < len(w) {
				s.MultiPV, _ = strconv.Atoi(w[i+1])
				i++
			}
		case "score":
			if i+2 < len(w) {
				value, err := strconv.Atoi(w[i+2])
				if err != nil {
					return s, false
				}
				if w[i+1] == "cp" {
					s.CP = &value
				} else if w[i+1] == "mate" {
					s.Mate = &value
				}
				i += 2
			}
		case "lowerbound":
			s.Lower = true
		case "upperbound":
			s.Upper = true
		case "pv":
			s.PV = append([]string(nil), w[i+1:]...)
			i = len(w)
		}
	}
	return s, (s.CP != nil || s.Mate != nil) && !s.Lower && !s.Upper
}
func (e *engine) probe(position, command string, timeout time.Duration) (result probeQuery) {
	result.Position, result.Command = position, command
	start := time.Now()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	if err := e.send(position); err != nil {
		result.Error = err.Error()
		return
	}
	if err := e.send(command); err != nil {
		result.Error = err.Error()
		return
	}
	for {
		select {
		case line, ok := <-e.lines:
			if !ok {
				result.Error = "engine exited"
				return
			}
			result.Lines = append(result.Lines, line)
			if score, exact := parseProbeScore(line); exact && score.MultiPV == 1 {
				result.Score = score
			}
			if strings.HasPrefix(line, telemetry.Prefix) {
				var m telemetry.Metrics
				if err := json.Unmarshal([]byte(strings.TrimPrefix(line, telemetry.Prefix)), &m); err != nil {
					result.Error = err.Error()
					return
				}
				result.Metrics = &m
			}
			if strings.HasPrefix(line, "bestmove ") {
				fields := strings.Fields(line)
				if len(fields) < 2 {
					result.Error = "missing bestmove"
					return
				}
				result.Move = fields[1]
				result.MS = float64(time.Since(start)) / float64(time.Millisecond)
				return
			}
		case <-timer.C:
			result.Error = "query response timeout"
			result.MS = float64(time.Since(start)) / float64(time.Millisecond)
			return
		}
	}
}
func probeGame(c probeCase) (*chess.Game, error) {
	position := c.Position
	if position == "" {
		position = "position fen " + c.FEN
	}
	w := strings.Fields(strings.TrimPrefix(position, "position "))
	if len(w) == 0 {
		return nil, fmt.Errorf("empty position")
	}
	g := chess.NewGame()
	next := 1
	if w[0] == "fen" {
		if len(w) < 7 {
			return nil, fmt.Errorf("short FEN")
		}
		option, err := chess.FEN(strings.Join(w[1:7], " "))
		if err != nil {
			return nil, err
		}
		g = chess.NewGame(option)
		next = 7
	} else if w[0] != "startpos" {
		return nil, fmt.Errorf("unsupported position")
	}
	if next < len(w) {
		if w[next] != "moves" {
			return nil, fmt.Errorf("missing moves marker")
		}
		for _, move := range w[next+1:] {
			if err := g.PushNotationMove(move, chess.UCINotation{}, nil); err != nil {
				return nil, err
			}
		}
	}
	return g, nil
}
func validProbeMove(g *chess.Game, move string) bool {
	moves := g.Position().ValidMoves()
	if len(moves) == 0 {
		return move == "0000" || move == "(none)"
	}
	for _, m := range moves {
		if m.String() == move {
			return true
		}
	}
	return false
}
func writeProbeJSON(path string, value any) error {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}
func probeDigest(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return hash(b).SHA, nil
}

func probeChoices(base, candidate, teacher, input, out string, budget, trials, workers int, cache string, requireTurns bool) error {
	if base == "" || candidate == "" || out == "" || budget < 1 || trials < 1 || workers < 1 {
		return fmt.Errorf("executables/output and positive budget/trials/workers required")
	}
	if err := os.Mkdir(out, 0755); err != nil {
		return err
	}
	var cases struct{ Cases []probeCase }
	if err := load(input, &cases); err != nil {
		return err
	}
	for i := range cases.Cases {
		if cases.Cases[i].Position == "" {
			cases.Cases[i].Position = "position fen " + cases.Cases[i].FEN
		}
		if _, err := probeGame(cases.Cases[i]); err != nil {
			return err
		}
	}
	identities := map[string]string{}
	for label, path := range map[string]string{"baseline": base, "candidate": candidate, "reference": teacher, "input": input} {
		h, err := probeDigest(path)
		if err != nil {
			return err
		}
		identities[label] = h
	}
	if err := writeProbeJSON(filepath.Join(out, "method.json"), map[string]any{"identities_sha256": identities, "milliseconds": budget, "trials": trials, "workers": workers, "reference_nodes": 100000, "reference_options": map[string]string{"Threads": "1", "Hash": "64", "Skill Level": "20", "UCI_LimitStrength": "false", "MultiPV": "1"}, "quality": "Per-position mean clamp(referenceCP-selectedCP-25,0,1000); retain mate cases separately; no forced candidate roots; native phase completes before reference queries."}); err != nil {
		return err
	}
	journal, err := os.OpenFile(filepath.Join(out, "native.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer journal.Close()
	encoder := json.NewEncoder(journal)
	jobs := make(chan int)
	done := make(chan probeNative)
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				c := cases.Cases[index]
				row := probeNative{Case: c}
				g, _ := probeGame(c)
				engines := map[string]*engine{}
				for actor, path := range map[string]string{"baseline": base, "candidate": candidate} {
					e, err := startEngine(path, nil)
					if err != nil {
						row.Error = err.Error()
						break
					}
					engines[actor] = e
				}
				if row.Error == "" {
					for trial := 0; trial < trials; trial++ {
						order := []string{"baseline", "candidate"}
						if (index+trial)%2 != 0 {
							order[0], order[1] = order[1], order[0]
						}
						for _, actor := range order {
							e := engines[actor]
							if err := e.send("ucinewgame"); err != nil {
								row.Error = err.Error()
								break
							}
							if err := e.send("isready"); err != nil {
								row.Error = err.Error()
								break
							}
							if _, err := e.wait("readyok", 5*time.Second); err != nil {
								row.Error = err.Error()
								break
							}
							q := e.probe(c.Position, fmt.Sprintf("go movetime %d", budget), time.Duration(budget)*time.Millisecond+5*time.Second)
							if q.Error == "" && !validProbeMove(g, q.Move) {
								q.Error = "illegal or invalid null move"
							}
							row.Trials = append(row.Trials, probeTrial{actor, trial, q})
						}
						if row.Error != "" {
							break
						}
					}
				}
				for _, e := range engines {
					e.close()
				}
				done <- row
			}
		}()
	}
	go func() {
		for i := range cases.Cases {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
		close(done)
	}()
	native := []probeNative{}
	for row := range done {
		native = append(native, row)
		if err := encoder.Encode(row); err != nil {
			return err
		}
		fmt.Printf("Native %d/%d\n", len(native), len(cases.Cases))
	}
	sort.Slice(native, func(i, j int) bool { return native[i].Case.ID < native[j].Case.ID })
	if err := writeProbeJSON(filepath.Join(out, "native.json"), native); err != nil {
		return err
	}
	for _, row := range native {
		if row.Error != "" {
			return fmt.Errorf("native initialization failed: %s", row.Error)
		}
		for _, trial := range row.Trials {
			if trial.Query.Error != "" {
				return fmt.Errorf("native query failed %s: %s", row.Case.ID, trial.Query.Error)
			}
		}
	}
	return probeReferences(teacher, out, budget, native, cache, identities["reference"], requireTurns)
}

type referenceRecord struct {
	Position string          `json:"position"`
	Move     *string         `json:"move"`
	Score    probeScore      `json:"last_exact"`
	Query    *probeQuery     `json:"query,omitempty"`
	Reused   bool            `json:"reused"`
	Original json.RawMessage `json:"original_cache_record,omitempty"`
}

func referenceKey(position, move string) string { return position + "\n" + move }
func probeReferences(teacher, out string, budget int, native []probeNative, cache, teacherHash string, requireTurns bool) error {
	references := map[string]referenceRecord{}
	if cache != "" {
		cacheFile := cache
		var data struct {
			Hash    string            `json:"stockfish_sha256"`
			Queries []json.RawMessage `json:"queries"`
		}
		info, err := os.Stat(cache)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(cache, "quality.json")); err != nil {
				return fmt.Errorf("reference cache run is incomplete")
			}
			var method struct {
				IDs   map[string]string `json:"identities_sha256"`
				Nodes int               `json:"reference_nodes"`
			}
			if err := load(filepath.Join(cache, "method.json"), &method); err != nil {
				return err
			}
			if method.Nodes != 100000 {
				return fmt.Errorf("reference cache node budget differs")
			}
			data.Hash = method.IDs["reference"]
			cacheFile = filepath.Join(cache, "references.jsonl")
			file, err := os.Open(cacheFile)
			if err != nil {
				return err
			}
			decoder := json.NewDecoder(file)
			for {
				var raw json.RawMessage
				err := decoder.Decode(&raw)
				if err == io.EOF {
					break
				}
				if err != nil {
					file.Close()
					return err
				}
				data.Queries = append(data.Queries, raw)
			}
			if err := file.Close(); err != nil {
				return err
			}
		} else if err := load(cache, &data); err != nil {
			return err
		}
		if data.Hash != teacherHash {
			return fmt.Errorf("reference cache engine hash differs")
		}
		for _, raw := range data.Queries {
			var row referenceRecord
			if err := json.Unmarshal(raw, &row); err != nil {
				return err
			}
			move := ""
			if row.Move != nil {
				move = *row.Move
			}
			row.Reused = true
			row.Original = raw
			references[referenceKey(row.Position, move)] = row
		}
		h, err := probeDigest(cacheFile)
		if err != nil {
			return err
		}
		if err := writeProbeJSON(filepath.Join(out, "reference-cache.json"), map[string]string{"source": cacheFile, "sha256": h}); err != nil {
			return err
		}
	}
	journal, err := os.OpenFile(filepath.Join(out, "references.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer journal.Close()
	encoder := json.NewEncoder(journal)
	e, err := startEngine(teacher, map[string]string{"Threads": "1", "Hash": "64", "Skill Level": "20", "UCI_LimitStrength": "false", "MultiPV": "1"})
	if err != nil {
		return err
	}
	defer e.close()
	used := map[string]referenceRecord{}
	get := func(position, move string) (referenceRecord, error) {
		key := referenceKey(position, move)
		if row, ok := used[key]; ok {
			return row, nil
		}
		row, ok := references[key]
		if !ok {
			if err := e.send("ucinewgame"); err != nil {
				return row, err
			}
			if err := e.send("isready"); err != nil {
				return row, err
			}
			if _, err := e.wait("readyok", 5*time.Second); err != nil {
				return row, err
			}
			command := "go nodes 100000"
			if move != "" {
				command += " searchmoves " + move
			}
			q := e.probe(position, command, 30*time.Second)
			row = referenceRecord{Position: position, Score: q.Score, Query: &q}
			if move != "" {
				copy := move
				row.Move = &copy
			}
		}
		// Journal before validation so scoreless, bounded or failed queries are
		// retained just as fully as accepted observations.
		if err := encoder.Encode(row); err != nil {
			return row, err
		}
		if row.Query != nil && row.Query.Error != "" {
			return row, fmt.Errorf("reference failed: %s", row.Query.Error)
		}
		if row.Score.CP == nil && row.Score.Mate == nil {
			return row, fmt.Errorf("reference lacks an exact score")
		}
		if row.Score.Lower || row.Score.Upper {
			return row, fmt.Errorf("bounded reference")
		}
		if move != "" && (len(row.Score.PV) == 0 || row.Score.PV[0] != move) {
			return row, fmt.Errorf("reference score belongs to another move")
		}
		used[key] = row
		return row, nil
	}
	type qualityCase struct {
		ID        string             `json:"id"`
		Turn      string             `json:"turn"`
		Reference probeScore         `json:"reference"`
		Loss      map[string]float64 `json:"mean_loss_cp,omitempty"`
		Risk      map[string]bool    `json:"losing_mate"`
		Terminal  bool               `json:"terminal"`
	}
	quality := []qualityCase{}
	timing := []string{}
	finiteMeans := map[string]float64{"baseline": 0, "candidate": 0}
	severe := map[string]int{"baseline": 0, "candidate": 0}
	byTurn := map[string]map[string][]float64{"w": {"baseline": {}, "candidate": {}}, "b": {"baseline": {}, "candidate": {}}}
	newRisk := []string{}
	finiteCases := 0
	for index, row := range native {
		g, _ := probeGame(row.Case)
		turn := "w"
		if g.Position().Turn() == chess.Black {
			turn = "b"
		}
		q := qualityCase{ID: row.Case.ID, Turn: turn, Risk: map[string]bool{}}
		for _, trial := range row.Trials {
			if trial.Query.MS > float64(budget+50) {
				timing = append(timing, row.Case.ID+"/"+trial.Actor)
			}
		}
		if len(g.Position().ValidMoves()) == 0 {
			q.Terminal = true
			quality = append(quality, q)
			continue
		}
		root, err := get(row.Case.Position, "")
		if err != nil {
			return err
		}
		q.Reference = root.Score
		losses := map[string][]float64{"baseline": {}, "candidate": {}}
		finite := root.Score.CP != nil
		for _, trial := range row.Trials {
			ref, err := get(row.Case.Position, trial.Query.Move)
			if err != nil {
				return err
			}
			if ref.Score.Mate != nil && *ref.Score.Mate < 0 {
				q.Risk[trial.Actor] = true
			}
			if ref.Score.CP == nil {
				finite = false
			}
			if root.Score.CP != nil && ref.Score.CP != nil {
				losses[trial.Actor] = append(losses[trial.Actor], float64(min(1000, max(0, *root.Score.CP-*ref.Score.CP-25))))
			}
		}
		if (root.Score.CP != nil || (root.Score.Mate != nil && *root.Score.Mate >= 0)) && q.Risk["candidate"] && !q.Risk["baseline"] {
			newRisk = append(newRisk, row.Case.ID)
		}
		if finite {
			finiteCases++
			q.Loss = map[string]float64{}
			for _, actor := range []string{"baseline", "candidate"} {
				sum := 0.
				for _, loss := range losses[actor] {
					sum += loss
				}
				mean := sum / float64(len(losses[actor]))
				q.Loss[actor] = mean
				finiteMeans[actor] += mean
				if mean > 300 {
					severe[actor]++
				}
				byTurn[turn][actor] = append(byTurn[turn][actor], mean)
			}
		}
		quality = append(quality, q)
		if err := writeProbeJSON(filepath.Join(out, "quality-checkpoint.json"), quality); err != nil {
			return err
		}
		fmt.Printf("Reference %d/%d\n", index+1, len(native))
	}
	for actor := range finiteMeans {
		if finiteCases > 0 {
			finiteMeans[actor] /= float64(finiteCases)
		}
	}
	turnMeans := map[string]map[string]float64{}
	for turn, actors := range byTurn {
		turnMeans[turn] = map[string]float64{}
		for actor, losses := range actors {
			sum := 0.
			for _, v := range losses {
				sum += v
			}
			if len(losses) > 0 {
				turnMeans[turn][actor] = sum / float64(len(losses))
			}
		}
	}
	passed := finiteCases > 0 && finiteMeans["candidate"] <= finiteMeans["baseline"]+math.Max(5, .1*finiteMeans["baseline"]) && severe["candidate"] <= severe["baseline"]+2 && len(newRisk) == 0 && len(timing) == 0
	turnGates := map[string]bool{}
	for _, turn := range []string{"w", "b"} {
		values := turnMeans[turn]
		turnGates[turn] = len(byTurn[turn]["baseline"]) > 0 && values["candidate"] <= values["baseline"]+math.Max(5, .1*values["baseline"])
		if requireTurns && !turnGates[turn] {
			passed = false
		}
	}
	report := map[string]any{"finite_cases": finiteCases, "mean_loss_cp": finiteMeans, "severe_over_300": severe, "new_candidate_losing_mate_cases": newRisk, "timing_overruns": timing, "by_turn_mean_loss_cp": turnMeans, "turn_gates": turnGates, "turn_gates_required": requireTurns, "passed": passed, "cases": quality}
	if err := writeProbeJSON(filepath.Join(out, "quality.json"), report); err != nil {
		return err
	}
	return printJSON(map[string]any{"finite_cases": finiteCases, "mean_loss_cp": finiteMeans, "severe_over_300": severe, "new_candidate_losing_mate_cases": newRisk, "timing_overruns": len(timing), "passed": passed})
}
