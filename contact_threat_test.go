package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/corentings/chess/v2"
)

func contactReferenceMoves(p *chess.Position) []chess.Move {
	var king chess.Square
	count := 0
	pieces := p.Board().SquareMap()
	for sq, pc := range pieces {
		if pc.Type() == chess.King && pc.Color() != p.Turn() {
			king = sq
			count++
		}
	}
	if count != 1 {
		return nil
	}
	var out []chess.Move
	for _, m := range p.ValidMoves() {
		if p.Board().Piece(m.S1()).Type() != chess.Queen || m.HasTag(chess.Capture) || m.Promo() != chess.NoPieceType || !m.HasTag(chess.Check) || abs(int(m.S2().File())-int(king.File())) > 1 || abs(int(m.S2().Rank())-int(king.Rank())) > 1 {
			continue
		}
		protected := false
		for sq, pc := range pieces {
			if pc.Color() == p.Turn() && (pc.Type() == chess.Pawn || pc.Type() == chess.Knight || pc.Type() == chess.King) && slices.Contains(p.Board().AttacksFrom(sq), m.S2()) {
				protected = true
			}
		}
		if protected {
			out = append(out, m)
		}
	}
	return out
}

func TestContactThreatControls(t *testing.T) {
	for _, c := range []struct {
		fen  string
		mate bool
	}{
		{"7k/8/5KQ1/8/8/8/8/8 w - - 0 1", true},
		{"8/8/8/8/8/5kq1/8/7K b - - 0 1", true},
		{"7k/8/6Q1/8/8/8/8/K7 w - - 0 1", false},
		{"7k/5Q2/6K1/8/8/8/8/8 b - - 0 1", false},
		{"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1", false},
		{"7k/P7/8/8/8/8/8/7K w - - 0 1", false},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", false},
	} {
		p := evalFEN(t, c.fen)
		before, hash := p.String(), p.ZobristHash()
		cold := p.AppendSearchProtectedQueenContacts(nil)
		want := contactReferenceMoves(p)
		if !sameMoves(cold, want) || !sameMoves(p.AppendSearchProtectedQueenContacts(nil), want) {
			t.Fatal("subset mismatch", c.fen)
		}
		var child chess.SearchPosition
		m, found := contactMateProof(p, &child)
		if found != c.mate || found && (!m.HasTag(chess.Check) || p.Update(&m).Status() != chess.Checkmate) {
			t.Fatal("proof mismatch", c.fen, found, m)
		}
		if p.String() != before || p.ZobristHash() != hash {
			t.Fatal("parent mutated")
		}
	}
}

func TestContactThreatTrainingDiagnosis(t *testing.T) {
	output := os.Getenv("CONTACT_DIAG_OUTPUT")
	if output == "" {
		t.Skip("explicit fixed diagnosis required")
	}
	path := "tools/benchmark/research/v0.19-sharp-branches/training-quiet-threat-scan.json"
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var input struct {
		Rows []struct {
			ID, Choice string
			Witness    struct{ FEN string }
			Opponent   []json.RawMessage `json:"opponent_mates_after_hypothetical_pass"`
		}
	}
	if err = json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}
	type hit struct {
		Row                         int
		ID, Choice, Move, FENSHA256 string
	}
	hits := []hit{}
	distinct := map[string]bool{}
	oldPositive, checks := 0, 0
	for i, r := range input.Rows {
		p := evalFEN(t, r.Witness.FEN)
		before := p.String()
		for _, q := range []*chess.Position{p, p.Update(nil)} {
			cold := q.AppendSearchProtectedQueenContacts(nil)
			want := contactReferenceMoves(q)
			if got := q.AppendSearchProtectedQueenContacts(nil); !sameMoves(got, want) || !sameMoves(cold, want) {
				t.Fatal("training subset", i)
			}
			var child chess.SearchPosition
			if m, ok := contactMateProof(q, &child); ok && (!slices.ContainsFunc(want, func(n chess.Move) bool { return n.String() == m.String() }) || q.Update(&m).Status() != chess.Checkmate) {
				t.Fatal("false mating certificate", i)
			}
			checks++
		}
		if len(r.Opponent) > 0 {
			oldPositive++
		}
		var pass, child chess.SearchPosition
		if m, ok := passedContactThreat(p, &pass, &child); ok {
			q := p.Update(nil)
			if len(r.Opponent) == 0 || q.Update(&m).Status() != chess.Checkmate {
				t.Fatal("false opponent certificate", i)
			}
			distinct[r.Witness.FEN] = true
			hits = append(hits, hit{i, r.ID, r.Choice, m.String(), fmt.Sprintf("%x", sha256.Sum256([]byte(r.Witness.FEN)))})
		}
		if p.String() != before {
			t.Fatal("mutated training root")
		}
	}
	passed := len(distinct) >= 8
	result := map[string]any{"passed": passed, "source": path, "source_sha256": fmt.Sprintf("%x", sha256.Sum256(b)), "survey_rows": len(input.Rows), "old_opponent_mate_rows": oldPositive, "full_subset_checks": checks, "distinct_covered_threats": len(distinct), "covered_rows": hits, "false_certificates": 0}
	data, _ := json.MarshalIndent(result, "", "  ")
	if err = os.WriteFile(output, append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if !passed {
		t.Fatal("frozen coverage gate fails", len(distinct))
	}
	t.Logf("%d distinct threats, %d subset checks, zero false certificates", len(distinct), checks)
}

func TestContactGeneratorOracle(t *testing.T) {
	path := os.Getenv("ORACLE_INPUT")
	if path == "" {
		t.Skip("explicit independent vectors required")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var input struct{ Cases []oracleCase }
	if err = json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}
	if len(input.Cases) != 596 {
		t.Fatal("independent coverage")
	}
	count := 0
	for _, c := range input.Cases {
		p := evalFEN(t, c.FEN)
		cold := p.AppendSearchProtectedQueenContacts(nil)
		if !sameMoves(cold, contactReferenceMoves(p)) {
			t.Fatal("independent subset geometry", c.ID)
		}
		for _, m := range cold {
			i := slices.IndexFunc(c.Moves, func(n oracleMove) bool { return n.UCI == m.String() })
			if i < 0 || !c.Moves[i].Check || c.Moves[i].Capture || p.Update(&m).String() != c.Moves[i].ChildFEN {
				t.Fatal("independent contact move/state", c.ID, m)
			}
			count++
		}
	}
	if out := os.Getenv("CONTACT_ORACLE_OUTPUT"); out != "" {
		b, _ = json.MarshalIndent(map[string]any{"passed": true, "positions": 596, "contact_moves": count}, "", "  ")
		if err = os.WriteFile(out, append(b, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
}
