package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestScripts(t *testing.T) {
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "ugh")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %v\n%s", err, output)

	params := testscript.Params{
		Dir: filepath.Join("testdata", "script"),
		Setup: func(env *testscript.Env) error {
			path := binDir + string(os.PathListSeparator) + os.Getenv("PATH")
			env.Setenv("PATH", path)
			env.Setenv("HOME", filepath.Join(env.WorkDir, "home"))
			env.Setenv("XDG_CONFIG_HOME", filepath.Join(env.WorkDir, "home", ".config"))
			env.Setenv("XDG_DATA_HOME", filepath.Join(env.WorkDir, "home", ".local", "share"))
			env.Setenv("TZ", "UTC")
			return nil
		},
	}

	testscript.Run(t, params)
}
