package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("PANDA_RESEARCH_TEST_ENGINE") == "scoreless" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			switch {
			case line == "uci":
				fmt.Println("id name Research control")
				for _, name := range []string{"Threads", "Hash", "Skill Level", "UCI_LimitStrength", "MultiPV"} {
					fmt.Printf("option name %s type string\n", name)
				}
				fmt.Println("uciok")
			case line == "isready":
				fmt.Println("readyok")
			case strings.HasPrefix(line, "go "):
				fmt.Println("info string scoreless-reference-control\nbestmove e2e4")
			}
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestReferenceFailureIsJournalled(t *testing.T) {
	t.Setenv("PANDA_RESEARCH_TEST_ENGINE", "scoreless")
	program, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	native := []probeNative{{Case: probeCase{ID: "start", Position: "position startpos"}}}
	if err := probeReferences(program, out, 100, native, "", "control", false); err == nil {
		t.Fatal("scoreless reference accepted")
	}
	data, err := os.ReadFile(filepath.Join(out, "references.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var row referenceRecord
	if err := json.Unmarshal(bytes.TrimSpace(data), &row); err != nil {
		t.Fatal(err)
	}
	if row.Query == nil || !strings.Contains(strings.Join(row.Query.Lines, "\n"), "scoreless-reference-control") {
		t.Fatal("failed reference output discarded")
	}
}

func pointer(v int) *int { return &v }
func TestSharpDecisionControls(t *testing.T) {
	finite := alternative{Move: "safe", Base: 10, Teacher: label{CP: pointer(-1000)}}
	bad := alternative{Move: "bad", Base: 10, Teacher: label{Mate: pointer(-4)}}
	good := alternative{Move: "win", Base: 20, Teacher: label{Mate: pointer(3)}}
	for _, tc := range []struct {
		name    string
		choices []alternative
		penalty float64
		risk    bool
	}{
		{"unsafe tie", []alternative{finite, bad}, 500, true},
		{"winning mate acceptable", []alternative{finite, good}, 0, false},
		{"all lost", []alternative{bad, bad}, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m, err := scoreRoots([]rootCase{{ID: "test", Turn: "b", Alternatives: tc.choices}}, make([]int, 24))
			if err != nil {
				t.Fatal(err)
			}
			if m.Mean != tc.penalty || m.Decisions[0].Risk != tc.risk {
				t.Fatalf("%+v", m)
			}
		})
	}
	a := alternative{Learnable: true, Values: make([]int, 12), Phase: 1}
	a.Values[0] = -1
	w := make([]int, 24)
	w[0] = 1
	got, err := adjusted(a, w)
	if err != nil || got != 0 {
		t.Fatalf("truncation %d %v", got, err)
	}
	a.Learnable = false
	a.Base = -999999
	got, err = adjusted(a, w)
	if err != nil || got != a.Base {
		t.Fatal("terminal changed")
	}
}
func TestArchivePathsAndNoOverwrite(t *testing.T) {
	for _, p := range []string{"../bad", "/bad", "a/../bad", "a\\bad", "C:bad"} {
		if safePath(p) {
			t.Fatal(p)
		}
	}
	var b bytes.Buffer
	tw := tar.NewWriter(&b)
	_ = tw.WriteHeader(&tar.Header{Name: "../escaped.go", Mode: 0644, Size: 1})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	if err := extract(tar.NewReader(&b), t.TempDir()); err == nil {
		t.Fatal("traversal accepted")
	}
	if err := replay("", "v0.18", t.TempDir(), "", false); err == nil {
		t.Fatal("existing replay directory accepted")
	}
	var source bytes.Buffer
	writer := tar.NewWriter(&source)
	if err := writer.WriteHeader(&tar.Header{Name: "pax_global_header", Typeflag: tar.TypeXGlobalHeader, PAXRecords: map[string]string{"comment": "commit metadata"}}); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteHeader(&tar.Header{Name: "main.go", Mode: 0644, Size: 12}); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("package main")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := extract(tar.NewReader(&source), t.TempDir()); err != nil {
		t.Fatal(err)
	}
}
func TestRetirementCannotHideEvidenceChange(t *testing.T) {
	root := t.TempDir()
	old := hash([]byte("original"))
	m := migration{Originals: map[string]fileEntry{"report.md": old}}
	if _, err := migratedEntry(root, "report.md", old, m); err == nil {
		t.Fatal("missing data accepted")
	}
	if err := os.WriteFile(filepath.Join(root, "report.md"), []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := migratedEntry(root, "report.md", old, m); err == nil {
		t.Fatal("unrecorded change accepted")
	}
	m.Rewritten = map[string]fileEntry{"report.md": hash([]byte("changed"))}
	if _, err := migratedEntry(root, "report.md", old, m); err != nil {
		t.Fatal(err)
	}
	m.Originals["record.json"] = old
	m.Rewritten["record.json"] = hash([]byte("changed"))
	if err := os.WriteFile(filepath.Join(root, "record.json"), []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := migratedEntry(root, "record.json", old, m); err == nil {
		t.Fatal("raw data accepted as a document rewrite")
	}
}
func TestQuietMatePositiveAndEscapeNegative(t *testing.T) {
	for _, tc := range []struct {
		fen      string
		positive bool
	}{
		{"7k/8/5KQ1/8/8/8/8/8 w - - 0 1", true},
		{"7k/8/6Q1/8/8/8/8/K7 w - - 0 1", false},
		{"7k/5Q2/6K1/8/8/8/8/8 b - - 0 1", false},
	} {
		moves, err := legalQuietMates(tc.fen)
		if err != nil {
			t.Fatal(err)
		}
		if (len(moves) > 0) != tc.positive {
			t.Fatalf("%s %v", tc.fen, moves)
		}
	}
}

func TestPublicationWording(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notes.md")
	if err := os.WriteFile(path, []byte("Playing policy remains unchanged."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"wording", "--file", path}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("A Maintainer approved this."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"wording", "--file", path}); err == nil {
		t.Fatal("prohibited publication wording accepted")
	}
	for _, text := range []string{"I kept playing policy unchanged.", "We preserved all results.", "Our next recommendation is bounded."} {
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
		if err := run([]string{"wording", "--file", path}); err == nil {
			t.Fatal("first-person publication text accepted")
		}
	}
}

func TestPGNClaimedRepetitionAndResultMismatch(t *testing.T) {
	dir := t.TempDir()
	summary := `{"matches":[{"pgn_file":"game.pgn","games":[{"result":"1/2-1/2","plies":8,"termination":"ThreefoldRepetition"}]}]}`
	if err := os.WriteFile(filepath.Join(dir, "summary.json"), []byte(summary), 0644); err != nil {
		t.Fatal(err)
	}
	pgn := "[Result \"1/2-1/2\"]\n\n1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8 1/2-1/2\n"
	if err := os.WriteFile(filepath.Join(dir, "game.pgn"), []byte(pgn), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := checkGames(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "game.pgn"), []byte("[Result \"1-0\"]\n\n1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8 1-0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := checkGames(dir); err == nil {
		t.Fatal("changed PGN result accepted")
	}
}
