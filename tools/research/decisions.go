package main

import (
	"fmt"
	"math"
	"path/filepath"
	"reflect"

	"github.com/corentings/chess/v2"
)

type label struct {
	CP   *int `json:"cp"`
	Mate *int `json:"mate"`
}
type alternative struct {
	Move      string `json:"move"`
	Teacher   label  `json:"teacher"`
	Base      int    `json:"base"`
	Learnable bool   `json:"learnable"`
	Values    []int  `json:"values"`
	Phase     int    `json:"phase"`
}
type rootCase struct {
	ID           string        `json:"id"`
	Split        string        `json:"split"`
	Turn         string        `json:"turn"`
	Alternatives []alternative `json:"alternatives"`
}
type decision struct {
	ID        string   `json:"id"`
	Turn      string   `json:"turn"`
	Penalty   float64  `json:"penalty"`
	Risk      bool     `json:"risk"`
	Selected  []string `json:"selected"`
	Avoidable bool     `json:"avoidable_mate_root"`
}
type turnMetric struct {
	Roots int     `json:"roots"`
	Mean  float64 `json:"mean_penalty"`
}
type metrics struct {
	Roots     int                   `json:"decision_roots"`
	Mean      float64               `json:"mean_penalty"`
	Severe    int                   `json:"severe_over_300"`
	Risks     int                   `json:"avoidable_loss_roots"`
	Coverage  int                   `json:"mate_coverage_roots"`
	CPRoots   int                   `json:"cp_only_roots"`
	CPMean    float64               `json:"cp_only_mean_regret"`
	Turns     map[string]turnMetric `json:"by_turn"`
	Decisions []decision            `json:"decisions,omitempty"`
}

func losing(a alternative) bool { return a.Teacher.Mate != nil && *a.Teacher.Mate < 0 }
func adjusted(a alternative, w []int) (int, error) {
	if !a.Learnable {
		return a.Base, nil
	}
	if len(w) != 24 || len(a.Values) != 12 || a.Phase < 0 || a.Phase > 24 {
		return 0, fmt.Errorf("invalid feature/weight dimensions")
	}
	numerator := 0
	for i, v := range a.Values {
		numerator += v * (w[i]*a.Phase + w[i+12]*(24-a.Phase))
	}
	return a.Base + numerator/24, nil // Go signed integer division truncates toward zero.
}
func scoreRoots(roots []rootCase, w []int) (metrics, error) {
	m := metrics{Turns: map[string]turnMetric{}}
	for _, root := range roots {
		if len(root.Alternatives) < 2 {
			continue
		}
		best := -int(^uint(0) >> 1)
		chosen := []alternative{}
		bestFinite := -int(^uint(0) >> 1)
		finite := 0
		bad := 0
		for _, a := range root.Alternatives {
			s, err := adjusted(a, w)
			if err != nil {
				return m, err
			}
			if s > best {
				best = s
				chosen = nil
			}
			if s == best {
				chosen = append(chosen, a)
			}
			if a.Teacher.CP != nil {
				bestFinite = max(bestFinite, *a.Teacher.CP)
				finite++
			}
			if losing(a) {
				bad++
			}
		}
		row := decision{ID: root.ID, Turn: root.Turn, Avoidable: bad > 0 && bad < len(root.Alternatives)}
		for _, a := range chosen {
			row.Selected = append(row.Selected, a.Move)
			switch {
			case losing(a):
				if row.Avoidable {
					row.Penalty += 1000
					row.Risk = true
				}
			case a.Teacher.Mate != nil:
			case a.Teacher.CP != nil:
				row.Penalty += float64(min(1000, max(0, bestFinite-*a.Teacher.CP-25)))
			default:
				return m, fmt.Errorf("unlabelled branch %s", root.ID)
			}
		}
		row.Penalty /= float64(len(chosen))
		m.Decisions = append(m.Decisions, row)
		m.Roots++
		m.Mean += row.Penalty
		if row.Penalty > 300 {
			m.Severe++
		}
		if row.Risk {
			m.Risks++
		}
		if row.Avoidable {
			m.Coverage++
		}
		if finite == len(root.Alternatives) {
			m.CPRoots++
			m.CPMean += row.Penalty
		}
		t := m.Turns[root.Turn]
		t.Roots++
		t.Mean += row.Penalty
		m.Turns[root.Turn] = t
	}
	if m.Roots == 0 {
		return m, fmt.Errorf("no decision roots")
	}
	m.Mean /= float64(m.Roots)
	if m.CPRoots > 0 {
		m.CPMean /= float64(m.CPRoots)
	}
	for turn, t := range m.Turns {
		t.Mean /= float64(t.Roots)
		m.Turns[turn] = t
	}
	return m, nil
}
func decisionReport(dir string) (map[string]map[string]metrics, error) {
	var input struct {
		Roots []rootCase `json:"roots"`
	}
	var w struct {
		Middle []int `json:"middle"`
		Ending []int `json:"ending"`
	}
	var saved struct {
		Quality map[string]map[string]metrics `json:"quality"`
		Passed  bool                          `json:"surrogate_gate_passed"`
	}
	for name, target := range map[string]any{"fit-data.json": &input, "model-weights.json": &w, "fit-results.json": &saved} {
		if err := load(filepath.Join(dir, name), target); err != nil {
			return nil, err
		}
	}
	if len(w.Middle) != 12 || len(w.Ending) != 12 {
		return nil, fmt.Errorf("invalid frozen weights")
	}
	report := map[string]map[string]metrics{}
	for _, split := range []string{"train", "validation", "holdout"} {
		cases := []rootCase{}
		for _, root := range input.Roots {
			if root.Split == split {
				cases = append(cases, root)
			}
		}
		report[split] = map[string]metrics{}
		for actor, weights := range map[string][]int{"baseline": make([]int, 24), "candidate": append(w.Middle, w.Ending...)} {
			got, err := scoreRoots(cases, weights)
			if err != nil {
				return nil, err
			}
			want, ok := saved.Quality[split][actor]
			if !ok {
				return nil, fmt.Errorf("missing saved metrics")
			}
			if math.Abs(got.Mean-want.Mean) > 1e-10 || !reflect.DeepEqual(got.Decisions, want.Decisions) || got.Severe != want.Severe || got.Risks != want.Risks || got.Roots != want.Roots || got.Coverage != want.Coverage || got.CPRoots != want.CPRoots || math.Abs(got.CPMean-want.CPMean) > 1e-10 || !reflect.DeepEqual(got.Turns, want.Turns) {
				return nil, fmt.Errorf("saved decisions mismatch: %s %s", split, actor)
			}
			got.Decisions = nil
			report[split][actor] = got
		}
	}
	return report, nil
}

func legalQuietMates(fen string) (map[string]bool, error) {
	option, err := chess.FEN(fen)
	if err != nil {
		return nil, err
	}
	p := chess.NewGame(option).Position()
	found := map[string]bool{}
	for _, m := range p.ValidMoves() {
		if m.HasTag(chess.Capture) || m.HasTag(chess.EnPassant) || m.Promo() != chess.NoPieceType {
			continue
		}
		if p.Update(&m).SearchStatus() == chess.Checkmate {
			found[m.String()] = true
		}
	}
	return found, nil
}
func mateProofs(input string) (map[string]int, error) {
	var data struct {
		Rows []struct {
			ID      string `json:"id"`
			Witness struct{ FEN string }
			Mates   []struct{ Move string } `json:"side_to_move_mates_in_one"`
		} `json:"rows"`
	}
	if err := load(input, &data); err != nil {
		return nil, err
	}
	positives, negatives := 0, 0
	leaves := map[string]bool{}
	for _, row := range data.Rows {
		actual, err := legalQuietMates(row.Witness.FEN)
		if err != nil {
			return nil, err
		}
		expected := map[string]bool{}
		for _, move := range row.Mates {
			expected[move.Move] = true
		}
		if !reflect.DeepEqual(actual, expected) {
			return nil, fmt.Errorf("quiet mate proof mismatch: %s: got %v expected %v", row.ID, actual, expected)
		}
		if len(actual) > 0 {
			positives++
			leaves[row.Witness.FEN] = true
		} else {
			negatives++
		}
	}
	return map[string]int{"records_with_quiet_mate": positives, "distinct_positive_leaf_fens": len(leaves), "negative_records_checked": negatives}, nil
}
