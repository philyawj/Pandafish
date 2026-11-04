// Command benchmarks plays local UCI matches and saves their history.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/corentings/chess/v2"
)

func stockfishOptions(skill int) map[string]string {
	return map[string]string{
		"Skill Level": fmt.Sprint(skill), "UCI_LimitStrength": "false", "Threads": "1", "Hash": "16", "SyzygyPath": "",
	}
}

type identity struct {
	Label    string            `json:"label"`
	UCIName  string            `json:"uci_name"`
	SHA256   string            `json:"sha256"`
	Options  map[string]string `json:"options"`
	Revision string            `json:"source_revision,omitempty"`
}
type settings struct {
	Pairs          int    `json:"pairs"`
	MoveMS         int    `json:"move_ms"`
	GraceMS        int    `json:"timeout_grace_ms"`
	MaxPlies       int    `json:"max_plies"`
	OpeningsSHA256 string `json:"openings_sha256"`
	OpeningsName   string `json:"openings_name,omitempty"`
	TimingPolicy   string `json:"timing_policy,omitempty"`
}
type moveTiming struct {
	Candidate bool    `json:"candidate"`
	ElapsedMS float64 `json:"elapsed_ms"`
}
type gameRecord struct {
	CandidateWhite bool         `json:"candidate_white"`
	Opening        []string     `json:"opening"`
	Result         string       `json:"result"`
	Termination    string       `json:"termination"`
	Plies          int          `json:"plies"`
	Timings        []moveTiming `json:"response_times,omitempty"`
}
type stats struct {
	Wins       int      `json:"wins"`
	Draws      int      `json:"draws"`
	Losses     int      `json:"losses"`
	Unfinished int      `json:"unfinished"`
	Forfeits   int      `json:"forfeits"`
	Score      *float64 `json:"score_percent"`
	ScoreLower *float64 `json:"possible_score_lower_percent"`
	ScoreUpper *float64 `json:"possible_score_upper_percent"`
}
type match struct {
	Engine identity     `json:"engine"`
	Stats  stats        `json:"stats"`
	Games  []gameRecord `json:"games"`
}
type summary struct {
	Schema         int      `json:"schema"`
	Created        string   `json:"created_utc"`
	Platform       string   `json:"platform"`
	GoVersion      string   `json:"go_version"`
	Candidate      identity `json:"candidate"`
	Settings       settings `json:"settings"`
	Matches        []match  `json:"matches"`
	Complete       bool     `json:"complete"`
	PlannedMatches int      `json:"planned_matches,omitempty"`
	RunnerSHA256   string   `json:"runner_sha256,omitempty"`
	Machine        string   `json:"machine,omitempty"`
}

// Copy first, then identify and launch only the copy. Rebuilding the original
// executable during a long run cannot change which engine plays later games.
func snapshotEngine(source, root, role string) (string, error) {
	dir := filepath.Join(root, role)
	if err := os.Mkdir(dir, 0700); err != nil {
		return "", err
	}
	in, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer in.Close()
	path := filepath.Join(dir, filepath.Base(source))
	out, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0700)
	if err != nil {
		return "", err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return "", copyErr
	}
	return path, closeErr
}

// Replace the summary atomically so an interrupted write preserves the last
// completed matchup. Old summaries remain readable by the report generator.
func saveSummary(dir string, s summary) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, "summary-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = f.Chmod(0644); err != nil {
		f.Close()
		return err
	}
	if _, err = f.Write(append(data, '\n')); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filepath.Join(dir, "summary.json"))
}

func digest(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(b)), nil
}

func identify(path, label string, options map[string]string) (identity, error) {
	e, err := startEngine(path, options)
	if err != nil {
		return identity{}, err
	}
	defer e.close()
	hash, err := digest(path)
	return identity{Label: label, UCIName: e.name, SHA256: hash, Options: options}, err
}

func newGame(opening []string) (*chess.Game, error) {
	game := chess.NewGame()
	for _, move := range opening {
		if err := game.PushNotationMove(move, chess.UCINotation{}, nil); err != nil {
			return nil, err
		}
	}
	return game, nil
}

func loadOpenings(path string) ([][]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var openings [][]string
	for _, line := range strings.Split(string(b), "\n") {
		line, _, _ = strings.Cut(line, "#")
		moves := strings.Fields(line)
		if len(moves) == 0 {
			continue
		}
		game, err := newGame(moves)
		if err != nil {
			return nil, fmt.Errorf("opening: %w", err)
		}
		if game.Outcome() != chess.NoOutcome {
			return nil, fmt.Errorf("opening is terminal")
		}
		openings = append(openings, moves)
	}
	if len(openings) == 0 {
		return nil, fmt.Errorf("no openings")
	}
	return openings, nil
}

func playGame(candidate, opponent string, options map[string]string, opening []string, white bool, cfg settings) (*chess.Game, gameRecord, error) {
	game, err := newGame(opening)
	if err != nil {
		return nil, gameRecord{}, err
	}
	a, err := startEngine(candidate, nil)
	if err != nil {
		return nil, gameRecord{}, err
	}
	defer a.close()
	b, err := startEngine(opponent, options)
	if err != nil {
		return nil, gameRecord{}, err
	}
	defer b.close()
	history := append([]string(nil), opening...)
	var timings []moveTiming
	termination := "move cap"
	for game.Outcome() == chess.NoOutcome {
		for _, method := range game.EligibleDraws() {
			if method == chess.ThreefoldRepetition || method == chess.FiftyMoveRule {
				_ = game.Draw(method)
				break
			}
		}
		if game.Outcome() != chess.NoOutcome || len(history) >= cfg.MaxPlies {
			break
		}
		e := a
		if (game.Position().Turn() == chess.White) != white {
			e = b
		}
		move, elapsed, moveErr := e.move(history, time.Duration(cfg.MoveMS)*time.Millisecond, time.Duration(cfg.GraceMS)*time.Millisecond)
		timings = append(timings, moveTiming{Candidate: e == a, ElapsedMS: float64(elapsed) / float64(time.Millisecond)})
		if moveErr == nil {
			moveErr = game.PushNotationMove(move, chess.UCINotation{}, nil)
		}
		if moveErr != nil {
			termination = "forfeit: " + moveErr.Error()
			game.Resign(game.Position().Turn())
			break
		}
		history = append(history, move)
	}
	if game.Outcome() != chess.NoOutcome && !strings.HasPrefix(termination, "forfeit:") {
		termination = game.Method().String()
	}
	record := gameRecord{CandidateWhite: white, Opening: opening, Result: game.Outcome().String(), Termination: termination, Plies: len(history), Timings: timings}
	return game, record, nil
}

func summarize(games []gameRecord) stats {
	s := stats{}
	for _, game := range games {
		switch game.Result {
		case "*":
			s.Unfinished++
		case "1/2-1/2":
			s.Draws++
		default:
			if (game.Result == "1-0") == game.CandidateWhite {
				s.Wins++
			} else {
				s.Losses++
			}
		}
		if strings.HasPrefix(game.Termination, "forfeit:") {
			s.Forfeits++
		}
	}
	if completed := s.Wins + s.Draws + s.Losses; completed > 0 {
		value := 100 * (float64(s.Wins) + float64(s.Draws)/2) / float64(completed)
		s.Score = &value
	}
	if total := s.Wins + s.Draws + s.Losses + s.Unfinished; total > 0 {
		lower := 100 * (float64(s.Wins) + float64(s.Draws)/2) / float64(total)
		upper := lower + 100*float64(s.Unfinished)/float64(total)
		s.ScoreLower, s.ScoreUpper = &lower, &upper
	}
	return s
}

func folderLabel(label string) string {
	label = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, strings.ToLower(label))
	if label = strings.Trim(label, ".-_"); label == "" {
		return "engine"
	}
	return label
}

func createRunDir(root, version, previous string, skill int) (string, error) {
	base := folderLabel(version) + "-vs-"
	if previous != "" {
		base += folderLabel(previous) + "-and-"
	}
	base += fmt.Sprintf("stockfish-%d", skill)
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	for run := 1; ; run++ {
		name := base
		if run > 1 {
			name = fmt.Sprintf("%s-run-%d", base, run)
		}
		dir := filepath.Join(root, name)
		if err := os.Mkdir(dir, 0755); err == nil {
			return dir, nil
		} else if !os.IsExist(err) {
			return "", err
		}
	}
}

func run() error {
	candidate := flag.String("candidate", "", "Candidate UCI executable (required)")
	version := flag.String("version", "", "Candidate version, e.g. v0.2 (required)")
	previous := flag.String("previous", "", "Previous version's executable (optional)")
	previousVersion := flag.String("previous-version", "", "Previous version label")
	stockfish := flag.String("stockfish", "stockfish", "Stockfish executable")
	skill := flag.Int("stockfish-skill", 1, "Stockfish Skill Level (1 through 20)")
	pairs := flag.Int("pairs", 50, "Color pairs per opponent")
	moveMS := flag.Int("move-ms", 100, "Requested milliseconds per move")
	graceMS := flag.Int("grace-ms", 50, "Additional wall-clock allowance for protocol overhead")
	revision := flag.String("revision", "", "Candidate source revision; append -dirty for modified source")
	previousRevision := flag.String("previous-revision", "", "Previous engine source revision")
	machine := flag.String("machine", "", "Hardware and power-setting note for comparisons")
	maxPlies := flag.Int("max-plies", 400, "Total plies before marking unfinished")
	results := flag.String("results", "benchmarks/results", "Results directory")
	reportOnly := flag.Bool("report-only", false, "Rebuild the Markdown table from saved results without playing")
	openingsPath := flag.String("openings", "benchmarks/openings-screening-10-v1.txt", "Opening suite")
	flag.Parse()
	if *reportOnly {
		return report(*results)
	}
	if *candidate == "" || *version == "" || (*previous == "") != (*previousVersion == "") || *pairs < 1 || *moveMS < 1 || *maxPlies < 10 || *graceMS < 0 || *skill < 1 || *skill > 20 {
		return fmt.Errorf("candidate/version required; previous flags paired; pairs/move-ms positive; max-plies >= 10; grace-ms >= 0; stockfish-skill 1..20")
	}
	var err error
	if *candidate, err = exec.LookPath(*candidate); err != nil {
		return err
	}
	if *stockfish, err = exec.LookPath(*stockfish); err != nil {
		return err
	}
	paths, labels, options := []string{}, []string{}, []map[string]string{}
	if *previous != "" {
		if *previous, err = exec.LookPath(*previous); err != nil {
			return err
		}
		paths, labels, options = append(paths, *previous), append(labels, *previousVersion), append(options, nil)
	}
	paths, labels, options = append(paths, *stockfish), append(labels, fmt.Sprintf("Stockfish Skill %d", *skill)), append(options, stockfishOptions(*skill))
	openings, err := loadOpenings(*openingsPath)
	if err != nil {
		return err
	}
	for _, opening := range openings {
		if len(opening) >= *maxPlies {
			return fmt.Errorf("move cap must exceed opening length")
		}
	}
	hash, err := digest(*openingsPath)
	if err != nil {
		return err
	}
	cfg := settings{Pairs: *pairs, MoveMS: *moveMS, GraceMS: *graceMS, MaxPlies: *maxPlies, OpeningsSHA256: hash,
		OpeningsName: filepath.Base(*openingsPath), TimingPolicy: "wall-clock-v2"}
	snapshots, err := os.MkdirTemp("", "pandafish-benchmark-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(snapshots)
	if *candidate, err = snapshotEngine(*candidate, snapshots, "candidate"); err != nil {
		return err
	}
	for i, path := range paths {
		paths[i], err = snapshotEngine(path, snapshots, fmt.Sprint("opponent-", i))
		if err != nil {
			return err
		}
	}
	id, err := identify(*candidate, *version, nil)
	if err != nil {
		return err
	}
	id.Revision = *revision
	identities := make([]identity, len(paths))
	for i, path := range paths {
		identities[i], err = identify(path, labels[i], options[i])
		if err != nil {
			return err
		}
		if *previous != "" && i == 0 {
			identities[i].Revision = *previousRevision
		}
	}
	now := time.Now().UTC()
	runner, err := os.Executable()
	if err != nil {
		return err
	}
	runnerHash, err := digest(runner)
	if err != nil {
		return err
	}
	s := summary{Schema: 2, Created: now.Format(time.RFC3339Nano), Platform: runtime.GOOS + "/" + runtime.GOARCH,
		GoVersion: runtime.Version(), Candidate: id, Settings: cfg, PlannedMatches: len(paths), RunnerSHA256: runnerHash, Machine: *machine}
	dir, err := createRunDir(*results, *version, *previousVersion, *skill)
	if err != nil {
		return err
	}
	pgn, err := os.Create(filepath.Join(dir, "games.pgn"))
	if err != nil {
		return err
	}
	defer pgn.Close()
	if err = saveSummary(dir, s); err != nil {
		return err
	}
	pgnSeparator := ""
	for i, path := range paths {
		m := match{Engine: identities[i]}
		for pair := 0; pair < cfg.Pairs; pair++ {
			for _, white := range []bool{true, false} {
				game, record, err := playGame(*candidate, path, options[i], openings[pair%len(openings)], white, cfg)
				if err != nil {
					return err
				}
				w, b := id.UCIName, identities[i].UCIName
				if !white {
					w, b = b, w
				}
				for key, value := range map[string]string{"Event": "Pandafish local benchmark", "Site": "local",
					"Date": now.Format("2006.01.02"), "White": w, "Black": b, "Round": fmt.Sprint(len(m.Games) + 1),
					"Result": record.Result, "Termination": record.Termination, "TimeControl": "-"} {
					game.AddTagPair(key, value)
				}
				if _, err = fmt.Fprintf(pgn, "%s%s\n", pgnSeparator, game.String()); err != nil {
					return err
				}
				pgnSeparator = "\n"
				m.Games = append(m.Games, record)
			}
			fmt.Printf("%s vs %s: %d/%d games\n", *version, labels[i], len(m.Games), 2*cfg.Pairs)
		}
		m.Stats = summarize(m.Games)
		s.Matches = append(s.Matches, m)
		if err = pgn.Sync(); err != nil {
			return err
		}
		if err = saveSummary(dir, s); err != nil {
			return err
		}
		if err = report(*results); err != nil {
			return err
		}
	}
	if err = pgn.Close(); err != nil {
		return err
	}
	s.Complete = true
	if err = saveSummary(dir, s); err != nil {
		return err
	}
	if err = report(*results); err != nil {
		return err
	}
	fmt.Println("Saved", dir)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
