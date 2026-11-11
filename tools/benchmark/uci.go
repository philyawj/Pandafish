package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Each game owns its processes. Closing one kills and reaps it, including after
// a timeout, so a broken engine cannot stall the rest of the benchmark.
type engine struct {
	cmd   *exec.Cmd
	in    io.WriteCloser
	lines chan string
	done  chan struct{}
	name  string
}

func startEngine(path string, options map[string]string) (*engine, error) {
	e := &engine{cmd: exec.Command(path), lines: make(chan string), done: make(chan struct{})}
	var err error
	e.in, err = e.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := e.cmd.StdoutPipe()
	if err != nil {
		e.in.Close()
		return nil, err
	}
	if err = e.cmd.Start(); err != nil {
		e.in.Close()
		out.Close()
		return nil, err
	}
	go func() {
		defer close(e.lines)
		scanner := bufio.NewScanner(out)
		scanner.Buffer(make([]byte, 4096), 1024*1024)
		for scanner.Scan() {
			select {
			case e.lines <- scanner.Text():
			case <-e.done:
				return
			}
		}
	}()
	fail := func(err error) (*engine, error) { e.close(); return nil, err }
	if err = e.send("uci"); err != nil {
		return fail(err)
	}
	// Read advertised options as well as identity before configuring the engine.
	advertised := map[string]bool{}
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
handshake:
	for {
		select {
		case line, ok := <-e.lines:
			if !ok {
				return fail(fmt.Errorf("engine exited during handshake"))
			}
			if line == "uciok" {
				break handshake
			}
			if strings.HasPrefix(line, "id name ") {
				e.name = strings.TrimPrefix(line, "id name ")
			}
			if strings.HasPrefix(line, "option name ") {
				name, _, _ := strings.Cut(strings.TrimPrefix(line, "option name "), " type ")
				advertised[name] = true
			}
		case <-timer.C:
			return fail(fmt.Errorf("UCI handshake timed out"))
		}
	}
	keys := make([]string, 0, len(options))
	for key := range options {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !advertised[key] {
			return fail(fmt.Errorf("engine does not advertise %q", key))
		}
		if err = e.send("setoption name " + key + " value " + options[key]); err != nil {
			return fail(err)
		}
	}
	for _, command := range []string{"ucinewgame", "isready"} {
		if err = e.send(command); err != nil {
			return fail(err)
		}
	}
	if _, err = e.wait("readyok", 5*time.Second); err != nil {
		return fail(err)
	}
	return e, nil
}

func (e *engine) send(command string) error {
	_, err := fmt.Fprintln(e.in, command)
	return err
}

func (e *engine) wait(prefix string, timeout time.Duration) (string, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case line, ok := <-e.lines:
			if !ok {
				return "", fmt.Errorf("engine exited or output could not be read")
			}
			if line == prefix || strings.HasPrefix(line, prefix+" ") {
				return line, nil
			}
		case <-timer.C:
			return "", fmt.Errorf("response timeout")
		}
	}
}

func (e *engine) move(history []string, allowance time.Duration, grace time.Duration) (string, time.Duration, error) {
	started := time.Now()
	limit := allowance + grace
	type response struct {
		move    string
		elapsed time.Duration
		err     error
	}
	responses := make(chan response, 1)
	// The deadline includes pipe writes and position processing, not just the
	// read for bestmove. Closing the engine unblocks a stalled write as well.
	go func() {
		var line string
		err := e.send("position startpos moves " + strings.Join(history, " "))
		if err == nil {
			err = e.send(fmt.Sprintf("go movetime %d", allowance.Milliseconds()))
		}
		if err == nil {
			line, err = e.wait("bestmove", time.Until(started.Add(limit)))
		}
		move := ""
		if err == nil {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				err = fmt.Errorf("bestmove is missing a move")
			} else {
				move = fields[1]
			}
		}
		responses <- response{move, time.Since(started), err}
	}()
	timer := time.NewTimer(limit)
	defer timer.Stop()
	var r response
	select {
	case r = <-responses:
	case <-timer.C:
		// A response received before the deadline remains valid even if this
		// goroutine was scheduled late.
		select {
		case r = <-responses:
		default:
			return "", time.Since(started), fmt.Errorf("move response timeout")
		}
	}
	if r.elapsed > limit {
		return "", r.elapsed, fmt.Errorf("move response timeout")
	}
	return r.move, r.elapsed, r.err
}

func (e *engine) close() {
	close(e.done)
	e.in.Close()
	_ = e.cmd.Process.Kill()
	_ = e.cmd.Wait()
}
