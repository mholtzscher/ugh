package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/editor"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/urfave/cli/v3"
)

var editCmd = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit a task",
	Description: `Edit a task by ID.

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
	ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  flags.FlagDescription,
			Usage: "update description",
		},
		&cli.StringFlag{
			Name:    flags.FlagPriority,
			Aliases: []string{"p"},
			Usage:   "set priority (A-Z)",
		},
		&cli.BoolFlag{
			Name:  flags.FlagNoPriority,
			Usage: "remove priority",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"P"},
			Usage:   "add project (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "add context (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagMeta,
			Aliases: []string{"m"},
			Usage:   "set metadata key:value (repeatable)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "mark as done",
		},
		&cli.BoolFlag{
			Name:  flags.FlagUndone,
			Usage: "mark as not done",
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveProject,
			Usage: "remove project (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveContext,
			Usage: "remove context (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveMeta,
			Usage: "remove metadata key (repeatable)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagEditor,
			Aliases: []string{"e"},
			Usage:   "open in $VISUAL/$EDITOR",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("edit requires a task id")
		}
		ids, err := parseIDs(commandArgs(cmd))
		if err != nil {
			return err
		}
		id := ids[0]

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		var saved *store.Task
		hasFields := hasFieldFlags(cmd)
		if cmd.Bool(flags.FlagEditor) && hasFields {
			return errors.New("cannot combine field flags with --editor")
		}

		if cmd.Bool(flags.FlagEditor) || !hasFields {
			saved, err = runEditorMode(ctx, svc, id)
			if err != nil {
				return err
			}
		} else {
			saved, err = runFlagsMode(ctx, cmd, svc, id)
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

func hasFieldFlags(cmd *cli.Command) bool {
	return cmd.String(flags.FlagDescription) != "" ||
		cmd.String(flags.FlagPriority) != "" ||
		cmd.Bool(flags.FlagNoPriority) ||
		len(cmd.StringSlice(flags.FlagProject)) > 0 ||
		len(cmd.StringSlice(flags.FlagContext)) > 0 ||
		len(cmd.StringSlice(flags.FlagMeta)) > 0 ||
		cmd.Bool(flags.FlagDone) ||
		cmd.Bool(flags.FlagUndone) ||
		len(cmd.StringSlice(flags.FlagRemoveProject)) > 0 ||
		len(cmd.StringSlice(flags.FlagRemoveContext)) > 0 ||
		len(cmd.StringSlice(flags.FlagRemoveMeta)) > 0
}

func runEditorMode(ctx context.Context, svc service.Service, id int64) (*store.Task, error) {
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

func runFlagsMode(ctx context.Context, cmd *cli.Command, svc service.Service, id int64) (*store.Task, error) {
	priority := cmd.String(flags.FlagPriority)
	if cmd.Bool(flags.FlagDone) && cmd.Bool(flags.FlagUndone) {
		return nil, errors.New("cannot use both --done and --undone")
	}
	if priority != "" {
		p := strings.ToUpper(strings.TrimSpace(priority))
		if len(p) != 1 || p[0] < 'A' || p[0] > 'Z' {
			return nil, fmt.Errorf("invalid priority %q: must be A-Z", priority)
		}
		priority = p
	}

	meta, err := parseMetaFlags(cmd.StringSlice(flags.FlagMeta))
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	req := service.UpdateTaskRequest{
		ID:             id,
		AddProjects:    cmd.StringSlice(flags.FlagProject),
		AddContexts:    cmd.StringSlice(flags.FlagContext),
		SetMeta:        meta,
		RemoveProjects: cmd.StringSlice(flags.FlagRemoveProject),
		RemoveContexts: cmd.StringSlice(flags.FlagRemoveContext),
		RemoveMetaKeys: cmd.StringSlice(flags.FlagRemoveMeta),
		RemovePriority: cmd.Bool(flags.FlagNoPriority),
	}

	if description := cmd.String(flags.FlagDescription); description != "" {
		req.Description = &description
	}

	if priority != "" {
		req.Priority = &priority
	}

	if cmd.Bool(flags.FlagDone) {
		done := true
		req.Done = &done
	} else if cmd.Bool(flags.FlagUndone) {
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
