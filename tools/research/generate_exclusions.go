package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/corentings/chess/v2"
)

type pathFlags []string

func (p *pathFlags) String() string     { return fmt.Sprint([]string(*p)) }
func (p *pathFlags) Set(s string) error { *p = append(*p, s); return nil }

func blockAdditionalProbeSources(blocked map[string]bool, inputs, analysis []string) ([]map[string]any, error) {
	sources := []map[string]any{}
	for kind, patterns := range [][]string{inputs, analysis} {
		for _, pattern := range patterns {
			paths, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			if len(paths) == 0 {
				return nil, fmt.Errorf("exclusion has no matches: %s", pattern)
			}
			for _, path := range paths {
				b, err := os.ReadFile(path)
				if err != nil {
					return nil, err
				}
				count := 0
				ignored := false
				if kind == 0 {
					var data struct{ Cases []probeCase }
					if err = json.Unmarshal(b, &data); err != nil {
						return nil, err
					}
					if len(data.Cases) == 0 {
						return nil, fmt.Errorf("empty case exclusion %s", path)
					}
					for _, c := range data.Cases {
						g, e := probeGame(c)
						if e != nil {
							return nil, e
						}
						blocked[canonicalProbe(g.Position())] = true
						count++
					}
				} else {
					trimmed := bytes.TrimSpace(b)
					if len(trimmed) == 0 || !json.Valid(trimmed) {
						return nil, fmt.Errorf("invalid analysis JSON: %s", path)
					}
					if trimmed[0] != '{' {
						sources = append(sources, map[string]any{"path": path, "identity": hash(b), "positions": 0, "ignored_non_game_metadata": true})
						continue
					}
					var data struct {
						Game     string `json:"game"`
						Complete bool   `json:"complete"`
						Moves    []struct {
							FEN string `json:"fen"`
						} `json:"moves"`
					}
					if err = json.Unmarshal(b, &data); err != nil {
						return nil, err
					}
					if data.Game == "" {
						ignored = true
					} else {
						if len(data.Moves) == 0 && !data.Complete {
							return nil, fmt.Errorf("game analysis has no positions: %s", path)
						}
						if len(data.Moves) == 0 {
							sources = append(sources, map[string]any{"path": path, "identity": hash(b), "positions": 0, "empty_completed_game": true})
							continue
						}
						for _, m := range data.Moves {
							option, e := chess.FEN(m.FEN)
							if e != nil {
								return nil, e
							}
							blocked[canonicalProbe(chess.NewGame(option).Position())] = true
							count++
						}
					}
				}
				sources = append(sources, map[string]any{"path": path, "identity": hash(b), "positions": count, "ignored_non_game_metadata": ignored})
			}
		}
	}
	return sources, nil
}
