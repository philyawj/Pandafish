package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

func TestContextTargetOrientationAndFixedOutcomes(t *testing.T) {
	q := probeQuery{Score: probeScore{CP: contextIntValue(123)}}
	for _, c := range []struct {
		leaf, root string
		want       int
	}{{"w", "w", 123}, {"b", "b", 123}, {"w", "b", -123}, {"b", "w", -123}} {
		if got := labelFromContextQuery(q, c.leaf, c.root); got.CP == nil || *got.CP != c.want {
			t.Fatal(c, got)
		}
	}
	q.Score = probeScore{Mate: contextIntValue(-3)}
	if got := labelFromContextQuery(q, "b", "w"); got.Mate == nil || *got.Mate != 3 {
		t.Fatal("mate orientation")
	}
	for _, base := range []int{-999998, 0, 999995} {
		got := fixedContextLabel(base)
		if base == 0 && (got.CP == nil || *got.CP != 0) || base < 0 && (got.Mate == nil || *got.Mate >= 0) || base > 0 && (got.Mate == nil || *got.Mate <= 0) {
			t.Fatal(base, got)
		}
	}
}
func TestContextAlignmentPositiveAndNegativeControls(t *testing.T) {
	for _, c := range []struct {
		name          string
		difference    int
		reverse, want bool
	}{{"unchanged", 0, false, false}, {"target gap", 100, false, true}, {"opposite close preferences", 0, true, true}} {
		t.Run(c.name, func(t *testing.T) {
			rows := []contextTarget{}
			selected := []string{}
			for root := 0; root < 16; root++ {
				id := fmt.Sprint(root)
				selected = append(selected, id)
				for i := 0; i < 2; i++ {
					old, new := 0, c.difference
					if c.reverse {
						old = (2*i - 1) * 30
						new = -old
					}
					rows = append(rows, contextTarget{RootID: id, Learnable: true, Original: label{CP: contextIntValue(old)}, Direct: label{CP: contextIntValue(new)}})
				}
			}
			if got := contextAlignment(rows, selected); got["passed"] != c.want {
				t.Fatal(got)
			}
		})
	}
}
func TestContextTargetHistoriesAndSelection(t *testing.T) {
	rows, selected, err := buildContextTargets(filepath.Join("..", "benchmark", "research", "v0.19-sharp-branches"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4474 || len(selected) != 16 {
		t.Fatal("coverage", len(rows), len(selected))
	}
	turns := map[string]string{}
	components := map[string]string{}
	for _, r := range rows {
		if slices.Contains(selected, r.RootID) {
			if r.Split != "train" {
				t.Fatal("nontraining diagnostic")
			}
			turns[r.RootID] = r.Turn
			if old, ok := components[r.Component]; ok && old != r.RootID {
				t.Fatal("component duplicated")
			}
			components[r.Component] = r.RootID
		}
	}
	white, black := 0, 0
	for _, turn := range turns {
		if turn == "w" {
			white++
		} else {
			black++
		}
	}
	if white != 8 || black != 8 {
		t.Fatal("turn balance")
	}
	for _, i := range []int{0, 100, len(rows) - 1} {
		g, err := probeGame(probeCase{Position: rows[i].Position})
		if err != nil || g.Position().String() != rows[i].FEN {
			t.Fatal("history serialization", i, err)
		}
	}
}
