package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"time"
)

type depthReviewCase struct {
	Case   probeCase `json:"case"`
	Depths []int     `json:"observed_depths"`
}

// Include improvements and neutral changes as well as regressions. Identity is
// the complete position command/FEN, since historical cohorts can reuse IDs.
func depthReviewSelection(pattern string, flagsOnly ...bool) ([]depthReviewCase, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil || len(paths) == 0 {
		return nil, fmt.Errorf("native inputs required: %v", err)
	}
	selected := map[string]*depthReviewCase{}
	for _, path := range paths {
		flagged := map[string]bool{}
		if len(flagsOnly) > 0 && flagsOnly[0] {
			var quality struct {
				IDs []string `json:"new_candidate_losing_mate_cases"`
			}
			if err := load(filepath.Join(filepath.Dir(path), "quality.json"), &quality); err != nil {
				return nil, err
			}
			for _, id := range quality.IDs {
				flagged[id] = true
			}
		}
		var rows []probeNative
		if err := load(path, &rows); err != nil {
			return nil, err
		}
		for _, row := range rows {
			if len(flagsOnly) > 0 && flagsOnly[0] && !flagged[row.Case.ID] {
				continue
			}
			if row.Error != "" {
				return nil, fmt.Errorf("native row error: %s", row.Error)
			}
			moves := map[string]bool{}
			actors := map[string]bool{}
			for _, trial := range row.Trials {
				if trial.Query.Error != "" || trial.Query.Score.Depth < 1 {
					return nil, fmt.Errorf("incomplete native query")
				}
				moves[trial.Query.Move] = true
				actors[trial.Actor] = true
			}
			if !actors["baseline"] || !actors["candidate"] {
				return nil, fmt.Errorf("both actors required")
			}
			if len(moves) < 2 {
				continue
			}
			key := row.Case.Position + "\n" + row.Case.FEN
			c := selected[key]
			if c == nil {
				c = &depthReviewCase{Case: row.Case}
				selected[key] = c
			}
			for _, trial := range row.Trials {
				d := trial.Query.Score.Depth
				if !slices.Contains(c.Depths, d) {
					c.Depths = append(c.Depths, d)
				}
			}
		}
	}
	keys := make([]string, 0, len(selected))
	for key := range selected {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	result := make([]depthReviewCase, 0, len(keys))
	for _, key := range keys {
		c := selected[key]
		slices.Sort(c.Depths)
		result = append(result, *c)
	}
	return result, nil
}

func reviewDepths(base, candidate, pattern, out string, flagsOnly ...bool) error {
	cases, err := depthReviewSelection(pattern, flagsOnly...)
	if err != nil {
		return err
	}
	if len(cases) == 0 {
		return fmt.Errorf("no changed choices")
	}
	if err = os.Mkdir(out, 0755); err != nil {
		return err
	}
	if err = writeNewJSON(filepath.Join(out, "selection.json"), cases); err != nil {
		return err
	}
	journal, err := os.OpenFile(filepath.Join(out, "queries.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer journal.Close()
	engines := map[string]*engine{}
	defer func() {
		for _, e := range engines {
			e.close()
		}
	}()
	for actor, path := range map[string]string{"baseline": base, "candidate": candidate} {
		e, err := startEngine(path, nil)
		if err != nil {
			return err
		}
		engines[actor] = e
	}
	query := func(actor string, c probeCase, depth int, restriction string) (probeQuery, error) {
		e := engines[actor]
		if err := e.send("ucinewgame"); err != nil {
			return probeQuery{}, err
		}
		if err := e.send("isready"); err != nil {
			return probeQuery{}, err
		}
		if _, err := e.wait("readyok", 5*time.Second); err != nil {
			return probeQuery{}, err
		}
		command := fmt.Sprintf("go depth %d", depth)
		if restriction != "" {
			command += " searchmoves " + restriction
		}
		q := e.probe(c.Position, command, 120*time.Second)
		if err := json.NewEncoder(journal).Encode(map[string]any{"id": c.ID, "actor": actor, "query": q}); err != nil {
			return q, err
		}
		g, err := probeGame(c)
		if err != nil {
			return q, err
		}
		if q.Error != "" || q.Score.Depth != depth || !validProbeMove(g, q.Move) || (q.Score.CP == nil && q.Score.Mate == nil) {
			return q, fmt.Errorf("incomplete/illegal depth query: %s", q.Error)
		}
		return q, nil
	}
	rows := []map[string]any{}
	passed := true
	for _, c := range cases {
		for _, depth := range c.Depths {
			b, err := query("baseline", c.Case, depth, "")
			if err != nil {
				return err
			}
			a, err := query("candidate", c.Case, depth, "")
			if err != nil {
				return err
			}
			equal := reflect.DeepEqual(a.Score.CP, b.Score.CP) && reflect.DeepEqual(a.Score.Mate, b.Score.Mate)
			if equal && a.Move != b.Move {
				proof, err := query("baseline", c.Case, depth, a.Move)
				if err != nil {
					return err
				}
				equal = reflect.DeepEqual(proof.Score.CP, b.Score.CP) && reflect.DeepEqual(proof.Score.Mate, b.Score.Mate)
			}
			passed = passed && equal
			rows = append(rows, map[string]any{"case": c.Case, "depth": depth, "baseline": b, "candidate": a, "value_and_selected_move_valid": equal})
		}
	}
	if err := writeNewJSON(filepath.Join(out, "comparisons.json"), rows); err != nil {
		return err
	}
	return writeNewJSON(filepath.Join(out, "summary.json"), map[string]any{"passed": passed, "changed_cases": len(cases), "matched_depth_comparisons": len(rows), "timing_use": "diagnosis only; unrestricted roots except explicitly recorded serial co-best proof"})
}
