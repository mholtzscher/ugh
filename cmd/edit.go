package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/spf13/cobra"
)

var editOpts struct {
	Text string
}

var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit a task",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("edit requires a task id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}

		st, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer st.Close()

		current, err := st.GetTask(ctx, ids[0])
		if err != nil {
			return err
		}

		line := strings.TrimSpace(editOpts.Text)
		if line == "" {
			return errors.New("todo text required")
		}

		parsed := todotxt.ParseLine(line)
		if parsed.CreationDate == nil {
			if current.CreationDate != nil {
				parsed.CreationDate = current.CreationDate
			} else if parsed.Done && parsed.CompletionDate != nil {
				parsed.CreationDate = parsed.CompletionDate
			} else {
				parsed.CreationDate = nowDate()
			}
		}
		updated := &store.Task{
			ID:             current.ID,
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

		saved, err := st.UpdateTask(ctx, updated)
		if err != nil {
			return err
		}
		writer := outputWriter()
		return writer.WriteTask(saved)
	},
}

func init() {
	editCmd.Flags().StringVar(&editOpts.Text, "text", "", "todo.txt formatted line")
}
