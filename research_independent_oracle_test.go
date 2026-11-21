package main

import (
	"encoding/json"
	"github.com/corentings/chess/v2"
	"os"
	"testing"
)

type oracleMove struct {
	UCI                                  string
	Tactical, Capture, EP, Castle, Check bool
	ChildFEN                             string `json:"child_fen"`
	ChildHash                            uint64 `json:"child_hash"`
}
type oracleCase struct {
	ID          int
	FEN, Status string
	Hash        uint64
	Moves       []oracleMove
	Perft2      uint64
}

func TestSearchMoveOracle(t *testing.T) {
	if os.Getenv("ORACLE_INPUT") == "" {
		t.Skip("explicit independent vectors required")
	}
	worker := searcher{}
	data, err := os.ReadFile(os.Getenv("ORACLE_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var input struct{ Cases []oracleCase }
	if err = json.Unmarshal(data, &input); err != nil {
		t.Fatal(err)
	}
	type failure struct {
		ID     int
		FEN    string
		Errors []string
	}
	failures := []failure{}
	for _, c := range input.Cases {
		errors := []string{}
		p := evalFEN(t, c.FEN)
		if p.ZobristHash() != c.Hash {
			errors = append(errors, "root hash")
		}
		status := "normal"
		if p.SearchStatus() == chess.Checkmate {
			status = "mate"
		}
		if p.SearchStatus() == chess.Stalemate {
			status = "stalemate"
		}
		if status != c.Status || p.HasLegalMove() != (len(c.Moves) > 0) {
			errors = append(errors, "status/existence")
		}
		got := p.AppendSearchMoves(nil, false)
		tactics := p.AppendSearchMoves(nil, true)
		moves := map[string]chess.Move{}
		tactical := map[string]bool{}
		for _, m := range got {
			moves[m.String()] = m
		}
		for _, m := range tactics {
			tactical[m.String()] = true
		}
		if len(got) != len(c.Moves) || len(moves) != len(got) {
			errors = append(errors, "legal set size/duplicates")
		}
		var count uint64
		for _, expected := range c.Moves {
			m, ok := moves[expected.UCI]
			if !ok {
				errors = append(errors, "missing "+expected.UCI)
				continue
			}
			if m.HasTag(chess.Capture) != expected.Capture || m.HasTag(chess.EnPassant) != expected.EP || m.HasTag(chess.Check) != expected.Check || (m.HasTag(chess.KingSideCastle) || m.HasTag(chess.QueenSideCastle)) != expected.Castle {
				errors = append(errors, "tags "+expected.UCI)
			}
			if tactical[expected.UCI] != expected.Tactical {
				errors = append(errors, "tactical eligibility "+expected.UCI)
			}
			child := worker.childPosition(p, &m, 1)
			if child.String() != expected.ChildFEN || child.ZobristHash() != expected.ChildHash {
				errors = append(errors, "child FEN/hash "+expected.UCI)
			}
			count += uint64(len(child.AppendSearchMoves(nil, false)))
		}
		wantedTactical := 0
		for _, m := range c.Moves {
			if m.Tactical {
				wantedTactical++
			}
		}
		if len(tactics) != wantedTactical {
			errors = append(errors, "tactical size")
		}
		if count != c.Perft2 {
			errors = append(errors, "depth2 perft")
		}
		if p.String() != c.FEN || p.ZobristHash() != c.Hash {
			errors = append(errors, "root mutation")
		}
		if len(errors) > 0 {
			failures = append(failures, failure{c.ID, c.FEN, errors})
		}
	}
	result := map[string]any{"positions": len(input.Cases), "passed": len(input.Cases) - len(failures), "failures": failures}
	output, _ := json.MarshalIndent(result, "", "  ")
	if err = os.WriteFile(os.Getenv("ORACLE_OUTPUT"), append(output, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if len(failures) > 0 {
		t.Fatalf("%d oracle cases failed; every position retained in output", len(failures))
	}
	t.Logf("Independent moves/tags, child FEN/hash, status and depth2 perft pass on %d positions", len(input.Cases))
}
