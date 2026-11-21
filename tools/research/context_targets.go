package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/corentings/chess/v2"
)

type contextTarget struct {
	ID, RootID, Move, Split, Turn, Component, FEN, Position, Key string
	Base                                                         int
	Learnable                                                    bool
	Original, Direct                                             label
	position                                                     *chess.Position
}
type contextLabelMethod struct {
	Teacher        string               `json:"teacher_sha256"`
	Sources        map[string]fileEntry `json:"sources"`
	Diagnostic     []string             `json:"diagnostic_roots"`
	Nodes, Workers int
}
type contextQuery struct {
	Key   string
	Query probeQuery
}

func fixedContextLabel(base int) label {
	if base > 900000 {
		return label{Mate: contextIntValue(1)}
	}
	if base < -900000 {
		return label{Mate: contextIntValue(-1)}
	}
	return label{CP: contextIntValue(base)}
}
func contextIntValue(v int) *int { return &v }
func labelFromContextQuery(q probeQuery, turn, rootTurn string) label {
	sign := 1
	if turn != rootTurn {
		sign = -1
	}
	var result label
	if q.Score.CP != nil {
		result.CP = contextIntValue(sign * *q.Score.CP)
	}
	if q.Score.Mate != nil {
		result.Mate = contextIntValue(sign * *q.Score.Mate)
	}
	return result
}
func contextLabelLosing(v label) bool { return v.Mate != nil && *v.Mate < 0 }

func buildContextTargets(dir string) ([]contextTarget, []string, error) {
	var data struct {
		Roots []struct {
			ID, Split, Turn, Component string
			Alternatives               []alternative
		}
	}
	var inputs struct{ Cases []probeCase }
	var leaves []struct {
		ID, Move string
		Witness  struct {
			FEN  string
			Path []string
		}
	}
	for name, target := range map[string]any{"fit-data.json": &data, "unlabelled-roots.json": &inputs, "leaves.json": &leaves} {
		if err := load(filepath.Join(dir, name), target); err != nil {
			return nil, nil, err
		}
	}
	rootInputs := map[string]probeCase{}
	for _, c := range inputs.Cases {
		rootInputs[c.ID] = c
	}
	leafMap := map[string]struct {
		FEN  string
		Path []string
	}{}
	for _, l := range leaves {
		leafMap[l.ID+"/"+l.Move] = l.Witness
	}
	targets := []contextTarget{}
	eligible := []struct{ ID, Turn, Component, Hash string }{}
	for _, root := range data.Roots {
		input, ok := rootInputs[root.ID]
		if !ok {
			return nil, nil, fmt.Errorf("missing root history %s", root.ID)
		}
		g, err := probeGame(input)
		if err != nil {
			return nil, nil, err
		}
		learnable := 0
		for _, a := range root.Alternatives {
			id := root.ID + "/" + a.Move
			w, ok := leafMap[id]
			if !ok {
				return nil, nil, fmt.Errorf("missing witness %s", id)
			}
			p := g.Position()
			for _, step := range w.Path {
				found := false
				for _, move := range p.ValidMovesUnsafe() {
					if move.String() == step {
						p = p.UpdateSearch(&move)
						found = true
						break
					}
				}
				if !found {
					return nil, nil, fmt.Errorf("illegal witness %s/%s", id, step)
				}
			}
			if p.String() != w.FEN {
				return nil, nil, fmt.Errorf("witness FEN differs %s", id)
			}
			position := input.Position
			if position == "" {
				position = "position fen " + input.FEN
			}
			if !strings.HasPrefix(position, "position ") {
				position = "position " + position
			}
			if len(w.Path) > 0 {
				words := strings.Fields(position)
				if !slices.Contains(words, "moves") {
					position += " moves"
				}
				position += " " + strings.Join(w.Path, " ")
			}
			target := contextTarget{ID: id, RootID: root.ID, Move: a.Move, Split: root.Split, Turn: root.Turn, Component: root.Component, FEN: w.FEN, Position: position, Key: hash([]byte(position)).SHA, Base: a.Base, Learnable: a.Learnable, Original: a.Teacher, position: p}
			if !a.Learnable {
				target.Direct = fixedContextLabel(a.Base)
			} else {
				learnable++
			}
			targets = append(targets, target)
		}
		if root.Split == "train" && learnable >= 2 {
			eligible = append(eligible, struct{ ID, Turn, Component, Hash string }{root.ID, root.Turn, root.Component, hash([]byte("v023-direct-leaf-audit/" + root.ID)).SHA})
		}
	}
	slices.SortFunc(eligible, func(a, b struct{ ID, Turn, Component, Hash string }) int { return strings.Compare(a.Hash, b.Hash) })
	selected := []string{}
	turns := map[string]int{}
	components := map[string]bool{}
	for _, r := range eligible {
		if turns[r.Turn] >= 8 || components[r.Component] {
			continue
		}
		turns[r.Turn]++
		components[r.Component] = true
		selected = append(selected, r.ID)
	}
	if len(selected) != 16 || turns["w"] != 8 || turns["b"] != 8 {
		return nil, nil, fmt.Errorf("insufficient balanced distinct-component diagnostic roots")
	}
	return targets, selected, nil
}

func contextAlignment(targets []contextTarget, selected []string) map[string]any {
	diffs := []float64{}
	changedMate, clear, opposite := 0, 0, 0
	roots := map[string][]contextTarget{}
	for _, t := range targets {
		if !slices.Contains(selected, t.RootID) {
			continue
		}
		roots[t.RootID] = append(roots[t.RootID], t)
		if !t.Learnable {
			continue
		}
		if t.Original.CP != nil && t.Direct.CP != nil {
			diffs = append(diffs, math.Abs(float64(*t.Original.CP-*t.Direct.CP)))
		}
		if contextLabelLosing(t.Original) != contextLabelLosing(t.Direct) {
			changedMate++
		}
	}
	for _, rows := range roots {
		for i, a := range rows {
			for _, b := range rows[i+1:] {
				if !a.Learnable || !b.Learnable || a.Original.CP == nil || b.Original.CP == nil || a.Direct.CP == nil || b.Direct.CP == nil {
					continue
				}
				old, new := *a.Original.CP-*b.Original.CP, *a.Direct.CP-*b.Direct.CP
				if absContext(old) > 25 && absContext(new) > 25 {
					clear++
					if old*new < 0 {
						opposite++
					}
				}
			}
		}
	}
	slices.Sort(diffs)
	median := 0.
	if len(diffs) > 0 {
		median = diffs[len(diffs)/2]
		if len(diffs)%2 == 0 {
			median = (median + diffs[len(diffs)/2-1]) / 2
		}
	}
	fraction := 0.
	if clear > 0 {
		fraction = float64(opposite) / float64(clear)
	}
	passed := len(diffs) >= 32 && (median >= 75 || fraction >= .1 || changedMate > 0)
	return map[string]any{"passed": passed, "diagnostic_roots": selected, "finite_comparisons": len(diffs), "median_absolute_target_difference_cp": median, "clear_pairs": clear, "opposite_preferences": opposite, "opposite_fraction": fraction, "changed_losing_mate_classifications": changedMate, "interpretation": "Root restricted-search values and actual weak-witness values differ; finite-node horizon effects also contribute. This is not a playing-strength result."}
}

func collectContextTargets(dir, teacher, out, stage string) error {
	if stage != "diagnose" && stage != "complete" {
		return fmt.Errorf("stage must be diagnose or complete")
	}
	targets, selected, err := buildContextTargets(dir)
	if err != nil {
		return err
	}
	teacherHash, err := probeDigest(teacher)
	if err != nil {
		return err
	}
	sources, err := workflowIdentities([]string{filepath.Join(dir, "fit-data.json"), filepath.Join(dir, "leaves.json"), filepath.Join(dir, "unlabelled-roots.json"), filepath.Join(dir, "boundary-audit-after-leaves.json")})
	if err != nil {
		return err
	}
	method := contextLabelMethod{Teacher: teacherHash, Sources: sources, Diagnostic: selected, Nodes: 50000, Workers: 3}
	journal := filepath.Join(out, "diagnostic-labels.jsonl")
	cache := map[string]probeQuery{}
	if stage == "diagnose" {
		if err = os.Mkdir(out, 0755); err != nil {
			return err
		}
		if err = writeNewJSON(filepath.Join(out, "method.json"), method); err != nil {
			return err
		}
	} else {
		var saved contextLabelMethod
		if err = load(filepath.Join(out, "method.json"), &saved); err != nil {
			return err
		}
		if !reflect.DeepEqual(saved, method) {
			return fmt.Errorf("label source/teacher/protocol changed")
		}
		var alignment struct{ Passed bool }
		if err = load(filepath.Join(out, "alignment.json"), &alignment); err != nil {
			return err
		}
		if !alignment.Passed {
			return fmt.Errorf("failed alignment gate; full collection not allowed")
		}
		f, err := os.Open(journal)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 4096), 8*1024*1024)
		for scanner.Scan() {
			var row contextQuery
			if err = json.Unmarshal(scanner.Bytes(), &row); err != nil {
				f.Close()
				return err
			}
			if row.Query.Error != "" {
				f.Close()
				return fmt.Errorf("saved failed query cannot be implicitly retried")
			}
			if _, ok := cache[row.Key]; ok {
				f.Close()
				return fmt.Errorf("duplicate saved query")
			}
			cache[row.Key] = row.Query
		}
		err = scanner.Err()
		f.Close()
		if err != nil {
			return err
		}
		journal = filepath.Join(out, "remaining-labels.jsonl")
	}
	flags := os.O_CREATE | os.O_EXCL | os.O_WRONLY
	f, err := os.OpenFile(journal, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	jobs := make(chan contextTarget, len(targets))
	needed := map[string]bool{}
	for _, t := range targets {
		if !t.Learnable || stage == "diagnose" && !slices.Contains(selected, t.RootID) {
			continue
		}
		if _, ok := cache[t.Key]; ok || needed[t.Key] {
			continue
		}
		needed[t.Key] = true
		jobs <- t
	}
	close(jobs)
	results := make(chan contextQuery, 3)
	var wg sync.WaitGroup
	for worker := 0; worker < 3; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e, err := startEngine(teacher, map[string]string{"Threads": "1", "Hash": "64", "MultiPV": "1", "Skill Level": "20", "UCI_LimitStrength": "false"})
			if err != nil {
				results <- contextQuery{Query: probeQuery{Error: err.Error()}}
				return
			}
			defer e.close()
			for target := range jobs {
				q := probeQuery{Position: target.Position, Command: "go nodes 50000"}
				if err = e.send("ucinewgame"); err == nil {
					err = e.send("isready")
				}
				if err == nil {
					_, err = e.wait("readyok", 5*time.Second)
				}
				if err != nil {
					q.Error = err.Error()
				} else {
					q = e.probe(target.Position, "go nodes 50000", 30*time.Second)
				}
				if q.Error == "" {
					legal := false
					for _, m := range target.position.ValidMovesUnsafe() {
						if m.String() == q.Move {
							legal = true
							break
						}
					}
					if !legal || q.Score.CP == nil && q.Score.Mate == nil {
						q.Error = "missing legal move/exact score"
					}
				}
				results <- contextQuery{Key: target.Key, Query: q}
				if q.Error != "" {
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(results) }()
	failed := ""
	var writeErr error
	done := 0
	for row := range results {
		if writeErr == nil {
			writeErr = encoder.Encode(row)
		}
		if row.Query.Error != "" {
			failed = row.Query.Error
		} else {
			cache[row.Key] = row.Query
		}
		done++
		if done%100 == 0 {
			fmt.Printf("labels %s %d/%d\n", stage, done, len(needed))
		}
	}
	if writeErr != nil {
		return writeErr
	}
	if failed != "" {
		return fmt.Errorf("label collection failed: %s", failed)
	}
	if done != len(needed) {
		return fmt.Errorf("label collection incomplete: %d/%d", done, len(needed))
	}
	for i := range targets {
		t := &targets[i]
		if !t.Learnable {
			continue
		}
		q, ok := cache[t.Key]
		if !ok {
			if stage == "diagnose" {
				continue
			}
			return fmt.Errorf("missing completed label %s", t.ID)
		}
		if q.Position != t.Position || hash([]byte(q.Position)).SHA != t.Key {
			return fmt.Errorf("query identity differs %s", t.ID)
		}
		t.Direct = labelFromContextQuery(q, strings.Fields(t.FEN)[1], t.Turn)
	}
	if stage == "diagnose" {
		report := contextAlignment(targets, selected)
		report["queries"] = len(cache)
		if err = writeNewJSON(filepath.Join(out, "alignment.json"), report); err != nil {
			return err
		}
		return printJSON(report)
	}
	if err = writeNewJSON(filepath.Join(out, "targets.json"), targets); err != nil {
		return err
	}
	if err = writeNewJSON(filepath.Join(out, "complete.json"), map[string]any{"passed": true, "targets": len(targets), "queries": len(cache), "new_queries": done, "source_method": method}); err != nil {
		return err
	}
	return printJSON(map[string]any{"passed": true, "targets": len(targets), "queries": len(cache), "new_queries": done})
}
