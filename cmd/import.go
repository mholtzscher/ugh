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
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:     "import <path|->",
	Aliases: []string{"in"},
	Short:   "Import tasks from todo.txt",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("import requires a file path or -")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		path := args[0]
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

		svc, err := newTaskService(ctx)
		if err != nil {
			return err
		}
		defer svc.Close()

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

		writer := outputWriter()
		return writer.WriteSummary(output.ImportSummary{Action: "import", Added: added, Skipped: skipped, File: path})
	},
}
