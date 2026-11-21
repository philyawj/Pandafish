package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/corentings/chess/v2"
)

// Verify the reused board/base link independently of the fitter's new features.
// Old split names remain data provenance, not an unseen validation claim.
func TestContextCorpusOriginalBase(t *testing.T) {
	dir := "tools/benchmark/research/v0.19-sharp-branches/"
	var roots struct {
		Roots []struct {
			ID, Turn     string
			Alternatives []struct {
				Move        string
				Base, Phase int
				Learnable   bool
			}
		}
	}
	var leaves []struct {
		ID, Move string
		Witness  struct {
			FEN, Reason string
			BaseRoot    int
		}
	}
	for name, target := range map[string]any{"fit-data.json": &roots, "leaves.json": &leaves} {
		b, err := os.ReadFile(dir + name)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(b, target); err != nil {
			t.Fatal(err)
		}
	}
	positions := map[string]struct {
		FEN, Reason string
		BaseRoot    int
	}{}
	for _, row := range leaves {
		positions[row.ID+"/"+row.Move] = row.Witness
	}
	checked := 0
	for _, root := range roots.Roots {
		color := chess.White
		if root.Turn == "b" {
			color = chess.Black
		}
		for _, a := range root.Alternatives {
			if !a.Learnable {
				continue
			}
			w, ok := positions[root.ID+"/"+a.Move]
			if !ok || w.Reason != "static" || w.BaseRoot != a.Base {
				t.Fatal("witness/base link", root.ID, a.Move)
			}
			p := evalFEN(t, w.FEN)
			if got := referenceBaseStaticScore(p, color); got != a.Base {
				t.Fatal("original base differs", root.ID, a.Move, got, a.Base)
			}
			checked++
		}
	}
	if checked < 4000 {
		t.Fatal("unexpected corpus coverage", checked)
	}
	t.Log("matched original static witnesses", checked)
}
