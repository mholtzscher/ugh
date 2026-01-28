package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/editor"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v2"
)

type editOptions struct {
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

func editCommand() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Aliases:   []string{"e"},
		Usage:     "Edit a task",
		ArgsUsage: "<id>",
		Description: "Edit a task by ID.\n\n" +
			"Opens the task in your editor ($VISUAL or $EDITOR) by default.\n" +
			"Use flags for quick single-field changes without opening an editor.\n\n" +
			"Examples:\n" +
			"  ugh edit 1                          # Open in editor (default)\n" +
			"  ugh edit 1 -p A                     # Set priority to A\n" +
			"  ugh edit 1 --no-priority            # Remove priority\n" +
			"  ugh edit 1 --description \"New text\" # Change description\n" +
			"  ugh edit 1 -P urgent                # Add project 'urgent'\n" +
			"  ugh edit 1 --remove-project old     # Remove project 'old'\n" +
			"  ugh edit 1 -c work -m due:tomorrow  # Add context and metadata",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "description",
				Usage: "Update description",
			},
			&cli.StringFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "Set priority (A-Z)",
			},
			&cli.BoolFlag{
				Name:  "no-priority",
				Usage: "Remove priority",
			},
			&cli.StringSliceFlag{
				Name:    "project",
				Aliases: []string{"P"},
				Usage:   "Add project (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:    "context",
				Aliases: []string{"c"},
				Usage:   "Add context (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:    "meta",
				Aliases: []string{"m"},
				Usage:   "Set metadata key:value (repeatable)",
			},
			&cli.BoolFlag{
				Name:    "done",
				Aliases: []string{"x"},
				Usage:   "Mark as done",
			},
			&cli.BoolFlag{
				Name:  "undone",
				Usage: "Mark as not done",
			},
			&cli.StringSliceFlag{
				Name:  "remove-project",
				Usage: "Remove project (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:  "remove-context",
				Usage: "Remove context (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:  "remove-meta",
				Usage: "Remove metadata key (repeatable)",
			},
			&cli.BoolFlag{
				Name:    "editor",
				Aliases: []string{"e"},
				Usage:   "Open in $VISUAL/$EDITOR",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return errors.New("edit requires a task id")
			}

			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()
			ids, err := deps.ParseIDs(c.Args().Slice())
			if err != nil {
				return err
			}
			id := ids[0]

			svc, err := deps.NewService(c)
			if err != nil {
				return err
			}
			defer func() { _ = svc.Close() }()

			if err := deps.MaybeSyncBeforeWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync pull: %w", err)
			}

			opts := editOptions{
				Description:    c.String("description"),
				Priority:       c.String("priority"),
				NoPriority:     c.Bool("no-priority"),
				Projects:       c.StringSlice("project"),
				Contexts:       c.StringSlice("context"),
				Meta:           c.StringSlice("meta"),
				Done:           c.Bool("done"),
				Undone:         c.Bool("undone"),
				RemoveProjects: c.StringSlice("remove-project"),
				RemoveContexts: c.StringSlice("remove-context"),
				RemoveMeta:     c.StringSlice("remove-meta"),
				Editor:         c.Bool("editor"),
			}

			var saved *store.Task
			if opts.Editor && hasFieldFlags(opts) {
				return errors.New("cannot combine field flags with --editor")
			}

			if opts.Editor || !hasFieldFlags(opts) {
				saved, err = runEditorMode(ctx, svc, id)
				if err != nil {
					return err
				}
			} else {
				saved, err = runFlagsMode(ctx, svc, id, opts)
				if err != nil {
					return err
				}
			}

			if saved == nil {
				fmt.Println("No changes made")
				return nil
			}

			if err := deps.MaybeSyncAfterWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync push: %w", err)
			}

			writer := deps.OutputWriter(c)
			return writer.WriteTask(saved)
		},
	}
}

func hasFieldFlags(opts editOptions) bool {
	return opts.Description != "" ||
		opts.Priority != "" ||
		opts.NoPriority ||
		len(opts.Projects) > 0 ||
		len(opts.Contexts) > 0 ||
		len(opts.Meta) > 0 ||
		opts.Done ||
		opts.Undone ||
		len(opts.RemoveProjects) > 0 ||
		len(opts.RemoveContexts) > 0 ||
		len(opts.RemoveMeta) > 0
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

func runFlagsMode(ctx context.Context, svc service.Service, id int64, opts editOptions) (*store.Task, error) {
	if opts.Done && opts.Undone {
		return nil, errors.New("cannot use both --done and --undone")
	}
	if opts.Priority != "" {
		p := strings.ToUpper(strings.TrimSpace(opts.Priority))
		if len(p) != 1 || p[0] < 'A' || p[0] > 'Z' {
			return nil, fmt.Errorf("invalid priority %q: must be A-Z", opts.Priority)
		}
	}

	meta, err := parseMetaFlags(opts.Meta)
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	req := service.UpdateTaskRequest{
		ID:             id,
		AddProjects:    opts.Projects,
		AddContexts:    opts.Contexts,
		SetMeta:        meta,
		RemoveProjects: opts.RemoveProjects,
		RemoveContexts: opts.RemoveContexts,
		RemoveMetaKeys: opts.RemoveMeta,
		RemovePriority: opts.NoPriority,
	}

	if opts.Description != "" {
		req.Description = &opts.Description
	}

	if opts.Priority != "" {
		req.Priority = &opts.Priority
	}

	if opts.Done {
		done := true
		req.Done = &done
	} else if opts.Undone {
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
