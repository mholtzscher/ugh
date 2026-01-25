package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
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
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}

	params := testscript.Params{
		Dir: filepath.Join("testdata", "script"),
		Setup: func(env *testscript.Env) error {
			path := binDir + string(os.PathListSeparator) + os.Getenv("PATH")
			env.Setenv("PATH", path)
			env.Setenv("TZ", "UTC")
			return nil
		},
	}

	testscript.Run(t, params)
}
