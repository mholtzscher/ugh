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
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
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

		st, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = st.Close() }()

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
			parsed := todotxt.ParseLine(line)
			if parsed.CreationDate == nil {
				parsed.CreationDate = nowDate()
			}
			task := &store.Task{
				Done:           parsed.Done,
				Priority:       parsed.Priority,
				CompletionDate: parsed.CompletionDate,
				CreationDate:   parsed.CreationDate,
				Description:    parsed.Description,
				Projects:       parsed.Projects,
				Contexts:       parsed.Contexts,
				Meta:           parsed.Meta,
				Unknown:        parsed.Unknown,
			}
			if _, err := st.CreateTask(ctx, task); err != nil {
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
