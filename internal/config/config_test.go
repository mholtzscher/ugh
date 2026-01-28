package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPath(t *testing.T) {
	path, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath() error = %v", err)
	}
	if path == "" {
		t.Fatal("DefaultPath() returned empty string")
	}
	if filepath.IsAbs(path) != true {
		t.Fatalf("DefaultPath() returned relative path: %s", path)
	}
}

func TestLoad_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	missingPath := filepath.Join(tmpDir, "does-not-exist.toml")

	_, err := Load(missingPath, false)
	if err == nil {
		t.Fatal("Load() should return error for missing file")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load() wrong error: %v", err)
	}
}

func TestLoad_AllowMissing(t *testing.T) {
	tmpDir := t.TempDir()
	missingPath := filepath.Join(tmpDir, "does-not-exist.toml")

	result, err := Load(missingPath, true)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if result.WasLoaded {
		t.Fatal("Load() should report WasLoaded = false")
	}
	if result.Config.Version != DefaultVersion {
		t.Fatalf("Load() version = %d, want %d", result.Config.Version, DefaultVersion)
	}
}

func TestLoad_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.toml")
	cfgContent := `version = 1

[db]
path = "~/.local/share/ugh/ugh.sqlite"
`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	result, err := Load(cfgPath, false)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !result.WasLoaded {
		t.Fatal("Load() should report WasLoaded = true")
	}
	if result.Config.DB.Path == "" {
		t.Fatal("Load() DB.Path is empty")
	}
	if result.UsedPath != cfgPath {
		t.Fatalf("Load() UsedPath = %s, want %s", result.UsedPath, cfgPath)
	}
}

func TestLoad_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "invalid.toml")
	invalidContent := `[invalid section
something = "value"
`
	if err := os.WriteFile(cfgPath, []byte(invalidContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(cfgPath, false)
	if err == nil {
		t.Fatal("Load() should return error for invalid TOML")
	}
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("Load() wrong error: %v", err)
	}
}

func TestResolveDBPath_ExpandEnv(t *testing.T) {
	origHome, _ := os.UserHomeDir()
	_ = os.Setenv("HOME", "/test/home")
	defer func() {
		if origHome != "" {
			_ = os.Setenv("HOME", origHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	path, err := ResolveDBPath("", "$HOME/data/ugh.sqlite")
	if err != nil {
		t.Fatalf("ResolveDBPath() error = %v", err)
	}
	if path != "/test/home/data/ugh.sqlite" {
		t.Fatalf("ResolveDBPath() = %s, want /test/home/data/ugh.sqlite", path)
	}
}

func TestResolveDBPath_ExpandHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get home dir: %v", err)
	}

	path, err := ResolveDBPath("", "~/data/ugh.sqlite")
	if err != nil {
		t.Fatalf("ResolveDBPath() error = %v", err)
	}
	expected := filepath.Join(homeDir, "data/ugh.sqlite")
	if path != expected {
		t.Fatalf("ResolveDBPath() = %s, want %s", path, expected)
	}
}

func TestResolveDBPath_ExpandHomeAlone(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get home dir: %v", err)
	}

	path, err := ResolveDBPath("", "~")
	if err != nil {
		t.Fatalf("ResolveDBPath() error = %v", err)
	}
	if path != homeDir {
		t.Fatalf("ResolveDBPath() = %s, want %s", path, homeDir)
	}
}

func TestResolveDBPath_RelativeToConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config", "ugh-config.toml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	path, err := ResolveDBPath(cfgPath, "data/ugh.sqlite")
	if err != nil {
		t.Fatalf("ResolveDBPath() error = %v", err)
	}
	expected := filepath.Join(tmpDir, "config", "data", "ugh.sqlite")
	if path != expected {
		t.Fatalf("ResolveDBPath() = %s, want %s", path, expected)
	}
}

func TestResolveDBPath_RelativeNoConfig(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	if err := os.WriteFile("ugh.sqlite", []byte{}, 0o644); err != nil {
		t.Fatalf("create sqlite file: %v", err)
	}

	path, err := ResolveDBPath("", "ugh.sqlite")
	if err != nil {
		t.Fatalf("ResolveDBPath() error = %v", err)
	}

	expected := filepath.Join(tmpDir, "ugh.sqlite")
	expectedInfo, err := os.Stat(expected)
	if err != nil {
		t.Fatalf("stat expected sqlite file: %v", err)
	}
	gotInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat resolved sqlite file: %v", err)
	}
	if !os.SameFile(expectedInfo, gotInfo) {
		t.Fatalf("ResolveDBPath() = %s, want same file as %s", path, expected)
	}
}

func TestResolveDBPath_ErrorEmpty(t *testing.T) {
	_, err := ResolveDBPath("", "")
	if err == nil {
		t.Fatal("ResolveDBPath() should return error for empty path")
	}
}

func TestUserConfigDir(t *testing.T) {
	dir, err := userConfigDir()
	if err != nil {
		t.Fatalf("userConfigDir() error = %v", err)
	}
	if dir == "" {
		t.Fatal("userConfigDir() returned empty string")
	}
	if filepath.IsAbs(dir) != true {
		t.Fatalf("userConfigDir() returned relative path: %s", dir)
	}
}

func TestLoad_EnvVarOverride(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.toml")
	cfgContent := `version = 1

[db]
path = "~/.local/share/ugh/ugh.sqlite"
auth_token = "file-token"
`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Set env var to override auth_token
	_ = os.Setenv("UGH_DB_AUTH_TOKEN", "env-token")
	defer func() { _ = os.Unsetenv("UGH_DB_AUTH_TOKEN") }()

	result, err := Load(cfgPath, false)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if result.Config.DB.AuthToken != "env-token" {
		t.Fatalf("Load() DB.AuthToken = %s, want env-token (env var should override file)", result.Config.DB.AuthToken)
	}
}

func TestLoad_EnvVarWithNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	missingPath := filepath.Join(tmpDir, "does-not-exist.toml")

	// Set env var
	_ = os.Setenv("UGH_DB_AUTH_TOKEN", "env-only-token")
	_ = os.Setenv("UGH_DB_PATH", "/some/path.db")
	defer func() {
		_ = os.Unsetenv("UGH_DB_AUTH_TOKEN")
		_ = os.Unsetenv("UGH_DB_PATH")
	}()

	result, err := Load(missingPath, true) // allowMissing=true
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if result.Config.DB.AuthToken != "env-only-token" {
		t.Fatalf("Load() DB.AuthToken = %s, want env-only-token", result.Config.DB.AuthToken)
	}
	if result.Config.DB.Path != "/some/path.db" {
		t.Fatalf("Load() DB.Path = %s, want /some/path.db", result.Config.DB.Path)
	}
}

func TestLoad_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.toml")
	// Minimal config with no daemon section
	cfgContent := `version = 1

[db]
path = "/some/path.db"
`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	result, err := Load(cfgPath, false)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check defaults are applied
	if result.Config.Daemon.PeriodicSync != "5m" {
		t.Fatalf("Load() Daemon.PeriodicSync = %s, want 5m (default)", result.Config.Daemon.PeriodicSync)
	}
	if result.Config.Daemon.LogLevel != "info" {
		t.Fatalf("Load() Daemon.LogLevel = %s, want info (default)", result.Config.Daemon.LogLevel)
	}
	if result.Config.Daemon.SyncRetryMax != 3 {
		t.Fatalf("Load() Daemon.SyncRetryMax = %d, want 3 (default)", result.Config.Daemon.SyncRetryMax)
	}
}

func TestNewViper_EnvPrefix(t *testing.T) {
	v := NewViper()

	// Set env var
	_ = os.Setenv("UGH_DB_SYNC_URL", "libsql://test.turso.io")
	defer func() { _ = os.Unsetenv("UGH_DB_SYNC_URL") }()

	// AutomaticEnv should pick it up
	if v.GetString("db.sync_url") != "libsql://test.turso.io" {
		t.Fatalf("Viper GetString(db.sync_url) = %s, want libsql://test.turso.io", v.GetString("db.sync_url"))
	}
}
