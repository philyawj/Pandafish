package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type entry struct {
	SHA   string `json:"sha256"`
	Bytes int    `json:"bytes"`
}

func identity(b []byte) entry { return entry{fmt.Sprintf("%x", sha256.Sum256(b)), len(b)} }
func must(e error) {
	if e != nil {
		panic(e)
	}
}
func read(p string) []byte { b, e := os.ReadFile(p); must(e); return b }
func writeNew(p string, b []byte) {
	f, e := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	must(e)
	_, e = f.Write(b)
	must(e)
	must(f.Close())
}
func main() {
	root := "tools/benchmark/research/v0.23-state-kernel"
	out := "tools/research/validation/v0.23-release"
	inv := read(filepath.Join(root, "inventory.json"))
	var entries map[string]entry
	must(json.Unmarshal(inv, &entries))
	for name, want := range entries {
		if !filepath.IsLocal(name) || identity(read(filepath.Join(root, name))) != want {
			panic("frozen evidence mismatch: " + name)
		}
	}
	before := read(filepath.Join(root, "README.md"))
	old := []byte("| --- | --- | --- | --- | --- | --- |")
	if bytes.Count(before, old) != 1 {
		panic("expected one malformed five-column delimiter")
	}
	after := bytes.Replace(before, old, []byte("| --- | --- | --- | --- | --- |"), 1)
	writeNew(filepath.Join(out, "candidate-readme-before-formatting.md"), before)
	writeNew(filepath.Join(out, "candidate-inventory-before-formatting.json"), inv)
	must(os.WriteFile(filepath.Join(root, "README.md"), after, 0644))
	for name, want := range entries {
		if name != "README.md" && identity(read(filepath.Join(root, name))) != want {
			panic("raw evidence changed: " + name)
		}
	}
	report := map[string]any{"markdown_only": true, "reason": "five header cells require five delimiter cells", "original_entries_verified": len(entries), "raw_entries_unchanged": len(entries) - 1, "readme_before": identity(before), "readme_after": identity(after), "original_inventory": identity(inv), "originals_preserved_in": out, "new_inventory": "rebuilt explicitly after a presentation-only correction, never to conceal changed observations"}
	b, e := json.MarshalIndent(report, "", "  ")
	must(e)
	writeNew(filepath.Join(out, "readme-amendment.json"), append(b, '\n'))
	must(os.Remove(filepath.Join(root, "inventory.json")))
	fmt.Println("PASS original document/inventory preserved; raw entries unchanged; regenerate the explicit current inventory")
}
