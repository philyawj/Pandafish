package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Core overlays deliberately exclude unrelated tooling/docs. All nested module
// metadata is accepted by replay, so no second, repaired snapshot is necessary.
func coreSourcePatterns() []string {
	return []string{"*.go", "go.mod", "go.sum", "internal/telemetry/*.go", "third_party/chess/*.go", "third_party/chess/go.mod", "third_party/chess/go.sum"}
}
func createSnapshot(base, out string, extra ...string) error {
	if base == "" || out == "" {
		return fmt.Errorf("base tag and new output file required")
	}
	if strings.HasPrefix(base, "-") {
		return fmt.Errorf("invalid base tag")
	}
	commit, err := exec.Command("git", "rev-parse", "--verify", base+"^{commit}").Output()
	if err != nil {
		return err
	}
	identities, err := workflowIdentities(append(coreSourcePatterns(), extra...))
	if err != nil {
		return err
	}
	files := map[string]string{}
	for path := range identities {
		if !allowedOverlay(filepath.ToSlash(path)) {
			return fmt.Errorf("unsafe/non-Go source overlay: %s", path)
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		files[filepath.ToSlash(path)] = string(b)
	}
	b, err := json.MarshalIndent(map[string]any{"base_tag": base, "base_commit": strings.TrimSpace(string(commit)), "files": files}, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err = writeNewFile(out, b); err != nil {
		return err
	}
	return printJSON(map[string]any{"file": out, "files": len(files), "identity": hash(b)})
}
func createInventory(dir string) error {
	if dir == "" {
		return fmt.Errorf("evidence directory required")
	}
	entries := map[string]fileEntry{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "inventory.json" {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("not regular evidence: %s", path)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entries[filepath.ToSlash(rel)] = hash(b)
		return nil
	})
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("evidence directory is empty")
	}
	path := filepath.Join(dir, "inventory.json")
	if err = writeNewJSON(path, entries); err != nil {
		return err
	}
	return printJSON(map[string]any{"inventory": path, "files": len(entries)})
}
