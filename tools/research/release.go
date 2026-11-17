package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corentings/chess/v2"
)

func checkPolicy(base, version string) (map[string]any, error) {
	if strings.HasPrefix(base, "-") {
		return nil, fmt.Errorf("invalid base ref")
	}
	listing, err := exec.Command("git", "ls-tree", "-r", "--name-only", base).Output()
	if err != nil {
		return nil, err
	}
	checked := map[string]bool{}
	for _, name := range strings.Fields(string(listing)) {
		if !((!strings.Contains(name, "/") && strings.HasSuffix(name, ".go")) || strings.HasPrefix(name, "third_party/chess/") || name == "go.mod" || name == "go.sum") {
			continue
		}
		old, err := exec.Command("git", "show", base+":"+name).Output()
		if err != nil {
			return nil, err
		}
		current, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if name == "main.go" || name == "main_test.go" {
			current = bytes.ReplaceAll(current, []byte(version), []byte(base))
		}
		if !bytes.Equal(current, old) {
			return nil, fmt.Errorf("playing/base source differs: %s", name)
		}
		checked[name] = true
	}
	for _, pattern := range []string{"*.go", "third_party/chess/*.go"} {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, name := range files {
			if !checked[name] && !strings.HasSuffix(name, "_test.go") {
				return nil, fmt.Errorf("new playing runtime file: %s", name)
			}
		}
	}
	return map[string]any{"base": base, "current_version": version, "original_files_checked": len(checked), "playing_policy_unchanged": true}, nil
}

func checkGames(dir string) (map[string]any, error) {
	var summary struct {
		Matches []struct {
			PGN   string `json:"pgn_file"`
			Games []struct {
				Result      string `json:"result"`
				Plies       int    `json:"plies"`
				Termination string `json:"termination"`
			} `json:"games"`
		} `json:"matches"`
	}
	if err := load(filepath.Join(dir, "summary.json"), &summary); err != nil {
		return nil, err
	}
	count := 0
	for _, match := range summary.Matches {
		if !safePath(match.PGN) {
			return nil, fmt.Errorf("unsafe PGN path")
		}
		f, err := os.Open(filepath.Join(dir, match.PGN))
		if err != nil {
			return nil, err
		}
		scanner := chess.NewScanner(f)
		index := 0
		for scanner.HasNext() {
			game, err := scanner.ParseNext()
			if err != nil {
				f.Close()
				return nil, err
			}
			if index >= len(match.Games) {
				f.Close()
				return nil, fmt.Errorf("extra PGN game")
			}
			row := match.Games[index]
			// The pinned parser can leave a claimed draw as NoOutcome. Check the
			// declared result and independently validate the claim, without changing it.
			if game.GetTagPair("Result") != row.Result || len(game.Moves()) != row.Plies {
				f.Close()
				return nil, fmt.Errorf("PGN result/length mismatch: %s game %d: parsed %s/%d, saved %s/%d", match.PGN, index+1, game.Outcome(), len(game.Moves()), row.Result, row.Plies)
			}
			if row.Termination == "ThreefoldRepetition" {
				eligible := false
				for _, method := range game.EligibleDraws() {
					if method == chess.ThreefoldRepetition {
						eligible = true
					}
				}
				if !eligible || row.Result != "1/2-1/2" {
					f.Close()
					return nil, fmt.Errorf("invalid repetition draw claim")
				}
			} else if row.Termination == "Checkmate" && (game.Position().SearchStatus() != chess.Checkmate || game.Outcome().String() != row.Result) {
				f.Close()
				return nil, fmt.Errorf("checkmate result mismatch")
			}
			p := game.Positions()[0]
			for _, move := range game.Moves() {
				found := false
				for _, legal := range p.ValidMoves() {
					if legal.String() == move.String() {
						p = p.Update(&legal)
						found = true
						break
					}
				}
				if !found {
					f.Close()
					return nil, fmt.Errorf("illegal saved move %s", move)
				}
			}
			index++
			count++
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
		if index != len(match.Games) {
			return nil, fmt.Errorf("missing PGN games")
		}
	}
	if count == 0 {
		return nil, fmt.Errorf("no games")
	}
	return map[string]any{"legal_games_replayed": count, "matchups": len(summary.Matches), "summary_results_and_lengths_match": true}, nil
}
