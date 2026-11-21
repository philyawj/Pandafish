package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDocumentMapPreservesOriginalsAndRawEvidence(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("tools/research", 0755); err != nil {
		t.Fatal(err)
	}
	originals := map[string]fileEntry{"AGENTS.md": hash([]byte("old")), "raw.json": hash([]byte("raw"))}
	if err := writeNewJSON("tools/research/migration.json", map[string]any{"originals": originals, "rewritten": map[string]fileEntry{"AGENTS.md": originals["AGENTS.md"]}, "retired": []string{"old_source"}}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("AGENTS.md", []byte("new authorized guidance"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := mapDocuments("transition.json", []string{"AGENTS.md"}); err != nil {
		t.Fatal(err)
	}
	var got migration
	if err := load("tools/research/migration.json", &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Originals, originals) || got.Rewritten["AGENTS.md"] != hash([]byte("new authorized guidance")) || !reflect.DeepEqual(got.Retired, []string{"old_source"}) {
		t.Fatal("rewrote original evidence", got)
	}
	for _, file := range []string{"raw.json", "unknown.md", "../outside.md"} {
		if err := mapDocuments(filepath.Base(file)+"-record.json", []string{file}); err == nil {
			t.Fatal("accepted unauthorized mapping", file)
		}
	}
	if err := mapDocuments("transition.json", []string{"AGENTS.md"}); err == nil {
		t.Fatal("overwrote transition record")
	}
}

func TestAuditReportsAllOriginalMismatches(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("tools/research", 0755); err != nil {
		t.Fatal(err)
	}
	originals := map[string]fileEntry{}
	for _, path := range []string{"AGENTS.md", "README.md"} {
		originals[path] = hash([]byte("original"))
		if err := os.WriteFile(path, []byte("changed"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := writeNewJSON("tools/research/migration.json", map[string]any{"originals": originals, "rewritten": map[string]fileEntry{}}); err != nil {
		t.Fatal(err)
	}
	_, err := audit(".", "")
	if err == nil || !strings.Contains(err.Error(), "AGENTS.md") || !strings.Contains(err.Error(), "README.md") {
		t.Fatal("incomplete audit report", err)
	}
}
