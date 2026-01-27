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
