package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/editor"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/spf13/cobra"
)

var editOpts struct {
	Description    string
	Priority       string
	NoPriority     bool
	Projects       []string
	Contexts       []string
	Meta           []string
	Done           bool
	Undone         bool
	RemoveProjects []string
	RemoveContexts []string
	RemoveMeta     []string
	Editor         bool
}

var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit a task",
	Long: `Edit a task by ID.

Opens the task in your editor ($VISUAL or $EDITOR) by default.
Use flags for quick single-field changes without opening an editor.

Examples:
  ugh edit 1                         # Open in editor (default)
  ugh edit 1 -p A                    # Set priority to A
  ugh edit 1 --no-priority           # Remove priority
  ugh edit 1 --description "New text" # Change description
  ugh edit 1 -P urgent               # Add project 'urgent'
  ugh edit 1 --remove-project old    # Remove project 'old'
  ugh edit 1 -c work -m due:tomorrow # Add context and metadata`,
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
		id := ids[0]

		svc, err := newTaskService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		var saved *store.Task
		if editOpts.Editor && hasFieldFlags() {
			return errors.New("cannot combine field flags with --editor")
		}

		if editOpts.Editor || !hasFieldFlags() {
			saved, err = runEditorMode(ctx, svc, id)
			if err != nil {
				return err
			}
		} else {
			saved, err = runFlagsMode(ctx, svc, id)
			if err != nil {
				return err
			}
		}

		if saved == nil {
			fmt.Println("No changes made")
			return nil
		}

		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}

		writer := outputWriter()
		return writer.WriteTask(saved)
	},
}

func hasFieldFlags() bool {
	return editOpts.Description != "" ||
		editOpts.Priority != "" ||
		editOpts.NoPriority ||
		len(editOpts.Projects) > 0 ||
		len(editOpts.Contexts) > 0 ||
		len(editOpts.Meta) > 0 ||
		editOpts.Done ||
		editOpts.Undone ||
		len(editOpts.RemoveProjects) > 0 ||
		len(editOpts.RemoveContexts) > 0 ||
		len(editOpts.RemoveMeta) > 0
}

func runEditorMode(ctx context.Context, svc *service.TaskService, id int64) (*store.Task, error) {
	task, err := svc.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	edited, err := editor.Edit(task)
	if err != nil {
		return nil, fmt.Errorf("editor: %w", err)
	}

	if edited == nil {
		return nil, nil
	}

	return svc.FullUpdateTask(ctx, service.FullUpdateTaskRequest{
		ID:          id,
		Description: edited.Description,
		Priority:    edited.Priority,
		Done:        edited.Done,
		Projects:    edited.Projects,
		Contexts:    edited.Contexts,
		Meta:        edited.Meta,
	})
}

func runFlagsMode(ctx context.Context, svc *service.TaskService, id int64) (*store.Task, error) {
	if editOpts.Done && editOpts.Undone {
		return nil, errors.New("cannot use both --done and --undone")
	}
	if editOpts.Priority != "" {
		p := strings.ToUpper(strings.TrimSpace(editOpts.Priority))
		if len(p) != 1 || p[0] < 'A' || p[0] > 'Z' {
			return nil, fmt.Errorf("invalid priority %q: must be A-Z", editOpts.Priority)
		}
	}

	meta, err := parseMetaFlags(editOpts.Meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	req := service.UpdateTaskRequest{
		ID:             id,
		AddProjects:    editOpts.Projects,
		AddContexts:    editOpts.Contexts,
		SetMeta:        meta,
		RemoveProjects: editOpts.RemoveProjects,
		RemoveContexts: editOpts.RemoveContexts,
		RemoveMetaKeys: editOpts.RemoveMeta,
		RemovePriority: editOpts.NoPriority,
	}

	if editOpts.Description != "" {
		req.Description = &editOpts.Description
	}

	if editOpts.Priority != "" {
		req.Priority = &editOpts.Priority
	}

	if editOpts.Done {
		done := true
		req.Done = &done
	} else if editOpts.Undone {
		done := false
		req.Done = &done
	}

	return svc.UpdateTask(ctx, req)
}

func parseMetaFlags(meta []string) (map[string]string, error) {
	if len(meta) == 0 {
		return nil, nil
	}
	result := make(map[string]string, len(meta))
	for _, m := range meta {
		k, v, ok := strings.Cut(m, ":")
		if !ok {
			return nil, fmt.Errorf("invalid meta format: %s (expected key:value)", m)
		}
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result, nil
}

func init() {
	editCmd.Flags().StringVar(&editOpts.Description, "description", "", "update description")
	editCmd.Flags().StringVarP(&editOpts.Priority, "priority", "p", "", "set priority (A-Z)")
	editCmd.Flags().BoolVar(&editOpts.NoPriority, "no-priority", false, "remove priority")
	editCmd.Flags().StringSliceVarP(&editOpts.Projects, "project", "P", nil, "add project (repeatable)")
	editCmd.Flags().StringSliceVarP(&editOpts.Contexts, "context", "c", nil, "add context (repeatable)")
	editCmd.Flags().StringSliceVarP(&editOpts.Meta, "meta", "m", nil, "set metadata key:value (repeatable)")
	editCmd.Flags().BoolVarP(&editOpts.Done, "done", "x", false, "mark as done")
	editCmd.Flags().BoolVar(&editOpts.Undone, "undone", false, "mark as not done")
	editCmd.Flags().StringSliceVar(&editOpts.RemoveProjects, "remove-project", nil, "remove project (repeatable)")
	editCmd.Flags().StringSliceVar(&editOpts.RemoveContexts, "remove-context", nil, "remove context (repeatable)")
	editCmd.Flags().StringSliceVar(&editOpts.RemoveMeta, "remove-meta", nil, "remove metadata key (repeatable)")
	editCmd.Flags().BoolVarP(&editOpts.Editor, "editor", "e", false, "open in $VISUAL/$EDITOR")
}
