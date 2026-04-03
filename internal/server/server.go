package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// Config holds llama-server startup parameters.
type Config struct {
	Binary     string
	ModelPath  string
	Host       string
	Port       int
	CtxSize    int
	Threads    int
	BatchSize  int
	NGPULayers int
}

// Info describes a running server process.
type Info struct {
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	ModelPath string    `json:"model_path"`
	StartedAt time.Time `json:"started_at"`
}

// Manager manages the llama-server process lifecycle.
type Manager struct {
	Config  Config
	DataDir string // ~/.bonsai
}

// NewManager creates a manager with defaults from environment variables.
func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".bonsai")

	port := 8081
	if p := os.Getenv("BONSAI_PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			port = n
		}
	}

	threads := runtime.NumCPU()
	if t := os.Getenv("BONSAI_THREADS"); t != "" {
		if n, err := strconv.Atoi(t); err == nil {
			threads = n
		}
	}

	return &Manager{
		Config: Config{
			Host:       "127.0.0.1",
			Port:       port,
			CtxSize:    4096,
			Threads:    threads,
			BatchSize:  512,
			NGPULayers: 99,
		},
		DataDir: dataDir,
	}
}

func (m *Manager) pidFile() string {
	return filepath.Join(m.DataDir, "server.pid")
}

func (m *Manager) logFile() string {
	return filepath.Join(m.DataDir, "server.log")
}

// BaseURL returns the server URL.
func (m *Manager) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", m.Config.Host, m.Config.Port)
}

// readInfo reads the PID file. Returns nil if not found.
func (m *Manager) readInfo() *Info {
	data, err := os.ReadFile(m.pidFile())
	if err != nil {
		return nil
	}
	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return nil
	}
	return &info
}

func (m *Manager) writeInfo(info *Info) error {
	if err := os.MkdirAll(m.DataDir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return os.WriteFile(m.pidFile(), data, 0644)
}

// isProcessAlive checks if a process with the given PID is running.
func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

// IsRunning checks if the server process is alive.
func (m *Manager) IsRunning() bool {
	info := m.readInfo()
	if info == nil {
		return false
	}
	return isProcessAlive(info.PID)
}

// ProcessInfo returns information about the running server.
func (m *Manager) ProcessInfo() (*Info, error) {
	info := m.readInfo()
	if info == nil {
		return nil, fmt.Errorf("no server info found")
	}
	if !isProcessAlive(info.PID) {
		return nil, fmt.Errorf("server process (PID %d) is not running", info.PID)
	}
	return info, nil
}

// LoadedModel returns the model path of the running server.
func (m *Manager) LoadedModel() string {
	info := m.readInfo()
	if info == nil {
		return ""
	}
	return info.ModelPath
}

// Start launches llama-server in the background.
func (m *Manager) Start() error {
	if m.IsRunning() {
		return fmt.Errorf("server already running (PID %d)", m.readInfo().PID)
	}

	binary := m.Config.Binary
	if binary == "" {
		var err error
		binary, err = FindBinary()
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(m.DataDir, 0755); err != nil {
		return err
	}

	logF, err := os.Create(m.logFile())
	if err != nil {
		return fmt.Errorf("cannot create log file: %w", err)
	}

	cmd := exec.Command(binary,
		"-m", m.Config.ModelPath,
		"--host", m.Config.Host,
		"--port", strconv.Itoa(m.Config.Port),
		"--ctx-size", strconv.Itoa(m.Config.CtxSize),
		"--threads", strconv.Itoa(m.Config.Threads),
		"--batch-size", strconv.Itoa(m.Config.BatchSize),
		"-ngl", strconv.Itoa(m.Config.NGPULayers),
	)
	cmd.Stdout = logF
	cmd.Stderr = logF
	// Detach from parent process group so server survives CLI exit
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		logF.Close()
		return fmt.Errorf("failed to start llama-server: %w", err)
	}

	// Close log file handle in parent (child owns it now via inheritance)
	logF.Close()

	info := &Info{
		PID:       cmd.Process.Pid,
		Port:      m.Config.Port,
		ModelPath: m.Config.ModelPath,
		StartedAt: time.Now(),
	}
	if err := m.writeInfo(info); err != nil {
		return fmt.Errorf("server started but failed to write PID file: %w", err)
	}

	// Release the process so it doesn't become a zombie
	cmd.Process.Release()

	return nil
}

// StartForeground runs llama-server attached to the current terminal.
func (m *Manager) StartForeground() error {
	binary := m.Config.Binary
	if binary == "" {
		var err error
		binary, err = FindBinary()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(binary,
		"-m", m.Config.ModelPath,
		"--host", m.Config.Host,
		"--port", strconv.Itoa(m.Config.Port),
		"--ctx-size", strconv.Itoa(m.Config.CtxSize),
		"--threads", strconv.Itoa(m.Config.Threads),
		"--batch-size", strconv.Itoa(m.Config.BatchSize),
		"-ngl", strconv.Itoa(m.Config.NGPULayers),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// Stop kills the running llama-server process.
func (m *Manager) Stop() error {
	info := m.readInfo()
	if info == nil {
		return fmt.Errorf("no server PID file found")
	}

	proc, err := os.FindProcess(info.PID)
	if err != nil {
		os.Remove(m.pidFile())
		return nil
	}

	// Try graceful shutdown first
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		os.Remove(m.pidFile())
		return nil // process already dead
	}

	// Wait up to 5 seconds for graceful shutdown
	done := make(chan struct{})
	go func() {
		for i := 0; i < 50; i++ {
			if !isProcessAlive(info.PID) {
				close(done)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		// Force kill
		proc.Signal(syscall.SIGKILL)
		close(done)
	}()
	<-done

	os.Remove(m.pidFile())
	return nil
}

// WaitHealthy polls the server health endpoint until it returns 200 or timeout.
func (m *Manager) WaitHealthy(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	url := m.BaseURL() + "/health"

	for time.Now().Before(deadline) {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		// Check if the process died
		if !m.IsRunning() {
			return fmt.Errorf("server process died during startup (check logs: %s)", m.logFile())
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("server did not become healthy within %s (check logs: %s)", timeout, m.logFile())
}

// EnsureRunning ensures the server is running with the given model.
// If already running with the same model, returns immediately.
// If running with a different model, restarts. If not running, starts.
func (m *Manager) EnsureRunning(modelPath string) (string, error) {
	m.Config.ModelPath = modelPath

	info := m.readInfo()
	if info != nil && isProcessAlive(info.PID) {
		// Server is running
		if info.ModelPath == modelPath {
			// Same model — just return the URL
			m.Config.Port = info.Port
			return m.BaseURL(), nil
		}
		// Different model — restart
		fmt.Fprintf(os.Stderr, "Restarting server with new model...\n")
		m.Stop()
	}

	fmt.Fprintf(os.Stderr, "Starting llama-server...\n")
	if err := m.Start(); err != nil {
		return "", err
	}

	if err := m.WaitHealthy(60 * time.Second); err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "Server ready.\n")
	return m.BaseURL(), nil
}
