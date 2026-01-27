package service

import (
	"bufio"
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
	launchdLabel   = "com.ugh.daemon"
	launchdPlist   = "com.ugh.daemon.plist"
	defaultLogFile = "ughd.log"
)

// Launchd implements Manager for macOS launchd.
type Launchd struct {
	plistPath string // Path to plist file
	logDir    string // Directory for log files
}

// NewLaunchd creates a new launchd service manager.
func NewLaunchd() *Launchd {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Launchd{}
	}
	return &Launchd{
		plistPath: filepath.Join(home, "Library", "LaunchAgents", launchdPlist),
		logDir:    filepath.Join(home, "Library", "Logs", "ugh"),
	}
}

// Name implements Manager.
func (l *Launchd) Name() string {
	return "launchd"
}

// Install implements Manager.
func (l *Launchd) Install(cfg InstallConfig) error {
	if l.plistPath == "" {
		return errors.New("cannot determine plist path")
	}

	// Check if already installed
	if _, err := os.Stat(l.plistPath); err == nil {
		return ErrAlreadyInstalled
	}

	// Create log directory
	if err := os.MkdirAll(l.logDir, 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	// Generate plist content
	content := l.generatePlist(cfg)

	// Create LaunchAgents directory if needed
	dir := filepath.Dir(l.plistPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create LaunchAgents dir: %w", err)
	}

	// Write plist file
	if err := os.WriteFile(l.plistPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write plist file: %w", err)
	}

	// Load the service (but don't start it)
	if err := l.launchctl("load", l.plistPath); err != nil {
		return fmt.Errorf("load plist: %w", err)
	}

	return nil
}

// generatePlist creates the launchd plist content.
func (l *Launchd) generatePlist(cfg InstallConfig) string {
	logPath := filepath.Join(l.logDir, defaultLogFile)
	errPath := filepath.Join(l.logDir, "ughd.err.log")

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>` + launchdLabel + `</string>
    <key>ProgramArguments</key>
    <array>
        <string>` + cfg.BinaryPath + `</string>
        <string>daemon</string>
        <string>run</string>
`)
	if cfg.ConfigPath != "" {
		sb.WriteString(`        <string>--config</string>
        <string>` + cfg.ConfigPath + `</string>
`)
	}
	sb.WriteString(`    </array>
    <key>RunAtLoad</key>
    <false/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>StandardOutPath</key>
    <string>` + logPath + `</string>
    <key>StandardErrorPath</key>
    <string>` + errPath + `</string>
    <key>ProcessType</key>
    <string>Background</string>
</dict>
</plist>
`)
	return sb.String()
}

// Uninstall implements Manager.
func (l *Launchd) Uninstall() error {
	if l.plistPath == "" {
		return errors.New("cannot determine plist path")
	}

	// Check if installed
	if _, err := os.Stat(l.plistPath); os.IsNotExist(err) {
		return ErrNotInstalled
	}

	// Stop if running (ignore error if not running)
	_ = l.launchctl("stop", launchdLabel)

	// Unload the service (ignore errors - service might not be loaded)
	_ = l.launchctl("unload", l.plistPath)

	// Remove plist file
	if err := os.Remove(l.plistPath); err != nil {
		return fmt.Errorf("remove plist file: %w", err)
	}

	return nil
}

// Start implements Manager.
func (l *Launchd) Start() error {
	// Check if installed
	if _, err := os.Stat(l.plistPath); os.IsNotExist(err) {
		return ErrNotInstalled
	}

	if err := l.launchctl("start", launchdLabel); err != nil {
		return fmt.Errorf("start service: %w", err)
	}
	return nil
}

// Stop implements Manager.
func (l *Launchd) Stop() error {
	status, err := l.Status()
	if err != nil {
		return err
	}
	if !status.Running {
		return ErrNotRunning
	}

	if err := l.launchctl("stop", launchdLabel); err != nil {
		return fmt.Errorf("stop service: %w", err)
	}
	return nil
}

// Status implements Manager.
func (l *Launchd) Status() (Status, error) {
	var status Status
	status.ServicePath = l.plistPath

	// Check if installed
	if _, err := os.Stat(l.plistPath); os.IsNotExist(err) {
		return status, nil
	}
	status.Installed = true

	// Check if running using launchctl list
	cmd := exec.Command("launchctl", "list", launchdLabel)
	output, err := cmd.Output()
	if err != nil {
		// If list fails, service is not loaded/running
		return status, nil
	}

	// Parse output - format is "PID\tStatus\tLabel"
	// or just "Status\tLabel" if not running
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			// First field might be PID or "-"
			pid, err := strconv.Atoi(fields[0])
			if err == nil && pid > 0 {
				status.Running = true
				status.PID = pid
				break
			}
		}
	}

	// Alternative check using launchctl print
	if !status.Running {
		status.Running, status.PID = l.checkRunningWithPrint()
	}

	return status, nil
}

// checkRunningWithPrint uses launchctl print to check if service is running.
func (l *Launchd) checkRunningWithPrint() (bool, int) {
	cmd := exec.Command("launchctl", "print", "gui/"+strconv.Itoa(os.Getuid())+"/"+launchdLabel)
	output, err := cmd.Output()
	if err != nil {
		return false, 0
	}

	// Parse output for PID
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "pid = ") {
			pidStr := strings.TrimPrefix(line, "pid = ")
			pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
			if err == nil && pid > 0 {
				return true, pid
			}
		}
	}

	return false, 0
}

// LogPath implements Manager.
func (l *Launchd) LogPath() string {
	return filepath.Join(l.logDir, defaultLogFile)
}

// TailLogs implements Manager.
func (l *Launchd) TailLogs(ctx context.Context, follow bool, lines int, w io.Writer) error {
	logPath := l.LogPath()

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("log file not found: %s", logPath)
	}

	var args []string
	if follow {
		args = []string{"-n", strconv.Itoa(lines), "-f", logPath}
	} else {
		args = []string{"-n", strconv.Itoa(lines), logPath}
	}

	cmd := exec.CommandContext(ctx, "tail", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start tail: %w", err)
	}

	// Wait for command to finish or context to be cancelled
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil && ctx.Err() == nil {
			return fmt.Errorf("tail: %w", err)
		}
		return nil
	}
}

// launchctl runs launchctl with the given arguments.
func (l *Launchd) launchctl(args ...string) error {
	cmd := exec.Command("launchctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// ReadLogFile reads the last n lines from the log file.
func (l *Launchd) ReadLogFile(n int) ([]string, error) {
	logPath := l.LogPath()
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	return lines, scanner.Err()
}
