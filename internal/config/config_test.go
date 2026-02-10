//nolint:testpackage // Tests exercise unexported config helpers.
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPath(t *testing.T) {
	path, err := DefaultPath()
	require.NoError(t, err, "DefaultPath() error")
	require.NotEmpty(t, path, "DefaultPath() returned empty string")
	require.True(t, filepath.IsAbs(path), "DefaultPath() returned relative path: %s", path)
}

func TestLoad_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	missingPath := filepath.Join(tmpDir, "does-not-exist.toml")

	_, err := Load(missingPath, false)
	require.Error(t, err, "Load() should return error for missing file")
	assert.ErrorIs(t, err, ErrNotFound, "Load() wrong error")
}

func TestLoad_AllowMissing(t *testing.T) {
	tmpDir := t.TempDir()
	missingPath := filepath.Join(tmpDir, "does-not-exist.toml")

	result, err := Load(missingPath, true)
	require.NoError(t, err, "Load() error")
	assert.False(t, result.WasLoaded, "Load() should report WasLoaded = false")
	assert.Equal(t, DefaultVersion, result.Config.Version, "Load() version mismatch")
}

func TestLoad_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.toml")
	cfgContent := `version = 1

[db]
path = "~/.local/share/ugh/ugh.sqlite"
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0o644), "write config error")

	result, err := Load(cfgPath, false)
	require.NoError(t, err, "Load() error")
	assert.True(t, result.WasLoaded, "Load() should report WasLoaded = true")
	assert.NotEmpty(t, result.Config.DB.Path, "Load() DB.Path is empty")
	assert.Equal(t, cfgPath, result.UsedPath, "Load() UsedPath mismatch")
}

func TestLoad_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "invalid.toml")
	invalidContent := `[invalid section
something = "value"
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(invalidContent), 0o644), "write config error")

	_, err := Load(cfgPath, false)
	require.Error(t, err, "Load() should return error for invalid TOML")
	assert.ErrorIs(t, err, ErrInvalid, "Load() wrong error")
}

func TestResolveDBPath_ExpandEnv(t *testing.T) {
	t.Setenv("HOME", "/test/home")

	path, err := ResolveDBPath("", "$HOME/data/ugh.sqlite")
	require.NoError(t, err, "ResolveDBPath() error")
	assert.Equal(t, "/test/home/data/ugh.sqlite", path, "ResolveDBPath() mismatch")
}

func TestResolveDBPath_ExpandHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "get home dir error")

	path, err := ResolveDBPath("", "~/data/ugh.sqlite")
	require.NoError(t, err, "ResolveDBPath() error")
	assert.Equal(t, filepath.Join(homeDir, "data/ugh.sqlite"), path, "ResolveDBPath() mismatch")
}

func TestResolveDBPath_ExpandHomeAlone(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "get home dir error")

	path, err := ResolveDBPath("", "~")
	require.NoError(t, err, "ResolveDBPath() error")
	assert.Equal(t, homeDir, path, "ResolveDBPath() mismatch")
}

func TestResolveDBPath_RelativeToConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config", "ugh-config.toml")
	require.NoError(t, os.MkdirAll(filepath.Dir(cfgPath), 0o755), "mkdir error")

	path, err := ResolveDBPath(cfgPath, "data/ugh.sqlite")
	require.NoError(t, err, "ResolveDBPath() error")
	assert.Equal(t, filepath.Join(tmpDir, "config", "data", "ugh.sqlite"), path, "ResolveDBPath() mismatch")
}

func TestResolveDBPath_RelativeNoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)
	err := os.WriteFile("ugh.sqlite", []byte{}, 0o644)
	require.NoError(t, err, "create sqlite file error")

	path, err := ResolveDBPath("", "ugh.sqlite")
	require.NoError(t, err, "ResolveDBPath() error")

	expected := filepath.Join(tmpDir, "ugh.sqlite")
	expectedInfo, err := os.Stat(expected)
	require.NoError(t, err, "stat expected sqlite file error")
	gotInfo, err := os.Stat(path)
	require.NoError(t, err, "stat resolved sqlite file error")
	assert.True(t, os.SameFile(expectedInfo, gotInfo), "ResolveDBPath() should return same file")
}

func TestResolveDBPath_ErrorEmpty(t *testing.T) {
	_, err := ResolveDBPath("", "")
	require.Error(t, err, "ResolveDBPath() should return error for empty path")
}

func TestUserConfigDir(t *testing.T) {
	dir, err := userConfigDir()
	require.NoError(t, err, "userConfigDir() error")
	require.NotEmpty(t, dir, "userConfigDir() returned empty string")
	require.True(t, filepath.IsAbs(dir), "userConfigDir() returned relative path: %s", dir)
}
