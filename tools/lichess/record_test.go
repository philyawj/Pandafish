package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"go.yaml.in/yaml/v3"
)

// The child is a local stand-in for Python, never a network client.
func TestMain(m *testing.M) {
	if os.Getenv("PANDAFISH_TEST_CLIENT") == "1" {
		mode := os.Getenv("PANDAFISH_TEST_MODE")
		if mode == "failure" {
			os.Exit(7)
		}
		if mode == "signal" {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			parent, _ := os.FindProcess(os.Getppid())
			_ = parent.Signal(syscall.SIGTERM)
			select {
			case <-ch:
				fmt.Println("client received graceful interrupt")
			case <-time.After(5 * time.Second):
				os.Exit(8)
			}
		}
		secret := os.Getenv("LICHESS_BOT_TOKEN")
		fmt.Print(secret[:len(secret)/2])
		fmt.Print(secret[len(secret)/2:] + "\n")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
func recordingFixture(t *testing.T) recordOptions {
	t.Helper()
	root := t.TempDir()
	client := filepath.Join(root, "client")
	if err := os.Mkdir(client, 0755); err != nil {
		t.Fatal(err)
	}
	engine := filepath.Join(root, "engine-original")
	if err := os.WriteFile(engine, []byte("#!/bin/sh\ncat >/dev/null\nprintf 'id name Pandafish test\\nuciok\\nreadyok\\n'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(client, "lichess-bot.py"), []byte("offline fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	config := "token: test-token-never-publish\nurl: https://lichess.org\nengine:\n  dir: ../\n  name: engine-original\n  protocol: uci\n  ponder: false\nchallenge:\n  modes: [casual]\n  accept_bot: false\nmatchmaking:\n  allow_matchmaking: false\n"
	if err := os.WriteFile(filepath.Join(client, "config.yml"), []byte(config), 0644); err != nil {
		t.Fatal(err)
	}
	python, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("LICHESS_BOT_TOKEN", "")
	return recordOptions{Client: client, Config: "config.yml", Python: python, Session: filepath.Join(root, "session"), PrepareOnly: true}
}
func readManifest(t *testing.T, dir string) manifest {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m manifest
	if err = json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	return m
}
func TestPrepareSnapshotPreservesPolicyAndExcludesToken(t *testing.T) {
	o := recordingFixture(t)
	before, err := os.ReadFile(filepath.Join(o.Client, o.Config))
	if err != nil {
		t.Fatal(err)
	}
	if err := record(o); err != nil {
		t.Fatal(err)
	}
	m := readManifest(t, o.Session)
	if m.Status != "prepared" || m.EngineID != "Pandafish test" || len(m.EngineSHA256) != 64 || len(m.ConfigSHA256) != 64 {
		t.Fatalf("manifest %+v", m)
	}
	after, err := os.ReadFile(filepath.Join(o.Client, o.Config))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("original config changed")
	}
	raw, err := os.ReadFile(filepath.Join(o.Session, "config.redacted.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte("test-token-never-publish")) {
		t.Fatal("token in snapshot")
	}
	var config map[string]any
	if err = yaml.Unmarshal(raw, &config); err != nil {
		t.Fatal(err)
	}
	engine := section(config, "engine")
	if engine["dir"] != o.Session || engine["name"] != "pandafish" || engine["ponder"] != false || section(config, "matchmaking")["allow_matchmaking"] != false || section(config, "challenge")["accept_bot"] != false {
		t.Fatal("capture changed bot policy")
	}
	info, err := os.Stat(filepath.Join(o.Session, "config.redacted.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("mode %o", info.Mode().Perm())
	}
	if err = record(o); err == nil {
		t.Fatal("overwrote existing session")
	}
	// A rebuild of the source executable cannot change the copied session engine.
	if err = os.WriteFile(filepath.Join(filepath.Dir(o.Client), "engine-original"), []byte("changed"), 0755); err != nil {
		t.Fatal(err)
	}
	snapshot, err := os.ReadFile(filepath.Join(o.Session, "pandafish"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(snapshot, []byte("uciok")) {
		t.Fatal("session snapshot changed")
	}
}
func TestRecordedClientLifecycle(t *testing.T) {
	for _, mode := range []string{"success", "signal", "failure"} {
		t.Run(mode, func(t *testing.T) {
			o := recordingFixture(t)
			o.PrepareOnly = false
			t.Setenv("PANDAFISH_TEST_CLIENT", "1")
			t.Setenv("PANDAFISH_TEST_MODE", mode)
			err := record(o)
			m := readManifest(t, o.Session)
			if mode == "failure" {
				if err == nil || m.Status != "client_failed" || m.ExitCode == nil || *m.ExitCode != 7 {
					t.Fatalf("failure: %v %+v", err, m)
				}
			} else {
				if err != nil || m.Status != "finished" || m.Finished == nil {
					t.Fatalf("success: %v %+v", err, m)
				}
			}
			if _, err = os.Stat(filepath.Join(o.Session, ".client-runtime.yml")); !os.IsNotExist(err) {
				t.Fatal("temporary config retained")
			}
			log, err := os.ReadFile(filepath.Join(o.Session, "lichess-bot.log"))
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(log, []byte("test-token-never-publish")) {
				t.Fatal("token leaked into log")
			}
			if mode == "signal" && !bytes.Contains(log, []byte("graceful interrupt")) {
				t.Fatal("interrupt not forwarded")
			}
		})
	}
}
func TestRedactionAcrossWritesAndNestedConfig(t *testing.T) {
	var console, log bytes.Buffer
	w := sessionOutput{console: &console, log: &log, secrets: []string{"SECRET"}}
	for _, text := range []string{"token=SE", "CRET\nlast=SEC", "RET"} {
		if _, err := w.Write([]byte(text)); err != nil {
			t.Fatal(err)
		}
	}
	w.flush()
	if strings.Contains(console.String(), "SECRET") || strings.Contains(log.String(), "SECRET") || console.String() != "token=REDACTED\nlast=REDACTED\n" {
		t.Fatalf("redaction %q %q", console.String(), log.String())
	}
	raw, err := yaml.Marshal(redactConfig(map[string]any{"nested": map[string]any{"api_key": "secret", "password": "secret", "normal": "kept"}}))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "secret") || !strings.Contains(string(raw), "kept") {
		t.Fatal("nested redaction failed")
	}
}
func TestInvalidCommandAndPreparationInputs(t *testing.T) {
	for _, args := range [][]string{nil, {"other"}, {"record"}, {"report"}, {"report", "--bogus"}} {
		if execute(args) == nil {
			t.Fatalf("accepted %v", args)
		}
	}
	o := recordingFixture(t)
	if err := os.WriteFile(filepath.Join(o.Client, o.Config), []byte("token: [\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := record(o); err == nil || strings.Contains(err.Error(), "token:") {
		t.Fatal("invalid config was accepted or exposed")
	}
}
