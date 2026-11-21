package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Explicit authorized documentation changes are different from raw evidence.
// Original hashes and retired-source records are never modified here.
func mapDocuments(record string, files []string) error {
	if record == "" || len(files) == 0 {
		return fmt.Errorf("new transition record and explicit document files required")
	}
	path := "tools/research/migration.json"
	before, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err = json.Unmarshal(before, &raw); err != nil {
		return err
	}
	var originals, rewritten map[string]fileEntry
	if err = json.Unmarshal(raw["originals"], &originals); err != nil {
		return err
	}
	if err = json.Unmarshal(raw["rewritten"], &rewritten); err != nil {
		return err
	}
	if rewritten == nil {
		rewritten = map[string]fileEntry{}
	}
	changes := map[string]any{}
	for _, file := range files {
		if !safePath(file) || !strings.HasSuffix(file, ".md") {
			return fmt.Errorf("document mapping accepts relative Markdown paths only: %s", file)
		}
		if _, ok := originals[file]; !ok {
			return fmt.Errorf("not a historical mapped document: %s", file)
		}
		info, err := os.Lstat(file)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("not a regular document: %s", file)
		}
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		next := hash(b)
		if rewritten[file] != next {
			change := map[string]any{"before": rewritten[file], "after": next}
			if head, e := exec.Command("git", "show", "HEAD:"+file).Output(); e == nil {
				change["head_identity"] = hash(head)
				change["matches_head"] = hash(head) == next
			}
			changes[file] = change
			rewritten[file] = next
		}
	}
	raw["rewritten"], err = json.Marshal(rewritten)
	if err != nil {
		return err
	}
	next, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	next = append(next, '\n')
	if err = writeNewJSON(record, map[string]any{"migration_before": hash(before), "migration_after": hash(next), "documents": changes, "originals_unchanged": true, "raw_evidence_unchanged": true}); err != nil {
		return err
	}
	tmp := filepath.Join(filepath.Dir(path), "migration.json.tmp")
	if err = os.WriteFile(tmp, next, 0644); err != nil {
		return err
	}
	if err = os.Rename(tmp, path); err != nil {
		return err
	}
	return printJSON(map[string]any{"updated_documents": len(changes), "transition_record": record})
}
