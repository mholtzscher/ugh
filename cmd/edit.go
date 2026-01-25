package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/mholtzscher/ugh/internal/ui"
	"github.com/spf13/cobra"
)

var editOpts struct {
	Text string
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a task",
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
			if err := mustTTY(); err != nil {
				return errors.New("todo text required (use --text in non-interactive mode)")
			}
			currentLine := todotxt.Format(todotxt.Parsed{
				Done:           current.Done,
				Priority:       current.Priority,
				CompletionDate: current.CompletionDate,
				CreationDate:   current.CreationDate,
				Description:    current.Description,
				Projects:       current.Projects,
				Contexts:       current.Contexts,
				Meta:           current.Meta,
				Unknown:        current.Unknown,
			})
			value, err := ui.PromptTodoLine("Edit task", currentLine)
			if err != nil {
				return err
			}
			line = strings.TrimSpace(value)
		}
		if line == "" {
			return errors.New("todo text required")
		}

		parsed := todotxt.ParseLine(line)
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
