package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type opponent struct {
	Path, Label, Revision, Role, AssetPath, AssetSHA string
	Options                                          map[string]string
	Identity                                         identity
}

func validateSchedule(stage, previous, sacrifice string, validation bool) error {
	if stage != "all" && stage != "previous" && stage != "external" {
		return fmt.Errorf("stage must be all, previous, or external")
	}
	if (stage == "previous" || (validation && stage == "all")) && previous == "" {
		return fmt.Errorf("previous-release stage requires --previous")
	}
	if validation && stage != "previous" && sacrifice == "" {
		return fmt.Errorf("validation external comparisons require --sacrifice")
	}
	return nil
}

// Full suites default to three-way concurrency. Explicit serial mode retains
// the previous-first gate; constrained hosts fall back to that mode as well.
func scheduleConcurrency(stage string, previous, sacrifice, parallel, parallelAll, sequential bool, cpus int) (int, int, error) {
	selected := 0
	for _, enabled := range []bool{parallel, parallelAll, sequential} {
		if enabled {
			selected++
		}
	}
	if selected > 1 {
		return 0, 0, fmt.Errorf("choose only one of parallel, parallel-all, or sequential")
	}
	if parallelAll && (stage != "all" || !previous || !sacrifice || cpus < 4) {
		return 0, 0, fmt.Errorf("parallel-all requires all three opponents, stage all, and at least four logical CPUs")
	}
	if !sequential && !parallel && stage == "all" && previous && sacrifice && cpus >= 4 {
		return 3, 2, nil
	}
	if parallel && cpus >= 4 {
		return 0, 2, nil
	}
	return 0, 1, nil
}

// Conservative screen, not a significance test: unfinished games count as losses.
func favorablePrevious(s stats) bool {
	total := s.Wins + s.Draws + s.Losses + s.Unfinished
	return total > 0 && s.Forfeits == 0 && 2*s.Wins+s.Draws > total
}
func runMatch(candidate string, id identity, o opponent, openings [][]string, cfg settings, dir string, now time.Time, verbose bool) (match, error) {
	m := match{Engine: o.Identity, PGN: o.Role + ".pgn"}
	if o.AssetPath != "" {
		hash, err := digest(o.AssetPath)
		if err != nil || hash != o.AssetSHA {
			return m, fmt.Errorf("Sacrifice asset changed before match")
		}
	}
	f, err := os.OpenFile(filepath.Join(dir, m.PGN), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return m, err
	}
	defer f.Close()
	fmt.Printf("Starting %s vs %s: %d games\n", id.Label, o.Label, 2*cfg.Pairs)
	separator := ""
	for pair := 0; pair < cfg.Pairs; pair++ {
		for _, white := range []bool{true, false} {
			game, record, e := playGame(candidate, o.Path, o.Options, openings[pair%len(openings)], white, cfg)
			if e != nil {
				return m, e
			}
			w, b := id.UCIName, o.Identity.UCIName
			if !white {
				w, b = b, w
			}
			for k, v := range map[string]string{"Event": "Pandafish local benchmark", "Site": "local", "Date": now.Format("2006.01.02"), "White": w, "Black": b, "Round": fmt.Sprint(len(m.Games) + 1), "Result": record.Result, "Termination": record.Termination, "TimeControl": "-"} {
				game.AddTagPair(k, v)
			}
			if _, e = fmt.Fprintf(f, "%s%s\n", separator, game.String()); e != nil {
				return m, e
			}
			separator = "\n"
			m.Games = append(m.Games, record)
		}
		if err = f.Sync(); err != nil {
			return m, err
		}
		if verbose {
			fmt.Printf("%s vs %s: %d/%d games\n", id.Label, o.Label, len(m.Games), 2*cfg.Pairs)
		}
	}
	if o.AssetPath != "" {
		hash, e := digest(o.AssetPath)
		if e != nil || hash != o.AssetSHA {
			return m, fmt.Errorf("Sacrifice asset changed during match")
		}
	}
	m.Stats = summarize(m.Games)
	fmt.Printf("%s vs %s: W/D/L %d/%d/%d; unfinished %d; forfeits %d; possible score %s\n", id.Label, o.Label, m.Stats.Wins, m.Stats.Draws, m.Stats.Losses, m.Stats.Unfinished, m.Stats.Forfeits, rangeText(m.Stats))
	return m, nil
}
func combinePGNs(dir string, matches []match) error {
	var all strings.Builder
	for _, m := range matches {
		data, err := os.ReadFile(filepath.Join(dir, m.PGN))
		if err != nil {
			return err
		}
		if all.Len() > 0 {
			all.WriteByte('\n')
		}
		all.Write(data)
	}
	return os.WriteFile(filepath.Join(dir, "games.pgn"), []byte(all.String()), 0644)
}

type matchupResult struct {
	index int
	match match
	err   error
}

type matchPlayer func(string, identity, opponent, [][]string, settings, string, time.Time, bool) (match, error)

func playSchedule(candidate string, id identity, opponents []opponent, openings [][]string, cfg settings, dir string, now time.Time, verbose bool, save func([]match) error, gate func(string)) ([]match, error) {
	return playScheduleWith(candidate, id, opponents, openings, cfg, dir, now, verbose, save, gate, runMatch)
}
func playScheduleWith(candidate string, id identity, opponents []opponent, openings [][]string, cfg settings, dir string, now time.Time, verbose bool, save func([]match) error, gate func(string), play matchPlayer) ([]match, error) {
	concurrency := max(1, min(2, cfg.ExternalConcurrency))
	concurrentPrevious := cfg.MatchupConcurrency > 1
	if concurrentPrevious {
		concurrency = min(3, cfg.MatchupConcurrency)
		gate("pending: previous screen runs concurrently")
	}
	previousFailed := false
	matches := make([]*match, len(opponents))
	completed := func() []match {
		var out []match
		for _, m := range matches {
			if m != nil {
				out = append(out, *m)
			}
		}
		return out
	}
	first := 0
	if !concurrentPrevious && len(opponents) > 0 && opponents[0].Role == "previous" {
		m, err := play(candidate, id, opponents[0], openings, cfg, dir, now, verbose)
		if err != nil {
			return nil, err
		}
		matches[0] = &m
		first = 1
		if len(opponents) > 1 {
			if !favorablePrevious(m.Stats) {
				gate("blocked: no favorable previous-release screen or forfeits")
				if err = save(completed()); err != nil {
					return completed(), err
				}
				return completed(), fmt.Errorf("previous-release gate failed; external comparisons skipped")
			}
			gate("passed: score above 50% treating unfinished as losses, zero forfeits")
		}
		if err = save(completed()); err != nil {
			return completed(), err
		}
	}
	// Only this goroutine updates shared summaries; workers own distinct PGNs.
	results := make(chan matchupResult, len(opponents)-first)
	next, running, finished := first, 0, first
	for finished < len(opponents) {
		for next < len(opponents) && running < concurrency {
			i := next
			next++
			running++
			go func() {
				m, e := play(candidate, id, opponents[i], openings, cfg, dir, now, verbose)
				results <- matchupResult{i, m, e}
			}()
		}
		r := <-results
		running--
		finished++
		if r.err != nil {
			// Reap remaining jobs rather than leaving writers/processes after return.
			for running > 0 {
				other := <-results
				running--
				if other.err == nil {
					m := other.match
					matches[other.index] = &m
				}
			}
			if e := save(completed()); e != nil {
				return completed(), e
			}
			return completed(), r.err
		}
		m := r.match
		matches[r.index] = &m
		if concurrentPrevious && opponents[r.index].Role == "previous" {
			previousFailed = !favorablePrevious(m.Stats)
			if previousFailed {
				gate("blocked: no favorable previous-release screen or forfeits; concurrent results retained")
			} else {
				gate("passed: score above 50% treating unfinished as losses, zero forfeits")
			}
		}
		if err := save(completed()); err != nil {
			for running > 0 {
				<-results
				running--
			}
			return completed(), err
		}
	}
	if previousFailed {
		return completed(), fmt.Errorf("previous-release gate failed; all concurrent matchups retained")
	}
	return completed(), nil
}
