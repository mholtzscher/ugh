package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/shell"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var shellCmd = &cli.Command{
	Name:     "shell",
	Aliases:  []string{"sh", "repl"},
	Usage:    "Start interactive NLP shell",
	Category: "System",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Usage: "Execute commands from file (scripting mode)",
		},
		&cli.BoolFlag{
			Name:  "stdin",
			Usage: "Execute commands from stdin (scripting mode)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		err = maybeSyncBeforeWrite(ctx, svc)
		if err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		var opts shell.Options
		if cmd.String("file") != "" {
			opts.Mode = shell.ModeScriptFile
			opts.InputFile = cmd.String("file")
		} else if cmd.Bool("stdin") {
			opts.Mode = shell.ModeScriptStdin
		} else {
			opts.Mode = shell.ModeInteractive
		}

		repl := shell.NewREPL(svc, opts)
		err = repl.Run(ctx)
		if err != nil {
			return err
		}

		return maybeSyncAfterWrite(ctx, svc)
	},
}
