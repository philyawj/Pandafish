package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/philyawj/Pandafish/internal/telemetry"
	"go.yaml.in/yaml/v3"
)

type recordOptions struct {
	Client, Config, Engine, Session, Python string
	Hardware, Workload                      string
	PrepareOnly                             bool
}
type manifest struct {
	Hardware       string     `json:"hardware,omitempty"`
	Workload       string     `json:"workload,omitempty"`
	Schema         int        `json:"schema"`
	Started        time.Time  `json:"started"`
	Finished       *time.Time `json:"finished,omitempty"`
	Status         string     `json:"status"`
	EngineSHA256   string     `json:"engine_sha256"`
	EngineID       string     `json:"engine_id"`
	EngineRevision string     `json:"engine_revision,omitempty"`
	EngineDirty    string     `json:"engine_dirty,omitempty"`
	ConfigSHA256   string     `json:"config_sha256"`
	ClientRevision string     `json:"client_revision,omitempty"`
	Platform       string     `json:"platform"`
	ClientDirty    bool       `json:"client_dirty"`
	ExitCode       *int       `json:"exit_code,omitempty"`
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".json-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(0644); err == nil {
		_, err = tmp.Write(append(b, '\n'))
	}
	closeErr := tmp.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(name, path)
}
func absolute(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}
func section(config map[string]any, key string) map[string]any {
	m, _ := config[key].(map[string]any)
	return m
}
func str(m map[string]any, key string) string { s, _ := m[key].(string); return s }
func redactConfig(v any) any {
	switch x := v.(type) {
	case map[string]any:
		result := map[string]any{}
		for k, v := range x {
			lower := strings.ToLower(k)
			if sensitiveKey(lower) {
				result[k] = "REDACTED"
			} else {
				result[k] = redactConfig(v)
			}
		}
		return result
	case []any:
		result := make([]any, len(x))
		for i, v := range x {
			result[i] = redactConfig(v)
		}
		return result
	default:
		return v
	}
}

func copyEngine(source, destination string) (string, error) {
	input, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	_, err = io.Copy(io.MultiWriter(output, hash), input)
	closeErr := output.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
func engineIdentity(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader("uci\nisready\nquit\n")
	// Do not let an inherited capture directory contaminate another session.
	cmd.Env = replaceEnv(os.Environ(), telemetry.Environment, "")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("engine handshake failed")
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "id name ") && bytes.Contains(output, []byte("uciok")) {
			return strings.TrimPrefix(line, "id name "), nil
		}
	}
	return "", fmt.Errorf("engine did not supply a UCI identity")
}
func replaceEnv(env []string, key, value string) []string {
	result := make([]string, 0, len(env)+1)
	for _, v := range env {
		if !strings.HasPrefix(v, key+"=") {
			result = append(result, v)
		}
	}
	return append(result, key+"="+value)
}

func record(o recordOptions) (err error) {
	if o.Client == "" || o.Session == "" {
		return fmt.Errorf("record requires --client and --session")
	}
	client, err := filepath.Abs(o.Client)
	if err != nil {
		return err
	}
	session, err := filepath.Abs(o.Session)
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(absolute(client, o.Config))
	if err != nil {
		return fmt.Errorf("read client config: %w", err)
	}
	config := map[string]any{}
	if yaml.Unmarshal(raw, &config) != nil {
		return fmt.Errorf("invalid YAML config (contents omitted)")
	}
	engine := section(config, "engine")
	if engine == nil {
		return fmt.Errorf("config needs an engine section")
	}
	if str(engine, "protocol") != "uci" || str(engine, "interpreter") != "" || len(section(engine, "engine_options")) > 0 {
		return fmt.Errorf("record requires a native UCI engine without engine_options")
	}
	source := o.Engine
	if source == "" {
		source = absolute(client, filepath.Join(str(engine, "dir"), str(engine, "name")))
	}
	source, err = filepath.Abs(source)
	if err != nil {
		return err
	}
	originalToken := str(config, "token")
	token := os.Getenv("LICHESS_BOT_TOKEN")
	if token == "" {
		token = str(config, "token")
	}
	if token == "" {
		return fmt.Errorf("config or LICHESS_BOT_TOKEN must provide a token")
	}
	if o.Python == "" {
		o.Python = filepath.Join(client, "venv", "bin", "python")
	}
	o.Python = absolute(client, o.Python)
	if _, err = os.Stat(filepath.Join(client, "lichess-bot.py")); err != nil {
		return fmt.Errorf("client entry point: %w", err)
	}
	if _, err = os.Stat(o.Python); err != nil {
		return fmt.Errorf("client Python: %w", err)
	}
	if err = os.MkdirAll(filepath.Dir(session), 0755); err != nil {
		return err
	}
	if err = os.Mkdir(session, 0755); err != nil {
		return fmt.Errorf("session must be new: %w", err)
	}
	for _, dir := range []string{"games", "traces"} {
		if err = os.Mkdir(filepath.Join(session, dir), 0755); err != nil {
			return err
		}
	}
	m := manifest{Schema: 1, Started: time.Now().UTC(), Status: "preparing", Platform: runtime.GOOS + "/" + runtime.GOARCH, Hardware: o.Hardware, Workload: o.Workload}
	manifestPath := filepath.Join(session, "manifest.json")
	defer func() {
		if err != nil && m.Finished == nil {
			if m.Status != "start_failed" {
				m.Status = "preparation_failed"
			}
			_ = writeJSON(manifestPath, m)
		}
	}()
	binary := filepath.Join(session, "pandafish")
	m.EngineSHA256, err = copyEngine(source, binary)
	if err != nil {
		return err
	}
	m.EngineID, err = engineIdentity(binary)
	if err != nil {
		return err
	}
	if info, e := buildinfo.ReadFile(binary); e == nil {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				m.EngineRevision = s.Value
			}
			if s.Key == "vcs.modified" {
				m.EngineDirty = s.Value
			}
		}
	}
	revision := exec.Command("git", "-C", client, "rev-parse", "HEAD")
	if b, e := revision.Output(); e == nil {
		m.ClientRevision = strings.TrimSpace(string(b))
	}
	dirty := exec.Command("git", "-C", client, "status", "--porcelain")
	if b, e := dirty.Output(); e == nil {
		m.ClientDirty = len(bytes.TrimSpace(b)) > 0
	}
	engine["dir"], engine["name"] = session, "pandafish"
	if cwd := str(engine, "working_dir"); cwd != "" {
		engine["working_dir"] = absolute(client, cwd)
	}
	config["pgn_directory"], config["pgn_file_grouping"] = filepath.Join(session, "games"), "game"
	config["token"] = "provided-via-LICHESS_BOT_TOKEN"
	sanitized, err := yaml.Marshal(redactConfig(config))
	if err != nil {
		return fmt.Errorf("cannot encode config")
	}
	if err = os.WriteFile(filepath.Join(session, "config.redacted.yml"), sanitized, 0644); err != nil {
		return err
	}
	sum := sha256.Sum256(sanitized)
	m.ConfigSHA256 = hex.EncodeToString(sum[:])
	m.Status = "prepared"
	if err = writeJSON(manifestPath, m); err != nil {
		return err
	}
	if o.Hardware != "" || o.Workload != "" {
		if err = os.WriteFile(filepath.Join(session, "hardware.txt"), []byte("Hardware: "+o.Hardware+"\nWorkload: "+o.Workload+"\n"), 0644); err != nil {
			return err
		}
	}
	if o.PrepareOnly {
		fmt.Println("Prepared offline session:", session)
		return nil
	}
	// The temporary effective config is removed after shutdown, including failures.
	runtimeConfig, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot encode effective config")
	}
	configPath := filepath.Join(session, ".client-runtime.yml")
	if err = os.WriteFile(configPath, runtimeConfig, 0644); err != nil {
		return err
	}
	defer os.Remove(configPath)
	log, err := os.OpenFile(filepath.Join(session, "lichess-bot.log"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer log.Close()
	output := &sessionOutput{console: os.Stdout, log: log, secrets: append(secretValues(config), token, originalToken)}
	cmd := exec.Command(o.Python, filepath.Join(client, "lichess-bot.py"), "--config", configPath, "--disable_auto_logging")
	cmd.Dir = client
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Env = replaceEnv(replaceEnv(os.Environ(), "LICHESS_BOT_TOKEN", token), telemetry.Environment, filepath.Join(session, "traces"))
	cmd.Env = replaceEnv(cmd.Env, "PYTHONDONTWRITEBYTECODE", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	if err = cmd.Start(); err != nil {
		m.Status = "start_failed"
		return err
	}
	m.Status = "running"
	if err = writeJSON(manifestPath, m); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return err
	}
	fmt.Println("Recording session:", session, "(Ctrl+C requests client shutdown)")
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	for {
		select {
		case <-signals:
			_ = cmd.Process.Signal(os.Interrupt)
		case err = <-wait:
			flushErr := output.flush()
			now := time.Now().UTC()
			m.Finished = &now
			code := cmd.ProcessState.ExitCode()
			m.ExitCode = &code
			m.Status = "finished"
			if err != nil {
				m.Status = "client_failed"
			}
			if flushErr != nil {
				m.Status = "capture_failed"
			}
			if e := writeJSON(manifestPath, m); e != nil {
				return e
			}
			if reportErr := report(session); reportErr != nil {
				return fmt.Errorf("client exited; automatic report failed: %w", reportErr)
			}
			if flushErr != nil {
				return fmt.Errorf("write client log: %w", flushErr)
			}
			if err != nil {
				return fmt.Errorf("client exited with code %d; records preserved", code)
			}
			return nil
		}
	}
}

// Buffer complete lines so credentials split across process writes are redacted.
type sessionOutput struct {
	mu           sync.Mutex
	pending      []byte
	console, log io.Writer
	secrets      []string
}

func (w *sessionOutput) line(line []byte) error {
	text := string(line)
	for _, secret := range w.secrets {
		if secret != "" {
			text = strings.ReplaceAll(text, secret, "REDACTED")
		}
	}
	if _, err := io.WriteString(w.console, text); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w.log, "%s %s", time.Now().UTC().Format(time.RFC3339Nano), text)
	return err
}
func (w *sessionOutput) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.pending = append(w.pending, p...)
	for {
		i := bytes.IndexByte(w.pending, '\n')
		if i < 0 {
			break
		}
		if err := w.line(w.pending[:i+1]); err != nil {
			return 0, err
		}
		w.pending = w.pending[i+1:]
	}
	if len(w.pending) > 1024*1024 {
		return 0, fmt.Errorf("client output line exceeds 1 MiB")
	}
	return len(p), nil
}
func (w *sessionOutput) flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.pending) > 0 {
		if err := w.line(append(w.pending, '\n')); err != nil {
			return err
		}
		w.pending = nil
	}
	return nil
}

func sensitiveKey(key string) bool {
	key = strings.ToLower(key)
	for _, word := range []string{"token", "secret", "password", "api_key", "apikey", "authorization"} {
		if strings.Contains(key, word) {
			return true
		}
	}
	return false
}
func secretValues(v any) []string {
	var values []string
	switch x := v.(type) {
	case map[string]any:
		for key, value := range x {
			if sensitiveKey(key) {
				if text, ok := value.(string); ok && text != "" {
					values = append(values, text)
				}
			}
			values = append(values, secretValues(value)...)
		}
	case []any:
		for _, value := range x {
			values = append(values, secretValues(value)...)
		}
	}
	return values
}
