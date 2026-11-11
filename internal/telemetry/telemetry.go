// Package telemetry defines the local recording format shared by engine and tools.
package telemetry

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const Prefix = "info string pandafish "
const Environment = "PANDAFISH_TRACE_DIR"

// Metrics describes the last completed iteration, plus work in interrupted ones.
// Scores are from the root side's perspective. Mate values give a sign, not distance.
type Metrics struct {
	Depth     int      `json:"depth"`
	Nodes     uint64   `json:"nodes"`
	ElapsedMS float64  `json:"elapsed_ms"`
	Timed     bool     `json:"timed"`
	SoftMS    float64  `json:"soft_ms"`
	HardMS    float64  `json:"hard_ms"`
	Reason    string   `json:"reason"`
	Fallback  bool     `json:"fallback"`
	ScoreKind string   `json:"score_kind"`
	Score     int      `json:"score"`
	PV        []string `json:"pv,omitempty"` // Only the root move is currently retained.
}

type Event struct {
	Schema    int       `json:"schema"`
	Sequence  uint64    `json:"sequence"`
	Time      time.Time `json:"time"`
	ElapsedNS int64     `json:"elapsed_ns"`
	Kind      string    `json:"kind"`
	Line      string    `json:"line,omitempty"`
}

type Recorder struct {
	file     *os.File
	encoder  *json.Encoder
	start    time.Time
	sequence uint64
	err      error
}

func Open(dir string) (*Recorder, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.CreateTemp(dir, fmt.Sprintf("engine-%d-*.jsonl", os.Getpid()))
	if err != nil {
		return nil, err
	}
	if err := f.Chmod(0644); err != nil {
		f.Close()
		return nil, err
	}
	r := &Recorder{file: f, encoder: json.NewEncoder(f), start: time.Now()}
	if err := r.Write("start", filepath.Base(f.Name())); err != nil {
		f.Close()
		return nil, err
	}
	return r, nil
}

// Write is called only by the UCI command loop, never by the search worker.
func (r *Recorder) Write(kind, line string) error {
	if r == nil {
		return nil
	}
	if r.err != nil {
		return r.err
	}
	r.sequence++
	r.err = r.encoder.Encode(Event{Schema: 1, Sequence: r.sequence, Time: time.Now().UTC(), ElapsedNS: time.Since(r.start).Nanoseconds(), Kind: kind, Line: line})
	return r.err
}

func (r *Recorder) Close() error {
	if r == nil {
		return nil
	}
	return r.file.Close()
}

// Output records UCI writes. A failed recording must not suppress a legal move.
type Output struct {
	Destination io.Writer
	Recorder    *Recorder
}

func (o Output) Write(p []byte) (int, error) {
	_ = o.Recorder.Write("out", string(p))
	return o.Destination.Write(p)
}
