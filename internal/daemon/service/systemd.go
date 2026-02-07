package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	systemdServiceName = "ughd.service"
	systemdUnit        = "ughd"
)

// Systemd implements Manager for Linux systemd.
type Systemd struct {
	servicePath string // Computed path to service file
}

// NewSystemd creates a new systemd service manager.
func NewSystemd() *Systemd {
	return &Systemd{
		servicePath: systemdServicePath(),
	}
}

// systemdServicePath returns the path to the user service file.
func systemdServicePath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "systemd", "user", systemdServiceName)
}

// Name implements Manager.
func (s *Systemd) Name() string {
	return "systemd"
}

// isAvailable checks if systemctl is available.
func (s *Systemd) isAvailable() bool {
	_, err := exec.LookPath("systemctl")
	return err == nil
}

// Install implements Manager.
func (s *Systemd) Install(cfg InstallConfig) error {
	if s.servicePath == "" {
		return errors.New("cannot determine service path")
	}

	// Check if already installed
	if _, err := os.Stat(s.servicePath); err == nil {
		return ErrAlreadyInstalled
	}

	// Generate service content
	content := s.generateServiceFile(cfg)

	// Create directory if needed
	dir := filepath.Dir(s.servicePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create service dir: %w", err)
	}

	// Write service file
	if err := os.WriteFile(s.servicePath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write service file: %w", err)
	}

	// Reload systemd
	if err := s.systemctl("daemon-reload"); err != nil {
		return fmt.Errorf("daemon-reload: %w", err)
	}

	// Enable the service (so it starts on login)
	if err := s.systemctl("enable", systemdUnit); err != nil {
		return fmt.Errorf("enable service: %w", err)
	}

	return nil
}

// generateServiceFile creates the systemd unit file content.
func (s *Systemd) generateServiceFile(cfg InstallConfig) string {
	var sb strings.Builder
	sb.WriteString("[Unit]\n")
	sb.WriteString("Description=ugh daemon - background Turso sync\n")
	sb.WriteString("After=network.target\n")
	sb.WriteString("\n")
	sb.WriteString("[Service]\n")
	sb.WriteString("Type=simple\n")

	// Build ExecStart command
	execStart := cfg.BinaryPath + " daemon run"
	if cfg.ConfigPath != "" {
		execStart += " --config " + cfg.ConfigPath
	}
	sb.WriteString("ExecStart=" + execStart + "\n")

	sb.WriteString("Restart=on-failure\n")
	sb.WriteString("RestartSec=5\n")
	sb.WriteString("\n")
	sb.WriteString("[Install]\n")
	sb.WriteString("WantedBy=default.target\n")

	return sb.String()
}

// Uninstall implements Manager.
func (s *Systemd) Uninstall() error {
	if s.servicePath == "" {
		return errors.New("cannot determine service path")
	}

	// Check if installed
	if _, err := os.Stat(s.servicePath); os.IsNotExist(err) {
		return ErrNotInstalled
	}

	// Stop if running (ignore error if not running)
	_ = s.systemctl("stop", systemdUnit)

	// Disable the service (ignore errors - service might not be enabled)
	_ = s.systemctl("disable", systemdUnit)

	// Remove service file
	if err := os.Remove(s.servicePath); err != nil {
		return fmt.Errorf("remove service file: %w", err)
	}

	// Reload systemd
	if err := s.systemctl("daemon-reload"); err != nil {
		return fmt.Errorf("daemon-reload: %w", err)
	}

	return nil
}

// Start implements Manager.
func (s *Systemd) Start() error {
	// Check if installed
	if _, err := os.Stat(s.servicePath); os.IsNotExist(err) {
		return ErrNotInstalled
	}

	if err := s.systemctl("start", systemdUnit); err != nil {
		return fmt.Errorf("start service: %w", err)
	}
	return nil
}

// Stop implements Manager.
func (s *Systemd) Stop() error {
	status, err := s.Status()
	if err != nil {
		return err
	}
	if !status.Running {
		return ErrNotRunning
	}

	err = s.systemctl("stop", systemdUnit)
	if err != nil {
		return fmt.Errorf("stop service: %w", err)
	}
	return nil
}

// Status implements Manager.
func (s *Systemd) Status() (Status, error) {
	var status Status
	status.ServicePath = s.servicePath

	// Check if installed
	if _, err := os.Stat(s.servicePath); os.IsNotExist(err) {
		return status, nil
	}
	status.Installed = true

	// Check if running using is-active
	cmd := exec.CommandContext(context.Background(), "systemctl", "--user", "is-active", systemdUnit)
	output, _ := cmd.Output()
	state := strings.TrimSpace(string(output))
	status.Running = state == "active"

	// Get PID if running
	if status.Running {
		status.PID = s.getPID()
	}

	return status, nil
}

// getPID retrieves the main PID of the service.
func (s *Systemd) getPID() int {
	cmd := exec.CommandContext(
		context.Background(),
		"systemctl",
		"--user",
		"show",
		"-p",
		"MainPID",
		"--value",
		systemdUnit,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0
	}
	return pid
}

// LogPath implements Manager.
func (s *Systemd) LogPath() string {
	// systemd manages logs via journald, no file path
	return ""
}

// TailLogs implements Manager.
func (s *Systemd) TailLogs(ctx context.Context, follow bool, lines int, w io.Writer) error {
	args := []string{"--user", "-u", systemdUnit, "-n", strconv.Itoa(lines)}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.CommandContext(ctx, "journalctl", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start journalctl: %w", err)
	}

	// Wait for command to finish or context to be cancelled
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Context cancelled, process should be killed by CommandContext
		return ctx.Err()
	case err := <-errCh:
		if err != nil && ctx.Err() == nil {
			return fmt.Errorf("journalctl: %w", err)
		}
		return nil
	}
}

// systemctl runs systemctl with --user flag.
func (s *Systemd) systemctl(args ...string) error {
	fullArgs := append([]string{"--user"}, args...)
	cmd := exec.CommandContext(context.Background(), "systemctl", fullArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// GetUptime returns the uptime of the service in seconds, or 0 if not available.
// Note: Accurate uptime is provided by the daemon's health endpoint.
func (s *Systemd) GetUptime() int64 {
	// Uptime is better tracked by the daemon itself via the health endpoint.
	// This method is kept for potential future use.
	return 0
}
