package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestCommandSurfaceSnapshot(t *testing.T) {
	app := NewApp()
	root, err := app.RootCommand()
	if err != nil {
		t.Fatalf("build root command: %v", err)
	}

	got := snapshotCommands(root.Commands)
	wantPath := filepath.Join("testdata", "command_surface.golden")
	wantBytes, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read snapshot %s: %v", wantPath, err)
	}
	want := string(wantBytes)

	if got != want {
		t.Fatalf("command surface snapshot mismatch\n--- want\n%s\n--- got\n%s", want, got)
	}
}

func snapshotCommands(commands []*cli.Command) string {
	lines := make([]string, 0)
	for _, cmd := range commands {
		appendSnapshotLine(cmd, "", &lines)
	}
	return strings.Join(lines, "\n") + "\n"
}

func appendSnapshotLine(cmd *cli.Command, prefix string, lines *[]string) {
	if cmd == nil {
		return
	}
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return
	}

	path := name
	if prefix != "" {
		path = prefix + " " + name
	}

	aliases := "-"
	if len(cmd.Aliases) > 0 {
		aliases = strings.Join(cmd.Aliases, ",")
	}

	category := strings.TrimSpace(cmd.Category)
	if category == "" {
		category = "-"
	}

	*lines = append(*lines, fmt.Sprintf("%s|aliases=%s|category=%s", path, aliases, category))
	for _, child := range cmd.Commands {
		appendSnapshotLine(child, path, lines)
	}
}
