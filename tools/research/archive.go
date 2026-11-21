package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type fileEntry struct {
	SHA   string `json:"sha256"`
	Bytes int64  `json:"bytes"`
}
type migration struct {
	Asset     string               `json:"asset"`
	SHA       string               `json:"asset_sha256"`
	Originals map[string]fileEntry `json:"originals"`
	Retired   []string             `json:"retired"`
	Rewritten map[string]fileEntry `json:"rewritten"`
}

func hash(b []byte) fileEntry { return fileEntry{fmt.Sprintf("%x", sha256.Sum256(b)), int64(len(b))} }
func safePath(name string) bool {
	return name != "" && name != "." && !filepath.IsAbs(name) && !strings.Contains(name, "\\") &&
		!strings.Contains(name, ":") && filepath.ToSlash(filepath.Clean(name)) == name && name != ".." && !strings.HasPrefix(name, "../")
}
func migratedEntry(root, name string, want fileEntry, m migration) (string, error) {
	if !safePath(name) {
		return "", fmt.Errorf("unsafe evidence path %q", name)
	}
	b, err := os.ReadFile(filepath.Join(root, name))
	if err == nil && hash(b) == want {
		return "unchanged", nil
	}
	if m.Originals[name] != want {
		return "", fmt.Errorf("evidence changed without original manifest: %s", name)
	}
	if err == nil && strings.HasSuffix(name, ".md") && m.Rewritten[name] == hash(b) {
		return "documented rewrite", nil
	}
	if os.IsNotExist(err) && regexp.MustCompile(`\.(py|sh|bash|zsh|js|ts)(\.|$)`).MatchString(filepath.Base(name)) {
		for _, retired := range m.Retired {
			if name == retired {
				return "retired to release asset", nil
			}
		}
	}
	return "", fmt.Errorf("evidence mismatch: %s", name)
}
func audit(root, legacy string) (map[string]any, error) {
	var m migration
	if err := load(filepath.Join(root, "tools/research/migration.json"), &m); err != nil {
		return nil, err
	}
	var originalErrors []error
	for name, want := range m.Originals {
		if _, err := migratedEntry(root, name, want, m); err != nil {
			originalErrors = append(originalErrors, err)
		}
	}
	if len(originalErrors) > 0 {
		return nil, errors.Join(originalErrors...)
	}
	counts := map[string]int{}
	inventories, entries := 0, 0
	err := filepath.WalkDir(filepath.Join(root, "tools/benchmark/research"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != "inventory.json" {
			return nil
		}
		var inv map[string]fileEntry
		if err := load(path, &inv); err != nil {
			return err
		}
		inventories++
		for name, want := range inv {
			if !safePath(name) {
				return fmt.Errorf("unsafe inventory path %q", name)
			}
			p := filepath.Join(filepath.Dir(path), name)
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}
			state, err := migratedEntry(root, filepath.ToSlash(rel), want, m)
			if err != nil {
				return err
			}
			counts[state]++
			entries++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Also cover retired scripts outside inventory-bearing archives and every rewrite.
	for _, name := range m.Retired {
		state, err := migratedEntry(root, name, m.Originals[name], m)
		if err != nil {
			return nil, err
		}
		if state != "retired to release asset" {
			return nil, fmt.Errorf("retired source remains: %s", name)
		}
	}
	for name, want := range m.Rewritten {
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			return nil, err
		}
		if hash(b) != want {
			return nil, fmt.Errorf("rewrite hash mismatch: %s", name)
		}
	}
	verified := 0
	if legacy != "" {
		var err error
		verified, err = verifyLegacy(legacy, m)
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{"inventories": inventories, "inventory_entries": entries, "states": counts,
		"retired_scripts": len(m.Retired), "explicit_rewrites": len(m.Rewritten), "release_asset": m.Asset,
		"release_asset_verified": legacy != "", "original_files_verified_in_asset": verified}, nil
}
func verifyLegacy(path string, m migration) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if hash(b).SHA != m.SHA {
		return 0, fmt.Errorf("release asset hash mismatch")
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return 0, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	seen := map[string]bool{}
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		if !safePath(h.Name) || seen[h.Name] {
			return 0, fmt.Errorf("invalid/duplicate archive path %q", h.Name)
		}
		seen[h.Name] = true
		if h.Typeflag != tar.TypeReg {
			return 0, fmt.Errorf("non-file evidence entry %s", h.Name)
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			return 0, err
		}
		if h.Name == "MANIFEST.json" {
			var original map[string]fileEntry
			if err := json.Unmarshal(body, &original); err != nil {
				return 0, err
			}
			if len(original) != len(m.Originals) {
				return 0, fmt.Errorf("manifest count mismatch")
			}
			for p, e := range original {
				if m.Originals[p] != e {
					return 0, fmt.Errorf("manifest mismatch %s", p)
				}
			}
			continue
		}
		if m.Originals[h.Name] != hash(body) {
			return 0, fmt.Errorf("original file mismatch %s", h.Name)
		}
	}
	if !seen["MANIFEST.json"] || len(seen) != len(m.Originals)+1 {
		return 0, fmt.Errorf("incomplete original archive")
	}
	return len(m.Originals), nil
}

func extract(tr *tar.Reader, out string) error {
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if h.Typeflag == tar.TypeXGlobalHeader {
			continue
		} // git archive's commit metadata, not a file.
		name := strings.TrimSuffix(h.Name, "/")
		if !safePath(name) {
			return fmt.Errorf("unsafe source path %q", h.Name)
		}
		path := filepath.Join(out, name)
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fs.FileMode(h.Mode)&0777)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(f, tr)
			closeErr := f.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		default:
			return fmt.Errorf("unsupported source entry %s", h.Name)
		}
	}
}
func allowedOverlay(p string) bool {
	return safePath(p) && (strings.HasSuffix(p, ".go") || strings.HasSuffix(p, ".json") || p == "go.mod" || p == "go.sum" || p == "third_party/chess/go.mod" || p == "third_party/chess/go.sum")
}

func replay(source, base, out, test string, check bool) error {
	var snapshot struct {
		Base  string            `json:"base_tag"`
		Files map[string]string `json:"files"`
	}
	if source != "" {
		if err := load(source, &snapshot); err != nil {
			return err
		}
		if snapshot.Base == "" {
			return fmt.Errorf("source has no base_tag")
		}
		base = snapshot.Base
	}
	if strings.HasPrefix(base, "-") {
		return fmt.Errorf("invalid base ref")
	}
	for p := range snapshot.Files {
		if !allowedOverlay(p) {
			return fmt.Errorf("overlay must contain only Go source/module files or JSON fixtures: %s", p)
		}
	}
	if err := os.Mkdir(out, 0755); err != nil {
		return err
	}
	cmd := exec.Command("git", "archive", base)
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	err = extract(tar.NewReader(stream), out)
	if err != nil {
		_ = cmd.Process.Kill()
	}
	waitErr := cmd.Wait()
	if err != nil {
		return err
	}
	if waitErr != nil {
		return waitErr
	}
	for p, text := range snapshot.Files {
		path := filepath.Join(out, p)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			return err
		}
	}
	if test != "" {
		b, err := os.ReadFile(test)
		if err != nil {
			return err
		}
		path := filepath.Join(out, "archived_proof_test.go")
		if err := os.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	if check {
		for _, step := range []struct {
			Name string
			Args []string
		}{{"tests", []string{"test", "-count=1", "./..."}}, {"build", []string{"build", "-o", "pandafish-replayed", "."}}} {
			cmd := exec.Command("go", step.Args...)
			cmd.Dir = out
			output, err := cmd.CombinedOutput()
			if writeErr := os.WriteFile(filepath.Join(out, step.Name+".txt"), output, 0644); writeErr != nil {
				return writeErr
			}
			if err != nil {
				return fmt.Errorf("%s failed (log retained): %w", step.Name, err)
			}
		}
	}
	return printJSON(map[string]any{"base": base, "directory": out, "overlay_files": len(snapshot.Files), "checks_ran": check})
}
