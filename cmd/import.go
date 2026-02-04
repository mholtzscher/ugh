package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var importCmd = &cli.Command{
	Name:      "import",
	Aliases:   []string{"in"},
	Usage:     "Import tasks from todo.txt",
	ArgsUsage: "<path|->",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("import requires a file path or -")
		}
		path := cmd.Args().Get(0)
		var reader io.Reader
		if path == "-" {
			reader = os.Stdin
		} else {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer func() { _ = file.Close() }()
			reader = file
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		scanner := bufio.NewScanner(reader)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		var added int64
		var skipped int64
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				skipped++
				continue
			}
			if _, err := svc.CreateTask(ctx, service.CreateTaskRequest{
				Line: line,
			}); err != nil {
				return err
			}
			added++
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read import: %w", err)
		}
		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}

		writer := outputWriter()
		return writer.WriteSummary(output.ImportSummary{Action: "import", Added: added, Skipped: skipped, File: path})
	},
}
