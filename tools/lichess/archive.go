package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

type archiveFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Bytes  int64  `json:"bytes"`
}
type archiveIndex struct {
	Schema int           `json:"schema"`
	Files  []archiveFile `json:"files"`
}

func digest(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }

var credentialPattern = regexp.MustCompile(`(?i)lip_[a-z0-9]{12,}|Bearer\s+[a-z0-9._-]{16,}`)

func inventory(dir string, scan bool) (archiveIndex, error) {
	index := archiveIndex{Schema: 1}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive refuses symlink: %s", d.Name())
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("archive requires regular files")
		}
		relative, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if scan && relative != "pandafish" {
			if credentialPattern.Match(data) {
				return fmt.Errorf("possible credential in %s; original preserved, archive stopped", relative)
			}
			if strings.HasSuffix(relative, ".yml") || strings.HasSuffix(relative, ".yaml") {
				var value any
				if yaml.Unmarshal(data, &value) != nil {
					return fmt.Errorf("invalid archive YAML: %s", relative)
				}
				if !reflect.DeepEqual(value, redactConfig(value)) {
					return fmt.Errorf("unredacted credential field in %s", relative)
				}
			}
		}
		index.Files = append(index.Files, archiveFile{relative, digest(data), int64(len(data))})
		return nil
	})
	return index, err
}
func archiveSession(source, destination string) error {
	source, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}
	if destination == source || strings.HasPrefix(destination, source+string(os.PathSeparator)) {
		return fmt.Errorf("archive cannot be inside source session")
	}
	data, err := os.ReadFile(filepath.Join(source, "manifest.json"))
	if err != nil {
		return err
	}
	var m manifest
	if json.Unmarshal(data, &m) != nil || m.Schema != 1 || (m.Status != "finished" && m.Status != "client_failed" && m.Status != "capture_failed") {
		return fmt.Errorf("archive requires a finalized session (finished or failed)")
	}
	if source == filepath.Join(destination, "raw") {
		return verifyArchive(destination)
	}
	for _, artifact := range []struct{ name, hash string }{{"pandafish", m.EngineSHA256}, {"config.redacted.yml", m.ConfigSHA256}} {
		b, e := os.ReadFile(filepath.Join(source, artifact.name))
		if os.IsNotExist(e) && artifact.name == "pandafish" {
			continue
		}
		if e != nil {
			return e
		}
		if digest(b) != artifact.hash {
			return fmt.Errorf("source %s does not match capture manifest", artifact.name)
		}
	}
	index, err := inventory(source, true)
	if err != nil {
		return err
	}
	if _, err = os.Stat(destination); err == nil {
		saved, e := os.ReadFile(filepath.Join(destination, "inventory.json"))
		if e != nil {
			return fmt.Errorf("existing archive has no inventory")
		}
		var previous archiveIndex
		if json.Unmarshal(saved, &previous) != nil || !reflect.DeepEqual(previous, index) {
			return fmt.Errorf("source differs from archive inventory; refusing overwrite")
		}
		return verifyArchive(destination)
	} else if !os.IsNotExist(err) {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(filepath.Dir(destination), ".archive-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	for _, f := range index.Files {
		data, e := os.ReadFile(filepath.Join(source, filepath.FromSlash(f.Path)))
		if e != nil {
			return e
		}
		if digest(data) != f.SHA256 {
			return fmt.Errorf("source changed while archiving")
		}
		target := filepath.Join(staging, "raw", filepath.FromSlash(f.Path))
		if e = os.MkdirAll(filepath.Dir(target), 0755); e != nil {
			return e
		}
		mode := os.FileMode(0644)
		if f.Path == "pandafish" {
			mode = 0755
		}
		if e = os.WriteFile(target, data, mode); e != nil {
			return e
		}
	}
	if err = writeJSON(filepath.Join(staging, "inventory.json"), index); err != nil {
		return err
	}
	final, e := inventory(source, true)
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(final, index) {
		return fmt.Errorf("source changed during archive copy")
	}
	return os.Rename(staging, destination)
}

// The executable is retained locally but excluded from Git. Every data file is mandatory.
func verifyArchive(destination string) error {
	saved, err := os.ReadFile(filepath.Join(destination, "inventory.json"))
	if err != nil {
		return err
	}
	var expected archiveIndex
	if json.Unmarshal(saved, &expected) != nil || expected.Schema != 1 {
		return fmt.Errorf("invalid archive inventory")
	}
	actual, err := inventory(filepath.Join(destination, "raw"), true)
	if err != nil {
		return err
	}
	byPath := map[string]archiveFile{}
	for _, f := range actual.Files {
		byPath[f.Path] = f
	}
	for _, f := range expected.Files {
		got, ok := byPath[f.Path]
		if !ok && f.Path == "pandafish" {
			continue
		}
		if !ok || got != f {
			return fmt.Errorf("archived raw file changed or missing: %s", f.Path)
		}
		delete(byPath, f.Path)
	}
	if len(byPath) > 0 {
		return fmt.Errorf("unrecorded file added to raw archive")
	}
	return nil
}
