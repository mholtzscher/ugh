package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mholtzscher/ugh/internal/config"
)

func TestEffectiveDBPathUsesOverride(t *testing.T) {
	r := New()
	override := filepath.Join(t.TempDir(), "override.sqlite")
	r.SetGlobalOptions("", override, false, false)

	path, err := r.EffectiveDBPath()
	if err != nil {
		t.Fatalf("effective db path: %v", err)
	}
	if path != override {
		t.Fatalf("expected override path %q, got %q", override, path)
	}
}

func TestEffectiveDBPathForConfigResolvesRelativePath(t *testing.T) {
	base := t.TempDir()
	cfgPath := filepath.Join(base, "config.toml")
	cfg := &config.Config{DB: config.DB{Path: "data/ugh.sqlite"}}

	path, err := EffectiveDBPathForConfig(cfgPath, cfg)
	if err != nil {
		t.Fatalf("effective db path for config: %v", err)
	}

	want, err := filepath.Abs(filepath.Join(base, "data", "ugh.sqlite"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if path != want {
		t.Fatalf("expected resolved path %q, got %q", want, path)
	}
}

func TestPrepareDBPathCreatesDirAndSetsCacheEnv(t *testing.T) {
	t.Setenv("TURSO_GO_CACHE_DIR", "")
	dbPath := filepath.Join(t.TempDir(), "nested", "dir", "ugh.sqlite")

	if err := PrepareDBPath(dbPath); err != nil {
		t.Fatalf("prepare db path: %v", err)
	}

	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); err != nil {
		t.Fatalf("expected db dir to exist: %v", err)
	}

	wantCache := filepath.Join(dbDir, ".cache")
	gotCache := os.Getenv("TURSO_GO_CACHE_DIR")
	if gotCache != wantCache {
		t.Fatalf("expected TURSO_GO_CACHE_DIR=%q, got %q", wantCache, gotCache)
	}
}
