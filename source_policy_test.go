package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// Cover current and archived source, including script-name.txt disguises.
// External clients and ignored local engine/checkouts are outside the source tree.
func TestRepositorySourceIsGo(t *testing.T) {
	foreign := map[string]bool{"py": true, "sh": true, "bash": true, "zsh": true,
		"js": true, "ts": true, "jsx": true, "tsx": true, "rs": true, "c": true,
		"cpp": true, "h": true, "rb": true, "pl": true, "lua": true}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "bin" ||
				d.Name() == "local-games" || d.Name() == "__pycache__" || d.Name() == "venv" ||
				d.Name() == "tmp" || d.Name() == "temp" || d.Name() == "build" || d.Name() == "dist") {
				return filepath.SkipDir
			}
			return nil
		}
		parts := strings.Split(strings.ToLower(d.Name()), ".")
		for _, part := range parts[1:] {
			if foreign[part] {
				t.Errorf("non-Go source file: %s", path)
				break
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
