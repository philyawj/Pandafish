package telemetry

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFailedRecordingDoesNotSuppressProtocol(t *testing.T) {
	recorder, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err = recorder.Close(); err != nil {
		t.Fatal(err)
	}
	var destination bytes.Buffer
	output := Output{Destination: &destination, Recorder: recorder}
	line := []byte("bestmove e2e4\n")
	n, err := output.Write(line)
	if err != nil || n != len(line) || !bytes.Equal(destination.Bytes(), line) {
		t.Fatal("recording failure suppressed output")
	}
	if recorder.Write("exit", "") == nil {
		t.Fatal("lost recording error")
	}
}
func TestRecorderPermissionsAndUniqueProcesses(t *testing.T) {
	dir := t.TempDir()
	a, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	files, err := filepath.Glob(filepath.Join(dir, "engine-*.jsonl"))
	if err != nil || len(files) != 2 {
		t.Fatal(files, err)
	}
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0644 {
			t.Fatalf("mode %o", info.Mode().Perm())
		}
	}
}
