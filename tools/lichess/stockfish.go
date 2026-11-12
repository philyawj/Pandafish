package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type sfScore struct {
	Kind  string `json:"kind"`
	Value int    `json:"value"`
}
type sfResult struct {
	Score              sfScore  `json:"score"`
	Depth              int      `json:"depth"`
	Nodes              int64    `json:"nodes"`
	TimeMS             int64    `json:"time_ms"`
	PV                 []string `json:"pv"`
	Bestmove           string   `json:"bestmove"`
	Output             []string `json:"output"`
	SelectionOutput    []string `json:"selection_output,omitempty"`
	ReferenceRechecked bool     `json:"reference_rechecked,omitempty"`
}
type stockfish struct {
	cmd    *exec.Cmd
	in     io.WriteCloser
	lines  chan string
	done   chan struct{}
	cancel context.CancelFunc
	name   string
}

func openStockfish(path string) (*stockfish, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stockfish{cmd: exec.CommandContext(ctx, path), lines: make(chan string), done: make(chan struct{}), cancel: cancel}
	var err error
	s.in, err = s.cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	out, err := s.cmd.StdoutPipe()
	if err != nil {
		s.in.Close()
		cancel()
		return nil, err
	}
	if err = s.cmd.Start(); err != nil {
		s.in.Close()
		cancel()
		return nil, err
	}
	go func() {
		defer close(s.lines)
		scanner := bufio.NewScanner(out)
		scanner.Buffer(make([]byte, 4096), 1024*1024)
		for scanner.Scan() {
			select {
			case s.lines <- scanner.Text():
			case <-s.done:
				return
			}
		}
	}()
	fail := func(err error) (*stockfish, error) { s.close(); return nil, err }
	if err = s.send("uci"); err != nil {
		return fail(err)
	}
	lines, err := s.until("uciok", 10*time.Second)
	if err != nil {
		return fail(err)
	}
	options := map[string]bool{}
	for _, line := range lines {
		if strings.HasPrefix(line, "id name ") {
			s.name = strings.TrimPrefix(line, "id name ")
		}
		if strings.HasPrefix(line, "option name ") {
			key, _, _ := strings.Cut(strings.TrimPrefix(line, "option name "), " type ")
			options[key] = true
		}
	}
	if !strings.HasPrefix(s.name, "Stockfish ") {
		return fail(fmt.Errorf("analysis requires Stockfish, got %q", s.name))
	}
	for _, option := range []struct{ key, value string }{{"Threads", "1"}, {"Hash", "64"}, {"MultiPV", "1"}, {"Skill Level", "20"}, {"UCI_LimitStrength", "false"}, {"Ponder", "false"}} {
		if !options[option.key] {
			return fail(fmt.Errorf("Stockfish does not advertise %s", option.key))
		}
		if err = s.send("setoption name " + option.key + " value " + option.value); err != nil {
			return fail(err)
		}
	}
	if err = s.ready(); err != nil {
		return fail(err)
	}
	return s, nil
}
func (s *stockfish) send(line string) error { _, err := fmt.Fprintln(s.in, line); return err }
func (s *stockfish) until(prefix string, timeout time.Duration) ([]string, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	var lines []string
	for {
		select {
		case line, ok := <-s.lines:
			if !ok {
				return nil, fmt.Errorf("Stockfish exited or output was unreadable")
			}
			lines = append(lines, line)
			if line == prefix || strings.HasPrefix(line, prefix+" ") {
				return lines, nil
			}
		case <-timer.C:
			s.cancel()
			return nil, fmt.Errorf("Stockfish timeout waiting for %s", prefix)
		}
	}
}
func (s *stockfish) ready() error {
	if err := s.send("isready"); err != nil {
		return err
	}
	_, err := s.until("readyok", 10*time.Second)
	return err
}
func (s *stockfish) close() { close(s.done); s.in.Close(); s.cancel(); _ = s.cmd.Wait() }
func parseSFInfo(line string) (sfResult, bool) {
	fields := strings.Fields(line)
	r := sfResult{}
	exact := true
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "lowerbound", "upperbound":
			exact = false
		case "score":
			if i+2 < len(fields) {
				v, e := strconv.Atoi(fields[i+2])
				if e == nil && (fields[i+1] == "cp" || fields[i+1] == "mate") {
					r.Score = sfScore{fields[i+1], v}
				}
				i += 2
			}
		case "depth", "nodes", "time":
			if i+1 < len(fields) {
				v, _ := strconv.ParseInt(fields[i+1], 10, 64)
				switch fields[i] {
				case "depth":
					r.Depth = int(v)
				case "nodes":
					r.Nodes = v
				case "time":
					r.TimeMS = v
				}
				i++
			}
		case "pv":
			r.PV = append([]string(nil), fields[i+1:]...)
			i = len(fields)
		}
	}
	return r, exact && r.Score.Kind != "" && len(r.PV) > 0
}
func (s *stockfish) analyse(position, move string, nodes int64) (sfResult, error) {
	// Clear both search state and hash between independent requests. Keep full history.
	for _, cmd := range []string{"ucinewgame", "setoption name Clear Hash"} {
		if err := s.send(cmd); err != nil {
			return sfResult{}, err
		}
	}
	if err := s.ready(); err != nil {
		return sfResult{}, err
	}
	if err := s.send(position); err != nil {
		return sfResult{}, err
	}
	command := fmt.Sprintf("go nodes %d", nodes)
	if move != "" {
		command += " searchmoves " + move
	}
	if err := s.send(command); err != nil {
		return sfResult{}, err
	}
	lines, err := s.until("bestmove", 2*time.Minute)
	if err != nil {
		return sfResult{}, err
	}
	var result sfResult
	for _, line := range lines {
		if r, ok := parseSFInfo(line); ok {
			result = r
		}
		if strings.HasPrefix(line, "bestmove ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				result.Bestmove = parts[1]
			}
		}
	}
	result.Output = lines
	if result.Score.Kind == "" || result.Bestmove == "" || len(result.PV) == 0 {
		return result, fmt.Errorf("Stockfish returned no exact score/PV")
	}
	if move == "" && result.PV[0] != result.Bestmove {
		rechecked, e := s.analyse(position, result.Bestmove, nodes)
		if e != nil {
			return result, e
		}
		rechecked.SelectionOutput = lines
		rechecked.ReferenceRechecked = true
		return rechecked, nil
	}
	if result.PV[0] != result.Bestmove {
		return result, fmt.Errorf("Stockfish score PV disagrees with selected move")
	}
	if move != "" && result.Bestmove != move {
		return result, fmt.Errorf("Stockfish ignored restricted root move")
	}
	return result, nil
}
