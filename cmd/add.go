package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/spf13/cobra"
)

var addOpts struct {
	Priority  string
	Projects  []string
	Contexts  []string
	Meta      []string
	Done      bool
	Created   string
	Completed string
}

var addCmd = &cobra.Command{
	Use:     "add [todo.txt line]",
	Aliases: []string{"a"},
	Short:   "Add a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		line := strings.TrimSpace(strings.Join(args, " "))
		if line == "" {
			return errors.New("todo text required")
		}

		parsed := todotxt.ParseLine(line)
		if addOpts.Priority != "" {
			parsed.Priority = normalizePriority(addOpts.Priority)
		}
		if len(addOpts.Projects) > 0 {
			parsed.Projects = append(parsed.Projects, addOpts.Projects...)
		}
		if len(addOpts.Contexts) > 0 {
			parsed.Contexts = append(parsed.Contexts, addOpts.Contexts...)
		}
		if addOpts.Done {
			parsed.Done = true
			if parsed.CompletionDate == nil {
				parsed.CompletionDate = nowDate()
			}
		}
		if addOpts.Created != "" {
			date, err := parseDate(addOpts.Created)
			if err != nil {
				return err
			}
			parsed.CreationDate = date
		}
		if addOpts.Completed != "" {
			date, err := parseDate(addOpts.Completed)
			if err != nil {
				return err
			}
			parsed.CompletionDate = date
		}
		if parsed.CreationDate == nil {
			if parsed.Done && parsed.CompletionDate != nil {
				parsed.CreationDate = parsed.CompletionDate
			} else {
				parsed.CreationDate = nowDate()
			}
		}

		meta, err := parseMetaFlags(addOpts.Meta)
		if err != nil {
			return err
		}
		if len(meta) > 0 {
			if parsed.Meta == nil {
				parsed.Meta = map[string]string{}
			}
			for key, value := range meta {
				parsed.Meta[key] = value
			}
		}

		storeTask := &store.Task{
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

		st, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer st.Close()

		created, err := st.CreateTask(ctx, storeTask)
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTask(created)
	},
}

func init() {
	addCmd.Flags().StringVarP(&addOpts.Priority, "priority", "p", "", "priority letter")
	addCmd.Flags().StringSliceVar(&addOpts.Projects, "project", nil, "project tag (repeatable)")
	addCmd.Flags().StringSliceVar(&addOpts.Contexts, "context", nil, "context tag (repeatable)")
	addCmd.Flags().StringSliceVar(&addOpts.Meta, "meta", nil, "metadata key:value (repeatable)")
	addCmd.Flags().BoolVar(&addOpts.Done, "done", false, "mark task done")
	addCmd.Flags().StringVar(&addOpts.Created, "created", "", "creation date (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addOpts.Completed, "completed", "", "completion date (YYYY-MM-DD)")
}
