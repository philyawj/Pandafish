package main

import (
	"bufio"
	"fmt"
	"github.com/corentings/chess/v2"
	"io"
	"strings"
	"testing"
	"time"
)

type uciSession struct {
	t     *testing.T
	input *io.PipeWriter
	lines chan string
	done  chan error
}

func newUCISession(t *testing.T) *uciSession {
	input, commands := io.Pipe()
	responses, output := io.Pipe()
	session := &uciSession{t: t, input: commands, lines: make(chan string, 128), done: make(chan error, 1)}
	go func() { session.done <- run(input, output); output.Close() }()
	go func() {
		scanner := bufio.NewScanner(responses)
		for scanner.Scan() {
			session.lines <- scanner.Text()
		}
		close(session.lines)
	}()
	t.Cleanup(func() {
		commands.Close() // EOF cancels any unfinished search.
		select {
		case err := <-session.done:
			if err != nil {
				t.Errorf("UCI loop: %v", err)
			}
		case <-time.After(3 * time.Second):
			t.Error("UCI loop did not exit")
		}
		input.Close()
		responses.Close()
	})
	return session
}

func (session *uciSession) send(command string) {
	session.t.Helper()
	if _, err := fmt.Fprintln(session.input, command); err != nil {
		session.t.Fatal(err)
	}
}

func (session *uciSession) read() string {
	session.t.Helper()
	select {
	case line, ok := <-session.lines:
		if !ok {
			session.t.Fatal("unexpected UCI EOF")
		}
		return line
	case <-time.After(3 * time.Second):
		session.t.Fatal("timed out waiting for UCI response")
	}
	return ""
}

func TestUCIReplacementAndStops(t *testing.T) {
	session := newUCISession(t)
	session.send("go infinite searchmoves e2e4")
	session.send("isready")
	if got := session.read(); got != "readyok" {
		t.Fatal(got)
	}
	session.send("position startpos moves e2e4")
	session.send("go depth 1 searchmoves e7e5")
	if got := session.read(); got != "bestmove e7e5" {
		t.Fatal(got)
	}
	session.send("go infinite depth 1 searchmoves e7e5")
	session.send("stop")
	if got := session.read(); got != "bestmove e7e5" {
		t.Fatal(got)
	}
	session.send("stop")
	session.send("isready")
	if got := session.read(); got != "readyok" {
		t.Fatal(got)
	}
}

func TestUCIInvalidLimits(t *testing.T) {
	session := newUCISession(t)
	for _, command := range []string{"go movetime -1", "go depth", "go depth 0", "go wtime 1 btime invalid"} {
		session.send(command)
		if got := session.read(); !strings.HasPrefix(got, "info string invalid go:") {
			t.Fatal(got)
		}
	}
	session.send("go movetime 0 searchmoves e2e4")
	if got := session.read(); got != "bestmove e2e4" {
		t.Fatal(got)
	}
}

func TestUCIQuitDuringSearch(t *testing.T) {
	session := newUCISession(t)
	session.send("go infinite")
	session.send("quit")
	select {
	case line, ok := <-session.lines:
		if ok {
			t.Fatalf("unexpected output after quit: %s", line)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("quit did not close output")
	}
}

func TestUCIClockSequence(t *testing.T) {
	session := newUCISession(t)
	game := chess.NewGame()
	clocks := [2]time.Duration{time.Second, time.Second}
	var history []string
	for range 20 {
		if game.Outcome() != chess.NoOutcome {
			break
		}
		side := 0
		if game.Position().Turn() == chess.Black {
			side = 1
		}
		session.send("position startpos moves " + strings.Join(history, " "))
		started := time.Now()
		session.send(fmt.Sprintf("go wtime %d btime %d winc 10 binc 10", clocks[0].Milliseconds(), clocks[1].Milliseconds()))
		line := session.read()
		elapsed := time.Since(started)
		if elapsed >= clocks[side] {
			t.Fatalf("exhausted game clock: %v of %v", elapsed, clocks[side])
		}
		if !strings.HasPrefix(line, "bestmove ") {
			t.Fatal(line)
		}
		move := strings.TrimPrefix(line, "bestmove ")
		if err := game.PushNotationMove(move, chess.UCINotation{}, nil); err != nil {
			t.Fatalf("illegal live-clock move %s: %v", move, err)
		}
		history = append(history, move)
		clocks[side] += 10*time.Millisecond - elapsed
	}
}

func TestUCINewGameCancelsPendingSearch(t *testing.T) {
	session := newUCISession(t)
	session.send("go infinite searchmoves e2e4")
	session.send("ucinewgame")
	session.send("go depth 1 searchmoves d2d4")
	if got := session.read(); got != "bestmove d2d4" {
		t.Fatal(got)
	}
	session.send("go depth 64 movetime 0 searchmoves e2e4")
	if got := session.read(); got != "bestmove e2e4" {
		t.Fatal(got)
	}
	session.send("go depth 1 searchmoves")
	if got := session.read(); got != "bestmove 0000" {
		t.Fatal(got)
	}
}

func TestUCIHeldResultsAndResponsiveStop(t *testing.T) {
	for _, mode := range []string{"infinite", "ponder"} {
		t.Run(mode, func(t *testing.T) {
			session := newUCISession(t)
			session.send("go " + mode + " depth 1 searchmoves e2e4")
			time.Sleep(20 * time.Millisecond)
			session.send("isready")
			if got := session.read(); got != "readyok" {
				t.Fatalf("held result emitted early: %s", got)
			}
			if mode == "ponder" {
				session.send("ponderhit")
			} else {
				session.send("stop")
			}
			if got := session.read(); got != "bestmove e2e4" {
				t.Fatal(got)
			}
		})
	}
	session := newUCISession(t)
	session.send("go infinite")
	time.Sleep(20 * time.Millisecond)
	started := time.Now()
	session.send("isready")
	if got := session.read(); got != "readyok" {
		t.Fatal(got)
	}
	session.send("stop")
	if got := session.read(); !strings.HasPrefix(got, "bestmove ") {
		t.Fatal(got)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("busy engine did not respond promptly: %v", elapsed)
	}
}
