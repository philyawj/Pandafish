package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

func confirmedMateFlags(row probeNative, scores map[string]probeScore) map[string]bool {
	flags := map[string]bool{"baseline": false, "candidate": false}
	for _, trial := range row.Trials {
		s := scores[trial.Query.Move]
		if s.Mate != nil && *s.Mate < 0 {
			flags[trial.Actor] = true
		}
	}
	return flags
}

// Diagnose every flagged case and both actors' complete move union. No actor
// timing is repeated, no favorable alternative is selected and old gates stay.
func confirmChoiceMates(dir, teacher, out string, nodes int) error {
	var method struct {
		Identities map[string]string `json:"identities_sha256"`
		Nodes      int               `json:"reference_nodes"`
	}
	if err := load(filepath.Join(dir, "method.json"), &method); err != nil {
		return err
	}
	identity, err := probeDigest(teacher)
	if err != nil {
		return err
	}
	if identity != method.Identities["reference"] || nodes <= method.Nodes {
		return fmt.Errorf("same teacher and larger fixed node budget required")
	}
	var quality struct {
		IDs []string `json:"new_candidate_losing_mate_cases"`
	}
	if err = load(filepath.Join(dir, "quality.json"), &quality); err != nil {
		return err
	}
	if len(quality.IDs) == 0 {
		return fmt.Errorf("no saved new mate flags")
	}
	var native []probeNative
	if err = load(filepath.Join(dir, "native.json"), &native); err != nil {
		return err
	}
	if err = os.Mkdir(out, 0755); err != nil {
		return err
	}
	options := map[string]string{"Threads": "1", "Hash": "64", "Skill Level": "20", "UCI_LimitStrength": "false", "MultiPV": "1"}
	if err = writeNewJSON(filepath.Join(out, "method.json"), map[string]any{"source": dir, "teacher_sha256": identity, "original_reference_nodes": method.Nodes, "confirmation_nodes": nodes, "options": options, "selection": "Every saved new mate flag; unrestricted root plus union of both actors' moves; no timed actor retry."}); err != nil {
		return err
	}
	journal, err := os.OpenFile(filepath.Join(out, "queries.jsonl"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer journal.Close()
	e, err := startEngine(teacher, options)
	if err != nil {
		return err
	}
	defer e.close()
	reports := []map[string]any{}
	for _, id := range quality.IDs {
		index := slices.IndexFunc(native, func(r probeNative) bool { return r.Case.ID == id })
		if index < 0 {
			return fmt.Errorf("missing flagged case %s", id)
		}
		row := native[index]
		g, err := probeGame(row.Case)
		if err != nil {
			return err
		}
		moves := []string{}
		for _, trial := range row.Trials {
			if !slices.Contains(moves, trial.Query.Move) {
				moves = append(moves, trial.Query.Move)
			}
		}
		slices.Sort(moves)
		scores := map[string]probeScore{}
		for _, move := range append([]string{""}, moves...) {
			if err = e.send("ucinewgame"); err != nil {
				return err
			}
			if err = e.send("isready"); err != nil {
				return err
			}
			if _, err = e.wait("readyok", 5*time.Second); err != nil {
				return err
			}
			command := fmt.Sprintf("go nodes %d", nodes)
			if move != "" {
				if !validProbeMove(g, move) {
					return fmt.Errorf("invalid saved move")
				}
				command += " searchmoves " + move
			}
			q := e.probe(row.Case.Position, command, 120*time.Second)
			if err = json.NewEncoder(journal).Encode(map[string]any{"id": id, "restricted_move": move, "query": q}); err != nil {
				return err
			}
			if q.Error != "" || !validProbeMove(g, q.Move) || (q.Score.CP == nil && q.Score.Mate == nil) {
				return fmt.Errorf("incomplete confirmation %s %s: %s", id, move, q.Error)
			}
			scores[move] = q.Score
		}
		flags := confirmedMateFlags(row, scores)
		reports = append(reports, map[string]any{"id": id, "case": row.Case, "reference_scores_by_move": scores, "losing_mate_reported": flags, "new_risk_still_reported": flags["candidate"] && !flags["baseline"]})
	}
	return writeNewJSON(filepath.Join(out, "summary.json"), map[string]any{"cases": reports, "completed": len(reports), "original_gate_rewritten": false, "interpretation": "Finite-node reference estimates, not a proof that a finite alternative avoids loss."})
}
